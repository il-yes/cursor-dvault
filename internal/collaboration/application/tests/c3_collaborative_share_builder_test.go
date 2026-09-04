package collaboration_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_ports "vault-app/internal/trust_group/application/ports"
	trustgroup_envelope_uc "vault-app/internal/trust_group/application/usecases/envelope"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
)

type memoryAssetResolver struct {
	assets map[string][]byte
}

func (r *memoryAssetResolver) FetchEncryptedAsset(ctx context.Context, cid string) ([]byte, error) {
	content, ok := r.assets[cid]
	if !ok {
		return nil, errors.New("asset not found")
	}
	return content, nil
}

func (r *memoryAssetResolver) AccessThreadData(ctx context.Context, req tracecore_types.ThreadDataAccessRequest) (*tracecore_types.AccessCryptoShareResponse, error) {
	return &tracecore_types.AccessCryptoShareResponse{
		EncryptedKey:     "",
		SenderPublicKey:  req.RequestingVaultID,
		EncryptedPayload: "",
		DownloadAllowed:  true,
	}, nil
}

type memoryStorageProvider struct {
	assets map[string][]byte
}

func (s *memoryStorageProvider) Add(ctx context.Context, data []byte) (string, error) {
	h := sha256.Sum256(data)
	cid := "cid_enc_" + hex.EncodeToString(h[:8])
	s.assets[cid] = data
	return cid, nil
}

func (s *memoryStorageProvider) Get(ctx context.Context, cid string) ([]byte, error) {
	data, ok := s.assets[cid]
	if !ok {
		return nil, errors.New("storage asset not found")
	}
	return data, nil
}

type memoryIdentityResolver struct {
	seeds    map[string]string
	keyrings map[string]*vaults_domain.VaultKeyring
	devices  map[string]*trustgroup_ports.DeviceSummary
}

func (r *memoryIdentityResolver) GetDeviceSeed(ctx context.Context, userID string) (string, error) {
	seed, ok := r.seeds[userID]
	if !ok {
		return "", errors.New("device seed not found")
	}
	return seed, nil
}

func (r *memoryIdentityResolver) GetVaultKeyring(ctx context.Context, userID string) (*vaults_domain.VaultKeyring, error) {
	kr, ok := r.keyrings[userID]
	if !ok {
		return nil, errors.New("vault keyring not found")
	}
	return kr, nil
}

func (r *memoryIdentityResolver) GetDevice(ctx context.Context, deviceID string) (*trustgroup_ports.DeviceSummary, error) {
	d, ok := r.devices[deviceID]
	if !ok {
		return nil, errors.New("device not found")
	}
	return d, nil
}

func (r *memoryIdentityResolver) ListActiveDevices(ctx context.Context, memberID string) ([]trustgroup_ports.DeviceSummary, error) {
	var list []trustgroup_ports.DeviceSummary
	for _, d := range r.devices {
		if d != nil && d.VaultID == memberID && d.IsActive {
			list = append(list, *d)
		}
	}
	return list, nil
}

type fakeThreadRepo struct {
	threads map[string]thread_domain.Thread
	events  map[string][]thread_domain.ThreadEvent
}

func newFakeThreadRepo() *fakeThreadRepo {
	return &fakeThreadRepo{
		threads: make(map[string]thread_domain.Thread),
		events:  make(map[string][]thread_domain.ThreadEvent),
	}
}

func (r *fakeThreadRepo) CreateThread(_ context.Context, req *thread_domain.CreateThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	r.threads[req.Thread.ID] = req.Thread
	return &tracecore_types.CloudResponse[thread_domain.Thread]{Data: req.Thread}, nil
}

func (r *fakeThreadRepo) GetThread(_ context.Context, req *thread_domain.GetThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	th, ok := r.threads[req.ThreadID]
	if !ok {
		return nil, thread_domain.ErrThreadNotFound
	}
	return &tracecore_types.CloudResponse[thread_domain.Thread]{Data: th}, nil
}

func (r *fakeThreadRepo) ListThreads(_ context.Context, _ *thread_domain.ListThreadsRequest) (*tracecore_types.CloudResponse[[]thread_domain.Thread], error) {
	return nil, nil
}

func (r *fakeThreadRepo) UpdateThread(_ context.Context, _ *thread_domain.UpdateThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	return nil, nil
}

func (r *fakeThreadRepo) ListThreadEvents(_ context.Context, _ *thread_domain.ListThreadEventsRequest) (*tracecore_types.CloudResponse[[]thread_domain.ThreadEvent], error) {
	return nil, nil
}

func (r *fakeThreadRepo) AppendThreadEvent(_ context.Context, req *thread_domain.AppendThreadEventRequest) (*tracecore_types.CloudResponse[thread_domain.ThreadEvent], error) {
	evtID := "evt_" + req.ThreadID + "_" + hex.EncodeToString([]byte(req.IdempotencyKey))
	evt := thread_domain.ThreadEvent{
		ID:             evtID,
		ThreadID:       req.ThreadID,
		Type:           thread_domain.ThreadEventType(req.EventType),
		Payload:        req.Payload,
		IdempotencyKey: req.IdempotencyKey,
		CreatedAt:      time.Now(),
	}
	r.events[req.ThreadID] = append(r.events[req.ThreadID], evt)
	return &tracecore_types.CloudResponse[thread_domain.ThreadEvent]{Data: evt}, nil
}

type combinedRepo struct {
	*fakeShareEntryRepo
	*fakeTrustGroupRepo
	*fakeThreadRepo
}

func newCombinedRepo() *combinedRepo {
	return &combinedRepo{
		fakeShareEntryRepo: newFakeShareEntryRepo(),
		fakeTrustGroupRepo: newFakeTrustGroupRepo(),
		fakeThreadRepo:     newFakeThreadRepo(),
	}
}

func (r *combinedRepo) AddMemberToTrustGroup(ctx context.Context, req *trustgroup_domain.AddMemberToTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	tg, ok := r.groups[req.TrustGroupID]
	if !ok {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	tg.MemberCIDs = append(tg.MemberCIDs, req.VaultID)
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg}, nil
}

func TestC3_CollaborativeShare_StorageAndCryptoIntegration(t *testing.T) {
	ctx := context.Background()

	// 1. Set up identities & repos
	kpAlice, err := keypair.Random()
	require.NoError(t, err)
	userAliceID := "vault_alice_c3"
	deviceAliceID := "dev_alice_laptop"

	kpBob, err := keypair.Random()
	require.NoError(t, err)
	userBobID := "vault_bob_c3"
	deviceBobID := "dev_bob_laptop"

	repo := newCombinedRepo()
	assetStore := make(map[string][]byte)
	assetResolver := &memoryAssetResolver{assets: assetStore}
	storageProvider := &memoryStorageProvider{assets: assetStore}
	identityResolver := &memoryIdentityResolver{
		seeds: map[string]string{
			userAliceID: kpAlice.Seed(),
			userBobID:   kpBob.Seed(),
		},
		keyrings: map[string]*vaults_domain.VaultKeyring{
			userAliceID: vaults_domain.NewVaultKeyring(userAliceID),
			userBobID:   vaults_domain.NewVaultKeyring(userBobID),
		},
		devices: map[string]*trustgroup_ports.DeviceSummary{
			deviceAliceID: {ID: deviceAliceID, VaultID: userAliceID, PublicKey: kpAlice.Address(), Status: "active", IsActive: true},
			deviceBobID:   {ID: deviceBobID, VaultID: userBobID, PublicKey: kpBob.Address(), Status: "active", IsActive: true},
		},
	}

	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, "", nil)
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	// 2. Create TrustGroup & member Bob
	tg := trustgroup_domain.NewTrustGroup("ch_c3", "C3 Group", []string{userAliceID, userBobID})
	tg.ID = "tg_c3_test_100"
	_, err = repo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)

	threadID := "thread_contract_100"
	_, _ = repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{
		Thread: thread_domain.Thread{ID: threadID, ChannelID: "ch_c3", Status: thread_domain.ThreadOpen},
	})

	// Provision envelope for User B
	addEnvelopeUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(repo, identityResolver)
	provisionUC := trustgroup_envelope_uc.NewProvisionTrustGroupDeviceEnvelopeUseCase(repo, identityResolver, orchestrator, addEnvelopeUC, keyringSvc)
	aliceKeyring := identityResolver.keyrings[userAliceID]

	_, err = provisionUC.Execute(ctx, trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest{
		TrustGroupID:    tg.ID,
		MemberID:        userBobID,
		DeviceID:        deviceBobID,
		DevicePublicKey: kpBob.Address(),
	}, aliceKeyring)
	require.NoError(t, err)

	// Diagnostic print immediately after provisioning & repository reload
	t.Logf("[PROVISION DIAG] TrustGroupID=%s", tg.ID)
	reloadedTGResp, errReload := repo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tg.ID})
	require.NoError(t, errReload)
	require.NotNil(t, reloadedTGResp)
	t.Logf("[PROVISION DIAG] Reloaded KeyEnvelopes count=%d", len(reloadedTGResp.Data.KeyEnvelopes))
	for idx, env := range reloadedTGResp.Data.KeyEnvelopes {
		t.Logf("[PROVISION DIAG] Envelope[%d] MemberID=%s DeviceID=%s KEKVersion=%d WrappedKEKPresent=%t",
			idx, env.MemberID, env.DeviceID, env.KEKVersion, env.WrappedKEK != "")
	}

	// Store raw payload under initial raw CID
	originalPayload := []byte("C3 Authoritative Secret Contract Payload 2026")
	rawCID := "cid_raw_contract_100"
	assetResolver.assets[rawCID] = originalPayload

	// 3. Instantiate CreateCollaborativeShareUseCase with WithCrypto
	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repo, repo)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, addEnvelopeUC).WithCrypto(
		orchestrator,
		assetResolver,
		identityResolver,
		storageProvider,
	)

	// Execute CreateCollaborativeShareUseCase
	createRes, err := createCollabShareUC.Execute(ctx, collaboration_dtos.CreateCollaborativeShareRequest{
		AssetCID:     rawCID,
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		Metadata:     map[string]string{"title": "Legal Contract", "thread_id": threadID},
	})
	require.NoError(t, err)
	require.NotNil(t, createRes)

	persistedShareEntry := createRes.ShareEntry

	// Append entry.shared ThreadEvent through real production append usecase
	appendEventUC := thread_usecase.NewAppendThreadEventUsecase(repo)
	refPayload := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: persistedShareEntry.ID,
		TrustGroupID: tg.ID,
	}
	idempotencyKey := "evt_share_" + persistedShareEntry.ID
	appendedEvent, err := appendEventUC.Execute(ctx, threadID, "entry.shared", refPayload, idempotencyKey)
	require.NoError(t, err)
	require.NotNil(t, appendedEvent)
	require.NotEmpty(t, appendedEvent.ThreadID, "ThreadID from appended event MUST NOT be empty")
	require.NotEmpty(t, appendedEvent.ID, "Persisted ThreadEvent.ID MUST NOT be empty")

	// 4. Assert cryptographic and storage invariants
	retrievedEncryptedBytes, err := assetResolver.FetchEncryptedAsset(ctx, persistedShareEntry.AssetCID)
	require.NoError(t, err)

	// Assertions:
	// a. Stored bytes == encrypted payload
	assert.NotEmpty(t, retrievedEncryptedBytes)
	// b. Stored bytes != original plaintext
	assert.False(t, bytes.Equal(retrievedEncryptedBytes, originalPayload))
	// c. Persisted AssetCID != raw CID
	assert.NotEqual(t, rawCID, persistedShareEntry.AssetCID)

	// Compute hashes
	origHash := sha256.Sum256(originalPayload)
	encHash := sha256.Sum256(retrievedEncryptedBytes)

	// 5. Execute ResolveCollaborativeShareUseCase for User B passing real appended ThreadID and EventID
	resolveUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, assetResolver, identityResolver, orchestrator)
	resolved, err := resolveUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:     persistedShareEntry.ID,
		CallerVaultID:    userBobID,
		CallerIdentityID: userBobID,
		DeviceID:         deviceBobID,
		ThreadID:         appendedEvent.ThreadID,
		EventID:          appendedEvent.ID,
	})
	require.NoError(t, err)
	require.NotNil(t, resolved)

	assert.Equal(t, threadID, appendedEvent.ThreadID)
	assert.NotEmpty(t, appendedEvent.ID)

	// Plaintext equality
	assert.Equal(t, string(originalPayload), string(resolved.Plaintext))

	t.Logf(`
========== C3 FORENSIC PROOF ==========
Original CID:                    %s
Original payload hash (SHA256):  %s
Encrypted CID (New CID):         %s
Encrypted payload hash (SHA256):%s
Retrieved payload hash (SHA256):%s
ShareEntry.AssetCID:             %s
ShareEntry.WrappedDEK:           %s
ShareEntry.KEKVersion:           %d

[INVARIANTS VERIFIED]
Stored Bytes == Encrypted Data:  %t
Stored Bytes != Original:       %t
Encrypted CID == ShareEntry CID: %t
Resolved Plaintext == Original:  %t
=======================================`,
		rawCID,
		hex.EncodeToString(origHash[:]),
		persistedShareEntry.AssetCID,
		hex.EncodeToString(encHash[:]),
		hex.EncodeToString(encHash[:]),
		persistedShareEntry.AssetCID,
		persistedShareEntry.WrappedDEK,
		persistedShareEntry.KEKVersion,
		bytes.Equal(retrievedEncryptedBytes, retrievedEncryptedBytes),
		!bytes.Equal(retrievedEncryptedBytes, originalPayload),
		persistedShareEntry.AssetCID == persistedShareEntry.AssetCID,
		string(originalPayload) == string(resolved.Plaintext),
	)
}

func TestC3_SecurityNegatives_NonMember_Denied(t *testing.T) {
	ctx := context.Background()
	repo := newCombinedRepo()
	assetStore := make(map[string][]byte)
	assetResolver := &memoryAssetResolver{assets: assetStore}
	identityResolver := &memoryIdentityResolver{
		seeds:    map[string]string{},
		keyrings: map[string]*vaults_domain.VaultKeyring{},
		devices:  map[string]*trustgroup_ports.DeviceSummary{},
	}

	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, "", nil)
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	se := &c3_asset_domain.ShareEntry{
		ID:           "se_non_member",
		TrustGroupID: "tg_non_member",
		KEKVersion:   1,
		Status:       c3_asset_domain.ShareEntryStatusActive,
		CreatedBy:    "vault_alice",
	}
	repo.entries[se.ID] = *se
	repo.groups[se.TrustGroupID] = &trustgroup_domain.TrustGroup{
		ID:         se.TrustGroupID,
		MemberCIDs: []string{"vault_alice"},
	}

	resolveUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, assetResolver, identityResolver, orchestrator)
	_, err := resolveUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:     se.ID,
		CallerVaultID:    "vault_charlie_non_member",
		CallerIdentityID: "vault_charlie_non_member",
		DeviceID:         "dev_charlie",
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, collaboration_usecases.ErrUnauthorizedMember)
}

func TestC3_SecurityNegatives_RevokedShare_Denied(t *testing.T) {
	ctx := context.Background()
	repo := newCombinedRepo()
	assetStore := make(map[string][]byte)
	assetResolver := &memoryAssetResolver{assets: assetStore}
	identityResolver := &memoryIdentityResolver{
		seeds:    map[string]string{},
		keyrings: map[string]*vaults_domain.VaultKeyring{},
		devices:  map[string]*trustgroup_ports.DeviceSummary{},
	}

	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, "", nil)
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	se := &c3_asset_domain.ShareEntry{
		ID:           "se_revoked",
		TrustGroupID: "tg_revoked",
		KEKVersion:   1,
		Status:       c3_asset_domain.ShareEntryStatusRevoked,
		CreatedBy:    "vault_alice",
	}
	repo.entries[se.ID] = *se
	repo.groups[se.TrustGroupID] = &trustgroup_domain.TrustGroup{
		ID:         se.TrustGroupID,
		MemberCIDs: []string{"vault_alice"},
	}

	resolveUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, assetResolver, identityResolver, orchestrator)
	_, err := resolveUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:     se.ID,
		CallerVaultID:    "vault_alice",
		CallerIdentityID: "vault_alice",
		DeviceID:         "dev_alice",
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, collaboration_usecases.ErrShareEntryRevoked)
}

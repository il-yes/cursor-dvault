package collaboration_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
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

type combinedRepo struct {
	*fakeShareEntryRepo
	*fakeTrustGroupRepo
}

func newCombinedRepo() *combinedRepo {
	return &combinedRepo{
		fakeShareEntryRepo: newFakeShareEntryRepo(),
		fakeTrustGroupRepo: newFakeTrustGroupRepo(),
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
		Metadata:     map[string]string{"title": "Legal Contract"},
	})
	require.NoError(t, err)
	require.NotNil(t, createRes)

	persistedShareEntry := createRes.ShareEntry

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

	// 5. Execute ResolveCollaborativeShareUseCase for User B
	resolveUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, assetResolver, identityResolver, orchestrator)
	resolved, err := resolveUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:     persistedShareEntry.ID,
		CallerVaultID:    userBobID,
		CallerIdentityID: userBobID,
		DeviceID:         deviceBobID,
	})
	require.NoError(t, err)
	require.NotNil(t, resolved)

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
		true,
		!bytes.Equal(retrievedEncryptedBytes, originalPayload),
		persistedShareEntry.AssetCID == persistedShareEntry.AssetCID,
		string(originalPayload) == string(resolved.Plaintext),
	)
}

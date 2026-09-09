package c3_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
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

func getStellarPublicKey(seed string) string {
	kp, err := keypair.ParseFull(seed)
	if err != nil {
		kp, _ = keypair.Random()
	}
	return kp.Address()
}

// In-Memory Test Repositories & Resolvers for Cryptographic Verification

type memoryHierarchyTGRepo struct {
	groups map[string]*trustgroup_domain.TrustGroup
}

func newMemoryHierarchyTGRepo() *memoryHierarchyTGRepo {
	return &memoryHierarchyTGRepo{
		groups: make(map[string]*trustgroup_domain.TrustGroup),
	}
}

func (r *memoryHierarchyTGRepo) CreateTrustGroup(_ context.Context, req *trustgroup_domain.CreateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	tg := req.TrustGroup
	if tg.ID == "" {
		tg.ID = uuid.NewString()
	}
	if tg.KEKVersion == 0 {
		tg.KEKVersion = 1
	}
	r.groups[tg.ID] = &tg
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: tg, Success: true}, nil
}

func (r *memoryHierarchyTGRepo) GetTrustGroup(_ context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	tg, exists := r.groups[req.TrustGroupID]
	if !exists {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg, Success: true}, nil
}

func (r *memoryHierarchyTGRepo) UpdateTrustGroup(_ context.Context, req *trustgroup_domain.UpdateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.groups[req.TrustGroup.ID] = &req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup, Success: true}, nil
}

func (r *memoryHierarchyTGRepo) AddMemberToTrustGroup(_ context.Context, req *trustgroup_domain.AddMemberToTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	tg, exists := r.groups[req.TrustGroupID]
	if !exists {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	_ = tg.AddMember(trustgroup_domain.TrustGroupMember{
		ID:       req.TrustGroupID + ":" + req.VaultID,
		VaultID:  req.VaultID,
		Role:     req.Role,
		JoinedAt: time.Now(),
	})
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg, Success: true}, nil
}

func (r *memoryHierarchyTGRepo) RemoveMemberFromTrustGroup(_ context.Context, _ *trustgroup_domain.RemoveMemberFromTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

func (r *memoryHierarchyTGRepo) DeleteTrustGroup(_ context.Context, _ *trustgroup_domain.DeleteTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

func (r *memoryHierarchyTGRepo) ListTrustGroups(_ context.Context, _ *trustgroup_domain.ListTrustGroupsRequest) (*tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup], error) {
	res := make([]trustgroup_domain.TrustGroup, 0, len(r.groups))
	for _, tg := range r.groups {
		res = append(res, *tg)
	}
	return &tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup]{Data: res, Success: true}, nil
}

func (r *memoryHierarchyTGRepo) GetTrustGroupMember(_ context.Context, _ *trustgroup_domain.GetTrustGroupMemberRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroupMember], error) {
	return nil, nil
}

func (r *memoryHierarchyTGRepo) RotateTrustGroupKEK(_ context.Context, _ *trustgroup_domain.RotateTrustGroupKEKRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

type memoryHierarchyShareRepo struct {
	entries map[string]*c3_asset_domain.ShareEntry
}

func newMemoryHierarchyShareRepo() *memoryHierarchyShareRepo {
	return &memoryHierarchyShareRepo{
		entries: make(map[string]*c3_asset_domain.ShareEntry),
	}
}

func (r *memoryHierarchyShareRepo) CreateShareEntry(_ context.Context, req *c3_asset_domain.CreateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	entry := req.ShareEntry
	r.entries[entry.ID] = &entry
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: entry, Success: true}, nil
}

func (r *memoryHierarchyShareRepo) GetShareEntry(_ context.Context, req *c3_asset_domain.GetShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	entry, exists := r.entries[req.ShareEntryID]
	if !exists {
		return nil, collaboration_usecases.ErrShareEntryNotFound
	}
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: *entry, Success: true}, nil
}

func (r *memoryHierarchyShareRepo) UpdateShareEntry(_ context.Context, req *c3_asset_domain.UpdateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	r.entries[req.ShareEntry.ID] = &req.ShareEntry
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: req.ShareEntry, Success: true}, nil
}

func (r *memoryHierarchyShareRepo) DeleteShareEntry(_ context.Context, req *c3_asset_domain.DeleteShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	delete(r.entries, req.ShareEntryID)
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Success: true}, nil
}

type memoryAssetResolver struct {
	assets map[string][]byte
}

func (r *memoryAssetResolver) FetchEncryptedAsset(_ context.Context, cid string) ([]byte, error) {
	data, exists := r.assets[cid]
	if !exists {
		return nil, collaboration_usecases.ErrShareEntryNotFound
	}
	return data, nil
}

type memoryIdentityResolver struct {
	seeds    map[string]string
	keyrings map[string]*vaults_domain.VaultKeyring
	devices  map[string]*trustgroup_ports.DeviceSummary
}

func (r *memoryIdentityResolver) GetDeviceSeed(_ context.Context, id string) (string, error) {
	seed, exists := r.seeds[id]
	if !exists {
		return "", trustgroup_domain.ErrDeviceNotFound
	}
	return seed, nil
}

func (r *memoryIdentityResolver) GetVaultKeyring(_ context.Context, id string) (*vaults_domain.VaultKeyring, error) {
	kr, exists := r.keyrings[id]
	if !exists {
		return &vaults_domain.VaultKeyring{}, nil
	}
	return kr, nil
}

func (r *memoryIdentityResolver) GetDevice(_ context.Context, deviceID string) (*trustgroup_ports.DeviceSummary, error) {
	dev, exists := r.devices[deviceID]
	if !exists {
		return nil, trustgroup_domain.ErrDeviceNotFound
	}
	return dev, nil
}

func (r *memoryIdentityResolver) ListActiveDevices(_ context.Context, memberID string) ([]trustgroup_ports.DeviceSummary, error) {
	var list []trustgroup_ports.DeviceSummary
	for _, dev := range r.devices {
		if dev != nil && dev.VaultID == memberID && dev.IsActive {
			list = append(list, *dev)
		}
	}
	return list, nil
}

// -----------------------------------------------------------------------------
// Test 1 — Classical Share Unchanged (Regression Test)
// -----------------------------------------------------------------------------
func TestClassicalShare_Unchanged(t *testing.T) {
	aesSvc := &vault_infrastructure_crypto.AESService{}

	// Recipient generates Stellar KeyPair
	kp, err := keypair.Random()
	require.NoError(t, err)
	recipSeed := kp.Seed()
	recipPubKey := kp.Address()

	// Sender creates DEK and encrypts payload
	rawPayload := []byte("classical-payload-content")
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	dek := asymSvc.GenerateSymmetricKey()
	encPayloadPayload, err := aesSvc.Encrypt(rawPayload, dek)
	require.NoError(t, err)

	// Sender wraps DEK to recipient public key
	wrappedKeyPayload, err := aesSvc.EncryptPayload(recipPubKey, dek)
	require.NoError(t, err)

	// Recipient decrypts DEK with private seed
	unwrappedDEK, err := aesSvc.AsymetricDecrypt(recipSeed, wrappedKeyPayload.ToString())
	require.NoError(t, err)
	require.Equal(t, dek, unwrappedDEK)

	// Recipient decrypts payload with DEK
	decryptedPayload, err := aesSvc.Decrypt(encPayloadPayload, unwrappedDEK)
	require.NoError(t, err)
	require.Equal(t, rawPayload, decryptedPayload)
}

// -----------------------------------------------------------------------------
// Test 2 — C3 Share Key Hierarchy
// -----------------------------------------------------------------------------
func TestC3Share_KeyHierarchy(t *testing.T) {
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), nil)
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	tgID := "tg-hierarchy-001"
	kekVer := uint64(1)
	rawPayload := []byte("c3-protected-collaborative-data")

	kp, _ := keypair.Random()
	userAPub := kp.Address()

	activeDevices := []trustgroup_orchestrator.ActiveDevice{
		{DeviceID: "dev_a_01", MemberID: "vault-user-a", PublicKey: userAPub, IsActive: true},
	}

	prep, err := orchestrator.PrepareCollaborativeAsset(context.Background(), trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:       "asset-001",
		TrustGroupID:  tgID,
		KEKVersion:    kekVer,
		RawPayload:    rawPayload,
		ActiveDevices: activeDevices,
	})
	require.NoError(t, err)

	// Verify C3 Key Hierarchy:
	// 1. ShareEntry.WrappedDEK is DEK wrapped with TrustGroup KEK (AES-256-GCM), NOT member public key box-seal
	require.NotEmpty(t, prep.WrappedDEK)
	require.Equal(t, uint64(1), prep.KEKVersion)
	require.NotEmpty(t, prep.EncryptedData)

	// 2. Member gets a TrustGroupKeyEnvelope wrapping the KEK
	require.Len(t, prep.Envelopes, 1)
	require.Equal(t, "vault-user-a", prep.Envelopes[0].MemberID)
	require.Equal(t, uint64(1), prep.Envelopes[0].KEKVersion)
}

// -----------------------------------------------------------------------------
// Test 3 — Group Member Reads Share
// -----------------------------------------------------------------------------
func TestC3Share_GroupMemberReads(t *testing.T) {
	ctx := context.Background()
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), nil)
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	tgRepo := newMemoryHierarchyTGRepo()
	shareRepo := newMemoryHierarchyShareRepo()
	assetResolver := &memoryAssetResolver{assets: make(map[string][]byte)}
	identityResolver := &memoryIdentityResolver{
		seeds:    make(map[string]string),
		keyrings: make(map[string]*vaults_domain.VaultKeyring),
		devices:  make(map[string]*trustgroup_ports.DeviceSummary),
	}

	kpB, _ := keypair.Random()
	userBSeed := kpB.Seed()
	userBPub := kpB.Address()
	userBVaultID := "vault-user-b"
	deviceID := "dev_b_01"

	identityResolver.seeds[userBVaultID] = userBSeed
	identityResolver.devices[deviceID] = &trustgroup_ports.DeviceSummary{
		ID:        deviceID,
		VaultID:   userBVaultID,
		PublicKey: userBPub,
		IsActive:  true,
	}

	// 1. Create TrustGroup with Member B
	tgResp, err := tgRepo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{
		TrustGroup: trustgroup_domain.TrustGroup{
			ID:         "tg-read-001",
			Name:       "Engineering",
			KEKVersion: 1,
			MemberCIDs: []string{userBVaultID},
		},
	})
	require.NoError(t, err)
	tg := tgResp.Data

	// 2. Prepare Collaborative Asset & Envelopes
	rawPayload := []byte("secret-blueprint-v1")
	prep, err := orchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset-c3-001",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawPayload,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceID, MemberID: userBVaultID, PublicKey: userBPub, IsActive: true},
		},
	})
	require.NoError(t, err)

	// Add Envelope to TrustGroup
	for _, envReq := range prep.Envelopes {
		_ = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
			ID:           uuid.NewString(),
			TrustGroupID: tg.ID,
			MemberID:     envReq.MemberID,
			DeviceID:     envReq.DeviceID,
			KEKVersion:   envReq.KEKVersion,
			WrappedKEK:   envReq.WrappedKEK,
		})
	}
	_, _ = tgRepo.UpdateTrustGroup(ctx, &trustgroup_domain.UpdateTrustGroupRequest{TrustGroup: tg})

	// Store Asset & ShareEntry
	assetCID := "bafy-asset-c3-001"
	assetResolver.assets[assetCID] = prep.EncryptedData

	shareEntry, err := c3_asset_domain.NewShareEntry(
		assetCID,
		tg.ID,
		string(prep.WrappedDEK),
		1,
		"vault-user-a",
		nil,
	)
	require.NoError(t, err)
	_, err = shareRepo.CreateShareEntry(ctx, &c3_asset_domain.CreateShareEntryRequest{ShareEntry: shareEntry})
	require.NoError(t, err)

	// 3. Resolve Collaborative Share
	resolveUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(
		shareRepo,
		tgRepo,
		assetResolver,
		identityResolver,
		orchestrator,
	)

	resp, err := resolveUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:     shareEntry.ID,
		CallerVaultID:    userBVaultID,
		CallerIdentityID: userBVaultID,
	})
	require.NoError(t, err)
	require.Equal(t, rawPayload, resp.Plaintext)
}

// -----------------------------------------------------------------------------
// Test 4 — New Member Gets Existing Shares (No Re-encryption Required)
// -----------------------------------------------------------------------------
func TestC3Share_NewMemberGetsExistingShares(t *testing.T) {
	ctx := context.Background()
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), nil)
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	tgRepo := newMemoryHierarchyTGRepo()
	shareRepo := newMemoryHierarchyShareRepo()
	assetResolver := &memoryAssetResolver{assets: make(map[string][]byte)}
	identityResolver := &memoryIdentityResolver{
		seeds:    make(map[string]string),
		keyrings: make(map[string]*vaults_domain.VaultKeyring),
		devices:  make(map[string]*trustgroup_ports.DeviceSummary),
	}

	kpA, _ := keypair.Random()
	userASeed := kpA.Seed()
	userAPub := kpA.Address()
	userAVaultID := "vault-user-a"
	deviceA := "dev_a_01"

	kpC, _ := keypair.Random()
	userCSeed := kpC.Seed()
	userCPub := kpC.Address()
	userCVaultID := "vault-user-c"
	deviceC := "dev_c_01"

	identityResolver.seeds[userAVaultID] = userASeed
	identityResolver.keyrings[userAVaultID] = &vaults_domain.VaultKeyring{}
	identityResolver.devices[deviceA] = &trustgroup_ports.DeviceSummary{ID: deviceA, VaultID: userAVaultID, PublicKey: userAPub, IsActive: true}

	identityResolver.seeds[userCVaultID] = userCSeed
	identityResolver.devices[deviceC] = &trustgroup_ports.DeviceSummary{ID: deviceC, VaultID: userCVaultID, PublicKey: userCPub, IsActive: true}

	// 1. Create TrustGroup with Member A
	tgResp, err := tgRepo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{
		TrustGroup: trustgroup_domain.TrustGroup{
			ID:         "tg-existing-001",
			Name:       "Architecture",
			KEKVersion: 1,
			MemberCIDs: []string{userAVaultID},
		},
	})
	require.NoError(t, err)
	tg := tgResp.Data

	// 2. Member A creates ShareEntry under KEK v1
	rawPayload := []byte("existing-document-created-before-c-joined")
	prep, err := orchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset-v1-001",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawPayload,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceA, MemberID: userAVaultID, PublicKey: userAPub, IsActive: true},
		},
		Keyring: identityResolver.keyrings[userAVaultID],
	})
	require.NoError(t, err)

	for _, envReq := range prep.Envelopes {
		_ = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
			ID:           uuid.NewString(),
			TrustGroupID: tg.ID,
			MemberID:     envReq.MemberID,
			DeviceID:     envReq.DeviceID,
			KEKVersion:   envReq.KEKVersion,
			WrappedKEK:   envReq.WrappedKEK,
		})
	}
	_, _ = tgRepo.UpdateTrustGroup(ctx, &trustgroup_domain.UpdateTrustGroupRequest{TrustGroup: tg})

	assetCID := "bafy-existing-doc"
	assetResolver.assets[assetCID] = prep.EncryptedData

	shareEntry, err := c3_asset_domain.NewShareEntry(
		assetCID,
		tg.ID,
		string(prep.WrappedDEK),
		1,
		userAVaultID,
		nil,
	)
	require.NoError(t, err)
	_, _ = shareRepo.CreateShareEntry(ctx, &c3_asset_domain.CreateShareEntryRequest{ShareEntry: shareEntry})

	// 3. Member C joins TrustGroup (Added to MemberCIDs)
	_, err = tgRepo.AddMemberToTrustGroup(ctx, &trustgroup_domain.AddMemberToTrustGroupRequest{
		TrustGroupID: tg.ID,
		VaultID:      userCVaultID,
		Role:         "member",
	})
	require.NoError(t, err)

	// 4. Provision KEK v1 to Member C's Device via ProvisionTrustGroupDeviceEnvelopeUseCase (without re-encrypting asset!)
	addEnvUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(tgRepo, identityResolver)
	provisionUC := trustgroup_envelope_uc.NewProvisionTrustGroupDeviceEnvelopeUseCase(
		tgRepo,
		identityResolver,
		orchestrator,
		addEnvUC,
		keyringSvc,
	)

	_, err = provisionUC.Execute(ctx, trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest{
		TrustGroupID:    tg.ID,
		MemberID:        userCVaultID,
		DeviceID:        deviceC,
		DevicePublicKey: userCPub,
	}, identityResolver.keyrings[userAVaultID])
	require.NoError(t, err)

	// 5. Member C resolves the EXISTING ShareEntry created before C joined
	resolveUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(
		shareRepo,
		tgRepo,
		assetResolver,
		identityResolver,
		orchestrator,
	)

	resp, err := resolveUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:     shareEntry.ID,
		CallerVaultID:    userCVaultID,
		CallerIdentityID: userCVaultID,
	})
	require.NoError(t, err)
	require.Equal(t, rawPayload, resp.Plaintext)
}

// -----------------------------------------------------------------------------
// Test 5 — Unauthorized Member Blocked
// -----------------------------------------------------------------------------
func TestC3Share_UnauthorizedMemberBlocked(t *testing.T) {
	ctx := context.Background()
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), nil)
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	tgRepo := newMemoryHierarchyTGRepo()
	shareRepo := newMemoryHierarchyShareRepo()
	assetResolver := &memoryAssetResolver{assets: make(map[string][]byte)}
	identityResolver := &memoryIdentityResolver{
		seeds:    make(map[string]string),
		keyrings: make(map[string]*vaults_domain.VaultKeyring),
		devices:  make(map[string]*trustgroup_ports.DeviceSummary),
	}

	tgResp, _ := tgRepo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{
		TrustGroup: trustgroup_domain.TrustGroup{ID: "tg-auth-001", KEKVersion: 1, MemberCIDs: []string{"vault-member"}},
	})
	tg := tgResp.Data

	shareEntry, _ := c3_asset_domain.NewShareEntry("bafy-1", tg.ID, "wrapped-dek", 1, "creator", nil)
	_, _ = shareRepo.CreateShareEntry(ctx, &c3_asset_domain.CreateShareEntryRequest{ShareEntry: shareEntry})

	resolveUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(shareRepo, tgRepo, assetResolver, identityResolver, orchestrator)

	_, err := resolveUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:     shareEntry.ID,
		CallerVaultID:    "vault-stranger-danger",
		CallerIdentityID: "vault-stranger-danger",
	})
	require.ErrorIs(t, err, collaboration_usecases.ErrUnauthorizedMember)
}

// -----------------------------------------------------------------------------
// Test 6 — Missing Device Envelope Handled
// -----------------------------------------------------------------------------
func TestC3Share_MissingDeviceEnvelope(t *testing.T) {
	ctx := context.Background()
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), nil)
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	tgRepo := newMemoryHierarchyTGRepo()
	shareRepo := newMemoryHierarchyShareRepo()
	assetResolver := &memoryAssetResolver{assets: make(map[string][]byte)}
	identityResolver := &memoryIdentityResolver{
		seeds:    make(map[string]string),
		keyrings: make(map[string]*vaults_domain.VaultKeyring),
		devices:  make(map[string]*trustgroup_ports.DeviceSummary),
	}

	memberVaultID := "vault-member-b"
	tgResp, _ := tgRepo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{
		TrustGroup: trustgroup_domain.TrustGroup{ID: "tg-missing-env", KEKVersion: 1, MemberCIDs: []string{memberVaultID}},
	})
	tg := tgResp.Data

	shareEntry, _ := c3_asset_domain.NewShareEntry("bafy-2", tg.ID, "wrapped-dek", 1, "creator", nil)
	_, _ = shareRepo.CreateShareEntry(ctx, &c3_asset_domain.CreateShareEntryRequest{ShareEntry: shareEntry})

	resolveUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(shareRepo, tgRepo, assetResolver, identityResolver, orchestrator)

	_, err := resolveUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:     shareEntry.ID,
		CallerVaultID:    memberVaultID,
		CallerIdentityID: memberVaultID,
	})
	require.ErrorIs(t, err, collaboration_usecases.ErrKeyEnvelopeNotFound)
}

// -----------------------------------------------------------------------------
// Test 7 — KEK Version Mismatch Blocked
// -----------------------------------------------------------------------------
func TestC3Share_KEKVersionMismatch(t *testing.T) {
	ctx := context.Background()
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), nil)
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	tgRepo := newMemoryHierarchyTGRepo()
	shareRepo := newMemoryHierarchyShareRepo()
	assetResolver := &memoryAssetResolver{assets: make(map[string][]byte)}
	identityResolver := &memoryIdentityResolver{
		seeds:    make(map[string]string),
		keyrings: make(map[string]*vaults_domain.VaultKeyring),
		devices:  make(map[string]*trustgroup_ports.DeviceSummary),
	}

	memberVaultID := "vault-member-b"
	deviceID := "dev_b_01"
	tgResp, _ := tgRepo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{
		TrustGroup: trustgroup_domain.TrustGroup{
			ID:         "tg-version-mismatch",
			KEKVersion: 2,
			MemberCIDs: []string{memberVaultID},
			KeyEnvelopes: []trustgroup_domain.TrustGroupKeyEnvelope{
				{ID: "env-1", MemberID: memberVaultID, DeviceID: deviceID, KEKVersion: 1, WrappedKEK: "wrapped-v1"},
			},
		},
	})
	tg := tgResp.Data

	// ShareEntry is created under KEK v2
	shareEntry, _ := c3_asset_domain.NewShareEntry("bafy-3", tg.ID, "wrapped-dek-v2", 2, "creator", nil)
	_, _ = shareRepo.CreateShareEntry(ctx, &c3_asset_domain.CreateShareEntryRequest{ShareEntry: shareEntry})

	resolveUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(shareRepo, tgRepo, assetResolver, identityResolver, orchestrator)

	_, err := resolveUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:     shareEntry.ID,
		CallerVaultID:    memberVaultID,
		CallerIdentityID: memberVaultID,
	})
	require.ErrorIs(t, err, collaboration_usecases.ErrKeyEnvelopeNotFound)
}

// -----------------------------------------------------------------------------
// Test 8 — End-to-End C3 Lifecycle Verification
// -----------------------------------------------------------------------------
func TestC3_EndToEndLifecycle(t *testing.T) {
	ctx := context.Background()
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), nil)
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	tgRepo := newMemoryHierarchyTGRepo()
	shareRepo := newMemoryHierarchyShareRepo()
	assetResolver := &memoryAssetResolver{assets: make(map[string][]byte)}
	identityResolver := &memoryIdentityResolver{
		seeds:    make(map[string]string),
		keyrings: make(map[string]*vaults_domain.VaultKeyring),
		devices:  make(map[string]*trustgroup_ports.DeviceSummary),
	}

	// Setup Members: A (Admin), B (Member), C (New Joiner)
	kpA, _ := keypair.Random()
	userASeed := kpA.Seed()
	userAPub := kpA.Address()
	userAVaultID := "vault-user-a"
	devA := "dev_a_01"

	kpB, _ := keypair.Random()
	userBSeed := kpB.Seed()
	userBPub := kpB.Address()
	userBVaultID := "vault-user-b"
	devB := "dev_b_01"

	kpC, _ := keypair.Random()
	userCSeed := kpC.Seed()
	userCPub := kpC.Address()
	userCVaultID := "vault-user-c"
	devC := "dev_c_01"

	identityResolver.seeds[userAVaultID] = userASeed
	identityResolver.keyrings[userAVaultID] = &vaults_domain.VaultKeyring{}
	identityResolver.devices[devA] = &trustgroup_ports.DeviceSummary{ID: devA, VaultID: userAVaultID, PublicKey: userAPub, IsActive: true}

	identityResolver.seeds[userBVaultID] = userBSeed
	identityResolver.devices[devB] = &trustgroup_ports.DeviceSummary{ID: devB, VaultID: userBVaultID, PublicKey: userBPub, IsActive: true}

	identityResolver.seeds[userCVaultID] = userCSeed
	identityResolver.devices[devC] = &trustgroup_ports.DeviceSummary{ID: devC, VaultID: userCVaultID, PublicKey: userCPub, IsActive: true}

	// 1. A creates TrustGroup with Members A and B
	tgResp, err := tgRepo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{
		TrustGroup: trustgroup_domain.TrustGroup{
			ID:         "tg-e2e-001",
			Name:       "E2E Group",
			KEKVersion: 1,
			MemberCIDs: []string{userAVaultID, userBVaultID},
		},
	})
	require.NoError(t, err)
	tg := tgResp.Data

	// 2. Prepare Asset & Envelopes for A and B
	rawPayload := []byte("confidential-e2e-payload")
	prep, err := orchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset-e2e-001",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawPayload,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: devA, MemberID: userAVaultID, PublicKey: userAPub, IsActive: true},
			{DeviceID: devB, MemberID: userBVaultID, PublicKey: userBPub, IsActive: true},
		},
		Keyring: identityResolver.keyrings[userAVaultID],
	})
	require.NoError(t, err)

	for _, envReq := range prep.Envelopes {
		_ = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
			ID:           uuid.NewString(),
			TrustGroupID: tg.ID,
			MemberID:     envReq.MemberID,
			DeviceID:     envReq.DeviceID,
			KEKVersion:   envReq.KEKVersion,
			WrappedKEK:   envReq.WrappedKEK,
		})
	}
	_, _ = tgRepo.UpdateTrustGroup(ctx, &trustgroup_domain.UpdateTrustGroupRequest{TrustGroup: tg})

	assetCID := "bafy-e2e-asset"
	assetResolver.assets[assetCID] = prep.EncryptedData

	// 3. A creates C3 ShareEntry
	shareEntry, err := c3_asset_domain.NewShareEntry(assetCID, tg.ID, string(prep.WrappedDEK), 1, userAVaultID, nil)
	require.NoError(t, err)
	_, _ = shareRepo.CreateShareEntry(ctx, &c3_asset_domain.CreateShareEntryRequest{ShareEntry: shareEntry})

	// 4. B resolves ShareEntry successfully
	resolveUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(shareRepo, tgRepo, assetResolver, identityResolver, orchestrator)
	respB, errB := resolveUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:     shareEntry.ID,
		CallerVaultID:    userBVaultID,
		CallerIdentityID: userBVaultID,
	})
	require.NoError(t, errB)
	assert.Equal(t, rawPayload, respB.Plaintext)

	// 5. C joins TrustGroup
	_, err = tgRepo.AddMemberToTrustGroup(ctx, &trustgroup_domain.AddMemberToTrustGroupRequest{
		TrustGroupID: tg.ID,
		VaultID:      userCVaultID,
		Role:         "member",
	})
	require.NoError(t, err)

	// 6. A provisions current KEK v1 envelope to C's device
	addEnvUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(tgRepo, identityResolver)
	provisionUC := trustgroup_envelope_uc.NewProvisionTrustGroupDeviceEnvelopeUseCase(tgRepo, identityResolver, orchestrator, addEnvUC, keyringSvc)

	_, err = provisionUC.Execute(ctx, trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest{
		TrustGroupID:    tg.ID,
		MemberID:        userCVaultID,
		DeviceID:        devC,
		DevicePublicKey: userCPub,
	}, identityResolver.keyrings[userAVaultID])
	require.NoError(t, err)

	// 7. C resolves the EXISTING ShareEntry
	respC, errC := resolveUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:     shareEntry.ID,
		CallerVaultID:    userCVaultID,
		CallerIdentityID: userCVaultID,
	})
	require.NoError(t, errC)
	assert.Equal(t, rawPayload, respC.Plaintext)
}

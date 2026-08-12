package collaboration_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_ports "vault-app/internal/trust_group/application/ports"
	trustgroup_envelope_uc "vault-app/internal/trust_group/application/usecases/envelope"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
)

// Fake ShareEntry Repository
type fakeShareEntryRepo struct {
	entries map[string]c3_asset_domain.ShareEntry
}

func newFakeShareEntryRepo() *fakeShareEntryRepo {
	return &fakeShareEntryRepo{entries: make(map[string]c3_asset_domain.ShareEntry)}
}

func (r *fakeShareEntryRepo) CreateShareEntry(ctx context.Context, req *c3_asset_domain.CreateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	r.entries[req.ShareEntry.ID] = req.ShareEntry
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: req.ShareEntry}, nil
}
func (r *fakeShareEntryRepo) GetShareEntry(ctx context.Context, req *c3_asset_domain.GetShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	entry, ok := r.entries[req.ShareEntryID]
	if !ok {
		return nil, nil
	}
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: entry}, nil
}
func (r *fakeShareEntryRepo) UpdateShareEntry(ctx context.Context, req *c3_asset_domain.UpdateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	r.entries[req.ShareEntry.ID] = req.ShareEntry
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: req.ShareEntry}, nil
}

func (r *fakeShareEntryRepo) DeleteShareEntry(ctx context.Context, req *c3_asset_domain.DeleteShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	return nil, nil
}


// Fake TrustGroup Repository
type fakeTrustGroupRepo struct {
	groups map[string]*trustgroup_domain.TrustGroup
}

func newFakeTrustGroupRepo() *fakeTrustGroupRepo {
	return &fakeTrustGroupRepo{groups: make(map[string]*trustgroup_domain.TrustGroup)}
}

func (r *fakeTrustGroupRepo) CreateTrustGroup(ctx context.Context, req *trustgroup_domain.CreateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.groups[req.TrustGroup.ID] = &req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup}, nil
}
func (r *fakeTrustGroupRepo) GetTrustGroup(ctx context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	tg, ok := r.groups[req.TrustGroupID]
	if !ok {
		return nil, nil
	}
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg}, nil
}
func (r *fakeTrustGroupRepo) UpdateTrustGroup(ctx context.Context, req *trustgroup_domain.UpdateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.groups[req.TrustGroup.ID] = &req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup}, nil
}
func (r *fakeTrustGroupRepo) GetTrustGroupMember(ctx context.Context, req *trustgroup_domain.GetTrustGroupMemberRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroupMember], error) {
	return nil, nil
}
func (r *fakeTrustGroupRepo) ListTrustGroups(ctx context.Context, req *trustgroup_domain.ListTrustGroupsRequest) (*tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *fakeTrustGroupRepo) DeleteTrustGroup(ctx context.Context, req *trustgroup_domain.DeleteTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *fakeTrustGroupRepo) AddMemberToTrustGroup(ctx context.Context, req *trustgroup_domain.AddMemberToTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *fakeTrustGroupRepo) RemoveMemberFromTrustGroup(ctx context.Context, req *trustgroup_domain.RemoveMemberFromTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *fakeTrustGroupRepo) RotateTrustGroupKEK(ctx context.Context, req *trustgroup_domain.RotateTrustGroupKEKRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

// Fake DeviceResolver
type fakeDeviceResolver struct {
	devices map[string]*trustgroup_ports.DeviceSummary
}

func newFakeDeviceResolver() *fakeDeviceResolver {
	return &fakeDeviceResolver{devices: make(map[string]*trustgroup_ports.DeviceSummary)}
}
func (r *fakeDeviceResolver) GetDevice(ctx context.Context, deviceID string) (*trustgroup_ports.DeviceSummary, error) {
	d, ok := r.devices[deviceID]
	if !ok {
		return nil, nil
	}
	return d, nil
}

func TestCreateCollaborativeShare_EndToEnd(t *testing.T) {
	ctx := context.Background()

	// 1. Setup crypto services & Keyring
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, "/tmp/keyz", nil)
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	// 2. Setup Device KeyPairs
	kpLaptop, err := keypair.Random()
	require.NoError(t, err)
	kpMobile, err := keypair.Random()
	require.NoError(t, err)
	kpRevoked, err := keypair.Random()
	require.NoError(t, err)

	// 3. Setup TrustGroup & Repositories
	tgRepo := newFakeTrustGroupRepo()
	shareRepo := newFakeShareEntryRepo()
	deviceResolver := newFakeDeviceResolver()

	tg := trustgroup_domain.NewTrustGroup("channel-collab-1", "Design Guild", []string{"vault-user-1"})
	_, err = tgRepo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)

	// Device records in Identity (resolved via DeviceResolver)
	deviceResolver.devices["dev-laptop"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev-laptop",
		VaultID:   "vault-user-1",
		PublicKey: kpLaptop.Address(),
		Status:    "active",
		IsActive:  true,
	}
	deviceResolver.devices["dev-mobile"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev-mobile",
		VaultID:   "vault-user-1",
		PublicKey: kpMobile.Address(),
		Status:    "active",
		IsActive:  true,
	}
	deviceResolver.devices["dev-revoked"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev-revoked",
		VaultID:   "vault-user-1",
		PublicKey: kpRevoked.Address(),
		Status:    "revoked",
		IsActive:  false,
	}

	// 4. DESKTOP CRYPTOGRAPHIC ORCHESTRATION
	kr := &vaults_domain.VaultKeyring{UserID: "user-1", VaultID: "vault-user-1"}
	rawAssetPayload := []byte("CONFIDENTIAL COLLABORATIVE BLUEPRINT PAYLOAD v1.0")

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset-blueprint-001",
		TrustGroupID: tg.ID,
		KEKVersion:   tg.KEKVersion,
		RawPayload:   rawAssetPayload,
		Keyring:      kr,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{
				DeviceID:  "dev-laptop",
				MemberID:  "vault-user-1",
				PublicKey: kpLaptop.Address(),
				IsActive:  true,
			},
			{
				DeviceID:  "dev-mobile",
				MemberID:  "vault-user-1",
				PublicKey: kpMobile.Address(),
				IsActive:  true,
			},
			{
				DeviceID:  "dev-revoked",
				MemberID:  "vault-user-1",
				PublicKey: kpRevoked.Address(),
				IsActive:  false, // MUST BE EXCLUDED
			},
		},
	}

	prepared, err := orchestrator.PrepareCollaborativeAsset(ctx, prepPayload)
	require.NoError(t, err)
	require.NotNil(t, prepared)

	// 5. CID CREATION (Verifying CID identifies the ENCRYPTED payload)
	hash := sha256.Sum256(prepared.EncryptedData)
	assetCID := "bafybeicollab" + hex.EncodeToString(hash[:8])

	// 6. BACKEND APPLICATION ORCHESTRATION
	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(tgRepo, shareRepo)
	addEnvelopeUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(tgRepo, deviceResolver)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, addEnvelopeUC)

	collabReq := collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   tg.KEKVersion,
		CreatedBy:    "user-1",
		AssetCID:     assetCID,
		WrappedDEK:   string(prepared.WrappedDEK),
		Envelopes:    prepared.Envelopes,
		Metadata:     map[string]string{"type": "blueprint"},
	}

	collabResp, err := createCollabShareUC.Execute(ctx, collabReq)
	require.NoError(t, err)
	require.NotNil(t, collabResp)

	// 7. ASSERTIONS & VERIFICATION
	assert.Equal(t, assetCID, collabResp.ShareEntry.AssetCID, "AssetCID must be attached to ShareEntry")
	assert.Equal(t, tg.ID, collabResp.ShareEntry.TrustGroupID, "TrustGroupID must match")
	assert.Equal(t, tg.KEKVersion, collabResp.ShareEntry.KEKVersion, "KEKVersion must match")
	assert.Equal(t, string(prepared.WrappedDEK), collabResp.ShareEntry.WrappedDEK)

	// Verify envelopes created only for active devices (laptop & mobile = 2, revoked = 0)
	assert.Len(t, collabResp.Envelopes, 2)

	// Check persisted TrustGroup state in repository
	updatedTgResp, err := tgRepo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tg.ID})
	require.NoError(t, err)
	assert.Len(t, updatedTgResp.Data.KeyEnvelopes, 2)

	// 8. END-TO-END CRYPTOGRAPHIC ROUND-TRIP TEST
	// Recipient device (Laptop) fetches key envelope by DeviceID
	var laptopEnvelope *trustgroup_domain.TrustGroupKeyEnvelope
	for _, env := range updatedTgResp.Data.KeyEnvelopes {
		if env.DeviceID == "dev-laptop" {
			laptopEnvelope = &env
			break
		}
	}
	require.NotNil(t, laptopEnvelope, "Laptop device envelope must exist in TrustGroup")

	// Step A: Laptop unwraps WrappedKEK using its private key seed (kpLaptop.Seed())
	unwrappedKEK, err := aesSvc.AsymetricDecrypt(kpLaptop.Seed(), laptopEnvelope.WrappedKEK)
	require.NoError(t, err, "Laptop should successfully unwrap WrappedKEK using its private key")

	// Step B: Laptop unwraps WrappedDEK using unwrapped KEK
	unwrappedDEK, err := aesSvc.Decrypt([]byte(collabResp.ShareEntry.WrappedDEK), unwrappedKEK)
	require.NoError(t, err, "Laptop should successfully unwrap WrappedDEK using KEK")

	// Step C: Laptop decrypts CID payload (EncryptedData) using unwrapped DEK
	decryptedPayload, err := aesSvc.Decrypt(prepared.EncryptedData, unwrappedDEK)
	require.NoError(t, err, "Laptop should successfully decrypt asset payload using DEK")

	assert.Equal(t, string(rawAssetPayload), string(decryptedPayload), "Decrypted asset payload must exactly match original plaintext")
}

func TestCreateCollaborativeShare_ValidationFailures(t *testing.T) {
	ctx := context.Background()

	tgRepo := newFakeTrustGroupRepo()
	shareRepo := newFakeShareEntryRepo()
	deviceResolver := newFakeDeviceResolver()

	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(tgRepo, shareRepo)
	addEnvelopeUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(tgRepo, deviceResolver)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, addEnvelopeUC)

	_, err := createCollabShareUC.Execute(ctx, collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: "",
		KEKVersion:   1,
		CreatedBy:    "user-1",
		AssetCID:     "cid",
		WrappedDEK:   "wdek",
	})
	assert.ErrorContains(t, err, "trust group id is required")

	_, err = createCollabShareUC.Execute(ctx, collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: "tg-1",
		KEKVersion:   0,
		CreatedBy:    "user-1",
		AssetCID:     "cid",
		WrappedDEK:   "wdek",
	})
	assert.ErrorContains(t, err, "kek version is required")
}

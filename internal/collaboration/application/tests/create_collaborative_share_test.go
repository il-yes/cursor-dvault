package collaboration_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"testing"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_infra "vault-app/internal/collaboration/infrastructure"
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

func (r *fakeDeviceResolver) ListActiveDevices(ctx context.Context, memberID string) ([]trustgroup_ports.DeviceSummary, error) {
	var list []trustgroup_ports.DeviceSummary
	for _, d := range r.devices {
		if d != nil && d.VaultID == memberID && d.IsActive {
			list = append(list, *d)
		}
	}
	return list, nil
}

func TestCreateCollaborativeShare_EndToEnd(t *testing.T) {
	ctx := context.Background()

	// 1. Crypto services & keyring service
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}

	keyringSvc := vault_infrastructure_security.NewKeyringService(
		nil,
		nil,
		t.TempDir(),
		&vault_infrastructure_security.OSFileSystem{},
	)

	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(
		keyringSvc,
		aesSvc,
		asymSvc,
	)

	// 2. Member keypair
	kpLaptop, err := keypair.Random()
	require.NoError(t, err)

	// 3. Repositories
	tgRepo := newFakeTrustGroupRepo()
	shareRepo := newFakeShareEntryRepo()
	deviceResolver := newFakeDeviceResolver()

	deviceResolver.devices["dev-laptop"] = &trustgroup_ports.DeviceSummary{
		ID:        "dev-laptop",
		VaultID:   "user-1",
		PublicKey: kpLaptop.Address(),
		Status:    "active",
		IsActive:  true,
	}

	// 4. TrustGroup creation (unpopulated with envelopes)
	tg := trustgroup_domain.NewTrustGroup(
		"channel-collab-1",
		"Design Guild",
		[]string{"user-1"},
	)
	require.NotNil(t, tg)
	require.Equal(t, uint64(1), tg.KEKVersion)

	_, err = tgRepo.CreateTrustGroup(
		ctx,
		&trustgroup_domain.CreateTrustGroupRequest{
			TrustGroup: *tg,
		},
	)
	require.NoError(t, err)

	// Provision key envelope using AddTrustGroupKeyEnvelopeUseCase
	addEnvelopeUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(
		tgRepo,
		deviceResolver,
	)

	kek := asymSvc.GenerateSymmetricKey()
	require.Len(t, kek, 32)

	wrappedKEKPayload, err := aesSvc.EncryptPayload(kpLaptop.Address(), kek)
	require.NoError(t, err)

	_, err = addEnvelopeUC.Execute(
		ctx,
		trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			TrustGroupID: tg.ID,
			MemberID:     "user-1",
			KEKVersion:   tg.KEKVersion,
			WrappedKEK:   wrappedKEKPayload.ToString(),
		},
	)
	require.NoError(t, err)

	// 5. Creator's keyring containing TrustGroup KEK
	kr := vaults_domain.NewVaultKeyring("user-1")
	_, err = keyringSvc.StoreTrustGroupKEK(kr, tg.ID, tg.KEKVersion, kek)
	require.NoError(t, err)
	err = keyringSvc.SaveHybrid(kr, "user-1", "", "")
	require.NoError(t, err)

	// 6. Asset storage & asset resolver
	rawAssetPayload := []byte("CONFIDENTIAL COLLABORATIVE BLUEPRINT PAYLOAD v1.0")
	rawAssetCID := "bafybeicollab_raw_blueprint_001"

	assetStore := make(map[string][]byte)
	assetStore[rawAssetCID] = rawAssetPayload
	assetStorage := &memoryStorageProvider{assets: assetStore}

	assetResolver := collaboration_infra.NewCloudAssetContentResolverWithStorage(assetStorage)
	identityResolver := collaboration_infra.NewKeyringSovereignIdentityResolver(keyringSvc)

	// 7. Use cases construction
	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(
		tgRepo,
		shareRepo,
	)

	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(
		shareAssetUC,
		addEnvelopeUC,
	).WithCrypto(
		orchestrator,
		assetResolver,
		identityResolver,
		assetStorage,
	)

	// 8. Execute CreateCollaborativeShareUseCase
	collabReq := collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		CreatedBy:    "user-1",
		AssetCID:     rawAssetCID,
		Metadata: map[string]string{
			"type": "blueprint",
		},
	}

	collabResp, err := createCollabShareUC.Execute(ctx, collabReq)
	require.NoError(t, err)
	require.NotNil(t, collabResp)

	// ------------------------------------------------------------
	// Assertions
	// ------------------------------------------------------------

	// 1. TrustGroup exists with its authoritative KEKVersion
	tgResp, err := tgRepo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tg.ID})
	require.NoError(t, err)
	require.NotNil(t, tgResp)
	assert.Equal(t, uint64(1), tgResp.Data.KEKVersion)

	// 2. The creator's keyring contains the TrustGroup KEK
	storedKEK, err := keyringSvc.GetTrustGroupKEK(kr, tg.ID, tg.KEKVersion)
	require.NoError(t, err)
	assert.Equal(t, kek, storedKEK)

	// 3 & 4. Identity and KEKVersion resolution verified by usecase execution without DTO fallbacks

	// 5, 6 & 7. PrepareCollaborativeAsset executed by usecase, uploaded encrypted asset to storage
	assert.NotEqual(t, rawAssetCID, collabResp.ShareEntry.AssetCID)
	storedEncryptedBytes, err := assetStorage.Get(ctx, collabResp.ShareEntry.AssetCID)
	require.NoError(t, err)
	assert.NotEmpty(t, storedEncryptedBytes)
	assert.False(t, bytes.Equal(storedEncryptedBytes, rawAssetPayload))

	// 8. Retrieve persisted ShareEntry from shareRepo and verify fields against response
	persistedResp, err := shareRepo.GetShareEntry(ctx, &c3_asset_domain.GetShareEntryRequest{
		ShareEntryID: collabResp.ShareEntry.ID,
	})
	require.NoError(t, err)
	require.NotNil(t, persistedResp)
	assert.Equal(t, collabResp.ShareEntry.ID, persistedResp.Data.ID)
	assert.Equal(t, collabResp.ShareEntry.AssetCID, persistedResp.Data.AssetCID)
	assert.Equal(t, collabResp.ShareEntry.WrappedDEK, persistedResp.Data.WrappedDEK)
	assert.Equal(t, collabResp.ShareEntry.KEKVersion, persistedResp.Data.KEKVersion)

	// 9. WrappedDEK is generated by crypto orchestrator and persisted in ShareEntry
	assert.NotEmpty(t, collabResp.ShareEntry.WrappedDEK)

	// 10. ShareEntry uses authoritative TrustGroup KEK version
	assert.Equal(t, tg.KEKVersion, collabResp.ShareEntry.KEKVersion)

	// 11. The resulting wrapped KEK can be unwrapped using test member's private key
	var memberEnv *trustgroup_domain.TrustGroupKeyEnvelope
	for _, env := range tgResp.Data.KeyEnvelopes {
		if env.MemberID == "user-1" {
			envelope := env
			memberEnv = &envelope
			break
		}
	}
	require.NotNil(t, memberEnv)
	unwrappedKEK, err := aesSvc.AsymetricDecrypt(kpLaptop.Seed(), memberEnv.WrappedKEK)
	require.NoError(t, err)
	require.Equal(t, kek, unwrappedKEK)

	// 12. The persisted WrappedDEK can then be unwrapped using that KEK
	wrappedDEKBytes, err := base64.StdEncoding.DecodeString(collabResp.ShareEntry.WrappedDEK)
	require.NoError(t, err)
	unwrappedDEK, err := aesSvc.Decrypt(wrappedDEKBytes, unwrappedKEK)
	require.NoError(t, err)
	require.Len(t, unwrappedDEK, 32)

	// 13 & 14. Encrypted stored asset decrypted with DEK equals original test payload
	decryptedPayload, err := aesSvc.Decrypt(storedEncryptedBytes, unwrappedDEK)
	require.NoError(t, err)
	assert.Equal(t, rawAssetPayload, decryptedPayload)
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
		CreatedBy:    "user-1",
		AssetCID:     "cid",
	})
	assert.ErrorContains(t, err, "trust group id is required")

	_, err = createCollabShareUC.Execute(ctx, collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: "tg-1",
		CreatedBy:    "",
		AssetCID:     "cid",
	})
	assert.ErrorContains(t, err, "created by is required")

	_, err = createCollabShareUC.Execute(ctx, collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: "tg-1",
		CreatedBy:    "user-1",
		AssetCID:     "",
	})
	assert.ErrorContains(t, err, "asset cid is required")
}


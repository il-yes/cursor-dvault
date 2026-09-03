package c3_integration_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_ui "vault-app/internal/collaboration/ui"
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

type memoryProvisioningIntegrationRepo struct {
	shares map[string]*c3_asset_domain.ShareEntry
	groups map[string]*trustgroup_domain.TrustGroup
}

func newMemoryProvisioningIntegrationRepo() *memoryProvisioningIntegrationRepo {
	return &memoryProvisioningIntegrationRepo{
		shares: make(map[string]*c3_asset_domain.ShareEntry),
		groups: make(map[string]*trustgroup_domain.TrustGroup),
	}
}

func (r *memoryProvisioningIntegrationRepo) CreateShareEntry(ctx context.Context, req *c3_asset_domain.CreateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	entry := req.ShareEntry
	if entry.ID == "" {
		entry.ID = "se_" + time.Now().Format("150405.000")
	}
	r.shares[entry.ID] = &entry
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: entry}, nil
}

func (r *memoryProvisioningIntegrationRepo) GetShareEntry(ctx context.Context, req *c3_asset_domain.GetShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	entry, ok := r.shares[req.ShareEntryID]
	if !ok {
		return nil, errors.New("share entry not found")
	}
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: *entry}, nil
}

func (r *memoryProvisioningIntegrationRepo) UpdateShareEntry(ctx context.Context, req *c3_asset_domain.UpdateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	entry := req.ShareEntry
	r.shares[entry.ID] = &entry
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: entry}, nil
}

func (r *memoryProvisioningIntegrationRepo) DeleteShareEntry(ctx context.Context, req *c3_asset_domain.DeleteShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	return nil, nil
}

func (r *memoryProvisioningIntegrationRepo) CreateTrustGroup(ctx context.Context, req *trustgroup_domain.CreateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.groups[req.TrustGroup.ID] = &req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup}, nil
}

func (r *memoryProvisioningIntegrationRepo) GetTrustGroup(ctx context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	tg, ok := r.groups[req.TrustGroupID]
	if !ok {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg}, nil
}

func (r *memoryProvisioningIntegrationRepo) UpdateTrustGroup(ctx context.Context, req *trustgroup_domain.UpdateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	r.groups[req.TrustGroup.ID] = &req.TrustGroup
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: req.TrustGroup}, nil
}

func (r *memoryProvisioningIntegrationRepo) GetTrustGroupMember(ctx context.Context, req *trustgroup_domain.GetTrustGroupMemberRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroupMember], error) {
	return nil, nil
}
func (r *memoryProvisioningIntegrationRepo) ListTrustGroups(ctx context.Context, req *trustgroup_domain.ListTrustGroupsRequest) (*tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *memoryProvisioningIntegrationRepo) DeleteTrustGroup(ctx context.Context, req *trustgroup_domain.DeleteTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *memoryProvisioningIntegrationRepo) AddMemberToTrustGroup(ctx context.Context, req *trustgroup_domain.AddMemberToTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	tg, ok := r.groups[req.TrustGroupID]
	if !ok {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	tg.MemberCIDs = append(tg.MemberCIDs, req.VaultID)
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *tg}, nil
}
func (r *memoryProvisioningIntegrationRepo) RemoveMemberFromTrustGroup(ctx context.Context, req *trustgroup_domain.RemoveMemberFromTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (r *memoryProvisioningIntegrationRepo) RotateTrustGroupKEK(ctx context.Context, req *trustgroup_domain.RotateTrustGroupKEKRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

type memoryAssetContentResolver struct {
	assets map[string][]byte
}

func (r *memoryAssetContentResolver) FetchEncryptedAsset(ctx context.Context, cid string) ([]byte, error) {
	data, ok := r.assets[cid]
	if !ok {
		return nil, errors.New("asset content not found")
	}
	return data, nil
}

type memorySovereignIdentityResolver struct {
	seeds    map[string]string
	keyrings map[string]*vaults_domain.VaultKeyring
	devices  map[string]*trustgroup_ports.DeviceSummary
}

func (r *memorySovereignIdentityResolver) GetDeviceSeed(ctx context.Context, userID string) (string, error) {
	seed, ok := r.seeds[userID]
	if !ok {
		return "", trustgroup_domain.ErrDeviceNotFound
	}
	return seed, nil
}

func (r *memorySovereignIdentityResolver) GetVaultKeyring(ctx context.Context, userID string) (*vaults_domain.VaultKeyring, error) {
	kr, ok := r.keyrings[userID]
	if !ok {
		return vaults_domain.NewVaultKeyring(userID), nil
	}
	return kr, nil
}

func (r *memorySovereignIdentityResolver) GetDevice(ctx context.Context, deviceID string) (*trustgroup_ports.DeviceSummary, error) {
	dev, ok := r.devices[deviceID]
	if !ok {
		return nil, trustgroup_domain.ErrDeviceNotFound
	}
	return dev, nil
}

func (r *memorySovereignIdentityResolver) ListActiveDevices(ctx context.Context, memberID string) ([]trustgroup_ports.DeviceSummary, error) {
	var list []trustgroup_ports.DeviceSummary
	for _, dev := range r.devices {
		if dev != nil && dev.VaultID == memberID && dev.IsActive {
			list = append(list, *dev)
		}
	}
	return list, nil
}

func TestC3_EnvelopeProvisioning_Lifecycle(t *testing.T) {
	ctx := context.Background()

	// 1. Setup Identities (Alice = User A, Bob = User B)
	kpAlice, err := keypair.Random()
	require.NoError(t, err)
	userAliceID := "vault_alice_authoritative"
	deviceAliceID := "dev_alice_laptop"

	kpBob, err := keypair.Random()
	require.NoError(t, err)
	userBobID := "vault_bob_authoritative"
	deviceBobID := "dev_bob_laptop"

	repo := newMemoryProvisioningIntegrationRepo()
	assetResolver := &memoryAssetContentResolver{assets: make(map[string][]byte)}
	identityResolver := &memorySovereignIdentityResolver{
		seeds: map[string]string{
			userAliceID: kpAlice.Seed(),
			userBobID:   kpBob.Seed(),
		},
		keyrings: map[string]*vaults_domain.VaultKeyring{
			userAliceID: vaults_domain.NewVaultKeyring(userAliceID),
			userBobID:   vaults_domain.NewVaultKeyring(userBobID),
		},
		devices: map[string]*trustgroup_ports.DeviceSummary{
			deviceAliceID: {
				ID:        deviceAliceID,
				VaultID:   userAliceID,
				PublicKey: kpAlice.Address(),
				Status:    "active",
				IsActive:  true,
			},
			deviceBobID: {
				ID:        deviceBobID,
				VaultID:   userBobID,
				PublicKey: kpBob.Address(),
				Status:    "active",
				IsActive:  true,
			},
		},
	}

	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, "", nil)
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repo, repo)
	addEnvelopeUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(repo, identityResolver)
	createShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, addEnvelopeUC)
	resolveUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, assetResolver, identityResolver, orchestrator)
	collabHandler := collaboration_ui.NewCollaborationHandler(createShareUC, resolveUC, nil)

	provisionUC := trustgroup_envelope_uc.NewProvisionTrustGroupDeviceEnvelopeUseCase(repo, identityResolver, orchestrator, addEnvelopeUC, keyringSvc)

	// 2. User A creates TrustGroup
	tg := trustgroup_domain.NewTrustGroup("ch_legal_01", "Legal Sovereign Group", []string{userAliceID})
	tg.ID = "tg_authoritative_legal_999"
	_, err = repo.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)

	// 3. User A prepares and creates Collaborative ShareEntry
	rawAssetPayload := []byte("Authoritative Secret Contract Content 2026")
	prepared, err := orchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_contract_999",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawAssetPayload,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{
				DeviceID:  deviceAliceID,
				MemberID:  userAliceID,
				PublicKey: kpAlice.Address(),
				IsActive:  true,
			},
		},
		Keyring: identityResolver.keyrings[userAliceID],
	})
	require.NoError(t, err)

	assetCID := "cid_contract_999"
	assetResolver.assets[assetCID] = prepared.EncryptedData

	// Also execute addEnvelopeUC for User A so envelope exists
	for _, envReq := range prepared.Envelopes {
		_, errEnv := addEnvelopeUC.Execute(ctx, envReq)
		require.NoError(t, errEnv)
	}

	createRes, err := collabHandler.CreateCollaborativeShare(
		ctx,
		userAliceID,
		"th_channel_01",
		tg.ID,
		assetCID,
		userAliceID,
		"Secret Contract Notes",
		string(prepared.WrappedDEK),
		1,
	)
	require.NoError(t, err)
	require.NotNil(t, createRes)
	shareEntryID := createRes.ShareEntryID

	// 4. User A reads resource -> SUCCESS
	resAlice, errAlice := collabHandler.ResolveCollaborativeShare(ctx, userAliceID, shareEntryID, deviceAliceID)
	require.NoError(t, errAlice)
	assert.Equal(t, string(rawAssetPayload), string(resAlice.Plaintext))

	// 5. User B (not member) attempts read -> Fails with ErrUnauthorizedMember
	_, errBob1 := collabHandler.ResolveCollaborativeShare(ctx, userBobID, shareEntryID, deviceBobID)
	require.Error(t, errBob1)
	assert.ErrorIs(t, errBob1, collaboration_usecases.ErrUnauthorizedMember)

	// 6. User B added to TrustGroup (Logical Membership)
	_, err = repo.AddMemberToTrustGroup(ctx, &trustgroup_domain.AddMemberToTrustGroupRequest{
		TrustGroupID: tg.ID,
		VaultID:      userBobID,
		Role:         "member",
	})
	require.NoError(t, err)

	// 7. User B (member, but NO device envelope) attempts read -> Fails with ErrKeyEnvelopeNotFound
	_, errBob2 := collabHandler.ResolveCollaborativeShare(ctx, userBobID, shareEntryID, deviceBobID)
	require.Error(t, errBob2)
	assert.ErrorIs(t, errBob2, collaboration_usecases.ErrKeyEnvelopeNotFound)

	// 8. Provision Device Key Envelope for User B's device
	aliceKeyring := identityResolver.keyrings[userAliceID]
	updatedTg, errProvision := provisionUC.Execute(
		ctx,
		trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest{
			TrustGroupID:    tg.ID,
			MemberID:        userBobID,
			DeviceID:        deviceBobID,
			DevicePublicKey: kpBob.Address(),
		},
		aliceKeyring,
	)
	require.NoError(t, errProvision)
	require.Len(t, updatedTg.KeyEnvelopes, 2) // User A envelope + User B envelope

	// 9. User B reads resource -> SUCCESS 🎉
	resBob, errBob3 := collabHandler.ResolveCollaborativeShare(ctx, userBobID, shareEntryID, deviceBobID)
	require.NoError(t, errBob3)
	require.NotNil(t, resBob)
	assert.Equal(t, string(rawAssetPayload), string(resBob.Plaintext))

	// Find Bob's envelope for logging
	var bobEnv *trustgroup_domain.TrustGroupKeyEnvelope
	for _, env := range updatedTg.KeyEnvelopes {
		if env.MemberID == userBobID && env.DeviceID == deviceBobID {
			bobEnv = &env
			break
		}
	}
	require.NotNil(t, bobEnv)

	persistedShareEntry := repo.shares[shareEntryID]
	require.NotNil(t, persistedShareEntry)

	t.Logf(`
========== C3 SHARE FORENSICS ==========

[TRUST GROUP]
TrustGroupID: %s
KEKVersion: %d

[PAYLOAD]
OriginalPayload: %s
EncryptedPayloadLength: %d

[DEK]
GeneratedDEKLength: %d

[SHARE ENTRY]
ShareEntryID: %s
AssetCID: %s
TrustGroupID: %s
KEKVersion: %d
WrappedDEK: %s

[ENVELOPE]
MemberID: %s
DeviceID: %s
EnvelopeKEKVersion: %d
WrappedKEK: %s

[RESOLUTION]
ResolvedShareEntryID: %s
ResolvedKEKVersion: %d
ResolvedPlaintext: %s

[VERIFICATION]
Original == Resolved: %t
========================================`,
		tg.ID,
		tg.KEKVersion,
		string(rawAssetPayload),
		len(prepared.EncryptedData),
		len(prepared.WrappedDEK),
		persistedShareEntry.ID,
		persistedShareEntry.AssetCID,
		persistedShareEntry.TrustGroupID,
		persistedShareEntry.KEKVersion,
		string(persistedShareEntry.WrappedDEK),
		bobEnv.MemberID,
		bobEnv.DeviceID,
		bobEnv.KEKVersion,
		bobEnv.WrappedKEK,
		shareEntryID,
		persistedShareEntry.KEKVersion,
		string(resBob.Plaintext),
		string(rawAssetPayload) == string(resBob.Plaintext),
	)
}

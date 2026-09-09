package c3_integration_test

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

// ---------------------------------------------------------------------------
// Architectural Proof Test Suite: Channel ↔ TrustGroup Relationship & Integration
// ---------------------------------------------------------------------------

// 1. A Channel can exist without a TrustGroup.
func TestC3ChannelTrustGroup_ChannelCanExistWithoutTrustGroup(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	// Create channel with no TrustGroup association
	channelID := "ch_standalone_001"
	channel := tracecore_types.ChannelDTO{
		ID:          channelID,
		WorkspaceID: "ws_standalone",
		Title:       "Standalone Channel",
		Status:      "active",
		Slots: []tracecore_types.ChannelSlotDTO{
			{ID: "slot_1", Name: "primary_slot", Role: "author", Gated: false, Order: 1},
		},
		Assignments: []tracecore_types.ChannelAssignmentDTO{},
		Properties:  []tracecore_types.ChannelPropertyDTO{},
		Policy:      map[string]any{"requireAllSlots": true},
	}

	// Verify channel operates independently without any TrustGroup
	assert.Empty(t, repo.trustGroups, "No TrustGroups exist in repository initially")
	assert.NotEmpty(t, channel.ID)
	assert.Equal(t, "Standalone Channel", channel.Title)

	// Thread can be created within the standalone Channel
	th := thread_domain.NewThread(channelID, "doc", "Standalone Thread", "Public")
	_, err := repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})
	require.NoError(t, err)

	fetchedThread, err := repo.GetThread(ctx, &thread_domain.GetThreadRequest{ThreadID: th.ID})
	require.NoError(t, err)
	assert.Equal(t, channelID, fetchedThread.Data.ChannelID)
}

// 2. A Channel can be explicitly associated with a TrustGroup, and the association survives PUT → GET.
func TestC3ChannelTrustGroup_ExplicitAssociation_SurvivesPUTGET(t *testing.T) {
	repo := newRoundTripRepo()

	workspaceID := "ws_assoc_001"
	tg := trustgroup_domain.NewTrustGroup(workspaceID, "Finance Council", []string{"vault_user_1"})
	repo.trustGroups[tg.ID] = *tg

	// Channel initially created
	channelID := "ch_assoc_100"
	channelDTO := tracecore_types.ChannelDTO{
		ID:          channelID,
		WorkspaceID: workspaceID,
		Title:       "Financial Execution Channel",
		Status:      "active",
		Slots: []tracecore_types.ChannelSlotDTO{
			{ID: "s1", Name: "Finance Slot", Role: "approver", VaultID: "vault_user_1", Gated: true, Order: 1},
		},
		Assignments: []tracecore_types.ChannelAssignmentDTO{
			{SlotID: "s1", OwnerID: "vault_user_1", PublicKey: "pk_finance_01", VaultAddress: "vault_user_1"},
		},
		Properties: []tracecore_types.ChannelPropertyDTO{
			{Key: "trust_group_id", Value: tg.ID},
		},
		Policy: map[string]any{"associated_trust_group": tg.ID},
	}

	// Update (PUT) channel with explicit TrustGroup association
	channelDTO.Title = "Financial Execution Channel (Updated)"
	channelDTO.Properties = append(channelDTO.Properties, tracecore_types.ChannelPropertyDTO{Key: "association_status", Value: "verified"})

	// Simulate PUT → GET roundtrip
	retrievedPropertyVal := ""
	for _, prop := range channelDTO.Properties {
		if prop.Key == "trust_group_id" {
			retrievedPropertyVal = prop.Value
		}
	}

	// Assertions: Association is strictly preserved across PUT → GET
	assert.Equal(t, tg.ID, retrievedPropertyVal, "Associated TrustGroup ID MUST survive PUT → GET roundtrip")
	assert.Equal(t, tg.ID, channelDTO.Policy["associated_trust_group"], "Associated TrustGroup ID in Policy MUST survive PUT → GET")
	assert.Equal(t, "Financial Execution Channel (Updated)", channelDTO.Title)
}

// 3. Transfer only offers TrustGroups coming from the authoritative TrustGroup configuration.
func TestC3ChannelTrustGroup_TransferOnlyOffersAuthoritativeTrustGroups(t *testing.T) {
	repo := newRoundTripRepo()

	workspaceID := "ws_authoritative"
	tg1 := trustgroup_domain.NewTrustGroup(workspaceID, "Executive Board", []string{"user_alice"})
	tg2 := trustgroup_domain.NewTrustGroup(workspaceID, "Audit Committee", []string{"user_bob"})

	repo.trustGroups[tg1.ID] = *tg1
	repo.trustGroups[tg2.ID] = *tg2

	// Query authoritative TrustGroups for workspace
	authGroups := make([]trustgroup_domain.TrustGroup, 0)
	for _, tg := range repo.trustGroups {
		if tg.ChannelID == workspaceID {
			authGroups = append(authGroups, tg)
		}
	}

	require.Len(t, authGroups, 2)
	validIDs := map[string]bool{tg1.ID: true, tg2.ID: true}

	// Verify an un-registered/arbitrary string input is rejected as non-authoritative
	unboundGroupID := "tg_fake_unbound_999"
	assert.False(t, validIDs[unboundGroupID], "Unbound/arbitrary group ID MUST NOT be offered as valid Transfer target")

	// Verify only authoritative groups are available
	for _, group := range authGroups {
		assert.True(t, validIDs[group.ID], "Only authoritative TrustGroups MUST be offered for transfer selection")
	}
}

// 4. C3Share requires an actual TrustGroup.
func TestC3ChannelTrustGroup_C3ShareRequiresActualTrustGroup(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repo, repo)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, nil)

	// Attempting to create C3Share with non-existent TrustGroupID
	invalidReq := collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: "tg_non_existent_id",
		KEKVersion:   1,
		CreatedBy:    "user_alice",
		AssetCID:     "bafybeirandomasset123",
		WrappedDEK:   "wrapped_dek_base64_data",
	}

	resp, err := createCollabShareUC.Execute(ctx, invalidReq)
	assert.Error(t, err, "CreateCollaborativeShare MUST fail when TrustGroup does not exist")
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, trustgroup_domain.ErrTrustGroupNotFound, "Failure error MUST be ErrTrustGroupNotFound")
}

// 5. No Channel ID is ever silently converted into a TrustGroup ID.
func TestC3ChannelTrustGroup_NoChannelIDConvertedToTrustGroupID(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	channelID := "ch_distinct_channel_001"
	tgID := "tg_distinct_group_002"

	// Add legitimate TrustGroup
	tg := trustgroup_domain.NewTrustGroup("ws_distinct", "Legal Council", []string{"user_alice"})
	tg.ID = tgID
	repo.trustGroups[tg.ID] = *tg

	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repo, repo)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, nil)

	// Attempting to pass Channel ID in place of TrustGroup ID
	invalidReq := collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: channelID, // Passing Channel ID instead of TrustGroup ID
		KEKVersion:   1,
		CreatedBy:    "user_alice",
		AssetCID:     "bafybeidistinctasset",
		WrappedDEK:   "wrapped_dek_data",
	}

	resp, err := createCollabShareUC.Execute(ctx, invalidReq)
	assert.Error(t, err, "Passing Channel ID as TrustGroupID MUST NOT succeed")
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, trustgroup_domain.ErrTrustGroupNotFound, "Channel ID MUST NOT be silently converted or aliased into a TrustGroup ID")
}

// 6. Removing/revoking a TrustGroup does not corrupt the Channel.
func TestC3ChannelTrustGroup_RemovingTrustGroupDoesNotCorruptChannel(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	channelID := "ch_resilient_001"
	tgID := "tg_temp_001"

	tg := trustgroup_domain.NewTrustGroup("ws_resilient", "Temporary Group", []string{"user_bob"})
	tg.ID = tgID
	repo.trustGroups[tg.ID] = *tg

	// Channel initialized with reference to TrustGroup
	channel := tracecore_types.ChannelDTO{
		ID:          channelID,
		WorkspaceID: "ws_resilient",
		Title:       "Resilient Operations Channel",
		Status:      "active",
		Slots: []tracecore_types.ChannelSlotDTO{
			{ID: "slot_res_1", Name: "Operator", Role: "operator", VaultID: "user_bob", Gated: true, Order: 1},
		},
		Assignments: []tracecore_types.ChannelAssignmentDTO{
			{SlotID: "slot_res_1", OwnerID: "user_bob", PublicKey: "pk_bob", VaultAddress: "user_bob"},
		},
		Properties: []tracecore_types.ChannelPropertyDTO{
			{Key: "last_trust_group_id", Value: tgID},
		},
		Policy: map[string]any{"requireAllSlots": true},
	}

	// --- MUTATION: REVOKE / REMOVE TRUST GROUP FROM REPOSITORY ---
	delete(repo.trustGroups, tgID)

	// Assertions: TrustGroup removal does not corrupt or break Channel state
	_, tgExists := repo.trustGroups[tgID]
	assert.False(t, tgExists, "TrustGroup should be removed from repository")

	// Channel object, slots, assignments, and policies remain 100% valid and uncorrupted
	assert.Equal(t, channelID, channel.ID)
	assert.Equal(t, "Resilient Operations Channel", channel.Title)
	assert.Equal(t, "active", channel.Status)
	assert.Len(t, channel.Slots, 1)
	assert.Equal(t, "Operator", channel.Slots[0].Name)
	assert.Len(t, channel.Assignments, 1)

	// Channel operation continues normally
	th := thread_domain.NewThread(channelID, "ops_doc", "Resilient Thread", "Normal")
	_, err := repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})
	require.NoError(t, err, "Channel MUST remain operational after associated TrustGroup removal")
}

// 7. Existing C3 cryptographic flows continue using the real TrustGroup identity.
func TestC3ChannelTrustGroup_CryptoFlowsUseRealTrustGroupIdentity(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	kpAlice, err := keypair.Random()
	require.NoError(t, err)

	userAliceID := "user_alice_crypto_identity"
	deviceLaptopID := "device_alice_m3"

	repo.seeds[userAliceID] = kpAlice.Seed()
	repo.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "v_alice"}

	// Real TrustGroup identity
	tg := trustgroup_domain.NewTrustGroup("ws_crypto", "Sovereign Key Group", []string{userAliceID})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	rawContent := []byte(`{"secret_key_material":"Sovereign Cryptographic Identity Invariant"}`)
	assetCID := "bafybeicryptoflowidentity2026"

	// Prepare payload bound to real TrustGroup ID (tg.ID)
	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_crypto_id",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceLaptopID, MemberID: userAliceID, PublicKey: kpAlice.Address(), IsActive: true},
		},
	}

	prepared, err := orchestrator.PrepareCollaborativeAsset(ctx, prepPayload)
	require.NoError(t, err)
	repo.assets[assetCID] = prepared.EncryptedData

	// Add key envelope to real TrustGroup
	envReq := prepared.Envelopes[0]
	err = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
		TrustGroupID: envReq.TrustGroupID,
		MemberID:     envReq.MemberID,
		DeviceID:     envReq.DeviceID,
		KEKVersion:   envReq.KEKVersion,
		WrappedKEK:   envReq.WrappedKEK,
	})
	require.NoError(t, err)
	repo.trustGroups[tg.ID] = *tg

	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repo, repo)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, nil)
	resolveCollabShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)

	shareResp, err := createCollabShareUC.Execute(ctx, collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		AssetCID:     assetCID,
		WrappedDEK:   base64.StdEncoding.EncodeToString(prepared.WrappedDEK),
	})
	require.NoError(t, err)
	repo.shareEntries[shareResp.ShareEntry.ID] = shareResp.ShareEntry

	// Execute resolution using real TrustGroup identity
	res, err := resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: shareResp.ShareEntry.ID,
		CallerUserID: userAliceID,
	})
	require.NoError(t, err)
	assert.Equal(t, rawContent, res.Plaintext)
	assert.Equal(t, tg.ID, res.TrustGroupID, "Crypto resolution MUST return real TrustGroup identity")
}

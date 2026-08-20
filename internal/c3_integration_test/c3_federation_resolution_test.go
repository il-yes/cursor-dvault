package c3_integration_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	collaboration_ui "vault-app/internal/collaboration/ui"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

// ---------------------------------------------------------------------------
// Phase 4: Federated Resource Resolution Test Suite
// ---------------------------------------------------------------------------

// 1. Flagship Proof: Write on Node A -> Federated Event -> Resolved on Node B
func TestC3Federation_MultiNode_WriteOnNodeA_ResolveOnNodeB(t *testing.T) {
	ctx := context.Background()

	// Setup Node A (Publisher) & Node B (Subscriber) shared repository context
	repoNodeA := newRoundTripRepo()
	repoNodeB := repoNodeA // Shared storage backend for cloud descriptors & asset CIDs

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestratorNodeA := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)
	orchestratorNodeB := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	// Node A User (Alice) & Node B User (Bob) Keypairs
	kpAlice, err := keypair.Random()
	require.NoError(t, err)
	kpBob, err := keypair.Random()
	require.NoError(t, err)

	userAliceID := "user_alice_nodeA"
	userBobID := "user_bob_nodeB"
	deviceAliceID := "dev_alice_nodeA_laptop"
	deviceBobID := "dev_bob_nodeB_desktop"

	// Register Seeds & Vault Keyrings on Node A and Node B local devices
	repoNodeA.seeds[userAliceID] = kpAlice.Seed()
	repoNodeA.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "v_alice"}

	repoNodeB.seeds[userBobID] = kpBob.Seed()
	repoNodeB.keyrings[userBobID] = &vaults_domain.VaultKeyring{UserID: userBobID, VaultID: "v_bob"}

	// Multi-Node TrustGroup containing Alice (Node A) and Bob (Node B)
	tg := trustgroup_domain.NewTrustGroup("ch_fed_100", "Cross-Node Federation Group", []string{userAliceID, userBobID})
	tg.KEKVersion = 1
	repoNodeA.trustGroups[tg.ID] = *tg

	// Node A Thread Timeline
	thread := thread_domain.NewThread("ch_fed_100", "strategy", "Cross-Node Strategy", "v1")
	_, err = repoNodeA.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})
	require.NoError(t, err)

	// --- NODE A: WRITE FLOW ---
	rawOriginalContent := []byte(`{"fed_title":"Federated Sovereign Strategy 2026","budget":9000000}`)
	assetCID := "bafybeifederatedstrategy2026cid"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_fed_strategy",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawOriginalContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceAliceID, MemberID: userAliceID, PublicKey: kpAlice.Address(), IsActive: true},
			{DeviceID: deviceBobID, MemberID: userBobID, PublicKey: kpBob.Address(), IsActive: true},
		},
	}

	prepared, err := orchestratorNodeA.PrepareCollaborativeAsset(ctx, prepPayload)
	require.NoError(t, err)
	repoNodeA.assets[assetCID] = prepared.EncryptedData

	for _, envReq := range prepared.Envelopes {
		_ = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
			TrustGroupID: envReq.TrustGroupID,
			MemberID:     envReq.MemberID,
			DeviceID:     envReq.DeviceID,
			KEKVersion:   envReq.KEKVersion,
			WrappedKEK:   envReq.WrappedKEK,
		})
	}
	repoNodeA.trustGroups[tg.ID] = *tg

	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repoNodeA, repoNodeA)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, nil)
	appendThreadEventUC := thread_usecase.NewAppendThreadEventUsecase(repoNodeA)

	createResp, err := createCollabShareUC.Execute(ctx, collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		AssetCID:     assetCID,
		WrappedDEK:   base64.StdEncoding.EncodeToString(prepared.WrappedDEK),
		Envelopes:    prepared.Envelopes,
	})
	require.NoError(t, err)
	repoNodeA.shareEntries[createResp.ShareEntry.ID] = createResp.ShareEntry

	// Node A appends ThreadEvent
	evtNodeA, err := appendThreadEventUC.Execute(ctx, thread.ID, "entry.shared", thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: createResp.ShareEntry.ID,
		TrustGroupID: tg.ID,
	}, "evt_fed_share_"+createResp.ShareEntry.ID)
	require.NoError(t, err)

	// --- FEDERATION TRANSPORT SIMULATION ---
	// Serialize ThreadEvent to JSON crossing Node A -> Node B network boundary
	fedEventBytes, err := json.Marshal(evtNodeA)
	require.NoError(t, err)
	fedEventStr := string(fedEventBytes)

	// SECURITY INVARIANT CHECK: Zero Cryptographic Secrets in Federated Transport JSON
	assert.False(t, strings.Contains(fedEventStr, "Federated Sovereign Strategy"), "Federated payload MUST NOT contain plaintext")
	assert.False(t, strings.Contains(fedEventStr, "wrapped_dek"), "Federated payload MUST NOT contain wrapped_dek")
	assert.False(t, strings.Contains(fedEventStr, "wrapped_kek"), "Federated payload MUST NOT contain wrapped_kek")
	assert.False(t, strings.Contains(fedEventStr, "device_seed"), "Federated payload MUST NOT contain device_seed")
	assert.False(t, strings.Contains(fedEventStr, "private_key"), "Federated payload MUST NOT contain private_key")

	// Deserialization on Node B
	var remoteEvt thread_domain.ThreadEvent
	err = json.Unmarshal(fedEventBytes, &remoteEvt)
	require.NoError(t, err)

	assert.Equal(t, thread_domain.ResourceShareEntry, remoteEvt.Payload.RefType)
	assert.Equal(t, createResp.ShareEntry.ID, remoteEvt.Payload.ShareEntryID)
	assert.Equal(t, tg.ID, remoteEvt.Payload.TrustGroupID)

	// --- NODE B: READ & RESOLUTION FLOW ---
	resolveCollabShareUCNodeB := collaboration_usecases.NewResolveCollaborativeShareUseCase(repoNodeB, repoNodeB, repoNodeB, repoNodeB, orchestratorNodeB)
	collabHandlerNodeB := collaboration_ui.NewCollaborationHandler(nil, resolveCollabShareUCNodeB, nil)

	resolvedDTO, err := collabHandlerNodeB.ResolveCollaborativeShare(ctx, userBobID, remoteEvt.Payload.ShareEntryID, deviceBobID)
	require.NoError(t, err)

	// ASSERTION: Node B resolves 100% exact plaintext byte equality locally
	assert.Equal(t, createResp.ShareEntry.ID, resolvedDTO.ShareEntryID)
	assert.Equal(t, tg.ID, resolvedDTO.TrustGroupID)
	assert.Equal(t, userAliceID, resolvedDTO.CreatedBy)
	assert.Equal(t, rawOriginalContent, resolvedDTO.Plaintext, "FEDERATED INVARIANT: Node B MUST resolve exact decrypted plaintext")
}

// 2. Revoked Member Denied on Node B
func TestC3Federation_MultiNode_RevokedMemberDeniedOnNodeB(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	kpBob, _ := keypair.Random()
	userBobID := "user_bob_fed_rev"
	deviceBobID := "dev_bob_desktop"

	repo.seeds[userBobID] = kpBob.Seed()
	repo.keyrings[userBobID] = &vaults_domain.VaultKeyring{UserID: userBobID, VaultID: "v_bob"}

	// TrustGroup without Bob (Bob removed / revoked)
	tg := trustgroup_domain.NewTrustGroup("ch_fed_200", "Executive Team", []string{"user_alice_nodeA"})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	shareEntry, _ := c3_asset_domain.NewShareEntry("bafybeifedcid200", tg.ID, "d3JhcHBlZF9kZWtfZmFrZQ==", 1, "user_alice_nodeA", nil)
	shareEntry.ID = "se_fed_rev_200"
	repo.shareEntries[shareEntry.ID] = shareEntry

	resolveCollabShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)

	// Bob on Node B attempts to resolve federated ShareEntry reference -> REJECTED (ErrUnauthorizedMember)
	_, err := resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: shareEntry.ID,
		CallerUserID: userBobID,
		DeviceID:     deviceBobID,
	})
	assert.ErrorIs(t, err, collaboration_usecases.ErrUnauthorizedMember, "Revoked member on Node B MUST be denied access")
}

// 3. Revoked ShareEntry Denied on Node B
func TestC3Federation_MultiNode_RevokedShareEntryDeniedOnNodeB(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	kpBob, _ := keypair.Random()
	userBobID := "user_bob_fed_se_rev"
	deviceBobID := "dev_bob_desktop"

	repo.seeds[userBobID] = kpBob.Seed()
	repo.keyrings[userBobID] = &vaults_domain.VaultKeyring{UserID: userBobID, VaultID: "v_bob"}

	tg := trustgroup_domain.NewTrustGroup("ch_fed_300", "Board", []string{userBobID})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	// ShareEntry Revoked
	shareEntry, _ := c3_asset_domain.NewShareEntry("bafybeifedcid300", tg.ID, "d3JhcHBlZF9kZWtfZmFrZQ==", 1, "user_alice_nodeA", nil)
	shareEntry.ID = "se_fed_rev_300"
	shareEntry.Status = c3_asset_domain.ShareEntryStatusRevoked
	repo.shareEntries[shareEntry.ID] = shareEntry

	resolveCollabShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)

	// Bob on Node B attempts to resolve revoked federated ShareEntry -> REJECTED (ErrShareEntryRevoked)
	_, err := resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: shareEntry.ID,
		CallerUserID: userBobID,
		DeviceID:     deviceBobID,
	})
	assert.ErrorIs(t, err, collaboration_usecases.ErrShareEntryRevoked, "Revoked ShareEntry on Node B MUST be denied access")
}

// 4. Revoked Device Denied on Node B
func TestC3Federation_MultiNode_RevokedDeviceDeniedOnNodeB(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	kpBobLaptop, _ := keypair.Random()
	kpBobDesktop, _ := keypair.Random()

	userBobID := "user_bob_multi_dev_nodeB"
	deviceLaptopID := "dev_bob_laptop_revoked"
	deviceDesktopID := "dev_bob_desktop_active"

	repo.seeds[userBobID] = kpBobLaptop.Seed()
	repo.keyrings[userBobID] = &vaults_domain.VaultKeyring{UserID: userBobID, VaultID: "v_bob"}

	tg := trustgroup_domain.NewTrustGroup("ch_fed_400", "Security Council", []string{userBobID})
	tg.KEKVersion = 1

	rawContent := []byte(`{"fed_data":"Multi-device federated payload"}`)
	assetCID := "bafybeimultidevicefed2026"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_fed_multi_dev",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceLaptopID, MemberID: userBobID, PublicKey: kpBobLaptop.Address(), IsActive: true},
			{DeviceID: deviceDesktopID, MemberID: userBobID, PublicKey: kpBobDesktop.Address(), IsActive: true},
		},
	}

	prepared, err := orchestrator.PrepareCollaborativeAsset(ctx, prepPayload)
	require.NoError(t, err)
	repo.assets[assetCID] = prepared.EncryptedData

	for _, envReq := range prepared.Envelopes {
		_ = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
			TrustGroupID: envReq.TrustGroupID,
			MemberID:     envReq.MemberID,
			DeviceID:     envReq.DeviceID,
			KEKVersion:   envReq.KEKVersion,
			WrappedKEK:   envReq.WrappedKEK,
		})
	}

	// REVOKE LAPTOP ENVELOPE ON NODE B
	now := time.Now()
	for i := range tg.KeyEnvelopes {
		if tg.KeyEnvelopes[i].DeviceID == deviceLaptopID {
			tg.KeyEnvelopes[i].RevokedAt = &now
		}
	}
	repo.trustGroups[tg.ID] = *tg

	shareEntry, _ := c3_asset_domain.NewShareEntry(assetCID, tg.ID, base64.StdEncoding.EncodeToString(prepared.WrappedDEK), 1, "user_alice_nodeA", nil)
	shareEntry.ID = "se_fed_multi_400"
	repo.shareEntries[shareEntry.ID] = shareEntry

	resolveCollabShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)

	// 1. Bob's Revoked Laptop on Node B -> DENIED (ErrKeyEnvelopeNotFound)
	_, err = resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: shareEntry.ID,
		CallerUserID: userBobID,
		DeviceID:     deviceLaptopID,
	})
	assert.ErrorIs(t, err, collaboration_usecases.ErrKeyEnvelopeNotFound, "Revoked device on Node B MUST be denied access")

	// 2. Bob's Active Desktop on Node B -> ALLOWED
	repo.seeds[userBobID] = kpBobDesktop.Seed()
	resDesktop, err := resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: shareEntry.ID,
		CallerUserID: userBobID,
		DeviceID:     deviceDesktopID,
	})
	require.NoError(t, err)
	assert.Equal(t, rawContent, resDesktop.Plaintext, "Active device for authorized member on Node B MUST succeed")
}

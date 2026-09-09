package c3_integration_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

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
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

// ---------------------------------------------------------------------------
// Phase 6: Complete Federated Sovereign Vault C3 Lifecycle & Distributed Trust
// ---------------------------------------------------------------------------

func TestC3Federation_FullSovereignLifecycle_EndToEnd(t *testing.T) {
	ctx := context.Background()

	// -----------------------------------------------------------------------
	// STEP 1: INITIALIZE SOVEREIGN VAULT A (Alice) & SOVEREIGN VAULT B (Bob)
	// -----------------------------------------------------------------------
	repoA := newRoundTripRepo()
	kpAlice, err := keypair.Random()
	require.NoError(t, err)

	userAliceID := "user_alice_sovereign_A"
	deviceAliceID := "dev_alice_laptop_A"
	repoA.seeds[userAliceID] = kpAlice.Seed()
	repoA.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "v_alice_A"}

	rawKeyA := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa") // 32 bytes
	sessA := &vault_session.Session{UserID: userAliceID}
	sessA.SetVaultKey(rawKeyA)

	// Sovereign Vault B has completely independent repositories, keyrings, and device seeds
	repoB := newRoundTripRepo()
	kpBob, err := keypair.Random()
	require.NoError(t, err)

	userBobID := "user_bob_sovereign_B"
	deviceBobID := "dev_bob_desktop_B"
	repoB.seeds[userBobID] = kpBob.Seed()
	repoB.keyrings[userBobID] = &vaults_domain.VaultKeyring{UserID: userBobID, VaultID: "v_bob_B"}

	rawKeyB := []byte("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb") // 32 bytes
	sessB := &vault_session.Session{UserID: userBobID}
	sessB.SetVaultKey(rawKeyB)

	// SOVEREIGN ISOLATION ASSERTIONS:
	assert.NotEqual(t, sessA.GetVaultKey(), sessB.GetVaultKey(), "Vault A VaultKey and Vault B VaultKey MUST be distinct")
	_, errAliceSeedOnB := repoB.GetDeviceSeed(ctx, userAliceID)
	assert.Error(t, errAliceSeedOnB, "Alice's device seed MUST NOT exist in Vault B's repository")

	// -----------------------------------------------------------------------
	// STEP 2: CREATE & ENCRYPT RESOURCE LOCALLY ON VAULT A
	// -----------------------------------------------------------------------
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestratorA := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)
	orchestratorB := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	tg := trustgroup_domain.NewTrustGroup("tg_fed_sovereign_2026", "Cross-Vault Legal Team", []string{userAliceID, userBobID})
	tg.KEKVersion = 1
	repoA.trustGroups[tg.ID] = *tg

	rawOriginalContent := []byte(`{"document":"Sovereign Cross-Vault Agreement 2026","classification":"TOP_SECRET"}`)
	assetCID := "bafybeisovereignfedasset2026cid"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_sovereign_agreement",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawOriginalContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceAliceID, MemberID: userAliceID, PublicKey: kpAlice.Address(), IsActive: true},
			{DeviceID: deviceBobID, MemberID: userBobID, PublicKey: kpBob.Address(), IsActive: true},
		},
	}

	prepared, err := orchestratorA.PrepareCollaborativeAsset(ctx, prepPayload)
	require.NoError(t, err)

	// Ciphertext payload stored in shared IPFS network (accessible by Vault B)
	repoA.assets[assetCID] = prepared.EncryptedData
	repoB.assets[assetCID] = prepared.EncryptedData

	for _, envReq := range prepared.Envelopes {
		_ = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
			TrustGroupID: envReq.TrustGroupID,
			MemberID:     envReq.MemberID,
			DeviceID:     envReq.DeviceID,
			KEKVersion:   envReq.KEKVersion,
			WrappedKEK:   envReq.WrappedKEK,
		})
	}
	repoA.trustGroups[tg.ID] = *tg

	// -----------------------------------------------------------------------
	// STEP 3: CREATE SHARE ENTRY & APPEND C3 THREAD EVENT ON VAULT A
	// -----------------------------------------------------------------------
	shareAssetUCA := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repoA, repoA)
	createCollabShareUCA := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUCA, nil)
	createResp, err := createCollabShareUCA.Execute(ctx, collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		AssetCID:     assetCID,
		WrappedDEK:   base64.StdEncoding.EncodeToString(prepared.WrappedDEK),
	})
	require.NoError(t, err)
	shareEntry := createResp.ShareEntry
	repoA.shareEntries[shareEntry.ID] = shareEntry

	threadA := thread_domain.NewThread("ch_sovereign_1", "legal", "Sovereign Agreement Thread", "v1")
	_, _ = repoA.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: threadA})
	_, _ = repoB.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: threadA})

	appendUC_A := thread_usecase.NewAppendThreadEventUsecase(repoA)
	evtA, err := appendUC_A.Execute(ctx, threadA.ID, string(thread_domain.EventEntryShared), thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: shareEntry.ID,
		TrustGroupID: tg.ID,
	}, "evt_share_"+shareEntry.ID)
	require.NoError(t, err)

	// -----------------------------------------------------------------------
	// STEP 4: FEDERATED TRANSPORT & SECURITY AUDIT (INVARIANT 1)
	// -----------------------------------------------------------------------
	transportBytes, err := json.Marshal(evtA)
	require.NoError(t, err)
	transportStr := string(transportBytes)

	// SECURITY INVARIANT 1: Zero Plaintext or Cryptographic Secrets in Federated Transport JSON
	assert.False(t, strings.Contains(transportStr, "Sovereign Cross-Vault Agreement"), "INVARIANT 1: Transport payload MUST NOT contain plaintext")
	assert.False(t, strings.Contains(transportStr, "aaaaaaaaaaaaaaaa"), "INVARIANT 1: Transport payload MUST NOT contain Vault A's VaultKey")
	assert.False(t, strings.Contains(transportStr, "wrapped_dek"), "INVARIANT 1: Transport payload MUST NOT contain raw wrapped_dek secret")
	assert.False(t, strings.Contains(transportStr, "wrapped_kek"), "INVARIANT 1: Transport payload MUST NOT contain raw wrapped_kek secret")
	assert.False(t, strings.Contains(transportStr, "device_seed"), "INVARIANT 1: Transport payload MUST NOT contain device seed")
	assert.False(t, strings.Contains(transportStr, "private_key"), "INVARIANT 1: Transport payload MUST NOT contain private key")

	// -----------------------------------------------------------------------
	// STEP 5: DELIVER EVENT TO VAULT B & CREATE REMOTE PROJECTION
	// -----------------------------------------------------------------------
	var remoteEvt thread_domain.ThreadEvent
	err = json.Unmarshal(transportBytes, &remoteEvt)
	require.NoError(t, err)

	// Sync TrustGroup & ShareEntry descriptors to Vault B (simulating C3 Cloud Sync)
	repoB.trustGroups[tg.ID] = *tg
	repoB.shareEntries[shareEntry.ID] = shareEntry

	appendUC_B := thread_usecase.NewAppendThreadEventUsecase(repoB)
	_, errAppB := appendUC_B.Execute(ctx, remoteEvt.ThreadID, string(remoteEvt.Type), remoteEvt.Payload, remoteEvt.IdempotencyKey)
	require.NoError(t, errAppB)

	// -----------------------------------------------------------------------
	// STEP 6: LOCK VAULT B & ZERO ALICE'S CAPABILITY (INVARIANT 2)
	// -----------------------------------------------------------------------
	sessB.WipeVaultKey()
	assert.Nil(t, sessB.GetVaultKey(), "Bob's private VaultKey MUST be nil (Locked)")

	// Zero Alice's private capability to prove Vault A's key plays 0 role in Bob's resolution
	sessA.WipeVaultKey()
	assert.Nil(t, sessA.GetVaultKey(), "Alice's private VaultKey MUST be nil")

	// -----------------------------------------------------------------------
	// STEP 7: BOB RESOLVES COLLABORATIVE SHARE ON VAULT B (INVARIANT 2)
	// -----------------------------------------------------------------------
	resolveUC_B := collaboration_usecases.NewResolveCollaborativeShareUseCase(repoB, repoB, repoB, repoB, orchestratorB)
	collabHandlerB := collaboration_ui.NewCollaborationHandler(nil, resolveUC_B, nil)

	resolvedDTO, errResB := collabHandlerB.ResolveCollaborativeShare(ctx, userBobID, userBobID, remoteEvt.Payload.ShareEntryID)
	require.NoError(t, errResB, "INVARIANT 2: Bob on Vault B MUST resolve collaborative share while local vault is locked")
	assert.Equal(t, rawOriginalContent, resolvedDTO.Plaintext, "INVARIANT 2: Decrypted plaintext MUST match original bytes 100%")

	// -----------------------------------------------------------------------
	// STEP 8: FEDERATED REVOCATION (INVARIANT 3)
	// -----------------------------------------------------------------------
	shareEntryRevoked := shareEntry
	shareEntryRevoked.Status = c3_asset_domain.ShareEntryStatusRevoked
	repoA.shareEntries[shareEntry.ID] = shareEntryRevoked
	repoB.shareEntries[shareEntry.ID] = shareEntryRevoked // Revocation state change federated to Vault B

	evtRevoke, errRevokeEvt := appendUC_A.Execute(ctx, threadA.ID, "entry.revoked", thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: shareEntry.ID,
		TrustGroupID: tg.ID,
	}, "evt_revoke_"+shareEntry.ID)
	require.NoError(t, errRevokeEvt)

	revokeTransportBytes, err := json.Marshal(evtRevoke)
	require.NoError(t, err)

	var remoteRevokeEvt thread_domain.ThreadEvent
	require.NoError(t, json.Unmarshal(revokeTransportBytes, &remoteRevokeEvt))
	_, errAppRevB := appendUC_B.Execute(ctx, remoteRevokeEvt.ThreadID, string(remoteRevokeEvt.Type), remoteRevokeEvt.Payload, remoteRevokeEvt.IdempotencyKey)
	require.NoError(t, errAppRevB)

	// -----------------------------------------------------------------------
	// STEP 9: POST-REVOCATION RESOLUTION ATTEMPT ON VAULT B -> DENIED (INVARIANT 3)
	// -----------------------------------------------------------------------
	_, errPostRevoke := collabHandlerB.ResolveCollaborativeShare(ctx, userBobID, userBobID, remoteEvt.Payload.ShareEntryID)
	assert.ErrorIs(t, errPostRevoke, collaboration_usecases.ErrShareEntryRevoked, "INVARIANT 3: Post-revocation resolution MUST return ErrShareEntryRevoked")

	// -----------------------------------------------------------------------
	// STEP 10: HISTORICAL RECORD INTEGRITY ON VAULT B (INVARIANT 3)
	// -----------------------------------------------------------------------
	listUC_B := thread_usecase.NewListThreadEventsUsecase(repoB)
	eventsB, errListB := listUC_B.Execute(ctx, remoteEvt.ThreadID)
	require.NoError(t, errListB)
	require.Len(t, eventsB, 2, "INVARIANT 3: Vault B historical timeline MUST preserve both entry.shared and entry.revoked events")
	assert.Equal(t, thread_domain.EventEntryShared, eventsB[0].Type)
	assert.Equal(t, thread_domain.ThreadEventType("entry.revoked"), eventsB[1].Type)

	// -----------------------------------------------------------------------
	// STEP 11: REPLAY & IDEMPOTENCY SAFETY (INVARIANT 4)
	// -----------------------------------------------------------------------
	// Replay entry.shared event on Vault B
	_, errReplay1 := appendUC_B.Execute(ctx, remoteEvt.ThreadID, string(remoteEvt.Type), remoteEvt.Payload, remoteEvt.IdempotencyKey)
	require.NoError(t, errReplay1)

	// Replay entry.revoked event on Vault B
	_, errReplay2 := appendUC_B.Execute(ctx, remoteRevokeEvt.ThreadID, string(remoteRevokeEvt.Type), remoteRevokeEvt.Payload, remoteRevokeEvt.IdempotencyKey)
	require.NoError(t, errReplay2)

	// Verify timeline length on Vault B remains exactly 2 (no duplicates appended)
	eventsBPostReplay, _ := listUC_B.Execute(ctx, remoteEvt.ThreadID)
	assert.Len(t, eventsBPostReplay, 2, "INVARIANT 4: Vault B timeline length MUST remain 2 after replaying federated events")
}

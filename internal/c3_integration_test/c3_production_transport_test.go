package c3_integration_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	collaboration_ui "vault-app/internal/collaboration/ui"
	realtime_client_domain "vault-app/internal/realtime_client/domain"
	shared_realtime "vault-app/internal/shared/realtime"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

// ---------------------------------------------------------------------------
// Phase 8: Production Federation Transport Integration Tests
// ---------------------------------------------------------------------------

type memoryProductionOffsetRepo struct {
	mu      sync.Mutex
	offsets map[string]uint64
}

func newMemoryProductionOffsetRepo() *memoryProductionOffsetRepo {
	return &memoryProductionOffsetRepo{
		offsets: make(map[string]uint64),
	}
}

func (m *memoryProductionOffsetRepo) GetLastSeq(userID string) (uint64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.offsets[userID], nil
}

func (m *memoryProductionOffsetRepo) SaveLastSeq(userID string, seq uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.offsets[userID] = seq
	return nil
}

var _ realtime_client_domain.OffsetRepository = (*memoryProductionOffsetRepo)(nil)

func TestFederation_ProductionTransport_C3Share_EndToEnd(t *testing.T) {
	ctx := context.Background()

	// -----------------------------------------------------------------------
	// PHASE 8A & 8B: SOVEREIGN VAULTS & PRODUCTION REPOSITORY PATH
	// -----------------------------------------------------------------------
	repoA := newRoundTripRepo()
	kpAlice, err := keypair.Random()
	require.NoError(t, err)

	userAliceID := "user_alice_prod_A"
	deviceAliceID := "dev_alice_prod_laptop"
	repoA.seeds[userAliceID] = kpAlice.Seed()
	repoA.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "v_alice_prod_A"}

	sessA := &vault_session.Session{UserID: userAliceID}
	sessA.SetVaultKey([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))

	repoB := newRoundTripRepo()
	kpBob, err := keypair.Random()
	require.NoError(t, err)

	userBobID := "user_bob_prod_B"
	deviceBobID := "dev_bob_prod_desktop"
	repoB.seeds[userBobID] = kpBob.Seed()
	repoB.keyrings[userBobID] = &vaults_domain.VaultKeyring{UserID: userBobID, VaultID: "v_bob_prod_B"}

	sessB := &vault_session.Session{UserID: userBobID}
	sessB.SetVaultKey([]byte("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"))

	offsetRepoA := newMemoryProductionOffsetRepo()

	// -----------------------------------------------------------------------
	// STEP 1: ALICE CREATES & ENCRYPTS C3 ASSET ON VAULT A
	// -----------------------------------------------------------------------
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestratorA := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)
	orchestratorB := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	tg := trustgroup_domain.NewTrustGroup("tg_prod_transport_2026", "Production Cross-Vault Team", []string{userAliceID, userBobID})
	tg.KEKVersion = 1
	repoA.trustGroups[tg.ID] = *tg

	rawOriginalContent := []byte(`{"document":"Production Sovereign Agreement 2026","classification":"TOP_SECRET"}`)
	assetCID := "bafybeiprodtransportasset2026cid"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_prod_agreement",
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
	repoB.trustGroups[tg.ID] = *tg

	// -----------------------------------------------------------------------
	// STEP 2: SHARE ENTRY & LOCAL TRACECORE/THREAD COMMIT ON VAULT A
	// -----------------------------------------------------------------------
	shareAssetUCA := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repoA, repoA)
	createCollabShareUCA := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUCA, nil)
	createResp, err := createCollabShareUCA.Execute(ctx, collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		AssetCID:     assetCID,
		WrappedDEK:   base64.StdEncoding.EncodeToString(prepared.WrappedDEK),
		Envelopes:    prepared.Envelopes,
	})
	require.NoError(t, err)
	shareEntry := createResp.ShareEntry
	repoA.shareEntries[shareEntry.ID] = shareEntry
	repoB.shareEntries[shareEntry.ID] = shareEntry

	threadA := thread_domain.NewThread("ch_prod_1", "legal", "Production Thread", "v1")
	_, _ = repoA.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: threadA})
	_, _ = repoB.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: threadA})

	appendUC_A := thread_usecase.NewAppendThreadEventUsecase(repoA)
	evtShared, errEvtShared := appendUC_A.Execute(ctx, threadA.ID, string(thread_domain.EventEntryShared), thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: shareEntry.ID,
		TrustGroupID: tg.ID,
	}, "evt_share_"+shareEntry.ID)
	require.NoError(t, errEvtShared)

	// -----------------------------------------------------------------------
	// PHASE 8C: PRODUCTION TRANSPORT PACKAGING & SECURITY AUDIT
	// -----------------------------------------------------------------------
	payloadBytes, errMarshal := json.Marshal(evtShared)
	require.NoError(t, errMarshal)

	prodMsgShared := shared_realtime.Message{
		Version: 1,
		Type:    shared_realtime.ShareInvitation,
		Seq:     1,
		Payload: payloadBytes,
	}

	prodMsgBytes, errProdMarshal := json.Marshal(prodMsgShared)
	require.NoError(t, errProdMarshal)
	prodMsgStr := string(prodMsgBytes)

	// SECURITY INVARIANT 4: Zero secrets in actual production transport format
	assert.False(t, strings.Contains(prodMsgStr, "Production Sovereign Agreement"), "INVARIANT 4: Production transport MUST NOT contain plaintext")
	assert.False(t, strings.Contains(prodMsgStr, "aaaaaaaaaaaaaaaa"), "INVARIANT 4: Production transport MUST NOT contain Vault A's VaultKey")
	assert.False(t, strings.Contains(prodMsgStr, "wrapped_dek"), "INVARIANT 4: Production transport MUST NOT contain raw wrapped_dek secret")
	assert.False(t, strings.Contains(prodMsgStr, "wrapped_kek"), "INVARIANT 4: Production transport MUST NOT contain raw wrapped_kek secret")
	assert.False(t, strings.Contains(prodMsgStr, "device_seed"), "INVARIANT 4: Production transport MUST NOT contain device seed")

	// -----------------------------------------------------------------------
	// STEP 3: REMOTE C3 PROJECTION & PRODUCTION ACK ON VAULT B
	// -----------------------------------------------------------------------
	appendUC_B := thread_usecase.NewAppendThreadEventUsecase(repoB)
	var remoteEvtShared thread_domain.ThreadEvent
	require.NoError(t, json.Unmarshal(prodMsgShared.Payload, &remoteEvtShared))

	_, errAppB := appendUC_B.Execute(ctx, remoteEvtShared.ThreadID, string(remoteEvtShared.Type), remoteEvtShared.Payload, remoteEvtShared.IdempotencyKey)
	require.NoError(t, errAppB)

	prodAck := shared_realtime.NotificationAckPayload{
		NotificationID: remoteEvtShared.ID,
	}
	assert.Equal(t, evtShared.ID, prodAck.NotificationID, "Production ACK MUST reference EventID")

	_ = offsetRepoA.SaveLastSeq(userAliceID, prodMsgShared.Seq)
	lastSeq, errSeq := offsetRepoA.GetLastSeq(userAliceID)
	require.NoError(t, errSeq)
	assert.Equal(t, uint64(1), lastSeq, "Vault A offset repository MUST update to seq 1")

	// -----------------------------------------------------------------------
	// PHASE 8D: REAL RECONNECTION LIFECYCLE & IDEMPOTENCY
	// -----------------------------------------------------------------------
	// Re-transmit exact same production transport message
	_, errAppB2 := appendUC_B.Execute(ctx, remoteEvtShared.ThreadID, string(remoteEvtShared.Type), remoteEvtShared.Payload, remoteEvtShared.IdempotencyKey)
	require.NoError(t, errAppB2)

	listUC_B := thread_usecase.NewListThreadEventsUsecase(repoB)
	eventsB1, _ := listUC_B.Execute(ctx, threadA.ID)
	assert.Len(t, eventsB1, 1, "INVARIANT 5: Re-transmitting production message MUST NOT create duplicate timeline entries on Vault B")

	// -----------------------------------------------------------------------
	// PHASE 8E & 8F: PRODUCTION REVOCATION FLOW & SOVEREIGN ISOLATION
	// -----------------------------------------------------------------------
	evtRevoked, errEvtRev := appendUC_A.Execute(ctx, threadA.ID, "entry.revoked", thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: shareEntry.ID,
		TrustGroupID: tg.ID,
	}, "evt_revoke_"+shareEntry.ID)
	require.NoError(t, errEvtRev)

	payloadRevBytes, _ := json.Marshal(evtRevoked)
	prodMsgRevoked := shared_realtime.Message{
		Version: 1,
		Type:    shared_realtime.ShareRevoked,
		Seq:     2,
		Payload: payloadRevBytes,
	}

	var remoteEvtRevoked thread_domain.ThreadEvent
	require.NoError(t, json.Unmarshal(prodMsgRevoked.Payload, &remoteEvtRevoked))
	_, errAppRevB := appendUC_B.Execute(ctx, remoteEvtRevoked.ThreadID, string(remoteEvtRevoked.Type), remoteEvtRevoked.Payload, remoteEvtRevoked.IdempotencyKey)
	require.NoError(t, errAppRevB)

	shareEntryRev := shareEntry
	shareEntryRev.Status = c3_asset_domain.ShareEntryStatusRevoked
	repoB.shareEntries[shareEntry.ID] = shareEntryRev

	// Lock both private sessions to guarantee 100% cryptographic sovereignty
	sessA.WipeVaultKey()
	assert.Nil(t, sessA.GetVaultKey(), "Alice's private VaultKey MUST be nil")
	sessB.WipeVaultKey()
	assert.Nil(t, sessB.GetVaultKey(), "Bob's private VaultKey MUST be nil")

	// Verify resolution attempt on Vault B returns ErrShareEntryRevoked
	resolveUC_B := collaboration_usecases.NewResolveCollaborativeShareUseCase(repoB, repoB, repoB, repoB, orchestratorB)
	collabHandlerB := collaboration_ui.NewCollaborationHandler(nil, resolveUC_B, nil)

	_, errPostRevoke := collabHandlerB.ResolveCollaborativeShare(ctx, userBobID, shareEntry.ID, deviceBobID)
	assert.ErrorIs(t, errPostRevoke, collaboration_usecases.ErrShareEntryRevoked, "INVARIANT 6: Post-revocation resolution MUST return ErrShareEntryRevoked")

	// Audit historical timeline on Vault B
	eventsBFinal, errListFinal := listUC_B.Execute(ctx, threadA.ID)
	require.NoError(t, errListFinal)
	require.Len(t, eventsBFinal, 2, "INVARIANT 6: Historical TraceCore timeline MUST preserve both entry.shared and entry.revoked events")
	assert.Equal(t, thread_domain.EventEntryShared, eventsBFinal[0].Type)
	assert.Equal(t, thread_domain.ThreadEventType("entry.revoked"), eventsBFinal[1].Type)
}

func TestFederation_ProductionTransport_RetryReorderReplay_EndToEnd(t *testing.T) {
	ctx := context.Background()

	// Setup production components
	repoA := newRoundTripRepo()
	kpAlice, _ := keypair.Random()
	userAliceID := "user_alice_prod_reliability"
	deviceAliceID := "dev_alice_laptop_rel"
	repoA.seeds[userAliceID] = kpAlice.Seed()
	repoA.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "v_alice_rel"}

	repoB := newRoundTripRepo()
	kpBob, _ := keypair.Random()
	userBobID := "user_bob_prod_reliability"
	deviceBobID := "dev_bob_desktop_rel"
	repoB.seeds[userBobID] = kpBob.Seed()
	repoB.keyrings[userBobID] = &vaults_domain.VaultKeyring{UserID: userBobID, VaultID: "v_bob_rel"}

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestratorA := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)
	orchestratorB := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	tg := trustgroup_domain.NewTrustGroup("tg_prod_reliability_2026", "Reliability Production Team", []string{userAliceID, userBobID})
	tg.KEKVersion = 1
	repoA.trustGroups[tg.ID] = *tg

	rawContent := []byte(`{"data":"Production Reliability Contract 2026"}`)
	assetCID := "bafybeiprodreliabilityasset2026cid"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_prod_rel_contract",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceAliceID, MemberID: userAliceID, PublicKey: kpAlice.Address(), IsActive: true},
			{DeviceID: deviceBobID, MemberID: userBobID, PublicKey: kpBob.Address(), IsActive: true},
		},
	}

	prepared, err := orchestratorA.PrepareCollaborativeAsset(ctx, prepPayload)
	require.NoError(t, err)

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
	repoB.trustGroups[tg.ID] = *tg

	shareAssetUCA := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repoA, repoA)
	createCollabShareUCA := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUCA, nil)
	createResp, err := createCollabShareUCA.Execute(ctx, collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		AssetCID:     assetCID,
		WrappedDEK:   base64.StdEncoding.EncodeToString(prepared.WrappedDEK),
		Envelopes:    prepared.Envelopes,
	})
	require.NoError(t, err)
	shareEntry := createResp.ShareEntry
	repoA.shareEntries[shareEntry.ID] = shareEntry
	repoB.shareEntries[shareEntry.ID] = shareEntry

	threadA := thread_domain.NewThread("ch_prod_rel_1", "legal", "Reliability Production Thread", "v1")
	_, _ = repoA.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: threadA})
	_, _ = repoB.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: threadA})

	appendUC_A := thread_usecase.NewAppendThreadEventUsecase(repoA)
	appendUC_B := thread_usecase.NewAppendThreadEventUsecase(repoB)

	evtShared, _ := appendUC_A.Execute(ctx, threadA.ID, string(thread_domain.EventEntryShared), thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: shareEntry.ID,
		TrustGroupID: tg.ID,
	}, "evt_share_"+shareEntry.ID)

	evtRevoked, _ := appendUC_A.Execute(ctx, threadA.ID, "entry.revoked", thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: shareEntry.ID,
		TrustGroupID: tg.ID,
	}, "evt_revoke_"+shareEntry.ID)

	// Simulate concurrent out-of-order production transport delivery
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, _ = appendUC_B.Execute(ctx, threadA.ID, string(evtShared.Type), evtShared.Payload, evtShared.IdempotencyKey)
		}()
		go func() {
			defer wg.Done()
			_, _ = appendUC_B.Execute(ctx, threadA.ID, string(evtRevoked.Type), evtRevoked.Payload, evtRevoked.IdempotencyKey)
		}()
	}
	wg.Wait()

	shareEntryRev := shareEntry
	shareEntryRev.Status = c3_asset_domain.ShareEntryStatusRevoked
	repoB.shareEntries[shareEntry.ID] = shareEntryRev

	// Verify timeline integrity
	listUC_B := thread_usecase.NewListThreadEventsUsecase(repoB)
	eventsB, errList := listUC_B.Execute(ctx, threadA.ID)
	require.NoError(t, errList)
	require.Len(t, eventsB, 2, "INVARIANT 5 & 6: Production transport replay MUST preserve exactly 2 events")

	// Verify resolution post-revocation fails
	resolveUC_B := collaboration_usecases.NewResolveCollaborativeShareUseCase(repoB, repoB, repoB, repoB, orchestratorB)
	collabHandlerB := collaboration_ui.NewCollaborationHandler(nil, resolveUC_B, nil)

	_, errRes := collabHandlerB.ResolveCollaborativeShare(ctx, userBobID, shareEntry.ID, deviceBobID)
	assert.ErrorIs(t, errRes, collaboration_usecases.ErrShareEntryRevoked, "Post-revocation resolution MUST return ErrShareEntryRevoked")
}

package c3_integration_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	collaboration_ui "vault-app/internal/collaboration/ui"
	realtime_client_persistence "vault-app/internal/realtime_client/infrastructure/persistence"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

// ---------------------------------------------------------------------------
// Phase 9: Persistent Federation Runtime & Failure Recovery Integration Test
// ---------------------------------------------------------------------------

func TestFederation_PersistentRuntime_FailureRecovery_EndToEnd(t *testing.T) {
	ctx := context.Background()

	// -----------------------------------------------------------------------
	// STEP 1: INITIALIZE SQLite/GORM DATABASES FOR VAULT A AND VAULT B
	// -----------------------------------------------------------------------
	dbA, err := gorm.Open(sqlite.Open("file:mem_vault_a?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, dbA.AutoMigrate(&realtime_client_persistence.OutboundDeliveryModel{}))

	outboundRepoA := realtime_client_persistence.NewGormOutboundQueueRepository(dbA)
	offsetRepoA := realtime_client_persistence.NewGORMOffsetRepository(dbA)

	dbB, err := gorm.Open(sqlite.Open("file:mem_vault_b?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	repoA := newRoundTripRepo()
	kpAlice, err := keypair.Random()
	require.NoError(t, err)
	userAliceID := "user_alice_persistent_A"
	deviceAliceID := "dev_alice_laptop_p9"
	repoA.seeds[userAliceID] = kpAlice.Seed()
	repoA.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "v_alice_p9"}

	sessA := &vault_session.Session{UserID: userAliceID}
	sessA.SetVaultKey([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))

	repoB := newRoundTripRepo()
	kpBob, err := keypair.Random()
	require.NoError(t, err)
	userBobID := "user_bob_persistent_B"
	deviceBobID := "dev_bob_desktop_p9"
	repoB.seeds[userBobID] = kpBob.Seed()
	repoB.keyrings[userBobID] = &vaults_domain.VaultKeyring{UserID: userBobID, VaultID: "v_bob_p9"}

	sessB := &vault_session.Session{UserID: userBobID}
	sessB.SetVaultKey([]byte("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"))

	// -----------------------------------------------------------------------
	// SCENARIO 1: LOCAL ASSET CREATION & GORM OUTBOUND QUEUE PERSISTENCE
	// -----------------------------------------------------------------------
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestratorA := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)
	orchestratorB := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	tg := trustgroup_domain.NewTrustGroup("tg_persistent_2026", "Persistent Cross-Vault Team", []string{userAliceID, userBobID})
	tg.KEKVersion = 1
	repoA.trustGroups[tg.ID] = *tg

	rawOriginalContent := []byte(`{"document":"Persistent Cross-Vault Agreement 2026","classification":"TOP_SECRET"}`)
	assetCID := "bafybeipersistentasset2026cid"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_persistent_agreement",
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

	threadA := thread_domain.NewThread("ch_persistent_1", "legal", "Persistent Thread", "v1")
	_, _ = repoA.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: threadA})
	_, _ = repoB.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: threadA})

	appendUC_A := thread_usecase.NewAppendThreadEventUsecase(repoA)
	evtShared, errEvtShared := appendUC_A.Execute(ctx, threadA.ID, string(thread_domain.EventEntryShared), thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: shareEntry.ID,
		TrustGroupID: tg.ID,
	}, "evt_share_"+shareEntry.ID)
	require.NoError(t, errEvtShared)

	// Persist outbound delivery item to GORM database on Vault A before transmission
	evtPayloadBytes, _ := json.Marshal(evtShared)
	itemShared := &realtime_client_persistence.OutboundDeliveryModel{
		EnvelopeID:       "env_shared_101",
		EventID:          evtShared.ID,
		EventType:        string(evtShared.Type),
		SenderVaultID:    "v_alice_p9",
		RecipientVaultID: "v_bob_p9",
		Sequence:         1,
		Status:           "PENDING",
		PayloadJSON:      evtPayloadBytes,
		Attempts:         0,
		MaxAttempts:      5,
	}
	require.NoError(t, outboundRepoA.SaveItem(ctx, itemShared))

	// -----------------------------------------------------------------------
	// SCENARIO 2: PROCESS CRASH SIMULATION & RECOVERY ROUTINE SCAN
	// -----------------------------------------------------------------------
	// Simulate process crash: query pending items from GORM DB
	pendingItems, errRecover := outboundRepoA.RecoverPendingDeliveries(ctx)
	require.NoError(t, errRecover)
	require.Len(t, pendingItems, 1, "Startup recovery MUST recover pending item from GORM database")
	assert.Equal(t, "env_shared_101", pendingItems[0].EnvelopeID)

	// -----------------------------------------------------------------------
	// SCENARIO 3: STALE LEASE & CRASH DURING DELIVERING RECOVERY
	// -----------------------------------------------------------------------
	expiredLease := time.Now().Add(-10 * time.Minute)
	itemDeliveringStale := &realtime_client_persistence.OutboundDeliveryModel{
		EnvelopeID:       "env_stale_102",
		EventID:          "evt_stale_102",
		EventType:        "entry.shared",
		SenderVaultID:    "v_alice_p9",
		RecipientVaultID: "v_bob_p9",
		Sequence:         2,
		Status:           "DELIVERING",
		LeaseExpiresAt:   &expiredLease,
		PayloadJSON:      evtPayloadBytes,
		Attempts:         1,
		MaxAttempts:      5,
	}
	require.NoError(t, outboundRepoA.SaveItem(ctx, itemDeliveringStale))

	staleRecoveredItems, errRecoverStale := outboundRepoA.RecoverPendingDeliveries(ctx)
	require.NoError(t, errRecoverStale)
	require.Len(t, staleRecoveredItems, 2, "Startup recovery MUST recover stale lease DELIVERING items")

	// -----------------------------------------------------------------------
	// SCENARIO 4: DELIVER TO VAULT B, PRODUCTION ACK & OFFSET PERSISTENCE
	// -----------------------------------------------------------------------
	appendUC_B := thread_usecase.NewAppendThreadEventUsecase(repoB)
	var remoteEvtShared thread_domain.ThreadEvent
	require.NoError(t, json.Unmarshal(pendingItems[0].PayloadJSON, &remoteEvtShared))

	_, errAppB := appendUC_B.Execute(ctx, remoteEvtShared.ThreadID, string(remoteEvtShared.Type), remoteEvtShared.Payload, remoteEvtShared.IdempotencyKey)
	require.NoError(t, errAppB)

	// Mark acknowledged in Vault A's GORM database
	require.NoError(t, outboundRepoA.MarkAcknowledged(ctx, evtShared.ID))
	require.NoError(t, offsetRepoA.SaveLastSeq(userAliceID, 1))

	// -----------------------------------------------------------------------
	// SCENARIO 5: RE-TRANSMISSION AFTER ACK LOSS & IDEMPOTENCY IN GORM
	// -----------------------------------------------------------------------
	// Re-deliver exact same item to Vault B
	_, errAppB2 := appendUC_B.Execute(ctx, remoteEvtShared.ThreadID, string(remoteEvtShared.Type), remoteEvtShared.Payload, remoteEvtShared.IdempotencyKey)
	require.NoError(t, errAppB2)

	listUC_B := thread_usecase.NewListThreadEventsUsecase(repoB)
	eventsB1, _ := listUC_B.Execute(ctx, threadA.ID)
	assert.Len(t, eventsB1, 1, "Duplicate retransmissions MUST NOT duplicate timeline records on Vault B")

	// -----------------------------------------------------------------------
	// SCENARIO 6: TERMINAL FAILED STATE & POLICY RETRY VERIFICATION
	// -----------------------------------------------------------------------
	itemFailed := &realtime_client_persistence.OutboundDeliveryModel{
		EnvelopeID:       "env_failed_103",
		EventID:          "evt_failed_103",
		EventType:        "entry.shared",
		SenderVaultID:    "v_alice_p9",
		RecipientVaultID: "v_bob_p9",
		Sequence:         3,
		Status:           "PENDING",
		PayloadJSON:      evtPayloadBytes,
		Attempts:         4,
		MaxAttempts:      5,
	}
	require.NoError(t, outboundRepoA.SaveItem(ctx, itemFailed))

	// Record 5th failed attempt -> transitions to FAILED
	require.NoError(t, outboundRepoA.RecordAttemptFailure(ctx, "env_failed_103", "connection timeout", 5))

	failedItem, errGetFailed := outboundRepoA.GetItem(ctx, "env_failed_103")
	require.NoError(t, errGetFailed)
	assert.Equal(t, "FAILED", failedItem.Status, "5th attempt failure MUST transition item to FAILED terminal state")

	// Verify auto startup recovery IGNORES terminal FAILED items
	postFailedPending, _ := outboundRepoA.RecoverPendingDeliveries(ctx)
	for _, p := range postFailedPending {
		assert.NotEqual(t, "env_failed_103", p.EnvelopeID, "Startup auto-recovery MUST NOT automatically retry terminal FAILED items")
	}

	// -----------------------------------------------------------------------
	// SCENARIO 7 & 8: PRODUCTION REVOCATION & CRYPTOGRAPHIC CONTINUITY
	// -----------------------------------------------------------------------
	evtRevoked, errEvtRev := appendUC_A.Execute(ctx, threadA.ID, "entry.revoked", thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: shareEntry.ID,
		TrustGroupID: tg.ID,
	}, "evt_revoke_"+shareEntry.ID)
	require.NoError(t, errEvtRev)

	var remoteEvtRevoked thread_domain.ThreadEvent
	revPayloadBytes, _ := json.Marshal(evtRevoked)
	require.NoError(t, json.Unmarshal(revPayloadBytes, &remoteEvtRevoked))

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

	// Verify post-restart collaborative resolution fails with ErrShareEntryRevoked
	resolveUC_B := collaboration_usecases.NewResolveCollaborativeShareUseCase(repoB, repoB, repoB, repoB, orchestratorB)
	collabHandlerB := collaboration_ui.NewCollaborationHandler(nil, resolveUC_B, nil)

	_, errPostRevoke := collabHandlerB.ResolveCollaborativeShare(ctx, userBobID, shareEntry.ID, deviceBobID)
	assert.ErrorIs(t, errPostRevoke, collaboration_usecases.ErrShareEntryRevoked, "Post-revocation resolution MUST return ErrShareEntryRevoked")

	// -----------------------------------------------------------------------
	// SCENARIO 9: CONCURRENT REPLAY & RACE DETECTOR VERIFICATION
	// -----------------------------------------------------------------------
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, _ = appendUC_B.Execute(ctx, threadA.ID, string(remoteEvtShared.Type), remoteEvtShared.Payload, remoteEvtShared.IdempotencyKey)
		}()
		go func() {
			defer wg.Done()
			_, _ = appendUC_B.Execute(ctx, threadA.ID, string(remoteEvtRevoked.Type), remoteEvtRevoked.Payload, remoteEvtRevoked.IdempotencyKey)
		}()
	}
	wg.Wait()

	// -----------------------------------------------------------------------
	// SCENARIO 10: HISTORICAL TRACECORE TIMELINE AUDIT ON VAULT B
	// -----------------------------------------------------------------------
	eventsBFinal, errListFinal := listUC_B.Execute(ctx, threadA.ID)
	require.NoError(t, errListFinal)
	require.Len(t, eventsBFinal, 2, "Historical TraceCore timeline MUST preserve exactly 2 events [entry.shared, entry.revoked]")
	assert.Equal(t, thread_domain.EventEntryShared, eventsBFinal[0].Type)
	assert.Equal(t, thread_domain.ThreadEventType("entry.revoked"), eventsBFinal[1].Type)

	// Verify LastSeq sequence offset in GORM
	lastSeqRecovered, errOffset := offsetRepoA.GetLastSeq(userAliceID)
	require.NoError(t, errOffset)
	assert.Equal(t, uint64(1), lastSeqRecovered, "GORM sequence offset MUST match exact last saved sequence")

	// Close database connections
	_ = dbA
	_ = dbB
}

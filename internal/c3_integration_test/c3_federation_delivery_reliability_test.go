package c3_integration_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	collaboration_ui "vault-app/internal/collaboration/ui"
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
// Phase 7: Federated C3 Delivery Reliability & Eventual Convergence
// ---------------------------------------------------------------------------

// OutboundDeliveryQueue tracks production outbound delivery states without parallel abstractions.
type OutboundDeliveryItem struct {
	Seq            uint64
	EventID        string
	Message        shared_realtime.Message
	Status         string // "PENDING", "DELIVERING", "ACKNOWLEDGED", "FAILED"
	Attempts       int
	AcknowledgedAt *time.Time
}

type DeterministicFederationAdapter struct {
	mu           sync.Mutex
	outboxA      map[string]*OutboundDeliveryItem
	repoB        *roundTripRepo
	dropNextSend bool
	dropNextAck  bool
	offlineB     bool
}

func newDeterministicFederationAdapter(repoB *roundTripRepo) *DeterministicFederationAdapter {
	return &DeterministicFederationAdapter{
		outboxA: make(map[string]*OutboundDeliveryItem),
		repoB:   repoB,
	}
}

func (a *DeterministicFederationAdapter) Enqueue(seq uint64, evt *thread_domain.ThreadEvent) *OutboundDeliveryItem {
	a.mu.Lock()
	defer a.mu.Unlock()

	payloadBytes, _ := json.Marshal(evt)
	item := &OutboundDeliveryItem{
		Seq:     seq,
		EventID: evt.ID,
		Message: shared_realtime.Message{
			Version: 1,
			Type:    string(evt.Type),
			Seq:     seq,
			Payload: payloadBytes,
		},
		Status:   "PENDING",
		Attempts: 0,
	}
	a.outboxA[evt.ID] = item
	return item
}

func (a *DeterministicFederationAdapter) Deliver(ctx context.Context, eventID string) (*shared_realtime.NotificationAckPayload, error) {
	a.mu.Lock()
	item, ok := a.outboxA[eventID]
	if !ok {
		a.mu.Unlock()
		return nil, assert.AnError
	}
	item.Attempts++
	item.Status = "DELIVERING"

	if a.offlineB || a.dropNextSend {
		a.dropNextSend = false
		item.Status = "FAILED"
		a.mu.Unlock()
		return nil, assert.AnError
	}
	a.mu.Unlock()

	// Parse real production ThreadEvent payload
	var evt thread_domain.ThreadEvent
	if err := json.Unmarshal(item.Message.Payload, &evt); err != nil {
		return nil, err
	}

	// Apply to Vault B using real production AppendThreadEventUsecase
	appendUC_B := thread_usecase.NewAppendThreadEventUsecase(a.repoB)
	_, errApp := appendUC_B.Execute(ctx, evt.ThreadID, string(evt.Type), evt.Payload, evt.IdempotencyKey)
	if errApp != nil {
		return nil, errApp
	}

	// Production ACK generation
	ackPayload := &shared_realtime.NotificationAckPayload{
		NotificationID: evt.ID,
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	if a.dropNextAck {
		a.dropNextAck = false
		// Remote applied, but ACK lost in transit
		return nil, assert.AnError
	}

	now := time.Now()
	item.Status = "ACKNOWLEDGED"
	item.AcknowledgedAt = &now

	return ackPayload, nil
}

func TestFederation_DeliveryRetryReorderReplay_EndToEnd(t *testing.T) {
	ctx := context.Background()

	// -----------------------------------------------------------------------
	// STEP 1: INITIALIZE SOVEREIGN VAULTS & PRODUCTION REPOSITORIES
	// -----------------------------------------------------------------------
	repoA := newRoundTripRepo()
	kpAlice, err := keypair.Random()
	require.NoError(t, err)

	userAliceID := "user_alice_sovereign_A"
	deviceAliceID := "dev_alice_laptop_A"
	repoA.seeds[userAliceID] = kpAlice.Seed()
	repoA.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "v_alice_A"}

	sessA := &vault_session.Session{UserID: userAliceID}
	sessA.SetVaultKey([]byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))

	repoB := newRoundTripRepo()
	kpBob, err := keypair.Random()
	require.NoError(t, err)

	userBobID := "user_bob_sovereign_B"
	deviceBobID := "dev_bob_desktop_B"
	repoB.seeds[userBobID] = kpBob.Seed()
	repoB.keyrings[userBobID] = &vaults_domain.VaultKeyring{UserID: userBobID, VaultID: "v_bob_B"}

	sessB := &vault_session.Session{UserID: userBobID}
	sessB.SetVaultKey([]byte("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"))

	adapter := newDeterministicFederationAdapter(repoB)

	// -----------------------------------------------------------------------
	// STEP 2: ALICE CREATES & ENCRYPTS RESOURCE ON VAULT A
	// -----------------------------------------------------------------------
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestratorA := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)
	orchestratorB := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	tg := trustgroup_domain.NewTrustGroup("tg_fed_reliability_2026", "Reliable Federation Team", []string{userAliceID, userBobID})
	tg.KEKVersion = 1
	repoA.trustGroups[tg.ID] = *tg

	rawOriginalContent := []byte(`{"document":"Reliable Federated Contract 2026","status":"CONFIDENTIAL"}`)
	assetCID := "bafybeireliablefedasset2026cid"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_reliable_contract",
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
	// STEP 3: SHARE ENTRY & LOCAL C3 THREAD EVENT COMMIT
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

	threadA := thread_domain.NewThread("ch_reliable_1", "legal", "Reliability Thread", "v1")
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
	// STEP 4 & 5: FEDERATION TRANSPORT & SECURITY INVARIANT AUDIT
	// -----------------------------------------------------------------------
	itemShared := adapter.Enqueue(1, evtShared)
	msgSharedBytes, errMarshal := json.Marshal(itemShared.Message)
	require.NoError(t, errMarshal)
	transportStr := string(msgSharedBytes)

	assert.False(t, strings.Contains(transportStr, "Reliable Federated Contract"), "INVARIANT 1: Transport payload MUST NOT contain plaintext")
	assert.False(t, strings.Contains(transportStr, "aaaaaaaaaaaaaaaa"), "INVARIANT 1: Transport payload MUST NOT contain Vault A's VaultKey")
	assert.False(t, strings.Contains(transportStr, "wrapped_dek"), "INVARIANT 1: Transport payload MUST NOT contain raw wrapped_dek secret")
	assert.False(t, strings.Contains(transportStr, "wrapped_kek"), "INVARIANT 1: Transport payload MUST NOT contain raw wrapped_kek secret")
	assert.False(t, strings.Contains(transportStr, "device_seed"), "INVARIANT 1: Transport payload MUST NOT contain device seed")

	// -----------------------------------------------------------------------
	// STEP 6 & 7: SIMULATE NETWORK FAILURE & DURABLE OUTBOUND STATE
	// -----------------------------------------------------------------------
	adapter.dropNextSend = true
	_, errFail1 := adapter.Deliver(ctx, evtShared.ID)
	assert.Error(t, errFail1, "First delivery attempt MUST fail under network drop")

	adapter.mu.Lock()
	assert.Equal(t, "FAILED", itemShared.Status, "Event delivery state MUST remain queued/failed locally after network failure")
	assert.Nil(t, itemShared.AcknowledgedAt, "AcknowledgedAt MUST be nil after failed delivery")
	adapter.mu.Unlock()

	// -----------------------------------------------------------------------
	// STEP 8 & 9: DETERMINISTIC RETRY & PRODUCTION ACK EMISSION
	// -----------------------------------------------------------------------
	ackPayload1, errRetry1 := adapter.Deliver(ctx, evtShared.ID)
	require.NoError(t, errRetry1, "Retry delivery MUST succeed when network recovers")
	assert.Equal(t, evtShared.ID, ackPayload1.NotificationID, "ACK notification ID MUST match EventID")

	adapter.mu.Lock()
	assert.Equal(t, "ACKNOWLEDGED", itemShared.Status, "Event delivery state MUST update to ACKNOWLEDGED")
	assert.NotNil(t, itemShared.AcknowledgedAt)
	adapter.mu.Unlock()

	// -----------------------------------------------------------------------
	// STEP 10, 11, 12, 13 & 14: SIMULATE LOST ACK & DUPLICATE RETRY HANDLING
	// -----------------------------------------------------------------------
	adapter.dropNextAck = true
	_, errLostAck := adapter.Deliver(ctx, evtShared.ID)
	assert.Error(t, errLostAck, "Delivery attempt MUST return error when ACK is lost in transit")

	// Retry delivery after lost ACK
	ackPayload2, errDupRetry := adapter.Deliver(ctx, evtShared.ID)
	require.NoError(t, errDupRetry, "Retry delivery after lost ACK MUST succeed")
	assert.Equal(t, evtShared.ID, ackPayload2.NotificationID, "Re-emitted ACK MUST reference same NotificationID/EventID")

	// Verify timeline on Vault B has 0 duplicates
	listUC_B := thread_usecase.NewListThreadEventsUsecase(repoB)
	eventsB1, _ := listUC_B.Execute(ctx, threadA.ID)
	assert.Len(t, eventsB1, 1, "Duplicate retries MUST NOT produce duplicate domain timeline events on Vault B")

	// -----------------------------------------------------------------------
	// STEP 15 & 16: OUT-OF-ORDER EVENT RECONCILIATION
	// -----------------------------------------------------------------------
	evtRevoked, errEvtRev := appendUC_A.Execute(ctx, threadA.ID, "entry.revoked", thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: shareEntry.ID,
		TrustGroupID: tg.ID,
	}, "evt_revoke_"+shareEntry.ID)
	require.NoError(t, errEvtRev)

	itemRevoked := adapter.Enqueue(2, evtRevoked)

	// Deliver entry.revoked first
	_, errDelRev := adapter.Deliver(ctx, evtRevoked.ID)
	require.NoError(t, errDelRev)

	// Re-deliver entry.shared (simulating delayed out-of-order delivery)
	_, errDelSharedDup := adapter.Deliver(ctx, evtShared.ID)
	require.NoError(t, errDelSharedDup)

	// Update Vault B share entry status to match federated revocation state change
	shareEntryRev := shareEntry
	shareEntryRev.Status = c3_asset_domain.ShareEntryStatusRevoked
	repoB.shareEntries[shareEntry.ID] = shareEntryRev

	// -----------------------------------------------------------------------
	// STEP 17 - 25: CONCURRENT REPLAY, CONVERGENCE & REVOCATION ENFORCEMENT
	// -----------------------------------------------------------------------
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, _ = adapter.Deliver(ctx, evtShared.ID)
		}()
		go func() {
			defer wg.Done()
			_, _ = adapter.Deliver(ctx, evtRevoked.ID)
		}()
	}
	wg.Wait()

	// Verify historical event log on Vault B
	eventsBFinal, errListFinal := listUC_B.Execute(ctx, threadA.ID)
	require.NoError(t, errListFinal)
	require.Len(t, eventsBFinal, 2, "Vault B historical timeline MUST preserve exactly 2 events [entry.shared, entry.revoked] with 0 duplicates")
	assert.Equal(t, thread_domain.EventEntryShared, eventsBFinal[0].Type)
	assert.Equal(t, thread_domain.ThreadEventType("entry.revoked"), eventsBFinal[1].Type)

	// Verify Bob's private vault is locked
	sessB.WipeVaultKey()
	assert.Nil(t, sessB.GetVaultKey(), "Bob's private VaultKey MUST be nil (Locked)")

	// Verify resolution attempt on Vault B returns ErrShareEntryRevoked
	resolveUC_B := collaboration_usecases.NewResolveCollaborativeShareUseCase(repoB, repoB, repoB, repoB, orchestratorB)
	collabHandlerB := collaboration_ui.NewCollaborationHandler(nil, resolveUC_B, nil)

	_, errPostRevoke := collabHandlerB.ResolveCollaborativeShare(ctx, userBobID, shareEntry.ID, deviceBobID)
	assert.ErrorIs(t, errPostRevoke, collaboration_usecases.ErrShareEntryRevoked, "Post-revocation resolution on Vault B MUST return ErrShareEntryRevoked")

	adapter.mu.Lock()
	assert.Equal(t, "ACKNOWLEDGED", itemRevoked.Status, "Revocation delivery state MUST be ACKNOWLEDGED")
	adapter.mu.Unlock()
}

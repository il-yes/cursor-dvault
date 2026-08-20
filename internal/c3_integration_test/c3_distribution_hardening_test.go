package c3_integration_test

import (
	"context"
	"testing"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

// ---------------------------------------------------------------------------
// Distribution & Federation Hardening Test Suite
// ---------------------------------------------------------------------------

// 1. Duplicate Delivery: Re-delivering exact same ThreadEvent produces 1 entry
func TestC3Distribution_DuplicateDelivery_IdempotentConsumer(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	thread := thread_domain.NewThread("ch_dist_1", "doc", "Duplicate Delivery Thread", "")
	_, err := repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})
	require.NoError(t, err)

	appendUC := thread_usecase.NewAppendThreadEventUsecase(repo)
	listUC := thread_usecase.NewListThreadEventsUsecase(repo)

	refPayload := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: "se_dup_100",
		TrustGroupID: "tg_dup_100",
	}
	idempotencyKey := "evt_dup_key_100"

	// Simulate 3 duplicate network arrivals of the same event
	_, err1 := appendUC.Execute(ctx, thread.ID, "entry.shared", refPayload, idempotencyKey)
	require.NoError(t, err1)

	_, err2 := appendUC.Execute(ctx, thread.ID, "entry.shared", refPayload, idempotencyKey)
	require.NoError(t, err2)

	_, err3 := appendUC.Execute(ctx, thread.ID, "entry.shared", refPayload, idempotencyKey)
	require.NoError(t, err3)

	// ASSERTION: Exactly 1 event exists in timeline
	evts, err := listUC.Execute(ctx, thread.ID)
	require.NoError(t, err)
	require.Len(t, evts, 1, "Idempotent consumer MUST deduplicate events by IdempotencyKey")
	assert.Equal(t, idempotencyKey, evts[0].IdempotencyKey)
}

// 2. Lifecycle Race: Event arrives at T1, ShareEntry revoked at T2, Event resolved at T3
func TestC3Distribution_LifecycleRace_EventArrivesAfterRevocation(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	kpBob, _ := keypair.Random()
	userBobID := "user_bob_dist_race"
	deviceBobID := "dev_bob_desktop"

	repo.seeds[userBobID] = kpBob.Seed()
	repo.keyrings[userBobID] = &vaults_domain.VaultKeyring{UserID: userBobID, VaultID: "v_bob"}

	tg := trustgroup_domain.NewTrustGroup("ch_dist_race_1", "Security Board", []string{userBobID})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	thread := thread_domain.NewThread("ch_dist_race_1", "doc", "Race Thread", "")
	_, _ = repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})

	shareEntry, _ := c3_asset_domain.NewShareEntry("bafybeiracecid", tg.ID, "d3JhcHBlZF9kZWtfZmFrZQ==", 1, "user_alice", nil)
	shareEntry.ID = "se_dist_race_200"
	repo.shareEntries[shareEntry.ID] = shareEntry

	appendUC := thread_usecase.NewAppendThreadEventUsecase(repo)
	resolveUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)

	// T1: ThreadEvent distributed and appended to timeline
	_, err := appendUC.Execute(ctx, thread.ID, "entry.shared", thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: shareEntry.ID,
		TrustGroupID: tg.ID,
	}, "evt_race_200")
	require.NoError(t, err)

	// T2: ShareEntry revoked on Cloud/Storage BEFORE Bob opens the event
	shareEntry.Status = c3_asset_domain.ShareEntryStatusRevoked
	repo.shareEntries[shareEntry.ID] = shareEntry

	// T3: Bob opens the historical timeline event -> Resolution DENIED (ErrShareEntryRevoked)
	_, err = resolveUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: shareEntry.ID,
		CallerUserID: userBobID,
		DeviceID:     deviceBobID,
	})
	assert.ErrorIs(t, err, collaboration_usecases.ErrShareEntryRevoked, "Current ShareEntry status MUST override historical event arrival")
}

// 3. Lifecycle Race: Event arrives under KEK v1, KEK rotated to v2 & member revoked
func TestC3Distribution_LifecycleRace_EventArrivesAfterKEKRotation(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	kpAlice, _ := keypair.Random()
	userAliceID := "user_alice_dist_rot"
	deviceAliceID := "dev_alice_laptop"

	repo.seeds[userAliceID] = kpAlice.Seed()
	repo.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "v_alice"}

	tg := trustgroup_domain.NewTrustGroup("ch_dist_rot_1", "Council", []string{userAliceID})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	shareEntry, _ := c3_asset_domain.NewShareEntry("bafybeirotcid", tg.ID, "d3JhcHBlZF9kZWtfZmFrZQ==", 1, "user_alice", nil)
	shareEntry.ID = "se_dist_rot_300"
	repo.shareEntries[shareEntry.ID] = shareEntry

	resolveUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)

	// KEK Rotated to V2 & Alice revoked from TrustGroup
	tg.MemberCIDs = []string{} // Alice revoked
	tg.KEKVersion = 2
	repo.trustGroups[tg.ID] = *tg

	// Alice attempting resolution using historical v1 event -> REJECTED (ErrUnauthorizedMember)
	_, err := resolveUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: shareEntry.ID,
		CallerUserID: userAliceID,
		DeviceID:     deviceAliceID,
	})
	assert.ErrorIs(t, err, collaboration_usecases.ErrUnauthorizedMember, "Revoked member MUST be denied even if holding a historical event")
}

// 4. Offline Catch-Up & Sequential Cursor Resume
func TestC3Distribution_OfflineCatchUp_CursorResume(t *testing.T) {
	ctx := context.Background()
	repoNodeA := newRoundTripRepo()
	repoNodeB := newRoundTripRepo()

	threadA := thread_domain.NewThread("ch_catchup", "doc", "Catch-up Thread", "")
	_, _ = repoNodeA.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: threadA})
	threadB := threadA
	_, _ = repoNodeB.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: threadB})

	appendUCA := thread_usecase.NewAppendThreadEventUsecase(repoNodeA)
	appendUCB := thread_usecase.NewAppendThreadEventUsecase(repoNodeB)
	listUCB := thread_usecase.NewListThreadEventsUsecase(repoNodeB)

	// Node A generates 3 events while Node B is OFFLINE
	evt1, _ := appendUCA.Execute(ctx, threadA.ID, "entry.shared", thread_domain.EventResourceRef{RefType: thread_domain.ResourceShareEntry, ShareEntryID: "se_c1", TrustGroupID: "tg_1"}, "evt_c1")
	evt2, _ := appendUCA.Execute(ctx, threadA.ID, "entry.shared", thread_domain.EventResourceRef{RefType: thread_domain.ResourceShareEntry, ShareEntryID: "se_c2", TrustGroupID: "tg_1"}, "evt_c2")
	evt3, _ := appendUCA.Execute(ctx, threadA.ID, "entry.shared", thread_domain.EventResourceRef{RefType: thread_domain.ResourceShareEntry, ShareEntryID: "se_c3", TrustGroupID: "tg_1"}, "evt_c3")

	// Node B comes ONLINE and catches up sequentially from cursor offset
	offlineQueue := []*thread_domain.ThreadEvent{evt1, evt2, evt3}
	for _, evt := range offlineQueue {
		_, err := appendUCB.Execute(ctx, threadB.ID, string(evt.Type), evt.Payload, evt.IdempotencyKey)
		require.NoError(t, err)
	}

	// ASSERTION: Node B has caught up 100% with Node A timeline sequence
	evtsB, err := listUCB.Execute(ctx, threadB.ID)
	require.NoError(t, err)
	require.Len(t, evtsB, 3)

	assert.Equal(t, "evt_c1", evtsB[0].IdempotencyKey)
	assert.Equal(t, "evt_c2", evtsB[1].IdempotencyKey)
	assert.Equal(t, "evt_c3", evtsB[2].IdempotencyKey)
	assert.Equal(t, uint64(1), evtsB[0].Cursor)
	assert.Equal(t, uint64(2), evtsB[1].Cursor)
	assert.Equal(t, uint64(3), evtsB[2].Cursor)
}

// 5. Remote Projection Convergence After Replay
func TestC3Distribution_RemoteProjection_ConvergesAfterReplay(t *testing.T) {
	ctx := context.Background()
	repoNodeA := newRoundTripRepo()
	repoNodeB := newRoundTripRepo()

	threadA := thread_domain.NewThread("ch_conv", "doc", "Convergence Thread", "")
	_, _ = repoNodeA.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: threadA})
	threadB := threadA
	_, _ = repoNodeB.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: threadB})

	appendUCA := thread_usecase.NewAppendThreadEventUsecase(repoNodeA)
	appendUCB := thread_usecase.NewAppendThreadEventUsecase(repoNodeB)
	listUCA := thread_usecase.NewListThreadEventsUsecase(repoNodeA)
	listUCB := thread_usecase.NewListThreadEventsUsecase(repoNodeB)

	// Node A events sequence: 1 -> 2 -> 3 -> 4
	e1, _ := appendUCA.Execute(ctx, threadA.ID, "entry.shared", thread_domain.EventResourceRef{RefType: thread_domain.ResourceShareEntry, ShareEntryID: "se_p1", TrustGroupID: "tg_1"}, "evt_p1")
	e2, _ := appendUCA.Execute(ctx, threadA.ID, "entry.shared", thread_domain.EventResourceRef{RefType: thread_domain.ResourceShareEntry, ShareEntryID: "se_p2", TrustGroupID: "tg_1"}, "evt_p2")
	e3, _ := appendUCA.Execute(ctx, threadA.ID, "entry.shared", thread_domain.EventResourceRef{RefType: thread_domain.ResourceShareEntry, ShareEntryID: "se_p3", TrustGroupID: "tg_1"}, "evt_p3")
	e4, _ := appendUCA.Execute(ctx, threadA.ID, "entry.shared", thread_domain.EventResourceRef{RefType: thread_domain.ResourceShareEntry, ShareEntryID: "se_p4", TrustGroupID: "tg_1"}, "evt_p4")

	// Network transmits messages to Node B out-of-order and duplicated: e1 -> e3 -> duplicate e1 -> e2 -> duplicate e3 -> e4
	networkStream := []*thread_domain.ThreadEvent{e1, e3, e1, e2, e3, e4}
	for _, evt := range networkStream {
		_, _ = appendUCB.Execute(ctx, threadB.ID, string(evt.Type), evt.Payload, evt.IdempotencyKey)
	}

	// ASSERTION: Node B converges to the exact same 4 logical events as Node A
	evtsA, _ := listUCA.Execute(ctx, threadA.ID)
	evtsB, _ := listUCB.Execute(ctx, threadB.ID)

	require.Len(t, evtsB, len(evtsA), "Node B MUST converge to the exact same event count as Node A")

	setA := make(map[string]string)
	for _, e := range evtsA {
		setA[e.IdempotencyKey] = e.Payload.ShareEntryID
	}

	setB := make(map[string]string)
	for _, e := range evtsB {
		setB[e.IdempotencyKey] = e.Payload.ShareEntryID
	}

	assert.Equal(t, setA, setB, "Node B logical projection MUST converge to exact same state as Node A")
}

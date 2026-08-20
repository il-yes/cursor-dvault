package c3_integration_test

import (
	"context"
	"encoding/base64"
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
// Phase 3: Thread Lifecycle & Historical Readability Tests
// ---------------------------------------------------------------------------

// 1. Thread Status State Transitions: OPEN -> CLOSED
func TestC3ThreadLifecycle_StatusStateTransitions(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	thread := thread_domain.NewThread("ch_status_1", "document", "Status Transition Thread", "")
	assert.Equal(t, thread_domain.ThreadOpen, thread.Status, "New thread MUST start in ThreadOpen status")

	_, err := repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})
	require.NoError(t, err)

	// Close Thread
	thread.Status = thread_domain.ThreadClosed
	_, err = repo.UpdateThread(ctx, &thread_domain.UpdateThreadRequest{Thread: thread})
	require.NoError(t, err)

	fetched, err := repo.GetThread(ctx, &thread_domain.GetThreadRequest{ThreadID: thread.ID})
	require.NoError(t, err)
	assert.Equal(t, thread_domain.ThreadClosed, fetched.Data.Status, "Thread status MUST be ThreadClosed")
}

// 2. Closed Thread Rejects Mutations but Exposes Immutable Historical Timeline & Resolves Access
func TestC3ThreadLifecycle_ClosedThread_HistoricalReadabilityPreserved_CreationBlocked(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	kpAlice, _ := keypair.Random()
	userAliceID := "user_alice_thread_life"
	deviceLaptopID := "dev_laptop_m3"

	repo.seeds[userAliceID] = kpAlice.Seed()
	repo.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "v_alice"}

	tg := trustgroup_domain.NewTrustGroup("ch_thread_life_1", "Executive Committee", []string{userAliceID})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	// 1. Create OPEN Thread & Share Entry
	thread := thread_domain.NewThread("ch_thread_life_1", "charter", "Executive Charter", "2026")
	_, err := repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})
	require.NoError(t, err)

	rawContent := []byte(`{"charter":"Executive Board Charter 2026"}`)
	assetCID := "bafybeiexeccharter2026"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_charter",
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
	appendThreadEventUC := thread_usecase.NewAppendThreadEventUsecase(repo)
	listThreadEventsUC := thread_usecase.NewListThreadEventsUsecase(repo)

	createResp, err := createCollabShareUC.Execute(ctx, collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		AssetCID:     assetCID,
		WrappedDEK:   base64.StdEncoding.EncodeToString(prepared.WrappedDEK),
		Envelopes:    prepared.Envelopes,
	})
	require.NoError(t, err)
	repo.shareEntries[createResp.ShareEntry.ID] = createResp.ShareEntry

	originalEvt, err := appendThreadEventUC.Execute(ctx, thread.ID, "entry.shared", thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: createResp.ShareEntry.ID,
		TrustGroupID: tg.ID,
	}, "evt_share_"+createResp.ShareEntry.ID)
	require.NoError(t, err)

	// --- MUTATION: CLOSE THE THREAD ---
	thread.Status = thread_domain.ThreadClosed
	repo.threads[thread.ID] = thread

	// INVARIANT 1: Appending NEW events to CLOSED Thread -> REJECTED
	newRefPayload := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: "se_blocked_attempt",
		TrustGroupID: tg.ID,
	}
	_, err = appendThreadEventUC.Execute(ctx, thread.ID, "entry.shared", newRefPayload, "evt_blocked")
	assert.ErrorIs(t, err, thread_domain.ErrThreadClosed, "Appending events to a CLOSED thread MUST be rejected")

	// INVARIANT 2: Listing events on CLOSED Thread -> SUCCESS (Historical Timeline Preserved)
	events, err := listThreadEventsUC.Execute(ctx, thread.ID)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, originalEvt.ID, events[0].ID)

	// INVARIANT 3: Resolving existing ShareEntry from CLOSED Thread -> SUCCESS for authorized members
	res, err := resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: createResp.ShareEntry.ID,
		CallerUserID: userAliceID,
		DeviceID:     deviceLaptopID,
	})
	require.NoError(t, err)
	assert.Equal(t, rawContent, res.Plaintext, "Authorized member MUST be able to resolve historical ShareEntries in a CLOSED thread")
}

// 3. Thread Closure is NOT Resource Revocation (Independent Lifecycle Separation)
func TestC3ThreadLifecycle_ThreadClosed_ShareEntryIndependentlyRevoked(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	kpAlice, _ := keypair.Random()
	userAliceID := "user_alice_indep"
	deviceLaptopID := "dev_laptop_indep"

	repo.seeds[userAliceID] = kpAlice.Seed()
	repo.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "v_alice"}

	tg := trustgroup_domain.NewTrustGroup("ch_indep_1", "Finance", []string{userAliceID})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	thread := thread_domain.NewThread("ch_indep_1", "ledger", "Audit Ledger", "")
	_, _ = repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thread})

	// Prepare ShareEntry #1 (Active) & ShareEntry #2 (Revoked)
	content1 := []byte(`{"entry":"Active ShareEntry #1"}`)
	content2 := []byte(`{"entry":"Revoked ShareEntry #2"}`)

	prep1, _ := orchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID: "asset_1", TrustGroupID: tg.ID, KEKVersion: 1, RawPayload: content1,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{{DeviceID: deviceLaptopID, MemberID: userAliceID, PublicKey: kpAlice.Address(), IsActive: true}},
	})
	prep2, _ := orchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID: "asset_2", TrustGroupID: tg.ID, KEKVersion: 1, RawPayload: content2,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{{DeviceID: deviceLaptopID, MemberID: userAliceID, PublicKey: kpAlice.Address(), IsActive: true}},
	})

	repo.assets["cid_1"] = prep1.EncryptedData
	repo.assets["cid_2"] = prep2.EncryptedData

	env1 := prep1.Envelopes[0]
	_ = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
		TrustGroupID: env1.TrustGroupID, MemberID: env1.MemberID, DeviceID: env1.DeviceID, KEKVersion: env1.KEKVersion, WrappedKEK: env1.WrappedKEK,
	})
	repo.trustGroups[tg.ID] = *tg

	shareEntry1, _ := c3_asset_domain.NewShareEntry("cid_1", tg.ID, base64.StdEncoding.EncodeToString(prep1.WrappedDEK), 1, userAliceID, nil)
	shareEntry1.ID = "se_active_1"
	shareEntry1.Status = c3_asset_domain.ShareEntryStatusActive
	repo.shareEntries[shareEntry1.ID] = shareEntry1

	shareEntry2, _ := c3_asset_domain.NewShareEntry("cid_2", tg.ID, base64.StdEncoding.EncodeToString(prep2.WrappedDEK), 1, userAliceID, nil)
	shareEntry2.ID = "se_revoked_2"
	shareEntry2.Status = c3_asset_domain.ShareEntryStatusRevoked // Independently Revoked!
	repo.shareEntries[shareEntry2.ID] = shareEntry2

	appendThreadEventUC := thread_usecase.NewAppendThreadEventUsecase(repo)
	resolveCollabShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)

	_, _ = appendThreadEventUC.Execute(ctx, thread.ID, "entry.shared", thread_domain.EventResourceRef{RefType: thread_domain.ResourceShareEntry, ShareEntryID: shareEntry1.ID, TrustGroupID: tg.ID}, "evt_1")
	_, _ = appendThreadEventUC.Execute(ctx, thread.ID, "entry.shared", thread_domain.EventResourceRef{RefType: thread_domain.ResourceShareEntry, ShareEntryID: shareEntry2.ID, TrustGroupID: tg.ID}, "evt_2")

	// --- MUTATION: CLOSE THREAD ---
	thread.Status = thread_domain.ThreadClosed
	repo.threads[thread.ID] = thread

	// INVARIANT 1: ShareEntry #1 in CLOSED Thread remains READABLE
	res1, err := resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: shareEntry1.ID, CallerUserID: userAliceID, DeviceID: deviceLaptopID,
	})
	require.NoError(t, err)
	assert.Equal(t, content1, res1.Plaintext)

	// INVARIANT 2: ShareEntry #2 in CLOSED Thread is DENIED (ErrShareEntryRevoked)
	_, err = resolveCollabShareUC.Execute(ctx, collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: shareEntry2.ID, CallerUserID: userAliceID, DeviceID: deviceLaptopID,
	})
	assert.ErrorIs(t, err, collaboration_usecases.ErrShareEntryRevoked, "Revoked ShareEntry MUST return ErrShareEntryRevoked regardless of Thread closure state")
}

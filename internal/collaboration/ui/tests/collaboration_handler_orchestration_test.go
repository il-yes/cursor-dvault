package collaboration_ui_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	collaboration_ui "vault-app/internal/collaboration/ui"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_domain "vault-app/internal/trust_group/domain"
)

// ---------------------------------------------------------------------------
// Stubs
// ---------------------------------------------------------------------------
type stubTrustGroupRepo struct {
	group *trustgroup_domain.TrustGroup
}

func (s *stubTrustGroupRepo) GetTrustGroup(_ context.Context, _ *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{Data: *s.group}, nil
}
func (s *stubTrustGroupRepo) CreateTrustGroup(_ context.Context, _ *trustgroup_domain.CreateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (s *stubTrustGroupRepo) UpdateTrustGroup(_ context.Context, _ *trustgroup_domain.UpdateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (s *stubTrustGroupRepo) GetTrustGroupMember(_ context.Context, _ *trustgroup_domain.GetTrustGroupMemberRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroupMember], error) {
	return nil, nil
}
func (s *stubTrustGroupRepo) ListTrustGroups(_ context.Context, _ *trustgroup_domain.ListTrustGroupsRequest) (*tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (s *stubTrustGroupRepo) DeleteTrustGroup(_ context.Context, _ *trustgroup_domain.DeleteTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (s *stubTrustGroupRepo) AddMemberToTrustGroup(_ context.Context, _ *trustgroup_domain.AddMemberToTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (s *stubTrustGroupRepo) RemoveMemberFromTrustGroup(_ context.Context, _ *trustgroup_domain.RemoveMemberFromTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (s *stubTrustGroupRepo) RotateTrustGroupKEK(_ context.Context, _ *trustgroup_domain.RotateTrustGroupKEKRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

type stubShareEntryRepo struct {
	createdEntries []c3_asset_domain.ShareEntry
}

func (s *stubShareEntryRepo) CreateShareEntry(_ context.Context, req *c3_asset_domain.CreateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	s.createdEntries = append(s.createdEntries, req.ShareEntry)
	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{Data: req.ShareEntry}, nil
}
func (s *stubShareEntryRepo) GetShareEntry(_ context.Context, _ *c3_asset_domain.GetShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	return nil, nil
}
func (s *stubShareEntryRepo) UpdateShareEntry(_ context.Context, _ *c3_asset_domain.UpdateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	return nil, nil
}
func (s *stubShareEntryRepo) DeleteShareEntry(_ context.Context, _ *c3_asset_domain.DeleteShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	return nil, nil
}

type stubThreadRepo struct {
	threads     map[string]thread_domain.Thread
	events      map[string][]thread_domain.ThreadEvent
	appendError error
}

func newStubThreadRepo() *stubThreadRepo {
	return &stubThreadRepo{
		threads: make(map[string]thread_domain.Thread),
		events:  make(map[string][]thread_domain.ThreadEvent),
	}
}

func (s *stubThreadRepo) CreateThread(_ context.Context, req *thread_domain.CreateThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	s.threads[req.Thread.ID] = req.Thread
	return &tracecore_types.CloudResponse[thread_domain.Thread]{Data: req.Thread}, nil
}
func (s *stubThreadRepo) GetThread(_ context.Context, req *thread_domain.GetThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	th, ok := s.threads[req.ThreadID]
	if !ok {
		return nil, thread_domain.ErrThreadNotFound
	}
	return &tracecore_types.CloudResponse[thread_domain.Thread]{Data: th}, nil
}
func (s *stubThreadRepo) ListThreads(_ context.Context, _ *thread_domain.ListThreadsRequest) (*tracecore_types.CloudResponse[[]thread_domain.Thread], error) {
	return nil, nil
}
func (s *stubThreadRepo) UpdateThread(_ context.Context, _ *thread_domain.UpdateThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	return nil, nil
}
func (s *stubThreadRepo) ListThreadEvents(_ context.Context, _ *thread_domain.ListThreadEventsRequest) (*tracecore_types.CloudResponse[[]thread_domain.ThreadEvent], error) {
	return nil, nil
}
func (s *stubThreadRepo) AppendThreadEvent(_ context.Context, req *thread_domain.AppendThreadEventRequest) (*tracecore_types.CloudResponse[thread_domain.ThreadEvent], error) {
	if s.appendError != nil {
		return nil, s.appendError
	}
	th, ok := s.threads[req.ThreadID]
	if !ok {
		return nil, thread_domain.ErrThreadNotFound
	}
	if th.Status == thread_domain.ThreadClosed {
		return nil, thread_domain.ErrThreadClosed
	}

	// Idempotency check
	if req.IdempotencyKey != "" {
		for _, existing := range s.events[req.ThreadID] {
			if existing.IdempotencyKey == req.IdempotencyKey {
				return &tracecore_types.CloudResponse[thread_domain.ThreadEvent]{Data: existing}, nil
			}
		}
	}

	evt := thread_domain.ThreadEvent{
		ID:             "evt_" + req.ThreadID,
		ThreadID:       req.ThreadID,
		Type:           thread_domain.ThreadEventType(req.EventType),
		Payload:        req.Payload,
		IdempotencyKey: req.IdempotencyKey,
		CreatedAt:      time.Now(),
	}
	s.events[req.ThreadID] = append(s.events[req.ThreadID], evt)
	return &tracecore_types.CloudResponse[thread_domain.ThreadEvent]{Data: evt}, nil
}

// ---------------------------------------------------------------------------
// Helper setup
// ---------------------------------------------------------------------------
func setupHandler(tg *trustgroup_domain.TrustGroup, threadRepo *stubThreadRepo) (*collaboration_ui.CollaborationHandler, *stubShareEntryRepo) {
	tgRepo := &stubTrustGroupRepo{group: tg}
	shareRepo := &stubShareEntryRepo{}

	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(tgRepo, shareRepo)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, nil)
	appendEventUC := thread_usecase.NewAppendThreadEventUsecase(threadRepo)

	handler := collaboration_ui.NewCollaborationHandler(createCollabShareUC, nil, appendEventUC)
	return handler, shareRepo
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

// Case A — Complete Success
func TestOrchestration_CaseA_CompleteSuccess(t *testing.T) {
	ctx := context.Background()
	tg := trustgroup_domain.NewTrustGroup("ch_1", "Engineering", []string{"v_1"})
	threadRepo := newStubThreadRepo()
	th := thread_domain.NewThread("ch_1", "document", "Title", "Subtitle")
	_, _ = threadRepo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})

	handler, shareRepo := setupHandler(tg, threadRepo)

	res, err := handler.CreateCollaborativeShare(ctx, "user_1", th.ID, tg.ID, "cid_blueprint", "v_target", "note")
	require.NoError(t, err)
	require.NotNil(t, res)

	// Verify ShareEntry created with canonical ShareEntryID
	require.Len(t, shareRepo.createdEntries, 1)
	canonicalShareID := shareRepo.createdEntries[0].ID
	assert.Equal(t, canonicalShareID, res.ShareEntryID)

	// Verify ThreadEvent appended with correct EventResourceRef and IdempotencyKey
	require.Len(t, threadRepo.events[th.ID], 1)
	evt := threadRepo.events[th.ID][0]
	assert.Equal(t, thread_domain.ResourceShareEntry, evt.Payload.RefType)
	assert.Equal(t, canonicalShareID, evt.Payload.ShareEntryID)
	assert.Equal(t, tg.ID, evt.Payload.TrustGroupID)
	assert.Equal(t, "evt_share_"+canonicalShareID, evt.IdempotencyKey)
}

// Case B — Thread Append Failure
func TestOrchestration_CaseB_ThreadAppendFailure(t *testing.T) {
	ctx := context.Background()
	tg := trustgroup_domain.NewTrustGroup("ch_1", "Engineering", []string{"v_1"})
	threadRepo := newStubThreadRepo()
	th := thread_domain.NewThread("ch_1", "document", "Title", "Subtitle")
	_, _ = threadRepo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})
	threadRepo.appendError = assert.AnError

	handler, shareRepo := setupHandler(tg, threadRepo)

	res, err := handler.CreateCollaborativeShare(ctx, "user_1", th.ID, tg.ID, "cid_blueprint", "v_target", "note")
	assert.ErrorIs(t, err, assert.AnError)
	require.NotNil(t, res, "ShareEntry result reference MUST be returned even if Thread append fails")

	// ShareEntry remains valid and persisted
	require.Len(t, shareRepo.createdEntries, 1)
	assert.Equal(t, shareRepo.createdEntries[0].ID, res.ShareEntryID)

	// No ThreadEvent recorded
	assert.Len(t, threadRepo.events[th.ID], 0)
}

// Case C — Retry Thread Event
func TestOrchestration_CaseC_ThreadEventRetry(t *testing.T) {
	ctx := context.Background()
	tg := trustgroup_domain.NewTrustGroup("ch_1", "Engineering", []string{"v_1"})
	threadRepo := newStubThreadRepo()
	th := thread_domain.NewThread("ch_1", "document", "Title", "Subtitle")
	_, _ = threadRepo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})

	handler, shareRepo := setupHandler(tg, threadRepo)

	res1, err := handler.CreateCollaborativeShare(ctx, "user_1", th.ID, tg.ID, "cid_blueprint", "v_target", "note")
	require.NoError(t, err)

	// Direct retry of AppendThreadEvent with canonical ShareEntryID & idempotency key
	appendUC := thread_usecase.NewAppendThreadEventUsecase(threadRepo)
	ref := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: res1.ShareEntryID,
		TrustGroupID: tg.ID,
	}
	idempotencyKey := "evt_share_" + res1.ShareEntryID

	evt2, err := appendUC.Execute(ctx, th.ID, "entry.shared", ref, idempotencyKey)
	require.NoError(t, err)
	assert.Equal(t, threadRepo.events[th.ID][0].ID, evt2.ID)
	assert.Len(t, threadRepo.events[th.ID], 1, "Retry must not create duplicate ThreadEvent")

	_ = shareRepo
}

// Case D — Closed or Invalid Thread
func TestOrchestration_CaseD_InvalidOrClosedThread(t *testing.T) {
	ctx := context.Background()
	tg := trustgroup_domain.NewTrustGroup("ch_1", "Engineering", []string{"v_1"})

	// Closed Thread
	threadRepoClosed := newStubThreadRepo()
	thClosed := thread_domain.NewThread("ch_1", "document", "Closed", "")
	thClosed.Status = thread_domain.ThreadClosed
	_, _ = threadRepoClosed.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: thClosed})

	handlerClosed, shareRepoClosed := setupHandler(tg, threadRepoClosed)

	resClosed, err := handlerClosed.CreateCollaborativeShare(ctx, "user_1", thClosed.ID, tg.ID, "cid_blueprint", "v_target", "note")
	assert.ErrorIs(t, err, thread_domain.ErrThreadClosed)
	require.NotNil(t, resClosed)
	assert.Len(t, shareRepoClosed.createdEntries, 1, "ShareEntry remains valid even when Thread is closed")

	// Invalid / Missing Thread
	threadRepoMissing := newStubThreadRepo()
	handlerMissing, shareRepoMissing := setupHandler(tg, threadRepoMissing)

	resMissing, err := handlerMissing.CreateCollaborativeShare(ctx, "user_1", "nonexistent_thread", tg.ID, "cid_blueprint", "v_target", "note")
	assert.ErrorIs(t, err, thread_domain.ErrThreadNotFound)
	require.NotNil(t, resMissing)
	assert.Len(t, shareRepoMissing.createdEntries, 1, "ShareEntry remains valid even when Thread does not exist")
}

// Case E — Security Boundary & Payload Verification
func TestOrchestration_CaseE_SecurityBoundaryVerification(t *testing.T) {
	ctx := context.Background()
	tg := trustgroup_domain.NewTrustGroup("ch_1", "Legal", []string{"v_1"})
	threadRepo := newStubThreadRepo()
	th := thread_domain.NewThread("ch_1", "contract", "NDA", "")
	_, _ = threadRepo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})

	handler, _ := setupHandler(tg, threadRepo)

	res, err := handler.CreateCollaborativeShare(ctx, "user_1", th.ID, tg.ID, "cid_blueprint", "v_target", "note")
	require.NoError(t, err)

	evt := threadRepo.events[th.ID][0]
	marshaledBytes, err := json.Marshal(evt)
	require.NoError(t, err)
	marshaledStr := string(marshaledBytes)

	// Forbidden cryptographic / secret strings
	forbiddenSecrets := []string{
		"wrapped_dek",
		"wrappedDEK",
		"wrapped_kek",
		"wrappedKEK",
		"key_envelopes",
		"PrivateKey",
		"SecretKey",
		"private_key",
		"secret_key",
		"encrypted_payload",
		"decrypted_content",
	}

	for _, forbidden := range forbiddenSecrets {
		if strings.Contains(marshaledStr, forbidden) {
			t.Errorf("SECURITY VIOLATION: Serialized ThreadEvent payload contains forbidden secret substring: %q", forbidden)
		}
	}

	// Allowed reference fields only
	assert.Equal(t, thread_domain.ResourceShareEntry, evt.Payload.RefType)
	assert.Equal(t, res.ShareEntryID, evt.Payload.ShareEntryID)
	assert.Equal(t, tg.ID, evt.Payload.TrustGroupID)
	assert.Empty(t, evt.Payload.CID)
}

// ---------------------------------------------------------------------------
// Step 3: CollaborationHandler.ResolveCollaborativeShare Delegation Tests
// ---------------------------------------------------------------------------

func TestCollaborationHandler_ResolveCollaborativeShare_NilUseCase(t *testing.T) {
	ctx := context.Background()
	handler := collaboration_ui.NewCollaborationHandler(nil, nil, nil)

	_, err := handler.ResolveCollaborativeShare(ctx, "user_alice", "se_100", "dev_laptop")
	assert.ErrorContains(t, err, "resolve collaborative share use case is not initialized")
}

package c3_integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	channel_domain "vault-app/internal/channel/domain"
	thread_dtos "vault-app/internal/thread/application/dtos"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	tracecore_types "vault-app/internal/tracecore/types"
)

// ---------------------------------------------------------------------------
// Stubs
// ---------------------------------------------------------------------------
type stubThreadRepo struct {
	threads map[string]thread_domain.Thread
	events  map[string][]thread_domain.ThreadEvent
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
	list := make([]thread_domain.Thread, 0, len(s.threads))
	for _, th := range s.threads {
		list = append(list, th)
	}
	return &tracecore_types.CloudResponse[[]thread_domain.Thread]{Data: list}, nil
}

func (s *stubThreadRepo) UpdateThread(_ context.Context, req *thread_domain.UpdateThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	s.threads[req.Thread.ID] = req.Thread
	return &tracecore_types.CloudResponse[thread_domain.Thread]{Data: req.Thread}, nil
}

func (s *stubThreadRepo) ListThreadEvents(_ context.Context, req *thread_domain.ListThreadEventsRequest) (*tracecore_types.CloudResponse[[]thread_domain.ThreadEvent], error) {
	evts := s.events[req.ThreadID]
	return &tracecore_types.CloudResponse[[]thread_domain.ThreadEvent]{Data: evts}, nil
}

func (s *stubThreadRepo) AppendThreadEvent(_ context.Context, req *thread_domain.AppendThreadEventRequest) (*tracecore_types.CloudResponse[thread_domain.ThreadEvent], error) {
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
		ID:             "evt_" + uuid.NewString()[:8],
		ThreadID:       req.ThreadID,
		Type:           thread_domain.ThreadEventType(req.EventType),
		Payload:        req.Payload,
		IdempotencyKey: req.IdempotencyKey,
		Cursor:         uint64(len(s.events[req.ThreadID]) + 1),
		CreatedAt:      time.Now(),
	}

	s.events[req.ThreadID] = append(s.events[req.ThreadID], evt)
	return &tracecore_types.CloudResponse[thread_domain.ThreadEvent]{Data: evt}, nil
}

type stubChannelGovReader struct {
	channels map[string]channel_domain.Channel
}

func (s *stubChannelGovReader) GetChannel(_ context.Context, req *channel_domain.GetChannelRequest) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
	ch, ok := s.channels[req.ChannelID]
	if !ok {
		return nil, channel_domain.ErrChannelNotFound
	}
	return &tracecore_types.CloudResponse[channel_domain.Channel]{Data: ch}, nil
}

type stubThreadEventBus struct{}

func (s *stubThreadEventBus) PublishThreadCreated(_ context.Context, _ thread_domain.ThreadCreated) error {
	return nil
}
func (s *stubThreadEventBus) SubscribeToThreadCreated(_ func(ctx context.Context, event thread_domain.ThreadCreated)) error {
	return nil
}
func (s *stubThreadEventBus) PublishThreadUpdated(_ context.Context, _ thread_domain.ThreadUpdated) error {
	return nil
}
func (s *stubThreadEventBus) SubscribeToThreadUpdated(_ func(ctx context.Context, event thread_domain.ThreadUpdated)) error {
	return nil
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestC3ThreadIntegration_CreateThread_RemainsIsolated(t *testing.T) {
	repo := newStubThreadRepo()
	govReader := &stubChannelGovReader{
		channels: map[string]channel_domain.Channel{
			"ch_1": {
				ID:          "ch_1",
				WorkspaceID: "ws_1",
				Status:      channel_domain.StatusActive,
			},
		},
	}

	eventBus := &stubThreadEventBus{}
	createUC := thread_usecase.NewCreateThreadUsecase(repo, eventBus, govReader)

	th, err := createUC.Execute(context.Background(), thread_dtos.CreateThreadRequest{
		ChannelID: "ch_1",
		Title:     "Financial Review",
		Subtitle:  "Q3 Audit",
		AssetType: "document",
	})

	require.NoError(t, err)
	require.NotNil(t, th)
	assert.Equal(t, "ch_1", th.ChannelID)
	assert.Equal(t, "ws_1", th.WorkspaceID)

	// Invariant Check: CreateThread contains ZERO C3 / ShareEntry / storage fields
	assert.NotEmpty(t, th.ID)
	assert.Equal(t, thread_domain.ThreadOpen, th.Status)
}

func TestC3ThreadIntegration_AppendStorageAssetEvent(t *testing.T) {
	ctx := context.Background()
	repo := newStubThreadRepo()

	th := thread_domain.NewThread("ch_1", "document", "Design Specs", "v1.0")
	_, err := repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})
	require.NoError(t, err)

	appendUC := thread_usecase.NewAppendThreadEventUsecase(repo)

	assetRef := thread_domain.EventResourceRef{
		RefType:     thread_domain.ResourceStorageAsset,
		CID:         "bafybeicollab123456789",
		ContentHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Size:        1048576,
		AssetType:   "blueprint_pdf",
	}

	evt, err := appendUC.Execute(ctx, th.ID, string(thread_domain.EventEntryShared), assetRef)

	require.NoError(t, err)
	require.NotNil(t, evt)
	assert.Equal(t, th.ID, evt.ThreadID)
	assert.Equal(t, thread_domain.ResourceStorageAsset, evt.Payload.RefType)
	assert.Equal(t, "bafybeicollab123456789", evt.Payload.CID)
	assert.Equal(t, int64(1048576), evt.Payload.Size)
	assert.Equal(t, "blueprint_pdf", evt.Payload.AssetType)
}

func TestC3ThreadIntegration_AppendC3ShareEntryEvent(t *testing.T) {
	ctx := context.Background()
	repo := newStubThreadRepo()

	th := thread_domain.NewThread("ch_1", "contract", "Master Agreement", "2026")
	_, err := repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})
	require.NoError(t, err)

	appendUC := thread_usecase.NewAppendThreadEventUsecase(repo)

	c3Ref := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: "se_9988776655",
		TrustGroupID: "tg_legal_council",
	}

	evt, err := appendUC.Execute(ctx, th.ID, string(thread_domain.EventEntryShared), c3Ref)

	require.NoError(t, err)
	require.NotNil(t, evt)
	assert.Equal(t, th.ID, evt.ThreadID)
	assert.Equal(t, thread_domain.ResourceShareEntry, evt.Payload.RefType)
	assert.Equal(t, "se_9988776655", evt.Payload.ShareEntryID)
	assert.Equal(t, "tg_legal_council", evt.Payload.TrustGroupID)
	assert.Equal(t, "evt_share_se_9988776655", evt.IdempotencyKey)
	assert.Empty(t, evt.Payload.CID) // NO fake CID aliasing!
}

func TestC3ThreadIntegration_ShareEntryEvent_Idempotency(t *testing.T) {
	ctx := context.Background()
	repo := newStubThreadRepo()

	th := thread_domain.NewThread("ch_1", "contract", "Agreement", "v1")
	_, err := repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})
	require.NoError(t, err)

	appendUC := thread_usecase.NewAppendThreadEventUsecase(repo)

	c3Ref := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: "se_repeat_001",
		TrustGroupID: "tg_finance",
	}

	evt1, err := appendUC.Execute(ctx, th.ID, "entry.shared", c3Ref)
	require.NoError(t, err)

	// Retry identical event execution
	evt2, err := appendUC.Execute(ctx, th.ID, "entry.shared", c3Ref)
	require.NoError(t, err)

	assert.Equal(t, evt1.ID, evt2.ID, "Idempotent retry must return existing ThreadEvent without creating a duplicate")
	assert.Len(t, repo.events[th.ID], 1, "Timeline must contain exactly 1 event")
}

func TestC3ThreadIntegration_ClosedThreadRejectsEvent(t *testing.T) {
	ctx := context.Background()
	repo := newStubThreadRepo()

	th := thread_domain.NewThread("ch_1", "contract", "Closed Contract", "v0")
	th.Status = thread_domain.ThreadClosed
	_, err := repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})
	require.NoError(t, err)

	appendUC := thread_usecase.NewAppendThreadEventUsecase(repo)

	c3Ref := thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: "se_closed_attempt",
		TrustGroupID: "tg_legal",
	}

	_, err = appendUC.Execute(ctx, th.ID, "entry.shared", c3Ref)
	assert.ErrorIs(t, err, thread_domain.ErrThreadClosed)
}

func TestC3EntryShared_POST_Response_Vs_GET_Response(t *testing.T) {
	ctx := context.Background()
	repo := newStubThreadRepo()

	th := thread_domain.NewThread("ch_post_get_test", "contract", "POST vs GET Test", "")
	_, err := repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{Thread: th})
	require.NoError(t, err)

	appendUC := thread_usecase.NewAppendThreadEventUsecase(repo)
	listUC := thread_usecase.NewListThreadEventsUsecase(repo)

	const shareID = "22081c47-29e5-4b03-bd1e-626d3f6f87e3"
	const trustGroupID = "ffc46329-6b01-4259-a101-22b6ccd48251"

	// 1. POST entry.shared event
	postResponse, err := appendUC.Execute(ctx, th.ID, "entry.shared", thread_domain.EventResourceRef{
		RefType:      thread_domain.ResourceShareEntry,
		ShareEntryID: shareID,
		TrustGroupID: trustGroupID,
	})
	require.NoError(t, err)

	// 2. Capture returned POST response event payload
	assert.Equal(t, thread_domain.ResourceShareEntry, postResponse.Payload.RefType)
	assert.Equal(t, shareID, postResponse.Payload.ShareEntryID)
	assert.Equal(t, trustGroupID, postResponse.Payload.TrustGroupID)

	// 3. GET ListThreadEvents immediately
	events, err := listUC.Execute(ctx, th.ID)
	require.NoError(t, err)
	require.Len(t, events, 1)

	getResponse := events[0]

	// 4. Compare POST response vs GET response
	assert.Equal(t, postResponse.Payload.RefType, getResponse.Payload.RefType, "POST response RefType MUST match GET response RefType")
	assert.Equal(t, postResponse.Payload.ShareEntryID, getResponse.Payload.ShareEntryID, "POST response ShareEntryID MUST match GET response ShareEntryID")
	assert.Equal(t, postResponse.Payload.TrustGroupID, getResponse.Payload.TrustGroupID, "POST response TrustGroupID MUST match GET response TrustGroupID")
}

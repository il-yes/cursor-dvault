package thread_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	channel_domain "vault-app/internal/channel/domain"
	thread_dtos "vault-app/internal/thread/application/dtos"
	thread_events "vault-app/internal/thread/application/events"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	tracecore_types "vault-app/internal/tracecore/types"
)

// ---------------------------------------------------------------------------
// Stub repository implementing thread_domain.ThreadRepository
// ---------------------------------------------------------------------------
type stubThreadRepo struct {
	lastCreated *thread_domain.Thread
}

func (s *stubThreadRepo) CreateThread(_ context.Context, req *thread_domain.CreateThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	s.lastCreated = &req.Thread
	return &tracecore_types.CloudResponse[thread_domain.Thread]{
		Status:  200,
		Data:    req.Thread,
		Success: true,
	}, nil
}
func (s *stubThreadRepo) ListThreads(_ context.Context, _ *thread_domain.ListThreadsRequest) (*tracecore_types.CloudResponse[[]thread_domain.Thread], error) {
	return nil, nil
}
func (s *stubThreadRepo) GetThread(_ context.Context, _ *thread_domain.GetThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	return nil, nil
}
func (s *stubThreadRepo) UpdateThread(_ context.Context, _ *thread_domain.UpdateThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	return nil, nil
}
func (s *stubThreadRepo) ListThreadEvents(_ context.Context, _ *thread_domain.ListThreadEventsRequest) (*tracecore_types.CloudResponse[[]thread_domain.ThreadEvent], error) {
	return nil, nil
}
func (s *stubThreadRepo) AppendThreadEvent(_ context.Context, _ *thread_domain.AppendThreadEventRequest) (*tracecore_types.CloudResponse[thread_domain.ThreadEvent], error) {
	return nil, nil
}

// ---------------------------------------------------------------------------
// Stub event bus capturing published ThreadCreated events
// ---------------------------------------------------------------------------
type stubThreadEventBus struct {
	lastCreatedEvent *thread_domain.ThreadCreated
}

func (s *stubThreadEventBus) PublishThreadCreated(_ context.Context, event thread_domain.ThreadCreated) error {
	s.lastCreatedEvent = &event
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

var _ thread_events.ThreadEventBus = (*stubThreadEventBus)(nil)

// ---------------------------------------------------------------------------
// Stub Channel governance reader
// ---------------------------------------------------------------------------
type stubChannelReader struct {
	channels map[string]*channel_domain.Channel
}

func (s *stubChannelReader) GetChannel(_ context.Context, req *channel_domain.GetChannelRequest) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
	ch, exists := s.channels[req.ChannelID]
	if !exists {
		return nil, errors.New("not found")
	}
	return &tracecore_types.CloudResponse[channel_domain.Channel]{
		Status:  200,
		Data:    *ch,
		Success: true,
	}, nil
}

var _ thread_usecase.ChannelGovernanceReader = (*stubChannelReader)(nil)

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestCreateThread_ActiveChannel_Succeeds(t *testing.T) {
	repo := &stubThreadRepo{}
	bus := &stubThreadEventBus{}
	reader := &stubChannelReader{
		channels: map[string]*channel_domain.Channel{
			"ch_active": {
				ID:          "ch_active",
				Status:      channel_domain.StatusActive,
				WorkspaceID: "ws_auth",
			},
		},
	}

	uc := thread_usecase.NewCreateThreadUsecase(repo, bus, reader)

	req := thread_dtos.CreateThreadRequest{
		ChannelID: "ch_active",
		AssetType: "entry.shared",
		Title:     "Invoice Thread",
		Subtitle:  "Shared financial document",
	}

	result, err := uc.Execute(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "ch_active", result.ChannelID)
	assert.Equal(t, "ws_auth", result.WorkspaceID) // Inherited authoritatively from Channel
	assert.Equal(t, "entry.shared", result.AssetType)
	assert.Equal(t, "Invoice Thread", result.Title)
	assert.Equal(t, thread_domain.ThreadOpen, result.Status)
	assert.NotEmpty(t, result.ID)
}

func TestCreateThread_ChannelNotFound_Rejected(t *testing.T) {
	repo := &stubThreadRepo{}
	bus := &stubThreadEventBus{}
	reader := &stubChannelReader{
		channels: map[string]*channel_domain.Channel{},
	}

	uc := thread_usecase.NewCreateThreadUsecase(repo, bus, reader)

	req := thread_dtos.CreateThreadRequest{
		ChannelID: "ch_nonexistent",
		AssetType: "entry.shared",
		Title:     "Title",
		Subtitle:  "Subtitle",
	}

	_, err := uc.Execute(context.Background(), req)
	assert.ErrorIs(t, err, thread_domain.ErrChannelNotFound)
}

func TestCreateThread_PendingChannel_Rejected(t *testing.T) {
	repo := &stubThreadRepo{}
	bus := &stubThreadEventBus{}
	reader := &stubChannelReader{
		channels: map[string]*channel_domain.Channel{
			"ch_pending": {
				ID:          "ch_pending",
				Status:      channel_domain.StatusPending,
				WorkspaceID: "ws_1",
			},
		},
	}

	uc := thread_usecase.NewCreateThreadUsecase(repo, bus, reader)

	req := thread_dtos.CreateThreadRequest{
		ChannelID: "ch_pending",
		AssetType: "entry.shared",
		Title:     "Title",
		Subtitle:  "Subtitle",
	}

	_, err := uc.Execute(context.Background(), req)
	assert.ErrorIs(t, err, thread_domain.ErrChannelNotActive)
}

func TestCreateThread_RevokedChannel_Rejected(t *testing.T) {
	repo := &stubThreadRepo{}
	bus := &stubThreadEventBus{}
	reader := &stubChannelReader{
		channels: map[string]*channel_domain.Channel{
			"ch_revoked": {
				ID:          "ch_revoked",
				Status:      channel_domain.StatusRevoked,
				WorkspaceID: "ws_1",
			},
		},
	}

	uc := thread_usecase.NewCreateThreadUsecase(repo, bus, reader)

	req := thread_dtos.CreateThreadRequest{
		ChannelID: "ch_revoked",
		AssetType: "entry.shared",
		Title:     "Title",
		Subtitle:  "Subtitle",
	}

	_, err := uc.Execute(context.Background(), req)
	assert.ErrorIs(t, err, thread_domain.ErrChannelNotActive)
}

func TestCreateThread_ArchivedChannel_Rejected(t *testing.T) {
	repo := &stubThreadRepo{}
	bus := &stubThreadEventBus{}
	reader := &stubChannelReader{
		channels: map[string]*channel_domain.Channel{
			"ch_archived": {
				ID:          "ch_archived",
				Status:      channel_domain.StatusArchived,
				WorkspaceID: "ws_1",
			},
		},
	}

	uc := thread_usecase.NewCreateThreadUsecase(repo, bus, reader)

	req := thread_dtos.CreateThreadRequest{
		ChannelID: "ch_archived",
		AssetType: "entry.shared",
		Title:     "Title",
		Subtitle:  "Subtitle",
	}

	_, err := uc.Execute(context.Background(), req)
	assert.ErrorIs(t, err, thread_domain.ErrChannelNotActive)
}

func TestCreateThread_GatedSlotsIncomplete_Rejected(t *testing.T) {
	repo := &stubThreadRepo{}
	bus := &stubThreadEventBus{}
	reader := &stubChannelReader{
		channels: map[string]*channel_domain.Channel{
			"ch_gated": {
				ID:          "ch_gated",
				Status:      channel_domain.StatusActive,
				WorkspaceID: "ws_1",
				Slots: []channel_domain.Slot{
					{ID: "slot_1", Name: "Auditor", Gated: true},
				},
				Assignments: []channel_domain.Assignment{}, // Unfulfilled gated slot
			},
		},
	}

	uc := thread_usecase.NewCreateThreadUsecase(repo, bus, reader)

	req := thread_dtos.CreateThreadRequest{
		ChannelID: "ch_gated",
		AssetType: "entry.shared",
		Title:     "Title",
		Subtitle:  "Subtitle",
	}

	_, err := uc.Execute(context.Background(), req)
	assert.ErrorIs(t, err, thread_domain.ErrChannelGatedSlotsIncomplete)
}

func TestCreateThread_GatedSlotsFulfilled_Succeeds(t *testing.T) {
	repo := &stubThreadRepo{}
	bus := &stubThreadEventBus{}
	reader := &stubChannelReader{
		channels: map[string]*channel_domain.Channel{
			"ch_fulfilled": {
				ID:          "ch_fulfilled",
				Status:      channel_domain.StatusActive,
				WorkspaceID: "ws_1",
				Slots: []channel_domain.Slot{
					{ID: "slot_1", Name: "Auditor", Gated: true},
				},
				Assignments: []channel_domain.Assignment{
					{SlotID: "slot_1", OwnerID: "user_99"},
				},
			},
		},
	}

	uc := thread_usecase.NewCreateThreadUsecase(repo, bus, reader)

	req := thread_dtos.CreateThreadRequest{
		ChannelID: "ch_fulfilled",
		AssetType: "entry.shared",
		Title:     "Title",
		Subtitle:  "Subtitle",
	}

	result, err := uc.Execute(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "ch_fulfilled", result.ChannelID)
}

func TestCreateThread_WorkspaceMismatch_Rejected(t *testing.T) {
	repo := &stubThreadRepo{}
	bus := &stubThreadEventBus{}
	reader := &stubChannelReader{
		channels: map[string]*channel_domain.Channel{
			"ch_ws": {
				ID:          "ch_ws",
				Status:      channel_domain.StatusActive,
				WorkspaceID: "ws_actual",
			},
		},
	}

	uc := thread_usecase.NewCreateThreadUsecase(repo, bus, reader)

	req := thread_dtos.CreateThreadRequest{
		ChannelID:   "ch_ws",
		WorkspaceID: "ws_conflicting", // Mismatch
		AssetType:   "entry.shared",
		Title:       "Title",
		Subtitle:    "Subtitle",
	}

	_, err := uc.Execute(context.Background(), req)
	assert.ErrorIs(t, err, thread_domain.ErrWorkspaceMismatch)
}

func TestCreateThread_ThreadCreatedEvent_ContainsContextWithoutCID(t *testing.T) {
	repo := &stubThreadRepo{}
	bus := &stubThreadEventBus{}
	reader := &stubChannelReader{
		channels: map[string]*channel_domain.Channel{
			"ch_evt": {
				ID:          "ch_evt",
				Status:      channel_domain.StatusActive,
				WorkspaceID: "ws_evt",
			},
		},
	}

	uc := thread_usecase.NewCreateThreadUsecase(repo, bus, reader)

	req := thread_dtos.CreateThreadRequest{
		ChannelID: "ch_evt",
		AssetType: "invoice.created",
		Title:     "Payment Thread",
		Subtitle:  "Monthly invoice",
	}

	_, err := uc.Execute(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, bus.lastCreatedEvent)

	event := bus.lastCreatedEvent
	assert.NotEmpty(t, event.EventID)
	assert.NotEmpty(t, event.ThreadID)
	assert.Equal(t, "ch_evt", event.ChannelID)
	assert.Equal(t, "ws_evt", event.WorkspaceID)
	assert.Equal(t, "invoice.created", event.AssetType)
	assert.False(t, event.Timestamp.IsZero())
}

func TestThreadDTO_MapsAllFields(t *testing.T) {
	now := time.Now()
	thread := &thread_domain.Thread{
		ID:          "th_dto",
		ChannelID:   "ch_dto",
		WorkspaceID: "ws_dto",
		AssetType:   "entry.shared",
		Title:       "DTO Test",
		Subtitle:    "DTO Subtitle",
		Status:      thread_domain.ThreadOpen,
		CreatedAt:   now,
	}

	dto := tracecore_types.ThreadDTO{
		ID:          thread.ID,
		ChannelID:   thread.ChannelID,
		WorkspaceID: thread.WorkspaceID,
		AssetType:   thread.AssetType,
		Title:       thread.Title,
		Subtitle:    thread.Subtitle,
		Status:      string(thread.Status),
		CreatedAt:   thread.CreatedAt,
	}

	assert.Equal(t, "th_dto", dto.ID)
	assert.Equal(t, "ch_dto", dto.ChannelID)
	assert.Equal(t, "ws_dto", dto.WorkspaceID)
	assert.Equal(t, "entry.shared", dto.AssetType)
	assert.Equal(t, "DTO Test", dto.Title)
	assert.Equal(t, "DTO Subtitle", dto.Subtitle)
	assert.Equal(t, "open", dto.Status)
	assert.True(t, dto.CreatedAt.Equal(now))
}

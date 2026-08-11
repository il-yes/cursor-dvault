package channel_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	channel_application "vault-app/internal/channel/application"
	channel_usecase "vault-app/internal/channel/application/channel_lifecycle_usecases"
	channel_events "vault-app/internal/channel/application/events"
	channel_domain "vault-app/internal/channel/domain"
	tracecore_types "vault-app/internal/tracecore/types"
)

// ------------------------------------------------------------------------------------------------------------
// HELPERS
// ------------------------------------------------------------------------------------------------------------

func activeChannel() channel_domain.Channel {
	now := time.Now().UTC()
	return channel_domain.Channel{
		ID:          "channel-001",
		TemplateID:  "tpl-001",
		Title:       "Battery Engineering Review",
		Status:      channel_domain.StatusActive,
		WorkspaceID: "workspace-001",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func revokedChannel() channel_domain.Channel {
	now := time.Now().UTC()
	revokedAt := now.Add(-1 * time.Hour)
	return channel_domain.Channel{
		ID:          "channel-002",
		TemplateID:  "tpl-001",
		Title:       "Revoked Channel",
		Status:      channel_domain.StatusRevoked,
		WorkspaceID: "workspace-001",
		RevokedAt:   &revokedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func pendingChannel() channel_domain.Channel {
	now := time.Now().UTC()
	return channel_domain.Channel{
		ID:          "channel-003",
		TemplateID:  "tpl-001",
		Title:       "Pending Channel",
		Status:      channel_domain.StatusPending,
		WorkspaceID: "workspace-001",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func validArchiveRequest() *channel_application.ArchiveChannelRequest {
	return &channel_application.ArchiveChannelRequest{
		ChannelID:   "channel-001",
		WorkspaceID: "workspace-001",
	}
}

// ------------------------------------------------------------------------------------------------------------
// TEST — Execute
// ------------------------------------------------------------------------------------------------------------

func TestArchiveChannelUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	ch := activeChannel()

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			require.Equal(t, "channel-001", req.ChannelID)

			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: ch,
			}, nil
		},
		updateFn: func(
			ctx context.Context,
			req *channel_domain.UpdateChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			require.Equal(t, channel_domain.StatusArchived, req.Channel.Status)
			require.NotNil(t, req.Channel.ArchivedAt)

			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: req.Channel,
			}, nil
		},
	}

	bus := &channelEventBusMock{}

	uc := channel_usecase.NewArchiveChannelUsecase(repo, bus)

	err := uc.Execute(ctx, validArchiveRequest())

	require.NoError(t, err)
	require.Len(t, bus.publishedArchivedEvents, 1)
	require.Equal(t, "channel-001", bus.publishedArchivedEvents[0].ChannelID)
	require.Equal(t, "workspace-001", bus.publishedArchivedEvents[0].WorkspaceID)
}

func TestArchiveChannelUsecase_Execute_RevokedChannelNotArchivable(t *testing.T) {
	ctx := context.Background()
	ch := revokedChannel()

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: ch,
			}, nil
		},
	}

	bus := &channelEventBusMock{}

	uc := channel_usecase.NewArchiveChannelUsecase(repo, bus)

	err := uc.Execute(ctx, validArchiveRequest())

	require.Error(t, err)
	require.EqualError(t, err, channel_domain.ErrChannelNotArchivable.Error())
	require.Empty(t, bus.publishedArchivedEvents)
}

func TestArchiveChannelUsecase_Execute_PendingChannelNotArchivable(t *testing.T) {
	ctx := context.Background()
	ch := pendingChannel()

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: ch,
			}, nil
		},
	}

	bus := &channelEventBusMock{}

	uc := channel_usecase.NewArchiveChannelUsecase(repo, bus)

	err := uc.Execute(ctx, validArchiveRequest())

	require.Error(t, err)
	require.EqualError(t, err, channel_domain.ErrChannelNotArchivable.Error())
	require.Empty(t, bus.publishedArchivedEvents)
}

func TestArchiveChannelUsecase_Execute_AlreadyArchivedNotArchivable(t *testing.T) {
	ctx := context.Background()
	archivedAt := time.Now().UTC()
	ch := channel_domain.Channel{
		ID:          "channel-004",
		Status:      channel_domain.StatusArchived,
		WorkspaceID: "workspace-001",
		ArchivedAt:  &archivedAt,
	}

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: ch,
			}, nil
		},
	}

	bus := &channelEventBusMock{}

	uc := channel_usecase.NewArchiveChannelUsecase(repo, bus)

	err := uc.Execute(ctx, validArchiveRequest())

	require.Error(t, err)
	require.EqualError(t, err, channel_domain.ErrChannelNotArchivable.Error())
	require.Empty(t, bus.publishedArchivedEvents)
}

func TestArchiveChannelUsecase_Execute_GetChannelError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("repository get error")

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return nil, expectedErr
		},
	}

	bus := &channelEventBusMock{}

	uc := channel_usecase.NewArchiveChannelUsecase(repo, bus)

	err := uc.Execute(ctx, validArchiveRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Empty(t, bus.publishedArchivedEvents)
}

func TestArchiveChannelUsecase_Execute_GetChannelNilResponse(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return nil, nil
		},
	}

	bus := &channelEventBusMock{}

	uc := channel_usecase.NewArchiveChannelUsecase(repo, bus)

	err := uc.Execute(ctx, validArchiveRequest())

	require.Error(t, err)
	require.EqualError(t, err, channel_domain.ErrRepositoryResponse.Error())
}

func TestArchiveChannelUsecase_Execute_UpdateChannelError(t *testing.T) {
	ctx := context.Background()
	ch := activeChannel()
	expectedErr := errors.New("repository update error")

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: ch,
			}, nil
		},
		updateFn: func(
			ctx context.Context,
			req *channel_domain.UpdateChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return nil, expectedErr
		},
	}

	bus := &channelEventBusMock{}

	uc := channel_usecase.NewArchiveChannelUsecase(repo, bus)

	err := uc.Execute(ctx, validArchiveRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Empty(t, bus.publishedArchivedEvents)
}

func TestArchiveChannelUsecase_Execute_EventBusError(t *testing.T) {
	ctx := context.Background()
	ch := activeChannel()
	expectedErr := errors.New("event bus error")

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: ch,
			}, nil
		},
		updateFn: func(
			ctx context.Context,
			req *channel_domain.UpdateChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: req.Channel,
			}, nil
		},
	}

	bus := &channelEventBusMock{
		publishArchivedFn: func(
			ctx context.Context,
			event channel_domain.ChannelArchived,
		) error {
			return expectedErr
		},
	}

	uc := channel_usecase.NewArchiveChannelUsecase(repo, bus)

	err := uc.Execute(ctx, validArchiveRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
}

// ------------------------------------------------------------------------------------------------------------
// TEST — ValidateRequest (table-driven)
// ------------------------------------------------------------------------------------------------------------

func TestArchiveChannelUsecase_ValidateRequest(t *testing.T) {
	repo := &channelRepositoryMock{}
	bus := &channelEventBusMock{}

	uc := channel_usecase.NewArchiveChannelUsecase(repo, bus)

	tests := []struct {
		name          string
		request       *channel_application.ArchiveChannelRequest
		expectedError string
	}{
		{
			name:          "nil request",
			request:       nil,
			expectedError: channel_domain.ErrRequestRequired.Error(),
		},
		{
			name: "missing channel id",
			request: &channel_application.ArchiveChannelRequest{
				ChannelID:   "",
				WorkspaceID: "workspace-001",
			},
			expectedError: channel_domain.ErrChannelIDRequired.Error(),
		},
		{
			name: "missing workspace id",
			request: &channel_application.ArchiveChannelRequest{
				ChannelID:   "channel-001",
				WorkspaceID: "",
			},
			expectedError: channel_domain.ErrWorkspaceIDRequired.Error(),
		},
		{
			name: "valid request",
			request: &channel_application.ArchiveChannelRequest{
				ChannelID:   "channel-001",
				WorkspaceID: "workspace-001",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.ValidateRequest(tt.request)

			if tt.expectedError == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.expectedError)
		})
	}
}

// ------------------------------------------------------------------------------------------------------------
// TEST — ValidateDependencies (table-driven)
// ------------------------------------------------------------------------------------------------------------

func TestArchiveChannelUsecase_ValidateDependencies(t *testing.T) {
	tests := []struct {
		name          string
		repo          channel_domain.ChannelRepository
		bus           channel_events.ChannelEventBus
		expectedError string
	}{
		{
			name:          "nil repository",
			repo:          nil,
			bus:           &channelEventBusMock{},
			expectedError: channel_domain.ErrRepositoryNil.Error(),
		},
		{
			name:          "nil event bus",
			repo:          &channelRepositoryMock{},
			bus:           nil,
			expectedError: channel_domain.ErrChannelBusRequired.Error(),
		},
		{
			name: "valid dependencies",
			repo: &channelRepositoryMock{},
			bus:  &channelEventBusMock{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := channel_usecase.NewArchiveChannelUsecase(tt.repo, tt.bus)

			err := uc.ValidateDependencies()

			if tt.expectedError == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.expectedError)
		})
	}
}

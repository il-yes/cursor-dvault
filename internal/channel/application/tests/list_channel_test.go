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
// MOCKS
// ------------------------------------------------------------------------------------------------------------

type channelRepositoryMock struct {
	createFn func(
		ctx context.Context,
		req *channel_domain.CreateChannelRequest,
	) (*tracecore_types.CloudResponse[channel_domain.Channel], error)

	listFn func(
		ctx context.Context,
		req *channel_domain.ListChannelsRequest,
	) (*tracecore_types.CloudResponse[[]channel_domain.Channel], error)

	getFn func(
		ctx context.Context,
		req *channel_domain.GetChannelRequest,
	) (*tracecore_types.CloudResponse[channel_domain.Channel], error)

	deleteFn func(
		ctx context.Context,
		req *channel_domain.DeleteChannelRequest,
	) error

	updateFn func(
		ctx context.Context,
		req *channel_domain.UpdateChannelRequest,
	) (*tracecore_types.CloudResponse[channel_domain.Channel], error)

	activateFn func(
		ctx context.Context,
		req *channel_domain.AcceptInvitationRequest,
	) error

	revokeFn func(
		ctx context.Context,
		req *channel_domain.RevokeInvitationRequest,
	) error
}

func (m *channelRepositoryMock) CreateChannel(
	ctx context.Context,
	req *channel_domain.CreateChannelRequest,
) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}

	return &tracecore_types.CloudResponse[channel_domain.Channel]{
		Data: req.Channel,
	}, nil
}

func (m *channelRepositoryMock) ListChannels(
	ctx context.Context,
	req *channel_domain.ListChannelsRequest,
) (*tracecore_types.CloudResponse[[]channel_domain.Channel], error) {
	if m.listFn != nil {
		return m.listFn(ctx, req)
	}

	return &tracecore_types.CloudResponse[[]channel_domain.Channel]{
		Data: []channel_domain.Channel{},
	}, nil
}

func (m *channelRepositoryMock) GetChannel(
	ctx context.Context,
	req *channel_domain.GetChannelRequest,
) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
	if m.getFn != nil {
		return m.getFn(ctx, req)
	}

	return &tracecore_types.CloudResponse[channel_domain.Channel]{}, nil
}

func (m *channelRepositoryMock) DeleteChannel(
	ctx context.Context,
	req *channel_domain.DeleteChannelRequest,
) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, req)
	}

	return nil
}

func (m *channelRepositoryMock) UpdateChannel(
	ctx context.Context,
	req *channel_domain.UpdateChannelRequest,
) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, req)
	}

	return &tracecore_types.CloudResponse[channel_domain.Channel]{}, nil
}

func (m *channelRepositoryMock) ActivateChannel(
	ctx context.Context,
	req *channel_domain.AcceptInvitationRequest,
) error {
	if m.activateFn != nil {
		return m.activateFn(ctx, req)
	}

	return nil
}

func (m *channelRepositoryMock) RevokeChannel(
	ctx context.Context,
	req *channel_domain.RevokeInvitationRequest,
) error {
	if m.revokeFn != nil {
		return m.revokeFn(ctx, req)
	}

	return nil
}

type channelEventBusMock struct {
	publishCreatedFn func(
		ctx context.Context,
		event channel_domain.ChannelCreated,
	) error

	publishRevokedFn func(
		ctx context.Context,
		event channel_domain.ChannelRevoked,
	) error

	publishDeletedFn func(
		ctx context.Context,
		event channel_domain.ChannelDeleted,
	) error

	publishArchivedFn func(
		ctx context.Context,
		event channel_domain.ChannelArchived,
	) error

	publishedCreatedEvents  []channel_domain.ChannelCreated
	publishedRevokedEvents  []channel_domain.ChannelRevoked
	publishedDeletedEvents  []channel_domain.ChannelDeleted
	publishedArchivedEvents []channel_domain.ChannelArchived
}

func (m *channelEventBusMock) PublishChannelCreated(
	ctx context.Context,
	event channel_domain.ChannelCreated,
) error {
	m.publishedCreatedEvents = append(m.publishedCreatedEvents, event)

	if m.publishCreatedFn != nil {
		return m.publishCreatedFn(ctx, event)
	}

	return nil
}

func (m *channelEventBusMock) SubscribeToChannelCreated(
	handler func(ctx context.Context, event channel_domain.ChannelCreated),
) error {
	return nil
}

func (m *channelEventBusMock) PublishChannelRevoked(
	ctx context.Context,
	event channel_domain.ChannelRevoked,
) error {
	m.publishedRevokedEvents = append(m.publishedRevokedEvents, event)

	if m.publishRevokedFn != nil {
		return m.publishRevokedFn(ctx, event)
	}

	return nil
}

func (m *channelEventBusMock) SubscribeToChannelRevoked(
	handler func(ctx context.Context, event channel_domain.ChannelRevoked),
) error {
	return nil
}

func (m *channelEventBusMock) PublishChannelDeleted(
	ctx context.Context,
	event channel_domain.ChannelDeleted,
) error {
	m.publishedDeletedEvents = append(m.publishedDeletedEvents, event)

	if m.publishDeletedFn != nil {
		return m.publishDeletedFn(ctx, event)
	}

	return nil
}

func (m *channelEventBusMock) SubscribeToChannelDeleted(
	handler func(ctx context.Context, event channel_domain.ChannelDeleted),
) error {
	return nil
}

func (m *channelEventBusMock) PublishChannelArchived(
	ctx context.Context,
	event channel_domain.ChannelArchived,
) error {
	m.publishedArchivedEvents = append(m.publishedArchivedEvents, event)

	if m.publishArchivedFn != nil {
		return m.publishArchivedFn(ctx, event)
	}

	return nil
}

func (m *channelEventBusMock) SubscribeToChannelArchived(
	handler func(ctx context.Context, event channel_domain.ChannelArchived),
) error {
	return nil
}

// Compile-time interface checks
var _ channel_domain.ChannelRepository = (*channelRepositoryMock)(nil)
var _ channel_events.ChannelEventBus = (*channelEventBusMock)(nil)

// ------------------------------------------------------------------------------------------------------------
// HELPERS
// ------------------------------------------------------------------------------------------------------------

func validListChannelsRequest() *channel_application.ListChannelsRequest {
	return &channel_application.ListChannelsRequest{
		WorkspaceID: "workspace-001",
	}
}

// ------------------------------------------------------------------------------------------------------------
// TEST — Execute
// ------------------------------------------------------------------------------------------------------------

func TestListChannelUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	now := time.Now().UTC()

	expectedChannels := []channel_domain.Channel{
		{
			ID:          "channel-001",
			TemplateID:  "tpl-001",
			Title:       "Battery Engineering Review",
			Status:      channel_domain.StatusActive,
			WorkspaceID: "workspace-001",
			CreatedAt:   now,
			UpdatedAt:   now,
			IsDraft:     false,
			IsDirty:     false,
		},
		{
			ID:          "channel-002",
			TemplateID:  "tpl-002",
			Title:       "Flight Control Discussion",
			Status:      channel_domain.StatusActive,
			WorkspaceID: "workspace-001",
			CreatedAt:   now,
			UpdatedAt:   now,
			IsDraft:     true,
			IsDirty:     false,
		},
	}

	repo := &channelRepositoryMock{
		listFn: func(
			ctx context.Context,
			req *channel_domain.ListChannelsRequest,
		) (*tracecore_types.CloudResponse[[]channel_domain.Channel], error) {
			require.Equal(t, "workspace-001", req.WorkspaceID)

			return &tracecore_types.CloudResponse[[]channel_domain.Channel]{
				Data: expectedChannels,
			}, nil
		},
	}

	uc := channel_usecase.NewListChannelUsecase(repo)

	result, err := uc.Execute(ctx, validListChannelsRequest())

	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, expectedChannels, result)
}

func TestListChannelUsecase_Execute_EmptyCollection(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{
		listFn: func(
			ctx context.Context,
			req *channel_domain.ListChannelsRequest,
		) (*tracecore_types.CloudResponse[[]channel_domain.Channel], error) {
			return &tracecore_types.CloudResponse[[]channel_domain.Channel]{
				Data: []channel_domain.Channel{},
			}, nil
		},
	}

	uc := channel_usecase.NewListChannelUsecase(repo)

	result, err := uc.Execute(ctx, validListChannelsRequest())

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Empty(t, result)
}

func TestListChannelUsecase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("failed to list channels")

	repo := &channelRepositoryMock{
		listFn: func(
			ctx context.Context,
			req *channel_domain.ListChannelsRequest,
		) (*tracecore_types.CloudResponse[[]channel_domain.Channel], error) {
			return nil, expectedErr
		},
	}

	uc := channel_usecase.NewListChannelUsecase(repo)

	result, err := uc.Execute(ctx, validListChannelsRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)
}

func TestListChannelUsecase_Execute_NilRepository(t *testing.T) {
	uc := channel_usecase.NewListChannelUsecase(nil)

	ctx := context.Background()

	result, err := uc.Execute(ctx, validListChannelsRequest())

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, channel_domain.ErrRepositoryNil.Error())
}

func TestListChannelUsecase_Execute_NilRequest(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{}

	uc := channel_usecase.NewListChannelUsecase(repo)

	result, err := uc.Execute(ctx, nil)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, channel_domain.ErrRequestRequired.Error())
}

func TestListChannelUsecase_Execute_MissingWorkspaceID(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{}

	uc := channel_usecase.NewListChannelUsecase(repo)

	req := &channel_application.ListChannelsRequest{
		WorkspaceID: "",
	}

	result, err := uc.Execute(ctx, req)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, channel_domain.ErrWorkspaceIDRequired.Error())
}

// ------------------------------------------------------------------------------------------------------------
// TEST — ValidateRequest (table-driven)
// ------------------------------------------------------------------------------------------------------------

func TestListChannelUsecase_ValidateRequest(t *testing.T) {
	repo := &channelRepositoryMock{}

	uc := channel_usecase.NewListChannelUsecase(repo)

	tests := []struct {
		name          string
		request       *channel_application.ListChannelsRequest
		expectedError string
	}{
		{
			name:          "nil request",
			request:       nil,
			expectedError: channel_domain.ErrRequestRequired.Error(),
		},
		{
			name: "missing workspace id",
			request: &channel_application.ListChannelsRequest{
				WorkspaceID: "",
			},
			expectedError: channel_domain.ErrWorkspaceIDRequired.Error(),
		},
		{
			name: "valid request",
			request: &channel_application.ListChannelsRequest{
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

func TestListChannelUsecase_ValidateDependencies(t *testing.T) {
	tests := []struct {
		name          string
		repo          channel_domain.ChannelRepository
		expectedError string
	}{
		{
			name:          "nil repository",
			repo:          nil,
			expectedError: channel_domain.ErrRepositoryNil.Error(),
		},
		{
			name: "valid repository",
			repo: &channelRepositoryMock{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := channel_usecase.NewListChannelUsecase(tt.repo)

			err := uc.ValidateDependencies()

			if tt.expectedError == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.expectedError)
		})
	}
}

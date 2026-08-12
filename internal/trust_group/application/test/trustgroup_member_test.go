package trustgroup_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_events "vault-app/internal/trust_group/application/events"
	trustgroup_member_usecases "vault-app/internal/trust_group/application/usecases/member"
	trustgroup_domain "vault-app/internal/trust_group/domain"
)

// --------------------------------------------------------------------------------------------------
// MOCKS
// --------------------------------------------------------------------------------------------------
func (m *trustGroupRepositoryMock) AddMemberToTrustGroup(
	ctx context.Context,
	req *trustgroup_domain.AddMemberToTrustGroupRequest,
) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	if m.addMemberFn != nil {
		return m.addMemberFn(ctx, req)
	}

	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{
		Data: trustgroup_domain.TrustGroup{
			ID: "trust-group-001",
		},
	}, nil
}

func validAddMemberRequest() trustgroup_dtos.AddMemberToTrustGroupRequest {
	return trustgroup_dtos.AddMemberToTrustGroupRequest{
		TrustGroupID: "trust-group-001",
		ChannelID:    "channel-001",
		MemberID:     "member-001",
	}
}

// --------------------------------------------------------------------------------------------------
// TESTS
// --------------------------------------------------------------------------------------------------
func TestAddMemberToTrustGroupUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	var receivedRequest *trustgroup_domain.AddMemberToTrustGroupRequest

	expectedGroup := trustgroup_domain.TrustGroup{
		ID:         "trust-group-001",
		ChannelID:  "workspace-001",
		Name:       "OEM Trust Group",
		KEKVersion: 1,
		MemberCIDs: []string{"cid-member-001"},
		IsDraft:    true,
		IsDirty:    true,
	}

	repo := &trustGroupRepositoryMock{
		addMemberFn: func(
			ctx context.Context,
			req *trustgroup_domain.AddMemberToTrustGroupRequest,
		) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
			receivedRequest = req

			return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{
				Data: expectedGroup,
			}, nil
		},
	}

	eventBus := &trustGroupEventBusMock{}

	uc := trustgroup_member_usecases.NewAddMemberToTrustGroupUsecase(
		repo,
		eventBus,
	)

	req := validAddMemberRequest()

	result, err := uc.Execute(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, result)

	require.Equal(t, expectedGroup, *result)

	require.NotNil(t, receivedRequest)
	require.Equal(t, req.TrustGroupID, receivedRequest.TrustGroupID)
	require.Equal(t, req.MemberID, receivedRequest.MemberID)

	require.Len(t, eventBus.publishedMemberAddedEvents, 1)

	event := eventBus.publishedMemberAddedEvents[0]

	require.NotEmpty(t, event.EventID)
	require.False(t, event.EventTimestamp.IsZero())
	require.Equal(t, req.TrustGroupID, event.TrustGroupID)
	require.Equal(t, req.ChannelID, event.ChannelID)
	require.Equal(t, req.MemberID, event.MemberID)
}

func TestAddMemberToTrustGroupUsecase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("failed to add member to trust group")

	repo := &trustGroupRepositoryMock{
		addMemberFn: func(
			ctx context.Context,
			req *trustgroup_domain.AddMemberToTrustGroupRequest,
		) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
			return nil, expectedErr
		},
	}

	eventBus := &trustGroupEventBusMock{}

	uc := trustgroup_member_usecases.NewAddMemberToTrustGroupUsecase(
		repo,
		eventBus,
	)

	result, err := uc.Execute(ctx, validAddMemberRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)

	require.Empty(t, eventBus.publishedMemberAddedEvents)
}

func TestAddMemberToTrustGroupUsecase_Execute_EventBusError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("failed to publish member added event")

	repo := &trustGroupRepositoryMock{}

	eventBus := &trustGroupEventBusMock{
		publishMemberAddedFn: func(
			ctx context.Context,
			event trustgroup_domain.MemberAddedToTrustGroup,
		) error {
			return expectedErr
		},
	}

	uc := trustgroup_member_usecases.NewAddMemberToTrustGroupUsecase(
		repo,
		eventBus,
	)

	result, err := uc.Execute(ctx, validAddMemberRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)

	require.Len(t, eventBus.publishedMemberAddedEvents, 1)
}

func TestAddMemberToTrustGroupUsecase_Execute_NilRepository(t *testing.T) {
	ctx := context.Background()

	uc := trustgroup_member_usecases.NewAddMemberToTrustGroupUsecase(
		nil,
		&trustGroupEventBusMock{},
	)

	result, err := uc.Execute(ctx, validAddMemberRequest())

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, "trust group repository is required")
}

func TestAddMemberToTrustGroupUsecase_Execute_NilEventBus(t *testing.T) {
	ctx := context.Background()

	uc := trustgroup_member_usecases.NewAddMemberToTrustGroupUsecase(
		&trustGroupRepositoryMock{},
		nil,
	)

	result, err := uc.Execute(ctx, validAddMemberRequest())

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, "trust group event bus is required")
}

func TestAddMemberToTrustGroupUsecase_Execute_InvalidRequest(t *testing.T) {
	tests := []struct {
		name          string
		modify        func(req *trustgroup_dtos.AddMemberToTrustGroupRequest)
		expectedError string
	}{
		{
			name: "missing trust group id",
			modify: func(req *trustgroup_dtos.AddMemberToTrustGroupRequest) {
				req.TrustGroupID = "   "
			},
			expectedError: "trust group id is required",
		},
		{
			name: "missing channel id",
			modify: func(req *trustgroup_dtos.AddMemberToTrustGroupRequest) {
				req.ChannelID = "   "
			},
			expectedError: "channel id is required",
		},
		{
			name: "missing member id",
			modify: func(req *trustgroup_dtos.AddMemberToTrustGroupRequest) {
				req.MemberID = "   "
			},
			expectedError: "member id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			repo := &trustGroupRepositoryMock{}
			eventBus := &trustGroupEventBusMock{}

			uc := trustgroup_member_usecases.NewAddMemberToTrustGroupUsecase(
				repo,
				eventBus,
			)

			req := validAddMemberRequest()
			tt.modify(&req)

			result, err := uc.Execute(ctx, req)

			require.Error(t, err)
			require.Nil(t, result)
			require.EqualError(t, err, tt.expectedError)

			require.Empty(t, eventBus.publishedMemberAddedEvents)
		})
	}
}
func TestAddMemberToTrustGroupUsecase_ValidateDependencies(t *testing.T) {
	tests := []struct {
		name          string
		repo          trustgroup_domain.TrustGroupRepository
		eventBus      trustgroup_events.TrustGroupEventBus
		expectedError string
	}{
		{
			name:          "nil repository",
			repo:          nil,
			eventBus:      &trustGroupEventBusMock{},
			expectedError: "trust group repository is required",
		},
		{
			name:          "nil event bus",
			repo:          &trustGroupRepositoryMock{},
			eventBus:      nil,
			expectedError: "trust group event bus is required",
		},
		{
			name:     "valid dependencies",
			repo:     &trustGroupRepositoryMock{},
			eventBus: &trustGroupEventBusMock{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := trustgroup_member_usecases.NewAddMemberToTrustGroupUsecase(
				tt.repo,
				tt.eventBus,
			)

			err := uc.ValidateDependencies()

			if tt.expectedError == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.expectedError)
		})
	}
}

func TestAddMemberToTrustGroupUsecase_ValidateRequest(t *testing.T) {
	repo := &trustGroupRepositoryMock{}
	eventBus := &trustGroupEventBusMock{}

	uc := trustgroup_member_usecases.NewAddMemberToTrustGroupUsecase(
		repo,
		eventBus,
	)

	tests := []struct {
		name          string
		request       trustgroup_dtos.AddMemberToTrustGroupRequest
		expectedError string
	}{
		{
			name: "missing trust group id",
			request: trustgroup_dtos.AddMemberToTrustGroupRequest{
				TrustGroupID: "",
				ChannelID:    "channel-001",
				MemberID:     "member-001",
			},
			expectedError: "trust group id is required",
		},
		{
			name: "missing channel id",
			request: trustgroup_dtos.AddMemberToTrustGroupRequest{
				TrustGroupID: "trust-group-001",
				ChannelID:    "",
				MemberID:     "member-001",
			},
			expectedError: "channel id is required",
		},
		{
			name: "missing member id",
			request: trustgroup_dtos.AddMemberToTrustGroupRequest{
				TrustGroupID: "trust-group-001",
				ChannelID:    "channel-001",
				MemberID:     "",
			},
			expectedError: "member id is required",
		},
		{
			name:    "valid request",
			request: validAddMemberRequest(),
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

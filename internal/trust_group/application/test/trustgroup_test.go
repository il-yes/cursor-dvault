package trustgroup_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_events "vault-app/internal/trust_group/application/events"
	trustgroup_usecases "vault-app/internal/trust_group/application/usecases/trust_group"
	trustgroup_domain "vault-app/internal/trust_group/domain"
)

// --------------------------------------------------------------------------------------------------
// MOCKS
// --------------------------------------------------------------------------------------------------
type trustGroupRepositoryMock struct {
	createFn func(ctx context.Context, req *trustgroup_domain.CreateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error)
	addMemberFn func(
		ctx context.Context,
		req *trustgroup_domain.AddMemberToTrustGroupRequest,
	) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error)
}

func (m *trustGroupRepositoryMock) CreateTrustGroup(ctx context.Context, req *trustgroup_domain.CreateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}

	return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{
		Data: req.TrustGroup,
	}, nil
}

func (m *trustGroupRepositoryMock) GetTrustGroupMember(ctx context.Context, req *trustgroup_domain.GetTrustGroupMemberRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroupMember], error) {
	return nil, nil
}

func (m *trustGroupRepositoryMock) ListTrustGroups(ctx context.Context, req *trustgroup_domain.ListTrustGroupsRequest) (*tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

func (m *trustGroupRepositoryMock) GetTrustGroup(ctx context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

func (m *trustGroupRepositoryMock) UpdateTrustGroup(ctx context.Context, req *trustgroup_domain.UpdateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

func (m *trustGroupRepositoryMock) DeleteTrustGroup(ctx context.Context, req *trustgroup_domain.DeleteTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}


func (m *trustGroupRepositoryMock) RemoveMemberFromTrustGroup(ctx context.Context, req *trustgroup_domain.RemoveMemberFromTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

func (m *trustGroupRepositoryMock) RotateTrustGroupKEK(ctx context.Context, req *trustgroup_domain.RotateTrustGroupKEKRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

type trustGroupEventBusMock struct {
	publishedCreatedEvents              []trustgroup_domain.TrustGroupCreated
	publishedMemberAddedEvents          []trustgroup_domain.MemberAddedToTrustGroup
	publishedMemberRemovedEvents        []trustgroup_domain.MemberRemovedFromTrustGroup
	publishedTrustGroupDeletedEvents    []trustgroup_domain.TrustGroupDeleted
	publishedTrustGroupKEKRotatedEvents []trustgroup_domain.TrustGroupKEKRotated
	publishedTrustGroupRenamedEvents    []trustgroup_domain.TrustGroupRenamed

	publishCreatedFn              func(ctx context.Context, event trustgroup_domain.TrustGroupCreated) error
	publishMemberAddedFn          func(ctx context.Context, event trustgroup_domain.MemberAddedToTrustGroup) error
	publishMemberRemovedFn        func(ctx context.Context, event trustgroup_domain.MemberRemovedFromTrustGroup) error
	publishTrustGroupDeletedFn    func(ctx context.Context, event trustgroup_domain.TrustGroupDeleted) error
	publishTrustGroupKEKRotatedFn func(ctx context.Context, event trustgroup_domain.TrustGroupKEKRotated) error
	publishTrustGroupRenamedFn    func(ctx context.Context, event trustgroup_domain.TrustGroupRenamed) error

	subscribeCreatedFn              func(context.Context, trustgroup_domain.TrustGroupCreated) error
	subscribeMemberAddedFn          func(context.Context, trustgroup_domain.MemberAddedToTrustGroup) error
	subscribeMemberRemovedFn        func(context.Context, trustgroup_domain.MemberRemovedFromTrustGroup) error
	subscribeTrustGroupDeletedFn    func(context.Context, trustgroup_domain.TrustGroupDeleted) error
	subscribeTrustGroupKEKRotatedFn func(context.Context, trustgroup_domain.TrustGroupKEKRotated) error
	subscribeTrustGroupRenamedFn    func(context.Context, trustgroup_domain.TrustGroupRenamed) error
}

func (m *trustGroupEventBusMock) PublishTrustGroupCreated(ctx context.Context, event trustgroup_domain.TrustGroupCreated) error {
	m.publishedCreatedEvents = append(m.publishedCreatedEvents, event)

	if m.publishCreatedFn != nil {
		return m.publishCreatedFn(ctx, event)
	}

	return nil
}

func (m *trustGroupEventBusMock) PublishMemberAddedToTrustGroup(ctx context.Context, event trustgroup_domain.MemberAddedToTrustGroup) error {
	m.publishedMemberAddedEvents = append(m.publishedMemberAddedEvents, event)
	if m.publishMemberAddedFn != nil {
		return m.publishMemberAddedFn(ctx, event)
	}

	return nil
}

func (m *trustGroupEventBusMock) PublishMemberRemovedFromTrustGroup(ctx context.Context, event trustgroup_domain.MemberRemovedFromTrustGroup) error {
	m.publishedMemberRemovedEvents = append(m.publishedMemberRemovedEvents, event)
	if m.publishMemberRemovedFn != nil {
		return m.publishMemberRemovedFn(ctx, event)
	}

	return nil
}

func (m *trustGroupEventBusMock) PublishTrustGroupDeleted(ctx context.Context, event trustgroup_domain.TrustGroupDeleted) error {
	m.publishedTrustGroupDeletedEvents = append(m.publishedTrustGroupDeletedEvents, event)
	if m.publishTrustGroupDeletedFn != nil {
		return m.publishTrustGroupDeletedFn(ctx, event)
	}

	return nil
}

func (m *trustGroupEventBusMock) PublishTrustGroupKEKRotated(ctx context.Context, event trustgroup_domain.TrustGroupKEKRotated) error {
	m.publishedTrustGroupKEKRotatedEvents = append(m.publishedTrustGroupKEKRotatedEvents, event)
	return m.publishTrustGroupKEKRotatedFn(ctx, event)
}

func (m *trustGroupEventBusMock) PublishTrustGroupRenamed(ctx context.Context, event trustgroup_domain.TrustGroupRenamed) error {
	m.publishedTrustGroupRenamedEvents = append(m.publishedTrustGroupRenamedEvents, event)
	return m.publishTrustGroupRenamedFn(ctx, event)
}

func (m *trustGroupEventBusMock) SubscribeToMemberAddedToTrustGroup(handler func(context.Context, trustgroup_domain.MemberAddedToTrustGroup) error) {
	//TODO implement this
}

func (m *trustGroupEventBusMock) SubscribeToMemberRemovedFromTrustGroup(handler func(context.Context, trustgroup_domain.MemberRemovedFromTrustGroup) error) {
	//TODO implement this
}

func (m *trustGroupEventBusMock) SubscribeToTrustGroupKEKRotated(handler func(context.Context, trustgroup_domain.TrustGroupKEKRotated) error) {
	//TODO implement this
}

func (m *trustGroupEventBusMock) SubscribeToTrustGroupCreated(handler func(context.Context, trustgroup_domain.TrustGroupCreated) error) {
	//TODO implement this
}

func (m *trustGroupEventBusMock) SubscribeToTrustGroupDeleted(handler func(context.Context, trustgroup_domain.TrustGroupDeleted) error) {
	//TODO implement this
}

func validCreateTrustGroupRequest() trustgroup_dtos.CreateTrustGroupRequest {
	return trustgroup_dtos.CreateTrustGroupRequest{
		ChannelID: "channel-001",
		VaultID:   "vault-owner-001",
		OwnerID:   "owner-001",
		Name:      "OEM Trust Group",
		MemberCIDs: []string{
			"cid-member-oem",
			"cid-member-supplier",
		},
	}
}

func (m *trustGroupEventBusMock) SubscribeToTrustGroupRenamed(handler func(context.Context, trustgroup_domain.TrustGroupRenamed) error) {
	//TODO implement this
}

// --------------------------------------------------------------------------------------------------
// TESTS
// --------------------------------------------------------------------------------------------------

func TestCreateTrustGroupUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	var receivedRequest *trustgroup_domain.CreateTrustGroupRequest

	repo := &trustGroupRepositoryMock{
		createFn: func(
			ctx context.Context,
			req *trustgroup_domain.CreateTrustGroupRequest,
		) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
			receivedRequest = req

			return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{
				Data: req.TrustGroup,
			}, nil
		},
	}

	eventBus := &trustGroupEventBusMock{}

	uc := trustgroup_usecases.NewCreateTrustGroupUsecase(
		repo,
		eventBus,
	)

	req := validCreateTrustGroupRequest()

	result, err := uc.Execute(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, result)

	require.NotEmpty(t, result.ID)
	require.Equal(t, req.ChannelID, result.ChannelID)
	require.Equal(t, req.Name, result.Name)
	require.Equal(t, uint64(1), result.KEKVersion)
	require.Equal(t, req.MemberCIDs, result.MemberCIDs)
	require.True(t, result.IsDraft)
	require.False(t, result.IsDirty)

	require.NotNil(t, receivedRequest)
	require.Equal(t, req.Name, receivedRequest.TrustGroup.Name)
	require.Equal(t, req.ChannelID, receivedRequest.TrustGroup.ChannelID)
	require.Equal(t, req.MemberCIDs, receivedRequest.TrustGroup.MemberCIDs)

	require.Len(t, eventBus.publishedCreatedEvents, 1)

	event := eventBus.publishedCreatedEvents[0]

	require.NotEmpty(t, event.EventID)
	require.False(t, event.EventTimestamp.IsZero())
	require.Equal(t, req.ChannelID, event.ChannelID)
	require.Equal(t, req.VaultID, event.VaultID)
	require.Equal(t, req.OwnerID, event.OwnerID)
	require.Equal(t, req.Name, event.Name)
}

func TestCreateTrustGroupUsecase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("failed to create trust group")

	repo := &trustGroupRepositoryMock{
		createFn: func(
			ctx context.Context,
			req *trustgroup_domain.CreateTrustGroupRequest,
		) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
			return nil, expectedErr
		},
	}

	eventBus := &trustGroupEventBusMock{}

	uc := trustgroup_usecases.NewCreateTrustGroupUsecase(
		repo,
		eventBus,
	)

	result, err := uc.Execute(
		ctx,
		validCreateTrustGroupRequest(),
	)

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)

	require.Empty(t, eventBus.publishedCreatedEvents)
}

func TestCreateTrustGroupUsecase_Execute_EventBusError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("failed to publish trust group event")

	repo := &trustGroupRepositoryMock{}

	eventBus := &trustGroupEventBusMock{
		publishCreatedFn: func(
			ctx context.Context,
			event trustgroup_domain.TrustGroupCreated,
		) error {
			return expectedErr
		},
	}

	uc := trustgroup_usecases.NewCreateTrustGroupUsecase(
		repo,
		eventBus,
	)

	result, err := uc.Execute(
		ctx,
		validCreateTrustGroupRequest(),
	)

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)

	require.Len(t, eventBus.publishedCreatedEvents, 1)
}

func TestCreateTrustGroupUsecase_Execute_NilRepository(t *testing.T) {
	ctx := context.Background()

	uc := trustgroup_usecases.NewCreateTrustGroupUsecase(
		nil,
		&trustGroupEventBusMock{},
	)

	result, err := uc.Execute(
		ctx,
		validCreateTrustGroupRequest(),
	)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(
		t,
		err,
		"trust group repository is required",
	)
}

func TestCreateTrustGroupUsecase_Execute_NilEventBus(t *testing.T) {
	ctx := context.Background()

	uc := trustgroup_usecases.NewCreateTrustGroupUsecase(
		&trustGroupRepositoryMock{},
		nil,
	)

	result, err := uc.Execute(
		ctx,
		validCreateTrustGroupRequest(),
	)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(
		t,
		err,
		"trust group event bus is required",
	)
}

func TestCreateTrustGroupUsecase_Execute_MissingChannelID(t *testing.T) {
	ctx := context.Background()

	repo := &trustGroupRepositoryMock{}
	eventBus := &trustGroupEventBusMock{}

	uc := trustgroup_usecases.NewCreateTrustGroupUsecase(
		repo,
		eventBus,
	)

	req := validCreateTrustGroupRequest()
	req.ChannelID = "   "

	result, err := uc.Execute(ctx, req)

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorIs(t, err, trustgroup_domain.ErrChannelIDRequired)

	require.Empty(t, eventBus.publishedCreatedEvents)
}

func TestCreateTrustGroupUsecase_Execute_MissingName(t *testing.T) {
	ctx := context.Background()

	repo := &trustGroupRepositoryMock{}
	eventBus := &trustGroupEventBusMock{}

	uc := trustgroup_usecases.NewCreateTrustGroupUsecase(
		repo,
		eventBus,
	)

	req := validCreateTrustGroupRequest()
	req.Name = "   "

	result, err := uc.Execute(ctx, req)

	require.Error(t, err)
	require.Nil(t, result)
	require.ErrorIs(t, err, trustgroup_domain.ErrTrustGroupNameRequired)

	require.Empty(t, eventBus.publishedCreatedEvents)
}

func TestCreateTrustGroupUsecase_ValidateDependencies(t *testing.T) {
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
			uc := trustgroup_usecases.NewCreateTrustGroupUsecase(
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

func TestCreateTrustGroupUsecase_ValidateRequest(t *testing.T) {
	repo := &trustGroupRepositoryMock{}
	eventBus := &trustGroupEventBusMock{}

	uc := trustgroup_usecases.NewCreateTrustGroupUsecase(
		repo,
		eventBus,
	)

	tests := []struct {
		name          string
		request       trustgroup_dtos.CreateTrustGroupRequest
		expectedError error
	}{
		{
			name: "missing channel id",
			request: trustgroup_dtos.CreateTrustGroupRequest{
				ChannelID: " ",
				Name:      "OEM Trust Group",
			},
			expectedError: trustgroup_domain.ErrChannelIDRequired,
		},
		{
			name: "missing name",
			request: trustgroup_dtos.CreateTrustGroupRequest{
				ChannelID: "channel-001",
				Name:      " ",
			},
			expectedError: trustgroup_domain.ErrTrustGroupNameRequired,
		},
		{
			name:    "valid request",
			request: validCreateTrustGroupRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.ValidateRequest(tt.request)

			if tt.expectedError == nil {
				require.NoError(t, err)
				return
			}

			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}

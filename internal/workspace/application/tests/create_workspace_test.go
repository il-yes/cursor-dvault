package workspace_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	tracecore_types "vault-app/internal/tracecore/types"
	workspace_application "vault-app/internal/workspace/application"
	workspace_events "vault-app/internal/workspace/application/events"
	workspace_usecase "vault-app/internal/workspace/application/usecases"
	workspace_domain "vault-app/internal/workspace/domain"
)

// ------------------------------------------------------------------------------------------------------------
// MOCKS
// ------------------------------------------------------------------------------------------------------------
type workspaceRepositoryMock struct {
	createFn func(
		ctx context.Context,
		req workspace_domain.CreateRequest,
	) (*tracecore_types.CloudResponse[workspace_domain.Workspace], error)

	updateFn func(
		ctx context.Context,
		req workspace_domain.UpdateRequest,
	) (*tracecore_types.CloudResponse[workspace_domain.Workspace], error)

	deleteFn func(
		ctx context.Context,
		req workspace_domain.DeleteRequest,
	) error

	getFn func(
		ctx context.Context,
		req workspace_domain.GetRequest,
	) (*workspace_domain.Workspace, error)

	listFn func(
		ctx context.Context,
		req workspace_domain.ListRequest,
	) ([]workspace_domain.Workspace, error)
}

func (m *workspaceRepositoryMock) CreateWorkspace(
	ctx context.Context,
	req workspace_domain.CreateRequest,
) (*tracecore_types.CloudResponse[workspace_domain.Workspace], error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}

	return &tracecore_types.CloudResponse[workspace_domain.Workspace]{
		Data: req.Workspace,
	}, nil
}

func (m *workspaceRepositoryMock) UpdateWorkspace(
	ctx context.Context,
	req workspace_domain.UpdateRequest,
) (*tracecore_types.CloudResponse[workspace_domain.Workspace], error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, req)
	}

	return &tracecore_types.CloudResponse[workspace_domain.Workspace]{
		Data: req.Workspace,
	}, nil
}

func (m *workspaceRepositoryMock) DeleteWorkspace(
	ctx context.Context,
	req workspace_domain.DeleteRequest,
) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, req)
	}

	return nil
}

func (m *workspaceRepositoryMock) GetWorkspace(
	ctx context.Context,
	req workspace_domain.GetRequest,
) (*workspace_domain.Workspace, error) {
	if m.getFn != nil {
		return m.getFn(ctx, req)
	}

	return nil, nil
}

func (m *workspaceRepositoryMock) ListWorkspace(
	ctx context.Context,
	req workspace_domain.ListRequest,
) ([]workspace_domain.Workspace, error) {
	if m.listFn != nil {
		return m.listFn(ctx, req)
	}

	return nil, nil
}

type workspaceEventBusMock struct {
	publishCreatedFn func(
		ctx context.Context,
		event workspace_domain.WorkspaceCreated,
	) error

	publishDeletedFn func(
		ctx context.Context,
		event workspace_domain.WorkspaceDeleted,
	) error
	publishRenamedFn func(
		ctx context.Context,
		event workspace_domain.WorkspaceRenamed,
	) error

	publishedRenamedEvents []workspace_domain.WorkspaceRenamed

	publishedCreatedEvents []workspace_domain.WorkspaceCreated
	publishedDeletedEvents []workspace_domain.WorkspaceDeleted
}

func (m *workspaceEventBusMock) PublishWorkspaceCreated(ctx context.Context, event workspace_domain.WorkspaceCreated) error {
	m.publishedCreatedEvents = append(m.publishedCreatedEvents, event)

	if m.publishCreatedFn != nil {
		return m.publishCreatedFn(ctx, event)
	}

	return nil
}
func (m *workspaceEventBusMock) SubscribeToWorkspaceCreated(handler func(ctx context.Context, event workspace_domain.WorkspaceCreated)) error {
	return nil
}

func (m *workspaceEventBusMock) SubscribeToWorkspaceRenamed(handler func(ctx context.Context, event workspace_domain.WorkspaceRenamed)) error {
	return nil
}
func (m *workspaceEventBusMock) PublishWorkspaceRenamed(
	ctx context.Context,
	event workspace_domain.WorkspaceRenamed,
) error {
	m.publishedRenamedEvents = append(
		m.publishedRenamedEvents,
		event,
	)

	if m.publishRenamedFn != nil {
		return m.publishRenamedFn(ctx, event)
	}

	return nil
}

func (m *workspaceEventBusMock) PublishWorkspaceDeleted(
	ctx context.Context,
	event workspace_domain.WorkspaceDeleted,
) error {
	m.publishedDeletedEvents = append(
		m.publishedDeletedEvents,
		event,
	)

	if m.publishDeletedFn != nil {
		return m.publishDeletedFn(ctx, event)
	}

	return nil
}
func (m *workspaceEventBusMock) SubscribeToWorkspaceDeleted(handler func(ctx context.Context, event workspace_domain.WorkspaceDeleted)) error {
	return nil
}

func validCreateWorkspaceRequest() *workspace_application.CreateWorkspaceRequest {
	return &workspace_application.CreateWorkspaceRequest{
		VaultID:     "vault-001",
		Name:        "EVTOL Development Program",
		Description: "Battery System Engineering",
		OwnerID:     "owner-001",
		Signature:   "valid-signature",
	}
}

func newCreateWorkspaceUsecase(
	repo workspace_domain.Repository,
	bus workspace_events.WorkspaceEventBus,
) *workspace_usecase.CreateWorkspaceUsecase {
	return workspace_usecase.NewCreateWorkspaceUsecase(repo, bus)
}

// ------------------------------------------------------------------------------------------------------------
// TEST
// ------------------------------------------------------------------------------------------------------------
func TestCreateWorkspaceUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	expectedWorkspace := workspace_domain.Workspace{
		ID:          "workspace-001",
		VaultID:     "vault-001",
		Name:        "EVTOL Development Program",
		Description: "Battery System Engineering",
		Status:      workspace_domain.WorkspaceActive,
		OwnerID:     "owner-001",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsDraft:     true,
		IsDirty:     false,
	}

	repo := &workspaceRepositoryMock{
		createFn: func(
			ctx context.Context,
			req workspace_domain.CreateRequest,
		) (*tracecore_types.CloudResponse[workspace_domain.Workspace], error) {
			return &tracecore_types.CloudResponse[workspace_domain.Workspace]{
				Data: expectedWorkspace,
			}, nil
		},
	}

	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewCreateWorkspaceUsecase(repo, bus)

	req := validCreateWorkspaceRequest()

	workspace, err := uc.Execute(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, workspace)

	require.NotEmpty(t, workspace.ID)
	require.Equal(t, req.VaultID, workspace.VaultID)
	require.Equal(t, req.Name, workspace.Name)
	require.Equal(t, req.Description, workspace.Description)
	require.Equal(t, workspace_domain.WorkspaceActive, workspace.Status)
	require.Equal(t, req.OwnerID, workspace.OwnerID)

	require.Len(t, bus.publishedCreatedEvents, 1)

	event := bus.publishedCreatedEvents[0]

	require.Equal(t, workspace.VaultID, event.VaultID)
	require.Equal(t, workspace.OwnerID, event.OwnerID)
	require.Equal(t, workspace.Name, event.Name)
}

func TestCreateWorkspaceUsecase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("repository failure")

	repo := &workspaceRepositoryMock{
		createFn: func(
			ctx context.Context,
			req workspace_domain.CreateRequest,
		) (*tracecore_types.CloudResponse[workspace_domain.Workspace], error) {
			return nil, expectedErr
		},
	}

	bus := &workspaceEventBusMock{}
	uc := newCreateWorkspaceUsecase(repo, bus)

	workspace, err := uc.Execute(ctx, validCreateWorkspaceRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, workspace)

	require.Empty(t, bus.publishedCreatedEvents)
}

func TestCreateWorkspaceUsecase_Execute_EventBusError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("event bus failure")

	repo := &workspaceRepositoryMock{}

	bus := &workspaceEventBusMock{
		publishCreatedFn: func(
			ctx context.Context,
			event workspace_domain.WorkspaceCreated,
		) error {
			return expectedErr
		},
	}

	uc := newCreateWorkspaceUsecase(repo, bus)

	workspace, err := uc.Execute(ctx, validCreateWorkspaceRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, workspace)

	require.Len(t, bus.publishedCreatedEvents, 1)
}

func TestCreateWorkspaceUsecase_Execute_NilRepository(t *testing.T) {
	ctx := context.Background()

	bus := &workspaceEventBusMock{}
	uc := newCreateWorkspaceUsecase(nil, bus)

	workspace, err := uc.Execute(ctx, validCreateWorkspaceRequest())

	require.Error(t, err)
	require.Nil(t, workspace)
	require.EqualError(t, err, workspace_domain.ErrRepositoryNil)

	require.Empty(t, bus.publishedCreatedEvents)
}

func TestCreateWorkspaceUsecase_Execute_NilEventBus(t *testing.T) {
	ctx := context.Background()

	repo := &workspaceRepositoryMock{}
	uc := newCreateWorkspaceUsecase(repo, nil)

	workspace, err := uc.Execute(ctx, validCreateWorkspaceRequest())

	require.Error(t, err)
	require.Nil(t, workspace)
	require.EqualError(t, err, workspace_domain.ErrWorkspaceBusRequired)
}

func TestCreateWorkspaceUsecase_Execute_NilRequest(t *testing.T) {
	ctx := context.Background()

	repo := &workspaceRepositoryMock{}
	bus := &workspaceEventBusMock{}

	uc := newCreateWorkspaceUsecase(repo, bus)

	workspace, err := uc.Execute(ctx, nil)

	require.Error(t, err)
	require.Nil(t, workspace)
	require.EqualError(t, err, workspace_domain.ErrRequestRequired)
}

func TestCreateWorkspaceUsecase_Execute_InvalidRequest(t *testing.T) {
	tests := []struct {
		name          string
		modify        func(req *workspace_application.CreateWorkspaceRequest)
		expectedError string
	}{
		{
			name: "missing vault id",
			modify: func(req *workspace_application.CreateWorkspaceRequest) {
				req.VaultID = ""
			},
			expectedError: workspace_domain.ErrVaultIDRequired,
		},
		{
			name: "missing owner id",
			modify: func(req *workspace_application.CreateWorkspaceRequest) {
				req.OwnerID = ""
			},
			expectedError: workspace_domain.ErrWorkspaceOwnerRequired,
		},
		{
			name: "missing workspace name",
			modify: func(req *workspace_application.CreateWorkspaceRequest) {
				req.Name = ""
			},
			expectedError: workspace_domain.ErrWorkspaceNameRequired,
		},
		{
			name: "missing signature",
			modify: func(req *workspace_application.CreateWorkspaceRequest) {
				req.Signature = ""
			},
			expectedError: workspace_domain.ErrSignatureMissing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			repo := &workspaceRepositoryMock{}
			bus := &workspaceEventBusMock{}
			uc := newCreateWorkspaceUsecase(repo, bus)

			req := validCreateWorkspaceRequest()
			tt.modify(req)

			workspace, err := uc.Execute(ctx, req)

			require.Error(t, err)
			require.Nil(t, workspace)
			require.EqualError(t, err, tt.expectedError)
			require.Empty(t, bus.publishedCreatedEvents)
		})
	}
}

func TestCreateWorkspaceUsecase_ValidateDependencies(t *testing.T) {
	tests := []struct {
		name          string
		repo          workspace_domain.Repository
		bus           workspace_events.WorkspaceEventBus
		expectedError string
	}{
		{
			name:          "nil repository",
			repo:          nil,
			bus:           &workspaceEventBusMock{},
			expectedError: workspace_domain.ErrRepositoryNil,
		},
		{
			name:          "nil event bus",
			repo:          &workspaceRepositoryMock{},
			bus:           nil,
			expectedError: workspace_domain.ErrWorkspaceBusRequired,
		},
		{
			name: "valid dependencies",
			repo: &workspaceRepositoryMock{},
			bus:  &workspaceEventBusMock{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := newCreateWorkspaceUsecase(tt.repo, tt.bus)

			err := uc.ValidateDependencies()

			if tt.expectedError == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.expectedError)
		})
	}
}

func TestCreateWorkspaceUsecase_ValidateRequest(t *testing.T) {
	repo := &workspaceRepositoryMock{}
	bus := &workspaceEventBusMock{}
	uc := newCreateWorkspaceUsecase(repo, bus)

	tests := []struct {
		name          string
		request       *workspace_application.CreateWorkspaceRequest
		expectedError string
	}{
		{
			name:          "nil request",
			request:       nil,
			expectedError: workspace_domain.ErrRequestRequired,
		},
		{
			name: "valid request",
			request: &workspace_application.CreateWorkspaceRequest{
				VaultID:   "vault-001",
				OwnerID:   "owner-001",
				Name:      "Engineering Workspace",
				Signature: "signature",
			},
		},
		{
			name: "missing vault id",
			request: &workspace_application.CreateWorkspaceRequest{
				OwnerID:   "owner-001",
				Name:      "Engineering Workspace",
				Signature: "signature",
			},
			expectedError: workspace_domain.ErrVaultIDRequired,
		},
		{
			name: "missing owner id",
			request: &workspace_application.CreateWorkspaceRequest{
				VaultID:   "vault-001",
				Name:      "Engineering Workspace",
				Signature: "signature",
			},
			expectedError: workspace_domain.ErrWorkspaceOwnerRequired,
		},
		{
			name: "missing name",
			request: &workspace_application.CreateWorkspaceRequest{
				VaultID:   "vault-001",
				OwnerID:   "owner-001",
				Signature: "signature",
			},
			expectedError: workspace_domain.ErrWorkspaceNameRequired,
		},
		{
			name: "missing signature",
			request: &workspace_application.CreateWorkspaceRequest{
				VaultID: "vault-001",
				OwnerID: "owner-001",
				Name:    "Engineering Workspace",
			},
			expectedError: workspace_domain.ErrSignatureMissing,
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

func NewWorkspace(
	vaultID string,
	name string,
	description string,
	ownerID string,
) workspace_domain.Workspace {
	now := time.Now()

	return workspace_domain.Workspace{
		ID:          uuid.NewString(),
		VaultID:     vaultID,
		Name:        name,
		Description: description,
		Status:      workspace_domain.WorkspaceActive,
		OwnerID:     ownerID,
		CreatedAt:   now,
		UpdatedAt:   now,
		IsDraft:     true,
		IsDirty:     false,
	}
}

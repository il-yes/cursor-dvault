package workspace_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	tracecore_types "vault-app/internal/tracecore/types"
	workspace_application "vault-app/internal/workspace/application"
	workspace_events "vault-app/internal/workspace/application/events"
	workspace_usecase "vault-app/internal/workspace/application/usecases"
	workspace_domain "vault-app/internal/workspace/domain"
)

func validRenameWorkspaceRequest() *workspace_application.RenameWorkspaceRequest {
	now := time.Now().UTC()

	return &workspace_application.RenameWorkspaceRequest{
		Workspace: &workspace_domain.Workspace{
			ID:          "workspace-001",
			VaultID:     "vault-001",
			Name:        "Old Workspace Name",
			Description: "Battery System Engineering",
			Status:      workspace_domain.WorkspaceActive,
			OwnerID:     "owner-001",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		Name:      "New Workspace Name",
		Signature: "valid-signature",
	}
}

func TestRenameWorkspaceUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	var receivedRequest workspace_domain.UpdateRequest

	repo := &workspaceRepositoryMock{
		updateFn: func(
			ctx context.Context,
			req workspace_domain.UpdateRequest,
		) (*tracecore_types.CloudResponse[workspace_domain.Workspace], error) {
			receivedRequest = req

			return &tracecore_types.CloudResponse[workspace_domain.Workspace]{
				Data: req.Workspace,
			}, nil
		},
	}

	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewRenameWorkspaceUsecase(repo, bus)

	req := validRenameWorkspaceRequest()
	oldName := req.Workspace.Name

	result, err := uc.Execute(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, result)

	require.Equal(t, "workspace-001", result.ID)
	require.Equal(t, "vault-001", result.VaultID)
	require.Equal(t, "New Workspace Name", result.Name)

	require.Equal(t, req.Workspace.OwnerID, receivedRequest.UserID)
	require.Equal(t, req.Workspace.VaultID, receivedRequest.VaultID)
	require.Equal(t, req.Workspace.ID, receivedRequest.Workspace.ID)
	require.Equal(t, req.Name, receivedRequest.Workspace.Name)
	require.Equal(t, req.Signature, receivedRequest.Signature)

	require.Len(t, bus.publishedRenamedEvents, 1)

	event := bus.publishedRenamedEvents[0]

	require.NotEmpty(t, event.EventID)
	require.False(t, event.EventTimestamp.IsZero())
	require.Equal(t, req.Workspace.ID, event.WorkspaceID)
	require.Equal(t, oldName, event.OldName)
	require.Equal(t, "New Workspace Name", event.NewName)
}

func TestRenameWorkspaceUsecase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("failed to update workspace")

	repo := &workspaceRepositoryMock{
		updateFn: func(
			ctx context.Context,
			req workspace_domain.UpdateRequest,
		) (*tracecore_types.CloudResponse[workspace_domain.Workspace], error) {
			return nil, expectedErr
		},
	}

	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewRenameWorkspaceUsecase(repo, bus)

	result, err := uc.Execute(ctx, validRenameWorkspaceRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)

	require.Empty(t, bus.publishedRenamedEvents)
}

func TestRenameWorkspaceUsecase_Execute_EventBusError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("failed to publish workspace renamed event")

	repo := &workspaceRepositoryMock{
		updateFn: func(
			ctx context.Context,
			req workspace_domain.UpdateRequest,
		) (*tracecore_types.CloudResponse[workspace_domain.Workspace], error) {
			return &tracecore_types.CloudResponse[workspace_domain.Workspace]{
				Data: req.Workspace,
			}, nil
		},
	}

	bus := &workspaceEventBusMock{
		publishRenamedFn: func(
			ctx context.Context,
			event workspace_domain.WorkspaceRenamed,
		) error {
			return expectedErr
		},
	}

	uc := workspace_usecase.NewRenameWorkspaceUsecase(repo, bus)

	result, err := uc.Execute(ctx, validRenameWorkspaceRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)

	require.Len(t, bus.publishedRenamedEvents, 1)
}

func TestRenameWorkspaceUsecase_Execute_NilRepository(t *testing.T) {
	ctx := context.Background()

	uc := workspace_usecase.NewRenameWorkspaceUsecase(
		nil,
		&workspaceEventBusMock{},
	)

	result, err := uc.Execute(ctx, validRenameWorkspaceRequest())

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, workspace_domain.ErrRepositoryNil)
}

func TestRenameWorkspaceUsecase_Execute_NilEventBus(t *testing.T) {
	ctx := context.Background()

	uc := workspace_usecase.NewRenameWorkspaceUsecase(
		&workspaceRepositoryMock{},
		nil,
	)

	result, err := uc.Execute(ctx, validRenameWorkspaceRequest())

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, workspace_domain.ErrWorkspaceBusRequired)
}

func TestRenameWorkspaceUsecase_Execute_NilRequest(t *testing.T) {
	ctx := context.Background()

	repo := &workspaceRepositoryMock{}
	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewRenameWorkspaceUsecase(repo, bus)

	result, err := uc.Execute(ctx, nil)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, workspace_domain.ErrRequestRequired)
}

func TestRenameWorkspaceUsecase_Execute_InvalidRequest(t *testing.T) {
	tests := []struct {
		name          string
		modify        func(req *workspace_application.RenameWorkspaceRequest)
		expectedError string
	}{
		{
			name: "missing name",
			modify: func(req *workspace_application.RenameWorkspaceRequest) {
				req.Name = ""
			},
			expectedError: workspace_domain.ErrWorkspaceNameRequired,
		},
		{
			name: "missing workspace",
			modify: func(req *workspace_application.RenameWorkspaceRequest) {
				req.Workspace = nil
			},
			expectedError: workspace_domain.ErrWorkspaceOwnerRequired,
		},
		{
			name: "missing signature",
			modify: func(req *workspace_application.RenameWorkspaceRequest) {
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

			uc := workspace_usecase.NewRenameWorkspaceUsecase(repo, bus)

			req := validRenameWorkspaceRequest()
			tt.modify(req)

			result, err := uc.Execute(ctx, req)

			require.Error(t, err)
			require.Nil(t, result)
			require.EqualError(t, err, tt.expectedError)
			require.Empty(t, bus.publishedRenamedEvents)
		})
	}
}

func TestRenameWorkspaceUsecase_ValidateRequest(t *testing.T) {
	repo := &workspaceRepositoryMock{}
	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewRenameWorkspaceUsecase(repo, bus)

	tests := []struct {
		name          string
		request       *workspace_application.RenameWorkspaceRequest
		expectedError string
	}{
		{
			name:          "nil request",
			request:       nil,
			expectedError: workspace_domain.ErrRequestRequired,
		},
		{
			name: "missing name",
			request: &workspace_application.RenameWorkspaceRequest{
				Workspace: validRenameWorkspaceRequest().Workspace,
				Signature: "signature",
			},
			expectedError: workspace_domain.ErrWorkspaceNameRequired,
		},
		{
			name: "missing workspace",
			request: &workspace_application.RenameWorkspaceRequest{
				Name:      "New Name",
				Signature: "signature",
			},
			expectedError: workspace_domain.ErrWorkspaceOwnerRequired,
		},
		{
			name: "missing signature",
			request: &workspace_application.RenameWorkspaceRequest{
				Workspace: validRenameWorkspaceRequest().Workspace,
				Name:      "New Name",
			},
			expectedError: workspace_domain.ErrSignatureMissing,
		},
		{
			name:    "valid request",
			request: validRenameWorkspaceRequest(),
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

func TestRenameWorkspaceUsecase_ValidateDependencies(t *testing.T) {
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
			uc := workspace_usecase.NewRenameWorkspaceUsecase(
				tt.repo,
				tt.bus,
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

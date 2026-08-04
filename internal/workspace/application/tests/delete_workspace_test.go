package workspace_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	workspace_application "vault-app/internal/workspace/application"
	workspace_events "vault-app/internal/workspace/application/events"
	workspace_usecase "vault-app/internal/workspace/application/usecases"
	workspace_domain "vault-app/internal/workspace/domain"
)

func validDeleteWorkspaceRequest() *workspace_application.DeleteWorkspaceRequest {
	return &workspace_application.DeleteWorkspaceRequest{
		VaultID:     "vault-001",
		WorkspaceID: "workspace-001",
		Signature:   "valid-signature",
	}
}

func TestDeleteWorkspaceUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	var receivedRequest workspace_domain.DeleteRequest

	repo := &workspaceRepositoryMock{
		deleteFn: func(
			ctx context.Context,
			req workspace_domain.DeleteRequest,
		) error {
			receivedRequest = req
			return nil
		},
	}

	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewDeleteWorkspaceUsecase(repo, bus)

	req := validDeleteWorkspaceRequest()

	err := uc.Execute(ctx, req)

	require.NoError(t, err)

	require.Equal(t, req.WorkspaceID, receivedRequest.WorkspaceID)
	require.Equal(t, req.Signature, receivedRequest.Signature)

	require.Len(t, bus.publishedDeletedEvents, 1)

	event := bus.publishedDeletedEvents[0]

	require.NotEmpty(t, event.EventID)
	require.False(t, event.EventTimestamp.IsZero())
	require.Equal(t, req.WorkspaceID, event.WorkspaceID)
	require.Equal(t, req.VaultID, event.VaultID)
}

func TestDeleteWorkspaceUsecase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("failed to delete workspace")

	repo := &workspaceRepositoryMock{
		deleteFn: func(
			ctx context.Context,
			req workspace_domain.DeleteRequest,
		) error {
			return expectedErr
		},
	}

	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewDeleteWorkspaceUsecase(repo, bus)

	err := uc.Execute(ctx, validDeleteWorkspaceRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)

	// The deleted event must not be published when deletion fails.
	require.Empty(t, bus.publishedDeletedEvents)
}

func TestDeleteWorkspaceUsecase_Execute_EventBusError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("failed to publish workspace deleted event")

	repo := &workspaceRepositoryMock{
		deleteFn: func(
			ctx context.Context,
			req workspace_domain.DeleteRequest,
		) error {
			return nil
		},
	}

	bus := &workspaceEventBusMock{
		publishDeletedFn: func(
			ctx context.Context,
			event workspace_domain.WorkspaceDeleted,
		) error {
			return expectedErr
		},
	}

	uc := workspace_usecase.NewDeleteWorkspaceUsecase(repo, bus)

	err := uc.Execute(ctx, validDeleteWorkspaceRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)

	require.Len(t, bus.publishedDeletedEvents, 1)
}

func TestDeleteWorkspaceUsecase_Execute_NilRepository(t *testing.T) {
	ctx := context.Background()

	uc := workspace_usecase.NewDeleteWorkspaceUsecase(
		nil,
		&workspaceEventBusMock{},
	)

	err := uc.Execute(ctx, validDeleteWorkspaceRequest())

	require.Error(t, err)
	require.EqualError(t, err, workspace_domain.ErrRepositoryNil)
}

func TestDeleteWorkspaceUsecase_Execute_NilEventBus(t *testing.T) {
	ctx := context.Background()

	uc := workspace_usecase.NewDeleteWorkspaceUsecase(
		&workspaceRepositoryMock{},
		nil,
	)

	err := uc.Execute(ctx, validDeleteWorkspaceRequest())

	require.Error(t, err)
	require.EqualError(t, err, workspace_domain.ErrWorkspaceBusRequired)
}

func TestDeleteWorkspaceUsecase_Execute_NilRequest(t *testing.T) {
	ctx := context.Background()

	repo := &workspaceRepositoryMock{}
	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewDeleteWorkspaceUsecase(repo, bus)

	err := uc.Execute(ctx, nil)

	require.Error(t, err)
	require.EqualError(t, err, workspace_domain.ErrRequestRequired)
}

func TestDeleteWorkspaceUsecase_Execute_InvalidRequest(t *testing.T) {
	tests := []struct {
		name          string
		modify        func(req *workspace_application.DeleteWorkspaceRequest)
		expectedError string
	}{
		{
			name: "missing vault id",
			modify: func(req *workspace_application.DeleteWorkspaceRequest) {
				req.VaultID = ""
			},
			expectedError: workspace_domain.ErrVaultIDRequired,
		},
		{
			name: "missing workspace id",
			modify: func(req *workspace_application.DeleteWorkspaceRequest) {
				req.WorkspaceID = ""
			},
			expectedError: workspace_domain.ErrWorkspaceIDRequired,
		},
		{
			name: "missing signature",
			modify: func(req *workspace_application.DeleteWorkspaceRequest) {
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

			uc := workspace_usecase.NewDeleteWorkspaceUsecase(repo, bus)

			req := validDeleteWorkspaceRequest()
			tt.modify(req)

			err := uc.Execute(ctx, req)

			require.Error(t, err)
			require.EqualError(t, err, tt.expectedError)
			require.Empty(t, bus.publishedDeletedEvents)
		})
	}
}

func TestDeleteWorkspaceUsecase_ValidateRequest(t *testing.T) {
	repo := &workspaceRepositoryMock{}
	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewDeleteWorkspaceUsecase(repo, bus)

	tests := []struct {
		name          string
		request       *workspace_application.DeleteWorkspaceRequest
		expectedError string
	}{
		{
			name:          "nil request",
			request:       nil,
			expectedError: workspace_domain.ErrRequestRequired,
		},
		{
			name: "missing vault id",
			request: &workspace_application.DeleteWorkspaceRequest{
				WorkspaceID: "workspace-001",
				Signature:   "signature",
			},
			expectedError: workspace_domain.ErrVaultIDRequired,
		},
		{
			name: "missing workspace id",
			request: &workspace_application.DeleteWorkspaceRequest{
				VaultID:   "vault-001",
				Signature: "signature",
			},
			expectedError: workspace_domain.ErrWorkspaceIDRequired,
		},
		{
			name: "missing signature",
			request: &workspace_application.DeleteWorkspaceRequest{
				VaultID:     "vault-001",
				WorkspaceID: "workspace-001",
			},
			expectedError: workspace_domain.ErrSignatureMissing,
		},
		{
			name: "valid request",
			request: &workspace_application.DeleteWorkspaceRequest{
				VaultID:     "vault-001",
				WorkspaceID: "workspace-001",
				Signature:   "signature",
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

func TestDeleteWorkspaceUsecase_ValidateDependencies(t *testing.T) {
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
			uc := workspace_usecase.NewDeleteWorkspaceUsecase(
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

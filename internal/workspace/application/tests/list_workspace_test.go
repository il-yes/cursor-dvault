package workspace_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	workspace_application "vault-app/internal/workspace/application"
	workspace_domain "vault-app/internal/workspace/domain"
	workspace_usecase "vault-app/internal/workspace/application/usecases"
)

func validListWorkspaceRequest() *workspace_application.ListWorkspacesRequest {
	return &workspace_application.ListWorkspacesRequest{
		VaultID: "vault-001",
	}
}

func TestListWorkspaceUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	now := time.Now().UTC()

	expectedWorkspaces := []workspace_domain.Workspace{
		{
			ID:          "workspace-001",
			VaultID:     "vault-001",
			Name:        "EVTOL Development Program",
			Description: "Battery System Engineering",
			Status:      workspace_domain.WorkspaceActive,
			OwnerID:     "owner-001",
			CreatedAt:   now,
			UpdatedAt:   now,
			IsDraft:     false,
			IsDirty:     false,
		},
		{
			ID:          "workspace-002",
			VaultID:     "vault-001",
			Name:        "Flight Control Program",
			Description: "Flight control engineering",
			Status:      workspace_domain.WorkspaceArchived,
			OwnerID:     "owner-001",
			CreatedAt:   now,
			UpdatedAt:   now,
			IsDraft:     false,
			IsDirty:     false,
		},
	}

	repo := &workspaceRepositoryMock{
		listFn: func(
			ctx context.Context,
			req workspace_domain.ListRequest,
		) ([]workspace_domain.Workspace, error) {
			require.Equal(t, "vault-001", req.VaultID)

			return expectedWorkspaces, nil
		},
	}

	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewListWorkspaceUsecase(repo, bus)

	result, err := uc.Execute(ctx, validListWorkspaceRequest())

	require.NoError(t, err)
	require.Len(t, result, 2)

	require.Equal(t, expectedWorkspaces, result)
	require.Empty(t, bus.publishedCreatedEvents)
}

func TestListWorkspaceUsecase_Execute_EmptyCollection(t *testing.T) {
	ctx := context.Background()

	repo := &workspaceRepositoryMock{
		listFn: func(
			ctx context.Context,
			req workspace_domain.ListRequest,
		) ([]workspace_domain.Workspace, error) {
			return []workspace_domain.Workspace{}, nil
		},
	}

	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewListWorkspaceUsecase(repo, bus)

	result, err := uc.Execute(ctx, validListWorkspaceRequest())

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Empty(t, result)
}

func TestListWorkspaceUsecase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("failed to list workspaces")

	repo := &workspaceRepositoryMock{
		listFn: func(
			ctx context.Context,
			req workspace_domain.ListRequest,
		) ([]workspace_domain.Workspace, error) {
			return nil, expectedErr
		},
	}

	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewListWorkspaceUsecase(repo, bus)

	result, err := uc.Execute(ctx, validListWorkspaceRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)
}

func TestListWorkspaceUsecase_Execute_NilRepository(t *testing.T) {
	ctx := context.Background()

	uc := workspace_usecase.NewListWorkspaceUsecase(
		nil,
		&workspaceEventBusMock{},
	)

	result, err := uc.Execute(ctx, validListWorkspaceRequest())

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, workspace_domain.ErrRepositoryNil)
}


func TestListWorkspaceUsecase_Execute_NilRequest(t *testing.T) {
	ctx := context.Background()

	repo := &workspaceRepositoryMock{}
	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewListWorkspaceUsecase(repo, bus)

	result, err := uc.Execute(ctx, nil)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, workspace_domain.ErrRequestRequired)
}

func TestListWorkspaceUsecase_Execute_MissingVaultID(t *testing.T) {
	ctx := context.Background()

	repo := &workspaceRepositoryMock{}
	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewListWorkspaceUsecase(repo, bus)

	req := &workspace_application.ListWorkspacesRequest{
		VaultID: "",
	}

	result, err := uc.Execute(ctx, req)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, workspace_domain.ErrVaultIDRequired)
}

func TestListWorkspaceUsecase_ValidateRequest(t *testing.T) {
	repo := &workspaceRepositoryMock{}
	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewListWorkspaceUsecase(repo, bus)

	tests := []struct {
		name          string
		request       *workspace_application.ListWorkspacesRequest
		expectedError string
	}{
		{
			name:          "nil request",
			request:       nil,
			expectedError: workspace_domain.ErrRequestRequired,
		},
		{
			name: "missing vault id",
			request: &workspace_application.ListWorkspacesRequest{
				VaultID: "",
			},
			expectedError: workspace_domain.ErrVaultIDRequired,
		},
		{
			name: "valid request",
			request: &workspace_application.ListWorkspacesRequest{
				VaultID: "vault-001",
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

func TestListWorkspaceUsecase_ValidateDependencies(t *testing.T) {
	bus := &workspaceEventBusMock{}

	tests := []struct {
		name          string
		repo          workspace_domain.Repository
		expectedError string
	}{
		{
			name:          "nil repository",
			repo:          nil,
			expectedError: workspace_domain.ErrRepositoryNil,
		},
		{
			name: "valid repository",
			repo: &workspaceRepositoryMock{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := workspace_usecase.NewListWorkspaceUsecase(
				tt.repo,
				bus,
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


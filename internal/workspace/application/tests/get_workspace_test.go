package workspace_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	workspace_application "vault-app/internal/workspace/application"
	workspace_usecase "vault-app/internal/workspace/application/usecases"
	workspace_domain "vault-app/internal/workspace/domain"
)

// ------------------------------------------------------------------------------------------------------------
// MOCKS
// ------------------------------------------------------------------------------------------------------------
func validGetWorkspaceRequest() *workspace_application.GetWorkspaceRequest {
	return &workspace_application.GetWorkspaceRequest{
		WorkspaceID: "workspace-001",
	}
}

func expectedWorkspace() *workspace_domain.Workspace {
	now := time.Now().UTC()

	return &workspace_domain.Workspace{
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
	}
}

// ------------------------------------------------------------------------------------------------------------
// TEST
// ------------------------------------------------------------------------------------------------------------
func TestGetWorkspaceUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	expected := expectedWorkspace()

	repo := &workspaceRepositoryMock{
		getFn: func(
			ctx context.Context,
			req workspace_domain.GetRequest,
		) (*workspace_domain.Workspace, error) {
			require.Equal(t, "workspace-001", req.WorkspaceID)

			return expected, nil
		},
	}

	// The Get use case does not use DomainBus, but the constructor accepts it.
	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewGetWorkspaceUsecase(repo, bus)

	result, err := uc.Execute(ctx, validGetWorkspaceRequest())

	require.NoError(t, err)
	require.NotNil(t, result)

	require.Equal(t, expected.ID, result.ID)
	require.Equal(t, expected.VaultID, result.VaultID)
	require.Equal(t, expected.Name, result.Name)
	require.Equal(t, expected.Description, result.Description)
	require.Equal(t, expected.Status, result.Status)
	require.Equal(t, expected.OwnerID, result.OwnerID)
	require.Equal(t, expected.CreatedAt, result.CreatedAt)
	require.Equal(t, expected.UpdatedAt, result.UpdatedAt)

	require.Empty(t, bus.publishedCreatedEvents)
}

func TestGetWorkspaceUsecase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("failed to get workspace")

	repo := &workspaceRepositoryMock{
		getFn: func(
			ctx context.Context,
			req workspace_domain.GetRequest,
		) (*workspace_domain.Workspace, error) {
			return nil, expectedErr
		},
	}

	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewGetWorkspaceUsecase(repo, bus)

	result, err := uc.Execute(ctx, validGetWorkspaceRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)

	require.Empty(t, bus.publishedCreatedEvents)
}
func TestGetWorkspaceUsecase_Execute_NilRepository(t *testing.T) {
	ctx := context.Background()

	uc := workspace_usecase.NewGetWorkspaceUsecase(
		nil,
		&workspaceEventBusMock{},
	)

	result, err := uc.Execute(ctx, validGetWorkspaceRequest())

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, workspace_domain.ErrRepositoryNil)
}

func TestGetWorkspaceUsecase_Execute_NilRequest(t *testing.T) {
	ctx := context.Background()

	repo := &workspaceRepositoryMock{}
	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewGetWorkspaceUsecase(repo, bus)

	result, err := uc.Execute(ctx, nil)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, workspace_domain.ErrRequestRequired)
}

func TestGetWorkspaceUsecase_Execute_MissingWorkspaceID(t *testing.T) {
	ctx := context.Background()

	repo := &workspaceRepositoryMock{}
	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewGetWorkspaceUsecase(repo, bus)

	req := &workspace_application.GetWorkspaceRequest{
		WorkspaceID: "",
	}

	result, err := uc.Execute(ctx, req)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, workspace_domain.ErrWorkspaceIDRequired)
}

func TestGetWorkspaceUsecase_ValidateRequest(t *testing.T) {
	repo := &workspaceRepositoryMock{}
	bus := &workspaceEventBusMock{}

	uc := workspace_usecase.NewGetWorkspaceUsecase(repo, bus)

	tests := []struct {
		name          string
		request       *workspace_application.GetWorkspaceRequest
		expectedError string
	}{
		{
			name:          "nil request",
			request:       nil,
			expectedError: workspace_domain.ErrRequestRequired,
		},
		{
			name: "missing workspace id",
			request: &workspace_application.GetWorkspaceRequest{
				WorkspaceID: "",
			},
			expectedError: workspace_domain.ErrWorkspaceIDRequired,
		},
		{
			name: "valid request",
			request: &workspace_application.GetWorkspaceRequest{
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

func TestGetWorkspaceUsecase_ValidateDependencies(t *testing.T) {
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
			uc := workspace_usecase.NewGetWorkspaceUsecase(
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

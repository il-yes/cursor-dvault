package workspace_usecase

import (
	"context"
	"errors"

	"vault-app/internal/utils"
	workspace_application "vault-app/internal/workspace/application"
	workspace_events "vault-app/internal/workspace/application/events"
	workspace_domain "vault-app/internal/workspace/domain"
)

type ListWorkspaceUsecase struct {
	Repo      workspace_domain.Repository
	DomainBus workspace_events.WorkspaceEventBus
}

func NewListWorkspaceUsecase(repo workspace_domain.Repository, workspaceBus workspace_events.WorkspaceEventBus) *ListWorkspaceUsecase {
	return &ListWorkspaceUsecase{
		Repo:      repo,
		DomainBus: workspaceBus,
	}
}

func (c *ListWorkspaceUsecase) Execute(ctx context.Context, req *workspace_application.ListWorkspacesRequest) ([]workspace_domain.Workspace, error) {
	if err := c.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := c.ValidateRequest(req); err != nil {
		return nil, err
	}

	collection, err := c.Repo.ListWorkspace(ctx, workspace_domain.ListRequest{
		VaultID:   req.VaultID,
	})
	if err != nil {
		utils.LogPretty("🚫 [Workspace] ListWorkspaceUsecase.Execute error", err)
		return nil, err
	}

	utils.LogPretty("[Workspace] ListWorkspaceUsecase.Execute result", collection)
	return collection, nil
}

func (c *ListWorkspaceUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return errors.New(workspace_domain.ErrRepositoryNil)
	}


	return nil
}

func (c *ListWorkspaceUsecase) ValidateRequest(req *workspace_application.ListWorkspacesRequest) error {
	if req == nil {
		return errors.New(workspace_domain.ErrRequestRequired)
	}

	if req.VaultID == "" {
		return errors.New(workspace_domain.ErrVaultIDRequired)
	}

	return nil
}

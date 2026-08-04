package workspace_usecase

import (
	"context"
	"errors"

	workspace_application "vault-app/internal/workspace/application"
	workspace_events "vault-app/internal/workspace/application/events"
	workspace_domain "vault-app/internal/workspace/domain"
)

type GetWorkspaceUsecase struct {
	Repo      workspace_domain.Repository
	DomainBus workspace_events.WorkspaceEventBus
}

func NewGetWorkspaceUsecase(repo workspace_domain.Repository, workspaceBus workspace_events.WorkspaceEventBus) *GetWorkspaceUsecase {
	return &GetWorkspaceUsecase{
		Repo:      repo,
		DomainBus: workspaceBus,
	}
}

func (c *GetWorkspaceUsecase) Execute(ctx context.Context, req *workspace_application.GetWorkspaceRequest) (*workspace_domain.Workspace, error) {
	if err := c.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := c.ValidateRequest(req); err != nil {
		return nil, err
	}

	worskpace, err := c.Repo.GetWorkspace(ctx, workspace_domain.GetRequest{
		WorkspaceID:   req.WorkspaceID,
	})
	if err != nil {
		return nil, err
	}

	return worskpace, nil
}

func (c *GetWorkspaceUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return errors.New(workspace_domain.ErrRepositoryNil)
	}
	return nil
}

func (c *GetWorkspaceUsecase) ValidateRequest(req *workspace_application.GetWorkspaceRequest) error {
	if req == nil {
		return errors.New(workspace_domain.ErrRequestRequired)
	}

	if req.WorkspaceID == "" {
		return errors.New(workspace_domain.ErrWorkspaceIDRequired)
	}

	return nil
}

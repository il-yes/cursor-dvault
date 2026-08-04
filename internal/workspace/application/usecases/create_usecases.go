package workspace_usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	workspace_application "vault-app/internal/workspace/application"
	workspace_events "vault-app/internal/workspace/application/events"
	workspace_domain "vault-app/internal/workspace/domain"
)

type CreateWorkspaceUsecase struct {
	Repo      workspace_domain.Repository
	DomainBus workspace_events.WorkspaceEventBus
}

func NewCreateWorkspaceUsecase(repo workspace_domain.Repository, workspaceBus workspace_events.WorkspaceEventBus) *CreateWorkspaceUsecase {
	return &CreateWorkspaceUsecase{
		Repo:      repo,
		DomainBus: workspaceBus,
	}
}

func (c *CreateWorkspaceUsecase) Execute(ctx context.Context, req *workspace_application.CreateWorkspaceRequest) (*workspace_domain.Workspace, error) {
	if err := c.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := c.ValidateRequest(req); err != nil {
		return nil, err
	}

	workspace := workspace_domain.NewWorkspace(req.VaultID, req.Name, req.Description, req.OwnerID)

	created, err := c.Repo.CreateWorkspace(ctx, workspace_domain.CreateRequest{
		UserID:    req.OwnerID,
		VaultID:   req.VaultID,
		Workspace: workspace,
		Signature: req.Signature,
	})
	if err != nil {
		return nil, err
	}
	if created == nil {
		return nil, errors.New(workspace_domain.ErrRepositoryResponse)
	}

	errEvent := c.DomainBus.PublishWorkspaceCreated(
		ctx,
		workspace_domain.WorkspaceCreated{
			EventID:        uuid.NewString(),
			EventTimestamp: time.Now(),
			WorkspaceID:    workspace.ID,
			VaultID:        workspace.VaultID,
			OwnerID:        workspace.OwnerID,
			Name:           workspace.Name,
		},
	)
	if errEvent != nil {
		return nil, errEvent
	}

	return &created.Data, nil
}

func (c *CreateWorkspaceUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return errors.New(workspace_domain.ErrRepositoryNil)
	}

	if c.DomainBus == nil {
		return errors.New(workspace_domain.ErrWorkspaceBusRequired)
	}

	return nil
}

func (c *CreateWorkspaceUsecase) ValidateRequest(req *workspace_application.CreateWorkspaceRequest) error {
	if req == nil {
		return errors.New(workspace_domain.ErrRequestRequired)
	}

	if req.VaultID == "" {
		return errors.New(workspace_domain.ErrVaultIDRequired)
	}

	if req.OwnerID == "" {
		return errors.New(workspace_domain.ErrWorkspaceOwnerRequired)
	}

	if req.Name == "" {
		return errors.New(workspace_domain.ErrWorkspaceNameRequired)
	}
	if req.Signature == "" {
		return errors.New(workspace_domain.ErrSignatureMissing)
	}
	return nil
}

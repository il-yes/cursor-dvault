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

type RenameWorkspaceUsecase struct {
	Repo      workspace_domain.Repository
	DomainBus workspace_events.WorkspaceEventBus
}

func NewRenameWorkspaceUsecase(repo workspace_domain.Repository, workspaceBus workspace_events.WorkspaceEventBus) *RenameWorkspaceUsecase {
	return &RenameWorkspaceUsecase{
		Repo:      repo,
		DomainBus: workspaceBus,
	}
}

func (c *RenameWorkspaceUsecase) Execute(ctx context.Context, req *workspace_application.RenameWorkspaceRequest) (*workspace_domain.Workspace, error) {
	if err := c.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := c.ValidateRequest(req); err != nil {
		return nil, err
	}

	workspace := req.Workspace
	oldName := workspace.Name
	workspace.Rename(req.Name)
	if err := workspace.Rename(req.Name); err != nil {
		return nil, err
	}

	renamed, err := c.Repo.UpdateWorkspace(ctx, workspace_domain.UpdateRequest{
		UserID:    req.Workspace.OwnerID,
		VaultID:   workspace.VaultID,
		Workspace: *workspace,
		Signature: req.Signature,
	})
	if err != nil {
		return nil, err
	}

	errEvent := c.DomainBus.PublishWorkspaceRenamed(
		ctx,
		workspace_domain.WorkspaceRenamed{
			EventID:        uuid.NewString(),
			EventTimestamp: time.Now(),
			WorkspaceID:    workspace.ID,
			OldName:        oldName,
			NewName:        workspace.Name,
		},
	)
	if errEvent != nil {
		return nil, errEvent
	}

	return &renamed.Data, nil
}

func (c *RenameWorkspaceUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return errors.New(workspace_domain.ErrRepositoryNil)
	}

	if c.DomainBus == nil {
		return errors.New(workspace_domain.ErrWorkspaceBusRequired)
	}

	return nil
}

func (c *RenameWorkspaceUsecase) ValidateRequest(req *workspace_application.RenameWorkspaceRequest) error {
	if req == nil {
		return errors.New(workspace_domain.ErrRequestRequired)
	}

	if req.Name == "" {
		return errors.New(workspace_domain.ErrWorkspaceNameRequired)
	}

	if req.Workspace == nil {
		return errors.New(workspace_domain.ErrWorkspaceOwnerRequired)
	}

	if req.Signature == "" {
		return errors.New(workspace_domain.ErrSignatureMissing)
	}
	return nil
}

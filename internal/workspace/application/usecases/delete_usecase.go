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

type DeleteWorkspaceUsecase struct {
	Repo      workspace_domain.Repository
	DomainBus workspace_events.WorkspaceEventBus
}

func NewDeleteWorkspaceUsecase(repo workspace_domain.Repository, workspaceBus workspace_events.WorkspaceEventBus) *DeleteWorkspaceUsecase {
	return &DeleteWorkspaceUsecase{
		Repo:      repo,
		DomainBus: workspaceBus,
	}
}

func (c *DeleteWorkspaceUsecase) Execute(ctx context.Context, req *workspace_application.DeleteWorkspaceRequest) error {
	if err := c.ValidateDependencies(); err != nil {
		return err
	}

	if err := c.ValidateRequest(req); err != nil {
		return err
	}

	err := c.Repo.DeleteWorkspace(ctx, workspace_domain.DeleteRequest{
		WorkspaceID: req.WorkspaceID,
		Signature: req.Signature,
	})
	if err != nil {
		return err
	}

	errEvent := c.DomainBus.PublishWorkspaceDeleted(
		ctx,
		workspace_domain.WorkspaceDeleted{
			EventID:        uuid.NewString(),
			EventTimestamp: time.Now(),
			WorkspaceID:    req.WorkspaceID,
			VaultID:        req.VaultID,
		},
	)
	if errEvent != nil {
		return errEvent
	}

	return nil
}

func (c *DeleteWorkspaceUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return errors.New(workspace_domain.ErrRepositoryNil)
	}
	if c.DomainBus == nil {
		return errors.New(workspace_domain.ErrWorkspaceBusRequired)
	}
	return nil
}

func (c *DeleteWorkspaceUsecase) ValidateRequest(req *workspace_application.DeleteWorkspaceRequest) error {
	if req == nil {
		return errors.New(workspace_domain.ErrRequestRequired)
	}

	if req.VaultID == "" {
		return errors.New(workspace_domain.ErrVaultIDRequired)
	}
	if req.WorkspaceID == "" {
		return errors.New(workspace_domain.ErrWorkspaceIDRequired)
	}
	if req.Signature == "" {
		return errors.New(workspace_domain.ErrSignatureMissing)
	}

	return nil
}

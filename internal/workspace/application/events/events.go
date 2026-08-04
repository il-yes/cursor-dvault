package workspace_events

import (
	"context"

	workspace_domain "vault-app/internal/workspace/domain"
)

// -------- EVENTS --------

// -------- EVENT BUS --------
type WorkspaceEventBus interface {
	PublishWorkspaceCreated(ctx context.Context, event workspace_domain.WorkspaceCreated) error
	SubscribeToWorkspaceCreated(handler func(ctx context.Context, event workspace_domain.WorkspaceCreated)) error


	PublishWorkspaceRenamed(ctx context.Context, event workspace_domain.WorkspaceRenamed) error
	SubscribeToWorkspaceRenamed(handler func(ctx context.Context, event workspace_domain.WorkspaceRenamed)) error


	PublishWorkspaceDeleted(ctx context.Context, event workspace_domain.WorkspaceDeleted) error
	SubscribeToWorkspaceDeleted(handler func(ctx context.Context, event workspace_domain.WorkspaceDeleted)) error
}

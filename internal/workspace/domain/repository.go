package workspace_domain

import (
	"context"

	tracecore_types "vault-app/internal/tracecore/types"
)

type CreateRequest struct {
	UserID    string
	VaultID   string
	Workspace Workspace
	Signature string
}
type UpdateRequest struct {
	UserID    string
	VaultID   string
	Workspace Workspace
	Signature string
}
type DeleteRequest struct {
	WorkspaceID string
	Signature string
}
type GetRequest struct {
	WorkspaceID string
}
type ListRequest struct {
	VaultID string
}
type Repository interface {
	CreateWorkspace(ctx context.Context, req CreateRequest) (*tracecore_types.CloudResponse[Workspace], error)
	UpdateWorkspace(ctx context.Context, req UpdateRequest) (*tracecore_types.CloudResponse[Workspace], error)
	DeleteWorkspace(ctx context.Context, req DeleteRequest) error
	GetWorkspace(ctx context.Context, req GetRequest) (*Workspace, error)
	ListWorkspace(ctx context.Context, req ListRequest) ([]Workspace, error)
}

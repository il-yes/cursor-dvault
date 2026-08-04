package workspace_application

import workspace_domain "vault-app/internal/workspace/domain"


type CreateWorkspaceRequest struct {
	VaultID string
	OwnerID string
	Name string
	Description string
	Signature string
}

type RenameWorkspaceRequest struct {
	Workspace *workspace_domain.Workspace
	Name string
	Signature string
}

type DeleteWorkspaceRequest struct {
	WorkspaceID string
	Signature string
	VaultID string
}

type GetWorkspaceRequest struct {
	WorkspaceID string
}

type ListWorkspacesRequest struct {
	VaultID string
}


package workspace_domain

var (
	ErrInvalidName   = "Invalid Workspace name"
	ErrRepositoryNil = "repository is nil"
	ErrSignatureMissing = "Signature is required"
	ErrWorkspaceOwnerRequired = "owner id is required"
	ErrWorkspaceIDRequired = "workspace id is required"
	ErrWorkspaceNameRequired = "workspace name is required"
	ErrVaultIDRequired = "vault id is required"
	ErrWorkspaceBusRequired = "Event bus is nil"
	ErrRequestRequired = "Request is nil"
	ErrRepositoryResponse = "workspace repository returned nil response"
)

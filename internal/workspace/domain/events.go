// internal/workspace/domain/events.go

package workspace_domain

import "time"

const (
	EventWorkspaceCreated = "workspace.created"
	EventWorkspaceRenamed = "workspace.renamed"
	EventWorkspaceDeleted = "workspace.deleted"
)

type WorkspaceCreated struct {
	EventID        string
	EventTimestamp time.Time

	WorkspaceID string
	VaultID     string
	OwnerID     string

	Name string
}

func (WorkspaceCreated) EventType() string {
	return EventWorkspaceCreated
}

type WorkspaceRenamed struct {
	EventID        string
	EventTimestamp time.Time

	WorkspaceID string

	OldName string
	NewName string
}

func (WorkspaceRenamed) EventType() string {
	return EventWorkspaceRenamed
}

type WorkspaceDeleted struct {
	EventID        string
	EventTimestamp time.Time

	WorkspaceID string
	VaultID     string
}

func (WorkspaceDeleted) EventType() string {
	return EventWorkspaceDeleted
}
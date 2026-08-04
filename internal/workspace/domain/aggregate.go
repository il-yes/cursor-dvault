package workspace_domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type WorkspaceStatus string

const (
	WorkspaceActive   WorkspaceStatus = "active"
	WorkspaceArchived WorkspaceStatus = "archived"
)

type Workspace struct {
	ID      string `json:"id"`
	VaultID string `json:"vault_id"`

	Name        string          `json:"name"`
	Description string          `json:"description"`
	Status      WorkspaceStatus `json:"status"`

	OwnerID string `json:"owner_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDraft   bool      `json:"is_draft"`
	IsDirty   bool      `json:"is_dirty" gorm:"boolean"`
}

func (w *Workspace) Rename(name string) error {

	if name == "" {
		return errors.New(ErrInvalidName)
	}

	w.Name = name

	w.UpdatedAt = time.Now()

	return nil
}

func NewWorkspace(vaultID string, name string, description string, ownerID string) Workspace {
	now := time.Now()
	return Workspace{
		ID:          uuid.NewString(),
		VaultID:     vaultID,
		Name:        name,
		Description: description,
		OwnerID: ownerID,
		Status:      WorkspaceActive,
		CreatedAt: now,
		UpdatedAt:   now,
		IsDraft: true,
		IsDirty: false,
	}
}

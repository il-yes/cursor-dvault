package tracecore_types

import "time"

// CloudWorkspaceDTO mirrors the Ankhora Cloud backend Workspace representation.
// Cloud serializes Workspace aggregates using default Go JSON field names (PascalCase).
// It MUST be decoded into this infrastructure DTO and explicitly mapped to
// tracecore_types.Workspace / workspace_domain.Workspace.
type CloudWorkspaceDTO struct {
	ID          string    `json:"ID"`
	VaultID     string    `json:"VaultID"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	Status      string    `json:"Status"`
	OwnerID     string    `json:"OwnerID"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	IsDraft     bool      `json:"IsDraft"`
	IsDirty     bool      `json:"IsDirty"`
}

// MapCloudWorkspaceToTypes converts a CloudWorkspaceDTO into a tracecore_types.Workspace.
func MapCloudWorkspaceToTypes(dto *CloudWorkspaceDTO) *Workspace {
	if dto == nil {
		return nil
	}
	return &Workspace{
		ID:          dto.ID,
		VaultID:     dto.VaultID,
		Name:        dto.Name,
		Description: dto.Description,
		Status:      dto.Status,
		OwnerID:     dto.OwnerID,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
		IsDraft:     dto.IsDraft,
		IsDirty:     dto.IsDirty,
	}
}

package collaboration_domain

import "errors"

var (
	ErrInvalidResourceType = errors.New("resource_type is required")
	ErrInvalidResourceID   = errors.New("resource_id is required")
)

const (
	ResourceTypeShareEntry = "share_entry"
	ResourceTypeVaultEntry = "vault_entry"
	ResourceTypeDocument   = "document"
	ResourceTypeMessage    = "message"
	ResourceTypeEvent      = "event"
)

type ResourceReference struct {
	ResourceType  string `json:"resource_type"`
	ResourceID    string `json:"resource_id"`
	SourceEventID string `json:"source_event_id,omitempty"`
}

func NewResourceReference(resourceType string, resourceID string, sourceEventID string) (ResourceReference, error) {
	if resourceType == "" {
		return ResourceReference{}, ErrInvalidResourceType
	}
	if resourceID == "" {
		return ResourceReference{}, ErrInvalidResourceID
	}
	return ResourceReference{
		ResourceType:  resourceType,
		ResourceID:    resourceID,
		SourceEventID: sourceEventID,
	}, nil
}

func (r ResourceReference) Validate() error {
	if r.ResourceType == "" {
		return ErrInvalidResourceType
	}
	if r.ResourceID == "" {
		return ErrInvalidResourceID
	}
	return nil
}

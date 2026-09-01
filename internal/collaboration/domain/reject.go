package collaboration_domain

import (
	"time"

	"github.com/google/uuid"
)

type RejectStatus string

const (
	RejectCreated RejectStatus = "created"
)

type Reject struct {
	ID                string            `json:"id"`
	ResourceReference ResourceReference `json:"resource_reference"`
	Message           string            `json:"message"`
	CreatedBy         string            `json:"created_by"`
	CreatedAt         time.Time         `json:"created_at"`
	Status            RejectStatus      `json:"status"`
}

func NewReject(ref ResourceReference, message string, createdBy string) (Reject, error) {
	if err := ref.Validate(); err != nil {
		return Reject{}, err
	}
	if message == "" {
		return Reject{}, ErrRejectionReasonRequired
	}
	if createdBy == "" {
		return Reject{}, ErrCreatedByRequired
	}

	return Reject{
		ID:                uuid.NewString(),
		ResourceReference: ref,
		Message:           message,
		CreatedBy:         createdBy,
		CreatedAt:         time.Now(),
		Status:            RejectCreated,
	}, nil
}

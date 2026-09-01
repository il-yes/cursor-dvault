package collaboration_domain

import (
	"time"

	"github.com/google/uuid"
)

type ApprovalStatus string

const (
	ApprovalRequested ApprovalStatus = "requested"
	Approved          ApprovalStatus = "approved"
	ApprovalRejected  ApprovalStatus = "rejected"
)

type Approval struct {
	ID                string            `json:"id"`
	ResourceReference ResourceReference `json:"resource_reference"`
	Message           string            `json:"message"`
	CreatedBy         string            `json:"created_by"`
	CreatedAt         time.Time         `json:"created_at"`
	Status            ApprovalStatus    `json:"status"`
	ApprovedBy        string            `json:"approved_by,omitempty"`
	ApprovedAt        *time.Time        `json:"approved_at,omitempty"`
}

func NewApproval(ref ResourceReference, message string, createdBy string) (Approval, error) {
	if err := ref.Validate(); err != nil {
		return Approval{}, err
	}
	if createdBy == "" {
		return Approval{}, ErrCreatedByRequired
	}

	return Approval{
		ID:                uuid.NewString(),
		ResourceReference: ref,
		Message:           message,
		CreatedBy:         createdBy,
		CreatedAt:         time.Now(),
		Status:            ApprovalRequested,
	}, nil
}

func (a *Approval) Approve(approvedBy string) error {
	if approvedBy == "" {
		return ErrCreatedByRequired
	}
	if a.Status == Approved {
		// Idempotent approval
		return nil
	}
	if a.Status != ApprovalRequested {
		return ErrApprovalAlreadyFinalized
	}

	now := time.Now()
	a.Status = Approved
	a.ApprovedBy = approvedBy
	a.ApprovedAt = &now
	return nil
}

func (a *Approval) RejectApproval(rejectedBy string) error {
	if rejectedBy == "" {
		return ErrCreatedByRequired
	}
	if a.Status == ApprovalRejected {
		// Idempotent reject
		return nil
	}
	if a.Status != ApprovalRequested {
		return ErrApprovalAlreadyFinalized
	}

	now := time.Now()
	a.Status = ApprovalRejected
	a.ApprovedBy = rejectedBy
	a.ApprovedAt = &now
	return nil
}

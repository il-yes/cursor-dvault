package collaboration_domain

import (
	"time"

	"github.com/google/uuid"
)

type TransferStatus string

const (
	TransferRequested TransferStatus = "transfer_requested"
	TransferApproved  TransferStatus = "transfer_approved"
	TransferRejected  TransferStatus = "transfer_rejected"
	TransferCompleted TransferStatus = "transfer_completed"
)

type Transfer struct {
	ID                 string            `json:"id"`
	ResourceReference  ResourceReference `json:"resource_reference"`
	TargetTrustGroupID string            `json:"target_trust_group_id"`
	Message            string            `json:"message"`
	CreatedBy          string            `json:"created_by"`
	CreatedAt          time.Time         `json:"created_at"`
	Status             TransferStatus    `json:"status"`
	ApprovedBy         string            `json:"approved_by,omitempty"`
	CompletedBy        string            `json:"completed_by,omitempty"`
	CompletedAt        *time.Time        `json:"completed_at,omitempty"`
}

func NewTransfer(ref ResourceReference, targetTrustGroupID string, message string, createdBy string) (Transfer, error) {
	if err := ref.Validate(); err != nil {
		return Transfer{}, err
	}
	if targetTrustGroupID == "" {
		return Transfer{}, ErrTargetTrustGroupRequired
	}
	if createdBy == "" {
		return Transfer{}, ErrCreatedByRequired
	}

	return Transfer{
		ID:                 uuid.NewString(),
		ResourceReference:  ref,
		TargetTrustGroupID: targetTrustGroupID,
		Message:            message,
		CreatedBy:          createdBy,
		CreatedAt:          time.Now(),
		Status:             TransferRequested,
	}, nil
}

func (t *Transfer) Approve(approvedBy string) error {
	if approvedBy == "" {
		return ErrCreatedByRequired
	}
	if t.Status == TransferApproved {
		return nil // Idempotent
	}
	if t.Status == TransferCompleted {
		return ErrTransferAlreadyCompleted
	}
	if t.Status != TransferRequested {
		return ErrInvalidTransferTransition
	}

	t.Status = TransferApproved
	t.ApprovedBy = approvedBy
	return nil
}

func (t *Transfer) Reject(rejectedBy string) error {
	if rejectedBy == "" {
		return ErrCreatedByRequired
	}
	if t.Status == TransferRejected {
		return nil // Idempotent
	}
	if t.Status == TransferCompleted {
		return ErrTransferAlreadyCompleted
	}
	if t.Status != TransferRequested {
		return ErrInvalidTransferTransition
	}

	t.Status = TransferRejected
	t.ApprovedBy = rejectedBy
	return nil
}

func (t *Transfer) Complete(completedBy string) error {
	if completedBy == "" {
		return ErrCreatedByRequired
	}
	if t.Status == TransferCompleted {
		return nil // Idempotent
	}
	if t.Status == TransferRequested {
		return ErrTransferNotApproved
	}
	if t.Status != TransferApproved {
		return ErrInvalidTransferTransition
	}

	now := time.Now()
	t.Status = TransferCompleted
	t.CompletedBy = completedBy
	t.CompletedAt = &now
	return nil
}

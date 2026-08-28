package collaboration_dtos

import (
	collaboration_domain "vault-app/internal/collaboration/domain"
)

type CreateApprovalRequest struct {
	ResourceType  string `json:"resource_type"`
	ResourceID    string `json:"resource_id"`
	SourceEventID string `json:"source_event_id,omitempty"`
	Message       string `json:"message,omitempty"`
	CreatedBy     string `json:"created_by"`
	ThreadID      string `json:"thread_id,omitempty"`
}

type CreateApprovalResponse struct {
	Approval collaboration_domain.Approval `json:"approval"`
}

type ApproveRequest struct {
	ApprovalID string `json:"approval_id"`
	ApprovedBy string `json:"approved_by"`
	ThreadID   string `json:"thread_id,omitempty"`
}

type ApproveResponse struct {
	Approval collaboration_domain.Approval `json:"approval"`
}

type CreateRejectRequest struct {
	ResourceType  string `json:"resource_type"`
	ResourceID    string `json:"resource_id"`
	SourceEventID string `json:"source_event_id,omitempty"`
	Message       string `json:"message"`
	CreatedBy     string `json:"created_by"`
	ThreadID      string `json:"thread_id,omitempty"`
}

type CreateRejectResponse struct {
	Reject collaboration_domain.Reject `json:"reject"`
}

type CreateTransferRequest struct {
	ResourceType       string `json:"resource_type"`
	ResourceID         string `json:"resource_id"`
	SourceEventID      string `json:"source_event_id,omitempty"`
	TargetTrustGroupID string `json:"target_trust_group_id"`
	Message            string `json:"message,omitempty"`
	CreatedBy          string `json:"created_by"`
	ThreadID           string `json:"thread_id,omitempty"`
}

type CreateTransferResponse struct {
	Transfer collaboration_domain.Transfer `json:"transfer"`
}

type ApproveTransferRequest struct {
	TransferID string `json:"transfer_id"`
	ApprovedBy string `json:"approved_by"`
	ThreadID   string `json:"thread_id,omitempty"`
}

type ApproveTransferResponse struct {
	Transfer collaboration_domain.Transfer `json:"transfer"`
}

type RejectTransferRequest struct {
	TransferID string `json:"transfer_id"`
	RejectedBy string `json:"rejected_by"`
	ThreadID   string `json:"thread_id,omitempty"`
}

type RejectTransferResponse struct {
	Transfer collaboration_domain.Transfer `json:"transfer"`
}

type CompleteTransferRequest struct {
	TransferID  string `json:"transfer_id"`
	CompletedBy string `json:"completed_by"`
	ThreadID    string `json:"thread_id,omitempty"`
}

type CompleteTransferResponse struct {
	Transfer collaboration_domain.Transfer `json:"transfer"`
}

type ListResourceActionsRequest struct {
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
}

type ListResourceActionsResponse struct {
	Approvals []collaboration_domain.Approval `json:"approvals"`
	Rejects   []collaboration_domain.Reject   `json:"rejects"`
	Transfers []collaboration_domain.Transfer `json:"transfers"`
}

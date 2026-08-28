package collaboration_domain

import "context"

type ApprovalRepository interface {
	SaveApproval(ctx context.Context, approval Approval) error
	GetApprovalByID(ctx context.Context, id string) (*Approval, error)
	ListApprovalsByResource(ctx context.Context, resourceType string, resourceID string) ([]Approval, error)
}

type RejectRepository interface {
	SaveReject(ctx context.Context, reject Reject) error
	GetRejectByID(ctx context.Context, id string) (*Reject, error)
	ListRejectsByResource(ctx context.Context, resourceType string, resourceID string) ([]Reject, error)
}

type TransferRepository interface {
	SaveTransfer(ctx context.Context, transfer Transfer) error
	GetTransferByID(ctx context.Context, id string) (*Transfer, error)
	ListTransfersByResource(ctx context.Context, resourceType string, resourceID string) ([]Transfer, error)
}

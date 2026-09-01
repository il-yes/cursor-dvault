package collaboration_domain

import "errors"

var (
	ErrTrustGroupMemberNotFound  = errors.New("trust group member not found")
	ErrApprovalNotFound          = errors.New("approval not found")
	ErrApprovalAlreadyFinalized  = errors.New("approval is already finalized")
	ErrInvalidApprovalTransition = errors.New("invalid approval status transition")
	ErrRejectNotFound            = errors.New("reject not found")
	ErrRejectionReasonRequired   = errors.New("rejection message/reason is required")
	ErrTransferNotFound          = errors.New("transfer not found")
	ErrTargetTrustGroupRequired  = errors.New("target trust group ID is required")
	ErrInvalidTransferTransition = errors.New("invalid transfer lifecycle transition")
	ErrTransferAlreadyCompleted  = errors.New("transfer is already completed")
	ErrTransferNotApproved       = errors.New("transfer must be approved before completion")
	ErrCreatedByRequired         = errors.New("created_by is required")
)

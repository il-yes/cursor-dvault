package collaboration_usecases

import (
	"context"
	"errors"

	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_domain "vault-app/internal/collaboration/domain"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
)

type ActionUseCases struct {
	approvalRepo  collaboration_domain.ApprovalRepository
	rejectRepo    collaboration_domain.RejectRepository
	transferRepo  collaboration_domain.TransferRepository
	appendEventUC *thread_usecase.AppendThreadEventUsecase
}

func NewActionUseCases(
	approvalRepo collaboration_domain.ApprovalRepository,
	rejectRepo collaboration_domain.RejectRepository,
	transferRepo collaboration_domain.TransferRepository,
	appendEventUC *thread_usecase.AppendThreadEventUsecase,
) *ActionUseCases {
	return &ActionUseCases{
		approvalRepo:  approvalRepo,
		rejectRepo:    rejectRepo,
		transferRepo:  transferRepo,
		appendEventUC: appendEventUC,
	}
}

// ---------------------------------------------------------------------------
// Approval Use Cases
// ---------------------------------------------------------------------------

func (uc *ActionUseCases) CreateApproval(ctx context.Context, req collaboration_dtos.CreateApprovalRequest) (*collaboration_dtos.CreateApprovalResponse, error) {
	if uc.approvalRepo == nil {
		return nil, errors.New("approval repository is not initialized")
	}

	ref, err := collaboration_domain.NewResourceReference(req.ResourceType, req.ResourceID, req.SourceEventID)
	if err != nil {
		return nil, err
	}

	approval, err := collaboration_domain.NewApproval(ref, req.Message, req.CreatedBy)
	if err != nil {
		return nil, err
	}

	if err := uc.approvalRepo.SaveApproval(ctx, approval); err != nil {
		return nil, err
	}

	// Emit Thread Event if thread_id provided
	if uc.appendEventUC != nil && req.ThreadID != "" {
		eventRef := thread_domain.EventResourceRef{
			RefType:         thread_domain.ResourceC3Action,
			ActionID:        approval.ID,
			ResourceTypeVal: approval.ResourceReference.ResourceType,
			ResourceIDVal:   approval.ResourceReference.ResourceID,
			SourceEventID:   approval.ResourceReference.SourceEventID,
		}
		idempotencyKey := "evt_appr_req_" + approval.ID
		_, _ = uc.appendEventUC.Execute(ctx, req.ThreadID, string(thread_domain.EventC3ApprovalRequested), eventRef, idempotencyKey)
	}

	return &collaboration_dtos.CreateApprovalResponse{Approval: approval}, nil
}

func (uc *ActionUseCases) Approve(ctx context.Context, req collaboration_dtos.ApproveRequest) (*collaboration_dtos.ApproveResponse, error) {
	if uc.approvalRepo == nil {
		return nil, errors.New("approval repository is not initialized")
	}

	approval, err := uc.approvalRepo.GetApprovalByID(ctx, req.ApprovalID)
	if err != nil {
		return nil, err
	}

	if err := approval.Approve(req.ApprovedBy); err != nil {
		return nil, err
	}

	if err := uc.approvalRepo.SaveApproval(ctx, *approval); err != nil {
		return nil, err
	}

	if uc.appendEventUC != nil && req.ThreadID != "" {
		eventRef := thread_domain.EventResourceRef{
			RefType:         thread_domain.ResourceC3Action,
			ActionID:        approval.ID,
			ResourceTypeVal: approval.ResourceReference.ResourceType,
			ResourceIDVal:   approval.ResourceReference.ResourceID,
			SourceEventID:   approval.ResourceReference.SourceEventID,
		}
		idempotencyKey := "evt_appr_done_" + approval.ID
		_, _ = uc.appendEventUC.Execute(ctx, req.ThreadID, string(thread_domain.EventC3ApprovalApproved), eventRef, idempotencyKey)
	}

	return &collaboration_dtos.ApproveResponse{Approval: *approval}, nil
}

// ---------------------------------------------------------------------------
// Reject Use Cases
// ---------------------------------------------------------------------------

func (uc *ActionUseCases) CreateReject(ctx context.Context, req collaboration_dtos.CreateRejectRequest) (*collaboration_dtos.CreateRejectResponse, error) {
	if uc.rejectRepo == nil {
		return nil, errors.New("reject repository is not initialized")
	}

	ref, err := collaboration_domain.NewResourceReference(req.ResourceType, req.ResourceID, req.SourceEventID)
	if err != nil {
		return nil, err
	}

	reject, err := collaboration_domain.NewReject(ref, req.Message, req.CreatedBy)
	if err != nil {
		return nil, err
	}

	if err := uc.rejectRepo.SaveReject(ctx, reject); err != nil {
		return nil, err
	}

	if uc.appendEventUC != nil && req.ThreadID != "" {
		eventRef := thread_domain.EventResourceRef{
			RefType:         thread_domain.ResourceC3Action,
			ActionID:        reject.ID,
			ResourceTypeVal: reject.ResourceReference.ResourceType,
			ResourceIDVal:   reject.ResourceReference.ResourceID,
			SourceEventID:   reject.ResourceReference.SourceEventID,
		}
		idempotencyKey := "evt_rej_created_" + reject.ID
		_, _ = uc.appendEventUC.Execute(ctx, req.ThreadID, string(thread_domain.EventC3RejectCreated), eventRef, idempotencyKey)
	}

	return &collaboration_dtos.CreateRejectResponse{Reject: reject}, nil
}

// ---------------------------------------------------------------------------
// Transfer Use Cases
// ---------------------------------------------------------------------------

func (uc *ActionUseCases) CreateTransfer(ctx context.Context, req collaboration_dtos.CreateTransferRequest) (*collaboration_dtos.CreateTransferResponse, error) {
	if uc.transferRepo == nil {
		return nil, errors.New("transfer repository is not initialized")
	}

	ref, err := collaboration_domain.NewResourceReference(req.ResourceType, req.ResourceID, req.SourceEventID)
	if err != nil {
		return nil, err
	}

	transfer, err := collaboration_domain.NewTransfer(ref, req.TargetTrustGroupID, req.Message, req.CreatedBy)
	if err != nil {
		return nil, err
	}

	if err := uc.transferRepo.SaveTransfer(ctx, transfer); err != nil {
		return nil, err
	}

	if uc.appendEventUC != nil && req.ThreadID != "" {
		eventRef := thread_domain.EventResourceRef{
			RefType:            thread_domain.ResourceC3Action,
			ActionID:           transfer.ID,
			ResourceTypeVal:    transfer.ResourceReference.ResourceType,
			ResourceIDVal:      transfer.ResourceReference.ResourceID,
			SourceEventID:      transfer.ResourceReference.SourceEventID,
			TargetTrustGroupID: transfer.TargetTrustGroupID,
		}
		idempotencyKey := "evt_tr_req_" + transfer.ID
		_, _ = uc.appendEventUC.Execute(ctx, req.ThreadID, string(thread_domain.EventC3TransferRequested), eventRef, idempotencyKey)
	}

	return &collaboration_dtos.CreateTransferResponse{Transfer: transfer}, nil
}

func (uc *ActionUseCases) ApproveTransfer(ctx context.Context, req collaboration_dtos.ApproveTransferRequest) (*collaboration_dtos.ApproveTransferResponse, error) {
	if uc.transferRepo == nil {
		return nil, errors.New("transfer repository is not initialized")
	}

	transfer, err := uc.transferRepo.GetTransferByID(ctx, req.TransferID)
	if err != nil {
		return nil, err
	}

	if err := transfer.Approve(req.ApprovedBy); err != nil {
		return nil, err
	}

	if err := uc.transferRepo.SaveTransfer(ctx, *transfer); err != nil {
		return nil, err
	}

	if uc.appendEventUC != nil && req.ThreadID != "" {
		eventRef := thread_domain.EventResourceRef{
			RefType:            thread_domain.ResourceC3Action,
			ActionID:           transfer.ID,
			ResourceTypeVal:    transfer.ResourceReference.ResourceType,
			ResourceIDVal:      transfer.ResourceReference.ResourceID,
			SourceEventID:      transfer.ResourceReference.SourceEventID,
			TargetTrustGroupID: transfer.TargetTrustGroupID,
		}
		idempotencyKey := "evt_tr_appr_" + transfer.ID
		_, _ = uc.appendEventUC.Execute(ctx, req.ThreadID, string(thread_domain.EventC3TransferApproved), eventRef, idempotencyKey)
	}

	return &collaboration_dtos.ApproveTransferResponse{Transfer: *transfer}, nil
}

func (uc *ActionUseCases) RejectTransfer(ctx context.Context, req collaboration_dtos.RejectTransferRequest) (*collaboration_dtos.RejectTransferResponse, error) {
	if uc.transferRepo == nil {
		return nil, errors.New("transfer repository is not initialized")
	}

	transfer, err := uc.transferRepo.GetTransferByID(ctx, req.TransferID)
	if err != nil {
		return nil, err
	}

	if err := transfer.Reject(req.RejectedBy); err != nil {
		return nil, err
	}

	if err := uc.transferRepo.SaveTransfer(ctx, *transfer); err != nil {
		return nil, err
	}

	if uc.appendEventUC != nil && req.ThreadID != "" {
		eventRef := thread_domain.EventResourceRef{
			RefType:            thread_domain.ResourceC3Action,
			ActionID:           transfer.ID,
			ResourceTypeVal:    transfer.ResourceReference.ResourceType,
			ResourceIDVal:      transfer.ResourceReference.ResourceID,
			SourceEventID:      transfer.ResourceReference.SourceEventID,
			TargetTrustGroupID: transfer.TargetTrustGroupID,
		}
		idempotencyKey := "evt_tr_rej_" + transfer.ID
		_, _ = uc.appendEventUC.Execute(ctx, req.ThreadID, string(thread_domain.EventC3TransferRejected), eventRef, idempotencyKey)
	}

	return &collaboration_dtos.RejectTransferResponse{Transfer: *transfer}, nil
}

func (uc *ActionUseCases) CompleteTransfer(ctx context.Context, req collaboration_dtos.CompleteTransferRequest) (*collaboration_dtos.CompleteTransferResponse, error) {
	if uc.transferRepo == nil {
		return nil, errors.New("transfer repository is not initialized")
	}

	transfer, err := uc.transferRepo.GetTransferByID(ctx, req.TransferID)
	if err != nil {
		return nil, err
	}

	if err := transfer.Complete(req.CompletedBy); err != nil {
		return nil, err
	}

	if err := uc.transferRepo.SaveTransfer(ctx, *transfer); err != nil {
		return nil, err
	}

	if uc.appendEventUC != nil && req.ThreadID != "" {
		eventRef := thread_domain.EventResourceRef{
			RefType:            thread_domain.ResourceC3Action,
			ActionID:           transfer.ID,
			ResourceTypeVal:    transfer.ResourceReference.ResourceType,
			ResourceIDVal:      transfer.ResourceReference.ResourceID,
			SourceEventID:      transfer.ResourceReference.SourceEventID,
			TargetTrustGroupID: transfer.TargetTrustGroupID,
		}
		idempotencyKey := "evt_tr_comp_" + transfer.ID
		_, _ = uc.appendEventUC.Execute(ctx, req.ThreadID, string(thread_domain.EventC3TransferCompleted), eventRef, idempotencyKey)
	}

	return &collaboration_dtos.CompleteTransferResponse{Transfer: *transfer}, nil
}

// ---------------------------------------------------------------------------
// Query Use Cases
// ---------------------------------------------------------------------------

func (uc *ActionUseCases) ListResourceActions(ctx context.Context, req collaboration_dtos.ListResourceActionsRequest) (*collaboration_dtos.ListResourceActionsResponse, error) {
	resp := &collaboration_dtos.ListResourceActionsResponse{
		Approvals: []collaboration_domain.Approval{},
		Rejects:   []collaboration_domain.Reject{},
		Transfers: []collaboration_domain.Transfer{},
	}

	if uc.approvalRepo != nil {
		apps, err := uc.approvalRepo.ListApprovalsByResource(ctx, req.ResourceType, req.ResourceID)
		if err == nil && apps != nil {
			resp.Approvals = apps
		}
	}

	if uc.rejectRepo != nil {
		rejs, err := uc.rejectRepo.ListRejectsByResource(ctx, req.ResourceType, req.ResourceID)
		if err == nil && rejs != nil {
			resp.Rejects = rejs
		}
	}

	if uc.transferRepo != nil {
		trs, err := uc.transferRepo.ListTransfersByResource(ctx, req.ResourceType, req.ResourceID)
		if err == nil && trs != nil {
			resp.Transfers = trs
		}
	}

	return resp, nil
}

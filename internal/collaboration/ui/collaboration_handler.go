package collaboration_ui

import (
	"context"
	"errors"
	"time"

	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_ports "vault-app/internal/collaboration/application/ports"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	app_config "vault-app/internal/config"
	app_config_domain "vault-app/internal/config/domain"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	tracecore_types "vault-app/internal/tracecore/types"
	vault_dto "vault-app/internal/vault/application/dto"
)

type CollaborationHandler struct {
	createCollabShareUC  *collaboration_usecases.CreateCollaborativeShareUseCase
	resolveCollabShareUC *collaboration_usecases.ResolveCollaborativeShareUseCase
	appendEventUC        *thread_usecase.AppendThreadEventUsecase
	actionUseCases       *collaboration_usecases.ActionUseCases
}

func NewCollaborationHandler(
	createCollabShareUC *collaboration_usecases.CreateCollaborativeShareUseCase,
	resolveCollabShareUC *collaboration_usecases.ResolveCollaborativeShareUseCase,
	appendEventUC *thread_usecase.AppendThreadEventUsecase,
) *CollaborationHandler {
	return &CollaborationHandler{
		createCollabShareUC:  createCollabShareUC,
		resolveCollabShareUC: resolveCollabShareUC,
		appendEventUC:        appendEventUC,
	}
}

func NewCollaborationHandlerWithActions(
	createCollabShareUC *collaboration_usecases.CreateCollaborativeShareUseCase,
	resolveCollabShareUC *collaboration_usecases.ResolveCollaborativeShareUseCase,
	appendEventUC *thread_usecase.AppendThreadEventUsecase,
	actionUseCases *collaboration_usecases.ActionUseCases,
) *CollaborationHandler {
	return &CollaborationHandler{
		createCollabShareUC:  createCollabShareUC,
		resolveCollabShareUC: resolveCollabShareUC,
		appendEventUC:        appendEventUC,
		actionUseCases:       actionUseCases,
	}
}

func (h *CollaborationHandler) SetActionUseCases(uc *collaboration_usecases.ActionUseCases) {
	h.actionUseCases = uc
}

func (h *CollaborationHandler) SetAssetStorage(storage app_config.StorageProvider) {
	if h.createCollabShareUC != nil {
		h.createCollabShareUC.WithStorageProvider(storage)
	}
}

func (h *CollaborationHandler) SetAssetResolver(resolver collaboration_ports.AssetContentResolver) {
    if h.createCollabShareUC != nil {
        h.createCollabShareUC.SetAssetResolver(resolver)
    }

    if h.resolveCollabShareUC != nil {
        h.resolveCollabShareUC.SetAssetResolver(resolver)
    }
}

// CreateCollaborativeShare persists a C3 share entry through the real
// Cloud persistence path and returns the authoritative ShareEntryRef.
func (h *CollaborationHandler) CreateCollaborativeShare(
	ctx context.Context,
	userID string,
	threadID string,
	trustGroupID string,
	assetCID string,
	notes string,
	password string,
	stellarSecret string,
) (*tracecore_types.ShareEntryRefDTO, error) {
	if h.createCollabShareUC == nil {
		return nil, errors.New("create collaborative share use case is not initialized")
	}

	req := collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: trustGroupID,
		CreatedBy:    userID,
		AssetCID:     assetCID,
		Metadata: map[string]string{
			"notes":     notes,
			"thread_id": threadID,
		},
		Password:     password,
		StellarSecret: stellarSecret,
	}

	resp, err := h.createCollabShareUC.Execute(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.ShareEntry.ID == "" {
		return nil, errors.New("cloud did not return a persisted share entry id")
	}

	createdShareID := resp.ShareEntry.ID
	createdTrustGroupID := resp.ShareEntry.TrustGroupID
	if createdTrustGroupID == "" {
		createdTrustGroupID = trustGroupID
	}

	shareRef := tracecore_types.ShareEntryRefDTO{
		ShareEntryID: createdShareID,
		TrustGroupID: createdTrustGroupID,
		AssetCID:     assetCID,
		CreatedBy:    userID,
		Status:       string(resp.ShareEntry.Status),
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	if h.appendEventUC != nil && threadID != "" {
		refPayload := thread_domain.EventResourceRef{
			RefType:      thread_domain.ResourceShareEntry,
			ShareEntryID: createdShareID,
			TrustGroupID: createdTrustGroupID,
		}
		idempotencyKey := "evt_share_" + createdShareID
		_, err := h.appendEventUC.Execute(ctx, threadID, "entry.shared", refPayload, idempotencyKey)
		if err != nil {
			return &shareRef, err
		}
	}

	return &shareRef, nil
}

func (h *CollaborationHandler) ResolveCollaborativeShare(
    ctx context.Context,
    callerIdentityID string,
    callerVaultID string,
    shareEntryID string,
    stellarAccount app_config_domain.StellarAccountConfig,
    getFilePayload vault_dto.GetFileFromIPFSRequest,
) (*collaboration_dtos.ResolveCollaborativeShareResponse, error) {

    req := collaboration_dtos.ResolveCollaborativeShareRequest{
        ShareEntryID:     shareEntryID,
        CallerIdentityID: callerIdentityID,
        CallerVaultID:    callerVaultID,
        CallerUserID:     callerIdentityID,
        StellarAccount:   stellarAccount,
        GetIPFSFile:      getFilePayload,
    }

    return h.resolveCollabShareUC.Execute(ctx, req)
}

// ---------------------------------------------------------------------------
// C3 Collaboration Actions API
// ---------------------------------------------------------------------------

func (h *CollaborationHandler) CreateApproval(
	ctx context.Context,
	req collaboration_dtos.CreateApprovalRequest,
) (*collaboration_dtos.CreateApprovalResponse, error) {
	if h.actionUseCases == nil {
		return nil, errors.New("action use cases are not initialized")
	}
	return h.actionUseCases.CreateApproval(ctx, req)
}

func (h *CollaborationHandler) ApproveAction(
	ctx context.Context,
	req collaboration_dtos.ApproveRequest,
) (*collaboration_dtos.ApproveResponse, error) {
	if h.actionUseCases == nil {
		return nil, errors.New("action use cases are not initialized")
	}
	return h.actionUseCases.Approve(ctx, req)
}

func (h *CollaborationHandler) CreateReject(
	ctx context.Context,
	req collaboration_dtos.CreateRejectRequest,
) (*collaboration_dtos.CreateRejectResponse, error) {
	if h.actionUseCases == nil {
		return nil, errors.New("action use cases are not initialized")
	}
	return h.actionUseCases.CreateReject(ctx, req)
}

func (h *CollaborationHandler) CreateTransfer(
	ctx context.Context,
	req collaboration_dtos.CreateTransferRequest,
) (*collaboration_dtos.CreateTransferResponse, error) {
	if h.actionUseCases == nil {
		return nil, errors.New("action use cases are not initialized")
	}
	return h.actionUseCases.CreateTransfer(ctx, req)
}

func (h *CollaborationHandler) ApproveTransferAction(
	ctx context.Context,
	req collaboration_dtos.ApproveTransferRequest,
) (*collaboration_dtos.ApproveTransferResponse, error) {
	if h.actionUseCases == nil {
		return nil, errors.New("action use cases are not initialized")
	}
	return h.actionUseCases.ApproveTransfer(ctx, req)
}

func (h *CollaborationHandler) RejectTransferAction(
	ctx context.Context,
	req collaboration_dtos.RejectTransferRequest,
) (*collaboration_dtos.RejectTransferResponse, error) {
	if h.actionUseCases == nil {
		return nil, errors.New("action use cases are not initialized")
	}
	return h.actionUseCases.RejectTransfer(ctx, req)
}

func (h *CollaborationHandler) CompleteTransferAction(
	ctx context.Context,
	req collaboration_dtos.CompleteTransferRequest,
) (*collaboration_dtos.CompleteTransferResponse, error) {
	if h.actionUseCases == nil {
		return nil, errors.New("action use cases are not initialized")
	}
	return h.actionUseCases.CompleteTransfer(ctx, req)
}

func (h *CollaborationHandler) ListResourceActions(
	ctx context.Context,
	req collaboration_dtos.ListResourceActionsRequest,
) (*collaboration_dtos.ListResourceActionsResponse, error) {
	if h.actionUseCases == nil {
		return nil, errors.New("action use cases are not initialized")
	}
	return h.actionUseCases.ListResourceActions(ctx, req)
}

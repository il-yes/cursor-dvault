package collaboration_ui

import (
	"context"
	"errors"
	"time"

	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	thread_usecase "vault-app/internal/thread/application/usecases"
	tracecore_types "vault-app/internal/tracecore/types"
)

type CollaborationHandler struct {
	createCollabShareUC  *collaboration_usecases.CreateCollaborativeShareUseCase
	resolveCollabShareUC *collaboration_usecases.ResolveCollaborativeShareUseCase
	appendEventUC        *thread_usecase.AppendThreadEventUsecase
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

// CreateCollaborativeShare persists a C3 share entry through the real
// Cloud persistence path and returns the authoritative ShareEntryRef.
//
// wrappedDEK and kekVersion are cryptographic material owned by the
// desktop crypto layer (TrustGroupCryptoOrchestrator.PrepareCollaborativeAsset).
// They must be supplied by the caller; this handler never invents them.
// A successful response only ever contains the share_entry_id returned by
// the Cloud C3 ShareEntry persistence path.
func (h *CollaborationHandler) CreateCollaborativeShare(
	ctx context.Context,
	userID string,
	threadID string,
	trustGroupID string,
	assetCID string,
	targetVaultID string,
	notes string,
	wrappedDEK string,
	kekVersion uint64,
) (*tracecore_types.ShareEntryRefDTO, error) {
	if h.createCollabShareUC == nil {
		return nil, errors.New("create collaborative share use case is not initialized")
	}
	if wrappedDEK == "" {
		return nil, errors.New("wrapped_dek is required: it must be produced by the desktop crypto orchestration path")
	}
	if kekVersion == 0 {
		return nil, errors.New("kek_version is required: it must come from the trust group key state")
	}

	req := collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: trustGroupID,
		KEKVersion:   kekVersion,
		CreatedBy:    userID,
		AssetCID:     assetCID,
		WrappedDEK:   wrappedDEK,
		Metadata: map[string]string{
			"target_vault_id": targetVaultID,
			"notes":           notes,
			"thread_id":       threadID,
		},
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
	userID string,
	shareEntryID string,
	deviceID string,
) (*collaboration_dtos.ResolveCollaborativeShareResponse, error) {
	if h.resolveCollabShareUC == nil {
		return nil, errors.New("resolve collaborative share use case is not initialized")
	}

	req := collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: shareEntryID,
		CallerUserID: userID,
		DeviceID:     deviceID,
	}

	return h.resolveCollabShareUC.Execute(ctx, req)
}

package collaboration_ui

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

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

func (h *CollaborationHandler) CreateCollaborativeShare(
	ctx context.Context,
	userID string,
	threadID string,
	trustGroupID string,
	assetCID string,
	targetVaultID string,
	notes string,
) (*tracecore_types.ShareEntryRefDTO, error) {
	nowStr := time.Now().Format(time.RFC3339)

	var createdShareID string
	var createdTrustGroupID string = trustGroupID

	if h.createCollabShareUC != nil {
		req := collaboration_dtos.CreateCollaborativeShareRequest{
			TrustGroupID: trustGroupID,
			KEKVersion:   1,
			CreatedBy:    userID,
			AssetCID:     assetCID,
			WrappedDEK:   "wrapped_dek_" + assetCID,
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
		if resp != nil {
			createdShareID = resp.ShareEntry.ID
			if resp.ShareEntry.TrustGroupID != "" {
				createdTrustGroupID = resp.ShareEntry.TrustGroupID
			}
		}
	}

	if createdShareID == "" {
		createdShareID = "se_" + uuid.NewString()[:12]
	}

	shareRef := tracecore_types.ShareEntryRefDTO{
		ShareEntryID: createdShareID,
		TrustGroupID: createdTrustGroupID,
		AssetCID:     assetCID,
		CreatedBy:    userID,
		Status:       "active",
		CreatedAt:    nowStr,
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

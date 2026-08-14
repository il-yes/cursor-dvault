package collaboration_ui

import (
	"context"
	"time"

	"github.com/google/uuid"

	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	thread_usecase "vault-app/internal/thread/application/usecases"
	tracecore_types "vault-app/internal/tracecore/types"
)

type CollaborationHandler struct {
	createCollabShareUC *collaboration_usecases.CreateCollaborativeShareUseCase
	appendEventUC       *thread_usecase.AppendThreadEventUsecase
}

func NewCollaborationHandler(
	createCollabShareUC *collaboration_usecases.CreateCollaborativeShareUseCase,
	appendEventUC *thread_usecase.AppendThreadEventUsecase,
) *CollaborationHandler {
	return &CollaborationHandler{
		createCollabShareUC: createCollabShareUC,
		appendEventUC:       appendEventUC,
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
	shareID := "se_" + uuid.NewString()[:12]
	nowStr := time.Now().Format(time.RFC3339)

	shareRef := tracecore_types.ShareEntryRefDTO{
		ShareEntryID: shareID,
		TrustGroupID: trustGroupID,
		AssetCID:     assetCID,
		CreatedBy:    userID,
		Status:       "active",
		CreatedAt:    nowStr,
	}

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
		if err == nil && resp != nil {
			shareRef.ShareEntryID = resp.ShareEntry.ID
			shareRef.TrustGroupID = resp.ShareEntry.TrustGroupID
			shareRef.AssetCID = resp.ShareEntry.AssetCID
			shareRef.CreatedBy = resp.ShareEntry.CreatedBy
		}
	}

	if h.appendEventUC != nil && threadID != "" {
		payload := map[string]interface{}{
			"notes":           notes,
			"target_vault_id": targetVaultID,
			"share_entry_ref": shareRef,
		}
		_, _ = h.appendEventUC.Execute(ctx, threadID, "share.created", payload)
	}

	return &shareRef, nil
}

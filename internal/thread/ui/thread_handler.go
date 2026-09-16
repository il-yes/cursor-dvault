package thread_ui

import (
	"context"
	"fmt"
	"log"

	thread_dtos "vault-app/internal/thread/application/dtos"
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	tracecore_types "vault-app/internal/tracecore/types"
)

type ThreadHandler struct {
	createUseCase      *thread_usecase.CreateThreadUsecase
	listUseCase        *thread_usecase.ListThreadsUsecase
	listEventsUseCase  *thread_usecase.ListThreadEventsUsecase
	appendEventUseCase *thread_usecase.AppendThreadEventUsecase
}

func NewThreadHandler(
	createUC *thread_usecase.CreateThreadUsecase,
	listUC *thread_usecase.ListThreadsUsecase,
	listEventsUC *thread_usecase.ListThreadEventsUsecase,
	appendEventUC *thread_usecase.AppendThreadEventUsecase,
) *ThreadHandler {
	return &ThreadHandler{
		createUseCase:      createUC,
		listUseCase:        listUC,
		listEventsUseCase:  listEventsUC,
		appendEventUseCase: appendEventUC,
	}
}

func (h *ThreadHandler) CreateThread(
	ctx context.Context,
	userID string,
	channelID string,
	title string,
	subtitle string,
	assetType string,
) (*tracecore_types.ThreadDTO, error) {
	if h.createUseCase == nil {
		return nil, fmt.Errorf("create thread use case is not initialized")
	}

	req := thread_dtos.CreateThreadRequest{
		ChannelID:  channelID,
		IdentityID: userID,
		AssetType:  assetType,
		Title:      title,
		Subtitle:   subtitle,
	}

	th, err := h.createUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	return toTracecoreThreadDTO(th), nil
}

func (h *ThreadHandler) ListThreads(
	ctx context.Context,
	userID string,
	channelID string,
) ([]tracecore_types.ThreadDTO, error) {
	log.Printf("[THREAD LIST HANDLER] userID=%s channelID=%s", userID, channelID)
	if h.listUseCase == nil {
		return nil, fmt.Errorf("list thread use case is not initialized")
	}

	threads, err := h.listUseCase.Execute(ctx, channelID)
	if err != nil {
		return nil, err
	}

	res := make([]tracecore_types.ThreadDTO, 0, len(threads))
	for _, th := range threads {
		res = append(res, *toTracecoreThreadDTO(&th))
	}

	log.Printf("[THREAD LIST HANDLER RETURN] count=%d", len(res))
	return res, nil
}

func (h *ThreadHandler) ListThreadEvents(
	ctx context.Context,
	userID string,
	threadID string,
) ([]tracecore_types.ThreadEventDTO, error) {
	if h.listEventsUseCase == nil {
		return nil, fmt.Errorf("list thread events use case is not initialized")
	}

	events, err := h.listEventsUseCase.Execute(ctx, threadID)
	if err != nil {
		return nil, err
	}

	res := make([]tracecore_types.ThreadEventDTO, 0, len(events))
	for _, evt := range events {
		res = append(res, *toTracecoreThreadEventDTO(&evt))
	}

	return res, nil
}

func (h *ThreadHandler) AppendThreadEvent(
	ctx context.Context,
	userID string,
	threadID string,
	eventType string,
	payload thread_domain.EventResourceRef,
) (*tracecore_types.ThreadEventDTO, error) {
	fmt.Printf("[APPEND][STEP=04] ThreadHandler.AppendThreadEvent enter threadID=%s eventType=%s refType=%s shareEntryID=%s trustGroupID=%s\n", threadID, eventType, payload.RefType, payload.ShareEntryID, payload.TrustGroupID)
	if h.appendEventUseCase == nil {
		fmt.Printf("[APPEND][STEP=04] ThreadHandler.AppendThreadEvent error: appendEventUseCase nil\n")
		return nil, fmt.Errorf("append thread event use case is not initialized")
	}

	evt, err := h.appendEventUseCase.Execute(ctx, threadID, eventType, payload)
	if err != nil {
		fmt.Printf("[APPEND][STEP=04] ThreadHandler.AppendThreadEvent execute error=%v\n", err)
		return nil, err
	}
	fmt.Printf("[APPEND][STEP=04] ThreadHandler.AppendThreadEvent success eventID=%s\n", evt.ID)
	return toTracecoreThreadEventDTO(evt), nil
}

func toTracecoreThreadDTO(th *thread_domain.Thread) *tracecore_types.ThreadDTO {
	if th == nil {
		return nil
	}
	return &tracecore_types.ThreadDTO{
		ID:          th.ID,
		ChannelID:   th.ChannelID,
		WorkspaceID: th.WorkspaceID,
		AssetType:   th.AssetType,
		Title:       th.Title,
		Subtitle:    th.Subtitle,
		Status:      string(th.Status),
		CreatedAt:   th.CreatedAt,
	}
}

func toTracecoreThreadEventDTO(evt *thread_domain.ThreadEvent) *tracecore_types.ThreadEventDTO {
	if evt == nil {
		return nil
	}
	var payloadMap map[string]any
	if evt.Payload.RefType == thread_domain.ResourceShareEntry || evt.Payload.ShareEntryID != "" || evt.Type == thread_domain.EventEntryShared {
		refTypeStr := string(evt.Payload.RefType)
		if refTypeStr == "" {
			refTypeStr = string(thread_domain.ResourceShareEntry)
		}
		payloadMap = map[string]any{
			"ref_type":       refTypeStr,
			"share_entry_id": evt.Payload.ShareEntryID,
			"trust_group_id": evt.Payload.TrustGroupID,
		}
	} else if evt.Payload.CID != "" || evt.Payload.AssetType != "" || evt.Payload.RefType == thread_domain.ResourceStorageAsset {
		refTypeStr := string(evt.Payload.RefType)
		if refTypeStr == "" {
			refTypeStr = string(thread_domain.ResourceStorageAsset)
		}
		payloadMap = map[string]any{
			"ref_type":     refTypeStr,
			"cid":          evt.Payload.CID,
			"content_hash": evt.Payload.ContentHash,
			"size":         evt.Payload.Size,
			"asset_type":   evt.Payload.AssetType,
		}
	} else {
		payloadMap = map[string]any{}
		if evt.Payload.RefType != "" {
			payloadMap["ref_type"] = string(evt.Payload.RefType)
		}
		if evt.Payload.ShareEntryID != "" {
			payloadMap["share_entry_id"] = evt.Payload.ShareEntryID
		}
		if evt.Payload.TrustGroupID != "" {
			payloadMap["trust_group_id"] = evt.Payload.TrustGroupID
		}
		if evt.Payload.CID != "" {
			payloadMap["cid"] = evt.Payload.CID
		}
		if evt.Payload.ContentHash != "" {
			payloadMap["content_hash"] = evt.Payload.ContentHash
		}
		if evt.Payload.Size != 0 {
			payloadMap["size"] = evt.Payload.Size
		}
		if evt.Payload.AssetType != "" {
			payloadMap["asset_type"] = evt.Payload.AssetType
		}
	}

	return &tracecore_types.ThreadEventDTO{
		ID:              evt.ID,
		ThreadID:        evt.ThreadID,
		PreviousEventID: evt.PreviousEventID,
		Type:            string(evt.Type),
		Payload:         payloadMap,
		IdempotencyKey:  evt.IdempotencyKey,
		Cursor:          evt.Cursor,
		Headers:         evt.Headers,
		Signature:       evt.Signature,
		CreatedAt:       evt.CreatedAt,
	}
}

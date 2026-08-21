package tracecore

import (
	"context"
	"fmt"
	"log"
	"time"


	thread_domain "vault-app/internal/thread/domain"
	tracecore_types "vault-app/internal/tracecore/types"
	"vault-app/internal/utils"
	workspace_domain "vault-app/internal/workspace/domain"
)

// WorkspaceRepository Implementation on TracecoreClient

func (c *TracecoreClient) UpdateWorkspace(ctx context.Context, req workspace_domain.UpdateRequest) (*tracecore_types.CloudResponse[workspace_domain.Workspace], error) {
	ws := req.Workspace
	ws.UpdatedAt = time.Now()
	return &tracecore_types.CloudResponse[workspace_domain.Workspace]{
		Status:  200,
		Data:    ws,
		Message: "success",
		Success: true,
	}, nil
}

func (c *TracecoreClient) DeleteWorkspace(ctx context.Context, req workspace_domain.DeleteRequest) error {
	return nil
}

func (c *TracecoreClient) GetWorkspace(ctx context.Context, req workspace_domain.GetRequest) (*workspace_domain.Workspace, error) {
	now := time.Now()
	ws := workspace_domain.NewWorkspace("v1", "Default Workspace", "Default workspace", "u1")
	ws.ID = req.WorkspaceID
	ws.CreatedAt = now
	ws.UpdatedAt = now
	return &ws, nil
}

func (c *TracecoreClient) ListWorkspace(ctx context.Context, req workspace_domain.ListRequest) ([]workspace_domain.Workspace, error) {
	utils.LogPretty("[Workspace] Repository.ListWorkspace vault_id", req.VaultID)
	cloudWorkspaces, err := c.ListWorkspaces(ctx, req.VaultID)
	if err != nil {
		utils.LogPretty("🚫 [Workspace] Repository.ListWorkspace error", err)
		return nil, err
	}

	result := make([]workspace_domain.Workspace, 0, len(cloudWorkspaces))
	for _, ws := range cloudWorkspaces {
		createdAt := ws.CreatedAt
		updatedAt := ws.UpdatedAt
		if createdAt.IsZero() {
			createdAt = time.Now()
		}
		if updatedAt.IsZero() {
			updatedAt = time.Now()
		}
		result = append(result, workspace_domain.Workspace{
			ID:          ws.ID,
			VaultID:     ws.VaultID,
			Name:        ws.Name,
			Description: ws.Description,
			Status:      workspace_domain.WorkspaceStatus(ws.Status),
			OwnerID:     ws.OwnerID,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			IsDraft:     ws.IsDraft,
			IsDirty:     ws.IsDirty,
		})
	}
	utils.LogPretty("[Workspace] Repository.ListWorkspace result", result)
	return result, nil
}

// ThreadRepository Implementation on TracecoreClient

func (c *TracecoreClient) CreateThread(ctx context.Context, req *thread_domain.CreateThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	th := req.Thread
	dto, err := c.CreateThreadDirect(ctx, req.IdentityID, th.ChannelID, th.Title, th.Subtitle, th.AssetType)
	if err != nil {
		return nil, fmt.Errorf("cloud create thread failed: %w", err)
	}

	createdAt := dto.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	domainThread := thread_domain.Thread{
		ID:          dto.ID,
		ChannelID:   dto.ChannelID,
		WorkspaceID: th.WorkspaceID,
		AssetType:   dto.AssetType,
		Title:       dto.Title,
		Subtitle:    dto.Subtitle,
		Status:      thread_domain.ThreadStatus(dto.Status),
		CreatedAt:   createdAt,
	}

	return &tracecore_types.CloudResponse[thread_domain.Thread]{
		Status:  200,
		Data:    domainThread,
		Message: "success",
		Success: true,
	}, nil
}

func (c *TracecoreClient) ListThreads(ctx context.Context, req *thread_domain.ListThreadsRequest) (*tracecore_types.CloudResponse[[]thread_domain.Thread], error) {
	log.Printf("[THREAD LIST REPO/TRACECORE] channelID=%s", req.ChannelID)
	dtos, err := c.ListThreadsDirect(ctx, "me", req.ChannelID)
	if err != nil {
		return nil, fmt.Errorf("cloud list threads failed: %w", err)
	}

	result := make([]thread_domain.Thread, 0, len(dtos))
	for _, dto := range dtos {
		createdAt := dto.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now()
		}
		result = append(result, thread_domain.Thread{
			ID:        dto.ID,
			ChannelID: dto.ChannelID,
			AssetType: dto.AssetType,
			Title:     dto.Title,
			Subtitle:  dto.Subtitle,
			Status:    thread_domain.ThreadStatus(dto.Status),
			CreatedAt: createdAt,
		})
	}

	return &tracecore_types.CloudResponse[[]thread_domain.Thread]{
		Status:  200,
		Data:    result,
		Message: "success",
		Success: true,
	}, nil
}

func (c *TracecoreClient) GetThread(ctx context.Context, req *thread_domain.GetThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	resp, err := c.ListThreads(ctx, &thread_domain.ListThreadsRequest{
		ChannelID: "all",
	})
	if err != nil {
		return nil, err
	}
	for _, th := range resp.Data {
		if th.ID == req.ThreadID {
			return &tracecore_types.CloudResponse[thread_domain.Thread]{
				Status:  200,
				Data:    th,
				Message: "success",
				Success: true,
			}, nil
		}
	}
	return nil, fmt.Errorf("thread not found: %s", req.ThreadID)
}

func (c *TracecoreClient) UpdateThread(ctx context.Context, req *thread_domain.UpdateThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	return &tracecore_types.CloudResponse[thread_domain.Thread]{
		Status:  200,
		Data:    req.Thread,
		Message: "success",
		Success: true,
	}, nil
}

func (c *TracecoreClient) ListThreadEvents(ctx context.Context, req *thread_domain.ListThreadEventsRequest) (*tracecore_types.CloudResponse[[]thread_domain.ThreadEvent], error) {
	dtos, err := c.ListThreadEventsDirect(ctx, "me", req.ThreadID)
	if err != nil {
		return nil, fmt.Errorf("cloud list thread events failed: %w", err)
	}

	result := make([]thread_domain.ThreadEvent, 0, len(dtos))
	for _, dto := range dtos {
		result = append(result, mapThreadEventDTO(dto))
	}

	return &tracecore_types.CloudResponse[[]thread_domain.ThreadEvent]{
		Status:  200,
		Data:    result,
		Message: "success",
		Success: true,
	}, nil
}

func (c *TracecoreClient) AppendThreadEvent(ctx context.Context, req *thread_domain.AppendThreadEventRequest) (*tracecore_types.CloudResponse[thread_domain.ThreadEvent], error) {
	payload := eventResourceRefToPayload(req.Payload)

	dto, err := c.AppendThreadEventDirect(ctx, "me", req.ThreadID, req.EventType, payload, req.IdempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("cloud append thread event failed: %w", err)
	}

	domainEvent := mapThreadEventDTO(*dto)

	return &tracecore_types.CloudResponse[thread_domain.ThreadEvent]{
		Status:  200,
		Data:    domainEvent,
		Message: "success",
		Success: true,
	}, nil
}

// eventResourceRefToPayload converts the domain EventResourceRef into the
// wire payload (map) expected by AppendThreadEventDirect.
func eventResourceRefToPayload(ref thread_domain.EventResourceRef) map[string]interface{} {
	payload := map[string]interface{}{}
	if ref.RefType != "" {
		payload["ref_type"] = string(ref.RefType)
	}
	if ref.RefType == thread_domain.ResourceShareEntry {
		if ref.ShareEntryID != "" {
			payload["share_entry_id"] = ref.ShareEntryID
		}
		if ref.TrustGroupID != "" {
			payload["trust_group_id"] = ref.TrustGroupID
		}
	} else {
		if ref.CID != "" {
			payload["cid"] = ref.CID
		}
		if ref.ContentHash != "" {
			payload["content_hash"] = ref.ContentHash
		}
		if ref.Size != 0 {
			payload["size"] = ref.Size
		}
		if ref.AssetType != "" {
			payload["asset_type"] = ref.AssetType
		}
	}
	return payload
}

// mapThreadEventDTO maps a Cloud ThreadEventDTO into the domain ThreadEvent.
func mapThreadEventDTO(dto tracecore_types.ThreadEventDTO) thread_domain.ThreadEvent {
	return thread_domain.ThreadEvent{
		ID:              dto.ID,
		ThreadID:        dto.ThreadID,
		PreviousEventID: dto.PreviousEventID,
		Type:            thread_domain.ThreadEventType(dto.Type),
		Payload:         payloadToEventResourceRef(dto.Payload),
		IdempotencyKey:  dto.IdempotencyKey,
		Cursor:          dto.Cursor,
		Headers:         dto.Headers,
		Signature:       dto.Signature,
		CreatedAt:       dto.CreatedAt,
	}
}

// payloadToEventResourceRef reconstructs the domain EventResourceRef from
// the wire payload map returned by the Cloud backend.
func payloadToEventResourceRef(payload map[string]any) thread_domain.EventResourceRef {
	ref := thread_domain.EventResourceRef{}
	if payload == nil {
		return ref
	}

	if rt, ok := payload["ref_type"].(string); ok {
		ref.RefType = thread_domain.ResourceType(rt)
	}

	// ShareEntry fields
	if v, ok := payload["share_entry_id"].(string); ok {
		ref.ShareEntryID = v
	}
	if v, ok := payload["trust_group_id"].(string); ok {
		ref.TrustGroupID = v
	}

	// Storage asset fields
	if v, ok := payload["cid"].(string); ok {
		ref.CID = v
	}
	if v, ok := payload["content_hash"].(string); ok {
		ref.ContentHash = v
	}
	// JSON numbers unmarshal as float64 in Go
	if v, ok := payload["size"].(float64); ok {
		ref.Size = int64(v)
	}
	if v, ok := payload["asset_type"].(string); ok {
		ref.AssetType = v
	}

	return ref
}

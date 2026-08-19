package tracecore

import (
	"context"
	"time"

	"github.com/google/uuid"

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
	if th.ID == "" {
		th.ID = uuid.NewString()
	}
	th.CreatedAt = time.Now()

	return &tracecore_types.CloudResponse[thread_domain.Thread]{
		Status:  200,
		Data:    th,
		Message: "success",
		Success: true,
	}, nil
}

func (c *TracecoreClient) ListThreads(ctx context.Context, req *thread_domain.ListThreadsRequest) (*tracecore_types.CloudResponse[[]thread_domain.Thread], error) {
	now := time.Now()
	th := thread_domain.NewThread(req.ChannelID, "vault_entry", "Default Thread", "Collaboration thread")
	th.CreatedAt = now

	return &tracecore_types.CloudResponse[[]thread_domain.Thread]{
		Status:  200,
		Data:    []thread_domain.Thread{th},
		Message: "success",
		Success: true,
	}, nil
}

func (c *TracecoreClient) GetThread(ctx context.Context, req *thread_domain.GetThreadRequest) (*tracecore_types.CloudResponse[thread_domain.Thread], error) {
	now := time.Now()
	th := thread_domain.NewThread("ch_1", "vault_entry", "Default Thread", "Collaboration thread")
	th.ID = req.ThreadID
	th.CreatedAt = now

	return &tracecore_types.CloudResponse[thread_domain.Thread]{
		Status:  200,
		Data:    th,
		Message: "success",
		Success: true,
	}, nil
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
	return &tracecore_types.CloudResponse[[]thread_domain.ThreadEvent]{
		Status:  200,
		Data:    []thread_domain.ThreadEvent{},
		Message: "success",
		Success: true,
	}, nil
}

func (c *TracecoreClient) AppendThreadEvent(ctx context.Context, req *thread_domain.AppendThreadEventRequest) (*tracecore_types.CloudResponse[thread_domain.ThreadEvent], error) {
	idempotencyKey := req.IdempotencyKey
	if idempotencyKey == "" {
		idempotencyKey = uuid.NewString()
	}

	evt := thread_domain.ThreadEvent{
		ID:             uuid.NewString(),
		ThreadID:       req.ThreadID,
		Type:           thread_domain.ThreadEventType(req.EventType),
		Payload:        req.Payload,
		IdempotencyKey: idempotencyKey,
		Cursor:         1,
		CreatedAt:      time.Now(),
	}

	return &tracecore_types.CloudResponse[thread_domain.ThreadEvent]{
		Status:  200,
		Data:    evt,
		Message: "success",
		Success: true,
	}, nil
}

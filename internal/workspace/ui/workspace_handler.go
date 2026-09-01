package workspace_ui

import (
	"context"
	"fmt"
	"sync"

	tracecore_types "vault-app/internal/tracecore/types"
	"vault-app/internal/utils"
	workspace_application "vault-app/internal/workspace/application"
	workspace_usecase "vault-app/internal/workspace/application/usecases"
	workspace_domain "vault-app/internal/workspace/domain"
)

type WorkspaceHandler struct {
	createUseCase *workspace_usecase.CreateWorkspaceUsecase
	listUseCase   *workspace_usecase.ListWorkspaceUsecase
	mu            sync.RWMutex
	snapshots     map[string]*tracecore_types.FederationSnapshotDTO
}

func NewWorkspaceHandler(
	createUC *workspace_usecase.CreateWorkspaceUsecase,
	listUC *workspace_usecase.ListWorkspaceUsecase,
) *WorkspaceHandler {
	return &WorkspaceHandler{
		createUseCase: createUC,
		listUseCase:   listUC,
		snapshots:     make(map[string]*tracecore_types.FederationSnapshotDTO),
	}
}

func (h *WorkspaceHandler) GetWorkspaceFederation(ctx context.Context, workspaceID string) (*tracecore_types.FederationSnapshotDTO, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if snap, found := h.snapshots[workspaceID]; found {
		return snap, nil
	}

	// Initial default snapshot for workspace
	initialSnap := &tracecore_types.FederationSnapshotDTO{
		WorkspaceID: workspaceID,
		RemoteVaults: []tracecore_types.RemoteVaultDTO{
			{
				ID:           "vault:supplier-x",
				Endpoint:     "https://afp.supplier-x.io/v1/sync",
				Status:       "trusted",
				LastSeen:     "2 min ago",
				Cursor:       "1247",
				Proto:        "AFP v1.2",
				PendingItems: "4 items",
				Alert:        true,
			},
			{
				ID:           "vault:ey-auditors",
				Endpoint:     "https://afp.ey.com/ankhora/v1/sync",
				Status:       "trusted",
				LastSeen:     "1 hr ago",
				Cursor:       "892",
				Proto:        "AFP v1.1",
				PendingItems: "0 pending",
				Alert:        false,
			},
			{
				ID:           "vault:partner-bank",
				Endpoint:     "https://afp.partnerbank.io/sync",
				Status:       "pending",
				LastSeen:     "1 day ago",
				Cursor:       "0",
				Proto:        "AFP v1.0",
				PendingItems: "2 pending",
				Alert:        false,
			},
		},
	}
	h.snapshots[workspaceID] = initialSnap
	return initialSnap, nil
}

func (h *WorkspaceHandler) AddRemoteVaultToWorkspace(ctx context.Context, workspaceID string, remoteVault tracecore_types.RemoteVaultDTO) (*tracecore_types.FederationSnapshotDTO, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	snap, found := h.snapshots[workspaceID]
	if !found {
		snap = &tracecore_types.FederationSnapshotDTO{
			WorkspaceID:  workspaceID,
			RemoteVaults: make([]tracecore_types.RemoteVaultDTO, 0),
		}
		h.snapshots[workspaceID] = snap
	}

	snap.RemoteVaults = append(snap.RemoteVaults, remoteVault)
	return snap, nil
}

func (h *WorkspaceHandler) CreateWorkspace(ctx context.Context,  userID string,vaultId string, name string, description string) (*tracecore_types.Workspace, error) {
	if h.createUseCase == nil {
		return nil, fmt.Errorf("create workspace use case is not initialized")
	}

	req := &workspace_application.CreateWorkspaceRequest{
		VaultID:     vaultId,
		OwnerID:     userID,
		Name:        name,
		Description: description,
		Signature:   "desktop_authenticated",
	}

	ws, err := h.createUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	return toTracecoreWorkspace(ws), nil
}

func (h *WorkspaceHandler) ListWorkspaces(ctx context.Context,vaultID string) ([]tracecore_types.Workspace, error) {
	if h.listUseCase == nil {
		return nil, fmt.Errorf("list workspace use case is not initialized")
	}

	req := &workspace_application.ListWorkspacesRequest{
		VaultID: vaultID,
	}

	workspaces, err := h.listUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	res := make([]tracecore_types.Workspace, 0, len(workspaces))
	for _, ws := range workspaces {
		res = append(res, *toTracecoreWorkspace(&ws))
	}
	utils.LogPretty("[Workspace] WorkspaceHandler.ListWorkspaces result", res)

	return res, nil
}

func toTracecoreWorkspace(ws *workspace_domain.Workspace) *tracecore_types.Workspace {
	if ws == nil {
		return nil
	}
	return &tracecore_types.Workspace{
		ID:          ws.ID,
		VaultID:     ws.VaultID,
		Name:        ws.Name,
		Description: ws.Description,
		Status:      string(ws.Status),
		OwnerID:     ws.OwnerID,
		CreatedAt:   ws.CreatedAt,
		UpdatedAt:   ws.UpdatedAt,
		IsDraft:     ws.IsDraft,
		IsDirty:     ws.IsDirty,
	}
}

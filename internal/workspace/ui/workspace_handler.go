package workspace_ui

import (
	"context"
	"fmt"

	tracecore_types "vault-app/internal/tracecore/types"
	"vault-app/internal/utils"
	workspace_application "vault-app/internal/workspace/application"
	workspace_usecase "vault-app/internal/workspace/application/usecases"
	workspace_domain "vault-app/internal/workspace/domain"
)

type WorkspaceHandler struct {
	createUseCase *workspace_usecase.CreateWorkspaceUsecase
	listUseCase   *workspace_usecase.ListWorkspaceUsecase
}

func NewWorkspaceHandler(
	createUC *workspace_usecase.CreateWorkspaceUsecase,
	listUC *workspace_usecase.ListWorkspaceUsecase,
) *WorkspaceHandler {
	return &WorkspaceHandler{
		createUseCase: createUC,
		listUseCase:   listUC,
	}
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

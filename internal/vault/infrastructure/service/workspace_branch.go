package vaults_service

import (
	"context"
	"encoding/json"
	"time"

	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
)

type WorkspaceNode struct {
	ID      string
	VaultID string

	Name        string
	Description string
	Status      vaults_domain.WorkspaceStatus

	OwnerID string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// =======================================================================================
// WRITE
// =======================================================================================
func (s *VaultService) BuildWorkspacesBranch(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, error) {
	// =========================
	// 1. BUILD ENTRIES
	// =========================
	workspaces := vp.Collaborative.Workspaces
	workspaceLinks, err := s.BuildWorkspaces(workspaces)
	if err != nil {
		return "", err
	}

	// =========================
	// 4. ENTRIES ROOT
	// =========================
	workspacesCID, _, err := s.BuildWorkspacesRoot(workspaceLinks)
	if err != nil {
		return "", err
	}

	return workspacesCID, nil
}

func (s *VaultService) BuildWorkspaces(workspaces []vaults_domain.Workspace) ([]vaults_domain.Link, error) {
	var links []vaults_domain.Link

	for _, workspace := range workspaces {

		node := WorkspaceNode{
			ID:          workspace.ID,
			VaultID:     workspace.VaultID,
			Name:        workspace.Name,
			Description: workspace.Description,
			Status:      workspace.Status,
			OwnerID:     workspace.OwnerID,
			CreatedAt:   workspace.CreatedAt,
			UpdatedAt:   workspace.UpdatedAt,
		}

		cid, _, err := s.putNode(node)
		if err != nil {
			return nil, err
		}

		links = append(links, vaults_domain.Link{CID: cid})
	}

	return links, nil
}

func (s *VaultService) BuildWorkspacesRoot(links []vaults_domain.Link) (string, int, error) {
	root := vaults_domain.WorkspacesRoot{Items: links}
	return s.putNode(root)
}

func (s *VaultService) RotateWorkspaceGraph(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, error) {

	for i := range vp.Collaborative.Workspaces {
		vp.Collaborative.Workspaces[i].IsDirty = true
	}
	// 	↓
	cid, err := s.BuildWorkspacesBranch(session, vp, mode)
	if err != nil {
		return "", err
	}
	return cid, nil
}



// =======================================================================================
// READ
// =======================================================================================
func (r *VaultReconstructor) resolveWorkspaces(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	workspacesRoot vaults_domain.WorkspacesRoot,
) ([]vaults_domain.Workspace, error) {

	var result []vaults_domain.Workspace

	for _, link := range workspacesRoot.Items {

		res, err := r.Query.Execute(ctx, cmd.WithCID(link.CID))
		if err != nil {
			return result, err
		}

		var workspace vaults_domain.Workspace
		if err := json.Unmarshal(res.Raw, &workspace); err != nil {
			return result, err
		}

		result = append(result, workspace)
	}

	return result, nil
}

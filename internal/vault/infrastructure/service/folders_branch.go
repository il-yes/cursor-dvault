package vaults_service

import (
	"context"
	"encoding/json"

	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
)

type FolderNode struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// =======================================================================================
// WRITE
// =======================================================================================
func (s *VaultService) BuildFoldersBranch(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, error) {
	// =========================
	// 1. BUILD ENTRIES
	// =========================
	folders := vp.Personal.Folders
	folderLinks, err := s.BuildFolders(folders)
	if err != nil {
		return "", err
	}

	// =========================
	// 4. ENTRIES ROOT
	// =========================
	foldersCID, _, err := s.BuildFoldersRoot(folderLinks)
	if err != nil {
		return "", err
	}

	return foldersCID, nil
}

func (s *VaultService) BuildFolders(folders []vaults_domain.Folder) ([]vaults_domain.Link, error) {
	var links []vaults_domain.Link

	for _, folder := range folders {
		if folder.IsDraft {
			continue
		}

		node := FolderNode{
			ID:   folder.ID,
			Name: folder.Name,
		}

		cid, _, err := s.putNode(node)
		if err != nil {
			return nil, err
		}

		links = append(links, vaults_domain.Link{CID: cid})
	}

	return links, nil
}

func (s *VaultService) BuildFoldersRoot(links []vaults_domain.Link) (string, int, error) {
	root := vaults_domain.FoldersRoot{Items: links}
	return s.putNode(root)
}

func (s *VaultService) RotateFolderGraph(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, error) {

	for i := range vp.Folders {
		vp.Folders[i].IsDirty = true
	}
	// 	↓
	cid, err := s.BuildFoldersBranch(session, vp, mode)
	if err != nil {
		return "", err
	}
	return cid, nil
}

// =======================================================================================
// READ
// =======================================================================================
func (r *VaultReconstructor) resolveFolders(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	foldersRoot vaults_domain.FoldersRoot,
) ([]vaults_domain.Folder, error) {

	var result []vaults_domain.Folder

	for _, link := range foldersRoot.Items {

		res, err := r.Query.Execute(ctx, cmd.WithCID(link.CID))
		if err != nil {
			return result, err
		}

		var folder vaults_domain.Folder
		if err := json.Unmarshal(res.Raw, &folder); err != nil {
			return result, err
		}

		result = append(result, folder)
	}

	return result, nil
}

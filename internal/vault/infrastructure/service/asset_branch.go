package vaults_service

import (
	"context"
	"encoding/json"

	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
)

type AssetNode struct {
	CID         string
	Size        int64
	ContentHash string
	Type string

	MediaType string
}

// =======================================================================================
// WRITE
// =======================================================================================
func (s *VaultService) BuildAssetsBranch(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, map[string][]vaults_domain.Link, map[string][]vaults_domain.Link, error) {
	// =========================
	// 1. BUILD ENTRIES
	// =========================
	assets := vp.Collaborative.Assets
	assetLinks, byHash, byType, err := s.BuildAssets(assets)
	if err != nil {
		return "", nil, nil,  err
	}

	// =========================
	// 4. ENTRIES ROOT
	// =========================
	assetsCID, _, err := s.BuildAssetsRoot(assetLinks)
	if err != nil {
		return "", nil, nil, err
	}

	return assetsCID, byHash, byType, nil
}

func (s *VaultService) BuildAssets(assets []vaults_domain.Asset) ([]vaults_domain.Link, map[string][]vaults_domain.Link, map[string][]vaults_domain.Link, error) {
	var links []vaults_domain.Link

	byHash := make(map[string][]vaults_domain.Link)
	byType := make(map[string][]vaults_domain.Link)

	addLink := func(base vaults_domain.Asset, cid string) {
		link := vaults_domain.Link{CID: cid}
		links = append(links, link)

		byHash[string(base.ContentHash)] = append(byHash[string(base.ContentHash)], link)

		if base.Type != "" {
			byType[base.Type] = append(byType[base.Type], link)
		}
	}

	for _, asset := range assets {

		node := AssetNode{
			CID:         asset.CID,
			Size:        asset.Size,
			ContentHash: asset.ContentHash,
		}

		cid, _, err := s.putNode(node)
		if err != nil {
			return nil, nil, nil, err
		}

		
		addLink(asset, cid)
	}

	return links, byHash, byType, nil
}

func (s *VaultService) BuildAssetsRoot(links []vaults_domain.Link) (string, int, error) {
	root := vaults_domain.AssetsRoot{Items: links}
	return s.putNode(root)
}

func (s *VaultService) RotateAssetGraph(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, map[string][]vaults_domain.Link, map[string][]vaults_domain.Link, error) {

	for i := range vp.Collaborative.Assets {
		vp.Collaborative.Assets[i].IsDirty = true
	}
	// 	↓
	cid, byHash, byType,  err := s.BuildAssetsBranch(session, vp, mode)
	if err != nil {
		return "", nil, nil, err
	}
	return cid, byHash, byType, nil
}

// =======================================================================================
// READ
// =======================================================================================
func (r *VaultReconstructor) resolveAssets(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	assetsRoot vaults_domain.AssetsRoot,
) ([]vaults_domain.Asset, error) {

	var result []vaults_domain.Asset

	for _, link := range assetsRoot.Items {

		res, err := r.Query.Execute(ctx, cmd.WithCID(link.CID))
		if err != nil {
			return result, err
		}
		// utils.LogPretty("VaultReconstructor - resolveAssets - res", res)

		var asset vaults_domain.Asset
		if err := json.Unmarshal(res.Raw, &asset); err != nil {
			return result, err
		}

		result = append(result, asset)
	}

	return result, nil
}

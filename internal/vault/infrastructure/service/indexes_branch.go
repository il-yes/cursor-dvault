package vaults_service

import (
	"context"
	"encoding/json"

	vault_queries "vault-app/internal/vault/application/queries"
	vaults_domain "vault-app/internal/vault/domain"
)

// =======================================================================================
// WRITE
// =======================================================================================
// Personnnal
func (s *VaultService) buildIndex(byType, byFolder map[string][]vaults_domain.Link) (string, int, error) {
	index := vaults_domain.Index{
		ByType:   byType,
		ByFolder: byFolder,
	}
	return s.putNode(index)
}

// - Collaborative
func (s *VaultService) buildCollaborativeIndex(
	threadCID string,
	assetCID string,
	federationCID string,
	trustGroupCID string,
) (string, int, error) {

	root := vaults_domain.CollaborativeIndexRoot{
		ThreadsIndex: vaults_domain.Link{
			CID: threadCID,
		},

		AssetsIndex: vaults_domain.Link{
			CID: assetCID,
		},

		FederationsIndex: vaults_domain.Link{
			CID: federationCID,
		},

		TrustGroupsIndex: vaults_domain.Link{
			CID: trustGroupCID,
		},
	}

	return s.putNode(root)
}

func (s *VaultService) buildThreadIndex(byChannel map[string][]vaults_domain.Link, byStatus map[string][]vaults_domain.Link) (string, int, error) {
	index := vaults_domain.ThreadsIndex{
		ByChannel: byChannel,
		ByStatus:  byStatus,
	}
	return s.putNode(index)
}
func (s *VaultService) buildAssetIndex(byChannel, byStatus map[string][]vaults_domain.Link) (string, int, error) {
	index := vaults_domain.AssetsIndex{
		ByHash: byChannel,
		ByType: byStatus,
	}
	return s.putNode(index)
}
func (s *VaultService) buildTrustGroupIndex(byChannel, byStatus map[string][]vaults_domain.Link) (string, int, error) {
	index := vaults_domain.TrustGroupsIndex{
		ByWorkspace: byChannel,
		ByMember:    byStatus,
	}
	return s.putNode(index)
}

	func (s *VaultService) buildFederationIndex(bv map[string][]vaults_domain.Link, bts map[string][]vaults_domain.Link) (string, int, error) {
		index := vaults_domain.FederationsIndex{
			ByVault: bv,
			ByTrustState:      bts,
		}
		return s.putNode(index)
	}
// func (s *VaultService) BuildFederationIndex(vaults []vaults_domain.RemoteVault) (string, int, error) {

// 	index := vaults_domain.FederationsIndex{
// 		ByVault:      make(map[string]vaults_domain.Link),
// 		ByTrustState: make(map[string][]vaults_domain.Link),
// 	}

// 	for _, remote := range vaults {
// 		// temporary reference,
// 		// populated later if you keep CID cache
// 		index.ByTrustState[string(remote.TrustState)] = append(index.ByTrustState[string(remote.TrustState)], vaults_domain.Link{})

// 	}

// 	return s.putNode(index)
// }

// =======================================================================================
// READ
// =======================================================================================
func (r *VaultReconstructor) resolveIndexC3s(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	indexC3sRoot vaults_domain.CollaborativeIndexRoot,
) (*vaults_domain.IndexC3, error) {

	threadRes, err := r.Query.Execute(ctx, cmd.WithCID(indexC3sRoot.ThreadsIndex.CID))
	if err != nil {
		return nil, err
	}
	var threadIndex vaults_domain.ThreadsIndex
	if err := json.Unmarshal(threadRes.Raw, &threadIndex); err != nil {
		return nil, err
	}

	assetRes, err := r.Query.Execute(ctx, cmd.WithCID(indexC3sRoot.AssetsIndex.CID))
	if err != nil {
		return nil, err
	}
	var assetIndex vaults_domain.AssetsIndex
	if err := json.Unmarshal(assetRes.Raw, &assetIndex); err != nil {
		return nil, err
	}

	federationRes, err := r.Query.Execute(ctx, cmd.WithCID(indexC3sRoot.FederationsIndex.CID))
	if err != nil {
		return nil, err
	}
	var federationIndex vaults_domain.FederationsIndex
	if err := json.Unmarshal(federationRes.Raw, &federationIndex); err != nil {
		return nil, err
	}

	trustGroupRes, err := r.Query.Execute(ctx, cmd.WithCID(indexC3sRoot.TrustGroupsIndex.CID))
	if err != nil {
		return nil, err
	}
	var trustGroupIndex vaults_domain.TrustGroupsIndex
	if err := json.Unmarshal(trustGroupRes.Raw, &trustGroupIndex); err != nil {
		return nil, err
	}

	index := vaults_domain.IndexC3{
		Threads:     threadIndex,
		Assets:      assetIndex,
		Federations: federationIndex,
		TrustGroups: trustGroupIndex,
	}

	return &index, nil
}

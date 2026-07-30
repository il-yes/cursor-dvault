package vaults_service

import vaults_domain "vault-app/internal/vault/domain"

func (s *VaultService) buildIndex(byType, byFolder map[string][]vaults_domain.Link) (string, int, error) {
	index := vaults_domain.Index{
		ByType:   byType,
		ByFolder: byFolder,
	}
	return s.putNode(index)
}

func (s *VaultService) buildCollaborativeIndex(
	indexThread vaults_domain.ThreadsIndex,
	indexAsset vaults_domain.AssetsIndex,
	indexFederation vaults_domain.FederationsIndex,
	indexTrustGroup vaults_domain.TrustGroupsIndex,
) (string, int, error) {

	index := vaults_domain.IndexC3{
		Thread:     indexThread,
		Asset:      indexAsset,
		Federation: indexFederation,
		TrustGroup: indexTrustGroup,
	}
	return s.putNode(index)
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

//	func (s *VaultService) buildFederationIndex(byChannel, byStatus map[string][]vaults_domain.Link) (string, int, error) {
//		index := vaults_domain.FederationsIndex{
//			ByRemoteVault: byChannel,
//			ByStatus:      byStatus,
//		}
//		return s.putNode(index)
//	}
func (s *VaultService) BuildFederationIndex(vaults []vaults_domain.RemoteVault) (string, int, error) {

	index := vaults_domain.FederationsIndex{
		ByVault:      make(map[string]vaults_domain.Link),
		ByTrustState: make(map[string][]vaults_domain.Link),
	}

	for _, remote := range vaults {

		// temporary reference,
		// populated later if you keep CID cache
		index.ByTrustState[string(remote.TrustState)] = append(index.ByTrustState[string(remote.TrustState)], vaults_domain.Link{})

	}

	return s.putNode(index)
}

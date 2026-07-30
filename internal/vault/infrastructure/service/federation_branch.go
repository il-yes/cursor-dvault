package vaults_service

import (
	"context"
	"encoding/json"
	"time"

	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
)

// ============================================================
// FEDERATION GRAPH NODES
// ============================================================

type FederationNode struct {
	Version string

	RemoteVaults vaults_domain.Link
	Index        vaults_domain.Link
}

type RemoteVaultNode struct {
	Version string

	VaultID    string
	LastCursor uint64
	LastSeen   time.Time
	TrustState vaults_domain.TrustState

	PendingSync vaults_domain.Link
}

type PendingSyncNode struct {
	Version string

	Items []vaults_domain.PendingSyncItem
}

// =======================================================================================
// WRITE
// =======================================================================================
func (s *VaultService) BuildFederationBranch(
	session vault_session.Session,
	vp vaults_domain.VaultPayload,
	mode SyncMode,
) (string, error) {

	federation := vp.Collaborative.Federation

	// =========================
	// 1. BUILD REMOTE VAULTS
	// =========================

	remoteVaultLinks, err := s.BuildRemoteVaults(
		federation.RemoteVaults,
	)

	if err != nil {
		return "", err
	}

	// =========================
	// 2. BUILD REMOTE VAULT ROOT
	// =========================

	remoteVaultRootCID, _, err := s.BuildRemoteVaultsRoot(
		remoteVaultLinks,
	)

	if err != nil {
		return "", err
	}

	// =========================
	// 3. BUILD INDEX
	// =========================

	indexCID, _, err := s.BuildFederationIndex(
		federation.RemoteVaults,
	)

	if err != nil {
		return "", err
	}

	// =========================
	// 4. BUILD FEDERATION ROOT
	// =========================

	node := FederationNode{
		Version: "1.0",
		RemoteVaults: vaults_domain.Link{
			CID: remoteVaultRootCID,
		},
		Index: vaults_domain.Link{
			CID: indexCID,
		},
	}

	cid, _, err := s.putNode(node)

	if err != nil {
		return "", err
	}

	return cid, nil
}

func (s *VaultService) BuildRemoteVaults(
	vaults []vaults_domain.RemoteVault,
) ([]vaults_domain.Link, error) {

	var links []vaults_domain.Link

	for _, remote := range vaults {

		pendingCID, err := s.BuildPendingSync(
			remote.Pending,
		)

		if err != nil {
			return nil, err
		}

		node := RemoteVaultNode{

			Version: "1.0",

			VaultID: remote.VaultID,

			LastCursor: remote.LastCursor,

			LastSeen: remote.LastSeen,

			TrustState: remote.TrustState,

			PendingSync: vaults_domain.Link{
				CID: pendingCID,
			},
		}

		cid, _, err := s.putNode(node)

		if err != nil {
			return nil, err
		}

		links = append(
			links,
			vaults_domain.Link{
				CID: cid,
			},
		)

	}

	return links, nil
}

func (s *VaultService) BuildRemoteVaultsRoot(links []vaults_domain.Link) (string, int, error) {

	root := vaults_domain.RemoteVaultsRoot{
		Items: links,
	}

	return s.putNode(root)
}

func (s *VaultService) BuildPendingSync(items []vaults_domain.PendingSyncItem) (string, error) {

	node := PendingSyncNode{
		Version: "1.0",
		Items:   items,
	}

	cid, _, err := s.putNode(node)

	if err != nil {
		return "", err
	}

	return cid, nil
}

func (s *VaultService) RotateFederationBranch(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, error) {
	vp.Collaborative.Federation.IsDirty = true

	return s.BuildFederationBranch(session, vp, mode)
}




// =======================================================================================
// READ
// =======================================================================================
func (r *VaultReconstructor) resolveFederations(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	federationsRoot vaults_domain.FederationRoot,
) ([]vaults_domain.FederationConfig, error) {

	var result []vaults_domain.FederationConfig

	for _, link := range federationsRoot.Items {

		res, err := r.Query.Execute(ctx, cmd.WithCID(link.CID))
		if err != nil {
			return result, err
		}

		var federation vaults_domain.FederationConfig
		if err := json.Unmarshal(res.Raw, &federation); err != nil {
			return result, err
		}

		result = append(result, federation)
	}

	return result, nil
}
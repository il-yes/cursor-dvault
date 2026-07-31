package vaults_service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"vault-app/internal/utils"
	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
)

// ============================================================
// FEDERATION GRAPH NODES
// ============================================================

// type FederationNode struct {
// 	Version string `json:"version"`

//		RemoteVaults vaults_domain.Link `json:"remote_vaults"`
//		Index        vaults_domain.Link
//	}
type FederationNode struct {
	Version      string             `json:"version"`
	RemoteVaults vaults_domain.Link `json:"remote_vaults"`
	Index        vaults_domain.Link `json:"index"`
}

type RemoteVaultNode struct {
	Version string `json:"version"`

	VaultID    string                   `json:"vault_id"`
	LastCursor uint64                   `json:"last_cursor"`
	LastSeen   time.Time                `json:"last_seen"`
	TrustState vaults_domain.TrustState `json:"trust_state"`

	PendingSync PendingSyncNode `json:"pending"`
}

type PendingSyncNode struct {
	Version string `json:"version"`

	Items []vaults_domain.PendingSyncItem `json:"items"`
}

// =======================================================================================
// WRITE
// =======================================================================================
func (s *VaultService) BuildFederationBranch(
	session vault_session.Session,
	vp vaults_domain.VaultPayload,
	mode SyncMode,
) (string, string, error) {

	federation := vp.Collaborative.Federation

	// =========================
	// 1. BUILD REMOTE VAULTS
	// =========================

	remoteVaultLinks, byVault, byTrustState, err := s.BuildRemoteVaults(
		federation.RemoteVaults,
	)

	if err != nil {
		return "", "", err
	}

	// =========================
	// 2. BUILD REMOTE VAULT ROOT
	// =========================

	remoteVaultRootCID, _, err := s.BuildRemoteVaultsRoot(
		remoteVaultLinks,
	)

	if err != nil {
		return "", "", err
	}

	// =========================
	// 3. BUILD INDEX
	// =========================

	indexCID, _, err := s.buildFederationIndex(byVault, byTrustState)

	if err != nil {
		return "", "", err
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
		return "", "", err
	}

	return cid, indexCID, nil
}

func (s *VaultService) BuildRemoteVaults(vaults []vaults_domain.RemoteVault) ([]vaults_domain.Link, map[string][]vaults_domain.Link, map[string][]vaults_domain.Link, error) {

	var links []vaults_domain.Link
	ByVault := make(map[string][]vaults_domain.Link)
	ByTrustState := make(map[string][]vaults_domain.Link)

	addLink := func(base vaults_domain.RemoteVault, cid string) {
		link := vaults_domain.Link{CID: cid}
		links = append(links, link)

		ByVault[string(base.VaultID)] = append(ByVault[string(base.VaultID)], link)

		if base.TrustState != "" {
			ByTrustState[string(base.TrustState)] = append(ByTrustState[string(base.TrustState)], link)
		}
	}

	for _, remote := range vaults {

		// pendingCID, err := s.BuildPendingSync(
		// 	remote.Pending,
		// )

		// if err != nil {
		// 	return nil, err
		// }

		node := RemoteVaultNode{
			Version:     "1.0",
			VaultID:     remote.VaultID,
			LastCursor:  remote.LastCursor,
			LastSeen:    remote.LastSeen,
			TrustState:  remote.TrustState,
			PendingSync: PendingSyncNode{Items: remote.Pending},
		}

		cid, _, err := s.putNode(node)

		if err != nil {
			return nil, nil, nil, err
		}

		addLink(remote, cid)

	}

	return links, ByVault, ByTrustState, nil
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

func (s *VaultService) RotateFederationBranch(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, string, error) {
	vp.Collaborative.Federation.IsDirty = true

	return s.BuildFederationBranch(session, vp, mode)
}

// =======================================================================================
// READ
// =======================================================================================
func (r *VaultReconstructor) resolveFederation(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	node FederationNode,
) (vaults_domain.FederationSnapshot, error) {
	var snapshot vaults_domain.FederationSnapshot

	// -----------------------------
	// Load FederationNode
	// -----------------------------
	// res, err := r.Query.Execute(ctx, cmd.WithCID(remoteVaultLink.CID))
	// if err != nil {
	// 	return snapshot, err
	// }

	// utils.LogPretty(
	// 	"FEDERATION RAW BEFORE UNMARSHAL",
	// 	string(res.Raw),
	// )

	// var node FederationNode
	// if err := json.Unmarshal(res.Raw, &node); err != nil {
	// 	return snapshot, err
	// }

	fmt.Printf("FederationNode type: %+v\n", node)
	utils.LogPretty(
		"FEDERATION NODE",
		node,
	)
	fmt.Printf(
		"REMOTE CID AFTER UNMARSHAL = '%s'\n",
		node.RemoteVaults.CID,
	)

	utils.LogPretty(
		"FEDERATION REMOTE VAULT LINK",
		node.RemoteVaults,
	)
	// -----------------------------
	// Resolve Remote Vaults
	// -----------------------------
	remoteVaults, err := r.resolveRemoteVaults(
		ctx,
		cmd,
		node.RemoteVaults,
	)
	if err != nil {
		return snapshot, err
	}

	snapshot.RemoteVaults = remoteVaults

	return snapshot, nil
}

func (r *VaultReconstructor) resolveRemoteVaults(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	root vaults_domain.Link,
) ([]vaults_domain.RemoteVault, error) {

	var result []vaults_domain.RemoteVault

	// -----------------------------
	// Load RemoteVaultRoot
	// -----------------------------
	res, err := r.Query.Execute(ctx, cmd.WithCID(root.CID))
	if err != nil {
		return result, err
	}

	var remoteRoot vaults_domain.RemoteVaultsRoot
	if err := json.Unmarshal(res.Raw, &remoteRoot); err != nil {
		return result, err
	}

	// -----------------------------
	// Load each RemoteVaultNode
	// -----------------------------
	for _, link := range remoteRoot.Items {

		vaultRes, err := r.Query.Execute(ctx, cmd.WithCID(link.CID))
		if err != nil {
			return result, err
		}

		var node RemoteVaultNode
		if err := json.Unmarshal(vaultRes.Raw, &node); err != nil {
			return result, err
		}
		utils.LogPretty(
			"REMOTE VAULT ROOT RAW",
			string(res.Raw),
		)

		utils.LogPretty(
			"REMOTE VAULT ROOT",
			remoteRoot,
		)

		fmt.Println(
			"REMOTE VAULT COUNT:",
			len(remoteRoot.Items),
		)

		result = append(result, vaults_domain.RemoteVault{
			VaultID:         node.VaultID,
			LastCursor:      node.LastCursor,
			LastSeen:        node.LastSeen,
			TrustState:      node.TrustState,
			Pending:         node.PendingSync.Items,
			ProtocolVersion: node.Version,
		})
	}

	return result, nil
}

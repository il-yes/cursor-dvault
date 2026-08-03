package vaults_service

import (
	"context"
	"encoding/json"
	"time"

	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
)

type TrustMemberNode struct {
	ID      string
	VaultID string
	Role    string
	KeyVersion uint64
	WrappedKEK string
	JoinedAt time.Time
}

type TrustMemberRequest struct {
	session vault_session.Session
	vp      vaults_domain.VaultPayload
	mode    SyncMode
	Members vaults_domain.Link
}

// =======================================================================================
// WRITE
// =======================================================================================
func (s *VaultService) BuildTrustMembersBranch(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, error) {
	// =========================
	// 1. BUILD TRUSTGROUP
	// =========================
	tgs := vp.Collaborative.TrustGroupMembers
	tgsLinks, err := s.BuildTrustMembers(tgs)
	if err != nil {
		return "", err
	}

	// =========================
	// 4. TRUSTGROUP ROOT
	// =========================
	trustGroupCID, _, err := s.BuildTrustMembersRoot(tgsLinks)
	if err != nil {
		return "", err
	}

	return trustGroupCID, nil
}

func (s *VaultService) BuildTrustMembers(groups []vaults_domain.TrustGroupMember) ([]vaults_domain.Link, error) {
	var links []vaults_domain.Link

	for _, group := range groups {
		if group.IsDraft {
			continue
		}

		node := TrustMemberNode{
			ID:      group.ID,
			VaultID: group.VaultID,
			Role:    group.Role,
			KeyVersion: group.WrappedKEK.Version,
			WrappedKEK: group.WrappedKEK.Value,
			JoinedAt: group.JoinedAt,
		}

		cid, _, err := s.putNode(node)
		if err != nil {
			return nil, err
		}

		links = append(links, vaults_domain.Link{CID: cid})
	}

	return links, nil
}

func (s *VaultService) BuildTrustMembersRoot(links []vaults_domain.Link) (string, int, error) {
	root := vaults_domain.TrustGroupMembersRoot{Items: links}
	return s.putNode(root)
}

func (s *VaultService) RotateTrustMembersGraph(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, error) {
	for i := range vp.Collaborative.TrustGroupMembers {
		vp.Collaborative.TrustGroupMembers[i].IsDirty = true
	}
	// 	↓
	cid, err := s.BuildTrustMembersBranch(session, vp, mode)
	if err != nil {
		return "", err
	}
	return cid, nil
}


// =======================================================================================
// READ
// =======================================================================================
func (r *VaultReconstructor) resolveTrustGroupMembers(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	trustGroupMembersRoot vaults_domain.TrustGroupMembersRoot,
) ([]vaults_domain.TrustGroupMember, error) {

	var result []vaults_domain.TrustGroupMember

	for _, link := range trustGroupMembersRoot.Items {

		res, err := r.Query.Execute(ctx, cmd.WithCID(link.CID))
		if err != nil {
			return result, err
		}
		// utils.LogPretty("VaultReconstructor - resolveTrustGroupMembers - res", res)

		var trustGroupMember vaults_domain.TrustGroupMember
		if err := json.Unmarshal(res.Raw, &trustGroupMember); err != nil {
			return result, err
		}

		result = append(result, trustGroupMember)
	}

	return result, nil
}

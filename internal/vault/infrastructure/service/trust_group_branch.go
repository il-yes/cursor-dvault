package vaults_service

import (
	"context"
	"encoding/json"

	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
)

type TrustGroupNode struct {
	ID                string
	TrustGroupID      string
	Name              string
	CurrentKEKVersion int64
	MembersCIDs       string
	UpdatedAt         string
}

type TrustGroupRequest struct {
	session vault_session.Session
	vp      vaults_domain.VaultPayload
	mode    SyncMode
}

// =======================================================================================
// WRITE
// =======================================================================================
func (s *VaultService) BuildTrustGroupsBranch(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, map[string][]vaults_domain.Link, map[string][]vaults_domain.Link, error) {
	// =========================
	// 1. BUILD TRUSTGROUP
	// =========================
	tgs := vp.Collaborative.TrustGroups
	tgsLinks, byWorkspace, byMember, err := s.BuildTrustGroups(tgs)
	if err != nil {
		return "", nil, nil, err
	}

	// =========================
	// 4. TRUSTGROUP ROOT
	// =========================
	trustGroupCID, _, err := s.BuildTrustGroupsRoot(tgsLinks)
	if err != nil {
		return "", nil, nil, err
	}

	return trustGroupCID, byWorkspace, byMember, nil
}

func (s *VaultService) BuildTrustGroups(groups []vaults_domain.TrustGroup) ([]vaults_domain.Link, map[string][]vaults_domain.Link, map[string][]vaults_domain.Link,  error) {
	var links []vaults_domain.Link

	byWorkspace := make(map[string][]vaults_domain.Link)
	byMember := make(map[string][]vaults_domain.Link)

	addLink := func(base vaults_domain.TrustGroup, cid string) {
		link := vaults_domain.Link{CID: cid}
		links = append(links, link)

		byWorkspace[string(base.WorkspaceID)] = append(byWorkspace[string(base.WorkspaceID)], link)

		if base.MemberCIDs != "" {
			byMember[base.MemberCIDs] = append(byMember[base.MemberCIDs], link)
		}
	}


	for _, group := range groups {
		if group.IsDraft {
			continue
		}

		node := TrustGroupNode{
			ID:                group.ID,
			TrustGroupID:      group.ID,
			Name:              group.Name,
			CurrentKEKVersion: int64(group.KEKVersion),
			MembersCIDs:       group.MemberCIDs, // TODO wrong, select the right member for the right trustgroup

		}

		cid, _, err := s.putNode(node)
		if err != nil {
			return nil, nil, nil, err
		}

		addLink(group, cid)
	}

	return links, byWorkspace, byMember, nil
}

func (s *VaultService) BuildTrustGroupsRoot(links []vaults_domain.Link) (string, int, error) {
	root := vaults_domain.TrustGroupsRoot{Items: links}
	return s.putNode(root)
}

func (s *VaultService) RotateTrustGroupsKeys(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, map[string][]vaults_domain.Link, map[string][]vaults_domain.Link, error) {
	for i := range vp.Collaborative.TrustGroups {
		vp.Collaborative.TrustGroups[i].IsDirty = true
	}
	// 	↓
	cid, byWorkspace, byMember, err := s.BuildTrustGroupsBranch(session, vp, mode)
	if err != nil {
		return "", nil, nil, err
	}
	return cid, byWorkspace, byMember, nil
}

// =======================================================================================
// READ
// =======================================================================================
func (r *VaultReconstructor) resolveTrustGroups(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	trustgroupsRoot vaults_domain.TrustGroupsRoot,
) ([]vaults_domain.TrustGroup, error) {

	var result []vaults_domain.TrustGroup

	for _, link := range trustgroupsRoot.Items {

		res, err := r.Query.Execute(ctx, cmd.WithCID(link.CID))
		if err != nil {
			return result, err
		}

		var trustgroup vaults_domain.TrustGroup
		if err := json.Unmarshal(res.Raw, &trustgroup); err != nil {
			return result, err
		}

		result = append(result, trustgroup)
	}

	return result, nil
}

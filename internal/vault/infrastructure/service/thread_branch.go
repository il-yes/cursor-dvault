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

type ThreadNode struct {
	ID        string
	ChannelID string

	AssetType string

	Title    string
	Subtitle string

	Status vaults_domain.ThreadStatus

	CreatedAt time.Time
	ClosedAt  *time.Time
}

// =======================================================================================
// WRITE
// =======================================================================================
func (s *VaultService) BuildThreadsBranch(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, map[string][]vaults_domain.Link, map[string][]vaults_domain.Link, error) {

	threads := vp.Collaborative.Threads

	threadLinks, indexByChannel, indexByStatus, err := s.BuildThreads(threads)
	if err != nil {
		return "", nil, nil, err
	}

	threadsCID, _, err := s.BuildThreadsRoot(threadLinks)
	if err != nil {
		return "", nil, nil, err
	}

	return threadsCID, indexByChannel, indexByStatus, nil
}

func (s *VaultService) BuildThreads(threads []vaults_domain.Thread) (
	[]vaults_domain.Link,
	map[string][]vaults_domain.Link,
	map[string][]vaults_domain.Link,
	error,
) {

	var links []vaults_domain.Link

	byChannel := make(map[string][]vaults_domain.Link)
	byStatus := make(map[string][]vaults_domain.Link)

	addLink := func(base vaults_domain.Thread, cid string) {
		link := vaults_domain.Link{CID: cid}
		links = append(links, link)

		byStatus[string(base.Status)] = append(byStatus[string(base.Status)], link)

		if base.ChannelID != "" {
			byChannel[base.ChannelID] = append(byChannel[base.ChannelID], link)
		}
	}

	for _, thread := range threads {

		if thread.IsDraft {
			continue
		}

		node := ThreadNode{
			ID:        thread.ID,
			ChannelID: thread.ChannelID,

			AssetType: thread.AssetType,

			Title:    thread.Title,
			Subtitle: thread.Subtitle,

			Status: thread.Status,

			CreatedAt: thread.CreatedAt,
			ClosedAt:  thread.ClosedAt,
		}

		cid, _, err := s.putNode(node)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("VaultService - BuildThreads - failed to provide ipfs cid: %v", err)
		}

		addLink(thread, cid)
	}

	return links, byChannel, byStatus, nil
}

func (s *VaultService) BuildThreadsRoot(links []vaults_domain.Link) (string, int, error) {

	root := vaults_domain.ThreadsRoot{
		Items: links,
	}

	return s.putNode(root)
}

func (s *VaultService) RotateThreadBranch(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, map[string][]vaults_domain.Link, map[string][]vaults_domain.Link, error) {

	for i := range vp.Collaborative.Threads {

		vp.Collaborative.Threads[i].IsDirty = true
	}

	return s.BuildThreadsBranch(session, vp, mode)
}

// =======================================================================================
// READ
// =======================================================================================
func (r *VaultReconstructor) resolveThreads(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	threadsRoot vaults_domain.ThreadsRoot,
) ([]vaults_domain.Thread, error) {

	var result []vaults_domain.Thread

	for _, link := range threadsRoot.Items {

		res, err := r.Query.Execute(ctx, cmd.WithCID(link.CID))
		if err != nil {
			return result, err
		}

		var thread vaults_domain.Thread
		if err := json.Unmarshal(res.Raw, &thread); err != nil {
			utils.LogPretty("VaultReconstructor - resolveThreads - error get ipfs data", err)
			return result, err
		}

		result = append(result, thread)
	}

	return result, nil
}

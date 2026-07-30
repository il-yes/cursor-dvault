package vaults_service

import (
	"time"

	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
)

type ShareEntryNode struct {
	ID           string
	AssetCID     string
	TrustGroupID string
	WrappedDEK   string
	CreatedBy    string
	CreatedAt    time.Time
	Metadata     map[string]string
}

func (s *VaultService) BuildShareEntriesBranch(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, error) {

	shareEntries := vp.Collaborative.ShareEntries

	links, err := s.BuildShareEntries(shareEntries, mode)
	if err != nil {
		return "", err
	}

	cid, _, err := s.BuildShareEntriesRoot(links)
	if err != nil {
		return "", err
	}

	return cid, nil
}

func (s *VaultService) BuildShareEntries(entries []vaults_domain.ShareEntry, mode SyncMode) ([]vaults_domain.Link, error) {

	var links []vaults_domain.Link

	for _, entry := range entries {
		node := ShareEntryNode{
			ID:           entry.ID,
			AssetCID:     entry.AssetCID,
			TrustGroupID: entry.TrustGroupID,
			WrappedDEK:   entry.WrappedDEK,
			CreatedBy:    entry.CreatedBy,
			CreatedAt:    entry.CreatedAt,
			Metadata:     entry.Metadata,
		}

		cid, _, err := s.putNode(node)

		if err != nil {
			return nil, err
		}

		links = append(links, vaults_domain.Link{CID: cid})
	}

	return links, nil
}

func (s *VaultService) BuildShareEntriesRoot(links []vaults_domain.Link) (string, int, error) {
	root := vaults_domain.ShareEntriesRoot{
		Items: links,
	}

	return s.putNode(root)
}

func (s *VaultService) RotateShareEntryKeys(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, error) {
	// mark all entries by type dirty
	for i := range vp.Collaborative.ShareEntries {
		vp.Collaborative.ShareEntries[i].IsDirty = true
	}
	return s.BuildShareEntriesBranch(session, vp, mode)
}

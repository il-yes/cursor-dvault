package vaults_service

import (
	"context"
	"encoding/json"

	"vault-app/internal/utils"
	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
)



func (s *VaultService) CommitAttachments(
	session vault_session.Session,
	vp vaults_domain.VaultPayload,
	mode SyncMode,
	foldersCID string,
	entriesCID string,
	indexCID string,
	entryUpdate []EntryUpdate,
) (string, []EntryUpdate, int, int, error) {
	utils.LogPretty("VaultService - CommitVault - ", "starting....")

	// =========================
	// 1. BUILD ATTACHEMENTS branch
	// =========================
	attachementRootCID, err := s.BuildAttachmentsBranch(session, vp, mode)
	if err != nil {
		return "", nil, 0, 0, err
	}

	// =========================
	// 2. Save ATTACHEMENTS ROOT
	// =========================
	return s.SaveVaultRoot(
		attachementRootCID,
		foldersCID,
		entriesCID,
		indexCID,
		entryUpdate,
		session,
	)

}


// =======================================================================================
// WRITE
// =======================================================================================
func (s *VaultService) BuildAttachmentsBranch(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, error) {
	utils.LogPretty("VaultService - CommitVault - ", "starting....")

	// =========================
	// 0. BUILD ATTACHEMENTS Nodes
	// =========================
	attachements := vp.GetAttachments()
	attachementLinks, err := s.buildAttachmentLinks(session.UserID, session.Runtime.VaultName, attachements, mode)
	if err != nil {
		return "", err
	}

	// =========================
	// 1. ATTACHEMENTS ROOT
	// =========================
	attachementCIDs, _, err := s.BuildAttachmentsRoot(attachementLinks)
	if err != nil {
		return "", err
	}

	return attachementCIDs, nil
}

func (s *VaultService) BuildAttachmentsRoot(links []vaults_domain.Link) (string, int, error) {
	root := vaults_domain.AttachementsRoot{Items: links}
	return s.putNode(root)
}

func (s *VaultService) buildAttachmentLinks(userID string, vaultName string, attachements []vaults_domain.Attachment, mode SyncMode) ([]vaults_domain.Link, error) {
	utils.LogPretty("VaultService - buildAttachmentLinks - ", "starting....")
	var links []vaults_domain.Link
	policy := resolvePolicy(mode)
	utils.LogPretty("VaultService - buildAttachmentLinks - policy", policy)

	for i := range attachements {
		attachement := &attachements[i]
		// =========================
		// 🟢 REUSE PATH (ONLY IF POLICY ALLOWS)
		// =========================
		if policy.AllowReuse && attachement.NodeCID != "" && !attachement.IsDirty {
			links = append(links, vaults_domain.Link{CID: attachement.NodeCID})
			continue
		}

		// Fetch local file
		res, err := s.VaultHandler.LoadAttachment(userID, vaultName, attachement.Hash, "bytes")
		if err != nil {
			return nil, err
		}

		// Upload to ipfs attachement file
		attachementCid, _, err := s.putRawFile(res.File)
		if err != nil {
			return nil, err
		}
		attachement.FileCID = attachementCid
		attachement.IsDirty = false

		// Get attachment Node link
		attachmentNodeLink, err := s.GetAttachmentNodeLink(*attachement)
		if err != nil {
			return nil, err
		}

		links = append(links, *attachmentNodeLink)
	}
	return links, nil
}

func (s *VaultService) GetAttachmentNodeLink(attachement vaults_domain.Attachment) (*vaults_domain.Link, error) {
	node := vaults_domain.AttachmentNode{
		Name:         attachement.Name,
		Size:         attachement.Size,
		Ext:          attachement.Ext,
		DownloadedAt: attachement.DownloadedAt,
		FileCID:      attachement.FileCID,
		Hash:         attachement.Hash,
	}

	cid, _, err := s.putNode(node)
	if err != nil {
		return nil, err
	}
	return &vaults_domain.Link{CID: cid}, nil
}


func (s *VaultService) RotateAttachmentGraph(
	session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode,
) {
	// mark all attachments dirty
	for i := range vp.Attachments {
		vp.Attachments[i].IsDirty = true
	}
	// 	↓
	s.BuildAttachmentsBranch(session, vp, mode)
}


// =======================================================================================
// READ
// =======================================================================================
func (r *VaultReconstructor) resolveAttachments(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	attachmentsRoot vaults_domain.AttachementsRoot,
) ([]vaults_domain.Attachment, error) {

	var result []vaults_domain.Attachment

	for _, link := range attachmentsRoot.Items {

		res, err := r.Query.Execute(ctx, cmd.WithCID(link.CID))
		if err != nil {
			return result, err
		}

		var attachment vaults_domain.Attachment
		if err := json.Unmarshal(res.Raw, &attachment); err != nil {
			return result, err
		}

		result = append(result, attachment)
	}

	return result, nil
}

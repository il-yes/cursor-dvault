package vaults_service

import (
	"vault-app/internal/utils"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
)



type EntryNode struct {
	ID        string
	Type      vaults_domain.EntryType
	EntryName string
	Data      vaults_domain.EntryInterface
	// Entries    []vaults_domain.Link
	AttachmentCIDs []string
}


func (s *VaultService) CommitEntries(
	session vault_session.Session,
	vp vaults_domain.VaultPayload,
	mode SyncMode,
	attachementRootCID string,
	foldersCID string,
) (string, []EntryUpdate, int, int, error) {
	utils.LogPretty("VaultService - CommitVault - ", "starting....")

	// =========================
	// 1. BUILD ATTACHEMENTS branch
	// =========================
	entriesCID, indexByType, indexByFolder, entryUpdates, err := s.BuildEntriesBranch(session, vp, mode)
	if err != nil {
		return "", nil, 0, 0, err
	}


	// =========================
	// Z. INDEX
	// =========================
	indexCID, _, err := s.buildIndex(indexByType, indexByFolder)
	if err != nil {
		return "", nil, 0, 0, err
	}


	// =========================
	// 3. Save ATTACHEMENTS ROOT
	// =========================
	return s.SaveVaultRoot(
		attachementRootCID,
		foldersCID,
		entriesCID,
		indexCID,
		entryUpdates,
		session,
	)

}

func (s *VaultService) BuildEntriesBranch(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, map[string][]vaults_domain.Link, map[string][]vaults_domain.Link, []EntryUpdate, error) {
	// =========================
	// 1. BUILD ENTRIES
	// =========================
	entries := vp.Personal.Entries
	entryLinks, indexByType, indexByFolder, entryUpdates, err := s.BuildEntries(entries, mode)
	if err != nil {
		return "", nil, nil, nil, err
	}

	
	// =========================
	// 4. ENTRIES ROOT
	// =========================
	entriesCID, _, err := s.BuildEntriesRoot(entryLinks, mode)
	if err != nil {
		return "", nil, nil, nil, err
	}

	return  entriesCID, indexByType, indexByFolder, entryUpdates, nil
}


func (s *VaultService) BuildEntries(entries vaults_domain.Entries, mode SyncMode) ([]vaults_domain.Link, map[string][]vaults_domain.Link, map[string][]vaults_domain.Link, []EntryUpdate, error) {

	var links []vaults_domain.Link
	byType := make(map[string][]vaults_domain.Link)
	byFolder := make(map[string][]vaults_domain.Link)
	var entriesUpdate []EntryUpdate
	policy := resolvePolicy(mode)

	s.processEntryList(loginToInterfaces(entries.Login), &links, byType, byFolder, &entriesUpdate, policy)
	s.processEntryList(cardToInterfaces(entries.Card), &links, byType, byFolder, &entriesUpdate, policy)
	s.processEntryList(identityToInterfaces(entries.Identity), &links, byType, byFolder, &entriesUpdate, policy)
	s.processEntryList(noteToInterfaces(entries.Note), &links, byType, byFolder, &entriesUpdate, policy)
	s.processEntryList(sshKeyToInterfaces(entries.SSHKey), &links, byType, byFolder, &entriesUpdate, policy)

	return links, byType, byFolder, entriesUpdate, nil
}


func (s *VaultService) BuildEntriesRoot(links []vaults_domain.Link, mode SyncMode) (string, int, error) {
	root := vaults_domain.EntriesRoot{Items: links}
	return s.putNode(root)
}


func (s *VaultService) processEntryList(
	list []vaults_domain.EntryInterface,
	links *[]vaults_domain.Link,
	byType map[string][]vaults_domain.Link,
	byFolder map[string][]vaults_domain.Link,
	eu *[]EntryUpdate,
	policy SyncPolicy,
) error {

	addLink := func(base vaults_domain.BaseEntry, cid string) {
		link := vaults_domain.Link{CID: cid}

		*links = append(*links, link)
		byType[string(base.Type)] = append(byType[string(base.Type)], link)

		if base.FolderID != "" {
			byFolder[base.FolderID] = append(byFolder[base.FolderID], link)
		}
	}

	for _, entry := range list {
		base := entry.GetBase()

		// 🚫 skip drafts always
		if base.IsDraft {
			continue
		}

		// =========================
		// 🟢 REUSE PATH (ONLY IF POLICY ALLOWS)
		// =========================
		if policy.AllowReuse && base.CID != "" && !base.IsDirty {

			addLink(*base, base.CID)

			*eu = append(*eu, EntryUpdate{
				ID:      base.ID,
				CID:     base.CID,
				IsDirty: false,
				Reused:  true,
			})

			continue
		}

		// =========================
		// 🔥 REBUILD PATH
		// =========================
		node := EntryNode{
			ID:             base.ID,
			Type:           base.Type,
			EntryName:      base.EntryName,
			Data:           entry,
			AttachmentCIDs: entry.GetBase().AttachmentCIDs, // s.buildAttachmentLinks(base.AttachmentCIDs),
		}

		cid, bytes, err := s.putNode(node)
		if err != nil {
			return err
		}

		base.CID = cid
		base.IsDirty = false

		addLink(*base, cid)

		*eu = append(*eu, EntryUpdate{
			ID:         base.ID,
			CID:        cid,
			TotalBytes: bytes,
			IsDirty:    false,
			Reused:     false,
		})
	}

	return nil
}


func loginToInterfaces(list []vaults_domain.LoginEntry) []vaults_domain.EntryInterface {
	result := make([]vaults_domain.EntryInterface, len(list))
	for i := range list {
		result[i] = &list[i]
	}
	return result
}

func cardToInterfaces(list []vaults_domain.CardEntry) []vaults_domain.EntryInterface {
	result := make([]vaults_domain.EntryInterface, len(list))
	for i := range list {
		result[i] = &list[i]
	}
	return result
}
func noteToInterfaces(list []vaults_domain.NoteEntry) []vaults_domain.EntryInterface {
	result := make([]vaults_domain.EntryInterface, len(list))
	for i := range list {
		result[i] = &list[i]
	}
	return result
}

func identityToInterfaces(list []vaults_domain.IdentityEntry) []vaults_domain.EntryInterface {
	result := make([]vaults_domain.EntryInterface, len(list))
	for i := range list {
		result[i] = &list[i]
	}
	return result
}
func sshKeyToInterfaces(list []vaults_domain.SSHKeyEntry) []vaults_domain.EntryInterface {
	result := make([]vaults_domain.EntryInterface, len(list))
	for i := range list {
		result[i] = &list[i]
	}
	return result
}

func (s *VaultService) RotateEntryKeys(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, map[string][]vaults_domain.Link, map[string][]vaults_domain.Link, []EntryUpdate, error) {
	// mark all entries by type dirty
	for i := range vp.Entries.Login {
		vp.Entries.Login[i].BaseEntry.IsDirty = true
	}
	for i := range vp.Entries.Card {
		vp.Entries.Card[i].BaseEntry.IsDirty = true
	}
	for i := range vp.Entries.Identity {
		vp.Entries.Identity[i].BaseEntry.IsDirty = true
	}
	for i := range  vp.Entries.Note {
		 vp.Entries.Note[i].BaseEntry.IsDirty = true
	}
	for i := range vp.Entries.SSHKey {
		vp.Entries.SSHKey[i].BaseEntry.IsDirty = true
	}
	// 	↓
	return s.BuildEntriesBranch(session, vp, mode)
}
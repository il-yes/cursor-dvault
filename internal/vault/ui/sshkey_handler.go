package vault_ui

import (
	"fmt"
	"time"
	"vault-app/internal/blockchain"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"

	"github.com/google/uuid"
)

type SSHKeyHandler struct {
	db     models.DBModel
	ipfs   blockchain.IPFSClient
	logger *logger.Logger
	NowUTC func() string
	Vault  vaults_domain.VaultPayload
	Session *vault_session.Session
}

func NewSSHKeyHandler(db models.DBModel, ipfs blockchain.IPFSClient, log *logger.Logger) *SSHKeyHandler {
	return &SSHKeyHandler{
		db:     db,
		ipfs:   ipfs,
		logger: log,
		NowUTC: func() string { return time.Now().Format(time.RFC3339) },
	}
}

func (h *SSHKeyHandler) Add(userID string, anEntry any) (*vaults_domain.VaultPayload, error) {
	entry, err := anEntry.(*vaults_domain.SSHKeyEntry)
	if !err {
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entry.ID = uuid.New().String() // Ensure entry has a UUID
	h.Vault.Entries.SSHKey = append(h.Vault.Entries.SSHKey, *entry)
	// session.LastUpdated = h.NowUTC()
	// session.Dirty = true

	h.logger.Info("✅ Added ssh key entry for user %s: %s\n", userID, entry.EntryName)

	return &h.Vault, nil

}
func (h *SSHKeyHandler) Edit(userID string, entry any) (*vaults_domain.VaultPayload, error) {
	updatedEntry, ok := entry.(*vaults_domain.SSHKeyEntry)
	if !ok {
		return nil, fmt.Errorf("invalid type: expected SSHKeyEntry")
	}

	entries := h.Vault.Entries.SSHKey
	updated := false

	for i, entry := range entries {
		if entry.ID == updatedEntry.ID {
			// Update the fields (you could also do a full replace)
			entries[i] = *updatedEntry
			updated = true
			break
		}
	}

	if !updated {
		return nil, fmt.Errorf("entry with ID %s not found for user %s", updatedEntry.ID, userID)
	}

	h.Vault.Entries.SSHKey = entries
	// h.MarkDirty(userID)

	h.logger.Info("✏️ Updated ssh key entry for user %s: %s\n", userID, updatedEntry.EntryName)

	return &h.Vault, nil
}
func (h *SSHKeyHandler) Trash(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashSSHKeyEntryAction(userID, entryID, true)
}
func (h *SSHKeyHandler) Restore(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashSSHKeyEntryAction(userID, entryID, false)
}
func (h *SSHKeyHandler) TrashSSHKeyEntryAction(userID string, entryID string, trashed bool) (*vaults_domain.VaultPayload, error) {

	for i, entry := range h.Vault.Entries.SSHKey {
		if entry.ID == entryID {
			h.Vault.Entries.SSHKey[i].Trashed = trashed
			// h.MarkDirty(userID)

			state := "restored"
			if trashed {
				state = "trashed"
			}
			h.logger.Info("🗑️ %s ssh key entry %s for user %s", state, entryID, userID)

			return &h.Vault, nil
		}
	}
	return nil, fmt.Errorf("entry with ID %s not found", entryID)
}

func (h *SSHKeyHandler) SetSession(session *vault_session.Session) {
	s := session
	h.Session = s
	payload, err := vault_session.DecodeSessionVault(s.Vault)
	if err != nil {
		return
	}
	h.Vault = *payload
}

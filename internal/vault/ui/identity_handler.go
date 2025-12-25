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

type IdentityHandler struct {
	db     models.DBModel
	ipfs   blockchain.IPFSClient
	logger *logger.Logger
	NowUTC func() string
	Vault  *vaults_domain.VaultPayload
	Session *vault_session.Session
}

func NewIdentityHandler(db models.DBModel, ipfs blockchain.IPFSClient, log *logger.Logger) *IdentityHandler {
	return &IdentityHandler{
		db:     db,
		ipfs:   ipfs,
		logger: log,
		NowUTC: func() string { return time.Now().Format(time.RFC3339) },
	}
}

func (h *IdentityHandler) Add(userID string, anEntry any) (*vaults_domain.VaultPayload, error) {
	if h.Vault == nil {
		return nil, fmt.Errorf("no active session for user %s", userID)
	}
	entry, err := anEntry.(*vaults_domain.IdentityEntry)
	if !err {
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entry.ID = uuid.New().String() // Ensure entry has a UUID
	h.Vault.Entries.Identity = append(h.Vault.Entries.Identity, *entry)
	// session.LastUpdated = h.NowUTC()
	// session.Dirty = true

	h.logger.Info("✅ Added identity entry for user %s: %s\n", userID, entry.EntryName)

	return h.Vault, nil

}
func (h *IdentityHandler) Edit(userID string, entry any) (*any, error) {
	if h.Vault == nil {
		return nil, fmt.Errorf("no active session for user %s", userID)
	}
	updatedEntry, ok := entry.(*vaults_domain.IdentityEntry)
	if !ok {
		return nil, fmt.Errorf("invalid type: expected IdentityEntry")
	}

	entries := h.Vault.Entries.Identity
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

	h.Vault.Entries.Identity = entries
	// h.MarkDirty(userID)

	h.logger.Info("✏️ Updated identity entry for user %s: %s\n", userID, updatedEntry.EntryName)
	// utils.LogPretty("session after update", session)

	var result any = updatedEntry
	return &result, nil
}
func (h *IdentityHandler) Trash(userID string, entryID string) error {
	return h.TrashIdentityEntryAction(userID, entryID, true)
}
func (h *IdentityHandler) Restore(userID string, entryID string) error {
	return h.TrashIdentityEntryAction(userID, entryID, false)
}
func (h *IdentityHandler) TrashIdentityEntryAction(userID string, entryID string, trashed bool) error {
	if h.Vault == nil {
		return fmt.Errorf("no active session for user %s", userID)
	}

	for i, entry := range h.Vault.Entries.Identity {
		if entry.ID == entryID {
			h.Vault.Entries.Identity[i].Trashed = trashed
			// h.MarkDirty(userID)

			state := "restored"
			if trashed {
				state = "trashed"
			}
			h.logger.Info("🗑️ %s identity entry %s for user %s", state, entryID, userID)

			return nil
		}
	}
	return fmt.Errorf("entry with ID %s not found", entryID)
}

func (h *IdentityHandler) SetVault(vault *vault_session.Session) {
	p := vault.Vault
	h.Vault = p
}
func (h *IdentityHandler) SetSession(session *vault_session.Session) {
	s := session
	h.Session = s
}

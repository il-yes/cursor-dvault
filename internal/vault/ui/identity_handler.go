package vault_ui

import (
	"fmt"
	"time"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"

	"github.com/google/uuid"
)

type IdentityHandler struct {
	db     models.DBModel
	logger *logger.Logger
	NowUTC func() string
	Vault  *vaults_domain.VaultPayload
	Session *vault_session.Session
}

func NewIdentityHandler(db models.DBModel, log *logger.Logger) *IdentityHandler {
	return &IdentityHandler{
		db:     db,
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

	h.logger.Info("✅ Added identity entry for user %s: %s\n", userID, entry.EntryName)

	return h.Vault, nil

}
func (h *IdentityHandler) Edit(userID string, entry any) (*vaults_domain.VaultPayload, error) {
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
			entries[i] = *updatedEntry
			updated = true
			break
		}
	}

	if !updated {
		return nil, fmt.Errorf("entry with ID %s not found for user %s", updatedEntry.ID, userID)
	}

	h.Vault.Entries.Identity = entries

	h.logger.Info("✏️ Updated identity entry for user %s: %s\n", userID, updatedEntry.EntryName)

	return h.Vault, nil
}
func (h *IdentityHandler) Trash(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashIdentityEntryAction(userID, entryID, true)
}
func (h *IdentityHandler) Restore(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashIdentityEntryAction(userID, entryID, false)
}
func (h *IdentityHandler) TrashIdentityEntryAction(userID string, entryID string, trashed bool) (*vaults_domain.VaultPayload, error) {
	if h.Vault == nil {
		return nil, fmt.Errorf("no active session for user %s", userID)
	}

	for i, entry := range h.Vault.Entries.Identity {
		if entry.ID == entryID {
			h.Vault.Entries.Identity[i].Trashed = trashed
			state := "restored"
			if trashed {
				state = "trashed"
			}
			h.logger.Info("🗑️ %s identity entry %s for user %s", state, entryID, userID)

			return h.Vault, nil
		}
	}
	return nil, fmt.Errorf("entry with ID %s not found", entryID)
}

func (h *IdentityHandler) SetSession(session *vault_session.Session) {
	s := session
	h.Session = s
	payload, err := vault_session.DecodeSessionVault(s.Vault)
	if err != nil {
		return
	}
	h.Vault = payload
}

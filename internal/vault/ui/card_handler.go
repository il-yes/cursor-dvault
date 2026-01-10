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

type CardHandler struct {
	db     models.DBModel
	logger *logger.Logger
	NowUTC func() string
	Vault  vaults_domain.VaultPayload
	Session *vault_session.Session
}

func NewCardHandler(db models.DBModel, log *logger.Logger) *CardHandler {
	return &CardHandler{
		db:     db,
		logger: log,
		NowUTC: func() string { return time.Now().Format(time.RFC3339) },
	}
}

func (h *CardHandler) Add(userID string, anEntry any) (*vaults_domain.VaultPayload, error) {
	entry, err := anEntry.(*vaults_domain.CardEntry)
	if !err {
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entry.ID = uuid.New().String() // Ensure entry has a UUID
	h.Vault.Entries.Card = append(h.Vault.Entries.Card, *entry)
	// session.LastUpdated = h.NowUTC()
	// session.Dirty = true

	h.logger.Info("✅ Added card entry for user %d: %s\n", userID, entry.EntryName)

	return &h.Vault, nil

}
func (h *CardHandler) Edit(userID string, entry any) (*vaults_domain.VaultPayload, error) {
	updatedEntry, ok := entry.(*vaults_domain.CardEntry)
	if !ok {
		return nil, fmt.Errorf("invalid type: expected CardEntry")
	}
	updatedEntry.IsDraft = true

	entries := h.Vault.Entries.Card
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
		return nil, fmt.Errorf("entry with ID %s not found for user %d", updatedEntry.ID, userID)
	}

	h.Vault.Entries.Card = entries
	// h.MarkDirty(userID)

	h.logger.Info("✏️ Updated card entry for user %d: %s\n", userID, updatedEntry.EntryName)

	return &h.Vault, nil
}
func (h *CardHandler) Trash(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashCardEntryAction(userID, entryID, true)
}
func (h *CardHandler) Restore(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashCardEntryAction(userID, entryID, false)
}
func (h *CardHandler) TrashCardEntryAction(userID string, entryID string, trashed bool) (*vaults_domain.VaultPayload, error) {
	for i, entry := range h.Vault.Entries.Card {
		if entry.ID == entryID {
			h.Vault.Entries.Card[i].Trashed = trashed

			state := "restored"
			if trashed {
				state = "trashed"
			}
			h.logger.Info("🗑️ %s card entry %s for user %d", state, entryID, userID)

			return &h.Vault, nil
		}
	}
	return nil, fmt.Errorf("entry with ID %s not found", entryID)
}
func (h *CardHandler) SetVault(vault *vault_session.Session) {
	p := vault.Vault
	payload, err := vault_session.DecodeSessionVault(p)
	if err != nil {
		return
	}
	h.Vault = *payload
}
func (h *CardHandler) SetSession(session *vault_session.Session) {
	s := session
	h.Session = s
	payload, err := vault_session.DecodeSessionVault(s.Vault)
	if err != nil {
		return
	}
	h.Vault = *payload
}

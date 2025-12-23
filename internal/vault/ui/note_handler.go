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

type NoteHandler struct {
	db       models.DBModel
	ipfs     blockchain.IPFSClient
	logger   *logger.Logger
	NowUTC   func() string
	Vault    vaults_domain.VaultPayload
}

func NewNoteHandler(db models.DBModel, ipfs blockchain.IPFSClient, log *logger.Logger) *NoteHandler {
	return &NoteHandler{
		db:       db,
		ipfs:     ipfs,
		logger:   log,
		NowUTC:   func() string { return time.Now().Format(time.RFC3339) },
	}
}

func (h *NoteHandler) Add(userID string, anEntry any) (*any, error) {
	entry, err := anEntry.(*vaults_domain.NoteEntry)
	if !err {
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entry.ID = uuid.New().String() // Ensure entry has a UUID
	h.Vault.Entries.Note = append(h.Vault.Entries.Note, *entry)
	// session.LastUpdated = h.NowUTC()
	// session.Dirty = true

	h.logger.Info("✅ Added note entry for user %s: %s\n", userID, entry.EntryName)

	var result any = entry
	return &result, nil

}
func (h *NoteHandler) Edit(userID string, entry any) (*any, error) {
	updatedEntry, ok := entry.(*vaults_domain.NoteEntry)
	if !ok {
		return nil, fmt.Errorf("invalid type: expected NoteEntry")
	}

	entries := h.Vault.Entries.Note
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

	h.Vault.Entries.Note = entries
	// h.MarkDirty(userID)

	h.logger.Info("✏️ Updated note entry for user %s: %s\n", userID, updatedEntry.EntryName)
	// utils.LogPretty("session after update", session)

	var result any = updatedEntry
	return &result, nil
}
func (h *NoteHandler) Trash(userID string, entryID string) error {
	return h.TrashNoteEntryAction(userID, entryID, true)
}
func (h *NoteHandler) Restore(userID string, entryID string) error {
	return h.TrashNoteEntryAction(userID, entryID, false)
}
func (h *NoteHandler) TrashNoteEntryAction(userID string, entryID string, trashed bool) error {

	for i, entry := range h	.Vault.Entries.Note {
		if entry.ID == entryID {
			h.Vault.Entries.Note[i].Trashed = trashed
			// h.MarkDirty(userID)

			state := "restored"
			if trashed {
				state = "trashed"
			}
			h.logger.Info("🗑️ %s note entry %s for user %s", state, entryID, userID)

			return nil
		}
	}
	return fmt.Errorf("entry with ID %s not found", entryID)
}

func (h *NoteHandler) SetVault(vault *vault_session.Session) {
	h.Vault = *vault.Vault
}
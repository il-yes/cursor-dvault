package vault_ui

import (
	"fmt"
	"time"
	utils "vault-app/internal/utils"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"

	"github.com/google/uuid"
)

type NoteHandler struct {
	db     models.DBModel
	logger *logger.Logger
	NowUTC func() string
	Vault  vaults_domain.VaultPayload
	Session *vault_session.Session	
}

func NewNoteHandler(db models.DBModel, log *logger.Logger) *NoteHandler {
	return &NoteHandler{
		db:     db,
		logger: log,
		NowUTC: func() string { return time.Now().Format(time.RFC3339) },
	}
}

func (h *NoteHandler) Add(userID string, anEntry any) (*vaults_domain.VaultPayload, error) {
	entry, err := anEntry.(*vaults_domain.NoteEntry)
	if !err {
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entry.ID = uuid.New().String() // Ensure entry has a UUID
	h.Vault.Entries.Note = append(h.Vault.Entries.Note, *entry)

	h.logger.Info("✅ Added note entry for user %s: %s\n", userID, entry.EntryName)

	return &h.Vault, nil

}
func (h *NoteHandler) Edit(userID string, entry any) (*vaults_domain.VaultPayload, error) {
	updatedEntry, ok := entry.(*vaults_domain.NoteEntry)
	if !ok {
		return nil, fmt.Errorf("invalid type: expected NoteEntry")
	}

	entries := h.Vault.Entries.Note
	updated := false
	utils.LogPretty("NoteHandler - Edit - Vault Before", h.Vault)
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

	h.Vault.Entries.Note = entries

	h.logger.Info("✏️ Updated note entry for user %s: %s\n", userID, updatedEntry.EntryName)

	return &h.Vault, nil
}
func (h *NoteHandler) Trash(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashNoteEntryAction(userID, entryID, true)
}
func (h *NoteHandler) Restore(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashNoteEntryAction(userID, entryID, false)
}
func (h *NoteHandler) TrashNoteEntryAction(userID string, entryID string, trashed bool) (*vaults_domain.VaultPayload, error) {

	for i, entry := range h.Vault.Entries.Note {
		if entry.ID == entryID {
			h.Vault.Entries.Note[i].Trashed = trashed
			state := "restored"
			if trashed {
				state = "trashed"
			}
			h.logger.Info("🗑️ %s note entry %s for user %s", state, entryID, userID)

			return &h.Vault, nil
		}
	}
	return nil, fmt.Errorf("entry with ID %s not found", entryID)
}

func (h *NoteHandler) SetSession(session *vault_session.Session) {
	s := session
	h.Session = s
	payload, err := vault_session.DecodeSessionVault(s.Vault)
	if err != nil {
		return
	}
	h.Vault = *payload
}
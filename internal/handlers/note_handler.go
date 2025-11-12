package handlers

import (
	"errors"
	"fmt"
	"time"
	utils "vault-app/internal"
	"vault-app/internal/blockchain"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"

	"github.com/google/uuid"
)

type NoteHandler struct {
	db         models.DBModel
	ipfs       blockchain.IPFSClient
	logger     *logger.Logger
	NowUTC     func() string
	sessions   map[int]*models.VaultSession
}

func NewNoteHandler(db models.DBModel, ipfs blockchain.IPFSClient, sessions map[int]*models.VaultSession, log *logger.Logger) *NoteHandler {
	return &NoteHandler{
		db:       db,
		ipfs:     ipfs,
		logger:   log,
		NowUTC:   func() string { return time.Now().Format(time.RFC3339) },
		sessions: sessions,
	}
}

func (h *NoteHandler) GetSession(userID int) (*models.VaultSession, error) {
	session, ok := h.sessions[userID]
	if !ok {
		return nil, errors.New("no vault session found")
	}
	return session, nil
}

func (h *NoteHandler) Add(userID int, anEntry any) (*any, error) {
	session, ok := h.GetSession(userID)
	if ok != nil {
		return nil, fmt.Errorf("no active session for user %d", userID)
	}
	entry, err := anEntry.(*models.NoteEntry)
	if !err {
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entry.ID = uuid.New().String() // Ensure entry has a UUID
	session.Vault.Entries.Note = append(session.Vault.Entries.Note, *entry)
	// session.LastUpdated = h.NowUTC()
	// session.Dirty = true

	h.logger.Info("✅ Added note entry for user %d: %s\n", userID, entry.EntryName)

	var result any = entry
	return &result, nil

}
func (h *NoteHandler) Edit(userID int, entry any) (*any, error) {
	session, err := h.GetSession(userID)
	if err != nil {
		return nil, fmt.Errorf("no active session for user %d", userID)
	}
	updatedEntry, ok := entry.(*models.NoteEntry)
	if !ok {
		return nil, fmt.Errorf("invalid type: expected NoteEntry")
	}

	entries := session.Vault.Entries.Note
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

	session.Vault.Entries.Note = entries
	// h.MarkDirty(userID)


	h.logger.Info("✏️ Updated note entry for user %d: %s\n", userID, updatedEntry.EntryName)
	utils.LogPretty("session after update", session)

	var result any = updatedEntry
	return &result, nil
}
func (h *NoteHandler) Trash(userID int, entryID string) error {
	return h.TrashNoteEntryAction(userID, entryID, true)
}
func (h *NoteHandler) Restore(userID int, entryID string) error {
	return h.TrashNoteEntryAction(userID, entryID, false)
}
func (h *NoteHandler) TrashNoteEntryAction(userID int, entryID string, trashed bool) error {
	session, err := h.GetSession(userID)
	if err != nil {
		return err
	}

	for i, entry := range session.Vault.Entries.Note {
		if entry.ID == entryID {
			session.Vault.Entries.Note[i].Trashed = trashed
			// h.MarkDirty(userID)

			state := "restored"
			if trashed {
				state = "trashed"
			}
			h.logger.Info("🗑️ %s note entry %s for user %d", state, entryID, userID)

			return nil
		}
	}
	return fmt.Errorf("entry with ID %s not found", entryID)
}
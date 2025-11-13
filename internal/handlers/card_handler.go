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

type CardHandler struct {
	db         models.DBModel
	ipfs       blockchain.IPFSClient
	logger     *logger.Logger
	NowUTC     func() string
	sessions   map[int]*models.VaultSession
}

func NewCardHandler(db models.DBModel, ipfs blockchain.IPFSClient, sessions map[int]*models.VaultSession, log *logger.Logger) *CardHandler {
	return &CardHandler{
		db:       db,
		ipfs:     ipfs,
		logger:   log,
		NowUTC:   func() string { return time.Now().Format(time.RFC3339) },
		sessions: sessions,
	}
}

func (h *CardHandler) GetSession(userID int) (*models.VaultSession, error) {
	session, ok := h.sessions[userID]
	if !ok {
		return nil, errors.New("no vault session found")
	}
	return session, nil
}

func (h *CardHandler) Add(userID int, anEntry any) (*any, error) {
	session, ok := h.GetSession(userID)
	if ok != nil {
		return nil, fmt.Errorf("no active session for user %d", userID)
	}
	entry, err := anEntry.(*models.CardEntry)
	if !err {
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entry.ID = uuid.New().String() // Ensure entry has a UUID
	session.Vault.Entries.Card = append(session.Vault.Entries.Card, *entry)
	// session.LastUpdated = h.NowUTC()
	// session.Dirty = true

	h.logger.Info("✅ Added card entry for user %d: %s\n", userID, entry.EntryName)


	var result any = entry
	return &result, nil

}
func (h *CardHandler) Edit(userID int, entry any) (*any, error) {
	session, err := h.GetSession(userID)
	if err != nil {
		return nil, fmt.Errorf("no active session for user %d", userID)
	}
	updatedEntry, ok := entry.(*models.CardEntry)
	if !ok {
		return nil, fmt.Errorf("invalid type: expected CardEntry")
	}
	updatedEntry.IsDraft = true

	entries := session.Vault.Entries.Card
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

	session.Vault.Entries.Card = entries
	// h.MarkDirty(userID)


	h.logger.Info("✏️ Updated card entry for user %d: %s\n", userID, updatedEntry.EntryName)
	utils.LogPretty("session after update", session)

	var result any = updatedEntry
	return &result, nil
}
func (h *CardHandler) Trash(userID int, entryID string) error {
	return h.TrashCardEntryAction(userID, entryID, true)
}
func (h *CardHandler) Restore(userID int, entryID string) error {
	return h.TrashCardEntryAction(userID, entryID, false)
}
func (h *CardHandler) TrashCardEntryAction(userID int, entryID string, trashed bool) error {
	session, err := h.GetSession(userID)
	if err != nil {
		return err
	}

	for i, entry := range session.Vault.Entries.Card {
		if entry.ID == entryID {
			session.Vault.Entries.Card[i].Trashed = trashed
			// h.MarkDirty(userID)

			state := "restored"
			if trashed {
				state = "trashed"
			}
			h.logger.Info("🗑️ %s card entry %s for user %d", state, entryID, userID)

			return nil
		}
	}
	return fmt.Errorf("entry with ID %s not found", entryID)
}
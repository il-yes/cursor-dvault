package handlers

import (
	"errors"
	"fmt"
	"time"
	"vault-app/internal/blockchain"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"

	"github.com/google/uuid"
)

type LoginHandler struct {
	db         models.DBModel
	ipfs       blockchain.IPFSClient
	logger     *logger.Logger
	NowUTC     func() string
	sessions   map[string]*models.VaultSession
}

func NewLoginHandler(db models.DBModel, ipfs blockchain.IPFSClient, sessions map[string]*models.VaultSession, log *logger.Logger) *LoginHandler {
	return &LoginHandler{
		db:       db,
		ipfs:     ipfs,
		logger:   log,
		NowUTC:   func() string { return time.Now().Format(time.RFC3339) },
		sessions: sessions,
	}
}

func (h *LoginHandler) GetSession(userID string) (*models.VaultSession, error) {
	session, ok := h.sessions[userID]
	if !ok {
		return nil, errors.New("no vault session found")
	}
	return session, nil
}

func (h *LoginHandler) Add(userID string, anEntry any) (*any, error) {
	session, ok := h.GetSession(userID)
	if ok != nil {
		return nil, fmt.Errorf("no active session for user %s", userID)
	}
	entry, err := anEntry.(*models.LoginEntry)
	if !err {
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entry.ID = uuid.New().String() // Ensure entry has a UUID
	session.Vault.Entries.Login = append(session.Vault.Entries.Login, *entry)
	h.logger.Info("✅ Added login entry for user %d: %s\n", userID, entry.EntryName)

	var result any = entry
	return &result, nil

}
func (h *LoginHandler) Edit(userID string, entry any) (*any, error) {
	session, err := h.GetSession(userID)
	if err != nil {
		return nil, fmt.Errorf("no active session for user %s", userID)
	}
	updatedEntry, ok := entry.(*models.LoginEntry)
	if !ok {
		return nil, fmt.Errorf("invalid type: expected LoginEntry")
	}
	updatedEntry.IsDraft = true

	entries := session.Vault.Entries.Login
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

	session.Vault.Entries.Login = entries
	// h.MarkDirty(userID)


	h.logger.Info("✏️ Updated login entry for user %s: %s\n", userID, updatedEntry.EntryName)
	// utils.LogPretty("session after update", session)

	var result any = updatedEntry
	return &result, nil
}
func (h *LoginHandler) Trash(userID string, entryID string) error {
	return h.TrashLoginEntryAction(userID, entryID, true)
}
func (h *LoginHandler) Restore(userID string, entryID string) error {
	return h.TrashLoginEntryAction(userID, entryID, false)
}
func (h *LoginHandler) TrashLoginEntryAction(userID string, entryID string, trashed bool) error {
	session, err := h.GetSession(userID)
	if err != nil {
		return err
	}

	for i, entry := range session.Vault.Entries.Login {
		if entry.ID == entryID {
			session.Vault.Entries.Login[i].Trashed = trashed
			// h.MarkDirty(userID)

			state := "restored"
			if trashed {
				state = "trashed"
			}
			h.logger.Info("🗑️ %s login entry %s for user %s", state, entryID, userID)	

			return nil
		}
	}
	return fmt.Errorf("entry with ID %s not found", entryID)
}
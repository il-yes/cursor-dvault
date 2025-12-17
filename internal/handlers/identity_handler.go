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

type IdentityHandler struct {
	db         models.DBModel
	ipfs       blockchain.IPFSClient
	logger     *logger.Logger
	NowUTC     func() string
	sessions   map[string]*models.VaultSession
}

func NewIdentityHandler(db models.DBModel, ipfs blockchain.IPFSClient, sessions map[string]*models.VaultSession, log *logger.Logger) *IdentityHandler {
	return &IdentityHandler{
		db:       db,
		ipfs:     ipfs,
		logger:   log,
		NowUTC:   func() string { return time.Now().Format(time.RFC3339) },
		sessions: sessions,
	}
}

func (h *IdentityHandler) GetSession(userID string) (*models.VaultSession, error) {
	session, ok := h.sessions[userID]
	if !ok {
		return nil, errors.New("no vault session found")
	}
	return session, nil
}

func (h *IdentityHandler) Add(userID string, anEntry any) (*any, error) {
	session, ok := h.GetSession(userID)
	if ok != nil {
		return nil, fmt.Errorf("no active session for user %s", userID)
	}
	entry, err := anEntry.(*models.IdentityEntry)
	if !err {
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entry.ID = uuid.New().String() // Ensure entry has a UUID
	session.Vault.Entries.Identity = append(session.Vault.Entries.Identity, *entry)
	// session.LastUpdated = h.NowUTC()
	// session.Dirty = true

	h.logger.Info("✅ Added identity entry for user %d: %s\n", userID, entry.EntryName)

	var result any = entry
	return &result, nil

}
func (h *IdentityHandler) Edit(userID string, entry any) (*any, error) {
	session, err := h.GetSession(userID)
	if err != nil {
		return nil, fmt.Errorf("no active session for user %s", userID)
	}
	updatedEntry, ok := entry.(*models.IdentityEntry)
	if !ok {
		return nil, fmt.Errorf("invalid type: expected IdentityEntry")
	}

	entries := session.Vault.Entries.Identity
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

	session.Vault.Entries.Identity = entries
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
	session, err := h.GetSession(userID)
	if err != nil {
		return err
	}

	for i, entry := range session.Vault.Entries.Identity {
		if entry.ID == entryID {
			session.Vault.Entries.Identity[i].Trashed = trashed
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
package vault_ui

import (
	"fmt"
	"time"
	utils "vault-app/internal"
	"vault-app/internal/blockchain"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"

	"github.com/google/uuid"
)

type LoginHandler struct {
	db     models.DBModel
	ipfs   blockchain.IPFSClient
	logger *logger.Logger
	NowUTC func() string
	Vault  *vaults_domain.VaultPayload
	Session *vault_session.Session
}

func NewLoginHandler(db models.DBModel, ipfs blockchain.IPFSClient, log *logger.Logger) *LoginHandler {
	return &LoginHandler{
		db:     db,
		ipfs:   ipfs,
		logger: log,
		NowUTC: func() string { return time.Now().Format(time.RFC3339) },
	}
}

func (h *LoginHandler) Add(userID string, anEntry any) (*vaults_domain.VaultPayload, error) {
	utils.LogPretty("LoginHandler - Add - anEntry request", anEntry)
	// 1. ---------- Unmarshal entry ----------
	entry, ok := anEntry.(*vaults_domain.LoginEntry)
	if !ok {
		h.logger.Error("LoginHandler - entry does not implement VaultEntry interface: %v", anEntry)
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}
	if h.Vault == nil {
		return nil, fmt.Errorf("vault not initialized for user %s", userID)
	}
	entry.ID = uuid.New().String() // Ensure entry has a UUID
	// 2. ---------- Add entry to vault ----------
	h.Vault.Entries.Login = append(h.Vault.Entries.Login, *entry)
	h.logger.Info("✅ Added login entry for user %s: %s\n", userID, entry.EntryName)

	return h.Vault, nil
}
func (h *LoginHandler) Edit(userID string, entry any) (*vaults_domain.VaultPayload, error) {
	// 1. ---------- Unmarshal entry ----------
	updatedEntry, ok := entry.(*vaults_domain.LoginEntry)
	if !ok {
		h.logger.Error("LoginHandler - invalid type: expected LoginEntry: %v", entry)
		return nil, fmt.Errorf("invalid type: expected LoginEntry")
	}
	// 2. ---------- Update entry ----------
	updatedEntry.IsDraft = true

	entries := h.Vault.Entries.Login
	updated := false

	for i, entry := range entries {
		if entry.ID == updatedEntry.ID {
			// Update the fields (you could also do a full replace)
			entries[i] = *updatedEntry
			updatedEntry.IsDraft = false
			updatedEntry.CreatedAt = h.NowUTC()
			updatedEntry.UpdatedAt = h.NowUTC()
			updated = true
			break
		}
	}

	if !updated {
		h.logger.Error("LoginHandler - entry with ID %s not found for user %s: %v", updatedEntry.ID, userID, entry)
		return nil, fmt.Errorf("entry with ID %s not found for user %s", updatedEntry.ID, userID)
	}
	// 3. ---------- Update vault ----------
	h.Vault.Entries.Login = entries
	h.logger.Info("✏️ Updated login entry for user %s: %s\n", userID, updatedEntry.EntryName)

	return h.Vault, nil
}
func (h *LoginHandler) Trash(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashLoginEntryAction(userID, entryID, true)
}
func (h *LoginHandler) Restore(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashLoginEntryAction(userID, entryID, false)
}
func (h *LoginHandler) TrashLoginEntryAction(userID string, entryID string, trashed bool) (*vaults_domain.VaultPayload, error) {
	for i, entry := range h.Vault.Entries.Login {
		if entry.ID == entryID {
			h.Vault.Entries.Login[i].Trashed = trashed
			// h.MarkDirty(userID)

			state := "restored"
			if trashed {
				state = "trashed"
			}
			h.logger.Info("🗑️ %s login entry %s for user %s", state, entryID, userID)

			return 	h.Vault, nil	
		}
	}

	return nil, fmt.Errorf("entry with ID %s not found", entryID)
}

func (h *LoginHandler) SetSession(session *vault_session.Session) {
	s := session
	h.Session = s
	payload, err := vault_session.DecodeSessionVault(s.Vault)
	if err != nil {
		return
	}
	h.Vault = payload
	
}

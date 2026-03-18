package vault_ui

import (
	"fmt"
	"time"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	utils "vault-app/internal/utils"
	vault_dto "vault-app/internal/vault/application/dto"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vaults_storage "vault-app/internal/vault/infrastructure/storage"

	"github.com/google/uuid"
)

type LoginHandler struct {
	db     models.DBModel
	logger *logger.Logger
	NowUTC func() string
	Vault  *vaults_domain.VaultPayload
	Session *vault_session.Session
	VaultRepository vaults_domain.VaultRepository
}

func NewLoginHandler(db models.DBModel, log *logger.Logger) *LoginHandler {
	return &LoginHandler{
		db:     db,
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

func (h *LoginHandler) EditWithAttachments(userID string, entry any, attachments []vault_dto.SelectedAttachment) (*vaults_domain.VaultPayload, error) {
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
	entryAttachments := []vaults_domain.Attachment{}

	// 3. ---------- Save attachments ----------
	for _, attachment := range attachments {
		hash, err := h.SaveAttachment(userID, attachment.Data)
		if err != nil {
			h.logger.Error("LoginHandler - SaveAttachment: failed to save attachment: %v", err)
			return nil, err
		}
		entryAttachments = append(entryAttachments, vaults_domain.Attachment{
			ID:   uuid.New().String(),
			EntryID: updatedEntry.ID,
			Hash: hash,
			Name: attachment.Name,
			Size: attachment.Size,
		})
		h.logger.LogPretty("✅ LoginHandler - EditWithAttachment - Attachment saved ", updatedEntry)
	}

	// 4. ---------- Update entry ----------
	for i, entry := range entries {
		if entry.ID == updatedEntry.ID {
			// Update the fields (you could also do a full replace)
			updatedEntry = updatedEntry.AddAttachments(entryAttachments)
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
	// 5. ---------- Update vault ----------
	h.Vault.Entries.Login = entries
	h.logger.Info("✏️ Updated login entry for user %s: %s\n", userID, updatedEntry.EntryName)

	return h.Vault, nil
}

func (h *LoginHandler) SaveAttachment(userID string, data []byte) (string, error) {
	// Get vault
	vault, err := h.VaultRepository.GetByUserIDAndName(userID, h.Vault.Name)
	if err != nil {	
		return "", fmt.Errorf("❌ VaultHandler - SaveAttachment: failed to get vault for user %s: %w", userID, err)
	}
	h.logger.Info("✅ VaultHandler - SaveAttachment: vault retrieved for user %s", userID)

	// Get vault attachement path
	vaultPath := vault.GetVaultPath()
	h.logger.Info("✅ VaultHandler - SaveAttachment: vault path: %s", vaultPath)

	// Create attachment store
	attachmentStore := vaults_storage.NewAttachmentStore(vaultPath)
	h.logger.Info("✅ VaultHandler - SaveAttachment: attachment store created")

	// Save attachment
	hash, err := attachmentStore.Save(data)
	if err != nil {
		return "", fmt.Errorf("❌ VaultHandler - SaveAttachment: failed to save attachment: %w", err)
	}
	h.logger.Info("✅ VaultHandler - SaveAttachment: attachment saved")

	return hash, nil

}

func (h *LoginHandler) SetVaultRepository(vaultRepository vaults_domain.VaultRepository) {
	h.VaultRepository = vaultRepository
}
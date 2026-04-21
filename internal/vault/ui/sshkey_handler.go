package vault_ui

import (
	"fmt"
	"time"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	vault_dto "vault-app/internal/vault/application/dto"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vaults_storage "vault-app/internal/vault/infrastructure/storage"

	"github.com/google/uuid"
)

type SSHKeyHandler struct {
	db     models.DBModel
	logger *logger.Logger
	NowUTC func() string
	Vault  vaults_domain.VaultPayload
	Session *vault_session.Session
	VaultRepository vaults_domain.VaultRepository
	SyncMode bool
}

func NewSSHKeyHandler(db models.DBModel, log *logger.Logger) *SSHKeyHandler {
	return &SSHKeyHandler{
		db:     db,
		logger: log,
		NowUTC: func() string { return time.Now().Format(time.RFC3339) },
	}
}

func (h *SSHKeyHandler) Find(userID string, entryName string) (vaults_domain.VaultEntry, error) {
	for i := range h.Vault.Entries.SSHKey {
		if h.Vault.Entries.SSHKey[i].EntryName == entryName {
			h.logger.Info("🗑️ ssh key entry %s for user %s found", entryName, userID)
			return &h.Vault.Entries.SSHKey[i], nil
		}
	}
	return nil, nil
}

func (h *SSHKeyHandler) Add(userID string, anEntry any) (*vaults_domain.VaultPayload, error) {
	entry, err := anEntry.(*vaults_domain.SSHKeyEntry)
	if !err {
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entry.ID = uuid.New().String() // Ensure entry has a UUID
	h.Vault.Entries.SSHKey = append(h.Vault.Entries.SSHKey, *entry)

	h.logger.Info("✅ Added ssh key entry for user %s: %s\n", userID, entry.EntryName)

	return &h.Vault, nil

}
func (h *SSHKeyHandler) Edit(userID string, entry any) (*vaults_domain.VaultPayload, error) {
	updatedEntry, ok := entry.(*vaults_domain.SSHKeyEntry)
	if !ok {
		return nil, fmt.Errorf("invalid type: expected SSHKeyEntry")
	}

	entries := h.Vault.Entries.SSHKey
	updated := false

	for i, entry := range entries {
		if entry.ID == updatedEntry.ID {
			entries[i] = *updatedEntry
			updatedEntry.IsDraft = true
			updatedEntry.IsDirty = true
			updatedEntry.UpdatedAt = h.NowUTC()
			updated = true
			break
		}
	}

	if !updated {
		return nil, fmt.Errorf("entry with ID %s not found for user %s", updatedEntry.ID, userID)
	}

	h.Vault.Entries.SSHKey = entries
	h.logger.Info("✏️ Updated ssh key entry for user %s: %s\n", userID, updatedEntry.EntryName)

	return &h.Vault, nil
}
func (h *SSHKeyHandler) Trash(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashSSHKeyEntryAction(userID, entryID, true)
}
func (h *SSHKeyHandler) Restore(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashSSHKeyEntryAction(userID, entryID, false)
}
func (h *SSHKeyHandler) TrashSSHKeyEntryAction(userID string, entryID string, trashed bool) (*vaults_domain.VaultPayload, error) {

	for i, entry := range h.Vault.Entries.SSHKey {
		if entry.ID == entryID {
			h.Vault.Entries.SSHKey[i].Trashed = trashed
			h.Vault.Entries.SSHKey[i].IsDirty = true

			state := "restored"
			if trashed {
				state = "trashed"
			}
			h.logger.Info("🗑️ %s ssh key entry %s for user %s", state, entryID, userID)

			return &h.Vault, nil
		}
	}
	return nil, fmt.Errorf("entry with ID %s not found", entryID)
}

func (h *SSHKeyHandler) SetSession(session *vault_session.Session) {
	s := session
	h.Session = s
	payload, err := vault_session.DecodeSessionVault(s.Vault)
	if err != nil {
		return
	}
	h.Vault = *payload
}

func (h *SSHKeyHandler) EditWithAttachments(userID string, entry any, attachments []vault_dto.SelectedAttachment) (*vaults_domain.VaultPayload, error) {
	// 1. ---------- Unmarshal entry ----------
	updatedEntry, ok := entry.(*vaults_domain.SSHKeyEntry)
	if !ok {
		h.logger.Error("SSHKeyHandler - invalid type: expected SSHKeyEntry: %v", entry)
		return nil, fmt.Errorf("invalid type: expected SSHKeyEntry")
	}
	
	// 2. ---------- Update entry ----------
	updatedEntry.IsDraft = true
	updatedEntry.IsDirty = true

	entries := h.Vault.Entries.SSHKey
	updated := false
	entryAttachments := []vaults_domain.Attachment{}

	// 3. ---------- Save attachments ----------
	for _, attachment := range attachments {
		hash, err := h.SaveAttachment(userID, attachment.Data)
		if err != nil {
			h.logger.Error("SSHKeyHandler - SaveAttachment: failed to save attachment: %v", err)
			return nil, err
		}
		entryAttachments = append(entryAttachments, vaults_domain.Attachment{
			ID:   uuid.New().String(),
			EntryID: updatedEntry.ID,
			Hash: hash,
			Name: attachment.Name,
			Size: attachment.Size,
			Ext: attachment.Ext,
		})
		h.logger.LogPretty("✅ SSHKeyHandler - EditWithAttachment - Attachment saved ", updatedEntry)
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
		h.logger.Error("SSHKeyHandler - entry with ID %s not found for user %s: %v", updatedEntry.ID, userID, entry)
		return nil, fmt.Errorf("entry with ID %s not found for user %s", updatedEntry.ID, userID)
	}
	// 5. ---------- Update vault ----------
	h.Vault.Entries.SSHKey = entries
	h.logger.Info("✏️ Updated ssh key entry for user %s: %s\n", userID, updatedEntry.EntryName)

	return &h.Vault, nil
}
func (h *SSHKeyHandler) SaveAttachment(userID string, data []byte) (string, error) {
	// Get vault
	vault, err := h.VaultRepository.GetVault(h.Session.Runtime.VaultID)
	if err != nil {
		return "", fmt.Errorf("❌ SSHKeyHandler - SaveAttachment: failed to get vault for user %s: %w", userID, err)
	}
	h.logger.Info("✅ SSHKeyHandler - SaveAttachment: vault retrieved for user %s", userID)

	// Get vault attachement path
	vaultPath := vault.GetVaultAttachmentPath()
	h.logger.Info("✅ SSHKeyHandler - SaveAttachment: vault path: %s", vaultPath)

	// Create attachment store
	attachmentStore := vaults_storage.NewAttachmentStore(vaultPath)
	h.logger.Info("✅ SSHKeyHandler - SaveAttachment: attachment store created")

	// Save attachment
	hash, err := attachmentStore.Save(data)
	if err != nil {
		return "", fmt.Errorf("❌ SSHKeyHandler - SaveAttachment: failed to save attachment: %w", err)
	}
	h.logger.Info("✅ SSHKeyHandler - SaveAttachment: attachment saved")

	return hash, nil

}

func (h *SSHKeyHandler) SetVaultRepository(vaultRepository vaults_domain.VaultRepository) {
	h.VaultRepository = vaultRepository
}

func (h *SSHKeyHandler) SetSyncMode(b bool) {
	h.SyncMode = b
}
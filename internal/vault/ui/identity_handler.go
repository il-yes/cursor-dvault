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

type IdentityHandler struct {
	db     models.DBModel
	logger *logger.Logger
	NowUTC func() string
	Vault  *vaults_domain.VaultPayload
	Session *vault_session.Session
	VaultRepository vaults_domain.VaultRepository
	SyncMode bool
}

func NewIdentityHandler(db models.DBModel, log *logger.Logger) *IdentityHandler {
	return &IdentityHandler{
		db:     db,
		logger: log,
		NowUTC: func() string { return time.Now().Format(time.RFC3339) },
	}
}

func (h *IdentityHandler) Add(userID string, anEntry any) (*vaults_domain.VaultPayload, error) {
	if h.Vault == nil {
		return nil, fmt.Errorf("no active session for user %s", userID)
	}
	entry, err := anEntry.(*vaults_domain.IdentityEntry)
	if !err {
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entry.ID = uuid.New().String() // Ensure entry has a UUID
	h.Vault.Entries.Identity = append(h.Vault.Entries.Identity, *entry)

	h.logger.Info("✅ Added identity entry for user %s: %s\n", userID, entry.EntryName)

	return h.Vault, nil

}
func (h *IdentityHandler) Edit(userID string, entry any) (*vaults_domain.VaultPayload, error) {
	if h.Vault == nil {
		return nil, fmt.Errorf("no active session for user %s", userID)
	}
	updatedEntry, ok := entry.(*vaults_domain.IdentityEntry)
	if !ok {
		return nil, fmt.Errorf("invalid type: expected IdentityEntry")
	}

	entries := h.Vault.Entries.Identity
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

	h.Vault.Entries.Identity = entries

	h.logger.Info("✏️ Updated identity entry for user %s: %s\n", userID, updatedEntry.EntryName)

	return h.Vault, nil
}
func (h *IdentityHandler) Trash(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashIdentityEntryAction(userID, entryID, true)
}
func (h *IdentityHandler) Restore(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashIdentityEntryAction(userID, entryID, false)
}
func (h *IdentityHandler) TrashIdentityEntryAction(userID string, entryID string, trashed bool) (*vaults_domain.VaultPayload, error) {
	if h.Vault == nil {
		return nil, fmt.Errorf("no active session for user %s", userID)
	}

	for i, entry := range h.Vault.Entries.Identity {
		if entry.ID == entryID {
			h.Vault.Entries.Identity[i].Trashed = trashed
			h.Vault.Entries.Identity[i].IsDirty = true

			state := "restored"
			if trashed {
				state = "trashed"
			}
			h.logger.Info("🗑️ %s identity entry %s for user %s", state, entryID, userID)

			return h.Vault, nil
		}
	}
	return nil, fmt.Errorf("entry with ID %s not found", entryID)
}

func (h *IdentityHandler) SetSession(session *vault_session.Session) {
	s := session
	h.Session = s
	payload, err := vault_session.DecodeSessionVault(s.Vault)
	if err != nil {
		return
	}
	h.Vault = payload
}

func (h *IdentityHandler) EditWithAttachments(userID string, entry any, attachments []vault_dto.SelectedAttachment) (*vaults_domain.VaultPayload, error) {
	// 1. ---------- Unmarshal entry ----------
	updatedEntry, ok := entry.(*vaults_domain.IdentityEntry)
	if !ok {
		h.logger.Error("IdentityHandler - invalid type: expected IdentityEntry: %v", entry)
		return nil, fmt.Errorf("invalid type: expected IdentityEntry")
	}
	
	// 2. ---------- Update entry ----------
	updatedEntry.IsDraft = true
	updatedEntry.IsDirty = true

	entries := h.Vault.Entries.Identity
	updated := false
	entryAttachments := []vaults_domain.Attachment{}

	// 3. ---------- Save attachments ----------
	for _, attachment := range attachments {
		hash, err := h.SaveAttachment(userID, attachment.Data)
		if err != nil {
			h.logger.Error("IdentityHandler - SaveAttachment: failed to save attachment: %v", err)
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
		h.logger.LogPretty("✅ IdentityHandler - EditWithAttachment - Attachment saved ", updatedEntry)
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
		h.logger.Error("CardHandler - entry with ID %s not found for user %s: %v", updatedEntry.ID, userID, entry)
		return nil, fmt.Errorf("entry with ID %s not found for user %s", updatedEntry.ID, userID)
	}
	// 5. ---------- Update vault ----------
	h.Vault.Entries.Identity = entries
	h.logger.Info("✏️ Updated identity entry for user %s: %s\n", userID, updatedEntry.EntryName)

	return h.Vault, nil
}
func (h *IdentityHandler) SaveAttachment(userID string, data []byte) (string, error) {
	// Get vault
	vault, err := h.VaultRepository.GetVault(h.Session.Runtime.VaultID)
	if err != nil {
		return "", fmt.Errorf("❌ IdentityHandler - SaveAttachment: failed to get vault for user %s: %w", userID, err)
	}
	h.logger.Info("✅ IdentityHandler - SaveAttachment: vault retrieved for user %s", userID)

	// Get vault attachement path
	vaultPath := vault.GetVaultAttachmentPath()
	h.logger.Info("✅ IdentityHandler - SaveAttachment: vault path: %s", vaultPath)

	// Create attachment store
	attachmentStore := vaults_storage.NewAttachmentStore(vaultPath)
	h.logger.Info("✅ IdentityHandler - SaveAttachment: attachment store created")

	// Save attachment
	hash, err := attachmentStore.Save(data)
	if err != nil {
		return "", fmt.Errorf("❌ IdentityHandler - SaveAttachment: failed to save attachment: %w", err)
	}
	h.logger.Info("✅ IdentityHandler - SaveAttachment: attachment saved")

	return hash, nil

}	

func (h *IdentityHandler) SetVaultRepository(vaultRepository vaults_domain.VaultRepository) {
	h.VaultRepository = vaultRepository
}

func (h *IdentityHandler) SetSyncMode(b bool) {
	h.SyncMode = b
}

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

type CardHandler struct {
	db     models.DBModel
	logger *logger.Logger
	NowUTC func() string
	Vault  vaults_domain.VaultPayload
	Session *vault_session.Session
	VaultRepository vaults_domain.VaultRepository
}

func NewCardHandler(db models.DBModel, log *logger.Logger) *CardHandler {
	return &CardHandler{
		db:     db,
		logger: log,
		NowUTC: func() string { return time.Now().Format(time.RFC3339) },
	}
}

func (h *CardHandler) Add(userID string, anEntry any) (*vaults_domain.VaultPayload, error) {
	entry, err := anEntry.(*vaults_domain.CardEntry)
	if !err {
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entry.ID = uuid.New().String() // Ensure entry has a UUID
	h.Vault.Entries.Card = append(h.Vault.Entries.Card, *entry)
	// session.LastUpdated = h.NowUTC()
	// session.Dirty = true

	h.logger.Info("✅ Added card entry for user %d: %s\n", userID, entry.EntryName)

	return &h.Vault, nil

}
func (h *CardHandler) Edit(userID string, entry any) (*vaults_domain.VaultPayload, error) {
	updatedEntry, ok := entry.(*vaults_domain.CardEntry)
	if !ok {
		return nil, fmt.Errorf("invalid type: expected CardEntry")
	}
	updatedEntry.IsDraft = true

	entries := h.Vault.Entries.Card
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

	h.Vault.Entries.Card = entries
	// h.MarkDirty(userID)

	h.logger.Info("✏️ Updated card entry for user %d: %s\n", userID, updatedEntry.EntryName)

	return &h.Vault, nil
}
func (h *CardHandler) Trash(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashCardEntryAction(userID, entryID, true)
}
func (h *CardHandler) Restore(userID string, entryID string) (*vaults_domain.VaultPayload, error) {
	return h.TrashCardEntryAction(userID, entryID, false)
}
func (h *CardHandler) TrashCardEntryAction(userID string, entryID string, trashed bool) (*vaults_domain.VaultPayload, error) {
	for i, entry := range h.Vault.Entries.Card {
		if entry.ID == entryID {
			h.Vault.Entries.Card[i].Trashed = trashed

			state := "restored"
			if trashed {
				state = "trashed"
			}
			h.logger.Info("🗑️ %s card entry %s for user %d", state, entryID, userID)

			return &h.Vault, nil
		}
	}
	return nil, fmt.Errorf("entry with ID %s not found", entryID)
}
func (h *CardHandler) SetVault(vault *vault_session.Session) {
	p := vault.Vault
	payload, err := vault_session.DecodeSessionVault(p)
	if err != nil {
		return
	}
	h.Vault = *payload
}
func (h *CardHandler) SetSession(session *vault_session.Session) {
	s := session
	h.Session = s
	payload, err := vault_session.DecodeSessionVault(s.Vault)
	if err != nil {
		return
	}
	h.Vault = *payload
}
func (h *CardHandler) EditWithAttachments(userID string, entry any, attachments []vault_dto.SelectedAttachment) (*vaults_domain.VaultPayload, error) {
	// 1. ---------- Unmarshal entry ----------
	updatedEntry, ok := entry.(*vaults_domain.CardEntry)
	if !ok {
		h.logger.Error("CardHandler - invalid type: expected CardEntry: %v", entry)
		return nil, fmt.Errorf("invalid type: expected CardEntry")
	}
	
	// 2. ---------- Update entry ----------
	updatedEntry.IsDraft = true

	entries := h.Vault.Entries.Card
	updated := false
	entryAttachments := []vaults_domain.Attachment{}

	// 3. ---------- Save attachments ----------
	for _, attachment := range attachments {
		hash, err := h.SaveAttachment(userID, attachment.Data)
		if err != nil {
			h.logger.Error("CardHandler - SaveAttachment: failed to save attachment: %v", err)
			return nil, err
		}
		entryAttachments = append(entryAttachments, vaults_domain.Attachment{
			ID:   uuid.New().String(),
			EntryID: updatedEntry.ID,
			Hash: hash,
			Name: attachment.Name,
			Size: attachment.Size,
		})
		h.logger.LogPretty("✅ CardHandler - EditWithAttachment - Attachment saved ", updatedEntry)
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
	h.Vault.Entries.Card = entries
	h.logger.Info("✏️ Updated card entry for user %s: %s\n", userID, updatedEntry.EntryName)

	return &h.Vault, nil
}	

func (h *CardHandler) SaveAttachment(userID string, data []byte) (string, error) {
	// Get vault
	vault, err := h.VaultRepository.GetByUserIDAndName(userID, h.Vault.Name)
	if err != nil {
		return "", fmt.Errorf("❌ CardHandler - SaveAttachment: failed to get vault for user %s: %w", userID, err)
	}
	h.logger.Info("✅ CardHandler - SaveAttachment: vault retrieved for user %s", userID)

	// Get vault attachement path
	vaultPath := vault.GetVaultPath()
	h.logger.Info("✅ CardHandler - SaveAttachment: vault path: %s", vaultPath)

	// Create attachment store
	attachmentStore := vaults_storage.NewAttachmentStore(vaultPath)
	h.logger.Info("✅ CardHandler - SaveAttachment: attachment store created")

	// Save attachment
	hash, err := attachmentStore.Save(data)
	if err != nil {
		return "", fmt.Errorf("❌ CardHandler - SaveAttachment: failed to save attachment: %w", err)
	}
	h.logger.Info("✅ CardHandler - SaveAttachment: attachment saved")

	return hash, nil

}

func (h *CardHandler) SetVaultRepository(vaultRepository vaults_domain.VaultRepository) {
	h.VaultRepository = vaultRepository
}
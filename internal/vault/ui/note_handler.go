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

type NoteHandler struct {
	db     models.DBModel
	logger *logger.Logger
	NowUTC func() string
	Vault  vaults_domain.VaultPayload
	Session *vault_session.Session
	VaultRepository vaults_domain.VaultRepository	
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

func (h *NoteHandler) EditWithAttachments(userID string, entry any, attachments []vault_dto.SelectedAttachment) (*vaults_domain.VaultPayload, error) {
	// 1. ---------- Unmarshal entry ----------
	updatedEntry, ok := entry.(*vaults_domain.NoteEntry)
	if !ok {
		h.logger.Error("NoteHandler - invalid type: expected NoteEntry: %v", entry)
		return nil, fmt.Errorf("invalid type: expected NoteEntry")
	}
	
	// 2. ---------- Update entry ----------
	updatedEntry.IsDraft = true

	entries := h.Vault.Entries.Note
	updated := false
	entryAttachments := []vaults_domain.Attachment{}

	// 3. ---------- Save attachments ----------
	for _, attachment := range attachments {
		hash, err := h.SaveAttachment(userID, attachment.Data)
		if err != nil {
			h.logger.Error("NoteHandler - SaveAttachment: failed to save attachment: %v", err)
			return nil, err
		}
		entryAttachments = append(entryAttachments, vaults_domain.Attachment{
			ID:   uuid.New().String(),
			EntryID: updatedEntry.ID,
			Hash: hash,
			Name: attachment.Name,
			Size: attachment.Size,
		})
		h.logger.LogPretty("✅ NoteHandler - EditWithAttachment - Attachment saved ", updatedEntry)
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
		h.logger.Error("NoteHandler - entry with ID %s not found for user %s: %v", updatedEntry.ID, userID, entry)
		return nil, fmt.Errorf("entry with ID %s not found for user %s", updatedEntry.ID, userID)
	}
	// 5. ---------- Update vault ----------
	h.Vault.Entries.Note = entries
	h.logger.Info("✏️ Updated note entry for user %s: %s\n", userID, updatedEntry.EntryName)

	return &h.Vault, nil
}
func (h *NoteHandler) SaveAttachment(userID string, data []byte) (string, error) {
	// Get vault
	vault, err := h.VaultRepository.GetByUserIDAndName(userID, h.Vault.Name)
	if err != nil {
		return "", fmt.Errorf("❌ NoteHandler - SaveAttachment: failed to get vault for user %s: %w", userID, err)
	}
	h.logger.Info("✅ NoteHandler - SaveAttachment: vault retrieved for user %s", userID)

	// Get vault attachement path
	vaultPath := vault.GetVaultPath()
	h.logger.Info("✅ NoteHandler - SaveAttachment: vault path: %s", vaultPath)

	// Create attachment store
	attachmentStore := vaults_storage.NewAttachmentStore(vaultPath)
	h.logger.Info("✅ NoteHandler - SaveAttachment: attachment store created")

	// Save attachment
	hash, err := attachmentStore.Save(data)
	if err != nil {
		return "", fmt.Errorf("❌ NoteHandler - SaveAttachment: failed to save attachment: %w", err)
	}
	h.logger.Info("✅ NoteHandler - SaveAttachment: attachment saved")

	return hash, nil

}

func (h *NoteHandler) SetVaultRepository(vaultRepository vaults_domain.VaultRepository) {
	h.VaultRepository = vaultRepository
}
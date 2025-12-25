package vault_ui

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
	utils "vault-app/internal"
	"vault-app/internal/blockchain"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	"vault-app/internal/registry"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vaults_persistence "vault-app/internal/vault/infrastructure/persistence"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VaultHandler struct {
	OpenVaultHandler *OpenVaultHandler

	SessionManager   *vault_session.Manager
	FolderRepository vaults_domain.FolderRepository
	VaultRepository  vaults_domain.VaultRepository

	IPFS                *blockchain.IPFSClient
	DB                  models.DBModel
	logger              logger.Logger
	NowUTC              func() string
	EntryRegistry       *registry.EntryRegistry
	VaultRuntimeContext vault_session.RuntimeContext

	// optionally keep in-memory pending commits per user
	SessionsMu sync.Mutex

	Ctx      context.Context
	Sessions map[string]*vault_session.Session // Offline vault session
}

func NewVaultHandler(
	openVaultHandler *OpenVaultHandler,
	folderRepository vaults_domain.FolderRepository,
	vaultRepository vaults_domain.VaultRepository,
	entryRegistry *registry.EntryRegistry,
	logger logger.Logger,
	ctx context.Context,
	ipfs *blockchain.IPFSClient,
	db *gorm.DB,
	sessions map[string]*vault_session.Session,
) *VaultHandler {
	sessionVaultRepo := vaults_persistence.NewGormSessionRepository(db)
	sessionManager := vault_session.NewManager(sessionVaultRepo, vaultRepository, &logger, ctx, ipfs, sessions)

	return &VaultHandler{
		OpenVaultHandler: openVaultHandler,
		FolderRepository: folderRepository,
		EntryRegistry:    entryRegistry,
		SessionManager:   sessionManager,
		Sessions:         sessions,
	}
}


// -----------------------------
// Session Management
// -----------------------------
func (vh *VaultHandler) PrepareSession(userID string) (*vault_session.Session, error) {
	// Fetch user session if exist, if not create a new one
	session, err := vh.SessionManager.Prepare(userID)
	if err != nil {
		vh.logger.Error("❌ VaultHandler v1 - PrepareSession - failed to prepare session for user %s: %v", userID, err)
		return nil, err
	}
	vh.logger.Info("✅ VaultHandler v1 - PrepareSession - prepared session for user %s", userID)
	vh.Sessions[userID] = session

	return session, nil
}
func (vh *VaultHandler) HasSession() bool {
	return vh.SessionManager.HasSession()
}
func (vh *VaultHandler) GetSession(userID string) (*vault_session.Session, error) {
	vh.SessionManager.SetSessions(vh.Sessions)
	return vh.SessionManager.GetSession(userID)
}
func (vh *VaultHandler) GetAllSessions() map[string]*vault_session.Session {
	return vh.SessionManager.GetSessions()
}
func (vh *VaultHandler) SaveSession(userID string) error {
	session, err := vh.SessionManager.GetSession(userID)
	if err != nil {
		return err
	}
	if err := vh.SessionManager.SessionRepository.SaveSession(userID, session); err != nil {
		vh.logger.Error("❌ VaultHandler v1 - Failed to save session for user %s: %v", userID, err)
		return err
	}
	vh.logger.Info("✅ VaultHandler v1 - Saved session for user %s", userID)
	return nil
}
func (vh *VaultHandler) IsMarkedDirty(userID string) bool {
	return vh.SessionManager.IsMarkedDirty(userID)
}
func (vh *VaultHandler) UpdateVaultPayload(userID string, vs *vaults_domain.VaultPayload) {
	vh.SessionsMu.Lock()
	vh.Sessions[userID].Vault = vs
	vh.SessionsMu.Unlock()
	vh.SessionManager.MarkDirty(userID)
}

// -----------------------------
// Vault - Crud
// -----------------------------
// AddEntryFor: resilient behaviour: will always add the entry; Tracecore commit best-effort
func (vh *VaultHandler) AddEntryFor(userID string, entry any) (*vaults_domain.VaultPayload, error) {
	defer func() {
		if r := recover(); r != nil {
			vh.logger.Error("🔥 Panic in AddEntryFor: %v\nStack:\n%s", r, debug.Stack())
		}
	}()
	// 1. ---------- Validate entry type ----------
	ve, ok := entry.(models.VaultEntry)
	if !ok {
		return nil, fmt.Errorf("❌ entry does not implement VaultEntry interface")
	}
	entryType := ve.GetTypeName()
	vh.logger.Info("✅ Adding %s entry for user %s", entryType, userID)
	// 2. ---------- Get handler for entry type ----------
	handler, err := vh.EntryRegistry.HandlerFor(entryType)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("AddEntryFor - entry", entry)
	// 3. ---------- Set vault & session to the handler ----------
	session, err := vh.GetSession(userID)
	if err != nil {
		return nil, err
	}
	handler.SetSession(session)
	vh.logger.Info("✅ Set vault for user %s", userID)
	// 4. ---------- Add entry to vault ----------
	created, err := handler.Add(userID, entry) // (vault, new_entry)
	vh.logger.Info("✅ Created %s entry for user %s", entryType, userID)

	// 5. ---------- Update session - Mark session as dirty ----------
	vh.UpdateVaultPayload(userID, created)

	return created, err
}
func (vh *VaultHandler) AddEntry(userID string, entryType string, raw json.RawMessage) (any, error) {
	// 1. ---------- Unmarshal entry ----------
	parsed, err := vh.EntryRegistry.UnmarshalEntry(entryType, raw)
	if err != nil {
		vh.logger.Error("❌ AddEntry - Failed to unmarshal %s entry for user %s: %v", entryType, userID, err)
		return nil, fmt.Errorf("failed to parse %s entry: %w", entryType, err)

	}
	// 2. ---------- Add entry (route to handler) ----------
	return vh.AddEntryFor(userID, parsed) // (vault, new_entry)
}
func (vh *VaultHandler) UpdateEntryFor(userID string, entry any) (any, error) {
	ve, ok := entry.(models.VaultEntry)
	if !ok {
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}

	entryType := ve.GetTypeName()
	handler, err := vh.EntryRegistry.HandlerFor(entryType)
	if err != nil {
		return nil, err
	}
	session, err := vh.GetSession(userID)
	if err != nil {
		return nil, err
	}
	handler.SetSession(session)
	updated, err := handler.Edit(userID, entry)
	if err != nil {
		return nil, err
	}

	vh.logger.Info("✅ Updated %s entry for user %s", entryType, userID)
	vh.SessionManager.MarkDirty(userID)

	return updated, nil
}
func (vh *VaultHandler) UpdateEntry(userID string, entryType string, raw json.RawMessage) (any, error) {
	parsed, err := vh.EntryRegistry.UnmarshalEntry(entryType, raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse entry: %w", err)
	}
	return vh.UpdateEntryFor(userID, parsed)
}
func (vh *VaultHandler) TrashEntryFor(userID string, entry any) error {
	ve, ok := entry.(models.VaultEntry)
	if !ok {
		return fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entryID := ve.GetId()
	entryType := ve.GetTypeName()

	handler, err := vh.EntryRegistry.HandlerFor(entryType)
	if err != nil {
		return fmt.Errorf("failed to find an entryRgistry for %s: %w", entryType, err)
	}
	session, err := vh.GetSession(userID)
	if err != nil {
		return err
	}
	handler.SetSession(session)
	err = handler.Trash(userID, entryID)
	vh.logger.Info("✅ trashed %s entry for user %s", entryType, userID)
	vh.SessionManager.MarkDirty(userID)

	return err
}
func (vh *VaultHandler) RestoreEntryFor(userID string, entry any) error {
	ve, ok := entry.(models.VaultEntry)
	if !ok {
		return fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entryID := ve.GetId()
	entryType := ve.GetTypeName()

	handler, err := vh.EntryRegistry.HandlerFor(entryType)
	if err != nil {
		return err
	}
	err = handler.Restore(userID, entryID)

	vh.logger.Info("✅ restored %s entry for user %s", entryType, userID)
	vh.SessionManager.MarkDirty(userID)

	return err
}
func (vh *VaultHandler) TrashEntry(userID string, entryType string, raw json.RawMessage) (any, error) {
	parsed, err := vh.EntryRegistry.UnmarshalEntry(entryType, raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse entry: %w", err)
	}
	err = vh.TrashEntryFor(userID, parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil // Return the restored entry
}
func (vh *VaultHandler) RestoreEntry(userID string, entryType string, raw json.RawMessage) (any, error) {
	parsed, err := vh.EntryRegistry.UnmarshalEntry(entryType, raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse entry: %w", err)
	}
	err = vh.RestoreEntryFor(userID, parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil // Return the restored entry
}

func (vh *VaultHandler) CreateFolder(userID string, name string) (*vaults_domain.VaultPayload, error) {
	// 1. ---------- Set session to the handler ----------
	vh.SessionManager.SetSessions(vh.Sessions)
	session, err := vh.SessionManager.GetSession(userID)
	if err != nil {
		return nil, fmt.Errorf("no active session for user %s", userID)
	}
	// 2. ---------- Create folder ----------	
	folder := &vaults_domain.Folder{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		IsDraft:   false,
	}
	// 3. ---------- Save folder ----------
	if err := vh.FolderRepository.SaveFolder(folder); err != nil {
		return nil, err
	}
	// 4. ---------- Add folder to vault ----------
	vault := session.Vault
	vault.Folders = append(vault.Folders, *folder)
	// 5. ---------- Update session ----------
	vh.logger.Info("✅ Added %s folder for user %s", folder.Name, userID)
	vh.UpdateVaultPayload(userID, vault)

	return vault, nil
}
func (vh *VaultHandler) GetFoldersByVault(vaultCID string) ([]vaults_domain.Folder, error) {
	return vh.FolderRepository.GetFoldersByVault(vaultCID)
}
func (vh *VaultHandler) UpdateFolder(id string, newName string, isDraft bool) (*vaults_domain.Folder, error) {
	folder, err := vh.FolderRepository.GetFolderById(id)
	if err != nil {
		return nil, err
	}

	folder.Name = newName
	folder.IsDraft = isDraft
	folder.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := vh.FolderRepository.UpdateFolder(folder); err != nil {
		return nil, err
	}
	return folder, nil
}
func (vh *VaultHandler) DeleteFolder(userID string, id string) error {
	session, err := vh.GetSession(userID)
	if err != nil {
		return fmt.Errorf("no active session for user %s", userID)
	}

	// 1. Find the folder
	folder, err := vh.FolderRepository.GetFolderById(id)
	if err != nil {
		return err
	}

	// 2. Move entries in this folder to unsorted
	vault := session.Vault
	moved := vault.MoveEntriesToUnsorted(folder.ID)

	// 3. Persist updates (keeping type safety)
	for _, e := range moved.Login {
		raw, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal entry %s: %w", e.ID, err)
		}

		if _, err := vh.UpdateEntry(userID, string(e.Type), raw); err != nil {
			return fmt.Errorf("failed to update login entry %s: %w", e.ID, err)
		}
	}
	for _, e := range moved.Card {
		raw, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal entry %s: %w", e.ID, err)
		}

		if _, err := vh.UpdateEntry(userID, string(e.Type), raw); err != nil {
			return fmt.Errorf("failed to update card entry %s: %w", e.ID, err)
		}
	}
	for _, e := range moved.Identity {
		raw, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal entry %s: %w", e.ID, err)
		}

		if _, err := vh.UpdateEntry(userID, string(e.Type), raw); err != nil {
			return fmt.Errorf("failed to update identity entry %s: %w", e.ID, err)
		}
	}
	for _, e := range moved.Note {
		raw, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal entry %s: %w", e.ID, err)
		}

		if _, err := vh.UpdateEntry(userID, string(e.Type), raw); err != nil {
			return fmt.Errorf("failed to update note entry %s: %w", e.ID, err)
		}
	}
	for _, e := range moved.SSHKey {
		raw, err := json.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal entry %s: %w", e.ID, err)
		}

		if _, err := vh.UpdateEntry(userID, string(e.Type), raw); err != nil {
			return fmt.Errorf("failed to update sshkey entry %s: %w", e.ID, err)
		}
	}

	// 4. Delete folder in DB
	if err := vh.FolderRepository.DeleteFolder(id); err != nil {
		return err
	}

	// 5. Remove from in-memory vault state
	newFolders := []vaults_domain.Folder{}
	for _, f := range vault.Folders {
		if f.ID != id {
			newFolders = append(newFolders, f)
		}
	}
	vault.Folders = newFolders

	return nil
}





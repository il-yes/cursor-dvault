package vault_ui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
	utils "vault-app/internal"
	"vault-app/internal/blockchain"
	app_config_ui "vault-app/internal/config/ui"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	"vault-app/internal/registry"
	vault_commands "vault-app/internal/vault/application/commands"
	vault_events "vault-app/internal/vault/application/events"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_eventbus "vault-app/internal/vault/infrastructure/eventbus"
	vaults_persistence "vault-app/internal/vault/infrastructure/persistence"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/gorm"
)

type VaultHandler struct {
	DB *gorm.DB

	FolderRepository    vaults_domain.FolderRepository
	VaultRepository     vaults_domain.VaultRepository
	EntryRegistry       *registry.EntryRegistry
	VaultRuntimeContext vault_session.RuntimeContext

	IPFS          *blockchain.IPFSClient
	CryptoService *blockchain.CryptoService
	logger        logger.Logger
	NowUTC        func() string

	SessionManager *vault_session.Manager
	SessionsMu     sync.Mutex

	Ctx      context.Context
	EventBus vault_events.VaultEventBus
}

func NewVaultHandler(
	entriesRegistry *registry.EntryRegistry,
	logger logger.Logger,
	ctx context.Context,
	ipfs *blockchain.IPFSClient,
	crypto *blockchain.CryptoService,
	db *gorm.DB,
) *VaultHandler {
	folderRepo := vaults_persistence.NewGormFolderRepository(db)
	vaultRepo := vaults_persistence.NewGormVaultRepository(db)
	sessionRepo := vaults_persistence.NewGormSessionRepository(db)
	sessionManager := vault_session.NewManager(sessionRepo, vaultRepo, &logger, ctx, ipfs, make(map[string]*vault_session.Session))
	eventBus := vault_infrastructure_eventbus.NewMemoryBus()

	return &VaultHandler{
		DB:               db,
		IPFS:             ipfs,
		CryptoService:    crypto,
		logger:           logger,
		NowUTC:           func() string { return time.Now().UTC().Format(time.RFC3339) },
		FolderRepository: folderRepo,
		EntryRegistry:    entriesRegistry,
		SessionManager:   sessionManager,
		EventBus:         eventBus,
	}
}

// -----------------------------
// Session Management
// -----------------------------
func (vh *VaultHandler) PrepareSession(userID string) (*vault_session.Session, error) {
	session, err := vh.SessionManager.Prepare(userID)
	if err != nil {
		vh.logger.Error(
			"❌ VaultHandler - PrepareSession - user %s: %v",
			userID, err,
		)
		return nil, err
	}

	vh.logger.Info("✅ Session prepared for user %s", userID)
	return session, nil
}

func (vh *VaultHandler) HasSession() bool {
	return vh.SessionManager.HasSession()
}

func (vh *VaultHandler) GetSession(userID string) (*vault_session.Session, error) {
	return vh.SessionManager.GetSession(userID)
}

func (vh *VaultHandler) GetAllSessions() map[string]*vault_session.Session {
	return vh.SessionManager.GetSessions()
}

func (vh *VaultHandler) IsMarkedDirty(userID string) bool {
	return vh.SessionManager.IsMarkedDirty(userID)
}

func (vh *VaultHandler) MarkDirty(userID string) {
	vh.SessionManager.MarkDirty(userID)
}

func (vh *VaultHandler) SessionAttachVault(
	ctx context.Context,
	req vault_commands.AttachVaultRequest,
) error {

	if req.VaultPayload == nil {
		return errors.New("vault payload is nil")
	}

	// 🔒 HARD INVARIANT
	req.VaultPayload.Normalize()

	session, err := vh.SessionManager.AttachVault(
		req.UserID,
		req.VaultPayload,
		req.Runtime,
		req.LastCID,
	)
	if err != nil {
		vh.logger.Error(
			"❌ VaultHandler - AttachVault - user %s: %v",
			req.UserID, err,
		)
		return err
	}	

	vh.logger.Info(
		"✅ Vault attached to session for user %s",
		session.UserID,
	)

	return nil
}
func (vh *VaultHandler) LogoutUser(userID string) error {
	return vh.SessionManager.LogoutUser(userID)
}


// -----------------------------
// Vault - Crud
// -----------------------------
func (vh *VaultHandler) Open(ctx context.Context, req vault_commands.OpenVaultCommand, appConfigHandler app_config_ui.AppConfigHandler) (*vault_commands.OpenVaultResult, error) {
	openHandler := NewOpenVaultHandler(
		vault_commands.NewOpenVaultCommandHandler(vh.DB),
		vh.IPFS,
		vh.CryptoService,
		vh.EventBus,
	)
	res, err := openHandler.OpenVault(ctx, req, appConfigHandler)
	if err != nil {
		return nil, err
	}
	
	return res, nil
}

// AddEntryFor: resilient behaviour: will always add the entry; Tracecore commit best-effort
func (vh *VaultHandler) AddEntryFor(userID string, entry any) (*models.VaultEntry, error) {
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
	// 2.1 ---------- Get handler for entry type ----------
	handler, err := vh.EntryRegistry.HandlerFor(entryType)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("AddEntryFor - entry", entry)
	// 2.2 ---------- Get session ----------
	session, err := vh.GetSession(userID)
	if err != nil {
		return nil, err
	}
	// 2.3 ---------- Add entry to vault ----------
	handler.SetSession(session)
	sessionWithNewEntry, err := handler.Add(userID, entry) // (vault, new_entry)
	vh.logger.Info("✅ Created %s entry for user %s", entryType, userID)
	// 4. ---------- Fire Vault stores entry event ----------
	vh.SessionManager.SetVault(userID, sessionWithNewEntry)
	// vh.VaultRuntimeContext.PublishVaultStoredEntry(userID, created)

	return &ve, err
}
func (vh *VaultHandler) AddEntry(userID string, entryType string, raw json.RawMessage) (*models.VaultEntry, error) {
	// 1. ---------- Unmarshal entry ----------
	parsed, err := vh.EntryRegistry.UnmarshalEntry(entryType, raw)
	if err != nil {
		vh.logger.Error("❌ AddEntry - Failed to unmarshal %s entry for user %s: %v", entryType, userID, err)
		return nil, fmt.Errorf("failed to parse %s entry: %w", entryType, err)

	}
	// 2. ---------- Add entry (route to handler) ----------
	res, err := vh.AddEntryFor(userID, parsed) // (vault, new_entry)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (vh *VaultHandler) UpdateEntryFor(userID string, entry any) (*models.VaultEntry, error) {
	vh.logger.Info("✅ Updating entry for user %s", userID)
	ve, ok := entry.(models.VaultEntry)
	if !ok {
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}
	// 1.1 ---------- Validate entry type ----------
	entryType := ve.GetTypeName()
	vh.logger.Info("✅ Updating %s entry for user %s", entryType, userID)
	// 1.2 ---------- Get handler for entry type ----------
	handler, err := vh.EntryRegistry.HandlerFor(entryType)
	if err != nil {
		return nil, err
	}
	// 1.3 ---------- Get session ----------
	session, err := vh.GetSession(userID)
	if err != nil {
		return nil, err
	}
	handler.SetSession(session)
	utils.LogPretty("VaultHandler - UpdateEntryFor - session Before", session)
	// 1.4 ---------- Update entry ----------
	vaultSessionWithUpdatedEntry, err := handler.Edit(userID, entry)
	if err != nil {
		return nil, err
	}
	// 1.5 ---------- Update session ----------
	vh.SessionManager.SetVault(userID, vaultSessionWithUpdatedEntry)
	utils.LogPretty("VaultHandler - UpdateEntryFor - vaultSessionWithUpdatedEntry", vaultSessionWithUpdatedEntry)
	// 1.6 ---------- Mark session as dirty ----------
	vh.logger.Info("✅ Updated %s entry for user %s", entryType, userID)
	// Fires vault stores Edit entry event
	vh.SessionManager.MarkDirty(userID)
	return &ve, nil
}
func (vh *VaultHandler) UpdateEntry(userID string, entryType string, raw json.RawMessage) (any, error) {
	utils.LogPretty("VaultHandler - UpdateEntry - raw", raw)
	parsed, err := vh.EntryRegistry.UnmarshalEntry(entryType, raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse entry: %w", err)
	}
	return vh.UpdateEntryFor(userID, parsed)
}
func (vh *VaultHandler) TrashEntryFor(userID string, entry any) error {
	// 1. ---------- Validate entry type ----------	
	ve, ok := entry.(models.VaultEntry)
	if !ok {
		return fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entryID := ve.GetId()
	entryType := ve.GetTypeName()
	// 2. ---------- Get handler for entry type ----------
	handler, err := vh.EntryRegistry.HandlerFor(entryType)
	if err != nil {
		return fmt.Errorf("failed to find an entryRgistry for %s: %w", entryType, err)
	}
	// 3. ---------- Get session ----------
	session, err := vh.GetSession(userID)
	if err != nil {
		return err
	}
	handler.SetSession(session)
	// 4. ---------- Trash entry ----------
	sessionWithTrash, err := handler.Trash(userID, entryID)
	if err != nil {
		return err
	}
	vh.logger.Info("✅ trashed %s entry for user %s", entryType, userID)
	// 5. ---------- Update session ----------
	vh.SessionManager.SetVault(userID, sessionWithTrash)
	// 6. ---------- Fires vault stores Trash entry event ----------
	
	return err
}
func (vh *VaultHandler) RestoreEntryFor(userID string, entry any) (*models.VaultEntry, error) {
	// 1. ---------- Validate entry type ----------	
	ve, ok := entry.(models.VaultEntry)
	if !ok {
		return nil, fmt.Errorf("entry does not implement VaultEntry interface")
	}
	entryID := ve.GetId()
	entryType := ve.GetTypeName()

	// 2. ---------- Get handler for entry type ----------
	handler, err := vh.EntryRegistry.HandlerFor(entryType)
	if err != nil {
		return nil, fmt.Errorf("failed to find an entryRgistry for %s: %w", entryType, err)
	}
	// 3. ---------- Restore entry ----------
	sessionWithRestored, err := handler.Restore(userID, entryID)
	if err != nil {
		return nil, err
	}
	vh.logger.Info("✅ restored %s entry for user %s", entryType, userID)
	// 5. ---------- Update session ----------
	vh.SessionManager.SetVault(userID, sessionWithRestored)
	// 6. ---------- Fires vault stores Restore entry event ----------

	return &ve, nil
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
	return &parsed, nil // Return the restored entry
}
func (vh *VaultHandler) RestoreEntry(userID string, entryType string, raw json.RawMessage) (*models.VaultEntry, error) {
	parsed, err := vh.EntryRegistry.UnmarshalEntry(entryType, raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse entry: %w", err)
	}
	res, err := vh.RestoreEntryFor(userID, parsed)
	if err != nil {
		return nil, err
	}
	return res, nil // Return the restored entry
}

func (vh *VaultHandler) CreateFolder(userID string, name string) (*vaults_domain.VaultPayload, error) {
	// 1.1 ---------- Create folder ----------
	folder := &vaults_domain.Folder{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		IsDraft:   false,
	}
	// 1.2 ---------- Save folder ----------
	if err := vh.FolderRepository.SaveFolder(folder); err != nil {
		return nil, err
	}
	// 2.1 ---------- Get session ----------
	session, err := vh.GetSession(userID)
	if err != nil {
		return nil, fmt.Errorf("no active session for user %s", userID)
	}
	// 2.2 ---------- Add folder to vault ----------
	vault := session.Vault
	v, err := vault_session.DecodeSessionVault(vault)			
	if err != nil {
		return nil, err
	}
	v.Folders = append(v.Folders, *folder)
	// 2.3 ---------- Update session ----------
	vh.SessionManager.SetVault(userID, v)
	// 2.4 ---------- Firers vault stores CreateFolder event ----------

	return v, nil
}
func (vh *VaultHandler) GetFoldersByVault(vaultCID string) ([]vaults_domain.Folder, error) {
	return vh.FolderRepository.GetFoldersByVault(vaultCID)
}
func (vh *VaultHandler) UpdateFolder(id string, newName string, isDraft bool) (*vaults_domain.Folder, error) {
	// 1.1 ---------- Get folder ----------
	folder, err := vh.FolderRepository.GetFolderById(id)
	if err != nil {
		return nil, err
	}
	// 1.2 ---------- Update folder ----------
	folder.Name = newName
	folder.IsDraft = isDraft
	folder.UpdatedAt = time.Now().Format(time.RFC3339)
	// 1.3 ---------- Save folder ----------
	if err := vh.FolderRepository.UpdateFolder(folder); err != nil {
		return nil, err
	}
	return folder, nil
}
func (vh *VaultHandler) DeleteFolder(userID string, id string) error {
	// 1.1 ---------- Get session ----------
	session, err := vh.GetSession(userID)
	if err != nil {
		return fmt.Errorf("no active session for user %s", userID)
	}

	// 2.1 ---------- Find the folder ----------
	folder, err := vh.FolderRepository.GetFolderById(id)
	if err != nil {
		return err
	}

	// 2.2 ---------- Move entries in this folder to unsorted ----------
	vault := session.Vault
	v, err := vault_session.DecodeSessionVault(vault)			
	if err != nil {
		return err
	}
	moved := v.MoveEntriesToUnsorted(folder.ID)

	// 3. ---------- Persist updates (keeping type safety) ----------
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

	// 4. ---------- Delete folder in DB ----------
	if err := vh.FolderRepository.DeleteFolder(id); err != nil {
		return err
	}

	// 5. ---------- Remove from in-memory vault state ----------
	newFolders := []vaults_domain.Folder{}
	for _, f := range v.Folders {
		if f.ID != id {
			newFolders = append(newFolders, f)
		}
	}
	v.Folders = newFolders

	// 6. ---------- Update session ----------
	vh.SessionManager.SetVault(userID, v)
	// 7. ---------- Fire vault stores DeleteFolder event ----------	

	return nil
}

func (vh *VaultHandler) SyncVault(userID string, password string) (string, error) {
	vh.logger.Info("🔄 Starting vault sync for UserID: %s", userID)

	// 1. ---------- Get session ----------
	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 10, "stage": "retrieving session"})
	session, err := vh.GetSession(userID)
	if err != nil {
		return "", fmt.Errorf("no active session: %w", err)
	}

	// 2. ---------- Marshal vault ----------
	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 20, "stage": "marshalling vault"})
	vaultBytes, err := json.Marshal(session.Vault)
	if err != nil {
		return "", fmt.Errorf("marshal failed: %w", err)
	}

	// 3. ---------- Encrypt vault ----------
	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 40, "stage": "encrypting vault"})
	encrypted, err := blockchain.Encrypt(vaultBytes, password)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	// 4. ---------- Upload to IPFS ----------
	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 70, "stage": "uploading to IPFS"})
	newCID, err := vh.IPFS.AddData(encrypted)
	if err != nil {
		return "", fmt.Errorf("IPFS upload failed: %w", err)
	}

	// 5. ---------- Submit to Stellar ----------
	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 90, "stage": "submitting to Stellar"})
	userCfg := session.Runtime.UserConfig
	txHash, err := blockchain.SubmitCID(userCfg.StellarAccount.PrivateKey, newCID)
	if err != nil {
		return "", fmt.Errorf("stellar submission failed: %w", err)
	}

	// 6. ---------- Create new vault ----------
	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 95, "stage": "saving metadata"})
	currentMeta, err := vh.VaultRepository.GetLatestByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get vault meta: %w", err)
	}
	newVault := vaults_domain.Vault{
		Name:      currentMeta.Name,
		Type:      currentMeta.Type,
		UserID:    userID,
		CID:       newCID,
		TxHash:    txHash,
		CreatedAt: vh.NowUTC(),
		UpdatedAt: vh.NowUTC(),
	}
	saved := vh.VaultRepository.SaveVault(&newVault)
	vh.logger.Info("💾 Vault saved for user %s: %v", userID, saved)
	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 100, "stage": "complete"})
	
	// 7. ---------- Update session ----------
	vh.SessionManager.Sync(userID, newCID)
	vh.logger.Info("✅ Vault sync complete for user %s", userID)
	// 8. ---------- Emit event ----------
	runtime.EventsEmit(vh.Ctx, "vault-synced", map[string]interface{}{"userID": userID, "newCID": newCID})
	// 9. ---------- Fires vault stores Sync event ----------
	
	return newCID, nil
}

func (vh *VaultHandler) UploadToIPFS(userID string, encrypted string) (string, error) {
	newCID, err := vh.IPFS.AddData([]byte(encrypted))
	if err != nil {
		return "", fmt.Errorf("❌ VaultHandler - UploadToIPFS: failed to upload to IPFS: %w", err)
	}
	vh.logger.Info("📤 VaultHandler - UploadToIPFS: Vault uploaded to IPFS (CID: %s)", newCID)
	// 3. ----------------- Fires vault stores UploadToIPFS event -----------------
	
	return newCID, nil
}

func (vh *VaultHandler) EncryptVault(userID string, password string) (string, error) {
	vh.logger.Info("🔄 VaultHandler - EncryptVault: Starting vault encryption for UserID: %s", userID)

	// 1. ----------------- Get session -----------------
	session, err := vh.GetSession(userID)
	if err != nil {
		return "", fmt.Errorf("❌ VaultHandler - EncryptVault: no active session for user %s: %w", userID, err)
	}
	// 2. ----------------- Marshal in-memory vault -----------------
	vaultBytes, err := json.Marshal(session.Vault) // session.Vault.Sync()
	if err != nil {
		return "", fmt.Errorf("❌ VaultHandler - EncryptVault: failed to marshal vault: %w", err)
	}
	vh.logger.Info("🧱 Vault marshalled (%d bytes)", len(vaultBytes))

	// 3. ----------------- Encrypt -----------------
	encrypted, err := blockchain.Encrypt(vaultBytes, password)
	if err != nil {
		return "", fmt.Errorf("❌ VaultHandler - EncryptVault: failed to encrypt vault: %w", err)
	}
	vh.logger.Info("🔐 VaultHandler - EncryptVault: Vault encrypted")		
	// 4. ----------------- Fires vault stores Encrypt event -----------------

	return string(encrypted), nil
}
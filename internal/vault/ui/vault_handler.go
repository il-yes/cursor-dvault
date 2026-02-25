package vault_ui

import (
	"context"

	// "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
	"vault-app/internal/blockchain"
	app_config "vault-app/internal/config"
	app_config_domain "vault-app/internal/config/domain"
	app_config_ui "vault-app/internal/config/ui"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	"vault-app/internal/registry"
	"vault-app/internal/tracecore"
	utils "vault-app/internal/utils"
	vault_commands "vault-app/internal/vault/application/commands"
	vault_dto "vault-app/internal/vault/application/dto"
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

	InitializeVaultCommandHandler   *vault_commands.InitializeVaultCommandHandler
	CreateIPFSPayloadCommandHandler *vault_commands.CreateIPFSPayloadCommandHandler
	CreateVaultCommandHandler       *vault_commands.CreateVaultCommandHandler
	VaultOpenedListener             *vault_commands.VaultOpenedListener

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
	appConfig app_config.AppConfig,
) *VaultHandler {
	folderRepo := vaults_persistence.NewGormFolderRepository(db)
	vaultRepo := vaults_persistence.NewGormVaultRepository(db)
	sessionRepo := vaults_persistence.NewGormSessionRepository(db)
	sessionManager := vault_session.NewManager(sessionRepo, vaultRepo, &logger, ctx, ipfs, make(map[string]*vault_session.Session))
	eventBus := vault_infrastructure_eventbus.NewMemoryBus()

	initializeVaultHandler := vault_commands.NewInitializeVaultCommandHandler(db)
	createIpfsCommandHandler := vault_commands.NewCreateIPFSPayloadCommandHandler(
		vaultRepo, crypto, ipfs,
	)
	createVaultCommand := vault_commands.NewCreateVaultCommandHandler(
		initializeVaultHandler, createIpfsCommandHandler, vaultRepo,
	)

	return &VaultHandler{
		DB:                              db,
		IPFS:                            ipfs,
		CryptoService:                   crypto,
		logger:                          logger,
		NowUTC:                          func() string { return time.Now().UTC().Format(time.RFC3339) },
		FolderRepository:                folderRepo,
		EntryRegistry:                   entriesRegistry,
		SessionManager:                  sessionManager,
		EventBus:                        eventBus,
		InitializeVaultCommandHandler:   initializeVaultHandler,
		CreateIPFSPayloadCommandHandler: createIpfsCommandHandler,
		CreateVaultCommandHandler:       createVaultCommand,
		VaultRepository:                 vaultRepo,
	}
}

func (vh *VaultHandler) InitializeVaultOpenedListener() {
	vh.logger.Info("Initializing vault opened listener")
	vh.VaultOpenedListener = vault_commands.NewVaultOpenedListener(&vh.logger, vh.EventBus, vh)
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
func (vh *VaultHandler) SessionAttachRuntime(ctx context.Context, req vault_commands.AttachRuntimeRequest) error {
	session, err := vh.SessionManager.AttachRuntime(req.UserID, req.Runtime)
	if err != nil {
		return err
	}
	utils.LogPretty("SessionAttachRuntime - session", session)
	return nil
}
func (vh *VaultHandler) LogoutUser(userID string) error {
	return vh.SessionManager.LogoutUser(userID)
}
func (vh *VaultHandler) GetVaultSession(userID string) (*vaults_domain.VaultPayload, error) {
	session, err := vh.SessionManager.GetSession(userID)
	if err != nil {
		return nil, err
	}
	vaultBytes := session.Vault
	utils.LogPretty("GetVaultSession - vaultPayload", vaultBytes)
	var vaultPayload vaults_domain.VaultPayload
	json.Unmarshal(vaultBytes, &vaultPayload)

	return &vaultPayload, nil
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
	// 4. ---------- Update session ----------
	vh.SessionManager.SetVault(userID, sessionWithNewEntry)
	// 5. ---------- Fire Vault stores entry event ----------
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
	utils.LogPretty("VaultHandler - DeleteFolder - session", session)

	// 2.1 ---------- Find the folder ----------
	folder, err := vh.FolderRepository.GetFolderById(id)
	if err != nil {
		return err
	}
	utils.LogPretty("VaultHandler - DeleteFolder - folder", folder)

	// 2.2 ---------- Move entries in this folder to unsorted ----------
	vault := session.Vault
	v, err := vault_session.DecodeSessionVault(vault)
	if err != nil {
		return err
	}
	moved := v.MoveEntriesToUnsorted(folder.ID)
	fmt.Println("VaultHandler - DeleteFolder - moved")

	// 3. ---------- Persist updates (keeping type safety) ----------
	if len(moved.Login) > 0 {
		fmt.Println("VaultHandler - DeleteFolder - moved entries detected")
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
	}
	fmt.Println("VaultHandler - DeleteFolder - moved entries updated")

	// 4. ---------- Delete folder in DB ----------
	if err := vh.FolderRepository.DeleteFolder(id); err != nil {
		return err
	}
	fmt.Println("VaultHandler - DeleteFolder - deleted folder")

	// 5. ---------- Remove from in-memory vault state ----------
	newFolders := []vaults_domain.Folder{}
	for _, f := range v.Folders {
		if f.ID != id {
			newFolders = append(newFolders, f)
		}
	}
	v.Folders = newFolders
	utils.LogPretty("VaultHandler - DeleteFolder - newFolders", newFolders)
	// 6. ---------- Update session ----------
	vh.SessionManager.SetVault(userID, v)
	// 7. ---------- Fire vault stores DeleteFolder event ----------

	return nil
}

func (vh *VaultHandler) SyncVault(ctx context.Context, input vault_dto.SynchronizeVaultRequest, tc tracecore.TracecoreClient) (string, error) {
	vh.logger.Info("🔄 Starting vault sync for UserID: %s", input.UserID)
	if ctx == nil {
		vh.logger.Error("❌ SyncVault aborted: ctx is nil")
		return "", errors.New("runtime context is nil")
	}
	vh.logger.Info("CTX = %#v", ctx)
	userID := input.UserID
	password := input.Password

	// 1. ---------- Get session ----------
	runtime.EventsEmit(ctx, "progress-update", map[string]interface{}{"percent": 10, "stage": "retrieving session"})
	vh.logger.Info("🔄 SyncVault - Retrieving session for UserID: %s", userID)
	session, err := vh.GetSession(userID)
	if err != nil {
		return "", fmt.Errorf("no active session: %w", err)
	}
	vh.logger.Info("🔄 SyncVault - Session retrieved for UserID: %s", userID)	
	vh.logger.LogPretty("SyncVault - session", session)

	// 2. ---------- Marshal vault ----------
	runtime.EventsEmit(ctx, "progress-update", map[string]interface{}{"percent": 20, "stage": "marshalling vault"})
	vaultBytes, err := json.Marshal(session.Vault)
	if err != nil {
		return "", fmt.Errorf("marshal failed: %w", err)
	}
	vh.logger.LogPretty("SyncVault - vaultBytes", vaultBytes)

	// 3. ---------- Encrypt vault ----------
	runtime.EventsEmit(ctx, "progress-update", map[string]interface{}{"percent": 40, "stage": "encrypting vault"})
	encrypted, err := blockchain.Encrypt(vaultBytes, password)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}
	vh.logger.LogPretty("SyncVault - encrypted", encrypted)

	// 4. ---------- Upload to IPFS ----------
	runtime.EventsEmit(ctx, "progress-update", map[string]interface{}{"percent": 70, "stage": "uploading to IPFS"})
	// newCID, err := vh.IPFS.AddData(encrypted)
	req := tracecore.SyncVaultStreamRequest{
		UserID:  input.Vault.UserSubscriptionID,
		VaultName: input.Vault.Name,
		Stream:  encrypted,
	}
	vh.logger.LogPretty("SyncVault - req ipfs", req)
	res, err := tc.UploadToIPFS(ctx, req)
	if err != nil {
		return "", fmt.Errorf("IPFS upload failed: %w", err)
	}
	vh.logger.Info("🔄 SyncVault - Vault uploaded to IPFS - cid: %s", res.Data.CID)
	newCID := res.Data.CID

	// 5. ---------- Submit to Stellar ----------
	runtime.EventsEmit(ctx, "progress-update", map[string]interface{}{"percent": 90, "stage": "submitting to Stellar"})
	userCfg := session.Runtime.UserConfig
	utils.LogPretty("SyncVault - userCfg", userCfg)

	// passwordCipherBytes, _ := base64.StdEncoding.DecodeString(string(userCfg.StellarAccount.EncPassword))
	// decryptedPassword, err := blockchain.Decrypt(passwordCipherBytes, ANKHORA_SECRET)
	// if err != nil {
	// 	return "", fmt.Errorf("stellar submission failed: %w", err)
	// }

	// cipherBytes, _ := base64.StdEncoding.DecodeString(userCfg.StellarAccount.PrivateKey)
	// decryptedPrivateKey, err := blockchain.Decrypt(cipherBytes, decryptedPassword)
	// if err != nil {
	// 	return "", fmt.Errorf("stellar submission failed: %w", err)
	// }
	// utils.LogPretty("SyncVault - decryptedPrivateKey", decryptedPrivateKey)
	txHash, err := blockchain.SubmitCID(userCfg.StellarAccount.PrivateKey, newCID)
	if err != nil {
		return "", fmt.Errorf("stellar submission failed: %w", err)
	}
	vh.logger.Info("🔄 SyncVault - Vault submitted to Stellar - txHash: %s", txHash)

	fmt.Println("SyncVault - CTX", ctx)
	utils.LogPretty("SyncVault - Vault repository", vh.VaultRepository)

	// 6. ---------- Create new vault ----------
	runtime.EventsEmit(ctx, "progress-update", map[string]interface{}{"percent": 95, "stage": "saving metadata"})
	currentMeta, err := vh.VaultRepository.GetLatestByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get vault meta: %w", err)
	}
	utils.LogPretty("SyncVault - currentMeta", currentMeta)
	newVault := vaults_domain.Vault{
		Name:      currentMeta.Name,
		Type:      currentMeta.Type,
		UserID:    userID,
		CID:       newCID,
		TxHash:    txHash,
		CreatedAt: vh.NowUTC(),
		UpdatedAt: vh.NowUTC(),
	}
	utils.LogPretty("SyncVault - newVault to save", newVault)
	saved := vh.VaultRepository.UpdateVault(&newVault)
	vh.logger.Info("💾 Vault saved for user %s: %v", userID, saved)
	runtime.EventsEmit(ctx, "progress-update", map[string]interface{}{"percent": 100, "stage": "complete"})

	// 7. ---------- Update session ----------
	vh.SessionManager.Sync(userID, newCID)
	vaultPayload, err := vault_session.DecodeSessionVault(vaultBytes)
	if err != nil {
		return "", fmt.Errorf("failed to decode session vault: %w", err)
	}
	vh.logger.LogPretty("SyncVault - vaultPayload", vaultPayload)
	vh.SessionManager.SetVault(userID, vaultPayload)
	vh.logger.Info("✅ Vault sync complete for user %s", userID)

	// 8. ---------- Emit event ----------
	runtime.EventsEmit(ctx, "vault-synced", map[string]interface{}{"userID": userID, "newCID": newCID})

	// 9. ---------- Fires vault stores Sync event ----------

	return newCID, nil
}

func (vh *VaultHandler) AccessEncryptedEntry(ctx context.Context, id string, req tracecore.AccessCryptoShareRequest, tc tracecore.TracecoreClient) (*tracecore.CloudResponse[tracecore.AccessCryptoShareResponse], error) {
	response, err := tc.AccessEncryptedEntry(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access encrypted entry: %w", err)
	}

	return response, nil
}
func (vh *VaultHandler) DecryptVaultEntry(ctx context.Context, req tracecore.DecryptCryptoShareRequest, tc tracecore.TracecoreClient) (*tracecore.CloudResponse[tracecore.DecryptCryptoShareResponse], error) {
	response, err := tc.DecryptVaultEntry(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault entry: %w", err)
	}

	return response, nil
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

type OnGenerateApiKeyParams struct {
	UserID     string
	UserConfig app_config_domain.UserConfig
}

func (vh *VaultHandler) OnGenerateApiKey(ctx context.Context, params OnGenerateApiKeyParams) error {
	vh.logger.Info("🔄 VaultHandler - OnGenerateApiKey: Starting vault encryption for UserID: %s", params.UserID)

	// 3. ----------------- Update UserConfig -----------------
	session, err := vh.GetSession(params.UserID)
	if err != nil {
		vh.logger.Error("❌ VaultHandler - OnGenerateApiKey: no active session for user %s: %v", params.UserID, err)
		return err
	}

	userCfg := app_config.UserConfig{
		ID:            params.UserConfig.ID,
		Role:          params.UserConfig.Role,
		Signature:     params.UserConfig.Signature,
		ConnectedOrgs: params.UserConfig.ConnectedOrgs,
		StellarAccount: app_config.StellarAccountConfig{
			PublicKey:  params.UserConfig.StellarAccount.PublicKey,
			PrivateKey: params.UserConfig.StellarAccount.PrivateKey,
		},
		TwoFactorEnabled: params.UserConfig.TwoFactorEnabled,
	}
	session.Runtime.UserConfig = userCfg
	vh.SessionManager.AttachRuntime(params.UserID, session.Runtime)

	return nil
}

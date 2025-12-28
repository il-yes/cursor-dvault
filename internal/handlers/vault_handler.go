package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	// "os"
	"runtime/debug"
	"sync"
	"time"
	utils "vault-app/internal"
	share_application_events "vault-app/internal/application/events/share"
	share_application_use_cases "vault-app/internal/application/use_cases"
	"vault-app/internal/blockchain"
	share_domain "vault-app/internal/domain/shared"
	share_infrastructure "vault-app/internal/infrastructure/share"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	"vault-app/internal/registry"
	"vault-app/internal/services"
	"vault-app/internal/tracecore"
	vault_session "vault-app/internal/vault/application/session"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type VaultHandler struct {
	Sessions            map[string]*models.VaultSession
	IPFS                *blockchain.IPFSClient
	DB                  models.DBModel
	vaultDirty          bool
	logger              logger.Logger
	NowUTC              func() string
	EntryRegistry       *registry.EntryRegistry
	TracecoreClient     *tracecore.TracecoreClient
	VaultRuntimeContext vault_session.RuntimeContext

	pendingMu sync.Mutex
	// optionally keep in-memory pending commits per user
	pendingCommits map[string][]tracecore.CommitEnvelope
	SessionsMu     sync.Mutex

	EventDispatcher share_application_events.EventDispatcher
	Ctx             context.Context
}

func NewVaultHandler(
	db models.DBModel,
	ipfs *blockchain.IPFSClient,
	registry *registry.EntryRegistry,
	Sessions map[string]*models.VaultSession,
	logger *logger.Logger,
	tc *tracecore.TracecoreClient,
	runtimeCtx vault_session.RuntimeContext,
) *VaultHandler {
	if tc == nil {
		logger.Error("❌ TracecoreClient is nil when initializing VaultHandler!")
	}

	return &VaultHandler{
		DB:                  db,
		Sessions:            Sessions,
		IPFS:                ipfs,
		logger:              *logger,
		NowUTC:              func() string { return time.Now().Format(time.RFC3339) },
		EntryRegistry:       registry,
		TracecoreClient:     tc,
		VaultRuntimeContext: runtimeCtx,
		pendingCommits:      make(map[string][]tracecore.CommitEnvelope),
		EventDispatcher:     share_infrastructure.InitializeEventDispatcher(),
	}
}

// -----------------------------
// Vault - Session
// -----------------------------
func (vh *VaultHandler) StartSession(userID string, vault models.VaultPayload, lastCID string, ctx *models.VaultRuntimeContext) {
	now := time.Now().Format(time.RFC3339)
	vh.Sessions[userID] = &models.VaultSession{
		UserID:              userID,
		Vault:               &vault,
		LastCID:             lastCID,
		LastSynced:          now,
		LastUpdated:         now,
		Dirty:               false,
		VaultRuntimeContext: *ctx,
	}
}
func (vh *VaultHandler) GetSession(userID string) (*models.VaultSession, error) {
	session, ok := vh.Sessions[userID]
	if !ok {
		return nil, errors.New("no vault session found")
	}
	return session, nil
}
func (vh *VaultHandler) EndSession(userID string) {
	if session, ok := vh.Sessions[userID]; ok {
		// utils.LogPretty("ssession saved", session)
		// Persist before deleting
		if err := vh.DB.SaveSession(userID, session); err != nil {
			vh.logger.Error("❌ Handlers v0 - failed to save session for user %s: %v", userID, err)
		} else {
			vh.logger.Info("💾 Session saved for user %s", userID)
		}
	}

	delete(vh.Sessions, userID)
}
func (vh *VaultHandler) LogoutUser(userID string) error {
	vh.SessionsMu.Lock()
	session, ok := vh.Sessions[userID]
	vh.SessionsMu.Unlock()

	if !ok {
		return fmt.Errorf("no active session for user %d", userID)
	}

	// Persist to DB
	if err := vh.DB.SaveSession(userID, session); err != nil {
		return fmt.Errorf("failed to save session for user %d: %w", userID, err)
	}

	vh.pendingMu.Lock()
	delete(vh.pendingCommits, userID)
	vh.pendingMu.Unlock()

	vh.SessionsMu.Lock()
	delete(vh.Sessions, userID)
	vh.SessionsMu.Unlock()

	vh.logger.Info("👋 User %s logged out and session saved", userID)
	return nil
}

func (vh *VaultHandler) MarkDirty(userID string) {
	if session, err := vh.GetSession(userID); err == nil {
		session.LastUpdated = utils.NowUTCString()
		session.Dirty = true
		vh.vaultDirty = true
	}
}
func (vh *VaultHandler) IsVaultDirty() bool {
	return vh.vaultDirty
}

func (vh *VaultHandler) SyncVault0(userID string, password string) (string, error) {
	vh.logger.Info("🔄 Starting vault sync for UserID: %s", userID)

	// 1. Get session
	session, err := vh.GetSession(userID)
	if err != nil {
		return "", fmt.Errorf("❌ no active session for user %s: %w", userID, err)	
	}
	// ✅ Removed noisy LogPretty - too verbose for production
	// 2. Marshal in-memory vault
	vaultBytes, err := json.Marshal(session.Vault) // session.Vault.Sync()
	if err != nil {
		return "", fmt.Errorf("❌ failed to marshal vault: %w", err)
	}
	vh.logger.Info("🧱 Vault marshalled (%d bytes)", len(vaultBytes))

	// 3. Encrypt
	encrypted, err := blockchain.Encrypt(vaultBytes, password)
	if err != nil {
		return "", fmt.Errorf("❌ failed to encrypt vault: %w", err)
	}
	vh.logger.Info("🔐 Vault encrypted")

	// GetBackendPlanParamForTransaction for managing plans from remote

	// 4. Upload to IPFS
	newCID, err := vh.IPFS.AddData(encrypted)
	if err != nil {
		return "", fmt.Errorf("❌ failed to upload to IPFS: %w", err)
	}
	vh.logger.Info("📤 Vault uploaded to IPFS (CID: %s)", newCID)

	// 5. Submit to Stellar
	userCfg := session.VaultRuntimeContext.CurrentUser
	txHash, err := blockchain.SubmitCID(userCfg.StellarAccount.PrivateKey, newCID)
	if err != nil {
		return "", fmt.Errorf("❌ failed to submit CID to Stellar: %w", err)
	}
	vh.logger.Info("🌐 CID submitted to Stellar (TX: %s)", txHash)

	// 5. Decrypt
	/*
	decrypted, err := blockchain.Decrypt([]byte(userCfg.StellarAccount.PrivateKey), os.Getenv("TRACECORE_PRIVATE_KEY"))
	if err != nil {
		return "", fmt.Errorf("❌ failed to decrypt vault: %w", err)
	}
	vh.logger.Info("🔓 Vault decrypted")
	// 6. Submit to Stellar
	txHash, errTrx := blockchain.SubmitCID(string(decrypted), newCID)
	if errTrx != nil {
		return "", fmt.Errorf("❌ failed to submit CID to Stellar: %w", errTrx)
	}
	*/

	// 6. Get latest metadata
	currentMeta, err := vh.DB.GetLatestVaultCIDByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("❌ failed to get current vault meta: %w", err)
	}

	// 7. Save new VaultCID
	newVault := models.VaultCID{
		Name:      currentMeta.Name,
		Type:      currentMeta.Type,
		UserID:    userID,
		CID:       newCID,
		TxHash:    txHash,
		CreatedAt: vh.NowUTC(),
		UpdatedAt: vh.NowUTC(),
	}
	if _, err := vh.DB.SaveVaultCID(newVault); err != nil {
		return "", fmt.Errorf("❌ failed to save new vault CID: %w", err)
	}
	vh.logger.Info("🗃️ VaultCID saved")

	// 8. Update session
	session.LastCID = newCID
	session.LastSynced = time.Now().Format(time.RFC3339)
	session.Dirty = false
	vh.vaultDirty = false

	vh.logger.Info("✅ Vault sync complete for user %d", userID)
	// utils.LogPretty("session after sync", session)

	return newCID, nil
}
func (vh *VaultHandler) SyncVault(userID string, password string) (string, error) {
	vh.logger.Info("🔄 Starting vault sync for UserID: %s", userID)

	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 10, "stage": "retrieving session"})
	session, err := vh.GetSession(userID)
	if err != nil {
		return "", fmt.Errorf("no active session: %w", err)
	}

	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 20, "stage": "marshalling vault"})
	vaultBytes, err := json.Marshal(session.Vault)
	if err != nil {
		return "", fmt.Errorf("marshal failed: %w", err)
	}

	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 40, "stage": "encrypting vault"})
	encrypted, err := blockchain.Encrypt(vaultBytes, password)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 70, "stage": "uploading to IPFS"})
	newCID, err := vh.IPFS.AddData(encrypted)
	if err != nil {
		return "", fmt.Errorf("IPFS upload failed: %w", err)
	}

	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 90, "stage": "submitting to Stellar"})
	userCfg := session.VaultRuntimeContext.CurrentUser
	txHash, err := blockchain.SubmitCID(userCfg.StellarAccount.PrivateKey, newCID)
	if err != nil {
		return "", fmt.Errorf("stellar submission failed: %w", err)
	}

	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 95, "stage": "saving metadata"})
	currentMeta, err := vh.DB.GetLatestVaultCIDByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get vault meta: %w", err)
	}
	newVault := models.VaultCID{
		Name:      currentMeta.Name,
		Type:      currentMeta.Type,
		UserID:    userID,
		CID:       newCID,
		TxHash:    txHash,
		CreatedAt: vh.NowUTC(),
		UpdatedAt: vh.NowUTC(),
	}
	if _, err := vh.DB.SaveVaultCID(newVault); err != nil {
		return "", fmt.Errorf("failed to save vault CID: %w", err)
	}

	runtime.EventsEmit(vh.Ctx, "progress-update", map[string]interface{}{"percent": 100, "stage": "complete"})

	// Update session
	session.LastCID = newCID
	session.LastSynced = time.Now().Format(time.RFC3339)
	session.Dirty = false
	vh.vaultDirty = false

	vh.logger.Info("✅ Vault sync complete for user %d", userID)
	return newCID, nil
}

func (vh *VaultHandler) EncryptFile(userID string, filePath []byte, password string) (string, error) {
	vh.logger.Info("🔄 Starting vault sync for UserID: %s", userID)

	// 1. Get session
	session, err := vh.GetSession(userID)
	if err != nil {
		return "", fmt.Errorf("❌ no active session for user %d: %w", userID, err)
	}
	// ✅ Removed noisy LogPretty - too verbose for production
	// 2. Marshal in-memory vault
	vaultBytes, err := json.Marshal(session.Vault) // session.Vault.Sync()
	if err != nil {
		return "", fmt.Errorf("❌ failed to marshal vault: %w", err)
	}
	vh.logger.Info("🧱 Vault marshalled (%d bytes)", len(vaultBytes))

	// 3. Encrypt
	encrypted, err := blockchain.Encrypt(vaultBytes, password)
	if err != nil {
		return "", fmt.Errorf("❌ failed to encrypt vault: %w", err)
	}
	vh.logger.Info("🔐 Vault encrypted")

	return string(encrypted), nil
}
func (vh *VaultHandler) UploadToIPFS(userID string, encrypted string) (string, error) {
	// GetBackendPlanParamForTransaction for managing plans from remote

	// Upload to IPFS
	newCID, err := vh.IPFS.AddData([]byte(encrypted))
	if err != nil {
		return "", fmt.Errorf("❌ failed to upload to IPFS: %w", err)
	}
	vh.logger.Info("📤 Vault uploaded to IPFS (CID: %s)", newCID)
	return newCID, nil
}
func (vh *VaultHandler) CreateStellarCommit(userID string, newCID string) (string, error) {
	// 1. Get session
	session, err := vh.GetSession(userID)
	if err != nil {
		return "", fmt.Errorf("❌ no active session for user %s: %w", userID, err)
	}

	userCfg := session.VaultRuntimeContext.CurrentUser
	txHash, err := blockchain.SubmitCID(userCfg.StellarAccount.PrivateKey, newCID)
	if err != nil {
		return "", fmt.Errorf("❌ failed to submit CID to Stellar: %w", err)
	}
	vh.logger.Info("🌐 CID submitted to Stellar (TX: %s)", txHash)

	// 6. Get latest metadata
	currentMeta, err := vh.DB.GetLatestVaultCIDByUserID(userID)
	if err != nil {
		return "", fmt.Errorf("❌ failed to get current vault meta: %w", err)
	}
	newVault := models.VaultCID{
		Name:      currentMeta.Name,
		Type:      currentMeta.Type,
		UserID:    userID,
		CID:       newCID,
		TxHash:    txHash,
		CreatedAt: vh.NowUTC(),
		UpdatedAt: vh.NowUTC(),
	}
	if _, err := vh.DB.SaveVaultCID(newVault); err != nil {
		return "", fmt.Errorf("❌ failed to save new vault CID: %w", err)
	}
	vh.logger.Info("🗃️ VaultCID saved")

	// 8. Update session
	session.LastCID = newCID
	session.LastSynced = time.Now().Format(time.RFC3339)
	session.Dirty = false
	vh.vaultDirty = false

	vh.logger.Info("✅ Vault sync complete for user %d", userID)
	// utils.LogPretty("session after sync", session)

	return newCID, nil
}

func (vh *VaultHandler) EncryptVault(userID string, password string) (string, error) {
	vh.logger.Info("🔄 Starting vault sync for UserID: %d", userID)

	// 1. Get session
	session, err := vh.GetSession(userID)
	if err != nil {
		return "", fmt.Errorf("❌ no active session for user %d: %w", userID, err)
	}
	// ✅ Removed noisy LogPretty - too verbose for production
	// 2. Marshal in-memory vault
	vaultBytes, err := json.Marshal(session.Vault) // session.Vault.Sync()
	if err != nil {
		return "", fmt.Errorf("❌ failed to marshal vault: %w", err)
	}
	vh.logger.Info("🧱 Vault marshalled (%d bytes)", len(vaultBytes))

	// 3. Encrypt
	encrypted, err := blockchain.Encrypt(vaultBytes, password)
	if err != nil {
		return "", fmt.Errorf("❌ failed to encrypt vault: %w", err)
	}
	vh.logger.Info("🔐 Vault encrypted")

	return string(encrypted), nil
}

// -----------------------------
// Vault - Crud
// -----------------------------
// AddEntryFor: resilient behaviour: will always add the entry; Tracecore commit best-effort
func (vh *VaultHandler) AddEntryFor(userID string, entry any) (*any, error) {
	defer func() {
		if r := recover(); r != nil {
			vh.logger.Error("🔥 Panic in AddEntryFor: %v\nStack:\n%s", r, debug.Stack())
		}
	}()

	ve, ok := entry.(models.VaultEntry)
	if !ok {
		return nil, fmt.Errorf("❌ entry does not implement VaultEntry interface")
	}
	entryType := ve.GetTypeName()

	// 1. Prepare tracecore commit envelope if applicable
	// options := map[string]any{
	// 	"action":      services.CREATE_ENTRY,
	// 	"receiver":    "",
	// 	"permissions": nil,
	// 	"expiry":      "",
	// }
	// var envelope *tracecore.CommitEnvelope
	// if env, err := vh.PrepareTracecoreEnvelope(userID, ve, options); err == nil && env != nil {
	// 	envelope = env
	// 	// Try committing immediately (best-effort)
	// 	if vh.TracecoreClient != nil {
	// 		traceResp, err := vh.TracecoreClient.Commit(*envelope)
	// 		if err != nil || traceResp == nil || traceResp.Status != 201 {
	// 			// commit failed -> enqueue and mark dirty, but proceed to add entry
	// 			vh.logger.Error("❌ Tracecore commit failed (best-effort): %v", err)
	// 			vh.EnqueuePendingCommit(userID, *envelope)
	// 			// Optionally emit event/metric here
	// 		} else {
	// 			vh.logger.Info("✅ Tracecore commit succeeded for entry %s user %d", entryType, userID)
	// 		}
	// 	} else {
	// 		// No tracecore client available, just enqueue
	// 		vh.logger.Warn("⚠️ TracecoreClient unavailable — enqueueing commit for later")
	// 		vh.EnqueuePendingCommit(userID, *envelope)
	// 	}
	// } else if err != nil {
	// 	// if preparing the envelope failed, log but continue adding entry
	// 	vh.logger.Warn("⚠️ Skip Tracecore commit prep: %v", err)
	// }

	// 2. Add entry (route to handler)
	handler, err := vh.EntryRegistry.HandlerFor(entryType)
	if err != nil {
		return nil, err
	}

	created, err := handler.Add(userID, entry)
	var result any = created

	vh.logger.Info("✅ Created %s entry for user %s", entryType, userID)
	vh.MarkDirty(userID)

	return &result, err
}
func (vh *VaultHandler) AddEntry(userID string, entryType string, raw json.RawMessage) (any, error) {
	parsed, err := vh.EntryRegistry.UnmarshalEntry(entryType, raw)
	if err != nil {
		vh.logger.Error("❌ Failed to unmarshal %s entry for user %d: %v", entryType, userID, err)
		return nil, fmt.Errorf("failed to parse %s entry: %w", entryType, err)

	}
	// ✅ Removed noisy debug log

	// 1. GetBackendPlanParamForTransaction for managing plans from remote

	return vh.AddEntryFor(userID, parsed)
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

	updated, err := handler.Edit(userID, entry)
	if err != nil {
		return nil, err
	}

	vh.logger.Info("✅ Updated %s entry for user %s", entryType, userID)
	vh.MarkDirty(userID)

	return updated, nil
}
func (vh *VaultHandler) UpdateEntry(userID string, entryType string, raw json.RawMessage) (any, error) {
	utils.LogPretty("VaultHandler.UpdateEntry raw", raw)
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
	_, err = handler.Trash(userID, entryID)
	vh.logger.Info("✅ trashed %s entry for user %s", entryType, userID)
	vh.MarkDirty(userID)

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
	_, err = handler.Restore(userID, entryID)

	vh.logger.Info("✅ restored %s entry for user %s", entryType, userID)
	vh.MarkDirty(userID)

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

func (vh *VaultHandler) CreateFolder(userID string, name string) (*models.VaultPayload, error) {
	session, err := vh.GetSession(userID)
	if err != nil {
		return nil, fmt.Errorf("no active session for user %s", userID)
	}

	folder := &models.Folder{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		IsDraft:   false,
	}

	newFolder, err := vh.DB.CreateFolder(*folder)
	if err != nil {
		return nil, err
	}

	vault := session.Vault
	vault.Folders = append(vault.Folders, *newFolder)

	vh.logger.Info("✅ Added %s folder for user %s", newFolder.Name, userID)
	vh.MarkDirty(userID)

	return vault, nil
}
func (vh *VaultHandler) GetFoldersByVault(vaultCID string) ([]models.Folder, error) {
	return vh.DB.GetFoldersByVault(vaultCID)
}
func (vh *VaultHandler) UpdateFolder(id string, newName string, isDraft bool) (*models.Folder, error) {
	folder, err := vh.DB.GetFolderById(id)
	if err != nil {
		return nil, err
	}

	folder.Name = newName
	folder.IsDraft = isDraft
	folder.UpdatedAt = time.Now().Format(time.RFC3339)

	saved, err := vh.DB.SaveFolder(*folder)
	if err != nil {
		return nil, err
	}
	return saved, nil
}
func (vh *VaultHandler) DeleteFolder(userID string, id string) error {	
	session, ok := vh.GetSession(userID)
	if ok != nil {
		return fmt.Errorf("no active session for user %s", userID)
	}

	// 1. Find the folder
	folder, err := vh.DB.GetFolderById(id)
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
	if err := vh.DB.DeleteFolder(id); err != nil {
		return err
	}

	// 5. Remove from in-memory vault state
	newFolders := []models.Folder{}
	for _, f := range vault.Folders {
		if f.ID != id {
			newFolders = append(newFolders, f)
		}
	}
	vault.Folders = newFolders

	return nil
}

func (vh *VaultHandler) ListSharedEntries(ctx context.Context, userID string) ([]share_domain.ShareEntry, error) {
	user, err := vh.DB.FindUserById(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found with ID %d: %w", userID, err)
	}
	utils.LogPretty("ListSharedEntries - user", user)
	existingSession, ok := vh.Sessions[userID]
	if !ok {
		return nil, fmt.Errorf("no active session for user %s", userID)
	}

	repo := share_infrastructure.NewGormShareRepository(vh.DB.DB)
	uc := share_application_use_cases.NewShareUseCase(repo, vh.TracecoreClient, vh.EventDispatcher)

	entries, err := uc.ListSharedEntries(ctx, userID, existingSession.VaultRuntimeContext.SessionSecrets["cloud_jwt"])
	if err != nil {
		return nil, fmt.Errorf("failed fetching shared entries: %w", err)
	}

	return entries, nil
}
func (vh *VaultHandler) ListReceivedShares(ctx context.Context, userID string) ([]share_domain.ShareEntry, error) {

	user, err := vh.DB.FindUserById(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	utils.LogPretty("ListReceivedEntries - user", user)
	existingSession, ok := vh.Sessions[user.ID]
	if !ok {
		return nil, fmt.Errorf("no active session for user %s", user.ID)
	}

	repo := share_infrastructure.NewGormShareRepository(vh.DB.DB)
	uc := share_application_use_cases.NewShareUseCase(repo, vh.TracecoreClient, vh.EventDispatcher)

	entries, err := uc.ListReceivedShares(ctx, userID, existingSession.VaultRuntimeContext.SessionSecrets["cloud_jwt"])
	if err != nil {
		return nil, fmt.Errorf("failed fetching shared entries: %w", err)
	}

	return entries, nil
}

// vaultHandler.go

func (vh *VaultHandler) GetShareForAccept(
	ctx context.Context,
	userID string,
	shareID string,
) (*share_domain.ShareAcceptData, error) {

	_, err := vh.DB.FindUserById(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	repo := share_infrastructure.NewGormShareRepository(vh.DB.DB)
	uc := share_application_use_cases.NewShareUseCase(repo, vh.TracecoreClient, vh.EventDispatcher)

	return uc.GetShareForAccept(ctx, shareID, userID)
}
func (vh *VaultHandler) AcceptShare(ctx context.Context, userID string, shareID string) (*share_application_use_cases.AcceptShareResult, error) {
	repo := share_infrastructure.NewGormShareRepository(vh.DB.DB)
	usecase := share_application_use_cases.NewShareUseCase(repo, vh.TracecoreClient, vh.EventDispatcher)
	result, err := usecase.AcceptShare(ctx, shareID, userID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (vh *VaultHandler) RejectShare(ctx context.Context, userID string, shareID string) (*share_application_use_cases.RejectShareResult, error) {
	repo := share_infrastructure.NewGormShareRepository(vh.DB.DB)
	usecase := share_application_use_cases.NewShareUseCase(repo, vh.TracecoreClient, vh.EventDispatcher)
	result, err := usecase.RejectShare(ctx, shareID, userID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (vh *VaultHandler) AddReceiver(ctx context.Context, userID string, in share_application_use_cases.AddReceiverInput) (*share_application_use_cases.AddReceiverResult, error) {

	repo := share_infrastructure.NewGormShareRepository(vh.DB.DB)
	usecase := share_application_use_cases.NewShareUseCase(repo, vh.TracecoreClient, vh.EventDispatcher)
	result, err := usecase.AddReceiver(ctx, userID, in)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type RecipientPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type CreateShareEntryPayload struct {
	EntryName  string `json:"entry_name"`
	EntryType  string `json:"entry_type"`
	EntryRef   string `json:"entry_ref"`
	Status     string `json:"status"`
	AccessMode string `json:"access_mode"`
	Encryption string `json:"encryption"`
	// Snapshot as JSON string from frontend
	EntrySnapshot string `json:"entry_snapshot"`

	ExpiresAt  string             `json:"expires_at"`
	Recipients []RecipientPayload `json:"recipients"`
}

func (vh *VaultHandler) CreateShareEntry(ctx context.Context, payload CreateShareEntryPayload, ownerID string) (*share_domain.ShareEntry, error) {
	// Convert JSON string -> domain struct
	var snapshot share_domain.EntrySnapshot
	if err := json.Unmarshal([]byte(payload.EntrySnapshot), &snapshot); err != nil {
		return nil, fmt.Errorf("invalid entry_snapshot: %w", err)
	}
	utils.LogPretty("payload", payload)

	// map payload -> domain.ShareEntry
	var s share_domain.ShareEntry
	s.OwnerID = ownerID
	s.EntryName = payload.EntryName
	s.EntryRef = payload.EntryRef
	s.EntryType = payload.EntryType
	s.Status = payload.Status
	s.AccessMode = payload.AccessMode
	s.Encryption = payload.Encryption
	s.EntrySnapshot = snapshot

	// parse ExpiresAt if present (payload.ExpiresAt is string)
	if payload.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, payload.ExpiresAt)
		if err == nil {
			s.ExpiresAt = &t
		} else {
			// optional: handle other formats or ignore
			s.ExpiresAt = nil
		}
	}

	// recipients
	recips := make([]share_domain.Recipient, 0, len(payload.Recipients))
	for _, r := range payload.Recipients {
		recips = append(recips, share_domain.Recipient{
			Name:  r.Name,
			Email: r.Email,
			Role:  r.Role,
			// IDs are assigned by DB
		})
	}
	s.Recipients = recips

	// create dependencies once and inject
	repo := share_infrastructure.NewGormShareRepository(vh.DB.DB)
	tcClient := vh.TracecoreClient   // ensure VaultHandler has TracecoreClient set
	dispatcher := vh.EventDispatcher // ensure VaultHandler has dispatcher reference (or nil)

	uc := share_application_use_cases.NewShareUseCase(repo, tcClient, dispatcher)
	// call usecase
	created, err := uc.CreateShare(ctx, s)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("Handler response - created", created)
	return created, nil
}

// -----------------------------
// Tracecore - connexion
// -----------------------------
// PrepareTracecoreEnvelope builds and signs envelope; returns nil,nil if not enabled/needed
func (vh *VaultHandler) PrepareTracecoreEnvelope(userID string, entry models.VaultEntry, options map[string]interface{}) (*tracecore.CommitEnvelope, error) {
	// load session
	session, err := vh.GetSession(userID)
	if err != nil {
		return nil, fmt.Errorf("❌ no active session for user %s: %w", userID, err)
	}
	log.Printf("🔍 Options in PrepareTracecoreEnvelope ---------- : %+v", options)

	// load app config (server source of truth)
	appCfg, err := vh.DB.GetAppConfigByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to load app config for user %d: %w", userID, err)
	}
	appCfg.TracecoreEnabled = true
	if !appCfg.TracecoreEnabled {
		return nil, nil // not enabled -> nothing to do
	}

	// prepare actor signature (from session user stellar private key)
	userCfg := session.VaultRuntimeContext.CurrentUser
	msgToSign := fmt.Sprintf("user-%d:%s:%d", userID, entry.GetTypeName(), time.Now().UnixNano())

	actorSig, err := blockchain.SignActorWithStellarPrivateKey(userCfg.StellarAccount.PrivateKey, msgToSign)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to sign actor: %w", err)
	}

	log.Printf("🔍 Options in PrepareTracecoreEnvelope: %+v", options)
	// Create a factory for creating/updating... a post entry payload
	var permissions []string
	if val, ok := options["permissions"]; ok && val != nil {
		switch v := val.(type) {
		case []string:
			permissions = v
		case []interface{}:
			for _, p := range v {
				if s, ok := p.(string); ok {
					permissions = append(permissions, s)
				}
			}
		default:
			log.Printf("⚠️ Unexpected permissions type: %T", v)
		}
	} else {
		log.Println("⚠️ No permissions found — defaulting to []")
	}

	factoryCreate := services.NewCommitPayloadFactory(
		options["action"].(string),
		userCfg.StellarAccount.PublicKey, // "alice_public_key",
		entry,
		options["receiver"].(string),
		permissions,
		options["expiry"].(string),
	)

	payloadMetadata, err := factoryCreate.Build()
	if err != nil {
		vh.logger.Error("❌ Failed to create commit tracecore payload: %v", err)
		return nil, fmt.Errorf("❌ failed to create commit tracecore payload - %s: %w", services.CREATE_ENTRY, err)
	}
	payloadMetadata.Actor = tracecore.Actor{
		ID:        userCfg.ID,
		Role:      "end_user",
		Signature: actorSig,
	}
	payloadMetadata.Signature = actorSig
	// ✅ Removed noisy fmt.Printf - use logger instead if needed

	// Build payload
	tracePayload := tracecore.CommitPayload{
		RepoID:          appCfg.RepoID,
		Branch:          session.VaultRuntimeContext.AppSettings.Branch,
		Metadata:        *payloadMetadata,
		ValidationRules: []string{"REQUIRES_SIGNATURE", "VALID_ACTORS_ONLY"}, // extractRuleNames(appCfg.CommitRules), // adapt to your types
	}

	// sign the payload with dvault private key (server-side signing)
	signature, err := blockchain.SignWithDvaultPrivateKey(tracePayload)
	if err != nil {
		return nil, fmt.Errorf("❌ signing commit as dvault failed: %w", err)
	}

	envelope := tracecore.CommitEnvelope{Commit: tracePayload, Signature: signature}
	return &envelope, nil
}

// -----------------------------
// Tracecore - Pending Commits
// -----------------------------
// Add a commit to the user's pending list and persist session
func (vh *VaultHandler) QueuePendingCommits(userID string, commit tracecore.CommitEnvelope) error {
	vh.SessionsMu.Lock()
	defer vh.SessionsMu.Unlock()

	session, ok := vh.Sessions[userID]
	if !ok {
		return fmt.Errorf("no active session for user %s", userID)
	}

	session.PendingCommits = append(session.PendingCommits, commit)

	// Persist so retries survive restart
	if err := vh.DB.SaveSession(userID, session); err != nil {
		return fmt.Errorf("failed to save session with queued commit: %w", err)
	}

	return nil
}
func (vh *VaultHandler) EnqueuePendingCommit(userID string, env tracecore.CommitEnvelope) {
	vh.pendingMu.Lock()
	defer vh.pendingMu.Unlock()
	vh.pendingCommits[userID] = append(vh.pendingCommits[userID], env)

	if s, ok := vh.Sessions[userID]; ok {
		s.PendingCommits = append(s.PendingCommits, env)
	}
	vh.logger.Info("🔁 Enqueued pending commit for user %s (total pending: %d)", userID, len(vh.pendingCommits[userID]))
}
func (vh *VaultHandler) RetryPendingCommits(userID string) {
	vh.pendingMu.Lock()
	pending := vh.pendingCommits[userID]
	vh.pendingMu.Unlock()

	if len(pending) == 0 {
		return
	}

	// ✅ Removed duplicate log - already logged in StartPendingCommitWorker
	remaining := make([]tracecore.CommitEnvelope, 0, len(pending))

	for _, env := range pending {
		resp, err := vh.TracecoreClient.Commit(env)
		if err != nil || resp == nil || resp.Status != 201 {
			vh.logger.Error("❌ Retry commit failed for user %d: %v", userID, err)
			remaining = append(remaining, env)
			continue
		}
		vh.logger.Info("✅ Pending commit succeeded for user %d", userID)
	}

	vh.pendingMu.Lock()
	if len(remaining) == 0 {
		delete(vh.pendingCommits, userID)
	} else {
		vh.pendingCommits[userID] = remaining
	}
	vh.pendingMu.Unlock()

	// update session's pending commits so they persist on logout
	if s, ok := vh.Sessions[userID]; ok {
		s.PendingCommits = remaining
	}
}

// StartPendingCommitWorker starts a background loop that retries pending commits periodically.
// Call this once during app startup.
func (vh *VaultHandler) StartPendingCommitWorker(ctx context.Context, interval time.Duration) {
	go func() {
		vh.logger.Info("🚀 PendingCommitWorker started (interval: %s)", interval)
		backoff := make(map[string]time.Duration) // userID -> backoff duration

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				vh.logger.Info("🛑 PendingCommitWorker stopped")
				return
			case <-ticker.C:
				vh.pendingMu.Lock()
				userIDs := make([]string, 0, len(vh.pendingCommits))
				for uid := range vh.pendingCommits {
					userIDs = append(userIDs, uid)
				}
				vh.pendingMu.Unlock()

				for _, uid := range userIDs {
					// Apply per-user backoff
					if delay, ok := backoff[uid]; ok && delay > interval {
						backoff[uid] -= interval
						continue
					}

					// Try retrying (only log if there are actually pending commits)
					before := len(vh.pendingCommits[uid])
					if before > 0 {	
						vh.logger.Info("ℹ️ Retrying %d pending commits for user %s", before, uid)
					}
					vh.RetryPendingCommits(uid)
					after := len(vh.pendingCommits[uid])

					if after == 0 {
						delete(backoff, uid) // reset backoff on success
						continue
					}

					// Increase backoff if still failing
					if _, ok := backoff[uid]; !ok {
						backoff[uid] = interval // start with interval delay
					}
					// simple exponential cap at 30m
					backoff[uid] = minDuration(backoff[uid]*2, 30*time.Minute)
				}
			}
		}
	}()
}

// helper
func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

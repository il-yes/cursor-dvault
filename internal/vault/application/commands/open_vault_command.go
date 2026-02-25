package vault_commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	utils "vault-app/internal/utils"
	app_config "vault-app/internal/config"
	app_config_commands "vault-app/internal/config/application/commands"
	app_config_domain "vault-app/internal/config/domain"
	app_config_ui "vault-app/internal/config/ui"
	vault_events "vault-app/internal/vault/application/events"
	vault_session "vault-app/internal/vault/application/session"
	vault_domain "vault-app/internal/vault/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vaults_persistence "vault-app/internal/vault/infrastructure/persistence"

	"gorm.io/gorm"
)

// -------- COMMAND --------

type OpenVaultCommand struct {
	UserID   string
	Password string
	Session  *vault_session.Session
}

// -------- RESULT --------

type OpenVaultResult struct {
	Vault          *vault_domain.Vault
	Content        *vault_domain.VaultPayload
	RuntimeContext *vault_session.RuntimeContext
	Session        *vault_session.Session
	LastCID        string
	ReusedExisting bool
}

// -------- HANDLER --------

type OpenVaultCommandHandler struct {
	vaultRepo vault_domain.VaultRepository
	now       func() string
}

// -------- CONSTRUCTOR --------

func NewOpenVaultCommandHandler(
	db *gorm.DB,
) *OpenVaultCommandHandler {
	vaultRepo := vaults_persistence.NewGormVaultRepository(db)

	return &OpenVaultCommandHandler{
		vaultRepo: vaultRepo,
		now:       func() string { return time.Now().UTC().Format(time.RFC3339) },
	}
}

// -------- EXECUTION --------
func (h *OpenVaultCommandHandler) Handle(
	ctx context.Context,
	cmd OpenVaultCommand,
	ipfs vault_domain.VaultStorage,
	crypto vault_domain.CryptoService,
	eventBus vault_events.VaultEventBus,
	configFacade app_config_ui.AppConfigHandler,
) (*OpenVaultResult, error) {

	// ------------------------------------------------------------
	// 0. ENFORCE INVARIANTS (NON-NEGOTIABLE)
	// ------------------------------------------------------------
	if cmd.Session == nil {
		cmd.Session = vault_session.InitNewSession(cmd.UserID)
	}

	runtimeCtx, err := h.GetRuntimeContext(ctx, cmd.UserID, configFacade)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("OpenVaultCommandHandler - Handle - runtimeCtx", runtimeCtx)

	// ------------------------------------------------------------
	// 1. REUSE SESSION VAULT IF POSSIBLE (SINGLE PATH)
	// ------------------------------------------------------------
	if cmd.Session.Vault != nil && cmd.Session.LastCID != "" {
		utils.LogPretty("OpenVaultCommandHandler - Handle - REUSE SESSION VAULT IF POSSIBLE (SINGLE PATH)", cmd.Session)
		payload := vaults_domain.ParseVaultPayload(cmd.Session.Vault)
		utils.LogPretty("OpenVaultCommandHandler - Handle - parsed payload", payload)

		evt := vault_events.VaultOpened{
			UserID:       cmd.UserID,
			VaultPayload: &payload,
			LastCID:      cmd.Session.LastCID,
			LastSynced:   cmd.Session.LastSynced,
			LastUpdated:  cmd.Session.LastUpdated,
			Runtime:      runtimeCtx,
			OccurredAt:   time.Now().Unix(),
		}

		eventBus.PublishVaultOpened(ctx, evt)

		return &OpenVaultResult{
			Vault:          nil,
			Content:        &payload,
			RuntimeContext: runtimeCtx,
			Session:        cmd.Session,
			LastCID:        cmd.Session.LastCID,
			ReusedExisting: true,
		}, nil
	}

	// ------------------------------------------------------------
	// 2. LOAD OR CREATE VAULT METADATA
	// ------------------------------------------------------------
	utils.LogPretty("OpenVaultCommandHandler - Handle - LOAD OR CREATE VAULT METADATA", cmd.UserID)
	vault, err := h.vaultRepo.GetLatestByUserID(cmd.UserID)
	if err != nil {
		if errors.Is(err, vault_domain.ErrVaultNotFound) {
			utils.LogPretty("OpenVaultCommandHandler - Handle - LOAD OR CREATE VAULT METADATA - VAULT NOT FOUND", cmd.UserID)
			vault = vault_domain.NewVault(cmd.UserID, "")
			if err := h.vaultRepo.SaveVault(vault); err != nil {
				utils.LogPretty("OpenVaultCommandHandler - Handle - LOAD OR CREATE VAULT METADATA - VAULT NOT FOUND", err)
				return nil, err
			}
			vault = vault_domain.NewVault(cmd.UserID, "")
			if err := h.vaultRepo.SaveVault(vault); err != nil {
				utils.LogPretty("OpenVaultCommandHandler - Handle - LOAD OR CREATE VAULT METADATA - NEW VAULT NOT FOUND", err)
				return nil, err
			}
		} else {
			utils.LogPretty("OpenVaultCommandHandler - Handle - LOAD OR CREATE VAULT METADATA", err)
			return nil, err
		}
	}

	// ------------------------------------------------------------
	// 3. LOAD ENCRYPTED VAULT CONTENT FROM IPFS
	// ------------------------------------------------------------
	encrypted, err := ipfs.GetData(vault.CID)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("OpenVaultCommandHandler - Handle - encrypted", encrypted)

	// ------------------------------------------------------------
	// 4. DECRYPT VAULT
	// ------------------------------------------------------------
	decrypted, err := crypto.Decrypt(encrypted, cmd.Password)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("Vault New blob JSON: %s\n", string(decrypted))
	// TODO: failed to retrieve stored data from ipfs returning an empty payload
	// ------------------------------------------------------------
	// 5. PARSE VAULT PAYLOAD
	// ------------------------------------------------------------
	payload := vaults_domain.ParseVaultPayload(decrypted)
	payload.Normalize()

	// ------------------------------------------------------------
	// 6. UPDATE SESSION (IN-MEMORY)
	// ------------------------------------------------------------
	cmd.Session.Vault = payload.ToBytes()
	cmd.Session.LastCID = vault.CID

	// ------------------------------------------------------------
	// 7. EMIT EVENT
	// ------------------------------------------------------------
	evt := vault_events.VaultOpened{
		UserID:       cmd.UserID,
		VaultPayload: &payload,
		LastCID:      vault.CID,
		LastSynced:   vault.UpdatedAt,
		LastUpdated:  vault.UpdatedAt,
		Runtime:      runtimeCtx,
		OccurredAt:   time.Now().Unix(),
	}

	eventBus.PublishVaultOpened(ctx, evt)

	// ------------------------------------------------------------
	// 8. RETURN RESULT
	// ------------------------------------------------------------
	return &OpenVaultResult{
		Vault:          vault,
		Content:        &payload,
		RuntimeContext: runtimeCtx,
		Session:        cmd.Session,
		LastCID:        vault.CID,
		ReusedExisting: false,
	}, nil
}

func (h *OpenVaultCommandHandler) GetRuntimeContext(ctx context.Context, userID string, configFacade app_config_ui.AppConfigHandler) (*vault_session.RuntimeContext, error) {
	// -----------------------------
	// 1. Load App & User Config
	// -----------------------------
	appCfg, userCfg, err := h.LoadConfigurationsForUserID(ctx, userID, configFacade)
	if err != nil {
		return nil, err
	}

	// -----------------------------
	// 2. Create runtime context
	// -----------------------------
	runtimeCtx := vault_session.NewRuntimeContext()
	runtimeCtx.SetAppConfig(*appCfg)
	runtimeCtx.SetUserConfig(*userCfg)

	return runtimeCtx, nil
}
func (h *OpenVaultCommandHandler) LoadConfigurationsForUserID(ctx context.Context, userID string, configFacade app_config_ui.AppConfigHandler) (*app_config.AppConfig, *app_config.UserConfig, error) {
	// -----------------------------
	// 1. Load App & User Config
	// -----------------------------
	// We use domain objects here because configFacade returns them
	domainAppCfg, _ := configFacade.GetAppConfigByUserID(userID)
	domainUserCfg, _ := configFacade.GetUserConfigByUserID(userID)
	utils.LogPretty("OpenVaultCommandHandler - Handle - domainAppCfg", domainAppCfg)
	utils.LogPretty("OpenVaultCommandHandler - Handle - domainUserCfg", domainUserCfg)

	// If either config missing → minimal config onboarding
	if domainAppCfg == nil || domainUserCfg == nil {
		// h.logger.Warn("⚠️ Missing configs for user %s — creating minimal config...", userID)
		
		// Fix: Use properly constructed input and extract AppConfig from output
		appInput := &app_config_commands.CreateAppConfigCommandInput{
			AppConfig: &app_config_domain.AppConfig{
				UserID: userID,
			 	Branch:           "main",
				TracecoreEnabled: false,
				EncryptionPolicy: "AES-256-GCM",
				VaultSettings: app_config_domain.VaultConfig{
					MaxEntries:       1000,
					AutoSyncEnabled:  false,
					EncryptionScheme: "AES-256-GCM",
				},
				Blockchain: app_config_domain.BlockchainConfig{
					Stellar: app_config_domain.StellarConfig{
						Network:    "testnet",
						HorizonURL: "https://horizon-testnet.stellar.org",
						Fee:        100,
					},
				},
			},
		}
		if err := configFacade.UpdateAppConfig(appInput.AppConfig); err != nil {
			return nil, nil, err
		}
		domainAppCfg = appInput.AppConfig

		// Fix: Use properly constructed input and extract UserConfig from output
		userConfig := &app_config_domain.UserConfig{
				ID: userID,	
				Role:           "user",
				Signature:      "",
				SharingRules:   []app_config_domain.SharingRule{},
				StellarAccount: app_config_domain.StellarAccountConfig{},
			}
		if err := configFacade.UpdateUserConfig(userConfig); err != nil {
			return nil, nil, err
		}
		domainUserCfg = userConfig
		utils.LogPretty("OpenVaultCommandHandler - Handle - domainUserCfg", domainUserCfg)
	}

	// Convert domain objects to local app_config objects (package mismatch workaround)
	appCfg := &app_config.AppConfig{}
	appBytes, _ := json.Marshal(domainAppCfg)
	if err := json.Unmarshal(appBytes, appCfg); err != nil {
		return nil, nil, fmt.Errorf("failed to convert app config: %w", err)
	}

	userCfg := &app_config.UserConfig{}
	userBytes, _ := json.Marshal(domainUserCfg)
	if err := json.Unmarshal(userBytes, userCfg); err != nil {
		return nil, nil, fmt.Errorf("failed to convert user config: %w", err)
	}

	return appCfg, userCfg, nil
}

type AttachRuntimeRequest struct {
	UserID  string
	Runtime *vault_session.RuntimeContext
}
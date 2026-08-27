package vault_commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	app_config_domain "vault-app/internal/config/domain"
	subscription_domain "vault-app/internal/subscription/domain"
	utils "vault-app/internal/utils"
	vault_events "vault-app/internal/vault/application/events"
	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vault_domain "vault-app/internal/vault/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vaults_persistence "vault-app/internal/vault/infrastructure/persistence"
)

// AppConfigFacade is a local interface for config operations needed by vault commands.
type AppConfigFacade interface {
	GetAppConfigByUserID(ctx context.Context, userID string) (*app_config_domain.AppConfig, error)
	GetUserConfigByUserID(userID string) (*app_config_domain.UserConfig, error)
	UpdateAppConfig(appConfig *app_config_domain.AppConfig) error
	UpdateUserConfig(userConfig *app_config_domain.UserConfig) error
	GetConfig(userID string, vault vaults_domain.Vault, sub *subscription_domain.Subscription) (*app_config_domain.Config, error)
}
type StorageEngineReconstructor interface {
	Reconstruct(
		ctx context.Context,
		cmd vault_queries.GetIPFSDataQuerry,
	) (vaults_domain.VaultPayload, error)
}

// -------- COMMAND --------

type OpenVaultCommand struct {
	UserID           string
	Password         string
	Session          *vault_session.Session
	UserOnboardingID string
	Configs          app_config_domain.Config
	Subscription     subscription_domain.Subscription
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
	VaultRepo     vault_domain.VaultRepository
	SessionRepo   vault_session.SessionRepository
	Now           func() string
	QueryHandler  vault_queries.GetIPFSDataQuerryHandler
	Reconstructor StorageEngineReconstructor
}

// -------- CONSTRUCTOR --------

func NewOpenVaultCommandHandler(
	db *gorm.DB,
	queryHandler vault_queries.GetIPFSDataQuerryHandler,
	reconstructor StorageEngineReconstructor,
) *OpenVaultCommandHandler {
	vaultRepo := vaults_persistence.NewGormVaultRepository(db)
	sessionRepo := vaults_persistence.NewGormSessionRepository(db)

	return &OpenVaultCommandHandler{
		VaultRepo:     vaultRepo,
		SessionRepo:   sessionRepo,
		Now:           func() string { return time.Now().UTC().Format(time.RFC3339) },
		QueryHandler:  queryHandler,
		Reconstructor: reconstructor,
	}
}

// -------- EXECUTION --------

func (h *OpenVaultCommandHandler) Handle(
	ctx context.Context,
	cmd OpenVaultCommand,
	eventBus vault_events.VaultEventBus,
	configFacade AppConfigFacade,
) (*OpenVaultResult, error) {
	// ------------------------------------------------------------
	// 0. EARLY DEPENDENCY & ARGUMENT GUARDS
	// ------------------------------------------------------------
	if h.VaultRepo == nil {
		return nil, errors.New("vault repository is nil")
	}
	if h.Reconstructor == nil {
		return nil, errors.New("StorageEngine Reconstructor is nil")
	}
	if configFacade == nil {
		return nil, errors.New("config facade is nil")
	}
	if eventBus == nil {
		return nil, errors.New("event bus is nil")
	}

	// ------------------------------------------------------------
	// 1. SESSION INVARIANT
	// ------------------------------------------------------------
	if cmd.Session == nil {
		cmd.Session = vault_session.InitNewSession(cmd.UserID)
	}

	runtimeCtx, err := h.GetRuntimeContext(ctx, cmd.UserID, configFacade)
	if err != nil {
		utils.LogPretty("OpenVaultCommandHandler - GetRuntimeContext error", err)
		return nil, fmt.Errorf("failed to get runtime context for user %s: %w", cmd.UserID, err)
	}

	// ------------------------------------------------------------
	// 2. FETCH AUTHORITATIVE VAULT RECORD
	// ------------------------------------------------------------
	authoritativeVault, err := h.VaultRepo.GetLatestByUserID(cmd.UserID)
	if err != nil {
		utils.LogPretty("OpenVaultCommandHandler - GetLatestByUserID err", err)
		return nil, fmt.Errorf("failed to load latest vault for user %s: %w", cmd.UserID, err)
	}

	authoritativeCID := ""
	if authoritativeVault != nil {
		authoritativeCID = authoritativeVault.CID
	}

	// ------------------------------------------------------------
	// 3. FAST PATH: SESSION IS FIRST SOURCE OF TRUTH
	// ------------------------------------------------------------
	sessionVaultPresent := cmd.Session != nil && len(cmd.Session.Vault) > 0
	sessionCIDPresent := cmd.Session != nil && cmd.Session.LastCID != ""
	cidsMatch := sessionCIDPresent && cmd.Session.LastCID == authoritativeCID && authoritativeCID != ""

	canReuseSession := sessionVaultPresent && sessionCIDPresent && cidsMatch

	if canReuseSession {
		payload := vaults_domain.ParseVaultPayload(cmd.Session.Vault)
		if authoritativeVault != nil {
			runtimeCtx.VaultID = authoritativeVault.ID
			runtimeCtx.AppConfig.RepoID = authoritativeVault.ID
			runtimeCtx.AppConfig.Branch = cmd.UserOnboardingID
			runtimeCtx.VaultName = authoritativeVault.Name
			payload.Name = authoritativeVault.Name
		}
		if runtimeCtx.SessionSecrets == nil {
			runtimeCtx.SessionSecrets = make(map[string]string)
		}
		if cmd.Password != "" {
			runtimeCtx.SessionSecrets["vault_password"] = cmd.Password
		}
		cmd.Session.Runtime = runtimeCtx

		if h.SessionRepo != nil {
			if err := h.SessionRepo.SaveSession(cmd.UserID, cmd.Session); err != nil {
				utils.LogPretty("OpenVaultCommandHandler - SaveSession fast-path err", err)
			}
		}

		eventBus.PublishVaultOpened(ctx, vault_events.VaultOpened{
			UserID:           cmd.UserID,
			UserOnboardingID: cmd.UserOnboardingID,
			VaultName:        payload.Name,
			VaultPayload:     &payload,
			LastCID:          cmd.Session.LastCID,
			Runtime:          runtimeCtx,
			OccurredAt:       time.Now().Unix(),
		})

		return &OpenVaultResult{
			Vault:          authoritativeVault,
			Content:        &payload,
			RuntimeContext: runtimeCtx,
			Session:        cmd.Session,
			LastCID:        cmd.Session.LastCID,
			ReusedExisting: true,
		}, nil
	}

	// ------------------------------------------------------------
	// 4. FALLBACK PATH (SESSION MISS): IPFS RECONSTRUCTION
	// ------------------------------------------------------------
	if authoritativeVault == nil {
		return nil, fmt.Errorf("vault record not found for user %s", cmd.UserID)
	}
	if authoritativeVault.CID == "" {
		return nil, fmt.Errorf("vault CID is empty for user %s; vault must be created first", cmd.UserID)
	}

	cfgs, err := configFacade.GetConfig(cmd.UserID, *authoritativeVault, &cmd.Subscription)
	if err != nil {
		utils.LogPretty("OpenVaultCommandHandler - Handle - Error getting config", err)
		return nil, fmt.Errorf("failed to get config for vault opening: %w", err)
	}

	reconstructedPayload, err := h.Reconstructor.Reconstruct(
		ctx,
		vault_queries.GetIPFSDataQuerry{
			CID:              authoritativeVault.CID,
			Password:         cmd.Password,
			Configs:          *cfgs,
			UserID:           cmd.UserID,
			VaultName:        authoritativeVault.Name,
			UserOnboardingID: cmd.UserOnboardingID,
		},
	)
	if err != nil {
		return nil, err
	}

	// ------------------------------------------------------------
	// 5. UPDATE & PERSIST REFRESHED SESSION
	// ------------------------------------------------------------
	cmd.Session.Vault = reconstructedPayload.ToBytes()
	cmd.Session.LastCID = authoritativeVault.CID
	runtimeCtx.VaultID = authoritativeVault.ID
	runtimeCtx.AppConfig.RepoID = cmd.UserOnboardingID
	reconstructedPayload.Name = authoritativeVault.Name
	if runtimeCtx.SessionSecrets == nil {
		runtimeCtx.SessionSecrets = make(map[string]string)
	}
	if cmd.Password != "" {
		runtimeCtx.SessionSecrets["vault_password"] = cmd.Password
	}
	cmd.Session.Runtime = runtimeCtx

	if h.SessionRepo != nil {
		if err := h.SessionRepo.SaveSession(cmd.UserID, cmd.Session); err != nil {
			utils.LogPretty("OpenVaultCommandHandler - SaveSession reconstruction err", err)
		}
	}

	// ------------------------------------------------------------
	// 6. EVENT
	// ------------------------------------------------------------
	eventBus.PublishVaultOpened(ctx, vault_events.VaultOpened{
		UserID:           cmd.UserID,
		UserOnboardingID: cmd.UserOnboardingID,
		VaultName:        authoritativeVault.Name,
		VaultPayload:     &reconstructedPayload,
		LastCID:          authoritativeVault.CID,
		Runtime:          runtimeCtx,
		OccurredAt:       time.Now().Unix(),
	})

	// ------------------------------------------------------------
	// 7. RETURN
	// ------------------------------------------------------------
	return &OpenVaultResult{
		Vault:          authoritativeVault,
		Content:        &reconstructedPayload,
		RuntimeContext: runtimeCtx,
		Session:        cmd.Session,
		LastCID:        authoritativeVault.CID,
		ReusedExisting: false,
	}, nil
}

func (h *OpenVaultCommandHandler) GetRuntimeContext(ctx context.Context, userID string, configFacade AppConfigFacade) (*vault_session.RuntimeContext, error) {
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
func (h *OpenVaultCommandHandler) LoadConfigurationsForUserID(ctx context.Context, userID string, configFacade AppConfigFacade) (*app_config_domain.AppConfig, *app_config_domain.UserConfig, error) {
	if userID == "" {
		return nil, nil, errors.New("LoadConfigurationsForUserID - user id is required")
	}
	if configFacade == nil {
		return nil, nil, errors.New("LoadConfigurationsForUserID - config facade is required")
	}
	// -----------------------------
	// 1. Load App & User Config
	// -----------------------------
	domainAppCfg, _ := configFacade.GetAppConfigByUserID(ctx, userID)
	domainUserCfg, _ := configFacade.GetUserConfigByUserID(userID)

	// If either config missing → minimal config onboarding
	if domainAppCfg == nil || domainUserCfg == nil {
		configs, err := app_config_domain.InitConfig(userID)
		if err != nil {
			utils.LogPretty("OpenVaultCommandHandler - LoadConfigurationsForUserID - internal error", err)
			return nil, nil, err
		}
		return configs.App, configs.User, nil
	}

	return domainAppCfg, domainUserCfg, nil
}

type AttachRuntimeRequest struct {
	UserID  string
	Runtime *vault_session.RuntimeContext
}

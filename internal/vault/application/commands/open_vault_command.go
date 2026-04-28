package vault_commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	app_config_domain "vault-app/internal/config/domain"
	utils "vault-app/internal/utils"
	vault_events "vault-app/internal/vault/application/events"
	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vault_domain "vault-app/internal/vault/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vaults_persistence "vault-app/internal/vault/infrastructure/persistence"

	"gorm.io/gorm"
)

// AppConfigFacade is a local interface for config operations needed by vault commands.
type AppConfigFacade interface {
	GetAppConfigByUserID(ctx context.Context, userID string) (*app_config_domain.AppConfig, error)
	GetUserConfigByUserID(userID string) (*app_config_domain.UserConfig, error)
	UpdateAppConfig(appConfig *app_config_domain.AppConfig) error
	UpdateUserConfig(userConfig *app_config_domain.UserConfig) error
}
type VaultReconstructorInterface interface {
	BuildFromRoot(
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
	VaultRepo          vault_domain.VaultRepository
	Now                func() string
	QueryHandler       vault_queries.GetIPFSDataQuerryHandler
	VaultReconstructor VaultReconstructorInterface
}

// -------- CONSTRUCTOR --------

func NewOpenVaultCommandHandler(
	db *gorm.DB,
	queryHandler vault_queries.GetIPFSDataQuerryHandler,
	vr VaultReconstructorInterface,
) *OpenVaultCommandHandler {
	vaultRepo := vaults_persistence.NewGormVaultRepository(db)

	return &OpenVaultCommandHandler{
		VaultRepo:          vaultRepo,
		Now:                func() string { return time.Now().UTC().Format(time.RFC3339) },
		QueryHandler:       queryHandler,
		VaultReconstructor: vr,
	}
}

// -------- EXECUTION --------

func (h *OpenVaultCommandHandler) Handle(
	ctx context.Context,
	cmd OpenVaultCommand,
	eventBus vault_events.VaultEventBus,
	configFacade AppConfigFacade,
) (*OpenVaultResult, error) {
	if configFacade == nil {
		return nil, errors.New("config facade is nil")
	}

	// ------------------------------------------------------------
	// 0. SESSION INVARIANT
	// ------------------------------------------------------------
	if cmd.Session == nil {
		utils.LogPretty("OpenVaultCommandHandler - session is nil", cmd.Session)
		cmd.Session = vault_session.InitNewSession(cmd.UserID)
	}

	runtimeCtx, err := h.GetRuntimeContext(ctx, cmd.UserID, configFacade)
	if err != nil {
		return nil, err
	}

	// ------------------------------------------------------------
	// 1. REUSE SESSION (FAST PATH)
	// ------------------------------------------------------------
	if cmd.Session.Vault != nil && cmd.Session.LastCID != "" {
		utils.LogPretty("OpenVaultCommandHandler - something is nil", cmd)
		payload := vaults_domain.ParseVaultPayload(cmd.Session.Vault)

		vault, err := h.VaultRepo.GetLatestByUserID(cmd.UserID)
		if err != nil {
			utils.LogPretty("OpenVaultCommandHandler - session is nil", cmd.Session)
		}

		// ------------------------------------------------------------
		// 1.a UPDATE SESSION
		// ------------------------------------------------------------
		cmd.Session.LastCID = vault.CID
		runtimeCtx.VaultID = vault.ID
		runtimeCtx.AppConfig.RepoID = vault.ID
		runtimeCtx.AppConfig.Branch = cmd.UserOnboardingID
		runtimeCtx.VaultID = vault.ID
		payload.Name = vault.Name
		cmd.Session.Runtime = runtimeCtx

		eventBus.PublishVaultOpened(ctx, vault_events.VaultOpened{
			UserID:       cmd.UserID,
			VaultPayload: &payload,
			LastCID:      cmd.Session.LastCID,
			Runtime:      runtimeCtx,
			OccurredAt:   time.Now().Unix(),
		})

		return &OpenVaultResult{
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
	vault, err := h.VaultRepo.GetLatestByUserID(cmd.UserID)
	if err != nil {
		fmt.Println("OpenVaultCommandHandler - Handle - 1st error")
		if errors.Is(err, vault_domain.ErrVaultNotFound) {
			fmt.Println("OpenVaultCommandHandler - Handle - 2nd error")
			vault = vault_domain.NewVault(cmd.UserID, "")
			utils.LogPretty("OpenVaultCommandHandler - Handle - newvault", vault)
			if err := h.VaultRepo.SaveVault(vault); err != nil {
				fmt.Println("OpenVaultCommandHandler - Handle - 3rd error")
				return nil, err
			}
		} else {
			fmt.Println("OpenVaultCommandHandler - Handle - 4th error", err)
			return nil, err
		}
	}

	if h.VaultReconstructor == nil {
		return nil, errors.New("VaultReconstructor is nil")
	}
	// ------------------------------------------------------------
	// 3. RECONSTRUCT VAULT FROM IPFS (ONLY PLACE)
	// ------------------------------------------------------------
	legacyPayload, err := h.VaultReconstructor.BuildFromRoot(
		ctx,
		vault_queries.GetIPFSDataQuerry{
			CID:              vault.CID,
			Password:         cmd.Password,
			AppCfg:           runtimeCtx.AppConfig,
			UserID:           cmd.UserID,
			VaultName:        vault.Name,
			UserOnboardingID: cmd.UserOnboardingID,
		},
	)
	if err != nil {
		return nil, err
	}

	// ------------------------------------------------------------
	// 4. UPDATE SESSION
	// ------------------------------------------------------------
	cmd.Session.Vault = legacyPayload.ToBytes()
	cmd.Session.LastCID = vault.CID
	runtimeCtx.VaultID = vault.ID
	runtimeCtx.AppConfig.RepoID = cmd.UserOnboardingID
	legacyPayload.Name = vault.Name
	cmd.Session.Runtime = runtimeCtx

	if eventBus == nil {
		return nil, errors.New("VaultReconstructor is nil")
	}
	// ------------------------------------------------------------
	// 5. EVENT
	// ------------------------------------------------------------
	eventBus.PublishVaultOpened(ctx, vault_events.VaultOpened{
		UserID:       cmd.UserID,
		VaultPayload: &legacyPayload,
		LastCID:      vault.CID,
		Runtime:      runtimeCtx,
		OccurredAt:   time.Now().Unix(),
	})

	// ------------------------------------------------------------
	// 6. RETURN
	// ------------------------------------------------------------
	return &OpenVaultResult{
		Vault:          vault,
		Content:        &legacyPayload,
		RuntimeContext: runtimeCtx,
		Session:        cmd.Session,
		LastCID:        vault.CID,
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

func (h *OpenVaultCommandHandler) Handle2(
	ctx context.Context,
	cmd OpenVaultCommand,
	eventBus vault_events.VaultEventBus,
	configFacade AppConfigFacade,
) (*OpenVaultResult, error) {
	utils.LogPretty("OpenVaultCommandHandler - Handle - cmd", cmd)

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

	// ------------------------------------------------------------
	// 1. REUSE SESSION VAULT IF POSSIBLE (SINGLE PATH)
	// ------------------------------------------------------------
	if cmd.Session.Vault != nil && cmd.Session.LastCID != "" {
		utils.LogPretty("OpenVaultCommandHandler - Handle - REUSE SESSION VAULT IF POSSIBLE (SINGLE PATH)", cmd.Session)
		payload := vaults_domain.ParseVaultPayload(cmd.Session.Vault)
		// utils.LogPretty("OpenVaultCommandHandler - Handle - parsed payload", payload)

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
	// 2. CREATE VAULT METADATA IF NOT EXISTS
	// ------------------------------------------------------------
	// utils.LogPretty("OpenVaultCommandHandler - Handle - LOAD OR CREATE VAULT METADATA", cmd.UserID)
	vault, err := h.VaultRepo.GetLatestByUserID(cmd.UserID)
	if err != nil {
		utils.LogPretty("OpenVaultCommandHandler - Handle - LOAD OR CREATE VAULT METADATA - VAULT NOT FOUND", err)
		if errors.Is(err, vault_domain.ErrVaultNotFound) {
			// utils.LogPretty("OpenVaultCommandHandler - Handle - LOAD OR CREATE VAULT METADATA - VAULT NOT FOUND", cmd.UserID)
			vault = vault_domain.NewVault(cmd.UserID, "")
			if err := h.VaultRepo.SaveVault(vault); err != nil {
				utils.LogPretty("OpenVaultCommandHandler - Handle - LOAD OR CREATE VAULT METADATA - VAULT NOT FOUND", err)
				return nil, err
			}
			vault = vault_domain.NewVault(cmd.UserID, "")
			if err := h.VaultRepo.SaveVault(vault); err != nil {
				// utils.LogPretty("OpenVaultCommandHandler - Handle - LOAD OR CREATE VAULT METADATA - NEW VAULT NOT FOUND", err)
				return nil, err
			}
		} else {
			// utils.LogPretty("OpenVaultCommandHandler - Handle - LOAD OR CREATE VAULT METADATA", err)
			return nil, err
		}
	}
	// ------------------------------------------------------------
	// 5. EMIT EVENT
	// ------------------------------------------------------------
	legacyPayload, err := h.VaultReconstructor.BuildFromRoot(ctx, vault_queries.GetIPFSDataQuerry{
		CID:       vault.CID,
		Password:  cmd.Password,
		AppCfg:    cmd.Session.Runtime.AppConfig,
		UserID:    cmd.UserID,
		VaultName: vault.Name,
	})
	if err != nil {
		return nil, err
	}
	evt := vault_events.VaultOpened{
		UserID:       cmd.UserID,
		VaultPayload: &legacyPayload,
		LastCID:      vault.CID,
		LastSynced:   vault.UpdatedAt,
		LastUpdated:  vault.UpdatedAt,
		Runtime:      runtimeCtx,
		OccurredAt:   time.Now().Unix(),
	}

	eventBus.PublishVaultOpened(ctx, evt)

	// ------------------------------------------------------------
	// 6. RETURN RESULT
	// ------------------------------------------------------------
	return &OpenVaultResult{
		Vault: vault,
		// Content:        &payload,
		RuntimeContext: runtimeCtx,
		Session:        cmd.Session,
		LastCID:        vault.CID,
		ReusedExisting: false,
	}, nil
}

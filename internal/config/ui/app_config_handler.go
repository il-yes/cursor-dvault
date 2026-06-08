package app_config_ui

import (
	"context"
	"errors"

	"gorm.io/gorm"

	app_config "vault-app/internal/config"
	app_config_commands "vault-app/internal/config/application/commands"
	app_config_dto "vault-app/internal/config/application/dto"
	app_config_worker "vault-app/internal/config/application/worker"
	app_config_domain "vault-app/internal/config/domain"
	app_config_persistence "vault-app/internal/config/infrastructure/persistence"
	"vault-app/internal/logger/logger"
	onboarding_ui_wails "vault-app/internal/onboarding/ui/wails"
	subscription_domain "vault-app/internal/subscription/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_ui "vault-app/internal/vault/ui"
)

type AppConfigHandler struct {
	DB                              *gorm.DB
	CreateAppConfigCommand          *app_config_commands.CreateAppConfigCommand
	CreateUserConfigCommand         *app_config_commands.CreateUserConfigCommand
	CreateVaultConfigCommand        *app_config_commands.CreateVaultConfigCommand
	CreateDeviceConfigCommand       *app_config_commands.CreateDeviceConfigCommand
	CreateSubscriptionConfigCommand *app_config_commands.CreateSubscriptionConfigCommand
	CreateOnboardingCommandHandler  *app_config_commands.CreateOnboardingCommandHandler
	UpdateVaultConfigCommand        *app_config_commands.UpdateVaultConfigCommand
	AppConfigRepository             app_config_domain.AppConfigRepository
	UserConfigRepository            app_config_domain.UserConfigRepository
	VaultConfigRepository           app_config_domain.VaultConfigRepository
	DeviceConfigRepository          app_config_domain.DeviceConfigRepository
	SubscriptionConfigRepository    app_config_domain.SubscriptionConfigRepository
	Logger                          logger.Logger
	VaultHandler                    *vault_ui.VaultHandler
	OnboardingConfigRepository      app_config_domain.OnboardingConfigRepository
	ApplyOnboardingPacksWorker      *app_config_worker.ApplyOnboardingPacksWorker
	OnboardingHandler               onboarding_ui_wails.OnBoardingHandler
}

func NewAppConfigHandler(
	DB *gorm.DB,
	Logger logger.Logger,
) *AppConfigHandler {
	appConfigRepository := app_config_persistence.NewGormAppConfigRepository(DB)
	userConfigRepository := app_config_persistence.NewGormUserConfigRepository(DB)
	vaultConfigRepository := app_config_persistence.NewGormVaultConfigRepository(DB)
	deviceConfigRepository := app_config_persistence.NewGormDeviceConfigRepository(DB)
	subscriptionConfigRepository := app_config_persistence.NewGormSubscriptionConfigRepository(DB)
	onboardingConfigRepository := app_config_persistence.NewGormOnboardingConfigRepo(DB)

	createAppConfigCommand := app_config_commands.NewCreateAppConfigCommand(appConfigRepository)
	createUserConfigCommand := app_config_commands.NewCreateUserConfigCommand(userConfigRepository)
	createVaultConfigCommand := app_config_commands.NewCreateVaultConfigCommand(vaultConfigRepository)
	createDeviceConfigCommand := app_config_commands.NewCreateDeviceConfigCommand(deviceConfigRepository)
	createSubscriptionConfigCommand := app_config_commands.NewCreateSubscriptionConfigCommand(subscriptionConfigRepository)
	createOnboardingCommandHandler := app_config_commands.NewCreateOnboardingCommandHandler(onboardingConfigRepository)
	updateVaultConfigCommand := app_config_commands.NewUpdateVaultConfigCommand(vaultConfigRepository, createVaultConfigCommand, &Logger)

	return &AppConfigHandler{
		DB:                              DB,
		CreateAppConfigCommand:          createAppConfigCommand,
		CreateUserConfigCommand:         createUserConfigCommand,
		CreateVaultConfigCommand:        createVaultConfigCommand,
		CreateDeviceConfigCommand:       createDeviceConfigCommand,
		CreateSubscriptionConfigCommand: createSubscriptionConfigCommand,
		UpdateVaultConfigCommand:        updateVaultConfigCommand,
		AppConfigRepository:             appConfigRepository,
		UserConfigRepository:            userConfigRepository,
		VaultConfigRepository:           vaultConfigRepository,
		DeviceConfigRepository:          deviceConfigRepository,
		SubscriptionConfigRepository:    subscriptionConfigRepository,
		Logger:                          Logger,
		OnboardingConfigRepository:      onboardingConfigRepository,
		CreateOnboardingCommandHandler:  createOnboardingCommandHandler,
	}
}

// -------- INITIALIZERS --------
func (vh *AppConfigHandler) InitAppConfig(input *app_config_commands.CreateAppConfigCommandInput) (*app_config_commands.CreateAppConfigCommandOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if vh.CreateAppConfigCommand == nil {
		return nil, errors.New("create app config command is required")
	}
	return vh.CreateAppConfigCommand.Execute(input)
}

func (vh *AppConfigHandler) InitUserConfig(input *app_config_commands.CreateUserConfigCommandInput) (*app_config_commands.CreateUserConfigCommandOutput, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if vh.CreateUserConfigCommand == nil {
		return nil, errors.New("create user config command is required")
	}
	return vh.CreateUserConfigCommand.Execute(input)
}

func (vh *AppConfigHandler) SaveConfigs(input *app_config_dto.CreateConfigCommandInput) (*app_config_dto.CreateConfigCommandOutput, error) {
	// Save AppCfg
	appConfigOutput, err := vh.CreateAppConfigCommand.Execute(&app_config_commands.CreateAppConfigCommandInput{
		AppConfig: input.Configs.App,
	})
	if err != nil {
		return nil, err
	}
	// Save UserCfg
	userConfigOutput, err := vh.CreateUserConfigCommand.Execute(&app_config_commands.CreateUserConfigCommandInput{
		UserConfig: input.Configs.User,
	})
	if err != nil {
		return nil, err
	}
	// Save VaultCfg
	vaultConfigOutput := vh.CreateVaultConfigCommand.Execute(app_config_commands.CreateVaultConfigInput{
		VaultConfig: input.Configs.Vaults,
	})
	if err != nil {
		return nil, err
	}
	// Save DeviceCfg
	_ = vh.CreateDeviceConfigCommand.Execute(app_config_commands.CreateDeviceConfigInput{
		Device: input.Configs.Devices[0],
	})
	if err != nil {
		return nil, err
	}
	// Save SubscriptionCfg
	subscriptionConfigOutput := vh.CreateSubscriptionConfigCommand.Execute(app_config_commands.CreateSubscriptionConfigInput{
		SubscriptionConfig: input.Configs.Subscription,
	})
	if err != nil {
		return nil, err
	}

	if err := vh.CreateOnboardingCommandHandler.Handle(&input.Configs.Onboarding); err != nil {
		return nil, err
	}

	return &app_config_dto.CreateConfigCommandOutput{
		Configs: app_config_domain.Config{
			App:          appConfigOutput.AppConfig,
			User:         userConfigOutput.UserConfig,
			Vaults:       vaultConfigOutput.VaultConfig,
			Subscription: subscriptionConfigOutput.SubscriptionConfig,
			Devices:      input.Configs.Devices,
		},
	}, nil
}

// -------- GETTERS --------
func (vh *AppConfigHandler) GetConfig(userID string, vault vaults_domain.Vault, sub *subscription_domain.Subscription) (*app_config_domain.Config, error) {
	vh.Logger.Info("AppConfigHandler: GetConfig - userID: %s, vaultName: %s", userID, vault.Name)

	appConfig, err := vh.GetAppConfigByUserID(context.Background(), userID)
	if err != nil {
		vh.Logger.Error("AppConfigHandler: GetConfig - Failed to get app config: %v", err)
	}
	userConfig, err := vh.GetUserConfigByUserID(userID)
	if err != nil {
		vh.Logger.Error("AppConfigHandler: GetConfig - Failed to get user config: %v", err)
	}
	vaultConfig, err := vh.GetVaultConfigByUserID(userID, vault.Name)
	if err != nil {
		vh.Logger.Error("AppConfigHandler: GetConfig - Failed to get vault config: %v", err)
	}
	vh.Logger.LogPretty("AppConfigHandler: GetConfig - appConfig", appConfig)
	vh.Logger.LogPretty("AppConfigHandler: GetConfig - appConfig", userConfig)

	subscriptionConfig, err := vh.GetSubscriptionConfigByUserID(sub.UserID, vault.Name)
	if err != nil {
		vh.Logger.Error("AppConfigHandler: GetConfig - Failed to get subscription config: %v", err)
	}
	// utils.LogPretty("AppConfigHandler: GetConfig - subscriptionConfig", subscriptionConfig)
	deviceConfigs, err := vh.GetDeviceConfigsByUserID(userID, vault.Name)
	if err != nil {
		vh.Logger.Error("AppConfigHandler: GetConfig - Failed to get device configs: %v", err)
	}
	if vh.VaultHandler == nil {
		vh.Logger.Error("AppConfigHandler: GetConfig - vaultHandler is nil: %v", err)
	}

	session, err := vh.VaultHandler.GetSession(userID)
	if err != nil {
		vh.Logger.Error("AppConfigHandler: GetConfig - Failed to get device configs: %v", err)
	}

	vh.Logger.LogPretty("AppConfigHandler: GetConfig - session", session)

	// vh.Logger.LogPretty("AppConfigHandler: GetConfig - AppConfig.Branch", session.Runtime.AppConfig.Branch)

	// Get User Oboarding
	var opt map[string]interface{}
	userOnboarding, err := vh.OnboardingHandler.UserRepo.FindByEmail(userConfig.Email)
	if err != nil {
		vh.Logger.Error("AppConfigHandler: GetConfig - Failed to get userOnboarding: %v", err)
	}
	// Get User Oboarding-Config
	onboardingConfig, err := vh.GetOnboardingConfigByUserID(context.Background(), userOnboarding.ID, opt)
	if err != nil {
		vh.Logger.Error("AppConfigHandler: GetConfig - Failed to get onboarding configs: %v", err)
	}
	if err != nil {
		vh.Logger.Error("AppConfigHandler: GetConfig - no onboarding configs ", "create a new one....")
		defaultConfig, err := app_config_domain.InitConfigFromVault(session.Runtime.AppConfig.Branch, vault.Name)
		defaultConfig.Onboarding.Packs = []string{"pack.client_data"}
		defaultConfig.Onboarding.UseCases = []string{"Personal Backups"}
		if err != nil {
			vh.Logger.Error("AppConfigHandler: GetConfig - Failed to init onboarding configs: %v", err)
		}
		vh.Logger.LogPretty("AppConfigHandler: GetConfig - onboarding configs created:", defaultConfig.Onboarding)
		onboardingConfig = &defaultConfig.Onboarding

		if err := vh.CreateOnboardingCommandHandler.Handle(onboardingConfig); err != nil {
			vh.Logger.Error("AppConfigHandler: GetConfig - Failed to create onboarding configs: %v", err)
		}
	}
	vh.Logger.LogPretty("AppConfigHandler: GetConfig - onboardingConfig", onboardingConfig)
	// vh.Logger.Info("AppConfigHandler: GetConfig - appConfig: %v", appConfig)
	// vh.Logger.Info("AppConfigHandler: GetConfig - userConfig: %v", userConfig)
	// vh.Logger.Info("AppConfigHandler: GetConfig - vaultConfig: %v", vaultConfig)
	// vh.Logger.Info("AppConfigHandler: GetConfig - subscriptionConfig: %v", subscriptionConfig)
	check, err := vh.GetOnboardingConfigByUserID(context.Background(), session.Runtime.AppConfig.Branch, opt)
	if err != nil {
		vh.Logger.Error("AppConfigHandler: GetConfig - Failed to check onboarding configs: %v", err)
	}
	vh.Logger.Info("AppConfigHandler: GetConfig - onboardingConfig checked:", check)

	return &app_config_domain.Config{
		App:          appConfig,
		User:         userConfig,
		Vaults:       vaultConfig,
		Subscription: &subscriptionConfig,
		Devices:      deviceConfigs,
		Onboarding:   *onboardingConfig,
	}, nil
}

func (vh *AppConfigHandler) GetAppConfigByUserID(ctx context.Context, userID string) (*app_config_domain.AppConfig, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}
	if vh.AppConfigRepository == nil {
		return nil, errors.New("app config repository is required")
	}
	appConfig, err := vh.AppConfigRepository.GetAppConfig(userID)
	if err != nil {
		vh.Logger.Error("AppConfigHandler: GetAppConfigByUserID - Failed to get app config: %v", err)
		return nil, err
	}

	return appConfig, nil
}

func (vh *AppConfigHandler) GetUserConfigByUserID(userID string) (*app_config_domain.UserConfig, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}
	if vh.UserConfigRepository == nil {
		return nil, errors.New("user config repository is required")
	}
	userConfig, err := vh.UserConfigRepository.GetUserConfig(userID)
	if err != nil {
		return nil, err
	}

	return userConfig, nil
}

func (vh *AppConfigHandler) GetDeviceConfigsByUserID(userID string, vaultName string) ([]app_config_domain.DeviceConfig, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}
	if vh.DeviceConfigRepository == nil {
		return nil, errors.New("device config repository is required")
	}
	deviceConfigs, err := vh.DeviceConfigRepository.FindByUserIDAndVaultName(userID, vaultName)
	if err != nil {
		return nil, err
	}

	return deviceConfigs, nil
}

func (vh *AppConfigHandler) GetVaultConfigByUserID(userID string, vaultName string) (app_config_domain.VaultConfigBeta, error) {
	if userID == "" {
		return app_config_domain.VaultConfigBeta{}, errors.New("user id is required")
	}
	if vaultName == "" {
		return app_config_domain.VaultConfigBeta{}, errors.New("vault name is required")
	}
	if vh.VaultConfigRepository == nil {
		return app_config_domain.VaultConfigBeta{}, errors.New("vault config repository is required")
	}
	vaultConfig, err := vh.VaultConfigRepository.FindByUserIDAndVaultName(userID, vaultName)
	if err != nil {
		return app_config_domain.VaultConfigBeta{}, nil
	}
	// utils.LogPretty("AppConfigHandler - GetVaultConfigByUserID - vaultConfig", vaultConfig)
	return vaultConfig, nil
}

func (vh *AppConfigHandler) GetSubscriptionConfigByUserID(userID string, vaultName string) (app_config_domain.SubscriptionConfig, error) {
	if userID == "" {
		return app_config_domain.SubscriptionConfig{}, errors.New("user id is required")
	}
	if vaultName == "" {
		return app_config_domain.SubscriptionConfig{}, errors.New("vault name is required")
	}
	if vh.SubscriptionConfigRepository == nil {
		return app_config_domain.SubscriptionConfig{}, errors.New("subscription config repository is required")
	}
	subscriptionConfig, err := vh.SubscriptionConfigRepository.FindByUserIDAndVaultName(userID, vaultName)
	if err != nil {
		return app_config_domain.SubscriptionConfig{}, err
	}

	return subscriptionConfig, nil
}

func (vh *AppConfigHandler) GetOnboardingConfigByUserID(ctx context.Context, userID string, optDB map[string]interface{}) (*app_config_domain.OnboardingConfig, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}
	if vh.OnboardingConfigRepository == nil {
		return nil, errors.New("onboarding config repository is required")
	}
	onboardingConfig, err := vh.OnboardingConfigRepository.FindByUserID(userID, optDB)
	if err != nil {
		vh.Logger.LogPretty("AppConfigHandler - GetOnboardingConfigByUserID - failed to find onboardingConfig", err)
		return nil, err
	}
	return onboardingConfig, nil
}

// -------- UPDATERS --------

func (vh *AppConfigHandler) EditSettings(userID string, vaultName string, settings *app_config_dto.Settings) error {

	if settings.UI != nil {
		userCfg, err := vh.GetUserConfigByUserID(userID)
		if err != nil {
			vh.Logger.Info("AppConfigHandler: OnChangingSettings - user config not found => create new one")
		}

		userCfg.UI = app_config_domain.UIConfig{
			Theme:             settings.UI.Theme,
			AnimationsEnabled: settings.UI.AnimationsEnabled,
		}
		// vh.Logger.LogPretty("userCfg", userCfg)
		if err := vh.UpdateUserConfig(userCfg); err != nil {
			vh.Logger.Error("AppConfigHandler: OnChangingSettings - Failed to update user config: %v", err)
			return err
		}
	}

	// vh.UpdateVaultConfigCommand.CreateIfNotExists(&settings.Vaults)

	existingVaultCfg, err := vh.GetVaultConfigByUserID(userID, vaultName)
	if err != nil {
		vh.Logger.Error("AppConfigHandler: OnChangingSettings - Failed to get vault config: %v", err)
	}
	if existingVaultCfg.ID == "" {
		vh.Logger.Info("AppConfigHandler: OnChangingSettings - vault config not found => create new one")
		output := vh.CreateVaultConfigCommand.Execute(app_config_commands.CreateVaultConfigInput{
			VaultConfig: settings.Vaults,
		})
		if output.Error != nil {
			vh.Logger.Error("AppConfigHandler: OnChangingSettings - Failed to create vault config: %v", output.Error)
			return output.Error
		}
		// vh.Logger.LogPretty("AppConfigHandler: OnChangingSettings - vaultCfg created", output.VaultConfig)
	} else {
		vh.Logger.Info("AppConfigHandler: OnChangingSettings - vault config found => update")
		if err := vh.UpdateVaultConfig(existingVaultCfg.BaseVaultConfig.ID, &settings.Vaults); err != nil {
			vh.Logger.Error("AppConfigHandler: OnChangingSettings - Failed to update vault config: %v", err)
			return err
		}
	}
	vh.Logger.Info("settings.Vaults passed....")

	if settings.Device != nil {
		existingDeviceCfg, _ := vh.GetDeviceConfigsByUserID(userID, vaultName)
		if len(existingDeviceCfg) == 0 {
			vh.Logger.Info("AppConfigHandler: OnChangingSettings - device config not found => create new one")
			output := vh.CreateDeviceConfigCommand.Execute(app_config_commands.CreateDeviceConfigInput{
				Device: *settings.Device,
			})
			if output.Error != nil {
				vh.Logger.Error("AppConfigHandler: OnChangingSettings - Failed to create device config: %v", output.Error)
				return output.Error
			}
		} else {
			vh.Logger.Info("AppConfigHandler: OnChangingSettings - device config found => update")
			if err := vh.UpdateDeviceConfigs(existingDeviceCfg[0].ID, settings.Device); err != nil {
				vh.Logger.Error("AppConfigHandler: OnChangingSettings - Failed to update device config: %v", err)
				return err
			}
		}
	}
	vh.Logger.Info("settings.Device passed....")

	existingSub, err := vh.GetSubscriptionConfigByUserID(userID, vaultName)
	if err != nil {
		vh.Logger.Info("AppConfigHandler: OnChangingSettings - subscription config not found")
	}
	if existingSub.ID == "" {
		vh.Logger.Info("AppConfigHandler: OnChangingSettings - subscription config not found => create new one")
		output := vh.CreateSubscriptionConfigCommand.Execute(app_config_commands.CreateSubscriptionConfigInput{
			SubscriptionConfig: settings.Subscription,
		})
		if output.Error != nil {
			vh.Logger.Error("AppConfigHandler: OnChangingSettings - Failed to create subscription config: %v", output.Error)
			return output.Error
		}
		// vh.Logger.LogPretty("AppConfigHandler: OnChangingSettings - subscriptionCfg created", output.SubscriptionConfig)
	} else {
		vh.Logger.Info("AppConfigHandler: OnChangingSettings - subscription config found => update")
		if err := vh.UpdateSubscriptionConfig(existingSub.ID, settings.Subscription); err != nil {
			vh.Logger.Error("AppConfigHandler: OnChangingSettings - Failed to update subscription config: %v", err)
			return err
		}
	}
	vh.Logger.Info("settings.Subscription passed....")

	if settings.Storage != nil {
		// vh.Logger.LogPretty("settings storage", settings.Storage)
		appCfg, err := vh.GetAppConfigByUserID(context.Background(), userID)
		if err != nil {
			vh.Logger.Error("AppConfigHandler: OnChangingSettings - Failed to get app config: %v", err)
			return err
		}
		appCfg.Storage.Mode = settings.Storage.Mode
		if settings.Storage.Mode == app_config.StorageLocal {
			appCfg.Storage.LocalIPFS.APIEndpoint = app_config_domain.DefaultStorageConfig().LocalIPFS.APIEndpoint
			appCfg.Storage.LocalIPFS.GatewayURL = app_config_domain.DefaultStorageConfig().LocalIPFS.GatewayURL
		}
		if settings.Storage.Mode == app_config.StorageCloud {
			appCfg.Storage.Cloud.BaseURL = app_config_domain.DefaultStorageConfig().Cloud.BaseURL
		}
		if err := vh.UpdateAppConfig(appCfg); err != nil {
			vh.Logger.Error("AppConfigHandler: OnChangingSettings - Failed to update Storage config: %v", err)
			return err
		}
		// vh.Logger.LogPretty("appCfg storage", appCfg.Storage)
	}
	vh.Logger.Info("settings.Storage passed....")

	var opt map[string]interface{}
	onboardingCfg := app_config_domain.NewOnboardingConfig(
		settings.Onboarding.UserID,
		settings.Onboarding.Packs,
		settings.Onboarding.UseCases,
		settings.Onboarding.InstalledSeeds,
		settings.Onboarding.Completed,
		settings.Onboarding.PacksApplied,
	)

	if err := vh.UpdateOnboardingConfig(settings.Onboarding.UserID, &onboardingCfg, opt); err != nil {
		vh.Logger.Error("AppConfigHandler: OnChangingSettings - Failed to update Onboarding config: %v", err)
		return err
	}
	// vh.Logger.LogPretty("AppConfigHandler: OnChangingSettings - onboardingCfg updated", onboardingCfg)

	vh.Logger.Info("settings.Onboarding passed....")

	return nil
}

func (vh *AppConfigHandler) UpdateAppConfig(appConfig *app_config_domain.AppConfig) error {
	if appConfig == nil {
		return errors.New("app config is required")
	}
	if vh.AppConfigRepository == nil {
		return errors.New("app config repository is required")
	}
	return vh.AppConfigRepository.UpdateAppConfig(appConfig)
}

func (vh *AppConfigHandler) UpdateUserConfig(userConfig *app_config_domain.UserConfig) error {
	return vh.UserConfigRepository.UpdateUserConfig(userConfig)
}

func (vh *AppConfigHandler) UpdateVaultConfig(id string, vaultConfig *app_config_domain.VaultConfigBeta) error {
	return vh.VaultConfigRepository.Update(id, vaultConfig)
}

func (vh *AppConfigHandler) UpdateSubscriptionConfig(id string, subscriptionConfig *app_config_domain.SubscriptionConfig) error {
	return vh.SubscriptionConfigRepository.Update(id, subscriptionConfig)
}

func (vh *AppConfigHandler) UpdateDeviceConfigs(id string, deviceConfigs *app_config_domain.DeviceConfig) error {
	if err := vh.DeviceConfigRepository.Update(id, deviceConfigs); err != nil {
		return err
	}
	return nil
}

func (vh *AppConfigHandler) UpdateOnboardingConfig(id string, onboardingConfig *app_config_domain.OnboardingConfig, opts map[string]interface{}) error {
	if err := vh.OnboardingConfigRepository.Update(id, onboardingConfig, opts); err != nil {
		return err
	}
	return nil
}

// -------- EVENT HANDLERS --------
func (vh *AppConfigHandler) OnGenerateApiKey(userID string, stellarAccount *app_config_domain.StellarAccountConfig) (*app_config_domain.UserConfig, error) {
	userConfig, err := vh.GetUserConfigByUserID(userID)
	if err != nil {
		vh.Logger.Error("AppConfigHandler: OnGenerateApiKey - Failed to get user config: %v", err)
		return nil, err
	}
	userConfig.OnGenerateApiKey(stellarAccount)

	return userConfig, vh.UpdateUserConfig(userConfig)
}

func (vh *AppConfigHandler) OnApplyOnboardingPacks(userID string, vaultName string, userOnboardingID string) {
	if vh.ApplyOnboardingPacksWorker == nil {
		return
	}
	go vh.ApplyOnboardingPacksWorker.HandleOnboardingCompleted(context.Background(), userID, vaultName, userOnboardingID)
}

func (vh *AppConfigHandler) SetVaultHandler(handler vault_ui.VaultHandler) {
	vh.VaultHandler = &handler
}
func (vh *AppConfigHandler) SetOnboardingHandler(handler onboarding_ui_wails.OnBoardingHandler) {
	vh.OnboardingHandler = handler
}

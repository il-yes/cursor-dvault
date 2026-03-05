package app_config_ui

import (
	"context"
	"errors"
	app_config_commands "vault-app/internal/config/application/commands"
	app_config_domain "vault-app/internal/config/domain"
	app_config_persistence "vault-app/internal/config/infrastructure/persistence"
	"vault-app/internal/logger/logger"
	"vault-app/internal/utils"

	"gorm.io/gorm"
)

type AppConfigHandler struct {
	DB                      *gorm.DB
	CreateAppConfigCommand  *app_config_commands.CreateAppConfigCommand
	CreateUserConfigCommand *app_config_commands.CreateUserConfigCommand
	AppConfigRepository     app_config_domain.AppConfigRepository
	UserConfigRepository    app_config_domain.UserConfigRepository
	Logger                  logger.Logger
}

func NewAppConfigHandler(
	DB *gorm.DB,
	Logger logger.Logger,
) *AppConfigHandler {
	appConfigRepository := app_config_persistence.NewGormAppConfigRepository(DB)
	utils.LogPretty("AppConfigHandler - NewAppConfigHandler - AppConfigRepository", appConfigRepository)
	userConfigRepository := app_config_persistence.NewGormUserConfigRepository(DB)
	createAppConfigCommand := app_config_commands.NewCreateAppConfigCommand(appConfigRepository)
	createUserConfigCommand := app_config_commands.NewCreateUserConfigCommand(userConfigRepository)

	return &AppConfigHandler{
		DB:                      DB,
		CreateAppConfigCommand:  createAppConfigCommand,
		CreateUserConfigCommand: createUserConfigCommand,
		AppConfigRepository:     appConfigRepository,
		UserConfigRepository:    userConfigRepository,
		Logger:                  Logger,
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

// -------- GETTERS --------
func (vh *AppConfigHandler) GetAppConfigByUserID(ctx context.Context, userID string) (*app_config_domain.AppConfig, error) {
	utils.LogPretty("AppConfigHandler - GetAppConfigByUserID - userID", userID)
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
	utils.LogPretty("AppConfigHandler - GetAppConfigByUserID - appConfig", appConfig)
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
	utils.LogPretty("AppConfigHandler - GetUserConfigByUserID - userConfig", userConfig)
	return userConfig, nil
}

// -------- UPDATERS --------
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

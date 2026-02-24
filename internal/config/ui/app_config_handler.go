package app_config_ui

import (
	app_config_commands "vault-app/internal/config/application/commands"
	app_config_domain "vault-app/internal/config/domain"
	app_config_persistence "vault-app/internal/config/infrastructure/persistence"
	"vault-app/internal/logger/logger"

	"gorm.io/gorm"
)

type AppConfigHandler struct {
	DB *gorm.DB
	CreateAppConfigCommand *app_config_commands.CreateAppConfigCommand
	CreateUserConfigCommand *app_config_commands.CreateUserConfigCommand
	AppConfigRepository app_config_domain.AppConfigRepository
	UserConfigRepository app_config_domain.UserConfigRepository
	Logger logger.Logger
}

func NewAppConfigHandler(
	DB *gorm.DB,
	Logger logger.Logger,
) *AppConfigHandler {
	appConfigRepository := app_config_persistence.NewGormAppConfigRepository(DB)
	userConfigRepository := app_config_persistence.NewGormUserConfigRepository(DB)
	createAppConfigCommand := app_config_commands.NewCreateAppConfigCommand(appConfigRepository)
	createUserConfigCommand := app_config_commands.NewCreateUserConfigCommand(userConfigRepository)

	return &AppConfigHandler{
		DB: DB,
		CreateAppConfigCommand: createAppConfigCommand,
		CreateUserConfigCommand: createUserConfigCommand,
		AppConfigRepository: appConfigRepository,
		UserConfigRepository: userConfigRepository,
		Logger: Logger,
	}
}


// -------- INITIALIZERS --------
func (vh *AppConfigHandler) InitAppConfig(input *app_config_commands.CreateAppConfigCommandInput) (*app_config_commands.CreateAppConfigCommandOutput, error) {
	
	return vh.CreateAppConfigCommand.Execute(input)
}

func (vh *AppConfigHandler) InitUserConfig(input *app_config_commands.CreateUserConfigCommandInput) (*app_config_commands.CreateUserConfigCommandOutput, error) {
	return vh.CreateUserConfigCommand.Execute(input)
}

// -------- GETTERS --------
func (vh *AppConfigHandler) GetAppConfigByUserID(userID string) (*app_config_domain.AppConfig, error) {
	return vh.AppConfigRepository.GetAppConfig(userID)
}

func (vh *AppConfigHandler) GetUserConfigByUserID(userID string) (*app_config_domain.UserConfig, error) {
	return vh.UserConfigRepository.GetUserConfig(userID)
}

// -------- UPDATERS --------
func (vh *AppConfigHandler) UpdateAppConfig(appConfig *app_config_domain.AppConfig) error {
	return vh.AppConfigRepository.UpdateAppConfig(appConfig)
}

func (vh *AppConfigHandler) UpdateUserConfig(userConfig *app_config_domain.UserConfig) error {
	return vh.UserConfigRepository.UpdateUserConfig(userConfig)
}

// -------- EVENT HANDLERS --------
func (vh *AppConfigHandler) OnGenerateApiKey(userID string, stellarAccount *app_config_domain.StellarAccountConfig) (*app_config_domain.UserConfig, error) {
	userConfig, err := vh.UserConfigRepository.GetUserConfig(userID)
	if err != nil {
		vh.Logger.Error("AppConfigHandler: OnGenerateApiKey - Failed to get user config: %v", err)
		return nil, err
	}
	userConfig.OnGenerateApiKey(stellarAccount)
	
	return userConfig , vh.UserConfigRepository.UpdateUserConfig(userConfig)
}
package app_config_domain

type AppConfigRepository interface {
	CreateAppConfig(appConfig *AppConfig) error
	GetAppConfig(id string) (*AppConfig, error)
	GetAppConfigByUserID(userID string) (*AppConfig, error)
	UpdateAppConfig(appConfig *AppConfig) error
	DeleteAppConfig(id string) error
}

type UserConfigRepository interface {
	CreateUserConfig(userConfig *UserConfig) error
	GetUserConfig(id string) (*UserConfig, error)
	UpdateUserConfig(userConfig *UserConfig) error
	DeleteUserConfig(id string) error
}

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

type VaultConfigRepository interface {
	Create(vaultConfig *VaultConfigBeta) (*VaultConfigBeta, error)
	Find(id string) (*VaultConfigBeta, error)
	FindByUserIDAndVaultName(userID string, vaultName string) (VaultConfigBeta, error)
	FindAll(userID string, vaultName string) ([]VaultConfigBeta, error)
	Update(id string, vaultConfig *VaultConfigBeta) error
	Delete(id string) error
}

type DeviceConfigRepository interface {
	Create(deviceConfig *DeviceConfig) error
	Find(id string) (*DeviceConfig, error)
	FindByUserIDAndVaultName(userID string, vaultName string) ([]DeviceConfig, error)
	Update(id string, deviceConfig *DeviceConfig) error
	Delete(id string) error
}

type SubscriptionConfigRepository interface {
	Create(subscriptionConfig *SubscriptionConfig) error
	Find(id string) (*SubscriptionConfig, error)
	FindByUserIDAndVaultName(userID string, vaultName string) (SubscriptionConfig, error)
	FindAll(userID string) ([]SubscriptionConfig, error)
	Update(id string, subscriptionConfig *SubscriptionConfig) error
	Delete(id string) error
}

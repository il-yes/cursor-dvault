package app_config_dto

import (
	app_config "vault-app/internal/config"
	app_config_domain "vault-app/internal/config/domain"
)

type Settings struct {
    UI         *app_config_domain.UIConfig `json:"ui"`
    Subscription *app_config_domain.SubscriptionConfig `json:"subscription"`
    Vaults       app_config_domain.VaultConfigBeta `json:"vaults"`
    Device      *app_config_domain.DeviceConfig `json:"devices"`
	Storage     *app_config.StorageConfig `json:"storage"`
	Onboarding *app_config_domain.OnboardingConfig `json:"onboarding"`
}

type CreateConfigCommandInput struct {
	Configs app_config_domain.Config `json:"configs"`
}

type CreateConfigCommandOutput struct {
	Configs app_config_domain.Config `json:"configs"`
}
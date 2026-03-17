package app_config_dto

import app_config_domain "vault-app/internal/config/domain"

type Settings struct {
    UI         *app_config_domain.UIConfig `json:"ui"`
    Subscription *app_config_domain.SubscriptionConfig `json:"subscription"`
    Vaults       app_config_domain.VaultConfigBeta `json:"vaults"`
    Device      *app_config_domain.DeviceConfig `json:"devices"`
}
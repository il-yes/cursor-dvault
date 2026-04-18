package blockchain_ipfs

import (
	"vault-app/internal/blockchain"
	app_config "vault-app/internal/config"
	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/tracecore"
)

type StorageFactory interface {
	New(ctx *app_config_domain.VaultContext) app_config.StorageProvider
}

type DefaultStorageFactory struct {}

func (f *DefaultStorageFactory) New(vaultCtx *app_config_domain.VaultContext) app_config.StorageProvider {
	return blockchain.NewStorageProvider(
		blockchain.Config{
			StorageConfig: vaultCtx.StorageConfig,
			UserID:        vaultCtx.UserID,
			VaultName:     vaultCtx.VaultName,
		},
		tracecore.NewTracecoreFromConfig(&vaultCtx.AppConfig, "token"),
	)
}
package share_entry_ports

import (
	app_config_domain "vault-app/internal/config/domain"
	vault_dto "vault-app/internal/vault/application/dto"
	vaults_domain "vault-app/internal/vault/domain"
	vault_ui "vault-app/internal/vault/ui"
)

type VaultHandlerInterface interface {
	LoadAttachment(userID string, vaultName string, hash string, formatReturned string) (*vault_dto.LoadAttachmentResponse, error)
	UploadAttachementToIPFSWithEncryption(userID string, ur vault_ui.UploadAttachRequest) (string, error)
	GetLatestByUserID(userID string) (*vaults_domain.Vault, error)
	GetVaultSession(userID string) (*vaults_domain.VaultPayload, error)
}

type AppConfigHandlerInterface interface {
	GetUserConfigByUserID(userID string) (*app_config_domain.UserConfig, error) 
}
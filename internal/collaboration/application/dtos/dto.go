package collaboration_dtos

import (
	c3_asset_domain "vault-app/internal/c3_asset/domain"
	app_config_domain "vault-app/internal/config/domain"
	vault_dto "vault-app/internal/vault/application/dto"
	vaults_domain "vault-app/internal/vault/domain"
)

type ShareAssetWithTrustGroupRequest struct {
	AssetCID string
	TrustGroupID string
	WrappedDEK string
	KEKVersion uint64
	CreatedBy string
	Metadata map[string]string
}

type CreateCollaborativeShareRequest struct {
	TrustGroupID string            `json:"trust_group_id"`
	KEKVersion   uint64            `json:"kek_version"`
	CreatedBy    string            `json:"created_by"`
	AssetCID     string            `json:"asset_cid"`
	WrappedDEK   string            `json:"wrapped_dek"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	UserID string  `json:"user_id"`
	Password     string            `json:"password,omitempty"`
	StellarSecret string          `json:"stellar_secret,omitempty"`
	Configs      app_config_domain.Config 
	Vault        vaults_domain.Vault
}

type CreateCollaborativeShareResponse struct {
	ShareEntry c3_asset_domain.ShareEntry `json:"share_entry"`
}

type ResolveCollaborativeShareRequest struct {
	ShareEntryID     string `json:"share_entry_id"`
	CallerIdentityID string `json:"caller_identity_id,omitempty"`
	CallerVaultID    string `json:"caller_vault_id"`
	CallerUserID     string `json:"caller_user_id,omitempty"`
	ThreadID         string `json:"thread_id,omitempty"`
	EventID          string `json:"event_id,omitempty"`
	StellarAccount   app_config_domain.StellarAccountConfig `json:"stellar_account"`
	GetIPFSFile      vault_dto.GetFileFromIPFSRequest `json:"get_ipfs_file"`
	WrappedKEK   string `json:"wrapped_kek"`
	DeviceSeed     string `json:"device_seed"`
}

type ResolveCollaborativeShareResponse struct {
	ShareEntryID string            `json:"share_entry_id"`
	TrustGroupID string            `json:"trust_group_id"`
	CreatedBy    string            `json:"created_by"`
	CreatedAt    string            `json:"created_at"`
	Metadata     map[string]string `json:"metadata"`
	Plaintext    []byte            `json:"plaintext"`
}




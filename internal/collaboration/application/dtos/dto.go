package collaboration_dtos

import (
	c3_asset_domain "vault-app/internal/c3_asset/domain"
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
}

type CreateCollaborativeShareResponse struct {
	ShareEntry c3_asset_domain.ShareEntry `json:"share_entry"`
}

type ResolveCollaborativeShareRequest struct {
	ShareEntryID     string `json:"share_entry_id"`
	CallerIdentityID string `json:"caller_identity_id,omitempty"`
	CallerVaultID    string `json:"caller_vault_id"`
	CallerUserID     string `json:"caller_user_id,omitempty"`
	DeviceID         string `json:"device_id"`
}

type ResolveCollaborativeShareResponse struct {
	ShareEntryID string            `json:"share_entry_id"`
	TrustGroupID string            `json:"trust_group_id"`
	CreatedBy    string            `json:"created_by"`
	CreatedAt    string            `json:"created_at"`
	Metadata     map[string]string `json:"metadata"`
	Plaintext    []byte            `json:"plaintext"`
}




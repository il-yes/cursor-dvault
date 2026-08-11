package collaboration_dtos

import (
	c3_asset_domain "vault-app/internal/c3_asset/domain"
	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
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
	TrustGroupID string                                             `json:"trust_group_id"`
	KEKVersion   uint64                                             `json:"kek_version"`
	CreatedBy    string                                             `json:"created_by"`
	AssetCID     string                                             `json:"asset_cid"`
	WrappedDEK   string                                             `json:"wrapped_dek"`
	Envelopes    []trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest `json:"envelopes"`
	Metadata     map[string]string                                  `json:"metadata,omitempty"`
}

type CreateCollaborativeShareResponse struct {
	ShareEntry c3_asset_domain.ShareEntry                           `json:"share_entry"`
	Envelopes  []trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest `json:"envelopes"`
}




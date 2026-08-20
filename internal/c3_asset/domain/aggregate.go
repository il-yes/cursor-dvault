package c3_asset_domain

import (
	"time"

	"github.com/google/uuid"
)




type ShareEntryStatus string

const (
	ShareEntryStatusActive  ShareEntryStatus = "active"
	ShareEntryStatusRevoked ShareEntryStatus = "revoked"
)

type ShareEntry struct {
	ID           string            `json:"id"`
	AssetCID     string            `json:"asset_cid"`
	TrustGroupID string            `json:"trust_group_id"`
	WrappedDEK   string            `json:"wrapped_dek"`
	KEKVersion   uint64            `json:"kek_version"`
	CreatedBy    string            `json:"created_by"`
	CreatedAt    time.Time         `json:"created_at"`
	Status       ShareEntryStatus  `json:"status"`
	Metadata     map[string]string `json:"metadata"`
	IsDraft      bool              `json:"is_draft"`
	IsDirty      bool              `json:"is_dirty" gorm:"boolean"`
}


type Asset struct {
	CID         string `json:"cid" `
	ContentHash string `json:"content_hash" `
	Size        int64 `json:"size" `
	Type string `json:"type" `
	IsDirty     bool `json:"is_dirty" gorm:"boolean"`
}


func NewShareEntry(
	assetCID string,
	 trustGroupID string,
	 wrappedDEK string,
    kekVersion uint64,
	 createdBy string,
	 metadata map[string]string,
) (ShareEntry, error){

	if assetCID == "" {
		return ShareEntry{}, ErrAssetCIDRequired
	}

	if trustGroupID == "" {
		return ShareEntry{}, ErrTrustGroupIDRequired
	}

	if wrappedDEK == "" {
		return ShareEntry{}, ErrWrappedDEKRequired
	}

	if kekVersion == 0 {
		return ShareEntry{}, ErrKEKVersionRequired
	}

	if createdBy == "" {
		return ShareEntry{}, ErrCreatedByRequired
	}
	return ShareEntry{
		ID:        uuid.NewString(),
		AssetCID:     assetCID,
		TrustGroupID: trustGroupID,
		WrappedDEK:   wrappedDEK,
		KEKVersion: kekVersion,
		CreatedBy:    createdBy,
		Status:       ShareEntryStatusActive,
		Metadata:     metadata,
		IsDraft:      true,
		IsDirty:      false,
	},nil
}


func NewAsset(
	cID string,
	contentHash string,
	size int64,
	assetType string,
) (Asset, error) {

	if cID == "" {
		return Asset{}, ErrCIDRequired
	}

	if contentHash == "" {
		return Asset{}, ErrContentHashRequired
	}

	if size == 0 {
		return Asset{}, ErrSizeRequired
	}

	if assetType == "" {
		return Asset{}, ErrAssetTypeRequired
	}

	return Asset{
		CID:         cID,
		ContentHash: contentHash,
		Size:        size,
		Type:        assetType,
		IsDirty:     false,
	}, nil
}
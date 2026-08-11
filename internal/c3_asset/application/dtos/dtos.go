package c3_asset_dtos

type CreateShareEntryRequest struct {
	AssetCID     string
	TrustGroupID string
	WrappedDEK   string
	KEKVersion   uint64
	CreatedBy    string
	Metadata     map[string]string
}

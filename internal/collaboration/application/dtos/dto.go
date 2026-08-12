package collaboration_dtos


type ShareAssetWithTrustGroupRequest struct {
	AssetCID string
	TrustGroupID string
	WrappedDEK string
	KEKVersion uint64
	CreatedBy string
	Metadata map[string]string
}



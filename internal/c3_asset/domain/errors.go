package c3_asset_domain


import "errors"


var (
	ErrCIDRequired         = errors.New("cid is required")
	ErrContentHashRequired = errors.New("content hash is required")
	ErrSizeRequired        = errors.New("size is required")
	ErrAssetTypeRequired   = errors.New("asset type is required")
	ErrAssetCIDRequired    = errors.New("asset cid is required")
	ErrTrustGroupIDRequired  = errors.New("trust group id is required")
	ErrWrappedDEKRequired    = errors.New("wrapped dek is required")
	ErrCreatedByRequired   = errors.New("created by is required")
	ErrKEKVersionRequired  = errors.New("kek version is required")
	ErrRepositoryNil       = errors.New("repository is nil")
)
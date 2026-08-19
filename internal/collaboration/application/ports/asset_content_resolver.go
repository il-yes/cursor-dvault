package collaboration_ports

import "context"

type AssetContentResolver interface {
	FetchEncryptedAsset(ctx context.Context, assetCID string) ([]byte, error)
}

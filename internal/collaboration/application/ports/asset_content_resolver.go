package collaboration_ports

import (
	"context"
	tracecore_types "vault-app/internal/tracecore/types"
)

type AssetContentResolver interface {
	FetchEncryptedAsset(ctx context.Context, assetCID string) ([]byte, error)
}

type ThreadDataAccessGate interface {
	AccessThreadData(ctx context.Context, req tracecore_types.ThreadDataAccessRequest) (*tracecore_types.AccessCryptoShareResponse, error)
}

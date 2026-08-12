package c3_asset_domain

import (
	"context"

	tracecore_types "vault-app/internal/tracecore/types"
)


type CreateShareEntryRequest struct {
	ShareEntry ShareEntry
}
type GetShareEntryRequest struct {
	ShareEntryID string
}
type UpdateShareEntryRequest struct {
	ShareEntry ShareEntry
}
type DeleteShareEntryRequest struct {
	ShareEntryID string
}

type ShareEntryRepository interface {
	CreateShareEntry(
		ctx context.Context,
		req *CreateShareEntryRequest,
	) (*tracecore_types.CloudResponse[ShareEntry], error)

	GetShareEntry(
		ctx context.Context,
		req *GetShareEntryRequest,
	) (*tracecore_types.CloudResponse[ShareEntry], error)

	UpdateShareEntry(
		ctx context.Context,
		req *UpdateShareEntryRequest,
	) (*tracecore_types.CloudResponse[ShareEntry], error)

	DeleteShareEntry(
		ctx context.Context,
		req *DeleteShareEntryRequest,
	) (*tracecore_types.CloudResponse[ShareEntry], error)
}
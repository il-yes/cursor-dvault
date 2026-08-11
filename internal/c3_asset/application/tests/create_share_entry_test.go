package c3_asset_tests

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	c3_asset_dtos "vault-app/internal/c3_asset/application/dtos"
	c3_asset_usecases "vault-app/internal/c3_asset/application/usecases"
	c3_asset_domain "vault-app/internal/c3_asset/domain"
	tracecore_types "vault-app/internal/tracecore/types"
)

// --------------------------------------------------------------------------------------------------
// TEST HELPERS
// --------------------------------------------------------------------------------------------------

type shareEntryRepositoryMock struct {
	createFn func(ctx context.Context, req *c3_asset_domain.CreateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error)
	getFn func(ctx context.Context, req *c3_asset_domain.GetShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error)
	updateFn func(ctx context.Context, req *c3_asset_domain.UpdateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error)
	deleteFn func(ctx context.Context, req *c3_asset_domain.DeleteShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error)
}

func (m *shareEntryRepositoryMock) CreateShareEntry(ctx context.Context, req *c3_asset_domain.CreateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}

	return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{
		Data: req.ShareEntry,
	}, nil
}

func (m *shareEntryRepositoryMock) GetShareEntry(ctx context.Context, req *c3_asset_domain.GetShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	if m.getFn != nil {
		return m.getFn(ctx, req)
	}

	return nil, nil
}

func (m *shareEntryRepositoryMock) UpdateShareEntry(ctx context.Context, req *c3_asset_domain.UpdateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, req)
	}

	return nil, nil
}

func (m *shareEntryRepositoryMock) DeleteShareEntry(ctx context.Context, req *c3_asset_domain.DeleteShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, req)
	}

	return nil, nil
}
func validCreateShareEntryRequest() c3_asset_dtos.CreateShareEntryRequest {
	return c3_asset_dtos.CreateShareEntryRequest{
		AssetCID:     "bafybeihasset001",
		TrustGroupID: "trust-group-001",
		WrappedDEK:   "wrapped-dek-value",
		KEKVersion:   1,
		CreatedBy:    "vault-owner-001",
		Metadata: map[string]string{
			"title":      "Engineering Package",
			"asset_type": "design-review",
		},
	}
}

// --------------------------------------------------------------------------------------------------
// TESTS
// --------------------------------------------------------------------------------------------------
func TestCreateShareEntryUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	var receivedRequest *c3_asset_domain.CreateShareEntryRequest

	repo := &shareEntryRepositoryMock{
		createFn: func(
			ctx context.Context,
			req *c3_asset_domain.CreateShareEntryRequest,
		) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
			receivedRequest = req

			return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{
				Data: req.ShareEntry,
			}, nil
		},
	}

	uc := c3_asset_usecases.NewCreateShareEntryUsecase(repo)

	req := validCreateShareEntryRequest()

	result, err := uc.Execute(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, result)

	require.NotEmpty(t, result.ID)
	require.Equal(t, req.AssetCID, result.AssetCID)
	require.Equal(t, req.TrustGroupID, result.TrustGroupID)
	require.Equal(t, req.WrappedDEK, result.WrappedDEK)
	require.Equal(t, req.KEKVersion, result.KEKVersion)
	require.Equal(t, req.CreatedBy, result.CreatedBy)
	require.Equal(t, req.Metadata, result.Metadata)

	require.True(t, result.IsDraft)
	require.False(t, result.IsDirty)

	require.NotNil(t, receivedRequest)
	require.Equal(t, req.AssetCID, receivedRequest.ShareEntry.AssetCID)
	require.Equal(t, req.TrustGroupID, receivedRequest.ShareEntry.TrustGroupID)
	require.Equal(t, req.WrappedDEK, receivedRequest.ShareEntry.WrappedDEK)
	require.Equal(t, req.KEKVersion, receivedRequest.ShareEntry.KEKVersion)
	require.Equal(t, req.CreatedBy, receivedRequest.ShareEntry.CreatedBy)
	require.Equal(t, req.Metadata, receivedRequest.ShareEntry.Metadata)
}

func TestCreateShareEntryUsecase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("failed to create share entry")

	repo := &shareEntryRepositoryMock{
		createFn: func(
			ctx context.Context,
			req *c3_asset_domain.CreateShareEntryRequest,
		) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
			return nil, expectedErr
		},
	}

	uc := c3_asset_usecases.NewCreateShareEntryUsecase(repo)

	result, err := uc.Execute(
		ctx,
		validCreateShareEntryRequest(),
	)

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)
}

func TestCreateShareEntryUsecase_Execute_NilRepository(t *testing.T) {
	ctx := context.Background()

	uc := c3_asset_usecases.NewCreateShareEntryUsecase(nil)

	result, err := uc.Execute(
		ctx,
		validCreateShareEntryRequest(),
	)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, "share entry repository is required")
}

func TestCreateShareEntryUsecase_Execute_InvalidRequest(t *testing.T) {
	tests := []struct {
		name          string
		modify        func(req *c3_asset_dtos.CreateShareEntryRequest)
		expectedError string
	}{
		{
			name: "missing asset cid",
			modify: func(req *c3_asset_dtos.CreateShareEntryRequest) {
				req.AssetCID = "   "
			},
			expectedError: "asset cid is required",
		},
		{
			name: "missing trust group id",
			modify: func(req *c3_asset_dtos.CreateShareEntryRequest) {
				req.TrustGroupID = "   "
			},
			expectedError: "trust group id is required",
		},
		{
			name: "missing wrapped dek",
			modify: func(req *c3_asset_dtos.CreateShareEntryRequest) {
				req.WrappedDEK = "   "
			},
			expectedError: "wrapped dek is required",
		},
		{
			name: "missing created by",
			modify: func(req *c3_asset_dtos.CreateShareEntryRequest) {
				req.CreatedBy = "   "
			},
			expectedError: "created by is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			repo := &shareEntryRepositoryMock{}
			uc := c3_asset_usecases.NewCreateShareEntryUsecase(repo)

			req := validCreateShareEntryRequest()
			tt.modify(&req)

			result, err := uc.Execute(ctx, req)

			require.Error(t, err)
			require.Nil(t, result)
			require.EqualError(t, err, tt.expectedError)
		})
	}
}

func TestCreateShareEntryUsecase_ValidateRequest(t *testing.T) {
	repo := &shareEntryRepositoryMock{}
	uc := c3_asset_usecases.NewCreateShareEntryUsecase(repo)

	tests := []struct {
		name          string
		request       c3_asset_dtos.CreateShareEntryRequest
		expectedError string
	}{
		{
			name: "missing asset cid",
			request: c3_asset_dtos.CreateShareEntryRequest{
				TrustGroupID: "trust-group-001",
				WrappedDEK:   "wrapped-dek",
				CreatedBy:    "owner-001",
			},
			expectedError: "asset cid is required",
		},
		{
			name: "missing trust group id",
			request: c3_asset_dtos.CreateShareEntryRequest{
				AssetCID:   "asset-cid",
				WrappedDEK: "wrapped-dek",
				CreatedBy:  "owner-001",
			},
			expectedError: "trust group id is required",
		},
		{
			name: "missing wrapped dek",
			request: c3_asset_dtos.CreateShareEntryRequest{
				AssetCID:     "asset-cid",
				TrustGroupID: "trust-group-001",
				CreatedBy:    "owner-001",
			},
			expectedError: "wrapped dek is required",
		},
		{
			name: "missing created by",
			request: c3_asset_dtos.CreateShareEntryRequest{
				AssetCID:     "asset-cid",
				TrustGroupID: "trust-group-001",
				WrappedDEK:   "wrapped-dek",
			},
			expectedError: "created by is required",
		},
		{
			name:    "valid request",
			request: validCreateShareEntryRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.ValidateRequest(tt.request)

			if tt.expectedError == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.expectedError)
		})
	}
}

func TestCreateShareEntryUsecase_ValidateDependencies(t *testing.T) {
	tests := []struct {
		name          string
		repo          c3_asset_domain.ShareEntryRepository
		expectedError string
	}{
		{
			name:          "nil repository",
			repo:          nil,
			expectedError: "share entry repository is required",
		},
		{
			name: "valid repository",
			repo: &shareEntryRepositoryMock{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := c3_asset_usecases.NewCreateShareEntryUsecase(tt.repo)

			err := uc.ValidateDependencies()

			if tt.expectedError == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.expectedError)
		})
	}
}

func TestNewShareEntry(t *testing.T) {
	tests := []struct {
		name          string
		assetCID      string
		trustGroupID  string
		wrappedDEK    string
		kekVersion    uint64
		createdBy     string
		expectedError error
	}{
		{
			name:          "missing asset cid",
			trustGroupID:  "trust-group-001",
			wrappedDEK:    "wrapped-dek",
			kekVersion:    1,
			createdBy:     "owner-001",
			expectedError: c3_asset_domain.ErrAssetCIDRequired,
		},
		{
			name:          "missing trust group id",
			assetCID:      "asset-cid",
			wrappedDEK:    "wrapped-dek",
			kekVersion:    1,
			createdBy:     "owner-001",
			expectedError: c3_asset_domain.ErrTrustGroupIDRequired,
		},
		{
			name:          "missing wrapped dek",
			assetCID:      "asset-cid",
			trustGroupID:  "trust-group-001",
			kekVersion:    1,
			createdBy:     "owner-001",
			expectedError: c3_asset_domain.ErrWrappedDEKRequired,
		},
		{
			name:          "missing created by",
			assetCID:      "asset-cid",
			trustGroupID:  "trust-group-001",
			wrappedDEK:    "wrapped-dek",
			kekVersion:    1,
			expectedError: c3_asset_domain.ErrCreatedByRequired,
		},
		{
			name:          "missing kek version",
			assetCID:      "asset-cid",
			trustGroupID:  "trust-group-001",
			wrappedDEK:    "wrapped-dek",
			kekVersion:    0,
			createdBy:     "owner-001",
			expectedError: c3_asset_domain.ErrKEKVersionRequired,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := c3_asset_domain.NewShareEntry(
				tt.assetCID,
				tt.trustGroupID,
				tt.wrappedDEK,
				tt.kekVersion,
				tt.createdBy,
				nil,
			)

			require.Error(t, err)
			require.Empty(t, entry)
			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}



package collaboration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_domain "vault-app/internal/trust_group/domain"
)

type mockTrustGroupRepo struct {
	getFn func(ctx context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error)
}

func (m *mockTrustGroupRepo) CreateTrustGroup(ctx context.Context, req *trustgroup_domain.CreateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (m *mockTrustGroupRepo) GetTrustGroup(ctx context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	if m.getFn != nil {
		return m.getFn(ctx, req)
	}
	return nil, nil
}
func (m *mockTrustGroupRepo) GetTrustGroupMember(ctx context.Context, req *trustgroup_domain.GetTrustGroupMemberRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroupMember], error) {
	return nil, nil
}
func (m *mockTrustGroupRepo) AddMemberToTrustGroup(ctx context.Context, req *trustgroup_domain.AddMemberToTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (m *mockTrustGroupRepo) RemoveMemberFromTrustGroup(ctx context.Context, req *trustgroup_domain.RemoveMemberFromTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (m *mockTrustGroupRepo) UpdateTrustGroup(ctx context.Context, req *trustgroup_domain.UpdateTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (m *mockTrustGroupRepo) DeleteTrustGroup(ctx context.Context, req *trustgroup_domain.DeleteTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (m *mockTrustGroupRepo) ListTrustGroups(ctx context.Context, req *trustgroup_domain.ListTrustGroupsRequest) (*tracecore_types.CloudResponse[[]trustgroup_domain.TrustGroup], error) {
	return nil, nil
}
func (m *mockTrustGroupRepo) RotateTrustGroupKEK(ctx context.Context, req *trustgroup_domain.RotateTrustGroupKEKRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
	return nil, nil
}

type mockShareEntryRepo struct {
	createFn func(ctx context.Context, req *c3_asset_domain.CreateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error)
}

func (m *mockShareEntryRepo) CreateShareEntry(ctx context.Context, req *c3_asset_domain.CreateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	if m.createFn != nil {
		return m.createFn(ctx, req)
	}
	return nil, nil
}
func (m *mockShareEntryRepo) GetShareEntry(ctx context.Context, req *c3_asset_domain.GetShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	return nil, nil
}
func (m *mockShareEntryRepo) UpdateShareEntry(ctx context.Context, req *c3_asset_domain.UpdateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	return nil, nil
}
func (m *mockShareEntryRepo) DeleteShareEntry(ctx context.Context, req *c3_asset_domain.DeleteShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
	return nil, nil
}

func TestShareAssetWithTrustGroupUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	tgRepo := &mockTrustGroupRepo{
		getFn: func(ctx context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
			return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{
				Data: trustgroup_domain.TrustGroup{
					ID:         "tg-001",
					KEKVersion: 1,
				},
			}, nil
		},
	}

	shareRepo := &mockShareEntryRepo{
		createFn: func(ctx context.Context, req *c3_asset_domain.CreateShareEntryRequest) (*tracecore_types.CloudResponse[c3_asset_domain.ShareEntry], error) {
			return &tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{
				Data: req.ShareEntry,
			}, nil
		},
	}

	uc := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(tgRepo, shareRepo)

	req := collaboration_dtos.ShareAssetWithTrustGroupRequest{
		AssetCID:     "bafy-asset-001",
		TrustGroupID: "tg-001",
		WrappedDEK:   "wrapped-dek-value",
		KEKVersion:   1,
		CreatedBy:    "user-001",
	}

	res, err := uc.Execute(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, "bafy-asset-001", res.AssetCID)
	require.Equal(t, "tg-001", res.TrustGroupID)
	require.Equal(t, uint64(1), res.KEKVersion)
}

func TestShareAssetWithTrustGroupUsecase_Execute_StaleKEKVersion(t *testing.T) {
	ctx := context.Background()

	tgRepo := &mockTrustGroupRepo{
		getFn: func(ctx context.Context, req *trustgroup_domain.GetTrustGroupRequest) (*tracecore_types.CloudResponse[trustgroup_domain.TrustGroup], error) {
			return &tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{
				Data: trustgroup_domain.TrustGroup{
					ID:         "tg-001",
					KEKVersion: 2, // Group is at KEKVersion 2
				},
			}, nil
		},
	}

	shareRepo := &mockShareEntryRepo{}
	uc := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(tgRepo, shareRepo)

	req := collaboration_dtos.ShareAssetWithTrustGroupRequest{
		AssetCID:     "bafy-asset-001",
		TrustGroupID: "tg-001",
		WrappedDEK:   "wrapped-dek-value",
		KEKVersion:   1, // Request attempts KEKVersion 1 (stale)
		CreatedBy:    "user-001",
	}

	res, err := uc.Execute(ctx, req)
	require.Error(t, err)
	require.ErrorIs(t, err, trustgroup_domain.ErrStaleKEKVersion)
	require.Nil(t, res)
}

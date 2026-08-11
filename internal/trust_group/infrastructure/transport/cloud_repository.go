package transport

import (
	"context"

	trustgroup_domain "vault-app/internal/trust_group/domain"
)

type CloudRepository struct {
	client trustgroup_domain.TrustGroupRepository
}

func NewCloudRepository(client trustgroup_domain.TrustGroupRepository) *CloudRepository {
	return &CloudRepository{
		client: client,
	}
}

func (r *CloudRepository) CreateTrustGroup(ctx context.Context, req trustgroup_domain.CreateTrustGroupRequest, signature string) (*trustgroup_domain.TrustGroup, error) {
	resp, err := r.client.CreateTrustGroup(ctx, &req)
	return &resp.Data, err
}

func (r *CloudRepository) UpdateTrustGroup(ctx context.Context, req trustgroup_domain.UpdateTrustGroupRequest) (*trustgroup_domain.TrustGroup, error) {

	resp, err := r.client.UpdateTrustGroup(ctx, &req)
	return &resp.Data, err
}

func (r *CloudRepository) DeleteTrustGroup( ctx context.Context, req trustgroup_domain.DeleteTrustGroupRequest) (*trustgroup_domain.TrustGroup, error) {
	resp, err := r.client.DeleteTrustGroup(ctx, &req)
	return &resp.Data, err
}
func (r *CloudRepository) RevokeMemberTrustGroup( ctx context.Context, req trustgroup_domain.RemoveMemberFromTrustGroupRequest) (*trustgroup_domain.TrustGroup, error) {
	resp, err := r.client.RemoveMemberFromTrustGroup(ctx, &req)
	return &resp.Data, err
}

func (r *CloudRepository) GetTrustGroup( ctx context.Context, req trustgroup_domain.GetTrustGroupRequest) (*trustgroup_domain.TrustGroup, error) {
	trustGroup, err := r.client.GetTrustGroup(ctx, &req)
	return &trustGroup.Data, err
}
func (r *CloudRepository) ListTrustGroup( ctx context.Context, req trustgroup_domain.ListTrustGroupsRequest) ([]trustgroup_domain.TrustGroup, error) {
	trustGroup, err := r.client.ListTrustGroups(ctx, &req)
	return trustGroup.Data, err
}

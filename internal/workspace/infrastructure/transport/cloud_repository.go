package transport

import (
	"context"

	workspace_domain "vault-app/internal/workspace/domain"
)

type CloudRepository struct {
	client workspace_domain.Repository
}

func NewCloudRepository(client workspace_domain.Repository) *CloudRepository {
	return &CloudRepository{
		client: client,
	}
}

func (r *CloudRepository) CreateWorkspace(ctx context.Context, req workspace_domain.CreateRequest, signature string) (*workspace_domain.Workspace, error) {
	resp, err := r.client.CreateWorkspace(ctx, req)
	return &resp.Data, err
}

func (r *CloudRepository) UpdateWorkspace(
	ctx context.Context, req workspace_domain.UpdateRequest) (*workspace_domain.Workspace, error) {

	resp, err := r.client.UpdateWorkspace(ctx, req)
	return &resp.Data, err
}

func (r *CloudRepository) DeleteWorkspace( ctx context.Context, req workspace_domain.DeleteRequest) (*workspace_domain.Workspace, error) {
	err := r.client.DeleteWorkspace(ctx, req)
	return nil, err
}

func (r *CloudRepository) GetWorkspace( ctx context.Context, req workspace_domain.GetRequest) (*workspace_domain.Workspace, error) {
	workspace, err := r.client.GetWorkspace(ctx, req)
	return workspace, err
}
func (r *CloudRepository) ListWorkspace( ctx context.Context, req workspace_domain.ListRequest) ([]workspace_domain.Workspace, error) {
	workspace, err := r.client.ListWorkspace(ctx, req)
	return workspace, err
}

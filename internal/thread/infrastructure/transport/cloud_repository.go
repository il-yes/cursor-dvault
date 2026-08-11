package transport

import (
	"context"

	thread_domain "vault-app/internal/thread/domain"
)

type CloudRepository struct {
	client thread_domain.ThreadRepository
}

func NewCloudRepository(client thread_domain.ThreadRepository) *CloudRepository {
	return &CloudRepository{
		client: client,
	}
}

func (r *CloudRepository) CreateThread(ctx context.Context, req thread_domain.CreateThreadRequest, signature string) (*thread_domain.Thread, error) {
	resp, err := r.client.CreateThread(ctx, &req)
	return &resp.Data, err
}

func (r *CloudRepository) UpdateThread(
	ctx context.Context, req thread_domain.UpdateThreadRequest) (*thread_domain.Thread, error) {

	resp, err := r.client.UpdateThread(ctx, &req)
	return &resp.Data, err
}

func (r *CloudRepository) GetThread( ctx context.Context, req thread_domain.GetThreadRequest) (*thread_domain.Thread, error) {
	thread, err := r.client.GetThread(ctx, &req)
	return &thread.Data, err
}
func (r *CloudRepository) ListThread( ctx context.Context, req thread_domain.ListThreadsRequest) ([]thread_domain.Thread, error) {
	thread, err := r.client.ListThreads(ctx, &req)
	return thread.Data, err
}

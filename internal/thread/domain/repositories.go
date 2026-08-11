package thread_domain

import (
	"context"

	tracecore_types "vault-app/internal/tracecore/types"
)


type CreateThreadRequest struct {
	Thread Thread
}

type ListThreadsRequest struct {
	ChannelID string
}
type GetThreadRequest struct {
	ThreadID string
}
type UpdateThreadRequest struct {
	Thread Thread
}


type ThreadRepository interface {
	CreateThread(ctx context.Context, req *CreateThreadRequest) (*tracecore_types.CloudResponse[Thread], error)
	ListThreads(ctx context.Context, req *ListThreadsRequest) (*tracecore_types.CloudResponse[[]Thread], error)
	GetThread(ctx context.Context, req *GetThreadRequest) (*tracecore_types.CloudResponse[Thread], error)
	UpdateThread(ctx context.Context, req *UpdateThreadRequest) (*tracecore_types.CloudResponse[Thread], error)
}

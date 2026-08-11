package thread_events

import (
	"context"

	thread_domain "vault-app/internal/thread/domain"
)

// -------- EVENTS --------

// -------- EVENT BUS --------
type ThreadEventBus interface {
	PublishThreadCreated(ctx context.Context, event thread_domain.ThreadCreated) error
	SubscribeToThreadCreated(handler func(ctx context.Context, event thread_domain.ThreadCreated)) error

	PublishThreadUpdated(ctx context.Context, event thread_domain.ThreadUpdated) error
	SubscribeToThreadUpdated(handler func(ctx context.Context, event thread_domain.ThreadUpdated)) error
}

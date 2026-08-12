package thread_infrastructure_eventbus

import (
	"context"
	"sync"

	thread_domain "vault-app/internal/thread/domain"
)

type MemoryBus struct {
	threadCreatedSubscribers []func(ctx context.Context, event thread_domain.ThreadCreated)
	threadUpdatedSubscribers []func(ctx context.Context, event thread_domain.ThreadUpdated)
	lock        sync.RWMutex
}

func NewMemoryBus() *MemoryBus {
	return &MemoryBus{
		threadCreatedSubscribers: make([]func(ctx context.Context, event thread_domain.ThreadCreated), 0),
		threadUpdatedSubscribers: make([]func(ctx context.Context, event thread_domain.ThreadUpdated), 0),
	}
}

// ThreadCreated Event
func (mb *MemoryBus) PublishThreadCreated(ctx context.Context, event thread_domain.ThreadCreated) error {
	mb.lock.RLock()
	defer mb.lock.RUnlock()
	for _, h := range mb.threadCreatedSubscribers {
		go h(ctx, event)
	}
	return nil
}
func (mb *MemoryBus) SubscribeToThreadCreated(handler func(ctx context.Context, event thread_domain.ThreadCreated)) error {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.threadCreatedSubscribers = append(mb.threadCreatedSubscribers, handler)
	return nil
}

func (mb *MemoryBus) PublishThreadUpdated(ctx context.Context, event thread_domain.ThreadUpdated) error {
	mb.lock.RLock()
	defer mb.lock.RUnlock()
	for _, h := range mb.threadUpdatedSubscribers {
		go h(ctx, event)
	}
	return nil
}
func (mb *MemoryBus) SubscribeToThreadUpdated(handler func(ctx context.Context, event thread_domain.ThreadUpdated)) error {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.threadUpdatedSubscribers = append(mb.threadUpdatedSubscribers, handler)
	return nil
}

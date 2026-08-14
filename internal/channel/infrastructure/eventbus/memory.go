package channel_eventbus

import (
	"context"
	"sync"

	channel_events "vault-app/internal/channel/application/events"
	channel_domain "vault-app/internal/channel/domain"
)

type MemoryEventBus struct {
	mu                         sync.RWMutex
	createdHandlers            []func(ctx context.Context, event channel_domain.ChannelCreated)
	revokedHandlers            []func(ctx context.Context, event channel_domain.ChannelRevoked)
	deletedHandlers            []func(ctx context.Context, event channel_domain.ChannelDeleted)
	archivedHandlers           []func(ctx context.Context, event channel_domain.ChannelArchived)
}

func NewMemoryEventBus() channel_events.ChannelEventBus {
	return &MemoryEventBus{}
}

func (b *MemoryEventBus) PublishChannelCreated(ctx context.Context, event channel_domain.ChannelCreated) error {
	b.mu.RLock()
	handlers := append([]func(ctx context.Context, event channel_domain.ChannelCreated){}, b.createdHandlers...)
	b.mu.RUnlock()

	for _, h := range handlers {
		h(ctx, event)
	}
	return nil
}

func (b *MemoryEventBus) SubscribeToChannelCreated(handler func(ctx context.Context, event channel_domain.ChannelCreated)) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.createdHandlers = append(b.createdHandlers, handler)
	return nil
}

func (b *MemoryEventBus) PublishChannelRevoked(ctx context.Context, event channel_domain.ChannelRevoked) error {
	b.mu.RLock()
	handlers := append([]func(ctx context.Context, event channel_domain.ChannelRevoked){}, b.revokedHandlers...)
	b.mu.RUnlock()

	for _, h := range handlers {
		h(ctx, event)
	}
	return nil
}

func (b *MemoryEventBus) SubscribeToChannelRevoked(handler func(ctx context.Context, event channel_domain.ChannelRevoked)) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.revokedHandlers = append(b.revokedHandlers, handler)
	return nil
}

func (b *MemoryEventBus) PublishChannelDeleted(ctx context.Context, event channel_domain.ChannelDeleted) error {
	b.mu.RLock()
	handlers := append([]func(ctx context.Context, event channel_domain.ChannelDeleted){}, b.deletedHandlers...)
	b.mu.RUnlock()

	for _, h := range handlers {
		h(ctx, event)
	}
	return nil
}

func (b *MemoryEventBus) SubscribeToChannelDeleted(handler func(ctx context.Context, event channel_domain.ChannelDeleted)) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.deletedHandlers = append(b.deletedHandlers, handler)
	return nil
}

func (b *MemoryEventBus) PublishChannelArchived(ctx context.Context, event channel_domain.ChannelArchived) error {
	b.mu.RLock()
	handlers := append([]func(ctx context.Context, event channel_domain.ChannelArchived){}, b.archivedHandlers...)
	b.mu.RUnlock()

	for _, h := range handlers {
		h(ctx, event)
	}
	return nil
}

func (b *MemoryEventBus) SubscribeToChannelArchived(handler func(ctx context.Context, event channel_domain.ChannelArchived)) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.archivedHandlers = append(b.archivedHandlers, handler)
	return nil
}

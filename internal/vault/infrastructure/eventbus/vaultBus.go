package vault_infrastructure_eventbus

import 	(
	"context"
	"sync"
	"vault-app/internal/vault/application/events"	
)

type MemoryBus struct {
	subscribers []func(ctx context.Context, event vault_events.VaultOpened)
	lock        sync.RWMutex
}

func NewMemoryBus() *MemoryBus {
	return &MemoryBus{
		subscribers: make([]func(ctx context.Context, event vault_events.VaultOpened), 0),
	}
}

func (mb *MemoryBus) PublishVaultOpened(ctx context.Context, event vault_events.VaultOpened) error {
	mb.lock.RLock()
	defer mb.lock.RUnlock()
	for _, h := range mb.subscribers {	
		go h(ctx, event)
	}
	return nil
}

func (mb *MemoryBus) SubscribeToVaultOpened(handler func(ctx context.Context, event vault_events.VaultOpened)) error {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.subscribers = append(mb.subscribers, handler)
	return nil
}


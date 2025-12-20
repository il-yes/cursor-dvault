package subscription_infrastructure_eventbus

import (
	"context"
	"sync"

	subscription_application_eventbus "vault-app/internal/subscription/application"
)

type MemoryBus struct {
	activationHandlers []func(context.Context, subscription_application_eventbus.SubscriptionActivated)
	creationHandlers   []func(context.Context, subscription_application_eventbus.SubscriptionCreated)
	lock               sync.RWMutex
}

func NewMemoryBus() *MemoryBus {
	return &MemoryBus{
		activationHandlers: []func(context.Context, subscription_application_eventbus.SubscriptionActivated){},
		creationHandlers:   []func(context.Context, subscription_application_eventbus.SubscriptionCreated){},
	}
}

func (bus *MemoryBus) PublishActivated(ctx context.Context, event subscription_application_eventbus.SubscriptionActivated) error {
	bus.lock.RLock()
	defer bus.lock.RUnlock()
	for _, h := range bus.activationHandlers {
		go h(ctx, event)
	}
	return nil
}

func (bus *MemoryBus) SubscribeToActivation(handler func(context.Context, subscription_application_eventbus.SubscriptionActivated)) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	bus.activationHandlers = append(bus.activationHandlers, handler)
	return nil
}

func (bus *MemoryBus) PublishCreated(ctx context.Context, event subscription_application_eventbus.SubscriptionCreated) error {
	bus.lock.RLock()
	defer bus.lock.RUnlock()
	for _, h := range bus.creationHandlers {
		go h(ctx, event)
	}
	return nil
}

func (bus *MemoryBus) SubscribeToCreation(handler func(ctx context.Context, event subscription_application_eventbus.SubscriptionCreated)) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	bus.creationHandlers = append(bus.creationHandlers, handler)
	return nil
}
package onboarding_infrastructure_eventbus

import (
	"context"
	"sync"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
)



type MemoryBus struct {
	creationHandlers []func(onboarding_application_events.AccountCreatedEvent)
	lock             sync.Mutex	
	activatedHandlers []func(onboarding_application_events.SubscriptionActivatedEvent)
}

func NewMemoryBus() *MemoryBus {
	return &MemoryBus{
		creationHandlers: make([]func(onboarding_application_events.AccountCreatedEvent), 0),
		activatedHandlers: make([]func(onboarding_application_events.SubscriptionActivatedEvent), 0),
	}
}

func (bus *MemoryBus) PublishCreated(ctx context.Context, event onboarding_application_events.AccountCreatedEvent) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	for _, handler := range bus.creationHandlers {
		go handler(event)
	}
	return nil
}

func (bus *MemoryBus) SubscribeToAccountCreation(handler func(onboarding_application_events.AccountCreatedEvent)) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	bus.creationHandlers = append(bus.creationHandlers, handler)
	return nil
}
func (bus *MemoryBus) PublishSubscriptionActivated(ctx context.Context, evt onboarding_application_events.SubscriptionActivatedEvent) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	for _, handler := range bus.activatedHandlers {
		go handler(evt)
	}
	return nil
}
func (bus *MemoryBus) SubscribeToSubscriptionActivation(handler func(onboarding_application_events.SubscriptionActivatedEvent)) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	bus.activatedHandlers = append(bus.activatedHandlers, handler)
	return nil
}
	



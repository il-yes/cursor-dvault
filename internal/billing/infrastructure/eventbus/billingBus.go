package billing_infrastructure_eventbus

import (
	"context"
	"sync"
	billing_eventbus "vault-app/internal/billing/application"
)



type MemoryBus struct {
	creationHandlers[]func(context.Context, billing_eventbus.PaymentMethodAddedEvent)
	lock             sync.Mutex		
}

func NewMemoryBus() *MemoryBus {
	return &MemoryBus{
		creationHandlers: make([]func(context.Context, billing_eventbus.PaymentMethodAddedEvent), 0),
	}
}

func (bus *MemoryBus) PublishPaymentMethodAdded(ctx context.Context, event billing_eventbus.PaymentMethodAddedEvent) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	for _, handler := range bus.creationHandlers {
		go handler(ctx, event)
	}
	return nil
}

func (bus *MemoryBus) SubscribeToPaymentMethodAdded(handler func(context.Context, billing_eventbus.PaymentMethodAddedEvent)) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	bus.creationHandlers = append(bus.creationHandlers, handler)
	return nil
}
	



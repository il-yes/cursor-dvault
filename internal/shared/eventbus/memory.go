package eventbus

import (
	"context"
	"sync"
	identity_eventbus "vault-app/internal/identity/application"
	subscription "vault-app/internal/subscription/domain"
)

// A tiny cross-context memory eventbus used by application wiring in cmd for local use/tests
// It implements identity.EventBus and subscription.EventBus and billing.EventBus signatures as needed.

// identity handlers
var (
	userRegisteredHandlersMu sync.RWMutex
	userRegisteredHandlers   []func(context.Context, identity_eventbus.UserRegistered)
)

func PublishUserRegistered(ctx context.Context, e identity_eventbus.UserRegistered) error {
	userRegisteredHandlersMu.RLock()
	hs := append([]func(context.Context, identity_eventbus.UserRegistered){}, userRegisteredHandlers...)
	userRegisteredHandlersMu.RUnlock()
	for _, h := range hs {
		h(ctx, e)
	}
	return nil
}

func SubscribeToUserRegistered(handler func(context.Context, identity_eventbus.UserRegistered)) error {
	userRegisteredHandlersMu.Lock()
	defer userRegisteredHandlersMu.Unlock()
	userRegisteredHandlers = append(userRegisteredHandlers, handler)
	return nil
}

// subscription handlers (example)
var (
	subscriptionCreatedHandlersMu sync.RWMutex
	subscriptionCreatedHandlers   []func(context.Context, subscription.SubscriptionCreated)
)

func PublishSubscriptionCreated(ctx context.Context, e subscription.SubscriptionCreated) error {
	subscriptionCreatedHandlersMu.RLock()
	hs := append([]func(context.Context, subscription.SubscriptionCreated){}, subscriptionCreatedHandlers...)
	subscriptionCreatedHandlersMu.RUnlock()
	for _, h := range hs {
		h(ctx, e)
	}
	return nil
}

func SubscribeToSubscriptionCreated(handler func(context.Context, subscription.SubscriptionCreated)) error {
	subscriptionCreatedHandlersMu.Lock()
	defer subscriptionCreatedHandlersMu.Unlock()
	subscriptionCreatedHandlers = append(subscriptionCreatedHandlers, handler)
	return nil
}



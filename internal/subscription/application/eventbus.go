package subscription_eventbus

import (
	"context"
	"sync"
	subscription_domain "vault-app/internal/subscription/domain"
)

var (
	subscriptionCreatedHandlersMu sync.RWMutex
	subscriptionCreatedHandlers   []func(context.Context, subscription_domain.SubscriptionCreated)
)

type EventBus interface {
	PublishSubscriptionCreated(ctx context.Context, e subscription_domain.SubscriptionCreated) error
	SubscribeToSubscriptionCreated(handler func(context.Context, subscription_domain.SubscriptionCreated)) error
}
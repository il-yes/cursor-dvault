package identity_eventbus

import "context"



type UserRegistered struct {
	UserID    string
	IsAnonymous bool
	OccurredAt int64
}


type UserRegisteredHandler func(context.Context, UserRegistered)

// EventBus inbound port for identity events (application layer)
type EventBus interface {
	PublishUserRegistered(ctx context.Context, e UserRegistered) error
	SubscribeToUserRegistered(handler UserRegisteredHandler) error
}
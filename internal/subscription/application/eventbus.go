package subscription_application_eventbus

import "context"
// Fired when a subscription is created (persisted)
type SubscriptionCreated struct {
	SubscriptionID string
	UserID         string
	UserEmail      string
	Password       string
	Tier           string
	TxHash         string
	OccurredAt     int64
}

// Fired when a subscription becomes active and tier features are enabled
type SubscriptionActivated struct {
	SubscriptionID string
	UserID         string
	UserEmail      string
	Password       string
	Tier           string
	TxHash         string
	Ledger         int32
	OccurredAt     int64
}

// Event bus interface with domain-friendly wording
type SubscriptionEventBus interface {
	PublishActivated(ctx context.Context, event SubscriptionActivated) error
	PublishCreated(ctx context.Context, event SubscriptionCreated) error

	SubscribeToActivation(handler func(ctx context.Context, event SubscriptionActivated)) error
	SubscribeToCreation(handler func(ctx context.Context, event SubscriptionCreated)) error
}

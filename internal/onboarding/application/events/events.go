package onboarding_application_events

import (
	"context"
	"time"
)


type AccountCreatedEvent struct {
	UserID string
	StellarPublicKey string
	OccurredAt time.Time
}

type SubscriptionActivatedEvent struct {
	UserID string
	SubscriptionID string
	Tier string
	OccurredAt time.Time
}


type OnboardingEventBus interface {
	PublishCreated(ctx context.Context, event AccountCreatedEvent) error
	SubscribeToAccountCreation(handler func(event AccountCreatedEvent)) error

    PublishSubscriptionActivated(ctx context.Context, event SubscriptionActivatedEvent) error
	SubscribeToSubscriptionActivation(handler func(event SubscriptionActivatedEvent)) error
}

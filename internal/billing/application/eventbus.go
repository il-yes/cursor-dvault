package billing_eventbus

import "context"

type PaymentMethodAddedEvent struct {
	InstrumentID string
	UserID       string
	Method       string
	OccurredAt   int64
}

type PaymentMethodAddedHandler func(context.Context, PaymentMethodAddedEvent)

type EventBus interface {
	PublishPaymentMethodAdded(ctx context.Context, e PaymentMethodAddedEvent) error
	SubscribeToPaymentMethodAdded(handler func(ctx context.Context, event PaymentMethodAddedEvent)) error
}
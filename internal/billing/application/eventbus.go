package billing_eventbus

import "context"

type PaymentMethodAdded struct {
	InstrumentID string
	UserID       string
	Method       string
	OccurredAt   int64
}

type PaymentMethodAddedHandler func(context.Context, PaymentMethodAdded)

type EventBus interface {
	PublishPaymentMethodAdded(ctx context.Context, e PaymentMethodAdded) error
	SubscribeToPaymentMethodAdded(handler PaymentMethodAddedHandler) error
}
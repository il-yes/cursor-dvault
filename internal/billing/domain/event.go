package billing_domain

type PaymentMethodAdded struct {
	InstrumentID string
	UserID       string
	Method       string
	OccurredAt   int64
}
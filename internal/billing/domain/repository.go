package billing_domain

import "context"

type Repository interface {
	Save(ctx context.Context, b *BillingInstrument) error
	FindByUserID(ctx context.Context, userID string) ([]*BillingInstrument, error)
}
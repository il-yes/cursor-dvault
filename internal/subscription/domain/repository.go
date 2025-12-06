package subscription_domain

import "context"

type Repository interface {
	Save(ctx context.Context, s *Subscription) error
	FindByUserID(ctx context.Context, userID string) (*Subscription, error)
	
}
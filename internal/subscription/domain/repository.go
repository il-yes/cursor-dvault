package subscription_domain

import "context"

type SubscriptionRepository interface {
	Save(ctx context.Context, s *Subscription) error
	FindByUserID(ctx context.Context, userID string) (*Subscription, error)
	GetByID(ctx context.Context, id string) (*Subscription, error)
	
}

type UserRepository interface {
	Save(ctx context.Context, s *UserSubscription) error
	FindByUserID(ctx context.Context, userID string) (*UserSubscription, error)
}
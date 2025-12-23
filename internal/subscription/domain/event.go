package subscription_domain

type SubscriptionCreated struct {
	SubscriptionID string
	UserID         string
	Tier           SubscriptionTier
	OccurredAt     int64
}

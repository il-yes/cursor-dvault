package subscription_domain

import "time"

// Subscription aggregate
type SubscriptionTier string

const (
	TierFree    SubscriptionTier = "free"
	TierPro     SubscriptionTier = "pro"
	TierProPlus SubscriptionTier = "pro_plus"
	TierBusiness SubscriptionTier = "business"
)

type Subscription struct {
	ID        string
	UserID    string
	Tier      SubscriptionTier
	Active    bool
	CreatedAt time.Time
}
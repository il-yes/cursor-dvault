package stellar_recovery_domain

import "context"


type UserRepository interface {
	GetByStellarPublicKey(ctx context.Context, publicKey string) (*User, error)
}

type VaultRepository interface {
	GetByUserID(ctx context.Context, userID string) (*Vault, error)
}

type SubscriptionRepository interface {
	GetActiveByUserID(ctx context.Context, userID string) (*Subscription, error)
}





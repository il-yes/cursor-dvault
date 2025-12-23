package stellar_recovery_domain

import "time"

type Vault struct {
	ID             string
	CreatedAt      time.Time
	StorageUsedGB  float64
	StorageQuotaGB int
	LastSyncedAt   *time.Time
	IPFSNodeID     string
	PinataPinID    string
}

type User struct {
	ID                       string
	Email                    string
	IsAnonymous              bool
	StellarPublicKey         string
	EncryptedSecretKey       string
	SubscriptionTier         string
	SubscriptionID           *string
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type Subscription struct {
	ID        string
	UserID    string
	Status    string
	Tier      string
	CreatedAt time.Time
}

type VaultCheckResult struct {
	VaultExists      bool
	StellarPublicKey string
	Vault            *Vault
	Subscription     *Subscription
}

type RecoveredVault struct {
	User         *User
	Vault        *Vault
	Subscription *Subscription
	SessionToken string
}

type ImportedKey struct {
	StellarPublic string
	StellarSecret string
	CanCreate     bool
}

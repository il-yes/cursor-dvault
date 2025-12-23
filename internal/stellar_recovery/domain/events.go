package stellar_recovery_domain

import "time"

type Event struct {
	AggregateID string
	Type        string
	OccurredAt  time.Time
}


type VaultRecovered struct {
	UserID         string
	VaultID        string
	StellarPubKey  string
	OccurredAt     time.Time
}

func NewVaultRecovered(userID, vaultID, pub string) VaultRecovered {
	return VaultRecovered{
		UserID:        userID,
		VaultID:       vaultID,
		StellarPubKey: pub,
		OccurredAt:    time.Now(),
	}
}

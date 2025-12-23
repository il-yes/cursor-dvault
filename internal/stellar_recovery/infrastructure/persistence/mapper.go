package stellar_recovery_persistence

import (
	"time"
	stellar_recovery_domain "vault-app/internal/stellar_recovery/domain"
)

type UserDB struct {
	ID                       string `db:"id"`
	Email                    string `db:"email"`
	StellarPublicKey         string `db:"stellar_public_key"`
	StellarSecretKeyEncrypted string `db:"stellar_secret_key_encrypted"`
	SubscriptionTier         string `db:"subscription_tier"`
	CreatedAt                time.Time `db:"created_at"`
}

func UserDBToDomain(u *UserDB) *stellar_recovery_domain.User {
	return &stellar_recovery_domain.User{
		ID:                   u.ID,
		Email:                u.Email,
		StellarPublicKey:     u.StellarPublicKey,
		EncryptedSecretKey:   u.StellarSecretKeyEncrypted,
		SubscriptionTier:     u.SubscriptionTier,
		CreatedAt:            u.CreatedAt,
	}
}






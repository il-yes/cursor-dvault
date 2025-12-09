package stellar_recovery_persistence

import (
	"context"
	"errors"
	"gorm.io/gorm"
	stellar_recovery_domain "vault-app/internal/stellar_recovery/domain"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) GetByStellarPublicKey(ctx context.Context, publicKey string) (*stellar_recovery_domain.User, error) {
	var user stellar_recovery_domain.User
	result := r.db.WithContext(ctx).Where("stellar_public_key = ?", publicKey).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

package stellar_recovery_persistence

import (
	"context"
	"errors"
	"gorm.io/gorm"
	stellar_recovery_domain "vault-app/internal/stellar_recovery/domain"
)

type GormVaultRepository struct {
	db *gorm.DB
}

func NewGormVaultRepository(db *gorm.DB) *GormVaultRepository {
	return &GormVaultRepository{db: db}
}

func (r *GormVaultRepository) GetByUserID(ctx context.Context, userID string) (*stellar_recovery_domain.Vault, error) {
	var vault stellar_recovery_domain.Vault
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&vault)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("vault not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &vault, nil
}

package stellar_recovery_persistence

import (
	"context"
	"errors"
	"gorm.io/gorm"
	stellar_recovery_domain "vault-app/internal/stellar_recovery/domain"
)

type GormSubscriptionRepository struct {
	db *gorm.DB
}

func NewGormSubscriptionRepository(db *gorm.DB) *GormSubscriptionRepository {
	return &GormSubscriptionRepository{db: db}
}

func (r *GormSubscriptionRepository) GetActiveByUserID(ctx context.Context, userID string) (*stellar_recovery_domain.Subscription, error) {
	var sub stellar_recovery_domain.Subscription
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, "active").
		First(&sub)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // optional subscription
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &sub, nil
}

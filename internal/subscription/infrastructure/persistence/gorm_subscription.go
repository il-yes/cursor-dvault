package subscription_persistence

import (
	"context"
	"vault-app/internal/logger/logger"
	"vault-app/internal/subscription/domain"

	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	DB *gorm.DB
	Logger *logger.Logger	
}

func NewSubscriptionRepository(db *gorm.DB, logger *logger.Logger) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db, Logger: logger}	
}

func (r *SubscriptionRepository) Save(ctx context.Context, s *subscription_domain.Subscription) error {
	sDB := SubscriptionToDB(s)
	if err := r.DB.Create(sDB).Error; err != nil {
		return err
	}
	return nil
}

func (r *SubscriptionRepository) FindByUserID(ctx context.Context, userID string) (*subscription_domain.Subscription, error) {
	var s subscription_domain.Subscription
	if err := r.DB.Where("user_id = ?", userID).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}
func (r *SubscriptionRepository) GetByID(ctx context.Context, id string) (*subscription_domain.Subscription, error) {
	var s subscription_domain.Subscription
	if err := r.DB.Where("id = ?", id).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}	
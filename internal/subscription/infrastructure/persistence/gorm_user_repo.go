package subscription_persistence

import (
	"context"
	"vault-app/internal/logger/logger"
	subscription_domain "vault-app/internal/subscription/domain"

	"gorm.io/gorm"
)

type UserSubscriptionRepository struct {
	DB *gorm.DB
	Logger *logger.Logger	
}

func NewUserSubscriptionRepository(db *gorm.DB, logger *logger.Logger) *UserSubscriptionRepository {
	return &UserSubscriptionRepository{DB: db, Logger: logger}	
}

func (r *UserSubscriptionRepository) Save(ctx context.Context, us *subscription_domain.UserSubscription) error {
	usDB := UserSubscriptionToDB(us)
	if err := r.DB.Create(usDB).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserSubscriptionRepository) FindByEmail(ctx context.Context, email string) (*subscription_domain.UserSubscription, error) {
	var s subscription_domain.UserSubscription
	if err := r.DB.Where("email = ?", email).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *UserSubscriptionRepository) FindByUserID(ctx context.Context, userID string) (*subscription_domain.UserSubscription, error) {
	var s subscription_domain.UserSubscription
	if err := r.DB.Where("user_id = ?", userID).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}
	
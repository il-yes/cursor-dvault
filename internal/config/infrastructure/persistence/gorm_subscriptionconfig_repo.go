package persistence

import (
	app_config_domain "vault-app/internal/config/domain"

	"gorm.io/gorm"
)


type GormSubscriptionConfigRepository struct {
	db *gorm.DB
}

func NewGormSubscriptionConfigRepository(db *gorm.DB) *GormSubscriptionConfigRepository {
	return &GormSubscriptionConfigRepository{db: db}
}

func (r *GormSubscriptionConfigRepository) Create(subscriptionConfig *app_config_domain.SubscriptionConfig) error {
	return r.db.Create(&subscriptionConfig).Error
}

func (r *GormSubscriptionConfigRepository) Find(id string) (*app_config_domain.SubscriptionConfig, error) {
	var subscriptionConfig app_config_domain.SubscriptionConfig
	err := r.db.First(&subscriptionConfig, "id = ?", id).Error
	return &subscriptionConfig, err
}

func (r *GormSubscriptionConfigRepository) FindByUserIDAndVaultName(userID string, vaultName string) (app_config_domain.SubscriptionConfig, error) {
	var subscriptionConfig app_config_domain.SubscriptionConfig
	err := r.db.Where("user_id = ? AND vault_name = ?", userID, vaultName).First(&subscriptionConfig).Error
	return subscriptionConfig, err
}		

func (r *GormSubscriptionConfigRepository) FindAll(userID string) ([]app_config_domain.SubscriptionConfig, error) {
	var subscriptionConfigs []app_config_domain.SubscriptionConfig
	err := r.db.Where("user_id = ?", userID).Find(&subscriptionConfigs).Error
	return subscriptionConfigs, err
}

func (r *GormSubscriptionConfigRepository) Update(id string, subscriptionConfig *app_config_domain.SubscriptionConfig) error {
	return r.db.Where("id = ?", id).Updates(&subscriptionConfig).Error
}

func (r *GormSubscriptionConfigRepository) Delete(id string) error {
	return r.db.Delete(&app_config_domain.SubscriptionConfig{}, "id = ?", id).Error
}
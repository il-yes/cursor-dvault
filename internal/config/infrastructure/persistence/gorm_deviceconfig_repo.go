package persistence

import (
	app_config_domain "vault-app/internal/config/domain"

	"gorm.io/gorm"
)




type GormDeviceConfigRepository struct {
	db *gorm.DB
}

func NewGormDeviceConfigRepository(db *gorm.DB) *GormDeviceConfigRepository {
	return &GormDeviceConfigRepository{db: db}
} 	

func (r *GormDeviceConfigRepository) Create(dc *app_config_domain.DeviceConfig) error {
	return r.db.Create(&dc).Error
}

func (r *GormDeviceConfigRepository) Find(id string) (*app_config_domain.DeviceConfig, error) {
	var dc *app_config_domain.DeviceConfig
	err := r.db.First(&dc, "id = ?", id).Error
	return dc, err
}

func (r *GormDeviceConfigRepository) FindByUserIDAndVaultName(userID string, vaultName string) ([]app_config_domain.DeviceConfig, error) {
	var deviceConfigs []app_config_domain.DeviceConfig
	err := r.db.Where("user_id = ? AND vault_name = ?", userID, vaultName).Find(&deviceConfigs).Error
	return deviceConfigs, err
}

func (r *GormDeviceConfigRepository) Update(id string, dc *app_config_domain.DeviceConfig) error {
	return r.db.Where("id = ?", id).Updates(&dc).Error
} 	

func (r *GormDeviceConfigRepository) Delete(id string) error {
	return r.db.Delete(&app_config_domain.DeviceConfig{}, "id = ?", id).Error
}

func (r *GormDeviceConfigRepository) FindAll(userID string) ([]*app_config_domain.DeviceConfig, error) {
	var deviceConfigs []*app_config_domain.DeviceConfig
	err := r.db.Where("user_id = ?", userID).Find(&deviceConfigs).Error
	return deviceConfigs, err
}
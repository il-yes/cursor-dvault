package persistence	

import (
	"gorm.io/gorm"
	app_config_domain "vault-app/internal/config/domain"
)

type GormUserConfigRepository struct {
	db *gorm.DB
}

func NewGormUserConfigRepository(db *gorm.DB) *GormUserConfigRepository {
	return &GormUserConfigRepository{db: db}
}	

func (r *GormUserConfigRepository) CreateUserConfig(userConfig *app_config_domain.UserConfig) error {
	return r.db.Create(userConfig).Error
}

func (r *GormUserConfigRepository) GetUserConfig(id string) (*app_config_domain.UserConfig, error) {
	var userConfig app_config_domain.UserConfig
	if err := r.db.First(&userConfig, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &userConfig, nil
}	

func (r *GormUserConfigRepository) UpdateUserConfig(userConfig *app_config_domain.UserConfig) error {
	return r.db.Save(userConfig).Error
}

func (r *GormUserConfigRepository) DeleteUserConfig(id string) error {
	return r.db.Delete(&app_config_domain.UserConfig{}, "id = ?", id).Error
}	
	
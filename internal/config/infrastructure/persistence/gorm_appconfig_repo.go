package persistence

import (
	"gorm.io/gorm"
	app_config_domain "vault-app/internal/config/domain"
)

type GormAppConfigRepository struct {
	db *gorm.DB
}

func NewGormAppConfigRepository(db *gorm.DB) *GormAppConfigRepository {
	return &GormAppConfigRepository{db: db}
}	

func (r *GormAppConfigRepository) CreateAppConfig(appConfig *app_config_domain.AppConfig) error {
	return r.db.Create(appConfig).Error
}

func (r *GormAppConfigRepository) GetAppConfig(id string) (*app_config_domain.AppConfig, error) {
	var appConfig app_config_domain.AppConfig
	if err := r.db.First(&appConfig, "user_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &appConfig, nil
}	

func (r *GormAppConfigRepository) UpdateAppConfig(appConfig *app_config_domain.AppConfig) error {
	return r.db.Save(appConfig).Error
}

func (r *GormAppConfigRepository) DeleteAppConfig(id string) error {
	return r.db.Delete(&app_config_domain.AppConfig{}, "user_id = ?", id).Error
}	

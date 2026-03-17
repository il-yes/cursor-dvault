package persistence

import (
	app_config_domain "vault-app/internal/config/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormVaultConfigRepository struct {
	db *gorm.DB
}

func NewGormVaultConfigRepository(db *gorm.DB) *GormVaultConfigRepository {
	return &GormVaultConfigRepository{db: db}
} 	

func (r *GormVaultConfigRepository) Create(vaultConfig *app_config_domain.VaultConfigBeta) (*app_config_domain.VaultConfigBeta, error) {
	if vaultConfig.ID == "" {
		vaultConfig.ID = uuid.New().String()
	}
	
	err := r.db.Create(&vaultConfig).Error
	return vaultConfig, err
}

func (r *GormVaultConfigRepository) Find(id string) (*app_config_domain.VaultConfigBeta, error) {
	var vaultConfig app_config_domain.VaultConfigBeta
	err := r.db.First(&vaultConfig, "id = ?", id).Error
	return &vaultConfig, err
}

func (r *GormVaultConfigRepository) FindByUserIDAndVaultName(userID string, vaultName string) (app_config_domain.VaultConfigBeta, error) {
	var vaultConfig app_config_domain.VaultConfigBeta
	err := r.db.Where("user_id = ? AND vault_name = ?", userID, vaultName).First(&vaultConfig).Error
	return vaultConfig, err
}	

func (r *GormVaultConfigRepository) Update(id string, vaultConfig *app_config_domain.VaultConfigBeta) error {
	return r.db.Where("id = ?", id).Updates(&vaultConfig).Error
} 

func (r *GormVaultConfigRepository) Delete(id string) error {
	return r.db.Delete(&app_config_domain.VaultConfigBeta{}, "id = ?", id).Error
}

func (r *GormVaultConfigRepository) FindAll(userID string, vaultName string) ([]app_config_domain.VaultConfigBeta, error) {
	var vaultConfigs []app_config_domain.VaultConfigBeta
	err := r.db.Where("user_id = ?", userID).Where("vault_name = ?", vaultName).Find(&vaultConfigs).Error
	return vaultConfigs, err
}

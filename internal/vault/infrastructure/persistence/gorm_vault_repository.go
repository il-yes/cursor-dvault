package vaults_persistence

import (
	vaults_domain "vault-app/internal/vault/domain"

	"gorm.io/gorm"
)

type GormVaultRepository struct {
	db *gorm.DB
}

func NewGormVaultRepository(db *gorm.DB) *GormVaultRepository {
	return &GormVaultRepository{db: db}
}

func (r *GormVaultRepository) SaveVault(vault *vaults_domain.Vault) error {
	vdb := VaultDomainToMapper(vault)
	return r.db.Create(&vdb).Error
}

func (r *GormVaultRepository) GetVault(vaultID string) (*vaults_domain.Vault, error) {
	var vault VaultMapper
	if err := r.db.First(&vault, vaultID).Error; err != nil {
		return nil, err
	}
	return vault.ToDomain(), nil
}	

func (r *GormVaultRepository) UpdateVault(vault *vaults_domain.Vault) error {
	vdb := VaultDomainToMapper(vault)
	return r.db.Save(&vdb).Error
}

func (r *GormVaultRepository) DeleteVault(vaultID string) error {
	return r.db.Delete(&VaultMapper{}, vaultID).Error
}	

func (r *GormVaultRepository) GetLatestByUserID(id string) (*vaults_domain.Vault, error) {
	var record VaultMapper

	if err := r.db.Order("created_at DESC").First(&record, "user_id = ?", id).Error; err != nil {
		return nil, err
	}
	return record.ToDomain(), nil
}
	

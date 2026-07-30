package vaults_persistence

import (
	vaults_domain "vault-app/internal/vault/domain"

	"gorm.io/gorm"
)

type GormVaultKeyringRepository struct {
	db *gorm.DB
}

func NewGormVaultKeyringRepository(db *gorm.DB) *GormVaultKeyringRepository {
	return &GormVaultKeyringRepository{db: db}
}

func (vk *GormVaultKeyringRepository) Save(model vaults_domain.VaultKeyring) error {
	mapper := KeyringDomainToMapper(model)
	return vk.db.Create(&mapper).Error
}

func (vk *GormVaultKeyringRepository) GetByUserID(uID string) (*vaults_domain.VaultKeyring, error) {
	var record KeyringMapper
	if err := vk.db.Order("created_at DESC").First(&record, "user_id = ?", uID).Error; err != nil {
		return nil, err
	}
	return record.ToDomain(), nil
}

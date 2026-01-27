package vaults_persistence

import (
	"fmt"
	vaults_domain "vault-app/internal/vault/domain"

	"gorm.io/gorm"
)

type GormFolderRepository struct {
	db *gorm.DB
}	

func NewGormFolderRepository(db *gorm.DB) *GormFolderRepository {
	return &GormFolderRepository{db: db}
}

func (r *GormFolderRepository) SaveFolder(folder *vaults_domain.Folder) error {
	return r.db.Create(folder).Error
}
func (r *GormFolderRepository) GetFolder(folderID string) (*vaults_domain.Folder, error) {
	var folder vaults_domain.Folder
	if err := r.db.Where("id = ?", folderID).Find(&folder).Error; err != nil {
		return nil, fmt.Errorf("❌ server internal error: %w", err)
	}
	return &folder, nil
}	
func (r *GormFolderRepository) UpdateFolder(folder *vaults_domain.Folder) error {
	return r.db.Save(folder).Error
}	
func (r *GormFolderRepository) DeleteFolder(folderID string) error {
	return r.db.Delete(&vaults_domain.Folder{}, "id = ?", folderID).Error
}	
func (r *GormFolderRepository) GetFoldersByUserID(userID string) ([]vaults_domain.Folder, error) {
	var folders []vaults_domain.Folder
	if err := r.db.Where("user_id = ?", userID).Find(&folders).Error; err != nil {
		return nil, fmt.Errorf("❌ server internal error: %w", err)
	}
	return folders, nil
} 

func (r *GormFolderRepository) GetFoldersByVault(vaultCID string) ([]vaults_domain.Folder, error) {
	var folders []vaults_domain.Folder
	if err := r.db.Where("vault_cid = ?", vaultCID).Find(&folders).Error; err != nil {
		return nil, fmt.Errorf("❌ server internal error: %w", err)
	}
	return folders, nil
}
func (r *GormFolderRepository) GetFolderById(id string) (*vaults_domain.Folder, error) {
	var folder vaults_domain.Folder
	if err := r.db.Where("id = ?", id).Find(&folder).Error; err != nil {
		return nil, fmt.Errorf("❌ server internal error: %w", err)
	}
	return &folder, nil
}
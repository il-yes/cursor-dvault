package auth_persistence

import (
	auth_domain "vault-app/internal/auth/domain"

	"gorm.io/gorm"
)

type GormPrincipalRepository struct {
	db *gorm.DB
}

func NewGormPrincipalRepository(db *gorm.DB) *GormPrincipalRepository {
	return &GormPrincipalRepository{db: db}
}

func (r *GormPrincipalRepository) TableName() string {
	return "auth_principals"
} 

func (r *GormPrincipalRepository) FindByID(id string) (*auth_domain.Principal, error) {
	var user auth_domain.Principal
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}	

func (r *GormPrincipalRepository) FindByEmail(email string) (*auth_domain.Principal, error) {
	var user auth_domain.Principal
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}	

func (r *GormPrincipalRepository) FindByUsername(username string) (*auth_domain.Principal, error) {
	var user auth_domain.Principal
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}		

func (r *GormPrincipalRepository) Save(user *auth_domain.Principal) error {
	return r.db.Save(user).Error
}


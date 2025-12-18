package identity_persistence

import (
	"context"
	identity_domain "vault-app/internal/identity/domain"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}


func (GormUserRepository) TableName() string {
	return "identity_users"
}

func (r *GormUserRepository) Save(ctx context.Context, u *identity_domain.User) error {
	return r.db.Create(u).Error
}

func (r *GormUserRepository) FindByID(ctx context.Context, id string) (*identity_domain.User, error) {
	var u identity_domain.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*identity_domain.User, error) {
	var u identity_domain.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// Ensure interface satisfaction at compile-time
var _ identity_domain.UserRepository = (*GormUserRepository)(nil)

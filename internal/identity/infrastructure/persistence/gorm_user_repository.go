package identity_persistence

import (
	"context"
	utils "vault-app/internal"
	identity_usecase "vault-app/internal/identity/application/usecase"
	identity_domain "vault-app/internal/identity/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IDGen struct {
}

func (g *IDGen) Generate() string {
	return uuid.New().String()
}
func NewIDGenerator() identity_usecase.IDGen {
	var idGen IDGen
	return idGen.Generate
}

type GormUserRepository struct {
	db *gorm.DB
}	

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}


func (r *GormUserRepository) Save(ctx context.Context, u *identity_domain.User) error {
	return r.db.Create(u).Error
}

func (r *GormUserRepository) FindByID(ctx context.Context, id string) (*identity_domain.User, error) {
	var u identity_domain.User
	utils.LogPretty("GormUserRepository - FindByID - id, processing...", id)
	if err := r.db.Model(&identity_domain.User{}).Where("id = ?", id).First(&u).Error; err != nil {
		return nil, err
	}
	utils.LogPretty("GormUserRepository - FindByID - u, processed", u)
	return &u, nil
}

func (r *GormUserRepository) Update(ctx context.Context, u *identity_domain.User) error {
	return r.db.Save(u).Error
}	

func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*identity_domain.User, error) {
	var u identity_domain.User
	if err := r.db.Model(&identity_domain.User{}).Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// Ensure interface satisfaction at compile-time
var _ identity_domain.UserRepository = (*GormUserRepository)(nil)

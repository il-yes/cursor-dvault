package identity_persistence

import (
	"context"
	"fmt"
	utils "vault-app/internal/utils"
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
	if err := r.db.WithContext(ctx).Model(&identity_domain.User{}).Where("id = ?", id).First(&u).Error; err == nil {
		return &u, nil
	}

	type customerRecord struct {
		ID        uint   `gorm:"column:id"`
		Email     string `gorm:"column:email"`
		PublicKey string `gorm:"column:public_key"`
	}
	var cust customerRecord
	if err := r.db.WithContext(ctx).Table("customers").Where("id = ?", id).First(&cust).Error; err == nil {
		user := identity_domain.NewStandardUser(fmt.Sprintf("%d", cust.ID), cust.Email, "")
		user.StellarPublicKey = cust.PublicKey
		return user, nil
	}

	return nil, fmt.Errorf("user not found for id: %s", id)
}

func (r *GormUserRepository) Update(ctx context.Context, u *identity_domain.User) error {
	return r.db.Save(u).Error
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*identity_domain.User, error) {
	var u identity_domain.User
	if err := r.db.WithContext(ctx).Model(&identity_domain.User{}).Where("email = ?", email).First(&u).Error; err == nil {
		return &u, nil
	}

	type customerRecord struct {
		ID        uint   `gorm:"column:id"`
		Email     string `gorm:"column:email"`
		PublicKey string `gorm:"column:public_key"`
	}
	var cust customerRecord
	if err := r.db.WithContext(ctx).Table("customers").Where("email = ?", email).First(&cust).Error; err == nil {
		user := identity_domain.NewStandardUser(fmt.Sprintf("%d", cust.ID), cust.Email, "")
		user.StellarPublicKey = cust.PublicKey
		return user, nil
	}

	return nil, fmt.Errorf("user not found for email: %s", email)
}

func (r *GormUserRepository) FindByPublicKey(ctx context.Context, publicKey string) (*identity_domain.User, error) {
	var u identity_domain.User
	if err := r.db.WithContext(ctx).Model(&identity_domain.User{}).Where("stellar_public_key = ?", publicKey).First(&u).Error; err == nil {
		return &u, nil
	}

	type customerRecord struct {
		ID        uint   `gorm:"column:id"`
		Email     string `gorm:"column:email"`
		PublicKey string `gorm:"column:public_key"`
	}
	var cust customerRecord
	if err := r.db.WithContext(ctx).Table("customers").Where("public_key = ?", publicKey).First(&cust).Error; err == nil {
		user := identity_domain.NewStandardUser(fmt.Sprintf("%d", cust.ID), cust.Email, "")
		user.StellarPublicKey = cust.PublicKey
		return user, nil
	}

	type legacyUserRecord struct {
		ID        uint   `gorm:"column:id"`
		Email     string `gorm:"column:email"`
		PublicKey string `gorm:"column:public_key"`
	}
	var legacyUser legacyUserRecord
	if err := r.db.WithContext(ctx).Table("users").Where("public_key = ?", publicKey).First(&legacyUser).Error; err == nil {
		user := identity_domain.NewStandardUser(fmt.Sprintf("%d", legacyUser.ID), legacyUser.Email, "")
		user.StellarPublicKey = legacyUser.PublicKey
		return user, nil
	}

	return nil, fmt.Errorf("identity user not found for public key: %s", publicKey)
}

// Ensure interface satisfaction at compile-time
var _ identity_domain.UserRepository = (*GormUserRepository)(nil)

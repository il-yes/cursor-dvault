package identity_persistence

import (
	"context"
	"fmt"
	"log"
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
		return u.EnsureAliases(), nil
	}

	type legacyUserRecord struct {
		ID        string `gorm:"column:id"`
		Email     string `gorm:"column:email"`
		Username  string `gorm:"column:username"`
		FirstName string `gorm:"column:first_name"`
		LastName  string `gorm:"column:last_name"`
	}
	var legacyUser legacyUserRecord
	if err := r.db.WithContext(ctx).Table("users").Where("id = ?", id).First(&legacyUser).Error; err == nil {
		user := identity_domain.NewStandardUser(legacyUser.ID, legacyUser.Email, "")
		user.Username = legacyUser.Username
		user.FirstName = legacyUser.FirstName
		user.LastName = legacyUser.LastName
		return user.EnsureAliases(), nil
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
		return user.EnsureAliases(), nil
	}

	return nil, fmt.Errorf("user not found for id: %s", id)
}

func (r *GormUserRepository) Update(ctx context.Context, u *identity_domain.User) error {
	if u == nil || u.ID == "" {
		return fmt.Errorf("cannot update user without valid ID")
	}
	u.EnsureAliases()
	log.Printf("[GormUserRepository.Update] START userID=%s username=%s first_name=%s last_name=%s", u.ID, u.Username, u.FirstName, u.LastName)

	updates := map[string]interface{}{
		"last_connected_at":  u.LastConnectedAt,
		"is_anonymous":       u.IsAnonymous,
		"identity":           u.Identity,
		"stellar_public_key": u.StellarPublicKey,
	}
	if u.Username != "" {
		updates["username"] = u.Username
	}
	if u.FirstName != "" {
		updates["first_name"] = u.FirstName
	}
	if u.LastName != "" {
		updates["last_name"] = u.LastName
	}

	tx := r.db.WithContext(ctx).
		Model(&identity_domain.User{}).
		Where("id = ?", u.ID).
		Updates(updates)

	log.Printf("[GormUserRepository.Update] RESULT identity_users RowsAffected=%d err=%v", tx.RowsAffected, tx.Error)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		log.Printf("[GormUserRepository.Update] RowsAffected=0, saving user row in identity_users table for userID=%s", u.ID)
		if saveErr := r.db.WithContext(ctx).Save(u).Error; saveErr != nil {
			log.Printf("❌ [GormUserRepository.Update] Save fallback error: %v", saveErr)
			return saveErr
		}
	}

	// Also sync legacy users table if present
	_ = r.db.WithContext(ctx).Table("users").Where("id = ?", u.ID).Updates(updates)

	return nil
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*identity_domain.User, error) {
	var u identity_domain.User
	if err := r.db.WithContext(ctx).Model(&identity_domain.User{}).Where("email = ?", email).First(&u).Error; err == nil {
		return u.EnsureAliases(), nil
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
		return user.EnsureAliases(), nil
	}

	return nil, fmt.Errorf("user not found for email: %s", email)
}

func (r *GormUserRepository) FindByPublicKey(ctx context.Context, publicKey string) (*identity_domain.User, error) {
	var u identity_domain.User
	if err := r.db.WithContext(ctx).Model(&identity_domain.User{}).Where("stellar_public_key = ?", publicKey).First(&u).Error; err == nil {
		return u.EnsureAliases(), nil
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
		return user.EnsureAliases(), nil
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
		return user.EnsureAliases(), nil
	}

	return nil, fmt.Errorf("identity user not found for public key: %s", publicKey)
}

// Ensure interface satisfaction at compile-time
var _ identity_domain.UserRepository = (*GormUserRepository)(nil)

package onboarding_persistence


import (
	"gorm.io/gorm"
	"vault-app/internal/onboarding/domain"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) onboarding_domain.UserRepository {
	return &GormUserRepository{db: db}
}

// Ensure GORM uses the exact table name "subscription_db" or "subscriptions".
// Replace "subscriptions" below with your actual table name if different.
func (GormUserRepository) TableName() string {
	return "users"
}

func (r *GormUserRepository) Create(user *onboarding_domain.User) (*onboarding_domain.User, error) {
	userDB := ToUserDB(user)
	if err := r.db.Create(userDB).Error; err != nil {
		return nil, err
	}
	created, err := r.GetByID(userDB.ID)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (r *GormUserRepository) GetByID(id string) (*onboarding_domain.User, error) {
	var userDB UserDB // ← correct GORM model
    if err := r.db.First(&userDB, "id = ?", id).Error; err != nil {
        return nil, err
    }
    user := userDB.ToUser()
	return user, nil
}

func (r *GormUserRepository) List() ([]onboarding_domain.User, error) {
	var users []onboarding_domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *GormUserRepository) Update(user *onboarding_domain.User) error {
	userDB := ToUserDB(user)
	return r.db.Save(userDB).Error
}

func (r *GormUserRepository) Delete(id string) error {
	return r.db.Delete(&UserDB{ID: id}).Error
}

package onboarding_persistence

import (
	"vault-app/internal/onboarding/domain"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) onboarding_domain.UserRepository {
	return &GormUserRepository{db: db}
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
	var usersDB []UserDB
	if err := r.db.Find(&usersDB).Error; err != nil {
		return nil, err
	}
	var users []onboarding_domain.User
	for _, userDB := range usersDB {
		users = append(users, *userDB.ToUser())
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
func (r *GormUserRepository) FindByEmail(email string) (*onboarding_domain.User, error) {
	var userDB UserDB // ← correct GORM model
    if err := r.db.First(&userDB, "email = ?", email).Error; err != nil {
        return nil, err
    }
    user := userDB.ToUser()
	return user, nil
}	

func (r *GormUserRepository) FindUserByEmail(email string) (*onboarding_domain.User, error) {
	return r.FindByEmail(email)
}
func (r *GormUserRepository) FindAll() ([]onboarding_domain.User, error) {
	return r.List()
}	
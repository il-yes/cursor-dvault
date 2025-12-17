package onboarding_persistence

import (
	"time"
	onboarding_domain "vault-app/internal/onboarding/domain"
	"gorm.io/gorm"
	uuid "github.com/google/uuid"
)

type UserDB struct {
	ID string `json:"id" gorm:"primarykey"`
	IsAnonymous bool `json:"is_anonymous"`
	StellarPublicKey string `json:"stellar_public_key"`
	CreatedAt time.Time `json:"created_at"`
	Email string `json:"email"`
	Password string `json:"password"`	
}
func (u *UserDB) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().String()	
	return nil
}

func (u *UserDB) TableName() string {
	return "user_onboarding"
}

func (u *UserDB) ToUser() *onboarding_domain.User {
	return &onboarding_domain.User{
		ID: u.ID,
		IsAnonymous: u.IsAnonymous,
		StellarPublicKey: u.StellarPublicKey,
		CreatedAt: u.CreatedAt,
		Email: u.Email,
		Password: u.Password,
	}
}
func ToUserDB(user *onboarding_domain.User) *UserDB {
	return &UserDB{
		ID: user.ID,
		IsAnonymous: user.IsAnonymous,
		StellarPublicKey: user.StellarPublicKey,
		CreatedAt: user.CreatedAt,
		Email: user.Email,
		Password: user.Password,	
	}
}



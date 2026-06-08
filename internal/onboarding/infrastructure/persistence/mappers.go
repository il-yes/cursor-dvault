package onboarding_persistence

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"

	onboarding_domain "vault-app/internal/onboarding/domain"
)

type UserDB struct {
	ID               string      `json:"id" gorm:"primarykey"`
	IsAnonymous      bool        `json:"is_anonymous"`
	StellarPublicKey string      `json:"stellar_public_key"`
	CreatedAt        time.Time   `json:"created_at"`
	LastConnectedAt  time.Time   `json:"last_connected_at"`
	Email            string      `json:"email"`
	Password         string      `json:"password"`
	UseCases         StringSlice `json:"use_cases" gorm:"type:json"`
}

type StringSlice []string
func (s *StringSlice) Scan(value any) error {
    if value == nil {
        *s = nil
        return nil
    }

    var b []byte
    switch v := value.(type) {
    case []byte:
        b = v
    case string:
        b = []byte(v)
    default:
        return fmt.Errorf("unsupported type %T", value)
    }

    if len(b) == 0 {
        *s = nil
        return nil
    }

    if b[0] != '[' && b[0] != '{' {
        *s = StringSlice{string(b)}
        return nil
    }

    return json.Unmarshal(b, s)
}

func (s StringSlice) Value() (driver.Value, error) {
    if s == nil {
        return []byte("[]"), nil
    }
    return json.Marshal([]string(s))
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
        ID:              u.ID,
        IsAnonymous:     u.IsAnonymous,
        StellarPublicKey: u.StellarPublicKey,
        CreatedAt:       u.CreatedAt,
        Email:           u.Email,
        Password:        u.Password,
        LastConnectedAt: u.LastConnectedAt,
        UseCases:        u.UseCases,
    }
}
func ToUserDB(user *onboarding_domain.User) *UserDB {
    return &UserDB{
        ID:              user.ID,
        IsAnonymous:     user.IsAnonymous,
        StellarPublicKey: user.StellarPublicKey,
        CreatedAt:       user.CreatedAt,
        Email:           user.Email,
        Password:        user.Password,
        LastConnectedAt: user.LastConnectedAt,
        UseCases:        user.UseCases,
    }
}

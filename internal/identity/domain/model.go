package identity_domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IdentityChoice string

const (
    IdentityPersonal    IdentityChoice = "personal"
    IdentityAnonymous   IdentityChoice = "anonymous"
    IdentityTeam        IdentityChoice = "team"
    IdentityCompliance  IdentityChoice = "compliance"
)


// User aggregate
type User struct {
	ID               string
	Email            string
	PasswordHash     string
	IsAnonymous      bool
	StellarPublicKey string
	CreatedAt        time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}
func (u *User) IsStandard() bool {
	return !u.IsAnonymous
}
func NewAnonymousUser(id, stellarPublicKey string) *User {
	return &User{ID: id, IsAnonymous: true, StellarPublicKey: stellarPublicKey, CreatedAt: time.Now()}
}

func NewStandardUser(id, email, passwordHash string) *User {
	return &User{ID: id, Email: email, PasswordHash: passwordHash, CreatedAt: time.Now()}
}

var ErrUserExists = errors.New("user already exists")

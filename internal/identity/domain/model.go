package identity_domain

import (
	"time"
	auth_domain "vault-app/internal/auth/domain"
	"vault-app/internal/models"
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
	LastConnectedAt  string
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}
func (u *User) IsStandard() bool {
	return !u.IsAnonymous
}
func (u *User) ToJwtUser() *auth_domain.JwtUser {
	return &auth_domain.JwtUser{
		ID:       u.ID,
		Username: u.Email,
		Email:    u.Email,
	}
}
func (u *User) ToFormerUser() *models.User {
	return &models.User{
		ID:       u.ID,
		Username: u.Email,
		Email:    u.Email,
		Password: u.PasswordHash,
		CreatedAt: u.CreatedAt,
	}
}

func NewAnonymousUser(id, stellarPublicKey string) *User {
	return &User{ID: id, IsAnonymous: true, StellarPublicKey: stellarPublicKey, CreatedAt: time.Now()}
}

func NewStandardUser(id, email, passwordHash string) *User {
	return &User{ID: id, Email: email, PasswordHash: passwordHash, CreatedAt: time.Now(), LastConnectedAt: time.Now().Format(time.RFC3339)}
}


// UserLoggedIn represents a successful login event
type UserLoggedIn struct {
	UserID    string
	Email     string
	OccurredAt time.Time
}


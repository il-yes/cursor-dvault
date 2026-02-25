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

// ------------------- Standard -------------------------
// User aggregate
type User struct {
	ID               string    `gorm:"primaryKey"`
	Email            string    `gorm:"uniqueIndex;not null"`
	PasswordHash     string    `gorm:"not null"`
	IsAnonymous      bool
	StellarPublicKey string
	CreatedAt        time.Time
	LastConnectedAt  time.Time
}
func (User) TableName() string {
	return "identity_users"
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
func (u *User) OnGenerateApiKey(pk string) *User {
	u.StellarPublicKey = pk
	return u
}

// ------------------- Anonymous -------------------------
func NewAnonymousUser(id, stellarPublicKey string) *User {
	return &User{ID: id, IsAnonymous: true, StellarPublicKey: stellarPublicKey, CreatedAt: time.Now()}
}

func NewStandardUser(id, email, passwordHash string) *User {
	return &User{ID: id, Email: email, PasswordHash: passwordHash, CreatedAt: time.Now(), LastConnectedAt: time.Now()}
}


// UserLoggedIn represents a successful login event
type UserLoggedIn struct {
	UserID    string
	Email     string
	OccurredAt time.Time
}


		
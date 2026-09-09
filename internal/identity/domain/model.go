package identity_domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	auth_domain "vault-app/internal/auth/domain"
	"vault-app/internal/models"
)

type IdentityChoice string

const (
	IdentityPersonal   IdentityChoice = "personal"
	IdentityAnonymous  IdentityChoice = "anonymous"
	IdentityTeam       IdentityChoice = "team"
	IdentityCompliance IdentityChoice = "compliance"
)

// ------------------- Standard -------------------------
// User aggregate
type User struct {
	ID               string         `gorm:"primaryKey" json:"id"`
	Email            string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash     string         `gorm:"not null" json:"-"`
	IsAnonymous      bool           `json:"is_anonymous"`
	Identity         IdentityChoice `gorm:"not null,default:personal" json:"identity"`
	Username         string         `gorm:"column:username" json:"username"`
	UserName         string         `gorm:"-" json:"user_name,omitempty"`
	FirstName        string         `gorm:"column:first_name" json:"first_name"`
	LastName         string         `gorm:"column:last_name" json:"last_name"`
	StellarPublicKey string         `json:"stellar_public_key"`
	CreatedAt        time.Time      `json:"created_at"`
	LastConnectedAt  time.Time      `json:"last_connected_at"`
}

func (User) TableName() string {
	return "identity_users"
}
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
func (u *User) EnsureAliases() *User {
	if u == nil {
		return nil
	}
	if u.UserName == "" {
		u.UserName = u.Username
	}
	if u.Username == "" && u.UserName != "" {
		u.Username = u.UserName
	}
	return u
}
func (u *User) IsStandard() bool {
	return !u.IsAnonymous
}
func (u *User) ToJwtUser() *auth_domain.JwtUser {
	username := u.Username
	if username == "" {
		username = u.Email
	}
	return &auth_domain.JwtUser{
		ID:       u.ID,
		Username: username,
		Email:    u.Email,
	}
}
func (u *User) ToFormerUser() *models.User {
	username := u.Username
	if username == "" {
		username = u.Email
	}
	return &models.User{
		ID:        u.ID,
		Username:  username,
		Email:     u.Email,
		Password:  u.PasswordHash,
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
	UserID     string
	Email      string
	OccurredAt time.Time
}

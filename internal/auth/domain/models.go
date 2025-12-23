package auth_domain

import (
	"time"
	"vault-app/internal/auth"

	"github.com/golang-jwt/jwt/v4"
)

type Credentials struct {
    Email    string
    Password string
}	

type Principal struct {
    UserID   string
    Email    string
    Role     string
    Claims   Claims
}

type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	CookieDomain  string
	CookiePath    string
	CookieName    string
}

type JwtUser struct {
	ID       string    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type TokenPairs struct {
    Token        string `gorm:"column:token;type:text" json:"access_token"`
    RefreshToken string `gorm:"column:refresh_token;type:text" json:"refresh_token"`
    UserID       string    `gorm:"column:user_id" json:"user_id"`
}
func (t *TokenPairs) ToFormerModel() *auth.TokenPairs {
	return &auth.TokenPairs{
		Token:        t.Token,
		RefreshToken: t.RefreshToken,
		UserID:       t.UserID,
	}
}

type Claims struct {
    UserID       string    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`

	jwt.RegisteredClaims
}
func (c *Claims) ToFormerModel() *auth.Claims {
	return &auth.Claims{
		UserID:       c.UserID,
		Username: c.Username,
		Email:    c.Email,
		RegisteredClaims: c.RegisteredClaims,
	}
}	





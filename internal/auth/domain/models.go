package auth_domain

import (
	"time"

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


type Claims struct {
    UserID       string    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`

	jwt.RegisteredClaims
}






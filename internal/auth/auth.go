package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

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
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type TokenPairs struct {
    Token        string `gorm:"column:token;type:text" json:"access_token"`
    RefreshToken string `gorm:"column:refresh_token;type:text" json:"refresh_token"`
    UserID       int    `gorm:"column:user_id" json:"user_id"`
}


type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Email    string `json:"email"`

	jwt.RegisteredClaims
}

// GenerateTokenPair generates both access & refresh tokens
func (j *Auth) GenerateTokenPair(user *JwtUser) (TokenPairs, error) {
    // --- Access Token ---
    accessClaims := &Claims{
        UserID:   user.ID,
        Username: user.Username,
        Email:    user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   fmt.Sprint(user.ID),
            Audience:  jwt.ClaimStrings{j.Audience},
            Issuer:    j.Issuer,
            IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
            ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(j.TokenExpiry)),
        },
    }

    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    signedAccessToken, err := accessToken.SignedString([]byte(j.Secret))
    if err != nil {
        return TokenPairs{}, err
    }

    // --- Refresh Token (simpler claims) ---
    refreshClaims := &Claims{
        UserID: user.ID,
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   fmt.Sprint(user.ID),
            IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
            ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(j.RefreshExpiry)),
        },
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    signedRefreshToken, err := refreshToken.SignedString([]byte(j.Secret))
    if err != nil {
        return TokenPairs{}, err
    }

    // --- Return tokens ---
    return TokenPairs{
        Token:        signedAccessToken,
        RefreshToken: signedRefreshToken,
        UserID:       user.ID,
    }, nil
}

func (j *Auth) VerifyToken(tokenStr string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(j.Secret), nil
    })

    if err != nil {
        if strings.HasPrefix(err.Error(), "token is expired by") {
            return nil, errors.New("expired token")
        }
        return nil, err
    }

    if !token.Valid {
        return nil, errors.New("invalid token")
    }

    if claims.Issuer != j.Issuer {
        return nil, errors.New("invalid issuer")
    }

    return claims, nil
}



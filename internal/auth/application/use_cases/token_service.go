package auth_usecases

import (
	"errors"
	"fmt"
	"strings"
	"time"
	utils "vault-app/internal"
	auth_domain "vault-app/internal/auth/domain"
    "gorm.io/gorm"  
	"github.com/golang-jwt/jwt/v4"
)



type TokenService struct {
	token auth_domain.Auth
	repo auth_domain.AuthRepository
    DB *gorm.DB
}

func NewTokenService(token auth_domain.Auth, repo auth_domain.AuthRepository, db *gorm.DB) *TokenService {
	return &TokenService{token: token, repo: repo, DB: db}
}

// GenerateTokenPair generates both access & refresh tokens
func (j *TokenService) GenerateTokenPair(user *auth_domain.JwtUser) (auth_domain.TokenPairs, error) {
    // --- Access Token ---
    accessClaims := &auth_domain.Claims{
        UserID:   user.ID,
        Username: user.Username,
        Email:    user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   fmt.Sprint(user.ID),
            Audience:  jwt.ClaimStrings{j.token.Audience},
            Issuer:    j.token.Issuer,
            IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
            ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(j.token.TokenExpiry)),
        },
    }

    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    signedAccessToken, err := accessToken.SignedString([]byte(j.token.Secret))
    if err != nil {
        return auth_domain.TokenPairs{}, err
    }

    // --- Refresh Token (simpler claims) ---
    refreshClaims := &auth_domain.Claims{
        UserID: user.ID,
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   fmt.Sprint(user.ID),
            IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
            ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(j.token.RefreshExpiry)),
        },
    }

    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    signedRefreshToken, err := refreshToken.SignedString([]byte(j.token.Secret))
    if err != nil {
        return auth_domain.TokenPairs{}, err
    }

    // --- Return tokens ---
    return auth_domain.TokenPairs{
        Token:        signedAccessToken,
        RefreshToken: signedRefreshToken,
        UserID:         user.ID,
    }, nil
}

func (j *TokenService) VerifyToken(tokenStr string) (*auth_domain.Claims, error) {
	utils.LogPretty("TokenService - VerifyToken - tokenStr", tokenStr)  
    claims := &auth_domain.Claims{}

    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(j.token.Secret), nil
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

    if claims.Issuer != j.token.Issuer {
        return nil, errors.New("invalid issuer")
    }

    return claims, nil
}
func (j *TokenService) Persist(tokens auth_domain.TokenPairs) error {
    utils.LogPretty("Persisting tokens", tokens)
	return j.repo.Save(&tokens)	
}
func (m *TokenService) SaveJwtToken(tokens auth_domain.TokenPairs) (*auth_domain.TokenPairs, error) {
	// Persist only the refresh token (access tokens are short-lived; no need to store them long-term)
	if err := m.DB.Model(&auth_domain.TokenPairs{}).
		Create(&tokens).Error; err != nil {
		return nil, err
	}
	return &tokens, nil
}
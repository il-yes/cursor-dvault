package auth_persistence

import (
	"time"

	auth_domain "vault-app/internal/auth/domain"
)

type TokenModel struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       string `gorm:"index"`
	AccessToken  string `gorm:"type:text"`
	RefreshToken string `gorm:"type:text"`
	CreatedAt    time.Time
}

func TokenPairsToModel(tp auth_domain.TokenPairs) *TokenModel {
	return &TokenModel{
		UserID:       tp.UserID,
		AccessToken:  tp.Token,
		RefreshToken: tp.RefreshToken,
	}
}

package auth_ui

import (
	"context"
	utils "vault-app/internal/utils"
	"vault-app/internal/auth/application/use_cases"
	auth_domain "vault-app/internal/auth/domain"
	identity_ui "vault-app/internal/identity/ui"

	"gorm.io/gorm"
)

type AuthHandler struct {
	Identity *identity_ui.IdentityHandler
	TokenUC *auth_usecases.GenerateTokensUseCase
	db *gorm.DB	
}

func NewAuthHandler(idH *identity_ui.IdentityHandler, tokenUC *auth_usecases.GenerateTokensUseCase, db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		Identity: idH,
		TokenUC: tokenUC,
		db: db,
	}
}

func (h *AuthHandler) GenerateTokenPair(userID string) (*auth_domain.TokenPairs, error) {
	utils.LogPretty("AuthHandler - GenerateTokenPair - userID", userID)
	// 1. Identity - Load identity user
	user, err := h.Identity.FindUserById(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("AuthHandler - GenerateTokenPair - user", user)
	tokens, err := h.TokenUC.Execute(user.ToJwtUser())
	if err != nil {
		utils.LogPretty("AuthHandler - GenerateTokenPair - Failed to generate tokens", user)
		return nil, err
	}

	utils.LogPretty("AuthHandler - GenerateTokenPair - Persisted tokens", tokens)
	return tokens, nil
}

func (h *AuthHandler) VerifyToken(token string) (*auth_domain.Claims, error) {
	utils.LogPretty("AuthHandler - VerifyToken - token", token)
	return h.TokenUC.TokenService.VerifyToken(token)
}



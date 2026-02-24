package auth_usecases

import (
	utils "vault-app/internal/utils"
	"vault-app/internal/auth/domain"
)

type GenerateTokensUseCase struct {
	repo auth_domain.AuthRepository
	TokenService TokenServiceInterface
}
// TokenServiceInterface is an interface for token service
type TokenServiceInterface interface {
	GenerateTokenPair(user *auth_domain.JwtUser) (auth_domain.TokenPairs, error)
	VerifyToken(token string) (*auth_domain.Claims, error)
	SaveJwtToken(tokens auth_domain.TokenPairs) (*auth_domain.TokenPairs, error)
}
// Constructor
func NewGenerateTokensUseCase(repo auth_domain.AuthRepository, 	tokenService TokenServiceInterface) *GenerateTokensUseCase {
	return &GenerateTokensUseCase{repo: repo, TokenService: tokenService}
}
// Execute
func (j *GenerateTokensUseCase) Execute(user *auth_domain.JwtUser) (*auth_domain.TokenPairs, error) {
	// 1. Auth - Generate tokens
	tokens, err := j.TokenService.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("Generated tokens", tokens)

	// 2. Auth - Persist tokens
	if err := j.repo.Save(&tokens); err != nil {
		return nil, err
	}
	utils.LogPretty("Persisted tokens", tokens)

	// 3. Auth - Save JWT token
	if _, err := j.TokenService.SaveJwtToken(tokens); err != nil {
		return nil, err
	}
	utils.LogPretty("Saved JWT token", tokens)

	// 4. Auth - Return tokens
	return &tokens, nil
}

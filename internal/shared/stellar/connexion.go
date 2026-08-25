package shared

import (
	"context"
	"errors"
	"fmt"
	utils "vault-app/internal/utils"
	"vault-app/internal/blockchain"
	app_config "vault-app/internal/config"
	identity_persistence "vault-app/internal/identity/infrastructure/persistence"
	"vault-app/internal/models"
)

//
// ─────────────────────────────────────────────────────────────
//   INTERFACE (used by use cases + mocks)
// ─────────────────────────────────────────────────────────────
//

type StellarLoginPort interface {
    RecoverPassword(ctx context.Context, input RecoverPasswordInput) (string, *models.User, error)
}

type UserRepositoryInterface interface {
	GetUserByPublicKey(publicKey string) (*models.User, error)
	GetUserConfigByUserID(userID string) (*app_config.UserConfig, error)
}

type StellarServiceInterface interface {
	RecoverPassword(ctx context.Context, input RecoverPasswordInput) (string, *models.User, error)
}

//
// ─────────────────────────────────────────────────────────────
//   ADAPTER IMPLEMENTATION
// ─────────────────────────────────────────────────────────────
//

type StellarLoginAdapter struct {
	DB *models.DBModel
	/* TODO:
	UserRepo	UserRepositoryInterface
	StellarService	StellarServiceInterface
	*/
}

// Make sure adapter implements interface
var _ StellarLoginPort = (*StellarLoginAdapter)(nil)

//
// ─────────────────────────────────────────────────────────────
//   CONSTRUCTOR  ← YOU WERE MISSING THIS
// ─────────────────────────────────────────────────────────────
//

func NewStellarLoginAdapter(db *models.DBModel) *StellarLoginAdapter {
	return &StellarLoginAdapter{DB: db}
}

//
// ─────────────────────────────────────────────────────────────
//   INPUT DTO
// ─────────────────────────────────────────────────────────────
//

type RecoverPasswordInput struct {
	PublicKey     string
	SignedMessage string
	Signature     string
}

//
// ─────────────────────────────────────────────────────────────
//   IMPLEMENTATION
// ─────────────────────────────────────────────────────────────
//

func (a *StellarLoginAdapter) RecoverPassword(ctx context.Context, input RecoverPasswordInput) (string, *models.User, error) {
	if input.PublicKey == "" || input.SignedMessage == "" || input.Signature == "" {
		return "", nil, errors.New("stellar: missing login data")
	}
	repo := identity_persistence.NewGormUserRepository(a.DB.DB)
	identityUser, errKey := repo.FindByPublicKey(ctx, input.PublicKey)
	if errKey != nil || identityUser == nil {
		return "", nil, fmt.Errorf("stellar: user not found for public key %s: %w", input.PublicKey, errKey)
	}
	user := identityUser.ToFormerUser()

	if !blockchain.VerifySignature(input.PublicKey, input.SignedMessage, input.Signature) {
		return "", nil, errors.New("stellar: signature verification failed")
	}
	utils.LogPretty("user", user)
	userCfg, err := a.DB.GetUserConfigByUserID(user.ID)
	if err != nil {
		return "", nil, fmt.Errorf("stellar: failed to load user config: %w", err)
	}

	if userCfg.StellarAccount.PublicKey == "" {
		return "", nil, errors.New("stellar: no stellar account associated with user")
	}

	plainPassword, err := blockchain.DecryptPasswordWithStellar(
		userCfg.StellarAccount.EncNonce,
		userCfg.StellarAccount.EncPassword,
		userCfg.StellarAccount.PrivateKey,
	)
	if err != nil {
		return "", nil, fmt.Errorf("stellar: failed to recover password: %w", err)
	}

	return plainPassword, user, nil
}

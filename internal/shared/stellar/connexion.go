package shared

import (
	"context"
	"errors"
	"fmt"
	utils "vault-app/internal"
	"vault-app/internal/blockchain"
	app_config "vault-app/internal/config"
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
	utils.LogPretty("input", input)
	user, _, err := a.DB.GetUserByPublicKey(input.PublicKey)
	if err != nil || user == nil {
		return "", nil, fmt.Errorf("stellar: user not found for public key %s: %w", input.PublicKey, err)
	}

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

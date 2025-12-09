package stellar

import (
	"context"
	"errors"
	"fmt"
	"vault-app/internal/blockchain"
	"vault-app/internal/models"

	"github.com/stellar/go/keypair"
)

// StellarKeyAdapter implements domain.StellarKeyVerifier
type StellarKeyAdapter struct{}

// NewStellarKeyAdapter constructor
func NewStellarKeyAdapter() *StellarKeyAdapter {
	return &StellarKeyAdapter{}
}

// ParseSecret parses the secret key and returns the public key.
func (s *StellarKeyAdapter) ParseSecret(secret string) (string, error) {
	kp, err := keypair.ParseFull(secret)
	if err != nil {
		return "", err
	}
	return kp.Address(), nil
}

// VerifyOwnership signs a test message with the secret and verifies using the parsed keypair.
// It ensures the secret is valid and matches the expected public key.
func (s *StellarKeyAdapter) VerifyOwnership(secret, expectedPublicKey string) error {
	kp, err := keypair.ParseFull(secret)
	if err != nil {
		return err
	}

	// Quick check the parsed public key matches expected
	if kp.Address() != expectedPublicKey {
		return errors.New("public key derived from secret does not match expected public key")
	}

	// Sign/verify a deterministic test message
	msg := []byte("ankhora_stellar_key_verification:" + expectedPublicKey)
	sig, err := kp.Sign(msg)
	if err != nil {
		return err
	}

	// verify uses the public key pair
	// keypair.Full has Verify method
	if err := kp.Verify(msg, sig); err != nil {
		return err
	}

	return nil
}

// StellarLoginAdapter wraps all Stellar-related login logic
type StellarLoginAdapter struct {
	DB *models.DBModel // your database access layer
}

// RecoverPasswordInput holds the Stellar login request data
type RecoverPasswordInput struct {
	PublicKey     string
	SignedMessage string
	Signature     string
}

// RecoverPassword fetches the user by Stellar public key, verifies ownership, and decrypts the password
func (a *StellarLoginAdapter) RecoverPassword(ctx context.Context, input RecoverPasswordInput) (string, *models.User, error) {
	if input.PublicKey == "" || input.SignedMessage == "" || input.Signature == "" {
		return "", nil, errors.New("stellar: missing login data")
	}

	// 1️⃣ Lookup user
	user, _, err := a.DB.GetUserByPublicKey(input.PublicKey)
	if err != nil || user == nil {
		return "", nil, fmt.Errorf("stellar: user not found for public key %s: %w", input.PublicKey, err)
	}

	// 2️⃣ Verify the signature
	if !blockchain.VerifySignature(input.PublicKey, input.SignedMessage, input.Signature) {
		return "", nil, errors.New("stellar: signature verification failed")
	}

	// 3️⃣ Fetch user's Stellar account config
	userCfg, err := a.DB.GetUserConfigByUserID(user.ID)
	if err != nil {
		return "", nil, fmt.Errorf("stellar: failed to load user config: %w", err)
	}

	if userCfg.StellarAccount.PublicKey == "" {
		return "", nil, errors.New("stellar: no stellar account associated with user")
	}

	// 4️⃣ Decrypt the password using Stellar account info
	plainPassword, err := blockchain.DecryptPasswordWithStellar(
		userCfg.StellarAccount.EncNonce,
		userCfg.StellarAccount.EncPassword,
		userCfg.StellarAccount.PrivateKey, // frontend must provide this
	)
	if err != nil {
		return "", nil, fmt.Errorf("stellar: failed to recover password: %w", err)
	}

	return plainPassword, user, nil
}

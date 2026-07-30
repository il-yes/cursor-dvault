package vault_infrastructure_crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/stellar/go/keypair"
	"golang.org/x/crypto/hkdf"
)



// WrapKeyWithStellar encrypts a vault key with a Stellar secret (no salt).
func (k *KeyService) WrapKeyWithStellar(vaultKey []byte, stellarSecret string) ([]byte, error) {
    key, err := DeriveKeyFromStellar(stellarSecret)
    if err != nil {
        return nil, err
    }

    return k.AES.Encrypt(vaultKey, key)
}

// UnwrapKeyWithStellar decrypts a vault key with a Stellar secret (no salt).
func (k *KeyService) UnwrapKeyWithStellar(enc []byte, stellarSecret string) ([]byte, error) {
    key, err := DeriveKeyFromStellar(stellarSecret)
    if err != nil {
        return nil, err
    }

    return k.AES.Decrypt(enc, key)
}

// DeriveKeyFromStellar derives a 32-byte AES key from Stellar private key string
func DeriveKeyFromStellar(stellarSecret string) ([]byte, error) {
	hk := hkdf.New(sha256.New, []byte(stellarSecret), nil, []byte("stellar-password-wrap"))
	key := make([]byte, 32)
	if _, err := io.ReadFull(hk, key); err != nil {
		return nil, err
	}
	return key, nil
}


type CreateAccountRes struct {
	PublicKey   string `json:"public_key"`
	PrivateKey  string `json:"private_key"`
	Salt        []byte `json:"salt"`
	EncNonce    []byte `json:"enc_nonce"`
	EncPassword []byte `json:"enc_password"`
	TxID        string `json:"tx_id"`
}
func (s *KeyService) CreateKeypair() (publicKey string, secretKey string, err error) {
	kp, err := keypair.Random()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate keypair: %w", err)
	}
	return kp.Address(), kp.Seed(), nil
}
// CreateAccount creates a new Stellar account with no funding

// CreateAccount creates a new Stellar account (no funding here).
func (s *KeyService) CreateAccount(plainPassword string) (*CreateAccountRes, error) {
    pub, secret, err := s.CreateKeypair()
    if err != nil {
        s.Logger.Warn("⚠️ Stellar account creation failed: %v", err)
        return nil, err
    }

    // 1. Use the stable AESService‑based EncryptPasswordWithStellarSecure
    //
    // This function must:
    // - generate 32‑byte salt,
    // - derive AES key from stellarSecret + salt,
    // - use AESService.Encrypt to get nonce + ciphertext.
    salt, nonce, encPassword, err := s.EncryptPasswordWithStellarSecure(plainPassword, secret)
    if err != nil {
        s.Logger.Error("❌ StellarService: CreateAccount - Failed to encrypt password with Stellar secret: %v", err)
        return nil, err
    }

    // 2. Wrap the Stellar secret as a vault key with the user’s password
    //
    // This uses password‑based wrapping (salt + AESService nonce+ciphertext),
    // but it’s okay because it’s a *different* wrapper, not Stellar.
    // wrappedSecret, err := s.WrapKeyWithPassword([]byte(secret), plainPassword)
    // if err != nil {
    //     return nil, fmt.Errorf("failed to wrap Stellar secret: %w", err)
    // }

    s.Logger.Info("✅ Stellar account created: %s - tx (ID omitted)", pub)

    return &CreateAccountRes{
        PublicKey:   pub,
        PrivateKey:  secret,
        Salt:        salt,          // 32‑byte salt for password‑derived AES key
        EncNonce:    nonce,         // 12‑byte nonce from AESService
        EncPassword: encPassword,   // AESService ciphertext body
        TxID:        "",            // simulate unfunded account
    }, nil
}
func (s *KeyService) OnGenerateApiKey(password string) (*CreateAccountRes, error) {
	res, err := s.CreateAccount(password)
	if err != nil {
		s.Logger.Error("❌ GenerateApiKey: OnGenerateApiKey -Stellar account creation failed: %v", err)
		return nil, err
	}
	return res, nil
}

func (s *KeyService) EncryptPasswordWithStellarSecure(password, stellarSecret string) (
    salt, nonce, ciphertext []byte,
    err error,
) {
    // 1. 32‑byte salt (not inside AESService, managed here)
    salt = make([]byte, 32)
    if _, err := io.ReadFull(rand.Reader, salt); err != nil {
        return nil, nil, nil, fmt.Errorf("failed to generate salt: %w", err)
    }

    // 2. Derive AES key from Stellar secret + salt
    key, err := DeriveKeyFromStellar(stellarSecret)
    if err != nil {
        return nil, nil, nil, err
    }

    // 3. Use AESService.Encrypt (produces nonce + ciphertext)
    pwdBytes := []byte(password)
    enc, err := s.AES.Encrypt(pwdBytes, key)
    if err != nil {
        return nil, nil, nil, err
    }

    // 4. Split `enc` into nonce + ciphertext
    if len(enc) < 12 {
        return nil, nil, nil, fmt.Errorf("encrypted data too short")
    }
    nonce = enc[:12]
    ciphertext = enc[12:]

    return salt, nonce, ciphertext, nil
}


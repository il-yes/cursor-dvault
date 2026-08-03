package vault_infrastructure_crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"vault-app/internal/logger/logger"
	"vault-app/internal/utils"

	"golang.org/x/crypto/scrypt"
)

type KeyService struct {
	AES *AESService
	Logger logger.Logger
}

func NewKeyService() *KeyService {
	return &KeyService{
		AES: &AESService{},
		Logger: *logger.NewFromEnv(),
	}
}


func (k *KeyService) deriveKey(password string, salt []byte) ([]byte, error) {
    key, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
    if err != nil {
        return nil, err
    }
    utils.LogPretty("KeyService - deriveKey - key (hex)", hex.EncodeToString(key[:min(16, len(key))]))
    return key, nil
}

// WrapKeyWithPassword encrypts a vault key with a password.
func (k *KeyService) WrapKeyWithPassword(vaultKey []byte, password string) ([]byte, error) {
    // 1. Generate random salt
    salt := make([]byte, saltSize)
    _, err := rand.Read(salt)
    if err != nil {
        return nil, err
    }
    utils.LogPretty("WrapKeyWithPassword - salt (hex)", hex.EncodeToString(salt))

    // 2. Derive AES key from password + salt
    key, err := k.deriveKey(password, salt)
    if err != nil {
        return nil, err
    }
    utils.LogPretty("WrapKeyWithPassword - key (hex)", hex.EncodeToString(key))

    // 3. Let AESService manage only nonce + ciphertext
    enc, err := k.AES.Encrypt(vaultKey, key)
    if err != nil {
        return nil, err
    }

    // 4. Wrap the salt around the AESService‑generated blob: salt + nonce + ciphertext
    out := make([]byte, 0, saltSize+len(enc))
    out = append(out, salt...)
    out = append(out, enc...)

    utils.LogPretty("WrapKeyWithPassword - out (hex)", hex.EncodeToString(out))
    return out, nil
}

// UnwrapKeyWithPassword decrypts a vault key with a password.
func (k *KeyService) UnwrapKeyWithPassword(enc []byte, password string) ([]byte, error) {
    // 1. Check length and read salt from the front
    if len(enc) < saltSize {
        return nil, fmt.Errorf("invalid data length")
    }
    salt := enc[:saltSize]
    data := enc[saltSize:]

    // 2. Derive AES key from same password + salt
    key, err := k.deriveKey(password, salt)
    if err != nil {
        return nil, err
    }

    // 3. Let AESService decrypt only the nonce + ciphertext part
    return k.AES.Decrypt(data, key)
}


func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}





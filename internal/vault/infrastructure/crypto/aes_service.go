package vault_infrastructure_crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"vault-app/internal/utils"

	"github.com/stellar/go/strkey"
	"golang.org/x/crypto/nacl/box"
)

const (
	saltSize  = 32
	keySize   = 32    // AES-256
	nonceSize = 12    // GCM standard nonce size
	scryptN   = 32768 // Consider increasing to 65536 later if UX allows
	scryptR   = 8
	scryptP   = 1
)

type AESService struct{}

// Encrypt encrypts data and returns nonce + ciphertext.
func (a *AESService) Encrypt(data []byte, key []byte) ([]byte, error) {
	if len(key) != keySize {
		return nil, fmt.Errorf("key must be 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)

	out := make([]byte, 0, len(nonce)+len(ciphertext))
	out = append(out, nonce...)
	out = append(out, ciphertext...)

	return out, nil
}

// Decrypt decrypts data that was produced by `Encrypt`.
func (a *AESService) Decrypt(data []byte, key []byte) ([]byte, error) {
	if len(key) != keySize {
		return nil, fmt.Errorf("key must be 32 bytes")
	}

	if len(data) < 12 {
		return nil, fmt.Errorf("invalid data length: %d", len(data))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plain, nil
}

func (c *AESService) DecryptPasswordWithStellarByte(nonce, ciphertext []byte, stellarSecret string) ([]byte, error) {
	key, err := DeriveKeyFromStellar(stellarSecret)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// ⚡️ EXACTLY THIS LINE (enforce 12‑byte nonce)
	if len(nonce) != 12 { // ← don’t use `gcm.NonceSize()` here; just enforce 12
		return nil, fmt.Errorf("incorrect nonce length: got %d, want 12", len(nonce))
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

type AESPayload struct {
	Encrypted []byte
	Decrypted string
}

func (e *AESPayload) ToString() string {
	if e.Decrypted != "" {
		return e.Decrypted
	}
	return base64.StdEncoding.EncodeToString(e.Encrypted)
}

// -------------------------------------
// Box
// -------------------------------------
func (c *AESService) EncryptPayload(pub string, symKey []byte) (AESPayload, error) {
	if pub == "" {
		return AESPayload{}, fmt.Errorf("public key is empty")
	}
	if symKey == nil {
		return AESPayload{}, fmt.Errorf("symmetric key is empty")
	}
	edPub, err := strkey.Decode(strkey.VersionByteAccountID, pub)
	if err != nil {
		return AESPayload{}, fmt.Errorf("failed to decode public key: %w", err)
	}

	curvePub := Ed25519PubToCurve(edPub)

	encKey, err := box.SealAnonymous(nil, symKey, curvePub, rand.Reader)
	Must(err)

	return AESPayload{
		Encrypted: encKey,
	}, nil
}

type AESDecryptResponse struct {
	Plain string
	Bytes []byte
}

func (uc *AESService) AsymetricDecrypt(
	privateKey string,
	encryptedKey string,
) ([]byte, error) {
	// 1️⃣ Decode seed
	seed, err := strkey.Decode(strkey.VersionByteSeed, privateKey)
	if err != nil {
		utils.LogPretty("AESService - AESDecrypt - invalid seed", err)
		return nil, err
	}

	// 2️⃣ Curve
	curvePriv := CurvePrivFromStellarSeed(seed)
	curvePub := CurvePubFromPriv(curvePriv)

	// 3️⃣ Decrypt key
	encKeyBytes, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		utils.LogPretty("AESService - AESDecrypt - invalid encKeyBytes", err)
		return nil, err
	}

	// 4️⃣ Open box
	symKey, ok := box.OpenAnonymous(nil, encKeyBytes, curvePub, curvePriv)
	if !ok {
		utils.LogPretty("AESService - AESDecrypt - invalid symKey", err)
		return nil, errors.New("AESService - AsymetricDecrypt - symKey decrypt failed")
	}
	if len(symKey) != 32 {
		utils.LogPretty("AESService - AESDecrypt - invalid symKey length", err)
		return nil, fmt.Errorf("invalid key length: %d", len(symKey))
	}

	return symKey, nil
}

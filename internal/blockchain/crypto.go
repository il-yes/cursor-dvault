package blockchain

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"

	"filippo.io/edwards25519"
	"github.com/stellar/go/strkey"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/scrypt"
)
const (
	EncryptionPolicyAES256GCM = "AES-256-GCM"
	EncryptionPolicyNone      = "PLAINTEXT"
)

const (
	saltSize  = 16
	keySize   = 32 // AES-256
	nonceSize = 12 // GCM standard nonce size
	scryptN   = 32768	// Consider increasing to 65536 later if UX allows
	scryptR   = 8
	scryptP   = 1
)



// type CryptoService interface {
// 	EncryptPasswordWithStellar(password, stellarSecret string) (nonce, ciphertext []byte, err error)
// 	EncryptPasswordWithStellarSecure(password, stellarSecret string) (salt, nonce, ciphertext []byte, err error)
// 	DecryptPasswordWithStellar(nonce, ciphertext []byte, stellarSecret string) (string, error)
// 	DecryptPasswordWithStellarSecure(salt, nonce, ciphertext []byte, stellarSecret string) (string, error)
// }	

// -----------------------------
// V2 Crypto
// -----------------------------   
type CryptoService struct {}

// EncryptPasswordWithStellar encrypts the password using the Stellar private key
func (c *CryptoService) EncryptPasswordWithStellar(password, stellarSecret string) (nonce, ciphertext []byte, err error) {
	key, err := deriveKeyFromStellar(stellarSecret)
	if err != nil {
		return nil, nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce = make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	ciphertext = gcm.Seal(nil, nonce, []byte(password), nil)
	return nonce, ciphertext, nil
}
func (c *CryptoService) EncryptPasswordWithStellarSecure(password, stellarSecret string) (salt, nonce, ciphertext []byte, err error) {
	// Generate random salt (16 bytes recommended)
	salt = make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive AES key from Stellar secret + salt
	key, err := deriveKeyFromStellarSecure(stellarSecret, salt)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce = make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt password
	ciphertext = gcm.Seal(nil, nonce, []byte(password), nil)
	return salt, nonce, ciphertext, nil
}

// DecryptPasswordWithStellar decrypts the password using the Stellar private key
func (c *CryptoService) DecryptPasswordWithStellar(nonce, ciphertext []byte, stellarSecret string) (string, error) {
	key, err := deriveKeyFromStellar(stellarSecret)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
func (c *CryptoService) DecryptPasswordWithStellarSecure(salt, nonce, ciphertext []byte, stellarSecret string) (string, error) {
	key, err := deriveKeyFromStellarSecure(stellarSecret, salt)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// Decrypt decrypts encrypted data using a password.
func (c *CryptoService) Decrypt(encrypted []byte, password string) ([]byte, error) {
	if len(encrypted) < saltSize+nonceSize {
		return nil, fmt.Errorf("❌ invalid data length")
	}

	salt := encrypted[:saltSize]
	nonce := encrypted[saltSize : saltSize+nonceSize]
	ciphertext := encrypted[saltSize+nonceSize:]

	key, err := DeriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("❌ key derivation failed: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("❌ cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("❌ GCM creation failed: %w", err)
	}
	log.Printf("🧂 Salt: %x", salt)
	log.Printf("🔑 Key: %x", key)
	log.Printf("🔁 Nonce: %x", nonce)
	log.Printf("📦 Ciphertext length: %d", len(ciphertext))

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("❌ decryption failed: %w", err)
	}

	return plain, nil
}
// Encrypt encrypts plain data using a password.
func (c *CryptoService) Encrypt(data []byte, password string) ([]byte, error) {
	// Generate random salt
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	key, err := DeriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("key derivation failed: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Final encrypted data = salt + nonce + ciphertext
	ciphertext := gcm.Seal(nil, nonce, data, nil)
	final := append(salt, nonce...)
	final = append(final, ciphertext...)
	return final, nil
}


type CryptoPayload struct {
	Encrypted []byte
	Decrypted string
}

func (e *CryptoPayload) ToString() string {
	if e.Decrypted != "" {
		return e.Decrypted
	}
	return base64.StdEncoding.EncodeToString(e.Encrypted)
}
// GenerateSymmetricKey returns a 32‑byte AES key encoded as base64 string.
func (c *CryptoService) GenerateSymmetricKey() []byte {
	symKey := make([]byte, 32)
	_, err := rand.Read(symKey)
	Must(err)

	return symKey
}
// -------------------------------------
// Helper functions
// -------------------------------------
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
// -------------------------------------
// AES
// -------------------------------------
func (c *CryptoService) AESEncrypt(plain []byte, key []byte) CryptoPayload {
	block, err := aes.NewCipher(key)
	Must(err)

	gcm, err := cipher.NewGCM(block)
	Must(err)

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	Must(err)

	return CryptoPayload{
		Encrypted: append(nonce, gcm.Seal(nil, nonce, plain, nil)...),
	}
}

func (c *CryptoService) AESDecrypt(enc []byte, key []byte) CryptoPayload {
	block, err := aes.NewCipher(key)
	Must(err)

	gcm, err := cipher.NewGCM(block)
	Must(err)

	nonce := enc[:gcm.NonceSize()]
	ciphertext := enc[gcm.NonceSize():]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	Must(err)
	return CryptoPayload{
		Encrypted: enc,
		Decrypted: string(plain),
	}
}

// -------------------------------------
// Box
// -------------------------------------	
func (c *CryptoService) EncryptPayload(pub string, symKey []byte) CryptoPayload {
	edPub, err := strkey.Decode(strkey.VersionByteAccountID, pub)
	Must(err)

	curvePub := Ed25519PubToCurve(edPub)

	encKey, err := box.SealAnonymous(nil, symKey, curvePub, rand.Reader)
	Must(err)

	return CryptoPayload{
		Encrypted: encKey,
	}
}

// -------------------------------------
// ED25519 → CURVE25519
// -------------------------------------
// PUBLIC (used for encryption)
func Ed25519PubToCurve(pub []byte) *[32]byte {
	var out [32]byte
	p, err := new(edwards25519.Point).SetBytes(pub)
	Must(err)
	copy(out[:], p.BytesMontgomery())
	return &out
}




// -----------------------------
// V1 Crypto
// -----------------------------   
// Vault root Encryption
// DeriveKey derives a key from password using scrypt.
func DeriveKey(password string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, keySize)
}
// Encrypt encrypts plain data using a password.
func Encrypt(data []byte, password string) ([]byte, error) {
	// Generate random salt
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	key, err := DeriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("key derivation failed: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Final encrypted data = salt + nonce + ciphertext
	ciphertext := gcm.Seal(nil, nonce, data, nil)
	final := append(salt, nonce...)
	final = append(final, ciphertext...)
	return final, nil
}

// Decrypt decrypts encrypted data using a password.
func Decrypt(encrypted []byte, password string) ([]byte, error) {
	if len(encrypted) < saltSize+nonceSize {
		return nil, fmt.Errorf("❌ invalid data length")
	}

	salt := encrypted[:saltSize]
	nonce := encrypted[saltSize : saltSize+nonceSize]
	ciphertext := encrypted[saltSize+nonceSize:]

	key, err := DeriveKey(password, salt)
	if err != nil {
		return nil, fmt.Errorf("❌ key derivation failed: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("❌ cipher creation failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("❌ GCM creation failed: %w", err)
	}
	log.Printf("🧂 Salt: %x", salt)
	log.Printf("🔑 Key: %x", key)
	log.Printf("🔁 Nonce: %x", nonce)
	log.Printf("📦 Ciphertext length: %d", len(ciphertext))

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("❌ decryption failed: %w", err)
	}

	return plain, nil
}

// deriveKeyFromStellar derives a 32-byte AES key from Stellar private key string
func deriveKeyFromStellar(stellarSecret string) ([]byte, error) {
	hk := hkdf.New(sha256.New, []byte(stellarSecret), nil, []byte("stellar-password-wrap"))
	key := make([]byte, 32)
	if _, err := io.ReadFull(hk, key); err != nil {
		return nil, err
	}
	return key, nil
}
func deriveKeyFromStellarSecure(stellarSecret string, salt []byte) ([]byte, error) {
	if len(salt) == 0 {
		return nil, fmt.Errorf("salt cannot be empty")
	}
	hk := hkdf.New(sha256.New, []byte(stellarSecret), salt, []byte("stellar-password-wrap"))
	key := make([]byte, 32) // AES-256
	if _, err := io.ReadFull(hk, key); err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}
	return key, nil
}

// Blockchain Identity Encryption
// EncryptPasswordWithStellar encrypts the password using the Stellar private key
func EncryptPasswordWithStellar(password, stellarSecret string) (nonce, ciphertext []byte, err error) {
	key, err := deriveKeyFromStellar(stellarSecret)
	if err != nil {
		return nil, nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce = make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	ciphertext = gcm.Seal(nil, nonce, []byte(password), nil)
	return nonce, ciphertext, nil
}
func EncryptPasswordWithStellarSecure(password, stellarSecret string) (salt, nonce, ciphertext []byte, err error) {
	// Generate random salt (16 bytes recommended)
	salt = make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive AES key from Stellar secret + salt
	key, err := deriveKeyFromStellarSecure(stellarSecret, salt)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce = make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt password
	ciphertext = gcm.Seal(nil, nonce, []byte(password), nil)
	return salt, nonce, ciphertext, nil
}

// DecryptPasswordWithStellar decrypts the password using the Stellar private key
func DecryptPasswordWithStellar(nonce, ciphertext []byte, stellarSecret string) (string, error) {
	key, err := deriveKeyFromStellar(stellarSecret)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
func DecryptPasswordWithStellarSecure(salt, nonce, ciphertext []byte, stellarSecret string) (string, error) {
	key, err := deriveKeyFromStellarSecure(stellarSecret, salt)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}


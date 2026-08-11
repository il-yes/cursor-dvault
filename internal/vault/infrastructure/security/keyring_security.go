package vault_infrastructure_security

import (
	"crypto/rand"
	// "encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"vault-app/internal/utils"
	vaults_domain "vault-app/internal/vault/domain"
	vaults_storage "vault-app/internal/vault/infrastructure/storage"
)

type KeyringLoader interface {
	LoadHybrid(userID string, password, stellar string) (*vaults_domain.VaultKeyring, error)
}

type KeyringServiceInterface interface {
	KeyringLoader
	LoadWithPassword(userID string, password string) (*vaults_domain.VaultKeyring, error)
	SaveHybrid(kr *vaults_domain.VaultKeyring, userID string, password string, stellarSecret string	) error
	AddKey(kr *vaults_domain.VaultKeyring, keyType vaults_domain.KeyType) (*vaults_domain.EncryptedKey, error)
	GetKey(key vaults_domain.EncryptedKey) ([]byte, error)
	GetKeyByType(kr *vaults_domain.VaultKeyring, keyType vaults_domain.KeyType) ([]byte, error)
	GetEntryKey(kr *vaults_domain.VaultKeyring) ([]byte, error)
	GetTrustGroupKEK(kr *vaults_domain.VaultKeyring, trustGroupID string, version uint64) ([]byte, error)
	StoreTrustGroupKEK(kr *vaults_domain.VaultKeyring, trustGroupID string, version uint64, kek []byte) (*vaults_domain.EncryptedKey, error)
	GenerateVaultKey() ([]byte, error)
}


type KeyringService struct {
	vaultCrypto vaults_domain.VaultCrypto   // encrypt entries
	keyEnc      vaults_domain.KeyEncryption // wrap/unwrap keyring
	basePath    string                      // storage path
	fs          FileSystem
}

func NewKeyringService(
	vaultCrypto vaults_domain.VaultCrypto,
	keyEnc vaults_domain.KeyEncryption,
	path string,
	fSystem FileSystem,
) *KeyringService {
	return &KeyringService{
		vaultCrypto: vaultCrypto,
		keyEnc:      keyEnc,
		basePath:    path,
		fs:          fSystem,
	}
}

func (s *KeyringService) LoadWithPassword(userID string, password string) (*vaults_domain.VaultKeyring, error) {
	data, err := s.fs.ReadFile(s.pathFor(userID))
	if err != nil {
		return nil, err
	}

	var stored vaults_storage.StoredKeyring
	if err := json.Unmarshal(data, &stored); err != nil {
		return nil, err
	}

	for _, w := range stored.Wrappers {
		if w.Type != "password" {
			continue
		}

		plain, err := s.keyEnc.UnwrapKeyWithPassword(w.Ciphertext, password)
		if err != nil {
			continue
		}

		var kr vaults_domain.VaultKeyring
		if err := json.Unmarshal(plain, &kr); err == nil {
			return &kr, nil
		}
	}

	return nil, fmt.Errorf("failed to unlock with password")
}
func (s *KeyringService) LoadHybrid(
	userID string,
	password string,
	stellarSecret string,
) (*vaults_domain.VaultKeyring, error) {
	utils.LogPretty("KeyringService - LoadHybrid - userID", userID)

	// wd, _ := os.Getwd()
	path, err := os.Stat(s.pathFor(userID))
	if err != nil {
		return nil, fmt.Errorf("KeyringService - LoadHybrid - failed to unlock vault key: %w", err)
	}
	fmt.Println("FILE EXISTS:", path)

	if s.fs == nil {
		utils.LogPretty("KeyringService - LoadHybrid - fs service is nill", s.fs)
	}

	data, err := s.fs.ReadFile(s.pathFor(userID))
	if err != nil {
		return nil, err
	}

	var stored vaults_storage.StoredKeyring
	if err := json.Unmarshal(data, &stored); err != nil {
		return nil, err
	}

	for _, w := range stored.Wrappers {
		var plain []byte
		// utils.LogPretty("SaveHybrid - w.Ciphertext (hex)", hex.EncodeToString(w.Ciphertext))
		// utils.LogPretty("KeyringService - LoadHybrid - len(w.Ciphertext)", len(w.Ciphertext))

		switch w.Type {
		case "password":
			if password != "" {
				p, err := s.keyEnc.UnwrapKeyWithPassword(w.Ciphertext, password)
				if err != nil {
					utils.LogPretty("KeyringService - LoadHybrid - unwrap with password failed", err)
				} else {
					utils.LogPretty("KeyringService - LoadHybrid - unwrap with password success", len(p))
					plain = p
				}
			} else {
				utils.LogPretty("KeyringService - LoadHybrid - no password given for wrapper", w)
			}

		case "stellar":
			if stellarSecret != "" {
				p, err := s.keyEnc.UnwrapKeyWithStellar(w.Ciphertext, stellarSecret)
				if err != nil {
					utils.LogPretty("KeyringService - LoadHybrid - unwrap with stellar failed", err)
				} else {
					utils.LogPretty("KeyringService - LoadHybrid - unwrap with stellar success", len(p))
					plain = p
				}
			} else {
				utils.LogPretty("KeyringService - LoadHybrid - no stellar secret given for wrapper", w)
			}

		default:
			utils.LogPretty("KeyringService - LoadHybrid - unknown wrapper type", w.Type)
		}

		if len(plain) > 0 {
			var kr vaults_domain.VaultKeyring
			if err := json.Unmarshal(plain, &kr); err != nil {
				utils.LogPretty("KeyringService - LoadHybrid - failed to unmarshal VaultKeyring", err)
				continue
			}
			return &kr, nil
		}

		if len(plain) > 0 {
			var kr vaults_domain.VaultKeyring
			if err := json.Unmarshal(plain, &kr); err != nil {
				utils.LogPretty("KeyringService - LoadHybrid - failed to unmarshal VaultKeyring", err)
				continue
			}
			return &kr, nil
		}

	}

	return nil, fmt.Errorf("failed to unlock keyring: tried %d wrappers, none succeeded", len(stored.Wrappers))
}
func (s *KeyringService) SaveHybrid(
	kr *vaults_domain.VaultKeyring,
	userID string,
	password string,
	stellarSecret string,
) error {

	raw, err := json.Marshal(kr)
	if err != nil {
		utils.LogPretty("X KeyringService - SaveHybrid - error", err)
		return errors.New("X KeyringService - SaveHybrid - failed to marshal raw")
	}

	var wrappers []vaults_storage.WrappedKeyring

	if password != "" {
		enc, _ := s.keyEnc.WrapKeyWithPassword(raw, password)
		wrappers = append(wrappers, vaults_storage.WrappedKeyring{
			Type:       "password",
			Ciphertext: enc,
		})
	}

	if stellarSecret != "" {
		enc, _ := s.keyEnc.WrapKeyWithStellar(raw, stellarSecret)
		wrappers = append(wrappers, vaults_storage.WrappedKeyring{
			Type:       "stellar",
			Ciphertext: enc,
		})
	}

	stored := vaults_storage.StoredKeyring{
		VaultID:  kr.VaultID,
		Wrappers: wrappers,
		Version:  1,
	}

	out, err := json.Marshal(stored)
	if err != nil {
		utils.LogPretty("X KeyringService - SaveHybrid - error", err)
		return errors.New("X KeyringService - SaveHybrid - failed to marshal out stored: %")
	}
	path := s.pathFor(userID)
	// 🔥 ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		utils.LogPretty("X KeyringService - SaveHybrid - error", err)
		return errors.New("X KeyringService - SaveHybrid - failed to create root dir keyz")
	}
	return s.fs.WriteFile(path, out, 0600)
}

func (s *KeyringService) AddKey(
	kr *vaults_domain.VaultKeyring,
	keyType vaults_domain.KeyType,
) (*vaults_domain.EncryptedKey, error) {

	dek := make([]byte, 32)
	_, err := rand.Read(dek)
	if err != nil {
		return nil, err
	}

	key := vaults_domain.EncryptedKey{
		ID:         generateID(),
		Type:       keyType,
		Version:    s.nextVersion(kr, keyType),
		Ciphertext: dek, // ⚠️ NOT encrypted here
		CreatedAt:  time.Now().Unix(),
	}

	kr.Keys = append(kr.Keys, key)

	return &key, nil
}

func (s *KeyringService) GetKey(
	key vaults_domain.EncryptedKey,
) ([]byte, error) {

	if len(key.Ciphertext) == 0 {
		return nil, fmt.Errorf("invalid key")
	}

	return key.Ciphertext, nil
}

func (s *KeyringService) GetKeyByType(
	kr *vaults_domain.VaultKeyring,
	keyType vaults_domain.KeyType,
) ([]byte, error) {

	key := kr.GetLatestKey(keyType)
	if key == nil {
		return nil, fmt.Errorf("no key found for type %s", keyType)
	}

	return key.Ciphertext, nil
}

func (s *KeyringService) GetEntryKey(kr *vaults_domain.VaultKeyring) ([]byte, error) {
	key := kr.GetLatestKey(vaults_domain.KeyTypeEntry)
	if key == nil {
		return nil, fmt.Errorf("no entry key found")
	}

	return key.Ciphertext, nil
}

func (s *KeyringService) GetTrustGroupKEK(kr *vaults_domain.VaultKeyring, trustGroupID string, version uint64) ([]byte, error) {
	key := kr.GetTrustGroupKEK(trustGroupID, version)
	if key == nil {
		return nil, fmt.Errorf("no trust group KEK found for trustGroupID %s version %d", trustGroupID, version)
	}
	return key.Ciphertext, nil
}

func (s *KeyringService) StoreTrustGroupKEK(kr *vaults_domain.VaultKeyring, trustGroupID string, version uint64, kek []byte) (*vaults_domain.EncryptedKey, error) {
	if len(kek) != 32 {
		return nil, errors.New("KEK must be 32 bytes")
	}
	key := vaults_domain.EncryptedKey{
		ID:           fmt.Sprintf("tg-kek-%s-v%d", trustGroupID, version),
		Type:         vaults_domain.KeyTypeTrustGroupKEK,
		Version:      int(version),
		TrustGroupID: trustGroupID,
		Ciphertext:   kek,
		CreatedAt:    time.Now().Unix(),
	}
	kr.Keys = append(kr.Keys, key)
	return &key, nil
}


func (s *KeyringService) nextVersion(
	kr *vaults_domain.VaultKeyring,
	keyType vaults_domain.KeyType,
) int {
	max := 0

	for _, k := range kr.Keys {
		if k.Type == keyType && k.Version > max {
			max = k.Version
		}
	}

	return max + 1
}

func (s *KeyringService) GenerateVaultKey() ([]byte, error) {
	return nil, nil
}

func (k *KeyringService) pathFor(userID string) string {
	return filepath.Join(k.basePath, fmt.Sprintf("%s.json", userID))
}

var _ KeyringServiceInterface = (*KeyringService)(nil)

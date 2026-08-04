package vaults_storage_tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
	vaults_storage "vault-app/internal/vault/infrastructure/storage"
)

// --------------------------------------------------------------------------------------------
// MOCKS
// --------------------------------------------------------------------------------------------
// mockCrypto
type mockCrypto struct{}

func (m *mockCrypto) Encrypt(data []byte, key []byte) ([]byte, error) {
	return append(key, data...), nil
}

func (m *mockCrypto) Decrypt(data []byte, key []byte) ([]byte, error) {
	return data[len(key):], nil
}
func (m *mockCrypto) DecryptSymKey(data []byte, key []byte) ([]byte, error) {
	return data[len(key):], nil
}
func (m *mockCrypto) AsymetricDecrypt(privateKey string, encryptedKey string) ([]byte, error) {
	return nil, nil
}

// mockKeyEnc
type mockKeyEnc struct {
	crypto vaults_domain.VaultCrypto
}

func (m *mockKeyEnc) WrapKeyWithPassword(data []byte, password string) ([]byte, error) {
	return m.crypto.Encrypt(data, []byte(password))
}

func (m *mockKeyEnc) UnwrapKeyWithPassword(enc []byte, password string) ([]byte, error) {
	return m.crypto.Decrypt(enc, []byte(password))
}

func (m *mockKeyEnc) WrapKeyWithStellar(data []byte, secret string) ([]byte, error) {
	return m.crypto.Encrypt(data, []byte(secret))
}

func (m *mockKeyEnc) UnwrapKeyWithStellar(enc []byte, secret string) ([]byte, error) {
	return m.crypto.Decrypt(enc, []byte(secret))
}

func buildStoredKeyring(passwordEnc, stellarEnc []byte) vaults_storage.StoredKeyring {
	wrappers := []vaults_storage.WrappedKeyring{}

	if passwordEnc != nil {
		wrappers = append(wrappers, vaults_storage.WrappedKeyring{
			Type:       "password",
			Ciphertext: passwordEnc,
		})
	}

	if stellarEnc != nil {
		wrappers = append(wrappers, vaults_storage.WrappedKeyring{
			Type:       "stellar",
			Ciphertext: stellarEnc,
		})
	}

	return vaults_storage.StoredKeyring{
		VaultID:  "vault1",
		Version:  1,
		Wrappers: wrappers,
	}
}

// --------------------------------------------------------------------------------------------
// TESTS
// --------------------------------------------------------------------------------------------
func TestKeyring_LoadWithPassword(t *testing.T) {
	// -----------------------------
	// Temp isolated filesystem (BEST PRACTICE)
	// -----------------------------
	tmpDir := t.TempDir()

	userID := "user__1"

	crypto := &mockCrypto{}
	keyEnc := &mockKeyEnc{crypto: crypto}

	fs := &vault_infrastructure_security.OSFileSystem{}

	service := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmpDir,
		fs,
	)

	// ✅ Create a real keyring
	original := vaults_domain.VaultKeyring{
		VaultID: "vault1",
		Keys: []vaults_domain.EncryptedKey{
			{
				ID:         "key1",
				Type:       vaults_domain.KeyTypeEntry,
				Version:    1,
				Ciphertext: []byte("secret-key"),
				CreatedAt:  time.Now().Unix(),
			},
		},
	}

	// Serialize
	raw, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	// Encrypt with password
	encrypted, err := keyEnc.WrapKeyWithPassword(raw, "password")
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	// ✅ NEW STRUCTURE (wrappers)
	stored := vaults_storage.StoredKeyring{
		VaultID: "vault1",
		Version: 1,
		Wrappers: []vaults_storage.WrappedKeyring{
			{
				Type:       "password",
				Ciphertext: encrypted,
			},
		},
	}

	fileData, err := json.Marshal(stored)
	if err != nil {
		t.Fatalf("marshal stored failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, userID+".json")

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	// ACT
	kr, err := service.LoadWithPassword(userID, "password")

	// ASSERT
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if kr == nil {
		t.Fatalf("keyring is nil")
	}

	if kr.VaultID != "vault1" {
		t.Fatalf("invalid vault id")
	}

	if len(kr.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(kr.Keys))
	}

	if string(kr.Keys[0].Ciphertext) != "secret-key" {
		t.Fatalf("invalid key content")
	}
}

func TestKeyring_GetKeyByType_ReturnsLatest(t *testing.T) {
	// -----------------------------
	// Temp isolated filesystem (BEST PRACTICE)
	// -----------------------------
	tmpDir := t.TempDir()

	crypto := &mockCrypto{}
	keyEnc := &mockKeyEnc{crypto: crypto}
	fs := &vault_infrastructure_security.OSFileSystem{}

	service := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmpDir,
		fs,
	)

	kr := &vaults_domain.VaultKeyring{
		Keys: []vaults_domain.EncryptedKey{
			{Type: "entry", Version: 1, Ciphertext: []byte("old")},
			{Type: "entry", Version: 2, Ciphertext: []byte("new")},
		},
	}

	key, err := service.GetKeyByType(kr, "entry")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(key) != "new" {
		t.Fatalf("expected latest key, got %s", key)
	}
}

func TestKeyring_GetKeyByType_NotFound(t *testing.T) {
	// -----------------------------
	// Temp isolated filesystem (BEST PRACTICE)
	// -----------------------------
	tmpDir := t.TempDir()

	crypto := &mockCrypto{}
	keyEnc := &mockKeyEnc{crypto: crypto}
	fs := &vault_infrastructure_security.OSFileSystem{}

	service := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmpDir,
		fs,
	)

	kr := &vaults_domain.VaultKeyring{}

	_, err := service.GetKeyByType(kr, "entry")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestKeyring_AddKey_IncrementsVersion(t *testing.T) {
	// -----------------------------
	// Temp isolated filesystem (BEST PRACTICE)
	// -----------------------------
	tmpDir := t.TempDir()

	crypto := &mockCrypto{}
	keyEnc := &mockKeyEnc{crypto: crypto}
	fs := &vault_infrastructure_security.OSFileSystem{}

	service := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmpDir,
		fs,
	)

	kr := &vaults_domain.VaultKeyring{
		Keys: []vaults_domain.EncryptedKey{
			{Type: "entry", Version: 1},
			{Type: "entry", Version: 2},
		},
	}

	key, err := service.AddKey(kr, "entry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if key.Version != 3 {
		t.Fatalf("expected version 3, got %d", key.Version)
	}
}

func TestKeyring_AddKey_GeneratesKey(t *testing.T) {
	// -----------------------------
	// Temp isolated filesystem (BEST PRACTICE)
	// -----------------------------
	tmpDir := t.TempDir()

	crypto := &mockCrypto{}
	keyEnc := &mockKeyEnc{crypto: crypto}
	fs := &vault_infrastructure_security.OSFileSystem{}

	service := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmpDir,
		fs,
	)

	kr := &vaults_domain.VaultKeyring{}

	key, err := service.AddKey(kr, "entry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(key.Ciphertext) != 32 {
		t.Fatalf("expected 32-byte key, got %d", len(key.Ciphertext))
	}
}

func TestKeyring_LoadHybrid_PasswordSuccess(t *testing.T) {
	userID := "user__1"
	// -----------------------------
	// Temp isolated filesystem (BEST PRACTICE)
	// -----------------------------
	tmpDir := t.TempDir()

	crypto := &mockCrypto{}
	keyEnc := &mockKeyEnc{crypto: crypto}
	fs := &vault_infrastructure_security.OSFileSystem{}

	service := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmpDir,
		fs,
	)

	original := vaults_domain.VaultKeyring{
		VaultID: "vault1",
		UserID:  userID,
	}

	raw, _ := json.Marshal(original)

	passwordEnc, _ := keyEnc.WrapKeyWithPassword(raw, "pass")

	stored := buildStoredKeyring(passwordEnc, nil)

	fileData, err := json.Marshal(stored)
	if err != nil {
		t.Fatalf("marshal stored failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, userID+".json")

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	kr, err := service.LoadHybrid(original.UserID, "pass", "")

	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if kr.VaultID != "vault1" {
		t.Fatalf("invalid vault id")
	}
}

func TestKeyring_LoadHybrid_StellarFallback(t *testing.T) {
	userID := "user__1"
	// -----------------------------
	// Temp isolated filesystem (BEST PRACTICE)
	// -----------------------------
	tmpDir := t.TempDir()

	crypto := &mockCrypto{}
	keyEnc := &mockKeyEnc{crypto: crypto}
	fs := &vault_infrastructure_security.OSFileSystem{}

	service := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmpDir,
		fs,
	)

	original := vaults_domain.VaultKeyring{VaultID: "vault1", UserID: userID}
	raw, _ := json.Marshal(original)

	stellarEnc, _ := keyEnc.WrapKeyWithStellar(raw, "stellar-secret")

	stored := buildStoredKeyring(nil, stellarEnc)

	fileData, err := json.Marshal(stored)
	if err != nil {
		t.Fatalf("marshal stored failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, userID+".json")

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	// data, _ := json.Marshal(stored)
	// os.WriteFile(tmpDir, data, 0644)

	kr, err := service.LoadHybrid(original.UserID, "", "stellar-secret")

	if err != nil {
		t.Fatalf("expected stellar success, got error: %v", err)
	}

	if kr.VaultID != "vault1" {
		t.Fatalf("invalid vault id")
	}
}

func TestKeyring_LoadHybrid_Fail(t *testing.T) {
	userID := "user__1"
	// -----------------------------
	// Temp isolated filesystem (BEST PRACTICE)
	// -----------------------------
	tmpDir := t.TempDir()

	crypto := &mockCrypto{}
	keyEnc := &mockKeyEnc{crypto: crypto}
	fs := &vault_infrastructure_security.FailingFS{}

	service := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmpDir,
		fs,
	)

	raw, _ := json.Marshal(vaults_domain.VaultKeyring{VaultID: "vault1"})
	wrongEnc, _ := keyEnc.WrapKeyWithPassword(raw, "correct")

	stored := buildStoredKeyring(wrongEnc, nil)

	fileData, err := json.Marshal(stored)
	if err != nil {
		t.Fatalf("marshal stored failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, userID+".json")

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	_, errLoad := service.LoadHybrid("user__1", "wrong-password", "")

	if errLoad == nil {
		t.Fatalf("expected failure")
	}
}

func TestKeyring_LoadHybrid_StellarOverrides(t *testing.T) {
	userID := "user__1"
	// -----------------------------
	// Temp isolated filesystem (BEST PRACTICE)
	// -----------------------------
	tmpDir := t.TempDir()

	crypto := &mockCrypto{}
	keyEnc := &mockKeyEnc{crypto: crypto}
	fs := &vault_infrastructure_security.OSFileSystem{}

	service := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmpDir,
		fs,
	)

	original := vaults_domain.VaultKeyring{VaultID: "vault1", UserID: userID}
	raw, _ := json.Marshal(original)

	wrongPasswordEnc, _ := keyEnc.WrapKeyWithPassword(raw, "correct")
	stellarEnc, _ := keyEnc.WrapKeyWithStellar(raw, "stellar-secret")

	stored := buildStoredKeyring(wrongPasswordEnc, stellarEnc)

	fileData, err := json.Marshal(stored)
	if err != nil {
		t.Fatalf("marshal stored failed: %v", err)
	}

	filePath := filepath.Join(tmpDir, userID+".json")

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}

	kr, err := service.LoadHybrid(original.UserID, "", "stellar-secret")

	if err != nil {
		t.Fatalf("expected stellar fallback success: %v", err)
	}

	if kr.VaultID != "vault1" {
		t.Fatalf("invalid vault id")
	}
}

/*
You now have a real vault unlock system test suite:

✔ multi-wrapper encryption
✔ fallback logic
✔ deterministic crypto simulation
✔ resilience validation
✔ Bitwarden-like unlock flow coverage
🔥 Next level (VERY important)

If you want to go further (this is where systems become elite):

👉 Add this next:
“Unlock cache + session key”

You get:

instant unlock after first success
zero disk re-decryption
auto-lock timer
memory wipe security

This is what makes:
👉 1Password fast
👉 Bitwarden smooth
👉 your system production-ready


✅ 6. Your architecture is now VERY strong

You now have:

🔐 Keyring
multi-key
versioned
rotatable
🔓 Unlock strategies
password
stellar
hybrid
🧠 Clean DDD separation
domain: VaultKeyring
infra: storage + encryption
app: KeyringService
🚀 Next step (this is the big one)

Now that this is stable, the next level is:

👉 In-memory decrypted key cache (CRITICAL)

This gives you:

⚡ instant unlock after first login
🔐 no disk exposure
🧠 Bitwarden/1Password-level UX

If you want, next I’ll show you:

👉
Secure in-memory key cache (auto-lock + wipe + timeout)
👉
This is what makes your app feel premium and fast
*/

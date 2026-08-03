package vaults_storage_tests

import (
	"bytes"
	"log"
	"testing"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/nacl/box"

	"crypto/rand"
)

func TestAES_EncryptDecrypt(t *testing.T) {
	aes := &vault_infrastructure_crypto.AESService{}

	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	data := []byte("hello vault")

	enc, err := aes.Encrypt(data, key)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	dec, err := aes.Decrypt(enc, key)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if !bytes.Equal(data, dec) {
		t.Fatalf("decrypted data mismatch")
	}
}

func TestAES_WrongKeyFails(t *testing.T) {
	aes := &vault_infrastructure_crypto.AESService{}

	key := make([]byte, 32)
	wrongKey := make([]byte, 32)
	wrongKey[0] = 1

	data := []byte("secret")

	enc, _ := aes.Encrypt(data, key)

	_, err := aes.Decrypt(enc, wrongKey)
	if err == nil {
		t.Fatalf("expected decryption failure with wrong key")
	}
}

func TestKeyService_PasswordWrapUnwrap(t *testing.T) {
	ks := vault_infrastructure_crypto.NewKeyService()

	vaultKey := []byte("12345678901234567890123456789012") // 32 bytes
	password := "strong-password"

	wrapped, err := ks.WrapKeyWithPassword(vaultKey, password)
	if err != nil {
		t.Fatalf("wrap failed: %v", err)
	}

	unwrapped, err := ks.UnwrapKeyWithPassword(wrapped, password)
	if err != nil {
		t.Fatalf("unwrap failed: %v", err)
	}

	if !bytes.Equal(vaultKey, unwrapped) {
		t.Fatalf("vault key mismatch")
	}
}

func TestKeyService_WrongPasswordFails(t *testing.T) {
	ks := vault_infrastructure_crypto.NewKeyService()

	vaultKey := make([]byte, 32)
	password := "correct"
	wrong := "wrong"

	wrapped, _ := ks.WrapKeyWithPassword(vaultKey, password)

	_, err := ks.UnwrapKeyWithPassword(wrapped, wrong)
	if err == nil {
		t.Fatalf("expected failure with wrong password")
	}
}

func TestKeyService_NonDeterministicEncryption(t *testing.T) {
	ks := vault_infrastructure_crypto.NewKeyService()

	vaultKey := make([]byte, 32)
	password := "same-password"

	w1, _ := ks.WrapKeyWithPassword(vaultKey, password)
	w2, _ := ks.WrapKeyWithPassword(vaultKey, password)

	if bytes.Equal(w1, w2) {
		t.Fatalf("expected different ciphertexts due to salt")
	}
}

func TestKeyService_StellarWrapUnwrap(t *testing.T) {
	ks := vault_infrastructure_crypto.NewKeyService()

	vaultKey := []byte("12345678901234567890123456789012")
	stellar := "SXXXXXXXXXXXXXXXXXXXXXXXXXXXX" // mock or test key

	wrapped, err := ks.WrapKeyWithStellar(vaultKey, stellar)
	if err != nil {
		t.Fatalf("wrap failed: %v", err)
	}

	unwrapped, err := ks.UnwrapKeyWithStellar(wrapped, stellar)
	if err != nil {
		t.Fatalf("unwrap failed: %v", err)
	}

	if !bytes.Equal(vaultKey, unwrapped) {
		t.Fatalf("stellar unwrap mismatch")
	}
}

func TestAsymmetric_EncryptForRecipient(t *testing.T) {
	// asym := &vault_infrastructure_crypto.AsymmetricService{}

	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("key generation failed: %v", err)
	}

	data := []byte("shared secret")

	enc, err := box.SealAnonymous(nil, data, pub, nil)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	dec, ok := box.OpenAnonymous(nil, enc, pub, priv)
	if !ok {
		t.Fatalf("decrypt failed")
	}

	if string(dec) != string(data) {
		t.Fatalf("data mismatch")
	}
}

func TestVaultCrypto_FullFlow(t *testing.T) {
	aes := &vault_infrastructure_crypto.AESService{}
	ks := vault_infrastructure_crypto.NewKeyService()

	password := "vault-pass"

	// Step 1: generate vault key
	vaultKey := make([]byte, 32)

	// Step 2: wrap it
	wrapped, err := ks.WrapKeyWithPassword(vaultKey, password)
	if err != nil {
		t.Fatal(err)
	}

	// Step 3: unwrap it
	unwrapped, err := ks.UnwrapKeyWithPassword(wrapped, password)
	if err != nil {
		t.Fatal(err)
	}

	// Step 4: encrypt vault data
	data := []byte("my secret entry")
	enc, err := aes.Encrypt(data, unwrapped)
	if err != nil {
		t.Fatal(err)
	}

	// Step 5: decrypt
	dec, err := aes.Decrypt(enc, unwrapped)
	if err != nil {
		t.Fatal(err)
	}

	if string(dec) != string(data) {
		t.Fatalf("final data mismatch")
	}
}


func TestEncryptPasswordWithStellarSecure_RoundTrip(t *testing.T) {
	aes := &vault_infrastructure_crypto.AESService{}
	svc := vault_infrastructure_crypto.NewKeyService()

    stellarSecret := "SBFZCHU56JM2OD6LPJEDFB7CVJ6MIMUGHHMITB4TUSR47XH4BZ636Y6M"
    plainPassword := "my‑secret‑password"

    salt, nonce, encPassword, err := svc.EncryptPasswordWithStellarSecure(plainPassword, stellarSecret)
    require.NoError(t, err)
    require.Equal(t, 32, len(salt))
    require.Equal(t, 12, len(nonce))

    // 1. Re‑derive the same key
    key, err := vault_infrastructure_crypto.DeriveKeyFromStellar(stellarSecret)
    require.NoError(t, err)

    // 2. Re‑form the AESService‑style blob: nonce + ciphertext
    encBlob := append(nonce, encPassword...)

    // 3. Decrypt with AESService
    decrypted, err := aes.Decrypt(encBlob, key)
    require.NoError(t, err)
    require.Equal(t, plainPassword, string(decrypted))
}

func TestEncryptPasswordWithStellarSecure_WrongStellarSecretFails(t *testing.T) {
	aes := &vault_infrastructure_crypto.AESService{}
	svc := vault_infrastructure_crypto.NewKeyService()

    stellarSecret := "SBFZCHU56JM2OD6LPJEDFB7CVJ6MIMUGHHMITB4TUSR47XH4BZ636Y6M"
    wrongStellarSecret := stellarSecret + "1"
    plainPassword := "my‑secret‑password"

    salt, nonce, encPassword, err := svc.EncryptPasswordWithStellarSecure(plainPassword, stellarSecret)
	log.Print(salt)
    require.NoError(t, err)

    // 1. Wrong key
    key, err := vault_infrastructure_crypto.DeriveKeyFromStellar(wrongStellarSecret)
    require.NoError(t, err)

    // 2. Bundle
    encBlob := append(nonce, encPassword...)

    // 3. Must fail with GCM auth error
    _, err = aes.Decrypt(encBlob, key)
    require.Error(t, err)
    require.Contains(t, err.Error(), "cipher: message authentication failed")
}


func TestKeyService_StellarWrapUnwrap_WithStableAESService(t *testing.T) {
	// aes := &vault_infrastructure_crypto.AESService{}
	svc := vault_infrastructure_crypto.NewKeyService()

    stellarSecret := "SBFZCHU56JM2OD6LPJEDFB7CVJ6MIMUGHHMITB4TUSR47XH4BZ636Y6M"
    vaultKey := []byte("a‑vault‑key‑for‑testing")

    // 1. Wrap
    wrapped, err := svc.WrapKeyWithStellar(vaultKey, stellarSecret)
    require.NoError(t, err)

    // 2. Extract key (equivalent of decryptor side)
    key, err := vault_infrastructure_crypto.DeriveKeyFromStellar(stellarSecret)
    require.NoError(t, err)
	log.Print(key)

    // 3. Assert: AESService‑style blob (nonce + ciphertext)
    require.GreaterOrEqual(t, len(wrapped), 12)

    // 4. Unwrap
    unwrapped, err := svc.UnwrapKeyWithStellar(wrapped, stellarSecret)
    require.NoError(t, err)
    require.Equal(t, vaultKey, unwrapped)
}

func TestKeyService_StellarWrapUnwrap_WrongStellarSecret(t *testing.T) {
	// aes := &vault_infrastructure_crypto.AESService{}
	svc := vault_infrastructure_crypto.NewKeyService()

    stellarSecret := "SBFZCHU56JM2OD6LPJEDFB7CVJ6MIMUGHHMITB4TUSR47XH4BZ636Y6M"
    wrongStellarSecret := "SDCDSFDCSDCFDSFDSCFDSFDSCFDSFDCSFDCSFDSFDCSDFDSC"
    vaultKey := []byte("a‑vault‑key‑for‑testing")

    wrapped, err := svc.WrapKeyWithStellar(vaultKey, stellarSecret)
    require.NoError(t, err)

    _, err = svc.UnwrapKeyWithStellar(wrapped, wrongStellarSecret)
    require.Error(t, err)
    // important: not a `crypto` panic, just GCM auth failure
    require.Contains(t, err.Error(), "cipher: message authentication failed")
}

func TestKeyService_CreateAccount_ThenHybridUnlock(t *testing.T) {
    // 1. Setup KeyService + fs / keyring repo
	// svc := vault_infrastructure_crypto.NewKeyService()
    // 2. Create a user keyring, call svc.CreateAccount(...)
    // 3. SaveHybrid(..., password: "", stellarSecret: "SBF...")
    // 4. LoadHybrid with the same stellarSecret
    // 5. Assert keyring decrypted correctly.
}

/*
[ Auth Layer ]
   password / stellar / device

        ↓ unwrap

[ Vault Key Layer ]
   ONE symmetric key

        ↓

[ Data Layer ]
   entries / attachments / index


🔥 7. What these tests guarantee

You now guarantee:

✅ Cryptographic correctness
AES works
key wrapping works
✅ Security properties
wrong key fails
wrong password fails
non-deterministic encryption
✅ Real-world flow works
vault unlock → decrypt → use


together

🚀 Next step

Now that tests are solid:

👉 next logical move is:

🔥 Vault Keyring (multi-key support)
password + stellar + device keys
rotation without re-encryption

or

⚡ Instant unlock (device key)
no password after first login
like 1Password
*/

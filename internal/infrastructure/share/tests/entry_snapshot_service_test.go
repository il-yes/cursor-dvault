package share_tests

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	share_domain "vault-app/internal/domain/shared"
	share_infrastructure "vault-app/internal/infrastructure/share"
	"vault-app/internal/logger/logger"
	vaults_domain "vault-app/internal/vault/domain"
	"vault-app/internal/vault/ui"

	"github.com/stretchr/testify/require"
)

type MockVaultHandler struct {
    LoadAttachmentFn                           func(userID, vaultName, hash string) (string, error)
    UploadAttachementToIPFSWithEncryptionFn    func(userID string, req vault_ui.UploadAttachRequest) (string, error)
}

// Implement the real method
func (m *MockVaultHandler) LoadAttachment(userID, vaultName, hash string) (string, error) {
    return m.LoadAttachmentFn(userID, vaultName, hash)
}

func (m *MockVaultHandler) UploadAttachementToIPFSWithEncryption(
    userID string,
    req vault_ui.UploadAttachRequest,
) (string, error) {
    return m.UploadAttachementToIPFSWithEncryptionFn(userID, req)
}

type MockIPFSDownloader struct {
    GetFromIpfsFn func(ctx context.Context, cid string) ([]byte, error)
}

func (m *MockIPFSDownloader) GetFromIpfs(ctx context.Context, cid string) ([]byte, error) {
    return m.GetFromIpfsFn(ctx, cid)
}	




func TestEntrySnapshotService_Build_CIDShared(t *testing.T) {
    // 1. Set up mock
    mockVault := &MockVaultHandler{
        LoadAttachmentFn: func(userID, vaultName, hash string) (string, error) {
            // For example: return the "raw file" content for a test attachment
            if hash == "hash1" {
                return "file data for hash1", nil
            }
            return "", fmt.Errorf("unknown hash: %s", hash)
        },
        UploadAttachementToIPFSWithEncryptionFn: func(userID string, req vault_ui.UploadAttachRequest) (string, error) {
            // simulate IPFS upload success with deterministic CID
            if string(req.Data) == "file data for hash1" {
                return "QmCID123", nil
            }
            return "", fmt.Errorf("unexpected data")
        },
    }

    logger := logger.NewFromEnv() // or your minimal logger

    service := share_infrastructure.NewEntrySnapshotService(*logger, mockVault)

    // 2. Input: share with one attachment
    entrySnapshot := share_domain.EntrySnapshot{
        EntryName: "test entry",
        Type:      "login",
        Attachements: []vaults_domain.Attachment{
            {
                Hash: "hash1", // this will be loaded from Vault
            },
        },
    }

    shareEntry := share_domain.ShareEntry{
        OwnerID:      "owner1",
        EntryName:    "test entry",
        EntryType:    "login",
        AccessMode:   "read",
        EntrySnapshot: entrySnapshot,
    }

    // 3. Build request
    req := share_infrastructure.BuildRequest{
        Share:              &shareEntry,
        Recipient:          share_domain.Recipient{Email: "alice@example.com"},
        UserID:             "user123",
        UserSubscriptionID: "usersub123",
        VaultName:          "testvault",
        Password:           "testpassword",
    }

    // 4. Call the service
    res, err := service.Build(context.Background(), req)
    require.NoError(t, err)
    require.NotEmpty(t, res)

    // 5. Unmarshal and verify
    var snapshot share_domain.EntrySnapshot
    err = json.Unmarshal(res, &snapshot)
    require.NoError(t, err)

    require.Len(t, snapshot.Attachements, 1)
    attach := snapshot.Attachements[0]
    require.Equal(t, "hash1", attach.Hash)
    require.Equal(t, "QmCID123", attach.CIDShared) // ← CIDShared set and marshaled
}

func TestRecipientDecryptsSharedAttachment(t *testing.T) {
    const testdata = "test attachment content"

    // 1. Generate a symmetric key (32 bytes)
    symKey := make([]byte, 32)
    _, err := rand.Read(symKey)
    require.NoError(t, err)

    // 2. AES‑GCM setup
    block, err := aes.NewCipher(symKey)
    require.NoError(t, err)
    gcm, err := cipher.NewGCM(block)
    require.NoError(t, err)

    nonce := make([]byte, 12)
    _, err = rand.Read(nonce)
    require.NoError(t, err)

    encryptedCiphertext := gcm.Seal(nil, nonce, []byte(testdata), nil)

    // 3. Simulate IPFS-style storage: CIDShared → ciphertext
    mockIPFS := &MockIPFSDownloader{
        GetFromIpfsFn: func(ctx context.Context, cid string) ([]byte, error) {
            if cid == "QmCIDShared123" {
                // store as nonce + ciphertext
                return append(nonce, encryptedCiphertext...), nil
            }
            return nil, errors.New("unknown CID")
        },
    }

    rawBytes, err := mockIPFS.GetFromIpfs(context.Background(), "QmCIDShared123")
    require.NoError(t, err)

    // 4. Recipient decrypts
    block2, err := aes.NewCipher(symKey)
    require.NoError(t, err)
    gcm2, err := cipher.NewGCM(block2)
    require.NoError(t, err)

    nonce2 := rawBytes[:gcm2.NonceSize()]
    ciphertext2 := rawBytes[gcm2.NonceSize():]

    plain, err := gcm2.Open(nil, nonce2, ciphertext2, nil)
    require.NoError(t, err)
    require.Equal(t, []byte(testdata), plain)
}
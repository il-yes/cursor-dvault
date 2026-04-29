package share_tests

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	share_domain "vault-app/internal/domain/shared"
	share_infrastructure "vault-app/internal/infrastructure/share"
	"vault-app/internal/logger/logger"
	vaults_domain "vault-app/internal/vault/domain"
	vault_ui "vault-app/internal/vault/ui"

	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------------------------
//
//	MOCKS
//
// -----------------------------------------------------------------------------------------------
type MockVaultHandler struct {
	LoadAttachmentFn                        func(userID, vaultName, hash string) (string, error)
	UploadAttachementToIPFSWithEncryptionFn func(userID string, req vault_ui.UploadAttachRequest) (string, error)
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

// MockVaultHandler is a simple mock that implements VaultHandler's interface methods
type mockVaultHandler struct {
	loadAttachmentFunc func(userID, vaultName, hash string) (string, error)
	uploadFunc         func(userID string, req vault_ui.UploadAttachRequest) (string, error)
}

func (m *mockVaultHandler) LoadAttachment(userID, vaultName, hash string) (string, error) {
	if m.loadAttachmentFunc != nil {
		return m.loadAttachmentFunc(userID, vaultName, hash)
	}
	panic("LoadAttachment not implemented")
}

func (m *mockVaultHandler) UploadAttachementToIPFSWithEncryption(userID string, req vault_ui.UploadAttachRequest) (string, error) {
	if m.uploadFunc != nil {
		return m.uploadFunc(userID, req)
	}
	panic("UploadAttachementToIPFSWithEncryption not implemented")
}

// -----------------------------------------------------------------------------------------------
//
//	TESTS
//
// -----------------------------------------------------------------------------------------------
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
		OwnerID:       "owner1",
		EntryName:     "test entry",
		EntryType:     "login",
		AccessMode:    "read",
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
	err = json.Unmarshal(res.Raw, &snapshot)
	require.NoError(t, err)

	require.Len(t, snapshot.Attachements, 1)
	attach := snapshot.Attachements[0]
	require.Equal(t, "hash1", attach.Hash)
	require.Equal(t, "QmCID123", attach.HashShare) // ← CIDShared set and marshaled
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


func TestEntrySnapshotService_Process_ShouldUploadAttachmentForOwnerAndRecipients(t *testing.T) {
	// 1. Prepare test data
	ownerUserID := "owner-1"
	recipientID1 := "recipient-1"
	recipientID2 := "recipient-2"
	attachmentHash := "f648264d19025c76671a970cdf3196cbf16698a94817ebc3ccf5bb2dba35ec7a"

	// Mock VaultPayload with one attachment
	mockVaultPayload := vaults_domain.InitEmptyVaultPayload("test", "1.0.0")
	mockVaultPayload.Entries.Note = []vaults_domain.NoteEntry{
		{
			BaseEntry: vaults_domain.BaseEntry{
				ID: "entry-1",
				Attachments: []vaults_domain.Attachment{
					{
						ID:   "att-1",
						Hash: attachmentHash,
						Name: "background_medium.jpg",
					},
				},
			},
		},
	}

	// 2. Mock VaultHandler (spy LoadAttachment + UploadAttachementToIPFSWithEncryption)
	uploads := 0
	spyCIDs := []string{}

	mockVaultHandler := &mockVaultHandler{
		loadAttachmentFunc: func(userID, vaultName, hash string) (string, error) {
			return "fake bytes for " + hash, nil
		},
		uploadFunc: func(userID string, req vault_ui.UploadAttachRequest) (string, error) {
			var modePart string
			if req.EncryptionMode == "public" && len(req.SymKey) > 0 {
				modePart = "public"
			} else {
				modePart = "private"
			}

			// Make CID per recipient differ by something
			cid := fmt.Sprintf("fake-cid-%s-%d", modePart, uploads+1)
			uploads++
			spyCIDs = append(spyCIDs, cid)
			return cid, nil
		},
	}

	// 3. Mock EntrySnapshotService
	svc := &share_infrastructure.EntrySnapshotService{
		VaultHandler: mockVaultHandler,
		Logger:       *logger.NewFromEnv(),
	}

	// 4. Build request with one attachment and two recipients
	req := share_infrastructure.BuildRequest{
		Share: &share_domain.ShareEntry{
			EntryName: "test-entry",
			EntryType: "note",
			EntrySnapshot: share_domain.EntrySnapshot{
				Attachements: []vaults_domain.Attachment{
					{ID: "att-1", Hash: attachmentHash},
				},
			},
			Recipients: []share_domain.Recipient{
				{ID: recipientID1},
				{ID: recipientID2},
			},
		},
		UserID:             ownerUserID,
		UserSubscriptionID: "sub-1",
		VaultName:          "Leeks",
		Password:           "password",
		SymKey:             []byte("shared-sym-key-12345"),
	}

	// 5. Run
	snapshot, cids, attachementIDs, err := svc.Process(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 6. Assert
	if len(snapshot.Attachements) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(snapshot.Attachements))
	}

	att := snapshot.Attachements[0]
	if att.ID != "att-1" {
		t.Fatalf("expected attachment ID = att-1, got %s", att.ID)
	}

	// entries should be: one per recipient
	if uploads != 2 {
		t.Fatalf("expected 2 uploads (one per recipient), got %d", uploads)
	}

	if len(cids) != 2 {
		t.Fatalf("expected 2 CIDs, got %d", len(cids))
	}

	if len(attachementIDs) != 2 {
		t.Fatalf("expected 2 attachementIDs, got %d", len(attachementIDs))
	}

	// RecipientCIDs must be set for each recipient
	if len(att.RecipientCIDs) != 2 {
		t.Fatalf("expected 2 RecipientCIDs entries, got %d", len(att.RecipientCIDs))
	}

	cid1 := att.RecipientCIDs[recipientID1]
	if cid1 == "" {
		t.Error("expected non-empty RecipientCIDs[recipientID1]")
	}

	cid2 := att.RecipientCIDs[recipientID2]
	if cid2 == "" {
		t.Error("expected non-empty RecipientCIDs[recipientID2]")
	}

	if cid1 == cid2 {
		t.Error("recipient CIDs must be different")
	}


	for _, cid := range spyCIDs {
		if !strings.Contains(cid, "public") {
			t.Errorf("expected public-mode CID, got %s", cid)
		}
	}
}

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

	"github.com/stretchr/testify/require"

	"vault-app/internal/logger/logger"
	share_entry_domain "vault-app/internal/share_entry/domain"
	share_entry_infrastructure "vault-app/internal/share_entry/infrastructure"
	"vault-app/internal/utils"
	vaults_domain "vault-app/internal/vault/domain"
	vaults_service "vault-app/internal/vault/infrastructure/service"
	vault_ui "vault-app/internal/vault/ui"
)

// -----------------------------------------------------------------------------------------------
//
//	MOCKS
//
// -----------------------------------------------------------------------------------------------
type MockVaultHandler struct {
	LoadAttachmentFn                        func(userID, vaultName, hash string, formatReturned string) (*vaults_service.LoadAttachmentResponse, error)
	UploadAttachementToIPFSWithEncryptionFn func(userID string, req vault_ui.UploadAttachRequest) (string, error)
}

// Implement the real method
func (m *MockVaultHandler) LoadAttachment(
	userID, vaultName, hash string, formatReturned string) (*vaults_service.LoadAttachmentResponse, error) {
	return m.LoadAttachmentFn(userID, vaultName, hash, formatReturned)
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
	loadAttachmentFunc func(userID, vaultName, hash string, formatReturned string) (*vaults_service.LoadAttachmentResponse, error)
	uploadFunc         func(userID string, req vault_ui.UploadAttachRequest) (string, error)
}

func (m *mockVaultHandler) LoadAttachment(userID, vaultName, hash string, formatReturned string) (*vaults_service.LoadAttachmentResponse, error) {
	if m.loadAttachmentFunc != nil {
		return m.loadAttachmentFunc(userID, vaultName, hash, formatReturned)
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
	// =====================================================
	// Mock Vault
	// =====================================================
	mockVault := &MockVaultHandler{
		LoadAttachmentFn: func(
			userID,
			vaultName,
			hash,
			formatReturned string,
		) (*vaults_service.LoadAttachmentResponse, error) {

			if hash == "hash1" {
				return &vaults_service.LoadAttachmentResponse{
					File: []byte("file data for hash1"),
					Hash: "base64-encoded-file",
				}, nil
			}

			return nil, fmt.Errorf(
				"unknown hash: %s",
				hash,
			)
		},

		UploadAttachementToIPFSWithEncryptionFn: func(
			userID string,
			req vault_ui.UploadAttachRequest,
		) (string, error) {

			require.Equal(
				t,
				"base64-encoded-file",
				string(req.Data),
			)

			return "QmCID123", nil
		},
	}

	log := logger.NewFromEnv()

	service :=
		share_entry_infrastructure.NewEntrySnapshotService(
			*log,
			mockVault,
		)

	// =====================================================
	// Session attachment
	// =====================================================
	attachment := vaults_domain.Attachment{
		ID:      "att1",
		NodeCID: "nodeCID1",
		Hash:    "hash1",

		RecipientCIDs: map[string]string{},
	}

	vaultSession := vaults_domain.VaultPayload{
		BaseVaultContent: vaults_domain.BaseVaultContent{
			Attachments: []vaults_domain.Attachment{
				attachment,
			},
		},
	}

	// =====================================================
	// Entry snapshot
	// =====================================================
	entrySnapshot :=
		share_entry_domain.EntrySnapshot{
			EntryName: "test entry",
			Type:      "login",

			AttachmentCIDs: []string{
				"nodeCID1",
			},
		}

	shareEntry :=
		share_entry_domain.ShareEntry{
			OwnerID:   "owner1",
			EntryName: "test entry",
			EntryType: "login",

			EntrySnapshot: entrySnapshot,

			Recipients: []share_entry_domain.Recipient{
				{
					Email: "alice@example.com",

					PublicKey: "pubkey1",
				},
			},
		}

	// =====================================================
	// Request
	// =====================================================
	req :=
		share_entry_infrastructure.BuildRequest{
			Share:              &shareEntry,
			UserID:             "user123",
			UserSubscriptionID: "sub123",
			VaultName:          "vault1",
			Password:           "pwd",

			SymKey: []byte("symkey"),

			VaultSession: vaultSession,
		}

	// =====================================================
	// Execute
	// =====================================================
	res, err :=
		service.Build(
			context.Background(),
			req,
		)

	require.NoError(t, err)

	// =====================================================
	// Validate marshaled snapshot
	// =====================================================
	var snapshot share_entry_domain.EntrySnapshot

	err =
		json.Unmarshal(
			res.Raw,
			&snapshot,
		)

	require.NoError(t, err)

	require.Len(
		t,
		snapshot.Attachments,
		1,
	)

	att :=
		snapshot.Attachments[0]

	require.Equal(
		t,
		"hash1",
		att.Hash,
	)

	require.True(
		t,
		att.IsDirty,
	)

	require.Equal(
		t,
		"QmCID123",
		att.RecipientCIDs["pubkey1"],
	)

	// =====================================================
	// Validate returned attachments
	// =====================================================
	require.Len(
		t,
		res.Attachments,
		1,
	)

	require.True(
		t,
		res.Attachments[0].IsDirty,
	)

	require.Equal(
		t,
		"QmCID123",
		res.Attachments[0].
			RecipientCIDs["pubkey1"],
	)
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
func TestEntrySnapshotService_Build_RecipientCIDMap(
	t *testing.T,
) {
	mockVault := &MockVaultHandler{
		LoadAttachmentFn: func(
			userID,
			vaultName,
			hash,
			formatReturned string,
		) (*vaults_service.LoadAttachmentResponse, error) {

			return &vaults_service.LoadAttachmentResponse{
				File: []byte("file"),
				Hash: "base64-encoded-file",
			}, nil
		},

		UploadAttachementToIPFSWithEncryptionFn: func(
			userID string,
			req vault_ui.UploadAttachRequest,
		) (string, error) {

			return "QmCID123", nil
		},
	}

	log := logger.NewFromEnv()

	service :=
		share_entry_infrastructure.
			NewEntrySnapshotService(
				*log,
				mockVault,
			)

	attachment :=
		vaults_domain.Attachment{
			ID:      "att1",
			NodeCID: "nodeCID1",
			Hash:    "hash1",

			RecipientCIDs: map[string]string{},
		}

	vaultSession :=
		vaults_domain.VaultPayload{
			BaseVaultContent: vaults_domain.BaseVaultContent{
				Attachments: []vaults_domain.Attachment{
					attachment,
				},
			},
		}

	shareEntry :=
		share_entry_domain.ShareEntry{
			OwnerID: "owner1",

			EntrySnapshot: share_entry_domain.EntrySnapshot{
				AttachmentCIDs: []string{
					"nodeCID1",
				},
			},

			Recipients: []share_entry_domain.Recipient{
				{
					PublicKey: "pubkey1",
				},
				{
					PublicKey: "pubkey2",
				},
			},
		}

	req :=
		share_entry_infrastructure.BuildRequest{
			Share: &shareEntry,

			UserID:             "user1",
			UserSubscriptionID: "sub1",
			VaultName:          "vault1",
			Password:           "pwd",

			VaultSession: vaultSession,
		}

	res, err :=
		service.Build(
			context.Background(),
			req,
		)

	require.NoError(t, err)

	require.Len(
		t,
		res.Attachments,
		1,
	)

	att :=
		res.Attachments[0]

	require.Len(
		t,
		att.RecipientCIDs,
		2,
	)

	require.Equal(
		t,
		"QmCID123",
		att.RecipientCIDs["pubkey1"],
	)

	require.Equal(
		t,
		"QmCID123",
		att.RecipientCIDs["pubkey2"],
	)

	require.True(
		t,
		att.IsDirty,
	)

	utils.LogPretty("res RecipientCIDs", res.Snapshot.RecipientCIDs)
	utils.LogPretty("res Attachments", res.Snapshot.Attachments)

}

func TestEntrySnapshotService_Process_ShouldUploadAttachmentForRecipients(
	t *testing.T,
) {
	// =====================================================
	// Test data
	// =====================================================
	ownerUserID := "owner-1"

	attachmentHash :=
		"f648264d19025c76671a970cdf3196cbf16698a94817ebc3ccf5bb2dba35ec7a"

	attachmentNodeCID := "node-att-1"

	// =====================================================
	// Vault session attachment
	// =====================================================
	attachment :=
		vaults_domain.Attachment{
			ID:      "att-1",
			NodeCID: attachmentNodeCID,
			Hash:    attachmentHash,
			Name:    "background_medium.jpg",

			RecipientCIDs: map[string]string{},
		}

	vaultSession :=
		vaults_domain.VaultPayload{
			BaseVaultContent: vaults_domain.BaseVaultContent{
				Attachments: []vaults_domain.Attachment{
					attachment,
				},
			},
		}

	// =====================================================
	// Mock VaultHandler
	// =====================================================
	uploads := 0

	mockVaultHandler :=
		&mockVaultHandler{
			loadAttachmentFunc: func(
				userID,
				vaultName,
				hash,
				formatReturned string,
			) (
				*vaults_service.LoadAttachmentResponse,
				error,
			) {

				return &vaults_service.
					LoadAttachmentResponse{
					File: []byte(
						"fake bytes",
					),

					Hash: "base64-file",
				}, nil
			},

			uploadFunc: func(
				userID string,
				req vault_ui.UploadAttachRequest,
			) (
				string,
				error,
			) {

				uploads++

				require.Equal(
					t,
					string(
						vaults_domain.
							EncryptionPublic,
					),
					req.EncryptionMode,
				)

				return "shared-cid",
					nil
			},
		}

	// =====================================================
	// Service
	// =====================================================
	svc :=
		&share_entry_infrastructure.
			EntrySnapshotService{
			VaultHandler: mockVaultHandler,

			Logger: *logger.NewFromEnv(),
		}

	// =====================================================
	// Build request
	// =====================================================
	req :=
		share_entry_infrastructure.
			BuildRequest{

			UserID: ownerUserID,

			UserSubscriptionID: "sub-1",

			VaultName: "Leeks",

			Password: "password",

			SymKey: []byte(
				"shared-key",
			),

			VaultSession: vaultSession,

			Share: &share_entry_domain.
				ShareEntry{

				EntryName: "test-entry",

				EntryType: "note",

				EntrySnapshot: share_entry_domain.
					EntrySnapshot{

					AttachmentCIDs: []string{
						attachmentNodeCID,
					},
				},

				Recipients: []share_entry_domain.
					Recipient{
					{
						PublicKey: "pub1",
					},
					{
						PublicKey: "pub2",
					},
				},
			},
		}

	// =====================================================
	// Execute
	// =====================================================
	snapshot,
		attachments,
		err :=
		svc.Process(
			context.Background(),
			req,
		)

	require.NoError(
		t,
		err,
	)

	// =====================================================
	// Assertions
	// =====================================================

	// ONE upload only
	require.Equal(
		t,
		1,
		uploads,
	)

	require.Len(
		t,
		snapshot.Attachments,
		1,
	)

	require.Len(
		t,
		attachments,
		1,
	)

	att :=
		snapshot.Attachments[0]

	require.Equal(
		t,
		"att-1",
		att.ID,
	)

	require.True(
		t,
		att.IsDirty,
	)

	require.Len(
		t,
		att.RecipientCIDs,
		2,
	)

	require.Equal(
		t,
		"shared-cid",
		att.RecipientCIDs["pub1"],
	)

	require.Equal(
		t,
		"shared-cid",
		att.RecipientCIDs["pub2"],
	)

	// same mapping returned
	require.Equal(
		t,
		"shared-cid",
		attachments[0].
			RecipientCIDs["pub1"],
	)

	require.Equal(
		t,
		"shared-cid",
		attachments[0].
			RecipientCIDs["pub2"],
	)
}

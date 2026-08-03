package vaults_storage_tests

import (
	"context"
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	blockchain_ipfs "vault-app/internal/blockchain/ipfs"
	app_config "vault-app/internal/config"
	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/tracecore"
	"vault-app/internal/utils"
	vault_commands "vault-app/internal/vault/application/commands"
	vault_dto "vault-app/internal/vault/application/dto"
	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vaults_service "vault-app/internal/vault/infrastructure/service"
	vaults_storage_engine_constructor "vault-app/internal/vault/infrastructure/storage/engine/constructor"
	vaults_storage_engine_nodestore "vault-app/internal/vault/infrastructure/storage/engine/node_store"
	vaults_storage_engine_serializer "vault-app/internal/vault/infrastructure/storage/engine/serializer"
	vaults_storage_engine_type "vault-app/internal/vault/infrastructure/storage/engine/types"
)

// ============= MockIPFS =======================================================
type MockIPFS struct {
	Store map[string][]byte
}

func NewMockIPFS() *MockIPFS {
	return &MockIPFS{
		Store: make(map[string][]byte),
	}
}

func (m *MockIPFS) Add(data []byte) (string, error) {
	cid := fmt.Sprintf("cid_%d", len(m.Store)+1)
	m.Store[cid] = data
	return cid, nil
}

func (m *MockIPFS) Get(cid string) ([]byte, error) {
	data, ok := m.Store[cid]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return data, nil
}
func (f *MockIPFS) AddData(data []byte) (string, error) {
	return "cid-123", nil
}

// ============= MockCIDIpfsBuilder =======================================================
type MockCIDIpfsBuilder struct {
	Store map[string][]byte
}

func NewMockCIDIpfsBuilder() *MockCIDIpfsBuilder {
	return &MockCIDIpfsBuilder{
		Store: make(map[string][]byte),
	}
}

func (m *MockCIDIpfsBuilder) ComputeCID(data []byte) (string, error) {
	return fmt.Sprintf("cid_%x", sha256.Sum256(data)), nil
}

func (m *MockCIDIpfsBuilder) Build(data []byte) (string, error) {
	cid, _ := m.ComputeCID(data)
	m.Store[cid] = data
	return cid, nil
}

// ============= fakeCreateIPFSPayloadCommandHandler =======================================================
type fakeCreateIPFSPayloadCommandHandler struct {
	result string
	err    error
	called bool
}

func (f *fakeCreateIPFSPayloadCommandHandler) StoreOnIpfs(ctx context.Context, req vault_commands.CreateIPFSPayloadCommand) (string, error) {
	f.called = true
	return f.result, f.err
}

// ============= MockVaultRepo =======================================================
type MockVaultRepo struct {
	Vault         *vaults_domain.Vault
	existingVault *vaults_domain.Vault
	saveCalled    bool
	saveError     error
	savedVault    *vaults_domain.Vault
	updateCalled  bool
	updateError   error
	deleteCalled  bool
	deleteError   error
}

func (m *MockVaultRepo) GetVault(vaultID string) (*vaults_domain.Vault, error) {
	return m.Vault, nil
}

func (m *MockVaultRepo) UpdateVaultCID(vaultID, cid string) error {
	m.Vault.CID = cid
	return nil
}
func (f *MockVaultRepo) DeleteVault(string) error {
	f.deleteCalled = true
	return f.deleteError
}
func (f *MockVaultRepo) GetByUserIDAndName(string, string) (*vaults_domain.Vault, error) {
	if f.existingVault != nil {
		return f.existingVault, nil
	}
	return nil, vaults_domain.ErrVaultNotFound
}
func (f *MockVaultRepo) GetLatestByUserID(string) (*vaults_domain.Vault, error) {
	if f.existingVault != nil {
		return f.existingVault, nil
	}
	return nil, vaults_domain.ErrVaultNotFound
}
func (f *MockVaultRepo) SaveVault(v *vaults_domain.Vault) error {
	f.saveCalled = true
	f.savedVault = v
	return f.saveError
}
func (f *MockVaultRepo) UpdateVault(*vaults_domain.Vault) error {
	f.updateCalled = true
	return f.updateError
}
func (f *MockVaultRepo) GetVaultByCID(vaultID string) (*vaults_domain.Vault, error) {
	return f.Vault, nil
}

// ============= MockNodeRepo =======================================================
type MockNodeRepo struct {
	Entries      vaults_domain.Entries
	Folders      []vaults_domain.Folder
	Attachements []vaults_domain.Attachment
}

func (m *MockNodeRepo) GetEntries(s vault_session.Session) (*vaults_domain.Entries, error) {
	return &m.Entries, nil
}

func (m *MockNodeRepo) GetFolders(vault_session.Session) ([]vaults_domain.Folder, error) {
	return m.Folders, nil
}
func (m *MockNodeRepo) GetAttachements(session vault_session.Session) ([]vaults_domain.Attachment, error) {
	return m.Attachements, nil
}

// ============= fakeCryptoService =======================================================
type fakeCryptoService struct{}

func (f *fakeCryptoService) EncryptData(data []byte, password string) ([]byte, error) {
	return nil, nil
}
func (f *fakeCryptoService) Encrypt(data []byte, password string) ([]byte, error) {
	return nil, nil
}

// ============= fakeTracecoreClient =======================================================
type fakeTracecoreClient struct{}

func (f *fakeTracecoreClient) AddDataToIPFS(data []byte) (string, error) {
	return "cid-123", nil
}

func (f *fakeTracecoreClient) SyncVaultToIPFS(vaultName string) (string, error) {
	return "cid-123", nil
}

// mockUnlock
type mockUnlock struct{}

func (m *mockUnlock) Execute(cmd vault_dto.UnlockVaultCommand) (*vault_dto.UnlockVaultResult, error) {
	return &vault_dto.UnlockVaultResult{
		VaultKey: vaults_domain.VaultKey{
			Key: []byte("test-dek-32-bytes-long-key!!!!"),
		},
	}, nil
}

// mockIPFS
type mockIPFS struct {
	cid  string
	data []byte
}

func (m *mockIPFS) StoreOnIpfs(ctx interface{}, params interface{}) (string, error) {
	return m.cid, nil
}
func (m *mockIPFS) Add(context.Context, []byte) (string, error) {
	return m.cid, nil
}
func (m *mockIPFS) Get(ctx context.Context, cid string) ([]byte, error) {
	return nil, nil
}

type mockUnlockVaultHandler struct {
	ExecuteFunc func(cmd vault_dto.UnlockVaultCommand) (*vault_dto.UnlockVaultResult, error)
}

func (m *mockUnlockVaultHandler) Execute(cmd vault_dto.UnlockVaultCommand) (*vault_dto.UnlockVaultResult, error) {
	return m.ExecuteFunc(cmd)
}

type mockCryptoService struct {
	EncryptFunc func(data, key []byte) ([]byte, error)
}

func (m *mockCryptoService) Encrypt(data, key []byte) ([]byte, error) {
	return m.EncryptFunc(data, key)
}

func (c *mockCryptoService) Decrypt(ciphertext []byte, b []byte) ([]byte, error) {
	// test‑only, not needed for CreateIPFSPayloadCommandHandler test
	return nil, nil
}

func (m *mockCryptoService) DecryptSymKey(data []byte, key []byte) ([]byte, error) {
	return data[len(key):], nil
}

func (m *mockCryptoService) AsymetricDecrypt(privateKey string, encryptedKey string) ([]byte, error) {
	return nil, nil
}

type mockStorageFactory struct {
	NewFunc func(ctx *app_config_domain.VaultContext) app_config.StorageProvider
}

func (m *mockStorageFactory) New(ctx *app_config_domain.VaultContext) app_config.StorageProvider {
	return m.NewFunc(ctx)
}

// mockStorageProvider stub
type mockStorageProvider struct {
	AddFunc func(ctx context.Context, data []byte) (string, error)
	GetFunc func(ctx context.Context, cid string) ([]byte, error)
}

func (m *mockStorageProvider) Add(ctx context.Context, data []byte) (string, error) {
	if m.AddFunc != nil {
		return m.AddFunc(ctx, data)
	}
	return "QmSGH4oAre11ktm4DvMsQ7XT2xfyxd9ERU68cUFRrt6FB7", nil
}

func (m *mockStorageProvider) Get(ctx context.Context, cid string) ([]byte, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, cid)
	}
	// encode encrypted blob as base64 string
	// here, encrypted = base64 of GCM‑encrypted vaultBytes
	encrypted := "eyJ2ZX..." // or from earlier test, or deterministic
	return []byte(encrypted), nil
}

// ============= mockVaultHandler =======================================================
type mockVaultHandler struct {
	vps *vaults_domain.VaultPayload
}

func (m *mockVaultHandler) GetVaultSession(userID string) (*vaults_domain.VaultPayload, error) {
	return m.vps, nil
}

func (m *mockVaultHandler) UpdateEntryFor(
	userID string,
	entry any,
	syncMode bool,
) (*vaults_domain.VaultEntry, error) {
	ve, ok := entry.(vaults_domain.VaultEntry)
	if !ok {
		return nil, nil
	}
	return &ve, nil
}
func (m *mockVaultHandler) LoadAttachment(userID string, vaultName string, hash string, formatReturned string) (*vaults_service.LoadAttachmentResponse, error) {
	return nil, nil
}

// ============= mockVaultHandler with storage engine =======================================================
type mockVaultHandler2 struct {
	vps *vaults_domain.VaultPayload
}

func (m *mockVaultHandler2) GetVaultSession(userID string) (*vaults_domain.VaultPayload, error) {
	return m.vps, nil
}

func (m *mockVaultHandler2) UpdateEntryFor(
	userID string,
	entry any,
	syncMode bool,
) (*vaults_domain.VaultEntry, error) {
	ve, ok := entry.(vaults_domain.VaultEntry)
	if !ok {
		return nil, nil
	}
	return &ve, nil
}
func (m *mockVaultHandler2) LoadAttachment(userID string, vaultName string, hash string, formatReturned string) (*vaults_storage_engine_serializer.LoadAttachmentResponse, error) {
	return nil, nil
}

// ============= Helper =======================================================
func appConfig(userID string) app_config_domain.AppConfig {
	return app_config_domain.AppConfig{
		UserID:       userID,
		RepoID:       "my-repo-id",
		Branch:       "DefaultBranch",
		DefaultPhase: "DefaultDefaultPhase",
		VaultSettings: app_config_domain.VaultConfig{
			MaxEntries:       app_config_domain.DefaultMaxEntries,
			AutoSyncEnabled:  true,
			EncryptionScheme: "AES-256-GCM",
		},
		Blockchain: app_config_domain.BlockchainConfig{
			Stellar: app_config_domain.StellarConfig{
				Network:    "testnet",
				HorizonURL: "https://horizon-testnet.stellar.org",
				Fee:        100,
			},
			IPFS: app_config_domain.IPFSConfig{
				APIEndpoint: "http://localhost:5001",
				GatewayURL:  "https://ipfs.io/ipfs/",
			},
		},
		Storage: app_config.StorageConfig{
			Mode: app_config.StorageLocal, // ← production default

			LocalIPFS: app_config.IPFSConfig{
				APIEndpoint: "http://localhost:5001",
				GatewayURL:  "https://ipfs.io/ipfs/",
			},

			PrivateIPFS: app_config.IPFSConfig{
				APIEndpoint: "http://192.168.1.10:5001",
				GatewayURL:  "http://192.168.1.10:8080/ipfs/",
			},

			Cloud: app_config.CloudConfig{
				BaseURL: "https://ankhora.io/back",
			},

			EnterpriseS3: app_config.S3Config{
				Region:   "us-east-1",
				Bucket:   "ankhora-enterprise",
				Endpoint: "https://s3.us-east-1.amazonaws.com",
			},
		},
	}

}
func GetSession(userID string, v vaults_domain.VaultPayload) vault_session.Session {
	cb, _ := v.GetContentBytes()
	session := vault_session.Session{
		UserID:  userID,
		Vault:   cb,
		Runtime: &vault_session.RuntimeContext{VaultName: v.Name},
	}

	return session
}
func fakeVaultPayload(userID string, vaultName string) vaults_domain.VaultPayload {
	vaultPayload := vaults_domain.VaultPayload{
		Version: "1.0.0",
		Name:    vaultName,
		BaseVaultContent: vaults_domain.BaseVaultContent{
			Entries: vaults_domain.Entries{
				Login: []vaults_domain.LoginEntry{
					{
						BaseEntry: vaults_domain.BaseEntry{
							ID:        "entry1",
							Type:      "login",
							EntryName: "GitHub",
							IsDraft:   false,
						},
					},
					{
						BaseEntry: vaults_domain.BaseEntry{
							ID:        "entry2",
							Type:      "login",
							EntryName: "DraftEntry",
							IsDraft:   true, // should be skipped
						},
					},
				},
				Card: []vaults_domain.CardEntry{
					{
						BaseEntry: vaults_domain.BaseEntry{
							ID:        "entry3",
							Type:      "card",
							EntryName: "Visa",
							IsDraft:   false,
						},
					},
				},
			},
			Folders: []vaults_domain.Folder{
				{
					ID:   "folder1",
					Name: "Work",
				},
			},
		},
		Personal: vaults_domain.BaseVaultContent{
			Entries: vaults_domain.Entries{
				Login: []vaults_domain.LoginEntry{
					{
						BaseEntry: vaults_domain.BaseEntry{
							ID:        "entry1",
							Type:      "login",
							EntryName: "GitHub",
							IsDraft:   false,
						},
					},
					{
						BaseEntry: vaults_domain.BaseEntry{
							ID:        "entry2",
							Type:      "login",
							EntryName: "DraftEntry",
							IsDraft:   true, // should be skipped
						},
					},
				},
				Card: []vaults_domain.CardEntry{
					{
						BaseEntry: vaults_domain.BaseEntry{
							ID:        "entry3",
							Type:      "card",
							EntryName: "Visa",
							IsDraft:   false,
						},
					},
				},
			},
			Folders: []vaults_domain.Folder{
				{
					ID:   "folder1",
					Name: "Work",
				},
			},
		},
		Collaborative: vaults_domain.C3VaultContent{
			TrustGroups: []vaults_domain.TrustGroup{
				{
					ID:          "tg-oem",
					WorkspaceID: "workspace-001",
					Name:        "OEM Engineering Trust",
				},
			},
			Assets: []vaults_domain.Asset{
				{
					CID:         "bafy-asset-001",
					ContentHash: "sha256-demo",
					Size:        1024,
				},
			},

			ShareEntries: []vaults_domain.ShareEntry{
				{
					ID:           "share-001",
					AssetCID:     "bafy-asset-001",
					TrustGroupID: "tg-oem",
					WrappedDEK:   "wrapped-demo-key",
					CreatedBy:    userID,
					CreatedAt:    time.Now(),
				},
			},
		},
	}
	return vaultPayload
}

// ====================================================================
// 		TESTS
// ====================================================================

func TestCommitVault(t *testing.T) {
	expectedKey := []byte("vault-key")
	encryptedData := []byte("encrypted-data")
	expectedCID := "Qm123"
	userID := "user_1"
	vaultName := "ocean"
	vaultRepo := &MockVaultRepo{
		Vault: &vaults_domain.Vault{
			ID:  "vault1",
			CID: "",
		},
	}

	nodeRepo := &MockNodeRepo{
		Entries: vaults_domain.Entries{
			Login: []vaults_domain.LoginEntry{
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:        "entry1",
						Type:      "login",
						EntryName: "GitHub",
						IsDraft:   false,
					},
				},
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:        "entry2",
						Type:      "login",
						EntryName: "DraftEntry",
						IsDraft:   true, // should be skipped
					},
				},
			},
			Card: []vaults_domain.CardEntry{
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:        "entry3",
						Type:      "card",
						EntryName: "Visa",
						IsDraft:   false,
					},
				},
			},
		},
		Folders: []vaults_domain.Folder{
			{
				ID:   "folder1",
				Name: "Work",
			},
		},
	}
	vp := fakeVaultPayload(userID, vaultName)
	session := GetSession(userID, vp)
	cfgs, err := GetConfig(userID, vaultName)
	if err != nil {
		utils.LogPretty("Vault service (WRITE) error", err)
	}

	vc := app_config_domain.VaultContext{
		Configs:       *cfgs,
		UserID:        userID,
		VaultName:     vaultName,
		StorageConfig: cfgs.App.Storage,
	}

	mockUnlock := &mockUnlockVaultHandler{
		ExecuteFunc: func(cmd vault_dto.UnlockVaultCommand) (*vault_dto.UnlockVaultResult, error) {
			return &vault_dto.UnlockVaultResult{
				VaultKey: vaults_domain.VaultKey{
					Key: expectedKey,
				},
			}, nil
		},
	}

	mockCrypto := &mockCryptoService{
		EncryptFunc: func(data, key []byte) ([]byte, error) {
			// ✅ only validate key (important)
			if string(key) != string(expectedKey) {
				t.Fatalf("unexpected key passed to Encrypt")
			}

			// ✅ optional: ensure data is not empty
			if len(data) == 0 {
				t.Fatal("expected non-empty data for encryption")
			}

			// return deterministic encrypted output
			return encryptedData, nil
		},
	}
	var capturedData []byte

	mockStorage := &mockStorageProvider{
		AddFunc: func(ctx context.Context, data []byte) (string, error) {
			// 🔥 CRITICAL ASSERT: ensure encrypted data is used
			if string(data) != string(encryptedData) {
				t.Fatalf("expected encrypted data, got %s", string(data))
			}
			capturedData = data
			return expectedCID, nil
		},
	}

	mockFactory := &mockStorageFactory{
		NewFunc: func(ctx *app_config_domain.VaultContext) app_config.StorageProvider {
			return mockStorage
		},
	}

	ipfsHandler := &vault_commands.CreateIPFSPayloadCommandHandler{
		UnlockVaultHandler: mockUnlock,
		CryptoService:      mockCrypto,
		StorageFactory:     mockFactory,
	}
	// mode := vaults_service.IncrementalSync
	service := &vaults_service.VaultService{
		VaultCtx:    vc,
		Encryptor:   &vaults_service.NoopEncryptor{},
		Repo:        vaultRepo,
		NodeRepo:    nodeRepo,
		Password:    "test", // unused
		IPFSHandler: ipfsHandler,
		IsDraftMode: true,
	}

	service.Password = "password"


	ipfs := mockIPFS{}
	ipfsQueryHandler := &vault_queries.GetIPFSDataQuerryHandler{
		UnlockVaultHandler: mockUnlock,
		CryptoService:      &vault_infrastructure_crypto.AESService{},
		IpfsService:        &ipfs,
		StorageFactory:     mockFactory,
	}
	mockVaultH := &mockVaultHandler2{}

	nodeStore := vaults_storage_engine_nodestore.NewNodeStore(
		*ipfsQueryHandler,
		ipfsHandler,
		vc,
		true,
	)

	serializer := vaults_storage_engine_serializer.NewSerializerDryRun(
		nodeStore , 
		vc, 
		mockVaultH,
	)

	vaultConstructor := vaults_storage_engine_constructor.NewVaultConstructor(
		serializer, 
		nodeStore,
		vaultRepo,
	)

	modeC := vaults_storage_engine_type.IncrementalSync
	rootCID, _, _, _, err := vaultConstructor.Execute(session, modeC)

	// utils.LogPretty("vaultConstructor", cid)

	
	// rootCID, _, _, _, err := service.CommitVault(session, mode)
	utils.LogPretty("TestCommitVault - rootCID", rootCID)
	require.NoError(t, err)
	require.NotEmpty(t, rootCID)
	// require.NotEmpty(t, capturedData)
	// ✅ Vault CID updated
	require.Equal(t, rootCID, vaultRepo.Vault.CID)

	// Root must be a valid CID
	require.NotEmpty(t, rootCID)
	fmt.Errorf("capturedData : %s", capturedData)

	assert.True(t, vaultConstructor.NodeStore.DraftStorage.Exists(rootCID))
	assert.True(t, vaultConstructor.NodeStore.DraftStorage.Exists(service.Personal))
	assert.True(t, vaultConstructor.NodeStore.DraftStorage.Exists(service.C3))
}

func TestCommitVault_Integration(t *testing.T) {
	userID := "user_1"
	userPassword := "password"
	vaultName := "ocean"
	vaultRepo := &MockVaultRepo{
		Vault: &vaults_domain.Vault{
			ID:  "vault1",
			CID: "",
		},
	}

	ipfs := mockIPFS{}

	// mockVaultH := &mockVaultHandler{}
	// mode := vaults_service.IncrementalSync
	nodeRepo := &MockNodeRepo{
		Entries: vaults_domain.Entries{
			Login: []vaults_domain.LoginEntry{
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:        "entry1",
						Type:      "login",
						EntryName: "GitHub",
						IsDraft:   false,
					},
				},
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:        "entry2",
						Type:      "login",
						EntryName: "DraftEntry",
						IsDraft:   true, // should be skipped
					},
				},
			},
			Card: []vaults_domain.CardEntry{
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:        "entry3",
						Type:      "card",
						EntryName: "Visa",
						IsDraft:   false,
					},
				},
			},
		},
		Folders: []vaults_domain.Folder{
			{
				ID:   "folder1",
				Name: "Work",
			},
		},
		Attachements: []vaults_domain.Attachment{},
	}
	cfgs, err := GetConfig(userID, vaultName)
	if err != nil {
		utils.LogPretty("Vault service (WRITE) error", err)
	}

	vc := app_config_domain.VaultContext{
		SessionID:          "session_123",
		Configs:            *cfgs,
		UserID:             userID,
		VaultName:          vaultName,
		StorageConfig:      cfgs.App.Storage,
		UserOnboarding:     "onboarding_123",
		UserSubscriptionID: "subscription_123",
	}

	store := map[string][]byte{}
	mockStorage := &mockStorageProvider{
		AddFunc: func(ctx context.Context, data []byte) (string, error) {
			cid := fmt.Sprintf("cid-%d", len(store)+1)
			store[cid] = data
			return cid, nil
		},
		GetFunc: func(ctx context.Context, cid string) ([]byte, error) {
			data, ok := store[cid]
			if !ok {
				return nil, fmt.Errorf("not found")
			}
			return data, nil
		},
	}
	mockFactory := &mockStorageFactory{
		NewFunc: func(ctx *app_config_domain.VaultContext) app_config.StorageProvider {
			return mockStorage
		},
	}

	mockUnlock := &mockUnlockVaultHandler{
		ExecuteFunc: func(cmd vault_dto.UnlockVaultCommand) (*vault_dto.UnlockVaultResult, error) {
			return &vault_dto.UnlockVaultResult{
				VaultKey: vaults_domain.VaultKey{
					Key: []byte("0123456789abcdef0123456789abcdef"), // ✅ 32 bytes
				},
			}, nil
		},
	}
	ipfsCreateHandler := &vault_commands.CreateIPFSPayloadCommandHandler{
		UnlockVaultHandler: mockUnlock,
		CryptoService:      &vault_infrastructure_crypto.AESService{},
		StorageFactory:     mockFactory,
	}

	ipfsQueryHandler := &vault_queries.GetIPFSDataQuerryHandler{
		UnlockVaultHandler: mockUnlock,
		CryptoService:      &vault_infrastructure_crypto.AESService{},
		IpfsService:        &ipfs,
		StorageFactory:     mockFactory,
	}

	// service := &vaults_service.VaultService{
	// 	VaultHandler: mockVaultH,
	// 	VaultCtx:     vc,
	// 	Encryptor:    &vaults_service.AESEncryptor{},
	// 	Repo:         vaultRepo,
	// 	NodeRepo:     nodeRepo,
	// 	Password:     userPassword,
	// 	IPFSHandler:  ipfsCreateHandler,
	// 	DraftStorage: *vaults_service.NewDraftStorage(),
	// 	IsDraftMode:  true,
	// }

	vp := fakeVaultPayload(userID, vaultName)


	session := GetSession(userID, vp)
	nodeStore := vaults_storage_engine_nodestore.NewNodeStore(
		*ipfsQueryHandler,
		ipfsCreateHandler,
		vc,
		true,
	)

	mockVaultH2 := &mockVaultHandler2{}
	serializer := vaults_storage_engine_serializer.NewSerializerDryRun(
		nodeStore , 
		vc, 
		mockVaultH2,
	)

	vaultConstructor := vaults_storage_engine_constructor.NewVaultConstructor(
		serializer, 
		nodeStore,
		vaultRepo,
	)
	vaultConstructor.NodeRepo = nodeRepo
	vaultConstructor.Encryptor = &vaults_service.AESEncryptor{}

	modeC := vaults_storage_engine_type.IncrementalSync
	rootCID, _, _, _, err := vaultConstructor.Execute(session, modeC)

	// rootCID, _, _, _, err := service.CommitVault(GetSession(userID, vp), mode)
	require.NoError(t, err)
	require.NotEmpty(t, rootCID)

	res, err := ipfsQueryHandler.Execute(
		context.Background(),
		vault_queries.GetIPFSDataQuerry{
			CID:              rootCID,
			Password:         userPassword,
			Configs:          *cfgs,
			UserID:           userID,
			VaultName:        vaultName,
			UserOnboardingID: "onboarding_123",
		},
	)
	require.NoError(t, err)
	require.NotNil(t, res)

	require.NotEmpty(t, res.Node.Entries)
	require.NotEmpty(t, res.Node.Index)
}

/*
	func TestBuildEntries_Indexing(t *testing.T) {
		ipfs := NewMockIPFS()

		mode := vaults_service.SyncMode(1)
		service := &vaults_service.VaultService{
			Ipfs: ipfs,
		}

		entries := vaults_domain.Entries{
			Login: []vaults_domain.LoginEntry{
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:        "e1",
						Type:      "login",
						EntryName: "GitHub",
						FolderID:  "f1",
						CID:       "cid_1",
						IsDraft:   false, // 👈 key
						IsDirty:   false,
					},
				},
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:        "e3",
						Type:      "login",
						EntryName: "Amazon",
						FolderID:  "f2",
						CID:       "cid_2",
						IsDraft:   false, // 👈 key
						IsDirty:   false,
					},
				},
			},
			Card: []vaults_domain.CardEntry{
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:        "e2",
						Type:      "card",
						EntryName: "Visa",
						FolderID:  "f1",
						CID:       "cid_3",
						IsDraft:   false, // 👈 key
						IsDirty:   false,
					},
				},
			},
		}

		links, byType, byFolder, _, err := service.BuildEntries(entries, mode)
		require.NoError(t, err)

		// ✅ 2 entries processed
		require.Equal(t, 2, len(links)) // 3 entries - 1 draft skipped

		// ✅ byType index
		require.Len(t, byType["login"], 2)
		require.Len(t, byType["card"], 1)

		// ✅ byFolder index
		require.Len(t, byFolder["f1"], 2)

		// ✅ Should reuse existing CID
		found := false
		for _, l := range links {
			if l.CID == "e1" {
				require.Equal(t, "cid_1", l.CID)
				found = true
			}
		}
		require.True(t, found)
	}

	func TestBuildEntries_SkipDrafts(t *testing.T) {
		ipfs := NewMockIPFS()
		existingCID := "cid_existing"

		mode := vaults_service.SyncMode(1)

		service := &vaults_service.VaultService{
			Ipfs: ipfs,
		}

		entries := vaults_domain.Entries{
			Login: []vaults_domain.LoginEntry{
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:      "e1",
						Type:    "login",
						IsDraft: true,
						IsDirty: true,
						CID:     existingCID,
					},
				},
			},
		}

		links, _, _, _, err := service.BuildEntries(entries, mode)
		require.NoError(t, err)

		// ✅ No entries processed
		require.Len(t, links, 0)
	}

	func TestBuildFolders(t *testing.T) {
		ipfs := NewMockIPFS()

		service := &vaults_service.VaultService{
			Ipfs: ipfs,
		}

		folders := []vaults_domain.Folder{
			{ID: "f1", Name: "Work"},
		}

		links, err := service.BuildFolders(folders)
		require.NoError(t, err)

		require.Len(t, links, 1)
	}

	func TestBuildEntries_PartialRebuild(t *testing.T) {
		ipfs := NewMockIPFS()

		mode := vaults_service.SyncMode(1)
		service := &vaults_service.VaultService{
			Ipfs: ipfs,
		}

		entries := vaults_domain.Entries{
			Login: []vaults_domain.LoginEntry{
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:      "e1",
						Type:    "login",
						CID:     "cid_1",
						IsDraft: false,
						IsDirty: false,
					},
				},
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:      "e2",
						Type:    "login",
						IsDraft: false,
						IsDirty: true,
					},
				},
			},
		}

		links, _, _, _, err := service.BuildEntries(entries, mode)
		require.NoError(t, err)

		// ✅ First reused
		require.Equal(t, "cid_1", links[0].CID)

		// ✅ Second rebuilt
		require.NotEqual(t, "", links[1].CID)

		// ✅ Only ONE new IPFS write
		// Only dirty node triggers new CID generation
		require.NotEqual(t, "", links[1].CID)
		require.Equal(t, "cid_1", links[0].CID)

		var rebuilt bool
		var reused bool

		for _, l := range links {
			if l.CID == "e1" && l.CID == "cid_1" {
				reused = true
			}
			if l.CID == "e2" && l.CID != "" {
				rebuilt = true
			}
		}

		require.True(t, reused)
		require.True(t, rebuilt)
	}

	func TestCommitVault_DAGStructure(t *testing.T) {
		ipfs := NewMockIPFS()

		mode := vaults_service.SyncMode(1)
		vaultRepo := &MockVaultRepo{
			Vault: &vaults_domain.Vault{ID: "vault1"},
		}

		nodeRepo := &MockNodeRepo{
			Entries: vaults_domain.Entries{
				Login: []vaults_domain.LoginEntry{
					{
						BaseEntry: vaults_domain.BaseEntry{
							ID: "e1", Type: "login",
						},
					},
				},
			},
		}
		service := &vaults_service.VaultService{
			Ipfs:     ipfs,
			Repo:     vaultRepo,
			NodeRepo: nodeRepo,
		}

		rootCID, _, err := service.CommitVault("vault1", mode)
		require.NoError(t, err)

		// Fetch root node
		data, err := ipfs.Get(rootCID)
		require.NoError(t, err)

		var root vaults_domain.VaultNode
		err = json.Unmarshal(data, &root)
		require.NoError(t, err)

		// ✅ Root links exist
		require.NotEmpty(t, root.Entries.CID)
		require.NotEmpty(t, root.Index.CID)

		require.NotEqual(t, root.Entries.CID, "")

require.NotEqual(t, root.Index.CID, "")
require.NotEqual(t, root.Entries.CID, root.Index.CID)
}
*/
func TestBuildEntries_DryRun_Estimation(t *testing.T) {
	userID := "user_1"
	vaultName := "ocean"

	vaultRepo := &MockVaultRepo{
		Vault: &vaults_domain.Vault{
			ID:  "vault1",
			CID: "",
		},
	}

	nodeRepo := &MockNodeRepo{
		Entries: vaults_domain.Entries{
			Login: []vaults_domain.LoginEntry{
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:        "entry1",
						Type:      "login",
						EntryName: "GitHub",
						IsDraft:   false,
					},
				},
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:        "entry2",
						Type:      "login",
						EntryName: "DraftEntry",
						IsDraft:   true, // should be skipped
					},
				},
			},
			Card: []vaults_domain.CardEntry{
				{
					BaseEntry: vaults_domain.BaseEntry{
						ID:        "entry3",
						Type:      "card",
						EntryName: "Visa",
						IsDraft:   false,
					},
				},
			},
		},
		Folders: []vaults_domain.Folder{
			{
				ID:   "folder1",
				Name: "Work",
			},
		},
	}

	tracecoreClient := tracecore.NewTracecoreClient("test", "test", "test", "test")
	mockUnlock := &mockUnlockVaultHandler{
		ExecuteFunc: func(cmd vault_dto.UnlockVaultCommand) (*vault_dto.UnlockVaultResult, error) {
			return &vault_dto.UnlockVaultResult{
				VaultKey: vaults_domain.VaultKey{
					Key: []byte("12345678901234567890123456789012"), // 32 bytes
				},
			}, nil
		},
	}
	cfgs, err := GetConfig(userID, vaultName)
	if err != nil {
		utils.LogPretty("Vault service (WRITE) error", err)
	}

	vc := app_config_domain.VaultContext{
		Configs:       *cfgs,
		UserID:        userID,
		VaultName:     vaultName,
		StorageConfig: cfgs.App.Storage,
	}
	ipfsHandler := vault_commands.NewCreateIPFSPayloadCommandHandler(
		vaultRepo,
		*tracecoreClient,
		&blockchain_ipfs.DefaultStorageFactory{},
		mockUnlock,
	)

	service := vaults_service.NewVaultServiceDryRun(
		vaultRepo,
		nodeRepo,
		&vaults_service.NoopEncryptor{},
		vc,
		ipfsHandler,
	)

	entries := vaults_domain.Entries{
		Login: []vaults_domain.LoginEntry{
			{
				BaseEntry: vaults_domain.BaseEntry{
					ID:        "e1",
					Type:      "login",
					EntryName: "GitHub",
					IsDraft:   false,
				},
			},
		},
	}

	// -----------------------------
	// Act (first run)
	// -----------------------------
	links, _, _, _, err := service.BuildEntries(
		entries,
		vaults_service.IncrementalSync,
	)
	require.NoError(t, err)
	require.Len(t, links, 1)

	link := links[0]

	// -----------------------------
	// Assert structure correctness
	// -----------------------------
	// require.Equal(t, "GitHub", link.CID)
	require.NotEmpty(t, link.CID)

	// CID should be stable format (not empty, not random panic)
	firstCID := link.CID

	// -----------------------------
	// Act (second run - determinism)
	// -----------------------------
	links2, _, _, _, err := service.BuildEntries(
		entries,
		vaults_service.IncrementalSync,
	)
	require.NoError(t, err)
	require.Len(t, links2, 1)

	secondCID := links2[0].CID

	// -----------------------------
	// Assert determinism
	// -----------------------------
	require.Equal(t, firstCID, secondCID)

	// -----------------------------
	// Assert draft filtering correctness
	// -----------------------------
	require.NotContains(t, func() []string {
		var cids []string
		for _, l := range links2 {
			cids = append(cids, l.CID)
		}
		return cids
	}(), "DraftEntry")
}

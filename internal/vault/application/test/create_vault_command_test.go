package vault_commands_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	blockchain_ipfs "vault-app/internal/blockchain/ipfs"
	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/tracecore"
	vault_commands "vault-app/internal/vault/application/commands"
	vault_dto "vault-app/internal/vault/application/dto"
	vault_domain "vault-app/internal/vault/domain"
	vaults_domain "vault-app/internal/vault/domain"
)


//	------------------------------------------------------------------------------------------
//  FAKES
//	------------------------------------------------------------------------------------------
// // ======= fakeInitVaultHandler ===============
type fakeInitVaultHandler struct {
	result *vault_commands.InitializeVaultResult
	err    error
	called bool
	executeFn func(cmd vault_commands.InitializeVaultCommand) (*vault_commands.InitializeVaultResult, error)
}

func (f *fakeInitVaultHandler) Execute(cmd vault_commands.InitializeVaultCommand) (*vault_commands.InitializeVaultResult, error) {
	f.called = true
	return f.result, f.err
}

// ======= fakeCreateIPFSPayloadHandler ===============
type fakeCreateIPFSPayloadHandler struct {
	ipfs      vault_commands.IpfsServiceInterface
	result    *vault_commands.CreateIPFSPayloadCommandResult
	err       error
	called    bool
	executeFn func(ctx context.Context, vc app_config_domain.VaultContext, cmd vault_commands.CreateIPFSPayloadCommand) (*vault_commands.CreateIPFSPayloadCommandResult, error)
}

func (f *fakeCreateIPFSPayloadHandler) Execute(ctx context.Context, vc app_config_domain.VaultContext, cmd vault_commands.CreateIPFSPayloadCommand) (*vault_commands.CreateIPFSPayloadCommandResult, error) {
	f.called = true
	// return f.result, f.err
	return f.executeFn(ctx, vc, cmd)
}
func (f *fakeCreateIPFSPayloadHandler) SetIpfsService(i vault_commands.IpfsServiceInterface) {
	f.ipfs = i
}

// ======= fakeCryptoService ===============
type fakeCryptoService struct {
}

func (f *fakeCryptoService) EncryptData(data []byte, b []byte) ([]byte, error) {
	return nil, nil
}
func (f *fakeCryptoService) Encrypt(data []byte, b []byte) ([]byte, error) {
	return nil, nil
}
func (f *fakeCryptoService) Decrypt(data []byte, b []byte) ([]byte, error) {
	return nil, nil
}

// ======= fakeTracecoreClient ===============
type fakeTracecoreClient struct{}

func (f *fakeTracecoreClient) AddDataToIPFS(data []byte) (string, error) {
	return "cid-123", nil
}

func (f *fakeTracecoreClient) SyncVaultToIPFS(vaultName string) (string, error) {
	return "cid-123", nil
}

// ======= fakeIPFSService ===============
type fakeIPFSService struct {
}

func (f *fakeIPFSService) Add(ctx context.Context, data []byte) (string, error) {
	return "cid-123", nil
}


// ======= fakeVaultRepo ===============
type fakeVaultRepo struct {
	Vault         *vault_domain.Vault
	existingVault *vault_domain.Vault
	saveCalled    bool
	saveError     error
	savedVault    *vault_domain.Vault
	updateCalled  bool
	updateError   error
	deleteCalled  bool
	deleteError   error
	saveFn        func(v *vault_domain.Vault) error
	updateFn      func(v *vault_domain.Vault) error
}

func (f *fakeVaultRepo) CreateVault(vault *vault_domain.Vault) error {
	return nil
}

func (m *fakeVaultRepo) SaveVault(v *vault_domain.Vault) error {
	if m.saveFn != nil {
		return m.saveFn(v)
	}
	m.existingVault = v
	return nil
}
//	func (f *fakeVaultRepo) SaveVault(v *vault_domain.Vault) error {
//		f.saveCalled = true
//		f.savedVault = v
//		return f.saveError
//	}
//
// vaultcommand
func (f *fakeVaultRepo) GetLatestByUserID(userID string) (*vault_domain.Vault, error) {
	if userID == "user-1" { //"test_user" {
		return &vault_domain.Vault{
			UserID: userID,
			Name:   "test_vault_name",
		}, nil
	}

	return nil, gorm.ErrRecordNotFound
}

func (f *fakeVaultRepo) GetVault(string) (*vault_domain.Vault, error) {
	panic("not used")
}
func (f *fakeVaultRepo) UpdateVault(v *vault_domain.Vault) error {
	f.updateCalled = true
	// return f.updateError
	if f.updateFn != nil {
		return f.updateFn(v)
	}
	return nil
}
func (f *fakeVaultRepo) DeleteVault(string) error {
	f.deleteCalled = true
	return f.deleteError
}
func (f *fakeVaultRepo) GetByUserIDAndName(string, string) (*vault_domain.Vault, error) {
	if f.existingVault != nil {
		return f.existingVault, nil
	}
	return nil, vault_domain.ErrVaultNotFound
}
func (f *fakeVaultRepo) UpdateVaultCID(vaultID, cid string) error {
	f.updateCalled = true
	return f.updateError
}

//	------------------------------------------------------------------------------------------
//  TESTS
//	------------------------------------------------------------------------------------------


func TestCreateVault_Success(t *testing.T) {
	// -----------------------------
	// Arrange
	// -----------------------------
	userSubscriptionID := "userSub-1"
	userID := "user-1"
	password := "password"
	vaultName := "my-vault"

	vault := vault_domain.NewVault(userID, vaultName)

	mockInit := &fakeInitVaultHandler{
		result: &vault_commands.InitializeVaultResult{
			Vault: vault,
		},
		executeFn: func(cmd vault_commands.InitializeVaultCommand) (*vault_commands.InitializeVaultResult, error) {
			require.Equal(t, userID, cmd.UserID)
			require.Equal(t, vaultName, cmd.VaultName)

			return &vault_commands.InitializeVaultResult{
				Vault: vault,
			}, nil
		},
	}

	expectedCID := "cid-123"

	mockIPFS := &fakeCreateIPFSPayloadHandler{
		executeFn: func(ctx context.Context, vc app_config_domain.VaultContext, cmd vault_commands.CreateIPFSPayloadCommand) (*vault_commands.CreateIPFSPayloadCommandResult, error) {

			// 🔥 critical assertions
			require.NotNil(t, cmd.Data)
			require.Equal(t, password, cmd.Password)

			return &vault_commands.CreateIPFSPayloadCommandResult{
				CID: expectedCID,
			}, nil
		},
	}

	var savedVault *vault_domain.Vault

	mockRepo := &fakeVaultRepo{
		updateFn: func(v *vault_domain.Vault) error {
			savedVault = v
			return nil
		},
	}

	handler := vault_commands.NewCreateVaultCommandHandler(mockInit, mockIPFS, mockRepo)

	cmd := vault_commands.CreateVaultCommand{
		UserID:             userID,
		VaultName:          vaultName,
		Password:           password,
		UserSubscriptionID: userSubscriptionID,
		AppConfig:          app_config_domain.AppConfig{},
	}

	// -----------------------------
	// Act
	// -----------------------------
	res, err := handler.CreateVault(cmd)

	// -----------------------------
	// Assert
	// -----------------------------
	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotNil(t, res.Vault)

	// ✅ CID attached
	require.Equal(t, expectedCID, res.Vault.CID)

	// ✅ subscription attached
	require.Equal(t, userSubscriptionID, res.Vault.UserSubscriptionID)

	// ✅ persisted vault
	require.NotNil(t, savedVault)
	require.Equal(t, expectedCID, savedVault.CID)
	require.Equal(t, userSubscriptionID, res.Vault.UserSubscriptionID)
	require.Greater(t, len(cmd.UserID), 0)
	require.Greater(t, len(cmd.VaultName), 0)
}

func CreateVault_FailsIfInitializeFails(t *testing.T) {
	repo := &fakeVaultRepo{}
	ipfsService := &fakeIPFSService{}
	initHandler := &fakeInitVaultHandler{
		err: errors.New("init failed"),
	}
	ipfsHandler := &fakeCreateIPFSPayloadHandler{}
	ipfsHandler.SetIpfsService(ipfsService)

	handler := vault_commands.NewCreateVaultCommandHandler(
		initHandler,
		ipfsHandler,
		repo,
	)

	cmd := vault_commands.CreateVaultCommand{
		UserID: "user-1",
	}

	result, err := handler.CreateVault(cmd)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.True(t, initHandler.called)
	assert.False(t, ipfsHandler.called)
}

func CreateVault_FailsIfIPFSFails(t *testing.T) {
	repo := &fakeVaultRepo{}

	initResult := &vault_commands.InitializeVaultResult{
		Vault: vault_domain.NewVault("user-1", "vault-name"),
	}
	initHandler := &fakeInitVaultHandler{
		result: initResult,
	}

	ipfsHandler := &fakeCreateIPFSPayloadHandler{
		err: errors.New("ipfs failed"),
	}

	handler := vault_commands.NewCreateVaultCommandHandler(
		initHandler,
		ipfsHandler,
		repo,
	)

	cmd := vault_commands.CreateVaultCommand{
		UserID:   "user-1",
		Password: "secret",
	}

	result, err := handler.CreateVault(cmd)

	require.Error(t, err)
	assert.Nil(t, result)

	assert.True(t, initHandler.called)
	assert.True(t, ipfsHandler.called)
}

func CreateVault_AttachesCIDToVault(t *testing.T) {
	expectedKey := []byte("vault-key")

	mockUnlock := &mockUnlockVaultHandler{
		ExecuteFunc: func(cmd vault_dto.UnlockVaultCommand) (*vault_dto.UnlockVaultResult, error) {
			return &vault_dto.UnlockVaultResult{
				VaultKey: vaults_domain.VaultKey{
					Key: expectedKey,
				},
			}, nil
		},
	}
	repo := &fakeVaultRepo{}
	ipfsService := &fakeIPFSService{}
	initHandler := vault_commands.NewInitializeVaultCommandHandler(&gorm.DB{})
	tracecoreClient := tracecore.NewTracecoreClient("test", "test", "test", "test")

	ipfsHandler := vault_commands.NewCreateIPFSPayloadCommandHandler(repo, *tracecoreClient, &blockchain_ipfs.DefaultStorageFactory{}, mockUnlock)
	ipfsHandler.SetIpfsService(ipfsService)

	handler := vault_commands.NewCreateVaultCommandHandler(
		initHandler,
		ipfsHandler,
		repo,
	)

	cmd := vault_commands.CreateVaultCommand{
		UserID:   "user-1",
		Password: "secret",
	}

	result, err := handler.CreateVault(cmd)

	require.NoError(t, err)
	// We expect cid-123 because fakeIPFSService returns it.
	assert.Equal(t, "cid-123", result.Vault.CID)
}

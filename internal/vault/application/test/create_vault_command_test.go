package vault_commands_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"vault-app/internal/tracecore"
	vault_commands "vault-app/internal/vault/application/commands"
	vault_domain "vault-app/internal/vault/domain"
)

//
// ---------- FAKES ----------
//
type fakeInitVaultHandler struct {
	result *vault_commands.InitializeVaultResult
	err    error
	called bool
}

func (f *fakeInitVaultHandler) Execute(cmd vault_commands.InitializeVaultCommand) (*vault_commands.InitializeVaultResult, error) {
	f.called = true
	return f.result, f.err
}

type fakeCreateIPFSPayloadHandler struct {
	ipfs vault_commands.IpfsServiceInterface
	result *vault_commands.CreateIPFSPayloadCommandResult
	err    error
	called bool
}

func (f *fakeCreateIPFSPayloadHandler) Execute(cmd vault_commands.CreateIPFSPayloadCommand) (*vault_commands.CreateIPFSPayloadCommandResult, error) {
	f.called = true
	return f.result, f.err
}
func (f *fakeCreateIPFSPayloadHandler) SetIpfsService(i vault_commands.IpfsServiceInterface) {
	f.ipfs = i
}

type fakeCryptoService struct{}

func (f *fakeCryptoService) EncryptData(data []byte, password string) ([]byte, error) {
	return nil, nil
}
func (f *fakeCryptoService) Encrypt(data []byte, password string) ([]byte, error) {
	return nil, nil
}

type fakeTracecoreClient struct{}

func (f *fakeTracecoreClient) AddDataToIPFS(data []byte) (string, error) {
	return "cid-123", nil
}

func (f *fakeTracecoreClient) SyncVaultToIPFS(vaultName string) (string, error) {
	return "cid-123", nil
}

type fakeIPFSService struct{}

func (f *fakeIPFSService) Add(ctx context.Context, data []byte) (string, error) {
	return "cid-123", nil
}

func (f *fakeVaultRepo) SaveVault(v *vault_domain.Vault) error {
	f.saveCalled = true
	f.savedVault = v
	return f.saveError
}

func (f *fakeVaultRepo) GetLatestByUserID(userID string) (*vault_domain.Vault, error) {
    if userID == "test_user" {
        return &vault_domain.Vault{
            UserID: userID,
            Name:   "test_vault_name",
        }, nil
    }
    return nil, nil
}
func (f *fakeVaultRepo) GetVault(string) (*vault_domain.Vault, error) {
	panic("not used")
}
func (f *fakeVaultRepo) UpdateVault(*vault_domain.Vault) error {
	f.updateCalled = true
	return f.updateError
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

//
// ---------- TESTS ----------
//

func CreateVault_Success(t *testing.T) {
	repo := &fakeVaultRepo{}
	cryptoService := &fakeCryptoService{}
	ipfsService := &fakeIPFSService{}
	tracecoreClient := tracecore.NewTracecoreClient("test", "test", "test", "test")

	initHandler := vault_commands.NewInitializeVaultCommandHandler(&gorm.DB{})
	ipfsHandler := vault_commands.NewCreateIPFSPayloadCommandHandler(repo, cryptoService, *tracecoreClient)
	ipfsHandler.SetIpfsService(ipfsService)

	handler := vault_commands.NewCreateVaultCommandHandler(
		initHandler,
		ipfsHandler,
		repo,
	)

	cmd := vault_commands.CreateVaultCommand{
		UserID:    "user-1",
		VaultName: "My Vault",
		Password:  "secret",
	}

	result, err := handler.CreateVault(cmd)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "cid-123", result.Vault.CID)
	assert.False(t, result.ReusedExisting)
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
	repo := &fakeVaultRepo{}
	ipfsService := &fakeIPFSService{}
	initHandler := vault_commands.NewInitializeVaultCommandHandler(&gorm.DB{})
	tracecoreClient := tracecore.NewTracecoreClient("test", "test", "test", "test")

	ipfsHandler := vault_commands.NewCreateIPFSPayloadCommandHandler(repo, &fakeCryptoService{}, *tracecoreClient)
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

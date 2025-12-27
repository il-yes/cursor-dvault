package vault_commands_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

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
	result *vault_commands.CreateIPFSPayloadCommandResult
	err    error
	called bool
}

func (f *fakeCreateIPFSPayloadHandler) Execute(cmd vault_commands.CreateIPFSPayloadCommand) (*vault_commands.CreateIPFSPayloadCommandResult, error) {
	f.called = true
	return f.result, f.err
}

type fakeCryptoService struct{}

func (f *fakeCryptoService) EncryptData(data []byte, password string) ([]byte, error) {
	return nil, nil
}
func (f *fakeCryptoService) Encrypt(data []byte, password string) ([]byte, error) {
	return nil, nil
}

type fakeIPFSService struct{}

func (f *fakeIPFSService) AddData(data []byte) (string, error) {
	return "cid-123", nil
}

func (f *fakeVaultRepo) SaveVault(v *vault_domain.Vault) error {
	f.saveCalled = true
	f.savedVault = v
	return f.saveError
}

func (f *fakeVaultRepo) GetLatestByUserID(string) (*vault_domain.Vault, error) {
	if f.existingVault != nil {
		return f.existingVault, nil
	}
	return nil, vault_domain.ErrVaultNotFound
}
func (f *fakeVaultRepo) GetVault(string) (*vault_domain.Vault, error) {
	panic("not used")
}
func (f *fakeVaultRepo) UpdateVault(*vault_domain.Vault) error {
	panic("not used")
}
func (f *fakeVaultRepo) DeleteVault(string) error {
	panic("not used")
}

//
// ---------- TESTS ----------
//

func TestCreateVault_Success(t *testing.T) {
	repo := &fakeVaultRepo{}
	cryptoService := &fakeCryptoService{}
	ipfsService := &fakeIPFSService{}

	initHandler := vault_commands.NewInitializeVaultCommandHandler(&gorm.DB{})
	ipfsHandler := vault_commands.NewCreateIPFSPayloadCommandHandler(repo, cryptoService, ipfsService)

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

func TestCreateVault_FailsIfInitializeFails(t *testing.T) {
	repo := &fakeVaultRepo{}
	initHandler := &fakeInitVaultHandler{
		err: errors.New("init failed"),
	}
	ipfsHandler := &fakeCreateIPFSPayloadHandler{}

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

func TestCreateVault_FailsIfIPFSFails(t *testing.T) {
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

func TestCreateVault_AttachesCIDToVault(t *testing.T) {
	repo := &fakeVaultRepo{}
	initHandler := vault_commands.NewInitializeVaultCommandHandler(&gorm.DB{})

	ipfsHandler := vault_commands.NewCreateIPFSPayloadCommandHandler(repo, &fakeCryptoService{}, &fakeIPFSService{})

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

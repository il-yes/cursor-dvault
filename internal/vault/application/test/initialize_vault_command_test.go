package vault_commands_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"vault-app/internal/utils"
	vault_commands "vault-app/internal/vault/application/commands"
	vault_domain "vault-app/internal/vault/domain"
)

// --------------------------------------------------------------------------------------------------
// FAKES
// --------------------------------------------------------------------------------------------------
// ============= fakeVaultRepoMock =======================================================
type fakeVaultRepoMock struct {
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
	DefaultVaultNameTest bool
	isCreate	bool
}

func (f *fakeVaultRepoMock) CreateVault(vault *vault_domain.Vault) error {
	return nil
}

func (m *fakeVaultRepoMock) SaveVault(v *vault_domain.Vault) error {
	if m.saveFn != nil {
		m.Vault = v
		return m.saveFn(v)
	}
	if m.DefaultVaultNameTest {
		utils.LogPretty("fakeVaultRepoMock - SaveVault - Vault", v)
		m.Vault = v
		return nil
	}
	m.existingVault = v
	return nil
}

func (f *fakeVaultRepoMock) DeleteVault(string) error {
	f.deleteCalled = true
	return f.deleteError
}
func (f *fakeVaultRepoMock) GetByUserIDAndName(string, string) (*vault_domain.Vault, error) {
	if f.existingVault != nil {
		return f.existingVault, nil
	}
	return nil, vault_domain.ErrVaultNotFound
}
func (f *fakeVaultRepoMock) GetLatestByUserID(userID string) (*vault_domain.Vault, error) {
	if userID == "user-1" && !f.isCreate { 
		return f.existingVault, nil
	}

	return nil, gorm.ErrRecordNotFound
}

func (f *fakeVaultRepoMock) GetVault(string) (*vault_domain.Vault, error) {
	panic("not used")
}
func (f *fakeVaultRepoMock) UpdateVault(v *vault_domain.Vault) error {
	f.updateCalled = true
	// return f.updateError
	if f.updateFn != nil {
		return f.updateFn(v)
	}
	return nil
}
func (f *fakeVaultRepoMock) UpdateVaultCID(vaultID, cid string) error {
	f.updateCalled = true
	return f.updateError
}


// --------------------------------------------------------------------------------------------------
// TESTS
// --------------------------------------------------------------------------------------------------
func TestInitializeVault_CreateNewVault(t *testing.T) {
	repo := &fakeVaultRepoMock{
		isCreate: true,
	}

	h := &vault_commands.InitializeVaultCommandHandler{
		VaultRepo: repo,
	}

	cmd := vault_commands.InitializeVaultCommand{
		UserID:    "user-1",
		VaultName: "my-vault",
	}

	res, err := h.Execute(cmd)

	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotNil(t, res.Vault)

	require.Equal(t, "user-1", res.Vault.UserID)
	require.Equal(t, "my-vault", res.Vault.Name)
}

func TestInitializeVault_ReturnExistingVault(t *testing.T) {
	existing := &vault_domain.Vault{
		ID:     "vault-1",
		UserID: "user-1",
		Name:   "existing-vault",
		Type:   "default",
	}

	saveCalled := false

	repo := &fakeVaultRepoMock{
		existingVault: existing,
		saveFn: func(v *vault_domain.Vault) error {
			return nil
		},
	}

	h := &vault_commands.InitializeVaultCommandHandler{
		VaultRepo: repo,
	}

	cmd := vault_commands.InitializeVaultCommand{
		UserID:    "user-1",
		VaultName: "existing-vault",
	}
	repo.saveCalled = false

	res, err := h.Execute(cmd)
	utils.LogPretty("res", res)

	require.NoError(t, err)
	require.NotNil(t, res)

	// 🔥 MUST return existing vault
	require.Equal(t, existing.Name, res.Vault.Name)

	// 🔥 MUST NOT create new vault
	require.False(t, saveCalled)
}

func TestInitializeVault_DefaultVaultName(t *testing.T) {
	repo := &fakeVaultRepoMock{
		DefaultVaultNameTest: true,
	}

	var saved *vault_domain.Vault

	repo.saveFn = func(v *vault_domain.Vault) error {
		saved = v
		return nil
	}

	h := &vault_commands.InitializeVaultCommandHandler{
		VaultRepo: repo,
	}

	cmd := vault_commands.InitializeVaultCommand{
		UserID:    "user-1",
		VaultName: "",
	}

	res, err := h.Execute(cmd)
	utils.LogPretty("res", res)

	require.NoError(t, err)
	require.NotNil(t, res)

	require.Equal(t, "user-1-vault", saved.Name)
	require.Equal(t, "user-1-vault", res.Vault.Name)
}

func TestInitializeVault_PropagatesSaveError(t *testing.T) {
	handler := &vault_commands.InitializeVaultCommandHandler{
		VaultRepo: nil,
	}
	// handler := vault_commands.NewInitializeVaultCommandHandler(&gorm.DB{})

	cmd := vault_commands.InitializeVaultCommand{
		UserID: "user-1",
	}

	result, err := handler.Execute(cmd)

	require.Error(t, err)
	assert.Nil(t, result)
}
func TestInitializeVault_SaveError(t *testing.T) {
	repo := &fakeVaultRepoMock{
		saveFn: func(v *vault_domain.Vault) error {
			return errors.New("db error")
		},
	}

	h := &vault_commands.InitializeVaultCommandHandler{
		VaultRepo: repo,
	}

	cmd := vault_commands.InitializeVaultCommand{
		UserID:    "user-1",
		VaultName: "vault",
	}

	res, err := h.Execute(cmd)

	require.Error(t, err)
	require.Nil(t, res)
	require.Contains(t, err.Error(), "failed to persist vault metadata")
}

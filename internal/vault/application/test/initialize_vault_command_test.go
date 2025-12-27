package vault_commands_test

import (
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

type fakeVaultRepo struct {
	existingVault *vault_domain.Vault
	saveCalled    bool
	saveError     error
	savedVault *vault_domain.Vault
}





//
// ---------- TESTS ----------
//

func TestInitializeVault_CreatesNewVault(t *testing.T) {
	repo := &fakeVaultRepo{}
	handler := vault_commands.NewInitializeVaultCommandHandler(&gorm.DB{})

	cmd := vault_commands.InitializeVaultCommand{
		UserID:    "user-1",
		VaultName: "My Vault",
	}

	result, err := handler.Execute(cmd)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Vault)

	assert.Equal(t, "user-1", result.Vault.UserID)
	assert.Equal(t, "My Vault", result.Vault.Name)
	assert.True(t, repo.saveCalled)
}

func TestInitializeVault_Idempotent_ReturnsExisting(t *testing.T) {
	repo := &fakeVaultRepo{}
	existing := vault_domain.NewVault("user-1", "Existing Vault")

	handler := vault_commands.NewInitializeVaultCommandHandler(&gorm.DB{})

	cmd := vault_commands.InitializeVaultCommand{
		UserID: "user-1",
	}

	result, err := handler.Execute(cmd)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, existing, result.Vault)
	assert.False(t, repo.saveCalled, "should not save when vault exists")
}

func TestInitializeVault_DefaultName(t *testing.T) {
	handler := vault_commands.NewInitializeVaultCommandHandler(&gorm.DB{})

	cmd := vault_commands.InitializeVaultCommand{
		UserID: "user-42",
	}

	result, err := handler.Execute(cmd)

	require.NoError(t, err)
	require.NotNil(t, result.Vault)

	assert.Equal(t, "user-42-vault", result.Vault.Name)
}

func TestInitializeVault_PropagatesSaveError(t *testing.T) {


	handler := vault_commands.NewInitializeVaultCommandHandler(&gorm.DB{})

	cmd := vault_commands.InitializeVaultCommand{
		UserID: "user-1",
	}

	result, err := handler.Execute(cmd)

	require.Error(t, err)
	assert.Nil(t, result)
}

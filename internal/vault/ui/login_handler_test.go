package vault_ui

import (
	"testing"

	"github.com/stretchr/testify/require"

	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	vaults_domain "vault-app/internal/vault/domain"
)

func TestLoginHandlerAddSuccess(t *testing.T) {
	t.Parallel()

	userID := "101"
	handler := newTestLoginHandler()
	handler.Vault = &vaults_domain.VaultPayload{}

	entry := &vaults_domain.LoginEntry{
		BaseEntry: vaults_domain.BaseEntry{
			EntryName: "Primary Account",
		},
		UserName: "user@example.com",
		Password: "super-secret",
	}

	result, err := handler.Add(userID, entry)
	require.NoError(t, err)
	require.NotNil(t, result)

	addedEntry := result
	require.NotEmpty(t, addedEntry.Entries.Login)
	require.Equal(t, "Primary Account", addedEntry.Entries.Login[0].EntryName)

	require.Len(t, handler.Vault.Entries.Login, 1)
	require.Equal(t, addedEntry.Entries.Login[0].ID, handler.Vault.Entries.Login[0].ID)
}

func TestLoginHandlerAddNoSession(t *testing.T) {
	t.Parallel()

	handler := newTestLoginHandler()
	entry := &vaults_domain.LoginEntry{}

	result, err := handler.Add("55", entry)
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "vault not initialized for user 55")
}

func TestLoginHandlerAddInvalidType(t *testing.T) {
	t.Parallel()

	handler := newTestLoginHandler()

	invalidEntry := &vaults_domain.CardEntry{}

	result, err := handler.Add("1", invalidEntry)
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "entry does not implement VaultEntry interface")
}

func TestLoginHandlerEditSuccess(t *testing.T) {
	t.Parallel()

	handler := newTestLoginHandler()
	handler.Vault = &vaults_domain.VaultPayload{}

	entry := vaults_domain.LoginEntry{
		BaseEntry: vaults_domain.BaseEntry{
			ID:        "entry-1",
			EntryName: "Old Name",
		},
	}
	handler.Vault.Entries.Login = append(handler.Vault.Entries.Login, entry)

	updated := &vaults_domain.LoginEntry{
		BaseEntry: vaults_domain.BaseEntry{
			ID:        "entry-1",
			EntryName: "Updated Name",
		},
		UserName: "new@example.com",
		Password: "new-pass",
	}

	result, err := handler.Edit("1", updated)
	require.NoError(t, err)
	require.NotNil(t, result)

	modifiedEntry := result.Entries.Login[0]
	require.True(t, modifiedEntry.IsDraft)
	require.Equal(t, "Updated Name", modifiedEntry.EntryName)
	require.Equal(t, "new@example.com", modifiedEntry.UserName)
	require.True(t, modifiedEntry.IsDraft)
}

func TestLoginHandlerEditEntryNotFound(t *testing.T) {
	t.Parallel()

	handler := newTestLoginHandler()
	handler.Vault = &vaults_domain.VaultPayload{}

	entry := &vaults_domain.LoginEntry{
		BaseEntry: vaults_domain.BaseEntry{
			ID:        "missing",
			EntryName: "Ghost",
		},
	}

	result, err := handler.Edit("1", entry)
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "entry with ID missing not found for user 1")
}

func TestLoginHandlerEditInvalidType(t *testing.T) {
	t.Parallel()

	handler := newTestLoginHandler()

	result, err := handler.Edit("1", struct{}{})
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "invalid type: expected LoginEntry")
}

func TestLoginHandlerTrashAndRestore(t *testing.T) {
	t.Parallel()

	entry := vaults_domain.LoginEntry{
		BaseEntry: vaults_domain.BaseEntry{
			ID:        "entry-trashed",
			EntryName: "Disposable",
		},
	}

	handler := newTestLoginHandler()
	handler.Vault = &vaults_domain.VaultPayload{}
	handler.Vault.Entries.Login = append(handler.Vault.Entries.Login, entry)

	_, err := handler.Trash("1", "entry-trashed")
	require.NoError(t, err)
	require.True(t, handler.Vault.Entries.Login[0].Trashed)

	_, err = handler.Restore("1", "entry-trashed")
	require.NoError(t, err)
	require.False(t, handler.Vault.Entries.Login[0].Trashed)
}

func TestLoginHandlerTrashEntryNotFound(t *testing.T) {
	t.Parallel()

	handler := newTestLoginHandler()
	handler.Vault = &vaults_domain.VaultPayload{}

	_, err := handler.Trash("1", "does-not-exist")
	require.Error(t, err)
	require.Contains(t, err.Error(), "entry with ID does-not-exist not found")
}

func newTestLoginHandler() *LoginHandler {
	return NewLoginHandler(
		models.DBModel{},
		logger.New(logger.ERROR),
	)
}

func newVaultSession(userID string, entries []vaults_domain.LoginEntry) *models.VaultSession {
	loginEntries := append([]vaults_domain.LoginEntry(nil), entries...)
	return &models.VaultSession{
		UserID: userID,
		Vault: &models.VaultPayload{
			BaseVaultContent: models.BaseVaultContent{
				Entries: models.Entries{
					Login: make([]models.LoginEntry, len(loginEntries)),
				},
			},
		},
	}
}

package handlers

import (
	"testing"

	"github.com/stretchr/testify/require"

	"vault-app/internal/blockchain"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
)

func TestLoginHandlerAddSuccess(t *testing.T) {
	t.Parallel()

	userID := "101"
	sessions := map[string]*models.VaultSession{
		userID: newVaultSession(userID, nil),
	}
	handler := newTestLoginHandler(sessions)

	entry := &models.LoginEntry{
		BaseEntry: models.BaseEntry{
			EntryName: "Primary Account",
		},
		UserName: "user@example.com",
		Password: "super-secret",
	}

	result, err := handler.Add(userID, entry)
	require.NoError(t, err)
	require.NotNil(t, result)

	addedEntry, ok := (*result).(*models.LoginEntry)
	require.True(t, ok)
	require.NotEmpty(t, addedEntry.ID)
	require.Equal(t, "Primary Account", addedEntry.EntryName)

	require.Len(t, sessions[userID].Vault.Entries.Login, 1)
	require.Equal(t, addedEntry.ID, sessions[userID].Vault.Entries.Login[0].ID)
}

func TestLoginHandlerAddNoSession(t *testing.T) {
	t.Parallel()

	handler := newTestLoginHandler(map[string]*models.VaultSession{})
	entry := &models.LoginEntry{}

	result, err := handler.Add("55", entry)
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "no active session for user 55")
}

func TestLoginHandlerAddInvalidType(t *testing.T) {
	t.Parallel()

	userID := "7"
	sessions := map[string]*models.VaultSession{
		userID: newVaultSession(userID, nil),
	}
	handler := newTestLoginHandler(sessions)

	invalidEntry := &models.CardEntry{}

	result, err := handler.Add(userID, invalidEntry)
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "entry does not implement VaultEntry interface")
}

func TestLoginHandlerEditSuccess(t *testing.T) {
	t.Parallel()

	userID := "88"
	existingEntry := models.LoginEntry{
		BaseEntry: models.BaseEntry{
			ID:        "entry-1",
			EntryName: "Old Name",
		},
		UserName: "old@example.com",
	}

	sessions := map[string]*models.VaultSession{
		userID: newVaultSession(userID, []models.LoginEntry{existingEntry}),
	}
	handler := newTestLoginHandler(sessions)

	updated := &models.LoginEntry{
		BaseEntry: models.BaseEntry{
			ID:        "entry-1",
			EntryName: "Updated Name",
		},
		UserName: "new@example.com",
		Password: "new-pass",
	}

	result, err := handler.Edit(userID, updated)
	require.NoError(t, err)
	require.NotNil(t, result)

	modifiedEntry, ok := (*result).(*models.LoginEntry)
	require.True(t, ok)
	require.True(t, modifiedEntry.IsDraft)
	require.Equal(t, "Updated Name", sessions[userID].Vault.Entries.Login[0].EntryName)
	require.Equal(t, "new@example.com", sessions[userID].Vault.Entries.Login[0].UserName)
	require.True(t, sessions[userID].Vault.Entries.Login[0].IsDraft)
}

func TestLoginHandlerEditEntryNotFound(t *testing.T) {
	t.Parallel()

	userID := "90"
	sessions := map[string]*models.VaultSession{
		userID: newVaultSession(userID, nil),
	}
	handler := newTestLoginHandler(sessions)

	entry := &models.LoginEntry{
		BaseEntry: models.BaseEntry{
			ID:        "missing",
			EntryName: "Ghost",
		},
	}

	result, err := handler.Edit(userID, entry)
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "entry with ID missing not found for user 90")
}

func TestLoginHandlerEditInvalidType(t *testing.T) {
	t.Parallel()

	userID := "43"
	sessions := map[string]*models.VaultSession{
		userID: newVaultSession(userID, nil),
	}
	handler := newTestLoginHandler(sessions)

	result, err := handler.Edit(userID, struct{}{})
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "invalid type: expected LoginEntry")
}

func TestLoginHandlerTrashAndRestore(t *testing.T) {
	t.Parallel()

	userID := 	"11"
	entry := models.LoginEntry{
		BaseEntry: models.BaseEntry{
			ID:        "entry-trashed",
			EntryName: "Disposable",
		},
	}

	sessions := map[string]*models.VaultSession{
		userID: newVaultSession(userID, []models.LoginEntry{entry}),
	}
	handler := newTestLoginHandler(sessions)

	err := handler.Trash(userID, "entry-trashed")
	require.NoError(t, err)
	require.True(t, sessions[userID].Vault.Entries.Login[0].Trashed)

	err = handler.Restore(userID, "entry-trashed")
	require.NoError(t, err)
	require.False(t, sessions[userID].Vault.Entries.Login[0].Trashed)
}

func TestLoginHandlerTrashEntryNotFound(t *testing.T) {
	t.Parallel()

	userID := 	"12"
	sessions := map[string]*models.VaultSession{
		userID: newVaultSession(userID, nil),
	}
	handler := newTestLoginHandler(sessions)

	err := handler.Trash(userID, "does-not-exist")
	require.Error(t, err)
	require.Contains(t, err.Error(), "entry with ID does-not-exist not found")
}

func newTestLoginHandler(sessions map[string]*models.VaultSession) *LoginHandler {
	return NewLoginHandler(
		models.DBModel{},
		blockchain.IPFSClient{},
		sessions,
		logger.New(logger.ERROR),
	)
}

func newVaultSession(userID string, entries []models.LoginEntry) *models.VaultSession {
	loginEntries := append([]models.LoginEntry(nil), entries...)
	return &models.VaultSession{
		UserID: userID,
		Vault: &models.VaultPayload{
			BaseVaultContent: models.BaseVaultContent{
				Entries: models.Entries{
					Login: loginEntries,
				},
			},
		},
	}
}

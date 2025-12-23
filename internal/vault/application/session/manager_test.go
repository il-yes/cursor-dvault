package vault_session_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	vault_session "vault-app/internal/vault/application/session"
	vault_domain "vault-app/internal/vault/domain"
)

// ---------- FAKES ----------

type fakeSessionRepo struct {
	saved bool
}

func (f *fakeSessionRepo) CreateSession(*vault_session.Session) error {
	f.saved = true
	return nil
}
func (f *fakeSessionRepo) UpdateSession(*vault_session.Session) error {
	f.saved = true
	return nil
}
func (f *fakeSessionRepo) DeleteSession(string) error {
	f.saved = true
	return nil
}
func (f *fakeSessionRepo) GetSession(string) (*vault_session.Session, error) {
	return nil, nil
}
func (f *fakeSessionRepo) SaveSession(string, *vault_session.Session) error {
	f.saved = true
	return nil
}
func (f *fakeSessionRepo) GetLatestByUserID(string) (*vault_session.Session, error) {
	return nil, nil
}

type fakeVaultRepo struct{}

func (f *fakeVaultRepo) GetLatestByUserID(userID string) (*vault_domain.Vault, error) {
	return nil, vault_domain.ErrVaultNotFound
}
func (f *fakeVaultRepo) CreateVault(v *vault_domain.Vault) error { return nil }
func (f *fakeVaultRepo) GetVault(id string) (*vault_domain.Vault, error) {
	return nil, nil
}
func (f *fakeVaultRepo) UpdateVault(v *vault_domain.Vault) error { return nil }
func (f *fakeVaultRepo) DeleteVault(id string) error             { return nil }
func (f *fakeVaultRepo) SaveVault(v *vault_domain.Vault) error   { return nil }

type fakeLogger struct{}

func (f *fakeLogger) Info(string, ...interface{})  {}
func (f *fakeLogger) Error(string, ...interface{}) {}

// ---------- TESTS ----------

func TestManager_PrepareAndGet(t *testing.T) {
	manager := vault_session.NewManager(
		&fakeSessionRepo{},
		&fakeVaultRepo{},
		&fakeLogger{},
		context.Background(),
		nil,
		nil,	
	)

	userID := "user-1"

	id := manager.Prepare(userID)
	assert.Equal(t, userID, id)

	session, ok := manager.Get(userID)
	assert.True(t, ok)
	assert.Equal(t, userID, session.UserID)
}

func TestManager_AttachVault(t *testing.T) {
	manager := vault_session.NewManager(
		&fakeSessionRepo{},
		&fakeVaultRepo{},
		&fakeLogger{},
		context.Background(),
		nil,
		nil,	
	)

	userID := "user-2"
	manager.Prepare(userID)

	payload := &vault_domain.VaultPayload{}
	runtime := &vault_session.RuntimeContext{}
	lastCID := "cid-123"

	session := manager.AttachVault(userID, payload, runtime, lastCID)

	assert.NotNil(t, session)
	assert.Equal(t, payload, session.Vault)
	assert.Equal(t, runtime, session.Runtime)
	assert.Equal(t, lastCID, session.LastCID)
}

func TestManager_StartSession(t *testing.T) {
	manager := vault_session.NewManager(
		&fakeSessionRepo{},
		&fakeVaultRepo{},
		&fakeLogger{},
		context.Background(),
		nil,
		nil,	
	)

	userID := "user-3"
	payload := vault_domain.VaultPayload{}
	runtime := &vault_session.RuntimeContext{}
	lastCID := "cid-xyz"

	session := manager.StartSession(userID, payload, lastCID, runtime)

	assert.Equal(t, userID, session.UserID)
	assert.Equal(t, lastCID, session.LastCID)
	assert.False(t, session.Dirty)
	assert.NotNil(t, session.Runtime)

	got, ok := manager.Get(userID)
	assert.True(t, ok)
	assert.Equal(t, session, got)
}

func TestManager_MarkDirty(t *testing.T) {
	manager := vault_session.NewManager(
		&fakeSessionRepo{},
		&fakeVaultRepo{},
		&fakeLogger{},
		context.Background(),
		nil,
		nil,	
	)

	manager.NowUTC = func() string { return "now" }

	userID := "user-4"
	manager.Prepare(userID)

	manager.MarkDirty(userID)

	session, _ := manager.Get(userID)
	assert.True(t, session.Dirty)
	assert.Equal(t, "now", session.LastUpdated)
	assert.True(t, manager.IsVaultDirty())
}

func TestManager_Close(t *testing.T) {
	manager := vault_session.NewManager(
		&fakeSessionRepo{},
		&fakeVaultRepo{},
		&fakeLogger{},
		context.Background(),
		nil,
		nil,	
		)

	userID := "user-5"
	manager.Prepare(userID)

	manager.Close(userID)

	_, ok := manager.Get(userID)
	assert.False(t, ok)
}

func TestManager_GetSession_Error(t *testing.T) {
	manager := vault_session.NewManager(
		&fakeSessionRepo{},
		&fakeVaultRepo{},
		&fakeLogger{},
		context.Background(),
		nil,
		nil,	
	)

	_, err := manager.GetSession("missing-user")
	assert.Error(t, err)
}

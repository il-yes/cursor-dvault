package identity_domain_test

import (
	"testing"
	"time"

	identity_domain "vault-app/internal/identity/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDevice_Valid(t *testing.T) {
	vaultID := "vault-123"
	pubKey := "ed25519-public-key-data"
	keyType := identity_domain.DeviceKeyTypeEd25519

	dev, err := identity_domain.NewDevice(vaultID, pubKey, keyType)
	require.NoError(t, err)
	require.NotNil(t, dev)

	assert.NotEmpty(t, dev.ID)
	assert.Equal(t, vaultID, dev.VaultID)
	assert.Equal(t, pubKey, dev.PublicKey)
	assert.Equal(t, keyType, dev.KeyType)
	assert.Equal(t, identity_domain.DeviceStatusActive, dev.Status)
	assert.True(t, dev.IsActive())
	assert.Nil(t, dev.RevokedAt)
	assert.False(t, dev.CreatedAt.IsZero())
	assert.Equal(t, "identity_devices", dev.TableName())
}

func TestNewDevice_MissingVaultID(t *testing.T) {
	dev, err := identity_domain.NewDevice("", "pubkey", identity_domain.DeviceKeyTypeEd25519)
	assert.Nil(t, dev)
	assert.ErrorIs(t, err, identity_domain.ErrDeviceVaultIDRequired)
}

func TestNewDevice_MissingPublicKey(t *testing.T) {
	dev, err := identity_domain.NewDevice("vault-1", "", identity_domain.DeviceKeyTypeEd25519)
	assert.Nil(t, dev)
	assert.ErrorIs(t, err, identity_domain.ErrDevicePublicKeyRequired)
}

func TestNewDevice_MissingKeyType(t *testing.T) {
	dev, err := identity_domain.NewDevice("vault-1", "pubkey", "")
	assert.Nil(t, dev)
	assert.ErrorIs(t, err, identity_domain.ErrDeviceKeyTypeRequired)
}

func TestDevice_RevocationLifecycle(t *testing.T) {
	dev, err := identity_domain.NewDevice("vault-1", "pubkey", identity_domain.DeviceKeyTypeEd25519)
	require.NoError(t, err)
	assert.True(t, dev.IsActive())

	// First revocation must succeed
	err = dev.Revoke()
	require.NoError(t, err)
	assert.False(t, dev.IsActive())
	assert.Equal(t, identity_domain.DeviceStatusRevoked, dev.Status)
	require.NotNil(t, dev.RevokedAt)
	assert.WithinDuration(t, time.Now(), *dev.RevokedAt, 2*time.Second)

	// Second revocation must fail with ErrDeviceAlreadyRevoked
	err = dev.Revoke()
	assert.ErrorIs(t, err, identity_domain.ErrDeviceAlreadyRevoked)
}

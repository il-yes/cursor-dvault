package trustgroup_adapters_test

import (
	"context"
	"testing"

	identity_domain "vault-app/internal/identity/domain"
	identity_persistence "vault-app/internal/identity/infrastructure/persistence"
	trustgroup_adapters "vault-app/internal/trust_group/infrastructure/adapters"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentityDeviceAdapter(t *testing.T) {
	ctx := context.Background()
	devRepo := identity_persistence.NewMemoryDeviceRepository()

	adapter := trustgroup_adapters.NewIdentityDeviceAdapter(devRepo)

	// Device not found
	summary, err := adapter.GetDevice(ctx, "nonexistent")
	require.NoError(t, err)
	assert.Nil(t, summary)

	// Create device in Identity repo
	dev, err := identity_domain.NewDevice("vault-100", "pubkey-100", identity_domain.DeviceKeyTypeEd25519)
	require.NoError(t, err)
	err = devRepo.Save(ctx, dev)
	require.NoError(t, err)

	// Fetch via adapter
	summary, err = adapter.GetDevice(ctx, dev.ID)
	require.NoError(t, err)
	require.NotNil(t, summary)

	assert.Equal(t, dev.ID, summary.ID)
	assert.Equal(t, "vault-100", summary.VaultID)
	assert.Equal(t, "pubkey-100", summary.PublicKey)
	assert.True(t, summary.IsActive)

	// Revoke device and check adapter output
	dev.Revoke()
	err = devRepo.Update(ctx, dev)
	require.NoError(t, err)

	summary, err = adapter.GetDevice(ctx, dev.ID)
	require.NoError(t, err)
	require.NotNil(t, summary)
	assert.False(t, summary.IsActive)
}

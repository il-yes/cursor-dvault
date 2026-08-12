package identity_usecase_test

import (
	"context"
	"sync"
	"testing"

	identity_eventbus "vault-app/internal/identity/application"
	identity_usecase "vault-app/internal/identity/application/usecase"
	identity_domain "vault-app/internal/identity/domain"
	identity_infrastructure_eventbus "vault-app/internal/identity/infrastructure/eventbus"
	identity_persistence "vault-app/internal/identity/infrastructure/persistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeviceUseCases_FullLifecycle(t *testing.T) {
	ctx := context.Background()
	repo := identity_persistence.NewMemoryDeviceRepository()
	bus := identity_infrastructure_eventbus.NewMemoryEventBus()

	createUC := identity_usecase.NewCreateDeviceUseCase(repo, bus)
	getUC := identity_usecase.NewGetDeviceUseCase(repo)
	listUC := identity_usecase.NewListDevicesUseCase(repo)
	revokeUC := identity_usecase.NewRevokeDeviceUseCase(repo, bus)

	var wg sync.WaitGroup
	var publishedCreated identity_eventbus.DeviceCreated
	var publishedRevoked identity_eventbus.DeviceRevoked

	wg.Add(2)
	err := bus.SubscribeToDeviceCreated(func(c context.Context, e identity_eventbus.DeviceCreated) {
		publishedCreated = e
		wg.Done()
	})
	require.NoError(t, err)

	err = bus.SubscribeToDeviceRevoked(func(c context.Context, e identity_eventbus.DeviceRevoked) {
		publishedRevoked = e
		wg.Done()
	})
	require.NoError(t, err)

	// 1. Create Device
	req := identity_usecase.CreateDeviceRequest{
		VaultID:   "vault-99",
		PublicKey: "pubkey-abc-123",
		KeyType:   identity_domain.DeviceKeyTypeEd25519,
	}

	dev, err := createUC.Execute(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, dev)
	assert.Equal(t, "vault-99", dev.VaultID)
	assert.True(t, dev.IsActive())

	// 2. Get Device
	fetched, err := getUC.Execute(ctx, dev.ID)
	require.NoError(t, err)
	assert.Equal(t, dev.ID, fetched.ID)
	assert.Equal(t, dev.PublicKey, fetched.PublicKey)

	// 3. List Devices
	devices, err := listUC.Execute(ctx, "vault-99")
	require.NoError(t, err)
	assert.Len(t, devices, 1)
	assert.Equal(t, dev.ID, devices[0].ID)

	// 4. Revoke Device
	err = revokeUC.Execute(ctx, dev.ID)
	require.NoError(t, err)

	wg.Wait()

	// Verify events received via bus
	assert.Equal(t, dev.ID, publishedCreated.DeviceID)
	assert.Equal(t, "pubkey-abc-123", publishedCreated.PublicKey)
	assert.Equal(t, dev.ID, publishedRevoked.DeviceID)

	// Verify device status updated in repo
	updated, err := getUC.Execute(ctx, dev.ID)
	require.NoError(t, err)
	assert.False(t, updated.IsActive())
	assert.Equal(t, identity_domain.DeviceStatusRevoked, updated.Status)

	// 5. Duplicate revocation fails
	err = revokeUC.Execute(ctx, dev.ID)
	assert.ErrorIs(t, err, identity_domain.ErrDeviceAlreadyRevoked)
}

func TestCreateDeviceUseCase_ValidationErrors(t *testing.T) {
	ctx := context.Background()
	repo := identity_persistence.NewMemoryDeviceRepository()
	createUC := identity_usecase.NewCreateDeviceUseCase(repo, nil)

	_, err := createUC.Execute(ctx, identity_usecase.CreateDeviceRequest{
		VaultID:   "",
		PublicKey: "pk",
		KeyType:   "ed25519",
	})
	assert.ErrorIs(t, err, identity_domain.ErrDeviceVaultIDRequired)

	_, err = createUC.Execute(ctx, identity_usecase.CreateDeviceRequest{
		VaultID:   "v",
		PublicKey: "",
		KeyType:   "ed25519",
	})
	assert.ErrorIs(t, err, identity_domain.ErrDevicePublicKeyRequired)

	_, err = createUC.Execute(ctx, identity_usecase.CreateDeviceRequest{
		VaultID:   "v",
		PublicKey: "pk",
		KeyType:   "",
	})
	assert.ErrorIs(t, err, identity_domain.ErrDeviceKeyTypeRequired)
}

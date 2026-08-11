package identity_persistence_test

import (
	"context"
	"testing"

	identity_domain "vault-app/internal/identity/domain"
	identity_persistence "vault-app/internal/identity/infrastructure/persistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMemoryDeviceRepository(t *testing.T) {
	ctx := context.Background()
	repo := identity_persistence.NewMemoryDeviceRepository()

	dev, err := identity_domain.NewDevice("vault-1", "pub-1", identity_domain.DeviceKeyTypeEd25519)
	require.NoError(t, err)

	err = repo.Save(ctx, dev)
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, dev.ID)
	require.NoError(t, err)
	assert.Equal(t, dev.ID, found.ID)

	list, err := repo.ListByVaultID(ctx, "vault-1")
	require.NoError(t, err)
	assert.Len(t, list, 1)

	dev.Revoke()
	err = repo.Update(ctx, dev)
	require.NoError(t, err)

	found2, err := repo.FindByID(ctx, dev.ID)
	require.NoError(t, err)
	assert.False(t, found2.IsActive())
}

func TestGormDeviceRepository(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&identity_domain.Device{})
	require.NoError(t, err)

	ctx := context.Background()
	repo := identity_persistence.NewGormDeviceRepository(db)

	dev, err := identity_domain.NewDevice("vault-2", "pub-2", identity_domain.DeviceKeyTypeEd25519)
	require.NoError(t, err)

	err = repo.Save(ctx, dev)
	require.NoError(t, err)

	found, err := repo.FindByID(ctx, dev.ID)
	require.NoError(t, err)
	assert.Equal(t, dev.ID, found.ID)
	assert.Equal(t, "pub-2", found.PublicKey)

	list, err := repo.ListByVaultID(ctx, "vault-2")
	require.NoError(t, err)
	assert.Len(t, list, 1)

	dev.Revoke()
	err = repo.Update(ctx, dev)
	require.NoError(t, err)

	found2, err := repo.FindByID(ctx, dev.ID)
	require.NoError(t, err)
	assert.False(t, found2.IsActive())
}

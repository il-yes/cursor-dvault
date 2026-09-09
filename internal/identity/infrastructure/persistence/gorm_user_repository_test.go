package identity_persistence_test

import (
	"context"
	"testing"
	"time"

	identity_domain "vault-app/internal/identity/domain"
	identity_persistence "vault-app/internal/identity/infrastructure/persistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGormUserRepository_Update_DoesNotDuplicateRowOrViolateUniqueEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&identity_domain.User{})
	require.NoError(t, err)

	ctx := context.Background()
	repo := identity_persistence.NewGormUserRepository(db)

	u := identity_domain.NewStandardUser("user-id-123", "testuser@ankhora.test", "hash123")
	err = repo.Save(ctx, u)
	require.NoError(t, err)

	// Verify only 1 user exists
	var count int64
	db.Model(&identity_domain.User{}).Count(&count)
	assert.Equal(t, int64(1), count)

	// Update LastConnectedAt on existing user
	u.LastConnectedAt = time.Now().Add(10 * time.Minute)
	err = repo.Update(ctx, u)
	require.NoError(t, err, "Update on existing user MUST NOT fail with unique constraint error")

	// Verify count is STILL 1
	db.Model(&identity_domain.User{}).Count(&count)
	assert.Equal(t, int64(1), count, "Update MUST NOT create a second identity_users row")

	// Reload user and verify updated LastConnectedAt and unchanged Email / ID
	reloaded, err := repo.FindByEmail(ctx, "testuser@ankhora.test")
	require.NoError(t, err)
	assert.Equal(t, "user-id-123", reloaded.ID)
	assert.Equal(t, "testuser@ankhora.test", reloaded.Email)
	assert.WithinDuration(t, u.LastConnectedAt, reloaded.LastConnectedAt, 2*time.Second)
}

func TestGormUserRepository_Update_UsernameFirstLastNames(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&identity_domain.User{})
	require.NoError(t, err)

	ctx := context.Background()
	repo := identity_persistence.NewGormUserRepository(db)

	u := identity_domain.NewStandardUser("user-id-456", "sovereign@ankhora.test", "hash456")
	err = repo.Save(ctx, u)
	require.NoError(t, err)

	// Update Username, FirstName, LastName
	u.Username = "sovereign_user"
	u.FirstName = "Alice"
	u.LastName = "Smith"
	err = repo.Update(ctx, u)
	require.NoError(t, err)

	// Reload and verify
	reloaded, err := repo.FindByID(ctx, "user-id-456")
	require.NoError(t, err)
	assert.Equal(t, "sovereign_user", reloaded.Username)
	assert.Equal(t, "sovereign_user", reloaded.UserName)
	assert.Equal(t, "Alice", reloaded.FirstName)
	assert.Equal(t, "Smith", reloaded.LastName)
	assert.Equal(t, "sovereign@ankhora.test", reloaded.Email)
}

func TestGormUserRepository_Update_UpsertFallback_And_EmailUntouched(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&identity_domain.User{})
	require.NoError(t, err)

	ctx := context.Background()
	repo := identity_persistence.NewGormUserRepository(db)

	// User initially created outside identity_users (e.g. in legacy users table)
	u := identity_domain.NewStandardUser("60496dae-7df6-4441-93f7-60d54e53c55e", "persistent@ankhora.test", "hash789")
	u.Username = "initial_alias"

	// Call Update directly on user not yet in identity_users (RowsAffected will be 0 initially)
	u.Username = "updated_alias"
	u.FirstName = "Bob"
	u.LastName = "Builder"
	err = repo.Update(ctx, u)
	require.NoError(t, err, "Update on un-indexed user must succeed via upsert fallback")

	// Fresh read from database
	freshRepo := identity_persistence.NewGormUserRepository(db)
	reloaded, err := freshRepo.FindByID(ctx, "60496dae-7df6-4441-93f7-60d54e53c55e")
	require.NoError(t, err)

	assert.Equal(t, "updated_alias", reloaded.Username)
	assert.Equal(t, "updated_alias", reloaded.UserName)
	assert.Equal(t, "Bob", reloaded.FirstName)
	assert.Equal(t, "Builder", reloaded.LastName)
	assert.Equal(t, "persistent@ankhora.test", reloaded.Email, "Email MUST remain untouched")
}

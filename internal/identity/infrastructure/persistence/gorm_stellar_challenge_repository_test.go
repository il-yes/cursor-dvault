package identity_persistence_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	identity_persistence "vault-app/internal/identity/infrastructure/persistence"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&identity_persistence.StellarAuthChallenge{})
	require.NoError(t, err)
	return db
}

func TestGormStellarChallengeRepository_Persistence_And_Lifecycle(t *testing.T) {
	db := setupTestDB(t)
	repo := identity_persistence.NewGormStellarChallengeRepository(db)
	pubKey := "GBV35PVNE77KMVFBK3JS4OXXQPHSVEYEDYNSSKPIFNJZH2EJNC5O4THV"

	// 1. Create Challenge
	record, err := repo.CreateOrReplace(context.Background(), pubKey, "challenge_payload_123", 5*time.Minute)
	require.NoError(t, err)
	require.NotNil(t, record)
	assert.Equal(t, pubKey, record.PublicKey)
	assert.Equal(t, "challenge_payload_123", record.Challenge)

	// 2. Restart app (re-create repo instance on same DB connection)
	repo2 := identity_persistence.NewGormStellarChallengeRepository(db)

	// 3. Atomically Consume
	challengeStr, errConsume := repo2.Consume(context.Background(), pubKey)
	require.NoError(t, errConsume)
	assert.Equal(t, "challenge_payload_123", challengeStr)

	// 4. Replay Attack Attempt (must fail)
	_, errReplay := repo2.Consume(context.Background(), pubKey)
	require.Error(t, errReplay)
	assert.Contains(t, errReplay.Error(), "already consumed")
}

func TestGormStellarChallengeRepository_Replacement(t *testing.T) {
	db := setupTestDB(t)
	repo := identity_persistence.NewGormStellarChallengeRepository(db)
	pubKey := "GBV35PVNE77KMVFBK3JS4OXXQPHSVEYEDYNSSKPIFNJZH2EJNC5O4THV"

	// 1. Create first challenge
	_, err1 := repo.CreateOrReplace(context.Background(), pubKey, "first_challenge", 5*time.Minute)
	require.NoError(t, err1)

	// 2. Create second challenge for same public key (replaces first)
	record2, err2 := repo.CreateOrReplace(context.Background(), pubKey, "second_challenge", 5*time.Minute)
	require.NoError(t, err2)
	assert.Equal(t, "second_challenge", record2.Challenge)

	// 3. Consume should return second challenge
	challengeStr, errConsume := repo.Consume(context.Background(), pubKey)
	require.NoError(t, errConsume)
	assert.Equal(t, "second_challenge", challengeStr)
}

func TestGormStellarChallengeRepository_Expiration(t *testing.T) {
	db := setupTestDB(t)
	repo := identity_persistence.NewGormStellarChallengeRepository(db)
	pubKey := "GBV35PVNE77KMVFBK3JS4OXXQPHSVEYEDYNSSKPIFNJZH2EJNC5O4THV"

	// Create challenge with -1 second TTL (already expired)
	_, err := repo.CreateOrReplace(context.Background(), pubKey, "expired_challenge", -1*time.Second)
	require.NoError(t, err)

	// Consume should fail due to expiry
	_, errConsume := repo.Consume(context.Background(), pubKey)
	require.Error(t, errConsume)
	assert.Contains(t, errConsume.Error(), "expired")
}

func TestGormStellarChallengeRepository_UnknownPublicKey(t *testing.T) {
	db := setupTestDB(t)
	repo := identity_persistence.NewGormStellarChallengeRepository(db)

	_, errConsume := repo.Consume(context.Background(), "UNKNOWN_PUBLIC_KEY")
	require.Error(t, errConsume)
	assert.Contains(t, errConsume.Error(), "already consumed")
}

func TestGormStellarChallengeRepository_AtomicConcurrentConsumption(t *testing.T) {
	db := setupTestDB(t)
	repo := identity_persistence.NewGormStellarChallengeRepository(db)
	pubKey := "GBV35PVNE77KMVFBK3JS4OXXQPHSVEYEDYNSSKPIFNJZH2EJNC5O4THV"

	_, err := repo.CreateOrReplace(context.Background(), pubKey, "concurrent_challenge", 5*time.Minute)
	require.NoError(t, err)

	concurrentRequests := 10
	var wg sync.WaitGroup
	successCount := 0
	failureCount := 0
	var mu sync.Mutex

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cStr, errConsume := repo.Consume(context.Background(), pubKey)
			mu.Lock()
			defer mu.Unlock()
			if errConsume == nil && cStr == "concurrent_challenge" {
				successCount++
			} else {
				failureCount++
			}
		}()
	}

	wg.Wait()

	// Exactly 1 request MUST succeed, remaining 9 MUST fail
	assert.Equal(t, 1, successCount, "Exactly one concurrent request must succeed in consuming the challenge")
	assert.Equal(t, concurrentRequests-1, failureCount, "All other concurrent requests must fail")
}

package stellar_recovery_usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"vault-app/internal/handlers"
	"vault-app/internal/models"
	shared "vault-app/internal/shared/stellar"
	stellar_recovery_usecase "vault-app/internal/stellar_recovery/application/use_case"
	stellar_recovery_domain "vault-app/internal/stellar_recovery/domain"
)

//
// ──────────────────────────────────────────
// M O C K S
// ──────────────────────────────────────────
//

// --- Stellar Port Mock ---
type mockStellarPort struct {
	RecoverPasswordFn func(ctx context.Context, in shared.RecoverPasswordInput) (string, *models.User, error)
}

func (m *mockStellarPort) RecoverPassword(ctx context.Context, in shared.RecoverPasswordInput) (string, *models.User, error) {
	return m.RecoverPasswordFn(ctx, in)
}


// --- Subscription Repo Mock ---
type mockSubscriptionRepo struct {
	GetActiveByUserIDFn func(ctx context.Context, userID string) (*stellar_recovery_domain.Subscription, error)
}

func (m *mockSubscriptionRepo) GetActiveByUserID(ctx context.Context, userID string) (*stellar_recovery_domain.Subscription, error) {
	return m.GetActiveByUserIDFn(ctx, userID)
}

//
// ──────────────────────────────────────────
// T E S T S
// ──────────────────────────────────────────
//

func TestConnectWithStellarUseCase_Execute(t *testing.T) {

	//
	// ❌ CASE 1 — RecoverPassword fails
	//
	t.Run("❌ RecoverPassword fails → return error", func(t *testing.T) {

		uc := stellar_recovery_usecase.ConnectWithStellarUseCase{
			StellarPort: &mockStellarPort{
				RecoverPasswordFn: func(ctx context.Context, in shared.RecoverPasswordInput) (string, *models.User, error) {
					return "", nil, errors.New("bad signature")
				},
			},
			VaultRepo: &mockVaultRepo{
				getByUser: func(ctx context.Context, userID string) (*stellar_recovery_domain.Vault, error) {
					return nil, nil
				},
			},
			SubRepo: &mockSubscriptionRepo{
				GetActiveByUserIDFn: func(ctx context.Context, userID string) (*stellar_recovery_domain.Subscription, error) {
					return nil, nil
				},
			},
		}

		res, err := uc.Execute(context.Background(), handlers.LoginRequest{
			PublicKey: "GKEY123",
			SignedMessage: "...",
			Signature: "...",
		})
		assert.Nil(t, res)
		assert.EqualError(t, err, "bad signature")
	})

	//
	// ✔ CASE 2 — Everything exists
	//
	t.Run("✔ RecoverPassword succeeds, vault + subscription exist", func(t *testing.T) {

		user := &models.User{
			ID: 42,
		}

		vault := &stellar_recovery_domain.Vault{
			ID:            "vault-123",
			StorageUsedGB: 5.5,
		}

		sub := &stellar_recovery_domain.Subscription{
			ID:   "sub-1",
			Tier: "Pro",
		}

		uc := stellar_recovery_usecase.ConnectWithStellarUseCase{
			StellarPort: &mockStellarPort{
				RecoverPasswordFn: func(ctx context.Context, in shared.RecoverPasswordInput) (string, *models.User, error) {
					return "DECRYPTED_PASSWORD", user, nil
				},
			},
			VaultRepo: &mockVaultRepo{
				getByUser: func(ctx context.Context, userID string) (*stellar_recovery_domain.Vault, error) {
					assert.Equal(t, "42", userID)
					return vault, nil
				},
			},
			SubRepo: &mockSubscriptionRepo{
				GetActiveByUserIDFn: func(ctx context.Context, userID string) (*stellar_recovery_domain.Subscription, error) {
					assert.Equal(t, "42", userID)
					return sub, nil
				},
			},
		}

		res, err := uc.Execute(context.Background(), handlers.LoginRequest{
			PublicKey: "GKEY123",
			SignedMessage: "...",
			Signature: "...",
		})

		assert.NoError(t, err)
		assert.NotNil(t, res)

		assert.Equal(t, "DECRYPTED_PASSWORD", res.Password)
		assert.Equal(t, user, res.User)
		assert.Equal(t, vault, res.Vault)
		assert.Equal(t, sub, res.Subscription)
	})

	//
	// ✔ CASE 3 — RecoverPassword OK, but vault/sub missing
	//
	t.Run("✔ RecoverPassword succeeds, but no vault or subscription", func(t *testing.T) {

		user := &models.User{ID: 99}

		uc := stellar_recovery_usecase.ConnectWithStellarUseCase{
			StellarPort: &mockStellarPort{
				RecoverPasswordFn: func(ctx context.Context, in shared.RecoverPasswordInput) (string, *models.User, error) {
					return "PASSWORD99", user, nil
				},
			},
			VaultRepo: &mockVaultRepo{
				getByUser: func(ctx context.Context, userID string) (*stellar_recovery_domain.Vault, error) {
					return nil, nil
				},
			},
			SubRepo: &mockSubscriptionRepo{
				GetActiveByUserIDFn: func(ctx context.Context, userID string) (*stellar_recovery_domain.Subscription, error) {
					return nil, nil
				},
			},
		}

		res, err := uc.Execute(context.Background(), handlers.LoginRequest{
			PublicKey: "GKEY99",
			SignedMessage: "...",
			Signature: "...",
		})

		assert.NoError(t, err)
		assert.NotNil(t, res)

		assert.Equal(t, "PASSWORD99", res.Password)
		assert.Equal(t, user, res.User)
		assert.Nil(t, res.Vault)
		assert.Nil(t, res.Subscription)
	})
}

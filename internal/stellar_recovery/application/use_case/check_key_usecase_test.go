package stellar_recovery_usecase_test

import (
	"context"
	"errors"
	"testing"

	usecase "vault-app/internal/stellar_recovery/application/use_case"
	domain "vault-app/internal/stellar_recovery/domain"
)

//
// ────────────────────────────────────────────────
//   MOCKS
// ────────────────────────────────────────────────
//

func (m *mockUserRepo) GetByStellarPublicKey(ctx context.Context, pub string) (*domain.User, error) {
	if m.getByPub == nil {
		return nil, errors.New("not implemented")
	}
	return m.getByPub(ctx, pub)
}

func (m *mockVaultRepo) GetByUserID(ctx context.Context, userID string) (*domain.Vault, error) {
	if m.getByUser == nil {
		return nil, errors.New("not implemented")
	}
	return m.getByUser(ctx, userID)
}

func (m *mockSubRepo) GetActiveByUserID(ctx context.Context, userID string) (*domain.Subscription, error) {
	if m.getActive == nil {
		return nil, nil // subscription optional
	}
	return m.getActive(ctx, userID)
}

func (m *mockVerifier) ParseSecret(secret string) (string, error) {
	if m.parse == nil {
		return "", errors.New("not implemented")
	}
	return m.parse(secret)
}

//
// ────────────────────────────────────────────────
//   TEST SUITE
// ────────────────────────────────────────────────
//

func TestCheckKeyUseCase_Execute(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		mockSetup func() *usecase.CheckKeyUseCase
		wantExist bool
		wantPub   string
	}{
		{
			name:   "invalid secret → VaultExists=false & pub empty",
			secret: "BAD_SECRET",
			mockSetup: func() *usecase.CheckKeyUseCase {
				return usecase.NewCheckKeyUseCase(
					&mockUserRepo{},
					&mockVaultRepo{},
					&mockSubRepo{},
					&mockVerifier{
						parse: func(s string) (string, error) {
							return "", errors.New("invalid format")
						},
					},
				)
			},
			wantExist: false,
			wantPub:   "",
		},
		{
			name:   "valid secret but user not found → VaultExists=false & pub returned",
			secret: "GOOD_SECRET",
			mockSetup: func() *usecase.CheckKeyUseCase {
				return usecase.NewCheckKeyUseCase(
					&mockUserRepo{
						getByPub: func(ctx context.Context, pub string) (*domain.User, error) {
							return nil, errors.New("not found")
						},
					},
					&mockVaultRepo{},
					&mockSubRepo{},
					&mockVerifier{
						parse: func(s string) (string, error) {
							return "GPUB12345", nil
						},
					},
				)
			},
			wantExist: false,
			wantPub:   "GPUB12345",
		},
		{
			name:   "valid secret + user exists → VaultExists=true",
			secret: "GOOD_SECRET",
			mockSetup: func() *usecase.CheckKeyUseCase {
				return usecase.NewCheckKeyUseCase(
					&mockUserRepo{
						getByPub: func(ctx context.Context, pub string) (*domain.User, error) {
							return &domain.User{ID: "10", StellarPublicKey: pub}, nil
						},
					},
					&mockVaultRepo{
						getByUser: func(ctx context.Context, userID string) (*domain.Vault, error) {
							return &domain.Vault{ID: "VAULT55"}, nil
						},
					},
					&mockSubRepo{
						getActive: func(ctx context.Context, userID string) (*domain.Subscription, error) {
							return &domain.Subscription{ID: "SUB200"}, nil
						},
					},
					&mockVerifier{
						parse: func(s string) (string, error) {
							return "GVALIDPUBKEY", nil
						},
					},
				)
			},
			wantExist: true,
			wantPub:   "GVALIDPUBKEY",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("▶ TEST: %s", tt.name)
			t.Logf("→ Input secret: %s", tt.secret)

			uc := tt.mockSetup()
			res, err := uc.Execute(context.Background(), tt.secret)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res.VaultExists != tt.wantExist {
				t.Errorf("VaultExists = %v; want %v", res.VaultExists, tt.wantExist)
			}

			if res.PublicKey != tt.wantPub {
				t.Errorf("PublicKey = %v; want %v", res.PublicKey, tt.wantPub)
			}

			if res.Vault != nil {
				t.Logf("✓ Vault found: VaultID=%s", res.Vault.ID)
			} else {
				t.Logf("✓ Vault not found")
			}

			if res.Subscription != nil {
				t.Logf("✓ Subscription found: SubscriptionID=%s", res.Subscription.ID)
			} else {
				t.Logf("✓ Subscription not found or optional")
			}
		})
	}
}

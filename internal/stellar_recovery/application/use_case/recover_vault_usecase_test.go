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
//   MOCKS — MATCH YOUR DOMAIN EXACTLY
// ────────────────────────────────────────────────
//

// ---------- UserRepository ----------
type mockUserRepo struct {
	getByPub func(ctx context.Context, pub string) (*domain.User, error)
}



// ---------- VaultRepository ----------
type mockVaultRepo struct {
	getByUser func(ctx context.Context, userID string) (*domain.Vault, error)
}

// ---------- SubscriptionRepository ----------
type mockSubRepo struct {
	getActive func(ctx context.Context, userID string) (*domain.Subscription, error)
}



// ---------- StellarKeyVerifier ----------
type mockVerifier struct {
	parse  func(secret string) (string, error)
	verify func(secret, pub string) error
}


func (m *mockVerifier) VerifyOwnership(secret, pub string) error {
	if m.verify == nil {
		return nil
	}
	return m.verify(secret, pub)
}

// ---------- EventDispatcher ----------
type mockEventDispatcher struct {
	dispatched []any
}

func (m *mockEventDispatcher) Dispatch(evt any) error {
	m.dispatched = append(m.dispatched, evt)
	return nil
}

// ---------- TokenGenerator ----------
type mockTokenGen struct {
	gen func(userID string) (string, error)
}

func (m *mockTokenGen) NewSessionToken(userID string) (string, error) {
	if m.gen == nil {
		return "", nil
	}
	return m.gen(userID)
}

//
// ────────────────────────────────────────────────
//   TEST SUITE
// ────────────────────────────────────────────────
//

func TestRecoverVaultUseCase_Execute(t *testing.T) {

	tests := []struct {
		name      string
		secret    string
		mockSetup func() *usecase.RecoverVaultUseCase
		wantErr   error
	}{
		{
			name:   "invalid secret → ErrInvalidKey",
			secret: "BAD_SECRET",
			mockSetup: func() *usecase.RecoverVaultUseCase {
				return usecase.NewRecoverVaultUseCase(
					&mockUserRepo{},
					&mockVaultRepo{},
					&mockSubRepo{},
					&mockVerifier{parse: func(s string) (string, error) {
						return "", errors.New("invalid")
					}},
					&mockEventDispatcher{},
					&mockTokenGen{},
				)
			},
			wantErr: domain.ErrInvalidKey,
		},
		{
			name:   "user not found → ErrUserNotFound",
			secret: "GOOD_SECRET",
			mockSetup: func() *usecase.RecoverVaultUseCase {
				return usecase.NewRecoverVaultUseCase(
					&mockUserRepo{getByPub: func(ctx context.Context, pub string) (*domain.User, error) {
						return nil, errors.New("not found")
					}},
					&mockVaultRepo{},
					&mockSubRepo{},
					&mockVerifier{parse: func(s string) (string, error) { return "GPUB", nil }},
					&mockEventDispatcher{},
					&mockTokenGen{},
				)
			},
			wantErr: domain.ErrUserNotFound,
		},
		{
			name:   "ownership verification fails → ErrKeyVerification",
			secret: "GOOD_SECRET",
			mockSetup: func() *usecase.RecoverVaultUseCase {
				return usecase.NewRecoverVaultUseCase(
					&mockUserRepo{getByPub: func(ctx context.Context, pub string) (*domain.User, error) {
						return &domain.User{ID: "10", StellarPublicKey: "GPUB"}, nil
					}},
					&mockVaultRepo{},
					&mockSubRepo{},
					&mockVerifier{
						parse:  func(s string) (string, error) { return "GPUB", nil },
						verify: func(_, _ string) error { return errors.New("no match") },
					},
					&mockEventDispatcher{},
					&mockTokenGen{},
				)
			},
			wantErr: domain.ErrKeyVerification,
		},
		{
			name:   "vault not found → ErrVaultNotFound",
			secret: "GOOD_SECRET",
			mockSetup: func() *usecase.RecoverVaultUseCase {
				return usecase.NewRecoverVaultUseCase(
					&mockUserRepo{getByPub: func(ctx context.Context, pub string) (*domain.User, error) {
						return &domain.User{ID: "10", StellarPublicKey: "GPUB"}, nil
					}},
					&mockVaultRepo{getByUser: func(ctx context.Context, userID string) (*domain.Vault, error) {
						return nil, errors.New("no vault")
					}},
					&mockSubRepo{},
					&mockVerifier{parse: func(s string) (string, error) { return "GPUB", nil }},
					&mockEventDispatcher{},
					&mockTokenGen{},
				)
			},
			wantErr: domain.ErrVaultNotFound,
		},
		{
			name:   "success → RecoveredVault returned",
			secret: "GOOD_SECRET",
			mockSetup: func() *usecase.RecoverVaultUseCase {

				return usecase.NewRecoverVaultUseCase(
					&mockUserRepo{getByPub: func(ctx context.Context, pub string) (*domain.User, error) {
						return &domain.User{ID: "10", StellarPublicKey: "GPUB"}, nil
					}},
					&mockVaultRepo{getByUser: func(ctx context.Context, userID string) (*domain.Vault, error) {
						return &domain.Vault{ID: "VAULT1"}, nil
					}},
					&mockSubRepo{getActive: func(ctx context.Context, userID string) (*domain.Subscription, error) {
						return &domain.Subscription{ID: "SUB1"}, nil
					}},
					&mockVerifier{
						parse:  func(s string) (string, error) { return "GPUB", nil },
						verify: func(_, _ string) error { return nil },
					},
					&mockEventDispatcher{},
					&mockTokenGen{gen: func(userID string) (string, error) {
						return "SESSION_123", nil
					}},
				)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			uc := tt.mockSetup()
			res, err := uc.Execute(context.Background(), tt.secret)

			// Error case
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}

			// Success case
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res == nil {
				t.Fatalf("expected non-nil result")
			}
			if res.User == nil || res.Vault == nil {
				t.Fatalf("expected user & vault")
			}
			if res.SessionToken == "" {
				t.Fatalf("missing session token")
			}
		})
	}
}

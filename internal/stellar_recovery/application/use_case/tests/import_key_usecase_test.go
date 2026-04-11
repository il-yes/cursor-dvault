package stellar_recovery_usecase_test

import (
	"context"
	"errors"
	"testing"

	stellar_recovery_usecase "vault-app/internal/stellar_recovery/application/use_case"
	stellar_recovery_domain "vault-app/internal/stellar_recovery/domain"
)

//
// ────────────────────────────────────────────────
//   MOCKS
// ────────────────────────────────────────────────
//

type mockKeyValidator struct {
	parse func(secret string) (string, error)
}

func (m *mockKeyValidator) ParseSecret(secret string) (string, error) {
	return m.parse(secret)
}

func (m *mockKeyValidator) VerifyOwnership(publicKey string, secret string) error {
	// Not used in ImportKeyUseCase tests, so just return nil
	return nil
}

//
// ────────────────────────────────────────────────
//   TESTS
// ────────────────────────────────────────────────
//

func TestImportKeyUseCase_Execute(t *testing.T) {
	tests := []struct {
		name       string
		secret     string
		mockSetup  func() *stellar_recovery_usecase.ImportKeyUseCase
		wantErr    error
		wantPub    string
		wantCreate bool
	}{
		{
			name:   "invalid secret → ErrInvalidKey",
			secret: "BAD_SECRET",
			mockSetup: func() *stellar_recovery_usecase.ImportKeyUseCase {
				return &stellar_recovery_usecase.ImportKeyUseCase{
					UserRepo: &mockUserRepo{},
					KeyService: &mockKeyValidator{
						parse: func(s string) (string, error) {
							return "", errors.New("invalid secret format")
						},
					},
				}
			},
			wantErr:    stellar_recovery_domain.ErrInvalidKey,
			wantPub:    "",
			wantCreate: false,
		},
		{
			name:   "secret valid but key already in use → ErrKeyAlreadyUsed",
			secret: "GOOD_SECRET",
			mockSetup: func() *stellar_recovery_usecase.ImportKeyUseCase {
				return &stellar_recovery_usecase.ImportKeyUseCase{
					UserRepo: &mockUserRepo{
						getByPub: func(ctx context.Context, pub string) (*stellar_recovery_domain.User, error) {
							return &stellar_recovery_domain.User{ID: "55"}, nil
						},
					},
					KeyService: &mockKeyValidator{
						parse: func(s string) (string, error) {
							return "GPUB12345", nil
						},
					},
				}
			},
			wantErr:    stellar_recovery_domain.ErrKeyAlreadyUsed,
			wantPub:    "",
			wantCreate: false,
		},
		{
			name:   "valid secret, unused → success",
			secret: "GOOD_SECRET",
			mockSetup: func() *stellar_recovery_usecase.ImportKeyUseCase {
				return &stellar_recovery_usecase.ImportKeyUseCase{
					UserRepo: &mockUserRepo{
						getByPub: func(ctx context.Context, pub string) (*stellar_recovery_domain.User, error) {
							return nil, errors.New("not found")
						},
					},
					KeyService: &mockKeyValidator{
						parse: func(s string) (string, error) {
							return "GVALIDPUBLICKEY", nil
						},
					},
				}
			},
			wantErr:    nil,
			wantPub:    "GVALIDPUBLICKEY",
			wantCreate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Logf("▶ TEST: %s", tt.name)
			t.Logf("→ Input secret: %s", tt.secret)

			uc := tt.mockSetup()

			res, err := uc.Execute(context.Background(), tt.secret)

			// Check error case
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error '%v', got '%v'", tt.wantErr, err)
				}
				t.Logf("✓ Received expected error: %v", err)
				return
			}

			// Should not return an error
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check success response
			if res.StellarPublic != tt.wantPub {
				t.Errorf("PublicKey = %v; want %v", res.StellarPublic, tt.wantPub)
			}

			if res.CanCreate != tt.wantCreate {
				t.Errorf("CanCreate = %v; want %v", res.CanCreate, tt.wantCreate)
			}

			t.Logf("✓ Success: ImportedKey → Pub=%s CanCreate=%v", res.StellarPublic, res.CanCreate)
		})
	}
}

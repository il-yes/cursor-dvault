package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	app_config_domain "vault-app/internal/config/domain"
	app_config_ui "vault-app/internal/config/ui"
	"vault-app/internal/tracecore"
	vault_ui "vault-app/internal/vault/ui"
)

func setupTestAppWithServer(t *testing.T, handler http.HandlerFunc) (*App, *httptest.Server, *app_config_domain.UserConfig) {
	ts := httptest.NewServer(handler)

	tc := tracecore.NewTracecoreClient(ts.URL, "test_token", ts.URL, ts.URL)

	app := &App{
		Vault: &vault_ui.VaultHandler{
			TracecoreClient: tc,
		},
	}

	userCfg := &app_config_domain.UserConfig{
		ID: "user_test_123",
		StellarAccount: app_config_domain.StellarAccountConfig{
			PublicKey:  "GBV35PVNE77KMVFBK3JS4OXXQPHSVEYEDYNSSKPIFNJZH2EJNC5O4THV",
			PrivateKey: "SADBPTCJHNVQ4KUZEDV4NDWQKGTOCYEKITLG4E2DZWKXJBBUTCVM7EYK",
		},
	}

	return app, ts, userCfg
}

type stubUserConfigRepo struct {
	userCfg *app_config_domain.UserConfig
}

func (s *stubUserConfigRepo) CreateUserConfig(userConfig *app_config_domain.UserConfig) error {
	return nil
}
func (s *stubUserConfigRepo) GetUserConfig(id string) (*app_config_domain.UserConfig, error) {
	if s.userCfg != nil {
		return s.userCfg, nil
	}
	return nil, errors.New("user config not found")
}
func (s *stubUserConfigRepo) UpdateUserConfig(userConfig *app_config_domain.UserConfig) error {
	return nil
}
func (s *stubUserConfigRepo) DeleteUserConfig(id string) error {
	return nil
}

// 1. New delegation -> Success (200 / 201)
func TestConnectVault_NewDelegation_Success(t *testing.T) {
	app, ts, userCfg := setupTestAppWithServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/identity/challenge":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(tracecore.VaultChallengeResponse{
				ChallengeID:    "chal_123",
				SigningPayload: "payload_to_sign",
			})
		case "/identity/", "/identity":
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(tracecore.VaultRegisterResponse{
				VaultID:      "vault_123",
				DelegationID: "del_123",
				Status:       "active",
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer ts.Close()

	app.AppConfigHandler = &app_config_ui.AppConfigHandler{
		UserConfigRepository: &stubUserConfigRepo{userCfg: userCfg},
	}

	err := app.ConnectVault("user_test_123", "vault_123")
	assert.NoError(t, err)
}

// 2. Existing active delegation -> Idempotent Success (409 active delegation already exists)
func TestConnectVault_ExistingActiveDelegation_IdempotentSuccess(t *testing.T) {
	app, ts, userCfg := setupTestAppWithServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/identity/challenge":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(tracecore.VaultChallengeResponse{
				ChallengeID:    "chal_123",
				SigningPayload: "payload_to_sign",
			})
		case "/identity/", "/identity":
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"message": "active delegation already exists"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer ts.Close()

	app.AppConfigHandler = &app_config_ui.AppConfigHandler{
		UserConfigRepository: &stubUserConfigRepo{userCfg: userCfg},
	}

	err := app.ConnectVault("user_test_123", "vault_123")
	assert.NoError(t, err, "ConnectVault MUST succeed idempotently when delegation already exists")
}

// 3. 409 Conflict with another reason -> Failure
func TestConnectVault_ConflictOtherReason_Failure(t *testing.T) {
	app, ts, userCfg := setupTestAppWithServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/identity/challenge":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(tracecore.VaultChallengeResponse{
				ChallengeID:    "chal_123",
				SigningPayload: "payload_to_sign",
			})
		case "/identity/", "/identity":
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"message": "vault_id already bound to another organization"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer ts.Close()

	app.AppConfigHandler = &app_config_ui.AppConfigHandler{
		UserConfigRepository: &stubUserConfigRepo{userCfg: userCfg},
	}

	err := app.ConnectVault("user_test_123", "vault_123")
	require.Error(t, err)
	assert.False(t, errors.Is(err, tracecore.ErrDelegationAlreadyExists))
	assert.True(t, stringsContains(err.Error(), "vault identity registration failed"))
}

// 4. 401 Unauthorized / 403 Forbidden -> Failure
func TestConnectVault_UnauthorizedOrForbidden_Failure(t *testing.T) {
	app, ts, userCfg := setupTestAppWithServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/identity/challenge":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(tracecore.VaultChallengeResponse{
				ChallengeID:    "chal_123",
				SigningPayload: "payload_to_sign",
			})
		case "/identity/", "/identity":
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message": "unauthorized"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer ts.Close()

	app.AppConfigHandler = &app_config_ui.AppConfigHandler{
		UserConfigRepository: &stubUserConfigRepo{userCfg: userCfg},
	}

	err := app.ConnectVault("user_test_123", "vault_123")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, tracecore.ErrCloudUnauthorized) || stringsContains(err.Error(), "cloud authentication required"))
}

// 5. 500 Server Error -> Failure
func TestConnectVault_ServerError_Failure(t *testing.T) {
	app, ts, userCfg := setupTestAppWithServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/identity/challenge":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(tracecore.VaultChallengeResponse{
				ChallengeID:    "chal_123",
				SigningPayload: "payload_to_sign",
			})
		case "/identity/", "/identity":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message": "internal server error"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer ts.Close()

	app.AppConfigHandler = &app_config_ui.AppConfigHandler{
		UserConfigRepository: &stubUserConfigRepo{userCfg: userCfg},
	}

	err := app.ConnectVault("user_test_123", "vault_123")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, tracecore.ErrCloudServerError) || stringsContains(err.Error(), "cloud server error"))
}

// 6. Repeated ConnectVault invocations for registered vault -> Safe & Idempotent
func TestConnectVault_MultipleInvocations_Idempotent(t *testing.T) {
	registrationCount := 0
	app, ts, userCfg := setupTestAppWithServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/identity/challenge":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(tracecore.VaultChallengeResponse{
				ChallengeID:    "chal_123",
				SigningPayload: "payload_to_sign",
			})
		case "/identity/", "/identity":
			registrationCount++
			if registrationCount == 1 {
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(tracecore.VaultRegisterResponse{
					VaultID:      "vault_123",
					DelegationID: "del_123",
					Status:       "active",
				})
			} else {
				w.WriteHeader(http.StatusConflict)
				_, _ = w.Write([]byte(`{"message": "active delegation already exists"}`))
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer ts.Close()

	app.AppConfigHandler = &app_config_ui.AppConfigHandler{
		UserConfigRepository: &stubUserConfigRepo{userCfg: userCfg},
	}

	// First call -> 201 Created
	err1 := app.ConnectVault("user_test_123", "vault_123")
	require.NoError(t, err1)

	// Second call -> 409 Conflict (active delegation already exists) -> Idempotent Success
	err2 := app.ConnectVault("user_test_123", "vault_123")
	require.NoError(t, err2, "Repeated ConnectVault call for registered vault MUST succeed idempotently")

	assert.Equal(t, 2, registrationCount)
}

func stringsContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && findSubstr(s, substr)))
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

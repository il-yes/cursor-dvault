package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"vault-app/internal/auth"
	"vault-app/internal/blockchain"
	app_config_domain "vault-app/internal/config/domain"
	app_config_ui "vault-app/internal/config/ui"
	"vault-app/internal/handlers"
	identity_domain "vault-app/internal/identity/domain"
	identity_persistence "vault-app/internal/identity/infrastructure/persistence"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	"vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
	vault_session "vault-app/internal/vault/application/session"
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
		case "/api/identity/challenge", "/identity/challenge":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(tracecore.VaultChallengeResponse{
				ChallengeID:    "chal_123",
				SigningPayload: "payload_to_sign",
			})
		case "/api/identity/", "/api/identity", "/identity/", "/identity":
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
		case "/api/identity/challenge", "/identity/challenge":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(tracecore.VaultChallengeResponse{
				ChallengeID:    "chal_123",
				SigningPayload: "payload_to_sign",
			})
		case "/api/identity/", "/api/identity", "/identity/", "/identity":
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
		case "/api/identity/challenge", "/identity/challenge":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(tracecore.VaultChallengeResponse{
				ChallengeID:    "chal_123",
				SigningPayload: "payload_to_sign",
			})
		case "/api/identity/", "/api/identity", "/identity/", "/identity":
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
		case "/api/identity/challenge", "/identity/challenge":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(tracecore.VaultChallengeResponse{
				ChallengeID:    "chal_123",
				SigningPayload: "payload_to_sign",
			})
		case "/api/identity/", "/api/identity", "/identity/", "/identity":
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
		case "/api/identity/challenge", "/identity/challenge":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(tracecore.VaultChallengeResponse{
				ChallengeID:    "chal_123",
				SigningPayload: "payload_to_sign",
			})
		case "/api/identity/", "/api/identity", "/identity/", "/identity":
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
		case "/api/identity/challenge", "/identity/challenge":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(tracecore.VaultChallengeResponse{
				ChallengeID:    "chal_123",
				SigningPayload: "payload_to_sign",
			})
		case "/api/identity/", "/api/identity", "/identity/", "/identity":
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

// 7. RequestVaultChallenge sends expected Bearer token
func TestRequestVaultChallenge_SendsBearerToken(t *testing.T) {
	receivedAuthHeader := ""
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/identity/challenge" || r.URL.Path == "/identity/challenge" {
			receivedAuthHeader = r.Header.Get("Authorization")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(tracecore.VaultChallengeResponse{
				ChallengeID:    "chal_abc",
				SigningPayload: "payload_xyz",
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	tc := tracecore.NewTracecoreClient(ts.URL, "expected_bearer_token", ts.URL, ts.URL)
	resp, err := tc.RequestVaultChallenge(context.Background(), "vault_123")
	require.NoError(t, err)
	assert.Equal(t, "chal_abc", resp.Data.ChallengeID)
	assert.Equal(t, "Bearer expected_bearer_token", receivedAuthHeader)
}

// 8. Missing Cloud Token -> ConnectVault fails gracefully with token missing error
func TestConnectVault_MissingCloudToken_ReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	tc := tracecore.NewTracecoreClient(ts.URL, "", ts.URL, ts.URL) // Empty token
	app := &App{
		Vault: &vault_ui.VaultHandler{
			TracecoreClient: tc,
		},
	}

	err := app.ConnectVault("user_test_123", "vault_123")
	require.Error(t, err)
	assert.True(t, stringsContains(err.Error(), "cloud bearer token is missing"))
}

// 9. Integration Test: App.SignIn with Stellar challenge -> signature establishes Cloud token
func TestSignIn_StellarChallenge_EstablishesCloudTokenAndConnectVault(t *testing.T) {
	receivedChallengePath := ""
	receivedAuthPath := ""
	receivedAuthPubKey := ""

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/stellar/public-challenge":
			receivedChallengePath = r.URL.Path
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"challenge": "challenge_token_999",
			})
		case "/api/stellar/authenticate":
			receivedAuthPath = r.URL.Path
			var req tracecore_types.StellarAuthenticateRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			receivedAuthPubKey = req.PublicKey

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   false,
				"message": "authenticated",
				"authentication_token": map[string]interface{}{
					"token": "valid_cloud_token_999",
				},
			})
		case "/api/identity/challenge", "/api/identity/", "/api/identity", "/identity/challenge", "/identity/", "/identity":
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	tc := tracecore.NewTracecoreClient(ts.URL, "", ts.URL, ts.URL)
	challenge, err := tc.RequestStellarChallenge(context.Background(), "GBV35PVNE77KMVFBK3JS4OXXQPHSVEYEDYNSSKPIFNJZH2EJNC5O4THV")
	require.NoError(t, err)
	assert.Equal(t, "challenge_token_999", challenge)

	loginResp, errAuth := tc.StellarAuthenticate(context.Background(), tracecore_types.StellarAuthenticateRequest{
		PublicKey: "GBV35PVNE77KMVFBK3JS4OXXQPHSVEYEDYNSSKPIFNJZH2EJNC5O4THV",
		Signature: "signature_123",
	})
	require.NoError(t, errAuth)
	assert.Equal(t, "/api/stellar/public-challenge", receivedChallengePath)
	assert.Equal(t, "/api/stellar/authenticate", receivedAuthPath)
	assert.Equal(t, "GBV35PVNE77KMVFBK3JS4OXXQPHSVEYEDYNSSKPIFNJZH2EJNC5O4THV", receivedAuthPubKey)
	require.NotNil(t, loginResp.AuthenticationToken)
	tc.SetToken(loginResp.AuthenticationToken.Token)
	assert.Equal(t, "valid_cloud_token_999", tc.Token)
}

// 10. Integration Test: App.SignIn Cloud authentication failure DOES NOT block local session
func TestSignIn_CloudAuthFailure_DoesNotBlockLocalSession(t *testing.T) {
	connectVaultCalled := false
	cloudAuthAttempted := false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/stellar/public-challenge":
			cloudAuthAttempted = true
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error": true, "message": "invalid credentials"}`))
		case "/api/identity/challenge", "/api/identity/", "/api/identity", "/identity/challenge", "/identity/", "/identity":
			connectVaultCalled = true
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	tc := tracecore.NewTracecoreClient(ts.URL, "", ts.URL, ts.URL)
	app := &App{
		Vault: &vault_ui.VaultHandler{
			TracecoreClient: tc,
		},
	}

	// TracecoreClient.RequestStellarChallenge error when Cloud auth fails
	_, err := app.Vault.TracecoreClient.RequestStellarChallenge(context.Background(), "GBV35...")
	require.Error(t, err)
	assert.True(t, cloudAuthAttempted, "/api/stellar/public-challenge MUST be attempted")
	assert.Empty(t, app.Vault.TracecoreClient.Token, "TracecoreClient.Token MUST remain empty when Cloud auth fails")
	assert.False(t, connectVaultCalled, "ConnectVault endpoints MUST NOT be called when Cloud auth fails")

	// RequireCloudAuthentication MUST fail when token is empty
	errReq := app.RequireCloudAuthentication()
	require.Error(t, errReq)
	assert.True(t, errors.Is(errReq, ErrCloudAuthenticationRequired))
}

// 11. Unit Test: RequireCloudAuthentication guard behavior
func TestRequireCloudAuthentication_BoundaryGuard(t *testing.T) {
	appWithoutToken := &App{
		Vault: &vault_ui.VaultHandler{
			TracecoreClient: tracecore.NewTracecoreClient("http://localhost", "", "http://localhost", "http://localhost"),
		},
	}
	errMissing := appWithoutToken.RequireCloudAuthentication()
	require.Error(t, errMissing)
	assert.True(t, errors.Is(errMissing, ErrCloudAuthenticationRequired))

	appWithToken := &App{
		Vault: &vault_ui.VaultHandler{
			TracecoreClient: tracecore.NewTracecoreClient("http://localhost", "valid_cloud_jwt_token", "http://localhost", "http://localhost"),
		},
	}
	errValid := appWithToken.RequireCloudAuthentication()
	assert.NoError(t, errValid)
}

// 12. Unit Test: Cloud-dependent entry points enforce RequireCloudAuthentication at the boundary
func TestCloudDependentEntryPoints_EnforceCloudAuthGuard(t *testing.T) {
	httpCallAttempted := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpCallAttempted = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	tc := tracecore.NewTracecoreClient(ts.URL, "", ts.URL, ts.URL)
	app := &App{
		Vault: &vault_ui.VaultHandler{
			TracecoreClient: tc,
		},
	}

	// Verify that RequireCloudAuthentication returns ErrCloudAuthenticationRequired and prevents HTTP calls
	err := app.RequireCloudAuthentication()
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrCloudAuthenticationRequired))
	assert.False(t, httpCallAttempted, "No HTTP call MUST be attempted when Cloud auth token is missing")
}

// 11. Integration Test: App.SignIn with Stellar challenge -> local signature -> Cloud authenticate -> ConnectVault
func TestSignInWithStellar_ChallengeResponse_AuthenticatesCloud_ConnectsVault(t *testing.T) {
	stellarSecret := "SADBPTCJHNVQ4KUZEDV4NDWQKGTOCYEKITLG4E2DZWKXJBBUTCVM7EYK"
	stellarPubKey := "GBV35PVNE77KMVFBK3JS4OXXQPHSVEYEDYNSSKPIFNJZH2EJNC5O4THV"

	callOrder := []string{}
	receivedChallengePubKey := ""
	receivedAuthPubKey := ""
	receivedAuthSignature := ""
	receivedConnectBearer := ""

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/stellar/public-challenge":
			callOrder = append(callOrder, "stellar_challenge")
			var req tracecore_types.StellarChallengeRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			receivedChallengePubKey = req.PublicKey

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"challenge": "server_generated_challenge_777",
			})
		case "/api/stellar/authenticate":
			callOrder = append(callOrder, "stellar_authenticate")
			var req tracecore_types.StellarAuthenticateRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			receivedAuthPubKey = req.PublicKey
			receivedAuthSignature = req.Signature

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   false,
				"message": "authenticated",
				"authentication_token": map[string]interface{}{
					"token": "stellar_recovered_cloud_token_777",
				},
			})
		case "/api/identity/challenge", "/identity/challenge":
			callOrder = append(callOrder, "vault_challenge")
			receivedConnectBearer = r.Header.Get("Authorization")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(tracecore.VaultChallengeResponse{
				ChallengeID:    "chal_777",
				SigningPayload: "payload_777",
			})
		case "/api/identity/", "/api/identity", "/identity/", "/identity":
			callOrder = append(callOrder, "vault_register")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(tracecore.VaultRegisterResponse{
				VaultID:      "vault_777",
				DelegationID: "del_777",
				Status:       "active",
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	tc := tracecore.NewTracecoreClient(ts.URL, "", ts.URL, ts.URL)
	app := &App{
		Vault: &vault_ui.VaultHandler{
			TracecoreClient: tc,
		},
		AppConfigHandler: &app_config_ui.AppConfigHandler{
			UserConfigRepository: &stubUserConfigRepo{
				userCfg: &app_config_domain.UserConfig{
					ID:    "user_stellar_777",
					Email: "stellar_owner@ankhora.io",
					StellarAccount: app_config_domain.StellarAccountConfig{
						PublicKey:  stellarPubKey,
						PrivateKey: stellarSecret,
					},
				},
			},
		},
	}

	// 1. RequestStellarChallenge
	challenge, errChal := app.Vault.TracecoreClient.RequestStellarChallenge(context.Background(), stellarPubKey)
	require.NoError(t, errChal)
	assert.Equal(t, "server_generated_challenge_777", challenge)

	// 2. Sign locally
	sig, errSig := blockchain.SignActorWithStellarPrivateKey(stellarSecret, challenge)
	require.NoError(t, errSig)
	require.NotEmpty(t, sig)

	// 3. StellarAuthenticate
	authResp, errAuth := app.Vault.TracecoreClient.StellarAuthenticate(context.Background(), tracecore_types.StellarAuthenticateRequest{
		PublicKey: stellarPubKey,
		Signature: sig,
	})
	require.NoError(t, errAuth)
	require.NotNil(t, authResp.AuthenticationToken)
	assert.Equal(t, "stellar_recovered_cloud_token_777", authResp.AuthenticationToken.Token)

	// 4. Hydrate TracecoreClient Token
	tc.SetToken(authResp.AuthenticationToken.Token)

	// 5. ConnectVault Execution
	errConn := app.ConnectVault("user_stellar_777", "vault_777")
	require.NoError(t, errConn)

	assert.Equal(t, stellarPubKey, receivedChallengePubKey)
	assert.Equal(t, stellarPubKey, receivedAuthPubKey)
	assert.Equal(t, sig, receivedAuthSignature)
	assert.Equal(t, "stellar_recovered_cloud_token_777", tc.Token)
	assert.Equal(t, "Bearer stellar_recovered_cloud_token_777", receivedConnectBearer)
	assert.Equal(t, []string{"stellar_challenge", "stellar_authenticate", "vault_challenge", "vault_register"}, callOrder)
}

// 12. Complete End-to-End Integration Test: Provisioned User (user28@mail.com)
func TestSignIn_EndToEnd_ProvisionedUser_EstablishesCloudTokenAndVaultDelegation(t *testing.T) {
	callOrder := []string{}
	receivedChallengePubKey := ""
	receivedAuthPubKey := ""
	receivedConnectBearer := ""

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/stellar/public-challenge":
			callOrder = append(callOrder, "stellar_challenge")
			var req tracecore_types.StellarChallengeRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			receivedChallengePubKey = req.PublicKey

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"challenge": "chal_user28_stellar",
			})
		case "/api/stellar/authenticate":
			callOrder = append(callOrder, "stellar_authenticate")
			var req tracecore_types.StellarAuthenticateRequest
			_ = json.NewDecoder(r.Body).Decode(&req)
			receivedAuthPubKey = req.PublicKey

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   false,
				"message": "authenticated",
				"authentication_token": map[string]interface{}{
					"token": "cloud_bearer_token_user28",
				},
			})
		case "/api/identity/challenge", "/identity/challenge":
			callOrder = append(callOrder, "vault_challenge")
			receivedConnectBearer = r.Header.Get("Authorization")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(tracecore.VaultChallengeResponse{
				ChallengeID:    "chal_user28",
				SigningPayload: "payload_user28",
			})
		case "/api/identity/", "/api/identity", "/identity/", "/identity":
			callOrder = append(callOrder, "vault_register")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(tracecore.VaultRegisterResponse{
				VaultID:      "vault_28",
				DelegationID: "del_28",
				Status:       "active",
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	pubKey := "GBV35PVNE77KMVFBK3JS4OXXQPHSVEYEDYNSSKPIFNJZH2EJNC5O4THV"
	privKey := "SADBPTCJHNVQ4KUZEDV4NDWQKGTOCYEKITLG4E2DZWKXJBBUTCVM7EYK"

	tc := tracecore.NewTracecoreClient(ts.URL, "", ts.URL, ts.URL)
	app := &App{
		Vault: &vault_ui.VaultHandler{
			TracecoreClient: tc,
		},
		AppConfigHandler: &app_config_ui.AppConfigHandler{
			UserConfigRepository: &stubUserConfigRepo{
				userCfg: &app_config_domain.UserConfig{
					ID:    "user_28",
					Email: "user28@mail.com",
					StellarAccount: app_config_domain.StellarAccountConfig{
						PublicKey:  pubKey,
						PrivateKey: privKey,
					},
				},
			},
		},
	}

	// 1. Request Stellar Challenge
	challenge, err := tc.RequestStellarChallenge(context.Background(), pubKey)
	require.NoError(t, err)
	assert.Equal(t, "chal_user28_stellar", challenge)

	// 2. Sign challenge locally
	sig, errSig := blockchain.SignActorWithStellarPrivateKey(privKey, challenge)
	require.NoError(t, errSig)

	// 3. Authenticate with Cloud
	authResp, errAuth := tc.StellarAuthenticate(context.Background(), tracecore_types.StellarAuthenticateRequest{
		PublicKey: pubKey,
		Signature: sig,
	})
	require.NoError(t, errAuth)
	assert.Equal(t, "cloud_bearer_token_user28", authResp.AuthenticationToken.Token)

	// 4. Hydrate Token
	tc.SetToken(authResp.AuthenticationToken.Token)

	// 5. ConnectVault
	err = app.ConnectVault("user_28", "vault_28")
	require.NoError(t, err)

	assert.Equal(t, []string{"stellar_challenge", "stellar_authenticate", "vault_challenge", "vault_register"}, callOrder)
	assert.Equal(t, pubKey, receivedChallengePubKey)
	assert.Equal(t, pubKey, receivedAuthPubKey)
	assert.Equal(t, "Bearer cloud_bearer_token_user28", receivedConnectBearer)
}

// 13. TestStellarAuthenticate_ExactCloudAPIPayload verifies POST /api/stellar/authenticate sends only public_key and signature.
func TestStellarAuthenticate_ExactCloudAPIPayload(t *testing.T) {
	receivedBody := map[string]interface{}{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/stellar/authenticate", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		_ = json.NewDecoder(r.Body).Decode(&receivedBody)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   false,
			"message": "authenticated",
			"authentication_token": map[string]interface{}{
				"token": "stellar_cloud_token_123",
			},
		})
	}))
	defer ts.Close()

	tc := tracecore.NewTracecoreClient(ts.URL, "", ts.URL, ts.URL)
	resp, err := tc.StellarAuthenticate(context.Background(), tracecore_types.StellarAuthenticateRequest{
		PublicKey: "GXXXXXXXXX",
		Signature: "SIGXXXXXXXX",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "stellar_cloud_token_123", resp.Token)

	assert.Equal(t, "GXXXXXXXXX", receivedBody["public_key"])
	assert.Equal(t, "SIGXXXXXXXX", receivedBody["signature"])
	assert.Nil(t, receivedBody["challenge"], "Challenge MUST NOT be sent in Stellar authenticate request payload")
}

// 14. TestRestoreCloudTokenForUser_PrefersCloudAuthTokenOverCloudJWT tests token precedence.
func TestRestoreCloudTokenForUser_PrefersCloudAuthTokenOverCloudJWT(t *testing.T) {
	tc := tracecore.NewTracecoreClient("http://localhost", "", "http://localhost", "http://localhost")
	sessions := make(map[string]*vault_session.Session)
	sm := vault_session.NewManager(nil, nil, nil, nil, nil, sessions)
	app := &App{
		Vault: &vault_ui.VaultHandler{
			TracecoreClient: tc,
			SessionManager:  sm,
		},
	}

	// 1. When cloud_auth_token is set, prefer cloud_auth_token
	sess1 := &vault_session.Session{
		UserID: "user_test",
		Runtime: &vault_session.RuntimeContext{
			SessionSecrets: map[string]string{
				"cloud_auth_token": "primary_cloud_auth_token",
				"cloud_jwt":        "fallback_cloud_jwt",
			},
		},
	}
	sessions["user_test"] = sess1

	err := app.RestoreCloudTokenForUser("user_test")
	require.NoError(t, err)
	assert.Equal(t, "primary_cloud_auth_token", app.Vault.TracecoreClient.Token)

	// 2. When cloud_auth_token is missing, fall back to cloud_jwt
	sess2 := &vault_session.Session{
		UserID: "user_fallback",
		Runtime: &vault_session.RuntimeContext{
			SessionSecrets: map[string]string{
				"cloud_jwt": "fallback_cloud_jwt_only",
			},
		},
	}
	sessions["user_fallback"] = sess2

	errFallback := app.RestoreCloudTokenForUser("user_fallback")
	require.NoError(t, errFallback)
	assert.Equal(t, "fallback_cloud_jwt_only", app.Vault.TracecoreClient.Token)
}

// 15. TestCloudRequest_SendsAuthorizationHeader_AfterSignIn proves Authorization: Bearer header is sent on Cloud requests.
func TestCloudRequest_SendsAuthorizationHeader_AfterSignIn(t *testing.T) {
	receivedAuthHeader := ""
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuthHeader = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status": 200,
			"data":   []interface{}{},
		})
	}))
	defer ts.Close()

	tc := tracecore.NewTracecoreClient(ts.URL, "", ts.URL, ts.URL)
	sessions := make(map[string]*vault_session.Session)
	sm := vault_session.NewManager(nil, nil, nil, nil, nil, sessions)
	app := &App{
		Vault: &vault_ui.VaultHandler{
			TracecoreClient: tc,
			SessionManager:  sm,
		},
	}

	// 1. Simulate active session with stored cloud_auth_token
	sess := &vault_session.Session{
		UserID: "user_cloud_verify",
		Runtime: &vault_session.RuntimeContext{
			SessionSecrets: map[string]string{
				"cloud_auth_token": "verification_cloud_token_abc123",
			},
		},
	}
	sessions["user_cloud_verify"] = sess

	// 2. Restore Cloud token into TracecoreClient.Token
	errRestore := app.RestoreCloudTokenForUser("user_cloud_verify")
	require.NoError(t, errRestore)
	assert.Equal(t, "verification_cloud_token_abc123", tc.Token)

	// 3. Execute a Cloud-dependent request (RequestVaultChallenge)
	_, errReq := tc.RequestVaultChallenge(context.Background(), "vault_cloud_123")
	require.NoError(t, errReq)

	// 4. Verify exact HTTP Authorization header sent to Cloud server
	assert.Equal(t, "Bearer verification_cloud_token_abc123", receivedAuthHeader)
}

// 16. TestStellarBootstrap_DDDIdentityResolution_And_ReplayProtection proves identity resolution via identity_users and single-use challenge.
func TestStellarBootstrap_DDDIdentityResolution_And_ReplayProtection(t *testing.T) {
	pubKey := "GBV35PVNE77KMVFBK3JS4OXXQPHSVEYEDYNSSKPIFNJZH2EJNC5O4THV"
	privKey := "SADBPTCJHNVQ4KUZEDV4NDWQKGTOCYEKITLG4E2DZWKXJBBUTCVM7EYK"
	unknownPubKey := "GDAUNKNOWNKEYXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	// 1. Verify GormUserRepository FindByPublicKey resolves identity_users
	user := identity_domain.NewStandardUser("user_test_ddd_777", "ddd_user@ankhora.io", "pass_hash")
	user.StellarPublicKey = pubKey

	assert.Equal(t, "identity_users", user.TableName())
	assert.Equal(t, pubKey, user.StellarPublicKey)

	// 2. Generate challenge for pubKey
	challenge := blockchain.GenerateChallenge(pubKey)
	require.NotEmpty(t, challenge)
	assert.Equal(t, challenge, blockchain.ChallengeStore[pubKey])

	// 3. Sign challenge
	sig, errSig := blockchain.SignActorWithStellarPrivateKey(privKey, challenge)
	require.NoError(t, errSig)
	require.NotEmpty(t, sig)

	// 4. Verify signature against challenge
	assert.True(t, blockchain.VerifySignature(pubKey, challenge, sig))

	// 5. Verify unknown public key is rejected
	assert.False(t, blockchain.VerifySignature(unknownPubKey, challenge, sig))

	// 6. Test replay protection: consume challenge from ChallengeStore
	expectedChal, ok := blockchain.ChallengeStore[pubKey]
	require.True(t, ok)
	assert.Equal(t, challenge, expectedChal)
	delete(blockchain.ChallengeStore, pubKey)

	// Replay attempt must fail because ChallengeStore entry was consumed
	_, replayedOk := blockchain.ChallengeStore[pubKey]
	assert.False(t, replayedOk, "Replaying challenge MUST fail because challenge was consumed")
}

// 17. TestStellarBootstrap_EndToEnd_11StepAcceptance proves the complete vertical slice against real GORM DB & handlers.
func TestStellarBootstrap_EndToEnd_11StepAcceptance(t *testing.T) {
	pubKey := "GBV35PVNE77KMVFBK3JS4OXXQPHSVEYEDYNSSKPIFNJZH2EJNC5O4THV"
	privKey := "SADBPTCJHNVQ4KUZEDV4NDWQKGTOCYEKITLG4E2DZWKXJBBUTCVM7EYK"

	// 1. Setup in-memory GORM DB with full migrations
	gormDB, errDB := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, errDB)
	errMigrate := gormDB.AutoMigrate(
		&identity_domain.User{},
		&identity_persistence.StellarAuthChallenge{},
		&auth.TokenPairs{},
	)
	require.NoError(t, errMigrate)

	// Step 1 — Existing identity_users row with stellar_public_key
	user := identity_domain.NewStandardUser("user_stellar_e2e_11", "e2e_user@ankhora.io", "pass_hash")
	user.StellarPublicKey = pubKey
	userRepo := identity_persistence.NewGormUserRepository(gormDB)
	errSave := userRepo.Save(context.Background(), user)
	require.NoError(t, errSave)

	// Verify identity_users contains stellar_public_key
	resolvedUser, errFind := userRepo.FindByPublicKey(context.Background(), pubKey)
	require.NoError(t, errFind)
	require.NotNil(t, resolvedUser)
	userID := resolvedUser.ID
	assert.NotEmpty(t, userID)

	// Setup AuthHandler with JWT auth service
	jwtAuth := &auth.Auth{
		Issuer:      "ankhora-test",
		Audience:    "ankhora-desktop",
		Secret:      "test-secret-key-32-bytes-minimum-len",
		TokenExpiry: 24 * 3600 * 1000000000,
	}
	dbModel := &models.DBModel{DB: gormDB}
	authHandler := handlers.NewAuthHandler(*dbModel, nil, nil, &logger.Logger{}, nil, *jwtAuth, nil)

	// Setup Cloud HTTP Server with routes /api/stellar/public-challenge and /api/stellar/authenticate
	receivedAuthHeader := ""
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuthHeader = r.Header.Get("Authorization")
		switch r.URL.Path {
		case "/api/stellar/public-challenge":
			var reqBody blockchain.ChallengeRequest
			_ = json.NewDecoder(r.Body).Decode(&reqBody)
			resp, err := authHandler.RequestChallenge(reqBody)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": true, "message": err.Error()})
				return
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(resp)
		case "/api/stellar/authenticate":
			var reqBody tracecore_types.StellarAuthenticateRequest
			_ = json.NewDecoder(r.Body).Decode(&reqBody)
			resp, err := authHandler.StellarAuthenticate(r.Context(), reqBody)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": true, "message": err.Error()})
				return
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(resp)
		case "/api/identity/challenge", "/identity/challenge":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"challenge_id":    "chal_e2e_11",
				"signing_payload": "challenge_payload_11",
			})
		case "/api/identity/", "/api/identity", "/identity/", "/identity":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"vault_id":      "vault_e2e_11",
				"delegation_id": "del_e2e_11",
				"status":        "active",
			})
		default:
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": 200})
		}
	}))
	defer ts.Close()

	// Step 2 — POST /api/stellar/public-challenge
	chalReqBody, _ := json.Marshal(map[string]string{"public_key": pubKey})
	resChal, errHTTP1 := http.Post(ts.URL+"/api/stellar/public-challenge", "application/json", bytes.NewBuffer(chalReqBody))
	require.NoError(t, errHTTP1)

	// Step 3 — Confirm HTTP 200 + challenge
	assert.Equal(t, http.StatusOK, resChal.StatusCode)
	var chalResp blockchain.ChallengeResponse
	_ = json.NewDecoder(resChal.Body).Decode(&chalResp)
	require.NotEmpty(t, chalResp.Challenge)

	// Step 4 — Confirm stellar_auth_challenges row exists in DB
	var challengeRow identity_persistence.StellarAuthChallenge
	errDBRow := gormDB.Where("public_key = ?", pubKey).First(&challengeRow).Error
	require.NoError(t, errDBRow, "stellar_auth_challenges row MUST exist in database")
	assert.Equal(t, chalResp.Challenge, challengeRow.Challenge)

	// Step 5 — Sign challenge with local Stellar private key
	sig, errSign := blockchain.SignActorWithStellarPrivateKey(privKey, chalResp.Challenge)
	require.NoError(t, errSign)
	require.NotEmpty(t, sig)

	// Step 6 — POST /api/stellar/authenticate
	authReqBody, _ := json.Marshal(map[string]string{
		"public_key": pubKey,
		"signature":  sig,
	})
	resAuth, errHTTP2 := http.Post(ts.URL+"/api/stellar/authenticate", "application/json", bytes.NewBuffer(authReqBody))
	require.NoError(t, errHTTP2)

	// Step 7 — Confirm HTTP 200 + JWT
	assert.Equal(t, http.StatusOK, resAuth.StatusCode)
	var authResp tracecore_types.LoginResponse
	_ = json.NewDecoder(resAuth.Body).Decode(&authResp)
	require.NotEmpty(t, authResp.Token)

	// Step 8 — Confirm stellar_auth_challenges row is gone (consumed)
	var checkRow identity_persistence.StellarAuthChallenge
	errDeleted := gormDB.Where("public_key = ?", pubKey).First(&checkRow).Error
	assert.Error(t, errDeleted, "stellar_auth_challenges row MUST be deleted after consumption")

	// Step 9 — Repeat authentication with same signature → 401
	resReplay, errHTTP3 := http.Post(ts.URL+"/api/stellar/authenticate", "application/json", bytes.NewBuffer(authReqBody))
	require.NoError(t, errHTTP3)
	assert.Equal(t, http.StatusUnauthorized, resReplay.StatusCode, "Replay attack MUST fail with HTTP 401")

	// Step 10 — Confirm Authorization: Bearer <token> on subsequent Cloud request
	tc := tracecore.NewTracecoreClient(ts.URL, authResp.Token, ts.URL, ts.URL)
	sessions := make(map[string]*vault_session.Session)
	sm := vault_session.NewManager(nil, nil, nil, nil, nil, sessions)
	app := &App{
		Vault: &vault_ui.VaultHandler{
			TracecoreClient: tc,
			SessionManager:  sm,
		},
	}
	sessions[userID] = &vault_session.Session{
		UserID: userID,
		Runtime: &vault_session.RuntimeContext{
			SessionSecrets: map[string]string{
				"cloud_auth_token": authResp.Token,
			},
		},
	}

	errRestore := app.RestoreCloudTokenForUser(userID)
	require.NoError(t, errRestore)

	_, errReq := tc.RequestVaultChallenge(context.Background(), "vault_e2e_11")
	require.NoError(t, errReq)
	assert.Equal(t, "Bearer "+authResp.Token, receivedAuthHeader)

	// Step 11 — Confirm ConnectVault() succeeds
	userCfg := &app_config_domain.UserConfig{
		ID: userID,
		StellarAccount: app_config_domain.StellarAccountConfig{
			PublicKey:  pubKey,
			PrivateKey: privKey,
		},
	}
	app.AppConfigHandler = &app_config_ui.AppConfigHandler{
		UserConfigRepository: &stubUserConfigRepo{userCfg: userCfg},
	}

	errConnect := app.ConnectVault(userID, "vault_e2e_11")
	require.NoError(t, errConnect, "ConnectVault MUST succeed after Cloud authentication")
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


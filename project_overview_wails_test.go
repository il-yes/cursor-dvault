package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	auth_usecases "vault-app/internal/auth/application/use_cases"
	auth_domain "vault-app/internal/auth/domain"
	auth_persistence "vault-app/internal/auth/infrastructure/persistence"
	auth_ui "vault-app/internal/auth/ui"
	identity_ui "vault-app/internal/identity/ui"
	"vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
	vault_session "vault-app/internal/vault/application/session"
	vault_ui "vault-app/internal/vault/ui"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestApp_GetProjectOverview_Unauthenticated_FailsGuard(t *testing.T) {
	app := &App{}

	// Invalid or unauthenticated call should fail RequireAuth
	_, err := app.GetProjectOverview("invalid-token", "PRJ-001")
	if err == nil {
		t.Fatal("Expected error for unauthenticated call, got nil")
	}
}

func TestApp_GetProjectOverview_CloudUnauthenticated_FailsCloudGuard(t *testing.T) {
	// Setup local in-memory DB and AuthHandler for valid local JWT
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test db: %v", err)
	}

	authV2 := auth_domain.Auth{
		Issuer:        "test-issuer",
		Audience:      "test-aud",
		Secret:        "test-secret-32-bytes-minimum-secret-key-123",
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
	}

	authRepo := auth_persistence.NewGormAuthRepository(db)
	tokenSvc := auth_usecases.NewTokenService(authV2, authRepo, db)
	tokenUC := auth_usecases.NewGenerateTokensUseCase(authRepo, tokenSvc)
	authH := auth_ui.NewAuthHandler(&identity_ui.IdentityHandler{}, tokenUC, db)

	tokens, err := tokenSvc.GenerateTokenPair(&auth_domain.JwtUser{ID: "user-123", Username: "User", Email: "test@example.com"})
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	traceClientWithoutToken := tracecore.NewTracecoreClient("http://localhost:4001/api", "", "", "")
	vaultH := &vault_ui.VaultHandler{
		TracecoreClient: traceClientWithoutToken,
	}

	app := &App{
		AuthHandler:     authH,
		Vault:           vaultH,
		tracecoreClient: traceClientWithoutToken,
	}

	// Local JWT passes RequireAuth, but RequireCloudAuthentication fails because Cloud token is empty
	_, err = app.GetProjectOverview(tokens.Token, "PRJ-001")
	if err == nil {
		t.Fatal("Expected error for missing Cloud token, got nil")
	}
	if !strings.Contains(err.Error(), "cloud") && !strings.Contains(err.Error(), "token") {
		t.Errorf("Expected cloud authentication error, got: %v", err)
	}
}

func TestApp_GetProjectOverview_FullyAuthenticated_Success(t *testing.T) {
	expectedProjectID := "PRJ-001"
	cloudToken := "valid-cloud-bearer-token"

	// Mock Cloud HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/construction/projects/"+expectedProjectID {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if r.Header.Get("Authorization") != "Bearer "+cloudToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		resp := tracecore_types.CloudResponse[tracecore_types.ProjectOverviewDTO]{
			Status: 200,
			Data: tracecore_types.ProjectOverviewDTO{
				ID:                 "PRJ-001",
				Code:               "PRJ-001",
				Name:               "Boulevard Haussmann Retrofit",
				Type:               "Infrastructure",
				Status:             "Active",
				ProgressPercent:    68,
				ScheduleDay:        142,
				BudgetSpentPercent: 68,
			},
			Message: "success",
			Success: true,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Setup local in-memory DB and AuthHandler for valid local JWT
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test db: %v", err)
	}

	authV2 := auth_domain.Auth{
		Issuer:        "test-issuer",
		Audience:      "test-aud",
		Secret:        "test-secret-32-bytes-minimum-secret-key-123",
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
	}

	authRepo := auth_persistence.NewGormAuthRepository(db)
	tokenSvc := auth_usecases.NewTokenService(authV2, authRepo, db)
	tokenUC := auth_usecases.NewGenerateTokensUseCase(authRepo, tokenSvc)
	authH := auth_ui.NewAuthHandler(&identity_ui.IdentityHandler{}, tokenUC, db)

	tokens, err := tokenSvc.GenerateTokenPair(&auth_domain.JwtUser{ID: "user-123", Username: "User", Email: "test@example.com"})
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	traceClient := tracecore.NewTracecoreClient(server.URL+"/api", cloudToken, server.URL, server.URL)
	vaultH := &vault_ui.VaultHandler{
		TracecoreClient: traceClient,
	}

	runtimeCtx := &vault_session.RuntimeContext{
		SessionSecrets: map[string]string{
			"cloud_auth_token": cloudToken,
		},
	}

	app := &App{
		ctx:             context.Background(),
		AuthHandler:     authH,
		Vault:           vaultH,
		tracecoreClient: traceClient,
		RuntimeContext:  runtimeCtx,
	}

	overview, err := app.GetProjectOverview(tokens.Token, expectedProjectID)
	if err != nil {
		t.Fatalf("App.GetProjectOverview returned unexpected error: %v", err)
	}

	if overview.ID != "PRJ-001" || overview.Name != "Boulevard Haussmann Retrofit" {
		t.Errorf("Unexpected ProjectOverview returned: %+v", overview)
	}
	if overview.ProgressPercent != 68 || overview.ScheduleDay != 142 {
		t.Errorf("Expected ProgressPercent=68, ScheduleDay=142, got: %+v", overview)
	}
}

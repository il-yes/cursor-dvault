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

func TestApp_ListConstructionProjects_FullyAuthenticated_Success(t *testing.T) {
	cloudToken := "valid-cloud-bearer-token"

	// Mock Cloud HTTP server returning exact real Cloud endpoint response JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/construction/projects" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if r.Header.Get("Authorization") != "Bearer "+cloudToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		rawCloudJSON := `{
			"status": 200,
			"data": [
				{
					"project_id": "PRJ-001",
					"workspace_id": "52e511f4-789d-4cf6-b5e3-bb5448d4714d",
					"project_reference": "PRJ-METRO-001",
					"project_name": "Metro Line 4 Expansion",
					"project_type": "Infrastructure",
					"status": "active",
					"location": {
						"address": "",
						"city": "",
						"country": "",
						"coordinates": {"lat": 0, "lng": 0}
					},
					"start_date": "2026-01-01",
					"planned_end_date": "2027-12-31",
					"actual_end_date": "",
					"budget": {"currency": "", "total": 0, "allocated": 0},
					"progress_percentage": 25,
					"current_phase": "Phase 2 Structural",
					"milestones": null,
					"stakeholder_ids": null,
					"requirement_ids": null,
					"site_id": "",
					"status_summary": "Civil works underway",
					"created_at": "2026-09-25T18:46:08.731Z",
					"updated_at": "2026-09-25T18:46:08.731Z"
				}
			],
			"message": "Enregistrement recupéré avec succès",
			"success": true
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(rawCloudJSON))
	}))
	defer server.Close()

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

	projects, err := app.ListConstructionProjects(tokens.Token, "")
	if err != nil {
		t.Fatalf("App.ListConstructionProjects returned unexpected error: %v", err)
	}

	if len(projects) != 1 {
		t.Fatalf("Expected 1 project from Cloud response, got %d", len(projects))
	}
	p := projects[0]
	if p.ID != "PRJ-001" || p.ProjectID != "PRJ-001" {
		t.Errorf("Expected ID/ProjectID 'PRJ-001', got ID='%s', ProjectID='%s'", p.ID, p.ProjectID)
	}
	if p.Code != "PRJ-METRO-001" || p.ProjectReference != "PRJ-METRO-001" {
		t.Errorf("Expected Code/ProjectReference 'PRJ-METRO-001', got Code='%s', ProjectReference='%s'", p.Code, p.ProjectReference)
	}
	if p.Name != "Metro Line 4 Expansion" || p.ProjectName != "Metro Line 4 Expansion" {
		t.Errorf("Expected Name/ProjectName 'Metro Line 4 Expansion', got Name='%s', ProjectName='%s'", p.Name, p.ProjectName)
	}
	if p.ProgressPercent != 25 {
		t.Errorf("Expected ProgressPercent=25, got %d", p.ProgressPercent)
	}
	if p.CurrentPhase != "Phase 2 Structural" {
		t.Errorf("Expected CurrentPhase 'Phase 2 Structural', got '%s'", p.CurrentPhase)
	}
	if p.RecentActivity != "Civil works underway" {
		t.Errorf("Expected RecentActivity 'Civil works underway', got '%s'", p.RecentActivity)
	}
}

func TestApp_ListConstructionProjects_WithScalarVaultID_Success(t *testing.T) {
	cloudToken := "valid-cloud-bearer-token"
	targetVaultID := "c59d1d77-a50f-49f6-886b-76bf740925c0"
	var queryVaultID string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/construction/projects" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		queryVaultID = r.URL.Query().Get("vault_id")

		rawCloudJSON := `{
			"status": 200,
			"data": [],
			"message": "Enregistrement recupéré avec succès",
			"success": true
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(rawCloudJSON))
	}))
	defer server.Close()

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

	tokens, err := tokenSvc.GenerateTokenPair(&auth_domain.JwtUser{ID: "user-70", Username: "User70", Email: "user70@example.com"})
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

	projects, err := app.ListConstructionProjects(tokens.Token, targetVaultID)
	if err != nil {
		t.Fatalf("App.ListConstructionProjects returned unexpected error: %v", err)
	}

	if queryVaultID != targetVaultID {
		t.Errorf("Expected query vault_id '%s', got '%s'", targetVaultID, queryVaultID)
	}
	if len(projects) != 0 {
		t.Errorf("Expected 0 projects for user 70 vault, got %d", len(projects))
	}
}


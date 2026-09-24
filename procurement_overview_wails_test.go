package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestApp_GetProcurementOverview_Unauthenticated_FailsGuard(t *testing.T) {
	app := &App{}

	_, err := app.GetProcurementOverview("invalid-token", "REQ-STRUCT-001")
	if err == nil {
		t.Fatal("Expected error for unauthenticated call, got nil")
	}
}

func TestApp_GetProcurementOverview_FullyAuthenticated_Success(t *testing.T) {
	expectedReqID := "REQ-STRUCT-001"
	cloudToken := "valid-cloud-bearer-token-456"

	// Mock Cloud HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/construction/requirements/"+expectedReqID {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if r.Header.Get("Authorization") != "Bearer "+cloudToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		resp := tracecore_types.CloudResponse[tracecore_types.ProcurementOverviewDTO]{
			Status: 200,
			Data: tracecore_types.ProcurementOverviewDTO{
				ID:                    "REQ-STRUCT-001",
				RequirementID:         "REQ-STRUCT-001",
				ProjectID:             "PRJ-001",
				ProjectCode:           "PRJ-001",
				ProjectName:           "Metro Line 4 Expansion",
				Code:                  "REQ-STRUCT-001",
				MaterialID:            "MAT-STRUCT-001",
				MaterialName:          "High-Strength Structural Steel Beams (Grade A992)",
				Specification:         "W18x86 Heavy Flange Beams...",
				Quantity:              120,
				Unit:                  "Metric Tons",
				TargetPhase:           "Phase 4 - Structure Pier Support",
				Status:                "fulfilled_delayed",
				RequiredDate:          "2026-08-15",
				Priority:              "CRITICAL",
				SiteID:                "SITE-SOUTH-01",
				SiteName:              "Site Alpha - South Pier Foundation",
				InvitedSuppliersCount: 3,
				OffersReceivedCount:   2,
				OfferID:               "OFF-1042",
				OfferReference:        "OFF-1042",
				SupplierID:            "SUP-001",
				SupplierName:          "Apex Steel Fabrication Ltd.",
				TotalPrice:            "$142,500.00",
				UnitPrice:             "$1,187.50 / Ton",
				PromisedDeliveryDate:  "2026-08-15",
				OfferStatus:           "ACCEPTED",
				IsVerifiedSupplier:    true,
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

	overview, err := app.GetProcurementOverview(tokens.Token, expectedReqID)
	if err != nil {
		t.Fatalf("App.GetProcurementOverview returned unexpected error: %v", err)
	}

	if overview.ID != "REQ-STRUCT-001" || overview.OfferReference != "OFF-1042" {
		t.Errorf("Unexpected ProcurementOverview returned: %+v", overview)
	}
	if overview.InvitedSuppliersCount != 3 || overview.OffersReceivedCount != 2 {
		t.Errorf("Expected InvitedSuppliersCount=3, OffersReceivedCount=2, got: %+v", overview)
	}
}

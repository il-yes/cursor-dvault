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

func TestApp_GetLogisticsOverview_Unauthenticated_FailsGuard(t *testing.T) {
	app := &App{}

	_, err := app.GetLogisticsOverview("invalid-token", "DEL-1042")
	if err == nil {
		t.Fatal("Expected error for unauthenticated call, got nil")
	}
}

func TestApp_GetLogisticsOverview_FullyAuthenticated_Success(t *testing.T) {
	expectedDelID := "DEL-1042"
	cloudToken := "valid-cloud-bearer-token-789"

	// Mock Cloud HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/construction/deliveries/"+expectedDelID+"/overview" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if r.Header.Get("Authorization") != "Bearer "+cloudToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		resp := tracecore_types.CloudResponse[tracecore_types.LogisticsOverviewDTO]{
			Status: 200,
			Data: tracecore_types.LogisticsOverviewDTO{
				Delivery: tracecore_types.DeliveryOverviewDTO{
					ID:                  "DEL-1042",
					Reference:           "DEL-1042",
					Status:              "delayed",
					PlannedDeliveryDate: "2026-08-15",
					ETA:                 "2026-08-16 07:30",
					ActualDeliveryDate:  "2026-08-16 07:51",
					Quantity:            120,
					Unit:                "Metric Tons",
					DeliveryNotes:       "Critical load for Pier 4 framework.",
				},
				Project: tracecore_types.ProjectReferenceDTO{
					ID:   "PRJ-001",
					Code: "PRJ-001",
					Name: "Metro Line 4 Expansion",
					Type: "INFRASTRUCTURE",
				},
				Requirement: tracecore_types.RequirementReferenceDTO{
					ID:          "REQ-STRUCT-001",
					Code:        "REQ-STRUCT-001",
					TargetPhase: "Phase 4 - Structure Pier Support",
				},
				Material: tracecore_types.MaterialReferenceDTO{
					ID:            "MAT-STRUCT-001",
					Name:          "High-Strength Structural Steel Beams (Grade A992)",
					Specification: "W18x86 Heavy Flange Beams",
				},
				Supplier: tracecore_types.SupplierReferenceDTO{
					ID:   "SUP-001",
					Name: "Apex Steel Fabrication Ltd.",
				},
				Offer: tracecore_types.OfferReferenceDTO{
					ID:                   "OFF-1042",
					Reference:            "OFF-1042",
					TotalPrice:           "$142,500.00",
					UnitPrice:            "$1,187.50 / Ton",
					PromisedDeliveryDate: "2026-08-15",
				},
				Site: tracecore_types.SiteReferenceDTO{
					ID:   "SITE-SOUTH-01",
					Name: "Site Alpha - South Pier Foundation",
				},
				Transport: tracecore_types.TransportOverviewDTO{
					ID:               "TR-1042",
					Reference:        "TR-1042",
					Vehicle:          "Heavy Hauler Truck #88",
					Driver:           "Mark Vance",
					Status:           "rerouted",
					Origin:           "Gary, IN",
					Destination:      "Chicago, IL",
					Route:            "M1 Highway Northbound -> Rerouted via Highway B",
					PlannedDeparture: "2026-08-15 06:00",
					ActualDeparture:  "2026-08-15 06:15",
					PlannedArrival:   "2026-08-15 10:00",
					ETA:              "2026-08-16 07:30",
					ActualArrival:    "2026-08-16 07:51",
					Constraints:      []string{"Axle weight limit max 35T on M1 Km 42 bridge structure."},
					DelayReason:      "Road Restriction on M1 Highway",
				},
				Issue: &tracecore_types.IssueOverviewDTO{
					ID:                 "ISS-1042",
					Reference:          "ISS-1042",
					Title:              "Road Restriction on M1 Highway - Axle Weight Limit",
					Description:        "Emergency weight reduction to 35T.",
					Severity:           "CRITICAL",
					Status:             "RESOLVED",
					ReportedBy:         "Mark Vance",
					ReportedAt:         "2026-08-15 11:41",
					Impact:             "Foundation phase delayed by 1 day",
					EvidenceReferences: []string{"DOC-EVID-001"},
				},
				Decision: &tracecore_types.DecisionOverviewDTO{
					ID:                    "DEC-1042",
					Reference:             "DEC-1042",
					Subject:               "Approve Alternative Route B",
					Context:               "Supplier reports road restriction on primary route M1",
					TechnicalAssessment:   "Alternative Route B remains compatible with vehicle constraints",
					RisksIdentified:       []string{"Additional transport cost", "6h delay"},
					ParticipantsConsulted: []string{"Alex Rivera", "David Chen", "Mark Vance"},
					OptionsConsidered:     []string{"Wait for original route", "Reroute via Highway B"},
					Decision:              "Reroute transport via Highway B with state police escort",
					DecidedBy:             "Alex Rivera",
					DecisionDate:          "2026-08-15 14:20",
					Consequence:           "Expected delivery restored to Aug 16 07:30",
					Status:                "EXECUTED",
				},
				Inspection: &tracecore_types.InspectionOverviewDTO{
					ID:             "INSP-1042",
					Reference:      "INSP-1042",
					Type:           "Structural Steel Weld & Dimensional Acceptance Inspection",
					Inspector:      "Sarah Jenkins, PE",
					InspectionDate: "2026-08-16 09:20",
					Status:         "PASSED",
					Result:         "ACCEPTED",
					Criteria:       []string{"Material conforms to spec", "Quantity matches delivery"},
					Findings:       []string{"100% Pass - Zero structural defects detected"},
					Notes:          "Approved for immediate hoisting on Pier 4 framework.",
				},
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

	overview, err := app.GetLogisticsOverview(tokens.Token, expectedDelID)
	if err != nil {
		t.Fatalf("App.GetLogisticsOverview returned unexpected error: %v", err)
	}

	if overview.Delivery.ID != "DEL-1042" || overview.Transport.ID != "TR-1042" {
		t.Errorf("Unexpected LogisticsOverview returned: %+v", overview)
	}
	if overview.Issue == nil || overview.Issue.ID != "ISS-1042" {
		t.Errorf("Expected Issue ID ISS-1042, got: %+v", overview.Issue)
	}
	if overview.Decision == nil || overview.Decision.ID != "DEC-1042" {
		t.Errorf("Expected Decision ID DEC-1042, got: %+v", overview.Decision)
	}
	if overview.Inspection == nil || overview.Inspection.Result != "ACCEPTED" {
		t.Errorf("Expected Inspection Result ACCEPTED, got: %+v", overview.Inspection)
	}
}

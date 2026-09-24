package tracecore_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
)

func TestGetProcurementOverview_Success_PathAndAuthPreserved(t *testing.T) {
	expectedToken := "test-bearer-token-456"
	expectedReqID := "REQ-STRUCT-001"
	pathHit := false
	authHeaderHit := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/api/construction/requirements/"+expectedReqID {
			pathHit = true
		}
		if r.Header.Get("Authorization") == "Bearer "+expectedToken {
			authHeaderHit = true
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
				Specification:         "W18x86 Heavy Flange Beams, ASTM A992 certified, anti-corrosion primer coated.",
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

	client := tracecore.NewTracecoreClient(server.URL+"/api", expectedToken, server.URL, server.URL)
	overview, err := client.GetProcurementOverview(context.Background(), expectedReqID)
	if err != nil {
		t.Fatalf("GetProcurementOverview returned unexpected error: %v", err)
	}

	if !pathHit {
		t.Errorf("Expected request path /api/construction/requirements/%s to be hit", expectedReqID)
	}
	if !authHeaderHit {
		t.Errorf("Expected Authorization header 'Bearer %s'", expectedToken)
	}
	if overview.ID != "REQ-STRUCT-001" || overview.OfferReference != "OFF-1042" {
		t.Errorf("Decoded overview mismatch, got: %+v", overview)
	}
	if overview.InvitedSuppliersCount != 3 || overview.OffersReceivedCount != 2 {
		t.Errorf("Expected InvitedSuppliersCount=3, OffersReceivedCount=2, got: %+v", overview)
	}
}

func TestGetProcurementOverview_EmptyReqID_ReturnsError(t *testing.T) {
	client := tracecore.NewTracecoreClient("http://localhost:4001/api", "token", "", "")
	_, err := client.GetProcurementOverview(context.Background(), "")
	if err == nil {
		t.Fatal("Expected error for empty requirementID, got nil")
	}
}

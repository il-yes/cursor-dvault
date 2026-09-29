package tracecore_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
)

func TestGetProjectOverview_Success_PathAndAuthPreserved(t *testing.T) {
	expectedToken := "test-bearer-token-123"
	expectedProjectID := "PRJ-001"
	pathHit := false
	authHeaderHit := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/api/construction/projects/"+expectedProjectID {
			pathHit = true
		}
		if r.Header.Get("Authorization") == "Bearer "+expectedToken {
			authHeaderHit = true
		}

		resp := tracecore_types.CloudResponse[tracecore_types.ProjectOverviewDTO]{
			Status: 200,
			Data: tracecore_types.ProjectOverviewDTO{
				ID:                    "PRJ-001",
				Code:                  "PRJ-001",
				Name:                  "Boulevard Haussmann Retrofit",
				ContractID:            "BFD-EUR-2024-099",
				Type:                  "Infrastructure",
				Sector:                "Infrastructure Sector • Transit Hub",
				Status:                "Active • On Schedule",
				Location:              tracecore_types.FlexLocation{Value: "Paris, France"},
				Description:           "Urban infrastructure renewal",
				CurrentPhase:          "Structure (Phase 4 of 7)",
				ProgressPercent:       68,
				OpenIssuesCount:       7,
				PendingDecisionsCount: 1,
				TargetCompletion:      "Oct 2026",
				BudgetSpentPercent:    68,
				ScheduleDay:           142,
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
	overview, err := client.GetProjectOverview(context.Background(), expectedProjectID)
	if err != nil {
		t.Fatalf("GetProjectOverview returned unexpected error: %v", err)
	}

	if !pathHit {
		t.Errorf("Expected request path /api/construction/projects/%s to be hit", expectedProjectID)
	}
	if !authHeaderHit {
		t.Errorf("Expected Authorization header 'Bearer %s'", expectedToken)
	}
	if overview.ID != "PRJ-001" || overview.Name != "Boulevard Haussmann Retrofit" {
		t.Errorf("Decoded overview mismatch, got: %+v", overview)
	}
	if overview.ProgressPercent != 68 || overview.ScheduleDay != 142 {
		t.Errorf("Expected ProgressPercent=68, ScheduleDay=142, got: %+v", overview)
	}
}

func TestGetProjectOverview_BareJSON_DecodesCorrectly(t *testing.T) {
	expectedProjectID := "PRJ-002"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dto := tracecore_types.ProjectOverviewDTO{
			ProjectID:          "PRJ-002",
			ProjectReference:   "PRJ-002",
			ProjectName:        "Commercial Plaza North",
			ProjectType:        "Commercial Mixed-Use",
			Status:             "Active",
			Location:           tracecore_types.FlexLocation{Value: "Lyon, France"},
			ProgressPercent:    42,
			BudgetSpentPercent: 42,
			ScheduleDay:        88,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(dto)
	}))
	defer server.Close()

	client := tracecore.NewTracecoreClient(server.URL+"/api", "token", server.URL, server.URL)
	overview, err := client.GetProjectOverview(context.Background(), expectedProjectID)
	if err != nil {
		t.Fatalf("GetProjectOverview returned unexpected error: %v", err)
	}

	if overview.ID != "PRJ-002" || overview.Code != "PRJ-002" || overview.Name != "Commercial Plaza North" {
		t.Errorf("Expected normalized fields (ID/Code/Name), got: %+v", overview)
	}
}

func TestGetProjectOverview_Cloud404_PropagatesError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"project not found"}`))
	}))
	defer server.Close()

	client := tracecore.NewTracecoreClient(server.URL+"/api", "token", server.URL, server.URL)
	_, err := client.GetProjectOverview(context.Background(), "PRJ-UNKNOWN")
	if err == nil {
		t.Fatal("Expected error for 404 response, got nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("Expected 404 in error message, got: %v", err)
	}
}

func TestGetProjectOverview_EmptyProjectID_ReturnsError(t *testing.T) {
	client := tracecore.NewTracecoreClient("http://localhost:4001/api", "token", "", "")
	_, err := client.GetProjectOverview(context.Background(), "")
	if err == nil {
		t.Fatal("Expected error for empty projectID, got nil")
	}
}

func TestListConstructionProjects_Success_CloudResponseEnvelope(t *testing.T) {
	expectedToken := "test-bearer-token-123"
	pathHit := false
	authHeaderHit := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/api/construction/projects" {
			pathHit = true
		}
		if r.Header.Get("Authorization") == "Bearer "+expectedToken {
			authHeaderHit = true
		}

		resp := tracecore_types.CloudResponse[[]tracecore_types.ProjectOverviewDTO]{
			Status: 200,
			Data: []tracecore_types.ProjectOverviewDTO{
				{
					ID:              "PRJ-001",
					Code:            "PRJ-001",
					Name:            "Boulevard Haussmann Retrofit",
					ProgressPercent: 68,
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

	client := tracecore.NewTracecoreClient(server.URL+"/api", expectedToken, server.URL, server.URL)
	projects, err := client.ListConstructionProjects(context.Background(), "")
	if err != nil {
		t.Fatalf("ListConstructionProjects returned unexpected error: %v", err)
	}

	if !pathHit {
		t.Errorf("Expected request path /api/construction/projects to be hit")
	}
	if !authHeaderHit {
		t.Errorf("Expected Authorization header 'Bearer %s'", expectedToken)
	}
	if len(projects) != 1 {
		t.Fatalf("Expected 1 project, got %d", len(projects))
	}
	if projects[0].ID != "PRJ-001" || projects[0].Name != "Boulevard Haussmann Retrofit" {
		t.Errorf("Decoded project mismatch, got: %+v", projects[0])
	}
}

func TestListConstructionProjects_BareArray_DecodesCorrectly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawList := []tracecore_types.ProjectOverviewDTO{
			{
				ProjectID:   "PRJ-001",
				ProjectName: "Boulevard Haussmann Retrofit",
			},
			{
				ProjectID:   "PRJ-002",
				ProjectName: "Commercial Plaza North",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(rawList)
	}))
	defer server.Close()

	client := tracecore.NewTracecoreClient(server.URL+"/api", "token", server.URL, server.URL)
	projects, err := client.ListConstructionProjects(context.Background(), "")
	if err != nil {
		t.Fatalf("ListConstructionProjects returned unexpected error: %v", err)
	}

	if len(projects) != 2 {
		t.Fatalf("Expected 2 projects, got %d", len(projects))
	}
	if projects[0].ID != "PRJ-001" || projects[1].ID != "PRJ-002" {
		t.Errorf("Expected normalized project IDs, got: %+v", projects)
	}
}

package tracecore_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestListThreadsDirect_CallsCorrectEndpointAndMapsResponse(t *testing.T) {
	channelID := "4ed4728f-458b-4457-9368-fbc4613062ab"
	expectedPath := "/threads/by-channel/" + channelID

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s, want %s", r.URL.Path, expectedPath)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing or incorrect Bearer token: %q", r.Header.Get("Authorization"))
		}

		responseJSON := fmt.Sprintf(`{
			"status": 200,
			"data": [
				{
					"id": "thread-123",
					"channel_id": "%s",
					"title": "TEST-PERSISTENCE",
					"subtitle": "Test Subtitle",
					"asset_type": "note",
					"status": "open",
					"created_at": "2026-08-20T10:00:00Z"
				}
			]
		}`, channelID)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, responseJSON)
	})

	threads, err := client.ListThreadsDirect(context.Background(), "user-1", channelID)
	if err != nil {
		t.Fatalf("ListThreadsDirect failed: %v", err)
	}

	if len(threads) != 1 {
		t.Fatalf("expected 1 thread, got %d", len(threads))
	}

	th := threads[0]
	if th.ID != "thread-123" {
		t.Errorf("thread.ID = %q, want thread-123", th.ID)
	}
	if th.ChannelID != channelID {
		t.Errorf("thread.ChannelID = %q, want %s", th.ChannelID, channelID)
	}
	if th.Title != "TEST-PERSISTENCE" {
		t.Errorf("thread.Title = %q, want TEST-PERSISTENCE", th.Title)
	}
}

func TestCreateThreadDirect_CallsCorrectEndpointAndMapsResponse(t *testing.T) {
	channelID := "4ed4728f-458b-4457-9368-fbc4613062ab"
	expectedPath := "/threads"

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("unexpected path: %s, want %s", r.URL.Path, expectedPath)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing or incorrect Bearer token: %q", r.Header.Get("Authorization"))
		}

		responseJSON := fmt.Sprintf(`{
			"status": 201,
			"data": {
				"id": "cloud-assigned-thread-999",
				"channel_id": "%s",
				"title": "TEST-PERSISTENCE",
				"subtitle": "Test Subtitle",
				"asset_type": "note",
				"status": "open",
				"created_at": "2026-08-20T10:00:00Z"
			}
		}`, channelID)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, responseJSON)
	})

	th, err := client.CreateThreadDirect(context.Background(), "user-1", channelID, "TEST-PERSISTENCE", "Test Subtitle", "note")
	if err != nil {
		t.Fatalf("CreateThreadDirect failed: %v", err)
	}

	if th == nil {
		t.Fatal("expected non-nil thread response")
	}

	if th.ID != "cloud-assigned-thread-999" {
		t.Errorf("th.ID = %q, want cloud-assigned-thread-999", th.ID)
	}
	if th.Title != "TEST-PERSISTENCE" {
		t.Errorf("th.Title = %q, want TEST-PERSISTENCE", th.Title)
	}
}

func TestListThreadsDirect_PrioritizesAnkhoraCloudUrl(t *testing.T) {
	cloudServerCalled := false
	cloudServer := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		cloudServerCalled = true
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":200,"data":[]}`)
	})

	dummyServer := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("request went to BaseURL dummy server instead of AnkhoraCloudUrl server!")
		w.WriteHeader(http.StatusInternalServerError)
	})

	// Set AnkhoraCloudUrl to cloudServer and BaseURL to dummyServer
	client := cloudServer
	client.BaseURL = dummyServer.BaseURL // BaseURL points to dummyServer

	channelID := "ch_url_test"
	_, err := client.ListThreadsDirect(context.Background(), "user-1", channelID)
	if err != nil {
		t.Fatalf("ListThreadsDirect returned unexpected error: %v", err)
	}

	if !cloudServerCalled {
		t.Errorf("expected request to hit AnkhoraCloudUrl server, but it was not called")
	}
}

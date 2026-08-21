package tracecore_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	thread_domain "vault-app/internal/thread/domain"
)

func TestThreadEvents_Repository_VerticalSlice(t *testing.T) {
	threadID := "thread-test-123"
	eventType := "entry.shared"
	idempotencyKey := "evt_share_se_test_key"

	var receivedBody map[string]any

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/threads/" + threadID + "/events":
			if r.Method == http.MethodPost {
				if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				respJSON := fmt.Sprintf(`{
					"status": 201,
					"data": {
						"id": "evt-cloud-456",
						"thread_id": "%s",
						"type": "%s",
						"payload": {
							"ref_type": "share_entry",
							"share_entry_id": "se_123",
							"trust_group_id": "tg_456"
						},
						"idempotency_key": "%s",
						"cursor": 1,
						"created_at": "2026-08-21T10:00:00Z"
					}
				}`, threadID, eventType, idempotencyKey)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				fmt.Fprint(w, respJSON)
				return
			} else if r.Method == http.MethodGet {
				respJSON := fmt.Sprintf(`{
					"status": 200,
					"data": [
						{
							"id": "evt-cloud-456",
							"thread_id": "%s",
							"type": "%s",
							"payload": {
								"ref_type": "share_entry",
								"share_entry_id": "se_123",
								"trust_group_id": "tg_456"
							},
							"idempotency_key": "%s",
							"cursor": 1,
							"created_at": "2026-08-21T10:00:00Z"
						}
					]
				}`, threadID, eventType, idempotencyKey)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, respJSON)
				return
			}
		}
		t.Errorf("unexpected endpoint hit: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	})

	// 1. Test AppendThreadEvent (Repository call)
	appendReq := &thread_domain.AppendThreadEventRequest{
		ThreadID:  threadID,
		EventType: eventType,
		Payload: thread_domain.EventResourceRef{
			RefType:      thread_domain.ResourceShareEntry,
			ShareEntryID: "se_123",
			TrustGroupID: "tg_456",
		},
		IdempotencyKey: idempotencyKey,
	}

	appendedResp, err := client.AppendThreadEvent(context.Background(), appendReq)
	if err != nil {
		t.Fatalf("AppendThreadEvent failed: %v", err)
	}

	if appendedResp == nil || !appendedResp.Success {
		t.Fatalf("expected successful response, got: %v", appendedResp)
	}

	appendedEvt := appendedResp.Data
	if appendedEvt.ID != "evt-cloud-456" {
		t.Errorf("appended event ID = %s, want evt-cloud-456", appendedEvt.ID)
	}
	if appendedEvt.ThreadID != threadID {
		t.Errorf("appended event ThreadID = %s, want %s", appendedEvt.ThreadID, threadID)
	}
	if string(appendedEvt.Type) != eventType {
		t.Errorf("appended event Type = %s, want %s", appendedEvt.Type, eventType)
	}
	if appendedEvt.IdempotencyKey != idempotencyKey {
		t.Errorf("appended event IdempotencyKey = %s, want %s", appendedEvt.IdempotencyKey, idempotencyKey)
	}
	if appendedEvt.Payload.ShareEntryID != "se_123" || appendedEvt.Payload.TrustGroupID != "tg_456" {
		t.Errorf("appended event Payload mismatch: %+v", appendedEvt.Payload)
	}

	// Verify wire payload sent to cloud
	if receivedBody["idempotency_key"] != idempotencyKey {
		t.Errorf("wire payload idempotency_key = %v, want %s", receivedBody["idempotency_key"], idempotencyKey)
	}
	if receivedBody["type"] != eventType {
		t.Errorf("wire payload type = %v, want %s", receivedBody["type"], eventType)
	}

	// 2. Test ListThreadEvents (Repository call)
	listReq := &thread_domain.ListThreadEventsRequest{
		ThreadID: threadID,
	}

	listResp, err := client.ListThreadEvents(context.Background(), listReq)
	if err != nil {
		t.Fatalf("ListThreadEvents failed: %v", err)
	}

	if listResp == nil || !listResp.Success {
		t.Fatalf("expected successful response, got: %v", listResp)
	}

	events := listResp.Data
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	listedEvt := events[0]
	if listedEvt.ID != appendedEvt.ID {
		t.Errorf("listed event ID %s != appended event ID %s", listedEvt.ID, appendedEvt.ID)
	}
	if listedEvt.ThreadID != threadID {
		t.Errorf("listed event ThreadID = %s, want %s", listedEvt.ThreadID, threadID)
	}
	if listedEvt.Payload.ShareEntryID != "se_123" || listedEvt.Payload.TrustGroupID != "tg_456" {
		t.Errorf("listed event Payload mismatch: %+v", listedEvt.Payload)
	}
}

package tracecore_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	channel_domain "vault-app/internal/channel/domain"
	tracecore "vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
)

const cloudChannelJSON = `{
  "status": 200,
  "data": [
    {
      "ID": "ch_7f3c2a",
      "TemplateID": "contract-execution",
      "Title": "contract-execution",
      "Status": "pending",
      "Federation": {
        "VaultAID": "vault_001",
        "VaultBID": "vault_002",
        "AllowedEventTypes": ["entry.shared"],
        "AllowedPaths": null,
        "AllowedDirections": "bidirectional"
      },
      "CreatedAt": "2026-08-14T23:12:53.097Z",
      "UpdatedAt": "2026-08-14T23:12:53.097Z",
      "WorkspaceID": "eb584235-00d1-4982-b087-5c8b4b8942ff"
    }
  ]
}`

func newTestClient(t *testing.T, handler http.HandlerFunc) *tracecore.TracecoreClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return &tracecore.TracecoreClient{
		BaseURL:         server.URL,
		AnkhoraCloudUrl: server.URL,
		Token:           "test-token",
		HTTPClient:      server.Client(),
	}
}

func TestListChannels_MapsCloudChannel(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/channels/workspace/eb584235-00d1-4982-b087-5c8b4b8942ff" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing Bearer token: %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, cloudChannelJSON)
	})

	resp, err := client.ListChannels(context.Background(), &channel_domain.ListChannelsRequest{
		WorkspaceID: "eb584235-00d1-4982-b087-5c8b4b8942ff",
	})
	if err != nil {
		t.Fatalf("ListChannels returned error: %v", err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatal("expected a valid response with channel data")
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 channel, got %d", len(resp.Data))
	}

	ch := resp.Data[0]
	if ch.ID != "ch_7f3c2a" {
		t.Errorf("ID = %q, want ch_7f3c2a", ch.ID)
	}
	if ch.TemplateID != "contract-execution" {
		t.Errorf("TemplateID = %q, want contract-execution", ch.TemplateID)
	}
	if ch.Title != "contract-execution" {
		t.Errorf("Title = %q, want contract-execution", ch.Title)
	}
	if ch.Status != channel_domain.StatusPending {
		t.Errorf("Status = %q, want pending", ch.Status)
	}
	if ch.WorkspaceID != "eb584235-00d1-4982-b087-5c8b4b8942ff" {
		t.Errorf("WorkspaceID = %q, want eb584235-00d1-4982-b087-5c8b4b8942ff", ch.WorkspaceID)
	}
	if ch.CreatedAt.IsZero() {
		t.Error("CreatedAt is zero time")
	}
	if ch.UpdatedAt.IsZero() {
		t.Error("UpdatedAt is zero time")
	}

	if ch.Federation == "" {
		t.Fatal("Federation was lost during mapping")
	}
	var fed tracecore_types.CloudChannelFederation
	if err := json.Unmarshal([]byte(ch.Federation), &fed); err != nil {
		t.Fatalf("Federation is not valid JSON: %v (raw: %s)", err, ch.Federation)
	}
	if fed.VaultAID != "vault_001" {
		t.Errorf("Federation.VaultAID = %q, want vault_001", fed.VaultAID)
	}
	if fed.VaultBID != "vault_002" {
		t.Errorf("Federation.VaultBID = %q, want vault_002", fed.VaultBID)
	}
	if len(fed.AllowedEventTypes) != 1 || fed.AllowedEventTypes[0] != "entry.shared" {
		t.Errorf("Federation.AllowedEventTypes = %v, want [entry.shared]", fed.AllowedEventTypes)
	}
	if fed.AllowedDirections != "bidirectional" {
		t.Errorf("Federation.AllowedDirections = %q, want bidirectional", fed.AllowedDirections)
	}
	if fed.AllowedPaths != nil {
		t.Errorf("Federation.AllowedPaths = %v, want nil", fed.AllowedPaths)
	}
}

func TestListChannels_EmptyFederationObject(t *testing.T) {
	payload := `{
  "status": 200,
  "data": [
    {
      "ID": "ch_empty_fed",
      "TemplateID": "payment",
      "Title": "Payment Channel",
      "Status": "active",
      "Federation": {},
      "CreatedAt": "2026-08-14T23:12:53.097Z",
      "UpdatedAt": "2026-08-14T23:12:53.097Z",
      "WorkspaceID": "ws_1"
    }
  ]
}`

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, payload)
	})

	resp, err := client.ListChannels(context.Background(), &channel_domain.ListChannelsRequest{WorkspaceID: "ws_1"})
	if err != nil {
		t.Fatalf("ListChannels returned error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 channel, got %d", len(resp.Data))
	}
	ch := resp.Data[0]
	if ch.ID != "ch_empty_fed" {
		t.Errorf("ID = %q, want ch_empty_fed", ch.ID)
	}
	if ch.Status != channel_domain.StatusActive {
		t.Errorf("Status = %q, want active", ch.Status)
	}
	// An empty Federation object must be preserved (not silently lost).
	var fed tracecore_types.CloudChannelFederation
	if err := json.Unmarshal([]byte(ch.Federation), &fed); err != nil {
		t.Fatalf("Federation is not valid JSON: %v (raw: %s)", err, ch.Federation)
	}
	if fed.VaultAID != "" || fed.VaultBID != "" || len(fed.AllowedEventTypes) != 0 {
		t.Errorf("unexpected fields in empty Federation: %+v", fed)
	}
}

func TestListChannels_NullFederation(t *testing.T) {
	payload := `{
  "status": 200,
  "data": [
    {
      "ID": "ch_null_fed",
      "TemplateID": "governance",
      "Title": "Budget",
      "Status": "pending",
      "Federation": null,
      "CreatedAt": "2026-08-14T23:12:53.097Z",
      "UpdatedAt": "2026-08-14T23:12:53.097Z",
      "WorkspaceID": "ws_1"
    }
  ]
}`

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, payload)
	})

	resp, err := client.ListChannels(context.Background(), &channel_domain.ListChannelsRequest{WorkspaceID: "ws_1"})
	if err != nil {
		t.Fatalf("ListChannels returned error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 channel, got %d", len(resp.Data))
	}
	if resp.Data[0].Federation != "" {
		t.Errorf("Federation = %q, want empty string for null", resp.Data[0].Federation)
	}
}

func TestListChannels_EmptyListIsValid(t *testing.T) {
	payload := `{ "status": 200, "data": [] }`

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, payload)
	})

	resp, err := client.ListChannels(context.Background(), &channel_domain.ListChannelsRequest{WorkspaceID: "ws_1"})
	if err != nil {
		t.Fatalf("ListChannels returned error: %v", err)
	}
	if resp == nil || resp.Data == nil || len(resp.Data) != 0 {
		t.Fatalf("expected empty channel list, got %+v", resp)
	}
}

func TestListChannels_MalformedElementReturnsError(t *testing.T) {
	payload := `{
  "status": 200,
  "data": [
    { "TemplateID": "contract-execution", "Title": "no-id", "Status": "pending", "WorkspaceID": "ws_1" }
  ]
}`

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, payload)
	})

	resp, err := client.ListChannels(context.Background(), &channel_domain.ListChannelsRequest{WorkspaceID: "ws_1"})
	if err == nil {
		t.Fatal("expected an error for a malformed channel, got nil")
	}
	if resp != nil {
		t.Fatalf("expected nil response on error, got %+v", resp)
	}
}

func TestListChannels_UnexpectedShapeReturnsError(t *testing.T) {
	payload := `{ "status": 200, "data": "not-an-array" }`

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, payload)
	})

	_, err := client.ListChannels(context.Background(), &channel_domain.ListChannelsRequest{WorkspaceID: "ws_1"})
	if err == nil {
		t.Fatal("expected an error for an unexpected response shape, got nil")
	}
}

func TestListChannels_MissingDataReturnsError(t *testing.T) {
	payload := `{ "status": 200, "message": "no data key" }`

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, payload)
	})

	_, err := client.ListChannels(context.Background(), &channel_domain.ListChannelsRequest{WorkspaceID: "ws_1"})
	if err == nil {
		t.Fatal("expected an error when Cloud omits data, got nil")
	}
}

func TestListChannels_HTTPErrorSurfaced(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{ "status": 500, "message": "boom" }`)
	})

	_, err := client.ListChannels(context.Background(), &channel_domain.ListChannelsRequest{WorkspaceID: "ws_1"})
	if err == nil {
		t.Fatal("expected an error for HTTP 500, got nil")
	}
}

func TestCreateChannel_MapsCloudChannel(t *testing.T) {
	created := `{
  "status": 201,
  "data": {
    "ID": "ch_9ab8cd",
    "TemplateID": "contract-execution",
    "Title": "Contract Execution",
    "Status": "pending",
    "Federation": {
      "VaultAID": "vault_001",
      "VaultBID": "vault_002",
      "AllowedEventTypes": ["entry.shared"],
      "AllowedPaths": null,
      "AllowedDirections": "bidirectional"
    },
    "CreatedAt": "2026-08-14T23:12:53.097Z",
    "UpdatedAt": "2026-08-14T23:12:53.097Z",
    "WorkspaceID": "eb584235-00d1-4982-b087-5c8b4b8942ff"
  }
}`

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/channels" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing Bearer token: %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, created)
	})

	resp, err := client.CreateChannel(context.Background(), &channel_domain.CreateChannelRequest{
		Channel: channel_domain.NewChannel("contract-execution", "Contract Execution", "eb584235-00d1-4982-b087-5c8b4b8942ff"),
	})
	if err != nil {
		t.Fatalf("CreateChannel returned error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected a non-nil response")
	}
	ch := resp.Data
	if ch.ID != "ch_9ab8cd" {
		t.Errorf("ID = %q, want ch_9ab8cd (server-created ID must be preserved)", ch.ID)
	}
	if ch.TemplateID != "contract-execution" {
		t.Errorf("TemplateID = %q, want contract-execution", ch.TemplateID)
	}
	if ch.Title != "Contract Execution" {
		t.Errorf("Title = %q, want Contract Execution", ch.Title)
	}
	if ch.Status != channel_domain.StatusPending {
		t.Errorf("Status = %q, want pending", ch.Status)
	}
	if ch.WorkspaceID != "eb584235-00d1-4982-b087-5c8b4b8942ff" {
		t.Errorf("WorkspaceID = %q, want eb584235-00d1-4982-b087-5c8b4b8942ff", ch.WorkspaceID)
	}
	if ch.CreatedAt.IsZero() {
		t.Error("CreatedAt is zero time")
	}
	if ch.UpdatedAt.IsZero() {
		t.Error("UpdatedAt is zero time")
	}

	var fed tracecore_types.CloudChannelFederation
	if err := json.Unmarshal([]byte(ch.Federation), &fed); err != nil {
		t.Fatalf("Federation is not valid JSON: %v (raw: %s)", err, ch.Federation)
	}
	if fed.VaultAID != "vault_001" || fed.VaultBID != "vault_002" {
		t.Errorf("Federation vault IDs not preserved: %+v", fed)
	}
	if len(fed.AllowedEventTypes) != 1 || fed.AllowedEventTypes[0] != "entry.shared" {
		t.Errorf("Federation.AllowedEventTypes = %v, want [entry.shared]", fed.AllowedEventTypes)
	}
	if fed.AllowedDirections != "bidirectional" {
		t.Errorf("Federation.AllowedDirections = %q, want bidirectional", fed.AllowedDirections)
	}
}

func TestCreateChannel_UnexpectedShapeReturnsError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{ "status": 201, "message": "created but no channel body" }`)
	})

	_, err := client.CreateChannel(context.Background(), &channel_domain.CreateChannelRequest{
		Channel: channel_domain.NewChannel("contract-execution", "Contract Execution", "ws_1"),
	})
	if err == nil {
		t.Fatal("expected an error when Cloud omits channel data, got nil")
	}
}

func TestCreateChannel_HTTPErrorSurfaced(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{ "status": 400, "message": "invalid payload" }`)
	})

	_, err := client.CreateChannel(context.Background(), &channel_domain.CreateChannelRequest{
		Channel: channel_domain.NewChannel("contract-execution", "Contract Execution", "ws_1"),
	})
	if err == nil {
		t.Fatal("expected an error for HTTP 400, got nil")
	}
}

const createdChannelJSON = `{
  "status": 201,
  "data": {
    "ID": "ch_slots1",
    "TemplateID": "contract-execution",
    "Title": "Channel A",
    "Status": "pending",
    "Slots": [
      { "ID": "draft", "Name": "contract_draft", "Role": "Author", "VaultID": "vault_legal", "Gated": false, "Order": 0 },
      { "ID": "finance", "Name": "financial_clearance", "Role": "Reviewer", "VaultID": "vault_finance", "Gated": true, "Order": 1 },
      { "ID": "signature", "Name": "executive_signature", "Role": "Approver", "VaultID": "vault_direction", "Gated": true, "Order": 2 }
    ],
    "Assignments": [
      { "SlotID": "draft", "OwnerID": "alice", "PublicKey": "GAALICE", "VaultAddress": "vault_legal" },
      { "SlotID": "finance", "OwnerID": "bob", "PublicKey": "GBBOB", "VaultAddress": "vault_finance" }
    ],
    "Properties": [
      { "Key": "jurisdiction", "Value": "EU" }
    ],
    "Policy": { "max_signatures": 2 },
    "Federation": {
      "VaultAID": "vault_001",
      "VaultBID": "vault_002",
      "AllowedEventTypes": ["entry.shared"],
      "AllowedPaths": null,
      "AllowedDirections": "bidirectional"
    },
    "CreatedAt": "2026-08-15T09:00:00.000Z",
    "UpdatedAt": "2026-08-15T09:00:00.000Z",
    "WorkspaceID": "eb584235-00d1-4982-b087-5c8b4b8942ff"
  }
}`

// channelWithConfig builds a Channel exactly as the create use case would:
// NewChannel + AddSlot/AddAssignment via the aggregate APIs.
func channelWithConfig(t *testing.T) channel_domain.Channel {
	t.Helper()
	ch := channel_domain.NewChannel("contract-execution", "Channel A", "eb584235-00d1-4982-b087-5c8b4b8942ff")
	for _, slot := range []channel_domain.Slot{
		{ID: "draft", Name: "contract_draft", Role: "Author", VaultID: "vault_legal", Gated: false, Order: 0},
		{ID: "finance", Name: "financial_clearance", Role: "Reviewer", VaultID: "vault_finance", Gated: true, Order: 1},
		{ID: "signature", Name: "executive_signature", Role: "Approver", VaultID: "vault_direction", Gated: true, Order: 2},
	} {
		if err := ch.AddSlot(slot); err != nil {
			t.Fatalf("AddSlot(%q) failed: %v", slot.ID, err)
		}
	}
	for _, assignment := range []channel_domain.Assignment{
		{SlotID: "draft", OwnerID: "alice", PublicKey: "GAALICE", VaultAddress: "vault_legal"},
		{SlotID: "finance", OwnerID: "bob", PublicKey: "GBBOB", VaultAddress: "vault_finance"},
	} {
		ch.AddAssignment(assignment)
	}
	return ch
}

func TestCreateChannel_SendsSlotsAndAssignments(t *testing.T) {
	var reqBody []byte
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/channels" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing Bearer token: %q", r.Header.Get("Authorization"))
		}
		var err error
		reqBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, createdChannelJSON)
	})

	ch := channelWithConfig(t)
	ch.Properties = []channel_domain.ChannelProperty{{Key: "jurisdiction", Value: "EU"}}
	ch.Policy = channel_domain.Policy{"max_signatures": 2}
	ch.Federation = `{"VaultAID":"vault_001","VaultBID":"vault_002","AllowedEventTypes":["entry.shared"],"AllowedPaths":null,"AllowedDirections":"bidirectional"}`
	resp, err := client.CreateChannel(context.Background(), &channel_domain.CreateChannelRequest{Channel: ch})
	if err != nil {
		t.Fatalf("CreateChannel returned error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected a non-nil response")
	}

	var body map[string]interface{}
	if err := json.Unmarshal(reqBody, &body); err != nil {
		t.Fatalf("request body is not valid JSON: %v (raw: %s)", err, string(reqBody))
	}

	if body["workspace_id"] != "eb584235-00d1-4982-b087-5c8b4b8942ff" {
		t.Errorf("workspace_id = %v, want eb584235-00d1-4982-b087-5c8b4b8942ff", body["workspace_id"])
	}
	if body["title"] != "Channel A" {
		t.Errorf("title = %v, want Channel A", body["title"])
	}
	if body["template_id"] != "contract-execution" {
		t.Errorf("template_id = %v, want contract-execution", body["template_id"])
	}
	if _, hasStatus := body["status"]; hasStatus {
		t.Error("payload must NOT contain status: Cloud CreateChannelRequest has no status field")
	}

	rawProps, ok := body["properties"].([]interface{})
	if !ok {
		t.Fatalf("properties missing from payload: %s", string(reqBody))
	}
	if len(rawProps) != 1 {
		t.Fatalf("expected 1 property, got %d", len(rawProps))
	}
	prop, ok := rawProps[0].(map[string]interface{})
	if !ok || prop["key"] != "jurisdiction" || prop["value"] != "EU" {
		t.Errorf("property mismatch: %v", prop)
	}

	rawPolicy, ok := body["policy"].(map[string]interface{})
	if !ok {
		t.Fatalf("policy missing from payload: %s", string(reqBody))
	}
	if rawPolicy["max_signatures"] != float64(2) {
		t.Errorf("policy.max_signatures = %v, want 2", rawPolicy["max_signatures"])
	}

	rawFed, ok := body["federation"].(map[string]interface{})
	if !ok {
		t.Fatalf("federation missing from payload: %s", string(reqBody))
	}
	if rawFed["vault_a_id"] != "vault_001" || rawFed["vault_b_id"] != "vault_002" {
		t.Errorf("federation vault ids mismatch: %v", rawFed)
	}
	if got, ok := rawFed["allowed_directions"]; !ok || got != "bidirectional" {
		t.Errorf("federation.allowed_directions = %v, want bidirectional", rawFed["allowed_directions"])
	}
	if got, ok := rawFed["allowed_event_types"].([]interface{}); !ok || len(got) != 1 || got[0] != "entry.shared" {
		t.Errorf("federation.allowed_event_types = %v, want [entry.shared]", rawFed["allowed_event_types"])
	}

	rawSlots, ok := body["slots"].([]interface{})
	if !ok {
		t.Fatalf("slots missing from payload: %s", string(reqBody))
	}
	if len(rawSlots) != 3 {
		t.Fatalf("expected 3 slots, got %d", len(rawSlots))
	}
	wantIDs := []string{"draft", "finance", "signature"}
	for i, raw := range rawSlots {
		slot, ok := raw.(map[string]interface{})
		if !ok {
			t.Fatalf("slot %d is not an object: %v", i, raw)
		}
		if got := slot["id"]; got != wantIDs[i] {
			t.Errorf("slot[%d].id = %v, want %s (semantic ID must be preserved verbatim)", i, got, wantIDs[i])
		}
		if got := slot["order"]; got != float64(i) {
			t.Errorf("slot[%d].order = %v, want %d", i, got, i)
		}
	}
	first := rawSlots[0].(map[string]interface{})
	if first["name"] != "contract_draft" || first["role"] != "Author" || first["vault_id"] != "vault_legal" || first["gated"] != false {
		t.Errorf("slot draft fields mismatch: %v", first)
	}

	rawAssign, ok := body["assignments"].([]interface{})
	if !ok {
		t.Fatalf("assignments missing from payload: %s", string(reqBody))
	}
	if len(rawAssign) != 2 {
		t.Fatalf("expected 2 assignments, got %d", len(rawAssign))
	}
	wantAssign := []map[string]interface{}{
		{"slot_id": "draft", "owner_id": "alice", "public_key": "GAALICE", "vault_address": "vault_legal"},
		{"slot_id": "finance", "owner_id": "bob", "public_key": "GBBOB", "vault_address": "vault_finance"},
	}
	for i, raw := range rawAssign {
		a, ok := raw.(map[string]interface{})
		if !ok {
			t.Fatalf("assignment %d is not an object: %v", i, raw)
		}
		for k, want := range wantAssign[i] {
			if got := a[k]; got != want {
				t.Errorf("assignment[%d].%s = %v, want %v", i, k, got, want)
			}
		}
	}
	if _, hasID := body["assignments"].([]interface{})[0].(map[string]interface{})["id"]; hasID {
		t.Error("assignments must NOT contain a top-level id field")
	}
}

func TestCreateChannel_MapsSlotsAndAssignments(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, createdChannelJSON)
	})

	resp, err := client.CreateChannel(context.Background(), &channel_domain.CreateChannelRequest{
		Channel: channelWithConfig(t),
	})
	if err != nil {
		t.Fatalf("CreateChannel returned error: %v", err)
	}
	ch := resp.Data

	if len(ch.Slots) != 3 {
		t.Fatalf("expected 3 slots mapped from Cloud, got %d", len(ch.Slots))
	}
	for i, wantID := range []string{"draft", "finance", "signature"} {
		s := ch.Slots[i]
		if s.ID != wantID {
			t.Errorf("Slots[%d].ID = %q, want %s", i, s.ID, wantID)
		}
		if s.Order != i {
			t.Errorf("Slots[%d].Order = %d, want %d", i, s.Order, i)
		}
	}
	if ch.Slots[1].Gated != true || ch.Slots[1].VaultID != "vault_finance" {
		t.Errorf("slot finance fields mismatch: %+v", ch.Slots[1])
	}

	if len(ch.Assignments) != 2 {
		t.Fatalf("expected 2 assignments mapped from Cloud, got %d", len(ch.Assignments))
	}
	if ch.Assignments[0].SlotID != "draft" || ch.Assignments[0].OwnerID != "alice" || ch.Assignments[0].PublicKey != "GAALICE" || ch.Assignments[0].VaultAddress != "vault_legal" {
		t.Errorf("assignment draft fields mismatch: %+v", ch.Assignments[0])
	}
	if ch.Assignments[1].SlotID != "finance" {
		t.Errorf("assignment[1].SlotID = %q, want finance", ch.Assignments[1].SlotID)
	}

	if len(ch.Properties) != 1 {
		t.Fatalf("expected 1 property mapped from Cloud, got %d", len(ch.Properties))
	}
	if ch.Properties[0].Key != "jurisdiction" || ch.Properties[0].Value != "EU" {
		t.Errorf("properties not preserved: %+v", ch.Properties)
	}
	if ch.Policy["max_signatures"] != float64(2) {
		t.Errorf("policy not preserved: %+v", ch.Policy)
	}

	var fed tracecore_types.CloudChannelFederation
	if err := json.Unmarshal([]byte(ch.Federation), &fed); err != nil {
		t.Fatalf("Federation is not valid JSON: %v (raw: %s)", err, ch.Federation)
	}
	if fed.VaultAID != "vault_001" || fed.AllowedDirections != "bidirectional" {
		t.Errorf("Federation not preserved: %+v", fed)
	}
}

func TestListChannels_MapsSlotsAndAssignments(t *testing.T) {
	payload := `{
  "status": 200,
  "data": [
    {
      "ID": "ch_list1",
      "TemplateID": "contract-execution",
      "Title": "Channel A",
      "Status": "pending",
      "Slots": [
        { "ID": "draft", "Name": "contract_draft", "Role": "Author", "VaultID": "vault_legal", "Gated": false, "Order": 0 },
        { "ID": "finance", "Name": "financial_clearance", "Role": "Reviewer", "VaultID": "vault_finance", "Gated": true, "Order": 1 }
      ],
      "Assignments": [
        { "SlotID": "draft", "OwnerID": "alice", "PublicKey": "GAALICE", "VaultAddress": "vault_legal" }
      ],
      "Federation": {
        "VaultAID": "vault_001",
        "VaultBID": "vault_002",
        "AllowedEventTypes": null,
        "AllowedPaths": null,
        "AllowedDirections": "unidirectional"
      },
      "CreatedAt": "2026-08-15T09:00:00.000Z",
      "UpdatedAt": "2026-08-15T09:00:00.000Z",
      "WorkspaceID": "eb584235-00d1-4982-b087-5c8b4b8942ff"
    }
  ]
}`

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, payload)
	})

	resp, err := client.ListChannels(context.Background(), &channel_domain.ListChannelsRequest{
		WorkspaceID: "eb584235-00d1-4982-b087-5c8b4b8942ff",
	})
	if err != nil {
		t.Fatalf("ListChannels returned error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 channel, got %d", len(resp.Data))
	}
	ch := resp.Data[0]

	if len(ch.Slots) != 2 {
		t.Fatalf("expected 2 slots mapped from Cloud list, got %d", len(ch.Slots))
	}
	if ch.Slots[0].ID != "draft" || ch.Slots[1].ID != "finance" {
		t.Errorf("slot IDs not preserved: %+v", ch.Slots)
	}
	if len(ch.Assignments) != 1 || ch.Assignments[0].SlotID != "draft" {
		t.Errorf("assignments not preserved: %+v", ch.Assignments)
	}

	var fed tracecore_types.CloudChannelFederation
	if err := json.Unmarshal([]byte(ch.Federation), &fed); err != nil {
		t.Fatalf("Federation is not valid JSON: %v (raw: %s)", err, ch.Federation)
	}
	if fed.AllowedDirections != "unidirectional" {
		t.Errorf("Federation.AllowedDirections = %q, want unidirectional", fed.AllowedDirections)
	}
}

const activatedChannelJSON = `{
  "status": 201,
  "data": {
    "ID": "ch_act1",
    "TemplateID": "contract-execution",
    "Title": "Channel B",
    "Status": "active",
    "Slots": [
      { "ID": "draft", "Name": "contract_draft", "Role": "Author", "VaultID": "vault_legal", "Gated": false, "Order": 0 },
      { "ID": "finance", "Name": "financial_clearance", "Role": "Reviewer", "VaultID": "vault_finance", "Gated": true, "Order": 1 },
      { "ID": "signature", "Name": "executive_signature", "Role": "Approver", "VaultID": "vault_direction", "Gated": true, "Order": 2 }
    ],
    "Assignments": [
      { "SlotID": "draft", "OwnerID": "alice", "PublicKey": "GAALICE", "VaultAddress": "vault_legal" },
      { "SlotID": "finance", "OwnerID": "bob", "PublicKey": "GBBOB", "VaultAddress": "vault_finance" }
    ],
    "Federation": {
      "VaultAID": "vault_001",
      "VaultBID": "vault_002",
      "AllowedEventTypes": ["entry.shared"],
      "AllowedPaths": null,
      "AllowedDirections": "bidirectional"
    },
    "CreatedAt": "2026-08-15T09:00:00.000Z",
    "UpdatedAt": "2026-08-15T09:00:00.000Z",
    "WorkspaceID": "eb584235-00d1-4982-b087-5c8b4b8942ff"
  }
}`

func TestActivateChannel_PostsToActivatePathWithoutBody(t *testing.T) {
	var reqBody []byte
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/channels/ch_act1/activate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing Bearer token: %q", r.Header.Get("Authorization"))
		}
		var err error
		reqBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, activatedChannelJSON)
	})

	resp, err := client.ActivateChannel(context.Background(), &channel_domain.ActivateChannelRequest{
		ChannelID: "ch_act1",
	})
	if err != nil {
		t.Fatalf("ActivateChannel returned error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected a non-nil response")
	}
	// The activation request contains only the channel id; there is no body.
	if len(reqBody) != 0 {
		t.Errorf("expected empty request body, got %q", string(reqBody))
	}
}

func TestActivateChannel_ReturnsActivatedChannel(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, activatedChannelJSON)
	})

	resp, err := client.ActivateChannel(context.Background(), &channel_domain.ActivateChannelRequest{
		ChannelID: "ch_act1",
	})
	if err != nil {
		t.Fatalf("ActivateChannel returned error: %v", err)
	}
	ch := resp.Data
	if ch.ID != "ch_act1" {
		t.Errorf("ID = %q, want ch_act1", ch.ID)
	}
	if ch.Status != channel_domain.StatusActive {
		t.Errorf("Status = %q, want active", ch.Status)
	}
	if len(ch.Slots) != 3 {
		t.Errorf("expected 3 slots mapped from Cloud, got %d", len(ch.Slots))
	}
	if len(ch.Assignments) != 2 || ch.Assignments[0].SlotID != "draft" {
		t.Errorf("assignments not preserved: %+v", ch.Assignments)
	}
	var fed tracecore_types.CloudChannelFederation
	if err := json.Unmarshal([]byte(ch.Federation), &fed); err != nil {
		t.Fatalf("Federation is not valid JSON: %v (raw: %s)", err, ch.Federation)
	}
	if fed.AllowedDirections != "bidirectional" {
		t.Errorf("Federation.AllowedDirections = %q, want bidirectional", fed.AllowedDirections)
	}
}

func TestActivateChannel_AlreadyActiveIsIdempotent(t *testing.T) {
	// Cloud is idempotent: activating an already-active channel returns the
	// channel (status active) without error. The client must respect that.
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, activatedChannelJSON)
	})

	resp, err := client.ActivateChannel(context.Background(), &channel_domain.ActivateChannelRequest{
		ChannelID: "ch_act1",
	})
	if err != nil {
		t.Fatalf("ActivateChannel on an already-active channel must not error: %v", err)
	}
	if resp.Data.Status != channel_domain.StatusActive {
		t.Errorf("Status = %q, want active", resp.Data.Status)
	}
}

func TestActivateChannel_RevokedSurfaced(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `channel revoked`)
	})

	_, err := client.ActivateChannel(context.Background(), &channel_domain.ActivateChannelRequest{
		ChannelID: "ch_revoked",
	})
	if err == nil {
		t.Fatal("expected an error for a revoked channel, got nil")
	}
	if got := err.Error(); got != "Cloud backend returned status 400: channel revoked" {
		t.Errorf("revoked error not surfaced cleanly: %s", got)
	}
}

func TestActivateChannel_GatedSlotsUnfulfilledSurfaced(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `validate activation: gated slots unfulfilled`)
	})

	_, err := client.ActivateChannel(context.Background(), &channel_domain.ActivateChannelRequest{
		ChannelID: "ch_unfulfilled",
	})
	if err == nil {
		t.Fatal("expected an error when gated slots are unfulfilled, got nil")
	}
	if got := err.Error(); got != "Cloud backend returned status 400: validate activation: gated slots unfulfilled" {
		t.Errorf("gated-slots error not surfaced cleanly: %s", got)
	}
}

func TestActivateChannel_HTTPErrorSurfaced(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{ "status": 500, "message": "boom" }`)
	})

	_, err := client.ActivateChannel(context.Background(), &channel_domain.ActivateChannelRequest{
		ChannelID: "ch_5xx",
	})
	if err == nil {
		t.Fatal("expected an error for HTTP 500, got nil")
	}
}

func TestActivateChannel_UnexpectedShapeReturnsError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{ "status": 201, "message": "no channel body" }`)
	})

	_, err := client.ActivateChannel(context.Background(), &channel_domain.ActivateChannelRequest{
		ChannelID: "ch_malformed",
	})
	if err == nil {
		t.Fatal("expected an error when Cloud omits channel data, got nil")
	}
}

func TestActivateChannel_MissingChannelID(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("server should not be called for a missing channel id")
	})

	_, err := client.ActivateChannel(context.Background(), &channel_domain.ActivateChannelRequest{})
	if err == nil {
		t.Fatal("expected an error for a missing channel id, got nil")
	}
}

// ------------------------------------------------------------------------------------------------------------
// RevokeChannel
// ------------------------------------------------------------------------------------------------------------

// The Cloud revoke response is an envelope without Channel data:
// { "status": 200, "code": "RECORD_UPDATED", "message": "...", "success": true }.
const channelRevokedJSON = `{
  "status": 200,
  "code": "RECORD_UPDATED",
  "message": "channel revoked",
  "success": true
}`

func TestRevokeChannel_PostsToRevokePathWithoutBody(t *testing.T) {
	var reqBody []byte
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/channels/ch_rev1/revoke" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing Bearer token: %q", r.Header.Get("Authorization"))
		}
		var err error
		reqBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, channelRevokedJSON)
	})

	err := client.RevokeChannel(context.Background(), &channel_domain.RevokeChannelRequest{
		ChannelID: "ch_rev1",
	})
	if err != nil {
		t.Fatalf("RevokeChannel returned error: %v", err)
	}
	// The revocation request contains only the channel id; there is no body.
	if len(reqBody) != 0 {
		t.Errorf("expected empty request body, got %q", string(reqBody))
	}
}

func TestRevokeChannel_HTTP200WithNilDataSucceeds(t *testing.T) {
	// The Cloud revoke response carries no Channel data. A 200 must still be
	// treated as success without requiring or fabricating a Channel.
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, channelRevokedJSON)
	})

	err := client.RevokeChannel(context.Background(), &channel_domain.RevokeChannelRequest{
		ChannelID: "ch_rev1",
	})
	if err != nil {
		t.Fatalf("RevokeChannel with a data-less 200 must not error: %v", err)
	}
}

func TestRevokeChannel_HTTPErrorSurfaced(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `channel already revoked`)
	})

	err := client.RevokeChannel(context.Background(), &channel_domain.RevokeChannelRequest{
		ChannelID: "ch_rev_already",
	})
	if err == nil {
		t.Fatal("expected an error for HTTP 400, got nil")
	}
	if got := err.Error(); got != "Cloud backend returned status 400: channel already revoked" {
		t.Errorf("Cloud error not surfaced cleanly: %s", got)
	}
}

func TestRevokeChannel_HTTP500Surfaced(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{ "status": 500, "message": "boom" }`)
	})

	err := client.RevokeChannel(context.Background(), &channel_domain.RevokeChannelRequest{
		ChannelID: "ch_rev_5xx",
	})
	if err == nil {
		t.Fatal("expected an error for HTTP 500, got nil")
	}
	if got := err.Error(); got != "Cloud backend returned status 500: { \"status\": 500, \"message\": \"boom\" }" {
		t.Errorf("Cloud error not surfaced cleanly: %s", got)
	}
}

func TestRevokeChannel_MissingChannelID(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("server should not be called for a missing channel id")
	})

	err := client.RevokeChannel(context.Background(), &channel_domain.RevokeChannelRequest{})
	if err == nil {
		t.Fatal("expected an error for a missing channel id, got nil")
	}
}

// ------------------------------------------------------------------------------------------------------------
// GetChannel
// ------------------------------------------------------------------------------------------------------------

const fetchedChannelJSON = `{
  "status": 200,
  "data": {
    "ID": "ch_get1",
    "TemplateID": "contract-execution",
    "Title": "Channel A",
    "Status": "revoked",
    "Slots": [
      { "ID": "draft", "Name": "contract_draft", "Role": "Author", "VaultID": "vault_legal", "Gated": false, "Order": 0 }
    ],
    "Assignments": [
      { "SlotID": "draft", "OwnerID": "alice", "PublicKey": "GAALICE", "VaultAddress": "vault_legal" }
    ],
    "Properties": [
      { "Key": "jurisdiction", "Value": "EU" }
    ],
    "Policy": { "retention_days": 30 },
    "Federation": {
      "VaultAID": "vault_001",
      "VaultBID": "vault_002",
      "AllowedEventTypes": null,
      "AllowedPaths": null,
      "AllowedDirections": "unidirectional"
    },
    "CreatedAt": "2026-08-15T09:00:00.000Z",
    "UpdatedAt": "2026-08-15T09:00:00.000Z",
    "RevokedAt": "2026-08-15T12:00:00.000Z",
    "WorkspaceID": "ws_1"
  }
}`

func TestGetChannel_GetsFromCloudPath(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/channels/ch_get1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing Bearer token: %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, fetchedChannelJSON)
	})

	resp, err := client.GetChannel(context.Background(), &channel_domain.GetChannelRequest{ChannelID: "ch_get1"})
	if err != nil {
		t.Fatalf("GetChannel returned error: %v", err)
	}
	if resp == nil || resp.Data.ID == "" {
		t.Fatal("expected a valid response with channel data")
	}
	ch := resp.Data
	if ch.ID != "ch_get1" {
		t.Errorf("ID = %q, want ch_get1", ch.ID)
	}
	if ch.Status != channel_domain.StatusRevoked {
		t.Errorf("Status = %q, want revoked", ch.Status)
	}
	if ch.RevokedAt == nil || ch.RevokedAt.IsZero() {
		t.Error("RevokedAt must be preserved from Cloud")
	}
	if len(ch.Properties) != 1 || ch.Properties[0].Key != "jurisdiction" || ch.Properties[0].Value != "EU" {
		t.Errorf("Properties not preserved: %+v", ch.Properties)
	}
	if ch.Policy["retention_days"] != float64(30) {
		t.Errorf("Policy not preserved: %+v", ch.Policy)
	}
	var fed tracecore_types.CloudChannelFederation
	if err := json.Unmarshal([]byte(ch.Federation), &fed); err != nil {
		t.Fatalf("Federation is not valid JSON: %v (raw: %s)", err, ch.Federation)
	}
	if fed.AllowedDirections != "unidirectional" {
		t.Errorf("Federation.AllowedDirections = %q, want unidirectional", fed.AllowedDirections)
	}
}

func TestGetChannel_HTTP404Surfaced(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `record not found`)
	})

	_, err := client.GetChannel(context.Background(), &channel_domain.GetChannelRequest{ChannelID: "ch_missing"})
	if err == nil {
		t.Fatal("expected an error for a missing channel, got nil")
	}
	if got := err.Error(); got != "Cloud backend returned status 404: record not found" {
		t.Errorf("Cloud 404 not surfaced verbatim: %s", got)
	}
}

func TestGetChannel_UnexpectedShapeReturnsError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{ "status": 200, "message": "no channel data" }`)
	})

	_, err := client.GetChannel(context.Background(), &channel_domain.GetChannelRequest{ChannelID: "ch_malformed"})
	if err == nil {
		t.Fatal("expected an error when Cloud omits channel data, got nil")
	}
}

func TestGetChannel_MissingChannelID(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("server should not be called for a missing channel id")
	})

	_, err := client.GetChannel(context.Background(), &channel_domain.GetChannelRequest{})
	if err == nil {
		t.Fatal("expected an error for a missing channel id, got nil")
	}
}

// ------------------------------------------------------------------------------------------------------------
// UpdateChannel
// ------------------------------------------------------------------------------------------------------------

const updatedChannelJSON = `{
  "status": 200,
  "data": {
    "ID": "ch_upd1",
    "TemplateID": "contract-execution",
    "Title": "Channel A (renamed)",
    "Status": "active",
    "Slots": [
      { "ID": "draft", "Name": "contract_draft", "Role": "Author", "VaultID": "vault_legal", "Gated": false, "Order": 0 }
    ],
    "Assignments": [
      { "SlotID": "draft", "OwnerID": "alice", "PublicKey": "GAALICE", "VaultAddress": "vault_legal" }
    ],
    "Properties": [
      { "Key": "jurisdiction", "Value": "EU" }
    ],
    "Policy": { "retention_days": 30 },
    "Federation": {
      "VaultAID": "vault_001",
      "VaultBID": "vault_002",
      "AllowedEventTypes": null,
      "AllowedPaths": null,
      "AllowedDirections": "unidirectional"
    },
    "CreatedAt": "2026-08-15T09:00:00.000Z",
    "UpdatedAt": "2026-08-15T10:00:00.000Z",
    "WorkspaceID": "ws_1"
  }
}`

func TestUpdateChannel_PutsExactContractPayload(t *testing.T) {
	var reqBody []byte
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/channels/ch_upd1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing Bearer token: %q", r.Header.Get("Authorization"))
		}
		var err error
		reqBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, updatedChannelJSON)
	})

	ch := channel_domain.Channel{
		ID:         "ch_upd1",
		TemplateID: "contract-execution",
		Title:      "Channel A (renamed)",
		Status:     channel_domain.StatusActive,
		Slots: []channel_domain.Slot{
			{ID: "draft", Name: "contract_draft", Role: "Author", VaultID: "vault_legal", Gated: false, Order: 0},
		},
		Assignments: []channel_domain.Assignment{
			{SlotID: "draft", OwnerID: "alice", PublicKey: "GAALICE", VaultAddress: "vault_legal"},
		},
		Properties:  []channel_domain.ChannelProperty{{Key: "jurisdiction", Value: "EU"}},
		Policy:      channel_domain.Policy{"retention_days": 30},
		WorkspaceID: "ws_1",
	}
	resp, err := client.UpdateChannel(context.Background(), &channel_domain.UpdateChannelRequest{Channel: ch})
	if err != nil {
		t.Fatalf("UpdateChannel returned error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected a non-nil response")
	}

	var body map[string]interface{}
	if err := json.Unmarshal(reqBody, &body); err != nil {
		t.Fatalf("request body is not valid JSON: %v (raw: %s)", err, string(reqBody))
	}

	if body["id"] != "ch_upd1" {
		t.Errorf("id = %v, want ch_upd1", body["id"])
	}
	if body["title"] != "Channel A (renamed)" {
		t.Errorf("title = %v, want Channel A (renamed)", body["title"])
	}
	if rawPolicy, ok := body["policy"].(map[string]interface{}); !ok || rawPolicy["retention_days"] != float64(30) {
		t.Errorf("policy not preserved: %v", body["policy"])
	}
	if rawProps, ok := body["properties"].([]interface{}); !ok || len(rawProps) != 1 {
		t.Errorf("properties not preserved: %v", body["properties"])
	}
	if rawSlots, ok := body["slots"].([]interface{}); !ok || len(rawSlots) != 1 {
		t.Errorf("slots not preserved: %v", body["slots"])
	}
	if rawAssign, ok := body["assignments"].([]interface{}); !ok || len(rawAssign) != 1 {
		t.Errorf("assignments not preserved: %v", body["assignments"])
	}

	// Cloud UpdateChannelRequest supports only id/title/slots/properties/
	// assignments/policy. Fields outside that contract must not leak through.
	for _, forbidden := range []string{"template_id", "status", "federation", "workspace_id", "created_at", "updated_at"} {
		if _, has := body[forbidden]; has {
			t.Errorf("payload must NOT contain %s: %s", forbidden, string(reqBody))
		}
	}
}

func TestUpdateChannel_MapsUpdatedChannel(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, updatedChannelJSON)
	})

	resp, err := client.UpdateChannel(context.Background(), &channel_domain.UpdateChannelRequest{
		Channel: channel_domain.Channel{ID: "ch_upd1", Title: "Channel A (renamed)"},
	})
	if err != nil {
		t.Fatalf("UpdateChannel returned error: %v", err)
	}
	ch := resp.Data
	if ch.ID != "ch_upd1" {
		t.Errorf("ID = %q, want ch_upd1", ch.ID)
	}
	if ch.Title != "Channel A (renamed)" {
		t.Errorf("Title = %q, want Channel A (renamed)", ch.Title)
	}
	if ch.Status != channel_domain.StatusActive {
		t.Errorf("Status = %q, want active", ch.Status)
	}
	if len(ch.Properties) != 1 || ch.Properties[0].Value != "EU" {
		t.Errorf("Properties not preserved: %+v", ch.Properties)
	}
	if ch.Policy["retention_days"] != float64(30) {
		t.Errorf("Policy not preserved: %+v", ch.Policy)
	}
}

func TestUpdateChannel_HTTPErrorSurfaced(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `channel revoked`)
	})

	_, err := client.UpdateChannel(context.Background(), &channel_domain.UpdateChannelRequest{
		Channel: channel_domain.Channel{ID: "ch_revoked", Title: "x"},
	})
	if err == nil {
		t.Fatal("expected an error for HTTP 400, got nil")
	}
	if got := err.Error(); got != "Cloud backend returned status 400: channel revoked" {
		t.Errorf("Cloud error not surfaced cleanly: %s", got)
	}
}

func TestUpdateChannel_MissingChannelID(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("server should not be called for a missing channel id")
	})

	_, err := client.UpdateChannel(context.Background(), &channel_domain.UpdateChannelRequest{})
	if err == nil {
		t.Fatal("expected an error for a missing channel id, got nil")
	}
}

// ------------------------------------------------------------------------------------------------------------
// DeleteChannel
// ------------------------------------------------------------------------------------------------------------

func TestDeleteChannel_DeletesFromCloudPath(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/channels/ch_del1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing Bearer token: %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{ "status": 200, "code": "RECORD_DELETED", "message": "channel deleted", "success": true }`)
	})

	err := client.DeleteChannel(context.Background(), &channel_domain.DeleteChannelRequest{ChannelID: "ch_del1"})
	if err != nil {
		t.Fatalf("DeleteChannel returned error: %v", err)
	}
}

func TestDeleteChannel_HTTPErrorSurfaced(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `record not found`)
	})

	err := client.DeleteChannel(context.Background(), &channel_domain.DeleteChannelRequest{ChannelID: "ch_missing"})
	if err == nil {
		t.Fatal("expected an error for a missing channel, got nil")
	}
	if got := err.Error(); got != "Cloud backend returned status 404: record not found" {
		t.Errorf("Cloud 404 not surfaced verbatim: %s", got)
	}
}

func TestDeleteChannel_HTTP500Surfaced(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{ "status": 500, "message": "boom" }`)
	})

	err := client.DeleteChannel(context.Background(), &channel_domain.DeleteChannelRequest{ChannelID: "ch_5xx"})
	if err == nil {
		t.Fatal("expected an error for HTTP 500, got nil")
	}
}

func TestDeleteChannel_MissingChannelID(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("server should not be called for a missing channel id")
	})

	err := client.DeleteChannel(context.Background(), &channel_domain.DeleteChannelRequest{})
	if err == nil {
		t.Fatal("expected an error for a missing channel id, got nil")
	}
}

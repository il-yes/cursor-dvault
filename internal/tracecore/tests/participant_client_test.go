package tracecore_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	channel_domain "vault-app/internal/channel/domain"
)

// Exact Cloud wire shapes: the Cloud Participant aggregate is marshalled with
// default Go JSON encoding (capitalized field names) inside the success envelope.

const cloudParticipantJSON = `{
  "status": 201,
  "data": {
    "ChannelID": "ch_7f3c2a",
    "VaultID": "vault_supplier",
    "PublicKey": "GABC...",
    "Direction": "bidirectional",
    "JoinedAt": 1750000000,
    "Role": "supplier",
    "Permissions": ["read", "write"]
  },
  "message": "Enregistrement créé avec succès",
  "success": true
}`

const cloudParticipantListJSON = `{
  "status": 200,
  "data": [
    {
      "ChannelID": "ch_7f3c2a",
      "VaultID": "vault_internal",
      "PublicKey": "GINT...",
      "Direction": "inbound",
      "JoinedAt": 1749999000,
      "Role": "owner",
      "Permissions": ["read", "write", "manage"]
    },
    {
      "ChannelID": "ch_7f3c2a",
      "VaultID": "vault_supplier",
      "PublicKey": "GABC...",
      "Direction": "bidirectional",
      "JoinedAt": 1750000000,
      "Role": "supplier",
      "Permissions": ["read"]
    }
  ],
  "message": "Enregistrement recupéré avec succès",
  "success": true
}`

func TestAddParticipant_PostsCorrectBodyAndAuth(t *testing.T) {
	var gotBody map[string]interface{}

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/channels/ch_7f3c2a/participants" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing Bearer token: %q", r.Header.Get("Authorization"))
		}
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, cloudParticipantJSON)
	})

	resp, err := client.AddParticipant(context.Background(), &channel_domain.JoinChannelRequest{
		ChannelID: "ch_7f3c2a",
		VaultID:   "vault_supplier",
		PublicKey: "GABC...",
		Direction: "bidirectional",
		SlotID:    "slot_supplier",
		Role:      "supplier",
	})
	if err != nil {
		t.Fatalf("AddParticipant returned error: %v", err)
	}

	if gotBody["channel_id"] != "ch_7f3c2a" {
		t.Errorf("body channel_id = %v", gotBody["channel_id"])
	}
	if gotBody["vault_id"] != "vault_supplier" {
		t.Errorf("body vault_id = %v", gotBody["vault_id"])
	}
	if gotBody["public_key"] != "GABC..." {
		t.Errorf("body public_key = %v", gotBody["public_key"])
	}
	if gotBody["direction"] != "bidirectional" {
		t.Errorf("body direction = %v", gotBody["direction"])
	}
	if gotBody["slot_id"] != "slot_supplier" {
		t.Errorf("body slot_id = %v", gotBody["slot_id"])
	}
	if gotBody["role"] != "supplier" {
		t.Errorf("body role = %v", gotBody["role"])
	}

	if resp == nil || resp.Data.VaultID == "" {
		t.Fatal("expected a valid response with participant data")
	}
}

func TestAddParticipant_OmitsEmptyOptionalFields(t *testing.T) {
	var gotBody map[string]interface{}

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, cloudParticipantJSON)
	})

	_, err := client.AddParticipant(context.Background(), &channel_domain.JoinChannelRequest{
		ChannelID: "ch_7f3c2a",
		VaultID:   "vault_supplier",
	})
	if err != nil {
		t.Fatalf("AddParticipant returned error: %v", err)
	}

	if _, ok := gotBody["slot_id"]; ok {
		t.Errorf("slot_id should be omitted when empty, got %v", gotBody["slot_id"])
	}
	if _, ok := gotBody["role"]; ok {
		t.Errorf("role should be omitted when empty, got %v", gotBody["role"])
	}
}

func TestAddParticipant_MapsCloudParticipant(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, cloudParticipantJSON)
	})

	resp, err := client.AddParticipant(context.Background(), &channel_domain.JoinChannelRequest{
		ChannelID: "ch_7f3c2a",
		VaultID:   "vault_supplier",
	})
	if err != nil {
		t.Fatalf("AddParticipant returned error: %v", err)
	}

	p := resp.Data
	if p.ChannelID != "ch_7f3c2a" {
		t.Errorf("ChannelID = %q, want ch_7f3c2a", p.ChannelID)
	}
	if p.VaultID != "vault_supplier" {
		t.Errorf("VaultID = %q, want vault_supplier", p.VaultID)
	}
	if p.PublicKey != "GABC..." {
		t.Errorf("PublicKey = %q, want GABC...", p.PublicKey)
	}
	if p.Direction != "bidirectional" {
		t.Errorf("Direction = %q, want bidirectional", p.Direction)
	}
	if p.JoinedAt != 1750000000 {
		t.Errorf("JoinedAt = %d, want 1750000000", p.JoinedAt)
	}
	if p.Role != "supplier" {
		t.Errorf("Role = %q, want supplier", p.Role)
	}
	if len(p.Permissions) != 2 || p.Permissions[0] != "read" || p.Permissions[1] != "write" {
		t.Errorf("Permissions = %v, want [read write]", p.Permissions)
	}
}

func TestAddParticipant_PropagatesCloudError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "channel revoked")
	})

	_, err := client.AddParticipant(context.Background(), &channel_domain.JoinChannelRequest{
		ChannelID: "ch_7f3c2a",
		VaultID:   "vault_supplier",
	})
	if err == nil {
		t.Fatal("expected error for HTTP 400")
	}
	want := "Cloud backend returned status 400: channel revoked"
	if err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}

func TestAddParticipant_RejectsMalformedResponse(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": 201, "data": {}}`)
	})

	_, err := client.AddParticipant(context.Background(), &channel_domain.JoinChannelRequest{
		ChannelID: "ch_7f3c2a",
		VaultID:   "vault_supplier",
	})
	if err == nil {
		t.Fatal("expected error for a response without participant data")
	}
}

func TestAddParticipant_MissingChannelID(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("request should not have been issued")
	})

	_, err := client.AddParticipant(context.Background(), &channel_domain.JoinChannelRequest{
		VaultID: "vault_supplier",
	})
	if err == nil {
		t.Fatal("expected channel id error")
	}
}

func TestListParticipants_MapsCloudParticipants(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/channels/ch_7f3c2a/participants" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing Bearer token: %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, cloudParticipantListJSON)
	})

	resp, err := client.ListParticipants(context.Background(), &channel_domain.ListParticipantsRequest{
		ChannelID: "ch_7f3c2a",
	})
	if err != nil {
		t.Fatalf("ListParticipants returned error: %v", err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatal("expected a valid response with participant data")
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 participants, got %d", len(resp.Data))
	}

	first := resp.Data[0]
	if first.VaultID != "vault_internal" {
		t.Errorf("VaultID = %q, want vault_internal", first.VaultID)
	}
	if first.Role != "owner" {
		t.Errorf("Role = %q, want owner", first.Role)
	}
	if len(first.Permissions) != 3 || first.Permissions[2] != "manage" {
		t.Errorf("Permissions = %v, want [read write manage]", first.Permissions)
	}
	if first.PublicKey != "GINT..." {
		t.Errorf("PublicKey = %q, want GINT...", first.PublicKey)
	}
	if first.JoinedAt != 1749999000 {
		t.Errorf("JoinedAt = %d, want 1749999000", first.JoinedAt)
	}

	second := resp.Data[1]
	if second.VaultID != "vault_supplier" {
		t.Errorf("VaultID = %q, want vault_supplier", second.VaultID)
	}
	if second.Role != "supplier" {
		t.Errorf("Role = %q, want supplier", second.Role)
	}
	if len(second.Permissions) != 1 || second.Permissions[0] != "read" {
		t.Errorf("Permissions = %v, want [read]", second.Permissions)
	}
	if second.Direction != "bidirectional" {
		t.Errorf("Direction = %q, want bidirectional", second.Direction)
	}
}

func TestListParticipants_EmptyListAccepted(t *testing.T) {
	// Cloud marshals a nil participant slice as null in the envelope.
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": 200, "data": null, "message": "ok", "success": true}`)
	})

	resp, err := client.ListParticipants(context.Background(), &channel_domain.ListParticipantsRequest{
		ChannelID: "ch_7f3c2a",
	})
	if err != nil {
		t.Fatalf("ListParticipants returned error for empty list: %v", err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatal("expected an empty, non-nil participant slice")
	}
	if len(resp.Data) != 0 {
		t.Fatalf("expected 0 participants, got %d", len(resp.Data))
	}
}

func TestListParticipants_EmptyListEnvelopeAccepted(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": 200, "data": [], "message": "ok", "success": true}`)
	})

	resp, err := client.ListParticipants(context.Background(), &channel_domain.ListParticipantsRequest{
		ChannelID: "ch_7f3c2a",
	})
	if err != nil {
		t.Fatalf("ListParticipants returned error for empty list: %v", err)
	}
	if resp == nil || resp.Data == nil || len(resp.Data) != 0 {
		t.Fatalf("expected an empty participant slice, got %+v", resp)
	}
}

func TestListParticipants_RejectsMalformedElement(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": 200, "data": [{"ChannelID": "ch_7f3c2a", "VaultID": "ok", "Role": "owner"}, {"ChannelID": "", "VaultID": ""}]}`)
	})

	_, err := client.ListParticipants(context.Background(), &channel_domain.ListParticipantsRequest{
		ChannelID: "ch_7f3c2a",
	})
	if err == nil {
		t.Fatal("expected error for a malformed participant element")
	}
}

func TestListParticipants_PropagatesCloudError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "channel not found")
	})

	_, err := client.ListParticipants(context.Background(), &channel_domain.ListParticipantsRequest{
		ChannelID: "missing",
	})
	if err == nil {
		t.Fatal("expected error for HTTP 400")
	}
	want := "Cloud backend returned status 400: channel not found"
	if err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}

func TestListParticipants_MissingChannelID(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("request should not have been issued")
	})

	_, err := client.ListParticipants(context.Background(), &channel_domain.ListParticipantsRequest{})
	if err == nil {
		t.Fatal("expected channel id error")
	}
}

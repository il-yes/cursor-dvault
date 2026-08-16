package tracecore_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	channel_domain "vault-app/internal/channel/domain"
)

// Exact Cloud wire shapes: the Cloud Invitation aggregate is marshalled with
// default Go JSON encoding (capitalized field names) inside the success
// envelope. AcceptedAt is null while the invitation is pending.

const cloudInvitationPendingJSON = `{
  "status": 201,
  "data": {
    "ID": "inv_9e24294d",
    "ChannelID": "ch_inv1",
    "InviterVaultID": "vault_alice",
    "InviteeVaultID": "vault_bob",
    "Status": "pending",
    "CreatedAt": "2026-08-15T10:00:00Z",
    "AcceptedAt": null
  },
  "message": "Enregistrement créé avec succès",
  "success": true
}`

const cloudInvitationAcceptedJSON = `{
  "status": 200,
  "data": {
    "ID": "inv_9e24294d",
    "ChannelID": "ch_inv1",
    "InviterVaultID": "vault_alice",
    "InviteeVaultID": "vault_bob",
    "Status": "accepted",
    "CreatedAt": "2026-08-15T10:00:00Z",
    "AcceptedAt": "2026-08-15T11:00:00Z"
  },
  "message": "Enregistrement recupéré avec succès",
  "success": true
}`

func TestInviteToChannel_PostsCorrectBodyAndAuth(t *testing.T) {
	var gotBody map[string]interface{}

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/channels/ch_inv1/invitations" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing Bearer token: %q", r.Header.Get("Authorization"))
		}
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, cloudInvitationPendingJSON)
	})

	resp, err := client.InviteToChannel(context.Background(), &channel_domain.InviteToChannelRequest{
		ChannelID:      "ch_inv1",
		InviterVaultID: "vault_alice",
		InviteeVaultID: "vault_bob",
	})
	if err != nil {
		t.Fatalf("InviteToChannel returned error: %v", err)
	}

	if gotBody["channel_id"] != "ch_inv1" {
		t.Errorf("body channel_id = %v", gotBody["channel_id"])
	}
	if gotBody["inviter_vault_id"] != "vault_alice" {
		t.Errorf("body inviter_vault_id = %v", gotBody["inviter_vault_id"])
	}
	if gotBody["invitee_vault_id"] != "vault_bob" {
		t.Errorf("body invitee_vault_id = %v", gotBody["invitee_vault_id"])
	}

	if resp == nil || resp.Data.ID == "" {
		t.Fatal("expected a valid response with invitation data")
	}
}

func TestInviteToChannel_MapsCloudInvitation(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, cloudInvitationPendingJSON)
	})

	resp, err := client.InviteToChannel(context.Background(), &channel_domain.InviteToChannelRequest{
		ChannelID:      "ch_inv1",
		InviterVaultID: "vault_alice",
		InviteeVaultID: "vault_bob",
	})
	if err != nil {
		t.Fatalf("InviteToChannel returned error: %v", err)
	}

	inv := resp.Data
	if inv.ID != "inv_9e24294d" {
		t.Errorf("ID = %q, want inv_9e24294d", inv.ID)
	}
	if inv.ChannelID != "ch_inv1" {
		t.Errorf("ChannelID = %q, want ch_inv1", inv.ChannelID)
	}
	if inv.InviterVaultID != "vault_alice" {
		t.Errorf("InviterVaultID = %q, want vault_alice", inv.InviterVaultID)
	}
	if inv.InviteeVaultID != "vault_bob" {
		t.Errorf("InviteeVaultID = %q, want vault_bob", inv.InviteeVaultID)
	}
	if inv.Status != channel_domain.InvitationStatusPending {
		t.Errorf("Status = %q, want pending", inv.Status)
	}
	if inv.AcceptedAt != nil {
		t.Errorf("AcceptedAt = %v, want nil while pending", inv.AcceptedAt)
	}
	if inv.CreatedAt.IsZero() {
		t.Error("CreatedAt should be parsed")
	}
}

func TestInviteToChannel_PropagatesCloudError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "record not found")
	})

	_, err := client.InviteToChannel(context.Background(), &channel_domain.InviteToChannelRequest{
		ChannelID:      "missing",
		InviterVaultID: "vault_alice",
		InviteeVaultID: "vault_bob",
	})
	if err == nil {
		t.Fatal("expected error for HTTP 400")
	}
	want := "Cloud backend returned status 400: record not found"
	if err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}

func TestInviteToChannel_RejectsMalformedResponse(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": 201, "data": {}}`)
	})

	_, err := client.InviteToChannel(context.Background(), &channel_domain.InviteToChannelRequest{
		ChannelID:      "ch_inv1",
		InviterVaultID: "vault_alice",
		InviteeVaultID: "vault_bob",
	})
	if err == nil {
		t.Fatal("expected error for a response without invitation data")
	}
}

func TestInviteToChannel_RequiresFields(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("request should not have been issued")
	})

	cases := []struct {
		name string
		req  *channel_domain.InviteToChannelRequest
	}{
		{"nil request", nil},
		{"missing channel", &channel_domain.InviteToChannelRequest{InviterVaultID: "vault_alice", InviteeVaultID: "vault_bob"}},
		{"missing inviter", &channel_domain.InviteToChannelRequest{ChannelID: "ch_inv1", InviteeVaultID: "vault_bob"}},
		{"missing invitee", &channel_domain.InviteToChannelRequest{ChannelID: "ch_inv1", InviterVaultID: "vault_alice"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.InviteToChannel(context.Background(), tc.req)
			if err == nil {
				t.Fatal("expected a validation error")
			}
		})
	}
}

func TestAcceptChannelInvitation_PostsCorrectBodyAndAuth(t *testing.T) {
	var gotBody map[string]interface{}

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/channels/invitations/inv_9e24294d/accept" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("missing Bearer token: %q", r.Header.Get("Authorization"))
		}
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, cloudInvitationAcceptedJSON)
	})

	resp, err := client.AcceptChannelInvitation(context.Background(), &channel_domain.AcceptInvitationRequest{
		InvitationID:     "inv_9e24294d",
		InviteeVaultID:   "vault_bob",
		InviteePublicKey: "GBOB...",
	})
	if err != nil {
		t.Fatalf("AcceptChannelInvitation returned error: %v", err)
	}

	if gotBody["invitation_id"] != "inv_9e24294d" {
		t.Errorf("body invitation_id = %v", gotBody["invitation_id"])
	}
	if gotBody["invitee_vault_id"] != "vault_bob" {
		t.Errorf("body invitee_vault_id = %v", gotBody["invitee_vault_id"])
	}
	if gotBody["invitee_public_key"] != "GBOB..." {
		t.Errorf("body invitee_public_key = %v", gotBody["invitee_public_key"])
	}

	if resp == nil || resp.Data.ID == "" {
		t.Fatal("expected a valid response with invitation data")
	}
}

func TestAcceptChannelInvitation_MapsAcceptedInvitation(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, cloudInvitationAcceptedJSON)
	})

	resp, err := client.AcceptChannelInvitation(context.Background(), &channel_domain.AcceptInvitationRequest{
		InvitationID:     "inv_9e24294d",
		InviteeVaultID:   "vault_bob",
		InviteePublicKey: "GBOB...",
	})
	if err != nil {
		t.Fatalf("AcceptChannelInvitation returned error: %v", err)
	}

	inv := resp.Data
	if inv.ID != "inv_9e24294d" {
		t.Errorf("ID = %q, want inv_9e24294d", inv.ID)
	}
	if inv.Status != channel_domain.InvitationStatusAccepted {
		t.Errorf("Status = %q, want accepted", inv.Status)
	}
	if inv.AcceptedAt == nil {
		t.Fatal("AcceptedAt should be set on an accepted invitation")
	}
	want := time.Date(2026, 8, 15, 11, 0, 0, 0, time.UTC)
	if !inv.AcceptedAt.Equal(want) {
		t.Errorf("AcceptedAt = %v, want %v", inv.AcceptedAt, want)
	}
	if inv.InviteeVaultID != "vault_bob" {
		t.Errorf("InviteeVaultID = %q, want vault_bob", inv.InviteeVaultID)
	}
}

func TestAcceptChannelInvitation_PropagatesCloudError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "invitation not for you")
	})

	_, err := client.AcceptChannelInvitation(context.Background(), &channel_domain.AcceptInvitationRequest{
		InvitationID:     "inv_9e24294d",
		InviteeVaultID:   "vault_mallory",
		InviteePublicKey: "GMAL...",
	})
	if err == nil {
		t.Fatal("expected error for HTTP 400")
	}
	want := "Cloud backend returned status 400: invitation not for you"
	if err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}

func TestAcceptChannelInvitation_RejectsMalformedResponse(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": 200, "data": {}}`)
	})

	_, err := client.AcceptChannelInvitation(context.Background(), &channel_domain.AcceptInvitationRequest{
		InvitationID:     "inv_9e24294d",
		InviteeVaultID:   "vault_bob",
		InviteePublicKey: "GBOB...",
	})
	if err == nil {
		t.Fatal("expected error for a response without invitation data")
	}
}

func TestAcceptChannelInvitation_RequiresFields(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("request should not have been issued")
	})

	cases := []struct {
		name string
		req  *channel_domain.AcceptInvitationRequest
	}{
		{"nil request", nil},
		{"missing invitation", &channel_domain.AcceptInvitationRequest{InviteeVaultID: "vault_bob", InviteePublicKey: "GBOB..."}},
		{"missing invitee", &channel_domain.AcceptInvitationRequest{InvitationID: "inv_9e24294d", InviteePublicKey: "GBOB..."}},
		{"missing public key", &channel_domain.AcceptInvitationRequest{InvitationID: "inv_9e24294d", InviteeVaultID: "vault_bob"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.AcceptChannelInvitation(context.Background(), tc.req)
			if err == nil {
				t.Fatal("expected a validation error")
			}
		})
	}
}

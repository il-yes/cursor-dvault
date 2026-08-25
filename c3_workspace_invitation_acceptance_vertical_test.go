package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	auth_usecases "vault-app/internal/auth/application/use_cases"
	auth_domain "vault-app/internal/auth/domain"
	auth_ui "vault-app/internal/auth/ui"
	channel_usecase "vault-app/internal/channel/application/channel_lifecycle_usecases"
	channel_domain "vault-app/internal/channel/domain"
	channel_ui "vault-app/internal/channel/ui"
	"vault-app/internal/tracecore"
	vault_ui "vault-app/internal/vault/ui"
)

type workspaceInvitationCloudStub struct {
	mu           sync.Mutex
	token        string
	invitations  map[string]channel_domain.Invitation
	participants map[string][]channel_domain.Participant
}

func newWorkspaceInvitationCloudStub(token string) *workspaceInvitationCloudStub {
	return &workspaceInvitationCloudStub{
		token:        token,
		invitations:  make(map[string]channel_domain.Invitation),
		participants: make(map[string][]channel_domain.Participant),
	}
}

func (s *workspaceInvitationCloudStub) seedInvitation(inv channel_domain.Invitation) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.invitations[inv.ID] = inv
}

func (s *workspaceInvitationCloudStub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer "+s.token {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// POST /channels/invitations/{id}/accept
	if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/channels/invitations/") && strings.HasSuffix(r.URL.Path, "/accept") {
		trimmed := strings.TrimPrefix(r.URL.Path, "/channels/invitations/")
		invID := strings.TrimSuffix(trimmed, "/accept")

		var reqBody map[string]string
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		inv, exists := s.invitations[invID]
		if !exists {
			http.Error(w, "invitation not found", http.StatusNotFound)
			return
		}

		inviteeVaultID := reqBody["invitee_vault_id"]
		inviteePubKey := reqBody["invitee_public_key"]

		if inviteeVaultID != inv.InviteeVaultID {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("invitation not for you"))
			return
		}

		now := time.Now().UTC()
		inv.Status = channel_domain.InvitationStatusAccepted
		inv.AcceptedAt = &now
		s.invitations[invID] = inv

		// Automatically create Participant on Cloud side
		pList := s.participants[inv.ChannelID]
		alreadyExists := false
		for _, p := range pList {
			if p.VaultID == inviteeVaultID {
				alreadyExists = true
				break
			}
		}
		if !alreadyExists {
			s.participants[inv.ChannelID] = append(s.participants[inv.ChannelID], channel_domain.Participant{
				ChannelID:   inv.ChannelID,
				VaultID:     inviteeVaultID,
				PublicKey:   inviteePubKey,
				Direction:   "bidirectional",
				JoinedAt:    now.Unix(),
				Role:        "",
				Permissions: []string{},
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status": 200,
			"data": map[string]interface{}{
				"ID":             inv.ID,
				"ChannelID":      inv.ChannelID,
				"InviterVaultID": inv.InviterVaultID,
				"InviteeVaultID": inv.InviteeVaultID,
				"Status":         string(inv.Status),
				"CreatedAt":      inv.CreatedAt.Format(time.RFC3339),
				"AcceptedAt":     now.Format(time.RFC3339),
			},
			"message": "Enregistrement recupéré avec succès",
			"success": true,
		})
		return
	}

	http.Error(w, "not found", http.StatusNotFound)
}

func TestWorkspaceInvitationAcceptance_VerticalFlow(t *testing.T) {
	ctx := context.Background()

	cloudToken := "test-bearer-token-1234567890"
	stub := newWorkspaceInvitationCloudStub(cloudToken)

	invID := "inv_ws_test_9e24"
	channelID := "ch_ws_test_001"
	inviteeVaultID := "vault_invitee_bob"
	inviteePubKey := "GBOBSTELLARPUBLICKEY..."

	stub.seedInvitation(channel_domain.Invitation{
		ID:             invID,
		ChannelID:      channelID,
		InviterVaultID: "vault_owner_alice",
		InviteeVaultID: inviteeVaultID,
		Status:         channel_domain.InvitationStatusPending,
		CreatedAt:      time.Now().UTC(),
	})

	ts := httptest.NewServer(stub)
	defer ts.Close()

	tc := tracecore.NewTracecoreClient(ts.URL, cloudToken, ts.URL, ts.URL)

	acceptUC := channel_usecase.NewAcceptChannelInvitationUsecase(tc)
	channelHdlr := channel_ui.NewChannelHandler(
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, acceptUC,
	)

	authV2 := auth_domain.Auth{
		Issuer:      "ankhora-test",
		Audience:    "ankhora-desktop",
		Secret:      "test-secret-key-ws-inv",
		TokenExpiry: 15 * time.Minute,
	}
	tokenService := auth_usecases.NewTokenService(authV2, nil, nil)
	tokenUC := auth_usecases.NewGenerateTokensUseCase(nil, tokenService)
	authHdlr := auth_ui.NewAuthHandler(nil, tokenUC, nil)

	vaultHdlr := &vault_ui.VaultHandler{
		TracecoreClient: tc,
	}

	app := &App{
		AuthHandler:    authHdlr,
		ChannelHandler: channelHdlr,
		Vault:          vaultHdlr,
		ctx:            ctx,
	}

	pairs, err := tokenService.GenerateTokenPair(&auth_domain.JwtUser{
		ID:       "user_bob",
		Username: "bob",
		Email:    "bob@ankhora.test",
	})
	require.NoError(t, err)

	// Execute AcceptChannelInvitation through App binding
	acceptedInv, err := app.AcceptChannelInvitation(pairs.Token, invID, inviteeVaultID, inviteePubKey)
	require.NoError(t, err)
	require.NotNil(t, acceptedInv)

	// Assert Invitation is accepted
	assert.Equal(t, invID, acceptedInv.ID)
	assert.Equal(t, string(channel_domain.InvitationStatusAccepted), acceptedInv.Status)
	assert.NotNil(t, acceptedInv.AcceptedAt)

	// Assert Cloud created the Channel Participant
	stub.mu.Lock()
	participants := stub.participants[channelID]
	stub.mu.Unlock()

	require.Len(t, participants, 1)
	assert.Equal(t, inviteeVaultID, participants[0].VaultID)
	assert.Equal(t, inviteePubKey, participants[0].PublicKey)
	assert.Equal(t, "bidirectional", participants[0].Direction)

	// Negative test: Accept with unauthorized invitee vault ID returns Cloud 400 error
	_, err = app.AcceptChannelInvitation(pairs.Token, invID, "vault_mallory", "GMALLORY...")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invitation not for you")
}

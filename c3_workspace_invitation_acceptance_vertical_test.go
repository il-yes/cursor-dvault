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
	tracecore_types "vault-app/internal/tracecore/types"
	vault_ui "vault-app/internal/vault/ui"
	workspace_usecase "vault-app/internal/workspace/application/usecases"
	workspace_ui "vault-app/internal/workspace/ui"
)

type workspaceInvitationCloudStub struct {
	mu           sync.Mutex
	token        string
	workspaces   map[string]tracecore_types.CloudWorkspaceDTO
	channels     map[string][]tracecore_types.CloudChannelDTO
	invitations  map[string]channel_domain.Invitation
	participants map[string][]tracecore_types.CloudChannelParticipant
}

func newWorkspaceInvitationCloudStub(token string) *workspaceInvitationCloudStub {
	return &workspaceInvitationCloudStub{
		token:        token,
		workspaces:   make(map[string]tracecore_types.CloudWorkspaceDTO),
		channels:     make(map[string][]tracecore_types.CloudChannelDTO),
		invitations:  make(map[string]channel_domain.Invitation),
		participants: make(map[string][]tracecore_types.CloudChannelParticipant),
	}
}

func (s *workspaceInvitationCloudStub) seedInvitation(inv channel_domain.Invitation) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.invitations[inv.ID] = inv
}

func (s *workspaceInvitationCloudStub) seedWorkspace(ws tracecore_types.CloudWorkspaceDTO) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.workspaces[ws.ID] = ws
}

func (s *workspaceInvitationCloudStub) seedChannel(ch tracecore_types.CloudChannelDTO) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.channels[ch.WorkspaceID] = append(s.channels[ch.WorkspaceID], ch)
}

func (s *workspaceInvitationCloudStub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer "+s.token {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// GET /workspaces?vault_id={vaultID} or GET /api/workspaces?vault_id={vaultID}
	if r.Method == http.MethodGet && (strings.HasPrefix(r.URL.Path, "/workspaces") || strings.HasPrefix(r.URL.Path, "/api/workspaces")) {
		vaultID := r.URL.Query().Get("vault_id")
		result := make([]tracecore_types.CloudWorkspaceDTO, 0)
		for _, ws := range s.workspaces {
			// User/vault is owner OR is participant in a channel of this workspace
			isOwner := ws.VaultID == vaultID || ws.OwnerID == vaultID
			isParticipant := false
			for wsID, chList := range s.channels {
				if wsID == ws.ID {
					for _, ch := range chList {
						for _, p := range s.participants[ch.ID] {
							if p.VaultID == vaultID {
								isParticipant = true
								break
							}
						}
					}
				}
			}
			if isOwner || isParticipant {
				result = append(result, ws)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  200,
			"data":    result,
			"message": "success",
			"success": true,
		})
		return
	}

	// GET /channels/workspace/{workspaceID}
	if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/channels/workspace/") {
		wsID := strings.TrimPrefix(r.URL.Path, "/channels/workspace/")
		chList := s.channels[wsID]

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  200,
			"data":    chList,
			"message": "success",
			"success": true,
		})
		return
	}

	// GET /channels/{channelID}/participants
	if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/channels/") && strings.HasSuffix(r.URL.Path, "/participants") {
		trimmed := strings.TrimPrefix(r.URL.Path, "/channels/")
		chID := strings.TrimSuffix(trimmed, "/participants")
		pList := s.participants[chID]

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  200,
			"data":    pList,
			"message": "success",
			"success": true,
		})
		return
	}

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
			s.participants[inv.ChannelID] = append(s.participants[inv.ChannelID], tracecore_types.CloudChannelParticipant{
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
	wsID := "ws_test_8f3c"
	channelID := "ch_ws_test_001"
	inviteeVaultID := "vault_invitee_bob"
	inviteePubKey := "GBOBSTELLARPUBLICKEY..."

	// Seed Cloud stub with initial Workspace, Channel, and pending Invitation
	stub.seedWorkspace(tracecore_types.CloudWorkspaceDTO{
		ID:          wsID,
		VaultID:     "vault_owner_alice",
		Name:        "Sovereign Collaboration Workspace",
		Description: "Workspace for cross-organization engineering",
		Status:      "active",
		OwnerID:     "user_alice",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})

	stub.seedChannel(tracecore_types.CloudChannelDTO{
		ID:          channelID,
		Title:       "General Architecture Channel",
		Status:      "active",
		WorkspaceID: wsID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})

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

	listChannelUC := channel_usecase.NewListChannelUsecase(tc)
	listParticipantsUC := channel_usecase.NewListParticipantsUsecase(tc)
	acceptUC := channel_usecase.NewAcceptChannelInvitationUsecase(tc)
	channelHdlr := channel_ui.NewChannelHandler(
		nil, listChannelUC, nil, nil, nil, nil, nil, nil, listParticipantsUC, nil, acceptUC,
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

	workspaceHdlr := workspace_ui.NewWorkspaceHandler(
		nil, workspace_usecase.NewListWorkspaceUsecase(tc, nil),
	)

	app := &App{
		AuthHandler:      authHdlr,
		ChannelHandler:   channelHdlr,
		Vault:            vaultHdlr,
		WorkspaceHandler: workspaceHdlr,
		ctx:              ctx,
	}

	pairs, err := tokenService.GenerateTokenPair(&auth_domain.JwtUser{
		ID:       "user_bob",
		Username: "bob",
		Email:    "bob@ankhora.test",
	})
	require.NoError(t, err)

	// Pre-condition check: Invitee cannot see workspace before accepting invitation
	preWorkspaces, err := app.ListWorkspaces(pairs.Token, inviteeVaultID)
	require.NoError(t, err)
	assert.Empty(t, preWorkspaces, "Invitee MUST NOT see workspace before accepting invitation")

	// Step 1: Execute AcceptChannelInvitation through App binding
	acceptedInv, err := app.AcceptChannelInvitation(pairs.Token, invID, inviteeVaultID, inviteePubKey)
	require.NoError(t, err)
	require.NotNil(t, acceptedInv)

	// Assert Invitation is accepted
	assert.Equal(t, invID, acceptedInv.ID)
	assert.Equal(t, string(channel_domain.InvitationStatusAccepted), acceptedInv.Status)
	assert.NotNil(t, acceptedInv.AcceptedAt)

	// Step 2: Verify Cloud created the Channel Participant
	stub.mu.Lock()
	participants := stub.participants[channelID]
	stub.mu.Unlock()

	require.Len(t, participants, 1)
	assert.Equal(t, inviteeVaultID, participants[0].VaultID)
	assert.Equal(t, inviteePubKey, participants[0].PublicKey)
	assert.Equal(t, "bidirectional", participants[0].Direction)

	// Step 3: Authoritative Read Flow — ListWorkspaces for Invitee Vault
	postWorkspaces, err := app.ListWorkspaces(pairs.Token, inviteeVaultID)
	require.NoError(t, err)
	require.Len(t, postWorkspaces, 1, "Authoritative read flow MUST return accepted workspace for invitee vault")
	assert.Equal(t, wsID, postWorkspaces[0].ID)
	assert.Equal(t, "Sovereign Collaboration Workspace", postWorkspaces[0].Name)

	// Step 4: Authoritative Read Flow — ListChannels for Workspace
	postChannels, err := app.ListChannels(pairs.Token, wsID)
	require.NoError(t, err)
	require.Len(t, postChannels, 1, "Authoritative read flow MUST return accessible channel in accepted workspace")
	assert.Equal(t, channelID, postChannels[0].ID)

	// Step 5: Authoritative Read Flow — ListParticipants for Channel
	postParticipants, err := app.ListParticipants(pairs.Token, channelID)
	require.NoError(t, err)
	require.Len(t, postParticipants, 1, "Authoritative read flow MUST return invitee vault as channel participant")
	assert.Equal(t, inviteeVaultID, postParticipants[0].VaultID)

	// Negative test: Accept with unauthorized invitee vault ID returns Cloud 400 error
	_, err = app.AcceptChannelInvitation(pairs.Token, invID, "vault_mallory", "GMALLORY...")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invitation not for you")
}

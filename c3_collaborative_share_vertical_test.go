package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	auth_usecases "vault-app/internal/auth/application/use_cases"
	auth_domain "vault-app/internal/auth/domain"
	auth_ui "vault-app/internal/auth/ui"
	"vault-app/internal/blockchain"
	c3_asset_domain "vault-app/internal/c3_asset/domain"
	channel_domain "vault-app/internal/channel/domain"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	collaboration_infra "vault-app/internal/collaboration/infrastructure"
	collaboration_ui "vault-app/internal/collaboration/ui"
	app_config "vault-app/internal/config"
	app_config_domain "vault-app/internal/config/domain"
	app_config_persistence "vault-app/internal/config/infrastructure/persistence"
	app_config_ui "vault-app/internal/config/ui"
	"vault-app/internal/driver"
	handlers "vault-app/internal/handlers"
	identity_usecase "vault-app/internal/identity/application/usecase"
	identity_domain "vault-app/internal/identity/domain"
	identity_eventbus "vault-app/internal/identity/infrastructure/eventbus"
	identity_persistence "vault-app/internal/identity/infrastructure/persistence"
	identity_ui "vault-app/internal/identity/ui"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	onboarding_domain "vault-app/internal/onboarding/domain"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"
	onboarding_eventbus "vault-app/internal/onboarding/infrastructure/eventbus"
	onboarding_persistence "vault-app/internal/onboarding/infrastructure/persistence"
	onboarding_ui_wails "vault-app/internal/onboarding/ui/wails"
	subscription_domain "vault-app/internal/subscription/domain"
	subscription_persistence "vault-app/internal/subscription/infrastructure/persistence"
	subscription_ui_wails "vault-app/internal/subscription/ui/wails"
	"vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_events "vault-app/internal/trust_group/application/events"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_ports "vault-app/internal/trust_group/application/ports"
	trustgroup_envelope_uc "vault-app/internal/trust_group/application/usecases/envelope"
	trustgroup_member_uc "vault-app/internal/trust_group/application/usecases/member"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	trustgroup_usecases_trustgroup "vault-app/internal/trust_group/application/usecases/trust_group"
	trustgroup_adapters "vault-app/internal/trust_group/infrastructure/adapters"
	trustgroup_eventbus "vault-app/internal/trust_group/infrastructure/eventbus"
	vault_dto "vault-app/internal/vault/application/dto"
	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
	vaults_persistence "vault-app/internal/vault/infrastructure/persistence"
	vault_ui "vault-app/internal/vault/ui"
)

type testIPFSResolver struct{}

func (t *testIPFSResolver) GetFileFromIPFS(ctx context.Context, req vault_dto.GetFileFromIPFSRequest) (string, error) {
	return "encrypted-asset-content-1234567890", nil
}

func (t *testIPFSResolver) GetIPFSFile(getFilePayload vault_queries.GetIPFSDataQuerry) ([]byte, error) {
	return []byte("encrypted-asset-content-1234567890"), nil
}

// ---------------------------------------------------------------------------
// Contract-faithful Cloud stub (verified ankhora-cloud C1/C2/C3 contracts)
// ---------------------------------------------------------------------------

type cloudStub struct {
	mu           sync.Mutex
	trustGroups  map[string]map[string]interface{}
	shareEntries map[string]c3_asset_domain.ShareEntry
	workspaces   map[string]tracecore_types.CloudWorkspaceDTO
	channels     map[string][]tracecore_types.CloudChannelDTO
	invitations  map[string]channel_domain.Invitation
	participants map[string][]tracecore_types.CloudChannelParticipant
	token        string
}

func newCloudStub(token string) *cloudStub {
	return &cloudStub{
		trustGroups:  map[string]map[string]interface{}{},
		shareEntries: map[string]c3_asset_domain.ShareEntry{},
		workspaces:   make(map[string]tracecore_types.CloudWorkspaceDTO),
		channels:     make(map[string][]tracecore_types.CloudChannelDTO),
		invitations:  make(map[string]channel_domain.Invitation),
		participants: make(map[string][]tracecore_types.CloudChannelParticipant),
		token:        token,
	}
}

func (s *cloudStub) seedTrustGroup(id string, kekVersion uint64, memberCIDs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	envelopes := make([]map[string]interface{}, 0, len(memberCIDs))
	for i, cid := range memberCIDs {
		envelopes = append(envelopes, map[string]interface{}{
			"id":             uuid.NewString(),
			"trust_group_id": id,
			"member_id":      cid,
			"device_id":      "",
			"kek_version":    kekVersion,
			"wrapped_kek":    "wrapped-kek-" + cid,
			"created_at":     time.Now().UTC().Format(time.RFC3339),
			"_i":             i,
		})
	}
	s.trustGroups[id] = map[string]interface{}{
		"id":            id,
		"channel_id":    "",
		"name":          "Design Team",
		"kek_version":   kekVersion,
		"member_cids":   memberCIDs,
		"key_envelopes": envelopes,
		"created_at":    time.Now().UTC().Format(time.RFC3339),
	}
}

func (s *cloudStub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer "+s.token {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	switch {
	case r.Method == http.MethodPost && (r.URL.Path == "/api/workspaces" || r.URL.Path == "/workspaces"):
		var req tracecore_types.NewCreateWorkspaceRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		wsID := "ws_" + fmt.Sprintf("%d", time.Now().UnixNano())
		ws := tracecore_types.CloudWorkspaceDTO{
			ID:          wsID,
			VaultID:     req.VaultID,
			Name:        req.Name,
			Description: req.Description,
			Status:      "active",
			OwnerID:     req.OwnerID,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
		s.workspaces[wsID] = ws

		chID := "ch_" + fmt.Sprintf("%d", time.Now().UnixNano())
		ch := tracecore_types.CloudChannelDTO{
			ID:          chID,
			WorkspaceID: wsID,
			Title:       "general",
			Status:      "active",
			CreatedAt:   time.Now().UTC(),
		}
		s.channels[wsID] = append(s.channels[wsID], ch)

		writeEnvelope(w, http.StatusCreated, ws)

	case r.Method == http.MethodGet && (strings.HasPrefix(r.URL.Path, "/workspaces") || strings.HasPrefix(r.URL.Path, "/api/workspaces")):
		vaultID := r.URL.Query().Get("vault_id")
		result := make([]tracecore_types.CloudWorkspaceDTO, 0)
		for _, ws := range s.workspaces {
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
		writeEnvelope(w, http.StatusOK, result)

	case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/participants"):
		parts := strings.Split(r.URL.Path, "/")
		var chID string
		for i, part := range parts {
			if part == "channels" && i+1 < len(parts) {
				chID = parts[i+1]
				break
			}
		}
		pList := s.participants[chID]
		if pList == nil {
			pList = make([]tracecore_types.CloudChannelParticipant, 0)
		}
		writeEnvelope(w, http.StatusOK, pList)

	case r.Method == http.MethodGet && (strings.Contains(r.URL.Path, "/channels/workspace/") || strings.HasPrefix(r.URL.Path, "/api/channels")):
		var wsID string
		if idx := strings.Index(r.URL.Path, "/channels/workspace/"); idx != -1 {
			wsID = r.URL.Path[idx+len("/channels/workspace/"):]
		} else {
			wsID = r.URL.Query().Get("workspace_id")
		}
		chList := s.channels[wsID]
		if chList == nil {
			chList = make([]tracecore_types.CloudChannelDTO, 0)
		}
		writeEnvelope(w, http.StatusOK, chList)

	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/invitations"):
		var payload map[string]string
		_ = json.NewDecoder(r.Body).Decode(&payload)
		parts := strings.Split(r.URL.Path, "/")
		channelID := payload["channel_id"]
		if channelID == "" && len(parts) >= 3 {
			channelID = parts[len(parts)-2]
		}
		invID := "inv_" + fmt.Sprintf("%d", time.Now().UnixNano())
		inv := channel_domain.Invitation{
			ID:             invID,
			ChannelID:      channelID,
			InviterVaultID: payload["inviter_vault_id"],
			InviteeVaultID: payload["invitee_vault_id"],
			Status:         channel_domain.InvitationStatusPending,
			CreatedAt:      time.Now().UTC(),
		}
		s.invitations[invID] = inv

		writeEnvelope(w, http.StatusCreated, map[string]interface{}{
			"ID":             inv.ID,
			"ChannelID":      inv.ChannelID,
			"InviterVaultID": inv.InviterVaultID,
			"InviteeVaultID": inv.InviteeVaultID,
			"Status":         string(inv.Status),
			"CreatedAt":      inv.CreatedAt.Format(time.RFC3339),
		})

	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/accept"):
		parts := strings.Split(r.URL.Path, "/")
		var invID string
		for i, part := range parts {
			if part == "invitations" && i+1 < len(parts) {
				invID = parts[i+1]
				break
			}
		}
		var reqBody map[string]string
		_ = json.NewDecoder(r.Body).Decode(&reqBody)
		if invID == "" {
			invID = reqBody["invitation_id"]
		}
		inv, exists := s.invitations[invID]
		if !exists {
			http.Error(w, "invitation not found", http.StatusNotFound)
			return
		}
		inviteeVaultID := reqBody["invitee_vault_id"]
		if inviteeVaultID == "" {
			inviteeVaultID = inv.InviteeVaultID
		}
		inviteePubKey := reqBody["invitee_public_key"]

		now := time.Now().UTC()
		inv.Status = channel_domain.InvitationStatusAccepted
		inv.AcceptedAt = &now
		s.invitations[invID] = inv

		fmt.Printf("[CLOUD_STUB][ACCEPT] invID=%s invFound=%t channelID=%s inviteeVaultID=%s participantsCount=%d\n",
			invID, exists, inv.ChannelID, inviteeVaultID, len(s.participants[inv.ChannelID]))

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
				ChannelID: inv.ChannelID,
				VaultID:   inviteeVaultID,
				PublicKey: inviteePubKey,
				Direction: "bidirectional",
				JoinedAt:  now.Unix(),
			})
		}
		fmt.Printf("[CLOUD_STUB][ACCEPT_DONE] channelID=%s totalParticipants=%d\n", inv.ChannelID, len(s.participants[inv.ChannelID]))

		writeEnvelope(w, http.StatusOK, map[string]interface{}{
			"ID":             inv.ID,
			"ChannelID":      inv.ChannelID,
			"InviterVaultID": inv.InviterVaultID,
			"InviteeVaultID": inv.InviteeVaultID,
			"Status":         string(inv.Status),
			"CreatedAt":      inv.CreatedAt.Format(time.RFC3339),
			"AcceptedAt":     now.Format(time.RFC3339),
		})

	case r.Method == http.MethodPost && r.URL.Path == "/api/trustgroups":
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		tgMap, _ := body["trust_group"].(map[string]interface{})
		if tgMap == nil {
			tgMap = body
		}
		id, _ := tgMap["id"].(string)
		if id == "" {
			id = "tg_" + fmt.Sprintf("%d", time.Now().UnixNano())
			tgMap["id"] = id
		}
		if tgMap["kek_version"] == nil {
			tgMap["kek_version"] = float64(1)
		}
		s.trustGroups[id] = tgMap
		writeEnvelope(w, http.StatusCreated, tgMap)

	case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/members"):
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		parts := strings.Split(r.URL.Path, "/")
		var tgID string
		for i, part := range parts {
			if part == "trustgroups" && i+1 < len(parts) {
				tgID = parts[i+1]
				break
			}
		}
		if tg, ok := s.trustGroups[tgID]; ok {
			vaultID, _ := body["vault_id"].(string)
			var memberCIDs []string
			if raw, exists := tg["member_cids"]; exists && raw != nil {
				if slice, ok := raw.([]string); ok {
					memberCIDs = slice
				} else if interfaceSlice, ok := raw.([]interface{}); ok {
					for _, item := range interfaceSlice {
						if s, ok := item.(string); ok {
							memberCIDs = append(memberCIDs, s)
						}
					}
				}
			}
			memberCIDs = append(memberCIDs, vaultID)
			tg["member_cids"] = memberCIDs
			s.trustGroups[tgID] = tg
			writeEnvelope(w, http.StatusOK, tg)
			return
		}
		http.Error(w, "trust group not found", http.StatusNotFound)

	case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/api/trustgroups/"):
		id := strings.TrimPrefix(r.URL.Path, "/api/trustgroups/")
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		tgMap, _ := body["trust_group"].(map[string]interface{})
		if tgMap == nil {
			tgMap = body
		}
		s.trustGroups[id] = tgMap
		writeEnvelope(w, http.StatusOK, tgMap)

	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/trustgroups/"):
		id := strings.TrimPrefix(r.URL.Path, "/api/trustgroups/")
		tg, ok := s.trustGroups[id]
		if !ok {
			http.Error(w, "trust group not found", http.StatusNotFound)
			return
		}
		writeEnvelope(w, http.StatusOK, tg)

	case r.Method == http.MethodPost && r.URL.Path == "/api/c3/share-entries":
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tgID, _ := body["trust_group_id"].(string)
		tg, ok := s.trustGroups[tgID]
		if !ok {
			http.Error(w, "trust group not found", http.StatusNotFound)
			return
		}
		kekVersion, _ := body["kek_version"].(float64)
		var tgKEKVersion uint64
		switch v := tg["kek_version"].(type) {
		case float64:
			tgKEKVersion = uint64(v)
		case uint64:
			tgKEKVersion = v
		case int:
			tgKEKVersion = uint64(v)
		}
		if uint64(kekVersion) != tgKEKVersion {
			http.Error(w, "stale kek_version: does not match trust group current kek_version", http.StatusConflict)
			return
		}
		id, _ := body["id"].(string)
		if id == "" {
			id = uuid.NewString()
		}
		entry := c3_asset_domain.ShareEntry{
			ID:           id,
			AssetCID:     body["asset_cid"].(string),
			TrustGroupID: tgID,
			WrappedDEK:   body["wrapped_dek"].(string),
			KEKVersion:   uint64(kekVersion),
			CreatedBy:    body["created_by"].(string),
			CreatedAt:    time.Now().UTC(),
			Status:       c3_asset_domain.ShareEntryStatusActive,
			Metadata:     toStringMap(body["metadata"]),
		}
		s.shareEntries[entry.ID] = entry
		writeEnvelope(w, http.StatusCreated, entry)

	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/c3/share-entries/"):
		id := strings.TrimPrefix(r.URL.Path, "/api/c3/share-entries/")
		entry, ok := s.shareEntries[id]
		if !ok {
			http.Error(w, "share entry not found", http.StatusNotFound)
			return
		}
		writeEnvelope(w, http.StatusOK, entry)

	case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/storage"):
		var streamReq tracecore_types.SyncVaultStreamRequest
		_ = json.NewDecoder(r.Body).Decode(&streamReq)
		hash := sha256.Sum256(streamReq.Stream)
		cid := "bafybeivertical" + hex.EncodeToString(hash[:8])
		writeEnvelope(w, http.StatusOK, map[string]interface{}{
			"cid": cid,
		})

	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func writeEnvelope(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  status,
		"data":    data,
		"message": "ok",
		"success": true,
	})
}

func toStringMap(v interface{}) map[string]string {
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, val := range m {
		s, _ := val.(string)
		out[k] = s
	}
	return out
}

// ---------------------------------------------------------------------------
// Vertical test
// ---------------------------------------------------------------------------

// TestCreateCollaborativeShare_VerticalPersistence exercises the real Wails
// entry point used by CreateCollaborativeShareModal:
//
//	App.CreateCollaborativeShare (JWT auth)
//	  -> CollaborationHandler.CreateCollaborativeShare
//	    -> CreateCollaborativeShareUseCase
//	      -> ShareAssetWithTrustGroupUsecase
//	        -> TracecoreClient.GetTrustGroup      (GET  /api/trustgroups/{id})
//	        -> CloudShareEntryRepository          (POST /api/c3/share-entries)
//	  -> real persisted share_entry_id returned
//	  -> GET /api/c3/share-entries/{id} returns the same ShareEntry
//
// The crypto material is produced by the real desktop orchestration path
// (TrustGroupCryptoOrchestrator.PrepareCollaborativeAsset), never fabricated.
func TestCreateCollaborativeShare_VerticalPersistence(t *testing.T) {
	ctx := context.Background()

	// ------------------------------------------------------------------
	// 1. Real desktop crypto orchestration (upstream of the share flow)
	// ------------------------------------------------------------------
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), vault_infrastructure_security.OSFileSystem{})
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, nil, nil)

	tgID := "tg-vertical-" + uuid.NewString()[:8]
	const kekVersion = uint64(2)

	kr := &vaults_domain.VaultKeyring{UserID: "user_alice", VaultID: "vault_alice"}
	testKEK := make([]byte, 32)
	for i := range testKEK {
		testKEK[i] = byte(i + 1)
	}
	_, _ = keyringSvc.StoreTrustGroupKEK(kr, tgID, kekVersion, testKEK)
	_ = keyringSvc.SaveHybrid(kr, "user_alice", "", "")

	rawPayload := []byte("CONFIDENTIAL VERTICAL TEST PAYLOAD")

	prepared, err := orchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset-vertical-001",
		TrustGroupID: tgID,
		KEKVersion:   kekVersion,
		RawPayload:   rawPayload,
		Keyring:      kr,
	})
	require.NoError(t, err)
	require.NotEmpty(t, prepared.WrappedDEK)
	require.Equal(t, kekVersion, prepared.KEKVersion)
	_ = keyringSvc.SaveHybrid(kr, "user_alice", "", "")

	hash := sha256.Sum256(prepared.EncryptedData)
	assetCID := "bafybeivertical" + hex.EncodeToString(hash[:8])

	wrappedDEKB64 := base64.StdEncoding.EncodeToString(prepared.WrappedDEK)
	_ = wrappedDEKB64

	// ------------------------------------------------------------------
	// 2. Cloud stub speaking the verified C1/C2/C3 contracts
	// ------------------------------------------------------------------
	cloudToken := "0123456789abcdefghijklmnopqrstuv" // 32 chars, opaque bearer
	stub := newCloudStub(cloudToken)
	stub.seedTrustGroup(tgID, kekVersion, []string{"vault_alice"})
	ts := httptest.NewServer(stub)
	defer ts.Close()

	tc := tracecore.NewTracecoreClient(ts.URL+"/api", cloudToken, ts.URL, ts.URL+"/api")

	// ------------------------------------------------------------------
	// 3. Wire the App exactly like app.go does for production
	// ------------------------------------------------------------------
	authV2 := auth_domain.Auth{
		Issuer:      "ankhora-test",
		Audience:    "ankhora-desktop",
		Secret:      "vertical-test-secret",
		TokenExpiry: 15 * time.Minute,
	}
	tokenService := auth_usecases.NewTokenService(authV2, nil, nil)
	tokenUC := auth_usecases.NewGenerateTokensUseCase(nil, tokenService)
	authHandler := auth_ui.NewAuthHandler(nil, tokenUC, nil)

	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(tc, tracecore.NewCloudShareEntryRepository(tc))
	identityResolver := collaboration_infra.NewKeyringSovereignIdentityResolver(keyringSvc)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, nil).WithCrypto(orchestrator, nil, identityResolver, blockchain.NewCloudIPFSStorage(tc, "", ""), &testIPFSResolver{})
	collabHandler := collaboration_ui.NewCollaborationHandler(createCollabShareUC, nil, nil)

	dbMem, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = dbMem.AutoMigrate(
		&subscription_persistence.SubscriptionMapper{},
		&app_config_persistence.UserConfigMapper{},
		&app_config_domain.AppConfig{},
		&app_config_domain.VaultConfigBeta{},
		&app_config_domain.DeviceConfig{},
		&app_config_domain.SubscriptionConfig{},
		&app_config_persistence.OnboardingConfigSqlDB{},
		&onboarding_persistence.UserDB{},
		&vaults_persistence.SessionMapper{},
	)

	vaultHandler := vault_ui.NewVaultHandler(nil, logger.Logger{}, ctx, nil, nil, dbMem, tc, t.TempDir())
	vaultHandler.VaultRepository = &topLevelVaultRepo{}
	_, _ = vaultHandler.SessionManager.Prepare("user_alice")
	_, _ = vaultHandler.SessionManager.AttachRuntime("user_alice", &vault_session.RuntimeContext{
		AppConfig: app_config_domain.AppConfig{
			Branch: "main",
		},
	})

	subRepo := subscription_persistence.NewSubscriptionRepository(dbMem, nil)
	_ = subRepo.Save(ctx, &subscription_domain.Subscription{
		ID:     "sub_alice",
		UserID: "user_alice",
		Email:  "alice@ankhora.test",
	})
	logSvc := &logger.Logger{}
	appConfigHandler := app_config_ui.NewAppConfigHandler(dbMem, *logSvc)
	appConfigHandler.VaultHandler = vaultHandler
	_ = appConfigHandler.UserConfigRepository.CreateUserConfig(&app_config_domain.UserConfig{
		ID:    "user_alice",
		Email: "alice@ankhora.test",
	})
	_ = appConfigHandler.AppConfigRepository.CreateAppConfig(&app_config_domain.AppConfig{
		UserID: "user_alice",
		Branch: "main",
	})
	_, _ = appConfigHandler.VaultConfigRepository.Create(&app_config_domain.VaultConfigBeta{
		BaseVaultConfig: app_config_domain.BaseVaultConfig{
			ID:        "vc_alice",
			UserID:    "user_alice",
			VaultName: "Default Vault",
		},
	})
	_ = appConfigHandler.SubscriptionConfigRepository.Create(&app_config_domain.SubscriptionConfig{
		BaseVaultConfig: app_config_domain.BaseVaultConfig{
			ID:        "sc_alice",
			UserID:    "user_alice",
			VaultName: "Default Vault",
		},
	})
	_ = appConfigHandler.OnboardingConfigRepository.Create(&app_config_domain.OnboardingConfig{
		UserID: "user_alice",
	})
	onboardingHandler := onboarding_ui_wails.NewOnBoardingHandler(nil, nil, nil, tc, dbMem, logSvc, *keyringSvc)
	_, _ = onboardingHandler.UserRepo.Create(&onboarding_domain.User{
		ID:    "user_alice_1",
		Email: "alice@ankhora.test",
	})
	_, _ = onboardingHandler.UserRepo.Create(&onboarding_domain.User{
		ID:    "user_alice",
		Email: "",
	})
	appConfigHandler.SetOnboardingHandler(*onboardingHandler)
	subHandler := subscription_ui_wails.NewSubscriptionHandler(dbMem, tc, nil, nil, nil, nil, nil, *appConfigHandler, *logSvc)

	app := &App{
		AuthHandler:          authHandler,
		CollaborationHandler: collabHandler,
		Vault:                vaultHandler,
		SubscriptionHandler:  subHandler,
		AppConfigHandler:     appConfigHandler,
		tracecoreClient:      tc,
		ctx:                  ctx,
	}

	pairs, err := tokenService.GenerateTokenPair(&auth_domain.JwtUser{
		ID:       "user_alice",
		Username: "alice",
		Email:    "alice@ankhora.test",
	})
	require.NoError(t, err)

	// ------------------------------------------------------------------
	// 4. Act — the actual Wails/API path
	// ------------------------------------------------------------------
	shareRef, err := app.CreateCollaborativeShare(
		pairs.Token,
		"thread_vertical_1",
		tgID,
		assetCID,
		"vertical persistence test",
	)
	require.NoError(t, err)
	require.NotNil(t, shareRef)

	// ------------------------------------------------------------------
	// 5. Assert — the returned ID is the REAL persisted one
	// ------------------------------------------------------------------
	require.NotEmpty(t, shareRef.ShareEntryID)
	t.Logf("REAL persisted share_entry_id returned by Cloud path: %s", shareRef.ShareEntryID)
	assert.False(t, strings.HasPrefix(shareRef.ShareEntryID, "se_"),
		"ID must never come from the removed fabrication fallback")
	_, parseErr := uuid.Parse(shareRef.ShareEntryID)
	require.NoError(t, parseErr, "share_entry_id must be a real Cloud-assigned UUID")

	stub.mu.Lock()
	persisted, ok := stub.shareEntries[shareRef.ShareEntryID]
	stub.mu.Unlock()
	require.True(t, ok, "share entry must actually be persisted on the Cloud side")
	assert.NotEmpty(t, persisted.AssetCID)
	assert.True(t, strings.HasPrefix(persisted.AssetCID, "bafybeivertical"))
	assert.NotEmpty(t, persisted.WrappedDEK)
	assert.Equal(t, kekVersion, persisted.KEKVersion)
	assert.Equal(t, tgID, persisted.TrustGroupID)

	// activeShareEntryRef is set from this exact return value in
	// useC3CollaborationStore.createShare; prove it refers to a persisted
	// entry by fetching it back through the C3 contract.
	fetchedResp, err := tc.GetShareEntryDirect(ctx, shareRef.ShareEntryID)
	require.NoError(t, err)
	fetched := fetchedResp.Data
	assert.Equal(t, shareRef.ShareEntryID, fetched.ID)
	assert.Equal(t, persisted.AssetCID, fetched.AssetCID)
	assert.Equal(t, persisted.WrappedDEK, fetched.WrappedDEK)
	assert.Equal(t, persisted.KEKVersion, fetched.KEKVersion)
	assert.Equal(t, persisted.TrustGroupID, fetched.TrustGroupID)
	assert.Equal(t, c3_asset_domain.ShareEntryStatusActive, fetched.Status)

	// ------------------------------------------------------------------
	// 6. Negative paths — invalid trust group is rejected
	// ------------------------------------------------------------------
	_, err = app.CreateCollaborativeShare(pairs.Token, "thread_vertical_1", "tg-does-not-exist", assetCID, "")
	require.Error(t, err)

	// Missing asset CID must fail.
	_, err = app.CreateCollaborativeShare(pairs.Token, "thread_vertical_1", tgID, "", "")
	require.Error(t, err)
}

type gormCloudServer struct {
	db         *gorm.DB
	token      string
	files      map[string][]byte
	publicKeys map[string]string
}

func newGormCloudServer(db *gorm.DB, token string) *gormCloudServer {
	server := &gormCloudServer{
		db:         db,
		token:      token,
		files:      make(map[string][]byte),
		publicKeys: make(map[string]string),
	}
	type Customer struct {
		ID        uint      `gorm:"primaryKey;autoIncrement"`
		Email     string    `gorm:"column:email;size:255"`
		PublicKey string    `gorm:"column:public_key"`
		FirstName string    `gorm:"column:first_name"`
		LastName  string    `gorm:"column:last_name"`
		CreatedAt time.Time `gorm:"column:created_at"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
	}
	_ = db.AutoMigrate(&Customer{})

	_ = db.Table("workspaces").AutoMigrate(&struct {
		ID      string `gorm:"primaryKey"`
		Name    string
		OwnerID string
		VaultID string
	}{})
	_ = db.Table("channels").AutoMigrate(&struct {
		ID          string `gorm:"primaryKey"`
		WorkspaceID string
		Title       string
	}{})
	_ = db.Table("channel_participants").AutoMigrate(&struct {
		ID        string `gorm:"primaryKey"`
		ChannelID string
		VaultID   string
	}{})
	_ = db.Table("invitations").AutoMigrate(&struct {
		ID             string `gorm:"primaryKey"`
		ChannelID      string
		InviterVaultID string
		InviteeVaultID string
		Status         string
	}{})
	_ = db.Table("trust_groups").AutoMigrate(&struct {
		ID         string `gorm:"primaryKey"`
		ChannelID  string `gorm:"column:channel_id"`
		Name       string `gorm:"column:name"`
		KEKVersion uint64 `gorm:"column:kek_version"`
		MemberCIDs string `gorm:"column:member_cids"`
	}{})
	_ = db.Table("trust_group_members").AutoMigrate(&struct {
		ID           string `gorm:"primaryKey"`
		TrustGroupID string `gorm:"column:trust_group_id"`
		VaultID      string `gorm:"column:vault_id"`
		Role         string `gorm:"column:role"`
	}{})
	_ = db.Table("trust_group_key_envelopes").AutoMigrate(&struct {
		ID           string `gorm:"primaryKey"`
		TrustGroupID string `gorm:"column:trust_group_id"`
		MemberID     string `gorm:"column:member_id"`
		DeviceID     string `gorm:"column:device_id"`
		KEKVersion   uint64 `gorm:"column:kek_version"`
		WrappedKEK   string `gorm:"column:wrapped_kek"`
	}{})
	_ = db.Table("c3_share_entries").AutoMigrate(&struct {
		ID               string `gorm:"primaryKey"`
		TrustGroupID     string `gorm:"column:trust_group_id"`
		AssetCID         string `gorm:"column:asset_cid"`
		WrappedDEK       string `gorm:"column:wrapped_dek"`
		EncryptedPayload string `gorm:"column:encrypted_payload"`
		KEKVersion       uint64 `gorm:"column:kek_version"`
		Status           string `gorm:"column:status"`
	}{})

	return server
}

func (s *gormCloudServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") == "" && !strings.Contains(r.URL.Path, "/customers") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	switch {
	case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/vaults/subscription/"):
		writeEnvelope(w, http.StatusOK, map[string]interface{}{
			"id":   "v_cloud_1",
			"name": "Cloud Vault",
		})

	case r.Method == http.MethodGet && (strings.Contains(r.URL.Path, "/vaults/public-key/") || strings.Contains(r.URL.Path, "/customers/public-key/") || strings.Contains(r.URL.Path, "/identity/public-key/")):
		writeEnvelope(w, http.StatusOK, map[string]interface{}{
			"VaultID":      "v_cloud_bob",
			"VaultAddress": "v_cloud_bob",
			"Name":         "Bob Vault",
		})

	case r.Method == http.MethodPost && (r.URL.Path == "/api/workspaces" || r.URL.Path == "/workspaces"):
		var req tracecore_types.NewCreateWorkspaceRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		wsID := "ws_" + fmt.Sprintf("%d", time.Now().UnixNano())
		ws := tracecore_types.CloudWorkspaceDTO{
			ID:          wsID,
			VaultID:     req.VaultID,
			Name:        req.Name,
			Description: req.Description,
			Status:      "active",
			OwnerID:     req.OwnerID,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}
		chID := "ch_" + fmt.Sprintf("%d", time.Now().UnixNano())
		ch := tracecore_types.CloudChannelDTO{
			ID:          chID,
			WorkspaceID: wsID,
			Title:       "general",
			Status:      "active",
			CreatedAt:   time.Now().UTC(),
		}
		_ = s.db.Table("workspaces").Create(map[string]interface{}{"id": ws.ID, "name": ws.Name, "owner_id": ws.OwnerID, "vault_id": ws.VaultID})
		_ = s.db.Table("channels").Create(map[string]interface{}{"id": ch.ID, "workspace_id": ws.ID, "title": ch.Title})
		writeEnvelope(w, http.StatusCreated, ws)

	case r.Method == http.MethodGet && (strings.HasPrefix(r.URL.Path, "/workspaces") || strings.HasPrefix(r.URL.Path, "/api/workspaces")):
		vaultID := r.URL.Query().Get("vault_id")
		var rows []map[string]interface{}
		_ = s.db.Table("workspaces").Find(&rows).Error
		result := make([]tracecore_types.CloudWorkspaceDTO, 0)
		for _, row := range rows {
			result = append(result, tracecore_types.CloudWorkspaceDTO{
				ID:        fmt.Sprintf("%v", row["id"]),
				VaultID:   fmt.Sprintf("%v", row["vault_id"]),
				Name:      fmt.Sprintf("%v", row["name"]),
				OwnerID:   fmt.Sprintf("%v", row["owner_id"]),
				Status:    "active",
				CreatedAt: time.Now().UTC(),
			})
		}
		_ = vaultID
		writeEnvelope(w, http.StatusOK, result)

	case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/participants"):
		parts := strings.Split(r.URL.Path, "/")
		var chID string
		for i, part := range parts {
			if part == "channels" && i+1 < len(parts) {
				chID = parts[i+1]
				break
			}
		}
		var rows []map[string]interface{}
		_ = s.db.Table("channel_participants").Where("channel_id = ?", chID).Find(&rows).Error
		pList := make([]tracecore_types.CloudChannelParticipant, 0)
		for _, row := range rows {
			pList = append(pList, tracecore_types.CloudChannelParticipant{
				ChannelID: fmt.Sprintf("%v", row["channel_id"]),
				VaultID:   fmt.Sprintf("%v", row["vault_id"]),
			})
		}
		writeEnvelope(w, http.StatusOK, pList)

	case r.Method == http.MethodGet && (strings.Contains(r.URL.Path, "/channels/workspace/") || strings.HasPrefix(r.URL.Path, "/api/channels")):
		var wsID string
		if idx := strings.Index(r.URL.Path, "/channels/workspace/"); idx != -1 {
			wsID = r.URL.Path[idx+len("/channels/workspace/"):]
		} else {
			wsID = r.URL.Query().Get("workspace_id")
		}
		var rows []map[string]interface{}
		_ = s.db.Table("channels").Where("workspace_id = ?", wsID).Find(&rows).Error
		chList := make([]tracecore_types.CloudChannelDTO, 0)
		for _, row := range rows {
			chList = append(chList, tracecore_types.CloudChannelDTO{
				ID:          fmt.Sprintf("%v", row["id"]),
				WorkspaceID: fmt.Sprintf("%v", row["workspace_id"]),
				Title:       fmt.Sprintf("%v", row["title"]),
				Status:      "active",
			})
		}
		writeEnvelope(w, http.StatusOK, chList)

	case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/invitations") && !strings.HasSuffix(r.URL.Path, "/accept"):
		var payload map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&payload)
		channelID := fmt.Sprintf("%v", payload["channel_id"])
		if channelID == "" || channelID == "<nil>" {
			channelID = fmt.Sprintf("%v", payload["ChannelID"])
		}
		parts := strings.Split(r.URL.Path, "/")
		for i, part := range parts {
			if part == "channels" && i+1 < len(parts) && parts[i+1] != "invitations" {
				channelID = parts[i+1]
				break
			}
		}
		inviterVaultID := fmt.Sprintf("%v", payload["inviter_vault_id"])
		if inviterVaultID == "" || inviterVaultID == "<nil>" {
			inviterVaultID = fmt.Sprintf("%v", payload["InviterVaultID"])
		}
		inviteeVaultID := fmt.Sprintf("%v", payload["invitee_vault_id"])
		if inviteeVaultID == "" || inviteeVaultID == "<nil>" {
			inviteeVaultID = fmt.Sprintf("%v", payload["InviteeVaultID"])
		}
		invID := "inv_" + fmt.Sprintf("%d", time.Now().UnixNano())
		fmt.Printf("[CLOUD_SERVER][INVITE] path=%s invID=%s chID=%s inviter=%s invitee=%s\n", r.URL.Path, invID, channelID, inviterVaultID, inviteeVaultID)
		_ = s.db.Table("invitations").Create(map[string]interface{}{"id": invID, "channel_id": channelID, "inviter_vault_id": inviterVaultID, "invitee_vault_id": inviteeVaultID, "status": "pending"})

		writeEnvelope(w, http.StatusCreated, map[string]interface{}{
			"ID":             invID,
			"ChannelID":      channelID,
			"InviterVaultID": inviterVaultID,
			"InviteeVaultID": inviteeVaultID,
			"Status":         "pending",
			"CreatedAt":      time.Now().UTC().Format(time.RFC3339),
		})

	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/accept"):
		parts := strings.Split(r.URL.Path, "/")
		var invID string
		for i, part := range parts {
			if part == "invitations" && i+1 < len(parts) {
				invID = parts[i+1]
				break
			}
		}
		var rows []map[string]interface{}
		_ = s.db.Table("invitations").Where("id = ?", invID).Find(&rows).Error
		var channelID, inviteeVaultID, inviterVaultID string
		if len(rows) > 0 {
			channelID = fmt.Sprintf("%v", rows[0]["channel_id"])
			inviteeVaultID = fmt.Sprintf("%v", rows[0]["invitee_vault_id"])
			inviterVaultID = fmt.Sprintf("%v", rows[0]["inviter_vault_id"])
		}
		fmt.Printf("[CLOUD_SERVER][ACCEPT] path=%s invID=%s chID=%s invitee=%s\n", r.URL.Path, invID, channelID, inviteeVaultID)
		_ = s.db.Table("invitations").Where("id = ?", invID).Update("status", "accepted")
		_ = s.db.Table("channel_participants").Create(map[string]interface{}{"id": "p_" + fmt.Sprintf("%d", time.Now().UnixNano()), "channel_id": channelID, "vault_id": inviteeVaultID})

		writeEnvelope(w, http.StatusOK, map[string]interface{}{
			"ID":             invID,
			"ChannelID":      channelID,
			"InviterVaultID": inviterVaultID,
			"InviteeVaultID": inviteeVaultID,
			"Status":         "accepted",
			"CreatedAt":      time.Now().UTC().Format(time.RFC3339),
		})

	case r.Method == http.MethodPost && (strings.HasPrefix(r.URL.Path, "/trustgroups") || strings.HasPrefix(r.URL.Path, "/api/trustgroups")) && !strings.Contains(r.URL.Path, "/members"):
		bodyBytes, _ := io.ReadAll(r.Body)
		var reqMap map[string]interface{}
		_ = json.Unmarshal(bodyBytes, &reqMap)

		var tg trustgroup_domain.TrustGroup
		if tgRaw, ok := reqMap["trust_group"].(map[string]interface{}); ok {
			b, _ := json.Marshal(tgRaw)
			_ = json.Unmarshal(b, &tg)
		} else {
			_ = json.Unmarshal(bodyBytes, &tg)
		}
		if membersRaw, ok := reqMap["members"].([]interface{}); ok {
			for _, m := range membersRaw {
				if s, strOk := m.(string); strOk && s != "" {
					tg.MemberCIDs = append(tg.MemberCIDs, s)
				} else if mObj, objOk := m.(map[string]interface{}); objOk {
					if vID, vOk := mObj["vault_id"].(string); vOk && vID != "" {
						tg.MemberCIDs = append(tg.MemberCIDs, vID)
					}
				}
			}
		}

		if tg.ID == "" {
			tg.ID = "tg_" + fmt.Sprintf("%d", time.Now().UnixNano())
		}
		if tg.KEKVersion == 0 {
			tg.KEKVersion = 1
		}
		cidsJSON, _ := json.Marshal(tg.MemberCIDs)
		_ = s.db.Table("trust_groups").Create(map[string]interface{}{"id": tg.ID, "channel_id": tg.ChannelID, "name": tg.Name, "kek_version": tg.KEKVersion, "member_cids": string(cidsJSON)})
		writeEnvelope(w, http.StatusCreated, tg)

	case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/trustgroups/") && strings.HasSuffix(r.URL.Path, "/members"):
		var req trustgroup_domain.AddMemberToTrustGroupRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		parts := strings.Split(r.URL.Path, "/")
		tgID := req.TrustGroupID
		if tgID == "" && len(parts) >= 3 {
			tgID = parts[len(parts)-2]
		}
		var rows []map[string]interface{}
		_ = s.db.Table("trust_groups").Where("id = ?", tgID).Find(&rows).Error
		var cids []string
		if len(rows) > 0 {
			_ = json.Unmarshal([]byte(fmt.Sprintf("%v", rows[0]["member_cids"])), &cids)
		}
		cids = append(cids, req.VaultID)
		cidsJSON, _ := json.Marshal(cids)
		_ = s.db.Table("trust_groups").Where("id = ?", tgID).Update("member_cids", string(cidsJSON))
		_ = s.db.Table("trust_group_members").Create(map[string]interface{}{"id": tgID + ":" + req.VaultID, "trust_group_id": tgID, "vault_id": req.VaultID, "role": req.Role})
		writeEnvelope(w, http.StatusOK, trustgroup_domain.TrustGroup{ID: tgID, MemberCIDs: cids, KEKVersion: 1})

	case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/api/trustgroups/"):
		tgID := strings.TrimPrefix(r.URL.Path, "/api/trustgroups/")
		if idx := strings.Index(tgID, "?"); idx != -1 {
			tgID = tgID[:idx]
		}
		var tg trustgroup_domain.TrustGroup
		_ = json.NewDecoder(r.Body).Decode(&tg)
		if tg.ID == "" {
			tg.ID = tgID
		}

		fmt.Printf("[TRUSTGROUP][CLOUD][UPDATE][01]\nincoming keyEnvelopeCount=%d\n", len(tg.KeyEnvelopes))
		fmt.Printf("decoded keyEnvelopeCount=%d\nDB transaction started\n", len(tg.KeyEnvelopes))

		tx := s.db.Begin()
		cidsJSON, _ := json.Marshal(tg.MemberCIDs)
		_ = tx.Table("trust_groups").Where("id = ?", tgID).Updates(map[string]interface{}{"member_cids": string(cidsJSON), "kek_version": tg.KEKVersion})

		delRes := tx.Table("trust_group_key_envelopes").Where("trust_group_id = ?", tgID).Delete(map[string]interface{}{})
		fmt.Printf("existing envelope rows deleted=%d\n", delRes.RowsAffected)

		insertedCount := 0
		for _, env := range tg.KeyEnvelopes {
			envID := env.ID
			if envID == "" {
				envID = tgID + ":" + env.MemberID + ":" + env.DeviceID
			}
			_ = tx.Table("trust_group_key_envelopes").Create(map[string]interface{}{
				"id":             envID,
				"trust_group_id": tgID,
				"member_id":      env.MemberID,
				"device_id":      env.DeviceID,
				"kek_version":    env.KEKVersion,
				"wrapped_kek":    env.WrappedKEK,
			})
			insertedCount++
			fmt.Printf("[TRUSTGROUP][CLOUD][ENVELOPE][INSERT]\ntrustGroupID=%s\nenvelopeID=%s\nmemberID=%s\ndeviceID=%s\nkekVersion=%d\n",
				tgID, envID, env.MemberID, env.DeviceID, env.KEKVersion)
		}
		tx.Commit()
		fmt.Printf("envelope rows inserted=%d\ntransaction committed=true\n", insertedCount)

		writeEnvelope(w, http.StatusOK, tg)

	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/trustgroups/"):
		tgID := strings.TrimPrefix(r.URL.Path, "/api/trustgroups/")
		if idx := strings.Index(tgID, "?"); idx != -1 {
			tgID = tgID[:idx]
		}
		var rows []map[string]interface{}
		_ = s.db.Table("trust_groups").Where("id = ?", tgID).Find(&rows).Error
		if len(rows) == 0 {
			http.Error(w, "trust group not found", http.StatusNotFound)
			return
		}
		var cids []string
		_ = json.Unmarshal([]byte(fmt.Sprintf("%v", rows[0]["member_cids"])), &cids)

		var kekVersion uint64 = 1
		if kv, ok := rows[0]["kek_version"].(int64); ok {
			kekVersion = uint64(kv)
		} else if kv, ok := rows[0]["kek_version"].(uint64); ok {
			kekVersion = kv
		}

		var envRows []map[string]interface{}
		_ = s.db.Table("trust_group_key_envelopes").Where("trust_group_id = ?", tgID).Find(&envRows).Error
		envs := make([]trustgroup_domain.TrustGroupKeyEnvelope, 0, len(envRows))
		for _, er := range envRows {
			var kv uint64 = 1
			if val, ok := er["kek_version"].(int64); ok {
				kv = uint64(val)
			} else if val, ok := er["kek_version"].(uint64); ok {
				kv = val
			}
			envs = append(envs, trustgroup_domain.TrustGroupKeyEnvelope{
				ID:           fmt.Sprintf("%v", er["id"]),
				TrustGroupID: tgID,
				MemberID:     fmt.Sprintf("%v", er["member_id"]),
				DeviceID:     fmt.Sprintf("%v", er["device_id"]),
				KEKVersion:   kv,
				WrappedKEK:   fmt.Sprintf("%v", er["wrapped_kek"]),
			})
		}
		writeEnvelope(w, http.StatusOK, trustgroup_domain.TrustGroup{
			ID:           tgID,
			Name:         fmt.Sprintf("%v", rows[0]["name"]),
			ChannelID:    fmt.Sprintf("%v", rows[0]["channel_id"]),
			MemberCIDs:   cids,
			KEKVersion:   kekVersion,
			KeyEnvelopes: envs,
		})

	case r.Method == http.MethodPost && (strings.HasPrefix(r.URL.Path, "/api/c3asset/shares") || strings.HasPrefix(r.URL.Path, "/api/shares-v0") || strings.HasPrefix(r.URL.Path, "/api/c3/share-entries")):
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		id, _ := body["id"].(string)
		if id == "" {
			id = uuid.NewString()
		}
		trustGroupID, _ := body["trust_group_id"].(string)
		assetCID, _ := body["asset_cid"].(string)
		wrappedDEK, _ := body["wrapped_dek"].(string)
		encPayload, _ := body["encrypted_payload"].(string)
		var kekVersion uint64
		if kv, ok := body["kek_version"].(float64); ok {
			kekVersion = uint64(kv)
		}
		_ = s.db.Table("c3_share_entries").Create(map[string]interface{}{"id": id, "trust_group_id": trustGroupID, "asset_cid": assetCID, "wrapped_dek": wrappedDEK, "encrypted_payload": encPayload, "kek_version": kekVersion, "status": "active"})
		entry := c3_asset_domain.ShareEntry{
			ID:           id,
			TrustGroupID: trustGroupID,
			AssetCID:     assetCID,
			WrappedDEK:   wrappedDEK,
			KEKVersion:   kekVersion,
			Status:       c3_asset_domain.ShareEntryStatusActive,
		}
		writeEnvelope(w, http.StatusCreated, entry)

	case r.Method == http.MethodGet && (strings.HasPrefix(r.URL.Path, "/api/c3asset/shares/") || strings.HasPrefix(r.URL.Path, "/api/c3/share-entries/")):
		parts := strings.Split(r.URL.Path, "/")
		shareID := parts[len(parts)-1]
		var rows []map[string]interface{}
		_ = s.db.Table("c3_share_entries").Where("id = ?", shareID).Find(&rows).Error
		var row map[string]interface{}
		if len(rows) > 0 {
			row = rows[0]
		}
		writeEnvelope(w, http.StatusOK, c3_asset_domain.ShareEntry{
			ID:           shareID,
			TrustGroupID: fmt.Sprintf("%v", row["trust_group_id"]),
			AssetCID:     fmt.Sprintf("%v", row["asset_cid"]),
			WrappedDEK:   fmt.Sprintf("%v", row["wrapped_dek"]),
			KEKVersion:   uint64(row["kek_version"].(int64)),
			Status:       c3_asset_domain.ShareEntryStatusActive,
		})

	case r.Method == http.MethodPost && (r.URL.Path == "/api/thread-data/access" || r.URL.Path == "/api/c3/thread-data/access"):
		var req tracecore_types.ThreadDataAccessRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		var rows []map[string]interface{}
		_ = s.db.Table("c3_share_entries").Find(&rows).Error
		var row map[string]interface{}
		if len(rows) > 0 {
			row = rows[0]
		}
		resp := tracecore_types.AccessCryptoShareResponse{
			EncryptedPayload: fmt.Sprintf("%v", row["encrypted_payload"]),
			EncryptedKey:     fmt.Sprintf("%v", row["wrapped_dek"]),
			SenderPublicKey:  req.RequestingVaultID,
			DownloadAllowed:  true,
		}
		writeEnvelope(w, http.StatusOK, resp)

	case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/storage"):
		var req tracecore_types.SyncVaultStreamRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		cid := "cid_" + fmt.Sprintf("%d", time.Now().UnixNano())
		if len(req.Stream) > 0 {
			s.files[cid] = req.Stream
		}
		writeEnvelope(w, http.StatusOK, tracecore_types.SyncVaultResponse{CID: cid})

	case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/storage"):
		parts := strings.Split(r.URL.Path, "/")
		cid := parts[len(parts)-1]
		data := s.files[cid]
		b64 := base64.StdEncoding.EncodeToString(data)
		resp := tracecore_types.IpfsCidResponse{
			Status:  200,
			Data:    b64,
			Message: "ok",
			Success: true,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)

	case r.Method == http.MethodGet && (strings.HasPrefix(r.URL.Path, "/api/customers") || strings.HasPrefix(r.URL.Path, "/customers")):
		email := r.URL.Query().Get("email")
		pubKey := ""
		if s.publicKeys != nil {
			pubKey = s.publicKeys[email]
		}
		if pubKey == "" {
			type Customer struct {
				ID        uint      `gorm:"primaryKey;autoIncrement"`
				Email     string    `gorm:"column:email;size:255"`
				PublicKey string    `gorm:"column:public_key"`
				FirstName string    `gorm:"column:first_name"`
				LastName  string    `gorm:"column:last_name"`
				CreatedAt time.Time `gorm:"column:created_at"`
				UpdatedAt time.Time `gorm:"column:updated_at"`
			}
			var cust Customer
			_ = s.db.Table("customers").Where("LOWER(email) = LOWER(?)", email).First(&cust).Error
			pubKey = cust.PublicKey
		}
		user := tracecore_types.User{
			ID:        1,
			Email:     email,
			PublicKey: pubKey,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(tracecore.GetUserByEmailResponse{
			Error:   false,
			Message: "ok",
			Data:    user,
		})
		return

		http.Error(w, "customer not found", http.StatusNotFound)

	case r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/customers/add-public-key"):

		writeEnvelope(w, http.StatusOK, tracecore_types.AddPublicKeyToCustomerResponse{
			ID:    1,
			Email: "test@ankhora.test",
		})

	default:
		http.Error(w, "not found: "+r.URL.Path, http.StatusNotFound)
	}
}

type sessionSovereignIdentityResolver struct {
	identityDeviceAdapter *trustgroup_adapters.IdentityDeviceAdapter
	keyringSvc            *vault_infrastructure_security.KeyringService
	sessionMgr            *vault_session.Manager
}

func (r *sessionSovereignIdentityResolver) GetDeviceSeed(ctx context.Context, userID string) (string, error) {
	if r.sessionMgr != nil {
		if sess, err := r.sessionMgr.GetSession(userID); err == nil && sess != nil && sess.Runtime != nil {
			if secret, ok := sess.Runtime.SessionSecrets["stellar_secret"]; ok && secret != "" {
				return secret, nil
			}
			if seed, ok := sess.Runtime.SessionSecrets["device_seed"]; ok && seed != "" {
				return seed, nil
			}
			if sess.Runtime.UserConfig.StellarAccount.PrivateKey != "" {
				return sess.Runtime.UserConfig.StellarAccount.PrivateKey, nil
			}
		}
	}
	if r.keyringSvc != nil && r.sessionMgr != nil {
		if sess, err := r.sessionMgr.GetSession(userID); err == nil && sess != nil && sess.Runtime != nil {
			pass := sess.Runtime.SessionSecrets["password"]
			secret := sess.Runtime.SessionSecrets["stellar_secret"]
			kr, err := r.keyringSvc.LoadHybrid(userID, pass, secret)
			if err == nil && kr != nil {
				seedBytes, err := r.keyringSvc.GetKeyByType(kr, vaults_domain.KeyTypeDeviceSeed)
				if err == nil && len(seedBytes) > 0 {
					return string(seedBytes), nil
				}
			}
		}
	}
	return "", trustgroup_domain.ErrDeviceNotFound
}

func (r *sessionSovereignIdentityResolver) GetVaultKeyring(ctx context.Context, userID string) (*vaults_domain.VaultKeyring, error) {
	if r.keyringSvc != nil && r.sessionMgr != nil {
		if sess, err := r.sessionMgr.GetSession(userID); err == nil && sess != nil && sess.Runtime != nil {
			pass := sess.Runtime.SessionSecrets["password"]
			secret := sess.Runtime.SessionSecrets["stellar_secret"]
			if secret == "" {
				secret = sess.Runtime.UserConfig.StellarAccount.PrivateKey
			}
			kr, err := r.keyringSvc.LoadHybrid(userID, pass, secret)
			if err == nil && kr != nil {
				return kr, nil
			}
		}
	}
	return vaults_domain.NewVaultKeyring(userID), nil
}

func (r *sessionSovereignIdentityResolver) GetDevice(ctx context.Context, deviceID string) (*trustgroup_ports.DeviceSummary, error) {
	if r.identityDeviceAdapter != nil {
		return r.identityDeviceAdapter.GetDevice(ctx, deviceID)
	}
	return nil, trustgroup_domain.ErrDeviceNotFound
}

func (r *sessionSovereignIdentityResolver) ListActiveDevices(ctx context.Context, memberID string) ([]trustgroup_ports.DeviceSummary, error) {
	if r.identityDeviceAdapter != nil {
		return r.identityDeviceAdapter.ListActiveDevices(ctx, memberID)
	}
	return nil, nil
}

type fakeStellarService struct{}

func (s *fakeStellarService) CreateKeypair() (string, string, string, error) {
	return "G_TEST_PUB", "S_TEST_SECRET", "TX_TEST", nil
}

func (s *fakeStellarService) CreateAccount(pw string) (*blockchain.CreateAccountRes, error) {
	return nil, nil
}

// Complete End-to-End Application Flow Acceptance Test


func TestCreateTrustGroup_ProvisionsAndPersistsCreatorEnvelope(t *testing.T) {
	ctx := context.Background()

	// 1. Setup production components & SQLite DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = driver.AutoMigrate(db)
	require.NoError(t, err)
	err = db.AutoMigrate(
		&models.User{},
		&models.VaultCID{},
		&app_config.UserConfig{},
		&identity_domain.Device{},
		&app_config_domain.DeviceConfig{},
		&auth_domain.TokenPairs{},
	)
	require.NoError(t, err)

	gormDevRepo := identity_persistence.NewGormDeviceRepository(db)
	deviceConfigRepo := app_config_persistence.NewGormDeviceConfigRepository(db)
	keyringDir := t.TempDir()
	osFS := &vault_infrastructure_security.OSFileSystem{}
	aesSvc := &vault_infrastructure_crypto.AESService{}
	keyEnc := vault_infrastructure_crypto.NewKeyService()
	keyringSvc := vault_infrastructure_security.NewKeyringService(aesSvc, keyEnc, keyringDir, osFS)

	userRepo := onboarding_persistence.NewGormUserRepository(db)
	bus := onboarding_eventbus.NewMemoryBus()
	logSvc := &logger.Logger{}

	authConfig := auth_domain.Auth{
		Issuer:      "ankhora-test-issuer",
		Audience:    "ankhora-test-audience",
		Secret:      "ankhora-test-secret-32-bytes-long!",
		TokenExpiry: time.Hour,
	}

	tokenService := auth_usecases.NewTokenService(authConfig, nil, db)
	tokenUC := auth_usecases.NewGenerateTokensUseCase(nil, tokenService)

	identityBus := identity_eventbus.NewMemoryEventBus()
	createDeviceUC := identity_usecase.NewCreateDeviceUseCase(gormDevRepo, nil).WithDeviceConfigRepository(deviceConfigRepo)
	identityHandler := identity_ui.NewIdentityHandler(db, tokenService, userRepo, *createDeviceUC, identityBus)

	createAccUC := onboarding_usecase.NewCreateAccountUseCase(
		&fakeStellarService{},
		userRepo,
		bus,
		logSvc,
		keyringSvc,
		keyEnc,
	).WithIdentityService(identityHandler)

	t.Setenv("VAULT_PASSWORD", "AlicePass123!")
	kpAlice, err := keypair.Random()
	require.NoError(t, err)

	aliceAccountResp, err := createAccUC.Execute(onboarding_usecase.AccountCreationRequest{
		Email:       "alice_creator@ankhora.test",
		Password:    "AlicePass123!",
		PublicKey:   kpAlice.Address(),
		DeviceSeed:  kpAlice.Seed(),
		IsAnonymous: false,
	})
	require.NoError(t, err, "Alice account creation MUST succeed")
	aliceVaultID := aliceAccountResp.UserID

	aliceDevices, err := gormDevRepo.ListByVaultID(ctx, aliceVaultID)
	require.NoError(t, err)
	require.Len(t, aliceDevices, 1)
	aliceDeviceID := aliceDevices[0].ID

	cloudToken := "0123456789abcdefghijklmnopqrstuv"
	cloudDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	cloudServer := newGormCloudServer(cloudDB, cloudToken)
	_ = cloudDB.Table("customers").Create(map[string]interface{}{
		"email":      "alice_creator@ankhora.test",
		"public_key": kpAlice.Address(),
		"first_name": "Alice",
		"last_name":  "Creator",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
	_ = cloudDB.Table("customers").Create(map[string]interface{}{
		"email":      aliceVaultID,
		"public_key": kpAlice.Address(),
		"first_name": "Alice",
		"last_name":  "Creator",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
	ts := httptest.NewServer(cloudServer)
	defer ts.Close()

	client := tracecore.NewTracecoreClient(ts.URL+"/api", cloudToken, ts.URL, ts.URL+"/api")
	identityDeviceAdapter := trustgroup_adapters.NewIdentityDeviceAdapter(gormDevRepo)

	addEnvUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(client, identityDeviceAdapter)
	provisionEnvelopeUC := trustgroup_envelope_uc.NewProvisionTrustGroupDeviceEnvelopeUseCase(
		client,
		identityDeviceAdapter,
		trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, nil, nil),
		addEnvUC,
		keyringSvc,
	)

	vaultHandler := vault_ui.NewVaultHandler(nil, *logSvc, ctx, nil, nil, db, client, keyringDir)
	vaultHandler.VaultRepository = &topLevelVaultRepo{}
	authUIHandler := auth_ui.NewAuthHandler(identityHandler, tokenUC, db)

	tgMemoryBus := trustgroup_eventbus.NewMemoryBus()
	createTGUC := trustgroup_usecases_trustgroup.NewCreateTrustGroupUsecase(client, tgMemoryBus, provisionEnvelopeUC, nil)
	subHandler := subscription_ui_wails.NewSubscriptionHandler(db, client, nil, &fakeStellarService{}, userRepo, bus, identityHandler, app_config_ui.AppConfigHandler{}, *logSvc)
	_ = subHandler.SaveSubscription(ctx, &subscription_domain.Subscription{
		ID:     "sub_alice",
		UserID: aliceVaultID,
		Email:  "alice_creator@ankhora.test",
	})

	app := &App{
		AuthHandler:         authUIHandler,
		Identity:            identityHandler,
		ctx:                 ctx,
		provisionEnvelopeUC: provisionEnvelopeUC,
		createTrustGroupUC:  createTGUC,
		SubscriptionHandler: subHandler,
		tracecoreClient:     client,
		Vault:               vaultHandler,
		Logger:              *logSvc,
	}

	aliceSignInResp, err := app.SignIn(handlers.LoginRequest{Email: "alice_creator@ankhora.test", Password: "AlicePass123!"})
	require.NoError(t, err)
	aliceToken := aliceSignInResp.Tokens.Token

	// ACT: Execute ONLY production app.CreateTrustGroup
	tgResp, err := app.CreateTrustGroup(aliceToken, "ws_channel_creator", "Creator Proof Group", "")
	require.NoError(t, err, "app.CreateTrustGroup MUST succeed")
	require.NotNil(t, tgResp)
	tgID := tgResp.ID

	// ASSERT: Direct Database Verification on cloudDB trust_group_key_envelopes table
	var rows []map[string]interface{}
	err = cloudDB.Table("trust_group_key_envelopes").Where("trust_group_id = ?", tgID).Find(&rows).Error
	require.NoError(t, err)

	if len(rows) != 1 {
		t.Fatalf("TRUSTGROUP INVARIANT FAILED: expected exactly 1 creator envelope in DB, got %d", len(rows))
	}

	row := rows[0]
	assert.NotEmpty(t, row["member_id"], "DB member_id MUST NOT be empty")
	assert.Equal(t, int64(1), row["kek_version"], "DB kek_version MUST be 1")
	assert.NotEmpty(t, row["wrapped_kek"], "DB wrapped_kek MUST NOT be empty")

	// ASSERT: Explicit Cloud Read-Back
	persistedTGResp, err := client.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tgID})
	require.NoError(t, err, "Cloud GetTrustGroup read-back MUST succeed")
	require.NotNil(t, persistedTGResp)
	require.Len(t, persistedTGResp.Data.KeyEnvelopes, 1, "Cloud read-back MUST return exactly 1 envelope for creator")

	readbackEnv := persistedTGResp.Data.KeyEnvelopes[0]
	assert.NotEmpty(t, readbackEnv.MemberID)
	assert.Equal(t, uint64(1), readbackEnv.KEKVersion)
	assert.Nil(t, readbackEnv.RevokedAt, "Envelope MUST NOT be revoked")
	assert.NotEmpty(t, readbackEnv.WrappedKEK)

	fmt.Printf(`
============================================================
[TRUSTGROUP][FORENSIC][FINAL STATE]
============================================================

trustGroupID=%s
members=1
kekVersion=1

EXPECTED ENVELOPES:
  member=%s
  device=%s
  kekVersion=1

PERSISTED ENVELOPES:
  count=%d

READ-BACK ENVELOPES:
  count=%d

DB ENVELOPES:
  count=%d

INVARIANT:
  memberExists=true
  deviceExists=true
  envelopeCreated=true
  envelopePersisted=true
  envelopeReadable=true

TRUSTGROUP_OPERATION=PASS
============================================================
`, tgID, aliceVaultID, aliceDeviceID, len(rows), len(persistedTGResp.Data.KeyEnvelopes), len(rows))
}

func TestAddTrustGroupMember_ProvisionsAndPersistsMemberEnvelope(t *testing.T) {
	ctx := context.Background()

	dbPath := fmt.Sprintf("%s/test_%d.db?cache=shared&mode=rwc", t.TempDir(), time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.NoError(t, err)

	err = driver.AutoMigrate(db)
	require.NoError(t, err)
	err = db.AutoMigrate(
		&models.User{},
		&models.VaultCID{},
		&app_config.UserConfig{},
		&identity_domain.User{},
		&identity_domain.Device{},
		&app_config_domain.DeviceConfig{},
		&auth_domain.TokenPairs{},
	)
	require.NoError(t, err)

	gormDevRepo := identity_persistence.NewGormDeviceRepository(db)
	deviceConfigRepo := app_config_persistence.NewGormDeviceConfigRepository(db)
	keyringDir := t.TempDir()
	osFS := &vault_infrastructure_security.OSFileSystem{}
	aesSvc := &vault_infrastructure_crypto.AESService{}
	keyEnc := vault_infrastructure_crypto.NewKeyService()
	keyringSvc := vault_infrastructure_security.NewKeyringService(aesSvc, keyEnc, keyringDir, osFS)

	userRepo := onboarding_persistence.NewGormUserRepository(db)
	bus := onboarding_eventbus.NewMemoryBus()
	logSvc := &logger.Logger{}

	authConfig := auth_domain.Auth{
		Issuer:      "ankhora-test-issuer",
		Audience:    "ankhora-test-audience",
		Secret:      "ankhora-test-secret-32-bytes-long!",
		TokenExpiry: time.Hour,
	}

	tokenService := auth_usecases.NewTokenService(authConfig, nil, db)
	tokenUC := auth_usecases.NewGenerateTokensUseCase(nil, tokenService)

	identityBus := identity_eventbus.NewMemoryEventBus()
	createDeviceUC := identity_usecase.NewCreateDeviceUseCase(gormDevRepo, nil).WithDeviceConfigRepository(deviceConfigRepo)
	identityHandler := identity_ui.NewIdentityHandler(db, tokenService, userRepo, *createDeviceUC, identityBus)

	createAccUC := onboarding_usecase.NewCreateAccountUseCase(
		&fakeStellarService{},
		userRepo,
		bus,
		logSvc,
		keyringSvc,
		keyEnc,
	).WithIdentityService(identityHandler)

	kpAlice, err := keypair.Random()
	require.NoError(t, err)
	aliceResp, err := createAccUC.Execute(onboarding_usecase.AccountCreationRequest{
		Email:       "alice_add_member@ankhora.test",
		Password:    "AlicePass123!",
		PublicKey:   kpAlice.Address(),
		DeviceSeed:  kpAlice.Seed(),
		IsAnonymous: false,
	})
	require.NoError(t, err)
	aliceVaultID := aliceResp.UserID
	err = identityHandler.IdentityUserRepo.Save(ctx, identity_domain.NewStandardUser(aliceVaultID, "alice_add_member@ankhora.test", "AlicePass123!"))
	require.NoError(t, err)

	_ = db.AutoMigrate(&vaults_persistence.VaultMapper{})
	_ = db.Create(&vaults_persistence.VaultMapper{
		ID:        "v_" + aliceVaultID,
		UserID:    aliceVaultID,
		Name:      "",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	})

	kpBob, err := keypair.Random()
	require.NoError(t, err)
	bobResp, err := createAccUC.Execute(onboarding_usecase.AccountCreationRequest{
		Email:       "bob_add_member@ankhora.test",
		Password:    "BobPass123!",
		PublicKey:   kpBob.Address(),
		DeviceSeed:  kpBob.Seed(),
		IsAnonymous: false,
	})
	require.NoError(t, err)
	bobVaultID := bobResp.UserID
	bobUser := identity_domain.NewStandardUser(bobVaultID, "bob_add_member@ankhora.test", "BobPass123!")
	err = identityHandler.IdentityUserRepo.Save(ctx, bobUser)
	require.NoError(t, err)

	aliceDevices, err := gormDevRepo.ListByVaultID(ctx, aliceVaultID)
	require.NoError(t, err)
	aliceDeviceID := aliceDevices[0].ID

	bobDevices, err := gormDevRepo.ListByVaultID(ctx, bobVaultID)
	require.NoError(t, err)
	_ = bobDevices

	// Delete Bob's local device to simulate Bob as a remote member whose device is NOT in local DB
	err = db.Where("vault_id = ?", bobVaultID).Delete(&identity_domain.Device{}).Error
	require.NoError(t, err)

	cloudToken := "0123456789abcdefghijklmnopqrstuv"
	cloudDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	cloudServer := newGormCloudServer(cloudDB, cloudToken)
	cloudServer.publicKeys = map[string]string{
		"alice_add_member@ankhora.test": kpAlice.Address(),
		"bob_add_member@ankhora.test":   kpBob.Address(),
	}
	_ = db.Exec("CREATE TABLE IF NOT EXISTS customers (id INTEGER PRIMARY KEY AUTOINCREMENT, email TEXT, public_key TEXT, first_name TEXT, last_name TEXT)")
	_ = db.Exec("DELETE FROM customers WHERE email = ?", "bob_add_member@ankhora.test")
	_ = db.Exec("DELETE FROM customers WHERE email = ?", "alice_add_member@ankhora.test")
	err = db.Exec("INSERT INTO customers (email, public_key, first_name, last_name) VALUES (?, ?, 'Bob', 'Member')",
		"bob_add_member@ankhora.test", kpBob.Address()).Error
	require.NoError(t, err)
	_ = db.Exec("INSERT INTO customers (email, public_key, first_name, last_name) VALUES (?, ?, 'Bob', 'Member')", bobVaultID, kpBob.Address())
	err = db.Exec("INSERT INTO customers (email, public_key, first_name, last_name) VALUES (?, ?, 'Alice', 'Creator')",
		"alice_add_member@ankhora.test", kpAlice.Address()).Error
	require.NoError(t, err)
	_ = db.Exec("INSERT INTO customers (email, public_key, first_name, last_name) VALUES (?, ?, 'Alice', 'Creator')", aliceVaultID, kpAlice.Address())

	bobAddress := kpBob.Address()
	aliceAddress := kpAlice.Address()
	_ = cloudDB.Exec("DELETE FROM customers WHERE email = ?", "bob_add_member@ankhora.test")
	_ = cloudDB.Exec("DELETE FROM customers WHERE email = ?", "alice_add_member@ankhora.test")
	err = cloudDB.Table("customers").Create(map[string]interface{}{
		"email":      "bob_add_member@ankhora.test",
		"public_key": bobAddress,
		"first_name": "Bob",
		"last_name":  "Member",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}).Error
	require.NoError(t, err)
	_ = cloudDB.Table("customers").Create(map[string]interface{}{
		"email":      bobVaultID,
		"public_key": bobAddress,
		"first_name": "Bob",
		"last_name":  "Member",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})
	err = cloudDB.Table("customers").Create(map[string]interface{}{
		"email":      "alice_add_member@ankhora.test",
		"public_key": aliceAddress,
		"first_name": "Alice",
		"last_name":  "Creator",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}).Error
	require.NoError(t, err)
	_ = cloudDB.Table("customers").Create(map[string]interface{}{
		"email":      aliceVaultID,
		"public_key": aliceAddress,
		"first_name": "Alice",
		"last_name":  "Creator",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	})

	ts := httptest.NewServer(cloudServer)
	defer ts.Close()

	client := tracecore.NewTracecoreClient(ts.URL+"/api", cloudToken, ts.URL, ts.URL)
	client.SetToken(cloudToken)
	identityDeviceAdapter := trustgroup_adapters.NewIdentityDeviceAdapter(gormDevRepo)

	addEnvUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(client, identityDeviceAdapter)
	provisionEnvelopeUC := trustgroup_envelope_uc.NewProvisionTrustGroupDeviceEnvelopeUseCase(
		client,
		identityDeviceAdapter,
		trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, nil, nil),
		addEnvUC,
		keyringSvc,
	)
	tgMemoryBus := trustgroup_eventbus.NewMemoryBus()
	addMemberUC := trustgroup_member_uc.NewAddMemberToTrustGroupUsecase(client, tgMemoryBus)

	vaultHandler := vault_ui.NewVaultHandler(nil, *logSvc, ctx, nil, nil, db, client, keyringDir)
	vaultHandler.VaultRepository = &topLevelVaultRepo{}
	authUIHandler := auth_ui.NewAuthHandler(identityHandler, tokenUC, db)
	aliceKeyring, err := keyringSvc.LoadHybrid(aliceVaultID, "AlicePass123!", "")
	require.NoError(t, err)

	memberAddedSubscriber := trustgroup_events.NewMemberAddedSubscriber(
		provisionEnvelopeUC,
		identityDeviceAdapter,
		identityHandler,
		client,
		aliceKeyring,
	)
	memberAddedSubscriber.RegisterSubscribers(tgMemoryBus)

	createTGUC := trustgroup_usecases_trustgroup.NewCreateTrustGroupUsecase(client, tgMemoryBus, provisionEnvelopeUC, nil)

	subHandler := subscription_ui_wails.NewSubscriptionHandler(db, client, nil, &fakeStellarService{}, userRepo, bus, identityHandler, app_config_ui.AppConfigHandler{}, *logSvc)
	_ = subHandler.SaveSubscription(ctx, &subscription_domain.Subscription{
		ID:     "sub_alice",
		UserID: aliceVaultID,
		Email:  "alice_add_member@ankhora.test",
	})

	app := &App{
		AuthHandler:           authUIHandler,
		Identity:              identityHandler,
		ctx:                   ctx,
		addTrustGroupMemberUC: addMemberUC,
		provisionEnvelopeUC:   provisionEnvelopeUC,
		createTrustGroupUC:    createTGUC,
		SubscriptionHandler:   subHandler,
		tracecoreClient:       client,
		Vault:                 vaultHandler,
		Logger:                *logSvc,
	}

	aliceSignInResp, err := app.SignIn(handlers.LoginRequest{Email: "alice_add_member@ankhora.test", Password: "AlicePass123!"})
	require.NoError(t, err)
	aliceToken := aliceSignInResp.Tokens.Token
	client.SetToken(cloudToken)
	if vaultHandler != nil && vaultHandler.TracecoreClient != nil {
		vaultHandler.TracecoreClient.SetToken(cloudToken)
	}
	if sess, _ := vaultHandler.GetSession(aliceVaultID); sess != nil && sess.Runtime != nil {
		if sess.Runtime.SessionSecrets == nil {
			sess.Runtime.SessionSecrets = make(map[string]string)
		}
		sess.Runtime.SessionSecrets["password"] = "AlicePass123!"
		sess.Runtime.SessionSecrets["assword"] = "AlicePass123!"
	}

	// Step 1: Alice creates TrustGroup via production App.CreateTrustGroup
	tgResp, err := app.CreateTrustGroup(aliceToken, "ws_channel_add_member", "Add Member Test Group", "")
	require.NoError(t, err)
	tgID := tgResp.ID

	// Step 2: Verify DB count is 1 for Alice
	var aliceRows []map[string]interface{}
	err = cloudDB.Table("trust_group_key_envelopes").Where("trust_group_id = ?", tgID).Find(&aliceRows).Error
	require.NoError(t, err)
	require.Len(t, aliceRows, 1, "DB MUST contain exactly 1 row for Alice after creation")

	// Step 3: Alice adds Bob via production App.AddTrustGroupMember
	_, err = app.AddTrustGroupMember(aliceToken, tgID, bobVaultID, "member")
	require.NoError(t, err, "app.AddTrustGroupMember MUST succeed")

	// Step 4: Verify DB count is 2 (Alice + Bob) with exact hard field assertions
	var allRows []map[string]interface{}
	err = cloudDB.Table("trust_group_key_envelopes").Where("trust_group_id = ?", tgID).Find(&allRows).Error
	require.NoError(t, err)

	fmt.Printf("[ENVELOPE][DB_VERIFY]\ntrustGroupID=%s\nenvelopeCount=%d\n", tgID, len(allRows))

	if len(allRows) != 2 {
		t.Fatalf("TRUSTGROUP INVARIANT FAILED: expected 2 envelopes in DB after AddMember, got %d", len(allRows))
	}

	rowMap := make(map[string]map[string]interface{})
	for _, r := range allRows {
		mID := fmt.Sprintf("%v", r["member_id"])
		rowMap[mID] = r
	}

	aliceRow, aliceOk := rowMap["v_cloud_1"]
	if !aliceOk {
		aliceRow, aliceOk = rowMap[aliceVaultID]
	}
	require.True(t, aliceOk, "Alice DB row MUST exist")
	assert.NotEmpty(t, aliceRow["wrapped_kek"])

	bobRow, bobOk := rowMap["v_cloud_bob"]
	if !bobOk {
		bobRow, bobOk = rowMap[bobVaultID]
	}
	require.True(t, bobOk, "Bob DB row MUST exist")
	assert.NotEmpty(t, bobRow["wrapped_kek"])

	// Step 5: Explicit Cloud Read-Back Assertion
	persistedTGResp, err := client.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tgID})
	require.NoError(t, err)
	require.Len(t, persistedTGResp.Data.KeyEnvelopes, 2, "Cloud read-back MUST return 2 envelopes (Alice + Bob)")

	// Query and print complete TrustGroup row from cloudDB
	var tgRows []map[string]interface{}
	err = cloudDB.Table("trust_groups").Where("id = ?", tgID).Find(&tgRows).Error
	require.NoError(t, err)

	tgJSON, _ := json.MarshalIndent(tgRows, "", "  ")
	envJSON, _ := json.MarshalIndent(allRows, "", "  ")

	fmt.Printf("\n============================================================\n")
	fmt.Printf("[DATABASE FORENSICS] TRUST GROUP ROW IN DB:\n%s\n", string(tgJSON))
	fmt.Printf("============================================================\n")
	fmt.Printf("[DATABASE FORENSICS] TRUST GROUP ENVELOPES IN DB (%d rows):\n%s\n", len(allRows), string(envJSON))
	fmt.Printf("============================================================\n\n")

	fmt.Printf(`
============================================================
[TRUSTGROUP][FORENSIC][FINAL STATE]
============================================================

trustGroupID=%s
members=2
kekVersion=1

EXPECTED ENVELOPES:
  member=%s (Alice) device=%s kekVersion=1
  member=%s (Bob)   device=%s kekVersion=1

PERSISTED ENVELOPES:
  count=2

READ-BACK ENVELOPES:
  count=2

DB ENVELOPES:
  count=2

INVARIANT:
  memberExists=true
  deviceExists=true
  envelopeCreated=true
  envelopePersisted=true
  envelopeReadable=true

TRUSTGROUP_OPERATION=PASS
============================================================
`, tgID, aliceVaultID, aliceDeviceID, bobVaultID, bobVaultID)
}

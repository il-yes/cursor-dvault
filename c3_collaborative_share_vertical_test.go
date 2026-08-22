package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	auth_usecases "vault-app/internal/auth/application/use_cases"
	auth_domain "vault-app/internal/auth/domain"
	auth_ui "vault-app/internal/auth/ui"
	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	collaboration_ui "vault-app/internal/collaboration/ui"
	"vault-app/internal/tracecore"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
)

// ---------------------------------------------------------------------------
// Contract-faithful Cloud stub (verified ankhora-cloud C1/C2/C3 contracts)
// ---------------------------------------------------------------------------

type cloudStub struct {
	mu           sync.Mutex
	trustGroups  map[string]map[string]interface{}
	shareEntries map[string]c3_asset_domain.ShareEntry
	token        string
}

func newCloudStub(token string) *cloudStub {
	return &cloudStub{
		trustGroups:  map[string]map[string]interface{}{},
		shareEntries: map[string]c3_asset_domain.ShareEntry{},
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
		if uint64(kekVersion) != tg["kek_version"].(uint64) {
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
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), nil)
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, nil, nil)

	tgID := "tg-vertical-" + uuid.NewString()[:8]
	const kekVersion = uint64(2)

	kr := &vaults_domain.VaultKeyring{UserID: "user_alice", VaultID: "vault_alice"}
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

	hash := sha256.Sum256(prepared.EncryptedData)
	assetCID := "bafybeivertical" + hex.EncodeToString(hash[:8])

	// Binary key material must travel base64-encoded over the JSON wire
	// contract (the resolve path decodes Base64-or-raw accordingly).
	wrappedDEKB64 := base64.StdEncoding.EncodeToString(prepared.WrappedDEK)

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
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, nil)
	collabHandler := collaboration_ui.NewCollaborationHandler(createCollabShareUC, nil, nil)

	app := &App{
		AuthHandler:          authHandler,
		CollaborationHandler: collabHandler,
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
		"vault_partner_02",
		"vertical persistence test",
		wrappedDEKB64,
		prepared.KEKVersion,
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
	assert.Equal(t, assetCID, persisted.AssetCID)
	assert.Equal(t, wrappedDEKB64, persisted.WrappedDEK)
	decoded, decErr := base64.StdEncoding.DecodeString(persisted.WrappedDEK)
	require.NoError(t, decErr)
	assert.Equal(t, prepared.WrappedDEK, decoded, "crypto material must survive the round-trip byte-exact")
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
	// 6. Negative paths — invalid trust group and stale KEK are rejected
	// ------------------------------------------------------------------
	_, err = app.CreateCollaborativeShare(pairs.Token, "thread_vertical_1", "tg-does-not-exist", assetCID, "v", "", wrappedDEKB64, kekVersion)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")

	_, err = app.CreateCollaborativeShare(pairs.Token, "thread_vertical_1", tgID, assetCID, "v", "", wrappedDEKB64, kekVersion+7)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "stale",
		"stale KEK must be rejected (desktop use case guard and/or Cloud 409)")

	// Missing crypto material must fail loudly instead of fabricating.
	_, err = app.CreateCollaborativeShare(pairs.Token, "thread_vertical_1", tgID, assetCID, "v", "", "", kekVersion)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "wrapped_dek is required")
}

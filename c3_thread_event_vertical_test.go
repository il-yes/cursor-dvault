package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
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
	thread_usecase "vault-app/internal/thread/application/usecases"
	thread_domain "vault-app/internal/thread/domain"
	thread_ui "vault-app/internal/thread/ui"
	"vault-app/internal/tracecore"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
	vault_ui "vault-app/internal/vault/ui"
)

// ---------------------------------------------------------------------------
// Cloud stub extending the verified C1/C2/C3 contracts with the ThreadEvent
// contract: POST /api/threads/{id}/events and GET /api/threads/{id}/events.
// ---------------------------------------------------------------------------

type threadEventStub struct {
	mu           sync.Mutex
	trustGroups  map[string]map[string]interface{}
	shareEntries map[string]c3_asset_domain.ShareEntry
	events       map[string][]map[string]interface{} // threadID -> persisted events
	appendBodies []string                            // raw captured POST bodies
	token        string
}

func newThreadEventStub(token string) *threadEventStub {
	return &threadEventStub{
		trustGroups:  map[string]map[string]interface{}{},
		shareEntries: map[string]c3_asset_domain.ShareEntry{},
		events:       map[string][]map[string]interface{}{},
		token:        token,
	}
}

func (s *threadEventStub) seedTrustGroup(id string, kekVersion uint64, memberCIDs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	envelopes := make([]map[string]interface{}, 0, len(memberCIDs))
	for _, cid := range memberCIDs {
		envelopes = append(envelopes, map[string]interface{}{
			"id":             uuid.NewString(),
			"trust_group_id": id,
			"member_id":      cid,
			"device_id":      "",
			"kek_version":    kekVersion,
			"wrapped_kek":    "wrapped-kek-" + cid,
			"created_at":     time.Now().UTC().Format(time.RFC3339),
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

func (s *threadEventStub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/events"):
		bodyBytes, _ := io.ReadAll(r.Body)
		s.appendBodies = append(s.appendBodies, string(bodyBytes))

		var body map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, ok := body["thread_id"].(string); !ok {
			http.Error(w, "thread_id is required", http.StatusBadRequest)
			return
		}

		parts := strings.Split(r.URL.Path, "/") // /api/threads/{id}/events
		storedThreadID := parts[len(parts)-2]

		eventID := uuid.NewString()
		prev := ""
		if existing := s.events[storedThreadID]; len(existing) > 0 {
			prev, _ = existing[len(existing)-1]["id"].(string)
		}
		event := map[string]interface{}{
			"id":              eventID,
			"thread_id":       storedThreadID,
			"type":            body["type"],
			"payload":         body["payload"],
			"idempotency_key": body["idempotency_key"],
			"cursor":          len(s.events[storedThreadID]),
			"created_at":      time.Now().UTC().Format(time.RFC3339Nano),
		}
		if prev != "" {
			event["previous_event_id"] = prev
		}
		s.events[storedThreadID] = append(s.events[storedThreadID], event)
		writeEnvelope(w, http.StatusCreated, event)

	case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/events"):
		parts := strings.Split(r.URL.Path, "/")
		threadID := parts[len(parts)-2]
		writeEnvelope(w, http.StatusOK, s.events[threadID])

	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

// ---------------------------------------------------------------------------
// Vertical test: persisted C3 ShareEntry -> entry.shared ThreadEvent
// ---------------------------------------------------------------------------

// TestAppendThreadEvent_ReferencesPersistedShareEntry exercises the real Wails
// entry point used by AppendThreadEventSlidingView:
//
//	App.AppendThreadEvent (JWT auth, flat EventResourceRef JSON)
//	  -> ThreadHandler.AppendThreadEvent
//	    -> AppendThreadEventUsecase.Execute (derives evt_share_<id>)
//	      -> TracecoreClient.AppendThreadEvent (eventResourceRefToPayload)
//	        -> POST /api/threads/{id}/events
//
// The ShareEntry reference submitted is the REAL one produced by the
// CreateCollaborativeShare flow in the same session — never fabricated.
func TestAppendThreadEvent_ReferencesPersistedShareEntry(t *testing.T) {
	ctx := context.Background()

	// ------------------------------------------------------------------
	// 1. Real crypto + real share creation (STEP 1-4 of the workflow)
	// ------------------------------------------------------------------
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, t.TempDir(), nil)
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, nil, nil)

	tgID := "tg-vertical-" + uuid.NewString()[:8]
	const kekVersion = uint64(3)

	prepared, err := orchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset-vertical-002",
		TrustGroupID: tgID,
		KEKVersion:   kekVersion,
		RawPayload:   []byte("CONFIDENTIAL THREAD EVENT VERTICAL PAYLOAD"),
		Keyring:      &vaults_domain.VaultKeyring{UserID: "user_alice", VaultID: "vault_alice"},
	})
	require.NoError(t, err)

	hash := sha256.Sum256(prepared.EncryptedData)
	assetCID := "bafybeievent" + hex.EncodeToString(hash[:8])
	wrappedDEKB64 := base64.StdEncoding.EncodeToString(prepared.WrappedDEK)

	cloudToken := "0123456789abcdefghijklmnopqrstuv"
	stub := newThreadEventStub(cloudToken)
	stub.seedTrustGroup(tgID, kekVersion, []string{"vault_alice"})
	ts := httptest.NewServer(stub)
	defer ts.Close()

	tc := tracecore.NewTracecoreClient(ts.URL+"/api", cloudToken, ts.URL, ts.URL+"/api")

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

	appendEventUC := thread_usecase.NewAppendThreadEventUsecase(tc)
	listEventsUC := thread_usecase.NewListThreadEventsUsecase(tc)
	threadHandler := thread_ui.NewThreadHandler(nil, nil, listEventsUC, appendEventUC)

	app := &App{
		AuthHandler:          authHandler,
		CollaborationHandler: collabHandler,
		ThreadHandler:        threadHandler,
		Vault:                &vault_ui.VaultHandler{}, // SessionManager nil -> RestoreCloudTokenForUser no-ops
		ctx:                  ctx,
	}

	pairs, err := tokenService.GenerateTokenPair(&auth_domain.JwtUser{
		ID:       "user_alice",
		Username: "alice",
		Email:    "alice@ankhora.test",
	})
	require.NoError(t, err)

	const threadID = "thread_vertical_events_1"

	shareRef, err := app.CreateCollaborativeShare(
		pairs.Token, threadID, tgID, assetCID, "vault_partner_02",
		"thread event vertical test", wrappedDEKB64, prepared.KEKVersion,
	)
	require.NoError(t, err)
	require.NotEmpty(t, shareRef.ShareEntryID)
	t.Logf("STEP 1-3 done: real persisted C3 ShareEntry %s", shareRef.ShareEntryID)

	// activeShareEntryRef now holds exactly this value (store.createShare).

	// ------------------------------------------------------------------
	// 2. STEP 8 — the exact flat payload the fixed frontend submits
	// ------------------------------------------------------------------
	frontendPayload, marshalErr := json.Marshal(map[string]string{
		"ref_type":       "share_entry",
		"share_entry_id": shareRef.ShareEntryID,
		"trust_group_id": tgID,
	})
	require.NoError(t, marshalErr)

	evtDTO, err := app.AppendThreadEvent(pairs.Token, threadID, "entry.shared", string(frontendPayload))
	require.NoError(t, err)
	require.NotNil(t, evtDTO)

	// ------------------------------------------------------------------
	// 3. Assert B — Wails boundary reconstructed a real EventResourceRef
	// ------------------------------------------------------------------
	assert.Equal(t, "entry.shared", evtDTO.Type)
	assert.Equal(t, thread_domain.ResourceShareEntry, thread_domain.ResourceType(evtDTO.Payload["ref_type"].(string)),
		"App.AppendThreadEvent must unmarshal the flat payload into RefType=share_entry")
	assert.Equal(t, shareRef.ShareEntryID, evtDTO.Payload["share_entry_id"],
		"boundary must carry the real ShareEntry ID, not a zero-value ref")
	assert.Equal(t, tgID, evtDTO.Payload["trust_group_id"])

	// ------------------------------------------------------------------
	// 4. Assert C + D + F — wire shape, idempotency key, no fabrication
	// ------------------------------------------------------------------
	stub.mu.Lock()
	require.Len(t, stub.appendBodies, 1)
	rawBody := stub.appendBodies[0]
	stub.mu.Unlock()

	var wireBody map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(rawBody), &wireBody))

	assert.Equal(t, "entry.shared", wireBody["type"], "C: event type on the wire")
	assert.Equal(t, threadID, wireBody["thread_id"], "C: thread id on the wire")

	wirePayload, ok := wireBody["payload"].(map[string]interface{})
	require.True(t, ok, "C: payload must be an object")
	assert.Equal(t, "share_entry", wirePayload["ref_type"], "C: flat ref_type")
	assert.Equal(t, shareRef.ShareEntryID, wirePayload["share_entry_id"], "C: flat share_entry_id")
	assert.Equal(t, tgID, wirePayload["trust_group_id"], "C: flat trust_group_id")
	assert.Nil(t, wireBody["payload_ref"], "no nested payload_ref may remain")
	assert.NotContains(t, rawBody, "share_entry_ref", "no nested share_entry_ref may remain")
	assert.NotContains(t, rawBody, "trust_group_ref", "no nested trust_group_ref may remain")

	expectedKey := "evt_share_" + shareRef.ShareEntryID
	assert.Equal(t, expectedKey, wireBody["idempotency_key"],
		"D: use case must derive evt_share_<real-share-entry-id>")

	assert.NotContains(t, rawBody, "\"se_", "F: no fabricated se_* identifier anywhere")
	assert.False(t, strings.Contains(rawBody, "tg_legal_counsel"), "F: no placeholder trust group")

	// ------------------------------------------------------------------
	// 5. Assert E — persisted event round-trips through GET
	// ------------------------------------------------------------------
	events, err := app.ListThreadEvents(pairs.Token, threadID)
	require.NoError(t, err)
	require.Len(t, events, 1)

	persistedEvt := events[0]
	assert.Equal(t, evtDTO.ID, persistedEvt.ID)
	assert.Equal(t, "entry.shared", persistedEvt.Type)
	assert.Equal(t, expectedKey, persistedEvt.IdempotencyKey)
	assert.Equal(t, shareRef.ShareEntryID, persistedEvt.Payload["share_entry_id"],
		"E: persisted event must reference the same ShareEntry")
	assert.Equal(t, tgID, persistedEvt.Payload["trust_group_id"])

	// Decisive invariant: ThreadEvent references the already-persisted ShareEntry.
	stub.mu.Lock()
	_, shareReallyPersisted := stub.shareEntries[persistedEvt.Payload["share_entry_id"].(string)]
	stub.mu.Unlock()
	assert.True(t, shareReallyPersisted,
		"the referenced ShareEntry must exist in persistence: createdShareEntry.ID == threadEvent.payload.share_entry_id")

	t.Logf("STEP 10-11 verified: ThreadEvent %s (idem=%s) references persisted ShareEntry %s",
		persistedEvt.ID, persistedEvt.IdempotencyKey, persistedEvt.Payload["share_entry_id"])
}

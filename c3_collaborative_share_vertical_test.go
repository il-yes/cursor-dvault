package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

	"vault-app/internal/auth"
	auth_usecases "vault-app/internal/auth/application/use_cases"
	auth_domain "vault-app/internal/auth/domain"
	auth_ui "vault-app/internal/auth/ui"
	"vault-app/internal/blockchain"
	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	collaboration_ui "vault-app/internal/collaboration/ui"
	app_config "vault-app/internal/config"
	"vault-app/internal/driver"
	handlers "vault-app/internal/handlers"
	identity_domain "vault-app/internal/identity/domain"
	identity_persistence "vault-app/internal/identity/infrastructure/persistence"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"
	onboarding_eventbus "vault-app/internal/onboarding/infrastructure/eventbus"
	onboarding_persistence "vault-app/internal/onboarding/infrastructure/persistence"
	"vault-app/internal/tracecore"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_ports "vault-app/internal/trust_group/application/ports"
	trustgroup_envelope_uc "vault-app/internal/trust_group/application/usecases/envelope"
	trustgroup_member_uc "vault-app/internal/trust_group/application/usecases/member"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	trustgroup_adapters "vault-app/internal/trust_group/infrastructure/adapters"
	trustgroup_eventbus "vault-app/internal/trust_group/infrastructure/eventbus"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
	vault_ui "vault-app/internal/vault/ui"
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

type memoryAssetStorageResolver struct {
	assets map[string][]byte
}

func (r *memoryAssetStorageResolver) FetchEncryptedAsset(ctx context.Context, cid string) ([]byte, error) {
	data, ok := r.assets[cid]
	if !ok {
		return nil, fmt.Errorf("asset content not found for CID %s", cid)
	}
	return data, nil
}

func (r *memoryAssetStorageResolver) Add(ctx context.Context, data []byte) (string, error) {
	cid := fmt.Sprintf("cid_enc_%d", time.Now().UnixNano())
	if r.assets == nil {
		r.assets = make(map[string][]byte)
	}
	r.assets[cid] = data
	return cid, nil
}

func (r *memoryAssetStorageResolver) Get(ctx context.Context, cid string) ([]byte, error) {
	return r.FetchEncryptedAsset(ctx, cid)
}

type memorySovereignIdentityResolver struct {
	seeds      map[string]string
	devices    map[string]*trustgroup_ports.DeviceSummary
	keyringSvc *vault_infrastructure_security.KeyringService
	sessionMgr *vault_session.Manager
}

func (r *memorySovereignIdentityResolver) GetDeviceSeed(ctx context.Context, userID string) (string, error) {
	seed, ok := r.seeds[userID]
	if !ok {
		return "", trustgroup_domain.ErrDeviceNotFound
	}
	return seed, nil
}

func (r *memorySovereignIdentityResolver) GetVaultKeyring(ctx context.Context, userID string) (*vaults_domain.VaultKeyring, error) {
	if r.keyringSvc != nil && r.sessionMgr != nil {
		if sess, err := r.sessionMgr.GetSession(userID); err == nil && sess != nil && sess.Runtime != nil {
			pass := sess.Runtime.SessionSecrets["password"]
			secret := sess.Runtime.SessionSecrets["stellar_secret"]
			kr, err := r.keyringSvc.LoadHybrid(userID, pass, secret)
			if err == nil && kr != nil {
				return kr, nil
			}
		}
	}
	return vaults_domain.NewVaultKeyring(userID), nil
}

func (r *memorySovereignIdentityResolver) GetDevice(ctx context.Context, deviceID string) (*trustgroup_ports.DeviceSummary, error) {
	dev, ok := r.devices[deviceID]
	if !ok {
		return nil, trustgroup_domain.ErrDeviceNotFound
	}
	return dev, nil
}

func (r *memorySovereignIdentityResolver) ListActiveDevices(ctx context.Context, memberID string) ([]trustgroup_ports.DeviceSummary, error) {
	var list []trustgroup_ports.DeviceSummary
	for _, dev := range r.devices {
		if dev != nil && dev.VaultID == memberID && dev.IsActive {
			list = append(list, *dev)
		}
	}
	return list, nil
}

type fakeStellarService struct{}

func (s *fakeStellarService) CreateKeypair() (string, string, string, error) {
	return "G_TEST_PUB", "S_TEST_SECRET", "TX_TEST", nil
}

func (s *fakeStellarService) CreateAccount(pw string) (*blockchain.CreateAccountRes, error) {
	return nil, nil
}

// Complete End-to-End Application Flow Acceptance Test
func TestC3_FullEndToEndUserAcceptanceFlow(t *testing.T) {
	ctx := context.Background()

	// 1. Setup User Identities & Keys
	kpAlice, err := keypair.Random()
	require.NoError(t, err)
	deviceAliceID := "dev_alice_macbook"

	kpBob, err := keypair.Random()
	require.NoError(t, err)
	deviceBobID := "dev_bob_workstation"

	// 2. Setup SQLite DB for repositories
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = driver.AutoMigrate(db)
	require.NoError(t, err)
	err = db.AutoMigrate(&models.User{}, &models.VaultCID{}, &app_config.UserConfig{}, &identity_domain.Device{})
	require.NoError(t, err)

	gormDevRepo := identity_persistence.NewGormDeviceRepository(db)
	identityDeviceAdapter := trustgroup_adapters.NewIdentityDeviceAdapter(gormDevRepo)

	// 3. Setup Keyring, Onboarding & Crypto Orchestration Services
	alicePassword := "AliceSecretPassword123!"
	bobPassword := "BobSecretPassword123!"

	osFS := &vault_infrastructure_security.OSFileSystem{}
	aesSvc := &vault_infrastructure_crypto.AESService{}
	keyEnc := vault_infrastructure_crypto.NewKeyService()
	keyringSvc := vault_infrastructure_security.NewKeyringService(aesSvc, keyEnc, t.TempDir(), osFS)
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	cryptoOrchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	userRepo := onboarding_persistence.NewGormUserRepository(db)
	bus := onboarding_eventbus.NewMemoryBus()
	logSvc := &logger.Logger{}

	createAccUC := onboarding_usecase.NewCreateAccountUseCase(
		&fakeStellarService{},
		userRepo,
		bus,
		logSvc,
		keyringSvc,
		keyEnc,
	).WithDeviceRepository(gormDevRepo)

	// Execute production Onboarding CreateAccountUseCase for Alice (creates real password-wrapped keyring file on disk)
	aliceAccountResp, err := createAccUC.Execute(onboarding_usecase.AccountCreationRequest{
		Email:       "alice@ankhora.test",
		Password:    alicePassword,
		PublicKey:   kpAlice.Address(),
		IsAnonymous: false,
	})
	require.NoError(t, err)
	aliceVaultID := aliceAccountResp.UserID

	// Execute production Onboarding CreateAccountUseCase for Bob (creates real password-wrapped keyring file on disk)
	bobAccountResp, err := createAccUC.Execute(onboarding_usecase.AccountCreationRequest{
		Email:       "bob@ankhora.test",
		Password:    bobPassword,
		PublicKey:   kpBob.Address(),
		IsAnonymous: false,
	})
	require.NoError(t, err)
	bobVaultID := bobAccountResp.UserID

	// 4. Cloud Stub & Client
	cloudToken := "0123456789abcdefghijklmnopqrstuv"
	stub := newCloudStub(cloudToken)
	ts := httptest.NewServer(stub)
	defer ts.Close()

	client := tracecore.NewTracecoreClient(ts.URL+"/api", cloudToken, ts.URL, ts.URL+"/api")
	cloudShareRepo := tracecore.NewCloudShareEntryRepository(client)

	addEnvUC := trustgroup_envelope_uc.NewAddTrustGroupKeyEnvelopeUseCase(client, identityDeviceAdapter)
	provisionEnvelopeUC := trustgroup_envelope_uc.NewProvisionTrustGroupDeviceEnvelopeUseCase(
		client,
		identityDeviceAdapter,
		cryptoOrchestrator,
		addEnvUC,
		keyringSvc,
	)
	addMemberUC := trustgroup_member_uc.NewAddMemberToTrustGroupUsecase(client, trustgroup_eventbus.NewMemoryBus())

	// 5. Setup Production Handlers, AuthHandler & Login Lifecycle
	dbModel := models.DBModel{DB: db}
	sessionMgr := vault_session.NewManager(nil, nil, logSvc, ctx, nil, make(map[string]*vault_session.Session))
	vaultHandler := &vault_ui.VaultHandler{
		KeyringService: keyringSvc,
		SessionManager: sessionMgr,
	}

	handlersVaultHandler := handlers.NewVaultHandler(dbModel, nil, nil, make(map[string]*models.VaultSession), logSvc, client, vault_session.RuntimeContext{})

	authConfig := auth.Auth{
		Issuer:      "ankhora",
		Audience:    "ankhora-desktop",
		Secret:      "e2e-acceptance-secret-key-32bytes",
		TokenExpiry: 15 * time.Minute,
	}

	prodAuthHandler := handlers.NewAuthHandler(dbModel, handlersVaultHandler, nil, logSvc, client, authConfig, userRepo).WithDeviceRepository(gormDevRepo)

	// Populate models.User records in database so AuthHandler.Login can construct models.User
	_, err = dbModel.CreateUser(&models.User{
		ID:       aliceVaultID,
		Email:    "alice@ankhora.test",
		Username: "alice",
	})
	require.NoError(t, err)

	_, err = dbModel.SaveVaultCID(models.VaultCID{
		ID:     uuid.New().String(),
		UserID: aliceVaultID,
		CID:    "cid_alice_vault",
		Name:   "alice-vault",
	})
	require.NoError(t, err)

	_, err = dbModel.CreateUser(&models.User{
		ID:       bobVaultID,
		Email:    "bob@ankhora.test",
		Username: "bob",
	})
	require.NoError(t, err)

	_, err = dbModel.SaveVaultCID(models.VaultCID{
		ID:     uuid.New().String(),
		UserID: bobVaultID,
		CID:    "cid_bob_vault",
		Name:   "bob-vault",
	})
	require.NoError(t, err)

	// Execute REAL production AuthHandler.Login for Alice
	aliceLoginResp, err := prodAuthHandler.Login(handlers.LoginRequest{
		Email:    "alice@ankhora.test",
		Password: alicePassword,
	})
	require.NoError(t, err, "Alice production AuthHandler.Login MUST succeed")
	require.NotNil(t, aliceLoginResp)
	require.NotNil(t, aliceLoginResp.Tokens)
	aliceRealToken := aliceLoginResp.Tokens.Token

	// Execute REAL production AuthHandler.Login for Bob
	bobLoginResp, err := prodAuthHandler.Login(handlers.LoginRequest{
		Email:    "bob@ankhora.test",
		Password: bobPassword,
	})
	require.NoError(t, err, "Bob production AuthHandler.Login MUST succeed")
	require.NotNil(t, bobLoginResp)

	// Architectural Boundary: handlers.AuthHandler (REST) populates handlers.VaultHandler.Sessions (models.VaultSession),
	// while App.AddTrustGroupMember uses vault_ui.VaultHandler (Wails UI) backed by SessionManager (vault_session.Session).
	// Pass the session populated by AuthHandler.Login directly to SessionManager for VaultHandler access.
	aliceHandlersSession, err := handlersVaultHandler.GetSession(aliceVaultID)
	require.NoError(t, err)
	_, err = sessionMgr.Prepare(aliceVaultID)
	require.NoError(t, err)
	_, err = sessionMgr.AttachRuntime(aliceVaultID, &vault_session.RuntimeContext{
		SessionSecrets: aliceHandlersSession.VaultRuntimeContext.SessionSecrets,
	})
	require.NoError(t, err)

	bobHandlersSession, err := handlersVaultHandler.GetSession(bobVaultID)
	require.NoError(t, err)
	_, err = sessionMgr.Prepare(bobVaultID)
	require.NoError(t, err)
	_, err = sessionMgr.AttachRuntime(bobVaultID, &vault_session.RuntimeContext{
		SessionSecrets: bobHandlersSession.VaultRuntimeContext.SessionSecrets,
	})
	require.NoError(t, err)

	identityResolver := &memorySovereignIdentityResolver{
		seeds: map[string]string{
			aliceVaultID: kpAlice.Seed(),
			bobVaultID:   kpBob.Seed(),
		},
		devices: map[string]*trustgroup_ports.DeviceSummary{
			deviceAliceID: {ID: deviceAliceID, VaultID: aliceVaultID, PublicKey: kpAlice.Address(), Status: "active", IsActive: true},
			deviceBobID:   {ID: deviceBobID, VaultID: bobVaultID, PublicKey: kpBob.Address(), Status: "active", IsActive: true},
		},
		keyringSvc: keyringSvc,
		sessionMgr: sessionMgr,
	}

	tokenService := auth_usecases.NewTokenService(auth_domain.Auth{
		Issuer:      authConfig.Issuer,
		Audience:    authConfig.Audience,
		Secret:      authConfig.Secret,
		TokenExpiry: authConfig.TokenExpiry,
	}, nil, nil)
	tokenUC := auth_usecases.NewGenerateTokensUseCase(nil, tokenService)
	authUIHandler := auth_ui.NewAuthHandler(nil, tokenUC, nil)

	app := &App{
		AuthHandler:           authUIHandler,
		ctx:                   ctx,
		addTrustGroupMemberUC: addMemberUC,
		provisionEnvelopeUC:   provisionEnvelopeUC,
		identityDeviceAdapter: identityDeviceAdapter,
		tracecoreClient:       client,
		Vault:                 vaultHandler,
	}

	originalPayload := []byte("TOP SECRET C3 E2E ACCEPTANCE TEST 2026")
	channelID := "ch_e2e_acceptance_room"

	// =========================================================================
	// ALICE STEPS
	// =========================================================================

	// [C3][E2E][01] Alice creates TrustGroup
	tg := trustgroup_domain.NewTrustGroup(channelID, "E2E Acceptance Group", []string{aliceVaultID})
	createTGResp, err := client.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)
	tgID := createTGResp.Data.ID
	require.NotEmpty(t, tgID)
	fmt.Printf("[C3][E2E][01] Alice creates TrustGroup id=%s\n", tgID)

	// Load Alice encrypted keyring from disk with her real password and store KEK v1
	aliceEncryptedKeyring, err := keyringSvc.LoadHybrid(aliceVaultID, alicePassword, "")
	require.NoError(t, err)
	_, err = keyringSvc.StoreTrustGroupKEK(aliceEncryptedKeyring, tgID, 1, []byte("32_byte_secret_kek_for_alice_32!"))
	require.NoError(t, err)
	err = keyringSvc.SaveHybrid(aliceEncryptedKeyring, aliceVaultID, alicePassword, "")
	require.NoError(t, err)

	// Ensure Alice admin member CID in Cloud
	_, err = client.AddMemberToTrustGroup(ctx, &trustgroup_domain.AddMemberToTrustGroupRequest{TrustGroupID: tgID, VaultID: aliceVaultID, Role: "admin"})
	require.NoError(t, err)

	// [C3][E2E][02] Alice adds Bob through production application flow
	// DIRECT CALL TO PRODUCTION ENTRY POINT USING ALICE'S REAL AUTH TOKEN FROM AuthHandler.Login!
	updatedTG, err := app.AddTrustGroupMember(aliceRealToken, tgID, bobVaultID, "member")
	require.NoError(t, err, "App.AddTrustGroupMember MUST succeed!")
	require.NotNil(t, updatedTG)
	fmt.Printf("[C3][E2E][02] Alice adds Bob through production application flow memberID=%s\n", bobVaultID)

	// [C3][E2E][03] Bob envelope exists as a RESULT of step 02
	reloadedTGResp, err := client.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{TrustGroupID: tgID})
	require.NoError(t, err)
	reloadedTG := reloadedTGResp.Data
	assert.Contains(t, reloadedTG.MemberCIDs, bobVaultID, "Bob must be persisted as TrustGroup member")

	var bobEnv *trustgroup_domain.TrustGroupKeyEnvelope
	for i := range reloadedTG.KeyEnvelopes {
		if reloadedTG.KeyEnvelopes[i].MemberID == bobVaultID {
			bobEnv = &reloadedTG.KeyEnvelopes[i]
			break
		}
	}
	require.NotNil(t, bobEnv, "Bob's envelope MUST exist as a RESULT of App.AddTrustGroupMember step 02")
	require.NotEmpty(t, bobEnv.WrappedKEK)
	fmt.Printf("[C3][E2E][03] Bob envelope exists as a RESULT of step 02 memberID=%s deviceID=%s kekVersion=%d wrappedKEKLen=%d\n",
		bobEnv.MemberID, bobEnv.DeviceID, bobEnv.KEKVersion, len(bobEnv.WrappedKEK))

	// [C3][E2E][04] Alice creates C3 ShareEntry through production write flow
	assetStore := &memoryAssetStorageResolver{assets: make(map[string][]byte)}

	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(client, cloudShareRepo)
	createShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, addEnvUC).WithCrypto(
		cryptoOrchestrator,
		assetStore,
		identityResolver,
		assetStore,
	)

	createShareRes, err := createShareUC.Execute(ctx, collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tgID,
		KEKVersion:   1,
		CreatedBy:    aliceVaultID,
		AssetCID:     string(originalPayload),
		Metadata: map[string]string{
			"notes": "E2E Acceptance Document",
		},
	})
	require.NoError(t, err, "Production C3 ShareEntry creation MUST succeed!")
	require.NotNil(t, createShareRes)
	shareEntryID := createShareRes.ShareEntry.ID
	require.NotEmpty(t, shareEntryID)

	fmt.Printf("[C3][E2E][04] Alice creates C3 ShareEntry through production write flow id=%s assetCID=%s kekVersion=%d\n",
		shareEntryID, createShareRes.ShareEntry.AssetCID, createShareRes.ShareEntry.KEKVersion)

	// =========================================================================
	// BOB STEPS
	// =========================================================================

	// [C3][E2E][05] Bob resolves ShareEntry through production read flow
	resolveShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(
		cloudShareRepo,
		client,
		assetStore,
		identityResolver,
		cryptoOrchestrator,
	)
	collabHandler := collaboration_ui.NewCollaborationHandler(createShareUC, resolveShareUC, nil)

	bobDeviceID := bobEnv.DeviceID
	fmt.Printf("[C3][E2E][05] Bob resolves ShareEntry through production read flow id=%s deviceID=%s\n", shareEntryID, bobDeviceID)
	resolvedShare, err := collabHandler.ResolveCollaborativeShare(ctx, bobVaultID, shareEntryID, bobDeviceID)
	require.NoError(t, err, "Bob's resolution of collaborative share MUST succeed!")
	require.NotNil(t, resolvedShare)

	// [C3][E2E][06] plaintext matches original payload
	assert.Equal(t, originalPayload, resolvedShare.Plaintext, "Plaintext returned to Bob MUST match Alice's original payload exactly!")
	fmt.Printf("[C3][E2E][06] plaintext matches original payload match=%t\n", true)
}

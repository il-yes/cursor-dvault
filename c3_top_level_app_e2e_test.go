package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	auth_usecases "vault-app/internal/auth/application/use_cases"
	auth_domain "vault-app/internal/auth/domain"
	auth_ui "vault-app/internal/auth/ui"
	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	collaboration_ui "vault-app/internal/collaboration/ui"
	identity_domain "vault-app/internal/identity/domain"
	identity_ui "vault-app/internal/identity/ui"
	"vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_ports "vault-app/internal/trust_group/application/ports"
	trustgroup_envelope_usecases "vault-app/internal/trust_group/application/usecases/envelope"
	trustgroup_member_usecases "vault-app/internal/trust_group/application/usecases/member"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	trustgroup_eventbus "vault-app/internal/trust_group/infrastructure/eventbus"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
	vault_ui "vault-app/internal/vault/ui"
)

// topLevelMockCloud backend simulates authoritative Ankhora Cloud HTTP endpoints
type topLevelMockCloud struct {
	mu          sync.Mutex
	trustGroups map[string]*trustgroup_domain.TrustGroup
	shares      map[string]*c3_asset_domain.ShareEntry
	users       map[string]*tracecore_types.User
}

func newTopLevelMockCloud() *topLevelMockCloud {
	return &topLevelMockCloud{
		trustGroups: make(map[string]*trustgroup_domain.TrustGroup),
		shares:      make(map[string]*c3_asset_domain.ShareEntry),
		users:       make(map[string]*tracecore_types.User),
	}
}

func (m *topLevelMockCloud) Server() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.mu.Lock()
		defer m.mu.Unlock()

		bodyBytes, _ := io.ReadAll(r.Body)

		// 1. Create TrustGroup: POST /api/trustgroups
		if r.Method == http.MethodPost && r.URL.Path == "/api/trustgroups" {
			var tg trustgroup_domain.TrustGroup
			_ = json.Unmarshal(bodyBytes, &tg)
			if tg.ID == "" {
				var reqWrapper struct {
					TrustGroup trustgroup_domain.TrustGroup `json:"trust_group"`
				}
				if err := json.Unmarshal(bodyBytes, &reqWrapper); err == nil && reqWrapper.TrustGroup.ID != "" {
					tg = reqWrapper.TrustGroup
				}
			}
			if tg.ID == "" {
				tg.ID = "tg_" + fmt.Sprintf("%d", time.Now().UnixNano())
			}
			if tg.KEKVersion == 0 {
				tg.KEKVersion = 1
			}
			m.trustGroups[tg.ID] = &tg

			resp := tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{
				Status:  201,
				Success: true,
				Data:    tg,
			}
			w.WriteHeader(201)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		// 2. Get TrustGroup: GET /api/trustgroups/{id}
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/trustgroups/") {
			tgID := strings.TrimPrefix(r.URL.Path, "/api/trustgroups/")
			tg, ok := m.trustGroups[tgID]
			if !ok {
				w.WriteHeader(404)
				return
			}
			resp := tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{
				Status:  200,
				Success: true,
				Data:    *tg,
			}
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		// 3. Update TrustGroup: PUT /api/trustgroups/{id}
		if r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/api/trustgroups/") {
			tgID := strings.TrimPrefix(r.URL.Path, "/api/trustgroups/")
			var updatedTG trustgroup_domain.TrustGroup
			_ = json.Unmarshal(bodyBytes, &updatedTG)
			m.trustGroups[tgID] = &updatedTG
			resp := tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{
				Status:  200,
				Success: true,
				Data:    updatedTG,
			}
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		// 4. Add Member: POST /api/trustgroups/{id}/members
		if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/trustgroups/") && strings.HasSuffix(r.URL.Path, "/members") {
			parts := strings.Split(r.URL.Path, "/")
			tgID := parts[3]
			tg, ok := m.trustGroups[tgID]
			if !ok {
				w.WriteHeader(404)
				return
			}
			var req map[string]string
			_ = json.Unmarshal(bodyBytes, &req)
			vaultID := req["vault_id"]
			alreadyPresent := false
			for _, cid := range tg.MemberCIDs {
				if cid == vaultID {
					alreadyPresent = true
					break
				}
			}
			if !alreadyPresent && vaultID != "" {
				tg.MemberCIDs = append(tg.MemberCIDs, vaultID)
			}
			m.trustGroups[tgID] = tg

			resp := tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{
				Status:  200,
				Success: true,
				Data:    *tg,
			}
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		// 5. Remove Member: DELETE /api/trustgroups/{id}/members/{vaultID}
		if r.Method == http.MethodDelete && strings.Contains(r.URL.Path, "/members/") {
			parts := strings.Split(r.URL.Path, "/")
			if len(parts) >= 6 {
				tgID := parts[3]
				vaultID := parts[5]
				tg, ok := m.trustGroups[tgID]
				if ok {
					var filteredCIDs []string
					for _, cid := range tg.MemberCIDs {
						if cid != vaultID {
							filteredCIDs = append(filteredCIDs, cid)
						}
					}
					tg.MemberCIDs = filteredCIDs
					m.trustGroups[tgID] = tg
				}
			}
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": 200, "success": true})
			return
		}

		// 6. Add Envelope: POST /api/trustgroups/{id}/envelopes
		if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/trustgroups/") && strings.HasSuffix(r.URL.Path, "/envelopes") {
			parts := strings.Split(r.URL.Path, "/")
			tgID := parts[3]
			tg, ok := m.trustGroups[tgID]
			if !ok {
				w.WriteHeader(404)
				return
			}
			var envReq trustgroup_domain.TrustGroupKeyEnvelope
			_ = json.Unmarshal(bodyBytes, &envReq)
			if envReq.ID == "" {
				envReq.ID = "env_" + fmt.Sprintf("%d", time.Now().UnixNano())
			}
			envReq.TrustGroupID = tgID
			envReq.CreatedAt = time.Now()
			tg.KeyEnvelopes = append(tg.KeyEnvelopes, envReq)
			m.trustGroups[tgID] = tg

			resp := tracecore_types.CloudResponse[trustgroup_domain.TrustGroup]{
				Status:  200,
				Success: true,
				Data:    *tg,
			}
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		// 7. Share Entry: POST /api/c3/share-entries & GET /api/c3/share-entries/{id}
		if r.Method == http.MethodPost && (r.URL.Path == "/api/c3/share-entries" || r.URL.Path == "/shares/cryptographic") {
			var se c3_asset_domain.ShareEntry
			_ = json.Unmarshal(bodyBytes, &se)
			if se.ID == "" {
				var reqWrapper struct {
					ShareEntry c3_asset_domain.ShareEntry `json:"share_entry"`
				}
				if err := json.Unmarshal(bodyBytes, &reqWrapper); err == nil && reqWrapper.ShareEntry.ID != "" {
					se = reqWrapper.ShareEntry
				}
			}
			if se.ID == "" {
				se.ID = "se_" + fmt.Sprintf("%d", time.Now().UnixNano())
			}
			m.shares[se.ID] = &se

			resp := tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{
				Status:  201,
				Success: true,
				Data:    se,
			}
			w.WriteHeader(201)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		if r.Method == http.MethodGet && (strings.HasPrefix(r.URL.Path, "/api/c3/share-entries/") || strings.HasPrefix(r.URL.Path, "/shares/cryptographic/")) {
			seID := strings.TrimPrefix(r.URL.Path, "/api/c3/share-entries/")
			seID = strings.TrimPrefix(seID, "/shares/cryptographic/")
			se, ok := m.shares[seID]
			if !ok {
				w.WriteHeader(404)
				return
			}
			resp := tracecore_types.CloudResponse[c3_asset_domain.ShareEntry]{
				Status:  200,
				Success: true,
				Data:    *se,
			}
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		// 8. Customer lookup by email/vaultID
		if r.Method == http.MethodGet && (strings.HasPrefix(r.URL.Path, "/customers") || strings.HasPrefix(r.URL.Path, "/api/customers")) {
			email := r.URL.Query().Get("email")
			vaultID := r.URL.Query().Get("vault_id")
			lookup := email
			if lookup == "" {
				lookup = vaultID
			}

			user, ok := m.users[lookup]
			if !ok {
				w.WriteHeader(404)
				return
			}

			resp := tracecore.GetUserByEmailResponse{
				Error:   false,
				Message: "found",
				Data:    *user,
			}
			w.WriteHeader(200)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(404)
	}))
}

type topLevelIdentityRepo struct {
	users map[string]*identity_domain.User
}

func (m *topLevelIdentityRepo) Save(ctx context.Context, u *identity_domain.User) error {
	m.users[u.ID] = u
	return nil
}

func (m *topLevelIdentityRepo) FindByID(ctx context.Context, id string) (*identity_domain.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	return u, nil
}

func (m *topLevelIdentityRepo) FindByEmail(ctx context.Context, email string) (*identity_domain.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user not found by email: %s", email)
}

func (m *topLevelIdentityRepo) Update(ctx context.Context, u *identity_domain.User) error {
	m.users[u.ID] = u
	return nil
}

func (m *topLevelIdentityRepo) FindByPublicKey(ctx context.Context, publicKey string) (*identity_domain.User, error) {
	for _, u := range m.users {
		if u.StellarPublicKey == publicKey {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user not found by public key: %s", publicKey)
}

type topLevelVaultRepo struct{}

func (r *topLevelVaultRepo) SaveVault(vault *vaults_domain.Vault) error              { return nil }
func (r *topLevelVaultRepo) GetVault(vaultID string) (*vaults_domain.Vault, error)  { return nil, nil }
func (r *topLevelVaultRepo) GetVaultByCID(vaultID string) (*vaults_domain.Vault, error) { return nil, nil }
func (r *topLevelVaultRepo) UpdateVault(vault *vaults_domain.Vault) error            { return nil }
func (r *topLevelVaultRepo) DeleteVault(vaultID string) error                         { return nil }
func (r *topLevelVaultRepo) GetLatestByUserID(userID string) (*vaults_domain.Vault, error) {
	return nil, nil
}
func (r *topLevelVaultRepo) GetByUserIDAndName(userID string, name string) (*vaults_domain.Vault, error) {
	return nil, nil
}
func (r *topLevelVaultRepo) UpdateVaultCID(vaultID, cid string) error { return nil }

type topLevelAssetContentResolver struct {
	assets map[string][]byte
}

func (r *topLevelAssetContentResolver) FetchEncryptedAsset(ctx context.Context, cid string) ([]byte, error) {
	data, ok := r.assets[cid]
	if !ok {
		return nil, fmt.Errorf("asset not found: %s", cid)
	}
	return data, nil
}

type topLevelSovereignIdentityResolver struct {
	seeds    map[string]string
	keyrings map[string]*vaults_domain.VaultKeyring
}

func (r *topLevelSovereignIdentityResolver) GetDeviceSeed(ctx context.Context, userID string) (string, error) {
	seed, ok := r.seeds[userID]
	if !ok {
		return "", trustgroup_domain.ErrDeviceNotFound
	}
	return seed, nil
}

func (r *topLevelSovereignIdentityResolver) GetVaultKeyring(ctx context.Context, userID string) (*vaults_domain.VaultKeyring, error) {
	kr, ok := r.keyrings[userID]
	if !ok {
		return vaults_domain.NewVaultKeyring(userID), nil
	}
	return kr, nil
}

func (r *topLevelSovereignIdentityResolver) GetDevice(ctx context.Context, deviceID string) (*trustgroup_ports.DeviceSummary, error) {
	return nil, trustgroup_domain.ErrDeviceNotFound
}

func (r *topLevelSovereignIdentityResolver) ListActiveDevices(ctx context.Context, memberID string) ([]trustgroup_ports.DeviceSummary, error) {
	return nil, nil
}

// TestC3_TopLevelApp_AddMember_CreateShare_ResolveShare_E2E exercises the complete top-level App workflow:
// App.AddTrustGroupMember → App.CreateCollaborativeShare → App.ResolveCollaborativeShare → assert decrypted plaintext.
func TestC3_TopLevelApp_AddMember_CreateShare_ResolveShare_E2E(t *testing.T) {
	ctx := context.Background()

	// 1. Setup mock Cloud server & tracecore client
	cloud := newTopLevelMockCloud()
	server := cloud.Server()
	defer server.Close()

	tcClient := tracecore.NewTracecoreClient(server.URL+"/api", "test-bearer-token", server.URL, server.URL+"/api")

	// 2. Keypairs for Creator Alice and Member Bob
	kpAlice, err := keypair.Random()
	require.NoError(t, err)

	kpBob, err := keypair.Random()
	require.NoError(t, err)

	aliceVaultID := "vault_alice_uuid_101"
	bobVaultID := "vault_bob_uuid_202"

	// 3. Register Bob on Cloud customer API & Identity repo
	cloud.users[bobVaultID] = &tracecore_types.User{
		ID:        202,
		Email:     "bob@example.com",
		PublicKey: kpBob.Address(),
	}

	identityRepo := &topLevelIdentityRepo{
		users: map[string]*identity_domain.User{
			aliceVaultID: {
				ID:               aliceVaultID,
				Email:            "alice@example.com",
				StellarPublicKey: kpAlice.Address(),
			},
			bobVaultID: {
				ID:               bobVaultID,
				Email:            "bob@example.com",
				StellarPublicKey: kpBob.Address(),
			},
		},
	}
	identityHandler := &identity_ui.IdentityHandler{
		IdentityUserRepo: identityRepo,
	}

	// 4. Setup Keyring Service & Keyrings for Alice & Bob
	tempDir := t.TempDir()
	aesSvc := &vault_infrastructure_crypto.AESService{}
	keyEnc := vault_infrastructure_crypto.NewKeyService()
	fs := &vault_infrastructure_security.OSFileSystem{}
	keyringSvc := vault_infrastructure_security.NewKeyringService(aesSvc, keyEnc, tempDir, fs)

	aliceKeyring := vaults_domain.NewVaultKeyring(aliceVaultID)
	aliceKeyring.Keys = append(aliceKeyring.Keys, vaults_domain.EncryptedKey{
		ID:        "k_sec_alice",
		Type:      vaults_domain.KeyTypeDeviceSeed,
		Version:   1,
		Ciphertext: []byte(kpAlice.Seed()),
		CreatedAt: time.Now().Unix(),
	})
	_ = keyringSvc.SaveHybrid(aliceKeyring, aliceVaultID, "pass_alice", kpAlice.Seed())

	bobKeyring := vaults_domain.NewVaultKeyring(bobVaultID)
	bobKeyring.Keys = append(bobKeyring.Keys, vaults_domain.EncryptedKey{
		ID:        "k_sec_bob",
		Type:      vaults_domain.KeyTypeDeviceSeed,
		Version:   1,
		Ciphertext: []byte(kpBob.Seed()),
		CreatedAt: time.Now().Unix(),
	})
	_ = keyringSvc.SaveHybrid(bobKeyring, bobVaultID, "pass_bob", kpBob.Seed())

	// 5. Setup Vault Handler & Sessions
	sessions := map[string]*vault_session.Session{
		aliceVaultID: {
			UserID:  aliceVaultID,
			Runtime: &vault_session.RuntimeContext{VaultID: aliceVaultID, SessionSecrets: map[string]string{"password": "pass_alice", "stellar_secret": kpAlice.Seed()}},
		},
		bobVaultID: {
			UserID:  bobVaultID,
			Runtime: &vault_session.RuntimeContext{VaultID: bobVaultID, SessionSecrets: map[string]string{"password": "pass_bob", "stellar_secret": kpBob.Seed()}},
		},
	}
	sessionMgr := vault_session.NewManager(nil, nil, nil, nil, nil, sessions)
	vaultHandler := &vault_ui.VaultHandler{
		KeyringService:  keyringSvc,
		SessionManager:  sessionMgr,
		VaultRepository: &topLevelVaultRepo{},
	}

	// 6. Setup Auth Token Service & Handler
	authCfg := auth_domain.Auth{
		Issuer:      "ankhora-test",
		Audience:    "ankhora-desktop",
		Secret:      "test-secret-key-32-bytes-length!!",
		TokenExpiry: 1 * time.Hour,
	}
	tokenSvc := auth_usecases.NewTokenService(authCfg, nil, nil)
	tokenUC := auth_usecases.NewGenerateTokensUseCase(nil, tokenSvc)
	authHandler := auth_ui.NewAuthHandler(nil, tokenUC, nil)

	tokenPairAlice, err := tokenSvc.GenerateTokenPair(&auth_domain.JwtUser{ID: aliceVaultID, Email: "alice@example.com"})
	require.NoError(t, err)
	tokenAlice := tokenPairAlice.Token

	tokenPairBob, err := tokenSvc.GenerateTokenPair(&auth_domain.JwtUser{ID: bobVaultID, Email: "bob@example.com"})
	require.NoError(t, err)
	tokenBob := tokenPairBob.Token

	// 7. Setup Crypto Orchestrator & Use Cases
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	cryptoOrchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	eventBus := trustgroup_eventbus.NewMemoryBus()
	addMemberUC := trustgroup_member_usecases.NewAddMemberToTrustGroupUsecase(tcClient, eventBus)
	addEnvelopeUC := trustgroup_envelope_usecases.NewAddTrustGroupKeyEnvelopeUseCase(tcClient, nil)
	provisionEnvelopeUC := trustgroup_envelope_usecases.NewProvisionTrustGroupMemberEnvelopeUseCase(
		tcClient,
		cryptoOrchestrator,
		addEnvelopeUC,
		keyringSvc,
	)

	cloudShareRepo := tracecore.NewCloudShareEntryRepository(tcClient)
	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(tcClient, cloudShareRepo)
	createShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, addEnvelopeUC)

	assetContentResolver := &topLevelAssetContentResolver{assets: make(map[string][]byte)}
	sovereignIdentityResolver := &topLevelSovereignIdentityResolver{
		seeds: map[string]string{
			aliceVaultID: kpAlice.Seed(),
			bobVaultID:   kpBob.Seed(),
		},
		keyrings: map[string]*vaults_domain.VaultKeyring{
			aliceVaultID: aliceKeyring,
			bobVaultID:   bobKeyring,
		},
	}
	resolveShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(
		cloudShareRepo,
		tcClient,
		assetContentResolver,
		sovereignIdentityResolver,
		cryptoOrchestrator,
	)
	collabHandler := collaboration_ui.NewCollaborationHandler(createShareUC, resolveShareUC, nil)

	// 8. Construct App Instance
	app := &App{
		ctx:                   ctx,
		AuthHandler:           authHandler,
		Identity:              identityHandler,
		Vault:                 vaultHandler,
		CollaborationHandler:  collabHandler,
		addTrustGroupMemberUC: addMemberUC,
		provisionEnvelopeUC:   provisionEnvelopeUC,
		tracecoreClient:       tcClient,
	}

	// -------------------------------------------------------------------------
	// Step A: Create initial TrustGroup aggregate for Alice
	// -------------------------------------------------------------------------
	tg := trustgroup_domain.NewTrustGroup("ch_e2e_channel", "E2E Production Group", []string{aliceVaultID})
	tg.ID = "tg-e2e-app-001"
	createTgResp, err := tcClient.CreateTrustGroup(ctx, &trustgroup_domain.CreateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)
	tgID := createTgResp.Data.ID

	// Prepare initial asset & KEK for Alice
	rawSecretPayload := []byte("TOP SECRET E2E APP PAYLOAD 2026")
	preparedAsset, err := cryptoOrchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_e2e_01",
		TrustGroupID: tgID,
		KEKVersion:   1,
		RawPayload:   rawSecretPayload,
		Keyring:      aliceKeyring,
	})
	require.NoError(t, err)

	// Save updated Alice keyring containing KEK_1 to disk so App.AddTrustGroupMember uses KEK_1
	err = keyringSvc.SaveHybrid(aliceKeyring, aliceVaultID, "pass_alice", kpAlice.Seed())
	require.NoError(t, err)

	assetCID := "bafybeie2eapppayload001"
	assetContentResolver.assets[assetCID] = preparedAsset.EncryptedData

	// Save Alice's envelope to Cloud
	rawKEK, err := keyringSvc.GetTrustGroupKEK(aliceKeyring, tgID, 1)
	require.NoError(t, err)
	alicePayload, err := aesSvc.EncryptPayload(kpAlice.Address(), rawKEK)
	require.NoError(t, err)
	aliceEnv := trustgroup_domain.TrustGroupKeyEnvelope{
		ID:           "env_alice_01",
		TrustGroupID: tgID,
		MemberID:     aliceVaultID,
		WrappedKEK:   alicePayload.ToString(),
		KEKVersion:   1,
	}
	tg.KeyEnvelopes = append(tg.KeyEnvelopes, aliceEnv)
	_, err = tcClient.UpdateTrustGroup(ctx, &trustgroup_domain.UpdateTrustGroupRequest{TrustGroup: *tg})
	require.NoError(t, err)

	// -------------------------------------------------------------------------
	// Step B: Top-Level App Call 1: App.AddTrustGroupMember
	// Exercises: target public-key resolution, caller keyring load, atomic member + envelope provisioning
	// -------------------------------------------------------------------------
	updatedTg, err := app.AddTrustGroupMember(tokenAlice, tgID, bobVaultID, "member")
	require.NoError(t, err, "App.AddTrustGroupMember top-level call must succeed")
	require.NotNil(t, updatedTg)
	require.Contains(t, updatedTg.MemberCIDs, bobVaultID, "Member B must be present in MemberCIDs")

	// Verify that App.AddTrustGroupMember produced an envelope for Bob with KEKVersion=1
	var bobEnvelope *trustgroup_domain.TrustGroupKeyEnvelope
	for i := range updatedTg.KeyEnvelopes {
		if updatedTg.KeyEnvelopes[i].MemberID == bobVaultID && updatedTg.KeyEnvelopes[i].KEKVersion == 1 {
			bobEnvelope = &updatedTg.KeyEnvelopes[i]
			break
		}
	}
	require.NotNil(t, bobEnvelope, "App.AddTrustGroupMember MUST produce a valid key envelope for Member B with KEKVersion=1")
	require.NotEmpty(t, bobEnvelope.WrappedKEK, "Produced envelope WrappedKEK must not be empty")

	// -------------------------------------------------------------------------
	// Step C: Top-Level App Call 2: App.CreateCollaborativeShare
	// Exercises: production ShareEntry creation using produced KEK version
	// -------------------------------------------------------------------------
	wrappedDEKB64 := base64.StdEncoding.EncodeToString(preparedAsset.WrappedDEK)
	createShareRes, err := app.CreateCollaborativeShare(
		tokenAlice,
		"th_thread_e2e",
		tgID,
		assetCID,
		bobVaultID,
		"E2E Production Share Notes",
		wrappedDEKB64,
		1,
	)
	require.NoError(t, err, "App.CreateCollaborativeShare top-level call must succeed")
	require.NotNil(t, createShareRes)
	require.NotEmpty(t, createShareRes.ShareEntryID)

	shareEntryID := createShareRes.ShareEntryID

	// -------------------------------------------------------------------------
	// Step D: Top-Level App Call 3: App.ResolveCollaborativeShare for Member B
	// Exercises: fetching ShareEntry, resolving TrustGroup, retrieving Bob's envelope (produced in Step B), unwrapping KEK + DEK, decrypting payload
	// -------------------------------------------------------------------------
	resolvedShare, err := app.ResolveCollaborativeShare(tokenBob, shareEntryID)
	require.NoError(t, err, "App.ResolveCollaborativeShare top-level call for Member B MUST succeed")
	require.NotNil(t, resolvedShare)

	// -------------------------------------------------------------------------
	// Step E: Assert Decrypted Plaintext Matches Original Secret
	// PROOF: The envelope produced by App.AddTrustGroupMember is the exact envelope consumed by App.ResolveCollaborativeShare
	// -------------------------------------------------------------------------
	assert.Equal(t, string(rawSecretPayload), string(resolvedShare.Plaintext),
		"Decrypted plaintext from App.ResolveCollaborativeShare MUST match original payload exactly")
}

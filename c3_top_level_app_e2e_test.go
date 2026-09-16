package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	identity_domain "vault-app/internal/identity/domain"
	"vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_ports "vault-app/internal/trust_group/application/ports"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vaults_domain "vault-app/internal/vault/domain"
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
	return &vaults_domain.Vault{ID: "v_" + userID, UserID: userID, Name: "Default Vault"}, nil
}
func (r *topLevelVaultRepo) GetByUserIDAndName(userID string, name string) (*vaults_domain.Vault, error) {
	return &vaults_domain.Vault{ID: "v_" + userID, UserID: userID, Name: name}, nil
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

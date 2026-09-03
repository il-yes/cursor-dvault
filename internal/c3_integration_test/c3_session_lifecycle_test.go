package c3_integration_test

import (
	"context"
	"encoding/base64"
	"errors"
	"sync"
	"testing"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vault-app/internal/blockchain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	collaboration_ui "vault-app/internal/collaboration/ui"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
	vaults_persistence "vault-app/internal/vault/infrastructure/persistence"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

type mockIpfsStorage struct {
	mu     sync.RWMutex
	assets map[string][]byte
	repo   *roundTripRepo
}

func (m *mockIpfsStorage) Get(ctx context.Context, cid string) ([]byte, error) {
	if m.repo != nil {
		return m.repo.FetchEncryptedAsset(ctx, cid)
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.assets == nil {
		return nil, errors.New("cid not found")
	}
	data, ok := m.assets[cid]
	if !ok {
		return nil, errors.New("cid not found")
	}
	return data, nil
}

// ---------------------------------------------------------------------------
// Phase 4B Test: Cryptographic Isolation of Collaborative Shares
// Proves a legitimate collaborative recipient resolves a shared asset via TrustGroup
// even when the recipient's private vault session is LOCKED / Session.VaultKey == nil.
// ---------------------------------------------------------------------------
func TestC3_CollaborativeShare_CryptographicIsolation(t *testing.T) {
	ctx := context.Background()
	repo := newRoundTripRepo()

	aesSvc := &vault_infrastructure_crypto.AESService{}
	cryptoSvc := &blockchain.CryptoService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	// 1. Setup Alice's Sovereign Device Keypair & TrustGroup
	kpAlice, err := keypair.Random()
	require.NoError(t, err)

	userAliceID := "user_alice_collaborator"
	deviceLaptopID := "device_laptop_alice"

	repo.seeds[userAliceID] = kpAlice.Seed()
	repo.keyrings[userAliceID] = &vaults_domain.VaultKeyring{UserID: userAliceID, VaultID: "vault_alice"}

	tg := trustgroup_domain.NewTrustGroup("tg_collab_2026", "Legal Review Group", []string{userAliceID})
	tg.KEKVersion = 1
	repo.trustGroups[tg.ID] = *tg

	// 2. Prepare & Encrypt Collaborative Asset
	rawOriginalContent := []byte(`{"document":"M&A Agreement Draft 2026","status":"CONFIDENTIAL"}`)
	assetCID := "bafybeicollabasset2026cid"

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset_ma_agreement",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawOriginalContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: deviceLaptopID, MemberID: userAliceID, PublicKey: kpAlice.Address(), IsActive: true},
		},
	}

	prepared, err := orchestrator.PrepareCollaborativeAsset(ctx, prepPayload)
	require.NoError(t, err)
	require.Len(t, prepared.Envelopes, 1)

	repo.assets[assetCID] = prepared.EncryptedData

	envReq := prepared.Envelopes[0]
	err = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
		TrustGroupID: envReq.TrustGroupID,
		MemberID:     envReq.MemberID,
		DeviceID:     envReq.DeviceID,
		KEKVersion:   envReq.KEKVersion,
		WrappedKEK:   envReq.WrappedKEK,
	})
	require.NoError(t, err)
	repo.trustGroups[tg.ID] = *tg

	// Wire Use Cases & Handler
	shareAssetUC := collaboration_usecases.NewShareAssetWithTrustGroupUsecase(repo, repo)
	createCollabShareUC := collaboration_usecases.NewCreateCollaborativeShareUseCase(shareAssetUC, nil)
	resolveCollabShareUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(repo, repo, repo, repo, orchestrator)

	collabHandler := collaboration_ui.NewCollaborationHandler(createCollabShareUC, resolveCollabShareUC, nil)

	// Create Share Entry
	wrappedDEKStr := base64.StdEncoding.EncodeToString(prepared.WrappedDEK)
	shareReq := collaboration_dtos.CreateCollaborativeShareRequest{
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		CreatedBy:    userAliceID,
		AssetCID:     assetCID,
		WrappedDEK:   wrappedDEKStr,
		Metadata:     map[string]string{"classification": "restricted"},
	}

	createResp, err := createCollabShareUC.Execute(ctx, shareReq)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	shareEntryID := createResp.ShareEntry.ID
	repo.shareEntries[shareEntryID] = createResp.ShareEntry

	// --- STEP 3: EXPLICITLY SIMULATE ALICE'S PRIVATE VAULT BEING LOCKED ---
	// Alice's private vault session is destroyed / locked / Session.VaultKey == nil
	alicePrivateSession := &vault_session.Session{
		UserID:   userAliceID,
		VaultKey: nil, // ⚠️ LOCKED: NO PRIVATE VAULT KEY AVAILABLE!
	}
	sessionRepo := vaults_persistence.NewGormSessionRepository(nil)
	sessionMgr := vault_session.NewManager(sessionRepo, nil, nil, nil, nil, map[string]*vault_session.Session{})
	sessionMgr.CloseSession(userAliceID) // Ensure Alice has zero active private session in memory

	// Verify Alice's private vault session is indeed non-existent/locked
	_, errSession := sessionMgr.GetSession(userAliceID)
	require.Error(t, errSession, "Alice's private vault session MUST be locked/non-existent")
	require.Nil(t, alicePrivateSession.VaultKey)

	// --- STEP 4: DECISIVE TEST INVARIANT — RESOLVE COLLABORATIVE SHARE WHILE LOCKED ---
	// Resolve collaborative share for Alice when her private vault session is locked
	resolvedDTO, errResolve := collabHandler.ResolveCollaborativeShare(ctx, userAliceID, shareEntryID, deviceLaptopID)
	require.NoError(t, errResolve, "DECISIVE INVARIANT: Collaborative share resolution MUST succeed even when recipient's private vault session is LOCKED!")
	require.NotNil(t, resolvedDTO)

	// Assert decrypted plaintext matches original payload exactly
	assert.Equal(t, rawOriginalContent, resolvedDTO.Plaintext, "DECISIVE INVARIANT: Collaborative asset plaintext MUST match original payload")
	assert.Equal(t, shareEntryID, resolvedDTO.ShareEntryID)
	assert.Equal(t, tg.ID, resolvedDTO.TrustGroupID)

	// --- STEP 5: PROVE PRIVATE VAULT OPERATIONS FOR ALICE AT THIS SAME MOMENT STILL FAIL ---
	// Attempting a private vault query without Session.VaultKey MUST fail
	privateQuery := vault_queries.GetIPFSDataQuerry{
		CID:      assetCID,
		Password: "",
		VaultKey: alicePrivateSession.VaultKey, // nil!
	}
	mockStorage := &mockIpfsStorage{repo: repo}
	queryHandlerNoUnlock := vault_queries.NewGetIPFSDataQuerryHandler(cryptoSvc, aesSvc, nil, nil)
	queryHandlerNoUnlock.SetIpfsService(mockStorage)
	_, errPrivateRead := queryHandlerNoUnlock.Execute(ctx, privateQuery)
	assert.Error(t, errPrivateRead, "Private vault operation for Alice MUST fail while locked")
	assert.Contains(t, errPrivateRead.Error(), "UnlockVaultHandler missing and no session VaultKey provided")
}

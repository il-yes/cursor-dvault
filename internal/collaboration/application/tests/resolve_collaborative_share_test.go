package collaboration_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

// ---------------------------------------------------------------------------
// Mock Application Ports
// ---------------------------------------------------------------------------
type mockAssetResolver struct {
	assets       map[string][]byte
	fetchCount   int
	failOnCID    string
}

func (m *mockAssetResolver) FetchEncryptedAsset(_ context.Context, cid string) ([]byte, error) {
	m.fetchCount++
	if m.failOnCID != "" && cid == m.failOnCID {
		return nil, errors.New("storage asset fetch failed")
	}
	data, ok := m.assets[cid]
	if !ok {
		return nil, errors.New("asset CID not found in storage")
	}
	return data, nil
}

type mockIdentityResolver struct {
	seeds       map[string]string
	keyrings    map[string]*vaults_domain.VaultKeyring
	seedCount   int
	failForUser string
}

func (m *mockIdentityResolver) GetDeviceSeed(_ context.Context, userID string) (string, error) {
	m.seedCount++
	if m.failForUser != "" && userID == m.failForUser {
		return "", errors.New("user device seed not found")
	}
	seed, ok := m.seeds[userID]
	if !ok {
		return "", errors.New("seed not found for user")
	}
	return seed, nil
}

func (m *mockIdentityResolver) GetVaultKeyring(_ context.Context, userID string) (*vaults_domain.VaultKeyring, error) {
	kr, ok := m.keyrings[userID]
	if !ok {
		return &vaults_domain.VaultKeyring{UserID: userID, VaultID: "vault_default"}, nil
	}
	return kr, nil
}

// ---------------------------------------------------------------------------
// Fixture Setup
// ---------------------------------------------------------------------------
type resolveTestFixture struct {
	orchestrator     *trustgroup_orchestrator.TrustGroupCryptoOrchestrator
	shareRepo        *fakeShareEntryRepo
	tgRepo           *fakeTrustGroupRepo
	assetResolver    *mockAssetResolver
	identityResolver *mockIdentityResolver
	useCase          *collaboration_usecases.ResolveCollaborativeShareUseCase
	kp               *keypair.Full
	trustGroup       *trustgroup_domain.TrustGroup
	shareEntry       c3_asset_domain.ShareEntry
	encryptedData    []byte
	rawContent       []byte
}

func setupResolveTestFixture(t *testing.T) *resolveTestFixture {
	ctx := context.Background()
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, aesSvc, asymSvc)

	kp, err := keypair.Random()
	require.NoError(t, err)

	rawContent := []byte("TOP SECRET FINANCIAL CONTRACT CONTENT")
	tg := trustgroup_domain.NewTrustGroup("ch_1", "Finance Council", []string{"user_alice"})
	tg.KEKVersion = 1

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset-100",
		TrustGroupID: tg.ID,
		KEKVersion:   1,
		RawPayload:   rawContent,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: "dev_laptop", MemberID: "user_alice", PublicKey: kp.Address(), IsActive: true},
		},
	}

	prepared, err := orchestrator.PrepareCollaborativeAsset(ctx, prepPayload)
	require.NoError(t, err)
	require.Len(t, prepared.Envelopes, 1)

	// Add key envelope to TrustGroup
	envReq := prepared.Envelopes[0]
	err = tg.AddEnvelope(trustgroup_domain.TrustGroupKeyEnvelope{
		TrustGroupID: envReq.TrustGroupID,
		MemberID:     envReq.MemberID,
		KEKVersion:   envReq.KEKVersion,
		WrappedKEK:   envReq.WrappedKEK,
	})
	require.NoError(t, err)

	wrappedDEKBase64 := base64.StdEncoding.EncodeToString(prepared.WrappedDEK)

	shareEntry, err := c3_asset_domain.NewShareEntry(
		"cid_blueprint_100",
		tg.ID,
		wrappedDEKBase64,
		1,
		"user_owner",
		map[string]string{"title": "Q3 Audit"},
	)
	require.NoError(t, err)
	shareEntry.ID = "se_100"

	shareRepo := newFakeShareEntryRepo()
	_, err = shareRepo.CreateShareEntry(ctx, &c3_asset_domain.CreateShareEntryRequest{ShareEntry: shareEntry})
	require.NoError(t, err)

	tgRepo := newFakeTrustGroupRepo()
	tgRepo.groups[tg.ID] = tg

	assetResolver := &mockAssetResolver{
		assets: map[string][]byte{
			"cid_blueprint_100": prepared.EncryptedData,
		},
	}

	identityResolver := &mockIdentityResolver{
		seeds: map[string]string{
			"user_alice": kp.Seed(),
		},
		keyrings: map[string]*vaults_domain.VaultKeyring{
			"user_alice": {UserID: "user_alice", VaultID: "vault_alice"},
		},
	}

	useCase := collaboration_usecases.NewResolveCollaborativeShareUseCase(shareRepo, tgRepo, assetResolver, identityResolver, orchestrator)

	return &resolveTestFixture{
		orchestrator:     orchestrator,
		shareRepo:        shareRepo,
		tgRepo:           tgRepo,
		assetResolver:    assetResolver,
		identityResolver: identityResolver,
		useCase:          useCase,
		kp:               kp,
		trustGroup:       tg,
		shareEntry:       shareEntry,
		encryptedData:    prepared.EncryptedData,
		rawContent:       rawContent,
	}
}

// ---------------------------------------------------------------------------
// Unit & Security Test Cases
// ---------------------------------------------------------------------------

// 1. Success
func TestResolveCollaborativeShare_Success(t *testing.T) {
	f := setupResolveTestFixture(t)

	res, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: f.shareEntry.ID,
		CallerUserID: "user_alice",
	})

	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, f.shareEntry.ID, res.ShareEntryID)
	assert.Equal(t, f.trustGroup.ID, res.TrustGroupID)
	assert.Equal(t, "user_owner", res.CreatedBy)
	assert.Equal(t, f.rawContent, res.Plaintext)
	assert.Equal(t, 1, f.assetResolver.fetchCount, "AssetContentResolver MUST be called exactly once on success")
	assert.Equal(t, 1, f.identityResolver.seedCount, "SovereignIdentityResolver MUST be called exactly once on success")
}

// 2. Missing ShareEntry
func TestResolveCollaborativeShare_MissingShareEntry(t *testing.T) {
	f := setupResolveTestFixture(t)

	_, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: "se_nonexistent",
		CallerUserID: "user_alice",
	})

	assert.ErrorIs(t, err, collaboration_usecases.ErrShareEntryNotFound)
	assert.Equal(t, 0, f.assetResolver.fetchCount, "MUST NOT fetch asset if ShareEntry missing")
	assert.Equal(t, 0, f.identityResolver.seedCount, "MUST NOT resolve identity if ShareEntry missing")
}

// 3. Missing TrustGroup
func TestResolveCollaborativeShare_MissingTrustGroup(t *testing.T) {
	f := setupResolveTestFixture(t)
	delete(f.tgRepo.groups, f.trustGroup.ID)

	_, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: f.shareEntry.ID,
		CallerUserID: "user_alice",
	})

	assert.ErrorIs(t, err, collaboration_usecases.ErrTrustGroupNotFound)
	assert.Equal(t, 0, f.assetResolver.fetchCount, "MUST NOT fetch asset if TrustGroup missing")
	assert.Equal(t, 0, f.identityResolver.seedCount, "MUST NOT resolve identity if TrustGroup missing")
}

// 4. Unauthorized Member — Ensures Stop Before Storage & Crypto Resolution
func TestResolveCollaborativeShare_UnauthorizedMember(t *testing.T) {
	f := setupResolveTestFixture(t)

	_, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: f.shareEntry.ID,
		CallerUserID: "user_eve", // Eve is not in MemberCIDs
	})

	assert.ErrorIs(t, err, collaboration_usecases.ErrUnauthorizedMember)
	assert.Equal(t, 0, f.assetResolver.fetchCount, "SECURITY INVARIANT: Unauthorized member MUST NOT trigger asset storage retrieval")
	assert.Equal(t, 0, f.identityResolver.seedCount, "SECURITY INVARIANT: Unauthorized member MUST NOT trigger identity/seed resolution")
}

// 5. Revoked Member — Ensures Stop Before Storage & Crypto Resolution
func TestResolveCollaborativeShare_RevokedMember(t *testing.T) {
	f := setupResolveTestFixture(t)
	f.trustGroup.MemberCIDs = []string{}
	f.tgRepo.groups[f.trustGroup.ID] = f.trustGroup

	_, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: f.shareEntry.ID,
		CallerUserID: "user_alice",
	})

	assert.ErrorIs(t, err, collaboration_usecases.ErrUnauthorizedMember)
	assert.Equal(t, 0, f.assetResolver.fetchCount, "SECURITY INVARIANT: Revoked member MUST NOT trigger asset storage retrieval")
	assert.Equal(t, 0, f.identityResolver.seedCount, "SECURITY INVARIANT: Revoked member MUST NOT trigger identity/seed resolution")
}

// 6. Missing Device Envelope — Ensures Stop Before Storage & Crypto Resolution
func TestResolveCollaborativeShare_MissingDeviceEnvelope(t *testing.T) {
	f := setupResolveTestFixture(t)
	f.trustGroup.KeyEnvelopes = nil
	f.tgRepo.groups[f.trustGroup.ID] = f.trustGroup

	_, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: f.shareEntry.ID,
		CallerUserID: "user_alice",
	})

	assert.ErrorIs(t, err, collaboration_usecases.ErrKeyEnvelopeNotFound)
	assert.Equal(t, 0, f.assetResolver.fetchCount, "SECURITY INVARIANT: Missing device envelope MUST NOT trigger asset storage retrieval")
	assert.Equal(t, 0, f.identityResolver.seedCount, "SECURITY INVARIANT: Missing device envelope MUST NOT trigger identity/seed resolution")
}

// 7. Revoked Device Envelope — Ensures Stop Before Storage & Crypto Resolution
func TestResolveCollaborativeShare_RevokedDeviceEnvelope(t *testing.T) {
	f := setupResolveTestFixture(t)
	now := f.trustGroup.KeyEnvelopes[0].CreatedAt
	f.trustGroup.KeyEnvelopes[0].RevokedAt = &now
	f.tgRepo.groups[f.trustGroup.ID] = f.trustGroup

	_, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: f.shareEntry.ID,
		CallerUserID: "user_alice",
	})

	assert.ErrorIs(t, err, collaboration_usecases.ErrKeyEnvelopeNotFound)
	assert.Equal(t, 0, f.assetResolver.fetchCount, "SECURITY INVARIANT: Revoked device envelope MUST NOT trigger asset storage retrieval")
	assert.Equal(t, 0, f.identityResolver.seedCount, "SECURITY INVARIANT: Revoked device envelope MUST NOT trigger identity/seed resolution")
}

// 8. KEK Version Mismatch — Ensures Stop Before Storage & Crypto Resolution
func TestResolveCollaborativeShare_KEKVersionMismatch(t *testing.T) {
	f := setupResolveTestFixture(t)
	f.trustGroup.KeyEnvelopes[0].KEKVersion = 2
	f.tgRepo.groups[f.trustGroup.ID] = f.trustGroup

	_, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: f.shareEntry.ID,
		CallerUserID: "user_alice",
	})

	assert.ErrorIs(t, err, collaboration_usecases.ErrKeyEnvelopeNotFound)
	assert.Equal(t, 0, f.assetResolver.fetchCount, "SECURITY INVARIANT: KEK version mismatch MUST NOT trigger asset storage retrieval")
	assert.Equal(t, 0, f.identityResolver.seedCount, "SECURITY INVARIANT: KEK version mismatch MUST NOT trigger identity/seed resolution")
}

// 9. Asset Storage Retrieval Failure
func TestResolveCollaborativeShare_StorageFetchFailure(t *testing.T) {
	f := setupResolveTestFixture(t)
	f.assetResolver.failOnCID = "cid_blueprint_100"

	_, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: f.shareEntry.ID,
		CallerUserID: "user_alice",
	})

	assert.ErrorContains(t, err, "storage asset fetch failed")
	assert.Equal(t, 1, f.assetResolver.fetchCount)
}

// 10. Crypto Resolution Failure (Wrong Device Seed)
func TestResolveCollaborativeShare_CryptoFailure(t *testing.T) {
	f := setupResolveTestFixture(t)
	wrongKp, _ := keypair.Random()
	f.identityResolver.seeds["user_alice"] = wrongKp.Seed()

	_, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: f.shareEntry.ID,
		CallerUserID: "user_alice",
	})

	assert.ErrorContains(t, err, "cryptographic resolution failed")
}

// 11. Zero Secret Leakage Security Assertion
func TestResolveCollaborativeShare_ZeroSecretLeakage(t *testing.T) {
	f := setupResolveTestFixture(t)

	res, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID: f.shareEntry.ID,
		CallerUserID: "user_alice",
	})

	require.NoError(t, err)
	marshaledBytes, err := json.Marshal(res)
	require.NoError(t, err)
	marshaledStr := string(marshaledBytes)

	forbiddenSecrets := []string{
		"wrapped_dek",
		"wrappedDEK",
		"wrapped_kek",
		"wrappedKEK",
		"key_envelopes",
		"PrivateKey",
		"SecretKey",
		"private_key",
		"secret_key",
		"device_seed",
		"DeviceSeed",
	}

	for _, forbidden := range forbiddenSecrets {
		if strings.Contains(marshaledStr, forbidden) {
			t.Errorf("SECURITY VIOLATION: Serialized ResolveCollaborativeShareResponse DTO contains forbidden secret: %q", forbidden)
		}
	}
}

// ---------------------------------------------------------------------------
// Explicit Identity Boundary Tests (IdentityID vs VaultID)
// ---------------------------------------------------------------------------

// TEST 1: IdentityID != VaultID, MemberCIDs contains VaultID, Caller resolves through session to VaultID -> AUTHORIZED
func TestResolveCollaborativeShare_IdentityBoundary_ResolvedVaultID_Authorized(t *testing.T) {
	f := setupResolveTestFixture(t)
	identityID := "user_identity_bob_100"
	vaultID := "55a0ced5-b246-477a-a3a7-abcc26b78ed8"
	deviceID := "dev_laptop_bob"

	// TrustGroup MemberCIDs contains VaultID, not IdentityID
	f.trustGroup.MemberCIDs = []string{vaultID}
	f.trustGroup.KeyEnvelopes[0].MemberID = vaultID
	f.trustGroup.KeyEnvelopes[0].DeviceID = deviceID
	f.tgRepo.groups[f.trustGroup.ID] = f.trustGroup
	f.identityResolver.seeds[identityID] = f.kp.Seed()
	f.identityResolver.seeds[vaultID] = f.kp.Seed()

	// Request has IdentityID != VaultID, with CallerVaultID properly resolved
	res, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:     f.shareEntry.ID,
		CallerIdentityID: identityID,
		CallerVaultID:    vaultID,
	})

	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, f.rawContent, res.Plaintext)
}

// TEST 2: IdentityID != VaultID, Caller incorrectly supplied as IdentityID -> ErrUnauthorizedMember
func TestResolveCollaborativeShare_IdentityBoundary_RawIdentityID_Unauthorized(t *testing.T) {
	f := setupResolveTestFixture(t)
	identityID := "user_identity_bob_100"
	vaultID := "55a0ced5-b246-477a-a3a7-abcc26b78ed8"

	// TrustGroup MemberCIDs contains VaultID
	f.trustGroup.MemberCIDs = []string{vaultID}
	f.tgRepo.groups[f.trustGroup.ID] = f.trustGroup

	// Request incorrectly supplies IdentityID as CallerVaultID
	_, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:     f.shareEntry.ID,
		CallerIdentityID: identityID,
		CallerVaultID:    identityID, // Incorrect: IdentityID passed into VaultID position
	})

	assert.ErrorIs(t, err, collaboration_usecases.ErrUnauthorizedMember)
	assert.Equal(t, 0, f.assetResolver.fetchCount, "SECURITY INVARIANT: Unauthorized member MUST NOT trigger asset storage retrieval")
}

// TEST 3: Caller VaultID is not a TrustGroup member -> ErrUnauthorizedMember
func TestResolveCollaborativeShare_IdentityBoundary_NonMemberVaultID_Unauthorized(t *testing.T) {
	f := setupResolveTestFixture(t)
	nonMemberVaultID := "vault_unauthorized_999"

	_, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:  f.shareEntry.ID,
		CallerVaultID: nonMemberVaultID,
	})

	assert.ErrorIs(t, err, collaboration_usecases.ErrUnauthorizedMember)
	assert.Equal(t, 0, f.assetResolver.fetchCount, "SECURITY INVARIANT: Non-member VaultID MUST NOT trigger asset storage retrieval")
}

// TEST 4: Correct member VaultID but wrong TrustGroup -> ErrUnauthorizedMember
func TestResolveCollaborativeShare_IdentityBoundary_WrongTrustGroup_Unauthorized(t *testing.T) {
	f := setupResolveTestFixture(t)
	memberVaultID := "55a0ced5-b246-477a-a3a7-abcc26b78ed8"

	// Create another TrustGroup that does NOT contain memberVaultID
	otherTG := trustgroup_domain.NewTrustGroup("ch_other", "Other Group", []string{"vault_other_only"})
	otherTG.KEKVersion = 1
	f.tgRepo.groups[otherTG.ID] = otherTG

	// Modify share entry to point to other TrustGroup
	wrongShareEntry := f.shareEntry
	wrongShareEntry.ID = "se_wrong_tg"
	wrongShareEntry.TrustGroupID = otherTG.ID
	f.shareRepo.entries[wrongShareEntry.ID] = wrongShareEntry

	_, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:  wrongShareEntry.ID,
		CallerVaultID: memberVaultID,
	})

	assert.ErrorIs(t, err, collaboration_usecases.ErrUnauthorizedMember)
	assert.Equal(t, 0, f.assetResolver.fetchCount, "SECURITY INVARIANT: Wrong TrustGroup MUST NOT trigger asset storage retrieval")
}

// TEST 5: Correct member + correct TrustGroup + valid envelope -> continue past membership authorization into envelope/decryption resolution
func TestResolveCollaborativeShare_IdentityBoundary_FullValidResolution(t *testing.T) {
	f := setupResolveTestFixture(t)
	memberVaultID := "55a0ced5-b246-477a-a3a7-abcc26b78ed8"
	deviceID := "dev_laptop_valid"

	f.trustGroup.MemberCIDs = []string{memberVaultID}
	f.trustGroup.KeyEnvelopes[0].MemberID = memberVaultID
	f.trustGroup.KeyEnvelopes[0].DeviceID = deviceID
	f.tgRepo.groups[f.trustGroup.ID] = f.trustGroup
	f.identityResolver.seeds[memberVaultID] = f.kp.Seed()

	res, err := f.useCase.Execute(context.Background(), collaboration_dtos.ResolveCollaborativeShareRequest{
		ShareEntryID:  f.shareEntry.ID,
		CallerVaultID: memberVaultID,
	})

	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, f.rawContent, res.Plaintext, "Valid member and envelope MUST proceed through decryption and return plaintext")
}

package trustgroup_orchestrator_test

import (
	"context"
	"os"
	"testing"


	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
)

type mockFileSystem struct{}

func (m *mockFileSystem) ReadFile(path string) ([]byte, error)                  { return nil, nil }
func (m *mockFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error { return nil }


func TestTrustGroupCryptoOrchestrator_FullFlow(t *testing.T) {
	ctx := context.Background()

	// 1. Setup Services & Devices
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, "/tmp/keyz", &mockFileSystem{})

	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	// Device 1 (Member 1 - Laptop)
	kp1, err := keypair.Random()
	require.NoError(t, err)

	// Device 2 (Member 1 - Mobile)
	kp2, err := keypair.Random()
	require.NoError(t, err)

	// Device 3 (Member 2 - Desktop)
	kp3, err := keypair.Random()
	require.NoError(t, err)

	// Device 4 (Member 2 - Revoked Device)
	kp4, err := keypair.Random()
	require.NoError(t, err)

	kr := &vaults_domain.VaultKeyring{
		UserID:  "user-1",
		VaultID: "vault-1",
	}

	payload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset-777",
		TrustGroupID: "tg-alpha",
		KEKVersion:   1,
		RawPayload:   []byte("TOP SECRET COLLABORATIVE ASSET CONTENT"),
		Keyring:      kr,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{
				DeviceID:  "dev-1-laptop",
				MemberID:  "member-1",
				PublicKey: kp1.Address(),
				IsActive:  true,
			},
			{
				DeviceID:  "dev-1-mobile",
				MemberID:  "member-1",
				PublicKey: kp2.Address(),
				IsActive:  true,
			},
			{
				DeviceID:  "dev-2-desktop",
				MemberID:  "member-2",
				PublicKey: kp3.Address(),
				IsActive:  true,
			},
			{
				DeviceID:  "dev-2-revoked",
				MemberID:  "member-2",
				PublicKey: kp4.Address(),
				IsActive:  false, // MUST BE EXCLUDED
			},
		},
	}

	// 2. Execute Orchestration
	prepared, err := orchestrator.PrepareCollaborativeAsset(ctx, payload)
	require.NoError(t, err)
	require.NotNil(t, prepared)

	// 3. Verify Structural Output
	assert.Equal(t, "asset-777", prepared.AssetID)
	assert.Equal(t, "tg-alpha", prepared.TrustGroupID)
	assert.Equal(t, uint64(1), prepared.KEKVersion)
	assert.NotEmpty(t, prepared.EncryptedData)
	assert.NotEmpty(t, prepared.WrappedDEK)

	// 4. Verify Active vs Revoked Device Envelopes
	// Only dev-1-laptop, dev-1-mobile, and dev-2-desktop (3 total active) should receive envelopes
	assert.Len(t, prepared.Envelopes, 3)

	var laptopEnv, mobileEnv, desktopEnv *trustgroup_orchestrator.PreparedCollaborativeAsset
	_ = laptopEnv
	_ = mobileEnv
	_ = desktopEnv

	envelopeMap := make(map[string]string)
	for _, env := range prepared.Envelopes {
		assert.Equal(t, "tg-alpha", env.TrustGroupID)
		assert.Equal(t, uint64(1), env.KEKVersion)
		assert.NotEmpty(t, env.WrappedKEK)
		envelopeMap[env.DeviceID] = env.WrappedKEK

		// Verify zero-knowledge boundary: DTO must contain no raw key material or private keys
		assert.NotContains(t, env.WrappedKEK, "TOP SECRET")
	}

	// Verify two devices belonging to member-1 receive distinct WrappedKEK envelopes
	assert.Contains(t, envelopeMap, "dev-1-laptop")
	assert.Contains(t, envelopeMap, "dev-1-mobile")
	assert.Contains(t, envelopeMap, "dev-2-desktop")
	assert.NotContains(t, envelopeMap, "dev-2-revoked", "Revoked device must not receive an envelope")
	assert.NotEqual(t, envelopeMap["dev-1-laptop"], envelopeMap["dev-1-mobile"], "Distinct devices must receive distinct wrapped envelopes")

	// 5. Verify KEK Storage in VaultKeyring
	storedKEKBytes, err := keyringSvc.GetTrustGroupKEK(kr, "tg-alpha", 1)
	require.NoError(t, err)
	assert.Len(t, storedKEKBytes, 32)

	// 6. LOCAL DEVICE UNWRAPPING VERIFICATION (DEVICE CRYPTO TEST)
	// Device 1 (laptop) unwraps its WrappedKEK envelope using its private key seed (kp1.Seed())
	wrappedKEKLaptop := envelopeMap["dev-1-laptop"]
	unwrappedKEK, err := aesSvc.AsymetricDecrypt(kp1.Seed(), wrappedKEKLaptop)
	require.NoError(t, err)
	assert.Equal(t, storedKEKBytes, unwrappedKEK, "Unwrapped KEK on device must match stored KEK")

	// Device 1 unwraps WrappedDEK using unwrapped KEK
	unwrappedDEK, err := aesSvc.Decrypt(prepared.WrappedDEK, unwrappedKEK)
	require.NoError(t, err)
	assert.Len(t, unwrappedDEK, 32)

	// Device 1 decrypts asset payload using unwrapped DEK
	decryptedAsset, err := aesSvc.Decrypt(prepared.EncryptedData, unwrappedDEK)
	require.NoError(t, err)
	assert.Equal(t, "TOP SECRET COLLABORATIVE ASSET CONTENT", string(decryptedAsset))

	// Device 3 (Member 2 Desktop) also unwraps successfully using its private key (kp3.Seed())
	wrappedKEKDesktop := envelopeMap["dev-2-desktop"]
	unwrappedKEK3, err := aesSvc.AsymetricDecrypt(kp3.Seed(), wrappedKEKDesktop)
	require.NoError(t, err)
	assert.Equal(t, storedKEKBytes, unwrappedKEK3)
}

func TestTrustGroupCryptoOrchestrator_ValidationErrors(t *testing.T) {
	ctx := context.Background()
	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, nil, nil)

	_, err := orchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		TrustGroupID: "",
		KEKVersion:   1,
		RawPayload:   []byte("data"),
	})
	assert.ErrorContains(t, err, "trust group ID is required")

	_, err = orchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		TrustGroupID: "tg-1",
		KEKVersion:   0,
		RawPayload:   []byte("data"),
	})
	assert.ErrorContains(t, err, "KEK version is required")

	_, err = orchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		TrustGroupID: "tg-1",
		KEKVersion:   1,
		RawPayload:   nil,
	})
	assert.ErrorContains(t, err, "raw payload cannot be empty")
}

func TestTrustGroupCryptoOrchestrator_ResolveCollaborativeAsset(t *testing.T) {
	ctx := context.Background()
	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}
	keyringSvc := vault_infrastructure_security.NewKeyringService(nil, nil, "/tmp/keyz2", &mockFileSystem{})

	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(keyringSvc, aesSvc, asymSvc)

	kp, err := keypair.Random()
	require.NoError(t, err)

	kr := &vaults_domain.VaultKeyring{UserID: "user-1", VaultID: "vault-1"}
	rawContent := []byte("CONFIDENTIAL AUDIT REPORT")

	prepPayload := trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
		AssetID:      "asset-999",
		TrustGroupID: "tg-finance",
		KEKVersion:   1,
		RawPayload:   rawContent,
		Keyring:      kr,
		ActiveDevices: []trustgroup_orchestrator.ActiveDevice{
			{DeviceID: "dev-laptop", MemberID: "member-1", PublicKey: kp.Address(), IsActive: true},
		},
	}

	prepared, err := orchestrator.PrepareCollaborativeAsset(ctx, prepPayload)
	require.NoError(t, err)
	require.Len(t, prepared.Envelopes, 1)

	wrappedKEK := prepared.Envelopes[0].WrappedKEK

	// 1. Resolve via Slow Path (No Keyring cache)
	emptyKr := &vaults_domain.VaultKeyring{UserID: "user-1", VaultID: "vault-1"}
	resolved, err := orchestrator.ResolveCollaborativeAsset(ctx, trustgroup_orchestrator.ResolveCollaborativeAssetPayload{
		AssetID:       prepared.AssetID,
		TrustGroupID:  prepared.TrustGroupID,
		KEKVersion:    prepared.KEKVersion,
		EncryptedData: prepared.EncryptedData,
		WrappedDEK:    prepared.WrappedDEK,
		WrappedKEK:    wrappedKEK,
		DeviceSeed:    kp.Seed(),
		Keyring:       emptyKr,
	})
	require.NoError(t, err)
	require.NotNil(t, resolved)
	assert.Equal(t, "asset-999", resolved.AssetID)
	assert.Equal(t, rawContent, resolved.Plaintext)

	// 2. Resolve via Fast Path (Keyring cache populated from step 1)
	resolvedFast, err := orchestrator.ResolveCollaborativeAsset(ctx, trustgroup_orchestrator.ResolveCollaborativeAssetPayload{
		AssetID:       prepared.AssetID,
		TrustGroupID:  prepared.TrustGroupID,
		KEKVersion:    prepared.KEKVersion,
		EncryptedData: prepared.EncryptedData,
		WrappedDEK:    prepared.WrappedDEK,
		Keyring:       emptyKr, // Now contains cached KEK
	})
	require.NoError(t, err)
	assert.Equal(t, rawContent, resolvedFast.Plaintext)

	// 3. Validation Errors & Invalid Seed
	_, err = orchestrator.ResolveCollaborativeAsset(ctx, trustgroup_orchestrator.ResolveCollaborativeAssetPayload{
		TrustGroupID:  "tg-finance",
		KEKVersion:    1,
		EncryptedData: prepared.EncryptedData,
		WrappedDEK:    prepared.WrappedDEK,
		WrappedKEK:    wrappedKEK,
		DeviceSeed:    "INVALID_SEED",
		Keyring:       &vaults_domain.VaultKeyring{},
	})
	assert.Error(t, err)
}

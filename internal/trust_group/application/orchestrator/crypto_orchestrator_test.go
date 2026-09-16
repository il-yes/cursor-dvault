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

func (m *mockFileSystem) ReadFile(path string) ([]byte, error)                       { return nil, nil }
func (m *mockFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error { return nil }

func TestTrustGroupCryptoOrchestrator_FullFlow(t *testing.T) {
	ctx := context.Background()

	keyringDir := t.TempDir()

	keyringSvc := vault_infrastructure_security.NewKeyringService(
		nil,
		nil,
		keyringDir,
		vault_infrastructure_security.OSFileSystem{}, // actual filesystem adapter used by production
	)

	aesSvc := &vault_infrastructure_crypto.AESService{}
	asymSvc := &vault_infrastructure_crypto.AsymmetricService{}

	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(
		keyringSvc,
		aesSvc,
		asymSvc,
	)

	kr := &vaults_domain.VaultKeyring{
		UserID:  "user-1",
		VaultID: "vault-1",
	}

	original := []byte("TOP SECRET COLLABORATIVE ASSET CONTENT")

	prepared, err := orchestrator.PrepareCollaborativeAsset(
		ctx,
		trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
			AssetID:       "asset-777",
			TrustGroupID:  "tg-alpha",
			KEKVersion:    1,
			RawPayload:    original,
			Keyring:       kr,
		},
	)
	require.NoError(t, err)

	require.NotEmpty(t, prepared.EncryptedData)
	require.NotEmpty(t, prepared.WrappedDEK)

	storedKEK, err := keyringSvc.GetTrustGroupKEK(
		kr,
		"tg-alpha",
		1,
	)
	require.NoError(t, err)
	require.Len(t, storedKEK, 32)

	resolved, err := orchestrator.ResolveCollaborativeAsset(
		ctx,
		trustgroup_orchestrator.ResolveCollaborativeAssetPayload{
			AssetID:       prepared.AssetID,
			TrustGroupID:  prepared.TrustGroupID,
			KEKVersion:    prepared.KEKVersion,
			EncryptedData: prepared.EncryptedData,
			WrappedDEK:    prepared.WrappedDEK,
			Keyring:       kr,
		},
	)
	require.NoError(t, err)

	require.Equal(t, original, resolved.Plaintext)
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

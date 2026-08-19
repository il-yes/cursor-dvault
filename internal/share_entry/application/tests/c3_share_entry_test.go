package share_entry_test

import (
	"context"
	"testing"

	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	app_config_domain "vault-app/internal/config/domain"
	share_entry_application_dto "vault-app/internal/share_entry/application"
	share_entry_application_events "vault-app/internal/share_entry/application/events"
	share_entry_application_use_cases "vault-app/internal/share_entry/application/use_cases"
	share_entry_domain "vault-app/internal/share_entry/domain"
	share_entry_infrastructure "vault-app/internal/share_entry/infrastructure"
	"vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

// ---------------------------------------------------------------------------
// Stubs
// ---------------------------------------------------------------------------
type stubTracecoreShareClient struct {
	lastCreatePayload *tracecore.ProdCreateCryptoShareRequest
}

func (s *stubTracecoreShareClient) CreateShare(_ context.Context, payload tracecore.ProdCreateCryptoShareRequest) (*tracecore.ProdCreateCryptoShareResponse, error) {
	s.lastCreatePayload = &payload
	return &tracecore.ProdCreateCryptoShareResponse{
		Status:  201,
		Code:    201,
		Message: "success",
		Data: tracecore.CloudCryptographicShare{
			ID: "share_cloud_001",
		},
	}, nil
}

func (s *stubTracecoreShareClient) AcceptShare(_ context.Context, _ tracecore_types.ShareAcceptedPayload) (*tracecore_types.CloudResponse[tracecore_types.PendingShareIntent], error) {
	return nil, nil
}
func (s *stubTracecoreShareClient) RejectShare(_ context.Context, _ tracecore_types.ShareRejectedPayload) (*tracecore_types.CloudResponse[tracecore_types.PendingShareIntent], error) {
	return nil, nil
}
func (s *stubTracecoreShareClient) GetShareByMe(_ context.Context, _ string) ([]share_entry_domain.ShareEntry, error) {
	return nil, nil
}
func (s *stubTracecoreShareClient) GetShareWithMe(_ context.Context, _ string) ([]share_entry_domain.ShareEntry, error) {
	return nil, nil
}
func (s *stubTracecoreShareClient) SetToken(_ string) {}
func (s *stubTracecoreShareClient) AddRecipient(_ context.Context, _ tracecore_types.AddRecipientRequest) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	return nil, nil
}
func (s *stubTracecoreShareClient) UpdateRecipient(_ context.Context, _ share_entry_application_dto.UpdateRecipientRequest) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	return nil, nil
}
func (s *stubTracecoreShareClient) RevokeShare(_ context.Context, _ tracecore_types.RevokeShareRequest) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	return nil, nil
}
func (s *stubTracecoreShareClient) ListPendingIntentSharesByMe(_ context.Context, _ string) (*tracecore_types.CloudResponse[[]tracecore_types.PendingShareIntent], error) {
	return nil, nil
}
func (s *stubTracecoreShareClient) ListPendingIntentSharesWithMe(_ context.Context, _ string) (*tracecore_types.CloudResponse[[]tracecore_types.PendingShareIntent], error) {
	return nil, nil
}

type stubEventDispatcher struct{}

func (s *stubEventDispatcher) Dispatch(_ share_entry_domain.DomainEvent) {}
func (s *stubEventDispatcher) Register(_ string, _ share_entry_application_events.EventHandler) {}

type stubSnapshotService struct{}

func (s *stubSnapshotService) Build(_ context.Context, _ share_entry_infrastructure.BuildRequest) (share_entry_infrastructure.BuildResponse, error) {
	snapshot := share_entry_domain.EntrySnapshot{
		EntryName: "Secret Contract",
		Type:      "secure_note",
		Note:      "Top secret contents",
	}
	return share_entry_infrastructure.BuildResponse{
		Raw:      []byte(`{"entry_name":"Secret Contract","type":"secure_note","note":"Top secret contents"}`),
		Snapshot: snapshot,
	}, nil
}

type stubConfigFacade struct{}

func (s *stubConfigFacade) GetUserConfigByUserID(_ string) (*app_config_domain.UserConfig, error) {
	kp, _ := keypair.Random()
	return &app_config_domain.UserConfig{
		StellarAccount: app_config_domain.StellarAccountConfig{
			PublicKey:  kp.Address(),
			PrivateKey: kp.Seed(),
		},
	}, nil
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestCreateShare_PersonalUserRecipient_Preserved(t *testing.T) {
	kp, err := keypair.Random()
	require.NoError(t, err)

	tc := &stubTracecoreShareClient{}
	aesSvc := &vault_infrastructure_crypto.AESService{}
	uc := share_entry_application_use_cases.NewShareUseCaseAES(
		nil,
		tc,
		&stubEventDispatcher{},
		aesSvc,
		&stubSnapshotService{},
	)

	share := share_entry_domain.ShareEntry{
		ID:         "se_personal",
		OwnerID:    "user_owner",
		EntryName:  "Personal Secret",
		EntryType:  "password",
		AccessMode: "read",
		Recipients: []share_entry_domain.Recipient{
			{
				ID:            "recip_1",
				Email:         "alice@example.com",
				PublicKey:     kp.Address(),
				Role:          "viewer",
				RecipientType: "user",
			},
		},
	}

	vault := &vaults_domain.Vault{Name: "default"}
	vp := vaults_domain.VaultPayload{}
	configFacade := &stubConfigFacade{}

	res, err := uc.Create(
		context.Background(),
		"user_owner",
		"owner@example.com",
		share,
		configFacade,
		"secret",
		vault,
		vp,
		"onboarding_1",
		app_config_domain.Config{},
		"sub_1",
	)

	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotNil(t, tc.lastCreatePayload)

	payload := tc.lastCreatePayload
	assert.Equal(t, "user_owner", payload.SenderID)
	assert.Equal(t, "owner@example.com", payload.SenderEmail)
	assert.Equal(t, "", payload.TrustGroupID) // Personal mode has no TrustGroupID
	assert.Contains(t, payload.Recipients, "alice@example.com")
	assert.Equal(t, "user", payload.Recipients["alice@example.com"].RecipientType)
	assert.NotEmpty(t, payload.VaultPayload) // Encrypted base64 payload
	assert.NotEmpty(t, payload.EncryptedKeys["alice@example.com"])
}

func TestCreateShare_TrustGroupRecipient_Success(t *testing.T) {
	kp, err := keypair.Random()
	require.NoError(t, err)

	tc := &stubTracecoreShareClient{}
	aesSvc := &vault_infrastructure_crypto.AESService{}
	uc := share_entry_application_use_cases.NewShareUseCaseAES(
		nil,
		tc,
		&stubEventDispatcher{},
		aesSvc,
		&stubSnapshotService{},
	)

	share := share_entry_domain.ShareEntry{
		ID:         "se_c3_group",
		OwnerID:    "user_owner",
		EntryName:  "Group Secret",
		EntryType:  "contract",
		AccessMode: "read",
		Recipients: []share_entry_domain.Recipient{
			{
				ID:            "recip_tg_1",
				TrustGroupID:  "tg_engineering_001",
				PublicKey:     kp.Address(),
				Role:          "member",
				RecipientType: "trust_group",
			},
		},
	}

	vault := &vaults_domain.Vault{Name: "default"}
	vp := vaults_domain.VaultPayload{}
	configFacade := &stubConfigFacade{}

	res, err := uc.Create(
		context.Background(),
		"user_owner",
		"owner@example.com",
		share,
		configFacade,
		"secret",
		vault,
		vp,
		"onboarding_1",
		app_config_domain.Config{},
		"sub_1",
	)

	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotNil(t, tc.lastCreatePayload)

	payload := tc.lastCreatePayload
	assert.Equal(t, "user_owner", payload.SenderID)
	assert.Equal(t, "tg_engineering_001", payload.TrustGroupID) // TrustGroupID correctly populated
	assert.Contains(t, payload.Recipients, "tg_engineering_001") // Keyed by TrustGroupID, not expanded
	assert.Equal(t, "trust_group", payload.Recipients["tg_engineering_001"].RecipientType)
	assert.Equal(t, "tg_engineering_001", payload.Recipients["tg_engineering_001"].TrustGroupID)
	assert.NotEmpty(t, payload.VaultPayload)
	assert.NotEmpty(t, payload.EncryptedKeys["tg_engineering_001"])
}

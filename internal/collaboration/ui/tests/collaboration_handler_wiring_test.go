package collaboration_ui_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_ui "vault-app/internal/collaboration/ui"
	collaboration_usecases "vault-app/internal/collaboration/application/usecases"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vaults_domain "vault-app/internal/vault/domain"
)

type dummyAssetResolver struct{}

func (d *dummyAssetResolver) FetchEncryptedAsset(_ context.Context, _ string) ([]byte, error) {
	return []byte("dummy_asset_data"), nil
}

type dummyIdentityResolver struct{}

func (d *dummyIdentityResolver) GetDeviceSeed(_ context.Context, userID string) (string, error) {
	return userID + "_seed", nil
}

func (d *dummyIdentityResolver) GetVaultKeyring(_ context.Context, userID string) (*vaults_domain.VaultKeyring, error) {
	return &vaults_domain.VaultKeyring{UserID: userID, VaultID: "v_1"}, nil
}

func TestCollaborationHandler_ResolveCollaborativeShare_InitializedUseCase(t *testing.T) {
	ctx := context.Background()

	tg := trustgroup_domain.NewTrustGroup("ch_1", "Test Group", []string{"user_alice"})
	tgRepo := &stubTrustGroupRepo{group: tg}

	shareEntry := c3_asset_domain.ShareEntry{
		ID:           "se_wired_100",
		TrustGroupID: tg.ID,
		AssetCID:     "bafybeiwired",
		WrappedDEK:   "wrapped_dek",
		KEKVersion:   1,
		CreatedBy:    "user_alice",
		Status:       c3_asset_domain.ShareEntryStatusActive,
	}

	shareRepo := &stubShareEntryRepo{
		createdEntries: []c3_asset_domain.ShareEntry{shareEntry},
	}

	orchestrator := trustgroup_orchestrator.NewTrustGroupCryptoOrchestrator(nil, nil, nil)
	resolveUC := collaboration_usecases.NewResolveCollaborativeShareUseCase(
		shareRepo,
		tgRepo,
		&dummyAssetResolver{},
		&dummyIdentityResolver{},
		orchestrator,
	)

	handler := collaboration_ui.NewCollaborationHandler(nil, resolveUC, nil)

	// Call handler directly -> MUST NOT return "resolve collaborative share use case is not initialized"
	_, err := handler.ResolveCollaborativeShare(ctx, "user_alice", "se_wired_100", "dev_100")
	if err != nil {
		assert.NotContains(t, err.Error(), "resolve collaborative share use case is not initialized")
	}
}

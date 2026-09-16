package collaboration_ui_test

import (
	"context"

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

package collaboration_infra

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	collaboration_ports "vault-app/internal/collaboration/application/ports"
	tracecore_types "vault-app/internal/tracecore/types"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
)

type TracecoreStorageClient interface {
	GetDataFromCloudStorage(ctx context.Context, req tracecore_types.IpfsCidRequest) (*tracecore_types.IpfsCidResponse, error)
}

type CloudAssetContentResolver struct {
	client TracecoreStorageClient
}

func NewCloudAssetContentResolver(client TracecoreStorageClient) *CloudAssetContentResolver {
	return &CloudAssetContentResolver{client: client}
}

var _ collaboration_ports.AssetContentResolver = (*CloudAssetContentResolver)(nil)

func (r *CloudAssetContentResolver) FetchEncryptedAsset(ctx context.Context, assetCID string) ([]byte, error) {
	if r.client == nil {
		return nil, fmt.Errorf("tracecore client is nil")
	}
	resp, err := r.client.GetDataFromCloudStorage(ctx, tracecore_types.IpfsCidRequest{CID: assetCID})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch asset content for CID %s: %w", assetCID, err)
	}
	if resp == nil || resp.Data == "" {
		return nil, fmt.Errorf("empty asset content returned for CID %s", assetCID)
	}
	rawBytes, err := base64.StdEncoding.DecodeString(resp.Data)
	if err != nil || len(rawBytes) == 0 {
		rawBytes = []byte(resp.Data)
	}
	return rawBytes, nil
}

type KeyringSovereignIdentityResolver struct {
	keyringService *vault_infrastructure_security.KeyringService
}

func NewKeyringSovereignIdentityResolver(keyringService *vault_infrastructure_security.KeyringService) *KeyringSovereignIdentityResolver {
	return &KeyringSovereignIdentityResolver{keyringService: keyringService}
}

var _ collaboration_ports.SovereignIdentityResolver = (*KeyringSovereignIdentityResolver)(nil)

func (r *KeyringSovereignIdentityResolver) GetDeviceSeed(ctx context.Context, userID string) (string, error) {
	if seed := os.Getenv("DEVICE_SEED"); seed != "" {
		return seed, nil
	}
	return userID + "_seed", nil
}

func (r *KeyringSovereignIdentityResolver) GetVaultKeyring(ctx context.Context, userID string) (*vaults_domain.VaultKeyring, error) {
	if r.keyringService != nil {
		kr, err := r.keyringService.LoadHybrid(userID, "", "")
		if err == nil && kr != nil {
			return kr, nil
		}
	}
	return &vaults_domain.VaultKeyring{UserID: userID, VaultID: "vault_" + userID}, nil
}

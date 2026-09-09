package collaboration_infra

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	collaboration_ports "vault-app/internal/collaboration/application/ports"
	app_config "vault-app/internal/config"
	tracecore_types "vault-app/internal/tracecore/types"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
)

type TracecoreStorageClient interface {
	GetDataFromCloudStorage(ctx context.Context, req tracecore_types.IpfsCidRequest) (*tracecore_types.IpfsCidResponse, error)
	AccessThreadData(ctx context.Context, req tracecore_types.ThreadDataAccessRequest) (*tracecore_types.CloudResponse[tracecore_types.AccessCryptoShareResponse], error)
}

type CloudAssetContentResolver struct {
	storage app_config.StorageProvider
	client  TracecoreStorageClient
}

func NewCloudAssetContentResolver(client TracecoreStorageClient) *CloudAssetContentResolver {
	return &CloudAssetContentResolver{client: client}
}

func (r *CloudAssetContentResolver) AccessThreadData(ctx context.Context, req tracecore_types.ThreadDataAccessRequest) (*tracecore_types.AccessCryptoShareResponse, error) {
	fmt.Printf("[C3-FORENSIC][09] (*CloudAssetContentResolver).AccessThreadData threadID=%s eventID=%s trustGroupID=%s requestingVaultID=%s sourceVaultID=%s\n", req.ThreadID, req.EventID, req.TrustGroupID, req.RequestingVaultID, req.SourceVaultID)
	if r.client != nil {
		resp, err := r.client.AccessThreadData(ctx, req)
		if err != nil {
			fmt.Printf("[C3-FORENSIC][09] (*CloudAssetContentResolver).AccessThreadData client call failed err=%v\n", err)
			return nil, fmt.Errorf("failed to access thread data: %w", err)
		}
		if resp == nil {
			fmt.Printf("[C3-FORENSIC][09] (*CloudAssetContentResolver).AccessThreadData nil response\n")
			return nil, fmt.Errorf("nil response from thread data access")
		}
		fmt.Printf("[C3-FORENSIC][09] (*CloudAssetContentResolver).AccessThreadData success downloadAllowed=%t encryptedPayloadLen=%d\n", resp.Data.DownloadAllowed, len(resp.Data.EncryptedPayload))
		return &resp.Data, nil
	}
	if r.storage != nil {
		return &tracecore_types.AccessCryptoShareResponse{
			EncryptedKey:     "",
			SenderPublicKey:  req.RequestingVaultID,
			EncryptedPayload: "",
			DownloadAllowed:  true,
		}, nil
	}
	return nil, fmt.Errorf("tracecore client and storage provider are both nil")
}

func NewCloudAssetContentResolverWithStorage(storage app_config.StorageProvider) *CloudAssetContentResolver {
	return &CloudAssetContentResolver{storage: storage}
}

var _ collaboration_ports.AssetContentResolver = (*CloudAssetContentResolver)(nil)

func (r *CloudAssetContentResolver) FetchEncryptedAsset(ctx context.Context, assetCID string) ([]byte, error) {
	if r.storage != nil {
		data, err := r.storage.Get(ctx, assetCID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch asset content for CID %s: %w", assetCID, err)
		}
		if len(data) == 0 {
			return nil, fmt.Errorf("empty asset content returned for CID %s", assetCID)
		}
		return data, nil
	}
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
	if r.keyringService != nil {
		kr, err := r.keyringService.LoadHybrid(userID, "", "")
		if err == nil && kr != nil {
			seedBytes, err := r.keyringService.GetKeyByType(kr, vaults_domain.KeyTypeDeviceSeed)
			if err == nil && len(seedBytes) > 0 {
				return string(seedBytes), nil
			}
		}
	}
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

package collaboration_usecases

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_ports "vault-app/internal/collaboration/application/ports"
	app_config "vault-app/internal/config"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_usecases "vault-app/internal/trust_group/application/usecases/envelope"
	vaults_domain "vault-app/internal/vault/domain"
)

type CreateCollaborativeShareUseCase struct {
	shareAssetUseCase  *ShareAssetWithTrustGroupUsecase
	addEnvelopeUseCase *trustgroup_usecases.AddTrustGroupKeyEnvelopeUseCase
	cryptoOrchestrator *trustgroup_orchestrator.TrustGroupCryptoOrchestrator
	assetResolver      collaboration_ports.AssetContentResolver
	identityResolver   collaboration_ports.SovereignIdentityResolver
	assetStorage       app_config.StorageProvider
}

func NewCreateCollaborativeShareUseCase(
	shareAssetUseCase *ShareAssetWithTrustGroupUsecase,
	addEnvelopeUseCase *trustgroup_usecases.AddTrustGroupKeyEnvelopeUseCase,
) *CreateCollaborativeShareUseCase {
	return &CreateCollaborativeShareUseCase{
		shareAssetUseCase:  shareAssetUseCase,
		addEnvelopeUseCase: addEnvelopeUseCase,
	}
}

func (u *CreateCollaborativeShareUseCase) WithCrypto(
	cryptoOrchestrator *trustgroup_orchestrator.TrustGroupCryptoOrchestrator,
	assetResolver collaboration_ports.AssetContentResolver,
	identityResolver collaboration_ports.SovereignIdentityResolver,
	assetStorage app_config.StorageProvider,
) *CreateCollaborativeShareUseCase {
	u.cryptoOrchestrator = cryptoOrchestrator
	u.assetResolver = assetResolver
	u.identityResolver = identityResolver
	u.assetStorage = assetStorage
	return u
}

func (u *CreateCollaborativeShareUseCase) WithStorageProvider(storage app_config.StorageProvider) *CreateCollaborativeShareUseCase {
	u.assetStorage = storage
	return u
}

func (u *CreateCollaborativeShareUseCase) ValidateDependencies() error {
	if u.shareAssetUseCase == nil {
		return errors.New("share asset use case is required")
	}
	return nil
}

func (u *CreateCollaborativeShareUseCase) ValidateRequest(req collaboration_dtos.CreateCollaborativeShareRequest) error {
	if strings.TrimSpace(req.TrustGroupID) == "" {
		return errors.New("trust group id is required")
	}
	if req.KEKVersion == 0 {
		return errors.New("kek version is required")
	}
	if strings.TrimSpace(req.CreatedBy) == "" {
		return errors.New("created by is required")
	}
	if strings.TrimSpace(req.AssetCID) == "" {
		return errors.New("asset cid is required")
	}
	return nil
}

func (u *CreateCollaborativeShareUseCase) Execute(
	ctx context.Context,
	req collaboration_dtos.CreateCollaborativeShareRequest,
) (*collaboration_dtos.CreateCollaborativeShareResponse, error) {
	if err := u.ValidateDependencies(); err != nil {
		return nil, err
	}
	if err := u.ValidateRequest(req); err != nil {
		return nil, err
	}

	wrappedDEKStr := req.WrappedDEK
	assetCIDToPersist := req.AssetCID

	// If cryptoOrchestrator is injected, execute single-DEK crypto preparation
	if u.cryptoOrchestrator != nil {
		var rawPayload []byte
		var err error

		if u.assetResolver != nil {
			rawPayload, err = u.assetResolver.FetchEncryptedAsset(ctx, req.AssetCID)
			if err != nil {
				// Fallback to raw CID bytes if fetcher does not hold CID
				rawPayload = []byte(req.AssetCID)
			}
		} else {
			rawPayload = []byte(req.AssetCID)
		}

		var keyring *vaults_domain.VaultKeyring
		if u.identityResolver != nil {
			keyring, _ = u.identityResolver.GetVaultKeyring(ctx, req.CreatedBy)
		}

		prepared, err := u.cryptoOrchestrator.PrepareCollaborativeAsset(ctx, trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
			AssetID:      req.AssetCID,
			TrustGroupID: req.TrustGroupID,
			KEKVersion:   req.KEKVersion,
			RawPayload:   rawPayload,
			Keyring:      keyring,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to prepare collaborative asset crypto: %w", err)
		}

		wrappedDEKStr = base64.StdEncoding.EncodeToString(prepared.WrappedDEK)

		// Upload prepared encrypted bytes to storage to obtain the new encrypted CID
		if u.assetStorage != nil {
			encCID, errUpload := u.assetStorage.Add(ctx, prepared.EncryptedData)
			if errUpload != nil {
				return nil, fmt.Errorf("failed to upload encrypted asset to cloud storage: %w", errUpload)
			}
			if encCID != "" {
				assetCIDToPersist = encCID
			}
		}
	}

	// Create and persist ShareEntry via ShareAssetWithTrustGroupUsecase
	shareEntry, err := u.shareAssetUseCase.Execute(ctx, collaboration_dtos.ShareAssetWithTrustGroupRequest{
		AssetCID:     assetCIDToPersist,
		TrustGroupID: req.TrustGroupID,
		WrappedDEK:   wrappedDEKStr,
		KEKVersion:   req.KEKVersion,
		CreatedBy:    req.CreatedBy,
		Metadata:     req.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create collaborative share entry: %w", err)
	}

	return &collaboration_dtos.CreateCollaborativeShareResponse{
		ShareEntry: *shareEntry,
	}, nil
}
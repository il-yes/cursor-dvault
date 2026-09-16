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
	trustgroup_domain "vault-app/internal/trust_group/domain"
	"vault-app/internal/utils"
	vault_dto "vault-app/internal/vault/application/dto"
)

type CreateCollaborativeShareUseCase struct {
	shareAssetUseCase  *ShareAssetWithTrustGroupUsecase
	addEnvelopeUseCase *trustgroup_usecases.AddTrustGroupKeyEnvelopeUseCase
	cryptoOrchestrator *trustgroup_orchestrator.TrustGroupCryptoOrchestrator
	assetResolver      collaboration_ports.AssetContentResolver
	identityResolver   collaboration_ports.SovereignIdentityResolver
	assetStorage       app_config.StorageProvider
	ipfsResolver	collaboration_ports.IPFSFileResolverInterface
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
	ipfsResolver collaboration_ports.IPFSFileResolverInterface,
) *CreateCollaborativeShareUseCase {
	u.cryptoOrchestrator = cryptoOrchestrator
	u.assetResolver = assetResolver
	u.identityResolver = identityResolver
	u.assetStorage = assetStorage
	u.ipfsResolver = ipfsResolver
	return u
}

func (u *CreateCollaborativeShareUseCase) Execute(ctx context.Context, req collaboration_dtos.CreateCollaborativeShareRequest) (*collaboration_dtos.CreateCollaborativeShareResponse, error) {

	if err := u.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := u.ValidateRequest(req); err != nil {
		return nil, err
	}

	tgResp, err := u.shareAssetUseCase.trustGroupRepo.GetTrustGroup(
		ctx,
		&trustgroup_domain.GetTrustGroupRequest{
			TrustGroupID: req.TrustGroupID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve trust group: %w", err)
	}
	if tgResp == nil || tgResp.Data.ID == "" {
		return nil, errors.New("trust group not found")
	}

	kekVersion := tgResp.Data.KEKVersion
	if kekVersion == 0 {
		return nil, errors.New("trust group KEK version is not initialized")
	}

	var assetCID = req.AssetCID
	var wrappedDEKBytes []byte

	// Get sharing data without decryption
	if u.cryptoOrchestrator != nil && u.assetResolver != nil && u.identityResolver != nil {
		keyring, err := u.identityResolver.GetVaultKeyring(ctx, req.UserID, req.Password, req.StellarSecret)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve vault keyring: %w", err)
		}
		if keyring == nil {
			return nil, errors.New("vault keyring is nil")
		}

		rawPayload, err := u.ipfsResolver.GetFileFromIPFS(ctx, vault_dto.GetFileFromIPFSRequest{
			CID:          req.AssetCID,
			UserID:       req.CreatedBy,
			Vault:        req.Vault,
			Password:     req.Password,
			PrivateKey:   req.StellarSecret,
			EncryptedKey: "",
			SymKey:    		[]byte{},
			Configs:      &req.Configs,
			IsShared:     true,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch encrypted asset: %w", err)
		}

		fmt.Printf(
			"[C3-FORENSIC][WRITE] rawPayload bytes=%d\n",
			len(rawPayload),
		)

		fmt.Printf(
			"[C3-FORENSIC][WRITE] rawPayload firstBytes=%x\n",
			rawPayload[:min(32, len(rawPayload))],
		)
		utils.LogPretty("CreateCollaborativeShareUseCase - Execute - rawPayload", rawPayload) // should be plain text

		prepared, err := u.cryptoOrchestrator.PrepareCollaborativeAsset(
			ctx,
			trustgroup_orchestrator.PrepareCollaborativeAssetPayload{
				AssetID:      req.AssetCID,
				TrustGroupID: req.TrustGroupID,
				KEKVersion:   kekVersion,
				RawPayload:   []byte(rawPayload),
				Keyring:      keyring,
			},
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to prepare collaborative asset crypto: %w",
				err,
			)
		}

		fmt.Printf(
			"[C3-FORENSIC][WRITE] prepared.EncryptedData bytes=%d\n",
			len(prepared.EncryptedData),
		)

		fmt.Printf(
			"[C3-FORENSIC][WRITE] prepared.EncryptedData firstBytes=%x\n",
			prepared.EncryptedData[:min(32, len(prepared.EncryptedData))],
		)

		wrappedDEKBytes = prepared.WrappedDEK
		utils.LogPretty("TrustGroupCryptoOrchestrator - PrepareCollaborativeAsset - assetStorage", u.assetStorage)

		if u.assetStorage != nil {
			assetCID, err = u.assetStorage.Add(ctx, prepared.EncryptedData)
			fmt.Printf(
				"[C3-FORENSIC][WRITE] storing encrypted asset bytes=%d\n",
				len(prepared.EncryptedData),
			)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to upload encrypted asset: %w",
					err,
				)
			}
		}
	}

	shareEntry, err := u.shareAssetUseCase.Execute(
		ctx,
		collaboration_dtos.ShareAssetWithTrustGroupRequest{
			AssetCID:     assetCID,
			TrustGroupID: req.TrustGroupID,
			WrappedDEK: base64.StdEncoding.EncodeToString(
				wrappedDEKBytes,
			),
			KEKVersion: kekVersion,
			CreatedBy:  req.CreatedBy,
			Metadata:   req.Metadata,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create collaborative share entry: %w",
			err,
		)
	}

	return &collaboration_dtos.CreateCollaborativeShareResponse{
		ShareEntry: *shareEntry,
	}, nil
}




func (u *CreateCollaborativeShareUseCase) SetAssetResolver(
	resolver collaboration_ports.AssetContentResolver,
) {
	u.assetResolver = resolver
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
	if strings.TrimSpace(req.CreatedBy) == "" {
		return errors.New("created by is required")
	}
	if strings.TrimSpace(req.AssetCID) == "" {
		return errors.New("asset cid is required")
	}
	return nil
}

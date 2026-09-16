package collaboration_usecases

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_ports "vault-app/internal/collaboration/application/ports"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	"vault-app/internal/utils"
)

var (
	ErrShareEntryNotFound  = errors.New("share entry not found")
	ErrShareEntryRevoked   = errors.New("share entry has been revoked")
	ErrTrustGroupNotFound  = errors.New("trust group not found")
	ErrUnauthorizedMember  = errors.New("caller is not an authorized member of trust group")
	ErrKeyEnvelopeNotFound = errors.New("no active key envelope found for member and KEK version")
)

type ResolveCollaborativeShareUseCase struct {
	shareEntryRepo     c3_asset_domain.ShareEntryRepository
	trustGroupRepo     trustgroup_domain.TrustGroupRepository
	assetResolver      collaboration_ports.AssetContentResolver
	identityResolver   collaboration_ports.SovereignIdentityResolver
	cryptoOrchestrator *trustgroup_orchestrator.TrustGroupCryptoOrchestrator
	ipfsFileResolver   collaboration_ports.IPFSFileResolverInterface
}

func NewResolveCollaborativeShareUseCase(
	shareEntryRepo c3_asset_domain.ShareEntryRepository,
	trustGroupRepo trustgroup_domain.TrustGroupRepository,
	assetResolver collaboration_ports.AssetContentResolver,
	identityResolver collaboration_ports.SovereignIdentityResolver,
	cryptoOrchestrator *trustgroup_orchestrator.TrustGroupCryptoOrchestrator,
	ipfsFileResolver collaboration_ports.IPFSFileResolverInterface,
) *ResolveCollaborativeShareUseCase {
	return &ResolveCollaborativeShareUseCase{
		shareEntryRepo:     shareEntryRepo,
		trustGroupRepo:     trustGroupRepo,
		assetResolver:      assetResolver,
		identityResolver:   identityResolver,
		cryptoOrchestrator: cryptoOrchestrator,
		ipfsFileResolver:   ipfsFileResolver,
	}
}

func (u *ResolveCollaborativeShareUseCase) SetAssetResolver(resolver collaboration_ports.AssetContentResolver) {
	u.assetResolver = resolver
}

func (u *ResolveCollaborativeShareUseCase) ValidateDependencies() error {
	if u.shareEntryRepo == nil {
		return c3_asset_domain.ErrRepositoryNil
	}
	if u.trustGroupRepo == nil {
		return trustgroup_domain.ErrRepositoryNil
	}
	if u.assetResolver == nil {
		return errors.New("asset content resolver is required")
	}
	if u.identityResolver == nil {
		return errors.New("sovereign identity resolver is required")
	}
	if u.cryptoOrchestrator == nil {
		return errors.New("crypto orchestrator is required")
	}
	return nil
}

func (u *ResolveCollaborativeShareUseCase) ValidateRequest(req collaboration_dtos.ResolveCollaborativeShareRequest) error {
	if strings.TrimSpace(req.ShareEntryID) == "" {
		return errors.New("share entry id is required")
	}
	if strings.TrimSpace(req.CallerVaultID) == "" && strings.TrimSpace(req.CallerUserID) == "" {
		return errors.New("caller vault id is required")
	}
	return nil
}


func (u *ResolveCollaborativeShareUseCase) Execute(
	ctx context.Context,
	req collaboration_dtos.ResolveCollaborativeShareRequest,
) (*collaboration_dtos.ResolveCollaborativeShareResponse, error) {

	if err := u.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := u.ValidateRequest(req); err != nil {
		return nil, err
	}

	fmt.Printf(
		"[C3][READ][01] start shareEntryID=%s callerIdentityID=%s callerVaultID=%s\n",
		req.ShareEntryID,
		req.CallerIdentityID,
		req.CallerVaultID,
	)

	// 1. Load ShareEntry
	shareResp, err := u.shareEntryRepo.GetShareEntry(
		ctx,
		&c3_asset_domain.GetShareEntryRequest{
			ShareEntryID: req.ShareEntryID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch share entry: %w", err)
	}

	if shareResp == nil {
		return nil, fmt.Errorf("failed to fetch share entry: repository returned nil response")
	}

	if shareResp.Data.ID == "" {
		return nil, ErrShareEntryNotFound
	}

	shareEntry := shareResp.Data

	fmt.Printf(
		"[C3][READ][02] share entry loaded id=%s trustGroupID=%s assetCID=%s kekVersion=%d wrappedDEK=%t status=%s\n",
		shareEntry.ID,
		shareEntry.TrustGroupID,
		shareEntry.AssetCID,
		shareEntry.KEKVersion,
		shareEntry.WrappedDEK != "",
		shareEntry.Status,
	)

	if shareEntry.Status == c3_asset_domain.ShareEntryStatusRevoked {
		return nil, ErrShareEntryRevoked
	}

	// 2. Load TrustGroup
	tgResp, err := u.trustGroupRepo.GetTrustGroup(
		ctx,
		&trustgroup_domain.GetTrustGroupRequest{
			TrustGroupID: shareEntry.TrustGroupID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trust group: %w", err)
	}

	if tgResp == nil {
		return nil, fmt.Errorf("failed to fetch trust group: repository returned nil response")
	}

	if tgResp.Data.ID == "" {
		return nil, ErrTrustGroupNotFound
	}

	trustGroup := tgResp.Data

	fmt.Printf(
		"[C3][READ][03] trust group loaded id=%s members=%d envelopes=%d kekVersion=%d\n",
		trustGroup.ID,
		len(trustGroup.MemberCIDs),
		len(trustGroup.KeyEnvelopes),
		trustGroup.KEKVersion,
	)

	// 3. Authorize caller
	if !trustGroup.HasMember(req.CallerVaultID) {
		return nil, ErrUnauthorizedMember
	}

	utils.LogPretty("[C3][READ][04] caller authorized", map[string]any{
		"callerVaultID": req.CallerVaultID,
		"trustGroupID":  trustGroup.ID,
	})

	// 4. Resolve caller's envelope
	envelope := resolveActiveTrustGroupEnvelope(
		trustGroup.KeyEnvelopes,
		req.CallerVaultID,
		shareEntry.KEKVersion,
	)

	utils.LogPretty("[C3][READ][05] envelope resolved", map[string]any{
		"envelopeNil":     envelope == nil,
		"callerVaultID":   req.CallerVaultID,
		"shareKEKVersion": shareEntry.KEKVersion,
	})

	if envelope == nil {
		return nil, ErrKeyEnvelopeNotFound
	}

	if envelope.WrappedKEK == "" {
		return nil, fmt.Errorf(
			"resolved key envelope has empty WrappedKEK: memberID=%s kekVersion=%d",
			envelope.MemberID,
			envelope.KEKVersion,
		)
	}

	utils.LogPretty("[C3][READ][05B] envelope usable", map[string]any{
		"memberID":    envelope.MemberID,
		"kekVersion":  envelope.KEKVersion,
		"wrappedKEK":  true,
		"envelopeID":  envelope.ID,
	})

	// 5. Retrieve encrypted asset
	fmt.Printf(
		"[C3][READ][06] retrieving encrypted asset cid=%s\n",
		shareEntry.AssetCID,
	)

	req.GetIPFSFile.CID = shareEntry.AssetCID

	encryptedPayload, err := u.ipfsFileResolver.GetFileFromIPFS(
		ctx,
		req.GetIPFSFile,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to retrieve encrypted asset: %w",
			err,
		)
	}

	fmt.Printf(
		"[C3][READ][07] encrypted asset retrieved bytes=%d\n",
		len(encryptedPayload),
	)

	// 6. Resolve existing caller keyring - I'M NOT SURE IT'S CORRECT
	fmt.Printf(
		"[C3][READ][08] resolving caller keyring identityID=%s\n",
		req.CallerIdentityID,
	)



	fmt.Printf(
		"[C3][READ][09] caller keyring resolved identityID=%s\n",
		req.CallerIdentityID,
	)

	// 7. Decode wrapped DEK
	fmt.Printf(
		"[C3][READ][10] decoding wrapped DEK shareEntryID=%s\n",
		shareEntry.ID,
	)

	wrappedDEK, err := base64.StdEncoding.DecodeString(shareEntry.WrappedDEK)
	if err != nil {
		return nil, fmt.Errorf("failed to decode wrapped DEK: %w", err)
	}
	utils.LogPretty("ResolveCollaborativeShareUseCase - wrappedDEK", utils.FingerprintKey(wrappedDEK))

	if len(wrappedDEK) == 0 {
		return nil, fmt.Errorf("failed to decode wrapped DEK: decoded value is empty")
	}

	fmt.Printf(
		"[C3][READ][11] wrapped DEK decoded bytes=%d\n",
		len(wrappedDEK),
	)

	// 8. Resolve collaborative asset
	fmt.Printf(
		"[C3][READ][12] resolving collaborative asset memberID=%s kekVersion=%d\n",
		envelope.MemberID,
		envelope.KEKVersion,
	)

	cryptoResult, err := u.cryptoOrchestrator.ResolveCollaborativeAsset(
		ctx,
		trustgroup_orchestrator.ResolveCollaborativeAssetPayload{
			AssetID:       shareEntry.ID,
			TrustGroupID:  shareEntry.TrustGroupID,
			KEKVersion:    int(shareEntry.KEKVersion),
			EncryptedData: []byte(encryptedPayload),
			WrappedDEK:    wrappedDEK,
			WrappedKEK:    envelope.WrappedKEK,
			PrivateKey:    req.StellarAccount.PrivateKey,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve collaborative asset: %w",
			err,
		)
	}

	if cryptoResult == nil {
		return nil, fmt.Errorf(
			"failed to resolve collaborative asset: crypto orchestrator returned nil result",
		)
	}

	fmt.Printf(
		"[C3][READ][13] collaborative asset resolved plaintextBytes=%d\n",
		len(cryptoResult.Plaintext),
	)



	// 9. Return plaintext
	return &collaboration_dtos.ResolveCollaborativeShareResponse{
		ShareEntryID: shareEntry.ID,
		TrustGroupID: shareEntry.TrustGroupID,
		CreatedBy:    shareEntry.CreatedBy,
		CreatedAt:    shareEntry.CreatedAt.Format(time.RFC3339),
		Metadata:     shareEntry.Metadata,
		Plaintext:    cryptoResult.Plaintext,
	}, nil
}

func decodeWrappedValue(value string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err == nil && len(decoded) > 0 {
		return decoded
	}

	return []byte(value)
}

func resolveActiveTrustGroupEnvelope(
	envelopes []trustgroup_domain.TrustGroupKeyEnvelope,
	callerVaultID string,
	kekVersion uint64,
) *trustgroup_domain.TrustGroupKeyEnvelope {
	for i := range envelopes {
		envelope := &envelopes[i]

		if envelope.MemberID != callerVaultID {
			continue
		}

		if envelope.KEKVersion != kekVersion {
			continue
		}

		if envelope.RevokedAt != nil {
			continue
		}

		return envelope
	}

	return nil
}

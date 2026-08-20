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
)

var (
	ErrShareEntryNotFound   = errors.New("share entry not found")
	ErrTrustGroupNotFound  = errors.New("trust group not found")
	ErrUnauthorizedMember  = errors.New("caller is not an authorized member of trust group")
	ErrKeyEnvelopeNotFound = errors.New("no active device key envelope found for member device and KEK version")
)

type ResolveCollaborativeShareUseCase struct {
	shareEntryRepo            c3_asset_domain.ShareEntryRepository
	trustGroupRepo            trustgroup_domain.TrustGroupRepository
	assetResolver             collaboration_ports.AssetContentResolver
	identityResolver          collaboration_ports.SovereignIdentityResolver
	cryptoOrchestrator        *trustgroup_orchestrator.TrustGroupCryptoOrchestrator
}

func NewResolveCollaborativeShareUseCase(
	shareEntryRepo c3_asset_domain.ShareEntryRepository,
	trustGroupRepo trustgroup_domain.TrustGroupRepository,
	assetResolver collaboration_ports.AssetContentResolver,
	identityResolver collaboration_ports.SovereignIdentityResolver,
	cryptoOrchestrator *trustgroup_orchestrator.TrustGroupCryptoOrchestrator,
) *ResolveCollaborativeShareUseCase {
	return &ResolveCollaborativeShareUseCase{
		shareEntryRepo:            shareEntryRepo,
		trustGroupRepo:            trustGroupRepo,
		assetResolver:             assetResolver,
		identityResolver:          identityResolver,
		cryptoOrchestrator:        cryptoOrchestrator,
	}
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
	if strings.TrimSpace(req.CallerUserID) == "" {
		return errors.New("caller user id is required")
	}
	if strings.TrimSpace(req.DeviceID) == "" {
		return errors.New("device id is required")
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

	// 1. Fetch ShareEntry (Access Descriptor)
	shareResp, err := u.shareEntryRepo.GetShareEntry(ctx, &c3_asset_domain.GetShareEntryRequest{
		ShareEntryID: req.ShareEntryID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch share entry: %w", err)
	}
	if shareResp == nil || shareResp.Data.ID == "" {
		return nil, ErrShareEntryNotFound
	}
	shareEntry := shareResp.Data

	// 2. Fetch TrustGroup
	tgResp, err := u.trustGroupRepo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{
		TrustGroupID: shareEntry.TrustGroupID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trust group: %w", err)
	}
	if tgResp == nil || tgResp.Data.ID == "" {
		return nil, ErrTrustGroupNotFound
	}
	trustGroup := tgResp.Data

	// 3. Authorize Member: Verify CallerUserID is in MemberCIDs BEFORE resolving assets or key material
	isMember := false
	for _, cid := range trustGroup.MemberCIDs {
		if cid == req.CallerUserID {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, ErrUnauthorizedMember
	}

	// 4. Authorize & Resolve Active Device Envelope BEFORE resolving assets or key material
	var activeEnvelope *trustgroup_domain.TrustGroupKeyEnvelope
	for i := range trustGroup.KeyEnvelopes {
		env := &trustGroup.KeyEnvelopes[i]
		if env.MemberID == req.CallerUserID &&
			env.DeviceID == req.DeviceID &&
			env.KEKVersion == shareEntry.KEKVersion &&
			env.RevokedAt == nil {
			activeEnvelope = env
			break
		}
	}

	if activeEnvelope == nil {
		return nil, ErrKeyEnvelopeNotFound
	}

	// 5. Fetch Encrypted Asset Content Bytes via AssetContentResolver (Only AFTER authorization)
	encryptedData, err := u.assetResolver.FetchEncryptedAsset(ctx, shareEntry.AssetCID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch encrypted asset content for CID %s: %w", shareEntry.AssetCID, err)
	}
	if len(encryptedData) == 0 {
		return nil, errors.New("encrypted asset payload data is empty")
	}

	// 6. Resolve Local Member Device Credentials via SovereignIdentityResolver (Only AFTER authorization)
	deviceSeed, err := u.identityResolver.GetDeviceSeed(ctx, req.CallerUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve local device seed for user %s: %w", req.CallerUserID, err)
	}
	keyring, err := u.identityResolver.GetVaultKeyring(ctx, req.CallerUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve local keyring for user %s: %w", req.CallerUserID, err)
	}

	// 7. Decode WrappedDEK (Base64 string or raw bytes)
	wrappedDEKBytes, err := base64.StdEncoding.DecodeString(shareEntry.WrappedDEK)
	if err != nil || len(wrappedDEKBytes) == 0 {
		wrappedDEKBytes = []byte(shareEntry.WrappedDEK)
	}

	// 8. Invoke Cryptographic Resolution (Local Sovereign Unwrapping & Decryption)
	cryptoResult, err := u.cryptoOrchestrator.ResolveCollaborativeAsset(ctx, trustgroup_orchestrator.ResolveCollaborativeAssetPayload{
		AssetID:       shareEntry.ID,
		TrustGroupID:  shareEntry.TrustGroupID,
		KEKVersion:    shareEntry.KEKVersion,
		EncryptedData: encryptedData,
		WrappedDEK:    wrappedDEKBytes,
		WrappedKEK:    activeEnvelope.WrappedKEK,
		DeviceSeed:    deviceSeed,
		Keyring:       keyring,
	})
	if err != nil {
		return nil, fmt.Errorf("cryptographic resolution failed: %w", err)
	}

	createdAtStr := shareEntry.CreatedAt.Format(time.RFC3339)

	// 9. Return Clean Response DTO (Zero Secret Leakage)
	return &collaboration_dtos.ResolveCollaborativeShareResponse{
		ShareEntryID: shareEntry.ID,
		TrustGroupID: shareEntry.TrustGroupID,
		CreatedBy:    shareEntry.CreatedBy,
		CreatedAt:    createdAtStr,
		Metadata:     shareEntry.Metadata,
		Plaintext:    cryptoResult.Plaintext,
	}, nil
}

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
	ErrShareEntryNotFound  = errors.New("share entry not found")
	ErrShareEntryRevoked   = errors.New("share entry has been revoked")
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
	if strings.TrimSpace(req.CallerVaultID) == "" && strings.TrimSpace(req.CallerUserID) == "" {
		return errors.New("caller vault id is required")
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

	if req.CallerVaultID == "" && req.CallerUserID != "" {
		req.CallerVaultID = req.CallerUserID
	}
	if req.CallerIdentityID == "" && req.CallerUserID != "" {
		req.CallerIdentityID = req.CallerUserID
	}
	if req.CallerIdentityID == "" {
		req.CallerIdentityID = req.CallerVaultID
	}

	// 1. Fetch ShareEntry (Access Descriptor)
	shareResp, err := u.shareEntryRepo.GetShareEntry(ctx, &c3_asset_domain.GetShareEntryRequest{
		ShareEntryID: req.ShareEntryID,
	})
	if err != nil {
		fmt.Printf("[C3][READ][DIAGNOSTIC] shareEntryID=%s callerIdentityID=%s deviceID=%s FETCH_SHARE_ENTRY_ERROR=%v\n", req.ShareEntryID, req.CallerIdentityID, req.DeviceID, err)
		return nil, fmt.Errorf("failed to fetch share entry: %w", err)
	}
	if shareResp == nil || shareResp.Data.ID == "" {
		fmt.Printf("[C3][READ][DIAGNOSTIC] shareEntryID=%s callerIdentityID=%s deviceID=%s ERR=ErrShareEntryNotFound\n", req.ShareEntryID, req.CallerIdentityID, req.DeviceID)
		return nil, ErrShareEntryNotFound
	}
	shareEntry := shareResp.Data
	if shareEntry.Status == c3_asset_domain.ShareEntryStatusRevoked {
		fmt.Printf("[C3][READ][DIAGNOSTIC] shareEntryID=%s trustGroupID=%s callerIdentityID=%s deviceID=%s shareEntryStatus=%s ERR=ErrShareEntryRevoked\n", shareEntry.ID, shareEntry.TrustGroupID, req.CallerIdentityID, req.DeviceID, shareEntry.Status)
		return nil, ErrShareEntryRevoked
	}

	// 2. Fetch TrustGroup
	tgResp, err := u.trustGroupRepo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{
		TrustGroupID: shareEntry.TrustGroupID,
	})
	if err != nil {
		fmt.Printf("[C3][READ][DIAGNOSTIC] shareEntryID=%s trustGroupID=%s callerIdentityID=%s deviceID=%s FETCH_TG_ERROR=%v\n", shareEntry.ID, shareEntry.TrustGroupID, req.CallerIdentityID, req.DeviceID, err)
		return nil, fmt.Errorf("failed to fetch trust group: %w", err)
	}
	if tgResp == nil || tgResp.Data.ID == "" {
		fmt.Printf("[C3][READ][DIAGNOSTIC] shareEntryID=%s trustGroupID=%s callerIdentityID=%s deviceID=%s ERR=ErrTrustGroupNotFound\n", shareEntry.ID, shareEntry.TrustGroupID, req.CallerIdentityID, req.DeviceID)
		return nil, ErrTrustGroupNotFound
	}
	trustGroup := tgResp.Data

	// 3. Authorize Member: Verify CallerVaultID is in MemberCIDs BEFORE resolving assets or key material
	isMember := false
	for _, cid := range trustGroup.MemberCIDs {
		if cid == req.CallerVaultID {
			isMember = true
			break
		}
	}

	membershipStatus := "unauthorized"
	if isMember {
		membershipStatus = "active"
	}

	memberIDsFormatted := strings.Join(trustGroup.MemberCIDs, ", ")
	fmt.Printf("[C3][AUTHZ][READ] identityID=%s resolvedVaultID=%s trustGroupID=%s memberVaultIDs=[%s] comparisonIdentifier=%s membershipFound=%t membershipStatus=%s\n",
		req.CallerIdentityID,
		req.CallerVaultID,
		trustGroup.ID,
		memberIDsFormatted,
		req.CallerVaultID,
		isMember,
		membershipStatus,
	)

	if !isMember {
		fmt.Printf("[C3][READ][DIAGNOSTIC] shareEntryID=%s trustGroupID=%s callerIdentityID=%q callerVaultID=%q deviceID=%s shareEntryStatus=%s shareEntryKEKVersion=%d trustGroupKEKVersion=%d membershipFound=false membershipStatus=%s envelopeFound=false envelopeKEKVersion=0 envelopeRevoked=false ERR=ErrUnauthorizedMember\n",
			shareEntry.ID, shareEntry.TrustGroupID, req.CallerIdentityID, req.CallerVaultID, req.DeviceID, shareEntry.Status, shareEntry.KEKVersion, trustGroup.KEKVersion, membershipStatus)
		return nil, ErrUnauthorizedMember
	}

	// 4. Authorize & Resolve Active Device Envelope BEFORE resolving assets or key material
	var activeEnvelope *trustgroup_domain.TrustGroupKeyEnvelope
	var inspectEnvKEKVer uint64
	var inspectEnvRevoked bool
	var matchingMemberID, matchingDeviceID string

	for i := range trustGroup.KeyEnvelopes {
		env := &trustGroup.KeyEnvelopes[i]
		if (env.MemberID == req.CallerVaultID || env.MemberID == req.CallerIdentityID) && env.DeviceID == req.DeviceID {
			inspectEnvKEKVer = env.KEKVersion
			inspectEnvRevoked = (env.RevokedAt != nil)
			matchingMemberID = env.MemberID
			matchingDeviceID = env.DeviceID
			if env.KEKVersion == shareEntry.KEKVersion && env.RevokedAt == nil {
				activeEnvelope = env
				break
			}
		}
	}

	resultStr := "NOT_FOUND"
	reasonStr := "no matching device envelope found"
	if activeEnvelope != nil {
		resultStr = "FOUND"
		reasonStr = "active envelope resolved"
	} else if inspectEnvRevoked {
		reasonStr = "device envelope is revoked"
	} else if inspectEnvKEKVer > 0 && inspectEnvKEKVer != shareEntry.KEKVersion {
		reasonStr = fmt.Sprintf("KEK version mismatch: envelope has %d, shareEntry needs %d", inspectEnvKEKVer, shareEntry.KEKVersion)
	} else if len(trustGroup.KeyEnvelopes) == 0 {
		reasonStr = "trust group has 0 key envelopes"
	}

	matchingRevokedAtStr := "none"
	if inspectEnvRevoked {
		matchingRevokedAtStr = "true"
	} else if matchingDeviceID != "" {
		matchingRevokedAtStr = "false"
	}

	fmt.Printf("[C3][ENVELOPE][READ] trustGroupID=%s shareEntryID=%s callerIdentityID=%s callerVaultID=%s callerDeviceID=%s shareEntryKEKVersion=%d trustGroupKEKVersion=%d envelopeCandidates=%d matchingMemberID=%s matchingDeviceID=%s matchingKEKVersion=%d matchingRevokedAt=%s result=%s reason=%q\n",
		trustGroup.ID, shareEntry.ID, req.CallerIdentityID, req.CallerVaultID, req.DeviceID, shareEntry.KEKVersion, trustGroup.KEKVersion, len(trustGroup.KeyEnvelopes), matchingMemberID, matchingDeviceID, inspectEnvKEKVer, matchingRevokedAtStr, resultStr, reasonStr)

	if activeEnvelope == nil {
		fmt.Printf("[C3][READ][DIAGNOSTIC] shareEntryID=%s trustGroupID=%s callerIdentityID=%s callerVaultID=%s deviceID=%s shareEntryStatus=%s shareEntryKEKVersion=%d trustGroupKEKVersion=%d membershipFound=true membershipStatus=%s envelopeFound=false envelopeKEKVersion=%d envelopeRevoked=%t ERR=ErrKeyEnvelopeNotFound\n",
			shareEntry.ID, shareEntry.TrustGroupID, req.CallerIdentityID, req.CallerVaultID, req.DeviceID, shareEntry.Status, shareEntry.KEKVersion, trustGroup.KEKVersion, membershipStatus, inspectEnvKEKVer, inspectEnvRevoked)
		return nil, ErrKeyEnvelopeNotFound
	}

	fmt.Printf("[C3][READ][DIAGNOSTIC] shareEntryID=%s trustGroupID=%s callerIdentityID=%s callerVaultID=%s deviceID=%s shareEntryStatus=%s shareEntryKEKVersion=%d trustGroupKEKVersion=%d membershipFound=true membershipStatus=%s envelopeFound=true envelopeKEKVersion=%d envelopeRevoked=false SUCCESS_AUTH=true\n",
		shareEntry.ID, shareEntry.TrustGroupID, req.CallerIdentityID, req.CallerVaultID, req.DeviceID, shareEntry.Status, shareEntry.KEKVersion, trustGroup.KEKVersion, membershipStatus, activeEnvelope.KEKVersion)

	// 5. Fetch Encrypted Asset Content Bytes via AssetContentResolver (Only AFTER authorization)
	encryptedData, err := u.assetResolver.FetchEncryptedAsset(ctx, shareEntry.AssetCID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch encrypted asset content for CID %s: %w", shareEntry.AssetCID, err)
	}
	if len(encryptedData) == 0 {
		return nil, errors.New("encrypted asset payload data is empty")
	}

	// 6. Resolve Local Member Device Credentials via SovereignIdentityResolver (Only AFTER authorization)
	deviceSeed, err := u.identityResolver.GetDeviceSeed(ctx, req.CallerIdentityID)
	if err != nil {
		deviceSeed, err = u.identityResolver.GetDeviceSeed(ctx, req.CallerVaultID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to resolve local device seed for user %s: %w", req.CallerVaultID, err)
	}
	keyring, err := u.identityResolver.GetVaultKeyring(ctx, req.CallerIdentityID)
	if err != nil {
		keyring, err = u.identityResolver.GetVaultKeyring(ctx, req.CallerVaultID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to resolve local keyring for user %s: %w", req.CallerVaultID, err)
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

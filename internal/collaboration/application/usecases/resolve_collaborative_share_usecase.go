package collaboration_usecases

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/stellar/go/keypair"

	blockchain "vault-app/internal/blockchain"
	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	collaboration_ports "vault-app/internal/collaboration/application/ports"
	tracecore_types "vault-app/internal/tracecore/types"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_domain "vault-app/internal/trust_group/domain"
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
}

func NewResolveCollaborativeShareUseCase(
	shareEntryRepo c3_asset_domain.ShareEntryRepository,
	trustGroupRepo trustgroup_domain.TrustGroupRepository,
	assetResolver collaboration_ports.AssetContentResolver,
	identityResolver collaboration_ports.SovereignIdentityResolver,
	cryptoOrchestrator *trustgroup_orchestrator.TrustGroupCryptoOrchestrator,
) *ResolveCollaborativeShareUseCase {
	return &ResolveCollaborativeShareUseCase{
		shareEntryRepo:     shareEntryRepo,
		trustGroupRepo:     trustGroupRepo,
		assetResolver:      assetResolver,
		identityResolver:   identityResolver,
		cryptoOrchestrator: cryptoOrchestrator,
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

	// Resolve caller identity context once.
	if req.CallerIdentityID == "" {
		req.CallerIdentityID = req.CallerUserID
	}

	if req.CallerVaultID == "" {
		req.CallerVaultID = req.CallerUserID
	}

	// -------------------------------------------------------------------------
	// 1. Load ShareEntry
	// -------------------------------------------------------------------------

	shareResp, err := u.shareEntryRepo.GetShareEntry(
		ctx,
		&c3_asset_domain.GetShareEntryRequest{
			ShareEntryID: req.ShareEntryID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch share entry: %w", err)
	}

	if shareResp == nil || shareResp.Data.ID == "" {
		return nil, ErrShareEntryNotFound
	}

	shareEntry := shareResp.Data

	if shareEntry.Status == c3_asset_domain.ShareEntryStatusRevoked {
		return nil, ErrShareEntryRevoked
	}

	// -------------------------------------------------------------------------
	// 2. Load TrustGroup
	// -------------------------------------------------------------------------

	tgResp, err := u.trustGroupRepo.GetTrustGroup(
		ctx,
		&trustgroup_domain.GetTrustGroupRequest{
			TrustGroupID: shareEntry.TrustGroupID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trust group: %w", err)
	}

	if tgResp == nil || tgResp.Data.ID == "" {
		return nil, ErrTrustGroupNotFound
	}

	trustGroup := tgResp.Data

	// -------------------------------------------------------------------------
	// 3. Authorize caller by VaultID
	// -------------------------------------------------------------------------

	if !trustGroup.HasMember(req.CallerVaultID) {
		return nil, ErrUnauthorizedMember
	}

	// -------------------------------------------------------------------------
	// 4. Resolve identity-level TrustGroup key envelope
	//
	// DeviceID is deliberately NOT part of this flow.
	// The envelope belongs to the TrustGroup member (VaultID).
	// -------------------------------------------------------------------------

	activeEnvelope := resolveActiveTrustGroupEnvelope(
		trustGroup.KeyEnvelopes,
		req.CallerVaultID,
		shareEntry.KEKVersion,
	)

	if activeEnvelope == nil {
		return nil, ErrKeyEnvelopeNotFound
	}

	// -------------------------------------------------------------------------
	// 5. Resolve caller key material
	//
	// The identity key is the cryptographic key used to unwrap the envelope.
	// -------------------------------------------------------------------------

	deviceSeed, err := u.identityResolver.GetDeviceSeed(
		ctx,
		req.CallerIdentityID,
	)
	if err != nil {
		deviceSeed, err = u.identityResolver.GetDeviceSeed(
			ctx,
			req.CallerVaultID,
		)
	}
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve caller key material for identity %s: %w",
			req.CallerIdentityID,
			err,
		)
	}

	keyring, err := u.identityResolver.GetVaultKeyring(
		ctx,
		req.CallerIdentityID,
	)
	if err != nil {
		keyring, err = u.identityResolver.GetVaultKeyring(
			ctx,
			req.CallerVaultID,
		)
	}
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve caller keyring for identity %s: %w",
			req.CallerIdentityID,
			err,
		)
	}

	// -------------------------------------------------------------------------
	// 6. Authenticate caller
	// -------------------------------------------------------------------------

	pubKey := req.CallerIdentityID

	if kp, err := keypair.ParseFull(deviceSeed); err == nil && kp != nil {
		pubKey = kp.Address()
	}

	challenge := blockchain.GenerateChallenge(pubKey)

	signature, err := blockchain.SignActorWithStellarPrivateKey(
		deviceSeed,
		challenge,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access challenge: %w", err)
	}

	// -------------------------------------------------------------------------
	// 7. Request encrypted material from Cloud
	// -------------------------------------------------------------------------

	threadID := req.ThreadID
	if threadID == "" && shareEntry.Metadata != nil {
		threadID = shareEntry.Metadata["thread_id"]
		if threadID == "" {
			threadID = shareEntry.Metadata["threadID"]
		}
	}

	eventID := req.EventID
	if eventID == "" && shareEntry.Metadata != nil {
		eventID = shareEntry.Metadata["event_id"]
		if eventID == "" {
			eventID = shareEntry.Metadata["eventID"]
		}
	}

	sourceVaultID := shareEntry.CreatedBy
	if sourceVaultID == "" {
		sourceVaultID = req.CallerVaultID
	}

	var encryptedData []byte
	var encryptedKey string

	if accessGate, ok := u.assetResolver.(collaboration_ports.ThreadDataAccessGate); ok && accessGate != nil {
		response, err := accessGate.AccessThreadData(
			ctx,
			tracecore_types.ThreadDataAccessRequest{
				ThreadID:          threadID,
				EventID:           eventID,
				RequestingVaultID: req.CallerVaultID,
				Challenge:         challenge,
				Signature:         signature,
				TrustGroupID:      shareEntry.TrustGroupID,
				SourceVaultID:     sourceVaultID,
			},
		)
		if err != nil {
			return nil, fmt.Errorf(
				"cloud authorization failed: %w",
				err,
			)
		}

		if response != nil {
			encryptedKey = response.EncryptedKey

			if response.EncryptedPayload != "" {
				decoded, err := base64.StdEncoding.DecodeString(
					response.EncryptedPayload,
				)
				if err == nil && len(decoded) > 0 {
					encryptedData = decoded
				} else {
					encryptedData = []byte(response.EncryptedPayload)
				}
			}
		}
	}

	// Test/local fallback.
	if len(encryptedData) == 0 && u.assetResolver != nil {
		encryptedData, err = u.assetResolver.FetchEncryptedAsset(
			ctx,
			shareEntry.AssetCID,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to fetch asset content: %w",
				err,
			)
		}
	}

	if len(encryptedData) == 0 {
		return nil, errors.New(
			"cloud access gate returned no encrypted asset payload",
		)
	}

	// -------------------------------------------------------------------------
	// 8. Resolve Wrapped DEK
	// -------------------------------------------------------------------------

	dekSource := encryptedKey
	if dekSource == "" {
		dekSource = shareEntry.WrappedDEK
	}

	wrappedDEK, err := base64.StdEncoding.DecodeString(dekSource)
	if err != nil || len(wrappedDEK) == 0 {
		wrappedDEK = []byte(dekSource)
	}

	// -------------------------------------------------------------------------
	// 9. Decrypt locally
	// -------------------------------------------------------------------------

	cryptoResult, err := u.cryptoOrchestrator.ResolveCollaborativeAsset(
		ctx,
		trustgroup_orchestrator.ResolveCollaborativeAssetPayload{
			AssetID:       shareEntry.ID,
			TrustGroupID:  shareEntry.TrustGroupID,
			KEKVersion:    shareEntry.KEKVersion,
			EncryptedData: encryptedData,
			WrappedDEK:    wrappedDEK,
			WrappedKEK:    activeEnvelope.WrappedKEK,
			DeviceSeed:    deviceSeed,
			Keyring:       keyring,
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"cryptographic resolution failed: %w",
			err,
		)
	}

	// -------------------------------------------------------------------------
	// 10. Return plaintext
	// -------------------------------------------------------------------------

	return &collaboration_dtos.ResolveCollaborativeShareResponse{
		ShareEntryID: shareEntry.ID,
		TrustGroupID: shareEntry.TrustGroupID,
		CreatedBy:    shareEntry.CreatedBy,
		CreatedAt:    shareEntry.CreatedAt.Format(time.RFC3339),
		Metadata:     shareEntry.Metadata,
		Plaintext:    cryptoResult.Plaintext,
	}, nil
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
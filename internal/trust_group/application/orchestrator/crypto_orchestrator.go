package trustgroup_orchestrator

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	"vault-app/internal/utils"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
)

type ActiveDevice struct {
	DeviceID  string
	MemberID  string
	PublicKey string
	IsActive  bool
}

type PrepareCollaborativeAssetPayload struct {
	AssetID       string
	TrustGroupID  string
	KEKVersion    uint64
	RawPayload    []byte
	ActiveDevices []ActiveDevice
	Keyring       *vaults_domain.VaultKeyring
}

type PreparedCollaborativeAsset struct {
	AssetID       string
	TrustGroupID  string
	KEKVersion    uint64
	EncryptedData []byte                                            // AES-256-GCM(Payload, DEK)
	WrappedDEK    []byte                                            // AES-256-GCM(DEK, KEK)
	Envelopes     []trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest // WrappedKEK envelopes per active device
}

type TrustGroupCryptoOrchestrator struct {
	keyringService vault_infrastructure_security.KeyringServiceInterface
	aesService     *vault_infrastructure_crypto.AESService
	asymService    *vault_infrastructure_crypto.AsymmetricService
}

func NewTrustGroupCryptoOrchestrator(
	keyringService vault_infrastructure_security.KeyringServiceInterface,
	aesService *vault_infrastructure_crypto.AESService,
	asymService *vault_infrastructure_crypto.AsymmetricService,
) *TrustGroupCryptoOrchestrator {
	if aesService == nil {
		aesService = &vault_infrastructure_crypto.AESService{}
	}
	if asymService == nil {
		asymService = &vault_infrastructure_crypto.AsymmetricService{}
	}
	return &TrustGroupCryptoOrchestrator{
		keyringService: keyringService,
		aesService:     aesService,
		asymService:    asymService,
	}
}

func (o *TrustGroupCryptoOrchestrator) PrepareCollaborativeAsset(
	ctx context.Context,
	req PrepareCollaborativeAssetPayload,
) (*PreparedCollaborativeAsset, error) {
	if req.TrustGroupID == "" {
		return nil, errors.New("trust group ID is required")
	}
	if req.KEKVersion == 0 {
		return nil, errors.New("KEK version is required")
	}
	if len(req.RawPayload) == 0 {
		return nil, errors.New("raw payload cannot be empty")
	}

	if req.Keyring == nil {
		return nil, errors.New("vault keyring is required for collaborative asset preparation")
	}

	if o.keyringService == nil {
		return nil, errors.New("keyring service is not initialized")
	}

	kek, err := o.keyringService.GetTrustGroupKEK(
		req.Keyring,
		req.TrustGroupID,
		req.KEKVersion,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve trust group KEK %s v%d: %w",
			req.TrustGroupID,
			req.KEKVersion,
			err,
		)
	}

	if len(kek) != 32 {
		return nil, fmt.Errorf(
			"invalid trust group KEK length: got %d, want 32",
			len(kek),
		)
	}

	fmt.Printf(
		"[C3][CRYPTO][WRITE][KEK_RESOLVED] trustGroupID=%s kekVersion=%d keyLen=%d\n",
		req.TrustGroupID,
		req.KEKVersion,
		len(kek),
	)
	utils.LogPretty(
		"TrustGroupCryptoOrchestrator - PrepareCollaborativeAsset - kek",
		utils.FingerprintKey(kek),
	)

	sum := sha256.Sum256(kek)

	fmt.Printf(
		"[C3][CRYPTO][WRITE] trustGroupID=%s kekVersion=%d kekFingerprint=%x\n",
		req.TrustGroupID,
		req.KEKVersion,
		sum[:8],
	)

	// 2. Generate Asset DEK (32 bytes)
	dek := o.asymService.GenerateSymmetricKey()

	// 3. Encrypt raw payload using DEK (AES-256-GCM)
	encryptedData, err := o.aesService.Encrypt(req.RawPayload, dek)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt asset payload: %w", err)
	}

	// 4. Wrap DEK with TrustGroup KEK (AES-256-GCM)
	wrappedDEK, err := o.aesService.Encrypt(dek, kek)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap DEK with KEK: %w", err)
	}
	utils.LogPretty("TrustGroupCryptoOrchestrator - PrepareCollaborativeAsset - wrappedDEK", utils.FingerprintKey(wrappedDEK))

	// 5. Wrap KEK per active member using Member public key
	envelopes := make([]trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest, 0)
	seenMembers := make(map[string]bool)
	for _, dev := range req.ActiveDevices {
		if !dev.IsActive || seenMembers[dev.MemberID] {
			continue
		}
		if dev.PublicKey == "" || dev.MemberID == "" {
			continue
		}
		seenMembers[dev.MemberID] = true

		wrappedKEKPayload, err := o.aesService.EncryptPayload(dev.PublicKey, kek)
		if err != nil {
			return nil, fmt.Errorf("failed to wrap KEK for member %s: %w", dev.MemberID, err)
		}

		envelopes = append(envelopes, trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			TrustGroupID: req.TrustGroupID,
			MemberID:     dev.MemberID,
			KEKVersion:   req.KEKVersion,
			WrappedKEK:   wrappedKEKPayload.ToString(),
		})
	}

	return &PreparedCollaborativeAsset{
		AssetID:       req.AssetID,
		TrustGroupID:  req.TrustGroupID,
		KEKVersion:    req.KEKVersion,
		EncryptedData: encryptedData,
		WrappedDEK:    wrappedDEK,
		Envelopes:     envelopes,
	}, nil
}

type ResolveCollaborativeAssetPayload struct {
	AssetID       string
	TrustGroupID  string
	KEKVersion    int
	EncryptedData []byte
	WrappedDEK    []byte
	WrappedKEK    string

	// Caller/member private key material used to unwrap WrappedKEK.
	PrivateKey string
}

type ResolvedCollaborativeAsset struct {
	AssetID      string
	TrustGroupID string
	KEKVersion   uint64
	Plaintext    []byte
}

func (o *TrustGroupCryptoOrchestrator) ResolveCollaborativeAsset(
	ctx context.Context,
	req ResolveCollaborativeAssetPayload,
) (*ResolvedCollaborativeAsset, error) {

	if req.TrustGroupID == "" {
		return nil, errors.New("trust group ID is required")
	}

	if req.KEKVersion == 0 {
		return nil, errors.New("KEK version is required")
	}

	if len(req.EncryptedData) == 0 {
		return nil, errors.New("encrypted data cannot be empty")
	}

	if len(req.WrappedDEK) == 0 {
		return nil, errors.New("wrapped DEK cannot be empty")
	}

	if req.WrappedKEK == "" {
		return nil, errors.New("wrapped KEK cannot be empty")
	}

	if req.PrivateKey == "" {
		return nil, errors.New("member private key is required")
	}

	// -------------------------------------------------------------------------
	// 1. Recover the TrustGroup KEK.
	//
	// WrappedKEK was created when the member was added:
	//
	//     WrappedKEK = EncryptPayload(memberPublicKey, KEK)
	//
	// Therefore the member's private key is required to recover KEK.
	// -------------------------------------------------------------------------
	kek, err := o.aesService.AsymetricDecrypt(
		req.PrivateKey,
		req.WrappedKEK,
	)
	utils.LogPretty("TrustGroupCryptoOrchestrator - ResolveCollaborativeAsset - kek", utils.FingerprintKey(kek))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to unwrap trust group KEK v%d: %w",
			req.KEKVersion,
			err,
		)
	}

	if len(kek) != 32 {
		return nil, fmt.Errorf(
			"unwrapped trust group KEK must be exactly 32 bytes, got %d",
			len(kek),
		)
	}

	// -------------------------------------------------------------------------
	// 2. Recover the asset DEK using the TrustGroup KEK.
	//
	//     WrappedDEK = Encrypt(DEK, KEK)
	// -------------------------------------------------------------------------
	dek, err := o.aesService.Decrypt(
		req.WrappedDEK,
		kek,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to unwrap DEK with trust group KEK v%d: %w",
			req.KEKVersion,
			err,
		)
	}

	if len(dek) != 32 {
		return nil, fmt.Errorf(
			"unwrapped DEK must be exactly 32 bytes, got %d",
			len(dek),
		)
	}

	// -------------------------------------------------------------------------
	// 3. Decrypt the actual asset using the DEK.
	// -------------------------------------------------------------------------
	plaintext, err := o.aesService.Decrypt(
		req.EncryptedData,
		dek,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to decrypt collaborative asset with DEK: %w",
			err,
		)
	}

	return &ResolvedCollaborativeAsset{
		AssetID:      req.AssetID,
		TrustGroupID: req.TrustGroupID,
		KEKVersion:   uint64(req.KEKVersion),
		Plaintext:    plaintext,
	}, nil
}

type RotateCollaborativeAssetInput struct {
	ShareEntryID string
	WrappedDEK   []byte // WrappedDEK under KEK vN
}

type RotateCollaborativeAssetOutput struct {
	ShareEntryID string
	ReWrappedDEK []byte // Re-wrapped DEK under KEK v(N+1)
}

type RotateTrustGroupKEKPayload struct {
	TrustGroupID  string
	OldVersion    uint64
	NewVersion    uint64
	ActiveDevices []ActiveDevice
	Assets        []RotateCollaborativeAssetInput
	Keyring       *vaults_domain.VaultKeyring
}

type RotateTrustGroupKEKResult struct {
	TrustGroupID  string
	NewVersion    uint64
	RotatedAssets []RotateCollaborativeAssetOutput
	NewEnvelopes  []trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest
}

func (o *TrustGroupCryptoOrchestrator) RotateTrustGroupKEK(
	ctx context.Context,
	req RotateTrustGroupKEKPayload,
) (*RotateTrustGroupKEKResult, error) {
	if req.TrustGroupID == "" {
		return nil, errors.New("trust group ID is required")
	}
	if req.OldVersion == 0 {
		return nil, errors.New("old version is required")
	}
	if req.NewVersion != req.OldVersion+1 {
		return nil, errors.New("new version must be old version + 1")
	}
	if req.Keyring == nil || o.keyringService == nil {
		return nil, errors.New("keyring and keyring service are required")
	}

	// 1. Retrieve old KEK vN from local VaultKeyring
	oldKEK, err := o.keyringService.GetTrustGroupKEK(req.Keyring, req.TrustGroupID, req.OldVersion)
	if err != nil || len(oldKEK) != 32 {
		return nil, fmt.Errorf("failed to retrieve old KEK v%d from keyring: %w", req.OldVersion, err)
	}

	// 2. Generate new KEK v(N+1) (32 bytes)
	newKEK := o.asymService.GenerateSymmetricKey()

	// 3. Store new KEK v(N+1) in local VaultKeyring
	_, err = o.keyringService.StoreTrustGroupKEK(req.Keyring, req.TrustGroupID, req.NewVersion, newKEK)
	if err != nil {
		return nil, fmt.Errorf("failed to store new KEK v%d in keyring: %w", req.NewVersion, err)
	}

	// 4. Re-wrap DEKs for every collaborative asset (KEK vN -> DEK -> KEK vN+1)
	rotatedAssets := make([]RotateCollaborativeAssetOutput, 0, len(req.Assets))
	for _, asset := range req.Assets {
		if len(asset.WrappedDEK) == 0 {
			continue
		}
		// Unwrap DEK using old KEK vN
		dek, err := o.aesService.Decrypt(asset.WrappedDEK, oldKEK)
		if err != nil {
			return nil, fmt.Errorf("failed to unwrap DEK for asset %s with old KEK v%d: %w", asset.ShareEntryID, req.OldVersion, err)
		}

		// Re-wrap DEK using new KEK v(N+1)
		reWrappedDEK, err := o.aesService.Encrypt(dek, newKEK)
		if err != nil {
			return nil, fmt.Errorf("failed to re-wrap DEK for asset %s with new KEK v%d: %w", asset.ShareEntryID, req.NewVersion, err)
		}

		rotatedAssets = append(rotatedAssets, RotateCollaborativeAssetOutput{
			ShareEntryID: asset.ShareEntryID,
			ReWrappedDEK: reWrappedDEK,
		})
	}

	// 5. Wrap new KEK v(N+1) for each remaining active member
	newEnvelopes := make([]trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest, 0)
	seenMembersRot := make(map[string]bool)
	for _, dev := range req.ActiveDevices {
		if !dev.IsActive || seenMembersRot[dev.MemberID] {
			continue
		}
		if dev.PublicKey == "" || dev.MemberID == "" {
			continue
		}
		seenMembersRot[dev.MemberID] = true

		wrappedKEKPayload, err := o.aesService.EncryptPayload(dev.PublicKey, newKEK)
		if err != nil {
			return nil, fmt.Errorf("failed to wrap new KEK for member %s: %w", dev.MemberID, err)
		}

		newEnvelopes = append(newEnvelopes, trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			TrustGroupID: req.TrustGroupID,
			MemberID:     dev.MemberID,
			KEKVersion:   req.NewVersion,
			WrappedKEK:   wrappedKEKPayload.ToString(),
		})
	}

	return &RotateTrustGroupKEKResult{
		TrustGroupID:  req.TrustGroupID,
		NewVersion:    req.NewVersion,
		RotatedAssets: rotatedAssets,
		NewEnvelopes:  newEnvelopes,
	}, nil
}

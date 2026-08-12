package trustgroup_orchestrator

import (
	"context"
	"errors"
	"fmt"

	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
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
	EncryptedData []byte                                              // AES-256-GCM(Payload, DEK)
	WrappedDEK    []byte                                              // AES-256-GCM(DEK, KEK)
	Envelopes     []trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest   // WrappedKEK envelopes per active device
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

	// 1. Resolve or Generate KEK for TrustGroup + KEKVersion
	var kek []byte
	if req.Keyring != nil && o.keyringService != nil {
		k, err := o.keyringService.GetTrustGroupKEK(req.Keyring, req.TrustGroupID, req.KEKVersion)
		if err == nil && len(k) == 32 {
			kek = k
		}
	}

	if len(kek) == 0 {
		kek = o.asymService.GenerateSymmetricKey()
		if req.Keyring != nil && o.keyringService != nil {
			_, _ = o.keyringService.StoreTrustGroupKEK(req.Keyring, req.TrustGroupID, req.KEKVersion, kek)
		}
	}

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

	// 5. Wrap KEK per active device using Device.PublicKey (nacl box anonymous seal)
	envelopes := make([]trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest, 0, len(req.ActiveDevices))
	for _, dev := range req.ActiveDevices {
		if !dev.IsActive {
			// Revoked/inactive devices MUST NOT receive key envelopes
			continue
		}
		if dev.PublicKey == "" || dev.DeviceID == "" || dev.MemberID == "" {
			continue
		}

		wrappedKEKPayload, err := o.aesService.EncryptPayload(dev.PublicKey, kek)
		if err != nil {
			return nil, fmt.Errorf("failed to wrap KEK for device %s: %w", dev.DeviceID, err)
		}

		envelopes = append(envelopes, trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			TrustGroupID: req.TrustGroupID,
			MemberID:     dev.MemberID,
			DeviceID:     dev.DeviceID,
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

package trustgroup_usecases_envelope

import (
	"context"
	"errors"
	"fmt"
	"strings"

	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_ports "vault-app/internal/trust_group/application/ports"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
)

type ProvisionTrustGroupDeviceEnvelopeUseCase struct {
	trustGroupRepo     trustgroup_domain.TrustGroupRepository
	deviceResolver     trustgroup_ports.DeviceResolver
	cryptoOrchestrator *trustgroup_orchestrator.TrustGroupCryptoOrchestrator
	addEnvelopeUseCase *AddTrustGroupKeyEnvelopeUseCase
	keyringService     vault_infrastructure_security.KeyringServiceInterface
	aesService         *vault_infrastructure_crypto.AESService
	asymService        *vault_infrastructure_crypto.AsymmetricService
}

func NewProvisionTrustGroupDeviceEnvelopeUseCase(
	trustGroupRepo trustgroup_domain.TrustGroupRepository,
	deviceResolver trustgroup_ports.DeviceResolver,
	cryptoOrchestrator *trustgroup_orchestrator.TrustGroupCryptoOrchestrator,
	addEnvelopeUseCase *AddTrustGroupKeyEnvelopeUseCase,
	keyringService vault_infrastructure_security.KeyringServiceInterface,
) *ProvisionTrustGroupDeviceEnvelopeUseCase {
	return &ProvisionTrustGroupDeviceEnvelopeUseCase{
		trustGroupRepo:     trustGroupRepo,
		deviceResolver:     deviceResolver,
		cryptoOrchestrator: cryptoOrchestrator,
		addEnvelopeUseCase: addEnvelopeUseCase,
		keyringService:     keyringService,
		aesService:         &vault_infrastructure_crypto.AESService{},
		asymService:        &vault_infrastructure_crypto.AsymmetricService{},
	}
}

func (uc *ProvisionTrustGroupDeviceEnvelopeUseCase) ValidateDependencies() error {
	if uc.trustGroupRepo == nil {
		return trustgroup_domain.ErrRepositoryNil
	}
	if uc.addEnvelopeUseCase == nil {
		return errors.New("add envelope use case is required")
	}
	return nil
}

func (uc *ProvisionTrustGroupDeviceEnvelopeUseCase) ValidateRequest(req trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest) error {
	if strings.TrimSpace(req.TrustGroupID) == "" {
		return trustgroup_domain.ErrTrustGroupIDRequired
	}
	if strings.TrimSpace(req.MemberID) == "" {
		return trustgroup_domain.ErrMemberIDRequired
	}
	if strings.TrimSpace(req.DeviceID) == "" {
		return trustgroup_domain.ErrDeviceIDRequired
	}
	return nil
}

func (uc *ProvisionTrustGroupDeviceEnvelopeUseCase) Execute(
	ctx context.Context,
	req trustgroup_dtos.ProvisionTrustGroupDeviceEnvelopeRequest,
	keyring *vaults_domain.VaultKeyring,
) (*trustgroup_domain.TrustGroup, error) {
	fmt.Printf("[C3][ADD_MEMBER][STEP_07] ProvisionTrustGroupDeviceEnvelopeUseCase.Execute enter trustGroupID=%s memberID=%s deviceID=%s hasPubKey=%t\n", req.TrustGroupID, req.MemberID, req.DeviceID, req.DevicePublicKey != "")

	if err := uc.ValidateDependencies(); err != nil {
		fmt.Printf("[C3][ADD_MEMBER][STEP_08] ValidateDependencies failed: %v\n", err)
		return nil, err
	}

	if err := uc.ValidateRequest(req); err != nil {
		fmt.Printf("[C3][ADD_MEMBER][STEP_08] ValidateRequest failed: %v\n", err)
		return nil, err
	}

	// 1. Fetch TrustGroup
	tgResp, err := uc.trustGroupRepo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{
		TrustGroupID: req.TrustGroupID,
	})
	if err != nil {
		fmt.Printf("[C3][ADD_MEMBER][STEP_08] GetTrustGroup failed: %v\n", err)
		return nil, fmt.Errorf("failed to fetch trust group %s: %w", req.TrustGroupID, err)
	}
	if tgResp == nil || tgResp.Data.ID == "" {
		fmt.Printf("[C3][ADD_MEMBER][STEP_08] GetTrustGroup returned empty\n")
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	tg := tgResp.Data
	fmt.Printf("[C3][ADD_MEMBER][STEP_09] GetTrustGroup fetched trustGroupID=%s kekVersion=%d envelopesCount=%d memberCIDsCount=%d\n", tg.ID, tg.KEKVersion, len(tg.KeyEnvelopes), len(tg.MemberCIDs))

	// 2. Verify Member belongs to TrustGroup
	memberFound := false
	for _, cid := range tg.MemberCIDs {
		if cid == req.MemberID {
			memberFound = true
			break
		}
	}
	fmt.Printf("[C3][ADD_MEMBER][STEP_10] MemberID membership validation memberFound=%t memberID=%s\n", memberFound, req.MemberID)
	if !memberFound {
		fmt.Printf("[C3][ADD_MEMBER][STEP_08] MemberID %s not in TrustGroup %s\n", req.MemberID, req.TrustGroupID)
		return nil, trustgroup_domain.ErrMemberNotInTrustGroup
	}

	// 2.1 Check if a non-revoked envelope already exists for (MemberID, DeviceID, KEKVersion)
	existingFound := false
	for _, env := range tg.KeyEnvelopes {
		if env.MemberID == req.MemberID && env.DeviceID == req.DeviceID && env.KEKVersion == tg.KEKVersion && env.RevokedAt == nil {
			existingFound = true
			break
		}
	}
	fmt.Printf("[C3][ADD_MEMBER][STEP_11] Current KEKVersion=%d\n", tg.KEKVersion)
	fmt.Printf("[C3][ADD_MEMBER][STEP_12] Existing active envelope check found=%t memberID=%s deviceID=%s\n", existingFound, req.MemberID, req.DeviceID)
	if existingFound {
		fmt.Printf("[C3][INVITE][ENVELOPE] Envelope already exists for trustGroupID=%s memberID=%s deviceID=%s kekVersion=%d (idempotent skip)\n", tg.ID, req.MemberID, req.DeviceID, tg.KEKVersion)
		return &tg, nil
	}

	// 3. Resolve Device Public Key (via DeviceResolver or req.DevicePublicKey)
	pubKeySource := "req.DevicePublicKey"
	targetPubKey := strings.TrimSpace(req.DevicePublicKey)
	if uc.deviceResolver != nil {
		dev, devErr := uc.deviceResolver.GetDevice(ctx, req.DeviceID)
		if devErr == nil && dev != nil && dev.IsActive {
			if dev.PublicKey != "" {
				targetPubKey = dev.PublicKey
				pubKeySource = "deviceResolver"
			}
		}
	}
	fmt.Printf("[C3][ADD_MEMBER][STEP_13] Resolved DevicePublicKey source=%s pubKeyLen=%d\n", pubKeySource, len(targetPubKey))
	if targetPubKey == "" {
		fmt.Printf("[C3][ADD_MEMBER][STEP_08] Device public key is empty\n")
		return nil, errors.New("device public key is required for key envelope provisioning")
	}

	// 4. Resolve current Group KEK (v1) inside sovereign crypto boundary
	var kek []byte
	if keyring != nil && uc.keyringService != nil {
		k, kErr := uc.keyringService.GetTrustGroupKEK(keyring, tg.ID, tg.KEKVersion)
		if kErr == nil && len(k) == 32 {
			kek = k
		}
	}

	if len(kek) == 0 {
		// Generate or initialize KEK if not cached, and cache it if keyring is present
		kek = uc.asymService.GenerateSymmetricKey()
		if keyring != nil && uc.keyringService != nil {
			_, _ = uc.keyringService.StoreTrustGroupKEK(keyring, tg.ID, tg.KEKVersion, kek)
		}
	}
	fmt.Printf("[C3][ADD_MEMBER][STEP_14] KEK resolution success=%t kekLen=%d\n", len(kek) == 32, len(kek))

	// 5. Wrap KEK using target device public key (asymmetric box seal)
	wrappedKEKPayload, err := uc.aesService.EncryptPayload(targetPubKey, kek)
	if err != nil {
		fmt.Printf("[C3][ADD_MEMBER][STEP_08] EncryptPayload failed: %v\n", err)
		return nil, fmt.Errorf("failed to wrap KEK for device %s: %w", req.DeviceID, err)
	}

	fmt.Printf("[C3][ENVELOPE][CLIENT][INPUT] trustGroupID=%s memberID=%s deviceID=%s kekVersion=%d publicKeyPresent=%t\n", tg.ID, req.MemberID, req.DeviceID, tg.KEKVersion, targetPubKey != "")
	fmt.Printf("[C3][ENVELOPE][CLIENT][CREATED] memberID=%s deviceID=%s kekVersion=%d wrappedKEKPresent=%t wrappedKEKLen=%d\n", req.MemberID, req.DeviceID, tg.KEKVersion, wrappedKEKPayload.ToString() != "", len(wrappedKEKPayload.ToString()))

	// 6. Delegate envelope attachment to AddTrustGroupKeyEnvelopeUseCase
	fmt.Printf("[C3][ADD_MEMBER][STEP_15] Calling AddTrustGroupKeyEnvelopeUseCase.Execute trustGroupID=%s memberID=%s deviceID=%s kekVersion=%d\n", tg.ID, req.MemberID, req.DeviceID, tg.KEKVersion)
	updatedTg, err := uc.addEnvelopeUseCase.Execute(ctx, trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
		TrustGroupID: tg.ID,
		MemberID:     req.MemberID,
		DeviceID:     req.DeviceID,
		KEKVersion:   tg.KEKVersion,
		WrappedKEK:   wrappedKEKPayload.ToString(),
	})
	fmt.Printf("[C3][ADD_MEMBER][STEP_16] AddTrustGroupKeyEnvelopeUseCase.Execute returned err=%v\n", err)
	if err != nil {
		return nil, fmt.Errorf("failed to attach device key envelope: %w", err)
	}

	if updatedTg != nil {
		fmt.Printf("[C3][ADD_MEMBER][STEP_17] Resulting TrustGroup.KeyEnvelopes count=%d\n", len(updatedTg.KeyEnvelopes))
	}

	return updatedTg, nil
}

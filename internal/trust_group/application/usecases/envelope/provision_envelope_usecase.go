package trustgroup_usecases_envelope

import (
	"context"
	"errors"
	"fmt"
	"strings"

	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_orchestrator "vault-app/internal/trust_group/application/orchestrator"
	trustgroup_domain "vault-app/internal/trust_group/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
)

type ProvisionTrustGroupMemberEnvelopeUseCase struct {
	trustGroupRepo     trustgroup_domain.TrustGroupRepository
	cryptoOrchestrator *trustgroup_orchestrator.TrustGroupCryptoOrchestrator
	addEnvelopeUseCase *AddTrustGroupKeyEnvelopeUseCase
	keyringService     vault_infrastructure_security.KeyringServiceInterface
	aesService         *vault_infrastructure_crypto.AESService
	asymService        *vault_infrastructure_crypto.AsymmetricService
}

type ProvisionTrustGroupDeviceEnvelopeUseCase = ProvisionTrustGroupMemberEnvelopeUseCase

func NewProvisionTrustGroupMemberEnvelopeUseCase(
	trustGroupRepo trustgroup_domain.TrustGroupRepository,
	cryptoOrchestrator *trustgroup_orchestrator.TrustGroupCryptoOrchestrator,
	addEnvelopeUseCase *AddTrustGroupKeyEnvelopeUseCase,
	keyringService vault_infrastructure_security.KeyringServiceInterface,
) *ProvisionTrustGroupMemberEnvelopeUseCase {
	return &ProvisionTrustGroupMemberEnvelopeUseCase{
		trustGroupRepo:     trustGroupRepo,
		cryptoOrchestrator: cryptoOrchestrator,
		addEnvelopeUseCase: addEnvelopeUseCase,
		keyringService:     keyringService,
		aesService:         &vault_infrastructure_crypto.AESService{},
		asymService:        &vault_infrastructure_crypto.AsymmetricService{},
	}
}

func NewProvisionTrustGroupDeviceEnvelopeUseCase(
	trustGroupRepo trustgroup_domain.TrustGroupRepository,
	args ...interface{},
) *ProvisionTrustGroupMemberEnvelopeUseCase {
	var cryptoOrchestrator *trustgroup_orchestrator.TrustGroupCryptoOrchestrator
	var addEnvelopeUseCase *AddTrustGroupKeyEnvelopeUseCase
	var keyringService vault_infrastructure_security.KeyringServiceInterface

	for _, arg := range args {
		switch v := arg.(type) {
		case *trustgroup_orchestrator.TrustGroupCryptoOrchestrator:
			cryptoOrchestrator = v
		case *AddTrustGroupKeyEnvelopeUseCase:
			addEnvelopeUseCase = v
		case vault_infrastructure_security.KeyringServiceInterface:
			keyringService = v
		}
	}
	return NewProvisionTrustGroupMemberEnvelopeUseCase(trustGroupRepo, cryptoOrchestrator, addEnvelopeUseCase, keyringService)
}

func (uc *ProvisionTrustGroupMemberEnvelopeUseCase) Execute(
	ctx context.Context,
	req trustgroup_dtos.ProvisionTrustGroupMemberEnvelopeRequest,
	keyring *vaults_domain.VaultKeyring,
) (*trustgroup_domain.TrustGroup, error) {
	pubKeyPresent := req.MemberPublicKey != "" || req.DevicePublicKey != ""
	fmt.Printf("[C3][ADD_MEMBER][PROVISION_MEMBER] ProvisionTrustGroupMemberEnvelopeUseCase.Execute enter trustGroupID=%s memberID=%s hasPubKey=%t\n", req.TrustGroupID, req.MemberID, pubKeyPresent)

	if err := uc.ValidateDependencies(); err != nil {
		fmt.Printf("[C3][ADD_MEMBER][PROVISION_MEMBER] ValidateDependencies failed: %v\n", err)
		return nil, err
	}

	if err := uc.ValidateRequest(req); err != nil {
		fmt.Printf("[C3][ADD_MEMBER][PROVISION_MEMBER] ValidateRequest failed: %v\n", err)
		return nil, err
	}

	// 1. Fetch TrustGroup
	tgResp, err := uc.trustGroupRepo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{
		TrustGroupID: req.TrustGroupID,
	})
	if err != nil {
		fmt.Printf("[C3][ADD_MEMBER][PROVISION_MEMBER] GetTrustGroup failed: %v\n", err)
		return nil, fmt.Errorf("failed to fetch trust group %s: %w", req.TrustGroupID, err)
	}
	if tgResp == nil || tgResp.Data.ID == "" {
		fmt.Printf("[C3][ADD_MEMBER][PROVISION_MEMBER] GetTrustGroup returned empty\n")
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	tg := tgResp.Data
	fmt.Printf("[C3][ENVELOPE_PROVISION_TRACE][04_PROVISION_CALL] trustGroupID=%s memberID=%s publicKeyPresent=%t currentKEKVersion=%d\n",
		req.TrustGroupID, req.MemberID, pubKeyPresent, tg.KEKVersion)

	// 2. Verify Member belongs to TrustGroup
	memberFound := false
	for _, cid := range tg.MemberCIDs {
		if cid == req.MemberID {
			memberFound = true
			break
		}
	}
	if !memberFound {
		fmt.Printf("[C3][ADD_MEMBER][PROVISION_MEMBER] MemberID %s not in TrustGroup %s\n", req.MemberID, req.TrustGroupID)
		return nil, trustgroup_domain.ErrMemberNotInTrustGroup
	}

	// 2.1 Check if a non-revoked envelope already exists for (MemberID, KEKVersion)
	existingFound := false
	for _, env := range tg.KeyEnvelopes {
		if env.MemberID == req.MemberID && env.KEKVersion == tg.KEKVersion && env.RevokedAt == nil {
			existingFound = true
			break
		}
	}
	if existingFound {
		fmt.Printf("[C3][INVITE][ENVELOPE] Envelope already exists for trustGroupID=%s memberID=%s kekVersion=%d (idempotent skip)\n", tg.ID, req.MemberID, tg.KEKVersion)
		return &tg, nil
	}

	// 3. Resolve Member Public Key
	targetPubKey := strings.TrimSpace(req.MemberPublicKey)
	if targetPubKey == "" {
		targetPubKey = strings.TrimSpace(req.DevicePublicKey)
	}
	if targetPubKey == "" {
		return nil, errors.New("member public key is required for key envelope provisioning")
	}

	// 4. Resolve current Group KEK from the sovereign keyring.
	if keyring == nil {
		return nil, errors.New("vault keyring is required for key envelope provisioning")
	}

	if uc.keyringService == nil {
		return nil, errors.New("keyring service is not initialized")
	}

	kek, err := uc.keyringService.GetTrustGroupKEK(
		keyring,
		tg.ID,
		tg.KEKVersion,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve trust group KEK %s v%d: %w",
			tg.ID,
			tg.KEKVersion,
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
		"[C3][ENVELOPE][KEK_RESOLVED] trustGroupID=%s memberID=%s kekVersion=%d keyLen=%d\n",
		tg.ID,
		req.MemberID,
		tg.KEKVersion,
		len(kek),
	)

	// 5. Wrap the existing TrustGroup KEK with the member's public key.
	wrappedKEKPayload, err := uc.aesService.EncryptPayload(
		targetPubKey,
		kek,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to wrap trust group KEK for member %s: %w",
			req.MemberID,
			err,
		)
	}

	fmt.Printf(
		"[C3][ENVELOPE][GENERATED] trustGroupID=%s memberID=%s kekVersion=%d\n",
		tg.ID,
		req.MemberID,
		tg.KEKVersion,
	)

	// 6. Delegate envelope attachment to AddTrustGroupKeyEnvelopeUseCase.
	updatedTg, err := uc.addEnvelopeUseCase.Execute(
		ctx,
		trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest{
			TrustGroupID: tg.ID,
			MemberID:     req.MemberID,
			KEKVersion:   tg.KEKVersion,
			WrappedKEK:   wrappedKEKPayload.ToString(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to attach member key envelope: %w",
			err,
		)
	}

	return updatedTg, nil
}

func (uc *ProvisionTrustGroupMemberEnvelopeUseCase) ValidateDependencies() error {
	if uc.trustGroupRepo == nil {
		return trustgroup_domain.ErrRepositoryNil
	}
	if uc.addEnvelopeUseCase == nil {
		return errors.New("add envelope use case is required")
	}
	return nil
}

func (uc *ProvisionTrustGroupMemberEnvelopeUseCase) ValidateRequest(req trustgroup_dtos.ProvisionTrustGroupMemberEnvelopeRequest) error {
	if strings.TrimSpace(req.TrustGroupID) == "" {
		return trustgroup_domain.ErrTrustGroupIDRequired
	}
	if strings.TrimSpace(req.MemberID) == "" {
		return trustgroup_domain.ErrMemberIDRequired
	}
	pubKey := strings.TrimSpace(req.MemberPublicKey)
	if pubKey == "" {
		pubKey = strings.TrimSpace(req.DevicePublicKey)
	}
	if pubKey == "" {
		return errors.New("member public key is required for key envelope provisioning")
	}
	return nil
}

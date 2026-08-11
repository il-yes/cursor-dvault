package trustgroup_usecases

import (
	"context"
	"errors"
	"strings"

	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_ports "vault-app/internal/trust_group/application/ports"
	trustgroup_domain "vault-app/internal/trust_group/domain"
)

type AddTrustGroupKeyEnvelopeUseCase struct {
	repo           trustgroup_domain.TrustGroupRepository
	deviceResolver trustgroup_ports.DeviceResolver
}

func NewAddTrustGroupKeyEnvelopeUseCase(
	repo trustgroup_domain.TrustGroupRepository,
	deviceResolver trustgroup_ports.DeviceResolver,
) *AddTrustGroupKeyEnvelopeUseCase {
	return &AddTrustGroupKeyEnvelopeUseCase{
		repo:           repo,
		deviceResolver: deviceResolver,
	}
}

func (uc *AddTrustGroupKeyEnvelopeUseCase) ValidateDependencies() error {
	if uc.repo == nil {
		return trustgroup_domain.ErrRepositoryNil
	}
	if uc.deviceResolver == nil {
		return errors.New("device resolver is required")
	}
	return nil
}

func (uc *AddTrustGroupKeyEnvelopeUseCase) ValidateRequest(req trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest) error {
	if strings.TrimSpace(req.TrustGroupID) == "" {
		return trustgroup_domain.ErrTrustGroupIDRequired
	}
	if strings.TrimSpace(req.MemberID) == "" {
		return trustgroup_domain.ErrMemberIDRequired
	}
	if strings.TrimSpace(req.DeviceID) == "" {
		return trustgroup_domain.ErrDeviceIDRequired
	}
	if strings.TrimSpace(req.WrappedKEK) == "" {
		return trustgroup_domain.ErrWrappedKEKRequired
	}
	if req.KEKVersion == 0 {
		return trustgroup_domain.ErrKEKVersionRequired
	}
	return nil
}

func (uc *AddTrustGroupKeyEnvelopeUseCase) Execute(
	ctx context.Context,
	req trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest,
) (*trustgroup_domain.TrustGroup, error) {
	if err := uc.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := uc.ValidateRequest(req); err != nil {
		return nil, err
	}

	// 1. Resolve TrustGroup
	resp, err := uc.repo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{
		TrustGroupID: req.TrustGroupID,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.Data.ID == "" {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}

	tg := &resp.Data

	// 2. Validate KEKVersion matches current TrustGroup version
	if req.KEKVersion != tg.KEKVersion {
		return nil, trustgroup_domain.ErrStaleKEKVersion
	}

	// 3. Verify Member belongs to TrustGroup
	memberFound := false
	for _, cid := range tg.MemberCIDs {
		if cid == req.MemberID {
			memberFound = true
			break
		}
	}
	if !memberFound {
		return nil, trustgroup_domain.ErrMemberNotInTrustGroup
	}

	// 4. Resolve & Validate Device via DeviceResolver port
	device, err := uc.deviceResolver.GetDevice(ctx, req.DeviceID)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, trustgroup_domain.ErrDeviceNotFound
	}
	if !device.IsActive {
		return nil, trustgroup_domain.ErrDeviceRevoked
	}
	if device.VaultID != req.MemberID {
		return nil, trustgroup_domain.ErrDeviceMemberMismatch
	}

	// 5. Add Key Envelope to TrustGroup aggregate
	envelope := trustgroup_domain.TrustGroupKeyEnvelope{
		TrustGroupID: tg.ID,
		MemberID:     req.MemberID,
		DeviceID:     req.DeviceID,
		KEKVersion:   req.KEKVersion,
		WrappedKEK:   req.WrappedKEK,
	}

	if err := tg.AddEnvelope(envelope); err != nil {
		return nil, err
	}

	// 6. Persist updated TrustGroup
	updatedResp, err := uc.repo.UpdateTrustGroup(ctx, &trustgroup_domain.UpdateTrustGroupRequest{
		TrustGroup: *tg,
	})
	if err != nil {
		return nil, err
	}
	if updatedResp == nil {
		return nil, trustgroup_domain.ErrRepositoryResponse
	}

	return &updatedResp.Data, nil
}

package collaboration_usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"

	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_usecases "vault-app/internal/trust_group/application/usecases/envelope"
)

type CreateCollaborativeShareUseCase struct {
	shareAssetUseCase  *ShareAssetWithTrustGroupUsecase
	addEnvelopeUseCase *trustgroup_usecases.AddTrustGroupKeyEnvelopeUseCase
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

func (u *CreateCollaborativeShareUseCase) ValidateDependencies() error {
	if u.shareAssetUseCase == nil {
		return errors.New("share asset use case is required")
	}
	if u.addEnvelopeUseCase == nil {
		return errors.New("add envelope use case is required")
	}
	return nil
}

func (u *CreateCollaborativeShareUseCase) ValidateRequest(req collaboration_dtos.CreateCollaborativeShareRequest) error {
	if strings.TrimSpace(req.TrustGroupID) == "" {
		return errors.New("trust group id is required")
	}
	if req.KEKVersion == 0 {
		return errors.New("kek version is required")
	}
	if strings.TrimSpace(req.CreatedBy) == "" {
		return errors.New("created by is required")
	}
	if strings.TrimSpace(req.AssetCID) == "" {
		return errors.New("asset cid is required")
	}
	if strings.TrimSpace(req.WrappedDEK) == "" {
		return errors.New("wrapped dek is required")
	}
	return nil
}

func (u *CreateCollaborativeShareUseCase) Execute(
	ctx context.Context,
	req collaboration_dtos.CreateCollaborativeShareRequest,
) (*collaboration_dtos.CreateCollaborativeShareResponse, error) {
	if err := u.ValidateDependencies(); err != nil {
		return nil, err
	}
	if err := u.ValidateRequest(req); err != nil {
		return nil, err
	}

	// 1. Create and persist ShareEntry via ShareAssetWithTrustGroupUsecase
	shareEntry, err := u.shareAssetUseCase.Execute(ctx, collaboration_dtos.ShareAssetWithTrustGroupRequest{
		AssetCID:     req.AssetCID,
		TrustGroupID: req.TrustGroupID,
		WrappedDEK:   req.WrappedDEK,
		KEKVersion:   req.KEKVersion,
		CreatedBy:    req.CreatedBy,
		Metadata:     req.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create collaborative share entry: %w", err)
	}

	// 2. Attach device key envelopes to TrustGroup
	attachedEnvelopes := make([]trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest, 0, len(req.Envelopes))
	for _, envReq := range req.Envelopes {
		// Ensure KEKVersion and TrustGroupID are populated from the request
		if envReq.TrustGroupID == "" {
			envReq.TrustGroupID = req.TrustGroupID
		}
		if envReq.KEKVersion == 0 {
			envReq.KEKVersion = req.KEKVersion
		}

		_, err := u.addEnvelopeUseCase.Execute(ctx, envReq)
		if err != nil {
			return nil, fmt.Errorf("failed to attach key envelope for device %s: %w", envReq.DeviceID, err)
		}
		attachedEnvelopes = append(attachedEnvelopes, envReq)
	}

	return &collaboration_dtos.CreateCollaborativeShareResponse{
		ShareEntry: *shareEntry,
		Envelopes:  attachedEnvelopes,
	}, nil
}

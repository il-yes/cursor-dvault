package collaboration_usecases

import (
	"context"
	"errors"
	"strings"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_dtos "vault-app/internal/collaboration/application/dtos"
	trustgroup_domain "vault-app/internal/trust_group/domain"
)

type ShareAssetWithTrustGroupUsecase struct {
	trustGroupRepo trustgroup_domain.TrustGroupRepository
	shareEntryRepo c3_asset_domain.ShareEntryRepository
}

func NewShareAssetWithTrustGroupUsecase(
	trustGroupRepo trustgroup_domain.TrustGroupRepository,
	shareEntryRepo c3_asset_domain.ShareEntryRepository,
) *ShareAssetWithTrustGroupUsecase {
	return &ShareAssetWithTrustGroupUsecase{
		trustGroupRepo: trustGroupRepo,
		shareEntryRepo: shareEntryRepo,
	}
}

func (u *ShareAssetWithTrustGroupUsecase) ValidateDependencies() error {
	if u.trustGroupRepo == nil {
		return trustgroup_domain.ErrRepositoryNil
	}

	if u.shareEntryRepo == nil {
		return c3_asset_domain.ErrRepositoryNil
	}

	return nil
}

func (u *ShareAssetWithTrustGroupUsecase) ValidateRequest(
	req collaboration_dtos.ShareAssetWithTrustGroupRequest,
) error {
	if strings.TrimSpace(req.AssetCID) == "" {
		return errors.New("asset cid is required")
	}

	if strings.TrimSpace(req.TrustGroupID) == "" {
		return errors.New("trust group id is required")
	}

	if strings.TrimSpace(req.WrappedDEK) == "" {
		return errors.New("wrapped dek is required")
	}

	if req.KEKVersion == 0 {
		return errors.New("kek version is required")
	}

	if strings.TrimSpace(req.CreatedBy) == "" {
		return errors.New("created by is required")
	}

	return nil
}

func (u *ShareAssetWithTrustGroupUsecase) Execute(
	ctx context.Context,
	req collaboration_dtos.ShareAssetWithTrustGroupRequest,
) (*c3_asset_domain.ShareEntry, error) {

	if err := u.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := u.ValidateRequest(req); err != nil {
		return nil, err
	}

	// ------------------------------------------------------------------
	// 1. Load TrustGroup
	// ------------------------------------------------------------------

	group, err := u.trustGroupRepo.GetTrustGroup(
		ctx,
		&trustgroup_domain.GetTrustGroupRequest{
			TrustGroupID: req.TrustGroupID,
		},
	)
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}

	// ------------------------------------------------------------------
	// 2. Validate KEKVersion matches current TrustGroup KEKVersion
	// ------------------------------------------------------------------

	if req.KEKVersion != group.Data.KEKVersion {
		return nil, trustgroup_domain.ErrStaleKEKVersion
	}

	// ------------------------------------------------------------------
	// 3. Create ShareEntry
	// ------------------------------------------------------------------

	shareEntry, err := c3_asset_domain.NewShareEntry(
		req.AssetCID,
		req.TrustGroupID,
		req.WrappedDEK,
		req.KEKVersion,
		req.CreatedBy,
		req.Metadata,
	)
	if err != nil {
		return nil, err
	}

	// ------------------------------------------------------------------
	// 4. Persist ShareEntry
	// ------------------------------------------------------------------

	created, err := u.shareEntryRepo.CreateShareEntry(
		ctx,
		&c3_asset_domain.CreateShareEntryRequest{
			ShareEntry: shareEntry,
		},
	)
	if err != nil {
		return nil, err
	}

	return &created.Data, nil
}


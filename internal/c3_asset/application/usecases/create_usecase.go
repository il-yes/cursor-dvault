package c3_asset_usecases

import (
	"context"
	"errors"
	"strings"

	c3_asset_dtos "vault-app/internal/c3_asset/application/dtos"
	c3_asset_domain "vault-app/internal/c3_asset/domain"
)

type CreateShareEntryUsecase struct {
	repo c3_asset_domain.ShareEntryRepository
}

func NewCreateShareEntryUsecase(
	repo c3_asset_domain.ShareEntryRepository,
) *CreateShareEntryUsecase {
	return &CreateShareEntryUsecase{
		repo: repo,
	}
}

func (u *CreateShareEntryUsecase) ValidateDependencies() error {
	if u.repo == nil {
		return errors.New("share entry repository is required")
	}

	return nil
}

func (u *CreateShareEntryUsecase) ValidateRequest(
	req c3_asset_dtos.CreateShareEntryRequest,
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

	if strings.TrimSpace(req.CreatedBy) == "" {
		return errors.New("created by is required")
	}

	return nil
}

func (u *CreateShareEntryUsecase) Execute(
	ctx context.Context,
	req c3_asset_dtos.CreateShareEntryRequest,
) (*c3_asset_domain.ShareEntry, error) {

	if err := u.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := u.ValidateRequest(req); err != nil {
		return nil, err
	}

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
	created, err := u.repo.CreateShareEntry(
		ctx,
		&c3_asset_domain.CreateShareEntryRequest{
			ShareEntry: shareEntry,
		},
	)
	if err != nil {
		return nil, err
	}
	if created == nil {
		return nil, errors.New("share entry repository returned nil response")
	}

	return &created.Data, nil
}

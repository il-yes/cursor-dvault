package share_application_use_cases

import (
	"context"
	"fmt"
	utils "vault-app/internal"
	share_application_dto "vault-app/internal/application"
	share_application_events "vault-app/internal/application/events/share"
	share_domain "vault-app/internal/domain/shared"
	"vault-app/internal/tracecore"
)


// ---------------------------------------------------------
//	Interfaces
// ---------------------------------------------------------
type AnkhoraClientInterface interface {
	CreateLinkShare(ctx context.Context, req share_application_dto.LinkShareCreateRequest) (*tracecore.CreateLinkShareResponse, error)
	ListLinkSharesByMe(ctx context.Context, userID string) (*tracecore.LinkShareResponse, error)
	ListLinkSharesWithMe(ctx context.Context, userID string) (*tracecore.LinkShareResponse, error)
}


// ---------------------------------------------------------
//
//	Link Share Use Case
//
// ---------------------------------------------------------
type LinkShareUseCase struct {
	repo       share_domain.Repository
	dispatcher share_application_events.EventDispatcher
	tc         AnkhoraClientInterface // new cloud client
	crypto     ClientCryptoService
}

func NewLinkShareUseCase(repo share_domain.Repository, tc AnkhoraClientInterface, d share_application_events.EventDispatcher, crypto ClientCryptoService) *LinkShareUseCase {
	return &LinkShareUseCase{
		repo:       repo,
		dispatcher: d,
		tc:         tc,
		crypto:     crypto,
	}
}

func (uc *LinkShareUseCase) CreateLinkShare(ctx context.Context, email string, req share_application_dto.LinkShareCreateRequest) (*share_domain.LinkShare, error) {
	// ----------------------------------------------------------
	// 1. Validate share request
	// ----------------------------------------------------------
	utils.LogPretty("LinkShareUseCase - CreateLinkShare", req)
	if req.CreatorEmail == "" || req.Title == "" || req.EntryType == "" || req.Payload == "" {
		return nil, fmt.Errorf("invalid share request")
	}	

	// ----------------------------------------------------------
	// 2. Send to Ankhora Cloud
	// ----------------------------------------------------------
	res, err := uc.tc.CreateLinkShare(ctx, req)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("LinkShareUseCase - CreateLinkShare", res)
	// ----------------------------------------------------------
	// 3. Create Link share
	// ----------------------------------------------------------
	share := share_domain.LinkShare{
		Payload: req.Payload,
		ExpiresAt: req.ExpiresAt,
		MaxViews: req.MaxViews,
		Password: req.Password,
		DownloadAllowed: req.DownloadAllowed,
		CreatorEmail: req.CreatorEmail,
		Metadata: share_domain.Metadata{
			EntryType: req.EntryType,
			Title: req.Title,
		},
	}

	return &share, nil
}

func (uc *LinkShareUseCase) ListLinkSharesByMe(ctx context.Context, email string) (*[]tracecore.WailsLinkShare, error) {
	// ----------------------------------------------------------
	// 1. Fetch from Ankhora Cloud
	// ----------------------------------------------------------
	res, err := uc.tc.ListLinkSharesByMe(ctx, email)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("LinkShareUseCase - ListLinkSharesByMe", res)
    shares := res.Data

	// ----------------------------------------------------------
	// 2. Map to domain
	// ----------------------------------------------------------
	var domainShares []tracecore.WailsLinkShare
	for _, share := range shares {
		domainShares = append(domainShares, *share.ToWailsLinkShare())
	}
	utils.LogPretty("LinkShareUseCase - ListLinkSharesByMe", domainShares)

	return &domainShares, nil
}

func (uc *LinkShareUseCase) ListLinkSharesWithMe(ctx context.Context, email string) (*[]tracecore.WailsLinkShare, error) {
	// ----------------------------------------------------------
	// 1. Fetch from Ankhora Cloud
	// ----------------------------------------------------------
	cloudShares, err := uc.tc.ListLinkSharesWithMe(ctx, email)
	if err != nil {
		return nil, err
	}
	shares := cloudShares.Data
	utils.LogPretty("LinkShareUseCase - ListLinkSharesWithMe", cloudShares)
	// ----------------------------------------------------------
	// 2. Map to domain
	// ----------------------------------------------------------
	var domainShares []tracecore.WailsLinkShare
	for _, share := range shares {
		domainShares = append(domainShares, *share.ToWailsLinkShare())
	}

	return &domainShares, nil
}

func (uc *LinkShareUseCase) DeleteLinkShare(ctx context.Context, shareID string) (string, error) {
	return "", nil
}

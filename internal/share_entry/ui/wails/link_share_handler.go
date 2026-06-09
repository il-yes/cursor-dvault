package sahre_entry_ui_wails

import (
	"context"
	"vault-app/internal/logger/logger"
	share_entry_application_dto "vault-app/internal/share_entry/application"
	share_entry_use_cases "vault-app/internal/share_entry/application/use_cases"
	share_entry_domain "vault-app/internal/share_entry/domain"
	"vault-app/internal/tracecore"
	"vault-app/internal/utils"
	"errors"
)




type LinkShareHandler struct {	
	LinkShareUseCase share_entry_use_cases.LinkShareUseCase
	Logger *logger.Logger
}

func NewLinkShareHandler(
	linkShareUseCase share_entry_use_cases.LinkShareUseCase,
	logger *logger.Logger,
) LinkShareHandler {

	return LinkShareHandler{
		LinkShareUseCase: linkShareUseCase,
		Logger: logger,
	}
}

func (vh *LinkShareHandler) CreateLinkShare(ctx context.Context, email string, payload share_entry_application_dto.LinkShareCreateRequest) (*share_entry_domain.LinkShare, error) {
	if vh.LinkShareUseCase != (share_entry_use_cases.LinkShareUseCase{}) {
		vh.Logger.LogPretty("share_entry_handler - CreateLinkShare - linkShareUseCase is nil", nil)
		return nil, errors.New("link share use case is not initialized")
	}
	created, err := vh.LinkShareUseCase.CreateLinkShare(ctx, email, payload)
	if err != nil {
		vh.Logger.LogPretty("share_entry_handler - CreateLinkShare - linkShareUseCase.CreateLinkShare error: %v\n", err)
		return nil, err
	}
	utils.LogPretty("Handler response - created", created)

	return created, nil
}


func (vh *LinkShareHandler) ListLinkSharesByMe(ctx context.Context, email string) (*[]tracecore.WailsLinkShare, error) {
	if vh.LinkShareUseCase != (share_entry_use_cases.LinkShareUseCase{}) {
		vh.Logger.LogPretty("share_entry_handler - ListLinkSharesByMe - linkShareUseCase is nil", nil)
		return nil, errors.New("link share use case is not initialized")
	}
	result, err := vh.LinkShareUseCase.ListLinkSharesByMe(ctx, email)
	if err != nil {
		vh.Logger.LogPretty("share_entry_handler - ListLinkSharesByMe - linkShareUseCase.ListLinkSharesByMe error: %v\n", err)
		return nil, err
	}
	return result, nil
}


func (vh *LinkShareHandler) ListLinkSharesWithMe(ctx context.Context, email string) (*[]tracecore.WailsLinkShare, error) {
	if vh.LinkShareUseCase != (share_entry_use_cases.LinkShareUseCase{}) {
		vh.Logger.LogPretty("share_entry_handler - ListLinkSharesWithMe - linkShareUseCase is nil", nil)
		return nil, errors.New("link share use case is not initialized")
	}
	result, err := vh.LinkShareUseCase.ListLinkSharesWithMe(ctx, email)
	if err != nil {
		vh.Logger.LogPretty("share_entry_handler - ListLinkSharesWithMe - linkShareUseCase.ListLinkSharesWithMe error: %v\n", err)
		return nil, err
	}
	return result, nil
}
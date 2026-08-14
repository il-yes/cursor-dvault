package channel_ui

import (
	"context"
	"fmt"

	channel_application "vault-app/internal/channel/application"
	channel_usecase "vault-app/internal/channel/application/channel_lifecycle_usecases"
	channel_domain "vault-app/internal/channel/domain"
	tracecore_types "vault-app/internal/tracecore/types"
)

type ChannelHandler struct {
	createUseCase *channel_usecase.CreateChannelUsecase
	listUseCase   *channel_usecase.ListChannelUsecase
}

func NewChannelHandler(
	createUC *channel_usecase.CreateChannelUsecase,
	listUC *channel_usecase.ListChannelUsecase,
) *ChannelHandler {
	return &ChannelHandler{
		createUseCase: createUC,
		listUseCase:   listUC,
	}
}

func (h *ChannelHandler) CreateChannel(ctx context.Context, userID string, workspaceID string, title string, templateID string) (*tracecore_types.ChannelDTO, error) {
	if h.createUseCase == nil {
		return nil, fmt.Errorf("create channel use case is not initialized")
	}

	req := &channel_application.CreateChannelRequest{
		TemplateID:  templateID,
		Title:       title,
		WorkspaceID: workspaceID,
	}

	ch, err := h.createUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	return toTracecoreChannelDTO(ch), nil
}

func (h *ChannelHandler) ListChannels(ctx context.Context, userID string, workspaceID string) ([]tracecore_types.ChannelDTO, error) {
	if h.listUseCase == nil {
		return nil, fmt.Errorf("list channel use case is not initialized")
	}

	req := &channel_application.ListChannelsRequest{
		WorkspaceID: workspaceID,
	}

	channels, err := h.listUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	res := make([]tracecore_types.ChannelDTO, 0, len(channels))
	for _, ch := range channels {
		res = append(res, *toTracecoreChannelDTO(&ch))
	}

	return res, nil
}

func toTracecoreChannelDTO(ch *channel_domain.Channel) *tracecore_types.ChannelDTO {
	if ch == nil {
		return nil
	}
	return &tracecore_types.ChannelDTO{
		ID:          ch.ID,
		WorkspaceID: ch.WorkspaceID,
		Title:       ch.Title,
		TemplateID:  ch.TemplateID,
		Status:      string(ch.Status),
		CreatedAt:   ch.CreatedAt,
		UpdatedAt:   ch.UpdatedAt,
	}
}

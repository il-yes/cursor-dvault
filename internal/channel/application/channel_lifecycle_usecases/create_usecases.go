package channel_usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	channel_application "vault-app/internal/channel/application"
	channel_events "vault-app/internal/channel/application/events"
	channel_domain "vault-app/internal/channel/domain"
)

type CreateChannelUsecase struct {
	Repo      channel_domain.ChannelRepository
	DomainBus channel_events.ChannelEventBus
}

func NewCreateChannelUsecase(repo channel_domain.ChannelRepository, channelBus channel_events.ChannelEventBus) *CreateChannelUsecase {
	return &CreateChannelUsecase{
		Repo:      repo,
		DomainBus: channelBus,
	}
}

func (c *CreateChannelUsecase) Execute(ctx context.Context, req *channel_application.CreateChannelRequest) (*channel_domain.Channel, error) {
	if err := c.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := c.ValidateRequest(req); err != nil {
		return nil, err
	}

	channel := channel_domain.NewChannel(req.TemplateID, req.Title, req.WorkspaceID)

	created, err := c.Repo.CreateChannel(ctx, &channel_domain.CreateChannelRequest{
		Channel  : channel,
	})
	if err != nil {
		return nil, err
	}
	if created == nil {
		return nil, channel_domain.ErrRepositoryResponse
	}

	errEvent := c.DomainBus.PublishChannelCreated(
		ctx,
		channel_domain.ChannelCreated{
			EventID:        uuid.NewString(),
			EventTimestamp: time.Now(),
			ChannelID:    channel.ID,
		},
	)
	if errEvent != nil {
		return nil, errEvent
	}

	return &created.Data, nil
}

func (c *CreateChannelUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return channel_domain.ErrRepositoryNil
	}

	if c.DomainBus == nil {
		return channel_domain.ErrChannelBusRequired
	}

	return nil
}

func (c *CreateChannelUsecase) ValidateRequest(req *channel_application.CreateChannelRequest) error {
	if req == nil {
		return channel_domain.ErrRequestRequired
	}

	if req.TemplateID == "" {
		return channel_domain.ErrVaultIDRequired
	}

	if req.Title == "" {
		return channel_domain.ErrChannelOwnerRequired
	}

	if req.WorkspaceID == "" {
		return channel_domain.ErrChannelNameRequired
	}
	return nil
}

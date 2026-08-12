package channel_usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	channel_application "vault-app/internal/channel/application"
	channel_events "vault-app/internal/channel/application/events"
	channel_domain "vault-app/internal/channel/domain"
)

type ArchiveChannelUsecase struct {
	Repo      channel_domain.ChannelRepository
	DomainBus channel_events.ChannelEventBus
}

func NewArchiveChannelUsecase(repo channel_domain.ChannelRepository, channelBus channel_events.ChannelEventBus) *ArchiveChannelUsecase {
	return &ArchiveChannelUsecase{
		Repo:      repo,
		DomainBus: channelBus,
	}
}

func (c *ArchiveChannelUsecase) Execute(ctx context.Context, req *channel_application.ArchiveChannelRequest) error {
	if err := c.ValidateDependencies(); err != nil {
		return err
	}

	if err := c.ValidateRequest(req); err != nil {
		return err
	}

	// 1. Fetch current aggregate from repository
	getResp, err := c.Repo.GetChannel(ctx, &channel_domain.GetChannelRequest{
		ChannelID: req.ChannelID,
	})
	if err != nil {
		return err
	}
	if getResp == nil {
		return channel_domain.ErrRepositoryResponse
	}

	// 2. Apply domain behavior — aggregate enforces invariants
	channel := getResp.Data
	if err := channel.Archive(); err != nil {
		return err
	}

	// 3. Persist updated aggregate
	_, err = c.Repo.UpdateChannel(ctx, &channel_domain.UpdateChannelRequest{
		Channel: channel,
	})
	if err != nil {
		return err
	}

	// 4. Publish domain event
	errEvent := c.DomainBus.PublishChannelArchived(
		ctx,
		channel_domain.ChannelArchived{
			EventID:        uuid.NewString(),
			EventTimestamp: time.Now(),
			ChannelID:      req.ChannelID,
			WorkspaceID:    req.WorkspaceID,
		},
	)
	if errEvent != nil {
		return errEvent
	}

	return nil
}

func (c *ArchiveChannelUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return channel_domain.ErrRepositoryNil
	}

	if c.DomainBus == nil {
		return channel_domain.ErrChannelBusRequired
	}

	return nil
}

func (c *ArchiveChannelUsecase) ValidateRequest(req *channel_application.ArchiveChannelRequest) error {
	if req == nil {
		return channel_domain.ErrRequestRequired
	}

	if req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	if req.WorkspaceID == "" {
		return channel_domain.ErrWorkspaceIDRequired
	}

	return nil
}

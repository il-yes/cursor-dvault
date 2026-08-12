package channel_usecase

import (
	"context"

	channel_application "vault-app/internal/channel/application"
	channel_domain "vault-app/internal/channel/domain"
)

type DeleteChannelUsecase struct {
	Repo channel_domain.ChannelRepository
}

func NewDeleteChannelUsecase(repo channel_domain.ChannelRepository) *DeleteChannelUsecase {
	return &DeleteChannelUsecase{
		Repo: repo,
	}
}

func (c *DeleteChannelUsecase) Execute(ctx context.Context, req *channel_application.DeleteChannelRequest) error {
	if err := c.ValidateDependencies(); err != nil {
		return err
	}

	if err := c.ValidateRequest(req); err != nil {
		return err
	}

	err := c.Repo.DeleteChannel(ctx, &channel_domain.DeleteChannelRequest{
		ChannelID: req.ChannelID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *DeleteChannelUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return channel_domain.ErrRepositoryNil
	}

	return nil
}

func (c *DeleteChannelUsecase) ValidateRequest(req *channel_application.DeleteChannelRequest) error {
	if req == nil {
		return channel_domain.ErrRequestRequired
	}

	if req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	return nil
}

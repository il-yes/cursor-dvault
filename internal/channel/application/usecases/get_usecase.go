package channel_usecase

import (
	"context"
	"errors"

	channel_application "vault-app/internal/channel/application"
	channel_domain "vault-app/internal/channel/domain"
)

type GetChannelUsecase struct {
	Repo channel_domain.ChannelRepository
}

func NewGetChannelUsecase(repo channel_domain.ChannelRepository) *GetChannelUsecase {
	return &GetChannelUsecase{
		Repo: repo,
	}
}

func (c *GetChannelUsecase) Execute(ctx context.Context, req *channel_application.GetChannelRequest) (*channel_domain.Channel, error) {
	if err := c.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := c.ValidateRequest(req); err != nil {
		return nil, err
	}

	channel, err := c.Repo.GetChannel(ctx, &channel_domain.GetChannelRequest{
		ChannelID: req.ChannelID,
	})
	if err != nil {
		return nil, err
	}

	return &channel.Data, nil
}

func (c *GetChannelUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return errors.New(channel_domain.ErrRepositoryNil)
	}

	return nil
}

func (c *GetChannelUsecase) ValidateRequest(req *channel_application.GetChannelRequest) error {
	if req == nil {
		return errors.New(channel_domain.ErrRequestRequired)
	}

	if req.ChannelID == "" {
		return errors.New(channel_domain.ErrChannelIDRequired)
	}

	return nil
}

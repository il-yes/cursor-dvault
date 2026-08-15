package channel_usecase

import (
	"context"

	channel_application "vault-app/internal/channel/application"
	channel_domain "vault-app/internal/channel/domain"
)

type ActivateChannelUsecase struct {
	Repo channel_domain.ChannelRepository
}

func NewActivateChannelUsecase(repo channel_domain.ChannelRepository) *ActivateChannelUsecase {
	return &ActivateChannelUsecase{
		Repo: repo,
	}
}

// Execute activates an existing Channel through the authoritative Cloud
// backend. The Cloud aggregate enforces the "every gated slot must be
// fulfilled" invariant; the Desktop client does not duplicate it.
func (c *ActivateChannelUsecase) Execute(ctx context.Context, req *channel_application.ActivateChannelRequest) (*channel_domain.Channel, error) {
	if err := c.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := c.ValidateRequest(req); err != nil {
		return nil, err
	}

	resp, err := c.Repo.ActivateChannel(ctx, &channel_domain.ActivateChannelRequest{
		ChannelID: req.ChannelID,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, channel_domain.ErrRepositoryResponse
	}

	return &resp.Data, nil
}

func (c *ActivateChannelUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return channel_domain.ErrRepositoryNil
	}

	return nil
}

func (c *ActivateChannelUsecase) ValidateRequest(req *channel_application.ActivateChannelRequest) error {
	if req == nil {
		return channel_domain.ErrRequestRequired
	}

	if req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	return nil
}

package channel_usecase

import (
	"context"

	channel_application "vault-app/internal/channel/application"
	channel_domain "vault-app/internal/channel/domain"
)

type ListChannelUsecase struct {
	Repo channel_domain.ChannelRepository
}

func NewListChannelUsecase(repo channel_domain.ChannelRepository) *ListChannelUsecase {
	return &ListChannelUsecase{
		Repo: repo,
	}
}

func (c *ListChannelUsecase) Execute(ctx context.Context, req *channel_application.ListChannelsRequest) ([]channel_domain.Channel, error) {
	if err := c.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := c.ValidateRequest(req); err != nil {
		return nil, err
	}

	collection, err := c.Repo.ListChannels(ctx, &channel_domain.ListChannelsRequest{
		WorkspaceID: req.WorkspaceID,
	})
	if err != nil {
		return nil, err
	}

	return collection.Data, nil
}

func (c *ListChannelUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return channel_domain.ErrRepositoryNil
	}

	return nil
}

func (c *ListChannelUsecase) ValidateRequest(req *channel_application.ListChannelsRequest) error {
	if req == nil {
		return channel_domain.ErrRequestRequired
	}

	if req.WorkspaceID == "" {
		return channel_domain.ErrWorkspaceIDRequired
	}

	return nil
}

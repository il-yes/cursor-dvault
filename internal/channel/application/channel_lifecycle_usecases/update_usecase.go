package channel_usecase

import (
	"context"

	channel_application "vault-app/internal/channel/application"
	channel_domain "vault-app/internal/channel/domain"
)

type UpdateChannelUsecase struct {
	Repo channel_domain.ChannelRepository
}

func NewUpdateChannelUsecase(repo channel_domain.ChannelRepository) *UpdateChannelUsecase {
	return &UpdateChannelUsecase{
		Repo: repo,
	}
}

// Execute forwards the caller's update intent to the authoritative Cloud
// backend (PUT /channels/{id}). The usecase is thin: it validates its
// dependencies and request, then delegates. Cloud applies the fields it
// supports and remains authoritative for all domain validation; the returned
// Channel is the Cloud-persisted aggregate.
func (c *UpdateChannelUsecase) Execute(ctx context.Context, req *channel_application.UpdateChannelRequest) (*channel_domain.Channel, error) {
	if err := c.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := c.ValidateRequest(req); err != nil {
		return nil, err
	}

	resp, err := c.Repo.UpdateChannel(ctx, &channel_domain.UpdateChannelRequest{
		Channel: channel_domain.Channel{
			ID:          req.ChannelID,
			Title:       req.Title,
			Slots:       req.Slots,
			Assignments: req.Assignments,
			Properties:  req.Properties,
			Policy:      req.Policy,
		},
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, channel_domain.ErrRepositoryResponse
	}

	return &resp.Data, nil
}

func (c *UpdateChannelUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return channel_domain.ErrRepositoryNil
	}

	return nil
}

func (c *UpdateChannelUsecase) ValidateRequest(req *channel_application.UpdateChannelRequest) error {
	if req == nil {
		return channel_domain.ErrRequestRequired
	}

	if req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	return nil
}
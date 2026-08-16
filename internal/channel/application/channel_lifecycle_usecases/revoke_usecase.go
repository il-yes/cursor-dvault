package channel_usecase

import (
	"context"

	channel_application "vault-app/internal/channel/application"
	channel_domain "vault-app/internal/channel/domain"
)

type RevokeChannelUsecase struct {
	Repo channel_domain.ChannelRepository
}

func NewRevokeChannelUsecase(repo channel_domain.ChannelRepository) *RevokeChannelUsecase {
	return &RevokeChannelUsecase{
		Repo: repo,
	}
}

// Execute revokes an existing Channel through the authoritative Cloud backend.
// The Desktop client only validates the request and forwards it; the Cloud
// aggregate owns the revocation invariants, the RevokedAt timestamp, and the
// channel.revoked lifecycle event. The Cloud revoke response carries no
// Channel data, so this usecase returns only an error; callers refresh the
// workspace channel list through ListChannels to observe the new status.
func (c *RevokeChannelUsecase) Execute(ctx context.Context, req *channel_application.RevokeChannelRequest) error {
	if err := c.ValidateDependencies(); err != nil {
		return err
	}

	if err := c.ValidateRequest(req); err != nil {
		return err
	}

	return c.Repo.RevokeChannel(ctx, &channel_domain.RevokeChannelRequest{
		ChannelID: req.ChannelID,
	})
}

func (c *RevokeChannelUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return channel_domain.ErrRepositoryNil
	}

	return nil
}

func (c *RevokeChannelUsecase) ValidateRequest(req *channel_application.RevokeChannelRequest) error {
	if req == nil {
		return channel_domain.ErrRequestRequired
	}

	if req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	return nil
}

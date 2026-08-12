package channel_assignment_usecases

import (
	"context"

	channel_application "vault-app/internal/channel/application"
	channel_domain    "vault-app/internal/channel/domain"
)

func (u *ChannelAssignmentUsecase) UpdateAssignment(ctx context.Context, req channel_application.UpdateAssignmentRequest) (*channel_domain.Channel, error) {
	if err := u.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := u.validateUpdateAssignmentRequest(req); err != nil {
		return nil, err
	}

	channel, err := u.repo.GetChannel(ctx, &channel_domain.GetChannelRequest{ChannelID: req.ChannelID})
	if err != nil {
		return nil, err
	}

	if channel == nil {
		return nil, channel_domain.ErrChannelNotFound
	}

	if ok := channel.Data.UpdateAssignment(req.Assignment); !ok {
		return nil, channel_domain.ErrChannelNotModifiable
	}

	updatedChannel, err := u.repo.UpdateChannel(ctx, &channel_domain.UpdateChannelRequest{Channel: channel.Data})
	if err != nil {
		return nil, err
	}

	return &updatedChannel.Data, nil
}

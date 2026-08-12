package channel_assignment_usecases

import (
	"context"

	channel_application "vault-app/internal/channel/application"
	channel_domain    "vault-app/internal/channel/domain"
)

func (u *ChannelAssignmentUsecase) RemoveAssignment(ctx context.Context, req channel_application.RemoveAssignmentRequest,
) (*channel_domain.Channel, error) {
	if err := u.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := u.validateRemoveAssignmentRequest(req); err != nil {
		return nil, err
	}

	channel, err := u.repo.GetChannel(ctx, &channel_domain.GetChannelRequest{ChannelID: req.ChannelID})
	if err != nil {
		return nil, err
	}

	if channel == nil {
		return nil, channel_domain.ErrChannelNotFound
	}

	// if err := channel.Data.RemoveAssignmentBySlotID(req.AssignmentID); err != nil {
	// 	return nil, err
	// }

	updatedChannel, err := u.repo.UpdateChannel(ctx, &channel_domain.UpdateChannelRequest{Channel: channel.Data})
	if err != nil {
		return nil, err
	}

	return &updatedChannel.Data, nil
}

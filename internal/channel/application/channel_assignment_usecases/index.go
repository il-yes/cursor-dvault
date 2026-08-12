package channel_assignment_usecases

import (
	channel_application "vault-app/internal/channel/application"
	channel_domain "vault-app/internal/channel/domain"
)

type ChannelAssignmentUsecase struct {
	repo channel_domain.ChannelRepository
}

func NewChannelAssignmentUsecase(repo channel_domain.ChannelRepository) *ChannelAssignmentUsecase {
	return &ChannelAssignmentUsecase{repo: repo}
}

func (u *ChannelAssignmentUsecase) ValidateDependencies() error {
	if u.repo == nil {
		return channel_domain.ErrChannelRepositoryRequired
	}
	return nil
}

func (u *ChannelAssignmentUsecase) validateAddAssignmentRequest(req channel_application.AddAssignmentRequest) error {
	if req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	if req.Assignment.SlotID == "" {
		return channel_domain.ErrAssignmentIDRequired
	}

	return nil
}

func (u *ChannelAssignmentUsecase) validateUpdateAssignmentRequest(req channel_application.UpdateAssignmentRequest) error {
	if req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	if req.Assignment.SlotID == "" {
		return channel_domain.ErrAssignmentIDRequired
	}

	return nil
}

func (u *ChannelAssignmentUsecase) validateRemoveAssignmentRequest(req channel_application.RemoveAssignmentRequest) error {
	if req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	if req.AssignmentID == "" {
		return channel_domain.ErrAssignmentIDRequired
	}

	return nil
}

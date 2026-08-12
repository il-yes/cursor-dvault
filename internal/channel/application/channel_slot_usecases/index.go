package channel_slot_usecases

import (
	channel_application "vault-app/internal/channel/application"
	channel_domain "vault-app/internal/channel/domain"
)

type ChannelSlotUsecase struct {
	repo channel_domain.ChannelRepository
}

func NewChannelSlotUsecase(repo channel_domain.ChannelRepository) *ChannelSlotUsecase {
	return &ChannelSlotUsecase{repo: repo}
}

func (u *ChannelSlotUsecase) ValidateDependencies() error {
	if u.repo == nil {
		return channel_domain.ErrChannelRepositoryRequired
	}
	return nil
}

func (u *ChannelSlotUsecase) validateAddSlotRequest(req channel_application.AddSlotRequest) error {
	if req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	if req.Slot.ID == "" {
		return channel_domain.ErrSlotIDRequired
	}

	return nil
}

func (u *ChannelSlotUsecase) validateUpdateSlotRequest(req channel_application.UpdateSlotRequest) error {
	if req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	if req.Slot.ID == "" {
		return channel_domain.ErrSlotIDRequired
	}

	return nil
}

func (u *ChannelSlotUsecase) validateRemoveSlotRequest(req channel_application.RemoveSlotRequest) error {
	if req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	if req.SlotID == "" {
		return channel_domain.ErrSlotIDRequired
	}

	return nil
}

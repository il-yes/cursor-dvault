package channel_usecase

import (
	"context"

	channel_application "vault-app/internal/channel/application"
	channel_domain "vault-app/internal/channel/domain"
)

// AddParticipantUsecase joins a vault to a channel through the authoritative
// Cloud backend. Cloud validates the join (revoked checks, slot existence,
// role derivation) and persists the participant; the Desktop never decides
// locally whether a vault may join.
type AddParticipantUsecase struct {
	Repo channel_domain.ChannelRepository
}

func NewAddParticipantUsecase(repo channel_domain.ChannelRepository) *AddParticipantUsecase {
	return &AddParticipantUsecase{
		Repo: repo,
	}
}

func (c *AddParticipantUsecase) Execute(ctx context.Context, req *channel_application.AddParticipantRequest) (*channel_domain.Participant, error) {
	if err := c.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := c.ValidateRequest(req); err != nil {
		return nil, err
	}

	resp, err := c.Repo.AddParticipant(ctx, &channel_domain.JoinChannelRequest{
		ChannelID: req.ChannelID,
		VaultID:   req.VaultID,
		PublicKey: req.PublicKey,
		Direction: req.Direction,
		SlotID:    req.SlotID,
		Role:      req.Role,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, channel_domain.ErrRepositoryResponse
	}

	return &resp.Data, nil
}

func (c *AddParticipantUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return channel_domain.ErrRepositoryNil
	}

	return nil
}

func (c *AddParticipantUsecase) ValidateRequest(req *channel_application.AddParticipantRequest) error {
	if req == nil {
		return channel_domain.ErrRequestRequired
	}

	if req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	if req.VaultID == "" {
		return channel_domain.ErrVaultIDRequired
	}

	return nil
}

// ListParticipantsUsecase returns the vaults Cloud has persisted as channel
// participants. The Desktop never fabricates participants; only what Cloud
// returns is surfaced to the caller.
type ListParticipantsUsecase struct {
	Repo channel_domain.ChannelRepository
}

func NewListParticipantsUsecase(repo channel_domain.ChannelRepository) *ListParticipantsUsecase {
	return &ListParticipantsUsecase{
		Repo: repo,
	}
}

func (c *ListParticipantsUsecase) Execute(ctx context.Context, req *channel_application.ListParticipantsRequest) ([]channel_domain.Participant, error) {
	if err := c.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := c.ValidateRequest(req); err != nil {
		return nil, err
	}

	collection, err := c.Repo.ListParticipants(ctx, &channel_domain.ListParticipantsRequest{
		ChannelID: req.ChannelID,
	})
	if err != nil {
		return nil, err
	}

	return collection.Data, nil
}

func (c *ListParticipantsUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return channel_domain.ErrRepositoryNil
	}

	return nil
}

func (c *ListParticipantsUsecase) ValidateRequest(req *channel_application.ListParticipantsRequest) error {
	if req == nil {
		return channel_domain.ErrRequestRequired
	}

	if req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	return nil
}

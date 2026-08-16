package channel_usecase

import (
	"context"

	channel_application "vault-app/internal/channel/application"
	channel_domain "vault-app/internal/channel/domain"
)

// InviteToChannelUsecase creates a channel invitation through the authoritative
// Cloud backend. Cloud persists the pending invitation and dedupes pending
// invitations for the same channel + invitee; the Desktop never fabricates an
// invitation locally.
type InviteToChannelUsecase struct {
	Repo channel_domain.ChannelRepository
}

func NewInviteToChannelUsecase(repo channel_domain.ChannelRepository) *InviteToChannelUsecase {
	return &InviteToChannelUsecase{
		Repo: repo,
	}
}

func (c *InviteToChannelUsecase) Execute(ctx context.Context, req *channel_application.InviteToChannelRequest) (*channel_domain.Invitation, error) {
	if err := c.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := c.ValidateRequest(req); err != nil {
		return nil, err
	}

	resp, err := c.Repo.InviteToChannel(ctx, &channel_domain.InviteToChannelRequest{
		ChannelID:      req.ChannelID,
		InviterVaultID: req.InviterVaultID,
		InviteeVaultID: req.InviteeVaultID,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, channel_domain.ErrRepositoryResponse
	}

	return &resp.Data, nil
}

func (c *InviteToChannelUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return channel_domain.ErrRepositoryNil
	}

	return nil
}

func (c *InviteToChannelUsecase) ValidateRequest(req *channel_application.InviteToChannelRequest) error {
	if req == nil {
		return channel_domain.ErrRequestRequired
	}

	if req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	if req.InviterVaultID == "" {
		return channel_domain.ErrVaultIDRequired
	}

	if req.InviteeVaultID == "" {
		return channel_domain.ErrVaultIDRequired
	}

	return nil
}

// AcceptChannelInvitationUsecase accepts a pending channel invitation through
// the authoritative Cloud backend. Cloud validates the acceptance (the
// accepting vault must be the invitation's invitee) and persists the resulting
// participant; the accept response carries the accepted Invitation, not the
// participant. Cloud is idempotent: accepting an already-accepted invitation
// returns the accepted invitation without a duplicate participant.
type AcceptChannelInvitationUsecase struct {
	Repo channel_domain.ChannelRepository
}

func NewAcceptChannelInvitationUsecase(repo channel_domain.ChannelRepository) *AcceptChannelInvitationUsecase {
	return &AcceptChannelInvitationUsecase{
		Repo: repo,
	}
}

func (c *AcceptChannelInvitationUsecase) Execute(ctx context.Context, req *channel_application.AcceptChannelInvitationRequest) (*channel_domain.Invitation, error) {
	if err := c.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := c.ValidateRequest(req); err != nil {
		return nil, err
	}

	resp, err := c.Repo.AcceptChannelInvitation(ctx, &channel_domain.AcceptInvitationRequest{
		InvitationID:     req.InvitationID,
		InviteeVaultID:   req.InviteeVaultID,
		InviteePublicKey: req.InviteePublicKey,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, channel_domain.ErrRepositoryResponse
	}

	return &resp.Data, nil
}

func (c *AcceptChannelInvitationUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return channel_domain.ErrRepositoryNil
	}

	return nil
}

func (c *AcceptChannelInvitationUsecase) ValidateRequest(req *channel_application.AcceptChannelInvitationRequest) error {
	if req == nil {
		return channel_domain.ErrRequestRequired
	}

	if req.InvitationID == "" {
		return channel_domain.ErrInvitationIDRequired
	}

	if req.InviteeVaultID == "" {
		return channel_domain.ErrVaultIDRequired
	}

	if req.InviteePublicKey == "" {
		return channel_domain.ErrInviteePublicKeyRequired
	}

	return nil
}

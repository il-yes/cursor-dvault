package channel_domain

import (
	"context"

	tracecore_types "vault-app/internal/tracecore/types"
)

type CreateChannelRequest struct {
	Channel Channel
}

type SlotRequest struct {
	ID      string
	Name    string
	Role    string
	VaultID string
	Gated   bool
	Order   int
}

type PropertyRequest struct {
	Key   string
	Value string
}

type AssignmentRequest struct {
	SlotID       string
	OwnerID      string
	PublicKey    string
	VaultAddress string
}

type FederationRequest struct {
	VaultAID          string
	VaultBID          string
	AllowedEventTypes []string
	AllowedPaths      []string
	AllowedDirections string
}

type GetChannelRequest struct {
	ChannelID string
}
type UpdateChannelRequest struct {
	Channel Channel
}

type DeleteChannelRequest struct {
	ChannelID string
}
type ListChannelsRequest struct {
	WorkspaceID string
}

type JoinChannelRequest struct {
	ChannelID string

	VaultID   string
	PublicKey string

	Direction string
}

type LeaveChannelRequest struct {
	ChannelID string
	VaultID   string
}

type ListParticipantsRequest struct {
	ChannelID string
}

type GetParticipantRequest struct {
	ChannelID string
	VaultID   string
}

type InviteToChannelRequest struct {
	ChannelID      string
	InviterVaultID string
	InviteeVaultID string
}

type AcceptInvitationRequest struct {
	InvitationID     string
	InviteeVaultID   string
	InviteePublicKey string
}

type RejectInvitationRequest struct {
	InvitationID string
}

type RevokeInvitationRequest struct {
	InvitationID   string
	InviterVaultID string
}
type ChannelRepository interface {
	CreateChannel(ctx context.Context, req *CreateChannelRequest) (*tracecore_types.CloudResponse[Channel], error)
	ListChannels(ctx context.Context, req *ListChannelsRequest) (*tracecore_types.CloudResponse[[]Channel], error)
	GetChannel(ctx context.Context, req *GetChannelRequest) (*tracecore_types.CloudResponse[Channel], error)
	DeleteChannel(ctx context.Context, req *DeleteChannelRequest) error
	UpdateChannel(ctx context.Context, req *UpdateChannelRequest) (*tracecore_types.CloudResponse[Channel], error)
	ActivateChannel(ctx context.Context, req *AcceptInvitationRequest) error
	RevokeChannel(ctx context.Context, req *RevokeInvitationRequest)  error
}

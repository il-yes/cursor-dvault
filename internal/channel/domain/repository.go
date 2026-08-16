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

type ActivateChannelRequest struct {
	ChannelID string
}

type JoinChannelRequest struct {
	ChannelID string

	VaultID   string
	PublicKey string

	Direction string

	// SlotID and Role are optional. When SlotID is provided the Cloud backend
	// verifies the slot exists and derives the participant Role from it when
	// Role is empty. Both are forwarded verbatim to the Cloud contract.
	SlotID string
	Role   string
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

// RevokeChannelRequest carries only the channel id. The Cloud backend is
// authoritative for revocation invariants; it returns no Channel data, so the
// client refreshes through ListChannels after a successful revoke.
type RevokeChannelRequest struct {
	ChannelID string
}

type ChannelRepository interface {
	CreateChannel(ctx context.Context, req *CreateChannelRequest) (*tracecore_types.CloudResponse[Channel], error)
	ListChannels(ctx context.Context, req *ListChannelsRequest) (*tracecore_types.CloudResponse[[]Channel], error)
	GetChannel(ctx context.Context, req *GetChannelRequest) (*tracecore_types.CloudResponse[Channel], error)
	DeleteChannel(ctx context.Context, req *DeleteChannelRequest) error
	UpdateChannel(ctx context.Context, req *UpdateChannelRequest) (*tracecore_types.CloudResponse[Channel], error)
	ActivateChannel(ctx context.Context, req *ActivateChannelRequest) (*tracecore_types.CloudResponse[Channel], error)
	RevokeChannel(ctx context.Context, req *RevokeChannelRequest) error
	// AddParticipant joins a vault to a channel through the authoritative Cloud
	// backend (POST /channels/{id}/participants). Cloud is responsible for
	// validating the join; the Desktop never decides locally whether a vault may
	// join. Cloud is idempotent: joining a vault that is already a participant
	// returns the persisted participant.
	AddParticipant(ctx context.Context, req *JoinChannelRequest) (*tracecore_types.CloudResponse[Participant], error)
	// ListParticipants returns the vaults that Cloud has persisted as channel
	// participants (GET /channels/{id}/participants). An empty result is valid.
	ListParticipants(ctx context.Context, req *ListParticipantsRequest) (*tracecore_types.CloudResponse[[]Participant], error)
	// InviteToChannel creates a channel invitation through the authoritative
	// Cloud backend (POST /channels/{id}/invitations). Cloud persists the
	// pending invitation and is authoritative for its lifecycle; the Desktop
	// never fabricates an invitation locally. Cloud dedupes pending invitations
	// for the same channel + invitee, returning the existing invitation.
	InviteToChannel(ctx context.Context, req *InviteToChannelRequest) (*tracecore_types.CloudResponse[Invitation], error)
	// AcceptChannelInvitation accepts a pending invitation through the
	// authoritative Cloud backend (POST /channels/invitations/{id}/accept).
	// Cloud validates the acceptance (the inviting vault may only accept an
	// invitation addressed to it) and persists the resulting participant. The
	// accept response carries the accepted Invitation, not the participant;
	// participants are observed through ListParticipants. Cloud is idempotent:
	// accepting an already-accepted invitation returns the accepted invitation
	// without creating a duplicate participant.
	AcceptChannelInvitation(ctx context.Context, req *AcceptInvitationRequest) (*tracecore_types.CloudResponse[Invitation], error)
}

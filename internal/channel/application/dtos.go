package channel_application

import channel_domain "vault-app/internal/channel/domain"

type CreateChannelRequest struct {
	TemplateID  string
	Title       string
	WorkspaceID string
	Slots       []channel_domain.Slot
	Assignments []channel_domain.Assignment
	Properties  []channel_domain.ChannelProperty
	Policy      channel_domain.Policy
	Federation  string
}

type UpdateChannelRequest struct {
	ChannelID   string
	Title       string
	Slots       []channel_domain.Slot
	Assignments []channel_domain.Assignment
	Properties  []channel_domain.ChannelProperty
	Policy      channel_domain.Policy
}

type ListChannelsRequest struct {
	WorkspaceID string
}

type ArchiveChannelRequest struct {
	ChannelID   string
	WorkspaceID string
}

type GetChannelRequest struct {
	ChannelID string
}

type DeleteChannelRequest struct {
	ChannelID string
}

type ActivateChannelRequest struct {
	ChannelID string
}

type RevokeChannelRequest struct {
	ChannelID string
}

type AddParticipantRequest struct {
	ChannelID string
	VaultID   string
	PublicKey string
	Direction string
	SlotID    string
	Role      string
}

type ListParticipantsRequest struct {
	ChannelID string
}

// InviteToChannelRequest mirrors the Cloud invitation contract
// (POST /channels/{id}/invitations). The invitation carries no slot or role
// information; role semantics are a channel participant concern.
type InviteToChannelRequest struct {
	ChannelID      string
	InviterVaultID string
	InviteeVaultID string
}

// AcceptChannelInvitationRequest mirrors the Cloud invitation-accept contract
// (POST /channels/invitations/{id}/accept). The accepting vault must be the
// invitation's invitee; the invitee public key is stored on the participant the
// Cloud creates on acceptance.
type AcceptChannelInvitationRequest struct {
	InvitationID     string
	InviteeVaultID   string
	InviteePublicKey string
}

type AddSlotRequest struct {
	ChannelID string
	Slot      channel_domain.Slot
}

type UpdateSlotRequest struct {
	ChannelID string
	Slot      channel_domain.Slot
}

type RemoveSlotRequest struct {
	ChannelID string
	SlotID    string
}

type AddAssignmentRequest struct {
	ChannelID  string
	Assignment channel_domain.Assignment
}

type UpdateAssignmentRequest struct {
	ChannelID  string
	Assignment channel_domain.Assignment
}

type RemoveAssignmentRequest struct {
	ChannelID    string
	AssignmentID string
}

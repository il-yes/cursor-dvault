package channel_application

import channel_domain "vault-app/internal/channel/domain"


type CreateChannelRequest struct {
	TemplateID  string
	Title       string
	WorkspaceID string
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


type InviteRequest struct {
	ChannelId string
	InviterVaultId string
	InviteeVaultId string
}


type AcceptChannelInviteRequest struct {
	InvitationID     string
	Direction        string
	TrustName        string
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
	ChannelID string
	Assignment      channel_domain.Assignment
}

type UpdateAssignmentRequest struct {
	ChannelID string
	Assignment      channel_domain.Assignment
}

type RemoveAssignmentRequest struct {
	ChannelID string
	AssignmentID    string
}
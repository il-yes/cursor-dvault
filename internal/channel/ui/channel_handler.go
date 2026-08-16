package channel_ui

import (
	"context"
	"fmt"

	channel_application "vault-app/internal/channel/application"
	channel_usecase "vault-app/internal/channel/application/channel_lifecycle_usecases"
	channel_domain "vault-app/internal/channel/domain"
	tracecore_types "vault-app/internal/tracecore/types"
)

type ChannelHandler struct {
	createUseCase           *channel_usecase.CreateChannelUsecase
	listUseCase             *channel_usecase.ListChannelUsecase
	getUseCase              *channel_usecase.GetChannelUsecase
	updateUseCase           *channel_usecase.UpdateChannelUsecase
	deleteUseCase           *channel_usecase.DeleteChannelUsecase
	activateUseCase         *channel_usecase.ActivateChannelUsecase
	revokeUseCase           *channel_usecase.RevokeChannelUsecase
	addParticipantUseCase   *channel_usecase.AddParticipantUsecase
	listParticipantsUseCase *channel_usecase.ListParticipantsUsecase
	inviteToChannelUseCase  *channel_usecase.InviteToChannelUsecase
	acceptInvitationUseCase *channel_usecase.AcceptChannelInvitationUsecase
}

func NewChannelHandler(
	createUC *channel_usecase.CreateChannelUsecase,
	listUC *channel_usecase.ListChannelUsecase,
	getUC *channel_usecase.GetChannelUsecase,
	updateUC *channel_usecase.UpdateChannelUsecase,
	deleteUC *channel_usecase.DeleteChannelUsecase,
	activateUC *channel_usecase.ActivateChannelUsecase,
	revokeUC *channel_usecase.RevokeChannelUsecase,
	addParticipantUC *channel_usecase.AddParticipantUsecase,
	listParticipantsUC *channel_usecase.ListParticipantsUsecase,
	inviteToChannelUC *channel_usecase.InviteToChannelUsecase,
	acceptInvitationUC *channel_usecase.AcceptChannelInvitationUsecase,
) *ChannelHandler {
	return &ChannelHandler{
		createUseCase:           createUC,
		listUseCase:             listUC,
		getUseCase:              getUC,
		updateUseCase:           updateUC,
		deleteUseCase:           deleteUC,
		activateUseCase:         activateUC,
		revokeUseCase:           revokeUC,
		addParticipantUseCase:   addParticipantUC,
		listParticipantsUseCase: listParticipantsUC,
		inviteToChannelUseCase:  inviteToChannelUC,
		acceptInvitationUseCase: acceptInvitationUC,
	}
}

func (h *ChannelHandler) CreateChannel(ctx context.Context, userID string, workspaceID string, title string, templateID string, slots []channel_domain.Slot, assignments []channel_domain.Assignment, properties []channel_domain.ChannelProperty, policy channel_domain.Policy, federation string) (*tracecore_types.ChannelDTO, error) {
	if h.createUseCase == nil {
		return nil, fmt.Errorf("create channel use case is not initialized")
	}

	req := &channel_application.CreateChannelRequest{
		TemplateID:  templateID,
		Title:       title,
		WorkspaceID: workspaceID,
		Slots:       slots,
		Assignments: assignments,
		Properties:  properties,
		Policy:      policy,
		Federation:  federation,
	}

	ch, err := h.createUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	return toTracecoreChannelDTO(ch), nil
}

// GetChannel fetches a single Channel from the authoritative Cloud backend
// (GET /channels/{id}). The returned Channel is the Cloud-persisted aggregate;
// no local channel is ever fabricated.
func (h *ChannelHandler) GetChannel(ctx context.Context, userID string, channelID string) (*tracecore_types.ChannelDTO, error) {
	if h.getUseCase == nil {
		return nil, fmt.Errorf("get channel use case is not initialized")
	}

	req := &channel_application.GetChannelRequest{
		ChannelID: channelID,
	}

	ch, err := h.getUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	return toTracecoreChannelDTO(ch), nil
}

// UpdateChannel updates an existing Channel through the authoritative Cloud
// backend (PUT /channels/{id}). The Cloud-persisted aggregate is returned; no
// local mutation is performed.
func (h *ChannelHandler) UpdateChannel(ctx context.Context, userID string, channelID string, title string, slots []channel_domain.Slot, assignments []channel_domain.Assignment, properties []channel_domain.ChannelProperty, policy channel_domain.Policy) (*tracecore_types.ChannelDTO, error) {
	if h.updateUseCase == nil {
		return nil, fmt.Errorf("update channel use case is not initialized")
	}

	req := &channel_application.UpdateChannelRequest{
		ChannelID:   channelID,
		Title:       title,
		Slots:       slots,
		Assignments: assignments,
		Properties:  properties,
		Policy:      policy,
	}

	ch, err := h.updateUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	return toTracecoreChannelDTO(ch), nil
}

// DeleteChannel deletes a Channel through the authoritative Cloud backend
// (DELETE /channels/{id}). Cloud is the single source of truth for channel
// existence; a 2xx response is success and HTTP >=400 is surfaced verbatim.
func (h *ChannelHandler) DeleteChannel(ctx context.Context, userID string, channelID string) error {
	if h.deleteUseCase == nil {
		return fmt.Errorf("delete channel use case is not initialized")
	}

	req := &channel_application.DeleteChannelRequest{
		ChannelID: channelID,
	}

	return h.deleteUseCase.Execute(ctx, req)
}

func (h *ChannelHandler) ActivateChannel(ctx context.Context, userID string, channelID string) (*tracecore_types.ChannelDTO, error) {
	if h.activateUseCase == nil {
		return nil, fmt.Errorf("activate channel use case is not initialized")
	}

	req := &channel_application.ActivateChannelRequest{
		ChannelID: channelID,
	}

	ch, err := h.activateUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	return toTracecoreChannelDTO(ch), nil
}

func (h *ChannelHandler) RevokeChannel(ctx context.Context, userID string, channelID string) error {
	if h.revokeUseCase == nil {
		return fmt.Errorf("revoke channel use case is not initialized")
	}

	req := &channel_application.RevokeChannelRequest{
		ChannelID: channelID,
	}

	return h.revokeUseCase.Execute(ctx, req)
}

// AddParticipant joins an external vault to a channel through the authoritative
// Cloud backend. slotID and role are optional: when provided they are forwarded
// verbatim to the Cloud JoinChannelRequest contract, which derives/validates the
// participant role. The persisted Cloud participant is returned.
func (h *ChannelHandler) AddParticipant(
	ctx context.Context,
	userID string,
	channelID string,
	vaultID string,
	publicKey string,
	direction string,
	slotID string,
	role string,
) (*tracecore_types.ChannelParticipantDTO, error) {
	if h.addParticipantUseCase == nil {
		return nil, fmt.Errorf("add participant use case is not initialized")
	}

	req := &channel_application.AddParticipantRequest{
		ChannelID: channelID,
		VaultID:   vaultID,
		PublicKey: publicKey,
		Direction: direction,
		SlotID:    slotID,
		Role:      role,
	}

	p, err := h.addParticipantUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	return toChannelParticipantDTO(p), nil
}

// ListParticipants returns the vaults Cloud has persisted as participants for
// the channel. An empty result is valid.
func (h *ChannelHandler) ListParticipants(ctx context.Context, userID string, channelID string) ([]tracecore_types.ChannelParticipantDTO, error) {
	if h.listParticipantsUseCase == nil {
		return nil, fmt.Errorf("list participants use case is not initialized")
	}

	req := &channel_application.ListParticipantsRequest{
		ChannelID: channelID,
	}

	participants, err := h.listParticipantsUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	res := make([]tracecore_types.ChannelParticipantDTO, 0, len(participants))
	for _, p := range participants {
		res = append(res, *toChannelParticipantDTO(&p))
	}

	return res, nil
}

func toChannelParticipantDTO(p *channel_domain.Participant) *tracecore_types.ChannelParticipantDTO {
	if p == nil {
		return nil
	}

	permissions := p.Permissions
	if permissions == nil {
		permissions = []string{}
	}

	return &tracecore_types.ChannelParticipantDTO{
		ChannelID:   p.ChannelID,
		VaultID:     p.VaultID,
		PublicKey:   p.PublicKey,
		Direction:   p.Direction,
		JoinedAt:    p.JoinedAt,
		Role:        p.Role,
		Permissions: permissions,
	}
}

// InviteToChannel creates a channel invitation through the authoritative Cloud
// backend (POST /channels/{id}/invitations). Cloud persists the pending
// invitation and dedupes pending invitations for the same channel + invitee.
// The invitation carries no slot or role information.
func (h *ChannelHandler) InviteToChannel(
	ctx context.Context,
	userID string,
	channelID string,
	inviterVaultID string,
	inviteeVaultID string,
) (*tracecore_types.ChannelInvitationDTO, error) {
	if h.inviteToChannelUseCase == nil {
		return nil, fmt.Errorf("invite to channel use case is not initialized")
	}

	req := &channel_application.InviteToChannelRequest{
		ChannelID:      channelID,
		InviterVaultID: inviterVaultID,
		InviteeVaultID: inviteeVaultID,
	}

	inv, err := h.inviteToChannelUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	return toChannelInvitationDTO(inv), nil
}

// AcceptChannelInvitation accepts a pending channel invitation through the
// authoritative Cloud backend (POST /channels/invitations/{id}/accept). Cloud
// validates the acceptance and persists the resulting participant; the accept
// response carries the accepted Invitation, not the participant. Cloud is
// idempotent: accepting an already-accepted invitation returns the accepted
// invitation without a duplicate participant.
func (h *ChannelHandler) AcceptChannelInvitation(
	ctx context.Context,
	userID string,
	invitationID string,
	inviteeVaultID string,
	inviteePublicKey string,
) (*tracecore_types.ChannelInvitationDTO, error) {
	if h.acceptInvitationUseCase == nil {
		return nil, fmt.Errorf("accept channel invitation use case is not initialized")
	}

	req := &channel_application.AcceptChannelInvitationRequest{
		InvitationID:     invitationID,
		InviteeVaultID:   inviteeVaultID,
		InviteePublicKey: inviteePublicKey,
	}

	inv, err := h.acceptInvitationUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	return toChannelInvitationDTO(inv), nil
}

func toChannelInvitationDTO(inv *channel_domain.Invitation) *tracecore_types.ChannelInvitationDTO {
	if inv == nil {
		return nil
	}

	return &tracecore_types.ChannelInvitationDTO{
		ID:             inv.ID,
		ChannelID:      inv.ChannelID,
		InviterVaultID: inv.InviterVaultID,
		InviteeVaultID: inv.InviteeVaultID,
		Status:         string(inv.Status),
		CreatedAt:      inv.CreatedAt,
		AcceptedAt:     inv.AcceptedAt,
	}
}

func (h *ChannelHandler) ListChannels(ctx context.Context, userID string, workspaceID string) ([]tracecore_types.ChannelDTO, error) {
	if h.listUseCase == nil {
		return nil, fmt.Errorf("list channel use case is not initialized")
	}

	req := &channel_application.ListChannelsRequest{
		WorkspaceID: workspaceID,
	}

	channels, err := h.listUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	res := make([]tracecore_types.ChannelDTO, 0, len(channels))
	for _, ch := range channels {
		res = append(res, *toTracecoreChannelDTO(&ch))
	}

	return res, nil
}

func toTracecoreChannelDTO(ch *channel_domain.Channel) *tracecore_types.ChannelDTO {
	if ch == nil {
		return nil
	}

	slots := make([]tracecore_types.ChannelSlotDTO, 0, len(ch.Slots))
	for _, s := range ch.Slots {
		slots = append(slots, tracecore_types.ChannelSlotDTO{
			ID:      s.ID,
			Name:    s.Name,
			Role:    s.Role,
			VaultID: s.VaultID,
			Gated:   s.Gated,
			Order:   s.Order,
		})
	}

	assignments := make([]tracecore_types.ChannelAssignmentDTO, 0, len(ch.Assignments))
	for _, a := range ch.Assignments {
		assignments = append(assignments, tracecore_types.ChannelAssignmentDTO{
			SlotID:       a.SlotID,
			OwnerID:      a.OwnerID,
			PublicKey:    a.PublicKey,
			VaultAddress: a.VaultAddress,
		})
	}

	properties := make([]tracecore_types.ChannelPropertyDTO, 0, len(ch.Properties))
	for _, p := range ch.Properties {
		properties = append(properties, tracecore_types.ChannelPropertyDTO{
			Key:   p.Key,
			Value: p.Value,
		})
	}

	return &tracecore_types.ChannelDTO{
		ID:          ch.ID,
		WorkspaceID: ch.WorkspaceID,
		Title:       ch.Title,
		TemplateID:  ch.TemplateID,
		Status:      string(ch.Status),
		Slots:       slots,
		Assignments: assignments,
		Properties:  properties,
		Policy:      ch.Policy,
		Federation:  ch.Federation,
		CreatedAt:   ch.CreatedAt,
		UpdatedAt:   ch.UpdatedAt,
		RevokedAt:   ch.RevokedAt,
	}
}

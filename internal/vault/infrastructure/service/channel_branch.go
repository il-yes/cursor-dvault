package vaults_service

import (
	"context"
	"encoding/json"
	"time"

	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
)

type ChannelNode struct {
	Type    string `json:"type"`
	Version string `json:"version"`

	ID           string `json:"id"`
	WorkspaceID  string `json:"workspace_id"`
	TemplateID   string `json:"template_id"`
	Title        string `json:"title"`
	Participants string `json:"participants"`

	Status vaults_domain.ChannelStatus `json:"status"`

	CreatedAt time.Time
	UpdatedAt time.Time

	RevokedAt     *time.Time
	PolicyCID     vaults_domain.Link `json:"policy_cid"`
	Slots         vaults_domain.Link `json:"slots"`
	Assignments   vaults_domain.Link `json:"assignments"`
	FederationCID vaults_domain.Link `json:"federation_cid"`
}
type ChannelSlotsNode struct {
	Slots []vaults_domain.Slot `json:"slots"`
}
type ChannelAssignmentsNode struct {
	Assignments []vaults_domain.Assignment `json:"assignments"`
}
type ChannelPolicyNode struct {
	Policy vaults_domain.Policy
}

// =======================================================================================
// WRITE
// =======================================================================================
func (s *VaultService) BuildChannelsBranch(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, error) {
	// =========================
	// 1. BUILD ENTRIES
	// =========================
	channels := vp.Collaborative.Channels
	channelLinks, err := s.BuildChannels(channels)
	if err != nil {
		return "", err
	}

	// =========================
	// 4. ENTRIES ROOT
	// =========================
	channelsCID, _, err := s.BuildChannelsRoot(channelLinks)
	if err != nil {
		return "", err
	}

	return channelsCID, nil
}

func (s *VaultService) BuildChannels(channels []vaults_domain.Channel) ([]vaults_domain.Link, error) {
	var links []vaults_domain.Link

	for _, channel := range channels {
		policyCID, err := s.BuildChannelPolicy(channel.Policy)
		if err != nil {
			return nil, err
		}

		slotsCID, err := s.BuildChannelSlots(channel.Slots)
		if err != nil {
			return nil, err
		}

		assignmentsCID, err := s.BuildChannelAssignments(channel.Assignments)
		if err != nil {
			return nil, err
		}

		node := ChannelNode{
			Version:       "1.0",
			ID:            channel.ID,
			WorkspaceID:   channel.WorkspaceID,
			TemplateID:    channel.TemplateID,
			Title:         channel.Title,
			Status:        channel.Status,
			Slots:         slotsCID,
			Assignments:   assignmentsCID,
			PolicyCID:     policyCID,
			FederationCID: vaults_domain.Link{channel.Federation},
			CreatedAt:     channel.CreatedAt,
			UpdatedAt:     channel.UpdatedAt,
			RevokedAt:     channel.RevokedAt,
		}

		cid, _, err := s.putNode(node)

		if err != nil {
			return nil, err
		}

		links = append(links, vaults_domain.Link{CID: cid})
	}

	return links, nil
}

func (s *VaultService) BuildChannelPolicy(policy vaults_domain.Policy) (vaults_domain.Link, error) {
	var link vaults_domain.Link

	node := ChannelPolicyNode{Policy: policy}

	cid, _, err := s.putNode(node)
	if err != nil {
		return link, err
	}

	link = vaults_domain.Link{CID: cid}

	return link, nil
}
func (s *VaultService) BuildChannelSlots(slots []vaults_domain.Slot) (vaults_domain.Link, error) {
	var link vaults_domain.Link

	node := ChannelSlotsNode{
		Slots: slots,
	}

	cid, _, err := s.putNode(node)
	if err != nil {
		return link, err
	}

	link = vaults_domain.Link{CID: cid}

	return link, nil
}

func (s *VaultService) BuildChannelAssignments(assignments []vaults_domain.Assignment) (vaults_domain.Link, error) {
	var link vaults_domain.Link

	node := ChannelAssignmentsNode{
		Assignments: assignments,
	}

	cid, _, err := s.putNode(node)
	if err != nil {
		return link, err
	}

	link = vaults_domain.Link{CID: cid}

	return link, nil
}

func (s *VaultService) BuildChannelsRoot(links []vaults_domain.Link) (string, int, error) {
	root := vaults_domain.ChannelsRoot{Items: links}
	return s.putNode(root)
}

func (s *VaultService) RotateChannelsKeys(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, error) {

	for i := range vp.Collaborative.Channels {
		vp.Collaborative.Channels[i].IsDirty = true
	}
	// 	↓
	cid, err := s.BuildChannelsBranch(session, vp, mode)
	if err != nil {
		return "", err
	}
	return cid, nil
}

// =======================================================================================
// READ
// =======================================================================================
func (r *VaultReconstructor) resolveChannels(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	channelsRoot vaults_domain.ChannelsRoot,
) ([]vaults_domain.Channel, error) {

	var result []vaults_domain.Channel

	for _, link := range channelsRoot.Items {

		res, err := r.Query.Execute(ctx, cmd.WithCID(link.CID))
		if err != nil {
			return nil, err
		}

		var channelNode ChannelNode
		errChannelNode := json.Unmarshal(res.Raw, &channelNode)
		if errChannelNode != nil {
			return nil, errChannelNode
		}

		slots, err := r.resolveSlots(ctx, cmd, channelNode.Slots)
		if err != nil {
			return nil, err
		}
		assignments, err := r.resolveAssignments(ctx, cmd, channelNode.Assignments)
		if err != nil {
			return nil, err
		}

		channel := vaults_domain.Channel{
			ID:          channelNode.ID,
			TemplateID:  channelNode.TemplateID,
			WorkspaceID: channelNode.WorkspaceID,
			Title:       channelNode.Title,
			Status:      channelNode.Status,
			Slots:       slots,
			Assignments: assignments,
			Federation:  channelNode.FederationCID.CID,
		}

		result = append(result, channel)
	}

	return result, nil
}

func (r *VaultReconstructor) resolveSlots(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	root vaults_domain.Link,
) ([]vaults_domain.Slot,error){

	var result []vaults_domain.Slot

	res, err := r.Query.Execute(ctx, cmd.WithCID(root.CID))
	if err != nil {
		return result, err
	}

	var slotsRoot vaults_domain.SlotsRoot

	if err := json.Unmarshal(res.Raw,&slotsRoot); err != nil {
		return result,err
	}


	for _,link := range slotsRoot.Items {

		slotRes,err := r.Query.Execute(
			ctx,
			cmd.WithCID(link.CID),
		)

		if err != nil {
			return result,err
		}


		var slot vaults_domain.Slot

		if err := json.Unmarshal(slotRes.Raw,&slot); err != nil {
			return result,err
		}

		result=append(result,slot)
	}

	return result,nil
}
func (r *VaultReconstructor) resolveAssignments(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	root vaults_domain.Link,
) ([]vaults_domain.Assignment,error){

	var result []vaults_domain.Assignment

	res, err := r.Query.Execute(ctx, cmd.WithCID(root.CID))
	if err != nil {
		return result, err
	}

	var assignmentsRoot vaults_domain.AssignmentsRoot

	if err := json.Unmarshal(res.Raw,&assignmentsRoot); err != nil {
		return result,err
	}


	for _,link := range assignmentsRoot.Items {

		assignmentRes,err := r.Query.Execute(
			ctx,
			cmd.WithCID(link.CID),
		)

		if err != nil {
			return result,err
		}


		var assignment vaults_domain.Assignment

		if err := json.Unmarshal(assignmentRes.Raw,&assignment); err != nil {
			return result,err
		}

		result=append(result,assignment)
	}

	return result,nil
}

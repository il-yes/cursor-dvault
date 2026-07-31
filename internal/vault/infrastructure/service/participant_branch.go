package vaults_service

import (
	"context"
	"encoding/json"

	vault_queries "vault-app/internal/vault/application/queries"
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
)

type ParticipantNode struct {
	Version     string
	ChannelID   string   `json:"channel_id"`
	VaultID     string   `json:"vault_id"`
	PublicKey   string   `json:"public_key"`
	Direction   string   `json:"direction"` // inbound | outbound | bidirectional
	Role        string   `json:"role"`
	Permissions []string `json:"joined_at"`
	JoinedAt    int64    `json:"permissions"`
}

// =======================================================================================
// WRITE
// =======================================================================================
func (s *VaultService) BuildParticipantsBranch(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, error) {
	// =========================
	// 1. BUILD ENTRIES
	// =========================
	participants := vp.Collaborative.Participants
	participantLinks, err := s.BuildParticipants(participants)
	if err != nil {
		return "", err
	}

	// =========================
	// 4. ENTRIES ROOT
	// =========================
	participantsCID, _, err := s.BuildParticipantsRoot(participantLinks)
	if err != nil {
		return "", err
	}

	return participantsCID, nil
}

func (s *VaultService) BuildParticipants(participants []vaults_domain.Participant) ([]vaults_domain.Link, error) {
	var links []vaults_domain.Link

	for _, participant := range participants {

		node := ParticipantNode{
			Version:     "v1.0.0",
			ChannelID:   participant.ChannelID,
			VaultID:     participant.VaultID,
			PublicKey:   participant.PublicKey,
			Direction:   participant.Direction,
			Role:        participant.Role,
			Permissions: participant.Permissions,
			JoinedAt:    participant.JoinedAt,
		}

		cid, _, err := s.putNode(node)
		if err != nil {
			return nil, err
		}

		links = append(links, vaults_domain.Link{CID: cid})
	}

	return links, nil
}

func (s *VaultService) BuildParticipantsRoot(links []vaults_domain.Link) (string, int, error) {
	root := vaults_domain.ParticipantsRoot{Items: links}
	return s.putNode(root)
}

func (s *VaultService) RotateParticipantKeys(session vault_session.Session, vp vaults_domain.VaultPayload, mode SyncMode) (string, error) {

	for i := range vp.Collaborative.Participants {
		vp.Collaborative.Participants[i].IsDirty = true
	}
	// 	↓
	cid, err := s.BuildParticipantsBranch(session, vp, mode)
	if err != nil {
		return "", err
	}
	return cid, nil
}

// =======================================================================================
// READ
// =======================================================================================
func (r *VaultReconstructor) resolveParticipants(
	ctx context.Context,
	cmd vault_queries.GetIPFSDataQuerry,
	participantsRoot vaults_domain.ParticipantsRoot,
) ([]vaults_domain.Participant, error) {

	var result []vaults_domain.Participant

	for _, link := range participantsRoot.Items {

		res, err := r.Query.Execute(ctx, cmd.WithCID(link.CID))
		if err != nil {
			return result, err
		}

		var participant vaults_domain.Participant
		if err := json.Unmarshal(res.Raw, &participant); err != nil {
			return result, err
		}

		result = append(result, participant)
	}

	return result, nil
}

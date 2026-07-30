package vaults_service

import (
	vault_session "vault-app/internal/vault/application/session"
	vaults_domain "vault-app/internal/vault/domain"
)

type ParticipantNode struct {
	Version string

	ChannelID string
	VaultID   string

	PublicKey string

	Direction string

	Role string

	Permissions []string

	JoinedAt int64
}

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
			Version: "v1.0.0",
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

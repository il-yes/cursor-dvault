package share_entry_infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/logger/logger"
	share_entry_domain "vault-app/internal/share_entry/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vaults_service "vault-app/internal/vault/infrastructure/service"
	vault_ui "vault-app/internal/vault/ui"
)

// ---------------------------------------------------------------------------------
// Interfaces
// ---------------------------------------------------------------------------------
type Vaulthandler interface {
	UploadAttachementToIPFSWithEncryption(userID string, ur vault_ui.UploadAttachRequest) (string, error)
	LoadAttachment(userID string, vaultName string, hash string, formatReturned string) (*vaults_service.LoadAttachmentResponse, error)
}

type EntrySnapshotService struct {
	Logger       logger.Logger
	VaultHandler Vaulthandler
}

func NewEntrySnapshotService(
	logger logger.Logger,
	vaultHandler Vaulthandler,
) *EntrySnapshotService {
	return &EntrySnapshotService{
		Logger:       logger,
		VaultHandler: vaultHandler,
	}
}

type BuildRequest struct {
	// Required for the share
	Share     *share_entry_domain.ShareEntry
	Recipient share_entry_domain.Recipient

	// Per‑user context (vault context)
	UserID             string
	UserSubscriptionID string
	UserOnboardingID   string
	VaultName          string
	Password           string
	SymKey             []byte
	VaultSession       vaults_domain.VaultPayload
	Configs            app_config_domain.Config
}
type BuildResponse struct {
	Raw         []byte
	Snapshot    share_entry_domain.EntrySnapshot
	Attachments []vaults_domain.Attachment
}

func (s *EntrySnapshotService) Build(
	ctx context.Context,
	req BuildRequest,
) (BuildResponse, error) {
	// 1. Process attachments under the given context
	updatedSnapshot, attachments, err := s.Process(ctx, req)
	if err != nil {
		return BuildResponse{}, err
	}

	// 2. Marshal the final, CID‑aware snapshot
	raw, err := json.Marshal(updatedSnapshot)
	if err != nil {
		return BuildResponse{}, err
	}

	return BuildResponse{
		Snapshot:    *updatedSnapshot,
		Raw:         raw,
		Attachments: attachments,
	}, nil
}

// Extract attachements from entry snapshot
func (s *EntrySnapshotService) Process(
	ctx context.Context,
	req BuildRequest,
) (*share_entry_domain.EntrySnapshot, []vaults_domain.Attachment, error) {
	entrySnapshot := req.Share.EntrySnapshot

	if len(entrySnapshot.AttachmentCIDs) == 0 {
		s.Logger.LogPretty("EntrySnapshotService - Process - AttachmentCIDs is nil --> create empty guard", entrySnapshot)
		entrySnapshot.AttachmentCIDs = make([]string, 0)
	}

	// Get attachments from cids
	attachmentsForShare := s.getAttachmentsForShare(entrySnapshot.AttachmentCIDs, req.VaultSession.Attachments)
	entrySnapshot.Attachments = nil
	s.Logger.LogPretty("EntrySnapshotService - Process - attachmentsForShare", attachmentsForShare)

	for i, attachment := range attachmentsForShare {
		// Get local file
		if attachment.Hash != "" {
			resp, err := s.VaultHandler.LoadAttachment(
				req.UserID,
				req.VaultName,
				attachment.Hash,
				"bytes",
			)
			if err != nil {
				s.Logger.Error(
					"❌ Failed to load attachment for user %s: %v",
					req.UserID,
					err,
				)
				return nil, nil, err
			}

			if attachment.RecipientCIDs == nil {
				attachment.RecipientCIDs = make(map[string]string)
			}
			if len(req.Share.Recipients) == 0 {
				return nil, nil, fmt.Errorf("share recipients empty")
			}
			if resp == nil {
				return nil, nil, errors.New("attachment load returned nil")
			}

			// Get cid for each attachment: one public cid is enough for sharing unless special conditions
			cid, err := s.VaultHandler.UploadAttachementToIPFSWithEncryption(
				req.UserID,
				vault_ui.UploadAttachRequest{
					Configs:            req.Configs,
					Data:               []byte(resp.File),
					UserSubscriptionID: req.UserSubscriptionID,
					VaultName:          req.VaultName,
					Password:           req.Password,
					EncryptionMode:     string(vaults_domain.EncryptionPublic),
					SymKey:             req.SymKey,
					UserOnboarding:     req.UserOnboardingID,
				},
			)
			if err != nil {
				s.Logger.Error(
					"❌ Failed to upload attachment for user %s: %v",
					req.UserID,
					err,
				)
				return nil, nil, err
			}

			for _, recipUser := range req.Share.Recipients {
				// Update each attachment with the new cid
				attachmentsForShare[i].RecipientCIDs[recipUser.PublicKey] = cid
				s.Logger.LogPretty("EntrySnapshotService - Process - attachmentsForShare 1", attachmentsForShare[i])

				// entrySnapshot.RecipientCIDs[recipUser.PublicKey] = cid
				s.Logger.LogPretty("EntrySnapshotService - Process - attachmentsForShare 2", "ok")
			}
			// Attachments must be marked as dirty to pass the buiding root node check
			attachmentsForShare[i].IsDirty = true

			// Update snapshot with the new cid
			entrySnapshot.Attachments = append(entrySnapshot.Attachments, attachmentsForShare[i])
			s.Logger.LogPretty("EntrySnapshotService - Process - attachmentsForShare", entrySnapshot.Attachments)
		}
	}
	// Return snapshot & attachmentsForShare
	return &entrySnapshot, attachmentsForShare, nil
}

func (s *EntrySnapshotService) getAttachmentByNodeCid(nodeCid string, attachments []vaults_domain.Attachment) vaults_domain.Attachment {
	for _, a := range attachments {
		if a.NodeCID == nodeCid {
			return a
		}
	}
	return vaults_domain.Attachment{}
}
func (s *EntrySnapshotService) getAttachmentsForShare(nodeCids []string, attachments []vaults_domain.Attachment) []vaults_domain.Attachment {
	attachmentsForShare := []vaults_domain.Attachment{}
	for _, nodeCid := range nodeCids {
		attachment := s.getAttachmentByNodeCid(nodeCid, attachments)
		if attachment.NodeCID != "" {
			attachmentsForShare = append(attachmentsForShare, attachment)
		}
	}
	s.Logger.LogPretty("EntrySnapshotService - getAttachmentsForShare - nodeCids", nodeCids)
	return attachmentsForShare
}

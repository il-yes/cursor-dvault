package share_infrastructure

import (
	"context"
	"encoding/json"
	share_domain "vault-app/internal/domain/shared"
	"vault-app/internal/logger/logger"
	vaults_domain "vault-app/internal/vault/domain"
	vault_ui "vault-app/internal/vault/ui"
)

// ---------------------------------------------------------------------------------
// Interfaces
// ---------------------------------------------------------------------------------
type Vaulthandler interface {
	UploadAttachementToIPFSWithEncryption(userID string, ur vault_ui.UploadAttachRequest) (string, error)
	LoadAttachment(userID string, vaultName string, hash string) (string, error)
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
	Share     *share_domain.ShareEntry
	Recipient share_domain.Recipient

	// Per‑user context (vault context)
	UserID             string
	UserSubscriptionID string
	VaultName          string
	Password           string
}
type BuildResponse struct {
	Raw            []byte
	Snapshot       share_domain.EntrySnapshot
	CIDs           []string
	AttachementIDs []string
}

func (s *EntrySnapshotService) Build(
	ctx context.Context,
	req BuildRequest,
) (BuildResponse, error) {
	// 1. Process attachments under the given context
	updatedSnapshot, cids, attachementIDs, err := s.Process(ctx, req)
	if err != nil {
		return BuildResponse{}, err
	}

	// 2. Marshal the final, CID‑aware snapshot
	raw, err := json.Marshal(updatedSnapshot)
	if err != nil {
		return BuildResponse{}, err
	}

	return BuildResponse{
		Snapshot:       *updatedSnapshot,
		Raw:            raw,
		CIDs:           cids,
		AttachementIDs: attachementIDs,
	}, nil
}




// Extract attachements from entry snapshot
func (s *EntrySnapshotService) Process(
    ctx context.Context,
    req BuildRequest,
) (*share_domain.EntrySnapshot, []string, []string, error) {
    entrySnapshot := req.Share.EntrySnapshot

    if entrySnapshot.Attachements == nil {
        entrySnapshot.Attachements = make([]vaults_domain.Attachment, 0, len(entrySnapshot.Attachements))
    }

    cids := []string{}
    attachementIDs := []string{}

    for i := range entrySnapshot.Attachements {
        attachment := &entrySnapshot.Attachements[i]
        if attachment.HashShare == "" {
            buffer, err := s.VaultHandler.LoadAttachment(
                req.UserID,
                req.VaultName,
                attachment.Hash,
            )
            if err != nil {
                s.Logger.Error(
                    "❌ Failed to load attachment for user %s: %v",
                    req.UserID,
                    err,
                )
                return nil, nil, nil, err
            }

            cid, err := s.VaultHandler.UploadAttachementToIPFSWithEncryption(
                req.UserID,
                vault_ui.UploadAttachRequest{
                    Data:               []byte(buffer),
                    UserSubscriptionID: req.UserSubscriptionID,
                    VaultName:          req.VaultName,
                    Password:           req.Password,
                },
            )
            if err != nil {
                s.Logger.Error(
                    "❌ Failed to upload attachment for user %s: %v",
                    req.UserID,
                    err,
                )
                return nil, nil, nil, err
            }

            attachment.HashShare = cid
            cids = append(cids, cid)
            attachementIDs = append(attachementIDs, attachment.ID)
        }
    }

    return &entrySnapshot, cids, attachementIDs, nil
}

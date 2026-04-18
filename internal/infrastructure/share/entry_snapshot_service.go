package share_infrastructure

import (
	"context"
	"encoding/json"
	share_domain "vault-app/internal/domain/shared"
	"vault-app/internal/logger/logger"
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
    Share              *share_domain.ShareEntry
    Recipient          share_domain.Recipient

    // Per‑user context (vault context)
    UserID             string
    UserSubscriptionID string
    VaultName          string
    Password           string
}

func (s *EntrySnapshotService) Build(
    ctx context.Context,
    req BuildRequest,
) ([]byte, error) {
    // 1. Process attachments under the given context
    updatedSnapshot, err := s.Process(ctx, req)
    if err != nil {
        return nil, err
    }

    // 2. Marshal the final, CID‑aware snapshot
    return json.Marshal(updatedSnapshot)
}

// Extract attachements from entry snapshot
func (s *EntrySnapshotService) Process(
    ctx context.Context,
    req BuildRequest,
) (*share_domain.EntrySnapshot, error) {
    entrySnapshot := req.Share.EntrySnapshot

    // range over attachments by index so we can mutate them
    for i := range entrySnapshot.Attachements {
        attachment := &entrySnapshot.Attachements[i]

        // Use request context:
        // - req.UserID, req.VaultName, req.Password
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
            return nil, err
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
            return nil, err
        }

        attachment.CIDShared = cid
    }

    return &entrySnapshot, nil
}


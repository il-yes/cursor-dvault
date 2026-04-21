package share_application_use_cases

import (
	"context"
	"fmt"

	// "log"
	"strconv"
	"time"
	share_application_dto "vault-app/internal/application"
	share_application_events "vault-app/internal/application/events/share"
	blockchain "vault-app/internal/blockchain"
	app_config_ui "vault-app/internal/config/ui"
	share_domain "vault-app/internal/domain/shared"
	share_infrastructure "vault-app/internal/infrastructure/share"
	"vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
	utils "vault-app/internal/utils"
	vaults_domain "vault-app/internal/vault/domain"
)

// ---------------------------------------------------------
//
//	Interfaces
//
// ---------------------------------------------------------
type TracecoreClientInterface interface {
	CreateShare(ctx context.Context, payload tracecore.ProdCreateCryptoShareRequest) (*tracecore.ProdCreateCryptoShareResponse, error)
	AcceptShare(ctx context.Context, shareID string) error
	RejectShare(ctx context.Context, shareID string) error
	GetShareByMe(ctx context.Context, email string) ([]share_domain.ShareEntry, error)
	GetShareWithMe(ctx context.Context, email string) ([]share_domain.ShareEntry, error)
	SetToken(token string)
	AddRecipient(ctx context.Context, req tracecore_types.AddRecipientRequest) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error)
	UpdateRecipient(ctx context.Context, req share_application_dto.UpdateRecipientRequest) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error)
	RevokeShare(ctx context.Context, req tracecore_types.RevokeShareRequest) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error)
}

type ClientCryptoService interface {
	GenerateSymmetricKey() []byte
	EncryptPayload(string, []byte) (blockchain.CryptoPayload, error)
	AESEncrypt(plain []byte, key []byte) blockchain.CryptoPayload
	AESDecrypt(enc []byte, key []byte) blockchain.CryptoPayload
}

type EntrySnapshotServiceInterface interface {
	Build(ctx context.Context, req share_infrastructure.BuildRequest) (share_infrastructure.BuildResponse, error)
}

// ---------------------------------------------------------
//
//	Cryptographic Share Use Case
//
// ---------------------------------------------------------
type ShareUseCase struct {
	repo                 share_domain.Repository
	dispatcher           share_application_events.EventDispatcher
	tc                   TracecoreClientInterface // new cloud client
	crypto               ClientCryptoService
	EntrySnapshotService EntrySnapshotServiceInterface
}

func NewShareUseCase(
	repo share_domain.Repository,
	tc TracecoreClientInterface,
	d share_application_events.EventDispatcher,
	crypto ClientCryptoService,
	entrySnapshotService EntrySnapshotServiceInterface,
) *ShareUseCase {
	return &ShareUseCase{
		repo:                 repo,
		tc:                   tc,
		dispatcher:           d,
		crypto:               crypto,
		EntrySnapshotService: entrySnapshotService,
	}
}

// ---------------------------------------------------------
// Create Share
// ---------------------------------------------------------
func (uc *ShareUseCase) CreateProdShareMode(
	ctx context.Context,
	userID string,
	ownerEmail string,
	share share_domain.ShareEntry,
	configFacade app_config_ui.AppConfigHandler,
	secret string,
	vault *vaults_domain.Vault,
) (*share_domain.ShareEntry, error) {
	// ---------------------------------------------------------
	// 1. Create share Request
	// ---------------------------------------------------------
	pcr, attachementsAdded, err := uc.BuildProdShareRequest(
		uc.crypto,
		userID,
		ownerEmail,
		share,
		configFacade,
		secret,
		vault,
	)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("share - ShareUseCase - pcr", pcr)

	// ---------------------------------------------------------
	// 2. send to Ankhora cloud
	// ---------------------------------------------------------
	createdRes, err := uc.tc.CreateShare(ctx, *pcr)
	if err != nil {
		return nil, fmt.Errorf("cloud CreateShare failed: %w", err)
	}
	utils.LogPretty("share - ShareUseCase - createdRes", createdRes)

	// ---------------------------------------------------------
	// 3. Publish event after commit
	// ---------------------------------------------------------
	uc.dispatcher.Dispatch(share_domain.ShareCreated{
		BaseEvent: share_domain.BaseEvent{
			Name: "ShareCreated",
			Time: time.Now(),
		},
		ShareID:        share.ID,
		EntryName:      share.EntryName,
		EntryType:      share.EntryType,
		OwnerID:        userID,
		CIDs:           attachementsAdded.CIDs,
		AttachementIDs: attachementsAdded.AttachementIDs,
	})

	return &share, nil
}

func (uc *ShareUseCase) BuildProdShareRequest(
	crypto ClientCryptoService,
	userID string,
	email string,
	share share_domain.ShareEntry,
	configFacade app_config_ui.AppConfigHandler,
	secret string,
	vault *vaults_domain.Vault,
) (*tracecore.ProdCreateCryptoShareRequest, *share_application_dto.AttachementCIDsAdded, error) {
	// ---------------------------------------------------------
	// 1. Generate symmetric key
	// ---------------------------------------------------------
	symKey := crypto.GenerateSymmetricKey()

	// ---------------------------------------------------------
	// 2. Build entry snapshot
	// ---------------------------------------------------------
	buildResponse, err := uc.EntrySnapshotService.Build(
		context.Background(),
		share_infrastructure.BuildRequest{
			Share:     &share,
			UserID:    userID,
			UserSubscriptionID: vault.UserSubscriptionID,
			VaultName: vault.Name,
			Password:  "password",
		})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build entry snapshot: %w", err)
	}
	utils.LogPretty("share - ShareUseCase - encryptedPayloadBeta", buildResponse)
	share.EntrySnapshot = buildResponse.Snapshot

	// ---------------------------------------------------------
	// 2. Encrypt payload
	// ---------------------------------------------------------
	// entrySnapshotRawJson, err := json.Marshal(buildResponse.Snapshot)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to marshal entry snapshot: %w", err)
	// }
	encryptedPayload := crypto.AESEncrypt(buildResponse.Raw, symKey)
	if encryptedPayload.Encrypted == nil {
		return nil, nil, fmt.Errorf("failed to encrypt payload")
	}

	// ---------------------------------------------------------
	// 3. Encrypt keys
	// ---------------------------------------------------------

	encryptedKeys := make(map[string]string)
	recipients := make(map[string]tracecore.CryptoRecipient, 0)

	for _, rid := range share.Recipients {
		encKey, err := crypto.EncryptPayload(rid.PublicKey, symKey)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to encrypt key: %w", err)
		}

		encryptedKeys[rid.Email] = encKey.ToString()
	}

	for _, rid := range share.Recipients {
		encKey, err := crypto.EncryptPayload(rid.PublicKey, symKey)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to encrypt key: %w", err)
		}

		recipients[rid.Email] = tracecore.CryptoRecipient{
			RevokedAt:     nil,
			EncryptedKeys: encKey.ToString(),
			Role:          rid.Role,
		}
	}
	// ---------------------------------------------------------
	// 4. Sign share
	// ---------------------------------------------------------
	// fetch userr private key from db
	userCfg, err := configFacade.GetUserConfigByUserID(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user config: %w", err)
	}

	message := "share.Message" // TODO: improve
	signature, err := blockchain.SignActorWithStellarPrivateKey(string(userCfg.StellarAccount.PrivateKey), message)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign share: %w", err)
	}

	// ---------------------------------------------------------
	// 5. Return request
	// ---------------------------------------------------------
	return &tracecore.ProdCreateCryptoShareRequest{
			SenderID:      share.OwnerID,
			SenderEmail:   email,
			Recipients:    recipients,
			VaultPayload:  encryptedPayload.ToString(),
			EncryptedKeys: encryptedKeys,
			Title:         share.EntryName,
			EntryType:     share.EntryType,
			AccessMode:    share.AccessMode,
			ExpiresAt:     share.ExpiresAt,
			// Metadata:      share.Metadata,
			PublicKey:       userCfg.StellarAccount.PublicKey,
			Signature:       signature,
			Message:         message,
			DownloadAllowed: share.DownloadAllowed,
		},
		&share_application_dto.AttachementCIDsAdded{
			CIDs:           buildResponse.CIDs,
			AttachementIDs: buildResponse.AttachementIDs,
		}, nil
}

// ------------------------------------------------
// Use case: list shared entries by the user
// ------------------------------------------------
func (s *ShareUseCase) ListSharedEntries(ctx context.Context, email string) ([]share_domain.ShareEntry, error) {
	// Mirror to cloud if client available
	cloudShares, err := s.tc.GetShareByMe(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("dvault ListReceivedShares failed: %w", err)
	}

	return cloudShares, nil
}

// ------------------------------------------------
// Use case: fetch shares *received* with the user
// ------------------------------------------------
func (s *ShareUseCase) ListReceivedShares(ctx context.Context, email string) ([]share_domain.ShareEntry, error) {
	cloudShares, err := s.tc.GetShareWithMe(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("dvault ListReceivedShares failed: %w", err)
	}

	return cloudShares, nil
}

func (uc *ShareUseCase) GetShareForAccept(
	ctx context.Context,
	shareID string,
	recipientUserID string,
) (*share_domain.ShareAcceptData, error) {

	share, recipient, blob, err :=
		uc.repo.GetShareForAccept(shareID, recipientUserID)

	if err != nil {
		return nil, fmt.Errorf("cannot get share for accept: %w", err)
	}

	return &share_domain.ShareAcceptData{
		Share:     *share,
		Recipient: *recipient,
		Blob:      blob,
	}, nil
}

type AcceptShareResult struct {
	Share     share_domain.ShareEntry
	Recipient share_domain.Recipient
	Blob      []byte // encrypted payload for this user
}

// ---------------------------------------------------------
// Accept Share Invitation
// ---------------------------------------------------------
func (uc *ShareUseCase) AcceptShare(ctx context.Context, shareID string, userID string) (*AcceptShareResult, error) {

	// 1. Load share entry + recipient-specific data
	share, recipient, err := uc.repo.GetShareAndRecipient(ctx, shareID, userID)
	if err != nil {
		return nil, err
	}

	// 2. Check expiration
	if share.ExpiresAt != nil && share.ExpiresAt.Before(time.Now()) {
		return nil, share_domain.ErrShareExpired
	}

	// 3. Mark recipient accepted
	if err := uc.repo.MarkRecipientAccepted(ctx, recipient.ID); err != nil {
		return nil, fmt.Errorf("failed to accept share: %w", err)
	}

	// 4. Return data to caller (VaultHandler → frontend)
	return &AcceptShareResult{
		Share:     *share,
		Recipient: *recipient,
		Blob:      recipient.EncryptedBlob,
	}, nil
}

type RejectShareResult struct {
	ShareID     string
	RecipientID string
	Message     string
}

// ---------------------------------------------------------
// Reject Share Invitation
// ---------------------------------------------------------
func (uc *ShareUseCase) RejectShare(ctx context.Context, shareID string, userID string) (*RejectShareResult, error) {

	// Load share + recipient
	_, recipient, err := uc.repo.GetShareAndRecipient(ctx, shareID, userID)
	if err != nil {
		return nil, err
	}

	// Mark the recipient invitation as "rejected"
	if err := uc.repo.MarkRecipientRejected(ctx, recipient.ID); err != nil {
		return nil, fmt.Errorf("failed to reject share: %w", err)
	}

	return &RejectShareResult{
		ShareID:     shareID,
		RecipientID: recipient.ID,
		Message:     "Share invitation rejected",
	}, nil
}

// ---------------------------------------------------------
// Add Receiver
// ---------------------------------------------------------
type AddReceiverInput struct {
	ShareID string
	Name    string
	Email   string
	Role    string
}

type AddReceiverResult struct {
	ShareID     string
	RecipientID string
	Message     string
}

func (uc *ShareUseCase) AddReceiver(ctx context.Context, requesterID string, in AddReceiverInput) (*AddReceiverResult, error) {

	// Load share
	share, err := uc.repo.GetShareByID(ctx, in.ShareID)

	if err != nil {
		return nil, fmt.Errorf("share not found: %w", err)
	}

	// Domain rule: only owner can add recipients
	if !share_domain.CanAddRecipient(share, requesterID) {
		return nil, fmt.Errorf("permission denied: not share owner")
	}

	// Create new recipient
	newRecipient := &share_domain.Recipient{
		ShareID:   in.ShareID,
		Name:      in.Name,
		Email:     in.Email,
		Role:      in.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		JoinedAt:  time.Now(),
	}

	if err := uc.repo.CreateRecipient(ctx, newRecipient); err != nil {
		return nil, fmt.Errorf("failed to add recipient: %w", err)
	}
	// share_domain.RecipientAdded event
	uc.dispatcher.Dispatch(share_domain.RecipientAdded{
		BaseEvent: share_domain.BaseEvent{
			Name: "RecipientAdded",
			Time: time.Now(),
		},
		ShareID:     share.ID,
		RecipientID: newRecipient.ID,
		Email:       newRecipient.Email,
	})

	return &AddReceiverResult{
		ShareID:     in.ShareID,
		RecipientID: newRecipient.ID,
		Message:     "Recipient added successfully",
	}, nil
}
func stringToUint(str string) uint {
	u64, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	return uint(u64)
}

// ---------------------------------------------------------
// Add Recipient (Cloud)
// ---------------------------------------------------------
func (uc *ShareUseCase) AddRecipient(
	ctx context.Context,
	userID string,
	in share_application_dto.AddRecipientRequest,
	configFacade app_config_ui.AppConfigHandler,
	secret string,
) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	// ---------------------------------------------------------
	// 1. Create add recipient Request
	// ---------------------------------------------------------
	addRecip, err := uc.BuildAddRecipientRequest(uc.crypto, userID, in, configFacade, secret)
	if err != nil {
		return nil, fmt.Errorf("❌ ShareUseCase - AddRecipient: failed to build add recipient request: %w", err)
	}
	// ---------------------------------------------------------
	// 2. Add recipient to cloud
	// ---------------------------------------------------------
	response, err := uc.tc.AddRecipient(ctx, *addRecip)
	if err != nil {
		return nil, fmt.Errorf("❌ ShareUseCase - AddRecipient: failed to add recipient: %w", err)
	}
	return response, nil
}

func (uc *ShareUseCase) BuildAddRecipientRequest(
	crypto ClientCryptoService,
	userID string,
	in share_application_dto.AddRecipientRequest,
	configFacade app_config_ui.AppConfigHandler,
	secret string,
) (*tracecore_types.AddRecipientRequest, error) {
	// ---------------------------------------------------------
	// 1. Generate symmetric key
	// ---------------------------------------------------------
	symKey := crypto.GenerateSymmetricKey()

	// ---------------------------------------------------------
	// 2. Encrypt keys
	// ---------------------------------------------------------

	encKey, err := crypto.EncryptPayload(in.PublicKey, symKey)
	if err != nil {
		return nil, fmt.Errorf("❌ ShareUseCase - BuildAddRecipientRequest: failed to encrypt key: %w", err)
	}
	encryptedKey := encKey.ToString()

	// ---------------------------------------------------------
	// 4. Sign share
	// ---------------------------------------------------------
	// fetch userr private key from db
	userCfg, err := configFacade.GetUserConfigByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("❌ ShareUseCase - BuildAddRecipientRequest: failed to get user config: %w", err)
	}

	message := "add.recipient" // TODO: improve
	signature, err := blockchain.SignActorWithStellarPrivateKey(string(userCfg.StellarAccount.PrivateKey), message)
	if err != nil {
		return nil, fmt.Errorf("❌ ShareUseCase - BuildAddRecipientRequest: failed to sign share: %w", err)
	}

	return &tracecore_types.AddRecipientRequest{
		ShareID:      in.ShareID,
		Email:        in.Email,
		Role:         in.Role,
		EncryptedKey: encryptedKey,
		RevokedAt:    in.RevokedAt,
		Signature:    signature,
	}, nil
}

// ---------------------------------------------------------
// Update Recipient (Cloud)
// ---------------------------------------------------------
func (uc *ShareUseCase) UpdateRecipient(ctx context.Context, requesterID string, in share_application_dto.UpdateRecipientRequest) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	utils.LogPretty("ShareUseCase - UpdateRecipient: updating recipient: %v", in)
	response, err := uc.tc.UpdateRecipient(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("❌ ShareUseCase - UpdateRecipient: failed to update recipient: %w", err)
	}
	return response, nil
}

func (uc *ShareUseCase) RevokeShare(ctx context.Context, requesterID string, in share_application_dto.UpdateRecipientRequest, configFacade app_config_ui.AppConfigHandler) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	// ---------------------------------------------------------
	// 1. Sign share
	// ---------------------------------------------------------
	// fetch userr private key from db
	userCfg, err := configFacade.GetUserConfigByUserID(requesterID)
	if err != nil {
		return nil, fmt.Errorf("❌ ShareUseCase - BuildAddRecipientRequest: failed to get user config: %w", err)
	}

	message := "revoke.share" // TODO: improve
	signature, err := blockchain.SignActorWithStellarPrivateKey(string(userCfg.StellarAccount.PrivateKey), message)
	if err != nil {
		return nil, fmt.Errorf("❌ ShareUseCase - BuildAddRecipientRequest: failed to sign share: %w", err)
	}

	input := tracecore_types.RevokeShareRequest{
		Challenge: message,
		Email:     in.Email,
		ShareID:   in.ShareID,
		Signature: signature,
	}
	response, err := uc.tc.RevokeShare(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("❌ ShareUseCase - RevokeShare: failed to revoke share: %w", err)
	}
	return response, nil
}

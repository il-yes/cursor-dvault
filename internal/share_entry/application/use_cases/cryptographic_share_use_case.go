package share_entry_use_cases

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	// "log"
	"strconv"
	"time"

	"vault-app/internal/blockchain"
	app_config_domain "vault-app/internal/config/domain"
	share_entry_application_dto "vault-app/internal/share_entry/application"
	share_entry_application_events "vault-app/internal/share_entry/application/events"
	share_entry_ports "vault-app/internal/share_entry/application/ports"
	share_entry_domain "vault-app/internal/share_entry/domain"
	share_entry_infrastructure "vault-app/internal/share_entry/infrastructure"
	"vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
	utils "vault-app/internal/utils"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

// ---------------------------------------------------------
//
//	Interfaces
//
// ---------------------------------------------------------
type TracecoreClientInterface interface {
	CreateShare(ctx context.Context, payload tracecore.ProdCreateCryptoShareRequest) (*tracecore.ProdCreateCryptoShareResponse, error)
	AcceptShare(ctx context.Context, req tracecore_types.ShareAcceptedPayload) (*tracecore_types.CloudResponse[tracecore_types.PendingShareIntent], error)
	RejectShare(ctx context.Context, req tracecore_types.ShareRejectedPayload) (*tracecore_types.CloudResponse[tracecore_types.PendingShareIntent], error)
	GetShareByMe(ctx context.Context, email string) ([]share_entry_domain.ShareEntry, error)
	GetShareWithMe(ctx context.Context, email string) ([]share_entry_domain.ShareEntry, error)
	SetToken(token string)
	AddRecipient(ctx context.Context, req tracecore_types.AddRecipientRequest) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error)
	UpdateRecipient(ctx context.Context, req share_entry_application_dto.UpdateRecipientRequest) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error)
	RevokeShare(ctx context.Context, req tracecore_types.RevokeShareRequest) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error)
	ListPendingIntentSharesByMe(ctx context.Context, email string) (*tracecore_types.CloudResponse[[]tracecore_types.PendingShareIntent], error)
	ListPendingIntentSharesWithMe(ctx context.Context, email string) (*tracecore_types.CloudResponse[[]tracecore_types.PendingShareIntent], error)
}

type ClientCryptoService interface {
	GenerateSymmetricKey() []byte
	EncryptPayload(string, []byte) (blockchain.CryptoPayload, error)
	AESEncrypt(plain []byte, key []byte) blockchain.CryptoPayload
	AESDecrypt(enc []byte, key []byte) blockchain.CryptoPayload
}

type EntrySnapshotServiceInterface interface {
	Build(ctx context.Context, req share_entry_infrastructure.BuildRequest) (share_entry_infrastructure.BuildResponse, error)
}

// ---------------------------------------------------------
//
//	Cryptographic Share Use Case
//
// ---------------------------------------------------------
type ShareUseCase struct {
	repo                 share_entry_domain.Repository
	dispatcher           share_entry_application_events.EventDispatcher
	tc                   TracecoreClientInterface // new cloud client
	crypto               ClientCryptoService
	aesService           *vault_infrastructure_crypto.AESService
	EntrySnapshotService EntrySnapshotServiceInterface
}

func NewShareUseCaseAES(
	repo share_entry_domain.Repository,
	tc TracecoreClientInterface,
	d share_entry_application_events.EventDispatcher,
	crypto *vault_infrastructure_crypto.AESService,
	entrySnapshotService EntrySnapshotServiceInterface,
) *ShareUseCase {
	return &ShareUseCase{
		repo:                 repo,
		tc:                   tc,
		dispatcher:           d,
		aesService:           crypto,
		EntrySnapshotService: entrySnapshotService,
	}
}

// ---------------------------------------------------------
// Create Share
// ---------------------------------------------------------
func (uc *ShareUseCase) Create(
	ctx context.Context,
	userID string,
	ownerEmail string,
	share share_entry_domain.ShareEntry,
	configFacade share_entry_ports.AppConfigHandlerInterface,
	secret string,
	vault *vaults_domain.Vault,
	vp vaults_domain.VaultPayload,
	userOnboardingID string,
	configs app_config_domain.Config,
	userSubscriptionID string,
) (*share_entry_domain.ShareEntry, error) {
	// ---------------------------------------------------------
	// 1. Create share Request
	// ---------------------------------------------------------
	pcr, attachementsAdded, err := uc.BuildProdShareRequest(
		uc.aesService,
		userID,
		ownerEmail,
		share,
		configFacade,
		secret,
		vault,
		vp,
		userOnboardingID,
		configs,
		userSubscriptionID,
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
	uc.dispatcher.Dispatch(share_entry_domain.ShareCreated{
		BaseEvent: share_entry_domain.BaseEvent{
			Name: "ShareCreated",
			Time: time.Now(),
		},
		ShareID:      share.ID,
		EntryName:    share.EntryName,
		EntryType:    share.EntryType,
		OwnerID:      userID,
		OwnerEmail: ownerEmail,
		CIDs:         attachementsAdded.CIDs,
		Attachements: attachementsAdded.Attachments,
	})

	return &share, nil
}

func (uc *ShareUseCase) BuildProdShareRequest(
	crypto *vault_infrastructure_crypto.AESService,
	userID string,
	email string,
	share share_entry_domain.ShareEntry,
	configFacade share_entry_ports.AppConfigHandlerInterface,
	secret string,
	vault *vaults_domain.Vault,
	vp vaults_domain.VaultPayload,
	userOnboardingID string,
	configs app_config_domain.Config,
	userSubscriptionID string,
) (*tracecore.ProdCreateCryptoShareRequest, *share_entry_application_dto.AttachementCIDsAdded, error) {
	// ---------------------------------------------------------
	// 1. Generate symmetric key
	// ---------------------------------------------------------
	as := vault_infrastructure_crypto.AsymmetricService{}
	symKey := as.GenerateSymmetricKey()


	utils.LogPretty("share - ShareUseCase - vault", vault)

	// ---------------------------------------------------------
	// 2. Build entry snapshot
	// ---------------------------------------------------------
	buildResponse, err := uc.EntrySnapshotService.Build(
		context.Background(),
		share_entry_infrastructure.BuildRequest{
			Share:              &share,
			UserID:             userID,
			UserSubscriptionID: userSubscriptionID,
			UserOnboardingID: userOnboardingID,
			VaultName:          vault.Name,
			Password:           "password",
			SymKey:             symKey,
			VaultSession: vp,
			Configs: configs,
		})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build entry snapshot: %w", err)
	}
	share.EntrySnapshot = buildResponse.Snapshot // for owner read only
	share.EntrySnapshot.Attachments = buildResponse.Attachments
	utils.LogPretty("share - ShareUseCase - EntrySnapshot", buildResponse)

	// ---------------------------------------------------------
	// 2. Encrypt payload
	// ---------------------------------------------------------
	encryptedPayload, err := crypto.Encrypt(buildResponse.Raw, symKey) // for recipient
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encrypt payload")
	}

	// ---------------------------------------------------------
	// 3. Encrypt keys
	// ---------------------------------------------------------
	encryptedKeys := make(map[string]string)
	recipients := make(map[string]tracecore.CryptoRecipient, 0)
	var primaryTrustGroupID string

	for _, rid := range share.Recipients {
		var str string

		if rid.RecipientType == "trust_group" || rid.TrustGroupID != "" {
			// TrustGroup mode: use TrustGroup ID as target key
			targetKey := rid.TrustGroupID
			if targetKey == "" {
				targetKey = rid.Email // Fallback if passed in Email field
			}
			primaryTrustGroupID = targetKey

			if rid.PublicKey != "" {
				encKey, err := crypto.EncryptPayload(rid.PublicKey, symKey)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to encrypt key for trust group: %w", err)
				}
				str = encKey.ToString()
			}

			encryptedKeys[targetKey] = str
			recipients[targetKey] = tracecore.CryptoRecipient{
				ID:            rid.ID,
				EncryptedKeys: str,
				Role:          rid.Role,
				TrustGroupID:  targetKey,
				RecipientType: "trust_group",
			}
		} else {
			// Personal User mode: existing user recipient logic
			if rid.PublicKey != "" {
				encKey, err := crypto.EncryptPayload(rid.PublicKey, symKey)
				if err != nil {
					return nil, nil, err
				}
				str = encKey.ToString()
			}
			encryptedKeys[rid.Email] = str

			recipients[rid.Email] = tracecore.CryptoRecipient{
				ID:            rid.ID,
				EncryptedKeys: str,
				Role:          rid.Role,
				RecipientType: "user",
			}
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
			SenderID:        share.OwnerID,
			SenderEmail:     email,
			Recipients:      recipients,
			VaultPayload:    base64.StdEncoding.EncodeToString(encryptedPayload),
			EncryptedKeys:   encryptedKeys,
			TrustGroupID:    primaryTrustGroupID,
			Title:           share.EntryName,
			EntryType:       share.EntryType,
			AccessMode:      share.AccessMode,
			ExpiresAt:       share.ExpiresAt,
			PublicKey:       userCfg.StellarAccount.PublicKey,
			Signature:       signature,
			Message:         message,
			DownloadAllowed: share.DownloadAllowed,
		},
		&share_entry_application_dto.AttachementCIDsAdded{
			Attachments: buildResponse.Attachments,
		}, nil
}

// ------------------------------------------------
// Use case: list shared entries by the user
// ------------------------------------------------
func (s *ShareUseCase) ListSharedEntries(ctx context.Context, email string) ([]share_entry_domain.ShareEntry, error) {
	if s.tc == nil {
		utils.LogPretty("share_entry_use_case - ListSharedEntries - tc is nil", nil)
		return nil, errors.New("tc is nil")
	}
	// Mirror to cloud if client available
	cloudShares, err := s.tc.GetShareByMe(ctx, email)
	if err != nil {
		utils.LogPretty("share_entry_use_case - ListSharedEntries - tc.GetShareByMe error: %v\n", err)
		return nil, fmt.Errorf("dvault ListReceivedShares failed: %w", err)
	}

	return cloudShares, nil
}

// ------------------------------------------------
// Use case: fetch shares *received* with the user
// ------------------------------------------------
func (s *ShareUseCase) ListReceivedShares(ctx context.Context, email string) ([]share_entry_domain.ShareEntry, error) {
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
) (*share_entry_domain.ShareAcceptData, error) {

	share, recipient, blob, err :=
		uc.repo.GetShareForAccept(shareID, recipientUserID)

	if err != nil {
		return nil, fmt.Errorf("cannot get share for accept: %w", err)
	}

	return &share_entry_domain.ShareAcceptData{
		Share:     *share,
		Recipient: *recipient,
		Blob:      blob,
	}, nil
}

type AcceptShareResult struct {
	Share     share_entry_domain.ShareEntry
	Recipient share_entry_domain.Recipient
	Blob      []byte // encrypted payload for this user
}

// ---------------------------------------------------------
// Accept Share Invitation
// ---------------------------------------------------------
func (uc *ShareUseCase) AcceptShare(ctx context.Context, shareID string, intentID string, email string) (*tracecore_types.CloudResponse[tracecore_types.PendingShareIntent], error) {
	acceptedResponse, err := uc.tc.AcceptShare(ctx, tracecore_types.ShareAcceptedPayload{
		IntentID:       intentID,
		ShareID:        shareID,
		RecipientEmail: email,
	})
	if err != nil {
		utils.LogPretty("share_entry_use_case - AcceptShare - tc.AcceptShare error: %v\n", err)
		return nil, fmt.Errorf("dvault AcceptShare failed: %w", err)
	}

	// 4. Return data to caller (VaultHandler → frontend)
	return acceptedResponse, nil
}

type RejectShareResult struct {
	ShareID     string
	RecipientID string
	Message     string
}

// ---------------------------------------------------------
// Reject Share Invitation
// ---------------------------------------------------------
func (uc *ShareUseCase) RejectShare(ctx context.Context, shareID string, intentID string, email string) (*tracecore_types.CloudResponse[tracecore_types.PendingShareIntent], error) {
	rejectResponse, err := uc.tc.RejectShare(ctx, tracecore_types.ShareRejectedPayload{
		IntentID:       intentID,
		ShareID:        shareID,
		RecipientEmail: email,
	})
	if err != nil {
		utils.LogPretty("share_entry_use_case - RejectShare - tc.RejectShare error: %v\n", err)
		return nil, fmt.Errorf("dvault RejectShare failed: %w", err)
	}

	return rejectResponse, nil
}

// ---------------------------------------------------------
// Revoke Share Invitation
// ---------------------------------------------------------
func (uc *ShareUseCase) RevokeShare(ctx context.Context, requesterID string, in share_entry_application_dto.UpdateRecipientRequest, configFacade share_entry_ports.AppConfigHandlerInterface) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	// ---------------------------------------------------------
	// 1. Sign share
	// ---------------------------------------------------------
	// fetch userr private key from db
	userCfg, err := configFacade.GetUserConfigByUserID(requesterID)
	if err != nil {
		return nil, fmt.Errorf("❌ ShareUseCase - RevokeShare: failed to get user config: %w", err)
	}

	message := "revoke.share" // TODO: improve
	signature, err := blockchain.SignActorWithStellarPrivateKey(string(userCfg.StellarAccount.PrivateKey), message)
	if err != nil {
		return nil, fmt.Errorf("❌ ShareUseCase - RevokeShare: failed to sign share: %w", err)
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
	if !share_entry_domain.CanAddRecipient(share, requesterID) {
		return nil, fmt.Errorf("permission denied: not share owner")
	}

	// Create new recipient
	newRecipient := &share_entry_domain.Recipient{
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
	// share_entry_domain.RecipientAdded event
	uc.dispatcher.Dispatch(share_entry_domain.RecipientAdded{
		BaseEvent: share_entry_domain.BaseEvent{
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
	in share_entry_application_dto.AddRecipientRequest,
	configFacade share_entry_ports.AppConfigHandlerInterface,
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
	in share_entry_application_dto.AddRecipientRequest,
	configFacade share_entry_ports.AppConfigHandlerInterface,
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
func (uc *ShareUseCase) UpdateRecipient(ctx context.Context, requesterID string, in share_entry_application_dto.UpdateRecipientRequest) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	utils.LogPretty("ShareUseCase - UpdateRecipient: updating recipient: %v", in)
	response, err := uc.tc.UpdateRecipient(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("❌ ShareUseCase - UpdateRecipient: failed to update recipient: %w", err)
	}
	return response, nil
}

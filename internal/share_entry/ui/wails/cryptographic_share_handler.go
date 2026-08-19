package sahre_entry_ui_wails

import (
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"

	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/logger/logger"
	share_entry_application_dto "vault-app/internal/share_entry/application"
	share_entry_application_events "vault-app/internal/share_entry/application/events"
	share_entry_ports "vault-app/internal/share_entry/application/ports"
	share_entry_application_use_cases "vault-app/internal/share_entry/application/use_cases"
	share_entry_domain "vault-app/internal/share_entry/domain"
	"vault-app/internal/tracecore"
	tracecore_types "vault-app/internal/tracecore/types"
	"vault-app/internal/utils"
)

type CryptographicShareHandler struct {
	CryptographicShareUseCase share_entry_application_use_cases.ShareUseCase
	Logger                    *logger.Logger
	Db                        gorm.DB
	TracecoreClient           *tracecore.TracecoreClient
	EventDispatcher           share_entry_application_events.EventDispatcher
}

func NewCryptographicShareHandler(
	cryptographicShareUseCase share_entry_application_use_cases.ShareUseCase,
	db gorm.DB,
	evtDispatcher share_entry_application_events.EventDispatcher,
	logger *logger.Logger,
	tracecoreClient *tracecore.TracecoreClient,
) CryptographicShareHandler {

	return CryptographicShareHandler{
		CryptographicShareUseCase: cryptographicShareUseCase,
		Db:                        db,
		EventDispatcher:           evtDispatcher,
		Logger:                    logger,
		TracecoreClient:           tracecoreClient,
	}
}

func (vh *CryptographicShareHandler) CreateShareEntry(
	ctx context.Context,
	payload share_entry_application_dto.CreateShareEntryPayload,
	ownerID string,
	ownerEmail string,
	configFacade share_entry_ports.AppConfigHandlerInterface,
	secret string,
	vaultHandler share_entry_ports.VaultHandlerInterface,
	tracecoreClient *tracecore.TracecoreClient,
	userOnboardingID string,
	configs app_config_domain.Config,
	userSubscriptionID string,
) (*share_entry_domain.ShareEntry, error) {
	// 1. Convert JSON string -> domain struct
	// ==================================================================================================================
	var snapshot share_entry_domain.EntrySnapshot
	if err := json.Unmarshal([]byte(payload.EntrySnapshot), &snapshot); err != nil {
		return nil, fmt.Errorf("invalid entry_snapshot: %w", err)
	}
	snapshot.AttachmentCIDs = payload.AttachmentCIDs

	// 2. create share entry via factory
	// ==================================================================================================================
	share := share_entry_domain.NewShareEntry(
		ownerID,
		payload.EntryID,
		payload.EntryName,
		payload.EntryRef,
		payload.EntryType,
		payload.Status,
		payload.AccessMode,
		payload.Encryption,
		snapshot,
		payload.DownloadAllowed,
		payload.ExpiresAt,
	)
	// vh.Logger.LogPretty("CryptographicShareHandler - CreateShareEntry - snapshot", snapshot)
	// vh.Logger.LogPretty("CryptographicShareHandler - CreateShareEntry - share", share)
	// 3. recipients
	// ==================================================================================================================
	recips := make([]share_entry_domain.Recipient, 0, len(payload.Recipients))
	for _, r := range payload.Recipients {
		recips = append(recips, share_entry_domain.Recipient{
			Name:          r.Name,
			Email:         r.Email,
			Role:          r.Role,
			PublicKey:     r.PublicKey,
			TrustGroupID:  r.TrustGroupID,
			RecipientType: r.Type,
		})
	}
	share.Recipients = recips

	// 4. get user's vault
	// ==================================================================================================================
	vault, err := vaultHandler.GetLatestByUserID(ownerID)
	if err != nil {
		return nil, err
	}

	vp, err := vaultHandler.GetVaultSession(ownerID)
	if err != nil {
		return nil, err
	}

	// 5. call usecase
	// ==================================================================================================================
	created, err := vh.CryptographicShareUseCase.Create(ctx, ownerID, ownerEmail, share, configFacade, secret, vault, *vp, userOnboardingID, configs, userSubscriptionID)
	if err != nil {
		return nil, err
	}
	// utils.LogPretty("Handler response - created", created)
	return created, nil
}

func (vh *CryptographicShareHandler) ListSharedEntries(ctx context.Context, email string) ([]share_entry_domain.ShareEntry, error) {
	if vh.TracecoreClient == nil {
		vh.Logger.LogPretty("share_entry_handler - ListSharedEntries - tracecoreClient is nil", nil)
		return nil, fmt.Errorf("tracecore client is not initialized")
	}
	utils.LogPretty("share_entry_handler - ListSharedEntries - email", email)
	res, err := vh.TracecoreClient.GetShareByMe(ctx, email)
	if err != nil {
		vh.Logger.LogPretty("share_entry_handler - ListSharedEntries - tracecoreClient.GetShareByMe error: %v\n", err)
		return nil, fmt.Errorf("failed fetching shared entries: %w", err)
	}

	return res, nil
}

func (vh *CryptographicShareHandler) ListReceivedShares(ctx context.Context, email string) ([]share_entry_domain.ShareEntry, error) {
	if vh.TracecoreClient == nil {
		vh.Logger.LogPretty("share_entry_handler - ListReceivedShares - tracecoreClient is nil", nil)
		return nil, fmt.Errorf("tracecore client is not initialized")
	}
	res, err := vh.TracecoreClient.GetShareWithMe(ctx, email)
	if err != nil {
		vh.Logger.LogPretty("share_entry_handler - ListReceivedShares - tracecoreClient.GetShareWithMe error: %v\n", err)
		return nil, fmt.Errorf("failed fetching shared entries: %w", err)
	}

	return res, nil
}

func (vh *CryptographicShareHandler) AddRecipient(
	ctx context.Context,
	userID string,
	in share_entry_application_dto.AddRecipientRequest,
	configFacade share_entry_ports.AppConfigHandlerInterface,
	secret string,
) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	result, err := vh.CryptographicShareUseCase.AddRecipient(ctx, userID, in, configFacade, secret)
	if err != nil {
		vh.Logger.Error("❌ CryptographicShareHandler - AddRecipient: Failed to add recipient: %v\n", err)
		return nil, err
	}

	vh.Logger.Info("✅ CryptographicShareHandler - AddRecipient: Successfully added recipient: %v\n", result)
	return result, nil
}

func (vh *CryptographicShareHandler) UpdateRecipient(
	ctx context.Context,
	userID string,
	in share_entry_application_dto.UpdateRecipientRequest,
) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	result, err := vh.CryptographicShareUseCase.UpdateRecipient(ctx, userID, in)
	if err != nil {
		vh.Logger.Error("❌ CryptographicShareHandler - UpdateRecipient: Failed to update recipient: %v\n", err)
		return nil, err
	}

	vh.Logger.Info("✅ CryptographicShareHandler - UpdateRecipient: Successfully updated recipient: %v\n", result)
	return result, nil
}
func (vh *CryptographicShareHandler) AcceptShare(ctx context.Context, shareID string, intentID string, email string) (*tracecore_types.CloudResponse[tracecore_types.PendingShareIntent], error) {
	result, err := vh.CryptographicShareUseCase.AcceptShare(ctx, shareID, intentID, email)
	if err != nil {
		vh.Logger.Error("❌ CryptographicShareHandler - AcceptShare: Failed to accept share: %v\n", err)
		return nil, err
	}

	vh.Logger.Info("✅ CryptographicShareHandler - AcceptShare: Successfully accepted share: %v\n", result)
	return result, nil
}
func (vh *CryptographicShareHandler) RejectShare(ctx context.Context, shareID string, intentID string, email string) (*tracecore_types.CloudResponse[tracecore_types.PendingShareIntent], error) {
	result, err := vh.CryptographicShareUseCase.RejectShare(ctx, shareID, intentID, email)
	if err != nil {
		vh.Logger.Error("❌ CryptographicShareHandler - RejectShare: Failed to reject share: %v\n", err)
		return nil, err
	}

	vh.Logger.Info("✅ CryptographicShareHandler - RejectShare: Successfully rejected share: %v\n", result)
	return result, nil
}
func (vh *CryptographicShareHandler) RevokeRecipient(ctx context.Context, userID string, in share_entry_application_dto.UpdateRecipientRequest) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	result, err := vh.CryptographicShareUseCase.UpdateRecipient(ctx, userID, in)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (vh *CryptographicShareHandler) RevokeShare(
	ctx context.Context,
	userID string,
	in share_entry_application_dto.UpdateRecipientRequest,
	configFacade share_entry_ports.AppConfigHandlerInterface,
) (*tracecore_types.CloudResponse[tracecore.CloudCryptographicShare], error) {
	result, err := vh.CryptographicShareUseCase.RevokeShare(ctx, userID, in, configFacade)
	if err != nil {
		vh.Logger.Error("❌ CryptographicShareHandler - RevokeShare: Failed to revoke share: %v\n", err)
		return nil, err
	}

	vh.Logger.Info("✅ CryptographicShareHandler - RevokeShare: Successfully revoked share: %v\n", result)
	return result, nil
}

// ---------------------------------------------------------
// List Pending Intent Shares (Cloud)
// ---------------------------------------------------------
func (vh *CryptographicShareHandler) ListPendingIntentSharesByMe(ctx context.Context, email string) (*tracecore_types.CloudResponse[[]tracecore_types.PendingShareIntent], error) {
	if vh.TracecoreClient == nil {
		vh.Logger.LogPretty("share_entry_handler - ListPendingIntentSharesByMe - tracecoreClient is nil", nil)
		return nil, fmt.Errorf("tracecore client is not initialized")
	}
	result, err := vh.TracecoreClient.ListPendingIntentSharesByMe(ctx, email)
	if err != nil {
		vh.Logger.LogPretty("share_entry_handler - ListPendingIntentSharesByMe - tracecoreClient.GetPendingShareIntents error: %v\n", err)
		return nil, fmt.Errorf("failed fetching shared entries: %w", err)
	}

	return result, nil
}
func (vh *CryptographicShareHandler) ListPendingIntentSharesWithMe(ctx context.Context, email string) (*tracecore_types.CloudResponse[[]tracecore_types.PendingShareIntent], error) {
	if vh.TracecoreClient == nil {
		vh.Logger.LogPretty("share_entry_handler - ListPendingIntentSharesWithMe - tracecoreClient is nil", nil)
		return nil, fmt.Errorf("tracecore client is not initialized")
	}
	result, err := vh.TracecoreClient.ListPendingIntentSharesWithMe(ctx, email)
	if err != nil {
		vh.Logger.LogPretty("share_entry_handler - ListPendingIntentSharesWithMe - tracecoreClient.GetPendingShareIntents error: %v\n", err)
		return nil, fmt.Errorf("failed fetching shared entries: %w", err)
	}

	return result, nil
}
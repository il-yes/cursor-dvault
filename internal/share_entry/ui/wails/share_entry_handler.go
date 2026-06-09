package sahre_entry_ui_wails

import (
	"context"

	"gorm.io/gorm"

	"vault-app/internal/blockchain"
	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/logger/logger"
	share_entry_application_dto "vault-app/internal/share_entry/application"
	share_entry_application_events "vault-app/internal/share_entry/application/events"
	share_entry_ports "vault-app/internal/share_entry/application/ports"
	share_entry_application_use_cases "vault-app/internal/share_entry/application/use_cases"
	share_entry_domain "vault-app/internal/share_entry/domain"
	share_entry_infrastructure "vault-app/internal/share_entry/infrastructure"
	"vault-app/internal/tracecore"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

type ShareEntryHandler struct {
	CryptographicShareHandler CryptographicShareHandler
	LinkShareHandler          LinkShareHandler
	Logger                    logger.Logger
	db                        *gorm.DB

	TracecoreClient *tracecore.TracecoreClient
	EventDispatcher share_entry_application_events.EventDispatcher
	VaultHandler    share_entry_ports.VaultHandlerInterface
}

func NewShareEntryHandler(
	tracecoreClient *tracecore.TracecoreClient,
	logger logger.Logger,
	db *gorm.DB,
	evtDispatcher share_entry_application_events.EventDispatcher,
	vaultHandler share_entry_ports.VaultHandlerInterface,
) ShareEntryHandler {
	if vaultHandler == nil {
		logger.LogPretty("share_entry_handler - NewShareEntryHandler - vaultHandler is nil", nil)
	}
	// Infra
	entrySnapshotService := share_entry_infrastructure.NewEntrySnapshotService(
		logger,
		vaultHandler,
	)
	cryptographicRepository := share_entry_infrastructure.NewGormShareRepository(db)

	linkShareRepository := share_entry_infrastructure.NewGormShareRepository(db)
	dispatcher := evtDispatcher

	if cryptographicRepository == nil {
		logger.LogPretty("share_entry_handler - NewShareEntryHandler - cryptographicRepository is nil", nil)
	}

	if dispatcher == nil {
		logger.LogPretty("share_entry_handler - NewShareEntryHandler - dispatcher is nil", nil)
	}
	if evtDispatcher == nil {
		logger.LogPretty("share_entry_handler - NewShareEntryHandler - evtDispatcher is nil", nil)
	}
	if entrySnapshotService == nil {
		logger.LogPretty("share_entry_handler - NewShareEntryHandler - entrySnapshotService is nil", nil)
	}

	// Application
	cryptoAESUC := share_entry_application_use_cases.NewShareUseCaseAES(cryptographicRepository, tracecoreClient, evtDispatcher, &vault_infrastructure_crypto.AESService{}, entrySnapshotService)

	if tracecoreClient == nil {
		logger.LogPretty("share_entry_handler - NewShareEntryHandler - tcClient is nil", nil)
	}
	if cryptoAESUC == nil {
		logger.LogPretty("share_entry_handler - NewShareEntryHandler - cryptoAESUC is nil", nil)
	}

	linkShareUseCase := share_entry_application_use_cases.NewLinkShareUseCase(linkShareRepository, tracecoreClient, dispatcher, &blockchain.CryptoService{})

	// Handler
	linkHandler := NewLinkShareHandler(*linkShareUseCase, &logger)
	cryptoHandler := NewCryptographicShareHandler(*cryptoAESUC, *db, evtDispatcher, &logger, tracecoreClient)

	if linkShareUseCase == nil {
		logger.LogPretty("share_entry_handler - NewShareEntryHandler - linkShareUseCase is nil", nil)
	}
	if linkShareRepository == nil {
		logger.LogPretty("share_entry_handler - NewShareEntryHandler - linkShareRepository is nil", nil)
	}

	return ShareEntryHandler{
		CryptographicShareHandler: cryptoHandler,
		LinkShareHandler:          linkHandler,
		VaultHandler:              vaultHandler,
		EventDispatcher:           evtDispatcher,
		db:                        db,
		Logger:                    logger,
		TracecoreClient:           tracecoreClient,
	}
}

func (ch *ShareEntryHandler) CreateCryptoShare(
	req share_entry_application_dto.CreateShareRequest,
	configFacade share_entry_ports.AppConfigHandlerInterface,
	vaultHandler share_entry_ports.VaultHandlerInterface,
	userOnboardingID string,
) (*share_entry_domain.ShareEntry, error) {

	return ch.CryptographicShareHandler.CreateShareEntry(
		context.Background(),
		req.Payload,
		req.UserID,
		req.OwnerEmail,
		configFacade,
		req.Secret,
		vaultHandler,
		ch.TracecoreClient,
		userOnboardingID,
		app_config_domain.Config{},
		"",
	)
}

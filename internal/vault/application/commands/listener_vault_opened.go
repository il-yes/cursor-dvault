package vault_commands

import (
	"context"
	"vault-app/internal/logger/logger"
	vault_events "vault-app/internal/vault/application/events"
	vault_session "vault-app/internal/vault/application/session"
	vault_domain "vault-app/internal/vault/domain"
)

// -------- REQUESTS --------
type AttachVaultRequest struct {
	UserID       string
	VaultPayload *vault_domain.VaultPayload
	Dirty        bool
	LastCID      string

	LastSynced  string
	LastUpdated string
	Runtime     *vault_session.RuntimeContext
}

// -------- INTERFACES --------
type VaultHandlerInterface interface {
	Open(ctx context.Context, req OpenVaultCommand, appConfigHandler AppConfigFacade) (*OpenVaultResult, error)
	SessionAttachVault(ctx context.Context, req AttachVaultRequest) error
	GetSession(userID string) (*vault_session.Session, error)
}

type ApplyOnboardingPacksWorkerInterface interface {
	OnApplyOnboardingPacks(userID, vaultName string, userOnboarding string)
}

// -------- LISTENER --------
type VaultOpenedListener struct {
	Logger       *logger.Logger
	Bus          vault_events.VaultEventBus
	VaultHandler VaultHandlerInterface
	PackWorker   ApplyOnboardingPacksWorkerInterface
}

func NewVaultOpenedListener(
	logger *logger.Logger,
	bus vault_events.VaultEventBus,
	vaultHandler VaultHandlerInterface,
	pw ApplyOnboardingPacksWorkerInterface,
) *VaultOpenedListener {
	return &VaultOpenedListener{
		Logger:       logger,
		Bus:          bus,
		VaultHandler: vaultHandler,
		PackWorker: pw,
	}
}

// -------- METHODS --------
func (l *VaultOpenedListener) Listen(ctx context.Context) {
	l.Logger.Info("Vault opened listener starting prrocessing...")
	l.Bus.SubscribeToVaultOpened(func(ctx context.Context, e vault_events.VaultOpened) {
		err := l.VaultHandler.SessionAttachVault(ctx, AttachVaultRequest{
			UserID:       e.UserID,
			VaultPayload: e.VaultPayload,
			Dirty:        e.Dirty,
			LastCID:      e.LastCID,
			LastSynced:   e.LastSynced,
			LastUpdated:  e.LastUpdated,
			Runtime:      e.Runtime,
		})
		if err != nil {
			l.Logger.Error("❌ VaultOpenedListener - failed to open vault for user %s: %v", e.UserID, err)
			return
		}
		l.Logger.Info("✅ VaultOpenedListener - vault opened for user: %s - %s", e.UserID, e.UserOnboardingID)

		// TODO: launch the packWorker with the runtime instead of the vaultHandler
		l.PackWorker.OnApplyOnboardingPacks(e.UserID, e.VaultName, e.UserOnboardingID)
	})

	<-ctx.Done()
	l.Logger.Warn("🛑 VaultOpenedListener stopped")
}

package vault_commands

import (
	"context"
	"vault-app/internal/logger/logger"
	vault_events "vault-app/internal/vault/application/events"
	vault_session "vault-app/internal/vault/application/session"
	vault_domain "vault-app/internal/vault/domain"
	app_config_ui "vault-app/internal/config/ui"
)

// -------- REQUESTS --------
type AttachVaultRequest struct {
	UserID  string
	VaultPayload *vault_domain.VaultPayload	
    Dirty   bool
    LastCID string

    LastSynced     string
    LastUpdated    string
    Runtime        *vault_session.RuntimeContext
}

// -------- INTERFACES --------
type VaultHandlerInterface interface {
	Open(ctx context.Context, req OpenVaultCommand, appConfigHandler app_config_ui.AppConfigHandler) (*OpenVaultResult, error)
	SessionAttachVault(ctx context.Context, req AttachVaultRequest) error
}


// -------- LISTENER --------
type VaultOpenedListener struct {
	Logger *logger.Logger
	Bus    vault_events.VaultEventBus	
	VaultHandler 	VaultHandlerInterface	
}

func NewVaultOpenedListener(
	logger *logger.Logger,
	bus vault_events.VaultEventBus,
	vaultHandler VaultHandlerInterface,
) *VaultOpenedListener {	
	return &VaultOpenedListener{
		Logger: logger,
		Bus: bus,
		VaultHandler: vaultHandler,
	}
}

// -------- METHODS --------
func (l *VaultOpenedListener) Listen(ctx context.Context) {
	l.Logger.Info("Vault opened listener starting prrocessing...")	
	l.Bus.SubscribeToVaultOpened(func(ctx context.Context, e vault_events.VaultOpened) {
		err := l.VaultHandler.SessionAttachVault(ctx, AttachVaultRequest{
			UserID: e.UserID,
			VaultPayload: e.VaultPayload,
			Dirty: e.Dirty,
			LastCID: e.LastCID,
			LastSynced: e.LastSynced,
			LastUpdated: e.LastUpdated,
			Runtime: e.Runtime,	
		})
		if err != nil {
			l.Logger.Error("❌ VaultOpenedListener - failed to open vault for user %s: %v", e.UserID, err)
			return
		}
		l.Logger.Info("✅ VaultOpenedListener - vault opened for user %s", e.UserID)
	})

	<-ctx.Done()
	l.Logger.Warn("🛑 VaultOpenedListener stopped")
}

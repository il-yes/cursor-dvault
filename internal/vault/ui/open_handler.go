package vault_ui

import (
	"context"
	app_config_ui "vault-app/internal/config/ui"
	vault_commands "vault-app/internal/vault/application/commands"
	vault_events "vault-app/internal/vault/application/events"
)


type OpenVaultHandler struct {
	openVaultCommandHandler *vault_commands.OpenVaultCommandHandler		
	EventBus                vault_events.VaultEventBus
}	

func NewOpenVaultHandler(openVaultCommandHandler *vault_commands.OpenVaultCommandHandler, eventBus vault_events.VaultEventBus) *OpenVaultHandler {
	return &OpenVaultHandler{
		openVaultCommandHandler: openVaultCommandHandler,
		EventBus:                eventBus,
	}
}

func (h *OpenVaultHandler) OpenVault(
	ctx context.Context, 
	req vault_commands.OpenVaultCommand, 
	appConfigHandler app_config_ui.AppConfigHandler,
) (*vault_commands.OpenVaultResult, error) {
	
	return h.openVaultCommandHandler.Handle(
		ctx, 
		req, 
		h.EventBus, 
		appConfigHandler,
	)
}


	
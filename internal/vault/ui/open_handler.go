package vault_ui

import (
	"context"
	"errors"
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
	appConfigHandler vault_commands.AppConfigFacade,
) (*vault_commands.OpenVaultResult, error) {
	if req.Session == nil {
		return nil, errors.New("session is required")
	}
	if req.UserID == "" {
		return nil, errors.New("user id is required")
	}
	if appConfigHandler == nil {
		return nil, errors.New("app config handler is required")
	}
	return h.openVaultCommandHandler.Handle(
		ctx,
		req,
		h.EventBus,
		appConfigHandler,
	)
}

package vault_ui

import (
	"context"
	"vault-app/internal/vault/application/commands"
)


type OpenVaultHandler struct {
	openVaultCommandHandler *vault_commands.OpenVaultCommandHandler		
}	

func NewOpenVaultHandler(openVaultCommandHandler *vault_commands.OpenVaultCommandHandler) *OpenVaultHandler {
	return &OpenVaultHandler{
		openVaultCommandHandler: openVaultCommandHandler,
	}
}

func (h *OpenVaultHandler) OpenVault(ctx context.Context, req vault_commands.OpenVaultCommand) (*vault_commands.OpenVaultResult, error) {
	return h.openVaultCommandHandler.Handle(ctx, req)
}



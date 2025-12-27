package vault_ui

import (
	"context"
	"vault-app/internal/blockchain"
	vault_commands "vault-app/internal/vault/application/commands"
	vault_events "vault-app/internal/vault/application/events"
)


type OpenVaultHandler struct {
	openVaultCommandHandler *vault_commands.OpenVaultCommandHandler		
	ipfs                    *blockchain.IPFSClient
	crypto                  *blockchain.CryptoService
	EventBus                vault_events.VaultEventBus
}	

func NewOpenVaultHandler(openVaultCommandHandler *vault_commands.OpenVaultCommandHandler, ipfs *blockchain.IPFSClient, crypto *blockchain.CryptoService, eventBus vault_events.VaultEventBus) *OpenVaultHandler {
	return &OpenVaultHandler{
		openVaultCommandHandler: openVaultCommandHandler,
		ipfs:                    ipfs,
		crypto:                  crypto,
		EventBus:                eventBus,
	}
}

func (h *OpenVaultHandler) OpenVault(ctx context.Context, req vault_commands.OpenVaultCommand) (*vault_commands.OpenVaultResult, error) {
	return h.openVaultCommandHandler.Handle(ctx, req, h.ipfs, h.crypto, h.EventBus)
}


	
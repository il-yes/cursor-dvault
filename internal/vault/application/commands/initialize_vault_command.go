package vault_commands

import (
	"fmt"
	utils "vault-app/internal"
	vault_domain "vault-app/internal/vault/domain"
)

// -------- COMMAND --------

type InitializeVaultCommand struct {
	UserID string
	VaultName string
}


// -------- RESULT --------
type InitializeVaultResult struct {
	Vault *vault_domain.Vault
}

// -------- HANDLER --------
type InitializeVaultCommandHandler struct {
	vaultRepo vault_domain.VaultRepository
}


// -------- CONSTRUCTOR --------
func NewInitializeVaultCommandHandler(vaultRepo vault_domain.VaultRepository) *InitializeVaultCommandHandler {
	return &InitializeVaultCommandHandler{
		vaultRepo: vaultRepo,
	}
}	

func (h *InitializeVaultCommandHandler) Execute(cmd InitializeVaultCommand) (*InitializeVaultResult, error) {
	// -----------------------------
	// 1. Validate & defaults -> can be done via value objects
	// -----------------------------
	if cmd.VaultName == "" {
		cmd.VaultName = cmd.UserID + "-vault"
	}
	
	// -----------------------------
	// 3. Init empty vault - idempotency	
	existing, err := h.vaultRepo.GetLatestByUserID(cmd.UserID)
    if err == nil && existing != nil {
		utils.LogPretty("InitializeVaultCommandHandler - vault found", existing)

        return &InitializeVaultResult{Vault: existing}, nil
    }
	utils.LogPretty("InitializeVaultCommandHandler - vault not found", cmd)

	// -----------------------------
	// 2. Save vault metadata to DB
	// -----------------------------	
	newVault := vault_domain.NewVault(cmd.UserID, cmd.VaultName)
	if err := h.vaultRepo.SaveVault(newVault); err != nil {
		return nil, fmt.Errorf("❌ failed to persist vault metadata: %w", err)
	}	
	utils.LogPretty("InitializeVaultCommandHandler - vault saved", newVault)
	
	return &InitializeVaultResult{Vault: newVault}, nil
}	



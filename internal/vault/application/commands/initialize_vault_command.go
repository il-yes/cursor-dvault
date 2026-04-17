package vault_commands

import (
	"errors"
	"fmt"
	utils "vault-app/internal/utils"
	vault_domain "vault-app/internal/vault/domain"
	vault_persistence "vault-app/internal/vault/infrastructure/persistence"

	"gorm.io/gorm"
)

// -------- COMMAND --------

type InitializeVaultCommand struct {
	UserID    string
	VaultName string
}

// -------- RESULT --------
type InitializeVaultResult struct {
	Vault *vault_domain.Vault
}

// -------- HANDLER --------
type InitializeVaultCommandHandler struct {
	VaultRepo vault_domain.VaultRepository
	DB        *gorm.DB
}

// -------- CONSTRUCTOR --------
func NewInitializeVaultCommandHandler(db *gorm.DB) *InitializeVaultCommandHandler {
	vaultRepo := vault_persistence.NewGormVaultRepository(db)

	return &InitializeVaultCommandHandler{
		VaultRepo: vaultRepo,
		DB:        db,
	}
}

func (h *InitializeVaultCommandHandler) Execute(cmd InitializeVaultCommand) (*InitializeVaultResult, error) {
	// -----------------------------
	// 1. Validate & defaults -> can be done via value objects
	// -----------------------------
	if cmd.VaultName == "" {
		cmd.VaultName = cmd.UserID + "-vault"
	}
	utils.LogPretty("InitializeVaultCommandHandler - Execute - cmd", cmd)

	// -----------------------------
	if h.VaultRepo == nil {
		return nil, errors.New("VaultRepo is nil")
	}
	// 3. Init empty vault - idempotency
	existing, err := h.VaultRepo.GetLatestByUserID(cmd.UserID)
	if existing != nil {
		utils.LogPretty("InitializeVaultCommandHandler - vault found", existing)

		return &InitializeVaultResult{Vault: existing}, nil
	}
	utils.LogPretty("InitializeVaultCommandHandler - vault not found mais il s'en tape !!!!", err)

	// -----------------------------
	// 2. Save vault metadata to DB
	// -----------------------------
	newVault := vault_domain.NewVault(cmd.UserID, cmd.VaultName)
	if err := h.VaultRepo.SaveVault(newVault); err != nil {
		return nil, fmt.Errorf("❌ failed to persist vault metadata: %w", err)
	}
	utils.LogPretty("InitializeVaultCommandHandler - vault saved", newVault)

	return &InitializeVaultResult{Vault: newVault}, nil
}

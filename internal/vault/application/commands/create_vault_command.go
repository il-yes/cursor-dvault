package vault_commands

import (
	"context"
	"errors"
	"fmt"
	"time"
	app_config_domain "vault-app/internal/config/domain"
	onboarding_domain "vault-app/internal/onboarding/domain"
	utils "vault-app/internal/utils"
	vault_domain "vault-app/internal/vault/domain"
	vaults_domain "vault-app/internal/vault/domain"
)

// -------- COMMAND query --------
type CreateVaultCommand struct {
	UserID             string
	VaultName          string
	Password           string
	UserSubscriptionID string
	AppConfig          app_config_domain.AppConfig
	UserOnboarding     *onboarding_domain.User
}

// -------- COMMAND result --------
type CreateVaultResult struct {
	Vault          *vault_domain.Vault
	ReusedExisting bool
}

// -------- COMMAND handler interfaces --------
type CryptoServiceInterface interface {
	Encrypt(data []byte, password string) ([]byte, error)
}

type IpfsServiceInterface interface {
	// AddData(data []byte) (string, error)
	Add(ctx context.Context, data []byte) (string, error)
}

type InitializeVaultHandler interface {
	Execute(cmd InitializeVaultCommand) (*InitializeVaultResult, error)
}

type CreateIPFSPayloadHandler interface {
	Execute(ctx context.Context, vc app_config_domain.VaultContext, cmd CreateIPFSPayloadCommand) (*CreateIPFSPayloadCommandResult, error)
	SetIpfsService(i IpfsServiceInterface)
}

// -------- COMMAND handler --------
type CreateVaultCommandHandler struct {
	initializeVaultHandler   InitializeVaultHandler
	createIPFSPayloadHandler CreateIPFSPayloadHandler
	vaultRepo                vault_domain.VaultRepository
}

type VaultInterfaceService interface {
	CreateVault(cmd CreateVaultCommand) (*CreateVaultResult, error)
}

// -------- COMMAND handler constructor --------
func NewCreateVaultCommandHandler(
	initializator InitializeVaultHandler,
	creator CreateIPFSPayloadHandler,
	vaultRepo vault_domain.VaultRepository,
) *CreateVaultCommandHandler {
	return &CreateVaultCommandHandler{
		initializeVaultHandler:   initializator,
		createIPFSPayloadHandler: creator,
		vaultRepo:                vaultRepo,
	}
}

func (h *CreateVaultCommandHandler) CreateVault(cmd CreateVaultCommand) (*CreateVaultResult, error) {
	utils.LogPretty("CreateVaultCommandHandler - CreateVault - cmd", cmd)
	if h.initializeVaultHandler == nil {
		return nil, errors.New("CreateVaultCommandHandler - CreateVault: initializeVaultHandler")
	}

	// -----------------------------
	// 1. Initialize vault
	// -----------------------------
	vault, err := h.initializeVaultHandler.Execute(InitializeVaultCommand{UserID: cmd.UserID, VaultName: cmd.VaultName})
	if err != nil {
		utils.LogPretty("CreateVaultCommandHandler - InitializeVaultHandler - Execute - 1st err", err)
		return nil, err
	}
	utils.LogPretty("CreateVaultCommandHandler - vault", vault)

	// -----------------------------
	// 1. Vault - Get vault content
	// -----------------------------
	const InitialVaultVersion = "1.0.0"
	vaultPayload := vault.Vault.BuildInitialPayload(InitialVaultVersion) // true for new user only
	utils.LogPretty("CreateIPFSPayloadCommandHandler - Execute - vaultPayload", vaultPayload)

	// -----------------------------
	// 2. Get vault content
	// -----------------------------
	vaultBytes, err := vaultPayload.GetContentBytes()
	if err != nil {
		utils.LogPretty("CreateVaultCommandHandler - InitializeVaultHandler - Execute - 2nd err", err)
		return nil, fmt.Errorf("❌ vault encryption failed: %w", err)
	}
	utils.LogPretty("CreateIPFSPayloadCommandHandler - Execute - vaultBytes", vaultBytes)

	// -----------------------------
	// 2. Create vault context
	// -----------------------------
	vc := app_config_domain.VaultContext{
		AppConfig:     cmd.AppConfig,
		StorageConfig: cmd.AppConfig.Storage,
		UserID:        cmd.UserSubscriptionID,
		VaultName:     cmd.VaultName,
	}

	if h.createIPFSPayloadHandler == nil {
		utils.LogPretty("CreateIPFSPayloadCommandHandler - Execute - createIPFSPayloadHandler", h.createIPFSPayloadHandler)
	}

	// -----------------------------
	// 2. Create IPFS CID
	// -----------------------------
	ipfsRecord, err := h.createIPFSPayloadHandler.Execute(
		context.Background(),
		vc,
		CreateIPFSPayloadCommand{
			Vault:    vault.Vault,
			Password: cmd.Password,
			Data:     vaultBytes,
			UserID: cmd.UserOnboarding.ID,
		})
	if err != nil {
		utils.LogPretty("CreateVaultCommandHandler - InitializeVaultHandler - Execute - 3rd err", err)
		return nil, err
	}
	utils.LogPretty("CreateVaultCommandHandler - ipfsRecord", ipfsRecord)

	// -----------------------------
	// 3. Update vault with IPFS CID
	// -----------------------------
	// vault.Vault.UserID = 
	vault.Vault.AttachCID(ipfsRecord.CID)
	vault.Vault.AttachUserSubscriptionID(cmd.UserSubscriptionID)
	vault.Vault.VaultMeta = vaults_domain.VaultMeta{
		Name: cmd.VaultName,
		UserID: cmd.UserID,
		CreatedAt: time.Now().Local().GoString(),
	}
	utils.LogPretty("CreateVaultCommandHandler - vault attached CID", vault.Vault)

	if vault.Vault == nil {
		utils.LogPretty("CreateVaultCommandHandler - InitializeVaultHandler - Execute - 4th err", err)
		return nil, errors.New("vault is nil before UpdateVault")
	}

	if err := h.vaultRepo.UpdateVault(vault.Vault); err != nil {
		utils.LogPretty("CreateVaultCommandHandler - InitializeVaultHandler - Execute - 5th err", err)
		return nil, err
	}

	// -----------------------------
	// 4. (Optional) Return event "VaultCreated"
	// -----------------------------

	return &CreateVaultResult{
		Vault:          vault.Vault,
		ReusedExisting: false,
	}, nil
}

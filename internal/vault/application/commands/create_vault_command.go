package vault_commands

import (
	"errors"
	utils "vault-app/internal"
	vault_domain "vault-app/internal/vault/domain"
)

// -------- COMMAND query --------
type CreateVaultCommand struct {
	UserID    string
	VaultName string
	Password  string
}

// -------- COMMAND result --------
type CreateVaultResult struct {
	Vault          *vault_domain.Vault
	ReusedExisting bool
}

// -------- COMMAND handler --------
// -------- COMMAND handler --------
type CreateVaultCommandHandler struct {
	initializeVaultHandler   InitializeVaultHandler
	createIPFSPayloadHandler CreateIPFSPayloadHandler
	vaultRepo                vault_domain.VaultRepository
}

// -------- COMMAND handler interfaces --------
type CryptoServiceInterface interface {
	Encrypt(data []byte, password string) ([]byte, error)
}

type IpfsServiceInterface interface {
	AddData(data []byte) (string, error)
}

type InitializeVaultHandler interface {
	Execute(cmd InitializeVaultCommand) (*InitializeVaultResult, error)
}

type CreateIPFSPayloadHandler interface {
	Execute(cmd CreateIPFSPayloadCommand) (*CreateIPFSPayloadCommandResult, error)
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
	// -----------------------------
	// 1. Initialize vault
	// -----------------------------
	vault, err := h.initializeVaultHandler.Execute(InitializeVaultCommand{UserID: cmd.UserID, VaultName: cmd.VaultName})
	if err != nil {
		return nil, err
	}
	utils.LogPretty("CreateVaultCommandHandler - vault", vault)

	// -----------------------------
	// 2. Create IPFS payload
	// -----------------------------
	ipfsRecord, err := h.createIPFSPayloadHandler.Execute(CreateIPFSPayloadCommand{Vault: vault.Vault, Password: cmd.Password})
	if err != nil {
		return nil, err
	}
	utils.LogPretty("CreateVaultCommandHandler - ipfsRecord", ipfsRecord)

	// -----------------------------
	// 3. Update vault with IPFS CID
	// -----------------------------
	vault.Vault.AttachCID(ipfsRecord.CID)
	utils.LogPretty("CreateVaultCommandHandler - vault attached CID", vault.Vault)

	if vault.Vault == nil {
		return nil, errors.New("vault is nil before UpdateVault")
	}

	if err := h.vaultRepo.UpdateVault(vault.Vault); err != nil {
		return nil, err
	}
	// utils.LogPretty("CreateVaultCommandHandler - vault", vault.Vault)

	// -----------------------------
	// 4. (Optional) Return event "VaultCreated"
	// -----------------------------

	return &CreateVaultResult{
		Vault:          vault.Vault,
		ReusedExisting: false,
	}, nil
}

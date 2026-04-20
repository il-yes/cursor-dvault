package vault_commands

import (
	"context"
	"fmt"
	blockchain_ipfs "vault-app/internal/blockchain/ipfs"
	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/tracecore"
	"vault-app/internal/utils"
	vault_dto "vault-app/internal/vault/application/dto"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

/*
// -------- COMMAND --------
type CreateIPFSPayloadCommand struct {
	Vault     *vaults_domain.Vault
	Password  string
	AppCfg    app_config_domain.AppConfig
	UserID    string
	VaultName string
	Data []byte
}

// -------- COMMAND response --------
type CreateIPFSPayloadCommandResult struct {
	CID string
}

// -------- COMMAND handler --------
type CreateIPFSPayloadCommandHandler struct {
	VaultRepo          vaults_domain.VaultRepository
	CryptoService      vaults_domain.VaultCrypto
	IpfsService        IpfsServiceInterface
	TracecoreClient    tracecore.TracecoreClient
	UnlockVaultHandler UnlockVaultHandlerInterface
	storageFactory StorageFactory
}

// -------- constructor --------
func NewCreateIPFSPayloadCommandHandler(vaultRepo vaults_domain.VaultRepository, tracecoreClient tracecore.TracecoreClient) *CreateIPFSPayloadCommandHandler {
	vc := &vault_infrastructure_crypto.AESService{}
	return &CreateIPFSPayloadCommandHandler{
		VaultRepo:       vaultRepo,
		CryptoService:   vc,
		TracecoreClient: tracecoreClient,
	}
}

func (h *CreateIPFSPayloadCommandHandler) Execute(cmd CreateIPFSPayloadCommand) (*CreateIPFSPayloadCommandResult, error) {
	// // -----------------------------
	// // 1. Vault - Get vault content
	// // -----------------------------
	// const InitialVaultVersion = "1.0.0"
	// vaultPayload := cmd.Vault.BuildInitialPayload(InitialVaultVersion) // true for new user only
	// utils.LogPretty("CreateIPFSPayloadCommandHandler - Execute - vaultPayload", vaultPayload)

	// // -----------------------------
	// // 2. CryptoEncrypt vault content
	// // -----------------------------
	// vaultBytes, err := vaultPayload.GetContentBytes()
	// if err != nil {
	// 	return nil, fmt.Errorf("❌ vault encryption failed: %w", err)
	// }
	// utils.LogPretty("CreateIPFSPayloadCommandHandler - Execute - vaultBytes", vaultBytes)

	// 1. unlock keyring
	// 2. get VaultKey (DEK)
	// 3. encrypt vault payload (entries, folders, index)
	unlockRes, err := h.UnlockVaultHandler.Execute(vault_dto.UnlockVaultCommand{
		Password: cmd.Password,
		// StellarSecret: cmd.StellarSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to unlock vault key: %w", err)
	}

	vaultKey := unlockRes.VaultKey.Key

	encrypted, err := h.CryptoService.Encrypt(
		cmd.Data ,
		vaultKey,
	)
	if err != nil {
		return nil, fmt.Errorf("❌ vault encryption failed: %w", err)
	}
	utils.LogPretty("CreateIPFSPayloadCommandHandler - encrypted length", len(encrypted))

	// -----------------------------
	// 3. IPFS - Add vault content to IPFS
	// -----------------------------
	cidFromIpfs, err := h.StoreOnIpfs(context.Background(), StoreIpfsParams{
		Data: encrypted,
	})
	if err != nil {
		return nil, fmt.Errorf("❌ failed to add vault to IPFS: %w", err)
	}
	utils.LogPretty("CreateIPFSPayloadCommandHandler - cidFromIpfs", cidFromIpfs)

	// -----------------------------
	// 4. Return result
	// -----------------------------
	return &CreateIPFSPayloadCommandResult{CID: cidFromIpfs}, nil
}

type StoreIpfsParams struct {
	Data []byte
}

func (h *CreateIPFSPayloadCommandHandler) StoreOnIpfs(ctx context.Context, req StoreIpfsParams) (string, error) {
	response, err := h.IpfsService.Add(ctx, req.Data)
	if err != nil {
		utils.LogPretty("CreateIPFSPayloadCommandHandler - StoreOnIpfs - error", err)
		return "", fmt.Errorf("failed to store ipfs: %w", err)
	}

	return response, nil
}
*/
func (h *CreateIPFSPayloadCommandHandler) SetIpfsService(ipfs IpfsServiceInterface) {
	h.IpfsService = ipfs
}




type CreateIPFSPayloadCommand struct {
	Vault    *vaults_domain.Vault
	Password string
	Data     []byte
	UserID	string // User app
}

// -------- COMMAND response --------
type CreateIPFSPayloadCommandResult struct {
	CID string
}

// -------- COMMAND handler --------
type CreateIPFSPayloadCommandHandler struct {
	VaultRepo          vaults_domain.VaultRepository
	CryptoService      vaults_domain.VaultCrypto
	IpfsService        IpfsServiceInterface
	TracecoreClient    tracecore.TracecoreClient
	UnlockVaultHandler UnlockVaultHandlerInterface
	StorageFactory blockchain_ipfs.StorageFactory
}

// -------- constructor --------
func NewCreateIPFSPayloadCommandHandler(
	vaultRepo vaults_domain.VaultRepository, 
	tracecoreClient tracecore.TracecoreClient,
	sf blockchain_ipfs.StorageFactory,
	uh UnlockVaultHandlerInterface,
) *CreateIPFSPayloadCommandHandler {
	vc := &vault_infrastructure_crypto.AESService{}
	return &CreateIPFSPayloadCommandHandler{
		VaultRepo:       vaultRepo,
		CryptoService:   vc,
		TracecoreClient: tracecoreClient,
		StorageFactory: sf,
		UnlockVaultHandler: uh,
	}
}
func (h *CreateIPFSPayloadCommandHandler) Execute(
	ctx context.Context,
	vaultCtx app_config_domain.VaultContext,
	cmd CreateIPFSPayloadCommand,
) (*CreateIPFSPayloadCommandResult, error) {
	// 1. Unlock vault key 
	// ==============================================
	unlockRes, err := h.UnlockVaultHandler.Execute(vault_dto.UnlockVaultCommand{
		Password: cmd.Password,
		UserID: vaultCtx.AppConfig.Branch,		// userOnboarding required
	})
	if err != nil {
		utils.LogPretty("CreateIPFSPayloadCommandHandler - Execute - - error", cmd)
		return nil, fmt.Errorf("CreateIPFSPayloadCommandHandler - Execute - failed to unlock vault key: %w", err)
	}

	vaultKey := unlockRes.VaultKey.Key

	// 2. Encryption 
	// ==============================================
	encrypted, err := h.CryptoService.Encrypt(cmd.Data, vaultKey)
	if err != nil {
		return nil, fmt.Errorf("CreateIPFSPayloadCommandHandler - Execute - vault encryption failed: %w", err)
	}

	// IPFS Upload 
	// ==============================================
	cidFromIpfs, err := h.StoreOnIpfs(ctx, vaultCtx, encrypted)
	if err != nil {
		return nil, fmt.Errorf("CreateIPFSPayloadCommandHandler - Execute - failed to add vault to IPFS: %w", err)
	}

	return &CreateIPFSPayloadCommandResult{CID: cidFromIpfs}, nil
}

func (h *CreateIPFSPayloadCommandHandler) StoreOnIpfs(
	ctx context.Context,
	vaultCtx app_config_domain.VaultContext,
	data []byte,
) (string, error) {
	storageProvider := h.StorageFactory.New(&vaultCtx)
	return storageProvider.Add(ctx, data)
}
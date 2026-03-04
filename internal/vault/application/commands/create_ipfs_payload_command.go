package vault_commands

import (
	"context"
	"fmt"
	"vault-app/internal/blockchain"
	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/tracecore"
	vaults_domain "vault-app/internal/vault/domain"
)

// -------- COMMAND --------
type CreateIPFSPayloadCommand struct {
	Vault *vaults_domain.Vault	
	Password string
}	

// -------- COMMAND response --------
type CreateIPFSPayloadCommandResult struct {
	CID string
}

// -------- COMMAND handler --------
type CreateIPFSPayloadCommandHandler struct {
	vaultRepo vaults_domain.VaultRepository
	CryptoService CryptoServiceInterface
	IpfsService IpfsServiceInterface
	TracecoreClient tracecore.TracecoreClient
}



// -------- constructor --------	
func NewCreateIPFSPayloadCommandHandler(vaultRepo vaults_domain.VaultRepository, cryptoService CryptoServiceInterface, ipfsService IpfsServiceInterface, tracecoreClient tracecore.TracecoreClient) *CreateIPFSPayloadCommandHandler {
	return &CreateIPFSPayloadCommandHandler{
		vaultRepo: vaultRepo,
		CryptoService: cryptoService,
		IpfsService: ipfsService,
		TracecoreClient: tracecoreClient,
	}
}

func (h *CreateIPFSPayloadCommandHandler) Execute(cmd CreateIPFSPayloadCommand) (*CreateIPFSPayloadCommandResult, error) {
	// -----------------------------
	// 1. Vault - Get vault content
	// -----------------------------
	const InitialVaultVersion = "1.0.0"
	vaultPayload := cmd.Vault.BuildInitialPayload(InitialVaultVersion) // true for new user only
		
	// -----------------------------
	// 2. CryptoEncrypt vault content
	// -----------------------------
	vaultBytes, err := vaultPayload.GetContentBytes()
	if err != nil {
		return nil, fmt.Errorf("❌ vault encryption failed: %w", err)
	}

	encrypted, err := h.CryptoService.Encrypt(vaultBytes, cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("❌ vault encryption failed: %w", err)
	}
	
	// -----------------------------
	// 3. IPFS - Add vault content to IPFS
	// -----------------------------
	cidFromIpfs, err := h.IpfsService.AddData(encrypted)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to add vault to IPFS: %w", err)
	}
	
	// -----------------------------
	// 4. Return result
	// -----------------------------
	return &CreateIPFSPayloadCommandResult{CID: cidFromIpfs}, nil	
}

type StoreIpfsParams struct {
	AppCfg app_config_domain.AppConfig
	UserID    string
	VaultName string
	Data      []byte
}
func (h *CreateIPFSPayloadCommandHandler) StoreOnIpfs(ctx context.Context, req StoreIpfsParams, tc tracecore.TracecoreClient) (string, error) {
	// ------------------------------------------------------------
	// 1. LOAD TRACECORE CLIENT
	// ------------------------------------------------------------
	tracecoreClient := tracecore.NewTracecoreFromConfig(&req.AppCfg, "token")	
	// ------------------------------------------------------------
	// 2. LOAD STORAGE PROVIDER
	// ------------------------------------------------------------
	storageProvider := blockchain.NewStorageProvider(blockchain.Config{
		StorageConfig: req.AppCfg.Storage,
		UserID:             req.UserID,
		VaultName:          req.VaultName,
	}, tracecoreClient)
	// ------------------------------------------------------------
	// 3. ADD TO IPFS
	// ------------------------------------------------------------
	response, err := storageProvider.Add(ctx, req.Data)
	if err != nil {
		return "", fmt.Errorf("failed to store ipfs: %w", err)
	}

	return response, nil
}
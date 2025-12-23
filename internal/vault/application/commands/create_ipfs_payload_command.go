package vault_commands

import (
	"fmt"
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
}



// -------- constructor --------	
func NewCreateIPFSPayloadCommandHandler(vaultRepo vaults_domain.VaultRepository, cryptoService CryptoServiceInterface, ipfsService IpfsServiceInterface) *CreateIPFSPayloadCommandHandler {
	return &CreateIPFSPayloadCommandHandler{
		vaultRepo: vaultRepo,
		CryptoService: cryptoService,
		IpfsService: ipfsService,
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


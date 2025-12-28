package vault_queries

import (
	"fmt"
	vault_domain "vault-app/internal/vault/domain"
	vaults_domain "vault-app/internal/vault/domain"
)

// -------- QUERRY --------
type GetIPFSDataQuerry struct {
	CID string
	Password string
}
// -------- RESPONSE --------
type GetIPFSDataResponse struct {
	Data vault_domain.VaultPayload
}
// -------- HANDLER --------
type GetIPFSDataQuerryHandler struct {
	ipfsService IpfsServiceInterface
	cryptoService CryptoServiceInterface
}

type IpfsServiceInterface interface {
	GetFile(cid string) ([]byte, error)
}

type CryptoServiceInterface interface {
	Decrypt(data []byte, password string) ([]byte, error)
}

// -------- CONSTRUCTOR --------
func NewGetIPFSDataQuerryHandler(ipfsService IpfsServiceInterface, cryptoService CryptoServiceInterface) *GetIPFSDataQuerryHandler {
	return &GetIPFSDataQuerryHandler{
		ipfsService: ipfsService,
		cryptoService: cryptoService,
	}
}

func (h *GetIPFSDataQuerryHandler) Execute(cmd GetIPFSDataQuerry) (*GetIPFSDataResponse, error) {
	
	// Fetch vault from IPFS
	vaultData, err := h.ipfsService.GetFile(cmd.CID)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to fetch vault from IPFS: %w", err)
	}
	if vaultData == nil || len(vaultData) == 0 {
		return nil, fmt.Errorf("❌ empty vault data for CID %s", cmd.CID)
	}

	// Decrypt vault
	decrypted, err := h.cryptoService.Decrypt(vaultData, cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("❌ failed to decrypt vault: %w", err)
	}
	if len(decrypted) == 0 {
		return nil, fmt.Errorf("❌ vault decryption returned empty result")
	}

	// Parse vault
	vaultPayload := vaults_domain.ParseVaultPayload(decrypted)

	return &GetIPFSDataResponse{Data: vaultPayload}, nil
}


package vault_queries

import (
	"context"
	"fmt"
	"vault-app/internal/blockchain"
	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/tracecore"
	"vault-app/internal/utils"
	vault_domain "vault-app/internal/vault/domain"
	vaults_domain "vault-app/internal/vault/domain"
)

// -------- QUERRY --------
type GetIPFSDataQuerry struct {
	CID string
	Password string
	AppCfg app_config_domain.AppConfig
	UserID string
	VaultName string
}
// -------- RESPONSE --------
type GetIPFSDataResponse struct {
	Data vault_domain.VaultPayload
}
// -------- HANDLER --------
type GetIPFSDataQuerryHandler struct {
	ipfsService IpfsServiceInterface
	cryptoService CryptoServiceInterface
	tracecoreClient tracecore.TracecoreClient
}

type IpfsServiceInterface interface {
	GetFile(cid string) ([]byte, error)
}

type CryptoServiceInterface interface {
	Decrypt(data []byte, password string) ([]byte, error)
}

// -------- CONSTRUCTOR --------
func NewGetIPFSDataQuerryHandler(ipfsService IpfsServiceInterface, cryptoService CryptoServiceInterface, tracecoreClient tracecore.TracecoreClient) *GetIPFSDataQuerryHandler {
	return &GetIPFSDataQuerryHandler{
		ipfsService: ipfsService,
		cryptoService: cryptoService,
		tracecoreClient: tracecoreClient,
	}
}

func (h *GetIPFSDataQuerryHandler) Execute(cmd GetIPFSDataQuerry) (*GetIPFSDataResponse, error) {
	
	// ------------------------------------------------------------
	// 1. LOAD ENCRYPTED VAULT CONTENT FROM IPFS
	// ------------------------------------------------------------
	// vaultData, err := h.ipfsService.GetFile(cmd.CID)
	encryptedData, err := h.GetFromIpfs(context.Background(), cmd)
	if err != nil {
		return nil, fmt.Errorf("❌ GetIPFSDataQuerryHandler - failed to fetch vault from IPFS: %w", err)
	}
	if encryptedData == nil || len(encryptedData) == 0 {
		return nil, fmt.Errorf("❌ GetIPFSDataQuerryHandler - empty vault data for CID %s", cmd.CID)
	}
	utils.LogPretty("GetIPFSDataQuerryHandler - Handle - encryptedData", encryptedData)

	// ------------------------------------------------------------
	// 2. DECRYPT VAULT
	// ------------------------------------------------------------
	decrypted, err := h.cryptoService.Decrypt(encryptedData, cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("❌ GetIPFSDataQuerryHandler - failed to decrypt vault: %w", err)
	}
	if len(decrypted) == 0 {
		return nil, fmt.Errorf("❌ GetIPFSDataQuerryHandler - vault decryption returned empty result")
	}
	utils.LogPretty("GetIPFSDataQuerryHandler - Handle - decrypted", decrypted)

	// ------------------------------------------------------------
	// 3. PARSE VAULT
	// ------------------------------------------------------------
	vaultPayload := vaults_domain.ParseVaultPayload(decrypted)
	utils.LogPretty("GetIPFSDataQuerryHandler - Handle - vaultPayload", vaultPayload)

	return &GetIPFSDataResponse{Data: vaultPayload}, nil
}


func (h *GetIPFSDataQuerryHandler) GetFromIpfs(ctx context.Context, req GetIPFSDataQuerry) ([]byte, error) {
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
	// 3. GET FROM IPFS
	// ------------------------------------------------------------
	response, err := storageProvider.Get(ctx, req.CID)
	if err != nil {
		return nil, fmt.Errorf("GetIPFSDataQuerryHandler - failed to get ipfs: %w", err)
	}

	return response, nil
}
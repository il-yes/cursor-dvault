package vault_queries

import (
	"context"
	"encoding/base64"
	"encoding/hex"
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
	CID       string
	Password  string
	AppCfg    app_config_domain.AppConfig
	UserID    string
	VaultName string
}

// -------- RESPONSE --------
type GetIPFSDataResponse struct {
	Data vault_domain.VaultPayload
}

// -------- INTERFACES --------
type CryptoServiceInterface interface {
	Decrypt(data []byte, password string) ([]byte, error)
}

// -------- HANDLER --------
type GetIPFSDataQuerryHandler struct {
	cryptoService CryptoServiceInterface
}

// -------- CONSTRUCTOR --------
func NewGetIPFSDataQuerryHandler(cryptoService CryptoServiceInterface) *GetIPFSDataQuerryHandler {
	return &GetIPFSDataQuerryHandler{
		cryptoService: cryptoService,
	}
}

func (h *GetIPFSDataQuerryHandler) Execute(cmd GetIPFSDataQuerry) (*GetIPFSDataResponse, error) {
    // 1. Get raw encrypted bytes from IPFS
    encryptedData, err := h.GetFromIpfs(context.Background(), cmd)
    if err != nil {
        return nil, fmt.Errorf("❌ failed to fetch vault from IPFS: %w", err)
    }
    utils.LogPretty("Raw encrypted bytes", len(encryptedData))
	utils.LogPretty("Raw encrypted bytes", hex.EncodeToString(encryptedData[:64]))

    // 2. SINGLE base64 decode (IPFS → binary)
    decryptedBytes, err := base64.StdEncoding.DecodeString(string(encryptedData))
    if err != nil {
        return nil, fmt.Errorf("❌ base64 decode failed: %w", err)
    }
    utils.LogPretty("✅ Base64 decoded", len(decryptedBytes))

    // 3. DECRYPT binary data
    plain, err := h.cryptoService.Decrypt(decryptedBytes, cmd.Password)
    if err != nil {
        return nil, fmt.Errorf("❌ decryption failed: %w", err)
    }

    // 4. Parse vault
    vaultPayload := vaults_domain.ParseVaultPayload(plain)
    return &GetIPFSDataResponse{Data: vaultPayload}, nil
}


func (h *GetIPFSDataQuerryHandler) GetFromIpfs(ctx context.Context, req GetIPFSDataQuerry) ([]byte, error) {
	// ------------------------------------------------------------
	// 1. LOAD TRACECORE CLIENT
	// ------------------------------------------------------------
	tracecoreClient := tracecore.NewTracecoreFromConfig(&req.AppCfg, "token")
	utils.LogPretty("GetIPFSDataQuerryHandler - GetFromIpfs - tracecoreClient", tracecoreClient.BaseURL)
	// ------------------------------------------------------------
	// 2. LOAD STORAGE PROVIDER
	// ------------------------------------------------------------
	storageProvider := blockchain.NewStorageProvider(blockchain.Config{
		StorageConfig: req.AppCfg.Storage,
		UserID:        req.UserID,
		VaultName:     req.VaultName,
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

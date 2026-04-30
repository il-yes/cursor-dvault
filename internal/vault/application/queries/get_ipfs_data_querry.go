package vault_queries

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	blockchain_ipfs "vault-app/internal/blockchain/ipfs"
	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/utils"
	vault_dto "vault-app/internal/vault/application/dto"
	vault_domain "vault-app/internal/vault/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
)

var (
	PRIVATE_MODE = "private"
	PUBLIC_MODE  = "public"
)

// -------- INTERFACES --------
type CryptoServiceInterface interface {
	Decrypt(data []byte, password string) ([]byte, error)
}
type IpfsServiceInterface interface {
	Get(ctx context.Context, cid string) ([]byte, error)
}

// -------- QUERRY --------
type GetIPFSDataQuerry struct {
	CID              string
	Password         string
	AppCfg           app_config_domain.AppConfig
	UserID           string
	VaultName        string
	UserOnboardingID string
	PrivateKey       string
	EncryptedKey     string
	SymKey           []byte
}

// -------- RESPONSE --------
type GetIPFSDataResponse struct {
	Raw  []byte
	Data vault_domain.VaultPayload
	Node vaults_domain.VaultNode
}

// -------- HANDLER --------
type GetIPFSDataQuerryHandler struct {
	cryptoService      CryptoServiceInterface
	IpfsService        IpfsServiceInterface
	UnlockVaultHandler vault_dto.UnlockVaultCommandInterface
	CryptoService      vaults_domain.VaultCrypto
	StorageFactory     blockchain_ipfs.StorageFactory
	EncryptionMode     string
}

// -------- CONSTRUCTOR --------
func NewGetIPFSDataQuerryHandler(
	cryptoService CryptoServiceInterface,
	vc *vault_infrastructure_crypto.AESService,
	sf *blockchain_ipfs.DefaultStorageFactory,
	unlockHandler vault_dto.UnlockVaultCommandInterface,
) *GetIPFSDataQuerryHandler {

	return &GetIPFSDataQuerryHandler{
		cryptoService:      cryptoService,
		CryptoService:      vc,
		StorageFactory:     sf,
		UnlockVaultHandler: unlockHandler,
		EncryptionMode:     PRIVATE_MODE,
	}
}

func (h *GetIPFSDataQuerryHandler) Execute(ctx context.Context, cmd GetIPFSDataQuerry) (*GetIPFSDataResponse, error) {
	if cmd.CID == "" {
		return nil, fmt.Errorf("CID is empty (invalid DAG state)")
	}
	// 1. Fetch raw data from IPFS
	// ==============================================
	rawBytes, err := h.GetFromIpfs(ctx, cmd)
	if err != nil {
		utils.LogPretty("GetIPFSDataQuerryHandler - Execute - rawBytes", err)
		return nil, err
	}
	if rawBytes == nil {
		utils.LogPretty("GetIPFSDataQuerryHandler - Execute - rawBytes is nil", err)
	}

	// Immediately after h.GetFromIpfs
	// utils.LogPretty("GetIPFSDataQuerryHandler - Execute - rawBytes len", len(rawBytes))
	// utils.LogPretty("GetIPFSDataQuerryHandler - Execute - rawBytes hex", hex.EncodeToString(rawBytes[:min(20, len(rawBytes))]))

	if h.UnlockVaultHandler == nil {
		utils.LogPretty("GetIPFSDataQuerryHandler - Execute - h.UnlockVaultHandler is nil", err)
	}

	// 2. Decrypt
	// ==============================================
	plain, err := h.HandleDecryption(cmd, rawBytes)
	if err != nil {
		return nil, fmt.Errorf("GetIPFSDataQuerryHandler - Execute - decrypt failed: %w", err)
	}

	// try to parse as VaultNode (optional)
	// ==============================================
	var node vaults_domain.VaultNode
	if err := json.Unmarshal(plain, &node); err == nil {
		return &GetIPFSDataResponse{
			Raw:  plain,
			Node: node,
		}, nil
	}

	// fallback → just raw
	return &GetIPFSDataResponse{
		Raw: plain,
	}, nil
}

func (h *GetIPFSDataQuerryHandler) GetFromIpfs(ctx context.Context, req GetIPFSDataQuerry) ([]byte, error) {
	utils.LogPretty("GetIPFSDataQuerryHandler - GetFromIpfs - req", req)
	vc := app_config_domain.VaultContext{
		AppConfig:     req.AppCfg,
		StorageConfig: req.AppCfg.Storage,
		UserID:        req.UserID,
		VaultName:     req.VaultName,
	}

	if h.StorageFactory == nil {
		utils.LogPretty("GetIPFSDataQuerryHandler - GetFromIpfs - fail h.StorageFactory is nil", h.StorageFactory)
	}

	storageProvider := h.StorageFactory.New(&vc)
	if storageProvider == nil {
		utils.LogPretty("GetIPFSDataQuerryHandler - GetFromIpfs - fail h.StorageFactory is nil", h.StorageFactory)
	}

	data, err := storageProvider.Get(ctx, req.CID)
	if err != nil {
		utils.LogPretty("GetIPFSDataQuerryHandler - GetFromIpfs - Get failed", err)
		return nil, fmt.Errorf("GetIPFSDataQuerryHandler - GetFromIpfs: %w", err)
	}

	if data == nil {
		return nil, fmt.Errorf("GetIPFSDataQuerryHandler - GetFromIpfs: Get returned nil data")
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("GetIPFSDataQuerryHandler - GetFromIpfs: Get returned empty data")
	}

	return data, nil
}

func (h GetIPFSDataQuerryHandler) PrivateDecryption(cmd GetIPFSDataQuerry, rawBytes []byte) ([]byte, error) {
	utils.LogPretty("GetIPFSDataQuerryHandler - ShareDecryption - ", "PrivateDecryption path")
	// 1. Unlock vault key
	// ==============================================
	unlockRes, err := h.UnlockVaultHandler.Execute(vault_dto.UnlockVaultCommand{
		Password: cmd.Password,
		UserID:   cmd.UserOnboardingID,
	})
	if err != nil {
		utils.LogPretty("GetIPFSDataQuerryHandler - Execute - unlockRes", err)
		return nil, fmt.Errorf("unlock failed: %w", err)
	}

	if h.CryptoService == nil {
		utils.LogPretty("GetIPFSDataQuerryHandler - Execute - fail CryptoService is nil", err)
	}

	// 2. Decrypt
	// ==============================================
	plain, err := h.CryptoService.Decrypt(rawBytes, unlockRes.VaultKey.Key)
	if err != nil {
		return nil, fmt.Errorf("decrypt failed: %w", err)
	}

	return plain, nil
}

func (h GetIPFSDataQuerryHandler) ShareDecryption(cmd GetIPFSDataQuerry, rawBytes []byte) ([]byte, error) {
	// 1. rawBytes is already base64‑decoded from CloudIPFSStorage → it's a *string* in bytes
    //    like "data:application/octet-stream;base64,..."
    //    Convert it to string to inspect and strip:
    s := string(rawBytes)

    // 2. If it's "data:...base64," prefixed, base64‑decode the payload:
    const prefix = "data:application/octet-stream;base64,"
    if strings.HasPrefix(s, prefix) {
        payload := s[len(prefix):]
        plain, err := base64.StdEncoding.DecodeString(payload)
        if err != nil {
            return nil, fmt.Errorf("failed to decode base64 payload: %w", err)
        }
        rawBytes = plain
    } else {
        // Already plain binary (no data: wrapper)
    }

	return rawBytes, nil
}
func (h GetIPFSDataQuerryHandler) HandleDecryption(cmd GetIPFSDataQuerry, rawBytes []byte) ([]byte, error) {
	if h.EncryptionMode == PUBLIC_MODE {
		utils.LogPretty("GetIPFSDataQuerryHandler - HandleDecryption - EncryptionMode", PUBLIC_MODE)
		return h.ShareDecryption(cmd, rawBytes)
	}
	utils.LogPretty("GetIPFSDataQuerryHandler - HandleDecryption - EncryptionMode", PRIVATE_MODE)
	return h.PrivateDecryption(cmd, rawBytes)
}

func (h *GetIPFSDataQuerryHandler) HydrateVaultNode(
	plainData []byte,
) (vaults_domain.VaultNode, error) {

	var node vaults_domain.VaultNode

	// 🔥 IMPORTANT: ensure it's JSON, not stringified JSON
	if err := json.Unmarshal(plainData, &node); err != nil {
		return vaults_domain.VaultNode{}, fmt.Errorf("hydrate vault node failed: %w", err)
	}

	return node, nil
}
func (q GetIPFSDataQuerry) WithCID(cid string) GetIPFSDataQuerry {
	q.CID = cid
	return q
}

func (h *GetIPFSDataQuerryHandler) SetIpfsService(ipfs IpfsServiceInterface) {
	h.IpfsService = ipfs
}

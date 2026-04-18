package vault_queries

import (
	"context"
	"encoding/json"
	"fmt"
	blockchain_ipfs "vault-app/internal/blockchain/ipfs"
	app_config_domain "vault-app/internal/config/domain"
	"vault-app/internal/utils"
	vault_dto "vault-app/internal/vault/application/dto"
	vault_domain "vault-app/internal/vault/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_crypto "vault-app/internal/vault/infrastructure/crypto"
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
	}
}

func (h *GetIPFSDataQuerryHandler) Execute(ctx context.Context, cmd GetIPFSDataQuerry) (*GetIPFSDataResponse, error) {
	if cmd.CID == "" {
		return nil, fmt.Errorf("CID is empty (invalid DAG state)")
	}
	// 1. Fetch raw data from IPFS
	rawBytes, err := h.GetFromIpfs(ctx, cmd)
	if err != nil {
		utils.LogPretty("GetIPFSDataQuerryHandler - Execute - rawBytes", err)
		return nil, err
	}

	if h.UnlockVaultHandler == nil {
		utils.LogPretty("GetIPFSDataQuerryHandler - Execute - h.UnlockVaultHandler", h.UnlockVaultHandler)
	}

	// 2. Unlock vault key
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

	// 3. Decrypt
	plain, err := h.CryptoService.Decrypt(rawBytes, unlockRes.VaultKey.Key)
	if err != nil {
		return nil, fmt.Errorf("decrypt failed: %w", err)
	}

	// try to parse as VaultNode (optional)
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
	vc := app_config_domain.VaultContext{
		AppConfig:     req.AppCfg,
		StorageConfig: req.AppCfg.Storage,
		UserID:        req.UserID,
		VaultName:     req.VaultName,
	}

	utils.LogPretty("GetIPFSDataQuerryHandler - GetFromIpfs - req", req)

	if h.StorageFactory == nil {
		utils.LogPretty("GetIPFSDataQuerryHandler - GetFromIpfs - fail h.StorageFactory is nil", h.StorageFactory)
	}

	storageProvider := h.StorageFactory.New(&vc)

	utils.LogPretty("GetIPFSDataQuerryHandler - GetFromIpfs - storageProvider", storageProvider)

	return storageProvider.Get(ctx, req.CID)
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

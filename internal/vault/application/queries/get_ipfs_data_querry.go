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
	VaultKey         []byte // 32-byte in-memory Vault DEK
	Configs          app_config_domain.Config
	UserID           string
	VaultName        string
	UserOnboardingID string
	PrivateKey       string
	EncryptedKey     string
	SymKey           []byte
}

// -------- RESPONSE --------
type GetIPFSDataResponse struct {
	Raw               []byte
	Data              vault_domain.VaultPayload
	Node              vaults_domain.VaultNode
	NodeBeta          vaults_domain.VaultNodeBeta
	CollaborativeNode vaults_domain.CollaborativeNode
	PersonalNode      vaults_domain.PersonalNode
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
	fmt.Printf("🔥 ENTERED GetIPFSDataQuerryHandler %p\n", h)
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
		return nil, fmt.Errorf("rawBytes nil after IPFS get")
	}

	if len(cmd.VaultKey) != 32 && cmd.Password == "" && h.UnlockVaultHandler == nil {
		utils.LogPretty("GetIPFSDataQuerryHandler - Execute - missing unlock credential and UnlockVaultHandler is nil", nil)
		return nil, fmt.Errorf("UnlockVaultHandler missing and no session VaultKey provided")
	}
	// utils.LogPretty("GetIPFSDataQuerryHandler - Execute - cmd", cmd)
	// utils.LogPretty("GetIPFSDataQuerryHandler - Execute - rawBytes", rawBytes)

	// 2. Decrypt
	// ==============================================
	plain, err := h.HandleDecryption(cmd, rawBytes)
	if err != nil {
		return nil, fmt.Errorf("GetIPFSDataQuerryHandler - Execute - decrypt failed: %w", err)
	}
	utils.LogPretty(
		"RAW ROOT BEFORE PARSE",
		string(plain),
	)
	fmt.Println("🔥🔥🔥 DECRYPTED ROOT LENGTH:", len(plain))
	fmt.Println("🔥🔥🔥 DECRYPTED ROOT:", string(plain))

	utils.LogPretty("GetIPFSDataQuerryHandler - Execute - Vault Merkle Tree", plain)
	utils.LogPretty(
		"RAW ROOT BEFORE PARSE",
		string(plain),
	)


	var node vaults_domain.VaultNode
	_ = json.Unmarshal(plain, &node)
	if node.Version == "" {
		var meta struct {
			Version string `json:"version"`
		}
		if err := json.Unmarshal(plain, &meta); err == nil && meta.Version != "" {
			node.Version = meta.Version
		}
	}

	// try to parse as VaultNodeBeta (optional)
	// ==============================================
	var nodeBeta vaults_domain.VaultNodeBeta
	if err := json.Unmarshal(plain, &nodeBeta); err == nil {
		if nodeBeta.Type == "vault" && 
			nodeBeta.Personal.CID != "" &&
			nodeBeta.Collaborative.CID != "" {

			return &GetIPFSDataResponse{
				Raw:      plain,
				NodeBeta: nodeBeta,
				Node:     node,
			}, nil
		}
	}
	utils.LogPretty(
		"NODE BETA DEBUG",
		nodeBeta,
	)


	// try to parse as CollaborativeNodeBet (optional)
	// ==============================================
	var collaborativeNode vaults_domain.CollaborativeNode
	if err := json.Unmarshal(plain, &collaborativeNode); err == nil {
		if collaborativeNode.Type == "collaborative_vault" {
			return &GetIPFSDataResponse{
				CollaborativeNode: collaborativeNode,
				Raw:               plain,
			}, nil

		}
	}
	
	// try to parse as VaultNode (optional) == PersonalNode
	// ==============================================
	if node.Type != "" || node.Version != "" {

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

	if h.IpfsService != nil {
		bytes, err := h.IpfsService.Get(ctx, req.CID)
		if err == nil && len(bytes) > 0 {
			return bytes, nil
		}
	}

	vc := app_config_domain.VaultContext{
		Configs:   req.Configs,
		UserID:    req.UserID,
		VaultName: req.VaultName,
	}
	if req.Configs.App != nil {
		vc.StorageConfig = req.Configs.App.Storage
	}
	if req.Configs.Subscription != nil {
		vc.UserSubscriptionID = req.Configs.Subscription.UserID
	}

	if h.StorageFactory == nil {
		return nil, fmt.Errorf("StorageFactory is nil")
	}

	storageProvider := h.StorageFactory.New(&vc)
	if storageProvider == nil {
		return nil, fmt.Errorf("storageProvider is nil")
	}

	data, err := storageProvider.Get(ctx, req.CID)
	if err != nil {
		utils.LogPretty("GetIPFSDataQuerryHandler - GetFromIpfs - Get failed", err)
		return nil, fmt.Errorf("GetIPFSDataQuerryHandler - GetFromIpfs: %w", err)
	}
	utils.LogPretty("GetIPFSDataQuerryHandler - GetFromIpfs - data", data)

	if data == nil {
		return nil, fmt.Errorf("GetIPFSDataQuerryHandler - GetFromIpfs: Get returned nil data")
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("GetIPFSDataQuerryHandler - GetFromIpfs: Get returned empty data")
	}

	return data, nil
}

func (h GetIPFSDataQuerryHandler) PrivateDecryption(cmd GetIPFSDataQuerry, rawBytes []byte) ([]byte, error) {
	utils.LogPretty("GetIPFSDataQuerryHandler - PrivateDecryption - ", "PrivateDecryption path")
	var vaultKey []byte

	if len(cmd.VaultKey) == 32 {
		vaultKey = cmd.VaultKey
	} else if cmd.Password != "" {
		unlockRes, err := h.UnlockVaultHandler.Execute(vault_dto.UnlockVaultCommand{
			Password: cmd.Password,
			UserID:   cmd.UserOnboardingID,
		})
		if err != nil {
			utils.LogPretty("GetIPFSDataQuerryHandler - Execute - unlockRes", err)
			return nil, fmt.Errorf("unlock failed: %w", err)
		}
		vaultKey = unlockRes.VaultKey.Key
	} else {
		return nil, fmt.Errorf("PrivateDecryption: missing vault DEK and no unlock password provided")
	}

	if len(vaultKey) != 32 {
		return nil, fmt.Errorf("PrivateDecryption: invalid vault DEK length (expected 32, got %d)", len(vaultKey))
	}

	plain, err := h.CryptoService.Decrypt(rawBytes, vaultKey)
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
	utils.LogPretty("GetIPFSDataQuerryHandler - HandleDecryption - EncryptionMode", h.EncryptionMode)
	if h.EncryptionMode == PUBLIC_MODE {
		utils.LogPretty("GetIPFSDataQuerryHandler - HandleDecryption - EncryptionMode", PUBLIC_MODE)
		return h.ShareDecryption(cmd, rawBytes)
	}
	utils.LogPretty("GetIPFSDataQuerryHandler - HandleDecryption - EncryptionMode", PRIVATE_MODE)
	return h.PrivateDecryption(cmd, rawBytes)
}

func (h *GetIPFSDataQuerryHandler) HydrateVaultNode(plainData []byte) (vaults_domain.VaultNode, error) {

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

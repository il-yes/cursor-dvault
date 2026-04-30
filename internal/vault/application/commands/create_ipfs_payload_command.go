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

var (
	PRIVATE_MODE = "private"
	PUBLIC_MODE  = "public"
)

type CreateIPFSPayloadCommand struct {
	Vault            *vaults_domain.Vault
	Password         string
	Data             []byte
	UserID           string // User app
	ShareKey         []byte
	UserOnboardingID string
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
	StorageFactory     blockchain_ipfs.StorageFactory
	EncryptionMode     string
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
		VaultRepo:          vaultRepo,
		CryptoService:      vc,
		TracecoreClient:    tracecoreClient,
		StorageFactory:     sf,
		UnlockVaultHandler: uh,
		EncryptionMode:     PRIVATE_MODE,
	}
}
func (h *CreateIPFSPayloadCommandHandler) Execute(
	ctx context.Context,
	vaultCtx app_config_domain.VaultContext,
	cmd CreateIPFSPayloadCommand,
) (*CreateIPFSPayloadCommandResult, error) {
	// 2. Encryption
	// ==============================================
	encrypted, err := h.HandleEcryption(cmd, vaultCtx)
	if err != nil {
		return nil, fmt.Errorf("CreateIPFSPayloadCommandHandler - Execute - vault encryption failed: %w", err)
	}

	// IPFS Upload
	// ==============================================
	cidFromIpfs, err := h.StoreOnIpfs(ctx, vaultCtx, encrypted)
	if err != nil {
		return nil, fmt.Errorf("CreateIPFSPayloadCommandHandler - Execute - failed to add vault to IPFS: %w", err)
	}
	utils.LogPretty("CreateIPFSPayloadCommandHandler - Execute - cidFromIpfs", cidFromIpfs)

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

func (h *CreateIPFSPayloadCommandHandler) PrivateEncryption(cmd CreateIPFSPayloadCommand, vaultCtx app_config_domain.VaultContext) ([]byte, error) {
	// 1. Unlock vault key
	// ==============================================
	unlockRes, err := h.UnlockVaultHandler.Execute(vault_dto.UnlockVaultCommand{
		Password: cmd.Password,
		UserID:   cmd.UserOnboardingID, // userOnboarding required
	})
	if err != nil {
		utils.LogPretty("CreateIPFSPayloadCommandHandler - PrivateEncryption - - error", cmd)
		return nil, fmt.Errorf("CreateIPFSPayloadCommandHandler - PrivateEncryption - failed to unlock vault key: %w", err)
	}
	vaultKey := unlockRes.VaultKey.Key

	encrypted, err := h.CryptoService.Encrypt(cmd.Data, vaultKey)
	if err != nil {
		return nil, fmt.Errorf("CreateIPFSPayloadCommandHandler - PrivateEncryption - vault encryption failed: %w", err)
	}
	return encrypted, nil
}

func (h *CreateIPFSPayloadCommandHandler) ShareEncryption(cmd CreateIPFSPayloadCommand, vaultCtx app_config_domain.VaultContext) ([]byte, error) {
	return cmd.Data, nil
}

func (h *CreateIPFSPayloadCommandHandler) HandleEcryption(cmd CreateIPFSPayloadCommand, vaultCtx app_config_domain.VaultContext) ([]byte, error) {
	if h.EncryptionMode == PUBLIC_MODE {
		utils.LogPretty("CreateIPFSPayloadCommandHandler - HandleEcryption - EncryptionMode", PUBLIC_MODE)
		return h.ShareEncryption(cmd, vaultCtx)
	}
	utils.LogPretty("CreateIPFSPayloadCommandHandler - HandleEcryption - EncryptionMode", PRIVATE_MODE)
	return h.PrivateEncryption(cmd, vaultCtx)
}


func (h *CreateIPFSPayloadCommandHandler) SetIpfsService(ipfs IpfsServiceInterface) {
	h.IpfsService = ipfs
}
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
	VaultKey         []byte // 32-byte in-memory Vault DEK from active session
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
	TracecoreClient    *tracecore.TracecoreClient
	UnlockVaultHandler UnlockVaultHandlerInterface
	StorageFactory     blockchain_ipfs.StorageFactory
	EncryptionMode     string
}

// -------- constructor --------
func NewCreateIPFSPayloadCommandHandler(
	vaultRepo vaults_domain.VaultRepository,
	tracecoreClient *tracecore.TracecoreClient,
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
	// utils.LogPretty("CreateIPFSPayloadCommandHandler - Execute - cidFromIpfs", cidFromIpfs)

	return &CreateIPFSPayloadCommandResult{CID: cidFromIpfs}, nil
}

func (h *CreateIPFSPayloadCommandHandler) StoreOnIpfs(
	ctx context.Context,
	vaultCtx app_config_domain.VaultContext,
	data []byte,
) (string, error) {
	if h.IpfsService != nil {
		return h.IpfsService.Add(ctx, data)
	}
	storageProvider := h.StorageFactory.New(&vaultCtx)
	return storageProvider.Add(ctx, data)
}

func (h *CreateIPFSPayloadCommandHandler) PrivateEncryption(cmd CreateIPFSPayloadCommand, vaultCtx app_config_domain.VaultContext) ([]byte, error) {
	var vaultKey []byte

	// 🟢 FAST-PATH: Use unlocked Vault DEK directly from in-memory session if available
	if len(cmd.VaultKey) == 32 {
		vaultKey = cmd.VaultKey
		utils.LogPretty("CreateIPFSPayloadCommandHandler - PrivateEncryption - using session VaultKey DEK", map[string]interface{}{
			"authenticatedUserID": vaultCtx.UserID,
			"keyLength":           len(vaultKey),
		})
	} else if cmd.Password != "" {
		// Fallback for explicit unlock boundary (e.g. initial unlock with password)
		utils.LogPretty("CreateIPFSPayloadCommandHandler - PrivateEncryption - unlocking via keyring password boundary", map[string]interface{}{
			"authenticatedUserID": vaultCtx.UserID,
			"userOnboardingID":    cmd.UserOnboardingID,
		})
		unlockRes, err := h.UnlockVaultHandler.Execute(vault_dto.UnlockVaultCommand{
			Password: cmd.Password,
			UserID:   cmd.UserOnboardingID,
		})
		if err != nil {
			return nil, fmt.Errorf("CreateIPFSPayloadCommandHandler - PrivateEncryption - failed to unlock vault keyring: %w", err)
		}
		vaultKey = unlockRes.VaultKey.Key
	} else {
		return nil, fmt.Errorf("CreateIPFSPayloadCommandHandler - PrivateEncryption: missing vault DEK and no unlock password provided")
	}

	if len(vaultKey) != 32 {
		return nil, fmt.Errorf("CreateIPFSPayloadCommandHandler - PrivateEncryption: invalid vault DEK length (expected 32, got %d)", len(vaultKey))
	}

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
	// utils.LogPretty("CreateIPFSPayloadCommandHandler - HandleEcryption - EncryptionMode", PRIVATE_MODE)
	return h.PrivateEncryption(cmd, vaultCtx)
}

func (h *CreateIPFSPayloadCommandHandler) SetIpfsService(ipfs IpfsServiceInterface) {
	h.IpfsService = ipfs
}

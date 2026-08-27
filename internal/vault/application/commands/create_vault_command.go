package vault_commands

import (
	"context"
	"errors"
	"fmt"

	app_config "vault-app/internal/config"
	app_config_domain "vault-app/internal/config/domain"
	onboarding_domain "vault-app/internal/onboarding/domain"
	utils "vault-app/internal/utils"
	vault_session "vault-app/internal/vault/application/session"
	vault_domain "vault-app/internal/vault/domain"
	vaults_storage_engine_type "vault-app/internal/vault/infrastructure/storage/engine/types"
)

func buildVaultContext(cmd CreateVaultCommand) app_config_domain.VaultContext {
	var storageCfg app_config.StorageConfig
	if cmd.Configs.App != nil {
		storageCfg = cmd.Configs.App.Storage
	}
	vc := app_config_domain.VaultContext{
		Configs:            cmd.Configs,
		StorageConfig:      storageCfg,
		UserID:             cmd.UserID,
		VaultName:          cmd.VaultName,
		UserSubscriptionID: cmd.UserSubscriptionID,
	}
	if cmd.UserOnboarding != nil {
		vc.UserOnboarding = cmd.UserOnboarding.ID
	}
	return vc
}

// -------- COMMAND query --------
type CreateVaultCommand struct {
	UserID             string
	VaultName          string
	Password           string
	UserSubscriptionID string
	Configs            app_config_domain.Config
	UserOnboarding     *onboarding_domain.User
}

// -------- RESULT --------
type CreateVaultResult struct {
	Vault          *vault_domain.Vault
	ReusedExisting bool
}

// -------- COMMAND handler interfaces --------
type CryptoServiceInterface interface {
	Encrypt(data []byte, password string) ([]byte, error)
}

type IpfsServiceInterface interface {
	Add(ctx context.Context, data []byte) (string, error)
}

type InitializeVaultHandler interface {
	Execute(cmd InitializeVaultCommand) (*InitializeVaultResult, error)
}

type CreateIPFSPayloadHandler interface {
	Execute(ctx context.Context, vc app_config_domain.VaultContext, cmd CreateIPFSPayloadCommand) (*CreateIPFSPayloadCommandResult, error)
	SetIpfsService(i IpfsServiceInterface)
}

type StorageEngineCommitter interface {
	Commit(session *vault_session.Session, mode vaults_storage_engine_type.SyncMode, opts ...vaults_storage_engine_type.CommitOptions) (string, []vaults_storage_engine_type.EntryUpdate, int, int, error)
}

// -------- COMMAND handler --------
type CreateVaultCommandHandler struct {
	initializeVaultHandler   InitializeVaultHandler
	createIPFSPayloadHandler CreateIPFSPayloadHandler
	vaultRepo                vault_domain.VaultRepository
	storageEngine            StorageEngineCommitter
}

type VaultInterfaceService interface {
	CreateVault(cmd CreateVaultCommand) (*CreateVaultResult, error)
}

// -------- COMMAND handler constructor --------
func NewCreateVaultCommandHandler(
	initializator InitializeVaultHandler,
	creator CreateIPFSPayloadHandler,
	vaultRepo vault_domain.VaultRepository,
) *CreateVaultCommandHandler {
	return &CreateVaultCommandHandler{
		initializeVaultHandler:   initializator,
		createIPFSPayloadHandler: creator,
		vaultRepo:                vaultRepo,
	}
}

func (h *CreateVaultCommandHandler) SetStorageEngine(se StorageEngineCommitter) {
	h.storageEngine = se
}

func (h *CreateVaultCommandHandler) CreateVault(cmd CreateVaultCommand) (*CreateVaultResult, error) {
	utils.LogPretty("CreateVaultCommandHandler - CreateVault - cmd", cmd)

	if cmd.VaultName == "" {
		cmd.VaultName = cmd.UserID + "-vault"
	}

	// 1. Idempotency check: reuse existing vault if already created
	if h.vaultRepo != nil {
		existing, _ := h.vaultRepo.GetLatestByUserID(cmd.UserID)
		if existing != nil && existing.CID != "" {
			utils.LogPretty("CreateVaultCommandHandler - Reusing existing vault", existing)
			return &CreateVaultResult{
				Vault:          existing,
				ReusedExisting: true,
			}, nil
		}
	}

	// 1.b Check optional InitVaultHandler
	if h.initializeVaultHandler != nil {
		_, err := h.initializeVaultHandler.Execute(InitializeVaultCommand{
			UserID:    cmd.UserID,
			VaultName: cmd.VaultName,
		})
		if err != nil {
			return nil, err
		}
	}

	// 2. Build initial domain Vault entity in memory (do NOT save to DB yet)
	vault := vault_domain.NewVault(cmd.UserID, cmd.VaultName)
	vault.AttachUserSubscriptionID(cmd.UserSubscriptionID)
	utils.LogPretty("CreateVaultCommandHandler - transient vault initialized", vault)

	// 3. Vault - Build initial payload
	const InitialVaultVersion = "1.0.0"
	vaultPayload := vault.BuildInitialPayload(InitialVaultVersion)
	utils.LogPretty("CreateVaultCommandHandler - Execute - vaultPayload", vaultPayload)

	// 4. Commit initial DAG via StorageEngine (if available) or fallback to flat IPFS payload
	if h.storageEngine != nil {
		session := vault_session.InitNewSession(cmd.UserID)
		session.Vault = vaultPayload.ToBytes()
		session.Runtime = vault_session.NewRuntimeContext()
		session.Runtime.SetVaultID(vault.ID)
		session.Runtime.SetVaultName(vault.Name)
		if cmd.UserOnboarding != nil {
			session.Runtime.UserConfig.ID = cmd.UserOnboarding.ID
		}

		vc := buildVaultContext(cmd)

		rootCID, _, _, _, err := h.storageEngine.Commit(session, vaults_storage_engine_type.FullSync, vaults_storage_engine_type.CommitOptions{
			VaultContext: vc,
			Password:     cmd.Password,
		})
		if err != nil {
			utils.LogPretty("CreateVaultCommandHandler - StorageEngine.Commit err", err)
			return nil, fmt.Errorf("failed to commit initial vault DAG: %w", err)
		}
		vault.AttachCID(rootCID)
	} else {
		// Fallback for tests/environments where StorageEngine is not injected
		vaultBytes, err := vaultPayload.GetContentBytes()
		if err != nil {
			return nil, fmt.Errorf("❌ vault encryption failed: %w", err)
		}
		vc := buildVaultContext(cmd)

		if h.createIPFSPayloadHandler == nil {
			return nil, errors.New("CreateVaultCommandHandler: createIPFSPayloadHandler is nil")
		}

		userOnboardingID := ""
		if cmd.UserOnboarding != nil {
			userOnboardingID = cmd.UserOnboarding.ID
		}

		ipfsRecord, err := h.createIPFSPayloadHandler.Execute(
			context.Background(),
			vc,
			CreateIPFSPayloadCommand{
				Vault:            vault,
				Password:         cmd.Password,
				Data:             vaultBytes,
				UserID:           cmd.UserID,
				UserOnboardingID: userOnboardingID,
			})
		if err != nil {
			return nil, err
		}
		if ipfsRecord == nil {
			return nil, errors.New("IPFS payload result is nil")
		}
		vault.AttachCID(ipfsRecord.CID)
	}

	if vault == nil || vault.CID == "" {
		return nil, errors.New("vault CID is empty before SaveVault")
	}

	// 5. Persist Vault metadata to DB ONLY AFTER DAG commit has succeeded with valid CID
	if err := h.vaultRepo.SaveVault(vault); err != nil {
		utils.LogPretty("CreateVaultCommandHandler - SaveVault err", err)
		return nil, fmt.Errorf("failed to persist initial vault metadata: %w", err)
	}

	return &CreateVaultResult{
		Vault:          vault,
		ReusedExisting: false,
	}, nil
}
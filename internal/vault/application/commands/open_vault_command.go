package vault_commands

import (
	"context"
	"errors"
	"time"

	utils "vault-app/internal"
	vault_events "vault-app/internal/vault/application/events"
	vault_session "vault-app/internal/vault/application/session"
	vault_domain "vault-app/internal/vault/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vaults_persistence "vault-app/internal/vault/infrastructure/persistence"

	"gorm.io/gorm"
)

// -------- COMMAND --------

type OpenVaultCommand struct {
	UserID   string
	Password string
	Session  *vault_session.Session
}

// -------- RESULT --------

type OpenVaultResult struct {
	Vault          *vault_domain.Vault
	Content        *vault_domain.VaultPayload
	RuntimeContext *vault_session.RuntimeContext
	Session        *vault_session.Session
	LastCID        string
	ReusedExisting bool
}

// -------- HANDLER --------

type OpenVaultCommandHandler struct {
	vaultRepo vault_domain.VaultRepository
	now       func() string
}

// -------- CONSTRUCTOR --------

func NewOpenVaultCommandHandler(
	db *gorm.DB,
) *OpenVaultCommandHandler {
	vaultRepo := vaults_persistence.NewGormVaultRepository(db)

	return &OpenVaultCommandHandler{
		vaultRepo: vaultRepo,
		now:       func() string { return time.Now().UTC().Format(time.RFC3339) },
	}
}

// -------- EXECUTION --------
func (h *OpenVaultCommandHandler) Handle(
	ctx context.Context,
	cmd OpenVaultCommand,
	ipfs vault_domain.VaultStorage,
	crypto vault_domain.CryptoService,
	eventBus vault_events.VaultEventBus,
) (*OpenVaultResult, error) {

	// ------------------------------------------------------------
	// 0. ENFORCE INVARIANTS (NON-NEGOTIABLE)
	// ------------------------------------------------------------
	if cmd.Session == nil {
		cmd.Session = vault_session.InitNewSession(cmd.UserID)
	}

	runtimeCtx := vault_session.NewRuntimeContext()

	// ------------------------------------------------------------
	// 1. REUSE SESSION VAULT IF POSSIBLE (SINGLE PATH)
	// ------------------------------------------------------------
	if cmd.Session.Vault != nil && cmd.Session.LastCID != "" {
		utils.LogPretty("OpenVaultCommandHandler - Handle - REUSE SESSION VAULT IF POSSIBLE (SINGLE PATH)", cmd.Session)	
		payload := vaults_domain.ParseVaultPayload(cmd.Session.Vault)
		utils.LogPretty("OpenVaultCommandHandler - Handle - parsed payload", payload)

		evt := vault_events.VaultOpened{
			UserID:       cmd.UserID,
			VaultPayload: &payload,
			LastCID:      cmd.Session.LastCID,
			LastSynced:   cmd.Session.LastSynced,
			LastUpdated:  cmd.Session.LastUpdated,
			Runtime:      runtimeCtx,
			OccurredAt:   time.Now().Unix(),
		}

		eventBus.PublishVaultOpened(ctx, evt)

		return &OpenVaultResult{
			Vault:          nil,
			Content:        &payload,
			RuntimeContext: runtimeCtx,
			Session:        cmd.Session,
			LastCID:        cmd.Session.LastCID,
			ReusedExisting: true,
		}, nil
	}

	// ------------------------------------------------------------
	// 2. LOAD OR CREATE VAULT METADATA
	// ------------------------------------------------------------
	utils.LogPretty("OpenVaultCommandHandler - Handle - LOAD OR CREATE VAULT METADATA", cmd.UserID)
	vault, err := h.vaultRepo.GetLatestByUserID(cmd.UserID)
	if err != nil {
		if errors.Is(err, vault_domain.ErrVaultNotFound) {
			vault = vault_domain.NewVault(cmd.UserID, "New Vault")
			if err := h.vaultRepo.SaveVault(vault); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// ------------------------------------------------------------
	// 3. LOAD ENCRYPTED VAULT CONTENT FROM IPFS
	// ------------------------------------------------------------
	encrypted, err := ipfs.GetData(vault.CID)
	if err != nil {
		return nil, err
	}

	// ------------------------------------------------------------
	// 4. DECRYPT VAULT
	// ------------------------------------------------------------
	decrypted, err := crypto.Decrypt(encrypted, cmd.Password)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("Vault New blob JSON: %s\n", string(decrypted))

	// ------------------------------------------------------------
	// 5. PARSE VAULT PAYLOAD
	// ------------------------------------------------------------
	payload := vaults_domain.ParseVaultPayload(decrypted)
	payload.Normalize()

	// ------------------------------------------------------------
	// 6. UPDATE SESSION (IN-MEMORY)
	// ------------------------------------------------------------
	cmd.Session.Vault = payload.ToBytes()
	cmd.Session.LastCID = vault.CID

	// ------------------------------------------------------------
	// 7. EMIT EVENT
	// ------------------------------------------------------------
	evt := vault_events.VaultOpened{
		UserID:       cmd.UserID,
		VaultPayload: &payload,
		LastCID:      vault.CID,
		LastSynced:   vault.UpdatedAt,
		LastUpdated:  vault.UpdatedAt,
		Runtime:      runtimeCtx,
		OccurredAt:   time.Now().Unix(),
	}

	eventBus.PublishVaultOpened(ctx, evt)

	// ------------------------------------------------------------
	// 8. RETURN RESULT
	// ------------------------------------------------------------
	return &OpenVaultResult{
		Vault:          vault,
		Content:        &payload,
		RuntimeContext: runtimeCtx,
		Session:        cmd.Session,
		LastCID:        vault.CID,
		ReusedExisting: false,
	}, nil
}

package vault_commands

import (
	"context"
	"errors"
	"time"

	utils "vault-app/internal"
	vault_session "vault-app/internal/vault/application/session"
	vault_domain "vault-app/internal/vault/domain"
)

// -------- COMMAND --------

type OpenVaultCommand struct {
	UserID   string
	Password string
}

// -------- RESULT --------

type OpenVaultResult struct {
	Vault               *vault_domain.Vault
	Content        *vault_domain.VaultPayload
	RuntimeContext      *vault_session.RuntimeContext
	Session             *vault_session.Session
	LastCID             string
	ReusedExisting      bool
}

// -------- HANDLER --------

type OpenVaultCommandHandler struct {
	vaultRepo   vault_domain.VaultRepository
	sessionMgr  *vault_session.Manager
	ipfs        vault_domain.VaultStorage // abstraction over IPFS
	crypto      vault_domain.CryptoService
	now         func() string
}

// -------- CONSTRUCTOR --------

func NewOpenVaultCommandHandler(
	vaultRepo vault_domain.VaultRepository,
	sessionMgr *vault_session.Manager,
	ipfs vault_domain.VaultStorage,
	crypto vault_domain.CryptoService,
) *OpenVaultCommandHandler {
	return &OpenVaultCommandHandler{
		vaultRepo:  vaultRepo,
		sessionMgr: sessionMgr,
		ipfs:       ipfs,
		crypto:     crypto,
		now:        func() string { return time.Now().UTC().Format(time.RFC3339) },
	}
}

// -------- EXECUTION --------

func (h *OpenVaultCommandHandler) Handle(
	ctx context.Context,
	cmd OpenVaultCommand,
) (*OpenVaultResult, error) {
	utils.LogPretty("OpenVaultCommandHandler - cmd", cmd)

	// 1️⃣ Reuse existing session if present
	if existing, ok := h.sessionMgr.Get(cmd.UserID); ok && existing.Vault != nil {
 
		utils.LogPretty("OpenVaultCommandHandler - existing", existing)
		return &OpenVaultResult{
			Vault:          nil,
			Content:        existing.Vault,
			RuntimeContext: existing.Runtime,
			Session:        existing,
			LastCID:        existing.LastCID,
			ReusedExisting: true,
		}, nil
	}

	// 2️⃣ Load latest vault metadata
	vault, err := h.vaultRepo.GetLatestByUserID(cmd.UserID)
	if err != nil {
		if errors.Is(err, vault_domain.ErrVaultNotFound) {
			// 3️⃣ Create minimal vault
			vault = vault_domain.NewVault(cmd.UserID, "New Vault")
			if err := h.vaultRepo.SaveVault(vault); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// 4️⃣ Fetch encrypted vault payload
	encrypted, err := h.ipfs.GetData(vault.CID)
	if err != nil {
		return nil, err
	}

	// 5️⃣ Decrypt vault
	decrypted, err := h.crypto.Decrypt(encrypted, cmd.Password)
	if err != nil {
		return nil, err
	}

	// 6️⃣ Parse vault payload
	payload := vault_domain.ParseVaultPayload(decrypted)

	// 7️⃣ Create runtime context
	runtimeCtx := vault_session.NewRuntimeContext()

	// 8️⃣ Start session
	session := h.sessionMgr.StartSession(
		cmd.UserID,
		payload,
		vault.CID,
		runtimeCtx,	
	)
	utils.LogPretty("payload", payload)

	return &OpenVaultResult{
		Vault:          vault,
		Content:        &payload,
		RuntimeContext: runtimeCtx,
		Session:        session,
		LastCID:        vault.CID,
		ReusedExisting: false,
	}, nil
}

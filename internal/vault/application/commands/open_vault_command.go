package vault_commands

import (
	"context"
	"errors"
	"log"
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

	// 1️⃣ Session - Reuse existing session if present
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

	// 2️⃣ Vault - Load latest vault metadata
	vault, err := h.vaultRepo.GetLatestByUserID(cmd.UserID)
	if err != nil {
		log.Println("OpenVaultCommandHandler - Error loading vault, go for minimal vault", err)
		if errors.Is(err, vault_domain.ErrVaultNotFound) {
			log.Println("OpenVaultCommandHandler - Create minimal vault")
			// 3️⃣ Create minimal vault
			vault = vault_domain.NewVault(cmd.UserID, "New Vault")
			utils.LogPretty("OpenVaultCommandHandler - new minimal vault", vault)
			if err := h.vaultRepo.SaveVault(vault); err != nil {
				log.Println("OpenVaultCommandHandler - Error saving minimal vault", err)
				return nil, err
			}
		} else {
			log.Println("OpenVaultCommandHandler - Error loading vault no minimal vault", err)
			return nil, err
		}
	}
	utils.LogPretty("OpenVaultCommandHandler - vault", vault)

	// 4️⃣ IPFS - Fetch encrypted vault payload
	encrypted, err := h.ipfs.GetData(vault.CID)
	if err != nil {
		log.Println("OpenVaultCommandHandler - Error fetching encrypted vault payload", err)
		return nil, err
	}
	utils.LogPretty("OpenVaultCommandHandler - encrypted", encrypted)

	// 5️⃣ Crypto - Decrypt vault	
	decrypted, err := h.crypto.Decrypt(encrypted, cmd.Password)
	if err != nil {
		return nil, err
	}

	// 6️⃣ Vault - Parse vault payload	
	payload := vault_domain.ParseVaultPayload(decrypted)
	log.Println("OpenVaultCommandHandler - payload", payload)

	// 7️⃣ Session - Create runtime context
	runtimeCtx := vault_session.NewRuntimeContext()
	log.Println("OpenVaultCommandHandler - Initialize new runtimeCtx")

	// 8️⃣ Start session
	session := h.sessionMgr.StartSession(
		cmd.UserID,
		payload,
		vault.CID,
		runtimeCtx,	
	)
	utils.LogPretty("OpenVaultCommandHandler - Start new session", session)

	return &OpenVaultResult{
		Vault:          vault,
		Content:        &payload,
		RuntimeContext: runtimeCtx,
		Session:        session,
		LastCID:        vault.CID,
		ReusedExisting: false,
	}, nil
}

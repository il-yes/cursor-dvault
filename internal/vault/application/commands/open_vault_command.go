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
	Session  *vault_session.Session
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
	reusedExisting := true

	// 1. Vault - Use existing if present or Load latest vault metadata
	vault, err := h.vaultRepo.GetLatestByUserID(cmd.UserID)
	if err != nil {
		log.Println("OpenVaultCommandHandler - Error loading vault, go for minimal vault", err)
		if errors.Is(err, vault_domain.ErrVaultNotFound) {
			log.Println("OpenVaultCommandHandler - Create minimal vault")
			// 3️⃣ Create minimal vault - secure by default
			vault = vault_domain.NewVault(cmd.UserID, "New Vault")
			utils.LogPretty("OpenVaultCommandHandler - new minimal vault", vault)
			if err := h.vaultRepo.SaveVault(vault); err != nil {
				log.Println("OpenVaultCommandHandler - Error saving minimal vault", err)
				return nil, err
			}
			reusedExisting = false
		} else {
			log.Println("OpenVaultCommandHandler - Error loading vault no minimal vault", err)
			return nil, err
		}
	}
	utils.LogPretty("OpenVaultCommandHandler - vault", vault)

	// 2. Session - Reuse existing session if present
	if cmd.Session.LastCID != "" {
		openedResult := &OpenVaultResult{
			Vault:          vault,
			Content:        cmd.Session.Vault,
			RuntimeContext: cmd.Session.Runtime,
			Session:        cmd.Session,
			LastCID:        cmd.Session.LastCID,
			ReusedExisting: reusedExisting,
		}
		utils.LogPretty("OpenVaultCommandHandler - reused existing session", cmd)
		return openedResult, nil
	}

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

	// 6. Firre event Vault Opened
	// if h.eventBus != nil {
	// 	h.eventBus.Publish(vault_domain.NewVaultOpenedEvent(vault.ID, vault.UserID, vault.Name))
	// }

	return &OpenVaultResult{
		Vault:          vault,
		Content:        &payload,
		RuntimeContext: runtimeCtx,
		LastCID:        vault.CID,
		ReusedExisting: false,
	}, nil
}

package vault_dto

import (
	auth_domain "vault-app/internal/auth/domain"
	identity_domain "vault-app/internal/identity/domain"
	vault_session "vault-app/internal/vault/application/session"
	vault_domain "vault-app/internal/vault/domain"

)

type LoginResponse struct {
	User                identity_domain.User
	Tokens              auth_domain.TokenPairs
	SessionID           string
	Vault               vault_domain.VaultPayload
	VaultRuntimeContext vault_session.RuntimeContext
	LastCID             string
	Dirty               bool
}

type SynchronizeVaultRequest struct {
	UserID string `json:"user_id"`
	Password string `json:"password"`
	Vault vault_domain.Vault `json:"vault"`
}
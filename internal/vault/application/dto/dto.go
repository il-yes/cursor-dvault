package vault_dto

import (
	auth_domain "vault-app/internal/auth/domain"
	identity_domain "vault-app/internal/identity/domain"
	vault_session "vault-app/internal/vault/application/session"
	vault_domain "vault-app/internal/vault/domain"
	vaults_domain "vault-app/internal/vault/domain"
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
type SelectedAttachment struct {
	Name string `json:"name"`
	Size int64 `json:"size"`
	Data []byte `json:"data"`
	Storage string `json:"storage"`
	Ext string `json:"ext"`
}
type SelectedAttachments []SelectedAttachment

type UnlockVaultCommandInterface interface {
	Execute(cmd UnlockVaultCommand) (*UnlockVaultResult, error) 
}
type UnlockVaultCommand struct {
	Password      string
	StellarSecret string
	UserID string
}

type UnlockVaultResult struct {
	VaultKey vaults_domain.VaultKey
}
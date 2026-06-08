package vault_dto

import (
	auth_domain "vault-app/internal/auth/domain"
	app_config_domain "vault-app/internal/config/domain"
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

type PrepareCommitRequest struct {
	UserID         string               `json:"user_id"`
	Password       string               `json:"password"`
	Vault          vault_domain.Vault   `json:"vault"`
	UserIdentity   identity_domain.User `json:"user_identity"`
	UserOnboarding string               `json:"user_onboarding"`
	Configs        app_config_domain.Config
	PrivateKey     string
}
type SynchronizeVaultRequest struct {
	UserID         string               `json:"user_id"`
	Password       string               `json:"password"`
	Vault          vault_domain.Vault   `json:"vault"`
	UserIdentity   identity_domain.User `json:"user_identity"`
	UserOnboarding string               `json:"user_onboarding"`
	Configs        app_config_domain.Config
	PrivateKey     string
}
type SynchronizeAttachmentRequest struct {
	UserID         string               `json:"user_id"`
	Password       string               `json:"password"`
	Vault          vault_domain.Vault   `json:"vault"`
	UserIdentity   identity_domain.User `json:"user_identity"`
	UserOnboarding string               `json:"user_onboarding"`
	Configs        app_config_domain.Config
	PrivateKey     string
}
type SelectedAttachment struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Data    []byte `json:"data"`
	Storage string `json:"storage"`
	Ext     string `json:"ext"`
}
type SelectedAttachments []SelectedAttachment

type UnlockVaultCommandInterface interface {
	Execute(cmd UnlockVaultCommand) (*UnlockVaultResult, error)
}
type UnlockVaultCommand struct {
	Password      string
	StellarSecret string
	UserID        string
}

type UnlockVaultResult struct {
	VaultKey vaults_domain.VaultKey
}

type AddAttachementsRequest struct {
	VaultName   string              `json:"vault_name"`
	EntryID     string              `json:"entry_id"`
	Password    string              `json:"password"`
	Attachments SelectedAttachments `json:"attachments"`
}
type AddAttachementRequest struct {
	UserID           string
	Data             []uint8
	Password         string
	EntryID          string
	VaultName        string
	UserOnboardingID string
	Configs          app_config_domain.Config
	Name             string
	Size             int64
	Ext              string
}

type DownloadShareAttachmentRequest struct {
	EncryptedKey  string
	AttachmentCID string
	FileExtension string
}

// internal/vault/application/usecase/onboarding_dto.go
type TemplateDTO struct {
	TemplateID    string
	RecordType    string
	SchemaVersion int
	Fields        map[string]any
}

type PackDTO struct {
	ID        string
	Templates []vaults_domain.TemplateRef `json:"templates"`
}

type VaultPayloadAddEntryRequest struct {
	UserID    string
	EntryType string
	Entry     any
}

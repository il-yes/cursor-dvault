package app_config

import (
	utils "vault-app/internal"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 🔧 Load this on boot from embedded JSON/YAML or system file (e.g., $HOME/.dvault/config.json)
type AppConfig struct {
	ID                 string           `json:"id" yaml:"id" gorm:"primaryKey"`
	RepoID             string           `json:"repo_id" yaml:"repo_id"`
	Branch             string           `json:"branch" yaml:"branch"`
	TracecoreEnabled   bool             `json:"tracecore_enabled" yaml:"tracecore_enabled"`
	CommitRules        []CommitRule     `json:"commit_rules" yaml:"commit_rules" gorm:"serializer:json"`
	BranchingModel     string           `json:"branching_model" yaml:"branching_model"`
	EncryptionPolicy   string           `json:"encryption_policy" yaml:"encryption_policy"`
	Actors             []string         `json:"actors" yaml:"actors" gorm:"type:json;serializer:json"`
	FederatedProviders []string         `json:"federated_providers" yaml:"federated_providers" gorm:"type:json;serializer:json"`
	DefaultPhase       string           `json:"default_phase" yaml:"default_phase"`
	DefaultVaultPath   string           `json:"default_vault_path" yaml:"default_vault_path"`
	VaultSettings      VaultConfig      `json:"vault_settings" yaml:"vault_settings" gorm:"embedded"`
	Blockchain         BlockchainConfig `json:"blockchain" yaml:"blockchain" gorm:"embedded"`
	UserID             int              `json:"user_id" gorm:"integer"`
}

func (a *AppConfig) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New().String()
	return
}

type CommitRule struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	AppConfigID string   `json:"-" gorm:"index"` // foreign key to AppConfig
	Rule        string   `json:"rule" yaml:"rule" gorm:"column:rule"`
	Actors      []string `json:"actors" yaml:"actors" gorm:"type:json"` // PostgreSQL array
}

func (c *CommitRule) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uint(utils.Uint64())
	return
}

// 🔐 Loaded after auth, encrypted at rest if persistent.
type UserConfig struct {
	ID             string               `json:"id" yaml:"id" gorm:"primaryKey"`
	Role           string               `json:"role" yaml:"role" gorm:"column:role"`
	Signature      string               `json:"signature" yaml:"signature" gorm:"column:signature"`
	ConnectedOrgs  []string             `json:"connected_orgs" yaml:"connected_orgs" gorm:"type:json;serializer:json"`
	StellarAccount StellarAccountConfig `json:"stellar_account" yaml:"stellar_account" gorm:"embedded;embeddedPrefix:stellar_"`
	SharingRules   []SharingRule        `json:"sharing_rules" yaml:"sharing_rules" gorm:"foreignKey:UserConfigID;constraint:OnDelete:CASCADE"`
}

type StellarAccountConfig struct {
	PublicKey  string `json:"public_key" yaml:"public_key" gorm:"column:public_key"`
	PrivateKey string `json:"private_key,omitempty" yaml:"private_key,omitempty" gorm:"column:private_key"` // should be encrypted
    EncPassword []byte `gorm:"column:enc_password"` // ciphertext
    EncNonce    []byte `gorm:"column:enc_nonce"`    // AES-GCM nonce

}

// internal/config/vault.go
type VaultConfig struct {
	MaxEntries       int    `json:"max_entries" yaml:"max_entries" gorm:"column:max_entries"`
	AutoSyncEnabled  bool   `json:"auto_sync_enabled" yaml:"auto_sync_enabled" gorm:"column:auto_sync_enabled"`
	EncryptionScheme string `json:"encryption_scheme" yaml:"encryption_scheme" gorm:"column:encryption_scheme"`
}

// internal/config/blockchain.go
type BlockchainConfig struct {
	Stellar StellarConfig `json:"stellar" yaml:"stellar" gorm:"embedded;embeddedPrefix:stellar_"`
	IPFS    IPFSConfig    `json:"ipfs" yaml:"ipfs" gorm:"embedded;embeddedPrefix:ipfs_"`
}

type StellarConfig struct {
	Network    string `json:"network" yaml:"network" gorm:"column:network"`
	HorizonURL string `json:"horizon_url" yaml:"horizon_url" gorm:"column:horizon_url"`
	Fee        int64  `json:"fee" yaml:"fee" gorm:"column:fee"`
}

type IPFSConfig struct {
	APIEndpoint string `json:"api_endpoint" yaml:"api_endpoint" gorm:"column:api_endpoint"`
	GatewayURL  string `json:"gateway_url" yaml:"gateway_url" gorm:"column:gateway_url"`
}

// 🔐 Loaded after auth, encrypted at rest if persistent.
type SharingRule struct {
	ID           uint     `json:"id" gorm:"primaryKey"`
	UserConfigID string   `json:"-" gorm:"index"` // foreign key
	EntryType    string   `json:"entry_type" yaml:"entry_type" gorm:"column:entry_type"`
	Targets      []string `json:"targets" yaml:"targets" gorm:"type:json;serializer:json"`
	Encrypted    bool     `json:"encrypted" yaml:"encrypted" gorm:"column:encrypted"`
}

func (s *SharingRule) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uint(utils.Uint64())
	return
}

type SharingConfig struct {
	ID              uint     `json:"id" gorm:"primaryKey"`
	RecipientID     string   `json:"recipient_id" yaml:"recipient_id"`
	Permissions     []string `json:"permissions" yaml:"permissions" gorm:"type:json;serializer:json"`
	Expiration      string   `json:"expiration" yaml:"expiration"`
	EncryptedFor    string   `json:"encrypted_for" yaml:"encrypted_for"`
	AnchorToChain   bool     `json:"anchor_to_chain" yaml:"anchor_to_chain"`
	NotifyRecipient bool     `json:"notify_recipient" yaml:"notify_recipient"`
	RedactedFields  []string `json:"redacted_fields,omitempty" yaml:"redacted_fields,omitempty" gorm:"type:json;serializer:json"`
	Revocable       bool     `json:"revocable" yaml:"revocable"`
}

func (s *SharingConfig) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uint(utils.Uint64())
	return
}

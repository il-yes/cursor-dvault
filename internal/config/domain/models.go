package app_config_domain

import (
	"vault-app/internal/blockchain"
	app_config "vault-app/internal/config"
	utils "vault-app/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

// -----------------------------
//
//	AppConfig
//
// -----------------------------
// 🔧 Load this on boot from embedded JSON/YAML or system file (e.g., $HOME/.dvault/config.json)
type AppConfig struct {
	ID                   string                   `json:"id" gorm:"primaryKey"`
	RepoID               string                   `json:"repo_id" gorm:"repo_id"`
	Branch               string                   `json:"branch" gorm:"branch"`
	TracecoreEnabled     bool                     `json:"tracecore_enabled" gorm:"tracecore_enabled"`
	CommitRules          []CommitRule             `json:"commit_rules" gorm:"serializer:json"`
	BranchingModel       string                   `json:"branching_model" gorm:"branching_model"`
	EncryptionPolicy     string                   `json:"encryption_policy" gorm:"encryption_policy"`
	Actors               []string                 `json:"actors" gorm:"type:json;serializer:json"`
	FederatedProviders   []string                 `json:"federated_providers" gorm:"type:json;serializer:json"`
	DefaultPhase         string                   `json:"default_phase" gorm:"default_phase"`
	DefaultVaultPath     string                   `json:"default_vault_path" gorm:"default_vault_path"`
	VaultSettings        VaultConfig              `json:"vault_settings" gorm:"embedded"`
	Blockchain           BlockchainConfig         `json:"blockchain" gorm:"embedded"`
	UserID               string                   `json:"user_id" gorm:"string"`
	AutoLockTimeout      string                   `json:"auto_lock_timeout" gorm:"auto_lock_timeout"`
	AccessPolicyDuration int64                    `json:"access_policy_duration" gorm:"access_policy_duration"`
	RemaskDelay          string                   `json:"remask_delay" gorm:"remask_delay"`
	Theme                string                   `json:"theme" gorm:"theme"`
	AnimationsEnabled    bool                     `json:"animations_enabled" gorm:"animations_enabled"`
	Storage              app_config.StorageConfig `json:"storage" yaml:"storage" gorm:"embedded"`
}

func (a *AppConfig) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New().String()
	return
}

// -----------------------------
//
//	UserConfig
//
// -----------------------------
// 🔐 Loaded after auth, encrypted at rest if persistent.
type UserConfig struct {
	ID               string               `json:"id" gorm:"primaryKey"` // -> UserId reference
	Role             string               `json:"role" gorm:"column:role"`
	Signature        string               `json:"signature" gorm:"column:signature"`
	ConnectedOrgs    []string             `json:"connected_orgs" gorm:"type:json;serializer:json"`
	StellarAccount   StellarAccountConfig `json:"stellar_account" gorm:"embedded;embeddedPrefix:stellar_"`
	SharingRules     []SharingRule        `json:"sharing_rules" gorm:"foreignKey:UserConfigID;constraint:OnDelete:CASCADE"`
	TwoFactorEnabled bool                 `json:"two_factor_enabled" yaml:"two_factor_enabled" gorm:"column:two_factor_enabled"`
	UI               UIConfig             `json:"ui" gorm:"embedded"`
}

func (u *UserConfig) OnGenerateApiKey(stellarAccount *StellarAccountConfig) {
	u.StellarAccount = *stellarAccount
}

type UIConfig struct {
	Theme             string `json:"theme" gorm:"theme"`
	AnimationsEnabled bool   `json:"animations_enabled" gorm:"animations_enabled"`
}

// -----------------------------
//
//	StellarAccountConfig
//
// -----------------------------
type StellarAccountConfig struct {
	PublicKey   string `json:"public_key" yaml:"public_key" gorm:"column:public_key"`
	PrivateKey  string `json:"private_key,omitempty" yaml:"private_key,omitempty" gorm:"column:private_key"` // should be encrypted
	EncPassword []byte `gorm:"column:enc_password"`                                                          // ciphertext
	EncNonce    []byte `gorm:"column:enc_nonce"`                                                             // AES-GCM nonce
	EncSalt     []byte `gorm:"column:enc_salt"`                                                              // salt
}

// -----------------------------
//
//	VaultConfig
//
// -----------------------------
type VaultConfig struct {
	MaxEntries       int    `json:"max_entries" yaml:"max_entries" gorm:"column:max_entries"`
	AutoSyncEnabled  bool   `json:"auto_sync_enabled" yaml:"auto_sync_enabled" gorm:"column:auto_sync_enabled"`
	EncryptionScheme string `json:"encryption_scheme" yaml:"encryption_scheme" gorm:"column:encryption_scheme"`
}

// -----------------------------
//
//	BlockchainConfig
//
// -----------------------------
// internal/config/blockchain.go
type BlockchainConfig struct {
	Stellar StellarConfig `json:"stellar" yaml:"stellar" gorm:"embedded;embeddedPrefix:stellar_"`
	IPFS    IPFSConfig    `json:"ipfs" yaml:"ipfs" gorm:"embedded;embeddedPrefix:ipfs_"`
}

// -----------------------------
//
//	StellarConfig
//
// -----------------------------
type StellarConfig struct {
	Network       string `json:"network" yaml:"network" gorm:"column:network"`
	HorizonURL    string `json:"horizon_url" yaml:"horizon_url" gorm:"column:horizon_url"`
	Fee           int64  `json:"fee" yaml:"fee" gorm:"column:fee"`
	SyncFrequency string `json:"sync_frequency" yaml:"sync_frequency" gorm:"column:sync_frequency"`
}

func NewStellarAccountConfigOnGeneratedApiKey(account *blockchain.CreateAccountRes) *StellarAccountConfig {
	return &StellarAccountConfig{
		PublicKey:   account.PublicKey,
		PrivateKey:  account.PrivateKey,
		EncSalt:     account.Salt,
		EncNonce:    account.EncNonce,
		EncPassword: account.EncPassword,
	}

}

// -----------------------------
//
//	IPFSConfig
//
// -----------------------------
type IPFSConfig struct {
	APIEndpoint string `json:"api_endpoint" yaml:"api_endpoint" gorm:"column:api_endpoint"`
	GatewayURL  string `json:"gateway_url" yaml:"gateway_url" gorm:"column:gateway_url"`
}

// -----------------------------
//
//	SharingRule
//
// -----------------------------
// 🔐 Loaded after auth, encrypted at rest if persistent.
type SharingRule struct {
	ID           string   `json:"id" gorm:"primaryKey;autoIncrement:false;size:36"`
	UserConfigID string   `json:"-" gorm:"index"` // foreign key
	EntryType    string   `json:"entry_type" yaml:"entry_type" gorm:"column:entry_type"`
	Targets      []string `json:"targets" yaml:"targets" gorm:"type:json;serializer:json"`
	Encrypted    bool     `json:"encrypted" yaml:"encrypted" gorm:"column:encrypted"`
}

func (s *SharingRule) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

// -----------------------------
//
//	SharingConfig
//
// -----------------------------
type SharingConfig struct {
	ID              string   `json:"id" gorm:"primaryKey"`
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
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

type BaseVaultConfig struct {
	ID        string `json:"id" gorm:"primaryKey"`
	UserID    string `json:"user_id" gorm:"column:user_id"`
	VaultName string `json:"vault_name" gorm:"column:vault_name"`
}

// -----------------------------
//
//	FeatureFlags
//
// -----------------------------
type FeatureFlags struct {
	TracecoreEnabled        bool `json:"tracecore_enabled"`
	CloudBackupEnabled      bool `json:"cloud_backup_enabled"`
	ThreatDetectionEnabled  bool `json:"threat_detection_enabled"`
	BrowserExtensionEnabled bool `json:"browser_extension_enabled"`
	GitCLIEnabled           bool `json:"git_cli_enabled"`
}

// -----------------------------
//
//	SyncConfig
//
// -----------------------------
type SyncConfig struct {
	BaseVaultConfig
	AutoSync            bool   `json:"auto_sync" gorm:"column:auto_sync"`
	SyncIntervalSeconds int    `json:"sync_interval_seconds" gorm:"column:sync_interval_seconds"`
	ConflictStrategy    string `json:"conflict_strategy" gorm:"column:conflict_strategy"`
	MaxRetries          int    `json:"max_retries" gorm:"column:max_retries"`
	StellarFrequency    string `json:"stellar_frequency" gorm:"column:stellar_frequency"`
}

func (s *SyncConfig) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

type VaultConfigBeta struct {
	BaseVaultConfig

	Features   FeatureFlags     `json:"features" gorm:"embedded;embeddedPrefix:features_"`
	Sync       SyncConfig       `json:"sync" gorm:"embedded;embeddedPrefix:sync_"`
	Backup     BackupConfig     `json:"backup" gorm:"embedded;embeddedPrefix:backup_"`
	Privacy    PrivacyConfig    `json:"privacy" gorm:"embedded;embeddedPrefix:privacy_"`
	Security   SecurityConfig   `json:"security" gorm:"embedded;embeddedPrefix:security_"`
	Sharing    SharingPolicy    `json:"sharing" gorm:"embedded;embeddedPrefix:sharing_"`
	Onboarding OnboardingConfig `json:"onboarding" gorm:"embedded;embeddedPrefix:onboarding_"`
}

func (dc *VaultConfigBeta) BeforeCreate(tx *gorm.DB) (err error) {
	if dc.ID == "" {
		dc.ID = uuid.New().String()
	}
	return nil
}

// -----------------------------
//
//	BackupConfig
//
// -----------------------------
type BackupConfig struct {
	BaseVaultConfig
	Enabled       bool   `json:"enabled" gorm:"column:enabled"`
	Schedule      string `json:"schedule" gorm:"column:schedule"`
	RetentionDays int    `json:"retention_days" gorm:"column:retention_days"`
	Encryption    bool   `json:"encryption" gorm:"column:encryption"`
}

func (b *BackupConfig) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return nil
}

// -----------------------------
//
//	PrivacyConfig
//
// -----------------------------
type PrivacyConfig struct {
	BaseVaultConfig
	TelemetryEnabled bool `json:"telemetry_enabled" gorm:"column:telemetry_enabled"`
	AnonymousMode    bool `json:"anonymous_mode" gorm:"column:anonymous_mode"`
}

func (p *PrivacyConfig) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// -----------------------------
//
//	NotificationConfig
//
// -----------------------------
type NotificationConfig struct {
	BaseVaultConfig
	BillingAlerts  bool `json:"billing_alerts" gorm:"column:billing_alerts"`
	SecurityAlerts bool `json:"security_alerts" gorm:"column:security_alerts"`
	ShareInvites   bool `json:"share_invites" gorm:"column:share_invites"`
	SystemUpdates  bool `json:"system_updates" gorm:"column:system_updates"`
}

func (nc *NotificationConfig) BeforeCreate(tx *gorm.DB) (err error) {
	if nc.ID == "" {
		nc.ID = uuid.New().String()
	}
	return nil
}

// -----------------------------
//
//	DeviceConfig
//
// -----------------------------
type DeviceConfig struct {
	BaseVaultConfig
	DeviceID   string `json:"device_id" gorm:"column:device_id"`
	DeviceName string `json:"device_name" gorm:"column:device_name"`
	Trusted    bool   `json:"trusted" gorm:"column:trusted"`
	LastSync   int64  `json:"last_sync" gorm:"column:last_sync"`
}

func (dc *DeviceConfig) BeforeCreate(tx *gorm.DB) (err error) {
	if dc.ID == "" {
		dc.ID = uuid.New().String()
	}
	return nil
}

// -----------------------------
//
//	SecurityConfig
//
// -----------------------------
type SecurityConfig struct {
	AutoLockSeconds     int
	SessionTimeout      int
	RequireBiometric    bool
	ClearClipboardAfter int
}

// -----------------------------
//
//	SharingPolicy
//
// -----------------------------
type SharingPolicy struct {
	AllowExternalSharing bool
	DefaultExpiryHours   int
	RequirePassword      bool
	MaxSharesPerEntry    int
}

// -----------------------------
//
//	DevicePolicy
//
// -----------------------------
type DevicePolicy struct {
	MaxDevices             int
	RequireDeviceApproval  bool
	AutoRevokeInactiveDays int
}

// -----------------------------
//
//	SubscriptionConfig
//
// -----------------------------
type SubscriptionConfig struct {
	BaseVaultConfig
	Plan     string             `json:"plan" gorm:"column:plan"`
	Features FeatureFlags       `json:"features" gorm:"embedded;embeddedPrefix:features_"`
	Limits   SubscriptionLimits `json:"limits" gorm:"embedded;embeddedPrefix:limits_"`
}

func (sc *SubscriptionConfig) BeforeCreate(tx *gorm.DB) (err error) {
	if sc.ID == "" {
		sc.ID = uuid.New().String()
	}
	return nil
}

// -----------------------------
//
//	SubscriptionLimits
//
// -----------------------------
type SubscriptionLimits struct {
	MaxVaults  int `json:"max_vaults"`
	MaxUsers   int `json:"max_users"`
	MaxDevices int `json:"max_devices"`
	MaxShares  int `json:"max_shares"`
	MaxStorage int `json:"maxstorage"`
}

// -----------------------------
//
//	OnboardingConfig
//
// -----------------------------
type OnboardingConfig struct {
	Packs              []string `json:"packs" gorm:"type:json"`               // ["compliance_critical","client_data"]
	UseCases           []string `json:"use_cases" gorm:"type:json"`           // optional UI labels if you want
	InstalledTemplates []string `json:"installed_templates" gorm:"type:json"` // ["devsecops.incident.v1", ...]
	Completed          bool     `json:"completed" gorm:"column:completed"`
}

package app_config_domain

import (
	app_config "vault-app/internal/config"

	"github.com/google/uuid"
)

const (
	TracecoreEnabled        = "tracecore_enabled"
	CloudBackupEnabled      = "cloud_backup_enabled"
	ThreatDetectionEnabled  = "threat_detection_enabled"
	BrowserExtensionEnabled = "browser_extension_enabled"
	GitCLIEnabled           = "git_cli_enabled"
)

const (
	DefaultSyncInterval         = 60
	DefaultStellarFrequency     = "manual"
	DefaultConflictStrategy     = "last_write_wins"
	DefaultMaxRetries           = 3
	DefaultBranch               = "main"
	DefaultDefaultPhase         = "vault_entry"
	DefaultMaxEntries           = 1000
	DefaultMaxShares            = 10
	DefaultAutoLockTimeout      = 300
	DefaultAccessPolicyDuration = 3600
	DefaultRemaskDelay          = 60
	DefaultTheme                = "dark"
	DefaultAnimationsEnabled    = true
	DefaultSessionTimeout       = 60
	DefaultClearClipboardAfter  = 60
	DefaultMaxVaults            = 1
	DefaultMaxUsers             = 1
	DefaultMaxDevices           = 1

	BackupSchedule      = "daily"
	BackupRetentionDays = 30

	DeviceName = "My Device"
)

type Config struct {
	App          *AppConfig
	User         *UserConfig
	Subscription *SubscriptionConfig
	Vaults       VaultConfigBeta
	Devices      []DeviceConfig
}

func InitConfig(userID string) (*Config, error) {
	// Load from environment variables, files, CLI flags, or defaults
	// For example, use `viper` or `envconfig` libraries for flexibility

	cfg := &Config{
		App: &AppConfig{
			UserID:       userID,
			RepoID:       "my-repo-id",
			Branch:       "main",
			DefaultPhase: "vault_entry",
			VaultSettings: VaultConfig{
				MaxEntries:       DefaultMaxRetries,
				AutoSyncEnabled:  true,
				EncryptionScheme: "AES-256-GCM",
			},
			Blockchain: BlockchainConfig{
				Stellar: StellarConfig{
					Network:    "testnet",
					HorizonURL: "https://horizon-testnet.stellar.org",
					Fee:        100,
				},
				IPFS: IPFSConfig{
					APIEndpoint: "http://localhost:5001",
					GatewayURL:  "https://ipfs.io/ipfs/",
				},
			},
			Storage: app_config.StorageConfig{
				Mode: app_config.StorageCloud, // ← production default

				LocalIPFS: app_config.IPFSConfig{
					APIEndpoint: "http://localhost:5001",
					GatewayURL:  "https://ipfs.io/ipfs/",
				},

				PrivateIPFS: app_config.IPFSConfig{
					APIEndpoint: "http://192.168.1.10:5001",
					GatewayURL:  "http://192.168.1.10:8080/ipfs/",
				},

				Cloud: app_config.CloudConfig{
					BaseURL: "https://ankhora.io/back",
				},

				EnterpriseS3: app_config.S3Config{
					Region:   "us-east-1",
					Bucket:   "ankhora-enterprise",
					Endpoint: "https://s3.us-east-1.amazonaws.com",
				},
			},
		},
		User: &UserConfig{
			ID:            userID,
			Role:          "user",
			Signature:     "",
			ConnectedOrgs: []string{""},
			StellarAccount: StellarAccountConfig{
				PublicKey:   "",
				PrivateKey:  "",
				EncSalt:     []byte(""),
				EncNonce:    []byte(""),
				EncPassword: []byte(""),
			},
			SharingRules: []SharingRule{
				{
					ID:           uuid.NewString(),
					UserConfigID: "",
					EntryType:    "",
					Targets:      []string{""},
					Encrypted:    false,
				},
			},
			TwoFactorEnabled: false,
		},
	}

	return cfg, nil
}

func InitConfigFromVault(userID string, vaultName string) (*Config, error) {
	cfg := &Config{
		App: &AppConfig{
			UserID:       userID,
			RepoID:       "my-repo-id",
			Branch:       DefaultBranch,
			DefaultPhase: DefaultDefaultPhase,
			VaultSettings: VaultConfig{
				MaxEntries:       DefaultMaxEntries,
				AutoSyncEnabled:  true,
				EncryptionScheme: "AES-256-GCM",
			},
			Blockchain: BlockchainConfig{
				Stellar: StellarConfig{
					Network:    "testnet",
					HorizonURL: "https://horizon-testnet.stellar.org",
					Fee:        100,
				},
				IPFS: IPFSConfig{
					APIEndpoint: "http://localhost:5001",
					GatewayURL:  "https://ipfs.io/ipfs/",
				},
			},
			Storage: app_config.StorageConfig{
				Mode: app_config.StorageCloud, // ← production default

				LocalIPFS: app_config.IPFSConfig{
					APIEndpoint: "http://localhost:5001",
					GatewayURL:  "https://ipfs.io/ipfs/",
				},

				PrivateIPFS: app_config.IPFSConfig{
					APIEndpoint: "http://192.168.1.10:5001",
					GatewayURL:  "http://192.168.1.10:8080/ipfs/",
				},

				Cloud: app_config.CloudConfig{
					BaseURL: "https://ankhora.io/back",
				},

				EnterpriseS3: app_config.S3Config{
					Region:   "us-east-1",
					Bucket:   "ankhora-enterprise",
					Endpoint: "https://s3.us-east-1.amazonaws.com",
				},
			},
		},
		User: &UserConfig{
			ID:            userID,
			Role:          "user",
			Signature:     "",
			ConnectedOrgs: []string{""},
			StellarAccount: StellarAccountConfig{
				PublicKey:   "",
				PrivateKey:  "",
				EncSalt:     []byte(""),
				EncNonce:    []byte(""),
				EncPassword: []byte(""),
			},
			SharingRules: []SharingRule{
				{
					ID:           uuid.NewString(),
					UserConfigID: "",
					EntryType:    "",
					Targets:      []string{""},
					Encrypted:    false,
				},
			},
			TwoFactorEnabled: false,
			UI: UIConfig{
				Theme:             DefaultTheme,
				AnimationsEnabled: DefaultAnimationsEnabled,
			},
		},
		Vaults: VaultConfigBeta{
			BaseVaultConfig: BaseVaultConfig{
				ID:        uuid.NewString(),
				UserID:    userID,
				VaultName: vaultName,
			},
			Features: FeatureFlags{
				TracecoreEnabled:        true,
				CloudBackupEnabled:      true,
				ThreatDetectionEnabled:  true,
				BrowserExtensionEnabled: true,
				GitCLIEnabled:           true,
			},
			Sync: SyncConfig{
				AutoSync:            true,
				SyncIntervalSeconds: DefaultSyncInterval,
				ConflictStrategy:    DefaultConflictStrategy,
				MaxRetries:          DefaultMaxRetries,
				StellarFrequency:    DefaultStellarFrequency,
			},
			Backup: BackupConfig{
				Enabled:       true,
				Schedule:      BackupSchedule,
				RetentionDays: BackupRetentionDays,
				Encryption:    true,
			},
			Privacy: PrivacyConfig{
				TelemetryEnabled: true,
				AnonymousMode:    false,
			},
			Security: SecurityConfig{
				AutoLockSeconds:     DefaultAutoLockTimeout,
				SessionTimeout:      DefaultSessionTimeout,
				RequireBiometric:    true,
				ClearClipboardAfter: DefaultClearClipboardAfter,
			},
			Sharing: SharingPolicy{
				AllowExternalSharing: true,
				DefaultExpiryHours:   24,
				RequirePassword:      true,
				MaxSharesPerEntry:    10,
			},
			Onboarding: OnboardingConfig{
				Packs:              []string{},
				UseCases:           []string{},
				InstalledTemplates: []string{},
				Completed:          false,
			},
		},
		Devices: []DeviceConfig{
			{
				BaseVaultConfig: BaseVaultConfig{
					ID:        uuid.NewString(),
					UserID:    userID,
					VaultName: vaultName,
				},
				DeviceID:   uuid.NewString(),
				DeviceName: DeviceName,
				LastSync:   0,
			},
		},
		Subscription: &SubscriptionConfig{
			BaseVaultConfig: BaseVaultConfig{
				ID:        uuid.NewString(),
				UserID:    userID,
				VaultName: vaultName,
			},
			Plan: "free",
			Features: FeatureFlags{
				TracecoreEnabled:        true,
				CloudBackupEnabled:      true,
				ThreatDetectionEnabled:  true,
				BrowserExtensionEnabled: true,
				GitCLIEnabled:           true,
			},
			Limits: SubscriptionLimits{
				MaxVaults:  DefaultMaxVaults,
				MaxUsers:   DefaultMaxUsers,
				MaxDevices: DefaultMaxDevices,
				MaxShares:  DefaultMaxShares,
			},
		},
	}

	return cfg, nil
}

// func (c *Config) ValidateFeatures(sub SubscriptionSnapshot) {

// 	if !sub.Tracecore {
// 		c.Features.Flags["tracecore"] = false
// 	}

// 	if !sub.CloudBackup {
// 		c.Backup.Enabled = false
// 	}

// 	if !sub.BrowserExtension {
// 		c.Features.Flags["browser_extension"] = false
// 	}
// }

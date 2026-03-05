package app_config_domain

import (
	app_config "vault-app/internal/config"

	"github.com/google/uuid"
)


type Config struct {
	App  AppConfig
	User UserConfig
}

func InitConfig(userID string) (*Config, error) {
	// Load from environment variables, files, CLI flags, or defaults
	// For example, use `viper` or `envconfig` libraries for flexibility

	cfg := &Config{
		App: AppConfig{
			UserID: userID,
			RepoID:       "my-repo-id",
			Branch:       "main",
			DefaultPhase: "vault_entry",
			VaultSettings: VaultConfig{
				MaxEntries:       1000,
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
		User: UserConfig{
			ID: userID,
			Role: "user",
			Signature: "",
			ConnectedOrgs: []string{""},
			StellarAccount: StellarAccountConfig{
				PublicKey: "",
				PrivateKey: "",
				EncSalt: []byte(""),							
				EncNonce: []byte(""),
				EncPassword: []byte(""),
			},	
			SharingRules: []SharingRule{
				{
					ID: uuid.NewString(),
					UserConfigID: "",
					EntryType: "",
					Targets: []string{""},
					Encrypted: false,
				},
			},
			TwoFactorEnabled: false,
		},
	}
	// Optional: override with env vars or config files here

	return cfg, nil
}
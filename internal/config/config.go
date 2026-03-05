package app_config

type Config struct {
	App  AppConfig
	User UserConfig
}

// func LoadConfig() (*Config, error) {
// 	// Load from environment variables, files, CLI flags, or defaults
// 	// For example, use `viper` or `envconfig` libraries for flexibility

// 	cfg := &Config{
// 		App: AppConfig{
// 			RepoID:       "my-repo-id",
// 			Branch:       "main",
// 			DefaultPhase: "vault_entry",
// 			VaultSettings: VaultConfig{
// 				MaxEntries:       1000,
// 				AutoSyncEnabled:  true,
// 				EncryptionScheme: "AES-256-GCM",
// 			},
// 			Blockchain: BlockchainConfig{
// 				Stellar: StellarConfig{
// 					Network:    "testnet",
// 					HorizonURL: "https://horizon-testnet.stellar.org",
// 					Fee:        100,
// 				},
// 				IPFS: IPFSConfig{
// 					APIEndpoint: "http://localhost:5001",
// 					GatewayURL:  "https://ipfs.io/ipfs/",
// 				},
// 			},
// 			Storage: StorageConfig{
// 				Mode: StorageCloud, // ← production default

// 				LocalIPFS: IPFSConfig{
// 					APIEndpoint: "http://localhost:5001",
// 					GatewayURL:  "https://ipfs.io/ipfs/",
// 				},

// 				PrivateIPFS: IPFSConfig{
// 					APIEndpoint: "http://192.168.1.10:5001",
// 					GatewayURL:  "http://192.168.1.10:8080/ipfs/",
// 				},

// 				Cloud: CloudConfig{
// 					BaseURL: "https://api.ankhora.io",
// 				},

// 				EnterpriseS3: S3Config{
// 					Region:   "us-east-1",
// 					Bucket:   "ankhora-enterprise",
// 					Endpoint: "https://s3.us-east-1.amazonaws.com",
// 				},
// 			},
// 		},
// 	}
// 	// Optional: override with env vars or config files here

// 	return cfg, nil
// }

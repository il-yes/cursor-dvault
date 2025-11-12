package app_config


type Config struct {
    App  AppConfig
    User UserConfig
}

func LoadConfig() (*Config, error) {
    // Load from environment variables, files, CLI flags, or defaults
    // For example, use `viper` or `envconfig` libraries for flexibility
    cfg := &Config{
        App: AppConfig{
            RepoID:       "my-repo-id",
            Branch:       "main",
            // CommitRules:  []string{"REQUIRES_SIGNATURE", "VALID_ACTORS_ONLY"},
            DefaultPhase: "vault_entry",
            VaultSettings: VaultConfig{
                MaxEntries:      1000,
                AutoSyncEnabled: true,
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
        },
        User: UserConfig{
            ID:        "",
            Role:      "end_user",
            Signature: "",
        },
    }
    // Optional: override with env vars or config files here

    return cfg, nil
}

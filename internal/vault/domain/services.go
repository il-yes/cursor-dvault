package vaults_domain

import "context"


type VaultStorage interface {
	GetData(cid string) ([]byte, error)
}

type CryptoService interface {
	Decrypt(data []byte, password string) ([]byte, error)
}

// In your vault_commands package or beneath it:

type StorageProvider interface {
    Add(data []byte) (string, error)
    // optional: Get(ctx context.Context, cid string) ([]byte, error)
	AddData(ctx context.Context, data []byte) (string, error)
}

package vaults_domain

import (
	"context"
)


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



type VaultCrypto interface {
	Encrypt(data []byte, key []byte) ([]byte, error)
	Decrypt(data []byte, key []byte) ([]byte, error)
}

type KeyEncryption interface {
	WrapKeyWithPassword(vaultKey []byte, password string) ([]byte, error)
	UnwrapKeyWithPassword(enc []byte, password string) ([]byte, error)

	WrapKeyWithStellar(vaultKey []byte, stellarSecret string) ([]byte, error)
	UnwrapKeyWithStellar(enc []byte, stellarSecret string) ([]byte, error)
}

type AsymmetricCrypto interface {
	EncryptForRecipient(pubKey string, data []byte) ([]byte, error)
}
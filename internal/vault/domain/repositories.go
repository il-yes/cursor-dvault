package vaults_domain

type VaultRepository interface {
	SaveVault(vault *Vault) error
	GetVault(vaultID string) (*Vault, error)
	UpdateVault(vault *Vault) error
	DeleteVault(vaultID string) error

	GetLatestByUserID(userID string) (*Vault, error)
}


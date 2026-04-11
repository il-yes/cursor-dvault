package vaults_domain

type VaultRepository interface {
	SaveVault(vault *Vault) error
	GetVault(vaultID string) (*Vault, error)
	UpdateVault(vault *Vault) error
	DeleteVault(vaultID string) error

	GetLatestByUserID(userID string) (*Vault, error)
	GetByUserIDAndName(userID string, name string) (*Vault, error)
	UpdateVaultCID(vaultID, cid string) error
}

type FolderRepository interface {
	SaveFolder(folder *Folder) error
	GetFolder(folderID string) (*Folder, error)
	UpdateFolder(folder *Folder) error
	DeleteFolder(folderID string) error

	GetFoldersByUserID(userID string) ([]Folder, error)
	GetFoldersByVault(vaultCID string) ([]Folder, error)	
	GetFolderById(id string) (*Folder, error)
}




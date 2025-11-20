package vault_domain

import (
	"vault-app/internal/models"
)

type Repository interface {
	CreateVault(vault models.VaultPayload) (*models.VaultPayload, error)
	UpdateVault(vault models.VaultPayload) (*models.VaultPayload, error)
	GetVaultByID(id int) (*models.VaultPayload, error)
	GetVaultByUserID(userID int) (*models.VaultPayload, error)
	DeleteVault(id int) error
}

package collaboration_ports

import (
	"context"

	vaults_domain "vault-app/internal/vault/domain"
)

type SovereignIdentityResolver interface {
	GetDeviceSeed(ctx context.Context, userID string) (string, error)
	GetVaultKeyring(ctx context.Context, userID string) (*vaults_domain.VaultKeyring, error)
}

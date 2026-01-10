package vault_events

import (
	"context"
	vault_domain "vault-app/internal/vault/domain"
	vault_session "vault-app/internal/vault/application/session"
)

// -------- EVENTS --------
type VaultOpened struct {
	UserID  string
	VaultPayload *vault_domain.VaultPayload
    Dirty   bool
    LastCID string

    LastSynced     string
    LastUpdated    string
    Runtime        *vault_session.RuntimeContext
	OccurredAt int64
}


// -------- EVENT BUS --------
type VaultEventBus interface {
	PublishVaultOpened(ctx context.Context, event VaultOpened) error
	SubscribeToVaultOpened(handler func(ctx context.Context, event VaultOpened)) error
}	
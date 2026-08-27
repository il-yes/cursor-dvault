package vault_commands

import (
	"context"
	vault_session "vault-app/internal/vault/application/session"
)

type CloseCommand struct {
	UserID string
}

type CloseCommandHandler struct {
	SessionManager *vault_session.Manager
}

func (c *CloseCommandHandler) Execute(ctx context.Context, cmd CloseCommand) error {
	if c.SessionManager != nil && cmd.UserID != "" {
		c.SessionManager.CloseSession(cmd.UserID)
	}
	return nil
}

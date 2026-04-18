package vault_commands

import "context"

type CloseCommand struct{}

type CloseCommandHandler struct {
}

func (c *CloseCommandHandler) Execute(ctx context.Context, cmd CloseCommand) error {
	// 1. wipe VaultKey from memory
	// 2. clear decrypted cache
	// 3. stop sync workers

	return nil
}

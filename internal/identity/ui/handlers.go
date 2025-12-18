package identity_ui

import (
	identity_commands "vault-app/internal/identity/application/commands"
)

type LoginHandler struct {
	loginCommandHandler *identity_commands.LoginCommandHandler
}

func NewLoginHandler(
	loginCommandHandler *identity_commands.LoginCommandHandler,
) *LoginHandler {
	return &LoginHandler{
		loginCommandHandler: loginCommandHandler,
	}
}

func (h *LoginHandler) Handle(cmd identity_commands.LoginCommand) (*identity_commands.LoginResult, error) {
	return h.loginCommandHandler.Handle(cmd)
}

package identity_ui

import (
	"context"
	identity_domain "vault-app/internal/identity/domain"
	// Added missing import for OnboardRequest
)

type IdentityHandler struct {
	LoginHandler        *LoginHandler
	RegistrationHandler *RegistrationHandler
}

func NewIdentityHandler(loginHandler *LoginHandler, registrationHandler *RegistrationHandler) *IdentityHandler {
	return &IdentityHandler{
		LoginHandler:        loginHandler,
		RegistrationHandler: registrationHandler,
	}
}

func (h *IdentityHandler) Registers(ctx context.Context, req OnboardRequest) (*identity_domain.User, error) {
	return h.RegistrationHandler.Registers(ctx, req)
}

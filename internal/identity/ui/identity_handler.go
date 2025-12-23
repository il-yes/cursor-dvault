package identity_ui

import (
	"context"
	identity_domain "vault-app/internal/identity/domain"
)

type IdentityHandler struct {
	LoginHandler        *LoginHandler
	RegistrationHandler *RegistrationHandler
	Finder *FinderHandler
}

func NewIdentityHandler(loginHandler *LoginHandler, registrationHandler *RegistrationHandler, finderHandler *FinderHandler) *IdentityHandler {
	return &IdentityHandler{
		LoginHandler:        loginHandler,
		RegistrationHandler: registrationHandler,
		Finder: finderHandler,
	}
}

func (h *IdentityHandler) Registers(ctx context.Context, req OnboardRequest) (*identity_domain.User, error) {
	return h.RegistrationHandler.Registers(ctx, req)
}


// finderQueryHandler := NewFinderQueryHandler(identity_persistence.NewUserRepository())
// loginCommandHandler := NewLoginCommandHandler()
// finderH := NewFinderHandler(finderQueryHandler)
// loginH := NewLoginHandler(loginCommandHandler)
// registerUC := NewRegisterIdentityUseCase()
// registrationH := NewRegistrationHandler(registerUC)

// identityH := NewIdentityHandler(loginH, registrationH, finderH)



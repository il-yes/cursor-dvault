package identity_ui

import (
	"context"
	identity_commands "vault-app/internal/identity/application/commands"
	identity_usecase "vault-app/internal/identity/application/usecase"
	identity_domain "vault-app/internal/identity/domain"
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


type RegistrationHandler struct {
	RegisterIdentityUC *identity_usecase.RegisterIdentityUseCase
}
type OnboardRequest struct {
	Identity             string
	Email                string
	Password             string
	IsAnonymous          bool
	StellarPublicKey     string
	Tier                 string
	PaymentMethod        string
	EncryptedPaymentData string
	SubscriptionID       string
}
func (h *RegistrationHandler) Registers(ctx context.Context, req OnboardRequest) (*identity_domain.User, error) {
	var publicKey string
	if req.StellarPublicKey != "" {
		publicKey = req.StellarPublicKey
	}
	registerRequest := identity_usecase.RegisterRequest{
		Email:            req.Email,
		Password:         req.Password,
		IsAnonymous:      req.IsAnonymous,
		StellarPublicKey: publicKey,
	}
	identity_user, err := h.RegisterIdentityUC.Execute(ctx, registerRequest)
	if err != nil {
		return nil, err
	}
	
	return identity_user, nil
}

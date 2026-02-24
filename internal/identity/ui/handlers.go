package identity_ui

import (
	"context"
	utils "vault-app/internal/utils"
	identity_commands "vault-app/internal/identity/application/commands"
	identity_queries "vault-app/internal/identity/application/queries"
	identity_usecase "vault-app/internal/identity/application/usecase"
	identity_domain "vault-app/internal/identity/domain"
	identity_eventbus "vault-app/internal/identity/application"
)

// ----------------------------------------
// - 	LoginHandler 
// ----------------------------------------
type LoginHandler struct {
	loginCommandHandler *identity_commands.LoginCommandHandler
	tokenService        identity_commands.TokenServiceInterface
	eventBus            identity_eventbus.EventBus
}

func NewLoginHandler(
	loginCommandHandler *identity_commands.LoginCommandHandler,
	tokenService        identity_commands.TokenServiceInterface,
	eventBus            identity_eventbus.EventBus,
) *LoginHandler {
	return &LoginHandler{
		loginCommandHandler: loginCommandHandler,
		tokenService:        tokenService,
		eventBus:            eventBus,
	}
}

func (h *LoginHandler) Handle(cmd identity_commands.LoginCommand) (*identity_commands.LoginResult, error) {
	return h.loginCommandHandler.Handle(cmd, h.tokenService, h.eventBus)
}

// ----------------------------------------
// - RegistrationHandler 
// ----------------------------------------
type RegistrationHandler struct {
	RegisterIdentityUC *identity_usecase.RegisterIdentityUseCase
}
func NewRegistrationHandler(
	RegisterIdentityUC *identity_usecase.RegisterIdentityUseCase,
) *RegistrationHandler {
	return &RegistrationHandler{
		RegisterIdentityUC: RegisterIdentityUC,
	}
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
	identity_user, err := h.RegisterIdentityUC.RegisterIdentity(ctx, registerRequest)
	if err != nil {
		return nil, err
	}
	
	return identity_user, nil
}

// ----------------------------------------
// - FinderHandler 
// ----------------------------------------
type FinderHandler struct {
	finderQueryHandler *identity_queries.FinderQueryHandler
}
func NewFinderHandler(
	finderQueryHandler *identity_queries.FinderQueryHandler,
) *FinderHandler {
	return &FinderHandler{
		finderQueryHandler: finderQueryHandler,
	}
}	
func (h *FinderHandler) FindByEmail(ctx context.Context, email string) (*identity_domain.User, error) {
	return h.finderQueryHandler.FindByEmail(ctx, email)
}

func (h *FinderHandler) FindById(ctx context.Context, id string) (*identity_domain.User, error) {
	utils.LogPretty("FinderHandler - FindById - id", id)
	return h.finderQueryHandler.FindById(ctx, id)
}
	

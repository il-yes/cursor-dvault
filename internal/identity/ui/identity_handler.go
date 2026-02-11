package identity_ui

import (
	"context"
	identity_eventbus "vault-app/internal/identity/application"
	identity_commands "vault-app/internal/identity/application/commands"
	identity_queries "vault-app/internal/identity/application/queries"
	identity_usecase "vault-app/internal/identity/application/usecase"
	identity_domain "vault-app/internal/identity/domain"
	identity_infrastructure_eventbus "vault-app/internal/identity/infrastructure/eventbus"
	identity_persistence "vault-app/internal/identity/infrastructure/persistence"
	onboarding_domain "vault-app/internal/onboarding/domain"

	gorm "gorm.io/gorm"
)

type IdentityHandler struct {
	DB *gorm.DB
	IdentityUserRepo identity_domain.UserRepository

	tokenService        identity_commands.TokenServiceInterface
	eventBus            identity_eventbus.EventBus

	loginHandler *LoginHandler
	registrationHandler *RegistrationHandler
	finderHandler *FinderHandler
	Bus identity_eventbus.EventBus
}

func NewIdentityHandler(
	db *gorm.DB, 
	tokenService identity_commands.TokenServiceInterface,
	onboardingUserRepo onboarding_domain.UserRepository,
	) *IdentityHandler {
	
	identityUserRepo := identity_persistence.NewGormUserRepository(db)
	identityMemoryBus := identity_infrastructure_eventbus.NewMemoryEventBus()	

	return &IdentityHandler{
		DB: db,
		IdentityUserRepo: identityUserRepo,
		eventBus: identityMemoryBus,
		tokenService:        tokenService,
		Bus: identityMemoryBus,
	}
}

func (h *IdentityHandler) Registers(req OnboardRequest) (*identity_domain.User, error) {
	identityIdGen := identity_persistence.NewIDGenerator()
	registerStandardUserUseCase := identity_usecase.NewRegisterStandardUserUseCase(h.IdentityUserRepo, h.eventBus, identityIdGen)
	registerAnonymousUserUseCase := identity_usecase.NewRegisterAnonymousUserUseCase(h.IdentityUserRepo, h.eventBus, identityIdGen)
	identityRegistrationHandler := NewRegistrationHandler(
		identity_usecase.NewRegisterIdentityUseCase(
			registerStandardUserUseCase, 
			registerAnonymousUserUseCase,
		),
	)
	
	return identityRegistrationHandler.Registers(context.Background(), req)
}


func (h *IdentityHandler) Login(req identity_commands.LoginCommand) (*identity_commands.LoginResult, error) {
	
	loginHandler := NewLoginHandler(
		identity_commands.NewLoginCommandHandler(h.DB),
		h.tokenService,
		h.eventBus,
	)
		
	return loginHandler.Handle(req)
}

func (h *IdentityHandler) FindUserByEmail(ctx context.Context, req identity_queries.ReqQuery) (*identity_domain.User, error) {
	identityFinderHandler := NewFinderHandler(
		identity_queries.NewFinderQueryHandler(h.IdentityUserRepo),
	)
	return identityFinderHandler.FindByEmail(ctx, req.Email)
}

func (h *IdentityHandler) FindUserById(ctx context.Context, req string) (*identity_domain.User, error) {
	identityFinderHandler := NewFinderHandler(
		identity_queries.NewFinderQueryHandler(h.IdentityUserRepo),
	)
	return identityFinderHandler.FindById(ctx, req)
}

// finderQueryHandler := NewFinderQueryHandler(identity_persistence.NewUserRepository())
// loginCommandHandler := NewLoginCommandHandler()
// finderH := NewFinderHandler(finderQueryHandler)
// loginH := NewLoginHandler(loginCommandHandler)
// registerUC := NewRegisterIdentityUseCase()
// registrationH := NewRegistrationHandler(registerUC)

// identityH := NewIdentityHandler(loginH, registrationH, finderH)



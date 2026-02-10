package subscription_ui_wails

import (
	"context"
	"time"
	billing_ui "vault-app/internal/billing/ui"
	identity_ui "vault-app/internal/identity/ui"
	"vault-app/internal/logger/logger"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"
	onboarding_domain "vault-app/internal/onboarding/domain"
	subscription_application_eventbus "vault-app/internal/subscription/application"
	subscription_usecase "vault-app/internal/subscription/application/usecase"
	subscription_domain "vault-app/internal/subscription/domain"
	subscription_infrastructure "vault-app/internal/subscription/infrastructure"
	subscription_infrastructure_eventbus "vault-app/internal/subscription/infrastructure/eventbus"
	subscription_persistence "vault-app/internal/subscription/infrastructure/persistence"
	"vault-app/internal/tracecore"
	vault_commands "vault-app/internal/vault/application/commands"

	"gorm.io/gorm"
)

type TracecoreInterface interface {
	CreateSubscription(ctx context.Context, id string, email string, plainPassword string) (*subscription_usecase.TraditionalSubscriptionResponse, error)
}


type SubscriptionHandler struct {
	DB *gorm.DB
	CreateListener subscription_usecase.SubscriptionCreatedListener
	CreateUC subscription_usecase.CreateSubscriptionUseCase
	SubscriptionRepository subscription_domain.SubscriptionRepository
	SubscriptionSyncService subscription_infrastructure.SubscriptionSyncService
	ActivatorService subscription_usecase.SubscriptionActivator
	MonitorActivationService subscription_usecase.SubscriptionActivationMonitor
	TracecoreClt *tracecore.TracecoreClient
	StellarService onboarding_usecase.StellarServiceInterface
	logger logger.Logger	
	Bus subscription_application_eventbus.SubscriptionEventBus
	UserSubscriptionRepository subscription_domain.UserRepository
	IdentityHandler *identity_ui.IdentityHandler
	BillingHandler *billing_ui.BillingHandler

	CreateVaultCommand *vault_commands.CreateVaultCommandHandler
	OnboardingUserRepo onboarding_domain.UserRepository
	OnboardingBus onboarding_application_events.OnboardingEventBus
}

func NewSubscriptionHandler(
	db *gorm.DB,
	tracecoreClt *tracecore.TracecoreClient,
	createVaultCommandHandler *vault_commands.CreateVaultCommandHandler,
	stellarService onboarding_usecase.StellarServiceInterface,
	onboardingUserRepo onboarding_domain.UserRepository,
	onboardingBus onboarding_application_events.OnboardingEventBus,
	identityHandler *identity_ui.IdentityHandler,
	logger logger.Logger,
) *SubscriptionHandler {
	subscriptionBus := subscription_infrastructure_eventbus.NewMemoryBus()
	// ===== New: subscription repository (in-memory for now) =====
	subscriptionSubRepo := subscription_persistence.NewSubscriptionRepository(db, &logger)
	userSubscriptionRepo := subscription_persistence.NewUserSubscriptionRepository(db, &logger)

	// ===== New: subscription sync service (cloud integration) =====
	subscriptionService := subscription_infrastructure.NewSubscriptionSyncService(db, tracecoreClt)
	
	// ===== New: core activator (business logic) =====
	// Note: pass a Stellar port implementation if you have one, otherwise nil
	activator := subscription_usecase.NewSubscriptionActivator(
		subscriptionSubRepo, // repo
		subscriptionBus,
		subscriptionService, // vault port (implements ActivationVaultPort)
	)

	// ===== New: listener which only forwards SubscriptionCreated -> activator =====
	createdListener := subscription_usecase.NewSubscriptionCreatedListener(&logger, activator, subscriptionBus)
	
	createSubscriptionUC := subscription_usecase.NewCreateSubscriptionUseCase(subscriptionSubRepo, subscriptionBus, tracecoreClt)
	


	return &SubscriptionHandler{
		CreateUC: *createSubscriptionUC, 
		SubscriptionSyncService: *subscriptionService,
		SubscriptionRepository: subscriptionSubRepo,
		CreateListener: *createdListener,
		ActivatorService: *activator,
		Bus: subscriptionBus,
		UserSubscriptionRepository: userSubscriptionRepo,
		CreateVaultCommand: createVaultCommandHandler,
		StellarService: stellarService,
		IdentityHandler: identityHandler,
		OnboardingUserRepo: onboardingUserRepo,
		OnboardingBus: onboardingBus,
	}
}

func (h *SubscriptionHandler) SetBillingHandler(billingHandler billing_ui.BillingHandler) {

	activationMonitor := subscription_usecase.NewSubscriptionActivationMonitor(
		&h.logger,
		h.Bus,
		h.UserSubscriptionRepository,
		h.SubscriptionRepository,
		h.CreateVaultCommand,
		h.StellarService,
		h.OnboardingUserRepo,
		h.OnboardingBus,
		h.IdentityHandler,
		&billingHandler,
	)
	h.BillingHandler = &billingHandler
	h.MonitorActivationService = *activationMonitor
}

func (h *SubscriptionHandler) CreateSubscription(ctx context.Context, id string, email string, plainPassword string) (*subscription_usecase.TraditionalSubscriptionResponse, error) {
	// create subscription
	response, err := h.CreateUC.Execute(ctx, id, email, plainPassword)
	if err != nil {
		return nil, err
	}

	return &subscription_usecase.TraditionalSubscriptionResponse{
		Subscription: response,
		OccurredAt:   time.Now().Unix(),
	}, nil
}



func (h *SubscriptionHandler) GetSubscription(ctx context.Context, id string) (*subscription_usecase.TraditionalSubscriptionResponse, error) {
	// get subscription
	response, err := h.SubscriptionRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &subscription_usecase.TraditionalSubscriptionResponse{
		Subscription: response,
		OccurredAt:   time.Now().Unix(),
	}, nil
}

func (h *SubscriptionHandler) SaveSubscription(ctx context.Context, s *subscription_domain.Subscription) error {
	// save subscription
	return h.SubscriptionRepository.Save(ctx, s)
}
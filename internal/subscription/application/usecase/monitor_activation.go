// internal/subscriptions/application/usecase/monitor_activation.go
package subscription_usecase

import (
	"context"
	utils "vault-app/internal"
	billing_ui_handlers "vault-app/internal/billing/ui/handlers"
	identity_usecase "vault-app/internal/identity/application/usecase"
	identity_domain "vault-app/internal/identity/domain"
	identity_ui "vault-app/internal/identity/ui"
	"vault-app/internal/logger/logger"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"
	subscription_eventbus "vault-app/internal/subscription/application"
	subscription_domain "vault-app/internal/subscription/domain"
)

// ----------- Interface -----------
type IdentityPortInterface interface {
	RegisterIdentity(ctx context.Context, req identity_usecase.RegisterRequest) (*identity_domain.User, error)
}
type BillingPortInterface interface {
	AddPaymentMethod(ctx context.Context, userID, method, payload string) (string, error)
}
type IdentityHandlerInterface interface {
	Registers(req identity_ui.OnboardRequest) (*identity_domain.User, error)
}
type BillingHandlerInterface interface {
	Onboard(ctx context.Context, req billing_ui_handlers.AddPaymentMethodRequest) (*billing_ui_handlers.AddPaymentMethodResponse, error)
}

// ----------- Struct -----------
type SubscriptionActivationMonitor struct {
	Logger                     *logger.Logger
	Bus                        subscription_eventbus.SubscriptionEventBus
	UserSubscriptionRepository subscription_domain.UserRepository
	SubscriptionRepository     subscription_domain.SubscriptionRepository
	VaultPort                  onboarding_usecase.VaultPort
	StellarService             onboarding_usecase.StellarServiceInterface
	OnBoardingUserRepository   onboarding_usecase.UserServiceInterface
	BusOnboarding              onboarding_application_events.OnboardingEventBus
	IdentityHandler            IdentityHandlerInterface
	BillingHandler             BillingHandlerInterface
}

// ----------- Constructor -----------
func NewSubscriptionActivationMonitor(
	log *logger.Logger,
	bus subscription_eventbus.SubscriptionEventBus,
	userSubRepo subscription_domain.UserRepository,
	subRepo subscription_domain.SubscriptionRepository,
	vaultPort onboarding_usecase.VaultPort,
	stellarService onboarding_usecase.StellarServiceInterface,
	userService onboarding_usecase.UserServiceInterface,
	busOnboarding onboarding_application_events.OnboardingEventBus,
	identityHandler IdentityHandlerInterface,
	billingHandler BillingHandlerInterface,
) *SubscriptionActivationMonitor {
	return &SubscriptionActivationMonitor{
		Logger:                     log,
		Bus:                        bus,
		UserSubscriptionRepository: userSubRepo,
		SubscriptionRepository:     subRepo,
		VaultPort:                  vaultPort,
		StellarService:             stellarService,
		OnBoardingUserRepository:   userService,
		BusOnboarding:              busOnboarding,
		IdentityHandler:            identityHandler,
		BillingHandler:             billingHandler,
	}
}


func (m *SubscriptionActivationMonitor) Listen(ctx context.Context) {
	m.Logger.Info("🛰️ Listening for subscription activations")

	m.Bus.SubscribeToActivation(func(ctx context.Context, event subscription_eventbus.SubscriptionActivated) {
		m.Logger.Info("🚀 Activated subscription=%s for user=%s tier=%s ledger=%d  email=%s password=%s",
			event.SubscriptionID, event.UserID, event.Tier, event.Ledger, event.UserEmail, event.Password)

		// 1. ------------ I.  Register User Subscription ------------
		userSubscription := subscription_domain.UserSubscription{
			ID:    event.SubscriptionID,
			Email: event.UserEmail,
			Role:  "user",
		}
		if err := m.UserSubscriptionRepository.Save(ctx, &userSubscription); err != nil {
			m.Logger.Error("Monitor - Failed to create user subscription: %v", err)
			return
		}
		utils.LogPretty("Monitor - User subscription created:", userSubscription)
		// 2. ------------ update subscription ------------
		subscription, err := m.SubscriptionRepository.GetByID(ctx, event.SubscriptionID)
		if err != nil {
			m.Logger.Error("Monitor - Failed to retrieve subscription: %v", err)
			return
		}
		subscription.UserID = userSubscription.ID
		if err := m.SubscriptionRepository.Update(ctx, subscription); err != nil {
			m.Logger.Error("Monitor - Failed to update subscription: %v", err)
			return
		}
		utils.LogPretty("Monitor - Subscription updated:", subscription)
		m.Logger.Info("Monitor - User subscription retrieved and updated with user ID %s: %v", userSubscription.ID, subscription)

		// 3. ------------ II. Close Onboarding ------------
		onboardingUC := onboarding_usecase.NewOnboardUseCase(
			m.VaultPort,
			m.StellarService,
			m.OnBoardingUserRepository,
			m.BusOnboarding,
			m.Logger,
			m.IdentityHandler,
			m.BillingHandler,
		)
		if _, err := onboardingUC.Execute(ctx, onboarding_usecase.OnboardRequest{
			Identity:             event.UserID,
			Email:                userSubscription.Email,
			IsAnonymous:          false,
			Password:             event.Password,
			Tier:                 event.Tier,
			PaymentMethod:        "",
			EncryptedPaymentData: "",
			SubscriptionID:       event.SubscriptionID,
		}); err != nil {
			m.Logger.Error("Onboarding failed for user=%s: %v", event.UserID, err)
		}

		// 4. ------------ III. Tier side effects (emails, notifications) ------------		
		switch event.Tier {
		case "free":
			m.Logger.Info("🧊 Free tier enabled")
		case "pro":
			m.Logger.Info("🔥 Pro features enabled")
		case "enterprise":
			m.Logger.Warn("🏢 Enterprise tier may need approval")
		default:
			m.Logger.Warn("⚠️ Unknown tier=%s", event.Tier)
		}
		m.Logger.Info("📧 Email queued for user=%s", event.UserID)
		m.Logger.Info("✅ Activation complete for subscription=%s", event.SubscriptionID)
	})

	<-ctx.Done()
	m.Logger.Warn("🛑 SubscriptionActivationMonitor stopped")
}

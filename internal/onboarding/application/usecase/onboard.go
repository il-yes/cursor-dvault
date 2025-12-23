package onboarding_usecase

import (
	"context"
	"errors"
	billing_usecase "vault-app/internal/billing/application/usecase"
	billing_domain "vault-app/internal/billing/domain"
	billing_ui_handlers "vault-app/internal/billing/ui/handlers"
	identity_usecase "vault-app/internal/identity/application/usecase"
	identity_domain "vault-app/internal/identity/domain"
	identity_ui "vault-app/internal/identity/ui"
	"vault-app/internal/logger/logger"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
	vault_commands "vault-app/internal/vault/application/commands"
)

type IdentityHandlerInterface interface {
	Registers(ctx context.Context, req identity_ui.OnboardRequest) (*identity_domain.User, error)
}

type BillingHandlerInterface interface {
	Onboard(ctx context.Context, req billing_ui_handlers.AddPaymentMethodRequest) (*billing_ui_handlers.AddPaymentMethodResponse, error)
}

// This use case orchestrates creating the user profiles, adding payment method, (creating subscription), and creating vault.
// It uses application ports (interfaces) to interact with other bounded contexts.

// ----------- Request -----------
type RegisterRequest struct {
	Email            string
	Password         string
	IsAnonymous      bool
	StellarPublicKey string
}

// ----------- Interfaces -----------
type IdentityRegisterPort interface {
	RegisterIdentity(ctx context.Context, req identity_usecase.RegisterRequest) (*identity_domain.User, error)
}

type BillingPort interface {
	AddPaymentMethod(ctx context.Context, userID string, method billing_domain.PaymentMethod, encryptedPayload string) (*billing_usecase.AddPaymentMethodResponse, error)
}

type SubscriptionPort interface {
	CreateSubscription(ctx context.Context, userID string, tier string) (string, error)
}

type VaultPort interface {
	CreateVault(v vault_commands.CreateVaultCommand) (*vault_commands.CreateVaultResult, error)
}

// ----------- UseCase Request -----------
type OnboardRequest struct {
	Identity             string
	Email                string
	Password             string
	IsAnonymous          bool
	Tier                 string
	PaymentMethod        string
	EncryptedPaymentData string
	SubscriptionID       string
}

// ----------- UseCase Result -----------
type OnboardResult struct {
	UserID         string
	StellarKey     string
	SubscriptionID string
}

type OnboardUseCase struct {
	StellarService  StellarServiceInterface
	OnBoardingUserRepository     UserServiceInterface
	Bus             onboarding_application_events.OnboardingEventBus
	Logger          *logger.Logger
	IdentityHandler IdentityHandlerInterface
	BillingHandler  BillingHandlerInterface
	Vault           VaultPort
}

// ----------- UseCase Constructor -----------
func NewOnboardUseCase(
	v VaultPort,
	stellarService StellarServiceInterface,
	userService UserServiceInterface,
	eventBus onboarding_application_events.OnboardingEventBus,
	logger *logger.Logger,
	identityHandler IdentityHandlerInterface,
	billingHandler BillingHandlerInterface) *OnboardUseCase {
	return &OnboardUseCase{
		Vault:           v,
		StellarService:  stellarService,
		OnBoardingUserRepository:     userService,
		Bus:             eventBus,
		Logger:          logger,
		IdentityHandler: identityHandler,
		BillingHandler:  billingHandler,
	}
}

func (uc *OnboardUseCase) Execute(ctx context.Context, req OnboardRequest) (*OnboardResult, error) {
	if(req.SubscriptionID == "") {
		return nil, errors.New("subscription ID is required")
	}
	if req.IsAnonymous {
		return nil, errors.New("anonymous requested but no stellar public key provided")
	}
	var secretKey string
	var err error

	// 1. ------------- Fetch onboarded account ------------------
	onboardUser, err := uc.OnBoardingUserRepository.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}		
	uc.Logger.Info("onboardUser", onboardUser)

	// 2. ------------- Identity registration ------------------
	userIdentity, err := uc.IdentityHandler.Registers(ctx, identity_ui.OnboardRequest{
		Email:            req.Email,
		Password:         onboardUser.Password,
		IsAnonymous:      req.IsAnonymous,
		StellarPublicKey: onboardUser.StellarPublicKey,
	})
	if err != nil {
		return nil, err
	}

	// 3. ------------- Vault creation ------------------
	result, err := uc.Vault.CreateVault(vault_commands.CreateVaultCommand{
		UserID:    userIdentity.ID,
		VaultName: "Default Vault",
		Password:  req.Password,
	})
	if err != nil {
		return nil, err
	}
	uc.Logger.Info("vault created", result)

	// 4. ------------- Billing registration optional) ------------------
	if req.PaymentMethod != "" {
		if _, err := uc.BillingHandler.Onboard(ctx, billing_ui_handlers.AddPaymentMethodRequest{
			UserID:           userIdentity.ID,
			Method:           req.PaymentMethod,
			EncryptedPayload: req.EncryptedPaymentData,
		}); err != nil {
			return nil, err
		}
	}

	// 5. ------------- (Optional event...) ------------------

	return &OnboardResult{UserID: userIdentity.ID, StellarKey: secretKey, SubscriptionID: req.SubscriptionID}, nil
}



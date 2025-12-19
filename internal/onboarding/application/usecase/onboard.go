package onboarding_usecase

import (
	"context"
	"errors"
	billing_ui_handlers "vault-app/internal/billing/ui/handlers"
	identity_usecase "vault-app/internal/identity/application/usecase"
	identity_domain "vault-app/internal/identity/domain"
	identity_ui "vault-app/internal/identity/ui"
	"vault-app/internal/logger/logger"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
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
	AddPaymentMethod(ctx context.Context, userID string, method string, encryptedPayload string) (string, error)
}

type SubscriptionPort interface {
	CreateSubscription(ctx context.Context, userID string, tier string) (string, error)
}

type VaultPort interface {
	CreateVault(ctx context.Context, userID string) error
}

// ----------- UseCase Request -----------
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

// ----------- UseCase Result -----------
type OnboardResult struct {
	UserID         string
	StellarKey     string
	SubscriptionID string
}

type OnboardUseCase struct {
	StellarService  StellarServiceInterface
	UserService     UserServiceInterface
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
		UserService:     userService,
		Bus:             eventBus,
		Logger:          logger,
		IdentityHandler: identityHandler,
		BillingHandler:  billingHandler,
	}
}

func (uc *OnboardUseCase) Execute(ctx context.Context, req OnboardRequest) (*OnboardResult, error) {
	if req.IsAnonymous && req.StellarPublicKey == "" {
		return nil, errors.New("anonymous requested but no stellar public key provided")
	}
	var secretKey string
	var err error

	// 1. ------------- Onboarding registration ------------------
	accountCreationResponse, err := uc.RegistrationOnboard(ctx, req)
	if err != nil {
		return nil, err
	}
	uc.Logger.Info("accountCreationResponse", accountCreationResponse)

	// 2. ------------- Identity registration ------------------
	userIdentity, err := uc.IdentityHandler.Registers(ctx, identity_ui.OnboardRequest{
		Email:            req.Email,
		Password:         req.Password,
		IsAnonymous:      req.IsAnonymous,
		StellarPublicKey: req.StellarPublicKey,
	})
	if err != nil {
		return nil, err
	}

	// 3. ------------- Vault creation ------------------
	if err := uc.Vault.CreateVault(ctx, userIdentity.ID); err != nil {
		return nil, err
	}

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

func (uc *OnboardUseCase) RegistrationOnboard(ctx context.Context, req OnboardRequest) (*AccountCreationResponse, error) {
	createUC := NewCreateAccountUseCase(uc.StellarService, uc.UserService, uc.Bus, uc.Logger)
	accountCreationResponse, err := createUC.Execute(AccountCreationRequest{
		Email:       req.Email,
		Password:    req.Password,
		IsAnonymous: req.IsAnonymous,
		StellarKey:  req.StellarPublicKey,
	})
	if err != nil {
		return nil, err
	}
	// utils.LogPretty("accountCreationResponse", accountCreationResponse)
	return accountCreationResponse, nil
}

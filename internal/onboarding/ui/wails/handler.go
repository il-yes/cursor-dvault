package onboarding_ui_wails

import (
	identity_domain "vault-app/internal/identity/domain"
	"vault-app/internal/logger/logger"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"
	onboarding_domain "vault-app/internal/onboarding/domain"
	onboarding_infrastructure_eventbus "vault-app/internal/onboarding/infrastructure/eventbus"
	onboarding_persistence "vault-app/internal/onboarding/infrastructure/persistence"
	subscription_domain "vault-app/internal/subscription/domain"
	"vault-app/internal/tracecore"

	"gorm.io/gorm"
)


type OnBoardingHandler struct {
	uc onboarding_usecase.GetRecommendedTierUseCase
	createAccountUseCase onboarding_usecase.CreateAccountUseCase   
    setupPaymentUseCase onboarding_usecase.SetupPaymentAndActivateUseCase
    GetRecommendedTierUseCase onboarding_usecase.GetRecommendedTierUseCase
    DB *gorm.DB
    appLogger *logger.Logger
    Bus onboarding_application_events.OnboardingEventBus
	UserRepo onboarding_domain.UserRepository
}

func  NewOnBoardingHandler(
    stellarService onboarding_usecase.StellarServiceInterface,
    userSubscriptionRepo subscription_domain.UserRepository,
    subscriptionSubRepo subscription_domain.SubscriptionRepository,
    tcClient *tracecore.TracecoreClient,
    DB *gorm.DB,
    appLogger *logger.Logger,
    ) *OnBoardingHandler {

	onboardingUserRepo := onboarding_persistence.NewGormUserRepository(DB)
	getRecommendedTierUC := onboarding_usecase.GetRecommendedTierUseCase{Db: DB}
	onboardingBus := onboarding_infrastructure_eventbus.NewMemoryBus()
	onboardingCreateAccountUC := onboarding_usecase.NewCreateAccountUseCase(stellarService, onboardingUserRepo, onboardingBus, appLogger)

	onboardingSetupPaymentUseCase := onboarding_usecase.NewSetupPaymentAndActivateUseCase(
        onboardingUserRepo, userSubscriptionRepo, subscriptionSubRepo, onboardingBus, *tcClient,
    )

	return &OnBoardingHandler{
        uc: getRecommendedTierUC, 
        createAccountUseCase: *onboardingCreateAccountUC, 
        setupPaymentUseCase: *onboardingSetupPaymentUseCase, 
        GetRecommendedTierUseCase: getRecommendedTierUC,
        DB: DB,
        appLogger: appLogger,
        Bus: onboardingBus,
		UserRepo: onboardingUserRepo,
    }
}


// 0. Get Tier Features
// GetTierFeatures returns feature comparison for pricing page
func (h *OnBoardingHandler) GetTierFeatures() map[string]onboarding_domain.SubscriptionFeatures {
    return map[string]onboarding_domain.SubscriptionFeatures{
        "free": {
            StorageGB:        5,
            CloudBackup:      false,
            MobileApps:       false,
            SharingLimit:     5,
            Support:          "community",
        },
        "pro": {
            StorageGB:        100,
            CloudBackup:      true,
            MobileApps:       true,
            UnlimitedSharing: true,
            Support:          "email_24_48h",
        },
        "pro_plus": {
            StorageGB:         200,
            CloudBackup:       true,
            MobileApps:        true,
            UnlimitedSharing:  true,
            VersionHistory:    true,
            Telemetry:         false,
            AnonymousAccount:  true,
            CryptoPayments:    true,
            EncryptedPayments: true,
            Support:           "encrypted_chat_12h",
        },
        "business": {
            StorageGB:        1024,
            CloudBackup:      true,
            MobileApps:       true,
            UnlimitedSharing: true,
            VersionHistory:   true,
            Telemetry:        false,
            CryptoPayments:   true,
            EncryptedPayments: true,
            APIAccess:        true,
            Tracecore:        true,
            SSO:              true,
            TeamFeatures:     true,
            Support:          "24_7_live",
        },
    }
}

// Step 2: Use Case (conditional based on Step 1)
type UseCaseResponse struct {
    UseCases []string `json:"use_cases"` // ["passwords", "financial", "medical", etc.]
}

// Step 3: Tier Selection
type TierSelectionResponse struct {
    Tier          subscription_domain.SubscriptionTier `json:"tier"`
    PaymentMethod subscription_domain.PaymentMethod    `json:"payment_method"`
}
func (h *OnBoardingHandler) GetRecommendedTier(identityChoice identity_domain.IdentityChoice) subscription_domain.SubscriptionTier {

	return h.uc.Execute(identityChoice)
}

type AccountCreationResponse struct {
    UserID      string `json:"user_id"`
    StellarKey  string `json:"stellar_key,omitempty"` // Generated for anonymous
    SecretKey   string `json:"secret_key,omitempty"`  // CRITICAL: User must save this
}	

// Step 4: Create Account
func (h *OnBoardingHandler) CreateAccount(req onboarding_usecase.AccountCreationRequest) (*AccountCreationResponse, error) {
    response, err := h.createAccountUseCase.Execute(req)
    if err != nil {
        return nil, err
    }

    return &AccountCreationResponse{
        UserID:      response.UserID,
        StellarKey:  response.StellarKey,
        SecretKey:   response.SecretKey,
    }, nil    
}

// Step 5: Setup Payment
func (h *OnBoardingHandler) SetupPaymentAndActivate(req onboarding_usecase.PaymentSetupRequest) (*subscription_domain.Subscription, error) {
    response, err := h.setupPaymentUseCase.Execute(req)
    if err != nil {
        return nil, err
    }

    return response, nil
}

func (h *OnBoardingHandler) FetchUsers() ([]onboarding_domain.User, error) {
    userRepository := onboarding_persistence.NewGormUserRepository(h.DB)
    findUserUC := onboarding_usecase.NewFindUsersUseCase(userRepository)
    return findUserUC.Execute()
}

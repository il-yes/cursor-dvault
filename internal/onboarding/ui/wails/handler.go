package onboarding_ui_wails

import (
	identity_domain "vault-app/internal/identity/domain"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"
	onboarding_domain "vault-app/internal/onboarding/domain"
	subscription_domain "vault-app/internal/subscription/domain"
)


type OnBoardingHandler struct {
	uc *onboarding_usecase.GetRecommendedTierUseCase
	createAccountUseCase *onboarding_usecase.CreateAccountUseCase   
    setupPaymentUseCase *onboarding_usecase.SetupPaymentAndActivateUseCase
}

func NewOnBoardingHandler(
    uc *onboarding_usecase.GetRecommendedTierUseCase, 
    createAccountUseCase *onboarding_usecase.CreateAccountUseCase,
    setupPaymentUseCase *onboarding_usecase.SetupPaymentAndActivateUseCase,
    ) *OnBoardingHandler {
	return &OnBoardingHandler{uc: uc, createAccountUseCase: createAccountUseCase, setupPaymentUseCase: setupPaymentUseCase}
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
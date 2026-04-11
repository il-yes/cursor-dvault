package onboarding_usecase

import (
	"context"
	"fmt"
	"vault-app/internal/onboarding/application/events"
	onboarding_domain "vault-app/internal/onboarding/domain"
	"vault-app/internal/subscription/domain"
	"vault-app/internal/tracecore"
	"vault-app/internal/utils"
)

type FreeSetupUseCase struct {
	UserOnboardingRepo   onboarding_domain.UserRepository
	UserSubscriptionRepo subscription_domain.UserRepository
	SubscriptionRepo     subscription_domain.SubscriptionRepository
	Bus                  onboarding_application_events.OnboardingEventBus
	// CloudAPI         cloud.CloudAPI
	TracecoreClient *tracecore.TracecoreClient
}

func NewFreeSetupUseCase(
	userRepo onboarding_domain.UserRepository,
	userSubscriptionRepo subscription_domain.UserRepository,
	subscriptionRepo subscription_domain.SubscriptionRepository,
	bus onboarding_application_events.OnboardingEventBus,
	tracecoreClient tracecore.TracecoreClient,
) *FreeSetupUseCase {
	return &FreeSetupUseCase{
		UserOnboardingRepo:   userRepo,
		UserSubscriptionRepo: userSubscriptionRepo,
		SubscriptionRepo:     subscriptionRepo,
		Bus:                  bus,
		TracecoreClient:      &tracecoreClient,
	}
}

type FreeSetupRequest struct {
	IsAnonymous bool                                 `json:"is_anonymous"`
	UserID      string                               `json:"user_id"`
	Email       string                               `json:"email"`
	Password    string                               `json:"password"`
	Tier        subscription_domain.SubscriptionTier `json:"tier"`
	Plan        string                               `json:"plan"`
	SessionID   string                               `json:"session_id"`
}

func (a *FreeSetupUseCase) Execute(req FreeSetupRequest) (*tracecore.FreeCheckoutResponse, error) {
	// -----------------------------------------------------
	// 1. create User Identity
	// -----------------------------------------------------

	// -----------------------------------------------------
	// 2. Collect metadata for the cloud
	// -----------------------------------------------------
	cloudReq := tracecore.PaymentSetupRequestBeta{
		Rail:   "standard", // crypto or traditional
		Tier:   string(req.Tier),
		Month:  1,
		UserID: req.UserID,
		Email:  req.Email,

		Plan:         req.Plan,
		PeriodMonths: "1",
		Amount:       "0",
		Currency:     "usd",
		IsAnonymous:  req.IsAnonymous,
		ProductID:    "bronze",
		SessionID:    req.SessionID,
	}
	utils.LogPretty("FreeSetupUseCase - cloudReq", cloudReq)

	// -----------------------------------------------------
	// 3. Call Ankhora Cloud
	// -----------------------------------------------------
	subscriptionRes, err := a.TracecoreClient.FreeCheckout(context.Background(), cloudReq)
	if err != nil {
		return nil, fmt.Errorf("cloud subscription setup failed: %w", err)
	}
	utils.LogPretty("FreeSetupUseCase - subscriptionRes", subscriptionRes)

	// -----------------------------------------------------
	// 6. All good → return
	// -----------------------------------------------------
	return subscriptionRes, nil
}

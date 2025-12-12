package onboarding_usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	// cloud "vault-app/internal/cloud"
	utils "vault-app/internal"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
	onboarding_domain "vault-app/internal/onboarding/domain"
	subscription_domain "vault-app/internal/subscription/domain"
	tracecore "vault-app/internal/tracecore"
)

type SetupPaymentAndActivateUseCase struct {
	// subscriptionService subscription_domain.SubscriptionService
	UserOnboardingRepo         onboarding_domain.UserRepository
	UserSubscriptionRepo subscription_domain.UserRepository
	SubscriptionRepo subscription_domain.SubscriptionRepository
	Bus              onboarding_application_events.OnboardingEventBus
	// CloudAPI         cloud.CloudAPI
	TracecoreClient *tracecore.TracecoreClient
}

func NewSetupPaymentAndActivateUseCase(
	userRepo onboarding_domain.UserRepository,
	userSubscriptionRepo subscription_domain.UserRepository,
	subscriptionRepo subscription_domain.SubscriptionRepository,
	bus onboarding_application_events.OnboardingEventBus,
	tracecoreClient tracecore.TracecoreClient,
) *SetupPaymentAndActivateUseCase {
	return &SetupPaymentAndActivateUseCase{
		UserOnboardingRepo:         userRepo,
		UserSubscriptionRepo: userSubscriptionRepo,
		SubscriptionRepo: subscriptionRepo,
		Bus:              bus,
		TracecoreClient:  &tracecoreClient,
	}
}	

type PaymentSetupRequest struct {
	UserID                string                               `json:"user_id"`
	Tier                  subscription_domain.SubscriptionTier `json:"tier"`
	PaymentMethod         subscription_domain.PaymentMethod    `json:"payment_method"`
	StripePaymentMethodID string                               `json:"stripe_payment_method_id,omitempty"`
	EncryptedPaymentData  string                               `json:"encrypted_payment_data,omitempty"` // Encrypted client-side
	StellarPublicKey      string                               `json:"stellar_public_key,omitempty"`
	CardNumber            string                               `json:"card_number,omitempty"`
	Exp                   string                               `json:"exp,omitempty"`
	CVC                   string                               `json:"cvc,omitempty"`	
}

// SetupPaymentAndActivate handles payment setup and subscription activation (Step 5)
func (a *SetupPaymentAndActivateUseCase) Execute(req PaymentSetupRequest) (*subscription_domain.Subscription, error) {

	// -----------------------------------------------------
	// 1. Retrieve user
	// -----------------------------------------------------
	user, err := a.UserOnboardingRepo.GetByID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// -----------------------------------------------------
	// 2. Collect metadata for the cloud
	// -----------------------------------------------------
	cloudReq := tracecore.PaymentSetupRequest{
		PaymentIntentID: "",
		Rail:            "traditional", // crypto or traditional
		Wallet:          "stripe",
		Month:           1,
		TxHash:          "", // confirmation Stellar payment hash

		UserID: user.ID,
		Email:  user.Email,
		FirstName: user.Email,
		LastName: user.Email,
		LastFour: "4242",
		CardNumber: req.CardNumber,
		Exp: req.Exp,
		CVC: req.CVC,

		StellarPublicKey:      user.StellarPublicKey,
		Tier:                  req.Tier,
		PaymentMethod:         req.PaymentMethod,
		StripePaymentMethodID: req.StripePaymentMethodID,
		EncryptedPaymentData:  req.EncryptedPaymentData,
	}
	utils.LogPretty("cloudReq", cloudReq)

	// -----------------------------------------------------
	// 3. Call Ankhora Cloud
	// -----------------------------------------------------
	subscriptionRes, err := a.TracecoreClient.SetupSubscription(context.Background(), cloudReq)
	if err != nil {
		return nil, fmt.Errorf("cloud subscription setup failed: %w", err)
	}
	utils.LogPretty("subscriptionRes", subscriptionRes)

	// -----------------------------------------------------
	// 4. Update DB with new subscription
	// -----------------------------------------------------
	var subscriptionCloud subscription_domain.Subscription
	if err := json.Unmarshal(subscriptionRes.Data, &subscriptionCloud); err != nil {
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}
	utils.LogPretty("subscriptionCloud", subscriptionCloud)
	subscription := &subscription_domain.Subscription{
		ID:        subscriptionCloud.ID,
		UserID:    user.ID,
		Tier:      string(req.Tier),
		Status:    string(subscription_domain.SubscriptionStatusActive),
		StartedAt: time.Now(),
	}
	utils.LogPretty("subscription", subscription)

	if err := a.SubscriptionRepo.Save(context.Background(), subscription); err != nil {
		return nil, fmt.Errorf("failed to save subscription: %w", err)
	}

	// -----------------------------------------------------
	// 5. Fire subscription activated event
	// -----------------------------------------------------
	evt := onboarding_application_events.SubscriptionActivatedEvent{
		UserID:         user.ID,
		SubscriptionID: subscription.ID,
		Tier:           subscription.Tier,
		OccurredAt:     time.Now(),
	}
	// Todo: remove this context for the app context
	if err := a.Bus.PublishSubscriptionActivated(context.Background(), evt); err != nil {
		return nil, fmt.Errorf("event-bus failure: %w", err)
	}

	// -----------------------------------------------------
	// 6. All good → return
	// -----------------------------------------------------
	return subscription, nil
}

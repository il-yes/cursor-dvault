package onboarding_usecase

import (
	"context"
	"errors"
)

// This use case orchestrates creating a user, adding payment method, creating subscription, and creating vault.
// It uses application ports (interfaces) to interact with other bounded contexts.

// Ports (kept minimal)
type IdentityRegisterPort interface {
	RegisterStandard(ctx context.Context, email, passwordHash string) (string, error) // returns userID
	RegisterAnonymous(ctx context.Context, stellarPublicKey string) (string, string, error) // userID, secretKey
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

// OnboardRequest represents input from UI
type OnboardRequest struct {
	Identity string
	Email string
	Password string
	IsAnonymous bool
	StellarPublicKey string
	Tier string
	PaymentMethod string
	EncryptedPaymentData string
}

// OnboardResult
type OnboardResult struct {
	UserID string
	StellarKey string
	SubscriptionID string
}

type OnboardUseCase struct {
	identity IdentityRegisterPort
	billing  BillingPort
	subscription SubscriptionPort
	vault    VaultPort
}

func NewOnboardUseCase(idp IdentityRegisterPort, b BillingPort, s SubscriptionPort, v VaultPort) *OnboardUseCase {
	return &OnboardUseCase{identity: idp, billing: b, subscription: s, vault: v}
}

func (uc *OnboardUseCase) Execute(ctx context.Context, req OnboardRequest) (*OnboardResult, error) {
	if req.IsAnonymous && req.StellarPublicKey == "" {
		return nil, errors.New("anonymous requested but no stellar public key provided")
	}
	var userID string
	var secretKey string
	var err error
	if req.IsAnonymous {
		userID, secretKey, err = uc.identity.RegisterAnonymous(ctx, req.StellarPublicKey)
		if err != nil { return nil, err }
	} else {
		userID, err = uc.identity.RegisterStandard(ctx, req.Email, req.Password)
		if err != nil { return nil, err }
	}
	// vault creation
	if err := uc.vault.CreateVault(ctx, userID); err != nil { return nil, err }
	// payment method (optional)
	if req.PaymentMethod != "" {
		if _, err := uc.billing.AddPaymentMethod(ctx, userID, req.PaymentMethod, req.EncryptedPaymentData); err != nil { return nil, err }
	}
	// create subscription
	subID, err := uc.subscription.CreateSubscription(ctx, userID, req.Tier)
	if err != nil { return nil, err }
	return &OnboardResult{UserID: userID, StellarKey: secretKey, SubscriptionID: subID}, nil
}
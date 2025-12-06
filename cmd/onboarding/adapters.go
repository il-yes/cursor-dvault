package main

import (
	"context"
	billing_eventbus "vault-app/internal/billing/application"
	billing "vault-app/internal/billing/application/usecase"
	billing_domain "vault-app/internal/billing/domain"
	identity_eventbus "vault-app/internal/identity/application"
	identity "vault-app/internal/identity/application/usecase"
	"vault-app/internal/shared/eventbus"
	subscription_usecase "vault-app/internal/subscription/application/usecase"
	subscription_domain "vault-app/internal/subscription/domain"
)

// Adapters to satisfy OnboardUseCase ports

type regAdapter struct {
	std  *identity.RegisterStandardUserUseCase
	anon *identity.RegisterAnonymousUserUseCase
}

func (r regAdapter) RegisterStandard(ctx context.Context, email, passwordHash string) (string, error) {
	u, err := r.std.Execute(ctx, email, passwordHash)
	if err != nil {
		return "", err
	}
	return u.ID, nil
}

func (r regAdapter) RegisterAnonymous(ctx context.Context, stellarPublicKey string) (string, string, error) {
	u, err := r.anon.Execute(ctx, stellarPublicKey)
	if err != nil {
		return "", "", err
	}
	// No secret key generation here — assumed returned elsewhere
	return u.ID, "", nil
}

type billingAdapter struct {
	uc *billing.AddPaymentMethodUseCase
}

func (b billingAdapter) AddPaymentMethod(ctx context.Context, userID string, method string, encryptedPayload string) (string, error) {
	m, err := b.uc.Execute(ctx, userID, billing_domain.PaymentMethod(method), encryptedPayload)
	if err != nil {
		return "", err
	}
	return m.ID, nil
}

type subscriptionAdapter struct {
	uc *subscription_usecase.CreateSubscriptionUseCase
}

func (s subscriptionAdapter) CreateSubscription(ctx context.Context, userID string, tier string) (string, error) {
	sub, err := s.uc.Execute(ctx, userID, subscription_domain.SubscriptionTier(tier))
	if err != nil {
		return "", err
	}
	return sub.ID, nil
}

// localVault implements VaultPort minimally
type localVault struct{}

func (v *localVault) CreateVault(ctx context.Context, userID string) error { return nil }

// sharedBusWrapper adapts internal/shared/eventbus global functions to the small EventBus interfaces
type sharedBusWrapper struct{}

func (s sharedBusWrapper) PublishUserRegistered(ctx context.Context, e identity_eventbus.UserRegistered) error {
	return eventbus.PublishUserRegistered(ctx, e)
}
func (s sharedBusWrapper) SubscribeToUserRegistered(handler identity_eventbus.UserRegisteredHandler) error {
	return eventbus.SubscribeToUserRegistered(handler)
}
func (s sharedBusWrapper) PublishSubscriptionCreated(ctx context.Context, e subscription_domain.SubscriptionCreated) error {
	return eventbus.PublishSubscriptionCreated(ctx, e)
}
func (s sharedBusWrapper) SubscribeToSubscriptionCreated(handler func(context.Context, subscription_domain.SubscriptionCreated)) error {
	return eventbus.SubscribeToSubscriptionCreated(handler)
}
func (s sharedBusWrapper) PublishPaymentMethodAdded(ctx context.Context, e billing_eventbus.PaymentMethodAdded) error {
	// Convert billing_domain.PaymentMethodAdded to billing_eventbus.PaymentMethodAdded
	evt := billing_eventbus.PaymentMethodAdded{
		InstrumentID: e.InstrumentID,
		UserID:       e.UserID,
		Method:       e.Method,
		OccurredAt:   e.OccurredAt,
	}
	return eventbus.PublishPaymentMethodAdded(ctx, evt)
}
func (s sharedBusWrapper) SubscribeToPaymentMethodAdded(handler billing_eventbus.PaymentMethodAddedHandler) error {
	// Pass the handler directly since it already expects billing_eventbus.PaymentMethodAdded
	return eventbus.SubscribeToPaymentMethodAdded(handler)
}

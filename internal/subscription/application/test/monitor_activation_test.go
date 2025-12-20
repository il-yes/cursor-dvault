// internal/subscription/application/usecase/monitor_activation_test.go
package subscription_usecase_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	billing_ui_handlers "vault-app/internal/billing/ui/handlers"
	"vault-app/internal/blockchain"
	identity_domain "vault-app/internal/identity/domain"
	identity_ui "vault-app/internal/identity/ui"
	"vault-app/internal/logger/logger"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
	onboarding_domain "vault-app/internal/onboarding/domain"
	subscription_eventbus "vault-app/internal/subscription/application"
	subscription_usecase "vault-app/internal/subscription/application/usecase"
	subscription_domain "vault-app/internal/subscription/domain"
	vault_commands "vault-app/internal/vault/application/commands"
)

/* ---------------------------------------------------
   FAKES
---------------------------------------------------*/

type fakeSubscriptionBus struct {
	mu                sync.Mutex
	activatedCallback func(context.Context, subscription_eventbus.SubscriptionActivated)
	createdCallback   func(context.Context, subscription_eventbus.SubscriptionCreated)
}

func (f *fakeSubscriptionBus) PublishActivated(ctx context.Context, event subscription_eventbus.SubscriptionActivated) error {
	f.Emit(event)
	return nil
}

func (f *fakeSubscriptionBus) PublishCreated(ctx context.Context, event subscription_eventbus.SubscriptionCreated) error {
	f.Emit(event)
	return nil
}

func (f *fakeSubscriptionBus) SubscribeToCreation(
	cb func(context.Context, subscription_eventbus.SubscriptionCreated),
) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.createdCallback = cb
	return nil
}

func (f *fakeSubscriptionBus) SubscribeToActivation(
	cb func(context.Context, subscription_eventbus.SubscriptionActivated),
) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.activatedCallback = cb
	return nil
}

func (f *fakeSubscriptionBus) Emit(event interface{}) {
	f.mu.Lock()
	activatedCb := f.activatedCallback
	createdCb := f.createdCallback
	f.mu.Unlock()

	switch e := event.(type) {
	case subscription_eventbus.SubscriptionActivated:
		if activatedCb != nil {
			activatedCb(context.Background(), e)
		}
	case subscription_eventbus.SubscriptionCreated:
		if createdCb != nil {
			createdCb(context.Background(), e)
		}
	}
}

type fakeUserSubscriptionRepo struct {
	sub *subscription_domain.UserSubscription
	err error
}

func (f *fakeUserSubscriptionRepo) FindByUserID(
	ctx context.Context,
	userID string,
) (*subscription_domain.Subscription, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &subscription_domain.Subscription{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}

func (f *fakeUserSubscriptionRepo) FindByEmail(
	ctx context.Context,
	email string,
) (*subscription_domain.UserSubscription, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.sub, nil
}

func (f *fakeUserSubscriptionRepo) Save(
	ctx context.Context,
	sub *subscription_domain.UserSubscription,
) error {
	f.sub = sub
	return nil
}
func (f *fakeUserSubscriptionRepo) Create(*onboarding_domain.User) (*onboarding_domain.User, error) {
	return nil, nil
}

type fakeVault struct {
	called bool
}

func (f *fakeVault) CreateVault(v vault_commands.CreateVaultCommand) (*vault_commands.CreateVaultResult, error) {
	f.called = true
	return &vault_commands.CreateVaultResult{}, nil
}

type fakeIdentityHandler struct {
	called bool
}

func (f *fakeIdentityHandler) Registers(
	ctx context.Context,
	req identity_ui.OnboardRequest,
) (*identity_domain.User, error) {
	f.called = true
	return &identity_domain.User{ID: req.Email}, nil
}

type fakeBillingHandler struct {
	called bool
}

func (f *fakeBillingHandler) Onboard(
	ctx context.Context,
	req billing_ui_handlers.AddPaymentMethodRequest,
) (*billing_ui_handlers.AddPaymentMethodResponse, error) {
	f.called = true
	return &billing_ui_handlers.AddPaymentMethodResponse{}, nil
}

type fakeStellarService struct{}

func (f *fakeStellarService) CreateAccount(plainPassword string) (*blockchain.CreateAccountRes, error) {
	return nil, nil
}
func (f *fakeStellarService) CreateKeypair() (string, string, string, error) {
	return "PUB", "SEC", "TX", nil
}

type fakeUserService struct{}

func (f *fakeUserService) Create(user *onboarding_domain.User) (*onboarding_domain.User, error) {
	return user, nil
}
func (f *fakeUserService) FindByEmail(email string) (*onboarding_domain.User, error) {
	return &onboarding_domain.User{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}

type fakeOnboardingBus struct{}

func (f *fakeOnboardingBus) PublishCreated(ctx context.Context, event onboarding_application_events.AccountCreatedEvent) error {
	return nil
}
func (f *fakeOnboardingBus) SubscribeToAccountCreation(handler func(event onboarding_application_events.AccountCreatedEvent)) error {
	return nil
}
func (f *fakeOnboardingBus) PublishSubscriptionActivated(ctx context.Context, event onboarding_application_events.SubscriptionActivatedEvent) error {
	return nil
}
func (f *fakeOnboardingBus) SubscribeToSubscriptionActivation(handler func(event onboarding_application_events.SubscriptionActivatedEvent)) error {
	return nil
}

type fakeSubRepo struct {
	err error
	sub *subscription_domain.Subscription
}

func (f *fakeSubRepo) FindByUserID(
	ctx context.Context,
	userID string,
) (*subscription_domain.Subscription, error) {
	if f.err != nil {
		return nil, f.err // ✅ THIS is missing
	}

	return &subscription_domain.Subscription{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}
func (f *fakeSubRepo) GetByID(
	ctx context.Context,
	id string,
) (*subscription_domain.Subscription, error) {
	if f.err != nil {
		return nil, f.err
	}

	return &subscription_domain.Subscription{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}
func (f *fakeSubRepo) Save(
	ctx context.Context,
	sub *subscription_domain.Subscription,
) error {
	return nil
}
func (f *fakeSubRepo) Update(
	ctx context.Context,
	sub *subscription_domain.Subscription,
) error {
	return nil
}
type fakeUserSubRepo struct {
	err error
}

func (f *fakeUserSubRepo) FindByUserID(
	ctx context.Context,
	userID string,
) (*subscription_domain.UserSubscription, error) {
	if f.err != nil {
		return nil, f.err // ✅ THIS is missing
	}
	return &subscription_domain.UserSubscription{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}
func (f *fakeUserSubRepo) Save(
	ctx context.Context,
	sub *subscription_domain.UserSubscription,
) error {
	return nil
}
func (f *fakeUserSubRepo) GetByID(
	ctx context.Context,
	id string,
) (*subscription_domain.UserSubscription, error) {
	if f.err != nil {
		return nil, f.err
	}

	return &subscription_domain.UserSubscription{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}
func (f *fakeUserSubRepo) FindByEmail(
	ctx context.Context,
	email string,
) (*subscription_domain.UserSubscription, error) {
	return &subscription_domain.UserSubscription{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}

/* ---------------------------------------------------
   TESTS
---------------------------------------------------*/

func TestSubscriptionActivationMonitor_Success(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := &fakeSubscriptionBus{}
	repo := &fakeSubRepo{
		sub: &subscription_domain.Subscription{
			ID:    "user-123",
			Email: "user@example.com",
		},
	}

	vault := &fakeVault{}
	identityHandler := &fakeIdentityHandler{}
	billingHandler := &fakeBillingHandler{}

	monitor := subscription_usecase.NewSubscriptionActivationMonitor(
		&logger.Logger{},
		bus,
		&fakeUserSubRepo{},
		repo,
		vault,
		&fakeStellarService{},
		&fakeUserService{},
		&fakeOnboardingBus{},
		identityHandler,
		billingHandler,
	)

	go monitor.Listen(ctx)
	time.Sleep(10 * time.Millisecond) // wait for subscription

	bus.Emit(subscription_eventbus.SubscriptionActivated{
		UserID:         "user-123",
		SubscriptionID: "sub-123",
		Tier:           "pro",
		Ledger:         1,
	})

	time.Sleep(100 * time.Millisecond)

	if !vault.called {
		t.Fatalf("expected vault creation to be triggered")
	}
}

func TestSubscriptionActivationMonitor_SubscriptionLookupFails(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := &fakeSubscriptionBus{}
	vault := &fakeVault{}
	subRepo := &fakeSubRepo{
		err: errors.New("db-fail"),
	}

	userSubRepo := &fakeUserSubRepo{}

	monitor := subscription_usecase.NewSubscriptionActivationMonitor(
		&logger.Logger{},
		bus,
		userSubRepo, // ✅ THIS must fail
		subRepo,
		vault,
		&fakeStellarService{},
		&fakeUserService{},
		&fakeOnboardingBus{},
		&fakeIdentityHandler{},
		&fakeBillingHandler{},
	)

	go monitor.Listen(ctx)
	time.Sleep(10 * time.Millisecond) // wait for subscription

	bus.Emit(subscription_eventbus.SubscriptionActivated{
		UserID:         "user-123",
		SubscriptionID: "sub-123",
		Tier:           "pro",
	})

	time.Sleep(50 * time.Millisecond)

	if vault.called {
		t.Fatalf("vault should NOT be created when subscription lookup fails")
	}
}

func TestSubscriptionActivationMonitor_UnknownTierDoesNotFail(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus := &fakeSubscriptionBus{}
	repo := &fakeSubRepo{
		sub: &subscription_domain.Subscription{
			ID:    "user-123",
			Email: "user@example.com",
		},
	}

	vault := &fakeVault{}
	identityHandler := &fakeIdentityHandler{}
	billingHandler := &fakeBillingHandler{}

	monitor := subscription_usecase.NewSubscriptionActivationMonitor(
		&logger.Logger{},
		bus,
		&fakeUserSubRepo{},
		repo,
		vault,
		&fakeStellarService{},
		&fakeUserService{},
		&fakeOnboardingBus{},
		identityHandler,
		billingHandler,
	)

	go monitor.Listen(ctx)
	time.Sleep(10 * time.Millisecond) // wait for subscription

	bus.Emit(subscription_eventbus.SubscriptionActivated{
		UserID:         "user-123",
		SubscriptionID: "sub-123",
		Tier:           "weird-tier",
	})

	time.Sleep(100 * time.Millisecond)

	if !vault.called {
		t.Fatalf("expected onboarding flow to continue even for unknown tier")
	}
}

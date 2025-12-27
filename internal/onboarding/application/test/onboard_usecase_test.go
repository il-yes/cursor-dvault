package onboarding_usecase_test

import (
	"context"
	"errors"
	"testing"

	billing_ui_handlers "vault-app/internal/billing/ui/handlers"
	"vault-app/internal/blockchain"
	identity_domain "vault-app/internal/identity/domain"
	identity_ui "vault-app/internal/identity/ui"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
	onboarding_events "vault-app/internal/onboarding/application/events"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"
	onboarding_domain "vault-app/internal/onboarding/domain"
	vault_commands "vault-app/internal/vault/application/commands"

	"vault-app/internal/logger/logger"
)

type fakeIdentity struct {
	called bool
	err    error
}
func (f *fakeIdentity) Registers(
	req identity_ui.OnboardRequest) (*identity_domain.User, error) {
	f.called = true
	if f.err != nil {
		return nil, f.err
	}
	return &identity_domain.User{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}

type fakeVault struct {
	called bool
	err    error
}

func (f *fakeVault) CreateVault(v vault_commands.CreateVaultCommand) (*vault_commands.CreateVaultResult, error) {
	f.called = true
	return &vault_commands.CreateVaultResult{}, f.err
}

type fakeBilling struct {
	called bool
	err    error
}
func (f *fakeBilling) Onboard(ctx context.Context,	
	req billing_ui_handlers.AddPaymentMethodRequest) (*billing_ui_handlers.AddPaymentMethodResponse, error) {
	f.called = true
	return &billing_ui_handlers.AddPaymentMethodResponse{}, f.err
}	
func (f *fakeBilling) AddPaymentMethod(
	ctx context.Context,
	userID string,
	method string,
	payload string,
) (string, error) {
	f.called = true
	return "pay-123", f.err
}
type fakeUserRepo struct{}

func (f *fakeUserRepo) Create(u *onboarding_domain.User) (*onboarding_domain.User, error) {
	return u, nil
}
func (f *fakeUserRepo) Update(u *onboarding_domain.User) error {
	return nil
}	
func (f *fakeUserRepo) FindByID(id string) (*onboarding_domain.User, error) {
	return &onboarding_domain.User{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}	
func (f *fakeUserRepo) Delete(id string) error {
	return nil
}	
func (f *fakeUserRepo) FindByEmail(email string) (*onboarding_domain.User, error) {
	return &onboarding_domain.User{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}
func (f *fakeUserRepo) GetByID(id string) (*onboarding_domain.User, error) {
	return &onboarding_domain.User{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}
func (f *fakeUserRepo) GetByEmail(email string) (*onboarding_domain.User, error) {
	return &onboarding_domain.User{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}
func (f *fakeUserRepo) List() ([]onboarding_domain.User, error) {
	return []onboarding_domain.User{
		{
			ID:    "user-123",
			Email: "user@example.com",
		},
	}, nil
}

type fakeStellarService struct{}

func (f *fakeStellarService) GenerateKeypair() (string, string, error) {
	return "PUB", "SECRET", nil
}
func (f *fakeStellarService) CreateAccount(stellarPublicKey string)(*blockchain.CreateAccountRes, error)	 {
	return &blockchain.CreateAccountRes{}, nil
}
func (f *fakeStellarService) GetPublicKey() (string, error) {
	return "PUB", nil
}
func (f *fakeStellarService) CreateKeypair() (string, string, string, error) {
	return "PUB", "SECRET", "", nil
}	
type fakeUserService struct{}
func (f *fakeUserService) Create(*onboarding_domain.User) (*onboarding_domain.User, error) {
	return &onboarding_domain.User{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}
func (f *fakeUserService) FindByEmail(email string) (*onboarding_domain.User, error) {
	return &onboarding_domain.User{
		ID:    "user-123",
		Email: "user@example.com",
	}, nil
}	
type fakeBus struct {
	called bool
}

func (f *fakeBus) PublishCreated(ctx context.Context, event onboarding_application_events.AccountCreatedEvent) error {
	f.called = true
	return nil
}
func (f *fakeBus) Publish(event onboarding_events.OnboardingEventBus) {
	f.called = true
}
func (f *fakeBus) PublishSubscriptionActivated(ctx context.Context, event onboarding_application_events.SubscriptionActivatedEvent) error {
	f.called = true
	return nil
}	
func (f *fakeBus) SubscribeToAccountCreation(func(onboarding_application_events.AccountCreatedEvent)) error {
	f.called = true
	return nil
}
func (f *fakeBus) SubscribeToSubscriptionActivated(onboarding_application_events.SubscriptionActivatedEvent) error {
	f.called = true
	return nil
}
func (f *fakeBus) SubscribeToSubscriptionActivation(func(onboarding_application_events.SubscriptionActivatedEvent)) error {
	f.called = true
	return nil
}	

func TestOnboardUseCase_Success(t *testing.T) {
	ctx := context.Background()

	uc := onboarding_usecase.NewOnboardUseCase(
		&fakeVault{},
		&fakeStellarService{},
		&fakeUserService{},
		&fakeBus{},
		&logger.Logger{},
		&fakeIdentity{},
		&fakeBilling{},
	)

	res, err := uc.Execute(ctx, onboarding_usecase.OnboardRequest{
		Email:    "user@example.com",
		Tier:     "pro",
	})

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if res.UserID != "user-123" {
		t.Fatalf("wrong userID")
	}
}

func TestOnboardUseCase_AnonymousMissingKey(t *testing.T) {
	ctx := context.Background()

	uc := onboarding_usecase.NewOnboardUseCase(
		&fakeVault{},
		&fakeStellarService{},
		&fakeUserService{},
		&fakeBus{},
		&logger.Logger{},
		&fakeIdentity{},
		&fakeBilling{},
	)

	_, err := uc.Execute(ctx, onboarding_usecase.OnboardRequest{
		IsAnonymous: true,
	})

	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestOnboardUseCase_IdentityFails(t *testing.T) {
	ctx := context.Background()

	uc := onboarding_usecase.NewOnboardUseCase(
		&fakeVault{},
		&fakeStellarService{},
		&fakeUserService{},
		&fakeBus{},
		&logger.Logger{},
		&fakeIdentity{err: errors.New("identity-fail")},
		&fakeBilling{},
	)

	_, err := uc.Execute(ctx, onboarding_usecase.OnboardRequest{
		Email: "x@x.com",
	})

	if err == nil || err.Error() != "identity-fail" {
		t.Fatalf("expected identity-fail")
	}
}

func TestOnboardUseCase_VaultFails(t *testing.T) {
	ctx := context.Background()

	uc := onboarding_usecase.NewOnboardUseCase(
		&fakeVault{err: errors.New("vault-fail")},
		&fakeStellarService{},
		&fakeUserService{},
		&fakeBus{},
		&logger.Logger{},
		&fakeIdentity{},
		&fakeBilling{},
	)

	_, err := uc.Execute(ctx, onboarding_usecase.OnboardRequest{
		Email: "x@x.com",
	})

	if err == nil || err.Error() != "vault-fail" {
		t.Fatalf("expected vault-fail")
	}
}
func TestOnboardUseCase_BillingFails(t *testing.T) {
	ctx := context.Background()

	uc := onboarding_usecase.NewOnboardUseCase(
		&fakeVault{},
		&fakeStellarService{},
		&fakeUserService{},
		&fakeBus{},
		&logger.Logger{},
		&fakeIdentity{},
		&fakeBilling{err: errors.New("billing-fail")},
	)

	_, err := uc.Execute(ctx, onboarding_usecase.OnboardRequest{
		Email:                "x@x.com",
		PaymentMethod:        "card",
		EncryptedPaymentData: "enc",
	})

	if err == nil || err.Error() != "billing-fail" {
		t.Fatalf("expected billing-fail")
	}
}

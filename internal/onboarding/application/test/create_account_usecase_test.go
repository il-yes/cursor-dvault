package onboarding_usecase_test

import (
	"context"
	"errors"
	"testing"

	"vault-app/internal/blockchain"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"
	onboarding_domain "vault-app/internal/onboarding/domain"
)

// ---------- MOCKS ----------
type MockUserService struct {
	CreateFunc func(u *onboarding_domain.User) (*onboarding_domain.User, error)	
}

func (m *MockUserService) Create(u *onboarding_domain.User) (*onboarding_domain.User, error) {
	return m.CreateFunc(u)
}

type MockStellarService struct {
	CreateKeypairFunc func() (string, string, string, error)
	CreateAccountFunc func(string) (*blockchain.CreateAccountRes, error)
	Called            bool
}

func (m *MockStellarService) CreateKeypair() (string, string, string, error) {
	m.Called = true
	return m.CreateKeypairFunc()
}

func (m *MockStellarService) CreateAccount(pw string) (*blockchain.CreateAccountRes, error) {
	return &blockchain.CreateAccountRes{}, nil
}


type MockBus struct {
	PublishFunc func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error
}

func (m *MockBus) PublishCreated(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
	return m.PublishFunc(ctx, evt)
}

func (m *MockBus) SubscribeToAccountCreation(handler func(event onboarding_application_events.AccountCreatedEvent)) error {
	return nil
}
func (m *MockBus) PublishSubscriptionActivated(ctx context.Context, evt onboarding_application_events.SubscriptionActivatedEvent) error {
	return nil
}
func (m *MockBus) SubscribeToSubscriptionActivation(handler func(event onboarding_application_events.SubscriptionActivatedEvent)) error {
	return nil
}

// ---------- TESTS ----------

// 1. Anonymous Success
func TestCreateAccount_Anonymous_Success(t *testing.T) {
	mockUser := &MockUserService{
		CreateFunc: func(u *onboarding_domain.User) (*onboarding_domain.User, error) {
			u.ID = "anon-001"
			return u, nil
		},
	}

	mockStellar := &MockStellarService{
		CreateKeypairFunc: func() (string, string, string, error) {
			return "GTESTPUB", "STELLARSECRETX", "TX999", nil
		},
		CreateAccountFunc: func(pw string) (*blockchain.CreateAccountRes, error) { return &blockchain.CreateAccountRes{}, nil },
	}

	mockBus := &MockBus{
		PublishFunc: func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
			return nil
		},
	}

	uc := onboarding_usecase.CreateAccountUseCase{
		UserService:    mockUser,
		StellarService: mockStellar,
		Bus:            mockBus,
	}

	req := onboarding_usecase.AccountCreationRequest{
		IsAnonymous: true,
	}

	res, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.UserID == "" {
		t.Fatalf("expected user ID to be set")
	}
	if res.StellarKey != "GTESTPUB" {
		t.Fatalf("unexpected stellar key: %s", res.StellarKey)
	}
	if !mockStellar.Called {
		t.Fatalf("stellar CreateKeypair should be called")
	}
}

// 2. Regular User Success
func TestCreateAccount_Regular_Success(t *testing.T) {
	mockUser := &MockUserService{
		CreateFunc: func(u *onboarding_domain.User) (*onboarding_domain.User, error) {
			u.ID = "user-001"
			return u, nil
		},
	}

	mockStellar := &MockStellarService{
		CreateKeypairFunc: func() (string, string, string, error) {
			t.Fatalf("stellar should NOT be called for regular accounts")
			return "", "", "", nil
		},
		CreateAccountFunc: func(pw string) (*blockchain.CreateAccountRes, error) { return &blockchain.CreateAccountRes{}, nil },	
	}

	mockBus := &MockBus{
		PublishFunc: func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
			return nil
		},
	}

	uc := onboarding_usecase.CreateAccountUseCase{
		UserService:    mockUser,
		StellarService: mockStellar,
		Bus:            mockBus,
	}

	req := onboarding_usecase.AccountCreationRequest{
		Email:       "user@example.com",
		IsAnonymous: false,
	}

	res, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.UserID != "user-001" {
		t.Fatalf("expected correct user ID, got %s", res.UserID)
	}

	if mockStellar.Called {
		t.Fatalf("stellar should NOT be called for regular users")
	}
}

// 3. Stellar Failure
func TestCreateAccount_StellarFailure(t *testing.T) {
	mockUser := &MockUserService{
		CreateFunc: func(u *onboarding_domain.User) (*onboarding_domain.User, error) {
			t.Fatalf("user creation should NOT happen on stellar failure")
			return nil, nil
		},
	}

	mockStellar := &MockStellarService{
		CreateKeypairFunc: func() (string, string, string, error) {
			return "", "", "", errors.New("stellar failed")
		},
		CreateAccountFunc: func(pw string) (*blockchain.CreateAccountRes, error) {
			return nil, errors.New("stellar failed")
		},
	}

	mockBus := &MockBus{
		PublishFunc: func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
			t.Fatalf("event bus should NOT be called on stellar failure")
			return nil
		},
	}

	uc := onboarding_usecase.CreateAccountUseCase{
		UserService:    mockUser,
		StellarService: mockStellar,
		Bus:            mockBus,
	}

	req := onboarding_usecase.AccountCreationRequest{IsAnonymous: true}

	_, err := uc.Execute(req)
	if err == nil {
		t.Fatalf("expected error from stellar service")
	}
}

// 4. UserService Failure
func TestCreateAccount_UserServiceFailure(t *testing.T) {
	mockUser := &MockUserService{
		CreateFunc: func(u *onboarding_domain.User) (*onboarding_domain.User, error) {
			return nil, errors.New("user svc fail")
		},
	}

	mockStellar := &MockStellarService{
		CreateKeypairFunc: func() (string, string, string, error) {
			return "PUB", "SECRET", "TX", nil
		},
		CreateAccountFunc: func(pw string) (*blockchain.CreateAccountRes, error) {
			return &blockchain.CreateAccountRes{}, nil
		},
	}

	mockBus := &MockBus{
		PublishFunc: func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
			t.Fatalf("event bus should NOT run when user creation fails")
			return nil
		},
	}

	uc := onboarding_usecase.CreateAccountUseCase{
		UserService:    mockUser,
		StellarService: mockStellar,
		Bus:            mockBus,
	}

	req := onboarding_usecase.AccountCreationRequest{IsAnonymous: true}

	_, err := uc.Execute(req)
	if err == nil {
		t.Fatalf("expected user service failure")
	}
}

// 5. Event Bus Failure
func TestCreateAccount_EventBusFailure(t *testing.T) {
	mockUser := &MockUserService{
		CreateFunc: func(u *onboarding_domain.User) (*onboarding_domain.User, error) {
			u.ID = "event-user"
			return u, nil
		},
	}

	mockStellar := &MockStellarService{
		CreateKeypairFunc: func() (string, string, string, error) {
			return "PUB", "SECRET", "TX", nil
		},
		CreateAccountFunc: func(pw string) (*blockchain.CreateAccountRes, error) { return &blockchain.CreateAccountRes{}, nil }	,
	}

	mockBus := &MockBus{
		PublishFunc: func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
			return errors.New("bus fail")
		},
	}

	uc := onboarding_usecase.CreateAccountUseCase{
		UserService:    mockUser,
		StellarService: mockStellar,
		Bus:            mockBus,
	}

	req := onboarding_usecase.AccountCreationRequest{IsAnonymous: true}

	_, err := uc.Execute(req)
	if err == nil {
		t.Fatalf("expected event bus failure")
	}
}

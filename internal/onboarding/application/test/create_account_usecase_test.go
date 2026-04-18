package onboarding_usecase_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"vault-app/internal/blockchain"
	"vault-app/internal/logger/logger"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"
	onboarding_domain "vault-app/internal/onboarding/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
)

// ----------------------------------------------------------------------------------------------------
// MOCKS
// ----------------------------------------------------------------------------------------------------
// ============= MockUserService =======================================================
type MockUserService struct {
	FindByEmailFunc func(string) (*onboarding_domain.User, error)
	CreateFunc      func(*onboarding_domain.User) (*onboarding_domain.User, error)
}

func (m *MockUserService) FindByEmail(email string) (*onboarding_domain.User, error) {
	if m.FindByEmailFunc == nil {
		panic("FindByEmailFunc is nil")
	}
	return m.FindByEmailFunc(email)
}

func (m *MockUserService) Create(u *onboarding_domain.User) (*onboarding_domain.User, error) {
	if m.CreateFunc == nil {
		panic("CreateFunc is nil")
	}
	return m.CreateFunc(u)
}

// ============= MockStellarService =======================================================
type MockStellarService struct {
	CreateKeypairFunc func() (string, string, string, error)
	CreateAccountFunc func(string) (*blockchain.CreateAccountRes, error)
	Called            bool
	PubKey            string
}

func (m *MockStellarService) CreateKeypair() (string, string, string, error) {
	if m.CreateKeypairFunc == nil {
		panic("CreateKeypairFunc is nil")
	}
	m.Called = true

	return m.CreateKeypairFunc()
}

func (m *MockStellarService) CreateAccount(pw string) (*blockchain.CreateAccountRes, error) {
	if m.CreateAccountFunc == nil {
		return nil, fmt.Errorf("CreateFunc not implemented")
	}
	return m.CreateAccountFunc(pw)
}

// ============= MockBus =======================================================
type MockBus struct {
	SubscribeToAccountCreationFunc        func(handler func(event onboarding_application_events.AccountCreatedEvent)) error
	PublishFunc                           func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error
	PublishSubscriptionActivatedFunc      func(ctx context.Context, evt onboarding_application_events.SubscriptionActivatedEvent) error
	SubscribeToSubscriptionActivationFunc func(handler func(event onboarding_application_events.SubscriptionActivatedEvent)) error
}

func (m *MockBus) PublishCreated(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
	if m.PublishFunc == nil {
		panic("PublishFunc is nil")
	}
	return m.PublishFunc(ctx, evt)
}

func (m *MockBus) SubscribeToAccountCreation(handler func(event onboarding_application_events.AccountCreatedEvent)) error {
	if m.SubscribeToAccountCreationFunc == nil {
		panic("SubscribeToAccountCreationFunc is nil")
	}
	return m.SubscribeToAccountCreationFunc(handler)
}
func (m *MockBus) PublishSubscriptionActivated(ctx context.Context, evt onboarding_application_events.SubscriptionActivatedEvent) error {
	if m.PublishSubscriptionActivatedFunc == nil {
		panic("PublishSubscriptionActivatedFunc is nil")
	}
	return m.PublishSubscriptionActivatedFunc(ctx, evt)
}
func (m *MockBus) SubscribeToSubscriptionActivation(handler func(event onboarding_application_events.SubscriptionActivatedEvent)) error {
	if m.SubscribeToSubscriptionActivationFunc == nil {
		panic("SubscribeToSubscriptionActivationFunc is nil")
	}
	return m.SubscribeToSubscriptionActivationFunc(handler)
}

// ============= MockKeyEncryption =======================================================
type MockKeyEncryption struct {
	WrapKeyWithPasswordFunc   func([]byte, string) ([]byte, error)
	UnwrapKeyWithPasswordFunc func([]byte, string) ([]byte, error)
	WrapKeyWithStellarFunc    func([]byte, string) ([]byte, error)
	UnwrapKeyWithStellarFunc  func([]byte, string) ([]byte, error)
}

func (m *MockKeyEncryption) WrapKeyWithPassword(vaultKey []byte, password string) ([]byte, error) {
	if m.WrapKeyWithPasswordFunc == nil {
		panic("WrapKeyWithPasswordFunc is nil")
	}
	return m.WrapKeyWithPasswordFunc(vaultKey, password)
}

func (m *MockKeyEncryption) UnwrapKeyWithPassword(enc []byte, password string) ([]byte, error) {
	if m.UnwrapKeyWithPasswordFunc == nil {
		panic("UnwrapKeyWithPasswordFunc is nil")
	}
	return m.UnwrapKeyWithPasswordFunc(enc, password)
}

func (m *MockKeyEncryption) WrapKeyWithStellar(vaultKey []byte, stellarKey string) ([]byte, error) {
	if m.WrapKeyWithStellarFunc == nil {
		panic("WrapKeyWithStellarFunc is nil")
	}
	return m.WrapKeyWithStellarFunc(vaultKey, stellarKey)
}

func (m *MockKeyEncryption) UnwrapKeyWithStellar(enc []byte, stellarKey string) ([]byte, error) {
	if m.UnwrapKeyWithStellarFunc == nil {
		panic("UnwrapKeyWithStellarFunc is nil")
	}
	return m.UnwrapKeyWithStellarFunc(enc, stellarKey)
}
func (m *MockKeyEncryption) Save(vaults_domain.VaultKeyring) error {
	return nil
}

// ============= MockKeyringRepo =======================================================
type MockKeyringRepo struct {
	SaveFunc func(kr vaults_domain.VaultKeyring) error
	Saved    vaults_domain.VaultKeyring
}

func (m *MockKeyringRepo) Save(kr vaults_domain.VaultKeyring) error {
	if m.SaveFunc == nil {
		panic("SaveFunc is nil")
	}
	return m.SaveFunc(kr)
}
func (m *MockKeyringRepo) Get(userID string) (vaults_domain.VaultKeyring, error) {
	return m.Saved, nil
}
func (m *MockKeyringRepo) UnwrapKeyWithPassword(enc []byte, password string) ([]byte, error) {
	return nil, nil
}
func (m *MockKeyringRepo) UnwrapKeyWithStellar(enc []byte, stellarKey string) ([]byte, error) {
	return nil, nil
}
func (m *MockKeyringRepo) WrapKeyWithPassword(vaultKey []byte, password string) ([]byte, error) {
	return nil, nil
}
func (m *MockKeyringRepo) WrapKeyWithStellar(vaultKey []byte, stellarKey string) ([]byte, error) {
	return nil, nil
}

// ============= MockLogger =======================================================
type MockLogger struct {
	InfoFunc  func(msg string, args ...any)
	ErrorFunc func(msg string, args ...any)
}

func (l *MockLogger) Info(msg string, args ...interface{})  {}
func (l *MockLogger) Error(msg string, args ...interface{}) {}

// ============= MockCrypto =======================================================
type mockCrypto struct{}

func (m *mockCrypto) Encrypt(data []byte, key []byte) ([]byte, error) {
	return append(key, data...), nil
}

func (m *mockCrypto) Decrypt(data []byte, key []byte) ([]byte, error) {
	return data[len(key):], nil
}

// ============= MockKeyEnc =======================================================
type MockKeyEnc struct {
	crypto vaults_domain.VaultCrypto
}

func (m *MockKeyEnc) WrapKeyWithPassword(data []byte, password string) ([]byte, error) {
	return m.crypto.Encrypt(data, []byte(password))
}

func (m *MockKeyEnc) UnwrapKeyWithPassword(enc []byte, password string) ([]byte, error) {
	return m.crypto.Decrypt(enc, []byte(password))
}

func (m *MockKeyEnc) WrapKeyWithStellar(data []byte, secret string) ([]byte, error) {
	return m.crypto.Encrypt(data, []byte(secret))
}

func (m *MockKeyEnc) UnwrapKeyWithStellar(enc []byte, secret string) ([]byte, error) {
	return m.crypto.Decrypt(enc, []byte(secret))
}

type MockKeyringService struct{}

func (m *MockKeyringService) SaveHybrid(kr *vaults_domain.VaultKeyring,
	userID string,
	password string,
	stellarSecret string,
) error {
	return errors.New("save failed")
}
// ----------------------------------------------------------------------------------------------------
// TESTS
// ----------------------------------------------------------------------------------------------------

// 1. Anonymous Success
func TestCreateAccount_Anonymous_Success(t *testing.T) {
	// -----------------------------
	// Temp isolated filesystem (BEST PRACTICE)
	// -----------------------------
	tmpDir := t.TempDir()

	// -----------------------------
	// Mocks
	// -----------------------------
	crypto := &mockCrypto{}

	keyEnc := &MockKeyEncryption{
		WrapKeyWithStellarFunc: func(data []byte, secret string) ([]byte, error) {
			return []byte("wrapped-" + string(data)), nil
		},
		WrapKeyWithPasswordFunc: func(data []byte, password string) ([]byte, error) {
			return []byte("wrapped-pw-" + string(data)), nil
		},
	}

	fs := &vault_infrastructure_security.OSFileSystem{}
	keyringService := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmpDir, // ✅ important: directory, not file
		fs,
	)

	logger := &logger.Logger{}

	mockUser := &MockUserService{
		CreateFunc: func(u *onboarding_domain.User) (*onboarding_domain.User, error) {
			u.ID = "anon-001"
			return u, nil
		},
		FindByEmailFunc: func(email string) (*onboarding_domain.User, error) {
			return nil, nil
		},
	}

	mockStellar := &MockStellarService{
		CreateKeypairFunc: func() (string, string, string, error) {
			return "GTESTPUB", "STELLARSECRETX", "TX999", nil
		},
	}

	mockBus := &MockBus{
		PublishFunc: func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
			return nil
		},
	}

	// -----------------------------
	// Use case
	// -----------------------------
	uc := onboarding_usecase.NewCreateAccountUseCase(
		mockStellar,
		mockUser,
		mockBus,
		logger,
		keyringService,
		keyEnc,
	)

	req := onboarding_usecase.AccountCreationRequest{
		IsAnonymous: true,
	}

	// -----------------------------
	// Execute
	// -----------------------------
	res, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// -----------------------------
	// Assertions
	// -----------------------------
	if res.UserID != "anon-001" {
		t.Fatalf("unexpected user ID: %s", res.UserID)
	}

	if res.StellarKey != "GTESTPUB" {
		t.Fatalf("unexpected stellar key: %s", res.StellarKey)
	}

	if res.SecretKey != "STELLARSECRETX" {
		t.Fatalf("unexpected secret key: %s", res.SecretKey)
	}

	if !mockStellar.Called {
		t.Fatalf("stellar CreateKeypair should be called")
	}

	// -----------------------------
	// Verify persistence
	// -----------------------------
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to read keyring dir: %v", err)
	}

	if len(files) == 0 {
		t.Fatalf("expected keyring file to be created")
	}
}

// 2. Regular User Success
func TestCreateAccount_Regular_Success(t *testing.T) {
	mockKeyEnc := &MockKeyEncryption{}

	// -----------------------------
	// Temp isolated filesystem (BEST PRACTICE)
	// -----------------------------
	tmp := t.TempDir()

	crypto := &mockCrypto{}
	keyEnc := &MockKeyEnc{crypto: crypto}

	fs := &vault_infrastructure_security.OSFileSystem{}
	keyringService := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmp,
		fs,
	)
	logger := &logger.Logger{} // IMPORTANT: NOT nil

	mockUser := &MockUserService{
		CreateFunc: func(u *onboarding_domain.User) (*onboarding_domain.User, error) {
			u.ID = "user-001"
			return u, nil
		},
		FindByEmailFunc: func(email string) (*onboarding_domain.User, error) {
			return nil, nil
		},
	}

	mockStellar := &MockStellarService{
		CreateKeypairFunc: func() (string, string, string, error) {
			t.Fatalf("stellar should NOT be called for regular accounts")
			return "", "", "", nil
		},
		CreateAccountFunc: func(pw string) (*blockchain.CreateAccountRes, error) {
			return &blockchain.CreateAccountRes{}, nil
		},
	}

	mockBus := &MockBus{
		PublishFunc: func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
			return nil
		},
	}

	uc := onboarding_usecase.NewCreateAccountUseCase(
		mockStellar,
		mockUser,
		mockBus,
		logger,
		keyringService,
		mockKeyEnc,
	)

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
	mockKeyEnc := &MockKeyEncryption{}

	// -----------------------------
	// Temp isolated filesystem (BEST PRACTICE)
	// -----------------------------
	tmp := t.TempDir()

	crypto := &mockCrypto{}
	keyEnc := &MockKeyEnc{crypto: crypto}

	fs := &vault_infrastructure_security.FailingFS{}
	keyringService := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmp,
		fs,
	)
	logger := &logger.Logger{}
	mockUser := &MockUserService{
		CreateFunc: func(u *onboarding_domain.User) (*onboarding_domain.User, error) {
			t.Fatalf("user creation should NOT happen on stellar failure")
			return nil, nil
		},
		FindByEmailFunc: func(email string) (*onboarding_domain.User, error) {
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

	uc := onboarding_usecase.NewCreateAccountUseCase(
		mockStellar,
		mockUser,
		mockBus,
		logger,
		keyringService,
		mockKeyEnc,
	)

	req := onboarding_usecase.AccountCreationRequest{IsAnonymous: true}

	_, err := uc.Execute(req)
	if err == nil {
		t.Fatalf("expected error from stellar service")
	}
}

// 4. UserService Failure
func TestCreateAccount_UserServiceFailure(t *testing.T) {
	mockKeyEnc := &MockKeyEncryption{}

	// -----------------------------
	// Temp isolated filesystem (BEST PRACTICE)
	// -----------------------------
	tmp := t.TempDir()

	crypto := &mockCrypto{}
	keyEnc := &MockKeyEnc{crypto: crypto}

	fs := &vault_infrastructure_security.OSFileSystem{}
	keyringService := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmp,
		fs,
	)
	logger := &logger.Logger{}
	mockUser := &MockUserService{
		CreateFunc: func(u *onboarding_domain.User) (*onboarding_domain.User, error) {
			return nil, errors.New("user svc fail")
		},
		FindByEmailFunc: func(email string) (*onboarding_domain.User, error) {
			return nil, nil
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

	uc := onboarding_usecase.NewCreateAccountUseCase(
		mockStellar,
		mockUser,
		mockBus,
		logger,
		keyringService,
		mockKeyEnc,
	)

	req := onboarding_usecase.AccountCreationRequest{IsAnonymous: true}

	_, err := uc.Execute(req)
	if err == nil {
		t.Fatalf("expected user service failure")
	}
}

// 5. Event Bus Failure
func TestCreateAccount_EventBusFailure(t *testing.T) {

	// -----------------------------
	// Temp isolated filesystem (BEST PRACTICE)
	// -----------------------------
	tmp := t.TempDir()

	crypto := &mockCrypto{}
	keyEnc := &MockKeyEnc{crypto: crypto}

	fs := &vault_infrastructure_security.OSFileSystem{}
	keyringService := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmp,
		fs,
	)
	logger := &logger.Logger{}
	mockUser := &MockUserService{
		CreateFunc: func(u *onboarding_domain.User) (*onboarding_domain.User, error) {
			u.ID = "event-user"
			return u, nil
		},
		FindByEmailFunc: func(email string) (*onboarding_domain.User, error) {
			return nil, nil
		},
	}
	mockKeyEnc := &MockKeyEncryption{
		WrapKeyWithStellarFunc: func(data []byte, secret string) ([]byte, error) {
			return []byte("wrapped"), nil
		},
		WrapKeyWithPasswordFunc: func(data []byte, password string) ([]byte, error) {
			return []byte("wrapped"), nil
		},
	}
	mockStellar := &MockStellarService{
		CreateKeypairFunc: func() (string, string, string, error) {
			return "PUB", "SECRET", "TX", nil
		},
		CreateAccountFunc: func(pw string) (*blockchain.CreateAccountRes, error) { return &blockchain.CreateAccountRes{}, nil },
	}

	mockBus := &MockBus{
		PublishFunc: func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
			return errors.New("bus fail")
		},
	}

	uc := onboarding_usecase.NewCreateAccountUseCase(
		mockStellar,
		mockUser,
		mockBus,
		logger,
		keyringService,
		mockKeyEnc,
	)

	req := onboarding_usecase.AccountCreationRequest{IsAnonymous: true}

	_, err := uc.Execute(req)
	if err == nil {
		t.Fatalf("expected event bus failure")
	}
}
func TestCreateAccount_Keyring_PasswordWrap(t *testing.T) {

	// -----------------------------
	// Temp isolated filesystem (BEST PRACTICE)
	// -----------------------------
	tmp := t.TempDir()

	crypto := &mockCrypto{}
	keyEnc := &MockKeyEnc{crypto: crypto}

	fs := &vault_infrastructure_security.OSFileSystem{}
	keyringService := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmp,
		fs,
	)
	logger := &logger.Logger{}

	mockUser := &MockUserService{
		FindByEmailFunc: func(email string) (*onboarding_domain.User, error) {
			return nil, nil
		},
		CreateFunc: func(u *onboarding_domain.User) (*onboarding_domain.User, error) {
			u.ID = "user-123"
			return u, nil
		},
	}

	mockStellar := &MockStellarService{
		CreateKeypairFunc: func() (string, string, string, error) {
			return "PUB", "SECRET", "TX", nil
		},
	}

	mockBus := &MockBus{
		PublishFunc: func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
			return nil
		},
	}

	mockKeyEnc := &MockKeyEncryption{
		WrapKeyWithStellarFunc: func(data []byte, secret string) ([]byte, error) {
			return []byte("wrapped"), nil
		},
		WrapKeyWithPasswordFunc: func(data []byte, password string) ([]byte, error) {
			return []byte("wrapped"), nil
		},
	}
	uc := onboarding_usecase.NewCreateAccountUseCase(
		mockStellar,
		mockUser,
		mockBus,
		logger,
		keyringService,
		mockKeyEnc,
	)

	req := onboarding_usecase.AccountCreationRequest{
		Email:    "test@test.com",
		Password: "testpass",
	}

	res, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.UserID == "" {
		t.Fatalf("expected user id")
	}
}

func TestCreateAccount_Keyring_PasswordAndStellarWrap(t *testing.T) {

	tmp := "test.json"
	defer os.Remove(tmp)

	crypto := &mockCrypto{}
	keyEnc := &MockKeyEnc{crypto: crypto}

	fs := &vault_infrastructure_security.OSFileSystem{}
	keyringService := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		tmp,
		fs,
	)

	mockUser := &MockUserService{
		FindByEmailFunc: func(email string) (*onboarding_domain.User, error) {
			return nil, nil
		},
		CreateFunc: func(u *onboarding_domain.User) (*onboarding_domain.User, error) {
			u.ID = "user-456"
			u.StellarPublicKey = "GPUBKEY"
			return u, nil
		},
	}

	mockStellar := &MockStellarService{
		CreateKeypairFunc: func() (string, string, string, error) {
			return "GPUBKEY", "SECRET", "TX123", nil
		},
	}

	mockBus := &MockBus{
		PublishFunc: func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
			return nil
		},
	}

	mockLogger := &MockLogger{
		InfoFunc:  func(msg string, args ...any) {},
		ErrorFunc: func(msg string, args ...any) {},
	}

	mockKeyEnc := &MockKeyEncryption{
		WrapKeyWithPasswordFunc: func(vaultKey []byte, password string) ([]byte, error) {
			return []byte("wrapped-password"), nil
		},
		WrapKeyWithStellarFunc: func(vaultKey []byte, pub string) ([]byte, error) {
			if pub != "GPUBKEY" {
				t.Fatalf("unexpected stellar key")
			}
			return []byte("wrapped-stellar"), nil
		},
	}

	uc := onboarding_usecase.CreateAccountUseCase{
		UserRepo:       mockUser,
		StellarService: mockStellar,
		Bus:            mockBus,
		Logger:         mockLogger,
		KeyringService: keyringService,
		KeyEncryption:  mockKeyEnc,
	}

	req := onboarding_usecase.AccountCreationRequest{
		Email:    "test@test.com",
		Password: "testpass",
	}

	res, err := uc.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res == nil {
		t.Fatalf("expected response")
	}
}

func TestCreateAccount_Keyring_SaveFailure(t *testing.T) {
	mockKeyring := &MockKeyringService{}

	mockUser := &MockUserService{
		FindByEmailFunc: func(email string) (*onboarding_domain.User, error) {
			return nil, nil
		},
		CreateFunc: func(u *onboarding_domain.User) (*onboarding_domain.User, error) {
			u.ID = "user-789"
			return u, nil
		},
	}

	mockStellar := &MockStellarService{
		CreateKeypairFunc: func() (string, string, string, error) {
			return "GPUBKEY", "SECRET", "TX", nil
		},
	}

	mockKeyEnc := &MockKeyEncryption{
		WrapKeyWithPasswordFunc: func(vaultKey []byte, password string) ([]byte, error) {
			return []byte("wrapped"), nil
		},
	}

	mockBus := &MockBus{
		PublishFunc: func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
			return nil
		},
	}

	uc := onboarding_usecase.NewCreateAccountUseCase(
		mockStellar,
		mockUser,
		mockBus,
		&MockLogger{},
		mockKeyring,
		mockKeyEnc,
	)

	req := onboarding_usecase.AccountCreationRequest{
		Email:    "fail@test.com",
		Password: "pass",
	}

	res, err := uc.Execute(req)

	if err == nil {
		t.Fatalf("expected error but got nil")
	}

	if res != nil {
		t.Fatalf("expected nil response on failure")
	}

	if !strings.Contains(err.Error(), "save failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

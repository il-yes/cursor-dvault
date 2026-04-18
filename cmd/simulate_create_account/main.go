package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"vault-app/internal/blockchain"
	"vault-app/internal/logger/logger"
	onboarding_application_events "vault-app/internal/onboarding/application/events"
	onboarding_usecase "vault-app/internal/onboarding/application/usecase"
	onboarding_domain "vault-app/internal/onboarding/domain"
	vaults_domain "vault-app/internal/vault/domain"
	vault_infrastructure_security "vault-app/internal/vault/infrastructure/security"
	vaults_storage "vault-app/internal/vault/infrastructure/storage"
)

// --------------------------------------------------------------------------------------------
// MOCKS
// --------------------------------------------------------------------------------------------
// mockCrypto
type mockCrypto struct{}

func (m *mockCrypto) Encrypt(data []byte, key []byte) ([]byte, error) {
	return append(key, data...), nil
}

func (m *mockCrypto) Decrypt(data []byte, key []byte) ([]byte, error) {
	return data[len(key):], nil
}

// mockKeyEnc
type mockKeyEnc struct {
	crypto vaults_domain.VaultCrypto
}

func (m *mockKeyEnc) WrapKeyWithPassword(data []byte, password string) ([]byte, error) {
	return m.crypto.Encrypt(data, []byte(password))
}

func (m *mockKeyEnc) UnwrapKeyWithPassword(enc []byte, password string) ([]byte, error) {
	return m.crypto.Decrypt(enc, []byte(password))
}

func (m *mockKeyEnc) WrapKeyWithStellar(data []byte, secret string) ([]byte, error) {
	return m.crypto.Encrypt(data, []byte(secret))
}

func (m *mockKeyEnc) UnwrapKeyWithStellar(enc []byte, secret string) ([]byte, error) {
	return m.crypto.Decrypt(enc, []byte(secret))
}

func buildStoredKeyring(passwordEnc, stellarEnc []byte) vaults_storage.StoredKeyring {
	wrappers := []vaults_storage.WrappedKeyring{}

	if passwordEnc != nil {
		wrappers = append(wrappers, vaults_storage.WrappedKeyring{
			Type:       "password",
			Ciphertext: passwordEnc,
		})
	}

	if stellarEnc != nil {
		wrappers = append(wrappers, vaults_storage.WrappedKeyring{
			Type:       "stellar",
			Ciphertext: stellarEnc,
		})
	}

	return vaults_storage.StoredKeyring{
		VaultID:  "vault1",
		Version:  1,
		Wrappers: wrappers,
	}
}
// ---- FAKE / LIGHT DEPENDENCIES ----

// Fake user repo (in-memory)
type InMemoryUserRepo struct {
	users map[string]*onboarding_domain.User
}

func (r *InMemoryUserRepo) Create(u *onboarding_domain.User) (*onboarding_domain.User, error) {
	u.ID = fmt.Sprintf("user-%d", len(r.users)+1)
	r.users[u.ID] = u
	return u, nil
}

func (r *InMemoryUserRepo) FindByEmail(email string) (*onboarding_domain.User, error) {
	return nil, nil
}

// Fake stellar
type FakeStellar struct{
	CreateAccountFunc func(pw string) (*blockchain.CreateAccountRes, error)
}

func (s *FakeStellar) CreateKeypair() (string, string, string, error) {
	return "G_SIM_PUB", "S_SIM_SECRET", "TX_SIM", nil
}

func (m *FakeStellar) CreateAccount(pw string) (*blockchain.CreateAccountRes, error) {
	if m.CreateAccountFunc == nil {
		return nil, fmt.Errorf("CreateFunc not implemented")
	}
	return m.CreateAccountFunc(pw)
}



// Fake bus
type FakeBus struct{
	PublishFunc func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error
	SubscribeToAccountCreationFunc func(handler func(event onboarding_application_events.AccountCreatedEvent)) error
	PublishSubscriptionActivatedFunc func(ctx context.Context, evt onboarding_application_events.SubscriptionActivatedEvent) error
	SubscribeToSubscriptionActivationFunc func(handler func(event onboarding_application_events.SubscriptionActivatedEvent)) error
}

func (m *FakeBus) PublishCreated(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
	if m.PublishFunc == nil {
		panic("PublishFunc is nil")
	}
	return m.PublishFunc(ctx, evt)
}
func (m *FakeBus) SubscribeToAccountCreation(handler func(event onboarding_application_events.AccountCreatedEvent)) error {
	if m.SubscribeToAccountCreationFunc == nil {
		panic("SubscribeToAccountCreationFunc is nil")
	}
	return m.SubscribeToAccountCreationFunc(handler)
}
func (m *FakeBus) PublishSubscriptionActivated(ctx context.Context, evt onboarding_application_events.SubscriptionActivatedEvent) error {
	if m.PublishSubscriptionActivatedFunc == nil {
		panic("PublishSubscriptionActivatedFunc is nil")
	}
	return m.PublishSubscriptionActivatedFunc(ctx, evt)
}
func (m *FakeBus) SubscribeToSubscriptionActivation(handler func(event onboarding_application_events.SubscriptionActivatedEvent)) error {
	if m.SubscribeToSubscriptionActivationFunc == nil {
		panic("SubscribeToSubscriptionActivationFunc is nil")
	}
	return m.SubscribeToSubscriptionActivationFunc(handler)
}



// ---- MAIN ----

func main() {
	fmt.Println("🚀 Simulating CreateAccountUseCase...")

	// temp path (like your tests)
	basePath := "./tmp-keyrings"

	// real crypto (or mock if you want speed)
	crypto := &mockCrypto{}
	keyEnc := &mockKeyEnc{crypto: crypto}

	fs := vault_infrastructure_security.OSFileSystem{}

	keyringService := vault_infrastructure_security.NewKeyringService(
		crypto,
		keyEnc,
		basePath,
		fs,
	)

	userRepo := &InMemoryUserRepo{
		users: make(map[string]*onboarding_domain.User),
	}

	stellar := &FakeStellar{}
	bus := &FakeBus{
		PublishFunc: func(ctx context.Context, evt onboarding_application_events.AccountCreatedEvent) error {
			return nil
		},
	}

	logger := &logger.Logger{}

	uc := onboarding_usecase.NewCreateAccountUseCase(
		stellar,
		userRepo,
		bus,
		logger,
		keyringService,
		keyEnc,
	)

	// ---- RUN ----
	req := onboarding_usecase.AccountCreationRequest{
		IsAnonymous: true,
	}

	res, err := uc.Execute(req)
	if err != nil {
		log.Fatalf("❌ error: %v", err)
	}

	fmt.Println("✅ RESULT:")
	fmt.Printf("UserID: %s\n", res.UserID)
	fmt.Printf("StellarPub: %s\n", res.StellarKey)
	fmt.Printf("Secret: %s\n", res.SecretKey)

	fmt.Println("📁 Check ./tmp-keyrings for output")

	time.Sleep(1 * time.Second)
}
package onboarding_usecase_test

import (
	"context"
	"errors"
	"testing"

	onboard "vault-app/internal/onboarding/application/usecase"
)

/* ---------------------------------------------------
   FAKE PORTS (minimal behavior + call logging)
---------------------------------------------------*/

type fakeIdentity struct {
	registerStandardCalled  bool
	registerAnonymousCalled bool

	returnUserID     string
	returnSecret     string
	returnError      error

	lastEmail        string
	lastPasswordHash string
	lastStellarKey   string
}

func (f *fakeIdentity) RegisterStandard(ctx context.Context, email, passwordHash string) (string, error) {
	f.registerStandardCalled = true
	f.lastEmail = email
	f.lastPasswordHash = passwordHash
	return f.returnUserID, f.returnError
}

func (f *fakeIdentity) RegisterAnonymous(ctx context.Context, stellarPublicKey string) (string, string, error) {
	f.registerAnonymousCalled = true
	f.lastStellarKey = stellarPublicKey
	return f.returnUserID, f.returnSecret, f.returnError
}

type fakeBilling struct {
	called            bool
	userID            string
	method            string
	payload           string
	returnPaymentID   string
	returnError       error
}

func (f *fakeBilling) AddPaymentMethod(ctx context.Context, userID string, method string, encryptedPayload string) (string, error) {
	f.called = true
	f.userID = userID
	f.method = method
	f.payload = encryptedPayload
	return f.returnPaymentID, f.returnError
}

type fakeSubscription struct {
	called      bool
	userID      string
	tier        string
	returnSubID string
	returnError error
}

func (f *fakeSubscription) CreateSubscription(ctx context.Context, userID string, tier string) (string, error) {
	f.called = true
	f.userID = userID
	f.tier = tier
	return f.returnSubID, f.returnError
}

type fakeVault struct {
	called    bool
	userID    string
	returnErr error
}

func (f *fakeVault) CreateVault(ctx context.Context, userID string) error {
	f.called = true
	f.userID = userID
	return f.returnErr
}

/* ---------------------------------------------------
   TESTS
---------------------------------------------------*/

func TestOnboard_StandardUser_Success(t *testing.T) {
	t.Log("➡️ START: TestOnboard_StandardUser_Success")

	ctx := context.Background()

	id := &fakeIdentity{returnUserID: "uid-123"}
	billing := &fakeBilling{returnPaymentID: "pay-001"}
	sub := &fakeSubscription{returnSubID: "sub-001"}
	vault := &fakeVault{}

	uc := onboard.NewOnboardUseCase(id, billing, sub, vault)

	req := onboard.OnboardRequest{
		Email:              "user@example.com",
		Password:           "hashedpw",
		IsAnonymous:        false,
		Tier:               "premium",
		PaymentMethod:      "visa",
		EncryptedPaymentData: "enc123",
	}

	t.Log("📩 Executing use case")
	res, err := uc.Execute(ctx, req)
	if err != nil {
		t.Fatalf("❌ expected success, got %v", err)
	}

	t.Logf("📦 Result: %+v", res)

	// Identity
	if !id.registerStandardCalled {
		t.Errorf("❌ RegisterStandard should have been called")
	}
	if id.lastEmail != "user@example.com" {
		t.Errorf("wrong email passed to identity")
	}

	// Vault
	if !vault.called {
		t.Errorf("❌ vault should have been created")
	}
	if vault.userID != "uid-123" {
		t.Errorf("vault created with wrong userID")
	}

	// Billing
	if !billing.called {
		t.Errorf("❌ payment method should have been added")
	}
	if billing.method != "visa" {
		t.Errorf("wrong payment method")
	}

	// Subscription
	if !sub.called {
		t.Fatalf("❌ subscription should have been created")
	}

	if res.SubscriptionID != "sub-001" {
		t.Errorf("subscription ID wrong")
	}

	t.Log("🎉 PASSED")
}

func TestOnboard_AnonymousUser_Success(t *testing.T) {
	t.Log("➡️ START: TestOnboard_AnonymousUser_Success")

	ctx := context.Background()

	id := &fakeIdentity{
		returnUserID: "anon-123",
		returnSecret: "secret-xyz",
	}
	billing := &fakeBilling{}
	sub := &fakeSubscription{returnSubID: "sub-002"}
	vault := &fakeVault{}

	uc := onboard.NewOnboardUseCase(id, billing, sub, vault)

	req := onboard.OnboardRequest{
		IsAnonymous:       true,
		StellarPublicKey:  "GABC",
		Tier:              "free",
	}

	res, err := uc.Execute(ctx, req)
	if err != nil {
		t.Fatalf("❌ expected success, got %v", err)
	}

	t.Logf("📦 Result: %+v", res)

	if !id.registerAnonymousCalled {
		t.Fatalf("❌ RegisterAnonymous should have been called")
	}

	if res.StellarKey != "secret-xyz" {
		t.Errorf("wrong returned stellar secret")
	}

	t.Log("🎉 PASSED")
}

func TestOnboard_AnonymousMissingStellarKey(t *testing.T) {
	ctx := context.Background()

	uc := onboard.NewOnboardUseCase(&fakeIdentity{}, &fakeBilling{}, &fakeSubscription{}, &fakeVault{})

	_, err := uc.Execute(ctx, onboard.OnboardRequest{
		IsAnonymous: true,
		StellarPublicKey: "",
	})
	if err == nil {
		t.Fatalf("❌ expected error for missing stellar key")
	}
}

func TestOnboard_IdentityFails(t *testing.T) {
	ctx := context.Background()

	id := &fakeIdentity{returnError: errors.New("identity-fail")}
	uc := onboard.NewOnboardUseCase(id, &fakeBilling{}, &fakeSubscription{}, &fakeVault{})

	_, err := uc.Execute(ctx, onboard.OnboardRequest{
		Email: "x@x.com", Password: "pw",
	})

	if err == nil || err.Error() != "identity-fail" {
		t.Fatalf("❌ expected identity-fail, got %v", err)
	}
}

func TestOnboard_VaultFails(t *testing.T) {
	ctx := context.Background()

	id := &fakeIdentity{returnUserID: "uid-1"}
	vault := &fakeVault{returnErr: errors.New("vault-fail")}
	uc := onboard.NewOnboardUseCase(id, &fakeBilling{}, &fakeSubscription{}, vault)

	_, err := uc.Execute(ctx, onboard.OnboardRequest{
		Email: "x@x.com", Password: "pw",
	})

	if err == nil || err.Error() != "vault-fail" {
		t.Fatalf("❌ expected vault-fail, got %v", err)
	}
}

func TestOnboard_PaymentFails(t *testing.T) {
	ctx := context.Background()

	id := &fakeIdentity{returnUserID: "uid-9"}
	billing := &fakeBilling{returnError: errors.New("payment-fail")}
	uc := onboard.NewOnboardUseCase(id, billing, &fakeSubscription{}, &fakeVault{})

	_, err := uc.Execute(ctx, onboard.OnboardRequest{
		Email: "x@x.com", Password: "pw",
		Tier: "pro",
		PaymentMethod: "visa",
		EncryptedPaymentData: "enc",
	})

	if err == nil || err.Error() != "payment-fail" {
		t.Fatalf("❌ expected payment-fail, got %v", err)
	}
}

func TestOnboard_SubscriptionFails(t *testing.T) {
	ctx := context.Background()

	id := &fakeIdentity{returnUserID: "uid-3"}
	sub := &fakeSubscription{returnError: errors.New("sub-fail")}
	uc := onboard.NewOnboardUseCase(id, &fakeBilling{}, sub, &fakeVault{})

	_, err := uc.Execute(ctx, onboard.OnboardRequest{
		Email: "x@x.com", Password: "pw",
		Tier: "premium",
	})

	if err == nil || err.Error() != "sub-fail" {
		t.Fatalf("❌ expected sub-fail, got %v", err)
	}
}

package subscription_usecase_test

import (
	"context"
	"errors"
	"testing"

	subscription_application_eventbus "vault-app/internal/subscription/application"
	sub_uc "vault-app/internal/subscription/application/usecase"
	sub_domain "vault-app/internal/subscription/domain"
	subscription_domain "vault-app/internal/subscription/domain"
)

/* ---------------------------------------------------
   FAKE REPO
---------------------------------------------------*/

type fakeRepo struct {
	saved     *sub_domain.Subscription
	returnErr error
}

func (r *fakeRepo) Save(ctx context.Context, s *sub_domain.Subscription) error {
	if r.returnErr != nil {
		return r.returnErr
	}
	r.saved = s
	return nil
}

func (r *fakeRepo) Update(ctx context.Context, s *sub_domain.Subscription) error {
	if r.returnErr != nil {
		return r.returnErr
	}
	r.saved = s
	return nil
}

func (r *fakeRepo) FindByUserID(ctx context.Context, userID string) (*sub_domain.Subscription, error) {
	if r.returnErr != nil {
		return nil, r.returnErr // ✅ THIS is missing
	}
	return r.saved, nil
}

func (r *fakeRepo) GetByID(ctx context.Context, id string) (*sub_domain.Subscription, error) {
	if r.returnErr != nil {
		return nil, r.returnErr
	}
	return r.saved, nil
}

/* ---------------------------------------------------
   FAKE EVENT BUS
---------------------------------------------------*/

type fakeBus struct {
	called bool
	event  subscription_application_eventbus.SubscriptionCreated
	event2 subscription_application_eventbus.SubscriptionActivated
}

func (b *fakeBus) PublishSubscriptionCreated(ctx context.Context, e subscription_application_eventbus.SubscriptionCreated) error {
	b.called = true
	b.event = e
	return nil
}

func (b *fakeBus) SubscribeToSubscriptionCreated(handler func(context.Context, subscription_application_eventbus.SubscriptionCreated)) error {
	// Not needed for these tests
	return nil
}
func (b *fakeBus) PublishActivated(ctx context.Context, e subscription_application_eventbus.SubscriptionActivated) error {
	b.called = true
	b.event2 = e
	return nil
}
func (b *fakeBus) PublishCreated(ctx context.Context, e subscription_application_eventbus.SubscriptionCreated) error {
	b.called = true
	b.event = e
	return nil
}
func (b *fakeBus) SubscribeToSubscriptionActivated(handler func(context.Context, subscription_application_eventbus.SubscriptionActivated)) error {
	// Not needed for these tests
	return nil
}
func (b *fakeBus) SubscribeToActivation(handler func(context.Context, subscription_application_eventbus.SubscriptionActivated)) error {
	// Not needed for these tests
	return nil
}
func (b *fakeBus) SubscribeToCreation(handler func(context.Context, subscription_application_eventbus.SubscriptionCreated)) error {
	// Not needed for these tests
	return nil
}

/* ---------------------------------------------------
   FAKE ID GENERATOR
---------------------------------------------------*/

func fakeIDGen() string {
	return "sub-123"
}
/* ---------------------------------------------------
   FAKE AnkhoraClient
---------------------------------------------------*/
type fakeTCClient struct {
	sub *subscription_domain.Subscription
}

func (tc *fakeTCClient) GetSubscriptionBySessionID(
	ctx context.Context,
	sessionID string,
) (*subscription_domain.Subscription, error) {

	if tc.sub != nil {
		return tc.sub, nil
	}

	// sensible default
	return NewValidSubscription(), nil
}

func NewValidSubscription() *subscription_domain.Subscription {
	return &subscription_domain.Subscription{
		ID:            "sub-123",
		UserID:        "user-123",
		Email:         "user@example.com",
		Tier:          "pro",
		Months:        12,
		Price:         1999,
		Rail:          "stripe",
	}
}	

func validCreateSubscriptionRequest() CreateSubscriptionRequest {
	return CreateSubscriptionRequest{
		ID:            "sub-123",
		UserID:        "user-123",
		Email:         "user@example.com",
		Password:      "password",
		Tier:          "pro",
		Months:        12,
		Price:         1999,
		Rail:          "stripe",
		PaymentMethod: "card",
	}
}

type CreateSubscriptionRequest struct {
	ID            string
	UserID        string
	Email         string
	Password      string
	Tier          string
	Months        int
	Price         int
	Rail          string
	PaymentMethod string
}
/* ---------------------------------------------------
   TESTS
---------------------------------------------------*/

func TestCreateSubscription_Success(t *testing.T) {
	t.Log("➡️ START: TestCreateSubscription_Success")

	ctx := context.Background()

	repo := &fakeRepo{}
	bus := &fakeBus{}
	tcClient := &fakeTCClient{}

	uc := sub_uc.NewCreateSubscriptionUseCase(repo, bus, tcClient)

	req := validCreateSubscriptionRequest()

	t.Log("📩 Executing use case")
	sub, err := uc.Execute(ctx, req.ID, req.Password)
	if err != nil {
		t.Fatalf("❌ expected no error, got %v", err)
	}

	t.Logf("📦 Result: %+v", sub)

	// Validate repo save
	if repo.saved == nil {
		t.Fatalf("❌ expected subscription saved in repo")
	}
	if repo.saved.ID != "sub-123" {
		t.Errorf("wrong subscription ID, got %s", repo.saved.ID)
	}
	if repo.saved.UserID != req.UserID {
		t.Errorf("wrong userID")
	}
	if repo.saved.Tier != req.Tier {
		t.Errorf("wrong tier")
	}
	if repo.saved.Active {
		t.Errorf("expected Active=false")
	}

	// Validate event bus publish
	if !bus.called {
		t.Fatalf("❌ expected event bus to publish")
	}

	ev := bus.event
	if ev.SubscriptionID != "sub-123" {
		t.Errorf("wrong event SubscriptionID")
	}
	if ev.UserID != req.UserID {
		t.Errorf("wrong event userID")
	}
	if ev.Tier != req.Tier {
		t.Errorf("wrong event tier")
	}

	t.Log("🎉 PASSED")
}

func TestCreateSubscription_SaveFails(t *testing.T) {
	t.Log("➡️ START: TestCreateSubscription_SaveFails")

	ctx := context.Background()

	repo := &fakeRepo{returnErr: errors.New("save-fail")}
	bus := &fakeBus{}
	uc := sub_uc.NewCreateSubscriptionUseCase(repo, bus, &fakeTCClient{})

	req := validCreateSubscriptionRequest()

	_, err := uc.Execute(ctx, req.ID, req.Password)
	if err == nil || err.Error() != "save-fail" {
		t.Fatalf("❌ expected save-fail, got %v", err)
	}

	// Event should NOT be published
	if bus.called {
		t.Fatalf("❌ event bus should NOT be called on failure")
	}

	t.Log("🎉 PASSED")
}

func TestCreateSubscription_NoEventBus(t *testing.T) {
	t.Log("➡️ START: TestCreateSubscription_NoEventBus")

	ctx := context.Background()

	repo := &fakeRepo{}
	uc := sub_uc.NewCreateSubscriptionUseCase(repo, nil, &fakeTCClient{})

	req := validCreateSubscriptionRequest()

	_, err := uc.Execute(ctx, req.ID, req.Password)
	if err != nil {
		t.Fatalf("❌ expected no error, got %v", err)
	}

	if repo.saved == nil {
		t.Fatalf("❌ repo should have saved subscription")
	}

	t.Log("🎉 PASSED")
}

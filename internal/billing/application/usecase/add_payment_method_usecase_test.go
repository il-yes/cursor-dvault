package billing_usecase_test

import (
	"context"
	"errors"
	"testing"

	billing_usecase "vault-app/internal/billing/application/usecase"
	billing_domain "vault-app/internal/billing/domain"
	billing_eventbus "vault-app/internal/billing/application"
)

//
// ─────────────────────────────────────────────────────────────
//  FAKES (Repo + EventBus)
// ─────────────────────────────────────────────────────────────
//

type fakeBillingRepo struct {
	saved       *billing_domain.BillingInstrument
	saveError   error
}
func (b *fakeBillingRepo) FindByUserID(_ context.Context, _ string) ([]*billing_domain.BillingInstrument, error) {
	return []*billing_domain.BillingInstrument{}, nil
}
func (r *fakeBillingRepo) Save(_ context.Context, b *billing_domain.BillingInstrument) error {
	if r.saveError != nil {
		return r.saveError
	}
	r.saved = b
	return nil
}

type fakeEventBus struct {
	received []billing_eventbus.PaymentMethodAdded
}

func (b *fakeEventBus) PublishPaymentMethodAdded(_ context.Context, e billing_eventbus.PaymentMethodAdded) error {
	b.received = append(b.received, e)
	return nil
}

func (b *fakeEventBus) SubscribeToPaymentMethodAdded(_ billing_eventbus.PaymentMethodAddedHandler) error {
	return nil
}

//
// ─────────────────────────────────────────────────────────────
//  TESTS
// ─────────────────────────────────────────────────────────────
//

func TestAddPaymentMethod_Success(t *testing.T) {
	t.Log("➡️ START: TestAddPaymentMethod_Success")

	ctx := context.Background()

	repo := &fakeBillingRepo{}
	bus := &fakeEventBus{}
	idGen := func() string { return "pm-123" }

	uc := billing_usecase.NewAddPaymentMethodUseCase(repo, bus, idGen)

	t.Log("📩 Executing use case")
	result, err := uc.Execute(ctx, "user-111", billing_domain.PaymentCard, "encrypted-abc")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Logf("📦 Result: %+v", result)

	//
	// Assertions
	//
	if result.Instrument.ID != "pm-123" {
		t.Fatalf("expected ID pm-123, got %s", result.Instrument.ID)
	}
	if result.Instrument.UserID != "user-111" {
		t.Fatalf("wrong user: %s", result.Instrument.UserID)
	}
	if result.Instrument.Type != billing_domain.PaymentCard {
		t.Fatalf("wrong payment method: %s", result.Instrument.Type)
	}
	if result.Instrument.EncryptedPayload != "encrypted-abc" {
		t.Fatalf("wrong encrypted payload")
	}

	// repo received save
	if repo.saved == nil {
		t.Fatalf("repo did not save instrument")
	}

	// event bus received event
	if len(bus.received) != 1 {
		t.Fatalf("expected 1 event, got %d", len(bus.received))
	}

	ev := bus.received[0]
	if ev.InstrumentID != "pm-123" || ev.UserID != "user-111" || ev.Method != string(billing_domain.PaymentCard) {
		t.Fatalf("event mismatch: %+v", ev)
	}

	t.Log("🎉 PASSED")
}

func TestAddPaymentMethod_SaveFails(t *testing.T) {
	t.Log("➡️ START: TestAddPaymentMethod_SaveFails")

	ctx := context.Background()
	repo := &fakeBillingRepo{saveError: errors.New("db down")}
	bus := &fakeEventBus{}

	uc := billing_usecase.NewAddPaymentMethodUseCase(repo, bus, func() string { return "pm-x" })

	_, err := uc.Execute(ctx, "user-890", billing_domain.PaymentCard, "payload")
	if err == nil {
		t.Fatalf("expected a failure, got nil")
	}

	// Should not publish event when save fails
	if len(bus.received) != 0 {
		t.Fatalf("event published even though repo failed")
	}

	t.Log("🎉 PASSED")
}

func TestAddPaymentMethod_NoEventBus(t *testing.T) {
	t.Log("➡️ START: TestAddPaymentMethod_NoEventBus")

	ctx := context.Background()

	repo := &fakeBillingRepo{}
	var bus billing_eventbus.EventBus = nil // ⬅ no event bus
	idGen := func() string { return "pm-333" }

	uc := billing_usecase.NewAddPaymentMethodUseCase(repo, bus, idGen)

	result, err := uc.Execute(ctx, "user-555", billing_domain.PaymentCard, "enc-777")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Instrument.ID != "pm-333" {
		t.Fatalf("expected ID pm-333, got %s", result.Instrument.ID)
	}

	if repo.saved == nil {
		t.Fatalf("repo did not save")
	}

	t.Log("🎉 PASSED")
}

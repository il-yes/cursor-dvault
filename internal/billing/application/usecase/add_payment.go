package billing_usecase

import (
	"context"
	billing_eventbus "vault-app/internal/billing/application"
	billing_domain "vault-app/internal/billing/domain"
)

type AddPaymentMethodRequest struct {
	UserID string
	Method billing_domain.PaymentMethod
	EncryptedPayload string
}

type AddPaymentMethodResponse struct {
	Instrument billing_domain.BillingInstrument
}


type AddPaymentMethodUseCase struct {
	repo billing_domain.Repository
}

func NewAddPaymentMethodUseCase(repo billing_domain.Repository) *AddPaymentMethodUseCase {
	return &AddPaymentMethodUseCase{repo: repo}
}

func (uc *AddPaymentMethodUseCase) AddPaymentMethod(
	ctx context.Context, 
	addPaymentMethodRequest AddPaymentMethodRequest,
	idGen func() string,
	bus billing_eventbus.EventBus,
) (*AddPaymentMethodResponse, error) {
	id := idGen()
	b := &billing_domain.BillingInstrument{ID: id, UserID: addPaymentMethodRequest.UserID, Type: addPaymentMethodRequest.Method, EncryptedPayload: addPaymentMethodRequest.EncryptedPayload}
	if err := uc.repo.Save(ctx, b); err != nil {
		return nil, err
	}
	if bus != nil {
		_ = bus.PublishPaymentMethodAdded(ctx, billing_eventbus.PaymentMethodAddedEvent{InstrumentID: id, UserID: addPaymentMethodRequest.UserID, Method: string(addPaymentMethodRequest.Method)})
	}
	return &AddPaymentMethodResponse{Instrument: *b}, nil
}
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
	bus  billing_eventbus.EventBus
	idGen func() string
}

func NewAddPaymentMethodUseCase(repo billing_domain.Repository, bus billing_eventbus.EventBus, idGen func() string) *AddPaymentMethodUseCase {
	return &AddPaymentMethodUseCase{repo: repo, bus: bus, idGen: idGen}
}

func (uc *AddPaymentMethodUseCase) Execute(ctx context.Context, userID string, method billing_domain.PaymentMethod, encryptedPayload string) (*AddPaymentMethodResponse, error) {
	id := uc.idGen()
	b := &billing_domain.BillingInstrument{ID: id, UserID: userID, Type: method, EncryptedPayload: encryptedPayload}
	if err := uc.repo.Save(ctx, b); err != nil {
		return nil, err
	}
	if uc.bus != nil {
		_ = uc.bus.PublishPaymentMethodAdded(ctx, billing_eventbus.PaymentMethodAdded{InstrumentID: id, UserID: userID, Method: string(method)})
	}
	return &AddPaymentMethodResponse{Instrument: *b}, nil
}
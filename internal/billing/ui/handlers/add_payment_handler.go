package billing_ui_handlers

import (
	"context"
	billing_eventbus "vault-app/internal/billing/application"
	billing_usecase "vault-app/internal/billing/application/usecase"
	billing_domain "vault-app/internal/billing/domain"
)


type AddPaymentHandler struct {
	AddPaymentMethodUseCase *billing_usecase.AddPaymentMethodUseCase
	EventBus billing_eventbus.EventBus
	IDGenerator func() string
}

type AddPaymentMethodRequest struct {
	UserID string
	Method string
	EncryptedPayload string
}
type AddPaymentMethodResponse struct {
	Instrument billing_domain.BillingInstrument
}

func NewAddPaymentHandler(
	addPaymentMethodUseCase *billing_usecase.AddPaymentMethodUseCase,
	eventBus billing_eventbus.EventBus,
	idGenerator func() string,
	) *AddPaymentHandler {
	
	return &AddPaymentHandler{AddPaymentMethodUseCase: addPaymentMethodUseCase, EventBus: eventBus, IDGenerator: idGenerator}
}

func (h *AddPaymentHandler) AddPaymentMethod(ctx context.Context, userID string, method string, encryptedPayload string) (*AddPaymentMethodResponse, error) {
	resp, err := h.AddPaymentMethodUseCase.AddPaymentMethod(
		ctx, 
		billing_usecase.AddPaymentMethodRequest{
			UserID: userID,
			Method: billing_domain.PaymentMethod(method),
			EncryptedPayload: encryptedPayload,
		},
		func() string { return h.IDGenerator() },
		h.EventBus,
	)
	if err != nil {
		return nil, err
	}
	return &AddPaymentMethodResponse{Instrument: resp.Instrument}, nil
}

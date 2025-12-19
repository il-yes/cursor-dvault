package billing_ui_handlers

import (
	"context"
	billing_usecase "vault-app/internal/billing/application/usecase"
	billing_domain "vault-app/internal/billing/domain"
)


type AddPaymentHandler struct {
	AddPaymentMethodUseCase *billing_usecase.AddPaymentMethodUseCase
}

type AddPaymentMethodRequest struct {
	UserID string
	Method string
	EncryptedPayload string
}
type AddPaymentMethodResponse struct {
	Instrument billing_domain.BillingInstrument
}

func NewAddPaymentHandler(addPaymentMethodUseCase *billing_usecase.AddPaymentMethodUseCase) *AddPaymentHandler {
	return &AddPaymentHandler{AddPaymentMethodUseCase: addPaymentMethodUseCase}
}

func (h *AddPaymentHandler) AddPaymentMethod(ctx context.Context, req AddPaymentMethodRequest) (*AddPaymentMethodResponse, error) {
	resp, err := h.AddPaymentMethodUseCase.Execute(ctx, req.UserID, billing_domain.PaymentMethod(req.Method), req.EncryptedPayload)
	if err != nil {
		return nil, err
	}
	return &AddPaymentMethodResponse{Instrument: resp.Instrument}, nil
}

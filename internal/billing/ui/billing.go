package billing_ui

import (
	"context"
	billing_ui_handlers "vault-app/internal/billing/ui/handlers"
)

type BillingHandler struct {
	PaymentHandler *billing_ui_handlers.AddPaymentHandler
}

func NewBillingHandler(addPaymentHandler *billing_ui_handlers.AddPaymentHandler) *BillingHandler {
	return &BillingHandler{PaymentHandler: addPaymentHandler}
}

func (h *BillingHandler) Onboard(ctx context.Context, req billing_ui_handlers.AddPaymentMethodRequest) (*billing_ui_handlers.AddPaymentMethodResponse, error) {
	return h.PaymentHandler.AddPaymentMethod(ctx, req)
}	

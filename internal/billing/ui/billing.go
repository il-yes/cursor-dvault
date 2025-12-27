package billing_ui

import (
	"context"
	billing_eventbus "vault-app/internal/billing/application"
	billing_persistence "vault-app/internal/billing/infrastructure/persistence"
	billing_usecase "vault-app/internal/billing/application/usecase"
	billing_ui_handlers "vault-app/internal/billing/ui/handlers"
	"gorm.io/gorm"
)

type BillingHandler struct {
	Db *gorm.DB
	EventBus billing_eventbus.EventBus
	PaymentHandler *billing_ui_handlers.AddPaymentHandler
}

func NewBillingHandler(db *gorm.DB, eventBus billing_eventbus.EventBus) *BillingHandler {
	return &BillingHandler{Db: db, EventBus: eventBus}
}

func (h *BillingHandler) Onboard(ctx context.Context, req billing_ui_handlers.AddPaymentMethodRequest) (*billing_ui_handlers.AddPaymentMethodResponse, error) {
	billingRepo := billing_persistence.NewGormBillingRepository(h.Db)
	addPaymentMethodUseCase := billing_usecase.NewAddPaymentMethodUseCase(billingRepo)
	paymentHandler := billing_ui_handlers.NewAddPaymentHandler(addPaymentMethodUseCase, h.EventBus, func() string { return "pm-x" })

	return paymentHandler.AddPaymentMethod(ctx, req.UserID, req.Method, req.EncryptedPayload)
}	

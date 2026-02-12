package billing_ui

import (
	"context"
	billing_eventbus "vault-app/internal/billing/application"
	billing_usecase "vault-app/internal/billing/application/usecase"
	billing_domain "vault-app/internal/billing/domain"
	billing_infrastructure_eventbus "vault-app/internal/billing/infrastructure/eventbus"
	billing_persistence "vault-app/internal/billing/infrastructure/persistence"
	billing_ui_handlers "vault-app/internal/billing/ui/handlers"

	"gorm.io/gorm"
)

type BillingHandler struct {
	Db *gorm.DB
	EventBus billing_eventbus.EventBus
	PaymentHandler *billing_ui_handlers.AddPaymentHandler
	BillingAppUC *billing_usecase.BillingApp
	SubscriptionService billing_usecase.SubscriptionServiceInterface
	AnkhoraService billing_usecase.AnkhoraCloudServiceInterface
	Bus billing_eventbus.EventBus
}

func NewBillingHandler(
	db *gorm.DB, 
	sub billing_usecase.SubscriptionServiceInterface, 
	ankh billing_usecase.AnkhoraCloudServiceInterface,
) *BillingHandler {

	billingBus := billing_infrastructure_eventbus.NewMemoryBus()
	
	return &BillingHandler{
		Db: db, 
		EventBus: billingBus, 
		SubscriptionService: sub, 
		AnkhoraService: ankh,
		Bus: billingBus,
	}
}

func (h *BillingHandler) Onboard(ctx context.Context, req billing_ui_handlers.AddPaymentMethodRequest) (*billing_ui_handlers.AddPaymentMethodResponse, error) {
	billingRepo := billing_persistence.NewGormBillingRepository(h.Db)
	addPaymentMethodUseCase := billing_usecase.NewAddPaymentMethodUseCase(billingRepo)
	paymentHandler := billing_ui_handlers.NewAddPaymentHandler(addPaymentMethodUseCase, h.EventBus, func() string { return "pm-x" })

	return paymentHandler.AddPaymentMethod(ctx, req.UserID, req.Method, req.EncryptedPayload)
}


func (h *BillingHandler) GetPendingPaymentRequestsByUserID(ctx context.Context, userID string) ([]*billing_domain.PaymentRequest, error) {
	billingAppH := billing_ui_handlers.NewBillingAppHandler(ctx, h.SubscriptionService, h.AnkhoraService)
	return  billingAppH.GetPendingPaymentRequests(userID)

}


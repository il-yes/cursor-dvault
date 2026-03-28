package billing_ui

import (
	"context"
	billing_eventbus "vault-app/internal/billing/application"
	billing_usecase "vault-app/internal/billing/application/usecase"
	billing_domain "vault-app/internal/billing/domain"
	billing_infrastructure_eventbus "vault-app/internal/billing/infrastructure/eventbus"
	billing_persistence "vault-app/internal/billing/infrastructure/persistence"
	billing_ui_handlers "vault-app/internal/billing/ui/handlers"
	tracecore_types "vault-app/internal/tracecore/types"
	"vault-app/internal/utils"

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
	utils.LogPretty("BillingHandler - Onboard - req", req)
	billingRepo := billing_persistence.NewGormBillingRepository(h.Db)
	utils.LogPretty("BillingHandler - billingRepo", billingRepo)
	addPaymentMethodUseCase := billing_usecase.NewAddPaymentMethodUseCase(billingRepo)
	utils.LogPretty("BillingHandler - addPaymentMethodUseCase", addPaymentMethodUseCase)
	paymentHandler := billing_ui_handlers.NewAddPaymentHandler(addPaymentMethodUseCase, h.EventBus, func() string { return "pm-x" })
	utils.LogPretty("BillingHandler - paymentHandler", paymentHandler)

	return paymentHandler.AddPaymentMethod(ctx, req.UserID, req.Method, req.EncryptedPayload)
}


func (h *BillingHandler) GetPendingPaymentRequestsByUserID(ctx context.Context, userID string) ([]*billing_domain.PaymentRequest, error) {
	billingAppH := billing_ui_handlers.NewBillingAppHandler(ctx, h.SubscriptionService, h.AnkhoraService)
	return  billingAppH.GetPendingPaymentRequests(userID)

}

func (h *BillingHandler) GetPaymentHistory(ctx context.Context, subID string, limit int) (*tracecore_types.CloudResponse[[]tracecore_types.PaymentHistory], error) {
	billingAppH := billing_ui_handlers.NewBillingAppHandler(ctx, h.SubscriptionService, h.AnkhoraService)
	return  billingAppH.GetBillingHistory(subID, limit)

}


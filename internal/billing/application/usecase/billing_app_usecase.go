package billing_usecase

import (
	"context"
	billing_domain "vault-app/internal/billing/domain"
	"vault-app/internal/subscription/domain"
	"vault-app/internal/tracecore"
)


// --------- Interfaces ---------
type SubscriptionServiceInterface interface {
	FetchFromCloud(ctx context.Context, userID string) (*subscription_domain.Subscription, error)
	GetCachedSubscription(ctx context.Context, userID string) (*subscription_domain.Subscription, error)
	SyncLocal(ctx context.Context, subscription *subscription_domain.Subscription) error
	CancelSubscription(ctx context.Context, userID string, reason string) error
	UpdatePaymentMethod(ctx context.Context, req *billing_domain.UpdatePaymentMethodRequest) error
	GetStorageUsage(ctx context.Context, userID string) (*billing_domain.StorageUsage, error)
	HandleUpgrade(ctx context.Context, userID string, newTier subscription_domain.SubscriptionTier, paymentMethod subscription_domain.PaymentMethod) error
	ReactivateSubscription(ctx context.Context, userID string, tier subscription_domain.SubscriptionTier, paymentMethod subscription_domain.PaymentMethod) error		
}

type AnkhoraCloudServiceInterface interface {
	GetPendingPaymentRequests(ctx context.Context, userID string) ([]*billing_domain.PaymentRequest, error)
	ProcessEncryptedPayment(ctx context.Context, req *billing_domain.ClientPaymentRequest) (*tracecore.ClientPaymentResponse, error)
	HandleClientInitiatedPayment(ctx context.Context, req *billing_domain.ClientPaymentRequest) (*tracecore.ClientPaymentResponse, error)
	GetBillingHistory(ctx context.Context, userID string, limit int) ([]*billing_domain.PaymentHistory, error)
	GenerateReceipt(ctx context.Context, userID string, paymentID string) (*billing_domain.Receipt, error)		
}

// --------- Billing App UseCase ---------
type BillingApp struct {
    ctx             context.Context
    subscriptionService SubscriptionServiceInterface
	AnkhoraCloudService AnkhoraCloudServiceInterface
}

func NewBillingApp(
	ctx context.Context, 
	subscriptionService SubscriptionServiceInterface,
	AnkhoraCloudService AnkhoraCloudServiceInterface,
) *BillingApp {
    return &BillingApp{
		ctx: ctx, 
		subscriptionService: subscriptionService,
		AnkhoraCloudService: AnkhoraCloudService,
	}
}

// GetPendingPaymentRequests returns all pending payment requests for current user
func (a *BillingApp) GetPendingPaymentRequests(userID string) ([]*billing_domain.PaymentRequest, error) {
    payments, err := a.AnkhoraCloudService.GetPendingPaymentRequests(a.ctx, userID)
	if err != nil {
		return nil, err
	}
	return payments, nil
}
	
// ProcessEncryptedPayment processes payment using decrypted card data
func (a *BillingApp) ProcessEncryptedPayment(req *billing_domain.ClientPaymentRequest) (*tracecore.ClientPaymentResponse, error) {
	response, err := a.AnkhoraCloudService.ProcessEncryptedPayment(a.ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// GetSubscriptionDetails returns current subscription details
func (a *BillingApp) GetSubscriptionDetails(userID string) (*subscription_domain.Subscription, error) {
    // 1. Always fetch cloud truth
    cloudSub, err := a.subscriptionService.FetchFromCloud(a.ctx, userID)
    if err != nil {
        // 2. Fallback to local cache
        return a.subscriptionService.GetCachedSubscription(a.ctx, userID)
    }

    // 3. Sync local store
    _ = a.subscriptionService.SyncLocal(a.ctx, cloudSub)

    return cloudSub, nil

}
func (a *BillingApp) HandleClientInitiatedPayment(req *billing_domain.ClientPaymentRequest) (*tracecore.ClientPaymentResponse, error) {
	response, err := a.AnkhoraCloudService.HandleClientInitiatedPayment(a.ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (a *BillingApp) GetBillingHistory(userID string, limit int) ([]*billing_domain.PaymentHistory, error) {
	response, err := a.AnkhoraCloudService.GetBillingHistory(a.ctx, userID, limit)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (a *BillingApp) GenerateReceipt(userID string, paymentID string) (*billing_domain.Receipt, error) {
	return a.AnkhoraCloudService.GenerateReceipt(a.ctx, userID, paymentID)
}

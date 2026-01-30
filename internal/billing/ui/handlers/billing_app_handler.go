package billing_ui_handlers

import (
	billing_usecase "vault-app/internal/billing/application/usecase"
	"vault-app/internal/tracecore"

	"context"
	billing_domain "vault-app/internal/billing/domain"
	subscription_domain "vault-app/internal/subscription/domain"
)

// --------- Billing Handler ---------
type BillingAppHandler struct {
	ctx context.Context
	BillingAppUC *billing_usecase.BillingApp
    subscriptionService billing_usecase.SubscriptionServiceInterface
	ankhoraService billing_usecase.AnkhoraCloudServiceInterface
}

func NewBillingAppHandler(
	ctx context.Context, 
	sub billing_usecase.SubscriptionServiceInterface, 
	ankh billing_usecase.AnkhoraCloudServiceInterface,
) *BillingAppHandler {
	bc := billing_usecase.NewBillingApp(ctx, sub, ankh)
	
	return &BillingAppHandler{
		ctx: ctx, 
		BillingAppUC: bc,
		subscriptionService: sub,
		ankhoraService: ankh,
	}
}	 

// GetPendingPaymentRequests returns all pending payment requests for current user
func (a *BillingAppHandler) GetPendingPaymentRequests(userID string) ([]*billing_domain.PaymentRequest, error) {
    return a.BillingAppUC.GetPendingPaymentRequests(userID)
}

// ProcessEncryptedPayment processes payment using decrypted card data
func (a *BillingAppHandler) ProcessEncryptedPayment(userID string, req *billing_domain.ClientPaymentRequest) (*tracecore.ClientPaymentResponse, error) {
    return a.BillingAppUC.HandleClientInitiatedPayment(req)
}

// GetSubscriptionDetails returns current subscription details
func (a *BillingAppHandler) GetSubscriptionDetails(userID string) (*subscription_domain.Subscription, error) {
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

// CancelSubscription cancels current subscription
func (a *BillingAppHandler) CancelSubscription(userID string, reason string) error {
    return a.subscriptionService.CancelSubscription(a.ctx, userID, reason)
}

// UpdatePaymentMethod updates payment method for subscription
func (a *BillingAppHandler) UpdatePaymentMethod(userID string, req *billing_domain.UpdatePaymentMethodRequest) error {
    req.UserID = userID
    return a.subscriptionService.UpdatePaymentMethod(a.ctx, req)
}

// GetBillingHistory returns payment history
func (a *BillingAppHandler) GetBillingHistory(userID string, limit int) ([]*billing_domain.PaymentHistory, error) {
    return a.BillingAppUC.GetBillingHistory(userID, limit)
}

// DownloadReceipt downloads blockchain-verified receipt
func (a *BillingAppHandler) DownloadReceipt(userID string, paymentID string) (*billing_domain.Receipt, error) {
    return a.BillingAppUC.GenerateReceipt(userID, paymentID)
}

// GetStorageUsage returns current storage usage vs quota
func (a *BillingAppHandler) GetStorageUsage(userID string) (*billing_domain.StorageUsage, error) {
    return a.subscriptionService.GetStorageUsage(a.ctx, userID)
}

// UpgradeSubscription upgrades to a higher tier
func (a *BillingAppHandler) UpgradeSubscription(userID string, req *billing_domain.UpgradeRequest) error {
    req.UserID = userID
    return a.subscriptionService.HandleUpgrade(a.ctx, req.UserID, req.NewTier, req.PaymentMethod)
}

// ReactivateSubscription reactivates a cancelled subscription
func (a *BillingAppHandler) ReactivateSubscription(userID string, tier subscription_domain.SubscriptionTier, paymentMethod subscription_domain.PaymentMethod) error {
    return a.subscriptionService.ReactivateSubscription(a.ctx, userID, tier, paymentMethod)
}


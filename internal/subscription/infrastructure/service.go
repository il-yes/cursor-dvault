package subscription_infrastructure

import (
	"context"
	"time"
	utils "vault-app/internal"
	billing_domain "vault-app/internal/billing/domain"
	subscription_domain "vault-app/internal/subscription/domain"
	"vault-app/internal/tracecore"

	"gorm.io/gorm"
)

type SubscriptionSyncService struct {
	Db            *gorm.DB
	AnkhoraClient *tracecore.TracecoreClient
}

func NewSubscriptionSyncService(db *gorm.DB, ankhoraClient *tracecore.TracecoreClient) *SubscriptionSyncService {
	return &SubscriptionSyncService{Db: db, AnkhoraClient: ankhoraClient}
}


// Fetch subscription from cloud
func (s *SubscriptionSyncService) FetchFromCloud(ctx context.Context, userID string) (*subscription_domain.Subscription, error) {
	response, err := s.AnkhoraClient.GetSubscriptionByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return response, nil
}
// Fetch subscription from local cache Database
func (s *SubscriptionSyncService) GetCachedSubscription(ctx context.Context, userID string) (*subscription_domain.Subscription, error) {
	var subscription subscription_domain.Subscription
	if err := s.Db.Where("user_id = ?", userID).First(&subscription).Error; err != nil {
		return nil, err
	}
	return &subscription, nil
}
// Sync local subscription with cloud subscription
func (s *SubscriptionSyncService) SyncLocal(ctx context.Context, cloud *subscription_domain.Subscription) error {
    var local subscription_domain.Subscription
    err := s.Db.Where("user_id = ?", cloud.UserID).First(&local).Error

    if err == gorm.ErrRecordNotFound {
        return s.Db.Create(cloud).Error
    }

    if cloud.UpdatedAt.After(local.UpdatedAt) {
        return s.Db.Save(cloud).Error
    }

    return nil
}

// Cancel subscription
func (s *SubscriptionSyncService) CancelSubscription(ctx context.Context, userID string, reason string) error {
	// call cloud cancel subscription
	err := s.AnkhoraClient.CancelSubscription(ctx, userID, reason)
	if err != nil {
		return err
	}
	utils.LogPretty("CancelSubscription - cloud response", err)

	// update local subscription
	var subscription subscription_domain.Subscription
	if err := s.Db.Where("user_id = ?", userID).First(&subscription).Error; err != nil {
		return err
	}
	subscription.Status = string(subscription_domain.SubscriptionStatusCancelled)
	subscription.CancelledAt = time.Now()
	subscription.CancelReason = reason
	if err := s.Db.Save(&subscription).Error; err != nil {
		return err
	}
	return nil
}

// Update payment method
func (s *SubscriptionSyncService) UpdatePaymentMethod(ctx context.Context, req *billing_domain.UpdatePaymentMethodRequest) error {
	// call cloud update payment method
	err := s.AnkhoraClient.UpdatePaymentMethod(ctx, req)
	if err != nil {
		return err
	}
	utils.LogPretty("UpdatePaymentMethod - cloud response", err)

	// update local subscription
	var subscription subscription_domain.Subscription
	if err := s.Db.Where("user_id = ?", req.UserID).First(&subscription).Error; err != nil {
		return err
	}
	subscription.PaymentMethod = subscription_domain.PaymentMethod(req.PaymentMethod)
	if err := s.Db.Save(&subscription).Error; err != nil {
		return err
	}
	return nil
}

// Get storage usage from cloud
func (s *SubscriptionSyncService) GetStorageUsage(ctx context.Context, userID string) (*billing_domain.StorageUsage, error) {
	// call cloud get storage usage
	response, err := s.AnkhoraClient.GetStorageUsage(ctx, userID)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("GetStorageUsage - cloud response", response)
	return response, nil
}

// Handle upgrade subscription
func (s *SubscriptionSyncService) HandleUpgrade(ctx context.Context, userID string, newTier subscription_domain.SubscriptionTier, paymentMethod subscription_domain.PaymentMethod) error {
	// call cloud handle upgrade subscription
	err := s.AnkhoraClient.HandleUpgrade(ctx, userID, newTier, paymentMethod)
	if err != nil {
		return err
	}
	utils.LogPretty("HandleUpgrade - cloud response", err)

	// update local subscription
	var subscription subscription_domain.Subscription
	if err := s.Db.Where("user_id = ?", userID).First(&subscription).Error; err != nil {
		return err
	}
	subscription.Tier = string(newTier)
	subscription.PaymentMethod = subscription_domain.PaymentMethod(paymentMethod)
	if err := s.Db.Save(&subscription).Error; err != nil {
		return err
	}
	return nil
}

// Reactivate subscription
func (s *SubscriptionSyncService) ReactivateSubscription(ctx context.Context, userID string, tier subscription_domain.SubscriptionTier, paymentMethod subscription_domain.PaymentMethod) error {
	// call cloud reactivate subscription
	err := s.AnkhoraClient.ReactivateSubscription(ctx, userID, tier, paymentMethod)
	if err != nil {
		return err
	}
	utils.LogPretty("ReactivateSubscription - cloud response", err)

	// update local subscription
	var subscription subscription_domain.Subscription
	if err := s.Db.Where("user_id = ?", userID).First(&subscription).Error; err != nil {
		return err
	}
	subscription.Tier = string(tier)
	subscription.PaymentMethod = subscription_domain.PaymentMethod(paymentMethod)
	if err := s.Db.Save(&subscription).Error; err != nil {
		return err
	}
	return nil
}


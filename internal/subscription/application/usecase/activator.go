package subscription_usecase

import (
	"context"
	"fmt"
	"time"

	subscription_application_eventbus "vault-app/internal/subscription/application"
	subscription_eventbus "vault-app/internal/subscription/application"
	subscription_domain "vault-app/internal/subscription/domain"
)

// ------------------------------------------------------
// Ports required by the activation process
// ------------------------------------------------------

type VaultPort interface {
    UpdateStorage(ctx context.Context, userID string, gb int) error
    EnableCloudBackup(ctx context.Context, userID string, enabled bool) error
    EnableVersionHistory(ctx context.Context, userID string, days int) error
    EnableTracecore(ctx context.Context, userID string) error
}

type StellarPort interface {
    LogSubscriptionActivated(ctx context.Context, userID, tier string) (string, error)
    VerifyAccount(ctx context.Context, publicKey string) (bool, error)
    CreatePaymentSchedule(ctx context.Context, publicKey string, amount float64, interval string) (string, error)
    RequestPayment(ctx context.Context, publicKey string, amount float64, description string) (string, error)
}


// ------------------------------------------------------
// Activator (main business logic)
// ------------------------------------------------------

type SubscriptionActivator struct {
    repo    subscription_domain.SubscriptionRepository
    bus     subscription_eventbus.SubscriptionEventBus
    vault   VaultPort
    stellar StellarPort
}

func NewSubscriptionActivator(
    repo subscription_domain.SubscriptionRepository,
    bus subscription_eventbus.SubscriptionEventBus,
    vault VaultPort,
) *SubscriptionActivator {
    return &SubscriptionActivator{
        repo:    repo,
        bus:     bus,
        vault:   vault,
    }
}

// Listen for the "SubscriptionCreated" event
func (s *SubscriptionActivator) Listen(ctx context.Context) error {
    return s.bus.SubscribeToCreation(func(ctx context.Context, e subscription_eventbus.SubscriptionCreated) {
        _ = s.Activate(ctx, e)
    })
}

// ------------------------------------------------------
// Main activation logic
// ------------------------------------------------------

func (s *SubscriptionActivator) Activate(ctx context.Context, evt subscription_eventbus.SubscriptionCreated) error {

    // Load subscription
    sub, err := s.repo.GetByID(ctx, evt.SubscriptionID)
    if err != nil {
        return err
    }
    

    // Already active?
    if sub.Active {
        return nil
    }

    // Reset – prevent inheritance bugs
    sub.Features = subscription_domain.SubscriptionFeatures{
        SubscriptionID: sub.ID,
        Compliance:     []string{},
    }

    // Apply tier features + storage + upgrade logic
    s.applyTierFeatures(ctx, sub)
    fmt.Printf("🚀 Activated subscription=%s for user=%s tier=%s ledger=%d", sub.ID, sub.UserID, sub.Tier, sub.Ledger)  

    // Mark it active
    sub.Active = true
    sub.ActivatedAt = time.Now().Unix()

    // Compute next billing date
    s.computeBillingCycle(sub)

    // Save updated subscription
    if err := s.repo.Save(ctx, sub); err != nil {
        return err
    }

    // Optional: log into Stellar ledger
    var ledger int32 = 0
    if s.stellar != nil {
        _, _ = s.stellar.LogSubscriptionActivated(ctx, sub.UserID, sub.Tier)
    }

    // Emit SubscriptionActivated event
    activated := subscription_application_eventbus.SubscriptionActivated{
        SubscriptionID: sub.ID,
        UserID:         sub.UserID,
        Tier:           sub.Tier,
        TxHash:         sub.TxHash,
        Ledger:         ledger,
        OccurredAt:     time.Now().Unix(),
    }
    fmt.Println("precessing to publish activated event...")
    _ = s.bus.PublishActivated(ctx, activated)

    return nil
}

// ------------------------------------------------------
// Apply all feature flags based on tier
// ------------------------------------------------------

func (s *SubscriptionActivator) applyTierFeatures(ctx context.Context, sub *subscription_domain.Subscription) {  
    
    
    switch sub.Tier {

    case "free":
        sub.Features.StorageGB = 5

    case "pro":
        sub.Features.StorageGB = 100
        sub.Features.CloudBackup = true

    case "pro_plus":    
        sub.Features.StorageGB = 200
        sub.Features.CloudBackup = true
        sub.Features.VersionHistory = true
        sub.Features.VersionHistoryDays = 30

    case "business":
        sub.Features.StorageGB = 1024
        sub.Features.CloudBackup = true
        sub.Features.VersionHistory = true
        sub.Features.VersionHistoryDays = 90
        sub.Features.Tracecore = true

    case "enterprise":
        sub.Features.StorageGB = 20480
        sub.Features.CloudBackup = true
        sub.Features.VersionHistory = true
        sub.Features.VersionHistoryDays = 365
        sub.Features.Tracecore = true
    }
}

// ------------------------------------------------------
// Compute billing cycle: trial, expiration, next billing
// ------------------------------------------------------

func (s *SubscriptionActivator) computeBillingCycle(sub *subscription_domain.Subscription) {
    fmt.Println("precessing to compute billing cycle ...")  
    // Trial? (only if defined in config or domain)
    if sub.TrialEndsAt > 0 {
        sub.NextBillingDate = sub.TrialEndsAt
    }
    
    // Normal paid plan
    if sub.Months > 0 {
        sub.NextBillingDate = time.Now().AddDate(0, sub.Months, 0).Unix()
        sub.BillingCycle = "monthly"
    }

    // Enterprise yearly billing?
    if sub.Tier == "enterprise" {
        sub.NextBillingDate = time.Now().AddDate(1, 0, 0).Unix()
        sub.BillingCycle = "yearly"
    }
}
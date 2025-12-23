package subscription_usecase

import (
	"context"
	"fmt"
	"time"

	utils "vault-app/internal"
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

    // 1. Load subscription
    sub, err := s.repo.GetByID(ctx, evt.SubscriptionID)
    if err != nil {
        return err
    }

    // 1.5. Optional: log into Stellar ledger
    var ledger int32 = 0
    if s.stellar != nil {
        _, _ = s.stellar.LogSubscriptionActivated(ctx, sub.UserID, sub.Tier)
    }

    // 2. Update Activation subscription status
    sub.Active = true
    sub.ActivatedAt = time.Now().Unix()
    if err := s.repo.Update(ctx, sub); err != nil {
        return err
    }
    utils.LogPretty("Subscription Activated", sub)
    

    // 3. Emit SubscriptionActivated event
    activated := subscription_application_eventbus.SubscriptionActivated{
        SubscriptionID: sub.ID,
        UserID:         sub.UserID,
        Password:       evt.Password,
        Tier:           sub.Tier,
        TxHash:         sub.TxHash,
        Ledger:         ledger,
        OccurredAt:     time.Now().Unix(),
    }
    // 4. Publish activated event
    fmt.Println("precessing to publish activated event...")
    _ = s.bus.PublishActivated(ctx, activated)

    return nil
}


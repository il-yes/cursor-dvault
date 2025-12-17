// internal/subscriptions/application/usecase/monitor_activation.go
package subscription_usecase

import (
	"context"
	"fmt"
	"time"
	utils "vault-app/internal"
	"vault-app/internal/logger/logger"
	"vault-app/internal/models"
	onboarding_domain "vault-app/internal/onboarding/domain"
	subscription_eventbus "vault-app/internal/subscription/application"

	"github.com/google/uuid"
)

type SubscriptionActivationMonitor struct {
	Logger                   *logger.Logger
	Bus                      subscription_eventbus.SubscriptionEventBus
	UserOnboardingRepository onboarding_domain.UserRepository
	Db                       models.DBModel
}

func NewSubscriptionActivationMonitor(log *logger.Logger, bus subscription_eventbus.SubscriptionEventBus, userOnboardingRepository onboarding_domain.UserRepository, db models.DBModel) *SubscriptionActivationMonitor {
	return &SubscriptionActivationMonitor{Logger: log, Bus: bus, UserOnboardingRepository: userOnboardingRepository, Db: db}
}

func (m *SubscriptionActivationMonitor) Listen(ctx context.Context) {
	m.Logger.Info("🛰️ Listening for subscription activations")

	m.Bus.SubscribeToActivation(func(ctx context.Context, event subscription_eventbus.SubscriptionActivated) {
		m.Logger.Info("🚀 Activated subscription=%s for user=%s tier=%s ledger=%d",
			event.SubscriptionID, event.UserID, event.Tier, event.Ledger)

		// Tier side effects (emails, notifications)
		switch event.Tier {
		case "free":
			m.Logger.Info("🧊 Free tier enabled")
		case "pro":
			m.Logger.Info("🔥 Pro features enabled")
		case "enterprise":
			m.Logger.Warn("🏢 Enterprise tier may need approval")
		default:
			m.Logger.Warn("⚠️ Unknown tier=%s", event.Tier)
		}

		m.Logger.Info("📧 Email queued for user=%s", event.UserID)
		m.Logger.Info("✅ Activation complete for subscription=%s", event.SubscriptionID)

		// 1. retrieve user subscription (user_dbs)
		userSubscription, err := m.UserOnboardingRepository.FindByEmail(event.UserID)
		if err != nil {
			m.Logger.Error("Monitor - Failed to retrieve user onboarding: %v", err)
			return
		}
		fmt.Println("Monitor - User onboarding retrieved: %v", userSubscription)

		// 2. Create and save user vault model
		userVault := models.User{
			ID:              uuid.NewString(),
			Role:            "user",
			Username:        event.UserID,
			Email:           event.UserID,
			Password:        userSubscription.Password,
			LastConnectedAt: time.Now(),
		}
		utils.LogPretty("userVault", userVault)
		user, err := m.Db.CreateUser(&userVault)
		if err != nil {
			m.Logger.Error("Monitor - Failed to create user vault: %v", err)
			return
		}
		m.Logger.Info("Monitor - User vault created: %v", user)

		// 3. Create and save user config model	
		
		
		m.Logger.Info("Monitor - User config created: %v", user)	
	})

	<-ctx.Done()
	m.Logger.Warn("🛑 SubscriptionActivationMonitor stopped")
}


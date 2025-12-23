// internal/subscriptions/application/usecase/listener_subscription_created.go
package subscription_usecase

import (
	"context"
	"vault-app/internal/logger/logger"
	subscription_application_eventbus "vault-app/internal/subscription/application"
)

type SubscriptionCreatedListener struct {
	log       *logger.Logger
	Activator *SubscriptionActivator
	Bus       subscription_application_eventbus.SubscriptionEventBus
}

func NewSubscriptionCreatedListener(
	log *logger.Logger,
	activator *SubscriptionActivator,
	bus subscription_application_eventbus.SubscriptionEventBus,
) *SubscriptionCreatedListener {
	return &SubscriptionCreatedListener{log: log, Activator: activator, Bus: bus}
}

func (l *SubscriptionCreatedListener) Listen(ctx context.Context) {
	l.Bus.SubscribeToCreation(func(ctx context.Context, e subscription_application_eventbus.SubscriptionCreated) {
		err := l.Activator.Activate(ctx, e)
		if err != nil {
			l.log.Warn("Failed to activate subscription %s: %v", e.SubscriptionID, err)
			return
		}
		l.log.Info("Subscription activated: %s", e.SubscriptionID)
	})

	<-ctx.Done()
	l.log.Warn("🛑 SubscriptionCreatedListener stopped")
}
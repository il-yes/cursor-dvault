package subscription_usecase

import (
	"context"
	"time"
	utils "vault-app/internal"
	subscription_eventbus "vault-app/internal/subscription/application"
	subscription_domain "vault-app/internal/subscription/domain"
)

type CreateSubscriptionUseCase struct {
	repo  subscription_domain.SubscriptionRepository
	bus   subscription_eventbus.SubscriptionEventBus
	AnkhorClient AnkhoraClient
}

func NewCreateSubscriptionUseCase(repo subscription_domain.SubscriptionRepository, bus subscription_eventbus.SubscriptionEventBus, ankhorClient AnkhoraClient) *CreateSubscriptionUseCase {
	return &CreateSubscriptionUseCase{repo: repo, bus: bus, AnkhorClient: ankhorClient}
}

type AnkhoraClient interface {
	GetSubscriptionBySessionID(ctx context.Context, sessionID string) (*subscription_domain.Subscription, error)
}
type TraditionalSubscriptionResponse struct {
	Subscription *subscription_domain.Subscription `json:"subscription"`
	OccurredAt     int64  `json:"occurred_at"`
}

func (uc *CreateSubscriptionUseCase) Execute(ctx context.Context, subID string, email string, plainPassword string) (*subscription_domain.Subscription, error) {
	// I. ---------- Fetch subscription from payment from Ankhora cloud ----------
	subscriptionFromPayment, err := uc.AnkhorClient.GetSubscriptionBySessionID(ctx, subID)
	if err != nil {
		return nil, err
	}

	// II. ---------- Validate subscription ----------
	if err := subscriptionFromPayment.Validate(); err != nil {
		return nil, err
	}
	// fmt.Println("Subscription from payment validated:", subscriptionFromPayment.Validate())

	subscriptionFromPayment.Email = email
	utils.LogPretty("Subscription from payment with email added:", subscriptionFromPayment)

	// III. ---------- Save subscription ----------
	if err := uc.repo.Save(ctx, subscriptionFromPayment); err != nil {
		return nil, err
	}

	// IV. ---------- Fire event immediately after saving ----------
	if uc.bus != nil {
		_ = uc.bus.PublishCreated(ctx, subscription_eventbus.SubscriptionCreated{
			SubscriptionID: subscriptionFromPayment.ID,
			UserID:         subscriptionFromPayment.UserID,
			UserEmail:      subscriptionFromPayment.Email,
			Password:       plainPassword,
			Tier:           subscriptionFromPayment.Tier,
			OccurredAt:     time.Now().Unix(),
		})
	}

	return subscriptionFromPayment, nil
}

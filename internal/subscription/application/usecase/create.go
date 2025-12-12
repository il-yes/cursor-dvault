package subscription_usecase

import (
	"context"
	subscription_eventbus "vault-app/internal/subscription/application"
	subscription_domain "vault-app/internal/subscription/domain"
)

type CreateSubscriptionUseCase struct {
	repo subscription_domain.SubscriptionRepository
	bus  subscription_eventbus.SubscriptionEventBus
	idGen func() string
}

func NewCreateSubscriptionUseCase(repo subscription_domain.SubscriptionRepository, bus subscription_eventbus.SubscriptionEventBus, idGen func() string) *CreateSubscriptionUseCase {
	return &CreateSubscriptionUseCase{repo: repo, bus: bus, idGen: idGen}
}	

func (uc *CreateSubscriptionUseCase) Execute(ctx context.Context, userID string, tier subscription_domain.SubscriptionTier) (*subscription_domain.Subscription, error) {
	id := uc.idGen()
	s := &subscription_domain.Subscription{ID: id, UserID: userID, Tier: string(tier), Active: true}
	if err := uc.repo.Save(ctx, s); err != nil {
		return nil, err
	}
	if uc.bus != nil {
		_ = uc.bus.PublishCreated(ctx, subscription_eventbus.SubscriptionCreated{SubscriptionID: id, UserID: userID, Tier: string(tier)})
	}
	return s, nil
}

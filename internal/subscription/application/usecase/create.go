package subscription_usecase

import (
	"context"
	subscription_eventbus "vault-app/internal/subscription/application"
	subscription_domain "vault-app/internal/subscription/domain"
)

type CreateSubscriptionUseCase struct {
	repo subscription_domain.Repository
	bus  subscription_eventbus.EventBus
	idGen func() string
}

func NewCreateSubscriptionUseCase(repo subscription_domain.Repository, bus subscription_eventbus.EventBus, idGen func() string) *CreateSubscriptionUseCase {
	return &CreateSubscriptionUseCase{repo: repo, bus: bus, idGen: idGen}
}	

func (uc *CreateSubscriptionUseCase) Execute(ctx context.Context, userID string, tier subscription_domain.SubscriptionTier) (*subscription_domain.Subscription, error) {
	id := uc.idGen()
	s := &subscription_domain.Subscription{ID: id, UserID: userID, Tier: tier, Active: true}
	if err := uc.repo.Save(ctx, s); err != nil {
		return nil, err
	}
	if uc.bus != nil {
		_ = uc.bus.PublishSubscriptionCreated(ctx, subscription_domain.SubscriptionCreated{SubscriptionID: id, UserID: userID, Tier: tier})
	}
	return s, nil
}

package subscription_ui_wails

import (
	"context"
	"time"
	subscription_usecase "vault-app/internal/subscription/application/usecase"
	subscription_domain "vault-app/internal/subscription/domain"
)

type SubscriptionHandler struct {
	CreateUC subscription_usecase.CreateSubscriptionUseCase
	SubscriptionRepository subscription_domain.SubscriptionRepository
}

func NewSubscriptionHandler(createUC subscription_usecase.CreateSubscriptionUseCase, subscriptionRepository subscription_domain.SubscriptionRepository) *SubscriptionHandler {
	return &SubscriptionHandler{CreateUC: createUC, SubscriptionRepository: subscriptionRepository}
}

func (h *SubscriptionHandler) CreateSubscription(ctx context.Context, id string) (*subscription_usecase.TraditionalSubscriptionResponse, error) {
	// create subscription
	response, err := h.CreateUC.Execute(ctx, id)
	if err != nil {
		return nil, err
	}

	return &subscription_usecase.TraditionalSubscriptionResponse{
		Subscription: response,
		OccurredAt:   time.Now().Unix(),
	}, nil
}



func (h *SubscriptionHandler) GetSubscription(ctx context.Context, id string) (*subscription_usecase.TraditionalSubscriptionResponse, error) {
	// get subscription
	response, err := h.SubscriptionRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &subscription_usecase.TraditionalSubscriptionResponse{
		Subscription: response,
		OccurredAt:   time.Now().Unix(),
	}, nil
}

func (h *SubscriptionHandler) SaveSubscription(ctx context.Context, s *subscription_domain.Subscription) error {
	// save subscription
	return h.SubscriptionRepository.Save(ctx, s)
}
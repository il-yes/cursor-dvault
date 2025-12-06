package persistence

import (
	"context"
	"sync"

	subscription_domain "vault-app/internal/subscription/domain"
)

type MemorySubscriptionRepository struct {
	mu sync.RWMutex
	byID map[string]*subscription_domain.Subscription
	byUser map[string]*subscription_domain.Subscription
}

func NewMemorySubscriptionRepository() *MemorySubscriptionRepository {
	return &MemorySubscriptionRepository{byID: make(map[string]*subscription_domain.Subscription), byUser: make(map[string]*subscription_domain.Subscription)}
}

func (r *MemorySubscriptionRepository) Save(ctx context.Context, s *subscription_domain.Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[s.ID] = s
	r.byUser[s.UserID] = s
	return nil
}

func (r *MemorySubscriptionRepository) FindByUserID(ctx context.Context, userID string) (*subscription_domain.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.byUser[userID]
	if !ok {
		return nil, nil
	}
	return s, nil
}

var _ subscription_domain.Repository = (*MemorySubscriptionRepository)(nil)

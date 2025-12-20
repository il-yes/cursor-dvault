package billing_persistence

import (
	"context"
	"sync"
	billing_domain "vault-app/internal/billing/domain"
)

type MemoryBillingRepository struct {
	mu sync.RWMutex
	byID map[string]*billing_domain.BillingInstrument
	byUser map[string][]*billing_domain.BillingInstrument
}

func NewMemoryBillingRepository() *MemoryBillingRepository {
	return &MemoryBillingRepository{byID: make(map[string]*billing_domain.BillingInstrument), byUser: make(map[string][]*billing_domain.BillingInstrument)}
}

func (r *MemoryBillingRepository) Save(ctx context.Context, b *billing_domain.BillingInstrument) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[b.ID] = b
	r.byUser[b.UserID] = append(r.byUser[b.UserID], b)
	return nil
}

func (r *MemoryBillingRepository) FindByUserID(ctx context.Context, userID string) (*billing_domain.BillingInstrument, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.byUser[userID][0], nil
}

var _ billing_domain.Repository = (*MemoryBillingRepository)(nil)

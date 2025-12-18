package identity_persistence

import (
	"context"
	"errors"
	"sync"
	identity_domain "vault-app/internal/identity/domain"
)

type MemoryUserRepository struct {
	mu    sync.RWMutex
	byID  map[string]*identity_domain.User
	byEmail map[string]*identity_domain.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{byID: make(map[string]*identity_domain.User), byEmail: make(map[string]*identity_domain.User)}
}

func (r *MemoryUserRepository) Save(ctx context.Context, u *identity_domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[u.ID]; ok {
		return errors.New("user id exists")
	}
	r.byID[u.ID] = u
	if u.Email != "" {
		r.byEmail[u.Email] = u
	}
	return nil
}

func (r *MemoryUserRepository) FindByID(ctx context.Context, id string) (*identity_domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (r *MemoryUserRepository) FindByEmail(ctx context.Context, email string) (*identity_domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byEmail[email]
	if !ok {
		return nil, nil
	}
	return u, nil
}

// Ensure interface satisfaction at compile-time
var _ identity_domain.UserRepository = (*MemoryUserRepository)(nil)
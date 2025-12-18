package identity_infrastructure_eventbus

import (
	"context"
	"errors"
	"sync"

	identity_eventbus "vault-app/internal/identity/application"
)

type MemoryEventBus struct {
	mu                     sync.RWMutex
	byID                   map[string]*identity_eventbus.UserRegistered
	byUser                 map[string]*identity_eventbus.UserRegistered
	userRegisteredHandlers []identity_eventbus.UserRegisteredHandler
	userLoggedInHandlers   []identity_eventbus.UserLoggedInHandler
}

func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		byID:   make(map[string]*identity_eventbus.UserRegistered),
		byUser: make(map[string]*identity_eventbus.UserRegistered),
	}
}

func (r *MemoryEventBus) PublishUserRegistered(ctx context.Context, e identity_eventbus.UserRegistered) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Store event (legacy behavior / for testing)
	r.byID[e.UserID] = &e
	r.byUser[e.UserID] = &e

	// Notify subscribers
	for _, handler := range r.userRegisteredHandlers {
		go handler(ctx, e)
	}

	return nil
}

func (r *MemoryEventBus) SubscribeToUserRegistered(handler identity_eventbus.UserRegisteredHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.userRegisteredHandlers = append(r.userRegisteredHandlers, handler)
	return nil
}

func (r *MemoryEventBus) PublishUserLoggedIn(ctx context.Context, e identity_eventbus.UserLoggedIn) error {
	r.mu.RLock()
	handlers := r.userLoggedInHandlers
	r.mu.RUnlock()

	for _, handler := range handlers {
		go handler(ctx, e)
	}
	return nil
}

func (r *MemoryEventBus) SubscribeToUserLoggedIn(handler identity_eventbus.UserLoggedInHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.userLoggedInHandlers = append(r.userLoggedInHandlers, handler)
	return nil
}

// Keep these for backward compatibility/testing if needed, or remove if unused.
// Since they were public methods, I'll keep them but attached to the new struct name.
func (r *MemoryEventBus) Save(ctx context.Context, s *identity_eventbus.UserRegistered) error {
	return r.PublishUserRegistered(ctx, *s)
}

func (r *MemoryEventBus) FindByUserID(ctx context.Context, userID string) (*identity_eventbus.UserRegistered, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.byUser[userID]
	if !ok {
		return nil, nil
	}
	return s, nil
}

func (r *MemoryEventBus) GetByID(ctx context.Context, id string) (*identity_eventbus.UserRegistered, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.byID[id]
	if !ok {
		return nil, errors.New(identity_eventbus.ErrUserRegisteredNotFound)
	}
	return s, nil
}

var _ identity_eventbus.EventBus = (*MemoryEventBus)(nil)

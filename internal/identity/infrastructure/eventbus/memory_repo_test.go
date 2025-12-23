package identity_infrastructure_eventbus

import (
	"context"
	"testing"
	"time"

	identity_eventbus "vault-app/internal/identity/application"

	"github.com/stretchr/testify/assert"
)

func TestMemoryEventBus_UserRegistered(t *testing.T) {
	bus := NewMemoryEventBus()
	ctx := context.Background()

	received := make(chan identity_eventbus.UserRegistered, 1)

	// Subscribe
	err := bus.SubscribeToUserRegistered(func(ctx context.Context, e identity_eventbus.UserRegistered) {
		received <- e
	})
	assert.NoError(t, err)

	// Publish
	event := identity_eventbus.UserRegistered{
		UserID:      "user-123",
		IsAnonymous: false,
		OccurredAt:  time.Now().Unix(),
	}
	err = bus.PublishUserRegistered(ctx, event)
	assert.NoError(t, err)

	// Verify receipt
	select {
	case e := <-received:
		assert.Equal(t, event.UserID, e.UserID)
		assert.Equal(t, event.IsAnonymous, e.IsAnonymous)
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for event")
	}

	// Verify storage (legacy behavior)
	saved, err := bus.FindByUserID(ctx, "user-123")
	assert.NoError(t, err)
	assert.NotNil(t, saved)
	assert.Equal(t, event.UserID, saved.UserID)
}

func TestMemoryEventBus_UserLoggedIn(t *testing.T) {
	bus := NewMemoryEventBus()
	ctx := context.Background()

	received := make(chan identity_eventbus.UserLoggedIn, 1)

	// Subscribe
	err := bus.SubscribeToUserLoggedIn(func(ctx context.Context, e identity_eventbus.UserLoggedIn) {
		received <- e
	})
	assert.NoError(t, err)

	// Publish
	event := identity_eventbus.UserLoggedIn{
		UserID:     "user-456",
		Email:      "test@example.com",
		OccurredAt: time.Now(),
	}
	err = bus.PublishUserLoggedIn(ctx, event)
	assert.NoError(t, err)

	// Verify receipt
	select {
	case e := <-received:
		assert.Equal(t, event.UserID, e.UserID)
		assert.Equal(t, event.Email, e.Email)
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for event")
	}
}

func TestMemoryEventBus_ImplementsInterface(t *testing.T) {
	var _ identity_eventbus.EventBus = (*MemoryEventBus)(nil)
}

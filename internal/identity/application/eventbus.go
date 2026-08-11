package identity_eventbus

import (
	"context"
	"time"
)

type UserRegistered struct {
	UserID      string
	IsAnonymous bool
	OccurredAt  int64
}

type UserLoggedIn struct {
	UserID     string
	Email      string
	OccurredAt time.Time
}

type DeviceCreated struct {
	DeviceID   string
	VaultID    string
	PublicKey  string
	KeyType    string
	OccurredAt time.Time
}

type DeviceRevoked struct {
	DeviceID   string
	VaultID    string
	RevokedAt  time.Time
	OccurredAt time.Time
}

const (
	ErrUserRegisteredNotFound = "user_registered_not_found"
	ErrUserLoggedInNotFound   = "user_logged_in_not_found"
)

type UserRegisteredHandler func(context.Context, UserRegistered)
type UserLoggedInHandler func(context.Context, UserLoggedIn)
type DeviceCreatedHandler func(context.Context, DeviceCreated)
type DeviceRevokedHandler func(context.Context, DeviceRevoked)

// EventBus inbound port for identity events (application layer)
type EventBus interface {
	PublishUserRegistered(ctx context.Context, e UserRegistered) error
	SubscribeToUserRegistered(handler UserRegisteredHandler) error

	PublishUserLoggedIn(ctx context.Context, e UserLoggedIn) error
	SubscribeToUserLoggedIn(handler UserLoggedInHandler) error

	PublishDeviceCreated(ctx context.Context, e DeviceCreated) error
	SubscribeToDeviceCreated(handler DeviceCreatedHandler) error

	PublishDeviceRevoked(ctx context.Context, e DeviceRevoked) error
	SubscribeToDeviceRevoked(handler DeviceRevokedHandler) error
}
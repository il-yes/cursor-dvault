package channel_events

import (
	"context"

	channel_domain "vault-app/internal/channel/domain"
)

// -------- EVENTS --------

// -------- EVENT BUS --------
type ChannelEventBus interface {
	PublishChannelCreated(ctx context.Context, event channel_domain.ChannelCreated) error
	SubscribeToChannelCreated(handler func(ctx context.Context, event channel_domain.ChannelCreated)) error


	PublishChannelRevoked(ctx context.Context, event channel_domain.ChannelRevoked) error
	SubscribeToChannelRevoked(handler func(ctx context.Context, event channel_domain.ChannelRevoked)) error


	PublishChannelDeleted(ctx context.Context, event channel_domain.ChannelDeleted) error
	SubscribeToChannelDeleted(handler func(ctx context.Context, event channel_domain.ChannelDeleted)) error
}

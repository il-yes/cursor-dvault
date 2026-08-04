// internal/channel/domain/events.go

package channel_domain

import "time"

const (
	EventChannelCreated = "channel.created"
	EventChannelRevoked = "channel.revoked"
	EventChannelDeleted = "channel.deleted"
)

type ChannelCreated struct {
	EventID        string
	EventTimestamp time.Time

	ChannelID string
	WorkspaceID     string
}

func (ChannelCreated) EventType() string {
	return EventChannelCreated
}

type ChannelRevoked struct {
	EventID        string
	EventTimestamp time.Time

	ChannelID string

	WorkspaceID string
}

func (ChannelRevoked) EventType() string {
	return EventChannelRevoked
}

type ChannelDeleted struct {
	EventID        string
	EventTimestamp time.Time

	ChannelID string
	WorkspaceID     string
}

func (ChannelDeleted) EventType() string {
	return EventChannelDeleted
}
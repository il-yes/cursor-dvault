package thread_domain

import "time"

type ThreadCreated struct {
	ThreadID string
	ChannelId string
	WorkspaceId string
	Timestamp time.Time `json:"timestamp"`
	EventID   string    `json:"event_id"`
}

type ThreadUpdated struct {
	Thread Thread
	WorkspaceId string
	ChannelId string
	Timestamp time.Time `json:"timestamp"`
	EventID   string    `json:"event_id"`
}

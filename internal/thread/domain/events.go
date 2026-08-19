package thread_domain

import "time"

type ThreadCreated struct {
	EventID     string    `json:"event_id"`
	ThreadID    string    `json:"thread_id"`
	ChannelID   string    `json:"channel_id"`
	WorkspaceID string    `json:"workspace_id,omitempty"`
	AssetType   string    `json:"asset_type"`
	Timestamp   time.Time `json:"timestamp"`
}

type ThreadUpdated struct {
	EventID     string    `json:"event_id"`
	Thread      Thread    `json:"thread"`
	WorkspaceID string    `json:"workspace_id"`
	ChannelID   string    `json:"channel_id"`
	Timestamp   time.Time `json:"timestamp"`
}


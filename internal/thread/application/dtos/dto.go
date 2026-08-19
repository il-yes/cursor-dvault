package thread_dtos

type CreateThreadRequest struct {
	ChannelID   string `json:"channel_id"`
	WorkspaceID string `json:"workspace_id,omitempty"`
	AssetType   string `json:"asset_type"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
}

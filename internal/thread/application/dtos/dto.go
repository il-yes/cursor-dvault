package thread_dtos

type CreateThreadRequest struct {
	ChannelID string                     `json:"channel_id"`
	AssetType string                     `json:"asset_type"`
	Title     string                     `json:"title"`
	Subtitle  string                     `json:"subtitle"`
}

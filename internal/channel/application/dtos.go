package channel_usecase


type CreateChannelRequest struct {
	TemplateID  string
	Title       string
	WorkspaceID string
}

type ListChannelsRequest struct {
	WorkspaceID string
}

type ArchiveChannelRequest struct {
	ChannelID   string
	WorkspaceID string
}
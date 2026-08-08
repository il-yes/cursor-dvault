package channel_application


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

type GetChannelRequest struct {
	ChannelID string
}

type DeleteChannelRequest struct {
	ChannelID string
}
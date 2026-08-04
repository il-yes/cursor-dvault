package channel_usecase


type CreateChannelRequest struct {
	TemplateID  string
	Title       string
	WorkspaceID string
}
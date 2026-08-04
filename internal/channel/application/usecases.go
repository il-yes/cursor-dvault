package channel_usecase

// import channel_domain "vault-app/internal/channel/domain"



/*
type ChannelUsecase struct {
	Repo channel_domain.RepositoryInterface
}

func NewChannelUsecase(repo channel_domain.RepositoryInterface) *ChannelUsecase {
	return &ChannelUsecase{
		Repo: repo,
	}
}

func (c *ChannelUsecase) CreateChannel(ctx context.Context, req *CreateChannelRequest) (*channel_domain.Channel, error) {
  channel := channel_domain.NewChannel(req.UserID, req.Name)

	ch, err := c.Repo.CreateChannel(ctx, channel_domain.NewCreateChannelRequest{
    VaultID: req.VaultID,
    Channel: channel,
    Signature: req.Signature,
  })
	if err != nil {
		return nil, err
	}

  eventBus.PublishChannelCreated(ctx, channel_events.ChannelCreated{
    VaultID:     req.VaultID,
    ChannelID:   ch.ID,
    WorkspaceID: ch.WorkspaceID,
  })

	return ch, nil
}


func (c *ChannelUsecase) ListChannels(ctx context.Context, req *ListChannelsRequest) (*channel_domain.Channel, error) {
	ch, err := c.Repo.ListChannels(ctx, req)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (c *ChannelUsecase) GetChannel(ctx context.Context, req *GetChannelRequest) (*channel_domain.Channel, error) {
	ch, err := c.Repo.GetChannel(ctx, req)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (c *ChannelUsecase) UpdateChannel(ctx context.Context, req *UpdateChannelRequest) (*channel_domain.Channel, error) {
	ch, err := c.Repo.UpdateChannel(ctx, req)
	if err != nil {
		return nil, err
	}
	return ch, nil
}
*/




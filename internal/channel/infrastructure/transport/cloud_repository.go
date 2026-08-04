package transport

import (
	"context"

	channel_domain "vault-app/internal/channel/domain"
)

type CloudRepository struct {
	client channel_domain.ChannelRepository
}

func NewCloudRepository(client channel_domain.ChannelRepository) *CloudRepository {
	return &CloudRepository{
		client: client,
	}
}

func (r *CloudRepository) CreateChannel(ctx context.Context, req channel_domain.CreateChannelRequest, signature string) (*channel_domain.Channel, error) {
	resp, err := r.client.CreateChannel(ctx, &req)
	return &resp.Data, err
}

func (r *CloudRepository) UpdateChannel(
	ctx context.Context, req channel_domain.UpdateChannelRequest) (*channel_domain.Channel, error) {

	resp, err := r.client.UpdateChannel(ctx, &req)
	return &resp.Data, err
}

func (r *CloudRepository) DeleteChannel( ctx context.Context, req channel_domain.DeleteChannelRequest) (*channel_domain.Channel, error) {
	err := r.client.DeleteChannel(ctx, &req)
	return nil, err
}
func (r *CloudRepository) RevokeChannel( ctx context.Context, req channel_domain.RevokeInvitationRequest) (*channel_domain.Channel, error) {
	err := r.client.RevokeChannel(ctx, &req)
	return nil, err
}

func (r *CloudRepository) GetChannel( ctx context.Context, req channel_domain.GetChannelRequest) (*channel_domain.Channel, error) {
	channel, err := r.client.GetChannel(ctx, &req)
	return &channel.Data, err
}
func (r *CloudRepository) ListChannel( ctx context.Context, req channel_domain.ListChannelsRequest) ([]channel_domain.Channel, error) {
	channel, err := r.client.ListChannels(ctx, &req)
	return channel.Data, err
}

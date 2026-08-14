package tracecore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	channel_domain "vault-app/internal/channel/domain"
	tracecore_types "vault-app/internal/tracecore/types"
)

// ChannelRepository Implementation on TracecoreClient

func (c *TracecoreClient) CreateChannel(ctx context.Context, req *channel_domain.CreateChannelRequest) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
	ch := req.Channel

	payload := map[string]interface{}{
		"workspace_id": ch.WorkspaceID,
		"title":        ch.Title,
		"template_id":  ch.TemplateID,
		"status":       string(ch.Status),
	}

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	baseUrl := c.AnkhoraCloudUrl
	if baseUrl == "" {
		baseUrl = c.BaseURL
	}
	url := baseUrl + "/channels"

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var cloudResp tracecore_types.CloudResponse[channel_domain.Channel]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && (cloudResp.Success || cloudResp.Data.ID != "") {
		return &cloudResp, nil
	}

	var channelDTO tracecore_types.ChannelDTO
	if errDTO := json.Unmarshal(respBytes, &channelDTO); errDTO == nil && channelDTO.ID != "" {
		createdChannel := channel_domain.NewChannel(channelDTO.TemplateID, channelDTO.Title, channelDTO.WorkspaceID)
		createdChannel.ID = channelDTO.ID
		if channelDTO.Status != "" {
			createdChannel.Status = channel_domain.ChannelStatus(channelDTO.Status)
		}
		if !channelDTO.CreatedAt.IsZero() {
			createdChannel.CreatedAt = channelDTO.CreatedAt
		}
		if !channelDTO.UpdatedAt.IsZero() {
			createdChannel.UpdatedAt = channelDTO.UpdatedAt
		}
		return &tracecore_types.CloudResponse[channel_domain.Channel]{
			Status:  200,
			Data:    createdChannel,
			Message: "success",
			Success: true,
		}, nil
	}

	return nil, fmt.Errorf("TracecoreClient - CreateChannel - cloud response unmarshal failed: %s", string(respBytes))
}

func (c *TracecoreClient) ListChannels(ctx context.Context, req *channel_domain.ListChannelsRequest) (*tracecore_types.CloudResponse[[]channel_domain.Channel], error) {
	baseUrl := c.AnkhoraCloudUrl
	if baseUrl == "" {
		baseUrl = c.BaseURL
	}
	url := baseUrl + "/channels/workspace/" + req.WorkspaceID 

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var cloudResp tracecore_types.CloudResponse[[]channel_domain.Channel]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && (cloudResp.Success || cloudResp.Data != nil) {
		return &cloudResp, nil
	}

	var dtos []tracecore_types.ChannelDTO
	if err := json.Unmarshal(respBytes, &dtos); err == nil {
		channels := make([]channel_domain.Channel, 0, len(dtos))
		for _, dto := range dtos {
			ch := channel_domain.NewChannel(dto.TemplateID, dto.Title, dto.WorkspaceID)
			ch.ID = dto.ID
			if dto.Status != "" {
				ch.Status = channel_domain.ChannelStatus(dto.Status)
			}
			if !dto.CreatedAt.IsZero() {
				ch.CreatedAt = dto.CreatedAt
			}
			if !dto.UpdatedAt.IsZero() {
				ch.UpdatedAt = dto.UpdatedAt
			}
			channels = append(channels, ch)
		}
		return &tracecore_types.CloudResponse[[]channel_domain.Channel]{
			Status:  200,
			Data:    channels,
			Message: "success",
			Success: true,
		}, nil
	}

	return &tracecore_types.CloudResponse[[]channel_domain.Channel]{
		Status:  200,
		Data:    []channel_domain.Channel{},
		Message: "success",
		Success: true,
	}, nil
}

func (c *TracecoreClient) GetChannel(ctx context.Context, req *channel_domain.GetChannelRequest) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
	now := time.Now()
	ch := channel_domain.NewChannel("tpl_general", "General Channel", "ws_1")
	ch.ID = req.ChannelID
	ch.CreatedAt = now
	ch.UpdatedAt = now

	return &tracecore_types.CloudResponse[channel_domain.Channel]{
		Status:  200,
		Data:    ch,
		Message: "success",
		Success: true,
	}, nil
}

func (c *TracecoreClient) DeleteChannel(ctx context.Context, req *channel_domain.DeleteChannelRequest) error {
	return nil
}

func (c *TracecoreClient) UpdateChannel(ctx context.Context, req *channel_domain.UpdateChannelRequest) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
	ch := req.Channel
	ch.UpdatedAt = time.Now()
	return &tracecore_types.CloudResponse[channel_domain.Channel]{
		Status:  200,
		Data:    ch,
		Message: "success",
		Success: true,
	}, nil
}

func (c *TracecoreClient) ActivateChannel(ctx context.Context, req *channel_domain.AcceptInvitationRequest) error {
	return nil
}

func (c *TracecoreClient) RevokeChannel(ctx context.Context, req *channel_domain.RevokeInvitationRequest) error {
	return nil
}

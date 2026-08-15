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
//
// The Cloud backend returns Channel aggregates using default Go JSON encoding
// (capitalized field names). We therefore decode into tracecore_types.CloudChannelDTO
// and explicitly map into channel_domain.Channel. This mapping is the canonical
// Cloud -> Desktop Channel boundary shared by CreateChannel and ListChannels.

func (c *TracecoreClient) CreateChannel(ctx context.Context, req *channel_domain.CreateChannelRequest) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
	ch := req.Channel

	payload := map[string]interface{}{
		"workspace_id": ch.WorkspaceID,
		"title":        ch.Title,
		"template_id":  ch.TemplateID,
		"status":       string(ch.Status),
		"slots":        ch.Slots,
		"assignments":  ch.Assignments,
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

	// Cloud returns a channel envelope: { "status": 201, "data": { ... } }.
	var cloudResp tracecore_types.CloudResponse[tracecore_types.CloudChannelDTO]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && cloudResp.Data.ID != "" {
		return &tracecore_types.CloudResponse[channel_domain.Channel]{
			Status:  200,
			Data:    mapCloudChannelDTO(cloudResp.Data),
			Message: "success",
			Success: true,
		}, nil
	}

	// Tolerate a bare channel object as a fallback shape.
	var dto tracecore_types.CloudChannelDTO
	if errDTO := json.Unmarshal(respBytes, &dto); errDTO != nil || dto.ID == "" {
		return nil, fmt.Errorf("TracecoreClient - CreateChannel - unexpected Cloud response shape (missing channel data): %s", string(respBytes))
	}

	return &tracecore_types.CloudResponse[channel_domain.Channel]{
		Status:  200,
		Data:    mapCloudChannelDTO(dto),
		Message: "success",
		Success: true,
	}, nil
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

	// Cloud returns a list envelope: { "status": 200, "data": [ { ... } ] }.
	var cloudResp tracecore_types.CloudResponse[[]tracecore_types.CloudChannelDTO]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil {
		if cloudResp.Data == nil {
			return nil, fmt.Errorf("TracecoreClient - ListChannels - Cloud response is missing data: %s", string(respBytes))
		}
		channels, mapErr := mapCloudChannels(cloudResp.Data)
		if mapErr != nil {
			return nil, mapErr
		}
		return &tracecore_types.CloudResponse[[]channel_domain.Channel]{
			Status:  200,
			Data:    channels,
			Message: "success",
			Success: true,
		}, nil
	}

	// Tolerate a bare array of channels as a fallback shape.
	var dtos []tracecore_types.CloudChannelDTO
	if errArr := json.Unmarshal(respBytes, &dtos); errArr != nil {
		return nil, fmt.Errorf("TracecoreClient - ListChannels - unexpected Cloud response shape: %s", string(respBytes))
	}
	channels, mapErr := mapCloudChannels(dtos)
	if mapErr != nil {
		return nil, mapErr
	}
	return &tracecore_types.CloudResponse[[]channel_domain.Channel]{
		Status:  200,
		Data:    channels,
		Message: "success",
		Success: true,
	}, nil
}

// mapCloudChannels maps a decoded Cloud channel list into the Channel aggregate.
// A non-empty list must contain well-formed channels; a malformed element is an
// error, never a fabricated empty result.
func mapCloudChannels(dtos []tracecore_types.CloudChannelDTO) ([]channel_domain.Channel, error) {
	channels := make([]channel_domain.Channel, 0, len(dtos))
	for i, dto := range dtos {
		if dto.ID == "" {
			return nil, fmt.Errorf("TracecoreClient - ListChannels - Cloud returned a malformed channel at index %d: missing ID", i)
		}
		channels = append(channels, mapCloudChannelDTO(dto))
	}
	return channels, nil
}

// mapCloudChannelDTO is the canonical Cloud -> Desktop Channel mapping. Every
// field exposed by the Cloud Channel aggregate is preserved.
func mapCloudChannelDTO(dto tracecore_types.CloudChannelDTO) channel_domain.Channel {
	status := channel_domain.ChannelStatus(dto.Status)
	if status == "" {
		status = channel_domain.StatusPending
	}

	slots := make([]channel_domain.Slot, 0, len(dto.Slots))
	for _, s := range dto.Slots {
		slots = append(slots, channel_domain.Slot{
			ID:      s.ID,
			Name:    s.Name,
			Role:    s.Role,
			VaultID: s.VaultID,
			Gated:   s.Gated,
			Order:   s.Order,
		})
	}

	assignments := make([]channel_domain.Assignment, 0, len(dto.Assignments))
	for _, a := range dto.Assignments {
		assignments = append(assignments, channel_domain.Assignment{
			SlotID:       a.SlotID,
			OwnerID:      a.OwnerID,
			PublicKey:    a.PublicKey,
			VaultAddress: a.VaultAddress,
		})
	}

	properties := make([]channel_domain.ChannelProperty, 0, len(dto.Properties))
	for _, p := range dto.Properties {
		properties = append(properties, channel_domain.ChannelProperty{
			Key:   p.Key,
			Value: p.Value,
		})
	}

	return channel_domain.Channel{
		ID:          dto.ID,
		TemplateID:  dto.TemplateID,
		Title:       dto.Title,
		Status:      status,
		Slots:       slots,
		Assignments: assignments,
		Properties:  properties,
		Policy:      dto.Policy,
		Federation:  encodeFederation(dto.Federation),
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
		RevokedAt:   dto.RevokedAt,
		ArchivedAt:  dto.ArchivedAt,
		WorkspaceID: dto.WorkspaceID,
		IsDraft:     dto.IsDraft,
		IsDirty:     dto.IsDirty,
	}
}

// encodeFederation preserves the full Cloud Federation object
// (VaultAID, VaultBID, AllowedEventTypes, AllowedPaths, AllowedDirections) in
// the Channel aggregate's string Federation field.
func encodeFederation(f *tracecore_types.CloudChannelFederation) string {
	if f == nil {
		return ""
	}
	raw, err := json.Marshal(f)
	if err != nil {
		return ""
	}
	return string(raw)
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

// ActivateChannel activates an existing Channel on the Cloud backend.
// The activation request carries only the channel id; there is no request body.
// The Cloud aggregate is authoritative for the "every gated slot must be
// fulfilled" invariant and returns the activated Channel aggregate, which we
// map through the same CloudChannelDTO -> channel_domain.Channel boundary used
// by CreateChannel and ListChannels.
func (c *TracecoreClient) ActivateChannel(ctx context.Context, req *channel_domain.ActivateChannelRequest) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
	if req == nil || req.ChannelID == "" {
		return nil, channel_domain.ErrChannelIDRequired
	}

	baseUrl := c.AnkhoraCloudUrl
	if baseUrl == "" {
		baseUrl = c.BaseURL
	}
	url := baseUrl + "/channels/" + req.ChannelID + "/activate"

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
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

	// Surface the Cloud/domain outcome (e.g. revoked, gated slots unfulfilled)
	// verbatim so the UI can present it cleanly. The Cloud remains authoritative.
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	// Cloud returns an envelope: { "status": 201, "data": { ... } }.
	var cloudResp tracecore_types.CloudResponse[tracecore_types.CloudChannelDTO]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && cloudResp.Data.ID != "" {
		return &tracecore_types.CloudResponse[channel_domain.Channel]{
			Status:  200,
			Data:    mapCloudChannelDTO(cloudResp.Data),
			Message: "success",
			Success: true,
		}, nil
	}

	// Tolerate a bare channel object as a fallback shape.
	var dto tracecore_types.CloudChannelDTO
	if errDTO := json.Unmarshal(respBytes, &dto); errDTO != nil || dto.ID == "" {
		return nil, fmt.Errorf("TracecoreClient - ActivateChannel - unexpected Cloud response shape (missing channel data): %s", string(respBytes))
	}

	return &tracecore_types.CloudResponse[channel_domain.Channel]{
		Status:  200,
		Data:    mapCloudChannelDTO(dto),
		Message: "success",
		Success: true,
	}, nil
}

func (c *TracecoreClient) RevokeChannel(ctx context.Context, req *channel_domain.RevokeInvitationRequest) error {
	return nil
}

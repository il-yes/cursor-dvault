package tracecore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
		"slots":        ch.Slots,
		"properties":   ch.Properties,
		"assignments":  ch.Assignments,
		"policy":       ch.Policy,
	}
	// federation is the Desktop string form of the Cloud FederationConfig. When
	// present it is converted into the snake_case federation object the Cloud
	// CreateChannelRequest contract requires. Cloud remains authoritative for
	// the federation semantics.
	if fed := federationRequest(ch.Federation); fed != nil {
		payload["federation"] = fed
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

// federationRequest converts the Desktop Channel's Federation string (the
// JSON-encoded Cloud FederationConfig with capitalized field names) into the
// snake_case federation object the Cloud CreateChannelRequest contract expects.
// It returns nil when the Desktop carries no federation, so the key is omitted
// from the outbound payload.
func federationRequest(federation string) map[string]interface{} {
	if federation == "" {
		return nil
	}

	var fed tracecore_types.CloudChannelFederation
	if err := json.Unmarshal([]byte(federation), &fed); err != nil {
		return nil
	}

	return map[string]interface{}{
		"vault_a_id":          fed.VaultAID,
		"vault_b_id":          fed.VaultBID,
		"allowed_event_types": fed.AllowedEventTypes,
		"allowed_paths":       fed.AllowedPaths,
		"allowed_directions":  fed.AllowedDirections,
	}
}

// GetChannel fetches a single Channel from the authoritative Cloud backend
// (GET /channels/{id}). The Cloud response is mapped through the same
// CloudChannelDTO -> channel_domain.Channel boundary used by CreateChannel and
// ListChannels. HTTP >=400 (including the Cloud 404 for an unknown channel) is
// surfaced verbatim; no local channel is ever fabricated.
func (c *TracecoreClient) GetChannel(ctx context.Context, req *channel_domain.GetChannelRequest) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
	if req == nil || req.ChannelID == "" {
		return nil, channel_domain.ErrChannelIDRequired
	}

	baseUrl := c.AnkhoraCloudUrl
	if baseUrl == "" {
		baseUrl = c.BaseURL
	}
	url := baseUrl + "/channels/" + req.ChannelID

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

	// Surface the Cloud outcome (e.g. record not found) verbatim. The Cloud is
	// the single source of truth for channel existence.
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	// Cloud returns an envelope: { "status": 200, "data": { ... } }.
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
		return nil, fmt.Errorf("TracecoreClient - GetChannel - unexpected Cloud response shape (missing channel data): %s", string(respBytes))
	}

	return &tracecore_types.CloudResponse[channel_domain.Channel]{
		Status:  200,
		Data:    mapCloudChannelDTO(dto),
		Message: "success",
		Success: true,
	}, nil
}

// DeleteChannel deletes a Channel through the authoritative Cloud backend
// (DELETE /channels/{id}). A 2xx response is success; HTTP >=400 is surfaced
// verbatim. No local deletion or fabrication occurs — Cloud is the single
// source of truth for channel existence.
func (c *TracecoreClient) DeleteChannel(ctx context.Context, req *channel_domain.DeleteChannelRequest) error {
	if req == nil || req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	baseUrl := c.AnkhoraCloudUrl
	if baseUrl == "" {
		baseUrl = c.BaseURL
	}
	url := baseUrl + "/channels/" + req.ChannelID

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body failed: %w", err)
	}

	// Surface the Cloud outcome (e.g. record not found) verbatim. The Cloud
	// remains authoritative for the delete decision.
	if resp.StatusCode >= 400 {
		return fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	// HTTP 2xx is success. The Cloud delete response carries no Channel data,
	// so nothing is decoded or fabricated here.
	return nil
}

// UpdateChannel updates an existing Channel through the authoritative Cloud
// backend (PUT /channels/{id}). The request body carries exactly the fields the
// Cloud UpdateChannelRequest contract supports (id, title, slots, properties,
// assignments, policy) in snake_case. Cloud is authoritative for which fields
// are applied and for all domain validation; the returned Cloud Channel is
// mapped through the same CloudChannelDTO -> channel_domain.Channel boundary.
// HTTP >=400 is surfaced verbatim; no local mutation is performed.
func (c *TracecoreClient) UpdateChannel(ctx context.Context, req *channel_domain.UpdateChannelRequest) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
	if req == nil {
		return nil, channel_domain.ErrRequestRequired
	}
	if req.Channel.ID == "" {
		return nil, channel_domain.ErrChannelIDRequired
	}

	payload := map[string]interface{}{
		"id":          req.Channel.ID,
		"title":       req.Channel.Title,
		"slots":       req.Channel.Slots,
		"properties":  req.Channel.Properties,
		"assignments": req.Channel.Assignments,
		"policy":      req.Channel.Policy,
	}

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	baseUrl := c.AnkhoraCloudUrl
	if baseUrl == "" {
		baseUrl = c.BaseURL
	}
	url := baseUrl + "/channels/" + req.Channel.ID

	request, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
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

	// Surface the Cloud/domain outcome verbatim. The Cloud remains authoritative
	// for the update.
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	// Cloud returns an envelope: { "status": 200, "data": { ... } }.
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
		return nil, fmt.Errorf("TracecoreClient - UpdateChannel - unexpected Cloud response shape (missing channel data): %s", string(respBytes))
	}

	return &tracecore_types.CloudResponse[channel_domain.Channel]{
		Status:  200,
		Data:    mapCloudChannelDTO(dto),
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

// RevokeChannel revokes an existing Channel on the Cloud backend.
// The revocation request carries only the channel id; there is no request body.
// The Cloud aggregate is authoritative for the revocation invariants and emits
// the channel.revoked lifecycle event. The Cloud response is an envelope
// without Channel data, so this method returns only an error on success;
// callers refresh through ListChannels to observe the new status.
func (c *TracecoreClient) RevokeChannel(ctx context.Context, req *channel_domain.RevokeChannelRequest) error {
	if req == nil || req.ChannelID == "" {
		return channel_domain.ErrChannelIDRequired
	}

	baseUrl := c.AnkhoraCloudUrl
	if baseUrl == "" {
		baseUrl = c.BaseURL
	}
	url := baseUrl + "/channels/" + req.ChannelID + "/revoke"

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body failed: %w", err)
	}

	// Surface the Cloud/domain outcome (e.g. channel already revoked) verbatim
	// so the UI can present it cleanly. The Cloud remains authoritative.
	if resp.StatusCode >= 400 {
		return fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	// HTTP 200 is success. The Cloud revoke response does not carry Channel
	// data, so nothing is decoded or fabricated here.
	return nil
}

// AddParticipant joins a vault to a channel through the authoritative Cloud
// backend (POST /channels/{id}/participants). The request body carries exactly
// the fields the Cloud JoinChannelRequest contract requires; Cloud validates
// the join and returns the persisted Participant, which is mapped through the
// same Cloud -> Desktop boundary used for Channel aggregates. Cloud is
// idempotent, so an already-joined vault returns the existing participant.
func (c *TracecoreClient) AddParticipant(ctx context.Context, req *channel_domain.JoinChannelRequest) (*tracecore_types.CloudResponse[channel_domain.Participant], error) {
	if req == nil || req.ChannelID == "" {
		return nil, channel_domain.ErrChannelIDRequired
	}
	if req.VaultID == "" {
		return nil, channel_domain.ErrVaultIDRequired
	}

	payload := map[string]interface{}{
		"channel_id": req.ChannelID,
		"vault_id":   req.VaultID,
		"public_key": req.PublicKey,
		"direction":  req.Direction,
	}
	if req.SlotID != "" {
		payload["slot_id"] = req.SlotID
	}
	if req.Role != "" {
		payload["role"] = req.Role
	}

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	baseUrl := c.AnkhoraCloudUrl
	if baseUrl == "" {
		baseUrl = c.BaseURL
	}
	url := baseUrl + "/channels/" + req.ChannelID + "/participants"

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

	// Surface the Cloud/domain outcome (e.g. channel revoked, slot not found)
	// verbatim. The Cloud remains authoritative for the join decision.
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	// Cloud returns an envelope: { "status": 201, "data": { ... } }.
	var cloudResp tracecore_types.CloudResponse[tracecore_types.CloudChannelParticipant]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && cloudResp.Data.VaultID != "" {
		return &tracecore_types.CloudResponse[channel_domain.Participant]{
			Status:  200,
			Data:    mapCloudChannelParticipant(cloudResp.Data),
			Message: "success",
			Success: true,
		}, nil
	}

	// Tolerate a bare participant object as a fallback shape.
	var dto tracecore_types.CloudChannelParticipant
	if errDTO := json.Unmarshal(respBytes, &dto); errDTO != nil || dto.VaultID == "" {
		return nil, fmt.Errorf("TracecoreClient - AddParticipant - unexpected Cloud response shape (missing participant data): %s", string(respBytes))
	}

	return &tracecore_types.CloudResponse[channel_domain.Participant]{
		Status:  200,
		Data:    mapCloudChannelParticipant(dto),
		Message: "success",
		Success: true,
	}, nil
}

// ListParticipants returns the vaults Cloud has persisted as channel
// participants (GET /channels/{id}/participants). An empty participant list is
// valid; a malformed element is an error, never a fabricated success.
func (c *TracecoreClient) ListParticipants(ctx context.Context, req *channel_domain.ListParticipantsRequest) (*tracecore_types.CloudResponse[[]channel_domain.Participant], error) {
	if req == nil || req.ChannelID == "" {
		return nil, channel_domain.ErrChannelIDRequired
	}

	baseUrl := c.AnkhoraCloudUrl
	if baseUrl == "" {
		baseUrl = c.BaseURL
	}
	url := baseUrl + "/channels/" + req.ChannelID + "/participants"

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
	// With no participants the Cloud marshals a nil slice as null, which is a
	// valid empty list.
	var cloudResp tracecore_types.CloudResponse[[]tracecore_types.CloudChannelParticipant]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil {
		participants, mapErr := mapCloudChannelParticipants(cloudResp.Data)
		if mapErr != nil {
			return nil, mapErr
		}
		return &tracecore_types.CloudResponse[[]channel_domain.Participant]{
			Status:  200,
			Data:    participants,
			Message: "success",
			Success: true,
		}, nil
	}

	// Tolerate a bare array of participants as a fallback shape.
	var dtos []tracecore_types.CloudChannelParticipant
	if errArr := json.Unmarshal(respBytes, &dtos); errArr != nil {
		return nil, fmt.Errorf("TracecoreClient - ListParticipants - unexpected Cloud response shape: %s", string(respBytes))
	}
	participants, mapErr := mapCloudChannelParticipants(dtos)
	if mapErr != nil {
		return nil, mapErr
	}
	return &tracecore_types.CloudResponse[[]channel_domain.Participant]{
		Status:  200,
		Data:    participants,
		Message: "success",
		Success: true,
	}, nil
}

// mapCloudChannelParticipants maps a decoded Cloud participant list into the
// Channel aggregate domain. A non-empty list must contain well-formed
// participants; a malformed element is an error. A nil (empty) list is valid
// and yields an empty, non-nil slice.
func mapCloudChannelParticipants(dtos []tracecore_types.CloudChannelParticipant) ([]channel_domain.Participant, error) {
	if dtos == nil {
		return []channel_domain.Participant{}, nil
	}

	participants := make([]channel_domain.Participant, 0, len(dtos))
	for i, dto := range dtos {
		if dto.ChannelID == "" || dto.VaultID == "" {
			return nil, fmt.Errorf("TracecoreClient - ListParticipants - Cloud returned a malformed participant at index %d: missing channel/vault identity", i)
		}
		participants = append(participants, mapCloudChannelParticipant(dto))
	}
	return participants, nil
}

// mapCloudChannelParticipant is the canonical Cloud -> Desktop Participant
// mapping. Every field exposed by the Cloud Participant aggregate is preserved.
func mapCloudChannelParticipant(dto tracecore_types.CloudChannelParticipant) channel_domain.Participant {
	permissions := dto.Permissions
	if permissions == nil {
		permissions = []string{}
	}

	return channel_domain.Participant{
		ChannelID:   dto.ChannelID,
		VaultID:     dto.VaultID,
		PublicKey:   dto.PublicKey,
		Direction:   dto.Direction,
		JoinedAt:    dto.JoinedAt,
		Role:        dto.Role,
		Permissions: permissions,
	}
}

// InviteToChannel creates a channel invitation through the authoritative Cloud
// backend (POST /channels/{id}/invitations). The request body carries exactly
// the fields the Cloud invitation contract requires; Cloud persists the pending
// invitation and dedupes pending invitations for the same channel + invitee.
// The response carries the persisted Invitation, never a participant.
func (c *TracecoreClient) InviteToChannel(ctx context.Context, req *channel_domain.InviteToChannelRequest) (*tracecore_types.CloudResponse[channel_domain.Invitation], error) {
	if req == nil {
		return nil, channel_domain.ErrRequestRequired
	}
	if req.ChannelID == "" {
		return nil, channel_domain.ErrChannelIDRequired
	}
	if req.InviterVaultID == "" {
		return nil, channel_domain.ErrVaultIDRequired
	}
	if req.InviteeVaultID == "" {
		return nil, channel_domain.ErrVaultIDRequired
	}

	payload := map[string]interface{}{
		"channel_id":       req.ChannelID,
		"inviter_vault_id": req.InviterVaultID,
		"invitee_vault_id": req.InviteeVaultID,
	}

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	baseUrl := c.AnkhoraCloudUrl
	if baseUrl == "" {
		baseUrl = c.BaseURL
	}
	url := baseUrl + "/channels/" + req.ChannelID + "/invitations"

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

	// Surface the Cloud/domain outcome verbatim. The Cloud remains
	// authoritative for the invitation decision.
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	// Cloud returns an envelope: { "status": 201, "data": { ... } }.
	var cloudResp tracecore_types.CloudResponse[tracecore_types.CloudChannelInvitation]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && cloudResp.Data.ID != "" {
		return &tracecore_types.CloudResponse[channel_domain.Invitation]{
			Status:  200,
			Data:    mapCloudChannelInvitation(cloudResp.Data),
			Message: "success",
			Success: true,
		}, nil
	}

	// Tolerate a bare invitation object as a fallback shape.
	var dto tracecore_types.CloudChannelInvitation
	if errDTO := json.Unmarshal(respBytes, &dto); errDTO != nil || dto.ID == "" {
		return nil, fmt.Errorf("TracecoreClient - InviteToChannel - unexpected Cloud response shape (missing invitation data): %s", string(respBytes))
	}

	return &tracecore_types.CloudResponse[channel_domain.Invitation]{
		Status:  200,
		Data:    mapCloudChannelInvitation(dto),
		Message: "success",
		Success: true,
	}, nil
}

// AcceptChannelInvitation accepts a pending channel invitation through the
// authoritative Cloud backend (POST /channels/invitations/{id}/accept). Cloud
// validates the acceptance (the accepting vault must be the invitation's
// invitee) and persists the resulting participant; the accept response carries
// the accepted Invitation, not the participant. Cloud is idempotent: accepting
// an already-accepted invitation returns the accepted invitation without a
// duplicate participant.
func (c *TracecoreClient) AcceptChannelInvitation(ctx context.Context, req *channel_domain.AcceptInvitationRequest) (*tracecore_types.CloudResponse[channel_domain.Invitation], error) {
	if req == nil {
		return nil, channel_domain.ErrRequestRequired
	}
	if req.InvitationID == "" {
		return nil, channel_domain.ErrInvitationIDRequired
	}
	if req.InviteeVaultID == "" {
		return nil, channel_domain.ErrVaultIDRequired
	}
	if req.InviteePublicKey == "" {
		return nil, channel_domain.ErrInviteePublicKeyRequired
	}

	payload := map[string]interface{}{
		"invitation_id":      req.InvitationID,
		"invitee_vault_id":   req.InviteeVaultID,
		"invitee_public_key": req.InviteePublicKey,
	}

	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	baseUrl := c.AnkhoraCloudUrl
	if baseUrl == "" {
		baseUrl = c.BaseURL
	}
	url := baseUrl + "/channels/invitations/" + req.InvitationID + "/accept"

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

	// Surface the Cloud/domain outcome (e.g. "invitation not for you", "record
	// not found") verbatim. The Cloud remains authoritative for the accept.
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	// Cloud returns an envelope: { "status": 200, "data": { ... } }.
	var cloudResp tracecore_types.CloudResponse[tracecore_types.CloudChannelInvitation]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && cloudResp.Data.ID != "" {
		return &tracecore_types.CloudResponse[channel_domain.Invitation]{
			Status:  200,
			Data:    mapCloudChannelInvitation(cloudResp.Data),
			Message: "success",
			Success: true,
		}, nil
	}

	// Tolerate a bare invitation object as a fallback shape.
	var dto tracecore_types.CloudChannelInvitation
	if errDTO := json.Unmarshal(respBytes, &dto); errDTO != nil || dto.ID == "" {
		return nil, fmt.Errorf("TracecoreClient - AcceptChannelInvitation - unexpected Cloud response shape (missing invitation data): %s", string(respBytes))
	}

	return &tracecore_types.CloudResponse[channel_domain.Invitation]{
		Status:  200,
		Data:    mapCloudChannelInvitation(dto),
		Message: "success",
		Success: true,
	}, nil
}

// mapCloudChannelInvitation is the canonical Cloud -> Desktop Invitation
// mapping. Every field exposed by the Cloud Invitation aggregate is preserved.
func mapCloudChannelInvitation(dto tracecore_types.CloudChannelInvitation) channel_domain.Invitation {
	return channel_domain.Invitation{
		ID:             dto.ID,
		ChannelID:      dto.ChannelID,
		InviterVaultID: dto.InviterVaultID,
		InviteeVaultID: dto.InviteeVaultID,
		Status:         channel_domain.InvitationStatus(dto.Status),
		CreatedAt:      dto.CreatedAt,
		AcceptedAt:     dto.AcceptedAt,
	}
}

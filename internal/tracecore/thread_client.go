package tracecore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	tracecore_types "vault-app/internal/tracecore/types"
)

func (c *TracecoreClient) CreateThreadDirect(ctx context.Context, userID string, channelID string, title string, subtitle string, assetType string) (*tracecore_types.ThreadDTO, error) {
	payload := map[string]interface{}{
		"user_id":    userID,
		"channel_id": channelID,
		"thread": map[string]interface{}{
			"title":      title,
			"subtitle":   subtitle,
			"asset_type": assetType,
			"channel_id": channelID,
			"status":     "open",
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	url := c.BaseURL + "/c3/" + userID + "/channel/" + channelID + "/thread/"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var cloudResp tracecore_types.CloudResponse[tracecore_types.ThreadDTO]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && cloudResp.Data.ID != "" {
		return &cloudResp.Data, nil
	}

	var thread tracecore_types.ThreadDTO
	if err := json.Unmarshal(respBytes, &thread); err == nil && thread.ID != "" {
		return &thread, nil
	}

	now := time.Now()
	return &tracecore_types.ThreadDTO{
		ID:        uuid.NewString(),
		ChannelID: channelID,
		AssetType: assetType,
		Title:     title,
		Subtitle:  subtitle,
		Status:    "open",
		CreatedAt: now,
	}, nil
}

func (c *TracecoreClient) ListThreadsDirect(ctx context.Context, userID string, channelID string) ([]tracecore_types.ThreadDTO, error) {
	url := c.BaseURL + "/c3/" + userID + "/channel/" + channelID + "/thread/"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var cloudResp tracecore_types.CloudResponse[[]tracecore_types.ThreadDTO]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && cloudResp.Data != nil {
		return cloudResp.Data, nil
	}

	var threads []tracecore_types.ThreadDTO
	if err := json.Unmarshal(respBytes, &threads); err == nil {
		return threads, nil
	}

	return []tracecore_types.ThreadDTO{}, nil
}

func (c *TracecoreClient) ListThreadEventsDirect(ctx context.Context, userID string, threadID string) ([]tracecore_types.ThreadEventDTO, error) {
	url := c.BaseURL + "/c3/" + userID + "/thread/" + threadID + "/event/"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var cloudResp tracecore_types.CloudResponse[[]tracecore_types.ThreadEventDTO]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && cloudResp.Data != nil {
		return cloudResp.Data, nil
	}

	var events []tracecore_types.ThreadEventDTO
	if err := json.Unmarshal(respBytes, &events); err == nil {
		return events, nil
	}

	return []tracecore_types.ThreadEventDTO{}, nil
}

func (c *TracecoreClient) AppendThreadEventDirect(ctx context.Context, userID string, threadID string, eventType string, payload map[string]interface{}) (*tracecore_types.ThreadEventDTO, error) {
	reqPayload := map[string]interface{}{
		"user_id":   userID,
		"thread_id": threadID,
		"event": map[string]interface{}{
			"type":      eventType,
			"thread_id": threadID,
			"payload":   payload,
		},
	}
	body, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}
	url := c.BaseURL + "/c3/" + userID + "/thread/" + threadID + "/event/"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var cloudResp tracecore_types.CloudResponse[tracecore_types.ThreadEventDTO]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && cloudResp.Data.ID != "" {
		return &cloudResp.Data, nil
	}

	var evt tracecore_types.ThreadEventDTO
	if err := json.Unmarshal(respBytes, &evt); err == nil && evt.ID != "" {
		return &evt, nil
	}

	now := time.Now()
	return &tracecore_types.ThreadEventDTO{
		ID:        uuid.NewString(),
		ThreadID:  threadID,
		Type:      eventType,
		Payload:   payload,
		Cursor:    1,
		CreatedAt: now,
	}, nil
}

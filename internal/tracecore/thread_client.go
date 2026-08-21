package tracecore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	tracecore_types "vault-app/internal/tracecore/types"
)

func (c *TracecoreClient) CreateThreadDirect(ctx context.Context, userID string, channelID string, title string, subtitle string, assetType string) (*tracecore_types.ThreadDTO, error) {
	payload := map[string]interface{}{
		"channel_id":  channelID,
		"identity_id": userID,
		"title":       title,
		"subtitle":    subtitle,
		"asset_type":  assetType,
		"status":      "open",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	baseURL := c.AnkhoraCloudUrl
	if baseURL == "" {
		baseURL = c.BaseURL
	}
	url := baseURL + "/threads"
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

	log.Printf("[CREATE RESPONSE RAW HTTP] status=%d body=%s", resp.StatusCode, string(respBytes))

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var cloudResp tracecore_types.CloudResponse[tracecore_types.ThreadDTO]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && cloudResp.Data.ID != "" {
		log.Printf("[CREATE RESPONSE DECODED] ID=%s ChannelID=%s WorkspaceID=%s Title=%s",
			cloudResp.Data.ID, cloudResp.Data.ChannelID, cloudResp.Data.WorkspaceID, cloudResp.Data.Title)
		return &cloudResp.Data, nil
	}

	var thread tracecore_types.ThreadDTO
	if err := json.Unmarshal(respBytes, &thread); err == nil && thread.ID != "" {
		log.Printf("[CREATE RESPONSE DECODED DIRECT] ID=%s ChannelID=%s WorkspaceID=%s Title=%s",
			thread.ID, thread.ChannelID, thread.WorkspaceID, thread.Title)
		return &thread, nil
	}

	return nil, fmt.Errorf("unexpected response shape from Cloud POST /api/threads: %s", string(respBytes))
}

func (c *TracecoreClient) ListThreadsDirect(
	ctx context.Context,
	userID string,
	channelID string,
) ([]tracecore_types.ThreadDTO, error) {
	baseURL := c.AnkhoraCloudUrl
	if baseURL == "" {
		baseURL = c.BaseURL
	}
	url := baseURL + "/threads/by-channel/" + channelID

	log.Printf("[THREAD LIST] AnkhoraCloudUrl = %s", c.AnkhoraCloudUrl)
	log.Printf("[THREAD LIST] BaseURL = %s", c.BaseURL)
	log.Printf("[THREAD LIST] FINAL URL = %s", url)
	log.Printf("[THREAD LIST] channelID parameter = %s", channelID)
	log.Printf("[THREAD LIST] Authorization present = %v", c.Token != "")

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
		log.Printf("[THREAD LIST HTTP ERROR] %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[THREAD LIST HTTP STATUS] = %d", resp.StatusCode)
	log.Printf("[THREAD LIST RAW RESPONSE] = %s", string(respBytes))

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf(
			"Cloud backend returned status %d: %s",
			resp.StatusCode,
			string(respBytes),
		)
	}

	var cloudResp tracecore_types.CloudResponse[[]tracecore_types.ThreadDTO]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && cloudResp.Data != nil {
		log.Printf("[THREAD LIST DECODED COUNT] = %d", len(cloudResp.Data))
		for idx, th := range cloudResp.Data {
			log.Printf("[THREAD LIST DECODED DTO #%d] id=%s channel_id=%s title=%s", idx, th.ID, th.ChannelID, th.Title)
		}
		return cloudResp.Data, nil
	}

	var threads []tracecore_types.ThreadDTO
	if err := json.Unmarshal(respBytes, &threads); err == nil {
		log.Printf("[THREAD LIST DECODED COUNT DIRECT] = %d", len(threads))
		for idx, th := range threads {
			log.Printf("[THREAD LIST DECODED DTO DIRECT #%d] id=%s channel_id=%s title=%s", idx, th.ID, th.ChannelID, th.Title)
		}
		return threads, nil
	}

	return []tracecore_types.ThreadDTO{}, nil
}

func (c *TracecoreClient) ListThreadEventsDirect(ctx context.Context, userID string, threadID string) ([]tracecore_types.ThreadEventDTO, error) {
	url := c.AnkhoraCloudUrl + "/threads/" + threadID + "/events"
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

func (c *TracecoreClient) AppendThreadEventDirect(ctx context.Context, userID string, threadID string, eventType string, payload map[string]interface{}, idempotencyKey string) (*tracecore_types.ThreadEventDTO, error) {
	reqPayload := map[string]interface{}{
		"type":      eventType,
		"thread_id": threadID,
		"payload":   payload,
	}
	if idempotencyKey != "" {
		reqPayload["idempotency_key"] = idempotencyKey
	}
	body, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}



	url := c.AnkhoraCloudUrl + "/threads/" + threadID + "/events"
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

	return nil, fmt.Errorf("unexpected response shape from Cloud POST /api/threads/%s/events: %s", threadID, string(respBytes))
}

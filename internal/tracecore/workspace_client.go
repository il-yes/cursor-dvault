package tracecore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	tracecore_types "vault-app/internal/tracecore/types"
	"vault-app/internal/utils"
	workspace_domain "vault-app/internal/workspace/domain"
)

func (c *TracecoreClient) CreateWorkspace(ctx context.Context, req workspace_domain.CreateRequest) (*tracecore_types.CloudResponse[workspace_domain.Workspace], error) {

	payload := tracecore_types.NewCreateWorkspaceRequest{
		VaultID: req.VaultID,
		// Workspace: tracecore_types.Workspace{
		// 	ID:          req.Workspace.ID,
		// 	VaultID:     req.Workspace.VaultID,
		// 	Name:        req.Workspace.Name,
		// 	Description: req.Workspace.Description,
		// 	Status:      string(req.Workspace.Status),
		// 	OwnerID:     req.Workspace.OwnerID,
		// 	CreatedAt:   req.Workspace.CreatedAt,
		// 	UpdatedAt:   req.Workspace.UpdatedAt,
		// },
		// Signature: req.Signature,
		Name:        req.Workspace.Name,
		Description: req.Workspace.Description,
		OwnerID:     req.Workspace.OwnerID,
	}

	// Step 1: build JSON body
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	// Step 2: build URL and request
	url := c.AnkhoraCloudUrl + "/workspaces"
	utils.LogPretty("TracecoreClient - CreateWorkspace - URL", url)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}

	utils.LogPretty("TracecoreClient - CreateWorkspace - URL", url)
	// utils.LogPretty("TracecoreClient - CreateWorkspace - request body", body.String())

	// Step 3: do the request
	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("TracecoreClient - CreateWorkspace - read body failed: %v", err)
		return nil, fmt.Errorf("read body failed: %w", err)
	}
	utils.LogPretty("TracecoreClient - CreateWorkspace - raw body (string)", string(respBytes))

	// if body is empty or very short, log raw bytes
	if len(respBytes) == 0 {
		utils.LogPretty("TracecoreClient - CreateWorkspace - raw body", "(empty)")
	} else {
		utils.LogPretty("TracecoreClient - CreateWorkspace - raw body (hex)", fmt.Sprintf("%x", respBytes))
	}

	if len(respBytes) == 0 {
		return nil, fmt.Errorf("TracecoreClient - CreateWorkspace - empty body")
	}

	var cloudResp tracecore_types.CloudResponse[workspace_domain.Workspace]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		utils.LogPretty("🚫 TracecoreClient - CreateWorkspace - json.Unmarshal error", err)
		utils.LogPretty("🚫 TracecoreClient - CreateWorkspace - raw body when unmarshal failed", string(respBytes))
		return nil, fmt.Errorf("TracecoreClient - CreateWorkspace - cloud response unmarshal failed: %w", err)
	}

	utils.LogPretty("✅ TracecoreClient - CreateWorkspace - cloudResp", cloudResp)
	return &cloudResp, nil
}

func (c *TracecoreClient) CreateWorkspaceDirect(ctx context.Context, vaultID string, userID string, name string, description string) (*tracecore_types.Workspace, error) {
	payload := map[string]interface{}{
		// "user_id": userID,
		// "workspace": map[string]interface{}{
		// 	"name":        name,
		// 	"description": description,
		// 	"vault_id":    "active",
		// 	"owner_id":    userID,
		// },
		"VaultID":     vaultID,
		"Name":        name,
		"Description": description,
		"OwnerID":     userID,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	url := c.BaseURL + "/workspaces/"
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

	var cloudResp tracecore_types.CloudResponse[tracecore_types.Workspace]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && cloudResp.Data.ID != "" {
		return &cloudResp.Data, nil
	}

	var ws tracecore_types.Workspace
	if err := json.Unmarshal(respBytes, &ws); err == nil && ws.ID != "" {
		return &ws, nil
	}

	now := time.Now()
	return &tracecore_types.Workspace{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		Status:      "active",
		OwnerID:     userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (c *TracecoreClient) ListWorkspaces(ctx context.Context, vaultID string) ([]tracecore_types.Workspace, error) {
	url := c.AnkhoraCloudUrl + "/workspaces?vault_id=" + vaultID
	utils.LogPretty("[Workspace] Cloud GET URL", url)

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
		utils.LogPretty("🚫 [Workspace] TracecoreClient.ListWorkspaces HTTP Do error", err)
		return nil, err
	}
	defer resp.Body.Close()

	utils.LogPretty("[Workspace] Cloud HTTP Status", resp.StatusCode)

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("[Workspace] Cloud Raw Response Body", string(respBytes))

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	var cloudResp tracecore_types.CloudResponse[[]tracecore_types.Workspace]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && cloudResp.Data != nil {
		utils.LogPretty("[Workspace] Decoded cloudResp.Data", cloudResp.Data)
		return cloudResp.Data, nil
	}

	var workspaces []tracecore_types.Workspace
	if err := json.Unmarshal(respBytes, &workspaces); err == nil {
		utils.LogPretty("[Workspace] Decoded workspaces array", workspaces)
		return workspaces, nil
	}
	utils.LogPretty("[Workspace] TracecoreClient - ListWorkspaces empty fallback", workspaces)
	return []tracecore_types.Workspace{}, nil
}

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

func (c *TracecoreClient) CreateCollaborativeShareDirect(ctx context.Context, userID string, threadID string, trustGroupID string, assetCID string, targetVaultID string, notes string) (*tracecore_types.ShareEntryRefDTO, error) {
	shareID := "se_" + uuid.NewString()[:12]
	nowStr := time.Now().Format(time.RFC3339)

	shareRef := tracecore_types.ShareEntryRefDTO{
		ShareEntryID: shareID,
		TrustGroupID: trustGroupID,
		AssetCID:     assetCID,
		CreatedBy:    userID,
		Status:       "active",
		CreatedAt:    nowStr,
	}

	payload := map[string]interface{}{
		"notes":           notes,
		"target_vault_id": targetVaultID,
		"share_entry_ref": shareRef,
	}

	_, err := c.AppendThreadEventDirect(ctx, userID, threadID, "share.created", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to record collaborative share thread event: %w", err)
	}

	return &shareRef, nil
}

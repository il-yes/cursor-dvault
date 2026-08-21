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

	var cloudResp tracecore_types.CloudResponse[tracecore_types.CloudWorkspaceDTO]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		utils.LogPretty("🚫 TracecoreClient - CreateWorkspace - json.Unmarshal error", err)
		utils.LogPretty("🚫 TracecoreClient - CreateWorkspace - raw body when unmarshal failed", string(respBytes))
		return nil, fmt.Errorf("TracecoreClient - CreateWorkspace - cloud response unmarshal failed: %w", err)
	}

	dto := cloudResp.Data
	domainWs := workspace_domain.Workspace{
		ID:          dto.ID,
		VaultID:     dto.VaultID,
		Name:        dto.Name,
		Description: dto.Description,
		Status:      workspace_domain.WorkspaceStatus(dto.Status),
		OwnerID:     dto.OwnerID,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
		IsDraft:     dto.IsDraft,
		IsDirty:     dto.IsDirty,
	}

	res := &tracecore_types.CloudResponse[workspace_domain.Workspace]{
		Status:  cloudResp.Status,
		Data:    domainWs,
		Message: cloudResp.Message,
		Success: cloudResp.Success,
	}

	utils.LogPretty("✅ TracecoreClient - CreateWorkspace - cloudResp", res)
	return res, nil
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

	var cloudResp tracecore_types.CloudResponse[tracecore_types.CloudWorkspaceDTO]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && cloudResp.Data.ID != "" {
		if mapped := tracecore_types.MapCloudWorkspaceToTypes(&cloudResp.Data); mapped != nil {
			return mapped, nil
		}
	}

	var dto tracecore_types.CloudWorkspaceDTO
	if err := json.Unmarshal(respBytes, &dto); err == nil && dto.ID != "" {
		if mapped := tracecore_types.MapCloudWorkspaceToTypes(&dto); mapped != nil {
			return mapped, nil
		}
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
	log.Printf("[CLOUD-TRACE] LEDGER REQUEST: client_pointer=%p token_length=%d token_fingerprint=%s authorization_header=%v vault_id=%s",
		c, len(c.Token), traceTokenFingerprint(c.Token), c.Token != "", vaultID)

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

	var cloudResp tracecore_types.CloudResponse[[]tracecore_types.CloudWorkspaceDTO]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && cloudResp.Data != nil {
		utils.LogPretty("[Workspace] Decoded cloudResp.Data DTOs", cloudResp.Data)
		workspaces := make([]tracecore_types.Workspace, 0, len(cloudResp.Data))
		for _, dto := range cloudResp.Data {
			if mapped := tracecore_types.MapCloudWorkspaceToTypes(&dto); mapped != nil {
				workspaces = append(workspaces, *mapped)
			}
		}
		return workspaces, nil
	}

	var dtos []tracecore_types.CloudWorkspaceDTO
	if err := json.Unmarshal(respBytes, &dtos); err == nil && len(dtos) > 0 {
		utils.LogPretty("[Workspace] Decoded dtos array", dtos)
		workspaces := make([]tracecore_types.Workspace, 0, len(dtos))
		for _, dto := range dtos {
			if mapped := tracecore_types.MapCloudWorkspaceToTypes(&dto); mapped != nil {
				workspaces = append(workspaces, *mapped)
			}
		}
		return workspaces, nil
	}
	utils.LogPretty("[Workspace] TracecoreClient - ListWorkspaces empty fallback", nil)
	return []tracecore_types.Workspace{}, nil
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
		"ref_type":       "share_entry",
		"share_entry_id": shareRef.ShareEntryID,
		"trust_group_id": shareRef.TrustGroupID,
	}

	_, err := c.AppendThreadEventDirect(ctx, userID, threadID, "entry.shared", payload, "")
	if err != nil {
		return nil, fmt.Errorf("failed to record collaborative share thread event: %w", err)
	}

	return &shareRef, nil
}

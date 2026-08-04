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
	"vault-app/internal/utils"
	workspace_domain "vault-app/internal/workspace/domain"
)

func (c *TracecoreClient) CreateWorkspace(ctx context.Context, req workspace_domain.CreateRequest) (*tracecore_types.CloudResponse[workspace_domain.Workspace], error) {

	payload := tracecore_types.NewCreateWorkspaceRequest{
		UserID:    req.UserID,
		VaultID:   req.VaultID,
		Workspace: tracecore_types.Workspace{
			ID:      req.Workspace.ID,
			VaultID: req.Workspace.VaultID,
			Name:        req.Workspace.Name,
			Description: req.Workspace.Description,
			Status:      string(req.Workspace.Status),
			OwnerID: req.Workspace.OwnerID,
			CreatedAt: req.Workspace.CreatedAt,
			UpdatedAt: req.Workspace.UpdatedAt,
		},
		Signature: req.Signature,
	}

	// Step 1: build JSON body
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return nil, err
	}

	// Step 2: build URL and request
	url := c.BaseURL + "/c3/" + payload.UserID + "/workspace/"
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

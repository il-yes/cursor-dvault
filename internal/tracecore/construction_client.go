package tracecore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	tracecore_types "vault-app/internal/tracecore/types"
)

// GetProjectOverview fetches the construction ProjectOverview read model from Cloud
// via GET /api/construction/projects/{id}.
func (c *TracecoreClient) GetProjectOverview(ctx context.Context, projectID string) (*tracecore_types.ProjectOverviewDTO, error) {
	if projectID == "" {
		return nil, errors.New("projectID is required")
	}

	baseUrl := strings.TrimRight(c.AnkhoraCloudUrl, "/")
	if baseUrl == "" {
		baseUrl = strings.TrimRight(c.BaseURL, "/")
	}

	cleanPath := "/construction/projects/" + projectID
	if !strings.HasSuffix(baseUrl, "/api") {
		cleanPath = "/api" + cleanPath
	}

	targetURL := baseUrl + cleanPath

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
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
		return nil, fmt.Errorf("read body failed: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Cloud backend returned status %d: %s", resp.StatusCode, string(respBytes))
	}

	// 1. Try CloudResponse envelope: { "status": 200, "data": { ... } }
	var cloudResp tracecore_types.CloudResponse[tracecore_types.ProjectOverviewDTO]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && (cloudResp.Data.ID != "" || cloudResp.Data.ProjectID != "" || cloudResp.Data.Name != "" || cloudResp.Data.ProjectName != "") {
		dto := normalizeProjectOverview(cloudResp.Data)
		return &dto, nil
	}

	// 2. Try bare ProjectOverviewDTO object
	var dto tracecore_types.ProjectOverviewDTO
	if errDTO := json.Unmarshal(respBytes, &dto); errDTO == nil && (dto.ID != "" || dto.ProjectID != "" || dto.Name != "" || dto.ProjectName != "") {
		dto = normalizeProjectOverview(dto)
		return &dto, nil
	}

	return nil, fmt.Errorf("TracecoreClient - GetProjectOverview - unexpected Cloud response shape: %s", string(respBytes))
}

func normalizeProjectOverview(dto tracecore_types.ProjectOverviewDTO) tracecore_types.ProjectOverviewDTO {
	if dto.ID == "" && dto.ProjectID != "" {
		dto.ID = dto.ProjectID
	}
	if dto.ProjectID == "" && dto.ID != "" {
		dto.ProjectID = dto.ID
	}
	if dto.Code == "" && dto.ProjectReference != "" {
		dto.Code = dto.ProjectReference
	}
	if dto.ProjectReference == "" && dto.Code != "" {
		dto.ProjectReference = dto.Code
	}
	if dto.Name == "" && dto.ProjectName != "" {
		dto.Name = dto.ProjectName
	}
	if dto.ProjectName == "" && dto.Name != "" {
		dto.ProjectName = dto.Name
	}
	if dto.Type == "" && dto.ProjectType != "" {
		dto.Type = dto.ProjectType
	}
	if dto.ProjectType == "" && dto.Type != "" {
		dto.ProjectType = dto.Type
	}
	return dto
}

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

// GetProcurementOverview fetches the construction procurement/requirement read model from Cloud
// via GET /api/construction/requirements/{id} or /api/construction/procurement/{id}.
func (c *TracecoreClient) GetProcurementOverview(ctx context.Context, requirementID string) (*tracecore_types.ProcurementOverviewDTO, error) {
	if requirementID == "" {
		return nil, errors.New("requirementID is required")
	}

	baseUrl := strings.TrimRight(c.AnkhoraCloudUrl, "/")
	if baseUrl == "" {
		baseUrl = strings.TrimRight(c.BaseURL, "/")
	}

	cleanPath := "/construction/requirements/" + requirementID
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
	var cloudResp tracecore_types.CloudResponse[tracecore_types.ProcurementOverviewDTO]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && (cloudResp.Data.ID != "" || cloudResp.Data.RequirementID != "" || cloudResp.Data.Code != "") {
		dto := normalizeProcurementOverview(cloudResp.Data)
		return &dto, nil
	}

	// 2. Try bare ProcurementOverviewDTO object
	var dto tracecore_types.ProcurementOverviewDTO
	if errDTO := json.Unmarshal(respBytes, &dto); errDTO == nil && (dto.ID != "" || dto.RequirementID != "" || dto.Code != "") {
		dto = normalizeProcurementOverview(dto)
		return &dto, nil
	}

	return nil, fmt.Errorf("TracecoreClient - GetProcurementOverview - unexpected Cloud response shape: %s", string(respBytes))
}

func normalizeProcurementOverview(dto tracecore_types.ProcurementOverviewDTO) tracecore_types.ProcurementOverviewDTO {
	if dto.ID == "" && dto.RequirementID != "" {
		dto.ID = dto.RequirementID
	}
	if dto.RequirementID == "" && dto.ID != "" {
		dto.RequirementID = dto.ID
	}
	if dto.Code == "" && dto.RequirementID != "" {
		dto.Code = dto.RequirementID
	}
	return dto
}

// GetLogisticsOverview fetches the construction logistics/delivery read model from Cloud
// via GET /api/construction/deliveries/{id}/overview.
func (c *TracecoreClient) GetLogisticsOverview(ctx context.Context, deliveryID string) (*tracecore_types.LogisticsOverviewDTO, error) {
	if deliveryID == "" {
		return nil, errors.New("deliveryID is required")
	}

	baseUrl := strings.TrimRight(c.AnkhoraCloudUrl, "/")
	if baseUrl == "" {
		baseUrl = strings.TrimRight(c.BaseURL, "/")
	}

	cleanPath := "/construction/deliveries/" + deliveryID + "/overview"
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
	var cloudResp tracecore_types.CloudResponse[tracecore_types.LogisticsOverviewDTO]
	if err := json.Unmarshal(respBytes, &cloudResp); err == nil && (cloudResp.Data.Delivery.ID != "" || cloudResp.Data.Delivery.Reference != "") {
		dto := normalizeLogisticsOverview(cloudResp.Data)
		return &dto, nil
	}

	// 2. Try bare LogisticsOverviewDTO object
	var dto tracecore_types.LogisticsOverviewDTO
	if errDTO := json.Unmarshal(respBytes, &dto); errDTO == nil && (dto.Delivery.ID != "" || dto.Delivery.Reference != "") {
		dto = normalizeLogisticsOverview(dto)
		return &dto, nil
	}

	return nil, fmt.Errorf("TracecoreClient - GetLogisticsOverview - unexpected Cloud response shape: %s", string(respBytes))
}

func normalizeLogisticsOverview(dto tracecore_types.LogisticsOverviewDTO) tracecore_types.LogisticsOverviewDTO {
	if dto.Delivery.ID == "" && dto.Delivery.Reference != "" {
		dto.Delivery.ID = dto.Delivery.Reference
	}
	if dto.Delivery.Reference == "" && dto.Delivery.ID != "" {
		dto.Delivery.Reference = dto.Delivery.ID
	}
	return dto
}



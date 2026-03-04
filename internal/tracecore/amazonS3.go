package tracecore

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	tracecore_types "vault-app/internal/tracecore/types"
	"vault-app/internal/utils"
)


func (c *TracecoreClient) AddToS3(ctx context.Context, req tracecore_types.SyncVaultStreamRequest) (*tracecore_types.CloudResponse[tracecore_types.SyncVaultResponse], error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/vaults/"+req.UserID+"/storage/"+req.VaultName, nil)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}
	request.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	respBytes, _ := io.ReadAll(resp.Body)
	var cloudResp tracecore_types.CloudResponse[tracecore_types.SyncVaultResponse]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		log.Printf("invalid cloud response: %w", err)
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}
	utils.LogPretty("AddToS3 - cloud response", cloudResp)	
	
	return &cloudResp, nil
}
func (c *TracecoreClient) GetDataFromS3(ctx context.Context, req tracecore_types.IpfsCidRequest) (*tracecore_types.CloudResponse[tracecore_types.IpfsCidResponse], error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/vaults/"+req.UserID+"/storage/"+req.VaultName+"/"+req.CID, nil)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}
	request.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	respBytes, _ := io.ReadAll(resp.Body)
	var cloudResp tracecore_types.CloudResponse[tracecore_types.IpfsCidResponse]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		log.Printf("invalid cloud response: %w", err)
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}
	utils.LogPretty("GetDataFromS3 - cloud response", cloudResp)	
	
	return &cloudResp, nil
}

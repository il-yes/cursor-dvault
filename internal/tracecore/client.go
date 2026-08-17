package tracecore

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	app_config "vault-app/internal/config"
	app_config_domain "vault-app/internal/config/domain"
	tracecore_models "vault-app/internal/tracecore/models"
)

// Error definitions
var (
	ErrUserNotFound = errors.New("user not found")
	ErrNotFound     = errors.New("endpoint not found")
)

// traceTokenFingerprint returns a safe fingerprint for a token (SHA-256 first 8 hex chars)
func traceTokenFingerprint(token string) string {
	if token == "" {
		return "absent"
	}
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h[:4])
}

// TraceTokenFingerprint is the exported version for use across packages
func TraceTokenFingerprint(token string) string {
	return traceTokenFingerprint(token)
}

type TracecoreClient struct {
	BaseURL         string
	Token           string
	HTTPClient      *http.Client
	AnkhoraFrontUrl string
	AnkhoraCloudUrl string
}

// NewTracecoreClient creates a new Tracecore client with default timeout.
func NewTracecoreClient(baseURL, token, ankhoraFrontUrl, ankhoraCloudUrl string) *TracecoreClient {
	if baseURL == "" {
		log.Println("⚠️ Cloud URL is empty — Cloud features will be disabled")
	}
	if token == "" {
		log.Println("ℹ️ Cloud token is empty — will be set after authentication")
	}

	return &TracecoreClient{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		AnkhoraFrontUrl: ankhoraFrontUrl,
		AnkhoraCloudUrl: ankhoraCloudUrl,
	}
}

// NewTracecoreClient creates a new Tracecore client with default timeout.
func NewTracecoreFromConfig(appCfg *app_config_domain.AppConfig, token string) *TracecoreClient {
	// ⚠️ Don't use log.Fatal during startup - it kills the app!
	if token == "" {
		log.Println("⚠️ TRACECORE_TOKEN is empty — Tracecore features will be disabled")
	}

	switch appCfg.Storage.Mode {
	case app_config.StorageCloud:
		return &TracecoreClient{
			BaseURL: appCfg.Storage.Cloud.BaseURL,
			Token:   token,
			HTTPClient: &http.Client{
				Timeout: 30 * time.Second,
			},
		}
	case app_config.StorageLocal:
		return &TracecoreClient{
			BaseURL: appCfg.Storage.LocalIPFS.APIEndpoint,
			Token:   token,
			HTTPClient: &http.Client{
				Timeout: 30 * time.Second,
			},
		}
	case app_config.StoragePrivateIPFS:
		return &TracecoreClient{
			BaseURL: appCfg.Storage.PrivateIPFS.APIEndpoint,
			Token:   token,
			HTTPClient: &http.Client{
				Timeout: 30 * time.Second,
			},
		}
	default:
		log.Println("⚠️ TRACECORE_STORAGE_MODE is unknown — Tracecore features will be disabled", appCfg.Storage.Mode)
		return &TracecoreClient{
			BaseURL: appCfg.Storage.Cloud.BaseURL, // "https://ankhora.io/back", // appCfg.Storage.Cloud.BaseURL,
			Token:   token,
			HTTPClient: &http.Client{
				Timeout: 30 * time.Second,
			},
		}
	}
}
func (c *TracecoreClient) SetToken(token string) {
	before := traceTokenFingerprint(c.Token)
	c.Token = token
	after := traceTokenFingerprint(c.Token)
	log.Printf("[CLOUD-TRACE] SetToken: before=%s after=%s length=%d", before, after, len(token))
}

func (c *TracecoreClient) doRequest(ctx context.Context, method, path string, body any, out any) error {
	authHeader := c.Token != ""
	log.Printf("[CLOUD-TRACE] F: doRequest: method=%s path=%s token_fingerprint=%s token_length=%d authorization_header=%v",
		method, path, traceTokenFingerprint(c.Token), len(c.Token), authHeader)

	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		buf = bytes.NewBuffer(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		// Return specific error for 404 Not Found
		if resp.StatusCode == http.StatusNotFound {
			return ErrUserNotFound
		}
		return fmt.Errorf("server error %d: %s", resp.StatusCode, string(b))
	}
	raw, _ := io.ReadAll(resp.Body)
	fmt.Println("RAW RESPONSE:", string(raw))
	resp.Body = io.NopCloser(bytes.NewBuffer(raw))

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

func (tc *TracecoreClient) Commit(payload tracecore_models.CommitEnvelope) (*tracecore_models.CommitResponse, error) {
	if tc == nil {
		return nil, fmt.Errorf("TracecoreClient is nil")
	}
	if tc.HTTPClient == nil {
		return nil, fmt.Errorf("HTTPClient is not initialized")
	}
	if tc.BaseURL == "" {
		return nil, fmt.Errorf("BaseURL is empty")
	}

	url := fmt.Sprintf("%s/d-vault/vaults", tc.BaseURL)

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode commit payload: %w", err)
	}
	// utils.LogPretty("commit enveloppe", data)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create commit request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tc.Token)

	var resp *http.Response
	for attempts := 0; attempts < 3; attempts++ {
		resp, err = tc.HTTPClient.Do(req)
		if err == nil && (resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
			break
		}
		time.Sleep(time.Duration(attempts+1) * 500 * time.Millisecond) // Backoff
	}

	// resp, err := tc.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("❌ Tracecore HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		// utils.LogPretty("Tracecore raw response", body)
		return nil, fmt.Errorf("❌ Tracecore returned status %d: %s", resp.StatusCode, body)
	}
	// utils.LogPretty("tracecore response", resp.Body)
	var commitResp tracecore_models.CommitResponse
	if err := json.NewDecoder(resp.Body).Decode(&commitResp); err != nil {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("❌ failed to decode Tracecore response: %w\nRaw body: %s", err, body)
	}
	// utils.LogPretty("commitResp", &commitResp)
	return &commitResp, nil

}

func (tc *TracecoreClient) CreateRepo() (*string, error) {
	return nil, fmt.Errorf("Methid Not implemented")
}

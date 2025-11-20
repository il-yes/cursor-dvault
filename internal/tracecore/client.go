package tracecore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Error definitions
var (
	ErrUserNotFound = errors.New("user not found")
)
var (
	ErrNotFound = errors.New("endpoint not found")
)

type TracecoreClient struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewTracecoreClient creates a new Tracecore client with default timeout.
func NewTracecoreClient(baseURL, token string) *TracecoreClient {
	// ⚠️ Don't use log.Fatal during startup - it kills the app!
	// Just log warnings and allow the app to start
	if baseURL == "" {
		log.Println("⚠️ TRACECORE_URL is empty — Tracecore features will be disabled")
	}
	if token == "" {
		log.Println("⚠️ TRACECORE_TOKEN is empty — Tracecore features will be disabled")
	}

	return &TracecoreClient{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *TracecoreClient) doRequest(ctx context.Context, method, path string, body any, out any) error {
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

func (tc *TracecoreClient) Commit(payload CommitEnvelope) (*CommitResponse, error) {
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
	var commitResp CommitResponse
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

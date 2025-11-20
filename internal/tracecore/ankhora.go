package tracecore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (c *TracecoreClient) GetVault(ctx context.Context) ([]byte, error) {
	var resp struct {
		Data []byte `json:"data"`
	}
	err := c.doRequest(ctx, "GET", "/api/vault", nil, &resp)
	return resp.Data, err
}

func (c *TracecoreClient) SaveVault(ctx context.Context, encryptedVault []byte) error {
	body := map[string][]byte{"data": encryptedVault}
	return c.doRequest(ctx, "POST", "/api/vault", body, nil)
}

type WrappedResponse[T any] struct {
	Result T   `json:"result"`
	Error  any `json:"error"`
}

// Ensure these imports exist:
//   "bytes"
//   "context"
//   "encoding/json"
//   "errors"
//   "fmt"
//   "io"
//   "net/http"
//   "net/url"
//   "strings"


type wrappedResultUser struct {
	Result *User       `json:"result"`
	Error  interface{} `json:"error"`
}

func (c *TracecoreClient) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// Build endpoint (ensure no double slashes)
	base := strings.TrimRight(c.BaseURL, "/")
	endpoint := fmt.Sprintf("%s/check-email?email=%s", base, url.QueryEscape(email))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle HTTP status codes
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrUserNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(b))
	}

	// Read full body (so we can try multiple unmarshals)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	// quick empty check
	if len(bytes.TrimSpace(bodyBytes)) == 0 {
		return nil, ErrUserNotFound
	}

	// 1) Try to decode directly into User
	var user User
	if err := json.Unmarshal(bodyBytes, &user); err == nil {
		// If ID or email present, assume success
		if user.ID != 0 || user.Email != "" {
			return &user, nil
		}
	}

	// 2) Try to decode into Wails-style wrapper { result: ..., error: ... }
	var wrapped wrappedResultUser
	if err := json.Unmarshal(bodyBytes, &wrapped); err == nil && wrapped.Result != nil {
		if wrapped.Result.ID != 0 || wrapped.Result.Email != "" {
			return wrapped.Result, nil
		}
	}

	// Nothing worked — return not found (or decoding error)
	return nil, ErrUserNotFound
}


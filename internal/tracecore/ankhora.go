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
	"time"
	utils "vault-app/internal"
	share_domain "vault-app/internal/domain/shared"
	subscription_domain "vault-app/internal/subscription/domain"
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

type CloudShareRequest struct {
	EntryName     string                     `json:"entry_name"`
	EntryType     string                     `json:"entry_type"`
	EntryRef      string                     `json:"entry_ref"`
	Status        string                     `json:"status"`
	AccessMode    string                     `json:"access_mode"`
	Encryption    string                     `json:"encryption"`
	EntrySnapshot share_domain.EntrySnapshot `json:"entry_snapshot"` // pass as object
	ExpiresAt     *time.Time                 `json:"expires_at,omitempty"`
	Recipients    []CloudRecipient           `json:"recipients"`
	SharedAt      time.Time                  `json:"shared_at"`
}
type CloudRecipient struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

// implement the interface

// CreateShare sends a share entry to the cloud backend
func (c *TracecoreClient) CreateShare(ctx context.Context, s share_domain.ShareEntry) (*share_domain.ShareEntry, error) {
	// Ensure SharedAt is set
	if s.SharedAt.IsZero() {
		s.SharedAt = time.Now().UTC()
	}

	// Convert recipients to cloud format
	recipients := make([]CloudRecipient, len(s.Recipients))
	for i, r := range s.Recipients {
		recipients[i] = CloudRecipient{
			ID:       r.Email, // cloud uses email as ID
			Name:     r.Name,
			Email:    r.Email,
			Role:     r.Role,
			JoinedAt: time.Now().UTC(), // pass time.Time, NOT string
		}
	}

	// Build cloud payload
	payload := CloudShareRequest{
		EntryName:     s.EntryName,
		EntryType:     s.EntryType,
		EntryRef:      s.EntryRef,
		Status:        s.Status,
		AccessMode:    s.AccessMode,
		Encryption:    s.Encryption,
		EntrySnapshot: s.EntrySnapshot,
		ExpiresAt:     s.ExpiresAt, // *time.Time, nullable
		Recipients:    recipients,
		SharedAt:      s.SharedAt,
	}

	
    bodyBytes, _ := json.Marshal(payload)
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/shares", bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    if c.Token != "" {
        req.Header.Set("Authorization", "Bearer "+c.Token)
    }

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    respBytes, _ := io.ReadAll(resp.Body)
    var cloudResp struct {
        Data       json.RawMessage `json:"data"`
        Status     string          `json:"status"`
        StatusCode int             `json:"status_code"`
        Message    string          `json:"message"`
    }
    if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
        return nil, fmt.Errorf("invalid cloud response: %w", err)
    }

    if cloudResp.Status != "ok" {
        return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
    }
	utils.LogPretty("cloud response", cloudResp)

    // unmarshal the data to domain struct
    var created share_domain.ShareEntry
    if err := json.Unmarshal(cloudResp.Data, &created); err != nil {
        return nil, fmt.Errorf("invalid cloud data: %w", err)
    }

		utils.LogPretty("created", &created)	

	return &created, nil
}
func (c *TracecoreClient) GetShareByMe(ctx context.Context) ([]share_domain.ShareEntry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/shares/by-me", nil)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("GetShareByMe - c.Token", c.Token)
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)

	var cloudResp struct {
		Data       json.RawMessage `json:"data"`
		Status     string          `json:"status"`
		StatusCode int             `json:"status_code"`
		Message    string          `json:"message"`
	}

	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}
	utils.LogPretty("cloud response", cloudResp)

	if cloudResp.Status != "ok" {
		return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}

	var list []share_domain.ShareEntry
	if err := json.Unmarshal(cloudResp.Data, &list); err != nil {
		return nil, fmt.Errorf("invalid cloud data: %w", err)
	}

	return list, nil
}
func (c *TracecoreClient) GetShareWithMe(ctx context.Context) ([]share_domain.ShareEntry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/shares/with-me", nil)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)

	var cloudResp struct {
		Data       json.RawMessage `json:"data"`
		Status     string          `json:"status"`
		StatusCode int             `json:"status_code"`
		Message    string          `json:"message"`
	}

	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}

	if cloudResp.Status != "ok" {
		return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}

	var list []share_domain.ShareEntry
	if err := json.Unmarshal(cloudResp.Data, &list); err != nil {
		return nil, fmt.Errorf("invalid cloud data: %w", err)
	}

	return list, nil
}


type PaymentSetupRequest struct {
	PaymentIntentID       string                               `json:"payment_intent_id"`
	Rail                  string                               `json:"rail"`
	Wallet                string                               `json:"wallet"`	
	Month                 int64                                `json:"month"`
	TxHash                string                               `json:"tx_hash"`	
	Email                 string 
	FirstName             string                               `json:"first_name"`
	LastName              string                               `json:"last_name"`
	UserID                string                               `json:"user_id"`
	Tier                  subscription_domain.SubscriptionTier `json:"tier"`
	LastFour              string                               `json:"last_four"`
	CardNumber            string                               `json:"card_number"`
	Exp                   string                               `json:"exp"`
	CVC                   string                               `json:"cvc"`
	PaymentMethod         subscription_domain.PaymentMethod    `json:"payment_method"`
	StripePaymentMethodID string                               `json:"stripe_payment_method_id,omitempty"`
	EncryptedPaymentData  string                               `json:"encrypted_payment_data,omitempty"` // Encrypted client-side
	StellarPublicKey      string                               `json:"stellar_public_key,omitempty"`
}
type PaymentSetupResponse struct {
	Data       json.RawMessage `json:"data"`
	Status     string          `json:"status"`
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
}
func (c *TracecoreClient) SetupSubscription(ctx context.Context, payload PaymentSetupRequest) (*PaymentSetupResponse, error) {
	bodyBytes, _ := json.Marshal(payload)
    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/subscriptions/stripe", bytes.NewReader(bodyBytes))
    req.Header.Set("Content-Type", "application/json")
    if c.Token != "" {
        req.Header.Set("Authorization", "Bearer "+c.Token)
    }

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    respBytes, _ := io.ReadAll(resp.Body)
    var cloudResp PaymentSetupResponse
    if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
        return nil, fmt.Errorf("invalid cloud response: %w", err)
    }

    if cloudResp.Status != "ok" {
        return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
    }
	utils.LogPretty("cloud response", cloudResp)

    return &cloudResp, nil
}	
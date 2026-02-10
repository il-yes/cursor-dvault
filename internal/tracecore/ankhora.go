package tracecore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	utils "vault-app/internal"
	share_application_dto "vault-app/internal/application"
	billing_domain "vault-app/internal/billing/domain"
	share_domain "vault-app/internal/domain/shared"
	subscription_domain "vault-app/internal/subscription/domain"

	"gorm.io/datatypes"
)

func (c *TracecoreClient) GetVault(ctx context.Context) ([]byte, error) {
	var resp struct {
		Data []byte `json:"data"`
	}
	err := c.doRequest(ctx, "GET", "/api/vault", nil, &resp)
	return resp.Data, err
}

type WrappedResponse[T any] struct {
	Result T   `json:"result"`
	Error  any `json:"error"`
}

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

// ---------------------------------------------------------
// Payment Setup
// ---------------------------------------------------------
type PaymentSetupRequest struct {
	PaymentIntentID       string `json:"payment_intent_id"`
	Rail                  string `json:"rail"`
	Wallet                string `json:"wallet"`
	Month                 int64  `json:"month"`
	TxHash                string `json:"tx_hash"`
	Email                 string
	FirstName             string                               `json:"first_name"`
	LastName              string                               `json:"last_name"`
	CardBrand             string                               `json:"card_brand"`
	ExpiryMonth           string                               `json:"exp_month"`
	ExpiryYear            string                               `json:"exp_year"`
	Currency              string                               `json:"currency"`
	Amount                string                               `json:"amount"`
	Plan                  string                               `json:"plan"`
	ProductID             string                               `json:"product_id"`
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
type PaymentSetupRequestBeta struct {
	Rail                  string `json:"rail"`
	Tier                  string `json:"tier"`
	Email             string `json:"email"`

	Plan                  string `json:"plan"`
	PeriodMonths          int    `json:"period_months"`
	SessionID             string `json:"session_id"`
}

type PaymentSetupResponse struct {
	Data       json.RawMessage `json:"data"`
	Status     string          `json:"status"`
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
}
// ---------------------------------------------------------
// Subscription
// ---------------------------------------------------------
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
// when onboarding new user, we dont have user id yet, so we use session id to get the subscription
func (c *TracecoreClient) GetSubscriptionBySessionID(ctx context.Context, sessionID string) (*subscription_domain.Subscription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/subscriptions?session_id="+sessionID, nil)
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

	var sub subscription_domain.Subscription
	if err := json.Unmarshal(cloudResp.Data, &sub); err != nil {
		return nil, fmt.Errorf("invalid cloud data: %w", err)
	}

	return &sub, nil
}
// after onboarding, we have user id, so we use user id to get the subscription
func (c *TracecoreClient) GetSubscriptionByUserID(ctx context.Context, userID string) (*subscription_domain.Subscription, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/subscriptions?user_id="+userID, nil)
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

	var sub subscription_domain.Subscription
	if err := json.Unmarshal(cloudResp.Data, &sub); err != nil {
		return nil, fmt.Errorf("invalid cloud data: %w", err)
	}

	return &sub, nil
}
func (c *TracecoreClient) GetPendingPaymentRequests(ctx context.Context, userID string) ([]*billing_domain.PaymentRequest, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/subscriptions?user_id="+userID, nil)
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

	var prs []*billing_domain.PaymentRequest
	if err := json.Unmarshal(cloudResp.Data, &prs); err != nil {
		return nil, fmt.Errorf("invalid cloud data: %w", err)
	}

	return prs, nil
}	
type ClientPaymentResponse struct {
	PaymentID string `json:"payment_id"`
	PaymentRequest *billing_domain.PaymentRequest `json:"payment_request"`
}
func (c *TracecoreClient) ProcessEncryptedPayment(ctx context.Context, request *billing_domain.ClientPaymentRequest) (*ClientPaymentResponse, error) {
	utils.LogPretty("request", request)
	reqBody, _ := json.Marshal(request)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/subscriptions/process-encrypted-payment", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	utils.LogPretty("request", req)
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
	var pr ClientPaymentResponse
	if err := json.Unmarshal(cloudResp.Data, &pr); err != nil {
		return nil, fmt.Errorf("invalid cloud data: %w", err)
	}
	return &pr, nil
}
func (c *TracecoreClient) HandleClientInitiatedPayment(ctx context.Context, request *billing_domain.ClientPaymentRequest) (*ClientPaymentResponse, error) {
	utils.LogPretty("request", request)
	reqBody, _ := json.Marshal(request)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/subscriptions/handle-client-initiated-payment", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	utils.LogPretty("request", req)
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
	var pr ClientPaymentResponse
	if err := json.Unmarshal(cloudResp.Data, &pr); err != nil {
		return nil, fmt.Errorf("invalid cloud data: %w", err)
	}
	return &pr, nil
}
func (c *TracecoreClient) GetBillingHistory(ctx context.Context, userID string, limit int) ([]*billing_domain.PaymentHistory, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/subscriptions/billing-history?user_id="+userID+"&limit="+strconv.Itoa(limit), nil)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	utils.LogPretty("request", req)
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
	var ph []*billing_domain.PaymentHistory
	if err := json.Unmarshal(cloudResp.Data, &ph); err != nil {
		return nil, fmt.Errorf("invalid cloud data: %w", err)
	}
	return ph, nil	
}
func (c *TracecoreClient) GenerateReceipt(ctx context.Context, userID string, paymentID string) (*billing_domain.Receipt, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/subscriptions/generate-receipt?user_id="+userID+"&payment_id="+paymentID, nil)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	utils.LogPretty("request", req)
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
	var r billing_domain.Receipt
	if err := json.Unmarshal(cloudResp.Data, &r); err != nil {
		return nil, fmt.Errorf("invalid cloud data: %w", err)
	}
	return &r, nil	
}
func (c *TracecoreClient) CancelSubscription(ctx context.Context, userID string, reason string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/subscriptions/cancel", nil)
	if err != nil {
		return err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	utils.LogPretty("request", req)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
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
		return fmt.Errorf("invalid cloud response: %w", err)
	}
	if cloudResp.Status != "ok" {
		return fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}
	return nil	
}
func (c *TracecoreClient) UpdatePaymentMethod(ctx context.Context, reqPaymentMethod *billing_domain.UpdatePaymentMethodRequest) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/subscriptions/update-payment-method", nil)
	if err != nil {
		return err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	utils.LogPretty("request", req)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
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
		return fmt.Errorf("invalid cloud response: %w", err)
	}
	if cloudResp.Status != "ok" {
		return fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}
	return nil	
}
func (c *TracecoreClient) GetStorageUsage(ctx context.Context, userID string) (*billing_domain.StorageUsage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/subscriptions/storage-usage?user_id="+userID, nil)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	utils.LogPretty("request", req)
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
	var su billing_domain.StorageUsage
	if err := json.Unmarshal(cloudResp.Data, &su); err != nil {
		return nil, fmt.Errorf("invalid cloud data: %w", err)
	}
	return &su, nil	
}	
func (c *TracecoreClient) HandleUpgrade(ctx context.Context, userID string, newTier subscription_domain.SubscriptionTier, paymentMethod subscription_domain.PaymentMethod) error {
	rerq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/subscriptions/handle-upgrade", nil)
	if err != nil {
		return err
	}
	if c.Token != "" {
		rerq.Header.Set("Authorization", "Bearer "+c.Token)
	}
	utils.LogPretty("request", rerq)
	resp, err := c.HTTPClient.Do(rerq)
	if err != nil {
		return err
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
		return fmt.Errorf("invalid cloud response: %w", err)
	}
	if cloudResp.Status != "ok" {
		return fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}
	return nil	
}
func (c *TracecoreClient) ReactivateSubscription(ctx context.Context, userID string, tier subscription_domain.SubscriptionTier, paymentMethod subscription_domain.PaymentMethod) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/subscriptions/reactivate", nil)
	if err != nil {
		return err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	utils.LogPretty("request", req)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
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
		return fmt.Errorf("invalid cloud response: %w", err)
	}
	if cloudResp.Status != "ok" {
		return fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}
	return nil	
}
// ---------------------------------------------------------
//
//	Crypto Share
//
// ---------------------------------------------------------
type WrappedShare struct {
	Data       CloudCryptographicShare  `json:"data"`
	Recipients []share_domain.Recipient `json:"recipients"`
}
type CloudCryptographicShare struct {
	ID               string
	EncryptedPayload string // AES-256-GCM encrypted vault entry
	SenderUserID     string
	SenderEmail      string
	SenderPublicKey  string // Stellar public key
	CreatedAt        time.Time
	RevokedAt        *time.Time
	AccessMode       string
	AccessLog        datatypes.JSON
	Signature        string

	Title           string         // optional
	EntryType       string         // optional
	DownloadAllowed bool           // optional
	Metadata        datatypes.JSON // flexible extra info

}
type CryptoRecipient struct {
	EncryptedKeys string     `json:"EncryptedKeys"`
	Role          string     `json:"Role"`
	RevokedAt     *time.Time `json:"RevokedAt"`
}
type ProdCreateCryptoShareRequest struct {
	SenderID      string                     `json:"SenderID"`
	SenderEmail   string                     `json:"SenderEmail"`
	Recipients    map[string]CryptoRecipient `json:"Recipients"`
	VaultPayload  string                     `json:"VaultPayload"`  // already encrypted
	EncryptedKeys map[string]string          `json:"EncryptedKeys"` // userID -> encrypted key
	Message       string                     `json:"Message"`
	PublicKey     string                     `json:"PublicKey"`
	Signature     string                     `json:"Signature"`
	Title         string                     `json:"Title"`
	EntryType     string                     `json:"EntryType"`
	Metadata      map[string]interface{}     `json:"Metadata,omitempty"`
	AccessMode    string                     `json:"AccessMode"`
	ExpiresAt    *time.Time                       `json:"ExpiresAt,omitempty"`
	DownloadAllowed bool                       `json:"DownloadAllowed,omitempty"`
}
type ProdCreateCryptoShareResponse struct {
	Data       json.RawMessage `json:"data"`
	Status     string          `json:"status"`
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
}

func (c *TracecoreClient) CreateShare(ctx context.Context, payload ProdCreateCryptoShareRequest) (*ProdCreateCryptoShareResponse, error) {
	utils.LogPretty("TracecoreClient - CreateCloudShare - payload", payload)
	bodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/shares/cryptographic", bytes.NewReader(bodyBytes))
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
	var cloudResp ProdCreateCryptoShareResponse
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}

	if cloudResp.Status != "ok" {
		return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}
	utils.LogPretty("TracecoreClient - CreateCloudShare - cloud response", cloudResp)

	return &cloudResp, nil
}

// ---------------------------------------------------------
// Get Share By Me
// ---------------------------------------------------------
func (c *TracecoreClient) GetShareByMe(ctx context.Context, email string) ([]share_domain.ShareEntry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/shares/cryptographic/by-me/"+email, nil)
	if err != nil {
		return nil, err
	}
	utils.LogPretty("GetShareByMe - email", email)
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	utils.LogPretty("GetShareByMe - respBytes", respBytes)
	var cloudResp struct {
		Data       []WrappedShare `json:"data"`
		Status     string         `json:"status"`
		StatusCode int            `json:"status_code"`
		Message    string         `json:"message"`
	}

	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}
	utils.LogPretty("GetShareByMe - cloud response", cloudResp)

	if cloudResp.Status != "ok" {
		return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}

	// convert to share_domain.ShareEntry
	list := CryptoShareConvertor(cloudResp.Data)

	return list, nil
}

// ---------------------------------------------------------
// Get Share With Me
// ---------------------------------------------------------
func (c *TracecoreClient) GetShareWithMe(ctx context.Context, email string) ([]share_domain.ShareEntry, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/shares/cryptographic/with-me/"+email, nil)
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
		Data       []WrappedShare `json:"data"`
		Status     string         `json:"status"`
		StatusCode int            `json:"status_code"`
		Message    string         `json:"message"`
	}

	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}

	if cloudResp.Status != "ok" {
		return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}

	list := CryptoShareConvertor(cloudResp.Data)

	return list, nil
}

// ---------------------------------------------------------
//
//	Link Share
//
// ---------------------------------------------------------
type LinkShare struct {
	ID        string    `json:"ID"`
	Payload   string    `json:"Payload"`
	CreatedAt time.Time `json:"CreatedAt"`

	ExpiresAt *time.Time `json:"ExpiresAt"`
	MaxViews  *int       `json:"MaxViews"`
	ViewCount int        `json:"ViewCount"`

	PasswordHash    *string `json:"PasswordHash"`
	DownloadAllowed bool    `json:"DownloadAllowed"`

	CreatorUserID string `json:"CreatorUserID"`
	CreatorEmail  string `json:"CreatorEmail"`

	Metadata Metadata `json:"Metadata"`
}

type Metadata struct {
	EntryType string `json:"entry_type"`
	Title     string `json:"title"`
}

func (l *LinkShare) ToWailsLinkShare() *WailsLinkShare {
	var expiry string
	if l.ExpiresAt != nil {
		expiry = l.ExpiresAt.Format("2006-01-02")
	} else {
		expiry = "never"
	}

	var usesLeft int
	if l.MaxViews != nil {
		usesLeft = *l.MaxViews - l.ViewCount
	} else {
		usesLeft = -1 // infinite
	}

	var password string
	if l.PasswordHash != nil {
		password = *l.PasswordHash
	}

	return &WailsLinkShare{
		ID:            l.ID,
		EntryName:     l.Metadata.Title,
		Status:        "active",
		Expiry:        expiry,
		UsesLeft:      usesLeft,
		Link:          "https://ankhora.app/shares/" + l.ID,
		AuditLog:      nil,
		Payload:       l.Payload,
		AllowDownload: l.DownloadAllowed,
		Password:      password,
	}
}

type WailsLinkShare struct {
	ID            string        `json:"id"`
	EntryName     string        `json:"entry_name"`
	Status        string        `json:"status"`
	Expiry        string        `json:"expiry"`
	UsesLeft      int           `json:"uses_left"`
	Link          string        `json:"link"`
	AuditLog      []interface{} `json:"audit_log"`
	Payload       string        `json:"payload"`
	AllowDownload bool          `json:"allow_download"`
	Password      string        `json:"password"`
}

type CreateLinkShareResponse struct {
	Data       LinkShare `json:"data"`
	Status     string    `json:"status"`
	StatusCode int       `json:"status_code"`
	Message    string    `json:"message"`
}
type LinkShareResponse struct {
	Data       []LinkShare `json:"data"`
	Status     string      `json:"status"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
}

func (c *TracecoreClient) CreateLinkShare(ctx context.Context, sh share_application_dto.LinkShareCreateRequest) (*CreateLinkShareResponse, error) {
	bodyBytes, _ := json.Marshal(sh)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/shares/link", bytes.NewReader(bodyBytes))
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
	var cloudResp CreateLinkShareResponse
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}

	if cloudResp.Status != "ok" {
		return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}
	utils.LogPretty("cloud response", cloudResp)

	return &cloudResp, nil
}

func (c *TracecoreClient) ListLinkSharesByMe(ctx context.Context, email string) (*LinkShareResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/shares/link/by-me/"+email, nil)
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
	utils.LogPretty("ListLinkSharesByMe - respBytes", respBytes)
	var cloudResp LinkShareResponse
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		log.Printf("invalid cloud response: %v", err)
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}
	utils.LogPretty("ListLinkSharesByMe - cloud response", cloudResp)

	return &cloudResp, nil
}
func (c *TracecoreClient) ListLinkSharesWithMe(ctx context.Context, email string) (*LinkShareResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/shares/link/with-me/"+email, nil)
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
	var cloudResp LinkShareResponse
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		log.Printf("invalid cloud response: %w", err)
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}
	utils.LogPretty("ListLinkSharesWithMe - cloud response", cloudResp)

	return &cloudResp, nil
}

// ---------------------------------------------------------
// Helper Functions
// ---------------------------------------------------------
func CryptoShareConvertor(cloudResp []WrappedShare) []share_domain.ShareEntry {
	var list []share_domain.ShareEntry
	for _, v := range cloudResp {
		list = append(list, share_domain.ShareEntry{
			ID:               v.Data.ID,
			OwnerID:          v.Data.SenderUserID,
			EntryName:        v.Data.Title,
			EntryType:        v.Data.EntryType,
			EntryRef:         v.Data.ID,
			Status:           "ok",
			AccessMode:       "cryptographic",
			AccessLog:        v.Data.AccessLog,
			Encryption:       "aes-256-gcm",
			EncryptedPayload: v.Data.EncryptedPayload,
			Recipients:       v.Recipients,
			ExpiresAt:        nil,
			SharedAt:         v.Data.CreatedAt,
			CreatedAt:        v.Data.CreatedAt,
		})
	}
	return list
}

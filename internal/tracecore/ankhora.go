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
	"os"
	"strings"
	"time"
	share_application_dto "vault-app/internal/application"
	billing_domain "vault-app/internal/billing/domain"
	share_domain "vault-app/internal/domain/shared"
	subscription_domain "vault-app/internal/subscription/domain"
	tracecore_types "vault-app/internal/tracecore/types"
	utils "vault-app/internal/utils"
	vaults_domain "vault-app/internal/vault/domain"

	"path"

	"gorm.io/datatypes"
)

func (c *TracecoreClient) GetVault(ctx context.Context) ([]byte, error) {
	var resp struct {
		Data []byte `json:"data"`
	}
	err := c.doRequest(ctx, "GET", "/api/vault", nil, &resp)
	return resp.Data, err
}

type GetUserByEmailResponse struct {
	Error   bool                 `json:"error"`
	Message string               `json:"message"`
	Data    tracecore_types.User `json:"data"`
}

// ---------------------------------------------------------
// User
// ---------------------------------------------------------
func (c *TracecoreClient) GetUserByEmail(ctx context.Context, email string) (*tracecore_types.User, error) {
	// Build endpoint (ensure no double slashes)
	base := strings.TrimRight(c.AnkhoraCloudUrl, "/")
	q := url.Values{}
	q.Set("email", email)

	endpoint := fmt.Sprintf("%s/customers?%s", base, q.Encode())
	fmt.Println("endpoint", endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		fmt.Println("TracecoreClient - failed to create request: %s", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		fmt.Printf("TracecoreClient - failed to create request: %v\n", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle HTTP status codes
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("TracecoreClient - User not found for email: %s", email)
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
	var result GetUserByEmailResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error {
		return nil, fmt.Errorf("API error: %s", result.Message)
	}

	if result.Data.ID == 0 {
		return nil, ErrUserNotFound
	}

	return &result.Data, nil
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
	Rail  string `json:"rail"`
	Tier  string `json:"tier"`
	Month int64  `json:"month"`
	UserID string `json:"user_id"`
	Email string `json:"email"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`

	Plan         string `json:"plan"`
	PeriodMonths string `json:"period_months"`
	Amount string `json:"amount"`
	ProductID string `json:"product_id"`
	Currency string `json:"currency"`
	IsAnonymous bool `json:"is_anonymous"`
	SessionID    string `json:"session_id"`
}

type PaymentSetupResponse struct {
	Data       json.RawMessage `json:"data"`
	Status     string          `json:"status"`
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
}
type FreeCheckoutResponse struct {
	Data       json.RawMessage `json:"data"`
	Status     int          `json:"status"`
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

func (c *TracecoreClient) FreeCheckout(ctx context.Context, req PaymentSetupRequestBeta) (*FreeCheckoutResponse, error) {
	utils.LogPretty("TracecoreClient - FreeCheckout - payload", req)
	bodyBytes, _ := json.Marshal(req)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.AnkhoraCloudUrl+"/subscriptions/activate", bytes.NewReader(bodyBytes))
	if err != nil {
		utils.LogPretty("TracecoreClient - FreeCheckout - error", err)
		return nil, err
	}
	defer request.Body.Close()

	request.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var cloudResp FreeCheckoutResponse
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}

	if cloudResp.Status != 201 {
		return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}
	utils.LogPretty("cloud response", cloudResp)

	return &cloudResp, nil
}
// when onboarding new user, we dont have user id yet, so we use session id to get the subscription
func (c *TracecoreClient) GetSubscriptionBySessionID(
	ctx context.Context,
	sessionID string,
) (*subscription_domain.Subscription, error) {

	q := url.Values{}
	q.Set("session_id", sessionID)
	utils.LogPretty("session id", sessionID)
	utils.LogPretty("ankhora cloud url", c.AnkhoraCloudUrl)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.AnkhoraCloudUrl+"/subscriptions?"+q.Encode(),
		nil,
	)
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type CloudResponse struct {
		Status  int                              `json:"status"`
		Data    subscription_domain.Subscription `json:"data"`
		Success bool                             `json:"success"`
		Message string                           `json:"message"`
	}

	var cloudResp CloudResponse
	if err := json.Unmarshal(body, &cloudResp); err != nil {
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}

	if !cloudResp.Success {
		return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}

	return &cloudResp.Data, nil
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
func (c *TracecoreClient) GetSubscriptionByID(ctx context.Context, subscriptionID string) (*subscription_domain.Subscription, error) {
	
	utils.LogPretty("subscription id", subscriptionID)
	utils.LogPretty("ankhora cloud url", c.AnkhoraCloudUrl)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.AnkhoraCloudUrl+"/subscriptions/"+subscriptionID,
		nil,
	)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type CloudResponse struct {
		Status  int                              `json:"status"`
		Data    subscription_domain.Subscription `json:"data"`
		Success bool                             `json:"success"`
		Message string                           `json:"message"`
	}

	var cloudResp CloudResponse
	if err := json.Unmarshal(body, &cloudResp); err != nil {
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}

	if !cloudResp.Success {
		return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}

	return &cloudResp.Data, nil
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
	PaymentID      string                         `json:"payment_id"`
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
func (c *TracecoreClient) GetBillingHistory(ctx context.Context, subID string, limit int) (*tracecore_types.CloudResponse[[]tracecore_types.PaymentHistory], error) {
	url := fmt.Sprintf("%s/billing/history/%s?limit=%d", c.AnkhoraCloudUrl, subID, limit)
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		utils.LogPretty("GetBillingHistory - request failed", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.LogPretty("GetBillingHistory - failed to read response body", err)
		return nil, err
	}

	var cloudResp tracecore_types.CloudResponse[[]tracecore_types.PaymentHistory]
	if err := json.Unmarshal(body, &cloudResp); err != nil {
		utils.LogPretty("GetBillingHistory - invalid cloud response", err)
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}

	if !cloudResp.Success {
		utils.LogPretty("GetBillingHistory - cloud returned error", cloudResp.Message)
		return nil, fmt.Errorf("GetBillingHistory - cloud returned error: %s", cloudResp.Message)
	}

	return &cloudResp, nil
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
//	Cryptographic Share
//
// ---------------------------------------------------------

type WrappedShare struct {
	Data       CloudCryptographicShare  `json:"data"`
	Recipients []share_domain.Recipient `json:"recipients"`
}
type CloudCryptographicShare struct {
	ID               string         `json:"ID"`
	EncryptedPayload string         `json:"EncryptedPayload"`
	SenderUserID     string         `json:"SenderUserID"`
	SenderEmail      string         `json:"SenderEmail"`
	SenderPublicKey  string         `json:"SenderPublicKey"`
	CreatedAt        time.Time      `json:"CreatedAt"`
	RevokedAt        *time.Time     `json:"RevokedAt"`
	AccessMode       string         `json:"AccessMode"`
	AccessLog        datatypes.JSON `json:"AccessLog"`
	Signature        string         `json:"Signature"`
	Title            string         `json:"Title"`
	EntryType        string         `json:"EntryType"`
	DownloadAllowed  bool           `json:"DownloadAllowed"`
	Metadata         datatypes.JSON `json:"Metadata"`
}

type CryptoRecipient struct {
	EncryptedKeys string     `json:"EncryptedKeys"`
	Role          string     `json:"Role"`
	RevokedAt     *time.Time `json:"RevokedAt"`
}
type ProdCreateCryptoShareRequest struct {
	SenderID        string                     `json:"SenderID"`
	SenderEmail     string                     `json:"SenderEmail"`
	Recipients      map[string]CryptoRecipient `json:"Recipients"`
	VaultPayload    string                     `json:"VaultPayload"`  // already encrypted
	EncryptedKeys   map[string]string          `json:"EncryptedKeys"` // userID -> encrypted key
	Message         string                     `json:"Message"`
	PublicKey       string                     `json:"PublicKey"`
	Signature       string                     `json:"Signature"`
	Title           string                     `json:"Title"`
	EntryType       string                     `json:"EntryType"`
	Metadata        map[string]interface{}     `json:"Metadata,omitempty"`
	AccessMode      string                     `json:"AccessMode"`
	ExpiresAt       *time.Time                 `json:"ExpiresAt,omitempty"`
	DownloadAllowed bool                       `json:"DownloadAllowed,omitempty"`
}
type ProdCreateCryptoShareResponse struct {
	Data    CloudCryptographicShare `json:"data"`
	Status  int                     `json:"status"`
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
}

func (c *TracecoreClient) CreateShare(ctx context.Context, payload ProdCreateCryptoShareRequest) (*ProdCreateCryptoShareResponse, error) {
	utils.LogPretty("TracecoreClient - CreateCloudShare - payload", payload)
	bodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.AnkhoraCloudUrl+"/shares/cryptographic", bytes.NewReader(bodyBytes))
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
	// var cloudResp CloudResponse[CloudCryptographicShare]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}

	if cloudResp.Status != 201 {
		return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}
	utils.LogPretty("TracecoreClient - CreateCloudShare - cloud response", cloudResp)

	return &cloudResp, nil
}

func (c *TracecoreClient) AccessEncryptedEntry(ctx context.Context, id string, req tracecore_types.AccessCryptoShareRequest) (*tracecore_types.CloudResponse[tracecore_types.AccessCryptoShareResponse], error) {
	utils.LogPretty("TracecoreClient - AccessEncryptedEntry - payload", req)
	bodyBytes, _ := json.Marshal(req)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.AnkhoraCloudUrl+"/shares/cryptographic/"+id+"/access", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	request.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var cloudResp tracecore_types.CloudResponse[tracecore_types.AccessCryptoShareResponse]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}

	if cloudResp.Status != 200 {
		return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}
	utils.LogPretty("TracecoreClient - AccessEncryptedEntry - cloud response", cloudResp)

	return &cloudResp, nil
}

func (c *TracecoreClient) DecryptVaultEntry(ctx context.Context, req tracecore_types.DecryptCryptoShareRequest) (*tracecore_types.CloudResponse[tracecore_types.DecryptCryptoShareResponse], error) {
	utils.LogPretty("TracecoreClient - DecryptVaultEntry - payload", req)
	bodyBytes, _ := json.Marshal(req)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.AnkhoraCloudUrl+"/shares/cryptographic/decrypt", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	request.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var cloudResp tracecore_types.CloudResponse[tracecore_types.DecryptCryptoShareResponse]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}

	if cloudResp.Status != 200 {
		return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}
	utils.LogPretty("TracecoreClient - DecryptVaultEntry - cloud response", cloudResp)

	return &cloudResp, nil
}

func (c *TracecoreClient) AddRecipient(ctx context.Context, req tracecore_types.AddRecipientRequest) (*tracecore_types.CloudResponse[CloudCryptographicShare], error) {
	utils.LogPretty("TracecoreClient - AddRecipient - payload", req)
	bodyBytes, _ := json.Marshal(req)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.AnkhoraCloudUrl+"/shares/cryptographic/"+req.ShareID+"/recipient", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	request.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var cloudResp tracecore_types.CloudResponse[CloudCryptographicShare]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("TracecoreClient - AddRecipient - invalid cloud response: %w", err)
	}

	if cloudResp.Status != 200 {
		return nil, fmt.Errorf("TracecoreClient - AddRecipient - cloud returned error: %s", cloudResp.Message)
	}
	utils.LogPretty("TracecoreClient - AddRecipient - cloud response", cloudResp)

	return &cloudResp, nil
}

func (c *TracecoreClient) UpdateRecipient(ctx context.Context, req share_application_dto.UpdateRecipientRequest) (*tracecore_types.CloudResponse[CloudCryptographicShare], error) {
	utils.LogPretty("TracecoreClient - UpdateRecipient - payload", req)
	bodyBytes, _ := json.Marshal(req)
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, c.AnkhoraCloudUrl+"/shares/cryptographic/"+req.ShareID+"/recipient/"+req.Email, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	request.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var cloudResp tracecore_types.CloudResponse[CloudCryptographicShare]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("TracecoreClient - UpdateRecipient - invalid cloud response: %w", err)
	}

	if cloudResp.Status != 200 {
		return nil, fmt.Errorf("TracecoreClient - UpdateRecipient - cloud returned error: %s", cloudResp.Message)
	}
	utils.LogPretty("TracecoreClient - UpdateRecipient - cloud response", cloudResp)

	return &cloudResp, nil
}

func (c *TracecoreClient) RevokeShare(ctx context.Context, req tracecore_types.RevokeShareRequest) (*tracecore_types.CloudResponse[CloudCryptographicShare], error) {
	utils.LogPretty("TracecoreClient - RevokeShare - payload", req)
	bodyBytes, _ := json.Marshal(req)
	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.AnkhoraCloudUrl+"/shares/cryptographic/"+req.ShareID+"/revoke", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	request.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var cloudResp tracecore_types.CloudResponse[CloudCryptographicShare]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		return nil, fmt.Errorf("TracecoreClient - RevokeRecipient - invalid cloud response: %w", err)
	}

	if cloudResp.Status != 200 {
		return nil, fmt.Errorf("TracecoreClient - RevokeRecipient - cloud returned error: %s", cloudResp.Message)
	}
	utils.LogPretty("TracecoreClient - RevokeRecipient - cloud response", cloudResp)

	return &cloudResp, nil
}

// ---------------------------------------------------------
// Get Share By Me
// ---------------------------------------------------------
func (c *TracecoreClient) GetShareByMe(
	ctx context.Context,
	email string,
) ([]share_domain.ShareEntry, error) {
	u, err := url.Parse(c.AnkhoraCloudUrl)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "shares", "cryptographic", "by-me", email)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cloudResp tracecore_types.CloudResponse[[]WrappedShare]
	if err := json.Unmarshal(body, &cloudResp); err != nil {
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}

	if !cloudResp.Success {
		return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}

	return CryptoShareConvertor(cloudResp.Data), nil
}

// ---------------------------------------------------------
// Get Share With Me
// ---------------------------------------------------------
func (c *TracecoreClient) GetShareWithMe(ctx context.Context, email string) ([]share_domain.ShareEntry, error) {

	u, err := url.Parse(c.AnkhoraCloudUrl)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "shares", "cryptographic", "with-me", email)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cloudResp tracecore_types.CloudResponse[[]WrappedShare]
	if err := json.Unmarshal(body, &cloudResp); err != nil {
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}

	if !cloudResp.Success {
		return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}

	return CryptoShareConvertor(cloudResp.Data), nil
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
		Link:          os.Getenv("CLOUD_FRONT_URL") + "/shares/" + l.ID,
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
	Data    LinkShare `json:"data"`
	Status  int       `json:"status"`
	Code    int       `json:"code"`
	Message string    `json:"message"`
}
type LinkShareResponse struct {
	Data    []LinkShare `json:"data"`
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
}

func (c *TracecoreClient) CreateLinkShare(ctx context.Context, sh share_application_dto.LinkShareCreateRequest) (*CreateLinkShareResponse, error) {
	bodyBytes, _ := json.Marshal(sh)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.AnkhoraCloudUrl+"/shares/link", bytes.NewReader(bodyBytes))
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

	if cloudResp.Status != 201 {
		return nil, fmt.Errorf("cloud returned error: %s", cloudResp.Message)
	}
	utils.LogPretty("cloud response", cloudResp)

	return &cloudResp, nil
}

func (c *TracecoreClient) ListLinkSharesByMe(ctx context.Context, email string) (*LinkShareResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.AnkhoraCloudUrl+"/shares/link/by-me/"+email, nil)
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.AnkhoraCloudUrl+"/shares/link/with-me/"+email, nil)
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
// VAULT
// ---------------------------------------------------------
func (c *TracecoreClient) AddToIPFS(ctx context.Context, req tracecore_types.SyncVaultStreamRequest) (*tracecore_types.CloudResponse[tracecore_types.SyncVaultResponse], error) {
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(req); err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/vaults/"+req.UserID+"/storage/"+req.VaultName, body)
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
	utils.LogPretty("AddToIPFS - cloud response", cloudResp)

	return &cloudResp, nil
}
func (c *TracecoreClient) GetDataFromCloudStorage(ctx context.Context, req tracecore_types.IpfsCidRequest) (*tracecore_types.IpfsCidResponse, error) {
	utils.LogPretty("GetDataFromCloudStorage - req", req)
	if c.BaseURL == "" {
		return nil, fmt.Errorf("TracecoreClient.BaseURL is empty")
	}
	base := strings.TrimRight(c.BaseURL, "/")
	path := fmt.Sprintf("/vaults/%s/storage/%s/%s",
		url.PathEscape(req.UserID), // ← Check why this is EMPTY!
		url.PathEscape(req.VaultName),
		req.CID)

	// path := fmt.Sprintf("/vaults/%s/storage/%s/%s", url.PathEscape(req.UserID), url.PathEscape(req.VaultName), req.CID)

	fullURL := base + path
	utils.LogPretty("GetDataFromCloudStorage - fullURL", fullURL)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		log.Printf("GetDataFromCloudStorage - invalid cloud response: %v", err)
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}
	if c.Token != "" {
		request.Header.Set("Authorization", "Bearer "+c.Token)
	}
	// request.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		if os.IsTimeout(err) {
			log.Printf("GetDataFromCloudStorage - request to %s timed out: %v", fullURL, err)
		} else {
			log.Printf("GetDataFromCloudStorage - request error: %v", err)
		}
		return nil, err
	}

	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	utils.LogPretty("RAW backend response", string(respBytes))

	// STEP 2
	var cloudResp tracecore_types.IpfsCidResponse
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		utils.LogPretty("🚫 STEP 2 FAILED", err) // ADD THIS
		return nil, fmt.Errorf("cloud response unmarshal failed: %w", err)
	}
	utils.LogPretty("✅ STEP 2 OK", cloudResp) // ADD THIS

	utils.LogPretty("✅ SUCCESS - cloud response", cloudResp)
	return &cloudResp, nil

}
func (c *TracecoreClient) SyncVaultToIPFS(ctx context.Context, req tracecore_types.SyncVaultStreamRequest) (*tracecore_types.CloudResponse[tracecore_types.SyncVaultResponse], error) {
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(req); err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.AnkhoraCloudUrl+"/vaults/"+req.UserID+"/sync/"+req.VaultName, body)
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
	utils.LogPretty("SyncVaultToIPFS - cloud response", cloudResp)

	return &cloudResp, nil
}
func (c *TracecoreClient) GetVaultByUserIDAndName(ctx context.Context, input tracecore_types.GetVaultInput) (*tracecore_types.CloudResponse[vaults_domain.Vault], error) {
	utils.LogPretty("GetVaultByUserIDAndName - req", input)
	u, err := url.Parse(c.AnkhoraCloudUrl)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "vaults", input.UserID, input.VaultName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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
	var cloudResp tracecore_types.CloudResponse[vaults_domain.Vault]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		log.Printf("invalid cloud response: %w", err)
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}
	utils.LogPretty("GetVaultByUserIDAndName - cloud response", cloudResp)

	return &cloudResp, nil
}
func (c *TracecoreClient) GetVaultBySubscription(ctx context.Context, subID string) (*tracecore_types.CloudResponse[tracecore_types.Vault], error) {
	utils.LogPretty("GetVaultBySubscription - req", subID)
	u, err := url.Parse(c.AnkhoraCloudUrl)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "vaults", "subscription", subID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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
	var cloudResp tracecore_types.CloudResponse[tracecore_types.Vault]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		log.Printf("invalid cloud response: %w", err)
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}
	utils.LogPretty("GetVaultByUserIDAndName - cloud response", cloudResp)

	return &cloudResp, nil
}
func (c *TracecoreClient) AddPublicKeyToCustomer(ctx context.Context, req tracecore_types.AddPublicKeyToCustomerRequest) (*tracecore_types.CloudResponse[tracecore_types.AddPublicKeyToCustomerResponse], error) {
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(req); err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.AnkhoraCloudUrl+"/customers/add-public-key", body)
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
	var cloudResp tracecore_types.CloudResponse[tracecore_types.AddPublicKeyToCustomerResponse]
	if err := json.Unmarshal(respBytes, &cloudResp); err != nil {
		log.Printf("invalid cloud response: %w", err)
		return nil, fmt.Errorf("invalid cloud response: %w", err)
	}
	utils.LogPretty("AddPublicKeyToCustomer - cloud response", cloudResp)

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
			OwnerID:          v.Data.SenderEmail,
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

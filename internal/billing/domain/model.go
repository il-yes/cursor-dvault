package billing_domain

import (
	"time"
	subscription_domain "vault-app/internal/subscription/domain"
)

type PaymentMethod string

const (
	PaymentCard PaymentMethod = "card"
	PaymentCrypto PaymentMethod = "crypto"
	PaymentEncryptedCard PaymentMethod = "encrypted_card"
)

type BillingInstrument struct {
	ID string
	UserID string
	Type PaymentMethod
	EncryptedPayload string // when applicable
	CreatedAt time.Time
}


type ClientPaymentRequest struct {
    PaymentRequestID       string `json:"payment_request_id"`
    StripePaymentMethodID  string `json:"stripe_payment_method_id"`
}

type UpdatePaymentMethodRequest struct {
    UserID                 string        `json:"user_id"`
    PaymentMethod          PaymentMethod `json:"payment_method"`
    StripePaymentMethodID  string        `json:"stripe_payment_method_id,omitempty"`
    EncryptedPaymentData   string        `json:"encrypted_payment_data,omitempty"`
    StellarPublicKey       string        `json:"stellar_public_key,omitempty"`
}

type UpgradeRequest struct {
    UserID        string                    `json:"user_id"`
    NewTier       subscription_domain.SubscriptionTier `json:"new_tier"`
    PaymentMethod subscription_domain.PaymentMethod    `json:"payment_method"`
}

type PaymentHistory struct {
    ID              string    `json:"id"`
    Amount          float64   `json:"amount"`
    Status          string    `json:"status"` // succeeded, failed, pending
    PaymentMethod   string    `json:"payment_method"`
    Description     string    `json:"description"`
    StellarTxHash   string    `json:"stellar_tx_hash,omitempty"`
    StripeIntentID  string    `json:"stripe_intent_id,omitempty"`
    CreatedAt       time.Time `json:"created_at"`
}

type Receipt struct {
    ID              string    `json:"id"`
    UserID          string    `json:"user_id"`
    Amount          float64   `json:"amount"`
    Tier            string    `json:"tier"`
    PaymentMethod   string    `json:"payment_method"`
    StellarTxHash   string    `json:"stellar_tx_hash"`
    VerificationURL string    `json:"verification_url"` // Link to Stellar explorer
    IssuedAt        time.Time `json:"issued_at"`
    PDFData         []byte    `json:"pdf_data,omitempty"`
}

type StorageUsage struct {
    UsedGB      float64 `json:"used_gb"`
    QuotaGB     int     `json:"quota_gb"`
    Percentage  float64 `json:"percentage"`
    CanUpload   bool    `json:"can_upload"`
}

type Notification struct {
    ID        string                 `json:"id"`
    UserID    string                 `json:"user_id"`
    Type      string                 `json:"type"`
    Title     string                 `json:"title"`
    Message   string                 `json:"message"`
    ActionURL string                 `json:"action_url,omitempty"`
    Priority  string                 `json:"priority"` // low, normal, high
    Status    string                 `json:"status"`   // unread, read
    Data      map[string]interface{} `json:"data,omitempty"`
    CreatedAt time.Time              `json:"created_at"`
    ReadAt    time.Time              `json:"read_at,omitempty"`
}

type PaymentRequest struct {
	ID          string `json:"id"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      string `json:"user_id"`
	EncryptedPaymentEntryID string `json:"encrypted_payment_entry_id"`
}
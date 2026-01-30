package subscription_domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Subscription aggregate
type SubscriptionTier string

const (
	TierFree     SubscriptionTier = "free"
	TierPro      SubscriptionTier = "pro"
	TierProPlus  SubscriptionTier = "pro_plus"
	TierBusiness SubscriptionTier = "business"
    TierEnterprise  SubscriptionTier = "enterprise"
)

type PaymentMethod string

const (
	PaymentStandard  PaymentMethod = "standard"  // Standard Stripe checkout
	PaymentEncrypted PaymentMethod = "encrypted" // Encrypted card in vault
	PaymentCrypto    PaymentMethod = "crypto"    // USDC/Bitcoin
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusInactive SubscriptionStatus = "inactive"
	SubscriptionStatusUpgraded SubscriptionStatus = "upgraded"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
)

type Subscription struct {
	ID            string        `json:"id"`
	Email         string        `json:"email"`
	Wallet        string        `json:"wallet,omitempty"` // only filled for crypto billing
	UserID        string        `json:"user_id"`          // optional, not required for crypto validation
	Tier          string        `json:"tier"`
	ExpiresAt     int64         `json:"expires_at"`
	Rail          string        `json:"rail"`              // "traditional" or "crypto"
	TxHash        string        `json:"tx_hash,omitempty"` // Tx hash confirmation payment
	Active        bool          `json:"active"`
	ActivatedAt   int64         `json:"activated_at"`
	Months        int           `json:"months"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	PaymentIntent string        `json:"payment_intent"`
	StartedAt     time.Time     `json:"started_at"`

	Features              SubscriptionFeatures `json:"features"`
	Ledger                int32                `json:"ledger"`
	BillingCycle          string               `json:"billing_cycle"`
	TrialEndsAt           int64                `json:"trial_ends_at"`
	NextBillingDate       int64                `json:"next_billing_date"`
	Price                 float64              `json:"price"`
	StripeSubscriptionID  string               `json:"stripe_subscription_id"`
	StellarPaymentAddress string               `json:"stellar_payment_address"`
	PaymentFlowType       string               `json:"payment_flow_type"`
	StellarScheduleID     string               `json:"stellar_schedule_id"`
	Status                string               `json:"status"`
	EndedAt               int64                `json:"ended_at"`
	CreatedAt             time.Time            `json:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at"`
	CancelReason          string               `json:"cancel_reason"`
	CancelledAt           time.Time            `json:"cancelled_at"`
	Version               int64                `json:"version"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New().String()
	return
}
func (f SubscriptionFeatures) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *SubscriptionFeatures) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan SubscriptionFeatures: %v", value)
	}
	return json.Unmarshal(b, f)
}

func (s *Subscription) Validate() error {
	if s.ID == "" || s.Tier == "" || s.Rail == "" {
		return ErrInvalidSubscription
	}
	return nil
}

type SubscriptionFeatures struct {
	SubscriptionID      string `json:"subscription_id"`
	StorageGB           int
	StorageType         string
	CloudBackup         bool
	MobileApps          bool
	SharingLimit        int
	UnlimitedSharing    bool
	VersionHistory      bool
	VersionHistoryDays  int
	StellarVerification bool
	Telemetry           bool
	AnonymousAccount    bool
	CryptoPayments      bool
	EncryptedPayments   bool
	Support             string
	APIAccess           bool
	Tracecore           bool
	SSO                 bool
	TeamFeatures        bool
	PaymentMethod       PaymentMethod               `json:"payment_method"`
	PaymentIntent       string                      `json:"payment_intent"`
	BrowserExtension    bool                        `json:"browser_extension"`
	ThreatDetection     bool                        `json:"threat_detection"`
	PriorityStellar     bool                        `json:"priority_stellar"` // Priority Stellar verification
	TeamSize            int                         `json:"team_size"`
	GitCLI              bool                        `json:"git_cli"`
	CustomStellar       bool                        `json:"custom_stellar"`
	OnPremise           bool                        `json:"on_premise"`
	MultiCloud          bool                        `json:"multi_cloud"`
	CustomIntegrations  bool                        `json:"custom_integrations"`
	AISovereignty       bool                        `json:"ai_sovereignty"`
	Compliance          datatypes.JSONSlice[string] `json:"compliance"`
	SLA                 string                      `json:"sla"`
}

type UserSubscription struct {
	ID              string       `json:"id" gorm:"primaryKey"`
	Username        string    `gorm:"column:username" json:"username"`
	Email           string    `gorm:"column:email" json:"email"`
	Password        string    `gorm:"column:password" json:"password"`
	Role            string    `gorm:"column:role" json:"role"`
	CreatedAt       time.Time `json:"created_at" gorm:"varchar(100)"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"varchar(100)"`
	LastConnectedAt time.Time `json:"last_connected_at" gorm:"last_connected_at"`
	StellarPublicKey string    `json:"stellar_public_key"`
}

type SubscriptionFeaturesV1 struct {
	// Storage
	StorageGB   int    `json:"storage_gb"`
	StorageType string `json:"storage_type"` // local_ipfs, pinata_ipfs, custom_multi_cloud
	CloudBackup bool   `json:"cloud_backup"`

	// Sharing
	SharingLimit     int  `json:"sharing_limit"` // 0 = unlimited
	UnlimitedSharing bool `json:"unlimited_sharing"`

	// Version control
	VersionHistory     bool `json:"version_history"`
	VersionHistoryDays int  `json:"version_history_days"`

	// Blockchain
	StellarVerification bool `json:"stellar_verification"`
	PriorityStellar     bool `json:"priority_stellar"`
	CustomStellar       bool `json:"custom_stellar"`

	// Privacy
	Telemetry        bool `json:"telemetry"`
	AnonymousAccount bool `json:"anonymous_account"`

	// Payments
	CryptoPayments    bool `json:"crypto_payments"`
	EncryptedPayments bool `json:"encrypted_payments"`

	// Apps and access
	MobileApps       bool `json:"mobile_apps"`
	BrowserExtension bool `json:"browser_extension"`
	ThreatDetection  bool `json:"threat_detection"`

	// Support
	Support string `json:"support"` // community, email_24_48h, encrypted_chat_12h, 24_7_live, dedicated_account_manager

	// Business features
	APIAccess     bool `json:"api_access"`
	Tracecore     bool `json:"tracecore"`
	SSO           bool `json:"sso"`
	TeamFeatures  bool `json:"team_features"`
	TeamSize      int  `json:"team_size,omitempty"`
	GitCLI        bool `json:"git_cli"`
	AISovereignty bool `json:"ai_sovereignty"`

	// Enterprise
	OnPremise          bool     `json:"on_premise"`
	MultiCloud         bool     `json:"multi_cloud"`
	Compliance         []string `json:"compliance,omitempty"`
	CustomIntegrations bool     `json:"custom_integrations"`
	SLA                string   `json:"sla,omitempty"` // 99.9, 99.99, etc.
}

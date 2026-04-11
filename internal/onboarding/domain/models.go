package onboarding_domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)





type User struct {
	ID string `json:"id"`
	IsAnonymous bool `json:"is_anonymous"`
	StellarPublicKey string `json:"stellar_public_key"`
	CreatedAt time.Time `json:"created_at"`
	LastConnectedAt  time.Time

	Email string `json:"email"`
	Password string `json:"password"`	
}
func NewUser(isAnonymous bool, email string, password string) User {
	return User{
		ID: uuid.New().String(),
		IsAnonymous: isAnonymous,
		CreatedAt: time.Now(),
		Email: email,
		Password: password,
		LastConnectedAt: time.Now(),
	}
}
func (u *User) AttachStellarPublicKey(publicKey string) {
	u.StellarPublicKey = publicKey
}

type SubscriptionFeatures struct {
    // Storage
    StorageGB           int      `json:"storage_gb"`
    StorageType         string   `json:"storage_type"` // local_ipfs, pinata_ipfs, custom_multi_cloud
    CloudBackup         bool     `json:"cloud_backup"`
    
    // Sharing
    SharingLimit        int      `json:"sharing_limit"` // 0 = unlimited
    UnlimitedSharing    bool     `json:"unlimited_sharing"`
    
    // Version control
    VersionHistory      bool     `json:"version_history"`
    VersionHistoryDays  int      `json:"version_history_days"`
    
    // Blockchain
    StellarVerification bool     `json:"stellar_verification"`
    PriorityStellar     bool     `json:"priority_stellar"`
    CustomStellar       bool     `json:"custom_stellar"`
    
    // Privacy
    Telemetry           bool     `json:"telemetry"`
    AnonymousAccount    bool     `json:"anonymous_account"`
    
    // Payments
    CryptoPayments      bool     `json:"crypto_payments"`
    EncryptedPayments   bool     `json:"encrypted_payments"`
    
    // Apps and access
    MobileApps          bool     `json:"mobile_apps"`
    BrowserExtension    bool     `json:"browser_extension"`
    ThreatDetection     bool     `json:"threat_detection"`
    
    // Support
    Support             string   `json:"support"` // community, email_24_48h, encrypted_chat_12h, 24_7_live, dedicated_account_manager
    
    // Business features
    APIAccess           bool     `json:"api_access"`
    Tracecore           bool     `json:"tracecore"`
    SSO                 bool     `json:"sso"`
    TeamFeatures        bool     `json:"team_features"`
    TeamSize            int      `json:"team_size,omitempty"`
    GitCLI              bool     `json:"git_cli"`
    AISovereignty       bool     `json:"ai_sovereignty"`
    
    // Enterprise
    OnPremise           bool     `json:"on_premise"`
    MultiCloud          bool     `json:"multi_cloud"`
    Compliance          []string `json:"compliance,omitempty"`
    CustomIntegrations  bool     `json:"custom_integrations"`
    SLA                 string   `json:"sla,omitempty"` // 99.9, 99.99, etc.
}



type AppState struct {
    ID string `json:"id" gorm:"column:id"`
	HasVault        bool `json:"has_vault" gorm:"column:has_vault"`
	HasSession      bool `json:"has_session" gorm:"column:has_session"`
	HasImportedKey  bool `json:"has_imported_key" gorm:"column:has_imported_key"`
	NeedsOnboarding bool `json:"needs_onboarding" gorm:"column:needs_onboarding"`
    IsRecordSupported bool `json:"is_record_supported" gorm:"column:is_record_supported"`
}
func (a *AppState) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New().String()
	return nil
}
func NewAppState() *AppState {
	return &AppState{
		ID: uuid.New().String(),
		HasVault:        false,
		HasSession:      false,
		HasImportedKey:  false,
		NeedsOnboarding: true,
        IsRecordSupported: false,
	}
}

func (a *AppState) SetHasVault(hasVault bool) {
	a.HasVault = hasVault
}

func (a *AppState) SetHasSession(hasSession bool) {
	a.HasSession = hasSession
}

func (a *AppState) SetHasImportedKey(hasImportedKey bool) {
	a.HasImportedKey = hasImportedKey
}

func (a *AppState) SetNeedsOnboarding(needsOnboarding bool) {
	a.NeedsOnboarding = needsOnboarding
}

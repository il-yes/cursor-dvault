package tracecore_types

import (
	"time"
)

type LoginRequest struct {
    Email         string `json:"email,omitempty"`
    Password      string `json:"password,omitempty"`
    PublicKey     string `json:"public_key,omitempty"`
    SignedMessage string `json:"signed_message,omitempty"`
    Signature     string `json:"signature,omitempty"`
}


type CloudLoginResponse struct {
	Error               bool   `json:"error"`
	Message             string `json:"message"`
	AuthenticationToken struct {
		Token  string    `json:"token"`
		Expiry time.Time `json:"expiry"`
	} `json:"authentication_token"`
}

type CloudAuthToken struct {
    Token  string    `json:"token"`
    Expiry time.Time `json:"expiry"`
}

type LoginResponse struct {
    Token string `json:"token"`

    AuthenticationToken *struct {
        Token  string    `json:"token"`
        Expiry time.Time `json:"expiry"`
    } `json:"authentication_token"`
}

type VaultEntry struct {
    ID        string    `json:"id"`
    EntryName string    `json:"entry_name"`
    Type      string    `json:"type"`
    UpdatedAt time.Time `json:"updated_at"`
}

// type ShareEntry struct {
//     ID        string    `json:"id"`
//     EntryRef  string    `json:"entry_ref"`
//     EntryName string    `json:"entry_name"`
//     Status    string    `json:"status"`
//     SharedAt  time.Time `json:"shared_at"`
// }
type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	PublicKey string `json:"public_key"`
	Password  string `json:"password"` // if needed
}

type SyncVaultStreamRequest struct {
    UserID   string
    VaultName  string
    // Metadata map[string]string
    Stream   []byte
}

type SyncVaultResponse struct {
	UserID    string
	CID       string
	CreatedAt string
	// Metadata  map[string]string `json:"metadata,omitempty"`
}
type VaultSync struct {
	ID string

	// Relations
	VaultID string
	UserID  string
	UserEmail string

	// Storage result
	CID        string
	SizeBytes int64
	Hash       string

	// Client metadata
	DeviceID   string
	ClientOS  string
	ClientVer string

	// Integrity
	Encrypted bool
	Algo      string // e.g. "XChaCha20-Poly1305"

	// Status
	Status string

	// Anchoring
	Anchored      bool
	AnchorTxHash  string
	AnchoredAt    *time.Time

	// Audit
	CreatedAt time.Time
	CompletedAt *time.Time
}

type IpfsCidRequest struct {
	UserID    string
	VaultName string
	CID       string
}
type IpfsCidResponse struct {
    Status  int    `json:"status"`
    Data    string `json:"data"`      // Direct base64 string
    Message string `json:"message"`
    Success bool   `json:"success"`
}
type CloudResponse[T any] struct {
    Status  int    `json:"status"`
    Data    T      `json:"data"`
    Message string `json:"message"`
    Success bool   `json:"success,omitempty"`
}

type WrappedResponse[T any] struct {
	Result T   `json:"result"`
	Error  any `json:"error"`
}

type WrappedResultUser struct {
	Result *User       `json:"result"`
	Error  interface{} `json:"error"`
}
type GetVaultInput struct {
	UserID    string
	VaultName string
}
// AccessCryptoShareRequest holds the parameters for accessing a cryptographic share.
type AccessCryptoShareRequest struct {
	ShareID        string `json:"share_id"`
	RecipientEmail string `json:"recipient_email"`
	Challenge      string `json:"challenge"`
	Signature      string `json:"signature"`
	IPAddress      string `json:"ip_address,omitempty"`
}

// AccessCryptoShareResponse holds the decrypted data returned after accessing a share.
type AccessCryptoShareResponse struct {
	EncryptedKey    string
	SenderPublicKey string
	EncryptedPayload        string
	DownloadAllowed bool
}
type DecryptCryptoShareRequest struct {
	// IPAddress string `json:"ip_address"`
	EncryptedKey        string `json:"encrypted_key"`
	EncryptedPayload    string `json:"encrypted_payload"`
	RecipientPrivateKey string `json:"recipient_private_key"`
}
type DecryptCryptoShareResponse struct {
	Payload string `json:"payload"` // decrypted vault payload
	ExpiresIn int64 `json:"expires_in,omitempty"`
}
type AddPublicKeyToCustomerRequest struct {
	PublicKey string `json:"public_key"`
	Email string `json:"email"`
}
type AddPublicKeyToCustomerResponse struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	PublicKey *string    `json:"public_key,omitempty"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
	

type AddRecipientRequest struct {
	ShareID string `json:"share_id"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	EncryptedKey string `json:"encrypted_key"`
	RevokedAt    *time.Time `json:"revoked_at"`
	Signature    string `json:"signature"`
}

type RevokeShareRequest struct {
	ShareID string `json:"share_id"`
	Challenge      string `json:"challenge"`
	Email    string `json:"email"`
	Signature    string `json:"signature"`
}	

type GetBillingHistory struct {
	ID string `json:"id"`
	UserID string `json:"user_id"`
	SubscriptionID string `json:"subscription_id"`
	Amount int64 `json:"amount"`
	Status string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	Description string `json:"description"`
	StripeIntentID string `json:"stripe_intent_id"`
	CreatedAt time.Time `json:"created_at"`
}

type PaymentHistory struct {
    ID              string    `json:"id" gorm:"primaryKey"`
	UserID          string    `json:"user_id" gorm:"not null"`
	SubscriptionID  string    `json:"subscription_id" gorm:"not null"`
    Amount          float64   `json:"amount" gorm:"not null"`
    Status          string    `json:"status" gorm:"not null"` // succeeded, failed, pending
    PaymentMethod   string    `json:"payment_method" gorm:"not null"`
    Description     string    `json:"description" gorm:"omitempty"`
    StellarTxHash   string    `json:"stellar_tx_hash,omitempty" gorm:"omitempty"`
    StripeIntentID  string    `json:"stripe_intent_id,omitempty" gorm:"omitempty"`
    CreatedAt       time.Time `json:"created_at" gorm:"not null"`
}
type Vault struct {
	ID        string
	OwnerID  string // user or team
	OwnerType string // user | team
	SubscriptionID string

	Name      string
	PlanID    string
	Active    bool

	// Storage
	QuotaBytes int64
	UsedBytes  int64    
	StorageBackend string // ipfs | s3 | hybrid

	// Features (resolved from subscription)
	Features VaultFeatures

	// Blockchain refs
	IPFSNodeID  string
	PinataPinID string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type VaultFeatures struct {
	CloudBackup bool
	Versioning  bool
	Sharing     bool
	Telemetry   bool
	Tracecore   bool
}

type StorageUsageRequest struct {
	UserID string `json:"user_id"`
	Tier string `json:"tier"`
	Challenge string `json:"challenge"`
	Signature string `json:"signature"`
	PublicKey string `json:"public_key"`	
}

type StorageUsageResponse struct {
	BytesUsed int64 `json:"bytes_used"`
	BytesLimit int64 `json:"bytes_limit"`
	UserID string `json:"user_id"`
}	
	
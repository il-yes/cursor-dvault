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
	Sync    VaultSync
	CreatedAt string
	Metadata  map[string]string `json:"metadata,omitempty"`
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
}

// AccessCryptoShareResponse holds the decrypted data returned after accessing a share.
type AccessCryptoShareResponse struct {
	EncryptedKey    string
	SenderPublicKey string
	EncryptedPayload        string
	DownloadAllowed bool
}
type DecryptCryptoShareRequest struct {
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
	

package share_entry_domain

import (
	"time"

	"gorm.io/datatypes"

	vaults_domain "vault-app/internal/vault/domain"
)

// --------------------------------------------------------------------
// Tier 1: Free - Link-Based Shares (Ephemeral)
// --------------------------------------------------------------------
type LinkShare struct {
	ID        string
	Payload   string
	CreatedAt time.Time

	ExpiresAt       *time.Time
	MaxViews        *int
	ViewCount       int
	Password        *string
	DownloadAllowed bool

	CreatorUserID string
	CreatorEmail  string
	Metadata      Metadata
}

type Metadata struct {
	EntryType string
	Title     string
}

// -----------------------
// Cryptography - Aggregate root
// -----------------------
type ShareEntry struct {
	ID               string         `json:"id"`
	OwnerID          string         `json:"owner_id"`
	EntryID			string 			`json:"entry_id"`
	EntryName        string         `json:"entry_name"`
	EntryType        string         `json:"entry_type"`
	EntryRef         string         `json:"entry_ref"`
	Status           string         `json:"status"`
	AccessMode       string         `json:"access_mode"`
	Encryption       string         `json:"encryption"`
	EntrySnapshot    EntrySnapshot  `json:"entry_snapshot"`
	EncryptedPayload string         `json:"encrypted_payload"`
	AccessLog        datatypes.JSON `json:"access_log"`
	ExpiresAt        *time.Time     `json:"expires_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	SharedAt         time.Time      `json:"shared_at"`
	DownloadAllowed  bool           `json:"download_allowed"`

	Recipients     []Recipient `gorm:"foreignKey:ShareID;constraint:OnDelete:CASCADE" json:"recipients"`
}

func NewShareEntry(
	ownerID string,
	entryID string,
	entryName string,
	entryRef string,
	entryType string,
	status string,
	accessMode string,
	encryption string,
	snapshot EntrySnapshot,
	downloadAllowed bool,
	expiresAt string,
) ShareEntry {
	var s ShareEntry
	s.OwnerID = ownerID
	s.EntryID = entryID
	s.EntryName = entryName
	s.EntryRef = entryRef
	s.EntryType = entryType
	s.Status = status
	s.AccessMode = accessMode
	s.Encryption = encryption
	s.EntrySnapshot = snapshot
	s.DownloadAllowed = downloadAllowed

	// parse ExpiresAt if present (payload.ExpiresAt is string)
	if expiresAt != "" {
		t, err := time.Parse(time.RFC3339, expiresAt)
		if err == nil {
			s.ExpiresAt = &t
		} else {
			// optional: handle other formats or ignore
			s.ExpiresAt = nil
		}
	}

	return s
}

// -----------------------
// Recipient
// -----------------------
type Recipient struct {
	ID        string    `json:"id"`
	ShareID   string    `json:"share_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	PublicKey string    `json:"public_key"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joined_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Blob containing encrypted vault snapshot (optional)
	EncryptedBlob []byte    `json:"encrypted_blob"`
	RevokedAt     time.Time `json:"revoked_at"`
}

// -----------------------
// ShareAcceptData
// -----------------------
type ShareAcceptData struct {
	Share     ShareEntry `json:"share"`
	Recipient Recipient  `json:"recipient"`
	Blob      []byte     `json:"blob"` // encrypted snapshot for this recipient
}

// -----------------------
// AuditLog
// -----------------------
type AuditLog struct {
	ID        string    `gorm:"primaryKey;size:64" json:"id"`
	ShareID   string    `gorm:"index;size:64" json:"share_id"`
	Action    string    `json:"action"`
	Actor     string    `json:"actor"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details"`
}

// -----------------------
// EntrySnapshot
// -----------------------
type EntrySnapshot struct {
	EntryName string `json:"entry_name" gorm:"size:255"`
	Type      string `json:"type" gorm:"size:64"`

	// Login
	UserName string `json:"user_name" gorm:"size:255"`
	Password string `json:"password" gorm:"size:255"`
	Website  string `json:"website" gorm:"size:255"`

	// Card
	CardholderName string `json:"cardholder_name" gorm:"size:255"`
	CardNumber     string `json:"card_number" gorm:"size:255"`
	ExpiryMonth    int    `json:"expiry_month"`
	ExpiryYear     int    `json:"expiry_year"`
	CVV            string `json:"cvv" gorm:"size:10"`

	// SSH Key
	PrivateKey string `json:"private_key" gorm:"type:text"`
	PublicKey  string `json:"public_key" gorm:"type:text"`

	// Secure Note
	Note string `json:"note" gorm:"type:text"`

	// Identity
	Genre                string `json:"genre"`
	FirstName            string `json:"firstname"`
	SecondFirstName      string `json:"second_firstname"`
	LastName             string `json:"lastname"`
	Username             string `json:"username"`
	Company              string `json:"company"`
	SocialSecurityNumber string `json:"social_security_number"`
	IDNumber             string `json:"ID_number"`
	DriverLicense        string `json:"driver_license"`
	Mail                 string `json:"mail"`
	Telephone            string `json:"telephone"`
	AddressOne           string `json:"address_one"`
	AddressTwo           string `json:"address_two"`
	AddressThree         string `json:"address_three"`
	City                 string `json:"city"`
	State                string `json:"state"`
	PostalCode           string `json:"postal_code"`
	Country              string `json:"country"`

	// Custom fields fallback
	ExtraFields    datatypes.JSON `json:"extra_fields" gorm:"type:json"`
	AttachmentCIDs []string       `json:"attachmentCIDs"`
	RecipientCIDs  map[string]string `json:"recipient_cids"`
	Attachments []vaults_domain.Attachment `json:"attachments" gorm:"-"`
}

// -----------------------
// RecipientCIDs
// -----------------------
type RecipientCIDs struct {
	FileCID string `json:"file_cid"`
	AttachmentID string `json:"attachment_id"`
}

// -----------------------
// PendingShareIntent
// -----------------------
type PendingShareIntent struct {
	ID         string
	ShareID    string
	OwnerID    string
	OwnerEmail string

	RecipientID    string
	RecipientEmail string

	Status string

	InvitationSentAt     *time.Time
	NotificationAttempts int
	AcceptedAt           *time.Time
	DeclinedAt           *time.Time
	ExpiresAt            *time.Time

	CreatedAt             time.Time
	LastNotificationError string
}
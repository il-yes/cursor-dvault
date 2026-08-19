package share_entry_application_dto

import (
	"time"

	share_entry_domain "vault-app/internal/share_entry/domain"
	"vault-app/internal/vault/domain"
)

type Metadata struct {
	EntryType string `json:"entry_type"`
	Title     string `json:"title"`
}

type RecipientPayload struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	PublicKey    string `json:"public_key"`
	TrustGroupID string `json:"trust_group_id,omitempty"`
	Type         string `json:"type,omitempty"` // "user" (default) or "trust_group"
}
type CreateShareEntryPayload struct {
	EntryID string	`json:"entry_id"`
	EntryName  string `json:"entry_name"`
	EntryType  string `json:"entry_type"`
	EntryRef   string `json:"entry_ref"`
	Status     string `json:"status"`
	AccessMode string `json:"access_mode"`
	Encryption string `json:"encryption"`
	// Snapshot as JSON string from frontend
	EntrySnapshot string `json:"entry_snapshot"`

	ExpiresAt       string             `json:"expires_at"`
	Recipients      []RecipientPayload `json:"recipients"`
	DownloadAllowed bool               `json:"download_allowed"`
	AttachmentCIDs  []string           `json:"attachmentCIDs"`
}
type CreateShareRequest struct {
	UserID     string
	OwnerEmail string
	Payload    CreateShareEntryPayload
	Share      share_entry_domain.ShareEntry
	Secret     string
}
type LinkShareDTO struct {
	ID              string     `json:"id"`
	Payload         string     `json:"payload"`
	CreatedAt       time.Time  `json:"created_at"`
	ExpiresAt       *time.Time `json:"expires_at"`
	MaxViews        *int       `json:"max_views"`
	ViewCount       int        `json:"view_count"`
	PasswordHash    *string    `json:"password_hash"`
	DownloadAllowed bool       `json:"download_allowed"`
	CreatorUserID   string     `json:"creator_user_id"`
	Metadata        Metadata   `json:"metadata"`
}

type LinkShareCreateRequest struct {
	Payload         string     `json:"payload"`
	ExpiresAt       *time.Time `json:"expires_at"`
	MaxViews        *int       `json:"max_views"`
	Password        *string    `json:"password"`
	DownloadAllowed bool       `json:"download_allowed"`
	CreatorEmail    string     `json:"creator_email"`
	EntryType       string     `json:"entry_type"`
	Title           string     `json:"title"`
}
type AddRecipientRequest struct {
	ShareID   string     `json:"share_id"`
	Email     string     `json:"email"`
	PublicKey string     `json:"public_key"`
	Role      string     `json:"role"`
	RevokedAt *time.Time `json:"revoked_at"`
}
type UpdateRecipientRequest struct {
	ShareID   string     `json:"share_id"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	RevokedAt *time.Time `json:"revoked_at"`
}

type AttachementCIDsAdded struct {
	CIDs           []string              `json:"cids"`
	Attachments    []vaults_domain.Attachment `json:"attachments"`
}

type ActivateCryptographicShareResponse struct {
	Cryptoshare share_entry_domain.ShareEntry `json:"cryptoshare"`
	PendingShare share_entry_domain.PendingShareIntent `json:"pending_share"`
}

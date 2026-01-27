package share_application_dto

import "time"

type Metadata struct {
	EntryType string `json:"entry_type"`
	Title string `json:"title"`
}

type LinkShareDTO struct {
	ID string `json:"id"`
	Payload string `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"`
	MaxViews *int `json:"max_views"`
	ViewCount int `json:"view_count"`
	PasswordHash *string `json:"password_hash"`
	DownloadAllowed bool `json:"download_allowed"`
	CreatorUserID string `json:"creator_user_id"`
	Metadata Metadata `json:"metadata"`
}

type LinkShareCreateRequest struct {
	Payload string `json:"payload"`
	ExpiresAt *time.Time `json:"expires_at"`
	MaxViews *int `json:"max_views"`
	Password *string `json:"password"`
	DownloadAllowed bool `json:"download_allowed"`
	CreatorEmail string `json:"creator_email"`
	EntryType string `json:"entry_type"`
	Title string `json:"title"`
}

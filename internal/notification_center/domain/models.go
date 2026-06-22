package notification_center_domain

import "time"

type NotificationStatus string

const (
	Unread  NotificationStatus = "unread"
	Read    NotificationStatus = "read"
	Archived NotificationStatus = "archived"
)

type Notification_ALPHA struct {
	ID string

	UserID string

	Type string

	Title string

	Message string

	Payload string

	Status NotificationStatus

	CreatedAt time.Time
	ReadAt *time.Time
	ArchivedAt *time.Time
}


// internal/notification_center/domain/models.go

type Notification_BETA struct {
	ID         string                 `json:"id"`
	UserID     string                 `json:"user_id"`

	Type       string                 `json:"type"`

	Title      string                 `json:"title"`
	Body       string                 `json:"body"`

	Status     string                 `json:"status"`

	CreatedAt  time.Time              `json:"created_at"`
	ReadAt     *time.Time             `json:"read_at"`
	ArchivedAt *time.Time             `json:"archived_at"`

	Payload    map[string]interface{} `json:"payload"`

	Sequence   uint64 `json:"sequence"`
}

type Category string

const (
	CategoryShare       Category = "share"
	CategorySubscription Category = "subscription"
	CategorySystem      Category = "system"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusSent      Status = "sent"
	StatusFailed    Status = "failed"
	StatusDelivered Status = "delivered"
	StatusRead      Status = "read"
	StatusUnread    Status = "unread"
	StatusArchived  Status = "archived"
)

type Notification struct {
	ID string `gorm:"primarykey;column:id" json:"id"`
	UserID string `gorm:"column:user_id" json:"user_id"`

	Type     string `gorm:"column:type" json:"type"`
	Category Category `gorm:"column:category" json:"category"`

	Title   string `gorm:"column:title" json:"title"`
	Body    string `gorm:"column:body" json:"body"`
	Payload any `gorm:"column:payload" json:"payload"` // raw JSON (your shared_realtime.Message.Payload)

	Status Status `gorm:"column:status" json:"status"`

	EventID string `gorm:"column:event_id" json:"event_id"` // idempotency key

	DeliveryAttempts int `gorm:"column:delivery_attempts" json:"delivery_attempts"`
	LastError        string `gorm:"column:last_error" json:"last_error"`

	DeliveredAt *time.Time `gorm:"column:delivered_at" json:"delivered_at"`
	ReadAt      *time.Time `gorm:"column:read_at" json:"read_at"`

	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	Sequence int64 `gorm:"column:sequence" json:"sequence"`
}
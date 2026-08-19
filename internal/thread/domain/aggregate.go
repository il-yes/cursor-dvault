package thread_domain

import (
	"time"

	"github.com/google/uuid"
)

type ThreadStatus string

const (
	ThreadOpen         ThreadStatus = "open"
	ThreadTransferring ThreadStatus = "transferring"
	ThreadClosed       ThreadStatus = "closed"
)

type Thread struct {
	ID          string       `json:"id"`
	ChannelID   string       `json:"channel_id"`
	WorkspaceID string       `json:"workspace_id,omitempty"`
	AssetType   string       `json:"asset_type"`
	Title       string       `json:"title"`
	Subtitle    string       `json:"subtitle"`
	Status      ThreadStatus `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	ClosedAt    *time.Time   `json:"closed_at,omitempty"`
	IsDraft     bool         `json:"is_draft"`
	IsDirty     bool         `json:"is_dirty" gorm:"boolean"`
}

func NewThread(channelID string, assetType string, title string, subtitle string) Thread {
	return Thread{
		ID:        uuid.NewString(),
		ChannelID: channelID,
		AssetType: assetType,
		Title:     title,
		Subtitle:  subtitle,
		Status:    ThreadOpen,
		CreatedAt: time.Now(),
		IsDraft:   true,
		IsDirty:   false,
	}
}	

type InsertStatus int

const (
	Inserted InsertStatus = iota
	Duplicate
)

type ThreadEventType string

const (
	EventEntryShared           ThreadEventType = "entry.shared"
	EventInvoiceCreated        ThreadEventType = "invoice.created"
	EventFinanceApproved       ThreadEventType = "finance.approved"
	EventPaymentReleased       ThreadEventType = "payment.released"
	EventReceiptIssued         ThreadEventType = "receipt.issued"
	EventFederationEntryShared ThreadEventType = "federation.shared.created"
	EventThreadEventAppended   ThreadEventType = "thread.event.appended"
)

type ResourceType string

const (
	ResourceStorageAsset ResourceType = "storage_asset"
	ResourceShareEntry   ResourceType = "share_entry"
)

type EventResourceRef struct {
	RefType ResourceType `json:"ref_type"`

	// Storage asset reference fields (populated when RefType == ResourceStorageAsset)
	CID         string `json:"cid,omitempty"`
	ContentHash string `json:"content_hash,omitempty"`
	Size        int64  `json:"size,omitempty"`
	AssetType   string `json:"asset_type,omitempty"`

	// C3 ShareEntry reference fields (populated when RefType == ResourceShareEntry)
	ShareEntryID string `json:"share_entry_id,omitempty"`
	TrustGroupID string `json:"trust_group_id,omitempty"`
}

type ThreadEvent struct {
	ID              string
	ThreadID        string
	PreviousEventID *string
	Type            ThreadEventType
	Payload         EventResourceRef
	IdempotencyKey  string
	Cursor          uint64
	Headers         map[string]string
	Signature       string
	CreatedAt       time.Time
}

type Asset struct {
	CID         string `json:"cid"`
	ContentHash string `json:"content_hash"`
	Size        int64  `json:"size"`
	Type        string `json:"type"`
	IsDirty     bool   `json:"is_dirty" gorm:"boolean"`
}
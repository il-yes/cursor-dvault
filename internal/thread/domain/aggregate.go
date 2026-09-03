package thread_domain

import (
	"encoding/json"
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

	// Generic C3 Collaboration Actions
	EventC3ApprovalRequested ThreadEventType = "c3.approval.requested"
	EventC3ApprovalApproved  ThreadEventType = "c3.approval.approved"
	EventC3RejectCreated     ThreadEventType = "c3.reject.created"
	EventC3TransferRequested ThreadEventType = "c3.transfer.requested"
	EventC3TransferApproved  ThreadEventType = "c3.transfer.approved"
	EventC3TransferRejected  ThreadEventType = "c3.transfer.rejected"
	EventC3TransferCompleted ThreadEventType = "c3.transfer.completed"
)

type ResourceType string

const (
	ResourceStorageAsset ResourceType = "storage_asset"
	ResourceShareEntry   ResourceType = "share_entry"
	ResourceC3Action     ResourceType = "c3_action"
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

	// Generic C3 Action reference fields (populated when RefType == ResourceC3Action or C3 Action event)
	ActionID           string `json:"action_id,omitempty"`
	ResourceTypeVal    string `json:"resource_type,omitempty"`
	ResourceIDVal      string `json:"resource_id,omitempty"`
	SourceEventID      string `json:"source_event_id,omitempty"`
	TargetTrustGroupID string `json:"target_trust_group_id,omitempty"`
}

func (r *EventResourceRef) UnmarshalJSON(data []byte) error {
	type Alias EventResourceRef
	var aux struct {
		Alias
		ShareEntryRef *struct {
			ShareEntryID string `json:"share_entry_id"`
			TrustGroupID string `json:"trust_group_id"`
		} `json:"share_entry_ref"`
		ResourceRef *struct {
			ShareEntryID string `json:"share_entry_id"`
			TrustGroupID string `json:"trust_group_id"`
		} `json:"resource_ref"`
		ShareID string `json:"share_id"`
		EntryID string `json:"entry_id"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	*r = EventResourceRef(aux.Alias)

	if r.ShareEntryID == "" {
		if aux.ShareEntryRef != nil && aux.ShareEntryRef.ShareEntryID != "" {
			r.ShareEntryID = aux.ShareEntryRef.ShareEntryID
			if r.TrustGroupID == "" && aux.ShareEntryRef.TrustGroupID != "" {
				r.TrustGroupID = aux.ShareEntryRef.TrustGroupID
			}
		} else if aux.ResourceRef != nil && aux.ResourceRef.ShareEntryID != "" {
			r.ShareEntryID = aux.ResourceRef.ShareEntryID
			if r.TrustGroupID == "" && aux.ResourceRef.TrustGroupID != "" {
				r.TrustGroupID = aux.ResourceRef.TrustGroupID
			}
		} else if aux.ShareID != "" {
			r.ShareEntryID = aux.ShareID
		} else if aux.EntryID != "" {
			r.ShareEntryID = aux.EntryID
		}
	}

	if r.ShareEntryID != "" && r.RefType == "" {
		r.RefType = ResourceShareEntry
	}

	return nil
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
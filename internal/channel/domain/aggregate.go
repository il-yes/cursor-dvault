package channel_domain

import (
	"time"

	"github.com/google/uuid"
)

type ChannelStatus string

const (
	StatusPending  ChannelStatus = "pending"
	StatusActive   ChannelStatus = "active"
	StatusRevoked  ChannelStatus = "revoked"
	StatusArchived ChannelStatus = "archived"
)

type Channel struct {
	ID          string            `json:"id"`
	TemplateID  string            `json:"template_id"`
	Title       string            `json:"title"`
	Status      ChannelStatus     `json:"status"`
	Slots       []Slot            `json:"slots"`	// "roles that exist"
	Assignments []Assignment      `json:"assignments"`	// "who occupies those roles"
	Properties  []ChannelProperty `json:"properties"`
	Policy      Policy            `json:"policy"`
	Federation  string            `json:"federation"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	RevokedAt   *time.Time        `json:"revoked_at,omitempty"`
	ArchivedAt  *time.Time        `json:"archived_at,omitempty"`
	WorkspaceID string            `json:"workspace_id"`
	IsDraft     bool              `json:"is_draft"`
	IsDirty     bool              `json:"is_dirty" gorm:"boolean"`
}

type Policy map[string]interface{}

type Slot struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Role    string `json:"role"`
	VaultID string `json:"vault_id"`
	Gated   bool   `json:"gated"`
	Order   int    `json:"order"`
}

type Assignment struct {
	SlotID       string `json:"slot_id"`
	OwnerID      string `json:"owner_id"`
	PublicKey    string `json:"public_key"`
	VaultAddress string `json:"vault_address"`
}

type ChannelProperty struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewChannel(tplID string, title string, workspaceID string) Channel {
	now := time.Now()
	return Channel{
		ID:          uuid.NewString(),
		TemplateID:  tplID,
		Title:       title,
		Status:      StatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
		WorkspaceID: workspaceID,
		IsDraft:     true,
		IsDirty:     false,
	}
}


type TrustState string

var (
	Trusted = "trusted"
)

type FederationSnapshot struct {
	RemoteVaults []RemoteVault `json:"remote_vaults"`
	IsDraft      bool `json:"is_draft"`
	IsDirty      bool `json:"is_dirty" gorm:"boolean"`
}
type RemoteVault struct {
	VaultID         string	`json:"vault_id"`
	LastCursor      uint64	`json:"last_cursor"`
	LastSeen        time.Time	`json:"last_seen"`
	Endpoint        string	`json:"endpoint"`
	TrustState      TrustState	`json:"trust_state"`
	Pending         []PendingSyncItem	`json:"pending"`
	ProtocolVersion string	`json:"protocol_version"`
}

type SyncStatus string

const (
	SyncPending    SyncStatus = "pending"
	SyncSending    SyncStatus = "sending"
	SyncWaitingAck SyncStatus = "waiting_ack"
	SyncCompleted  SyncStatus = "completed"
	SyncFailed     SyncStatus = "failed"
	SyncDeadLetter SyncStatus = "dead_letter"
)

type PendingSyncItem struct {
	ExchangeID string `json:"exchange_id"`
	Status     SyncStatus `json:"status"`
	Cursor     uint64 `json:"cursor"`
	RetryCount int `json:"retry_count"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Participant struct {
	ChannelID   string  `json:"channel_id"`
	VaultID     string `json:"vault_id"`
	PublicKey   string `json:"public_key"`
	Direction   string `json:"direction"` // inbound | outbound | bidirectional
	JoinedAt    int64 `json:"joined_at"`
	Role        string `json:"role"`
	Permissions []string `json:"permissions"`
	IsDraft     bool `json:"is_draft"`
	IsDirty     bool `json:"is_dirty" gorm:"boolean"`
}
// ==============================================================================
// Lifecycle
// ==============================================================================

// Archive transitions a Channel from active to archived.
// Only active channels can be archived. The aggregate enforces this invariant.
// NOTE: ChannelArchived is a candidate for promotion to integration event
// when TraceCore lifecycle recording is integrated.
func (c *Channel) Archive() error {
	if c.Status != StatusActive {
		return ErrChannelNotArchivable
	}

	now := time.Now().UTC()
	c.Status = StatusArchived
	c.ArchivedAt = &now
	c.UpdatedAt = now

	return nil
}

// ==============================================================================
// Slot CRUD
// ==============================================================================

func (c *Channel) canModify() bool {
	return c.Status == StatusPending ||
		c.Status == StatusActive
}

func (c *Channel) AddSlot(slot Slot) error {
	if !c.canModify() {
		return ErrChannelNotModifiable
	}

	if slot.ID == "" {
		return ErrSlotNotFound
	}

	if _, exists := c.GetSlotByID(slot.ID); exists {
		return ErrSlotAlreadyExists
	}

	c.Slots = append(c.Slots, slot)
	c.UpdatedAt = time.Now().UTC()

	return nil
}

func (c *Channel) UpdateSlot(updated Slot) error {
	if !c.canModify() {
		return ErrChannelNotModifiable
	}

	for i := range c.Slots {
		if c.Slots[i].ID != updated.ID {
			continue
		}

		// Slot identity is immutable.
		if c.Slots[i].ID != updated.ID {
			return ErrSlotIdentityMismatch
		}

		c.Slots[i] = updated
		c.UpdatedAt = time.Now().UTC()

		return nil
	}

	return ErrSlotNotFound
}

func (c *Channel) RemoveSlotByID(id string) error {
	if !c.canModify() {
		return ErrChannelNotModifiable
	}

	for i := range c.Slots {
		if c.Slots[i].ID != id {
			continue
		}

		// Removing a Slot also removes its Assignment(s).
		c.removeAssignmentsBySlotID(id)

		c.Slots = append(
			c.Slots[:i],
			c.Slots[i+1:]...,
		)

		c.UpdatedAt = time.Now().UTC()

		return nil
	}

	return ErrSlotNotFound
}

func (c *Channel) removeAssignmentsBySlotID(slotID string) {
	filtered := c.Assignments[:0]

	for _, assignment := range c.Assignments {
		if assignment.SlotID != slotID {
			filtered = append(filtered, assignment)
		}
	}

	c.Assignments = filtered
}
// func (c *Channel) AddSlot(slot Slot) Channel {
// 	c.Slots = append(c.Slots, slot)
// 	c.UpdatedAt = time.Now().UTC()

// 	return *c
// }

func (c *Channel) GetSlotByID(id string) (Slot, bool) {
	for i := range c.Slots {
		if c.Slots[i].ID == id {
			return c.Slots[i], true
		}
	}

	return Slot{}, false
}

func (c *Channel) GetSlotByVaultID(vaultID string) (*Slot, bool) {
	for i := range c.Slots {
		if c.Slots[i].VaultID == vaultID {
			return &c.Slots[i], true
		}
	}

	return nil, false
}

func (c *Channel) GetSlotsByRole(role string) []Slot {
	slots := make([]Slot, 0)

	for _, slot := range c.Slots {
		if slot.Role == role {
			slots = append(slots, slot)
		}
	}

	return slots
}

func (c *Channel) GetGatedSlots() []Slot {
	slots := make([]Slot, 0)

	for _, slot := range c.Slots {
		if slot.Gated {
			slots = append(slots, slot)
		}
	}

	return slots
}

// func (c *Channel) UpdateSlot(updated Slot) bool {
// 	for i := range c.Slots {
// 		if c.Slots[i].ID == updated.ID {
// 			c.Slots[i] = updated
// 			c.UpdatedAt = time.Now().UTC()

// 			return true
// 		}
// 	}

// 	return false
// }

// func (c *Channel) RemoveSlotByID(id string) bool {
// 	for i := range c.Slots {
// 		if c.Slots[i].ID == id {
// 			c.Slots = append(
// 				c.Slots[:i],
// 				c.Slots[i+1:]...,
// 			)

// 			c.UpdatedAt = time.Now().UTC()

// 			return true
// 		}
// 	}

// 	return false
// }

func (c *Channel) RemoveSlotByVaultID(vaultID string) bool {
	for i := range c.Slots {
		if c.Slots[i].VaultID == vaultID {
			c.Slots = append(
				c.Slots[:i],
				c.Slots[i+1:]...,
			)

			c.UpdatedAt = time.Now().UTC()

			return true
		}
	}

	return false
}

func (c *Channel) ListSlots() []Slot {
	slots := make([]Slot, len(c.Slots))
	copy(slots, c.Slots)

	return slots
}

// ==============================================================================
// Assignment CRUD
// ==============================================================================

func (c *Channel) AddAssignment(assignment Assignment) Channel {
	c.Assignments = append(c.Assignments, assignment)
	c.UpdatedAt = time.Now().UTC()

	return *c
}

func (c *Channel) GetAssignmentBySlotID(slotID string) (*Assignment, bool) {
	for i := range c.Assignments {
		if c.Assignments[i].SlotID == slotID {
			return &c.Assignments[i], true
		}
	}

	return nil, false
}

func (c *Channel) GetAssignmentsByOwnerID(ownerID string) []Assignment {
	assignments := make([]Assignment, 0)

	for _, assignment := range c.Assignments {
		if assignment.OwnerID == ownerID {
			assignments = append(assignments, assignment)
		}
	}

	return assignments
}

func (c *Channel) GetAssignmentByVaultAddress(
	vaultAddress string,
) (*Assignment, bool) {
	for i := range c.Assignments {
		if c.Assignments[i].VaultAddress == vaultAddress {
			return &c.Assignments[i], true
		}
	}

	return nil, false
}

func (c *Channel) UpdateAssignment(updated Assignment) bool {
	for i := range c.Assignments {
		if c.Assignments[i].SlotID == updated.SlotID {
			c.Assignments[i] = updated
			c.UpdatedAt = time.Now().UTC()

			return true
		}
	}

	return false
}

func (c *Channel) RemoveAssignmentBySlotID(slotID string) bool {
	for i := range c.Assignments {
		if c.Assignments[i].SlotID == slotID {
			c.Assignments = append(
				c.Assignments[:i],
				c.Assignments[i+1:]...,
			)

			c.UpdatedAt = time.Now().UTC()

			return true
		}
	}

	return false
}

func (c *Channel) RemoveAssignmentsByOwnerID(ownerID string) int {
	removed := 0
	filtered := make([]Assignment, 0, len(c.Assignments))

	for _, assignment := range c.Assignments {
		if assignment.OwnerID == ownerID {
			removed++
			continue
		}

		filtered = append(filtered, assignment)
	}

	if removed > 0 {
		c.Assignments = filtered
		c.UpdatedAt = time.Now().UTC()
	}

	return removed
}

func (c *Channel) ListAssignments() []Assignment {
	assignments := make([]Assignment, len(c.Assignments))
	copy(assignments, c.Assignments)

	return assignments
}

// ==============================================================================
// Policy CRUD
// ==============================================================================

func (c *Channel) SetPolicy(policy Policy) Channel {
	c.Policy = clonePolicy(policy)
	c.UpdatedAt = time.Now().UTC()

	return *c
}

func (c *Channel) GetPolicy() Policy {
	return clonePolicy(c.Policy)
}

func (c *Channel) HasPolicy() bool {
	return c.Policy != nil
}

func (c *Channel) SetPolicyValue(key string, value interface{}) bool {
	if key == "" {
		return false
	}

	if c.Policy == nil {
		c.Policy = make(Policy)
	}

	c.Policy[key] = value
	c.UpdatedAt = time.Now().UTC()

	return true
}

func (c *Channel) GetPolicyValue(key string) (interface{}, bool) {
	if c.Policy == nil {
		return nil, false
	}

	value, exists := c.Policy[key]
	return value, exists
}

func (c *Channel) UpdatePolicyValue(key string, value interface{}) bool {
	if c.Policy == nil {
		return false
	}

	if _, exists := c.Policy[key]; !exists {
		return false
	}

	c.Policy[key] = value
	c.UpdatedAt = time.Now().UTC()

	return true
}

func (c *Channel) RemovePolicyValue(key string) bool {
	if c.Policy == nil {
		return false
	}

	if _, exists := c.Policy[key]; !exists {
		return false
	}

	delete(c.Policy, key)
	c.UpdatedAt = time.Now().UTC()

	return true
}

func (c *Channel) ClearPolicy() bool {
	if c.Policy == nil {
		return false
	}

	c.Policy = nil
	c.UpdatedAt = time.Now().UTC()

	return true
}

func (c *Channel) ListPolicyKeys() []string {
	if c.Policy == nil {
		return []string{}
	}

	keys := make([]string, 0, len(c.Policy))

	for key := range c.Policy {
		keys = append(keys, key)
	}

	return keys
}

func clonePolicy(policy Policy) Policy {
	if policy == nil {
		return nil
	}

	cloned := make(Policy, len(policy))

	for key, value := range policy {
		cloned[key] = value
	}

	return cloned
}

// ==============================================================================
// ChannelProperty CRUD
// ==============================================================================

func (c *Channel) AddChannelProperty(property ChannelProperty) Channel {
	c.Properties = append(c.Properties, property)
	c.UpdatedAt = time.Now().UTC()

	return *c
}

func (c *Channel) GetChannelProperty(
	key string,
) (*ChannelProperty, bool) {
	for i := range c.Properties {
		if c.Properties[i].Key == key {
			return &c.Properties[i], true
		}
	}

	return nil, false
}

func (c *Channel) UpdateChannelProperty(
	updated ChannelProperty,
) bool {
	for i := range c.Properties {
		if c.Properties[i].Key == updated.Key {
			c.Properties[i] = updated
			c.UpdatedAt = time.Now().UTC()

			return true
		}
	}

	return false
}

func (c *Channel) RemoveChannelProperty(key string) bool {
	for i := range c.Properties {
		if c.Properties[i].Key == key {
			c.Properties = append(
				c.Properties[:i],
				c.Properties[i+1:]...,
			)

			c.UpdatedAt = time.Now().UTC()

			return true
		}
	}

	return false
}

func (c *Channel) ListChannelProperties() []ChannelProperty {
	properties := make(
		[]ChannelProperty,
		len(c.Properties),
	)

	copy(properties, c.Properties)

	return properties
}
func (c *Channel) AddChannelPropertyIfAbsent(
	property ChannelProperty,
) bool {
	if property.Key == "" {
		return false
	}

	if _, exists := c.GetChannelProperty(property.Key); exists {
		return false
	}

	c.Properties = append(c.Properties, property)
	c.UpdatedAt = time.Now().UTC()

	return true
}

// ==============================================================================
// RemoteVault CRUD
// ==============================================================================

func (f *FederationSnapshot) AddRemoteVault(remote RemoteVault) FederationSnapshot {
	f.RemoteVaults = append(f.RemoteVaults, remote)
	f.IsDirty = true

	return *f
}

func (f *FederationSnapshot) GetRemoteVaultByID(
	vaultID string,
) (*RemoteVault, bool) {
	for i := range f.RemoteVaults {
		if f.RemoteVaults[i].VaultID == vaultID {
			return &f.RemoteVaults[i], true
		}
	}

	return nil, false
}

func (f *FederationSnapshot) UpdateRemoteVault(
	updated RemoteVault,
) bool {
	for i := range f.RemoteVaults {
		if f.RemoteVaults[i].VaultID == updated.VaultID {
			f.RemoteVaults[i] = updated
			f.IsDirty = true

			return true
		}
	}

	return false
}

func (f *FederationSnapshot) RemoveRemoteVaultByID(
	vaultID string,
) bool {
	for i := range f.RemoteVaults {
		if f.RemoteVaults[i].VaultID == vaultID {
			f.RemoteVaults = append(
				f.RemoteVaults[:i],
				f.RemoteVaults[i+1:]...,
			)
			f.IsDirty = true

			return true
		}
	}

	return false
}

func (f *FederationSnapshot) ListRemoteVaults() []RemoteVault {
	remoteVaults := make(
		[]RemoteVault,
		len(f.RemoteVaults),
	)

	copy(remoteVaults, f.RemoteVaults)

	return remoteVaults
}

func (f *FederationSnapshot) GetRemoteVaultsByTrustState(state TrustState) []RemoteVault {
	remoteVaults := make([]RemoteVault, 0)

	for _, remote := range f.RemoteVaults {
		if remote.TrustState == state {
			remoteVaults = append(remoteVaults, remote)
		}
	}

	return remoteVaults
}

func (f *FederationSnapshot) GetTrustedRemoteVaults() []RemoteVault {
	return f.GetRemoteVaultsByTrustState(TrustState(Trusted))
}

func (f *FederationSnapshot) GetRemoteVaultsWithPendingSync() []RemoteVault {
	remoteVaults := make([]RemoteVault, 0)

	for _, remote := range f.RemoteVaults {
		if len(remote.Pending) > 0 {
			remoteVaults = append(remoteVaults, remote)
		}
	}

	return remoteVaults
}
func (f *FederationSnapshot) SetRemoteVaultTrustState(
	vaultID string,
	state TrustState,
) bool {
	remote, found := f.GetRemoteVaultByID(vaultID)
	if !found {
		return false
	}

	remote.TrustState = state
	remote.LastSeen = time.Now().UTC()
	f.IsDirty = true

	return true
}

func (f *FederationSnapshot) UpdateRemoteVaultCursor(
	vaultID string,
	cursor uint64,
) bool {
	remote, found := f.GetRemoteVaultByID(vaultID)
	if !found {
		return false
	}

	remote.LastCursor = cursor
	remote.LastSeen = time.Now().UTC()
	f.IsDirty = true

	return true
}

func (f *FederationSnapshot) UpdateRemoteVaultEndpoint(
	vaultID string,
	endpoint string,
) bool {
	remote, found := f.GetRemoteVaultByID(vaultID)
	if !found {
		return false
	}

	remote.Endpoint = endpoint
	remote.LastSeen = time.Now().UTC()
	f.IsDirty = true

	return true
}
// ==============================================================================
// PendingSyncItem CRUD
// ==============================================================================

func (r *RemoteVault) AddPendingSyncItem(item PendingSyncItem) RemoteVault {
	if item.UpdatedAt.IsZero() {
		item.UpdatedAt = time.Now().UTC()
	}

	r.Pending = append(r.Pending, item)

	return *r
}

func (r *RemoteVault) GetPendingSyncItemByID(exchangeID string) (*PendingSyncItem, bool) {
	for i := range r.Pending {
		if r.Pending[i].ExchangeID == exchangeID {
			return &r.Pending[i], true
		}
	}

	return nil, false
}

func (r *RemoteVault) UpdatePendingSyncItem(updated PendingSyncItem) bool {
	for i := range r.Pending {
		if r.Pending[i].ExchangeID == updated.ExchangeID {
			updated.UpdatedAt = time.Now().UTC()
			r.Pending[i] = updated

			return true
		}
	}

	return false
}

func (r *RemoteVault) RemovePendingSyncItem(exchangeID string) bool {
	for i := range r.Pending {
		if r.Pending[i].ExchangeID == exchangeID {
			r.Pending = append(
				r.Pending[:i],
				r.Pending[i+1:]...,
			)

			return true
		}
	}

	return false
}

func (r *RemoteVault) ListPendingSyncItems() []PendingSyncItem {
	items := make([]PendingSyncItem, len(r.Pending))
	copy(items, r.Pending)

	return items
}
func (r *RemoteVault) GetPendingSyncItemsByStatus(status SyncStatus) []PendingSyncItem {
	items := make([]PendingSyncItem, 0)

	for _, item := range r.Pending {
		if item.Status == status {
			items = append(items, item)
		}
	}

	return items
}

func (r *RemoteVault) GetFailedSyncItems() []PendingSyncItem {
	return r.GetPendingSyncItemsByStatus(SyncFailed)
}

func (r *RemoteVault) GetDeadLetterSyncItems() []PendingSyncItem {
	return r.GetPendingSyncItemsByStatus(SyncDeadLetter)
}

func (r *RemoteVault) HasPendingSyncItems() bool {
	return len(r.Pending) > 0
}
func (f *FederationSnapshot) GetPendingSyncItem(vaultID string, exchangeID string) (*PendingSyncItem, bool) {
	remote, found := f.GetRemoteVaultByID(vaultID)
	if !found {
		return nil, false
	}

	return remote.GetPendingSyncItemByID(exchangeID)
}

func (f *FederationSnapshot) AddPendingSyncItem(vaultID string, item PendingSyncItem) bool {
	remote, found := f.GetRemoteVaultByID(vaultID)
	if !found {
		return false
	}

	remote.AddPendingSyncItem(item)
	f.IsDirty = true

	return true
}

func (f *FederationSnapshot) UpdatePendingSyncItem(vaultID string, item PendingSyncItem) bool {
	remote, found := f.GetRemoteVaultByID(vaultID)
	if !found {
		return false
	}

	updated := remote.UpdatePendingSyncItem(item)
	if updated {
		f.IsDirty = true
	}

	return updated
}

func (f *FederationSnapshot) RemovePendingSyncItem(vaultID string, exchangeID string) bool {
	remote, found := f.GetRemoteVaultByID(vaultID)
	if !found {
		return false
	}

	removed := remote.RemovePendingSyncItem(exchangeID)
	if removed {
		f.IsDirty = true
	}

	return removed
}
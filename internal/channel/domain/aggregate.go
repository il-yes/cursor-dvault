package channel_domain

import (
	"errors"
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
	Slots       []Slot            `json:"slots"`
	Assignments []Assignment      `json:"assignments"`
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

// ==============================================================================
// Lifecycle
// ==============================================================================

// Archive transitions a Channel from active to archived.
// Only active channels can be archived. The aggregate enforces this invariant.
// NOTE: ChannelArchived is a candidate for promotion to integration event
// when TraceCore lifecycle recording is integrated.
func (c *Channel) Archive() error {
	if c.Status != StatusActive {
		return errors.New(ErrChannelNotArchivable)
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

func (c *Channel) AddSlot(slot Slot) Channel {
	c.Slots = append(c.Slots, slot)
	c.UpdatedAt = time.Now().UTC()

	return *c
}

func (c *Channel) GetSlotByID(id string) (*Slot, bool) {
	for i := range c.Slots {
		if c.Slots[i].ID == id {
			return &c.Slots[i], true
		}
	}

	return nil, false
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

func (c *Channel) UpdateSlot(updated Slot) bool {
	for i := range c.Slots {
		if c.Slots[i].ID == updated.ID {
			c.Slots[i] = updated
			c.UpdatedAt = time.Now().UTC()

			return true
		}
	}

	return false
}

func (c *Channel) RemoveSlotByID(id string) bool {
	for i := range c.Slots {
		if c.Slots[i].ID == id {
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
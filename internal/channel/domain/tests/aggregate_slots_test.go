package channel_tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	channel_domain "vault-app/internal/channel/domain"
)

func newTestChannel() channel_domain.Channel {
	return channel_domain.Channel{
		ID:          "channel-001",
		TemplateID:  "design-review",
		Title:       "Battery System Engineering",
		Status:      channel_domain.StatusActive,
		WorkspaceID: "workspace-001",
		Slots:       []channel_domain.Slot{},
	}
}

func testSlotOne() channel_domain.Slot {
	return channel_domain.Slot{
		ID:      "slot-001",
		Name:    "Battery Engineering",
		Role:    "lead-engineer",
		VaultID: "vault-oem",
		Gated:   true,
		Order:   1,
	}
}

func testSlotTwo() channel_domain.Slot {
	return channel_domain.Slot{
		ID:      "slot-002",
		Name:    "Safety Review",
		Role:    "reviewer",
		VaultID: "vault-regulator",
		Gated:   false,
		Order:   2,
	}
}

func TestChannel_AddSlot_GetSlotByID(t *testing.T) {
	channel := newTestChannel()
	slot := testSlotOne()

	result := channel.AddSlot(slot)

	require.Len(t, result.Slots, 1)
	require.Equal(t, slot, result.Slots[0])

	found, ok := result.GetSlotByID("slot-001")

	require.True(t, ok)
	require.NotNil(t, found)
	require.Equal(t, slot, *found)
}
func TestChannel_GetSlotByVaultID(t *testing.T) {
	channel := newTestChannel()

	channel = channel.AddSlot(testSlotOne())
	channel = channel.AddSlot(testSlotTwo())

	found, ok := channel.GetSlotByVaultID("vault-regulator")

	require.True(t, ok)
	require.NotNil(t, found)
	require.Equal(t, "slot-002", found.ID)
	require.Equal(t, "Safety Review", found.Name)
}
func TestChannel_GetSlotByID_NotFound(t *testing.T) {
	channel := newTestChannel()

	found, ok := channel.GetSlotByID("unknown-slot")

	require.False(t, ok)
	require.Nil(t, found)
}
func TestChannel_GetSlotByVaultID_NotFound(t *testing.T) {
	channel := newTestChannel()

	found, ok := channel.GetSlotByVaultID("unknown-vault")

	require.False(t, ok)
	require.Nil(t, found)
}
func TestChannel_GetSlotsByRole(t *testing.T) {
	channel := newTestChannel()

	channel = channel.AddSlot(channel_domain.Slot{
		ID:      "slot-001",
		Name:    "Thermal Engineering",
		Role:    "engineer",
		VaultID: "vault-oem",
		Order:   1,
	})

	channel = channel.AddSlot(channel_domain.Slot{
		ID:      "slot-002",
		Name:    "Electrical Engineering",
		Role:    "engineer",
		VaultID: "vault-supplier",
		Order:   2,
	})

	channel = channel.AddSlot(channel_domain.Slot{
		ID:      "slot-003",
		Name:    "Regulatory Review",
		Role:    "reviewer",
		VaultID: "vault-regulator",
		Order:   3,
	})

	slots := channel.GetSlotsByRole("engineer")

	require.Len(t, slots, 2)
	require.Equal(t, "slot-001", slots[0].ID)
	require.Equal(t, "slot-002", slots[1].ID)
}
func TestChannel_GetSlotsByRole_NotFound(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddSlot(testSlotOne())

	slots := channel.GetSlotsByRole("unknown-role")

	require.NotNil(t, slots)
	require.Empty(t, slots)
}
func TestChannel_GetGatedSlots(t *testing.T) {
	channel := newTestChannel()

	channel = channel.AddSlot(channel_domain.Slot{
		ID:    "slot-001",
		Name:  "Lead Engineer",
		Role:  "lead-engineer",
		Gated: true,
	})

	channel = channel.AddSlot(channel_domain.Slot{
		ID:    "slot-002",
		Name:  "Observer",
		Role:  "observer",
		Gated: false,
	})

	channel = channel.AddSlot(channel_domain.Slot{
		ID:    "slot-003",
		Name:  "Regulator",
		Role:  "reviewer",
		Gated: true,
	})

	slots := channel.GetGatedSlots()

	require.Len(t, slots, 2)
	require.Equal(t, "slot-001", slots[0].ID)
	require.Equal(t, "slot-003", slots[1].ID)
}
func TestChannel_GetGatedSlots_Empty(t *testing.T) {
	channel := newTestChannel()

	slots := channel.GetGatedSlots()

	require.NotNil(t, slots)
	require.Empty(t, slots)
}
func TestChannel_UpdateSlot(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddSlot(testSlotOne())

	updated := testSlotOne()
	updated.Name = "Principal Battery Engineer"
	updated.Role = "principal-engineer"
	updated.Gated = false
	updated.Order = 10

	updatedResult := channel.UpdateSlot(updated)

	require.True(t, updatedResult)

	found, ok := channel.GetSlotByID("slot-001")

	require.True(t, ok)
	require.Equal(t, "Principal Battery Engineer", found.Name)
	require.Equal(t, "principal-engineer", found.Role)
	require.False(t, found.Gated)
	require.Equal(t, 10, found.Order)
}
func TestChannel_UpdateSlot_NotFound(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddSlot(testSlotOne())

	updated := channel_domain.Slot{
		ID:      "unknown-slot",
		Name:    "Unknown",
		Role:    "unknown",
		VaultID: "unknown-vault",
	}

	updatedResult := channel.UpdateSlot(updated)

	require.False(t, updatedResult)
	require.Len(t, channel.Slots, 1)

	found, ok := channel.GetSlotByID("slot-001")

	require.True(t, ok)
	require.Equal(t, testSlotOne(), *found)
}
func TestChannel_RemoveSlotByID(t *testing.T) {
	channel := newTestChannel()

	channel = channel.AddSlot(testSlotOne())
	channel = channel.AddSlot(testSlotTwo())

	removed := channel.RemoveSlotByID("slot-001")

	require.True(t, removed)
	require.Len(t, channel.Slots, 1)
	require.Equal(t, "slot-002", channel.Slots[0].ID)

	_, found := channel.GetSlotByID("slot-001")
	require.False(t, found)
}
func TestChannel_RemoveSlotByID_NotFound(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddSlot(testSlotOne())

	removed := channel.RemoveSlotByID("unknown-slot")

	require.False(t, removed)
	require.Len(t, channel.Slots, 1)
}
func TestChannel_RemoveSlotByVaultID(t *testing.T) {
	channel := newTestChannel()

	channel = channel.AddSlot(testSlotOne())
	channel = channel.AddSlot(testSlotTwo())

	removed := channel.RemoveSlotByVaultID("vault-oem")

	require.True(t, removed)
	require.Len(t, channel.Slots, 1)
	require.Equal(t, "vault-regulator", channel.Slots[0].VaultID)
}
func TestChannel_RemoveSlotByVaultID_NotFound(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddSlot(testSlotOne())

	removed := channel.RemoveSlotByVaultID("unknown-vault")

	require.False(t, removed)
	require.Len(t, channel.Slots, 1)
}
func TestChannel_ListSlots(t *testing.T) {
	channel := newTestChannel()

	slotOne := testSlotOne()
	slotTwo := testSlotTwo()

	channel = channel.AddSlot(slotOne)
	channel = channel.AddSlot(slotTwo)

	slots := channel.ListSlots()

	require.Len(t, slots, 2)
	require.Equal(t, slotOne, slots[0])
	require.Equal(t, slotTwo, slots[1])
}
func TestChannel_ListSlots_ReturnsCopy(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddSlot(testSlotOne())

	slots := channel.ListSlots()

	require.Len(t, slots, 1)

	slots[0].Name = "Changed Outside Channel"
	slots = append(slots, testSlotTwo())

	require.Len(t, channel.Slots, 1)

	original, ok := channel.GetSlotByID("slot-001")

	require.True(t, ok)
	require.Equal(t, "Battery Engineering", original.Name)
}
func TestChannel_SlotCRUD_EmptyChannel(t *testing.T) {
	channel := newTestChannel()

	require.Empty(t, channel.ListSlots())
	require.Empty(t, channel.GetSlotsByRole("engineer"))
	require.Empty(t, channel.GetGatedSlots())

	slotByID, foundByID := channel.GetSlotByID("slot-001")
	require.False(t, foundByID)
	require.Nil(t, slotByID)

	slotByVault, foundByVault := channel.GetSlotByVaultID("vault-oem")
	require.False(t, foundByVault)
	require.Nil(t, slotByVault)

	require.False(t, channel.UpdateSlot(testSlotOne()))
	require.False(t, channel.RemoveSlotByID("slot-001"))
	require.False(t, channel.RemoveSlotByVaultID("vault-oem"))
}

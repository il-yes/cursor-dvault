package channel_tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	channel_domain "vault-app/internal/channel/domain"
)

func testAssignmentOne() channel_domain.Assignment {
	return channel_domain.Assignment{
		SlotID:       "slot-001",
		OwnerID:      "owner-001",
		PublicKey:    "GBUZX-OWNER-001",
		VaultAddress: "owner-vault@oem.ankhora",
	}
}

func testAssignmentTwo() channel_domain.Assignment {
	return channel_domain.Assignment{
		SlotID:       "slot-002",
		OwnerID:      "owner-002",
		PublicKey:    "GBUZX-OWNER-002",
		VaultAddress: "owner-vault@partner.ankhora",
	}
}

func testAssignmentThree() channel_domain.Assignment {
	return channel_domain.Assignment{
		SlotID:       "slot-003",
		OwnerID:      "owner-001",
		PublicKey:    "GBUZX-OWNER-001-SECOND",
		VaultAddress: "second-vault@oem.ankhora",
	}
}

func TestChannel_AddAssignment_GetAssignmentBySlotID(t *testing.T) {
	channel := newTestChannel()
	assignment := testAssignmentOne()

	result := channel.AddAssignment(assignment)

	require.Len(t, result.Assignments, 1)
	require.Equal(t, assignment, result.Assignments[0])

	found, ok := result.GetAssignmentBySlotID("slot-001")

	require.True(t, ok)
	require.NotNil(t, found)
	require.Equal(t, assignment, *found)
}

func TestChannel_GetAssignmentBySlotID_NotFound(t *testing.T) {
	channel := newTestChannel()

	found, ok := channel.GetAssignmentBySlotID("unknown-slot")

	require.False(t, ok)
	require.Nil(t, found)
}

func TestChannel_RemoveAssignmentsByOwnerID(t *testing.T) {
	channel := newTestChannel()

	channel = channel.AddAssignment(testAssignmentOne())
	channel = channel.AddAssignment(testAssignmentTwo())
	channel = channel.AddAssignment(testAssignmentThree())

	removed := channel.RemoveAssignmentsByOwnerID("owner-001")

	require.Equal(t, 2, removed)
	require.Len(t, channel.Assignments, 1)
	require.Equal(t, "slot-002", channel.Assignments[0].SlotID)
	require.Equal(t, "owner-002", channel.Assignments[0].OwnerID)
}

func TestChannel_GetAssignmentsByOwnerID_NotFound(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddAssignment(testAssignmentOne())

	assignments := channel.GetAssignmentsByOwnerID("unknown-owner")

	require.NotNil(t, assignments)
	require.Empty(t, assignments)
}

func TestChannel_GetAssignmentByVaultAddress(t *testing.T) {
	channel := newTestChannel()

	channel = channel.AddAssignment(testAssignmentOne())
	channel = channel.AddAssignment(testAssignmentTwo())

	found, ok := channel.GetAssignmentByVaultAddress(
		"owner-vault@partner.ankhora",
	)

	require.True(t, ok)
	require.NotNil(t, found)
	require.Equal(t, "slot-002", found.SlotID)
	require.Equal(t, "owner-002", found.OwnerID)
}

func TestChannel_GetAssignmentByVaultAddress_NotFound(t *testing.T) {
	channel := newTestChannel()

	found, ok := channel.GetAssignmentByVaultAddress(
		"unknown-vault@ankhora",
	)

	require.False(t, ok)
	require.Nil(t, found)
}

func TestChannel_UpdateAssignment(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddAssignment(testAssignmentOne())

	updated := testAssignmentOne()
	updated.OwnerID = "owner-updated"
	updated.PublicKey = "GBUZX-UPDATED"
	updated.VaultAddress = "updated-vault@partner.ankhora"

	updatedResult := channel.UpdateAssignment(updated)

	require.True(t, updatedResult)

	found, ok := channel.GetAssignmentBySlotID("slot-001")

	require.True(t, ok)
	require.NotNil(t, found)
	require.Equal(t, "owner-updated", found.OwnerID)
	require.Equal(t, "GBUZX-UPDATED", found.PublicKey)
	require.Equal(t, "updated-vault@partner.ankhora", found.VaultAddress)
}

func TestChannel_UpdateAssignment_NotFound(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddAssignment(testAssignmentOne())

	updated := channel_domain.Assignment{
		SlotID:       "unknown-slot",
		OwnerID:      "owner-unknown",
		PublicKey:    "unknown-key",
		VaultAddress: "unknown-vault",
	}

	updatedResult := channel.UpdateAssignment(updated)

	require.False(t, updatedResult)
	require.Len(t, channel.Assignments, 1)
	require.Equal(t, testAssignmentOne(), channel.Assignments[0])
}

func TestChannel_RemoveAssignmentBySlotID(t *testing.T) {
	channel := newTestChannel()

	channel = channel.AddAssignment(testAssignmentOne())
	channel = channel.AddAssignment(testAssignmentTwo())

	removed := channel.RemoveAssignmentBySlotID("slot-001")

	require.True(t, removed)
	require.Len(t, channel.Assignments, 1)
	require.Equal(t, "slot-002", channel.Assignments[0].SlotID)

	_, found := channel.GetAssignmentBySlotID("slot-001")
	require.False(t, found)
}

func TestChannel_RemoveAssignmentBySlotID_NotFound(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddAssignment(testAssignmentOne())

	removed := channel.RemoveAssignmentBySlotID("unknown-slot")

	require.False(t, removed)
	require.Len(t, channel.Assignments, 1)
	require.Equal(t, testAssignmentOne(), channel.Assignments[0])
}


func TestChannel_RemoveAssignmentsByOwnerID_NotFound(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddAssignment(testAssignmentOne())

	removed := channel.RemoveAssignmentsByOwnerID("unknown-owner")

	require.Equal(t, 0, removed)
	require.Len(t, channel.Assignments, 1)
	require.Equal(t, testAssignmentOne(), channel.Assignments[0])
}

func TestChannel_ListAssignments(t *testing.T) {
	channel := newTestChannel()

	assignmentOne := testAssignmentOne()
	assignmentTwo := testAssignmentTwo()

	channel = channel.AddAssignment(assignmentOne)
	channel = channel.AddAssignment(assignmentTwo)

	assignments := channel.ListAssignments()

	require.Len(t, assignments, 2)
	require.Equal(t, assignmentOne, assignments[0])
	require.Equal(t, assignmentTwo, assignments[1])
}

func TestChannel_ListAssignments_ReturnsCopy(t *testing.T) {
	channel := newTestChannel()
	channel = channel.AddAssignment(testAssignmentOne())

	assignments := channel.ListAssignments()

	require.Len(t, assignments, 1)

	assignments[0].OwnerID = "changed-outside-channel"
	assignments = append(assignments, testAssignmentTwo())

	require.Len(t, channel.Assignments, 1)

	original, ok := channel.GetAssignmentBySlotID("slot-001")

	require.True(t, ok)
	require.Equal(t, "owner-001", original.OwnerID)
}

func TestChannel_AssignmentCRUD_EmptyChannel(t *testing.T) {
	channel := newTestChannel()

	require.Empty(t, channel.ListAssignments())
	require.Empty(t, channel.GetAssignmentsByOwnerID("owner-001"))

	assignmentBySlot, foundBySlot :=
		channel.GetAssignmentBySlotID("slot-001")

	require.False(t, foundBySlot)
	require.Nil(t, assignmentBySlot)

	assignmentByVault, foundByVault :=
		channel.GetAssignmentByVaultAddress("vault@ankhora")

	require.False(t, foundByVault)
	require.Nil(t, assignmentByVault)

	require.False(t, channel.UpdateAssignment(testAssignmentOne()))
	require.False(t, channel.RemoveAssignmentBySlotID("slot-001"))
	require.Equal(t, 0, channel.RemoveAssignmentsByOwnerID("owner-001"))
}

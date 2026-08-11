// internal/workspace/domain/events.go

package trustgroup_domain

import "time"

const (
	EventTrustGroupCreated = "trustgroup.created"
	EventTrustGroupRenamed = "trustgroup.renamed"
	EventTrustGroupDeleted = "trustgroup.deleted"
	EventMemberAddedToTrustGroup = "trustgroup.member.added"
	EventMemberRemovedFromTrustGroup = "trustgroup.member.removed"
	EventTrustGroupKEKRotated = "trustgroup.kek.rotated"
)

type TrustGroupCreated struct {
	EventID        string
	EventTimestamp time.Time

	ChannelID string
	VaultID     string
	OwnerID     string

	Name string
}

func (TrustGroupCreated) EventType() string {
	return EventTrustGroupCreated
}

type TrustGroupRenamed struct {
	EventID        string
	EventTimestamp time.Time

	WorkspaceID string

	OldName string
	NewName string
}

func (TrustGroupRenamed) EventType() string {
	return EventTrustGroupRenamed
}

type TrustGroupDeleted struct {
	EventID        string
	EventTimestamp time.Time

	TrustGroupID string
	WorkspaceID     string
}

func (TrustGroupDeleted) EventType() string {
	return EventTrustGroupDeleted
}

type MemberAddedToTrustGroup struct {
	EventID        string
	EventTimestamp time.Time

	TrustGroupID string
	ChannelID     string
	MemberID     string
}

func (MemberAddedToTrustGroup) EventType() string {
	return EventMemberAddedToTrustGroup
}

type MemberRemovedFromTrustGroup struct {
	EventID        string
	EventTimestamp time.Time

	TrustGroupID string
	WorkspaceID     string
	OwnerID     string
	Name string
}

func (MemberRemovedFromTrustGroup) EventType() string {
	return EventMemberRemovedFromTrustGroup
}

type TrustGroupKEKRotated struct {
	EventID        string
	EventTimestamp time.Time

	TrustGroupID string
	WorkspaceID     string
	OwnerID     string
	Name string
}

func (TrustGroupKEKRotated) EventType() string {
	return EventTrustGroupKEKRotated
}		
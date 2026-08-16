package tracecore_types

import "time"

// The Cloud backend marshals its Channel aggregates with default Go JSON
// encoding (capitalized field names). These DTOs mirror that wire format so
// the Tracecore client can decode the Cloud response before mapping it into
// channel_domain.Channel. The mapping lives in internal/tracecore.

type CloudChannelFederation struct {
	VaultAID          string   `json:"VaultAID"`
	VaultBID          string   `json:"VaultBID"`
	AllowedEventTypes []string `json:"AllowedEventTypes"`
	AllowedPaths      []string `json:"AllowedPaths"`
	AllowedDirections string   `json:"AllowedDirections"`
}

type CloudChannelSlot struct {
	ID      string `json:"ID"`
	Name    string `json:"Name"`
	Role    string `json:"Role"`
	VaultID string `json:"VaultID"`
	Gated   bool   `json:"Gated"`
	Order   int    `json:"Order"`
}

type CloudChannelAssignment struct {
	SlotID       string `json:"SlotID"`
	OwnerID      string `json:"OwnerID"`
	PublicKey    string `json:"PublicKey"`
	VaultAddress string `json:"VaultAddress"`
}

type CloudChannelProperty struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type CloudChannelDTO struct {
	ID          string                   `json:"ID"`
	TemplateID  string                   `json:"TemplateID"`
	Title       string                   `json:"Title"`
	Status      string                   `json:"Status"`
	Slots       []CloudChannelSlot       `json:"Slots"`
	Assignments []CloudChannelAssignment `json:"Assignments"`
	Properties  []CloudChannelProperty   `json:"Properties"`
	Policy      map[string]any           `json:"Policy"`
	Federation  *CloudChannelFederation  `json:"Federation"`
	CreatedAt   time.Time                `json:"CreatedAt"`
	UpdatedAt   time.Time                `json:"UpdatedAt"`
	RevokedAt   *time.Time               `json:"RevokedAt"`
	ArchivedAt  *time.Time               `json:"ArchivedAt"`
	WorkspaceID string                   `json:"WorkspaceID"`
	IsDraft     bool                     `json:"IsDraft"`
	IsDirty     bool                     `json:"IsDirty"`
}

// CloudChannelParticipant mirrors the Cloud Participant aggregate, which is
// marshalled with default Go JSON encoding (capitalized field names). JoinedAt
// is a Unix timestamp in seconds.
type CloudChannelParticipant struct {
	ChannelID   string   `json:"ChannelID"`
	VaultID     string   `json:"VaultID"`
	PublicKey   string   `json:"PublicKey"`
	Direction   string   `json:"Direction"`
	JoinedAt    int64    `json:"JoinedAt"`
	Role        string   `json:"Role"`
	Permissions []string `json:"Permissions"`
}

// CloudChannelInvitation mirrors the Cloud Invitation aggregate, which is
// marshalled with default Go JSON encoding (capitalized field names). AcceptedAt
// is null while the invitation is pending. Invitations carry no slot or role
// information.
type CloudChannelInvitation struct {
	ID             string     `json:"ID"`
	ChannelID      string     `json:"ChannelID"`
	InviterVaultID string     `json:"InviterVaultID"`
	InviteeVaultID string     `json:"InviteeVaultID"`
	Status         string     `json:"Status"`
	CreatedAt      time.Time  `json:"CreatedAt"`
	AcceptedAt     *time.Time `json:"AcceptedAt"`
}

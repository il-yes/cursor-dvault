package tracecore_types

import (
	"encoding/json"
	"time"
)



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

func (s *CloudChannelSlot) UnmarshalJSON(data []byte) error {
	type Alias CloudChannelSlot
	aux := &struct {
		AltID      string `json:"id"`
		AltName    string `json:"name"`
		AltRole    string `json:"role"`
		AltVaultID string `json:"vault_id"`
		AltGated   *bool  `json:"gated"`
		AltOrder   *int   `json:"order"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if s.ID == "" && aux.AltID != "" {
		s.ID = aux.AltID
	}
	if s.Name == "" && aux.AltName != "" {
		s.Name = aux.AltName
	}
	if s.Role == "" && aux.AltRole != "" {
		s.Role = aux.AltRole
	}
	if s.VaultID == "" && aux.AltVaultID != "" {
		s.VaultID = aux.AltVaultID
	}
	if !s.Gated && aux.AltGated != nil {
		s.Gated = *aux.AltGated
	}
	if s.Order == 0 && aux.AltOrder != nil {
		s.Order = *aux.AltOrder
	}
	return nil
}

// UpdateChannelSlotDTO represents outbound channel slot serialization for HTTP PUT /channels/{id}.
// The Cloud backend UpdateChannel contract requires snake_case JSON property keys.
type UpdateChannelSlotDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Role    string `json:"role"`
	VaultID string `json:"vault_id"`
	Gated   bool   `json:"gated"`
	Order   int    `json:"order"`
}

type CloudChannelAssignment struct {
	SlotID       string `json:"SlotID"`
	OwnerID      string `json:"OwnerID"`
	PublicKey    string `json:"PublicKey"`
	VaultAddress string `json:"VaultAddress"`
}

func (a *CloudChannelAssignment) UnmarshalJSON(data []byte) error {
	type Alias CloudChannelAssignment
	aux := &struct {
		AltSlotID       string `json:"slot_id"`
		AltOwnerID      string `json:"owner_id"`
		AltPublicKey    string `json:"public_key"`
		AltVaultAddress string `json:"vault_address"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if a.SlotID == "" && aux.AltSlotID != "" {
		a.SlotID = aux.AltSlotID
	}
	if a.OwnerID == "" && aux.AltOwnerID != "" {
		a.OwnerID = aux.AltOwnerID
	}
	if a.PublicKey == "" && aux.AltPublicKey != "" {
		a.PublicKey = aux.AltPublicKey
	}
	if a.VaultAddress == "" && aux.AltVaultAddress != "" {
		a.VaultAddress = aux.AltVaultAddress
	}
	return nil
}

type CloudChannelProperty struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func (p *CloudChannelProperty) UnmarshalJSON(data []byte) error {
	type Alias CloudChannelProperty
	aux := &struct {
		AltKey   string `json:"key"`
		AltValue string `json:"value"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if p.Key == "" && aux.AltKey != "" {
		p.Key = aux.AltKey
	}
	if p.Value == "" && aux.AltValue != "" {
		p.Value = aux.AltValue
	}
	return nil
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

func (c *CloudChannelDTO) UnmarshalJSON(data []byte) error {
	type Alias CloudChannelDTO
	aux := &struct {
		AltID          string                   `json:"id"`
		AltTemplateID  string                   `json:"template_id"`
		AltTitle       string                   `json:"title"`
		AltStatus      string                   `json:"status"`
		AltSlots       []CloudChannelSlot       `json:"slots"`
		AltAssignments []CloudChannelAssignment `json:"assignments"`
		AltProperties  []CloudChannelProperty   `json:"properties"`
		AltPolicy      map[string]any           `json:"policy"`
		AltCreatedAt   *time.Time               `json:"created_at"`
		AltUpdatedAt   *time.Time               `json:"updated_at"`
		AltRevokedAt   *time.Time               `json:"revoked_at"`
		AltArchivedAt  *time.Time               `json:"archived_at"`
		AltWorkspaceID string                   `json:"workspace_id"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if c.ID == "" && aux.AltID != "" {
		c.ID = aux.AltID
	}
	if c.TemplateID == "" && aux.AltTemplateID != "" {
		c.TemplateID = aux.AltTemplateID
	}
	if c.Title == "" && aux.AltTitle != "" {
		c.Title = aux.AltTitle
	}
	if c.Status == "" && aux.AltStatus != "" {
		c.Status = aux.AltStatus
	}
	if len(c.Slots) == 0 && len(aux.AltSlots) > 0 {
		c.Slots = aux.AltSlots
	}
	if len(c.Assignments) == 0 && len(aux.AltAssignments) > 0 {
		c.Assignments = aux.AltAssignments
	}
	if len(c.Properties) == 0 && len(aux.AltProperties) > 0 {
		c.Properties = aux.AltProperties
	}
	if c.Policy == nil && aux.AltPolicy != nil {
		c.Policy = aux.AltPolicy
	}
	if c.CreatedAt.IsZero() && aux.AltCreatedAt != nil {
		c.CreatedAt = *aux.AltCreatedAt
	}
	if c.UpdatedAt.IsZero() && aux.AltUpdatedAt != nil {
		c.UpdatedAt = *aux.AltUpdatedAt
	}
	if c.RevokedAt == nil && aux.AltRevokedAt != nil {
		c.RevokedAt = aux.AltRevokedAt
	}
	if c.ArchivedAt == nil && aux.AltArchivedAt != nil {
		c.ArchivedAt = aux.AltArchivedAt
	}
	if c.WorkspaceID == "" && aux.AltWorkspaceID != "" {
		c.WorkspaceID = aux.AltWorkspaceID
	}
	return nil
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

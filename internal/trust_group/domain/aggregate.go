package trustgroup_domain

import (
	"time"

	"github.com/google/uuid"
)

type  TrustGroup struct {
	ID          string `json:"id"`
	ChannelID string `json:"channel_id"`
	Name        string `json:"name"`
	KEKVersion  uint64 `json:"kek_version"`
	MemberCIDs  []string `json:"member_cids"`
	KeyEnvelopes []TrustGroupKeyEnvelope   `json:"key_envelopes"`
	CreatedAt   string `json:"created_at"`
	IsDraft     bool   `json:"is_draft"`
	IsDirty     bool   `json:"is_dirty" gorm:"boolean"`
}


type TrustGroupMember struct {
	ID         string          `json:"id"`
	VaultID    string          `json:"vault_id"`
	Role       string          `json:"role"`
	WrappedKEK *WrappedGroupKey `json:"wrapped_kek,omitempty"`
	JoinedAt   time.Time       `json:"joined_at"`
	IsDraft    bool            `json:"is_draft"`
	IsDirty    bool            `json:"is_dirty" gorm:"boolean"`
}

type WrappedGroupKey struct {
	Version uint64 `json:"version"`
	Value   string `json:"value"`
}


func NewWrappedGroupKey(version uint64, value string) *WrappedGroupKey {
	if value == "" {
		return nil
	}
	return &WrappedGroupKey{
		Version: version,
		Value:   value,
	}
}

func NewTrustGroup(
	channelID string,
	name string,
	memberCIDs []string,
) *TrustGroup {
	return &TrustGroup{
		ID:          uuid.NewString(),
		ChannelID: channelID,
		Name:        name,
		KEKVersion:  1,
		MemberCIDs:  memberCIDs,
		CreatedAt:   time.Now().Format(time.RFC3339Nano),
		IsDraft:     true,
		IsDirty:     false,
	}
}

func (g *TrustGroup) AddMember(member TrustGroupMember) error {
	if member.VaultID == "" {
		return ErrVaultIDRequired
	}
	for _, cid := range g.MemberCIDs {
		if cid == member.VaultID {
			return nil
		}
	}
	g.MemberCIDs = append(g.MemberCIDs, member.VaultID)
	g.IsDirty = true
	return nil
}

func (g *TrustGroup) RemoveMember(vaultID string) error {
	if vaultID == "" {
		return ErrVaultIDRequired
	}
	found := false
	updatedCIDs := make([]string, 0, len(g.MemberCIDs))
	for _, cid := range g.MemberCIDs {
		if cid == vaultID {
			found = true
			continue
		}
		updatedCIDs = append(updatedCIDs, cid)
	}
	if !found {
		return ErrTrustGroupMemberNotFound
	}
	g.MemberCIDs = updatedCIDs
	g.IsDirty = true
	return nil
}

func (g *TrustGroup) RotateKEK(newVersion uint64, newEnvelopes []TrustGroupKeyEnvelope) error {
	if newVersion != g.KEKVersion+1 {
		return ErrInvalidKEKVersionIncrement
	}
	now := time.Now()
	// Mark all existing active envelopes as revoked
	for i := range g.KeyEnvelopes {
		if g.KeyEnvelopes[i].RevokedAt == nil {
			g.KeyEnvelopes[i].RevokedAt = &now
		}
	}
	g.KEKVersion = newVersion
	for _, env := range newEnvelopes {
		if env.KEKVersion != g.KEKVersion {
			return ErrStaleKEKVersion
		}
		if env.ID == "" {
			env.ID = uuid.NewString()
		}
		env.TrustGroupID = g.ID
		if env.CreatedAt.IsZero() {
			env.CreatedAt = now
		}
		g.KeyEnvelopes = append(g.KeyEnvelopes, env)
	}
	g.IsDirty = true
	return nil
}

func (g *TrustGroup) AddEnvelope(env TrustGroupKeyEnvelope) error {
	if env.MemberID == "" {
		return ErrMemberIDRequired
	}
	if env.DeviceID == "" {
		return ErrDeviceIDRequired
	}
	if env.WrappedKEK == "" {
		return ErrWrappedKEKRequired
	}
	if env.KEKVersion != g.KEKVersion {
		return ErrStaleKEKVersion
	}
	for _, existing := range g.KeyEnvelopes {
		if existing.DeviceID == env.DeviceID && existing.KEKVersion == env.KEKVersion && existing.RevokedAt == nil {
			return ErrDuplicateKeyEnvelope
		}
	}
	if env.ID == "" {
		env.ID = uuid.NewString()
	}
	if env.TrustGroupID == "" {
		env.TrustGroupID = g.ID
	}
	if env.CreatedAt.IsZero() {
		env.CreatedAt = time.Now()
	}
	g.KeyEnvelopes = append(g.KeyEnvelopes, env)
	g.IsDirty = true
	return nil
}


type TrustGroupKeyEnvelope struct {
	ID           string     `json:"id"`
	TrustGroupID string     `json:"trust_group_id"`
	MemberID     string     `json:"member_id"`
	DeviceID     string     `json:"device_id"`
	KEKVersion   uint64     `json:"kek_version"`
	WrappedKEK   string     `json:"wrapped_kek"`
	CreatedAt    time.Time  `json:"created_at"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
}
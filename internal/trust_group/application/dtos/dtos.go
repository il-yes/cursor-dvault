package trustgroup_dtos



type CreateTrustGroupRequest struct {
	ChannelID string
	OwnerID     string
	VaultID     string
	Name        string
	MemberCIDs  []string
}

type AddMemberToTrustGroupRequest struct {
	TrustGroupID string `json:"-"`
	VaultID      string `json:"vault_id"`
	Role         string `json:"role"`
}

type RemoveMemberFromTrustGroupRequest struct {
	TrustGroupID string
	MemberID     string
	ChannelID    string
}

type AddTrustGroupKeyEnvelopeRequest struct {
	TrustGroupID string `json:"trust_group_id"`
	MemberID     string `json:"member_id"`
	DeviceID     string `json:"device_id,omitempty"`
	KEKVersion   uint64 `json:"kek_version"`
	WrappedKEK   string `json:"wrapped_kek"`
}

type ProvisionTrustGroupMemberEnvelopeRequest struct {
	TrustGroupID    string `json:"trust_group_id"`
	MemberID        string `json:"member_id"`
	DeviceID        string `json:"device_id,omitempty"`
	DevicePublicKey string `json:"device_public_key,omitempty"`
	MemberPublicKey string `json:"member_public_key,omitempty"`
}

type ProvisionTrustGroupDeviceEnvelopeRequest = ProvisionTrustGroupMemberEnvelopeRequest
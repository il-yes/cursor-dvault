package trustgroup_dtos



type CreateTrustGroupRequest struct {
	ChannelID string
	OwnerID     string
	VaultID     string
	Name        string
	MemberCIDs  []string
}

type AddMemberToTrustGroupRequest struct {
	TrustGroupID string
	MemberID     string
	ChannelID    string
}

type RemoveMemberFromTrustGroupRequest struct {
	TrustGroupID string
	MemberID     string
	ChannelID    string
}

type AddTrustGroupKeyEnvelopeRequest struct {
	TrustGroupID string `json:"trust_group_id"`
	MemberID     string `json:"member_id"`
	DeviceID     string `json:"device_id"`
	KEKVersion   uint64 `json:"kek_version"`
	WrappedKEK   string `json:"wrapped_kek"`
}
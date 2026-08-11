package trustgroup_domain

import (
	"context"

	tracecore_types "vault-app/internal/tracecore/types"
)


type CreateTrustGroupRequest struct {
	TrustGroup TrustGroup
}
type GetTrustGroupMemberRequest struct {
	TrustGroupID string
	MemberID     string
}
type AddMemberToTrustGroupRequest struct {
	TrustGroupID string
	MemberID     string
}
type RemoveMemberFromTrustGroupRequest struct {
	TrustGroupID string
	MemberID     string
}
type RotateTrustGroupKEKRequest struct {
	TrustGroupID string
}
type GetTrustGroupRequest struct {
	TrustGroupID string
}
type ListTrustGroupsRequest struct {
	ChannelID string
}
type UpdateTrustGroupRequest struct {
	TrustGroup TrustGroup
}
type DeleteTrustGroupRequest struct {
	TrustGroupID string
}
type TrustGroupRepository interface {
	CreateTrustGroup(ctx context.Context, req *CreateTrustGroupRequest) (*tracecore_types.CloudResponse[TrustGroup], error)
	GetTrustGroupMember(ctx context.Context, req *GetTrustGroupMemberRequest) (*tracecore_types.CloudResponse[TrustGroupMember], error)
	ListTrustGroups(ctx context.Context, req *ListTrustGroupsRequest) (*tracecore_types.CloudResponse[[]TrustGroup], error)
	GetTrustGroup(ctx context.Context, req *GetTrustGroupRequest) (*tracecore_types.CloudResponse[TrustGroup], error)
	UpdateTrustGroup(ctx context.Context, req *UpdateTrustGroupRequest) (*tracecore_types.CloudResponse[TrustGroup], error)
	DeleteTrustGroup(ctx context.Context, req *DeleteTrustGroupRequest) (*tracecore_types.CloudResponse[TrustGroup], error)
	AddMemberToTrustGroup(ctx context.Context, req *AddMemberToTrustGroupRequest) (*tracecore_types.CloudResponse[TrustGroup], error)
	RemoveMemberFromTrustGroup(ctx context.Context, req *RemoveMemberFromTrustGroupRequest) (*tracecore_types.CloudResponse[TrustGroup], error)
	RotateTrustGroupKEK(ctx context.Context, req *RotateTrustGroupKEKRequest) (*tracecore_types.CloudResponse[TrustGroup], error)
}

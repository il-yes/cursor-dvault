package trustgroup_events


import (
	"context"

	trustgroup_domain "vault-app/internal/trust_group/domain"
)

type TrustGroupEventBus interface {
	PublishTrustGroupCreated(
		ctx context.Context,
		event trustgroup_domain.TrustGroupCreated,
	) error

	SubscribeToTrustGroupCreated(
		handler func(context.Context, trustgroup_domain.TrustGroupCreated) error,
	)

	PublishTrustGroupRenamed(
		ctx context.Context,
		event trustgroup_domain.TrustGroupRenamed,
	) error

	SubscribeToTrustGroupRenamed(
		handler func(context.Context, trustgroup_domain.TrustGroupRenamed) error,
	)

	PublishTrustGroupDeleted(
		ctx context.Context,
		event trustgroup_domain.TrustGroupDeleted,
	) error

	SubscribeToTrustGroupDeleted(
		handler func(context.Context, trustgroup_domain.TrustGroupDeleted) error,
	)

	PublishMemberAddedToTrustGroup(
		ctx context.Context,
		event trustgroup_domain.MemberAddedToTrustGroup,
	) error

	SubscribeToMemberAddedToTrustGroup(
		handler func(context.Context, trustgroup_domain.MemberAddedToTrustGroup) error,
	)

	PublishMemberRemovedFromTrustGroup(
		ctx context.Context,
		event trustgroup_domain.MemberRemovedFromTrustGroup,
	) error

	SubscribeToMemberRemovedFromTrustGroup(
		handler func(context.Context, trustgroup_domain.MemberRemovedFromTrustGroup) error,
	)

	PublishTrustGroupKEKRotated(
		ctx context.Context,
		event trustgroup_domain.TrustGroupKEKRotated,
	) error

	SubscribeToTrustGroupKEKRotated(
		handler func(context.Context, trustgroup_domain.TrustGroupKEKRotated) error,
	)
}
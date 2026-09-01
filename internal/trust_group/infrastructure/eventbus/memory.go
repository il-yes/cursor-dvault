package trustgroup_infrastructure_eventbus

import (
	"context"
	"sync"

	trustGroup_domain "vault-app/internal/trust_group/domain"
)

type MemoryBus struct {
	trustGroupCreatedSubscribers           []func(ctx context.Context, event trustGroup_domain.TrustGroupCreated) error
	trustGroupRenamedSubscribers           []func(ctx context.Context, event trustGroup_domain.TrustGroupRenamed) error
	trustGroupDeletedSubscribers           []func(ctx context.Context, event trustGroup_domain.TrustGroupDeleted) error
	memberAddedToTrustGroupSubscribers     []func(ctx context.Context, event trustGroup_domain.MemberAddedToTrustGroup) error
	memberRemovedFromTrustGroupSubscribers []func(ctx context.Context, event trustGroup_domain.MemberRemovedFromTrustGroup) error
	trustGroupKEKRotatedSubscribers        []func(ctx context.Context, event trustGroup_domain.TrustGroupKEKRotated) error
	lock                                   sync.RWMutex
}

func NewMemoryBus() *MemoryBus {
	return &MemoryBus{
		trustGroupCreatedSubscribers:           make([]func(ctx context.Context, event trustGroup_domain.TrustGroupCreated) error, 0),
		trustGroupRenamedSubscribers:           make([]func(ctx context.Context, event trustGroup_domain.TrustGroupRenamed) error, 0),
		trustGroupDeletedSubscribers:           make([]func(ctx context.Context, event trustGroup_domain.TrustGroupDeleted) error, 0),
		memberAddedToTrustGroupSubscribers:     make([]func(ctx context.Context, event trustGroup_domain.MemberAddedToTrustGroup) error, 0),
		memberRemovedFromTrustGroupSubscribers: make([]func(ctx context.Context, event trustGroup_domain.MemberRemovedFromTrustGroup) error, 0),
		trustGroupKEKRotatedSubscribers:        make([]func(ctx context.Context, event trustGroup_domain.TrustGroupKEKRotated) error, 0),
	}
}

// TrustGroupCreated Event
func (mb *MemoryBus) PublishTrustGroupCreated(ctx context.Context, event trustGroup_domain.TrustGroupCreated) error {
	mb.lock.RLock()
	defer mb.lock.RUnlock()
	for _, h := range mb.trustGroupCreatedSubscribers {
		go h(ctx, event)
	}
	return nil
}
func (mb *MemoryBus) SubscribeToTrustGroupCreated(handler func(ctx context.Context, event trustGroup_domain.TrustGroupCreated) error) {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.trustGroupCreatedSubscribers = append(mb.trustGroupCreatedSubscribers, handler)
}

func (mb *MemoryBus) PublishTrustGroupRenamed(ctx context.Context, event trustGroup_domain.TrustGroupRenamed) error {
	mb.lock.RLock()
	defer mb.lock.RUnlock()
	for _, h := range mb.trustGroupRenamedSubscribers {
		go h(ctx, event)
	}
	return nil
}
func (mb *MemoryBus) SubscribeToTrustGroupRenamed(handler func(ctx context.Context, event trustGroup_domain.TrustGroupRenamed) error) {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.trustGroupRenamedSubscribers = append(mb.trustGroupRenamedSubscribers, handler)
}

func (mb *MemoryBus) PublishTrustGroupDeleted(ctx context.Context, event trustGroup_domain.TrustGroupDeleted) error {
	mb.lock.RLock()
	defer mb.lock.RUnlock()
	for _, h := range mb.trustGroupDeletedSubscribers {
		go h(ctx, event)
	}
	return nil
}
func (mb *MemoryBus) SubscribeToTrustGroupDeleted(handler func(ctx context.Context, event trustGroup_domain.TrustGroupDeleted) error) {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.trustGroupDeletedSubscribers = append(mb.trustGroupDeletedSubscribers, handler)
}

func (mb *MemoryBus) PublishMemberAddedToTrustGroup(ctx context.Context, event trustGroup_domain.MemberAddedToTrustGroup) error {
	mb.lock.RLock()
	defer mb.lock.RUnlock()
	for _, h := range mb.memberAddedToTrustGroupSubscribers {
		go h(ctx, event)
	}
	return nil
}
func (mb *MemoryBus) SubscribeToMemberAddedToTrustGroup(handler func(ctx context.Context, event trustGroup_domain.MemberAddedToTrustGroup) error) {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.memberAddedToTrustGroupSubscribers = append(mb.memberAddedToTrustGroupSubscribers, handler)
}

func (mb *MemoryBus) PublishMemberRemovedFromTrustGroup(ctx context.Context, event trustGroup_domain.MemberRemovedFromTrustGroup) error {
	mb.lock.RLock()
	defer mb.lock.RUnlock()
	for _, h := range mb.memberRemovedFromTrustGroupSubscribers {
		go h(ctx, event)
	}
	return nil
}
func (mb *MemoryBus) SubscribeToMemberRemovedFromTrustGroup(handler func(ctx context.Context, event trustGroup_domain.MemberRemovedFromTrustGroup) error) {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.memberRemovedFromTrustGroupSubscribers = append(mb.memberRemovedFromTrustGroupSubscribers, handler)
}

func (mb *MemoryBus) PublishTrustGroupKEKRotated(ctx context.Context, event trustGroup_domain.TrustGroupKEKRotated) error {
	mb.lock.RLock()
	defer mb.lock.RUnlock()
	for _, h := range mb.trustGroupKEKRotatedSubscribers {
		go h(ctx, event)
	}
	return nil
}
func (mb *MemoryBus) SubscribeToTrustGroupKEKRotated(handler func(ctx context.Context, event trustGroup_domain.TrustGroupKEKRotated) error) {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.trustGroupKEKRotatedSubscribers = append(mb.trustGroupKEKRotatedSubscribers, handler)
}	
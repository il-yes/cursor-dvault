package workspace_infrastructure_eventbus

import (
	"context"
	"sync"

	workspace_domain "vault-app/internal/workspace/domain"
)

type MemoryBus struct {
	workspaceCreatedSubscribers []func(ctx context.Context, event workspace_domain.WorkspaceCreated)
	workspaceRenamedSubscribers []func(ctx context.Context, event workspace_domain.WorkspaceRenamed)
	workspaceDeletedSubscribers []func(ctx context.Context, event workspace_domain.WorkspaceDeleted)
	lock        sync.RWMutex
}

func NewMemoryBus() *MemoryBus {
	return &MemoryBus{
		workspaceCreatedSubscribers: make([]func(ctx context.Context, event workspace_domain.WorkspaceCreated), 0),
		workspaceRenamedSubscribers: make([]func(ctx context.Context, event workspace_domain.WorkspaceRenamed), 0),
		workspaceDeletedSubscribers: make([]func(ctx context.Context, event workspace_domain.WorkspaceDeleted), 0),
	}
}

// WorkspaceCreated Event
func (mb *MemoryBus) PublishWorkspaceCreated(ctx context.Context, event workspace_domain.WorkspaceCreated) error {
	mb.lock.RLock()
	defer mb.lock.RUnlock()
	for _, h := range mb.workspaceCreatedSubscribers {
		go h(ctx, event)
	}
	return nil
}
func (mb *MemoryBus) SubscribeToWorkspaceCreated(handler func(ctx context.Context, event workspace_domain.WorkspaceCreated)) error {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.workspaceCreatedSubscribers = append(mb.workspaceCreatedSubscribers, handler)
	return nil
}

func (mb *MemoryBus) PublishWorkspaceRenamed(ctx context.Context, event workspace_domain.WorkspaceRenamed) error {
	mb.lock.RLock()
	defer mb.lock.RUnlock()
	for _, h := range mb.workspaceRenamedSubscribers {
		go h(ctx, event)
	}
	return nil
}
func (mb *MemoryBus) SubscribeToWorkspaceRenamed(handler func(ctx context.Context, event workspace_domain.WorkspaceRenamed)) error {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.workspaceRenamedSubscribers = append(mb.workspaceRenamedSubscribers, handler)
	return nil
}

func (mb *MemoryBus) PublishWorkspaceDeleted(ctx context.Context, event workspace_domain.WorkspaceDeleted) error {
	mb.lock.RLock()
	defer mb.lock.RUnlock()
	for _, h := range mb.workspaceDeletedSubscribers {
		go h(ctx, event)
	}
	return nil
}
func (mb *MemoryBus) SubscribeToWorkspaceDeleted(handler func(ctx context.Context, event workspace_domain.WorkspaceDeleted)) error {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.workspaceDeletedSubscribers = append(mb.workspaceDeletedSubscribers, handler)
	return nil
}

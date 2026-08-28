package collaboration_infra

import (
	"context"
	"sync"

	collaboration_domain "vault-app/internal/collaboration/domain"
)

type MemoryActionRepository struct {
	mu        sync.RWMutex
	approvals map[string]collaboration_domain.Approval
	rejects   map[string]collaboration_domain.Reject
	transfers map[string]collaboration_domain.Transfer
}

func NewMemoryActionRepository() *MemoryActionRepository {
	return &MemoryActionRepository{
		approvals: make(map[string]collaboration_domain.Approval),
		rejects:   make(map[string]collaboration_domain.Reject),
		transfers: make(map[string]collaboration_domain.Transfer),
	}
}

// ApprovalRepository
func (r *MemoryActionRepository) SaveApproval(_ context.Context, approval collaboration_domain.Approval) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.approvals[approval.ID] = approval
	return nil
}

func (r *MemoryActionRepository) GetApprovalByID(_ context.Context, id string) (*collaboration_domain.Approval, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	appr, ok := r.approvals[id]
	if !ok {
		return nil, collaboration_domain.ErrApprovalNotFound
	}
	return &appr, nil
}

func (r *MemoryActionRepository) ListApprovalsByResource(_ context.Context, resourceType string, resourceID string) ([]collaboration_domain.Approval, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []collaboration_domain.Approval
	for _, appr := range r.approvals {
		if appr.ResourceReference.ResourceType == resourceType && appr.ResourceReference.ResourceID == resourceID {
			result = append(result, appr)
		}
	}
	return result, nil
}

// RejectRepository
func (r *MemoryActionRepository) SaveReject(_ context.Context, reject collaboration_domain.Reject) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rejects[reject.ID] = reject
	return nil
}

func (r *MemoryActionRepository) GetRejectByID(_ context.Context, id string) (*collaboration_domain.Reject, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rej, ok := r.rejects[id]
	if !ok {
		return nil, collaboration_domain.ErrRejectNotFound
	}
	return &rej, nil
}

func (r *MemoryActionRepository) ListRejectsByResource(_ context.Context, resourceType string, resourceID string) ([]collaboration_domain.Reject, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []collaboration_domain.Reject
	for _, rej := range r.rejects {
		if rej.ResourceReference.ResourceType == resourceType && rej.ResourceReference.ResourceID == resourceID {
			result = append(result, rej)
		}
	}
	return result, nil
}

// TransferRepository
func (r *MemoryActionRepository) SaveTransfer(_ context.Context, transfer collaboration_domain.Transfer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.transfers[transfer.ID] = transfer
	return nil
}

func (r *MemoryActionRepository) GetTransferByID(_ context.Context, id string) (*collaboration_domain.Transfer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tr, ok := r.transfers[id]
	if !ok {
		return nil, collaboration_domain.ErrTransferNotFound
	}
	return &tr, nil
}

func (r *MemoryActionRepository) ListTransfersByResource(_ context.Context, resourceType string, resourceID string) ([]collaboration_domain.Transfer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []collaboration_domain.Transfer
	for _, tr := range r.transfers {
		if tr.ResourceReference.ResourceType == resourceType && tr.ResourceReference.ResourceID == resourceID {
			result = append(result, tr)
		}
	}
	return result, nil
}

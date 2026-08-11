package collaboration_usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	c3_asset_domain "vault-app/internal/c3_asset/domain"
	collaboration_ports "vault-app/internal/collaboration/application/ports"
	trustgroup_dtos "vault-app/internal/trust_group/application/dtos"
	trustgroup_domain "vault-app/internal/trust_group/domain"
)

type RotatedShareEntryInput struct {
	ShareEntryID string
	ReWrappedDEK string
}

type RotateTrustGroupKEKRequest struct {
	RequestID           string // Idempotency key
	TrustGroupID        string
	OldVersion          uint64
	NewVersion          uint64
	RevokedMemberID     string // Optional: member to remove during rotation
	NewEnvelopes        []trustgroup_dtos.AddTrustGroupKeyEnvelopeRequest
	RotatedShareEntries []RotatedShareEntryInput
}

type RotateTrustGroupKEKResponse struct {
	TrustGroup trustgroup_domain.TrustGroup
	Count      int
}

type RotateTrustGroupKEKUseCase struct {
	trustGroupRepo    trustgroup_domain.TrustGroupRepository
	shareEntryRepo    c3_asset_domain.ShareEntryRepository
	txManager         collaboration_ports.TransactionManager
	mu                sync.Mutex
	processedRequests map[string]*RotateTrustGroupKEKResponse
}

func NewRotateTrustGroupKEKUseCase(
	trustGroupRepo trustgroup_domain.TrustGroupRepository,
	shareEntryRepo c3_asset_domain.ShareEntryRepository,
) *RotateTrustGroupKEKUseCase {
	return NewRotateTrustGroupKEKUseCaseWithTx(trustGroupRepo, shareEntryRepo, &collaboration_ports.NopTransactionManager{})
}

func NewRotateTrustGroupKEKUseCaseWithTx(
	trustGroupRepo trustgroup_domain.TrustGroupRepository,
	shareEntryRepo c3_asset_domain.ShareEntryRepository,
	txManager collaboration_ports.TransactionManager,
) *RotateTrustGroupKEKUseCase {
	if txManager == nil {
		txManager = &collaboration_ports.NopTransactionManager{}
	}
	return &RotateTrustGroupKEKUseCase{
		trustGroupRepo:    trustGroupRepo,
		shareEntryRepo:    shareEntryRepo,
		txManager:         txManager,
		processedRequests: make(map[string]*RotateTrustGroupKEKResponse),
	}
}

func (u *RotateTrustGroupKEKUseCase) ValidateDependencies() error {
	if u.trustGroupRepo == nil {
		return trustgroup_domain.ErrRepositoryNil
	}
	if u.shareEntryRepo == nil {
		return c3_asset_domain.ErrRepositoryNil
	}
	return nil
}

func (u *RotateTrustGroupKEKUseCase) Execute(
	ctx context.Context,
	req RotateTrustGroupKEKRequest,
) (*RotateTrustGroupKEKResponse, error) {
	if err := u.ValidateDependencies(); err != nil {
		return nil, err
	}

	if strings.TrimSpace(req.TrustGroupID) == "" {
		return nil, errors.New("trust group ID is required")
	}
	if req.OldVersion == 0 {
		return nil, errors.New("old version is required")
	}
	if req.NewVersion != req.OldVersion+1 {
		return nil, trustgroup_domain.ErrInvalidKEKVersionIncrement
	}

	// -------------------------------------------------------------------------
	// 1. Idempotency Cache Check (Duplicate / Retry Request Detection)
	// -------------------------------------------------------------------------
	if req.RequestID != "" {
		u.mu.Lock()
		if resp, exists := u.processedRequests[req.RequestID]; exists {
			u.mu.Unlock()
			return resp, nil
		}
		u.mu.Unlock()
	}

	// -------------------------------------------------------------------------
	// 2. Fetch TrustGroup & Validate Current Version
	// -------------------------------------------------------------------------
	tgResp, err := u.trustGroupRepo.GetTrustGroup(ctx, &trustgroup_domain.GetTrustGroupRequest{
		TrustGroupID: req.TrustGroupID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trust group: %w", err)
	}
	if tgResp == nil || tgResp.Data.ID == "" {
		return nil, trustgroup_domain.ErrTrustGroupNotFound
	}
	tg := tgResp.Data

	// Stale Client / Concurrency Protection
	if tg.KEKVersion != req.OldVersion {
		return nil, trustgroup_domain.ErrStaleKEKVersion
	}


	// -------------------------------------------------------------------------
	// 3. Prepare Domain Envelopes & Aggregate State Transition
	// -------------------------------------------------------------------------
	if req.RevokedMemberID != "" {
		_ = tg.RemoveMember(req.RevokedMemberID)
	}

	domainEnvelopes := make([]trustgroup_domain.TrustGroupKeyEnvelope, 0, len(req.NewEnvelopes))
	for _, envReq := range req.NewEnvelopes {
		domainEnvelopes = append(domainEnvelopes, trustgroup_domain.TrustGroupKeyEnvelope{
			TrustGroupID: req.TrustGroupID,
			MemberID:     envReq.MemberID,
			DeviceID:     envReq.DeviceID,
			KEKVersion:   req.NewVersion,
			WrappedKEK:   envReq.WrappedKEK,
		})
	}

	if err := tg.RotateKEK(req.NewVersion, domainEnvelopes); err != nil {
		return nil, fmt.Errorf("failed to rotate KEK in trust group aggregate: %w", err)
	}

	// -------------------------------------------------------------------------
	// 4. Atomic Unit-of-Work Persistence
	// -------------------------------------------------------------------------
	var finalTrustGroup trustgroup_domain.TrustGroup
	var updatedCount int

	txErr := u.txManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		// Save updated TrustGroup
		updatedTgResp, err := u.trustGroupRepo.UpdateTrustGroup(txCtx, &trustgroup_domain.UpdateTrustGroupRequest{
			TrustGroup: tg,
		})
		if err != nil {
			return fmt.Errorf("failed to persist rotated trust group: %w", err)
		}
		finalTrustGroup = updatedTgResp.Data

		// Update ShareEntries with re-wrapped DEKs and new KEKVersion
		for _, item := range req.RotatedShareEntries {
			entryResp, err := u.shareEntryRepo.GetShareEntry(txCtx, &c3_asset_domain.GetShareEntryRequest{
				ShareEntryID: item.ShareEntryID,
			})
			if err != nil || entryResp == nil || entryResp.Data.ID == "" {
				return fmt.Errorf("failed to find share entry %s during rotation: %w", item.ShareEntryID, err)
			}

			entry := entryResp.Data
			entry.WrappedDEK = item.ReWrappedDEK
			entry.KEKVersion = req.NewVersion

			_, err = u.shareEntryRepo.UpdateShareEntry(txCtx, &c3_asset_domain.UpdateShareEntryRequest{
				ShareEntry: entry,
			})
			if err != nil {
				return fmt.Errorf("failed to update share entry %s during rotation: %w", item.ShareEntryID, err)
			}
			updatedCount++
		}
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	finalResponse := &RotateTrustGroupKEKResponse{
		TrustGroup: finalTrustGroup,
		Count:      updatedCount,
	}

	// Cache idempotent response
	if req.RequestID != "" {
		u.mu.Lock()
		u.processedRequests[req.RequestID] = finalResponse
		u.mu.Unlock()
	}

	return finalResponse, nil
}

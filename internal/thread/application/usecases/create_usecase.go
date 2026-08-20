package thread_usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	channel_domain "vault-app/internal/channel/domain"
	thread_dtos "vault-app/internal/thread/application/dtos"
	thread_events "vault-app/internal/thread/application/events"
	thread_domain "vault-app/internal/thread/domain"
	tracecore_types "vault-app/internal/tracecore/types"
)

// ChannelGovernanceReader abstracts the read interface required by Thread creation
// to inspect parent Channel state without giving Thread full Channel persistence capabilities.
type ChannelGovernanceReader interface {
	GetChannel(ctx context.Context, req *channel_domain.GetChannelRequest) (*tracecore_types.CloudResponse[channel_domain.Channel], error)
}

type CreateThreadUseCase interface {
	Execute(ctx context.Context, cmd thread_dtos.CreateThreadRequest) (*thread_domain.Thread, error)
}

type CreateThreadUsecase struct {
	Repo        thread_domain.ThreadRepository
	DomainBus   thread_events.ThreadEventBus
	ChannelReader ChannelGovernanceReader
}

func NewCreateThreadUsecase(
	repo thread_domain.ThreadRepository,
	threadBus thread_events.ThreadEventBus,
	channelReader ChannelGovernanceReader,
) *CreateThreadUsecase {
	return &CreateThreadUsecase{
		Repo:          repo,
		DomainBus:     threadBus,
		ChannelReader: channelReader,
	}
}

func (uc *CreateThreadUsecase) Execute(ctx context.Context, req thread_dtos.CreateThreadRequest) (*thread_domain.Thread, error) {
	if err := uc.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := uc.ValidateRequest(req); err != nil {
		return nil, err
	}

	// Governance check: load parent channel state if reader is provided
	var workspaceID = req.WorkspaceID
	if uc.ChannelReader != nil {
		resp, err := uc.ChannelReader.GetChannel(ctx, &channel_domain.GetChannelRequest{
			ChannelID: req.ChannelID,
		})
		if err != nil || resp == nil {
			return nil, thread_domain.ErrChannelNotFound
		}

		channel := &resp.Data
		if channel.Status != channel_domain.StatusActive {
			return nil, thread_domain.ErrChannelNotActive
		}

		// Verify gated slots are fulfilled
		gatedSlots := channel.GetGatedSlots()
		for _, slot := range gatedSlots {
			if _, assigned := channel.GetAssignmentBySlotID(slot.ID); !assigned {
				return nil, thread_domain.ErrChannelGatedSlotsIncomplete
			}
		}

		// Channel is authoritative for WorkspaceID
		if req.WorkspaceID != "" && req.WorkspaceID != channel.WorkspaceID {
			return nil, thread_domain.ErrWorkspaceMismatch
		}
		workspaceID = channel.WorkspaceID
	}

	thread := thread_domain.NewThread(req.ChannelID, req.AssetType, req.Title, req.Subtitle)
	thread.WorkspaceID = workspaceID

	created, err := uc.Repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{
		Thread:     thread,
		IdentityID: req.IdentityID,
	})
	if err != nil {
		return nil, err
	}
	if created == nil {
		return nil, thread_domain.ErrRepositoryResponse
	}

	errEvent := uc.DomainBus.PublishThreadCreated(
		ctx,
		thread_domain.ThreadCreated{
			EventID:     uuid.NewString(),
			Timestamp:   time.Now(),
			ThreadID:    thread.ID,
			ChannelID:   req.ChannelID,
			WorkspaceID: workspaceID,
			AssetType:   req.AssetType,
		},
	)
	if errEvent != nil {
		return nil, errEvent
	}

	return &created.Data, nil
}

func (c *CreateThreadUsecase) ValidateDependencies() error {
	if c.Repo == nil {
		return thread_domain.ErrRepositoryNil
	}

	if c.DomainBus == nil {
		return thread_domain.ErrThreadBusRequired
	}

	return nil
}

func (c *CreateThreadUsecase) ValidateRequest(req thread_dtos.CreateThreadRequest) error {
	if req.ChannelID == "" {
		return thread_domain.ErrChannelIDRequired
	}

	if req.AssetType == "" {
		return thread_domain.ErrAssetTypeRequired
	}

	if req.Title == "" {
		return thread_domain.ErrThreadTitleRequired
	}

	return nil
}

package thread_usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	thread_dtos "vault-app/internal/thread/application/dtos"
	thread_events "vault-app/internal/thread/application/events"
	thread_domain "vault-app/internal/thread/domain"
)

type CreateThreadUseCase interface {
	Execute(ctx context.Context, cmd thread_dtos.CreateThreadRequest) (thread_domain.Thread, error)
}

type CreateThreadUsecase struct {
	Repo      thread_domain.ThreadRepository
	DomainBus thread_events.ThreadEventBus
}

func NewCreateThreadUsecase(repo thread_domain.ThreadRepository, threadBus thread_events.ThreadEventBus) *CreateThreadUsecase {
	return &CreateThreadUsecase{
		Repo:      repo,
		DomainBus: threadBus,
	}
}

func (uc *CreateThreadUsecase) Execute(ctx context.Context, req thread_dtos.CreateThreadRequest) (*thread_domain.Thread, error) {
	if err := uc.ValidateDependencies(); err != nil {
		return nil, err
	}

	if err := uc.ValidateRequest(req); err != nil {
		return nil, err
	}

	thread := thread_domain.NewThread(req.ChannelID, req.AssetType, req.Title, req.Subtitle)

	created, err := uc.Repo.CreateThread(ctx, &thread_domain.CreateThreadRequest{
		Thread: thread,
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
			EventID:   uuid.NewString(),
			Timestamp: time.Now(),
			ThreadID:    thread.ID,
			ChannelId: req.ChannelID,
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

	if req.Subtitle == "" {
		return thread_domain.ErrThreadSubtitleRequired
	}

	return nil
}

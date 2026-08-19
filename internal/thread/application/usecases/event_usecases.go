package thread_usecase

import (
	"context"
	"errors"

	thread_domain "vault-app/internal/thread/domain"
)

type ListThreadEventsUsecase struct {
	Repo thread_domain.ThreadRepository
}

func NewListThreadEventsUsecase(repo thread_domain.ThreadRepository) *ListThreadEventsUsecase {
	return &ListThreadEventsUsecase{
		Repo: repo,
	}
}

func (uc *ListThreadEventsUsecase) Execute(ctx context.Context, threadID string) ([]thread_domain.ThreadEvent, error) {
	if uc.Repo == nil {
		return nil, errors.New("repository is required")
	}
	if threadID == "" {
		return nil, errors.New("thread id is required")
	}

	resp, err := uc.Repo.ListThreadEvents(ctx, &thread_domain.ListThreadEventsRequest{
		ThreadID: threadID,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return []thread_domain.ThreadEvent{}, nil
	}

	return resp.Data, nil
}

type AppendThreadEventUsecase struct {
	Repo thread_domain.ThreadRepository
}

func NewAppendThreadEventUsecase(repo thread_domain.ThreadRepository) *AppendThreadEventUsecase {
	return &AppendThreadEventUsecase{
		Repo: repo,
	}
}

func (uc *AppendThreadEventUsecase) Execute(
	ctx context.Context,
	threadID string,
	eventType string,
	payload thread_domain.EventResourceRef,
	idempotencyKey ...string,
) (*thread_domain.ThreadEvent, error) {
	if uc.Repo == nil {
		return nil, errors.New("repository is required")
	}
	if threadID == "" {
		return nil, errors.New("thread id is required")
	}
	if eventType == "" {
		return nil, errors.New("event type is required")
	}

	key := ""
	if len(idempotencyKey) > 0 && idempotencyKey[0] != "" {
		key = idempotencyKey[0]
	} else if payload.RefType == thread_domain.ResourceShareEntry && payload.ShareEntryID != "" {
		key = "evt_share_" + payload.ShareEntryID
	}

	resp, err := uc.Repo.AppendThreadEvent(ctx, &thread_domain.AppendThreadEventRequest{
		ThreadID:       threadID,
		EventType:      eventType,
		Payload:        payload,
		IdempotencyKey: key,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("empty repository response")
	}

	return &resp.Data, nil
}

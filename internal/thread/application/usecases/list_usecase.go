package thread_usecase

import (
	"context"
	"errors"

	thread_domain "vault-app/internal/thread/domain"
)

type ListThreadsUsecase struct {
	Repo thread_domain.ThreadRepository
}

func NewListThreadsUsecase(repo thread_domain.ThreadRepository) *ListThreadsUsecase {
	return &ListThreadsUsecase{
		Repo: repo,
	}
}

func (uc *ListThreadsUsecase) Execute(ctx context.Context, channelID string) ([]thread_domain.Thread, error) {
	if uc.Repo == nil {
		return nil, errors.New("repository is required")
	}
	if channelID == "" {
		return nil, errors.New("channel id is required")
	}

	resp, err := uc.Repo.ListThreads(ctx, &thread_domain.ListThreadsRequest{
		ChannelID: channelID,
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return []thread_domain.Thread{}, nil
	}

	return resp.Data, nil
}

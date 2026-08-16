package channel_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	channel_application "vault-app/internal/channel/application"
	channel_usecase "vault-app/internal/channel/application/channel_lifecycle_usecases"
	channel_domain "vault-app/internal/channel/domain"
)

func TestRevokeChannelUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{
		revokeFn: func(
			ctx context.Context,
			req *channel_domain.RevokeChannelRequest,
		) error {
			require.Equal(t, "channel-001", req.ChannelID)
			return nil
		},
	}

	uc := channel_usecase.NewRevokeChannelUsecase(repo)

	err := uc.Execute(ctx, &channel_application.RevokeChannelRequest{
		ChannelID: "channel-001",
	})

	require.NoError(t, err)
}

func TestRevokeChannelUsecase_Execute_RepositoryErrorPropagates(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("Cloud backend returned status 400: channel already revoked")

	repo := &channelRepositoryMock{
		revokeFn: func(
			ctx context.Context,
			req *channel_domain.RevokeChannelRequest,
		) error {
			return expectedErr
		},
	}

	uc := channel_usecase.NewRevokeChannelUsecase(repo)

	err := uc.Execute(ctx, &channel_application.RevokeChannelRequest{
		ChannelID: "channel-001",
	})

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
}

func TestRevokeChannelUsecase_Execute_NilRepository(t *testing.T) {
	ctx := context.Background()

	uc := channel_usecase.NewRevokeChannelUsecase(nil)

	err := uc.Execute(ctx, &channel_application.RevokeChannelRequest{
		ChannelID: "channel-001",
	})

	require.Error(t, err)
	require.EqualError(t, err, channel_domain.ErrRepositoryNil.Error())
}

func TestRevokeChannelUsecase_Execute_NilRequest(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{}
	uc := channel_usecase.NewRevokeChannelUsecase(repo)

	err := uc.Execute(ctx, nil)

	require.Error(t, err)
	require.EqualError(t, err, channel_domain.ErrRequestRequired.Error())
}

func TestRevokeChannelUsecase_Execute_MissingChannelID(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{}
	uc := channel_usecase.NewRevokeChannelUsecase(repo)

	err := uc.Execute(ctx, &channel_application.RevokeChannelRequest{})

	require.Error(t, err)
	require.EqualError(t, err, channel_domain.ErrChannelIDRequired.Error())
}

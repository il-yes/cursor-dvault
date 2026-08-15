package channel_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	channel_application "vault-app/internal/channel/application"
	channel_usecase "vault-app/internal/channel/application/channel_lifecycle_usecases"
	channel_domain "vault-app/internal/channel/domain"
	tracecore_types "vault-app/internal/tracecore/types"
)

func TestActivateChannelUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	activated := channel_domain.Channel{
		ID:          "channel-001",
		TemplateID:  "tpl-001",
		Title:       "Contract Execution",
		Status:      channel_domain.StatusActive,
		WorkspaceID: "workspace-001",
	}

	repo := &channelRepositoryMock{
		activateFn: func(
			ctx context.Context,
			req *channel_domain.ActivateChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			require.Equal(t, "channel-001", req.ChannelID)

			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: activated,
			}, nil
		},
	}

	uc := channel_usecase.NewActivateChannelUsecase(repo)

	result, err := uc.Execute(ctx, &channel_application.ActivateChannelRequest{
		ChannelID: "channel-001",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, activated, *result)
}

func TestActivateChannelUsecase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("activate channel: validate activation: gated slots unfulfilled")

	repo := &channelRepositoryMock{
		activateFn: func(
			ctx context.Context,
			req *channel_domain.ActivateChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return nil, expectedErr
		},
	}

	uc := channel_usecase.NewActivateChannelUsecase(repo)

	result, err := uc.Execute(ctx, &channel_application.ActivateChannelRequest{
		ChannelID: "channel-001",
	})

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)
}

func TestActivateChannelUsecase_Execute_NilRepositoryResponse(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{
		activateFn: func(
			ctx context.Context,
			req *channel_domain.ActivateChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return nil, nil
		},
	}

	uc := channel_usecase.NewActivateChannelUsecase(repo)

	result, err := uc.Execute(ctx, &channel_application.ActivateChannelRequest{
		ChannelID: "channel-001",
	})

	require.Error(t, err)
	require.ErrorIs(t, err, channel_domain.ErrRepositoryResponse)
	require.Nil(t, result)
}

func TestActivateChannelUsecase_Execute_NilRepository(t *testing.T) {
	uc := channel_usecase.NewActivateChannelUsecase(nil)

	ctx := context.Background()

	result, err := uc.Execute(ctx, &channel_application.ActivateChannelRequest{
		ChannelID: "channel-001",
	})

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, channel_domain.ErrRepositoryNil.Error())
}

func TestActivateChannelUsecase_Execute_NilRequest(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{}

	uc := channel_usecase.NewActivateChannelUsecase(repo)

	result, err := uc.Execute(ctx, nil)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, channel_domain.ErrRequestRequired.Error())
}

func TestActivateChannelUsecase_Execute_MissingChannelID(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{}

	uc := channel_usecase.NewActivateChannelUsecase(repo)

	result, err := uc.Execute(ctx, &channel_application.ActivateChannelRequest{})

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, channel_domain.ErrChannelIDRequired.Error())
}

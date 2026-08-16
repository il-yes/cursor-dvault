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

func validAddParticipantRequest() *channel_application.AddParticipantRequest {
	return &channel_application.AddParticipantRequest{
		ChannelID: "channel-001",
		VaultID:   "vault_external",
		PublicKey: "GABC...",
		Direction: "bidirectional",
		SlotID:    "slot_supplier",
		Role:      "supplier",
	}
}

func TestAddParticipantUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	participant := channel_domain.Participant{
		ChannelID:   "channel-001",
		VaultID:     "vault_external",
		PublicKey:   "GABC...",
		Direction:   "bidirectional",
		JoinedAt:    1750000000,
		Role:        "supplier",
		Permissions: []string{"read", "write"},
	}

	repo := &channelRepositoryMock{
		addParticipantFn: func(
			ctx context.Context,
			req *channel_domain.JoinChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Participant], error) {
			require.Equal(t, "channel-001", req.ChannelID)
			require.Equal(t, "vault_external", req.VaultID)
			require.Equal(t, "GABC...", req.PublicKey)
			require.Equal(t, "bidirectional", req.Direction)
			require.Equal(t, "slot_supplier", req.SlotID)
			require.Equal(t, "supplier", req.Role)

			return &tracecore_types.CloudResponse[channel_domain.Participant]{
				Data: participant,
			}, nil
		},
	}

	uc := channel_usecase.NewAddParticipantUsecase(repo)

	result, err := uc.Execute(ctx, validAddParticipantRequest())

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, participant, *result)
}

func TestAddParticipantUsecase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("add participant: create participant: duplicate")

	repo := &channelRepositoryMock{
		addParticipantFn: func(
			ctx context.Context,
			req *channel_domain.JoinChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Participant], error) {
			return nil, expectedErr
		},
	}

	uc := channel_usecase.NewAddParticipantUsecase(repo)

	result, err := uc.Execute(ctx, validAddParticipantRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)
}

func TestAddParticipantUsecase_Execute_NilRepositoryResponse(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{
		addParticipantFn: func(
			ctx context.Context,
			req *channel_domain.JoinChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Participant], error) {
			return nil, nil
		},
	}

	uc := channel_usecase.NewAddParticipantUsecase(repo)

	result, err := uc.Execute(ctx, validAddParticipantRequest())

	require.Error(t, err)
	require.EqualError(t, err, channel_domain.ErrRepositoryResponse.Error())
	require.Nil(t, result)
}

func TestAddParticipantUsecase_Execute_NilRepository(t *testing.T) {
	uc := channel_usecase.NewAddParticipantUsecase(nil)

	ctx := context.Background()

	result, err := uc.Execute(ctx, validAddParticipantRequest())

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, channel_domain.ErrRepositoryNil.Error())
}

func TestAddParticipantUsecase_Execute_NilRequest(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{}

	uc := channel_usecase.NewAddParticipantUsecase(repo)

	result, err := uc.Execute(ctx, nil)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, channel_domain.ErrRequestRequired.Error())
}

func TestAddParticipantUsecase_Execute_MissingChannelID(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{}

	uc := channel_usecase.NewAddParticipantUsecase(repo)

	req := validAddParticipantRequest()
	req.ChannelID = ""

	result, err := uc.Execute(ctx, req)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, channel_domain.ErrChannelIDRequired.Error())
}

func TestAddParticipantUsecase_Execute_MissingVaultID(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{}

	uc := channel_usecase.NewAddParticipantUsecase(repo)

	req := validAddParticipantRequest()
	req.VaultID = ""

	result, err := uc.Execute(ctx, req)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, channel_domain.ErrVaultIDRequired.Error())
}

func TestAddParticipantUsecase_ValidateRequest(t *testing.T) {
	repo := &channelRepositoryMock{}

	uc := channel_usecase.NewAddParticipantUsecase(repo)

	tests := []struct {
		name          string
		request       *channel_application.AddParticipantRequest
		expectedError string
	}{
		{
			name:          "nil request",
			request:       nil,
			expectedError: channel_domain.ErrRequestRequired.Error(),
		},
		{
			name: "missing channel id",
			request: &channel_application.AddParticipantRequest{
				VaultID: "vault_external",
			},
			expectedError: channel_domain.ErrChannelIDRequired.Error(),
		},
		{
			name: "missing vault id",
			request: &channel_application.AddParticipantRequest{
				ChannelID: "channel-001",
			},
			expectedError: channel_domain.ErrVaultIDRequired.Error(),
		},
		{
			name:    "valid request",
			request: validAddParticipantRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.ValidateRequest(tt.request)

			if tt.expectedError == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.expectedError)
		})
	}
}

func TestAddParticipantUsecase_ValidateDependencies(t *testing.T) {
	tests := []struct {
		name          string
		repo          channel_domain.ChannelRepository
		expectedError string
	}{
		{
			name:          "nil repository",
			repo:          nil,
			expectedError: channel_domain.ErrRepositoryNil.Error(),
		},
		{
			name: "valid repository",
			repo: &channelRepositoryMock{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := channel_usecase.NewAddParticipantUsecase(tt.repo)

			err := uc.ValidateDependencies()

			if tt.expectedError == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.expectedError)
		})
	}
}

func validListParticipantsRequest() *channel_application.ListParticipantsRequest {
	return &channel_application.ListParticipantsRequest{
		ChannelID: "channel-001",
	}
}

func TestListParticipantsUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	expected := []channel_domain.Participant{
		{
			ChannelID:   "channel-001",
			VaultID:     "vault_internal",
			PublicKey:   "GINT...",
			Direction:   "inbound",
			JoinedAt:    1750000000,
			Role:        "owner",
			Permissions: []string{"read", "write"},
		},
		{
			ChannelID:   "channel-001",
			VaultID:     "vault_external",
			PublicKey:   "GABC...",
			Direction:   "bidirectional",
			JoinedAt:    1750000100,
			Role:        "supplier",
			Permissions: []string{"read"},
		},
	}

	repo := &channelRepositoryMock{
		listParticipantsFn: func(
			ctx context.Context,
			req *channel_domain.ListParticipantsRequest,
		) (*tracecore_types.CloudResponse[[]channel_domain.Participant], error) {
			require.Equal(t, "channel-001", req.ChannelID)

			return &tracecore_types.CloudResponse[[]channel_domain.Participant]{
				Data: expected,
			}, nil
		},
	}

	uc := channel_usecase.NewListParticipantsUsecase(repo)

	result, err := uc.Execute(ctx, validListParticipantsRequest())

	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, expected, result)
}

func TestListParticipantsUsecase_Execute_EmptyCollection(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{
		listParticipantsFn: func(
			ctx context.Context,
			req *channel_domain.ListParticipantsRequest,
		) (*tracecore_types.CloudResponse[[]channel_domain.Participant], error) {
			return &tracecore_types.CloudResponse[[]channel_domain.Participant]{
				Data: []channel_domain.Participant{},
			}, nil
		},
	}

	uc := channel_usecase.NewListParticipantsUsecase(repo)

	result, err := uc.Execute(ctx, validListParticipantsRequest())

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Empty(t, result)
}

func TestListParticipantsUsecase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("failed to list participants")

	repo := &channelRepositoryMock{
		listParticipantsFn: func(
			ctx context.Context,
			req *channel_domain.ListParticipantsRequest,
		) (*tracecore_types.CloudResponse[[]channel_domain.Participant], error) {
			return nil, expectedErr
		},
	}

	uc := channel_usecase.NewListParticipantsUsecase(repo)

	result, err := uc.Execute(ctx, validListParticipantsRequest())

	require.Error(t, err)
	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)
}

func TestListParticipantsUsecase_Execute_NilRepository(t *testing.T) {
	uc := channel_usecase.NewListParticipantsUsecase(nil)

	ctx := context.Background()

	result, err := uc.Execute(ctx, validListParticipantsRequest())

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, channel_domain.ErrRepositoryNil.Error())
}

func TestListParticipantsUsecase_Execute_NilRequest(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{}

	uc := channel_usecase.NewListParticipantsUsecase(repo)

	result, err := uc.Execute(ctx, nil)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, channel_domain.ErrRequestRequired.Error())
}

func TestListParticipantsUsecase_Execute_MissingChannelID(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{}

	uc := channel_usecase.NewListParticipantsUsecase(repo)

	req := validListParticipantsRequest()
	req.ChannelID = ""

	result, err := uc.Execute(ctx, req)

	require.Error(t, err)
	require.Nil(t, result)
	require.EqualError(t, err, channel_domain.ErrChannelIDRequired.Error())
}

func TestListParticipantsUsecase_ValidateRequest(t *testing.T) {
	repo := &channelRepositoryMock{}

	uc := channel_usecase.NewListParticipantsUsecase(repo)

	tests := []struct {
		name          string
		request       *channel_application.ListParticipantsRequest
		expectedError string
	}{
		{
			name:          "nil request",
			request:       nil,
			expectedError: channel_domain.ErrRequestRequired.Error(),
		},
		{
			name: "missing channel id",
			request: &channel_application.ListParticipantsRequest{
				ChannelID: "",
			},
			expectedError: channel_domain.ErrChannelIDRequired.Error(),
		},
		{
			name:    "valid request",
			request: validListParticipantsRequest(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.ValidateRequest(tt.request)

			if tt.expectedError == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.expectedError)
		})
	}
}

func TestListParticipantsUsecase_ValidateDependencies(t *testing.T) {
	tests := []struct {
		name          string
		repo          channel_domain.ChannelRepository
		expectedError string
	}{
		{
			name:          "nil repository",
			repo:          nil,
			expectedError: channel_domain.ErrRepositoryNil.Error(),
		},
		{
			name: "valid repository",
			repo: &channelRepositoryMock{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := channel_usecase.NewListParticipantsUsecase(tt.repo)

			err := uc.ValidateDependencies()

			if tt.expectedError == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.expectedError)
		})
	}
}

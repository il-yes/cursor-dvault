package channel_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	channel_application "vault-app/internal/channel/application"
	channel_usecase "vault-app/internal/channel/application/channel_lifecycle_usecases"
	channel_domain "vault-app/internal/channel/domain"
	tracecore_types "vault-app/internal/tracecore/types"
)

func validInviteToChannelRequest() *channel_application.InviteToChannelRequest {
	return &channel_application.InviteToChannelRequest{
		ChannelID:      "channel-001",
		InviterVaultID: "vault_owner",
		InviteeVaultID: "vault_external",
	}
}

func TestInviteToChannelUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	invitation := channel_domain.Invitation{
		ID:             "inv_9e24294d-2598-4be6-a649-9f0a56896a03",
		ChannelID:      "channel-001",
		InviterVaultID: "vault_owner",
		InviteeVaultID: "vault_external",
		Status:         channel_domain.InvitationStatusPending,
		CreatedAt:      time.Now().UTC(),
	}

	repo := &channelRepositoryMock{
		inviteToChannelFn: func(
			ctx context.Context,
			req *channel_domain.InviteToChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Invitation], error) {
			require.Equal(t, "channel-001", req.ChannelID)
			require.Equal(t, "vault_owner", req.InviterVaultID)
			require.Equal(t, "vault_external", req.InviteeVaultID)

			return &tracecore_types.CloudResponse[channel_domain.Invitation]{
				Data: invitation,
			}, nil
		},
	}

	uc := channel_usecase.NewInviteToChannelUsecase(repo)

	result, err := uc.Execute(ctx, validInviteToChannelRequest())

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, invitation, *result)
}

func TestInviteToChannelUsecase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("create invitation: channel revoked")

	repo := &channelRepositoryMock{
		inviteToChannelFn: func(
			ctx context.Context,
			req *channel_domain.InviteToChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Invitation], error) {
			return nil, expectedErr
		},
	}

	uc := channel_usecase.NewInviteToChannelUsecase(repo)

	_, err := uc.Execute(ctx, validInviteToChannelRequest())

	require.ErrorIs(t, err, expectedErr)
}

func TestInviteToChannelUsecase_Execute_RepositoryNilResponse(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{
		inviteToChannelFn: func(
			ctx context.Context,
			req *channel_domain.InviteToChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Invitation], error) {
			return nil, nil
		},
	}

	uc := channel_usecase.NewInviteToChannelUsecase(repo)

	_, err := uc.Execute(ctx, validInviteToChannelRequest())

	require.ErrorIs(t, err, channel_domain.ErrRepositoryResponse)
}

func TestInviteToChannelUsecase_Execute_NilRepo(t *testing.T) {
	ctx := context.Background()

	uc := channel_usecase.NewInviteToChannelUsecase(nil)

	_, err := uc.Execute(ctx, validInviteToChannelRequest())

	require.ErrorIs(t, err, channel_domain.ErrRepositoryNil)
}

func TestInviteToChannelUsecase_ValidateRequest(t *testing.T) {
	uc := channel_usecase.NewInviteToChannelUsecase(&channelRepositoryMock{})

	cases := []struct {
		name string
		req  *channel_application.InviteToChannelRequest
		want error
	}{
		{"nil request", nil, channel_domain.ErrRequestRequired},
		{"missing channel", &channel_application.InviteToChannelRequest{InviterVaultID: "vault_owner", InviteeVaultID: "vault_external"}, channel_domain.ErrChannelIDRequired},
		{"missing inviter", &channel_application.InviteToChannelRequest{ChannelID: "channel-001", InviteeVaultID: "vault_external"}, channel_domain.ErrVaultIDRequired},
		{"missing invitee", &channel_application.InviteToChannelRequest{ChannelID: "channel-001", InviterVaultID: "vault_owner"}, channel_domain.ErrVaultIDRequired},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := uc.ValidateRequest(tc.req)
			require.ErrorIs(t, err, tc.want)
		})
	}

	require.NoError(t, uc.ValidateRequest(validInviteToChannelRequest()))
}

func validAcceptChannelInvitationRequest() *channel_application.AcceptChannelInvitationRequest {
	return &channel_application.AcceptChannelInvitationRequest{
		InvitationID:     "inv_9e24294d-2598-4be6-a649-9f0a56896a03",
		InviteeVaultID:   "vault_external",
		InviteePublicKey: "GEXT...",
	}
}

func TestAcceptChannelInvitationUsecase_Execute_Success(t *testing.T) {
	ctx := context.Background()

	acceptedAt := time.Now().UTC()

	invitation := channel_domain.Invitation{
		ID:             "inv_9e24294d-2598-4be6-a649-9f0a56896a03",
		ChannelID:      "channel-001",
		InviterVaultID: "vault_owner",
		InviteeVaultID: "vault_external",
		Status:         channel_domain.InvitationStatusAccepted,
		CreatedAt:      time.Now().UTC(),
		AcceptedAt:     &acceptedAt,
	}

	repo := &channelRepositoryMock{
		acceptChannelInvitationFn: func(
			ctx context.Context,
			req *channel_domain.AcceptInvitationRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Invitation], error) {
			require.Equal(t, "inv_9e24294d-2598-4be6-a649-9f0a56896a03", req.InvitationID)
			require.Equal(t, "vault_external", req.InviteeVaultID)
			require.Equal(t, "GEXT...", req.InviteePublicKey)

			return &tracecore_types.CloudResponse[channel_domain.Invitation]{
				Data: invitation,
			}, nil
		},
	}

	uc := channel_usecase.NewAcceptChannelInvitationUsecase(repo)

	result, err := uc.Execute(ctx, validAcceptChannelInvitationRequest())

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, invitation, *result)
}

func TestAcceptChannelInvitationUsecase_Execute_RepositoryError(t *testing.T) {
	ctx := context.Background()

	expectedErr := errors.New("accept invitation: invitation not for you")

	repo := &channelRepositoryMock{
		acceptChannelInvitationFn: func(
			ctx context.Context,
			req *channel_domain.AcceptInvitationRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Invitation], error) {
			return nil, expectedErr
		},
	}

	uc := channel_usecase.NewAcceptChannelInvitationUsecase(repo)

	_, err := uc.Execute(ctx, validAcceptChannelInvitationRequest())

	require.ErrorIs(t, err, expectedErr)
}

func TestAcceptChannelInvitationUsecase_Execute_RepositoryNilResponse(t *testing.T) {
	ctx := context.Background()

	repo := &channelRepositoryMock{
		acceptChannelInvitationFn: func(
			ctx context.Context,
			req *channel_domain.AcceptInvitationRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Invitation], error) {
			return nil, nil
		},
	}

	uc := channel_usecase.NewAcceptChannelInvitationUsecase(repo)

	_, err := uc.Execute(ctx, validAcceptChannelInvitationRequest())

	require.ErrorIs(t, err, channel_domain.ErrRepositoryResponse)
}

func TestAcceptChannelInvitationUsecase_Execute_NilRepo(t *testing.T) {
	ctx := context.Background()

	uc := channel_usecase.NewAcceptChannelInvitationUsecase(nil)

	_, err := uc.Execute(ctx, validAcceptChannelInvitationRequest())

	require.ErrorIs(t, err, channel_domain.ErrRepositoryNil)
}

func TestAcceptChannelInvitationUsecase_ValidateRequest(t *testing.T) {
	uc := channel_usecase.NewAcceptChannelInvitationUsecase(&channelRepositoryMock{})

	cases := []struct {
		name string
		req  *channel_application.AcceptChannelInvitationRequest
		want error
	}{
		{"nil request", nil, channel_domain.ErrRequestRequired},
		{"missing invitation", &channel_application.AcceptChannelInvitationRequest{InviteeVaultID: "vault_external", InviteePublicKey: "GEXT..."}, channel_domain.ErrInvitationIDRequired},
		{"missing invitee", &channel_application.AcceptChannelInvitationRequest{InvitationID: "inv_1", InviteePublicKey: "GEXT..."}, channel_domain.ErrVaultIDRequired},
		{"missing public key", &channel_application.AcceptChannelInvitationRequest{InvitationID: "inv_1", InviteeVaultID: "vault_external"}, channel_domain.ErrInviteePublicKeyRequired},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := uc.ValidateRequest(tc.req)
			require.ErrorIs(t, err, tc.want)
		})
	}

	require.NoError(t, uc.ValidateRequest(validAcceptChannelInvitationRequest()))
}

package channel_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	channel_application "vault-app/internal/channel/application"
	"vault-app/internal/channel/application/channel_slot_usecases"
	channel_domain "vault-app/internal/channel/domain"
	tracecore_types "vault-app/internal/tracecore/types"
)




func newTestChannel() channel_domain.Channel {
	now := time.Now().UTC()

	return channel_domain.Channel{
		ID:        uuid.NewString(),
		Status:    channel_domain.StatusActive,
		Slots:     []channel_domain.Slot{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestChannelSlotUsecase_AddSlot_Success(t *testing.T) {
	channel := newTestChannel()

	slot := channel_domain.Slot{
		ID:      uuid.NewString(),
		Role:    "engineer",
		VaultID: uuid.NewString(),
		Gated:   false,
	}

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			require.Equal(t, channel.ID, req.ChannelID)

			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: channel,
			}, nil
		},

		updateFn: func(
			ctx context.Context,
			req *channel_domain.UpdateChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			require.Len(t, req.Channel.Slots, 1)
			require.Equal(t, slot.ID, req.Channel.Slots[0].ID)

			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: req.Channel,
			}, nil
		},
	}

	uc := channel_slot_usecases.NewChannelSlotUsecase(repo)

	result, err := uc.AddSlot(
		context.Background(),
		channel_application.AddSlotRequest{
			ChannelID: channel.ID,
			Slot:      slot,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, result.Slots, 1)
	require.Equal(t, slot.ID, result.Slots[0].ID)
}

func TestChannelSlotUsecase_AddSlot_GetChannelError(t *testing.T) {
	expectedErr := errors.New("repository error")

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return nil, expectedErr
		},
	}

	uc := channel_slot_usecases.NewChannelSlotUsecase(repo)

	result, err := uc.AddSlot(
		context.Background(),
		channel_application.AddSlotRequest{
			ChannelID: uuid.NewString(),
			Slot: channel_domain.Slot{
				ID: uuid.NewString(),
			},
		},
	)

	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)
}

func TestChannelSlotUsecase_AddSlot_ChannelNotFound(t *testing.T) {
	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: channel_domain.Channel{},
			}, nil
		},
	}

	uc := channel_slot_usecases.NewChannelSlotUsecase(repo)

	result, err := uc.AddSlot(
		context.Background(),
		channel_application.AddSlotRequest{
			ChannelID: uuid.NewString(),
			Slot: channel_domain.Slot{
				ID: uuid.NewString(),
			},
		},
	)

	require.ErrorIs(t, err, channel_domain.ErrChannelNotFound)
	require.Nil(t, result)
}

// func TestChannelSlotUsecase_AddSlot_AggregateError(t *testing.T) {
// 	channel := newTestChannel()

// 	slot := channel_domain.Slot{
// 		ID:   uuid.NewString(),
// 		Role: "engineer",
// 	}

// 	channel.Slots = append(channel.Slots, slot)

// 	repo := &channelRepositoryMock{
// 		getFn: func(
// 			ctx context.Context,
// 			req *channel_domain.GetChannelRequest,
// 		) (*channel_domain.GetChannelResponse, error) {
// 			return &channel_domain.GetChannelResponse{
// 				Data: channel,
// 			}, nil
// 		},
// 	}

// 	uc := channel_slot_usecases.NewChannelSlotUsecase(repo)

// 	result, err := uc.AddSlot(
// 		context.Background(),
// 		channel_application.AddSlotRequest{
// 			ChannelID: channel.ID,
// 			Slot:      slot,
// 		},
// 	)

// 	require.ErrorIs(t, err, channel_domain.ErrSlotAlreadyExists)
// 	require.Nil(t, result)
// }

// func TestChannelSlotUsecase_AddSlot_AggregateError(t *testing.T) {
// 	channel := newTestChannel()

// 	slot := channel_domain.Slot{
// 		ID:   uuid.NewString(),
// 		Role: "engineer",
// 	}

// 	channel.Slots = append(channel.Slots, slot)

// 	repo := &channelRepositoryMock{
// 		getFn: func(
// 			ctx context.Context,
// 			req *channel_domain.GetChannelRequest,
// 		) (*channel_domain.GetChannelResponse, error) {
// 			return &channel_domain.GetChannelResponse{
// 				Data: channel,
// 			}, nil
// 		},
// 	}

// 	uc := channel_slot_usecases.NewChannelSlotUsecase(repo)

// 	result, err := uc.AddSlot(
// 		context.Background(),
// 		channel_application.AddSlotRequest{
// 			ChannelID: channel.ID,
// 			Slot:      slot,
// 		},
// 	)

// 	require.ErrorIs(t, err, channel_domain.ErrSlotAlreadyExists)
// 	require.Nil(t, result)
// }

func TestChannelSlotUsecase_UpdateSlot_Success(t *testing.T) {
	channel := newTestChannel()

	slot := channel_domain.Slot{
		ID:      uuid.NewString(),
		Role:    "engineer",
		VaultID: uuid.NewString(),
	}

	channel.Slots = append(channel.Slots, slot)

	updatedSlot := slot
	updatedSlot.Role = "architect"
	updatedSlot.Gated = true

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: channel,
			}, nil
		},

		updateFn: func(
			ctx context.Context,
			req *channel_domain.UpdateChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			require.Len(t, req.Channel.Slots, 1)
			require.Equal(t, "architect", req.Channel.Slots[0].Role)
			require.True(t, req.Channel.Slots[0].Gated)

			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: req.Channel,
			}, nil
		},
	}

	uc := channel_slot_usecases.NewChannelSlotUsecase(repo)

	result, err := uc.UpdateSlot(
		context.Background(),
		channel_application.UpdateSlotRequest{
			ChannelID: channel.ID,
			Slot:      updatedSlot,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "architect", result.Slots[0].Role)
	require.True(t, result.Slots[0].Gated)
}

func TestChannelSlotUsecase_UpdateSlot_SlotNotFound(t *testing.T) {
	channel := newTestChannel()

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: channel,
			}, nil
		},
	}

	uc := channel_slot_usecases.NewChannelSlotUsecase(repo)

	result, err := uc.UpdateSlot(
		context.Background(),
		channel_application.UpdateSlotRequest{
			ChannelID: channel.ID,
			Slot: channel_domain.Slot{
				ID: uuid.NewString(),
			},
		},
	)

	require.ErrorIs(t, err, channel_domain.ErrSlotNotFound)
	require.Nil(t, result)
}

func TestChannelSlotUsecase_UpdateSlot_GetChannelError(t *testing.T) {
	expectedErr := errors.New("get failed")

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return nil, expectedErr
		},
	}

	uc := channel_slot_usecases.NewChannelSlotUsecase(repo)

	result, err := uc.UpdateSlot(
		context.Background(),
		channel_application.UpdateSlotRequest{
			ChannelID: uuid.NewString(),
			Slot: channel_domain.Slot{
				ID: uuid.NewString(),
			},
		},
	)

	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)
}

func TestChannelSlotUsecase_UpdateSlot_UpdateChannelError(t *testing.T) {
	expectedErr := errors.New("update failed")
	channel := newTestChannel()

	slot := channel_domain.Slot{
		ID: uuid.NewString(),
	}

	channel.Slots = append(channel.Slots, slot)

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: channel,
			}, nil
		},

		updateFn: func(
			ctx context.Context,
			req *channel_domain.UpdateChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return nil, expectedErr
		},
	}

	uc := channel_slot_usecases.NewChannelSlotUsecase(repo)

	result, err := uc.UpdateSlot(
		context.Background(),
		channel_application.UpdateSlotRequest{
			ChannelID: channel.ID,
			Slot:      slot,
		},
	)

	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)
}

func TestChannelSlotUsecase_RemoveSlot_Success(t *testing.T) {
	channel := newTestChannel()

	slot := channel_domain.Slot{
		ID:      uuid.NewString(),
		Role:    "engineer",
		VaultID: uuid.NewString(),
	}

	channel.Slots = append(channel.Slots, slot)

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: channel,
			}, nil
		},

		updateFn: func(
			ctx context.Context,
			req *channel_domain.UpdateChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			require.Empty(t, req.Channel.Slots)

			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: req.Channel,
			}, nil
		},
	}

	uc := channel_slot_usecases.NewChannelSlotUsecase(repo)

	result, err := uc.RemoveSlot(
		context.Background(),
		channel_application.RemoveSlotRequest{
			ChannelID: channel.ID,
			SlotID:    slot.ID,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Empty(t, result.Slots)
}


func TestChannelSlotUsecase_RemoveSlot_UpdateChannelError(t *testing.T) {
	expectedErr := errors.New("update failed")
	channel := newTestChannel()

	slot := channel_domain.Slot{
		ID: uuid.NewString(),
	}

	channel.Slots = append(channel.Slots, slot)

	repo := &channelRepositoryMock{
		getFn: func(
			ctx context.Context,
			req *channel_domain.GetChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return &tracecore_types.CloudResponse[channel_domain.Channel]{
				Data: channel,
			}, nil
		},

		updateFn: func(
			ctx context.Context,
			req *channel_domain.UpdateChannelRequest,
		) (*tracecore_types.CloudResponse[channel_domain.Channel], error) {
			return nil, expectedErr
		},
	}

	uc := channel_slot_usecases.NewChannelSlotUsecase(repo)

	result, err := uc.RemoveSlot(
		context.Background(),
		channel_application.RemoveSlotRequest{
			ChannelID: channel.ID,
			SlotID:    slot.ID,
		},
	)

	require.ErrorIs(t, err, expectedErr)
	require.Nil(t, result)
}

func TestChannelSlotUsecase_ValidateDependencies(t *testing.T) {
	uc := channel_slot_usecases.NewChannelSlotUsecase(nil)

	err := uc.ValidateDependencies()

	require.Error(t, err)
}

func TestChannelSlotUsecase_AddSlot_Validation(t *testing.T) {
	repo := &channelRepositoryMock{}
	uc := channel_slot_usecases.NewChannelSlotUsecase(repo)

	tests := []struct {
		name string
		req  channel_application.AddSlotRequest
	}{
		{
			name: "missing channel id",
			req: channel_application.AddSlotRequest{
				Slot: channel_domain.Slot{
					ID: uuid.NewString(),
				},
			},
		},
		{
			name: "missing slot id",
			req: channel_application.AddSlotRequest{
				ChannelID: uuid.NewString(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := uc.AddSlot(context.Background(), tt.req)

			require.Error(t, err)
			require.Nil(t, result)
		})
	}
}

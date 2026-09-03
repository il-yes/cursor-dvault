package realtime_client_handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	shared_realtime "vault-app/internal/shared/realtime"
	"vault-app/internal/utils"
)

type ChannelInvitationCreatedHandler struct {
	appCtx context.Context
}

func NewChannelInvitationCreatedHandler(
	appCtx context.Context,
) *ChannelInvitationCreatedHandler {
	return &ChannelInvitationCreatedHandler{
		appCtx: appCtx,
	}
}

func (h *ChannelInvitationCreatedHandler) Handle(
	ctx context.Context,
	msg shared_realtime.Message,
) error {
	utils.LogPretty("[C3][REALTIME] channel.invitation.created received by dVault", msg)
	fmt.Printf("[C3][REALTIME] channel.invitation.created received by dVault type=%s\n", msg.Type)

	var payload any
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		payload = string(msg.Payload)
	}

	if h.appCtx != nil {
		fmt.Printf("[C3][WAILS] channel.invitation.created emitted\n")
		runtime.EventsEmit(
			h.appCtx,
			"channel.invitation.created",
			payload,
		)
	}

	return nil
}

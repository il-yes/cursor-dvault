package realtime_client_handlers

import (
	"context"
	"encoding/json"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	shared_realtime "vault-app/internal/shared/realtime"
	"vault-app/internal/utils"
)

type ShareInvitationHandler struct {
	appCtx context.Context
}

func NewShareInvitationHandler(
	appCtx context.Context,
) *ShareInvitationHandler {
	return &ShareInvitationHandler{
		appCtx: appCtx,
	}
}

func (h *ShareInvitationHandler) Handle(
	ctx context.Context,
	msg shared_realtime.Message,
) error {

	var payload shared_realtime.ShareInvitationNotificationPayload

	data, err := json.Marshal(msg.Payload)
	if err != nil {
		return err
	}
	utils.LogPretty("MESSAGE TYPE =", msg.Type)

	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	runtime.EventsEmit(
		h.appCtx,
		msg.Type,
		payload,
	)

	return nil
}

package realtime_client_handlers

import (
	"context"
	"encoding/json"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	shared_realtime "vault-app/internal/shared/realtime"
	"vault-app/internal/utils"
)

type ShareRejectedHandler struct {
	appCtx context.Context
}

func NewShareRejectedHandler(
	appCtx context.Context,
) *ShareRejectedHandler {
	return &ShareRejectedHandler{
		appCtx: appCtx,
	}
}

func (h *ShareRejectedHandler) Handle(
	ctx context.Context,
	msg shared_realtime.Message,
) error {

	var payload shared_realtime.ShareRejectedPayload

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

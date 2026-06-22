package realtime_client_handlers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	identity_domain "vault-app/internal/identity/domain"
	shared_realtime "vault-app/internal/shared/realtime"
	"vault-app/internal/utils"
)

type ShareReadyHandler struct {
	appCtx context.Context
}

func NewShareReadyHandler(
	appCtx context.Context,
) *ShareReadyHandler {
	return &ShareReadyHandler{
		appCtx: appCtx,
	}
}

func (h *ShareReadyHandler) Handle(
	ctx context.Context,
	msg shared_realtime.Message,
) error {

	var payload shared_realtime.ShareReadyToAcceptPayload

	data, err := json.Marshal(msg.Payload)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	user, ok := ctx.Value("user").(identity_domain.User)
	if !ok {
		return errors.New("user not found in context")
	}
	utils.LogPretty("ShareReadyHandler - Handle - USER ", user)

	if payload.RecipientEmail != "" && user.Email == payload.RecipientEmail {
		runtime.EventsEmit(
			h.appCtx,
			msg.Type,
			payload,
		)
	}

	return nil
}

package realtime_client_handlers

import (
	context "context"

	realtime_client_domain "vault-app/internal/realtime_client/domain"
	shared_realtime "vault-app/internal/shared/realtime"
)

func RegisterHandlers(
	ctx context.Context,
) map[string]realtime_client_domain.MessageHandler {

	return map[string]realtime_client_domain.MessageHandler{
		shared_realtime.ShareInvitation: NewShareInvitationHandler(ctx),
		shared_realtime.ShareAccepted:   NewShareAcceptedHandler(ctx),
		shared_realtime.ShareRejected:   NewShareRejectedHandler(ctx),
		shared_realtime.ShareReady:      NewShareReadyHandler(ctx),
	}
}

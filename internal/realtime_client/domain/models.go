package realtime_client_domain

import (
	"context"

	shared_realtime "vault-app/internal/shared/realtime"
)

type MessageHandler interface {
	Handle(
		ctx context.Context,
		msg shared_realtime.Message,
	) error
}

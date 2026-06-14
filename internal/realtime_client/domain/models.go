package realtime_client_domain

import (
	"context"
	"time"

	shared_realtime "vault-app/internal/shared/realtime"
)

type MessageHandler interface {
	Handle(
		ctx context.Context,
		msg shared_realtime.Message,
	) error
}

type RealtimeOffset struct {
	UserID    string
	LastSeq   uint64
	UpdatedAt time.Time
}
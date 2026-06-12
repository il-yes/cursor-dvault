package realtime_client_application_events

import (
	shared_realtime "vault-app/internal/shared/realtime"
)

type EventBus interface {
	Publish(event shared_realtime.Message)
	Subscribe(handler any)
}
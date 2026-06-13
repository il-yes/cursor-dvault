package realtime_client_infrastructure

import (
	"encoding/json"

	"github.com/gorilla/websocket"

	shared_realtime "vault-app/internal/shared/realtime"
)

type AckSender struct{}

func NewAckSender() *AckSender {
	return &AckSender{}
}

func (a *AckSender) Send(
	conn *websocket.Conn,
	notificationID string,
) error {

	ack := shared_realtime.Message{
		Version: 1,
		Type:    shared_realtime.NotificationAck,
		Payload: mustMarshal(shared_realtime.NotificationAckPayload{
			NotificationID: notificationID,
		}),
	}

	return conn.WriteJSON(ack)
}

func mustMarshal(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
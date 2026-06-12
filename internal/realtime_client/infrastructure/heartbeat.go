package realtime_client_infrastructure

import (
	"time"

	"github.com/gorilla/websocket"
)


type Heartbeat struct {
	timeout time.Duration
}

func NewHeartbeat(timeout time.Duration) *Heartbeat {
	return &Heartbeat{timeout: timeout}
}

func (h *Heartbeat) Apply(conn *websocket.Conn) {
	// IMPORTANT: when pong arrives, extend deadline
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(h.timeout))
	})

	// initial deadline
	_ = conn.SetReadDeadline(time.Now().Add(h.timeout))
}
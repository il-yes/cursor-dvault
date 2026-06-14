package realtime_client_infrastructure_websocket

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"

	realtime_client_infrastructure "vault-app/internal/realtime_client/infrastructure"
	shared_realtime "vault-app/internal/shared/realtime"
	"vault-app/internal/utils"
)

type Client struct {
	url    string
	dialer *websocket.Dialer
	Ack    *realtime_client_infrastructure.AckSender
}

func NewClient(url string) *Client {
	return &Client{
		url:    url,
		dialer: websocket.DefaultDialer,
		Ack:    realtime_client_infrastructure.NewAckSender(),
	}
}

func (c *Client) Start(
	ctx context.Context,
	conn *websocket.Conn,
	onMessage func(shared_realtime.Message),
) error {

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		var msg shared_realtime.Message
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		onMessage(msg)
	}
}

func (c *Client) Run(
	ctx context.Context,
	onMessage func(shared_realtime.Message) error,
) error {
	utils.LogPretty("Attempting to connect to realtime server at: ", c.url)
	backoff := 1 * time.Second

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		conn, _, err := c.dialer.Dial(c.url, nil)
		if err != nil {
			log.Printf("WS dial failed: %v", err)
			time.Sleep(backoff)

			if backoff < 30*time.Second {
				backoff *= 2
			}

			continue
		}
		log.Printf("WS connected")

		heartbeat := realtime_client_infrastructure.NewHeartbeat(70 * time.Second)
		heartbeat.Apply(conn)

		backoff = 1 * time.Second

		readErr := c.readLoop(ctx, conn, onMessage)

		_ = conn.Close()

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if readErr != nil {
			time.Sleep(backoff)
			continue
		}
	}
}
func (c *Client) readLoop(
	ctx context.Context,
	conn *websocket.Conn,
	onMessage func(shared_realtime.Message) error,
) error {

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// IMPORTANT: extend deadline on each read
		_ = conn.SetReadDeadline(time.Now().Add(70 * time.Second))

		_, data, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		var msg shared_realtime.Message
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		log.Printf("WS message: %s", string(data))

		// Ignore ACK messages
		if msg.Type == shared_realtime.NotificationAck {
			continue
		}

		// 1. Execute business logic
		if err := onMessage(msg); err != nil {

			log.Printf(
				"WS handler failed seq=%d err=%v",
				msg.Seq,
				err,
			)

			// IMPORTANT:
			// No ACK sent.
			// Cloud will replay later.
			continue
		}

		// 2. Extract notification id
		var notif shared_realtime.NotificationPayload

		if err := json.Unmarshal(
			msg.Payload,
			&notif,
		); err != nil {

			log.Printf(
				"WS payload decode failed seq=%d err=%v",
				msg.Seq,
				err,
			)

			continue
		}

		// 3. ACK after successful processing
		if err := c.Ack.Send(
			conn,
			notif.ID,
		); err != nil {

			log.Printf(
				"WS ACK failed notification=%s err=%v",
				notif.ID,
				err,
			)
		}
	}
}

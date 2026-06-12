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
}

func NewClient(url string) *Client {
	return &Client{
		url:    url,
		dialer: websocket.DefaultDialer,
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
	onMessage func(shared_realtime.Message),
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
	onMessage func(shared_realtime.Message),
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

		onMessage(msg)
	}
}

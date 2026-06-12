package realtime_client_application_services

import (
	"context"
	"errors"
	"fmt"

	realtime_client_domain "vault-app/internal/realtime_client/domain"
	shared_realtime "vault-app/internal/shared/realtime"
)

type Transport interface {
	Start(ctx context.Context, onMessage func(msg shared_realtime.Message)) error
	Close() error
}

type Client struct {
	handlers map[string]realtime_client_domain.MessageHandler
}

func NewClient(handlers map[string]realtime_client_domain.MessageHandler) *Client {
	return &Client{handlers: handlers}
}

func (c *Client) Register(
	eventType string,
	handler realtime_client_domain.MessageHandler,
) {
	c.handlers[eventType] = handler
}

func (c *Client) Handle(
	ctx context.Context,
	msg shared_realtime.Message,
) error {
	handler, ok := c.handlers[msg.Type]
	if !ok {
		fmt.Println("NO HANDLER FOR =", msg.Type)
		return errors.New("no handler registered")
	}
	return handler.Handle(ctx, msg)
}

func (c *Client) Dispatch(
	ctx context.Context,
	msg shared_realtime.Message,
) error {

	handler, ok := c.handlers[msg.Type]
	if !ok {
		fmt.Println("NO HANDLER FOR =", msg.Type)
		return errors.New("no handler registered")
	}

	return handler.Handle(ctx, msg)
}

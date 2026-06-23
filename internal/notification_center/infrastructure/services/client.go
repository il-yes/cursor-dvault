package notification_center_infrastructure_services

import (
	"context"

	notification_center_domain "vault-app/internal/notification_center/domain"
)

type NotificationClient struct {
	client notification_center_domain.NotificationServiceInterface
}

func NewNotificationClient(
	client notification_center_domain.NotificationServiceInterface,
) *NotificationClient {
	return &NotificationClient{
		client: client,
	}
}

func (c *NotificationClient) ListByUser(
		ctx context.Context,
		userID string,
		limit int,
		offset int,
	) ([]notification_center_domain.Notification, error) {
		return c.client.ListByUser(ctx, userID, limit, offset)
	}

func (c *NotificationClient) CountUnread(
		ctx context.Context,
		userID string,
	) (int64, error) {
		return c.client.CountUnread(ctx, userID)
	}

func (c *NotificationClient) MarkRead(
		ctx context.Context,
		id string,
	) error {
		return c.client.MarkRead(ctx, id)
	}

func (c *NotificationClient) Archive(
		ctx context.Context,
		id string,
	) error {
		return c.client.Archive(ctx, id)
	}

func (c *NotificationClient) MarkAllRead(
		ctx context.Context,
		userID string,
	) error {
		return c.client.MarkAllRead(ctx, userID)
	}

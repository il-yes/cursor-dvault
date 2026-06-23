package notification_center_usecases

import (
	"context"

	notification_center_domain "vault-app/internal/notification_center/domain"
)

type NotificationUseCase struct {
	client notification_center_domain.NotificationServiceInterface
}

func NewNotificationUseCase(
	client notification_center_domain.NotificationServiceInterface,
) *NotificationUseCase {
	return &NotificationUseCase{
		client: client,
	}
}

func (uc *NotificationUseCase) ListByUser(
	ctx context.Context,
	userID string,
	limit int,
	offset int,
) ([]notification_center_domain.Notification, error) {
	return uc.client.ListByUser(ctx, userID, limit, offset)
}

func (uc *NotificationUseCase) CountUnread(
	ctx context.Context,
	userID string,
) (int64, error) {
	return uc.client.CountUnread(ctx, userID)
}

func (uc *NotificationUseCase) MarkRead(
	ctx context.Context,
	id string,
) error {
	return uc.client.MarkRead(ctx, id)
}

func (uc *NotificationUseCase) Archive(
	ctx context.Context,
	id string,
) error {
	return uc.client.Archive(ctx, id)
}

func (uc *NotificationUseCase) MarkAllRead(
	ctx context.Context,
	userID string,
) error {
	return uc.client.MarkAllRead(ctx, userID)
}
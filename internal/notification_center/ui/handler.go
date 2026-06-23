package notification_center_ui

import (
	"context"

	notification_center_usecases "vault-app/internal/notification_center/application/use_cases"
	notification_center_domain "vault-app/internal/notification_center/domain"
)

type NotificationHandler struct {
	notificationUseCase notification_center_usecases.NotificationUseCase
}

func NewNotificationHandler(
	notificationUseCase notification_center_usecases.NotificationUseCase,
) *NotificationHandler {
	return &NotificationHandler{
		notificationUseCase: notificationUseCase,
	}
}

func (h *NotificationHandler) ListByUser(
	ctx context.Context,
	userID string,
	limit int,
	offset int,
) ([]notification_center_domain.Notification, error) {
	return h.notificationUseCase.ListByUser(ctx, userID, limit, offset)
}

func (uc *NotificationHandler) CountUnread(
	ctx context.Context,
	userID string,
) (int64, error) {
	return uc.notificationUseCase.CountUnread(ctx, userID)
}

func (uc *NotificationHandler) MarkRead(
	ctx context.Context,
	id string,
) error {
	return uc.notificationUseCase.MarkRead(ctx, id)
}

func (uc *NotificationHandler) Archive(
	ctx context.Context,
	id string,
) error {
	return uc.notificationUseCase.Archive(ctx, id)
}

func (uc *NotificationHandler) MarkAllRead(
	ctx context.Context,
	userID string,
) error {
	return uc.notificationUseCase.MarkAllRead(ctx, userID)
}
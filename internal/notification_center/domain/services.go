package notification_center_domain

import "context"

type NotificationServiceInterface interface {

	ListByUser(
		ctx context.Context,
		userID string,
		limit int,
		offset int,
	) ([]Notification, error)

	CountUnread(
		ctx context.Context,
		userID string,
	) (int64, error)

	MarkRead(
		ctx context.Context,
		id string,
	) error

	Archive(
		ctx context.Context,
		id string,
	) error

	MarkAllRead(
		ctx context.Context,
		userID string,
	) error
}
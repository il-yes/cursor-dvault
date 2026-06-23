import { acceptShare, rejectShare, revokeShare } from "@/services/api";
import { parseNotificationPayload } from "@/services/utils";
import { create } from "zustand";
import { devtools } from "zustand/middleware";


export type NotificationStatus =
	| "unread"
	| "read"
	| "archived";

export type NotificationType =
	| "share.invitation"
	| "share.accepted"
	| "share.rejected"
	| "share.ready_to_accept"
	| "subscription.activated"
	| "subscription.expired"
	| string;

export type NotificationPayload = Record<string, any>;

export interface Notification {
	id: string;

	sequence: number;

	user_id: string;

	type: NotificationType;

	title: string;
	body: string;

	status: NotificationStatus;

	created_at: string;
	read_at: string | null;
	// archived_at: string | null;

	payload?: NotificationPayload;
}

interface NotificationsState {
	notifications: Notification[];
	getUnreadCount: () => number;

	isLoading: boolean;
	error: string | null;

	setNotifications: (
		items: Notification[]
	) => void;

	pushNotification: (
		item: Notification
	) => void;

	markRead: (
		id: string
	) => void;

	archive: (
		id: string
	) => void;

	markAllRead: () => void;

	removeNotification: (
		id: string
	) => void;

	clearArchived: () => void;

	acceptShare: (notification: Notification) => Promise<void>;

	rejectShare: (notification: Notification) => Promise<void>;

	revokeShare: (notification: Notification) => Promise<void>;

}

export const useNotificationsStore =
	create<NotificationsState>()(
		devtools(
			(set, get) => ({
				notifications: [],

				isLoading: false,
				error: null,

				setNotifications: (
					items
				) =>
					set(
						() => ({
							notifications: [...items].sort(
								(a, b) =>
									b.sequence - a.sequence
							),
						}),
						false,
						"notifications/setNotifications"
					),

				pushNotification: (
					item
				) =>
					set(
						(state) => {
							const exists =
								state.notifications.some(
									(n) =>
										n.sequence ===
										item.sequence
								);

							if (exists) {
								return state;
							}

							return {
								notifications: [
									item,
									...state.notifications,
								].sort(
									(a, b) =>
										b.sequence -
										a.sequence
								),
							};
						},
						false,
						"notifications/pushNotification"
					),

				markRead: (
					id
				) =>
					set(
						(state) => ({
							notifications:
								state.notifications.map(
									(n) =>
										n.id === id &&
											n.status === "unread"
											? {
												...n,
												status: "read",
												read_at:
													new Date().toISOString(),
											}
											: n
								),
						}),
						false,
						"notifications/markRead"
					),

				archive: async (notification: Notification) => {
					const snapshot = get().notifications;
					set(
						(state) => ({
							notifications:
								state.notifications.map(
									(n) =>
										n.id === notification.id
											? {
												...n,
												status:
													"archived",
												archived_at:
													new Date().toISOString(),
												read_at:
													n.read_at ??
													new Date().toISOString(),
											}
											: n
								),
						}),
						false,
						"notifications/archive"
					);

					// TODO
					// try {
					// 	const payload = parseNotificationPayload(notification.payload);
					// 	const shareId = payload?.share_id;

					// 	await archiveShare(shareId, notification.id);
					// } catch (err) {
					// 	set(
					// 		() => ({ notifications: snapshot, error: "Failed to archive notification" }),
					// 		false,
					// 		"notifications/archive:rollback"
					// 	);
					// 	throw err;
					// }

				},

				markAllRead: () =>
					set(
						(state) => ({
							notifications:
								state.notifications.map(
									(n) =>
										n.status ===
											"unread"
											? {
												...n,
												status: "read",
												read_at:
													new Date().toISOString(),
											}
											: n
								),
						}),
						false,
						"notifications/markAllRead"
					),

				removeNotification: (
					id
				) =>
					set(
						(state) => ({
							notifications:
								state.notifications.filter(
									(n) =>
										n.id !== id
								),
						}),
						false,
						"notifications/removeNotification"
					),

				clearArchived: () =>
					set(
						(state) => ({
							notifications:
								state.notifications.filter(
									(n) =>
										n.status !==
										"archived"
								),
						}),
						false,
						"notifications/clearArchived"
					),

				unreadCount: () =>
					get().notifications.filter(
						(n) =>
							n.status ===
							"unread"
					).length,

				acceptShare: async (notification: Notification) => {
					const snapshot = get().notifications;

					set(
						(state) => ({
							notifications: state.notifications.map((n) =>
								n.id === notification.id
									? {
										...n,
										status: "read",
										read_at: n.read_at ?? new Date().toISOString(),
									}
									: n
							),
						}),
						false,
						"notifications/acceptShare:optimistic"
					);

					try {
						const payload = parseNotificationPayload(notification.payload);
						const shareId = payload?.share_id;

						await acceptShare(shareId, notification.id);
					} catch (err) {
						set(
							() => ({ notifications: snapshot, error: "Failed to accept share" }),
							false,
							"notifications/acceptShare:rollback"
						);
						throw err;
					}
				},

				rejectShare: async (notification: Notification) => {
					const snapshot = get().notifications;

					set(
						(state) => ({
							notifications: state.notifications.map((n) =>
								n.id === notification.id
									? {
										...n,
										status: "read",
										read_at: n.read_at ?? new Date().toISOString(),
									}
									: n
							),
						}),
						false,
						"notifications/rejectShare:optimistic"
					);

					try {
						const payload = parseNotificationPayload(notification.payload);
						const shareId = payload?.share_id;

						if (!shareId) {
							throw new Error("Invalid notification payload for share rejection");
						}
						await rejectShare(shareId, notification.id);
					} catch (err) {
						set(
							() => ({ notifications: snapshot, error: "Failed to reject share" }),
							false,
							"notifications/rejectShare:rollback"
						);
						throw err;
					}
				},

				revokeShare: async (notificationId: string) => {
					const snapshot = get().notifications;

					set(
						(state) => ({
							notifications: state.notifications.map((n) =>
								n.id === notificationId
									? {
										...n,
										status: "archived",
										archived_at: new Date().toISOString(),
										read_at: n.read_at ?? new Date().toISOString(),
									}
									: n
							),
						}),
						false,
						"notifications/revokeShare:optimistic"
					);

					try {
						await revokeShare(notificationId);
					} catch (err) {
						set(
							() => ({ notifications: snapshot, error: "Failed to revoke share" }),
							false,
							"notifications/revokeShare:rollback"
						);
						throw err;
					}
				},
				getUnreadCount: () =>
					get().notifications.filter(
						n => n.status === "unread"
					).length
			}),
			{
				name: "notifications-store",
			}
		)
	);
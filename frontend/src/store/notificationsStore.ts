import { acceptShare, acceptChannelInvitation, rejectShare, revokeShare } from "@/services/api";
import { markRead as markReadApi, markAllRead as markAllReadApi } from "@/services/notificationsApi";
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
	| "workspace.invitation"
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
	) => Promise<void>;

	archive: (
		id: string
	) => void;

	markAllRead: () => Promise<void>;

	removeNotification: (
		id: string
	) => void;

	clearArchived: () => void;

	acceptShare: (notification: Notification) => Promise<void>;

	acceptWorkspaceInvitation: (notification: Notification) => Promise<void>;

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
						(state) => {
							const map = new Map<string, Notification>();
							for (const existing of state.notifications) {
								map.set(existing.id, existing);
							}
							for (const item of items) {
								map.set(item.id, item);
							}
							return {
								notifications: Array.from(map.values()).sort(
									(a, b) => b.sequence - a.sequence
								),
							};
						},
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
										n.id === item.id ||
										(item.sequence > 0 && n.sequence === item.sequence)
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

				markRead: async (
					id
				) => {
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
					);
					try {
						await markReadApi(id);
					} catch (err) {
						console.error("[NOTIFICATIONS_SYNC] Failed to mark read on backend:", err);
					}
				},

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

				markAllRead: async () => {
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
					);
					try {
						await markAllReadApi();
					} catch (err) {
						console.error("[NOTIFICATIONS_SYNC] Failed to mark all read on backend:", err);
					}
				},

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

				acceptWorkspaceInvitation: async (notification: Notification) => {
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
						"notifications/acceptWorkspaceInvitation:optimistic"
					);

					try {
						const payload = parseNotificationPayload(notification.payload);
						const invitationId = payload?.invitation_id || payload?.invitationId || payload?.InvitationID || payload?.intent_id || payload?.id || payload?.ID;
						console.log("invitationId", invitationId);

						if (!invitationId) {
							throw new Error("Invalid notification payload: missing invitation_id");
						}

						await acceptChannelInvitation(invitationId);

						// Re-query Cloud for accessible workspaces after invitation acceptance
						const workspaceStore = (await import("@/components/C3/infrastructure/store/useC3WorkspaceStore")).useC3WorkspaceStore;
						await workspaceStore.getState().fetchWorkspaces();

						const activeWorkspaceId = workspaceStore.getState().activeWorkspaceId;
						if (activeWorkspaceId) {
							const channelStore = (await import("@/components/C3/infrastructure/store/useC3ChannelStore")).useC3ChannelStore;
							await channelStore.getState().fetchChannels(activeWorkspaceId);
						}
					} catch (err) {
						set(
							() => ({ notifications: snapshot, error: "Failed to accept workspace invitation" }),
							false,
							"notifications/acceptWorkspaceInvitation:rollback"
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
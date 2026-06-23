import { useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Notification, useNotificationsStore } from "@/store/notificationsStore";
import { NotificationItem } from "@/components/Notification/NotificationItem";
import { DashboardLayout } from "@/components/DashboardLayout";
import { list } from "@/services/notificationsApi";
import { NotificationCenter } from "@/components/Notification/NotificationCenter";
import { acceptShare, rejectShare, revokeShare } from "@/services/api";
import { parseNotificationPayload } from "@/services/utils";
import { useAuthStore } from "@/store/useAuthStore";

export default function NotificationsPage() {
	const notifications = useNotificationsStore((s) => s.notifications);
	// const fetchNotifications = useNotificationsStore((s) => s.fetchNotifications);
	const markRead = useNotificationsStore((s) => s.markRead);
	const archive = useNotificationsStore((s) => s.archive) ;
	const markAllRead = useNotificationsStore((s) => s.markAllRead);

	const fetchNotifications = async () => {
		const notifications = await list();
		console.log("fetched notifications:", notifications);
		useNotificationsStore.getState().setNotifications(notifications);
	};

	useEffect(() => {
		fetchNotifications();
	}, []);



	return (
		<DashboardLayout>
			<div className="h-full overflow-y-auto scrollbar-glassmorphism thin-scrollbar bg-gradient-to-br from-white/50 via-white/30 to-zinc-50/20 dark:from-zinc-900/50 dark:via-zinc-900/30 dark:to-black/20 backdrop-blur-xl">
				<div className="max-w-5xl mx-auto p-8 space-y-8">
					{/* Hero Header */}
					<div className="text-center backdrop-blur-xl bg-white/40 dark:bg-zinc-900/40 rounded-3xl p-12 border border-white/30 dark:border-zinc-700/30 shadow-2xl">
						<h1 className="text-6xl font-black bg-gradient-to-r from-foreground via-primary to-amber-500/80 bg-clip-text text-transparent drop-shadow-2xl mb-4">
							Notifications
						</h1>
						<p className="text-2xl text-muted-foreground/90 max-w-2xl mx-auto leading-relaxed backdrop-blur-sm">
							Manage your notifications
						</p>
					</div>


					<div className="mx-auto max-w-3xl space-y-6 p-6">
						<div className="flex items-center justify-between">
							<div>
								<h1 className="text-2xl font-semibold">Notifications</h1>
								<p className="text-sm text-muted-foreground">
									Your recent activity and system updates.
								</p>

								<p className="text-sm text-muted-foreground">
									Total {notifications.length}
								</p>
							</div>
							<Button onClick={() => markAllRead()}>Mark all read</Button>
						</div>

						<div className="space-y-3">
							{notifications.length > 0 && notifications.map((item) => (
								<NotificationItem
									key={item.id}
									item={item}
									onMarkRead={markRead}
									onArchive={archive}
								/>
							))}
							{notifications.length === 0 && (
								<div className="text-sm text-muted-foreground">No notifications yet.</div>
							)}
						</div>
					</div>

				</div>
			</div>
		</DashboardLayout>
	);
}
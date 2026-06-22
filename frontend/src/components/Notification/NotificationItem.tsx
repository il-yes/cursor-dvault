import { Archive, CheckCheck } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { useMemo } from "react";
import { Notification, useNotificationsStore } from "@/store/notificationsStore";
import { parseNotificationPayload } from "@/services/utils";



type Props = {
	item: Notification;
	onMarkRead: (id: string) => void;
	onArchive: (id: string) => void;
};

function formatRelative(iso: string) {
	const diff = Date.now() - new Date(iso).getTime();
	const minutes = Math.floor(diff / 60000);
	if (minutes < 1) return "just now";
	if (minutes < 60) return `${minutes} min ago`;
	const hours = Math.floor(minutes / 60);
	if (hours < 24) return `${hours} h ago`;
	const days = Math.floor(hours / 24);
	return `${days} d ago`;
}
function isShareNotification(type: string) {
	// return type.startsWith("share.");
	console.log(type.startsWith("share_"))
	return type.startsWith("share_");
}
/**
 * A practical mapping is:
 * share.invitation → Accept / Reject.
 * share.ready_to_accept → Accept / Reject.
 * share.accepted → Open / Revoke depending on your product rules.
 * share.rejected → Archive only, or reopen if you support retries.
 */
export function NotificationItem({ item, onMarkRead, onArchive }: Props) {
	const isUnread = item.status === "unread";
	const timeLabel = useMemo(() => formatRelative(item.created_at), [item.created_at]);
	const acceptShare = useNotificationsStore((s) => s.acceptShare);
	const rejectShare = useNotificationsStore((s) => s.rejectShare);
	const revokeShare = useNotificationsStore((s) => s.revokeShare);
	console.log({item})

	return (
		<div
			className={cn(
				"flex items-start gap-3 rounded-xl border p-3 transition-colors",
				isUnread ? "bg-primary/5 border-primary/20" : "bg-background"
			)}
		>
			<div className={cn("mt-2 h-2 w-2 rounded-full", isUnread ? "bg-primary" : "bg-muted-foreground/40")} />
			<div className="min-w-0 flex-1">
				<div className="flex items-start justify-between gap-3">
					<div>
						<p className="font-medium leading-tight">{item.title}</p>
						<p className="mt-1 text-sm text-muted-foreground">{item.body}</p>
					</div>
					<span className="shrink-0 text-xs text-muted-foreground">{timeLabel}</span>
				</div>

				<div className="mt-3 flex flex-wrap gap-2">
					{item.status === "unread" ? (
						<Button size="sm" variant="outline" onClick={() => onMarkRead(item.id)}>
							<CheckCheck className="mr-2 h-4 w-4" />
							Mark read
						</Button>
					) : null}

					{item.status !== "archived" ? (
						<Button size="sm" variant="ghost" onClick={() => onArchive(item.id)}>
							<Archive className="mr-2 h-4 w-4" />
							Archive
						</Button>
					) : null}

					{isShareNotification(item.type) ? (
						<>
							{item.type === "share_invitation" || item.type === "share_ready_to_accept" ? (
								<>
									<Button size="sm" onClick={() => acceptShare(item)}>
										Accept
									</Button>
									<Button size="sm" variant="destructive" onClick={() => rejectShare(item)}>
										Reject
									</Button>
								</>
							) : null}

							{item.type === "share_accepted" ? (
								<Button size="sm" onClick={() => revokeShare(item)}>
									Revoke
								</Button>
							) : null}
						</>
					) : null}
				</div>
			</div>
		</div>
	);
}
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from "@/components/ui/sheet";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "@/components/ui/button";
import { useEffect } from "react";
import { NotificationBell } from "./NotificationBell";
import { NotificationItem } from "./NotificationItem";
import { useNotificationsStore } from "@/store/notificationsStore";
import { list } from "@/services/notificationsApi";
import { Notification } from "@/store/notificationsStore";

export function NotificationCenter({notifications}: {notifications: Notification[]}) {
  const isLoading = useNotificationsStore((state) => state.isLoading);
  const markRead = useNotificationsStore((state) => state.markRead);
  const archive = useNotificationsStore((state) => state.archive);
  const markAllRead = useNotificationsStore((state) => state.markAllRead);
  

  return (
    <Sheet>
      <SheetTrigger asChild>
        <span>
          <NotificationBell />
        </span>
      </SheetTrigger>

      <SheetContent className="w-full sm:max-w-md">
        <SheetHeader className="space-y-2">
          <SheetTitle>Notifications</SheetTitle>
          <div className="flex items-center justify-between">
            <p className="text-sm text-muted-foreground">
              {notifications.length} total
            </p>
            <Button variant="outline" size="sm" onClick={() => markAllRead()}>
              Mark all read
            </Button>
          </div>
        </SheetHeader>

        <div className="mt-6">
          { isLoading ? (
            <div className="text-sm text-muted-foreground">Loading…</div>
          ) : notifications.length === 0 ? (
            <div className="text-sm text-muted-foreground">No notifications yet.</div>
          ) : ( 
            <ScrollArea className="h-[calc(100vh-180px)] pr-4">
              <div className="space-y-3">
                {notifications.length > 0 && notifications.map((item) => (
                  <NotificationItem
                    key={item.id}
                    item={item}
                    onMarkRead={markRead}
                    onArchive={archive}
                  />
                ))}
              </div>
            </ScrollArea>
          )}
        </div>
      </SheetContent>
    </Sheet>
  );
}
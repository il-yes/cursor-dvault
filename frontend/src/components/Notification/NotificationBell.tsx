import React from 'react'
import { Bell } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useNotificationsStore } from '@/store/notificationsStore';

type Props = {
  onClick?: () => void;
} & React.ComponentProps<typeof Button>;


export function NotificationBell({ onClick, ...restProps }: Props) {
  const unreadCount = useNotificationsStore((state) => state.getUnreadCount());

  return (
    <Button
      variant="ghost"
      size="icon"
      onClick={onClick}
      className="relative"
      {...restProps}
    >
      <Bell className="h-5 w-5" />

      {unreadCount > 0 && (
        <Badge
          className="
            absolute
            -right-1
            -top-1
            h-5
            min-w-5
            rounded-full
            px-1
            text-[10px]
          "
        >
          {unreadCount > 99 ? "99+" : unreadCount}
        </Badge>
      )}

    </Button>
  );
}

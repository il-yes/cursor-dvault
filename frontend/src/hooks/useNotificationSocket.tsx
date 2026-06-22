import { Notification, useNotificationsStore } from "@/store/notificationsStore";
import { EVENTS, ShareAcceptedPayload, ShareInvitationNotificationPayload, ShareReadyToAcceptPayload, ShareRejectedPayload } from "@/types/sharing";
import { useEffect, useRef } from "react";
import { toast } from "./use-toast";
import { useAuthStore } from "@/store/useAuthStore";
import { User } from "@/types/vault";

type SocketEnvelope =
  | {
    type: "notification";
    notification: Notification;
  }
  | {
    type: "share.invitation";
    notification: Notification;
  }
  | {
    type: "share.accepted";
    notification: Notification;
  }
  | {
    type: "share.rejected";
    notification: Notification;
  }
  | {
    type: "share.ready_to_accept";
    notification: Notification;
  }
  | {
    type: string;
    [key: string]: any;
  };

type UseNotificationSocketOptions = {
  enabled?: boolean;
  socketUrl?: string;
  onError?: (error: unknown) => void;
};

function normalizeNotification(input: any): Notification | null {
  if (!input) return null;
  if (!input.id || !input.type || !input.title || !input.body) return null;

  return {
    id: String(input.id),
    user_id: String(input.user_id ?? ""),
    type: String(input.type),
    title: String(input.title),
    body: String(input.body),
    status: (input.status ?? "unread") as Notification["status"],
    created_at: String(input.created_at ?? new Date().toISOString()),
    read_at: input.read_at ?? null,
    // archived_at: input.archived_at ?? null,
    payload: input.payload ?? {},
    sequence: input.sequence ?? 0,
  };
}



// Filtered coming messages for the current session user from the backend wails
// Triggered Actions:
//    - Update pending intent share status
// active notification ui:
//    - popup (done)
//    - notification tabs menu
export function useNotificationsEvents() {
  const pushNotification = useNotificationsStore((s) => s.pushNotification);
  const currentUser = useAuthStore((s) => s.user);

  const handleShareInvitation = (payload: ShareInvitationNotificationPayload) => {

    pushNotification({
      id: payload.id,
      user_id: currentUser.id,  
      type: EVENTS.SHARE_INVITATION,
      title: payload.title,
      body: payload.body,

      status: "unread",

      created_at: payload.created_at,
      read_at: null,
      sequence: payload.seq,
      payload: payload,
    });

    toast({
      title: payload.title,
      description: payload.body,
    });

  };

  const handleShareAccepted = (payload: ShareAcceptedPayload) => {

    pushNotification({
      id: payload.id,
      user_id: currentUser.id,  
      type: EVENTS.SHARE_ACCEPTED,
      title: payload.title,
      body: payload.body,

      status: "unread",

      created_at: payload.created_at,
      read_at: null,
      sequence: payload.seq,
      payload: payload,
    });

    toast({
      title: payload.title,
      description: payload.body,
    });
  };

  const handleShareRejected = (payload: ShareRejectedPayload) => {

    pushNotification({
      id: payload.id,
      user_id: currentUser.id,  
      type: EVENTS.SHARE_REJECTED,
      title: payload.title,
      body: payload.body,

      status: "unread",

      created_at: payload.created_at,
      read_at: null,
      sequence: payload.seq,
      payload: payload,
    });

    toast({
      title: payload.title,
      description: payload.body,
      variant: "destructive",
    });

  };

  const handleShareReadyToAccept = (payload: ShareReadyToAcceptPayload) => {

    pushNotification({
      id: payload.id,
      user_id: currentUser.id,  
      type: EVENTS.SHARE_READY,
      title: payload.title,
      body: payload.body,

      status: "unread",

      created_at: payload.created_at,
      read_at: null,
      sequence: payload.seq,
      payload: payload,
    });

    toast({
      title: payload.title,
      description: payload.body,
    });

  };

  useEffect(() => {
    const unsubInvitation = window.runtime?.EventsOn(
      EVENTS.SHARE_INVITATION,
      handleShareInvitation
    );

    const unsubAccepted = window.runtime?.EventsOn(
      EVENTS.SHARE_ACCEPTED,
      handleShareAccepted
    );

    const unsubRejected = window.runtime?.EventsOn(
      EVENTS.SHARE_REJECTED,
      handleShareRejected
    );

    const unsubReady = window.runtime?.EventsOn(
      EVENTS.SHARE_READY,
      handleShareReadyToAccept
    );

    return () => {
      unsubInvitation?.();
      unsubAccepted?.();
      unsubRejected?.();
      unsubReady?.();
    };
  }, [pushNotification, currentUser]);
}
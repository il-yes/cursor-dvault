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
      user_id: currentUser?.id ?? "",  
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
      user_id: currentUser?.id ?? "",  
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
      user_id: currentUser?.id ?? "",  
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
      user_id: currentUser?.id ?? "",  
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

  const handleChannelInvitationCreated = (data: any) => {
    console.log("[C3][NOTIFICATION_TRACE][EVENT_IN] eventType=channel.invitation.created data=", data);

    const invId = data.InvitationID || data.invitation_id || data.ID || data.id;
    const channelId = data.ChannelID || data.channel_id;
    const inviterVaultId = data.InviterVaultID || data.inviter_vault_id;
    const inviteeVaultId = data.InviteeVaultID || data.invitee_vault_id;
    const occurredAt = data.OccurredAt || data.occurred_at || data.CreatedAt || new Date().toISOString();

    if (!invId) {
      console.warn("[C3][NOTIFICATION_TRACE] missing invitation id, skipping notification push");
      return;
    }

    const notificationId = `inv_${invId}`;
    const seq = data.seq || data.sequence || Date.now();

    const notification: Notification = {
      id: notificationId,
      user_id: currentUser?.id ?? "",
      type: "workspace.invitation",
      title: "Channel Invitation",
      body: `You have been invited to join channel ${channelId || invId}`,
      status: "unread",
      created_at: occurredAt,
      read_at: null,
      sequence: seq,
      payload: {
        invitation_id: invId,
        channel_id: channelId,
        inviter_vault_id: inviterVaultId,
        invitee_vault_id: inviteeVaultId,
        ...data,
      },
    };

    console.log("[C3][NOTIFICATION_TRACE][PUSH] notificationId=", notificationId, "invitationID=", invId);
    pushNotification(notification);

    toast({
      title: notification.title,
      description: notification.body,
    });

    console.log("[C3][NOTIFICATION_TRACE][COMPLETE] invitationID=", invId, "notificationCreated=true");
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

    const unsubChannelInv = window.runtime?.EventsOn(
      "channel.invitation.created",
      handleChannelInvitationCreated
    );

    return () => {
      unsubInvitation?.();
      unsubAccepted?.();
      unsubRejected?.();
      unsubReady?.();
      unsubChannelInv?.();
    };
  }, [pushNotification, currentUser]);
}
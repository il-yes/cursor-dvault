
import * as AppAPI from "../../wailsjs/go/main/App";
import { useAuthStore } from "@/store/useAuthStore";
import { Notification } from "@/store/notificationsStore";

export const list = async (): Promise<Notification[]> => {
    const { jwtToken } = useAuthStore.getState();
    const limit = 10
    const offset = 0
    const searchTerm = ""
    const sortBy = "createdAt"
    const sortOrder = "desc"
    const rawNotifications = await AppAPI.ListByUser(jwtToken, limit, offset);
    return rawNotifications.map((n) => ({
        id: String(n.id),
        user_id: String(n.user_id ?? ""),
        type: String(n.type),
        title: String(n.title),
        body: String(n.body),
        status: (n.status ?? "unread") as Notification["status"],
        created_at: String(n.created_at ?? new Date().toISOString()),
        read_at: n.read_at ?? null,
        payload: n.payload ?? {},
        sequence: n.sequence ?? 0,
    }));
}
export const markRead = async(id: string) => {
    const { jwtToken } = useAuthStore.getState();
    return await AppAPI.MarkRead(jwtToken, id)
}
export const archive = async(id: string) => {
    const { jwtToken } = useAuthStore.getState();
    return await AppAPI.Archive(jwtToken, id)
}
export const markAllRead = async() => {
    const { jwtToken } = useAuthStore.getState();
    return await AppAPI.MarkAllRead(jwtToken)
}
export const countUnread = async() => {
    const { jwtToken } = useAuthStore.getState();
    return await AppAPI.CountUnread(jwtToken)
}

// export const clearArchived = async() => {
//     const { jwtToken } = useAuthStore.getState();
//     return await AppAPI.ClearArchived(jwtToken)
// }

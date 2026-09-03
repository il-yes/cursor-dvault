import { describe, it, expect, vi, beforeEach } from "vitest";
import { useNotificationsStore, Notification } from "./notificationsStore";
import * as api from "@/services/api";

vi.mock("@/services/api", () => ({
	acceptShare: vi.fn(),
	acceptChannelInvitation: vi.fn(),
	rejectShare: vi.fn(),
	revokeShare: vi.fn(),
}));

vi.mock("@/components/C3/infrastructure/store/useC3WorkspaceStore", () => ({
	useC3WorkspaceStore: {
		getState: () => ({
			fetchWorkspaces: vi.fn().mockResolvedValue(undefined),
			activeWorkspaceId: "ws_1",
		}),
	},
}));

vi.mock("@/components/C3/infrastructure/store/useC3ChannelStore", () => ({
	useC3ChannelStore: {
		getState: () => ({
			fetchChannels: vi.fn().mockResolvedValue(undefined),
		}),
	},
}));

describe("notificationsStore - acceptWorkspaceInvitation", () => {
	beforeEach(() => {
		useNotificationsStore.setState({
			notifications: [],
			isLoading: false,
			error: null,
		});
		vi.clearAllMocks();
	});

	it("optimistically updates notification to read and calls acceptChannelInvitation with extracted invitation_id", async () => {
		const notification: Notification = {
			id: "notif-001",
			sequence: 1,
			user_id: "user-123",
			type: "workspace.invitation",
			title: "Workspace Invitation",
			body: "You have been invited to join Workspace",
			status: "unread",
			created_at: new Date().toISOString(),
			read_at: null,
			payload: { invitation_id: "inv_9e24294d" },
		};

		useNotificationsStore.setState({ notifications: [notification] });

		vi.mocked(api.acceptChannelInvitation).mockResolvedValueOnce({
			id: "inv_9e24294d",
			channel_id: "ch_1",
			inviter_vault_id: "v_1",
			invitee_vault_id: "v_2",
			status: "accepted",
			created_at: new Date().toISOString(),
		});

		await useNotificationsStore.getState().acceptWorkspaceInvitation(notification);

		expect(api.acceptChannelInvitation).toHaveBeenCalledWith("inv_9e24294d");

		const updated = useNotificationsStore.getState().notifications[0];
		expect(updated.status).toBe("read");
		expect(updated.read_at).not.toBeNull();
	});

	it("rolls back state if acceptChannelInvitation fails", async () => {
		const notification: Notification = {
			id: "notif-002",
			sequence: 2,
			user_id: "user-123",
			type: "workspace.invitation",
			title: "Workspace Invitation",
			body: "Invitation text",
			status: "unread",
			created_at: new Date().toISOString(),
			read_at: null,
			payload: { invitation_id: "inv_failed" },
		};

		useNotificationsStore.setState({ notifications: [notification] });

		vi.mocked(api.acceptChannelInvitation).mockRejectedValueOnce(
			new Error("Cloud backend returned status 400: invitation not for you")
		);

		await expect(
			useNotificationsStore.getState().acceptWorkspaceInvitation(notification)
		).rejects.toThrow("Cloud backend returned status 400: invitation not for you");

		const restored = useNotificationsStore.getState().notifications[0];
		expect(restored.status).toBe("unread");
		expect(restored.read_at).toBeNull();
		expect(useNotificationsStore.getState().error).toBe("Failed to accept workspace invitation");
	});

	it("deduplicates notifications when same item id is pushed twice", () => {
		const notification1: Notification = {
			id: "inv_97090042",
			sequence: 100,
			user_id: "user-123",
			type: "workspace.invitation",
			title: "Channel Invitation",
			body: "You have been invited to join channel c9fdacfb",
			status: "unread",
			created_at: new Date().toISOString(),
			read_at: null,
			payload: { invitation_id: "97090042", channel_id: "c9fdacfb" },
		};

		useNotificationsStore.getState().pushNotification(notification1);
		expect(useNotificationsStore.getState().notifications.length).toBe(1);

		// Replay exact same notification
		useNotificationsStore.getState().pushNotification(notification1);
		expect(useNotificationsStore.getState().notifications.length).toBe(1);
	});

	it("correctly adapts real channel.invitation.created payload and passes exact invitation_id on accept click", async () => {
		const rawRealtimePayload = {
			InvitationID: "real_inv_uuid_777888",
			ChannelID: "c9fdacfb-d3da-45a4-9c94-c18923dfd781",
			WorkspaceID: "ws_legal_deal_room_99",
			InviterVaultID: "v_alice_123",
			InviteeVaultID: "v_bob_456",
			OccurredAt: "2026-09-02T11:17:47.520519-07:00",
		};

		const invId = rawRealtimePayload.InvitationID;
		const notification: Notification = {
			id: `inv_${invId}`,
			user_id: "user_bob_123",
			type: "workspace.invitation",
			title: "Channel Invitation",
			body: `You have been invited to join channel ${rawRealtimePayload.ChannelID}`,
			status: "unread",
			created_at: rawRealtimePayload.OccurredAt,
			read_at: null,
			sequence: 1001,
			payload: {
				invitation_id: invId,
				channel_id: rawRealtimePayload.ChannelID,
				inviter_vault_id: rawRealtimePayload.InviterVaultID,
				invitee_vault_id: rawRealtimePayload.InviteeVaultID,
				...rawRealtimePayload,
			},
		};

		// 1. Push notification
		useNotificationsStore.getState().pushNotification(notification);

		const items = useNotificationsStore.getState().notifications;
		expect(items.length).toBe(1);
		expect(items[0].id).toBe("inv_real_inv_uuid_777888");
		expect(items[0].type).toBe("workspace.invitation");
		expect(items[0].status).toBe("unread");
		expect(items[0].payload?.invitation_id).toBe("real_inv_uuid_777888");

		// 2. Replay duplicate event delivery
		useNotificationsStore.getState().pushNotification(notification);
		expect(useNotificationsStore.getState().notifications.length).toBe(1);

		// 3. User clicks Accept button
		vi.mocked(api.acceptChannelInvitation).mockResolvedValueOnce({
			id: invId,
			channel_id: rawRealtimePayload.ChannelID,
			inviter_vault_id: rawRealtimePayload.InviterVaultID,
			invitee_vault_id: rawRealtimePayload.InviteeVaultID,
			status: "accepted",
			created_at: new Date().toISOString(),
		});

		await useNotificationsStore.getState().acceptWorkspaceInvitation(items[0]);
		expect(api.acceptChannelInvitation).toHaveBeenCalledWith("real_inv_uuid_777888");
		expect(useNotificationsStore.getState().notifications[0].status).toBe("read");
	});
});

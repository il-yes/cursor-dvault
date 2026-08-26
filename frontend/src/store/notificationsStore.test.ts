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
});

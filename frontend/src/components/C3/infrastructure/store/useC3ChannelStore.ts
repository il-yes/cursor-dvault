import { create } from "zustand";
import { listChannels, ChannelResponse } from "@/services/api";

export interface C3ChannelState {
	channels: ChannelResponse[];
	activeChannel: ChannelResponse | null;
	activeChannelId: string | null;
	activeWorkspaceId: string | null;
	isLoading: boolean;
	error: string | null;

	fetchChannels: (workspaceId: string) => Promise<void>;
	selectChannel: (channelId: string) => void;
	addChannel: (channel: ChannelResponse) => void;
	clearChannels: () => void;
}

export const useC3ChannelStore = create<C3ChannelState>((set, get) => ({
	channels: [],
	activeChannel: null,
	activeChannelId: null,
	activeWorkspaceId: null,
	isLoading: false,
	error: null,

	fetchChannels: async (workspaceId: string) => {
		if (!workspaceId) {
			set({
				channels: [],
				activeChannel: null,
				activeChannelId: null,
				activeWorkspaceId: null,
				isLoading: false,
				error: null,
			});
			return;
		}

		// Workspace -> Channel Invariant: Clear state on workspace change
		set({
			activeWorkspaceId: workspaceId,
			channels: [],
			activeChannel: null,
			activeChannelId: null,
			isLoading: true,
			error: null,
		});

		try {
			const fetched = await listChannels(workspaceId);
			if (get().activeWorkspaceId !== workspaceId) return;

			let active = fetched.length > 0 ? fetched[0] : null;
			set({
				channels: fetched,
				activeChannel: active,
				activeChannelId: active ? active.id : null,
				isLoading: false,
				error: null,
			});
		} catch (err: any) {
			if (get().activeWorkspaceId !== workspaceId) return;
			console.error(`Failed to fetch channels for workspace ${workspaceId}:`, err);
			set({
				channels: [],
				isLoading: false,
				error: err?.message || "Failed to load channels.",
			});
		}
	},

	selectChannel: (channelId: string) => {
		const found = get().channels.find((c) => c.id === channelId) || null;
		set({
			activeChannel: found,
			activeChannelId: found ? found.id : channelId,
		});
	},

	addChannel: (newChannel: ChannelResponse) => {
		const currentWorkspaceId = get().activeWorkspaceId;
		if (currentWorkspaceId && newChannel.workspace_id === currentWorkspaceId) {
			set((state) => ({
				channels: [newChannel, ...state.channels],
				activeChannel: newChannel,
				activeChannelId: newChannel.id,
			}));
		}
	},

	clearChannels: () => {
		set({
			channels: [],
			activeChannel: null,
			activeChannelId: null,
			activeWorkspaceId: null,
			isLoading: false,
			error: null,
		});
	},
}));

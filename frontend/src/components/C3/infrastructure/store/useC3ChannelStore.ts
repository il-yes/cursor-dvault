import { create } from "zustand";
import { listChannels, ChannelResponse } from "@/services/api";

export interface C3ChannelState {
	channels: ChannelResponse[];
	activeChannelId: string | null;
	activeChannel: ChannelResponse | null;
	activeWorkspaceId: string | null;
	isLoading: boolean;
	error: string | null;

	fetchChannels: (workspaceId: string) => Promise<void>;
	selectChannel: (channelId: string) => void;
	addChannel: (channel: ChannelResponse) => void;
	clearChannels: () => void;
	clearError: () => void;
}

export const useC3ChannelStore = create<C3ChannelState>((set, get) => ({
	channels: [],
	activeChannelId: null,
	activeChannel: null,
	activeWorkspaceId: null,
	isLoading: false,
	error: null,

	fetchChannels: async (workspaceId: string) => {
		if (!workspaceId) {
			set({
				channels: [],
				activeChannelId: null,
				activeChannel: null,
				activeWorkspaceId: null,
				isLoading: false,
				error: null,
			});
			return;
		}

		const previousWorkspaceId = get().activeWorkspaceId;
		const isWorkspaceChanged = previousWorkspaceId !== workspaceId;

		// Workspace isolation: reset active channel context if workspace changed
		if (isWorkspaceChanged) {
			set({
				activeWorkspaceId: workspaceId,
				activeChannelId: null,
				activeChannel: null,
				channels: [],
			});
		}

		set({ isLoading: true, error: null, activeWorkspaceId: workspaceId });

		try {
			const data = await listChannels(workspaceId);
			const currentActiveChannelId = get().activeChannelId;

			let nextActiveChannel: ChannelResponse | null = null;
			if (currentActiveChannelId && data.some((c) => c.id === currentActiveChannelId)) {
				nextActiveChannel = data.find((c) => c.id === currentActiveChannelId) || null;
			}

			set({
				channels: data,
				activeChannel: nextActiveChannel,
				activeChannelId: nextActiveChannel ? nextActiveChannel.id : null,
				isLoading: false,
				error: null,
			});
		} catch (err: any) {
			console.error(`Failed to fetch C3 channels for workspace ${workspaceId}:`, err);
			set({
				isLoading: false,
				error: err?.message || "Failed to load channels from cloud backend.",
			});
		}
	},

	selectChannel: (channelId: string) => {
		const { channels } = get();
		const target = channels.find((c) => c.id === channelId);
		if (target) {
			set({
				activeChannelId: target.id,
				activeChannel: target,
			});
		}
	},

	addChannel: (channel: ChannelResponse) => {
		set((state) => {
			const existingIndex = state.channels.findIndex((c) => c.id === channel.id);
			let updatedChannels: ChannelResponse[];

			if (existingIndex >= 0) {
				updatedChannels = [...state.channels];
				updatedChannels[existingIndex] = channel;
			} else {
				updatedChannels = [channel, ...state.channels];
			}

			return {
				channels: updatedChannels,
				activeChannelId: channel.id,
				activeChannel: channel,
			};
		});
	},

	clearChannels: () =>
		set({
			channels: [],
			activeChannelId: null,
			activeChannel: null,
			activeWorkspaceId: null,
			isLoading: false,
			error: null,
		}),

	clearError: () => set({ error: null }),
}));

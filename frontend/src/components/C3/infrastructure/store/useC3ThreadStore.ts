import { create } from "zustand";
import { listThreads, ThreadResponse } from "@/services/api";

export interface C3ThreadState {
	threads: ThreadResponse[];
	activeThread: ThreadResponse | null;
	activeThreadId: string | null;
	activeChannelId: string | null;
	isLoading: boolean;
	error: string | null;

	fetchThreads: (channelId: string) => Promise<void>;
	selectThread: (threadId: string) => void;
	addThread: (thread: ThreadResponse) => void;
	clearThreads: () => void;
}

export const useC3ThreadStore = create<C3ThreadState>((set, get) => ({
	threads: [],
	activeThread: null,
	activeThreadId: null,
	activeChannelId: null,
	isLoading: false,
	error: null,

	fetchThreads: async (channelId: string) => {
		console.log("[THREAD LIST] ThreadStore.fetchThreads channelId =", channelId);
		if (!channelId) {
			set({
				threads: [],
				activeThread: null,
				activeThreadId: null,
				activeChannelId: null,
				isLoading: false,
				error: null,
			});
			return;
		}

		// Channel -> Thread Invariant: Clear state on channel change
		set({
			activeChannelId: channelId,
			threads: [],
			activeThread: null,
			activeThreadId: null,
			isLoading: true,
			error: null,
		});

		try {
			const fetched = await listThreads(channelId);
			if (get().activeChannelId !== channelId) return;

			let active = fetched.length > 0 ? fetched[0] : null;
			set({
				threads: fetched,
				activeThread: active,
				activeThreadId: active ? active.id : null,
				isLoading: false,
				error: null,
			});
		} catch (err: any) {
			if (get().activeChannelId !== channelId) return;
			console.error(`Failed to fetch threads for channel ${channelId}:`, err);
			set({
				threads: [],
				isLoading: false,
				error: err?.message || "Failed to load threads.",
			});
		}
	},

	selectThread: (threadId: string) => {
		const found = get().threads.find((t) => t.id === threadId) || null;
		set({
			activeThread: found,
			activeThreadId: found ? found.id : threadId,
		});
	},

	addThread: (newThread: ThreadResponse) => {
		const currentChannelId = get().activeChannelId;
		const isMatch = !!currentChannelId && newThread.channel_id === currentChannelId;
		console.log("[BOUNDARY_LOG] inside addThread():", {
			currentChannelId,
			newThread_channel_id: newThread.channel_id,
			isMatch,
			newThread,
		});
		if (isMatch) {
			set((state) => ({
				threads: [newThread, ...state.threads],
				activeThread: newThread,
				activeThreadId: newThread.id,
			}));
		} else {
			console.warn("[BOUNDARY_LOG] addThread() SKIPPED set state because channel IDs did not match or currentChannelId was null!", {
				currentChannelId,
				newThread_channel_id: newThread.channel_id,
			});
		}
	},

	clearThreads: () => {
		set({
			threads: [],
			activeThread: null,
			activeThreadId: null,
			activeChannelId: null,
			isLoading: false,
			error: null,
		});
	},
}));

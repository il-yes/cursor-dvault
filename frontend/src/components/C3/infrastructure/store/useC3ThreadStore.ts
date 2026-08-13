import { create } from "zustand";
import { listThreads, ThreadResponse } from "@/services/api";

export interface C3ThreadState {
	threads: ThreadResponse[];
	activeThreadId: string | null;
	activeThread: ThreadResponse | null;
	activeChannelId: string | null;
	isLoading: boolean;
	error: string | null;

	fetchThreads: (channelId: string) => Promise<void>;
	selectThread: (threadId: string | null) => void;
	addThread: (thread: ThreadResponse) => void;
	clearThreads: () => void;
}

export const useC3ThreadStore = create<C3ThreadState>((set, get) => ({
	threads: [],
	activeThreadId: null,
	activeThread: null,
	activeChannelId: null,
	isLoading: false,
	error: null,

	fetchThreads: async (channelId: string) => {
		if (!channelId) {
			set({
				threads: [],
				activeThreadId: null,
				activeThread: null,
				activeChannelId: null,
				isLoading: false,
				error: null,
			});
			return;
		}

		// Enforce Invariant: Clear previous channel thread state when channel context changes
		set({
			activeChannelId: channelId,
			threads: [],
			activeThreadId: null,
			activeThread: null,
			isLoading: true,
			error: null,
		});

		try {
			const fetchedThreads = await listThreads(channelId);

			// Preserve channel isolation check
			if (get().activeChannelId !== channelId) {
				return;
			}

			const active = fetchedThreads.length > 0 ? fetchedThreads[0] : null;

			set({
				threads: fetchedThreads,
				activeThreadId: active ? active.id : null,
				activeThread: active,
				isLoading: false,
				error: null,
			});
		} catch (err: any) {
			if (get().activeChannelId !== channelId) {
				return;
			}

			console.error(`Failed to fetch threads for channel ${channelId}:`, err);
			set({
				threads: [],
				activeThreadId: null,
				activeThread: null,
				isLoading: false,
				error: err?.message || "Failed to load threads for the selected channel.",
			});
		}
	},

	selectThread: (threadId: string | null) => {
		if (!threadId) {
			set({ activeThreadId: null, activeThread: null });
			return;
		}

		const match = get().threads.find((t) => t.id === threadId) || null;
		set({
			activeThreadId: threadId,
			activeThread: match,
		});
	},

	addThread: (newThread: ThreadResponse) => {
		const currentChannelId = get().activeChannelId;

		// Only append thread if it matches active channel isolation context
		if (currentChannelId && newThread.channel_id === currentChannelId) {
			set((state) => ({
				threads: [newThread, ...state.threads],
				activeThreadId: newThread.id,
				activeThread: newThread,
			}));
		}
	},

	clearThreads: () => {
		set({
			threads: [],
			activeThreadId: null,
			activeThread: null,
			activeChannelId: null,
			isLoading: false,
			error: null,
		});
	},
}));

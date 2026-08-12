import { create } from "zustand";
import { listThreadEvents, ThreadEventResponse } from "@/services/api";

export interface C3ThreadEventState {
	events: ThreadEventResponse[];
	activeThreadId: string | null;
	isLoading: boolean;
	isAppending: boolean;
	error: string | null;

	fetchEvents: (threadId: string) => Promise<void>;
	addEvent: (event: ThreadEventResponse) => void;
	clearEvents: () => void;
}

export const useC3ThreadEventStore = create<C3ThreadEventState>((set, get) => ({
	events: [],
	activeThreadId: null,
	isLoading: false,
	isAppending: false,
	error: null,

	fetchEvents: async (threadId: string) => {
		if (!threadId) {
			set({
				events: [],
				activeThreadId: null,
				isLoading: false,
				error: null,
			});
			return;
		}

		// Enforce Thread Isolation Invariant: Clear previous events immediately
		set({
			activeThreadId: threadId,
			events: [],
			isLoading: true,
			error: null,
		});

		try {
			const fetchedEvents = await listThreadEvents(threadId);

			// Preserve thread isolation check
			if (get().activeThreadId !== threadId) {
				return;
			}

			// Ensure cursor ordering (ascending by cursor or created_at)
			const sorted = [...fetchedEvents].sort((a, b) => {
				if (a.cursor && b.cursor) return a.cursor - b.cursor;
				return new Date(a.created_at || 0).getTime() - new Date(b.created_at || 0).getTime();
			});

			set({
				events: sorted,
				isLoading: false,
				error: null,
			});
		} catch (err: any) {
			if (get().activeThreadId !== threadId) {
				return;
			}

			console.error(`Failed to fetch events for thread ${threadId}:`, err);
			set({
				events: [],
				isLoading: false,
				error: err?.message || "Failed to load thread events.",
			});
		}
	},

	addEvent: (newEvent: ThreadEventResponse) => {
		const currentThreadId = get().activeThreadId;

		// Only append event if it matches active thread isolation context
		if (currentThreadId && newEvent.thread_id === currentThreadId) {
			set((state) => {
				const updated = [...state.events, newEvent].sort((a, b) => {
					if (a.cursor && b.cursor) return a.cursor - b.cursor;
					return new Date(a.created_at || 0).getTime() - new Date(b.created_at || 0).getTime();
				});
				return { events: updated };
			});
		}
	},

	clearEvents: () => {
		set({
			events: [],
			activeThreadId: null,
			isLoading: false,
			isAppending: false,
			error: null,
		});
	},
}));

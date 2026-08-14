import { create } from 'zustand';
import { ShareEntryRefResponse, createCollaborativeShare, CreateCollaborativeSharePayload } from '@/services/api';

interface C3CollaborationState {
	activeShareEntryRef: ShareEntryRefResponse | null;
	isLoading: boolean;
	error: string | null;

	createShare: (payload: CreateCollaborativeSharePayload) => Promise<ShareEntryRefResponse>;
	clearCollaborationState: () => void;
}

export const useC3CollaborationStore = create<C3CollaborationState>((set) => ({
	activeShareEntryRef: null,
	isLoading: false,
	error: null,

	createShare: async (payload: CreateCollaborativeSharePayload) => {
		set({ isLoading: true, error: null });
		try {
			const shareRef = await createCollaborativeShare(payload);
			set({ activeShareEntryRef: shareRef, isLoading: false });
			return shareRef;
		} catch (err: any) {
			const errMsg = err?.message || 'Failed to create collaborative share';
			set({ error: errMsg, isLoading: false });
			throw err;
		}
	},

	clearCollaborationState: () => {
		set({
			activeShareEntryRef: null,
			isLoading: false,
			error: null,
		});
	},
}));

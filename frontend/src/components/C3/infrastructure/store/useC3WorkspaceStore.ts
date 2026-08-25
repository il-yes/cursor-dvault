import { create } from "zustand";
import { listWorkspaces, WorkspaceResponse } from "@/services/api";
import { useAppStore } from "@/store/appStore";

export interface C3WorkspaceState {
	workspaces: WorkspaceResponse[];
	activeWorkspace: WorkspaceResponse | null;
	activeWorkspaceId: string | null;
	isLoading: boolean;
	error: string | null;
	isCloudAuthRequired: boolean;

	fetchWorkspaces: () => Promise<void>;
	selectWorkspace: (workspaceId: string) => void;
	addWorkspace: (workspace: WorkspaceResponse) => void;
	clearCloudAuthRequired: () => void;
}

export const useC3WorkspaceStore = create<C3WorkspaceState>((set, get) => ({
	workspaces: [],
	activeWorkspace: null,
	activeWorkspaceId: null,
	isLoading: false,
	error: null,
	isCloudAuthRequired: false,

	fetchWorkspaces: async () => {
		const session = useAppStore.getState().session;
		

		set({ isLoading: true, error: null });
		try {
			const fetched = await listWorkspaces();
			const currentActiveId = get().activeWorkspaceId;
			let active = fetched.find((w) => w.id === currentActiveId) || null;

			if (!active && fetched.length > 0) {
				active = fetched[0];
			}

			set({
				workspaces: fetched,
				activeWorkspace: active,
				activeWorkspaceId: active ? active.id : null,
				isLoading: false,
				error: null,
				isCloudAuthRequired: false,
			});
		} catch (err: any) {
			const msg = err?.message || String(err) || "Failed to load workspaces.";
			if (msg.includes("cloud authentication required")) {
				console.info("ℹ️ C3 Workspace Store: Cloud authentication is required to access remote workspaces.");
				set({
					isLoading: false,
					error: null,
					isCloudAuthRequired: true,
				});
				return;
			}
			console.error("Failed to fetch workspaces:", err);
			set({
				isLoading: false,
				error: msg,
				isCloudAuthRequired: false,
			});
		}
	},

	selectWorkspace: (workspaceId: string) => {
		const found = get().workspaces.find((w) => w.id === workspaceId) || null;
		set({
			activeWorkspace: found,
			activeWorkspaceId: found ? found.id : workspaceId,
		});
	},

	addWorkspace: (newWorkspace: WorkspaceResponse) => {
		set((state) => ({
			workspaces: [newWorkspace, ...state.workspaces],
			activeWorkspace: newWorkspace,
			activeWorkspaceId: newWorkspace.id,
		}));
	},

	clearCloudAuthRequired: () => set({ isCloudAuthRequired: false }),
}));

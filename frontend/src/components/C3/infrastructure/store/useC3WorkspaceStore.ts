import { create } from "zustand";
import { listWorkspaces, WorkspaceResponse } from "@/services/api";

export interface C3WorkspaceState {
	workspaces: WorkspaceResponse[];
	activeWorkspaceId: string | null;
	activeWorkspace: WorkspaceResponse | null;
	isLoading: boolean;
	error: string | null;

	fetchWorkspaces: () => Promise<void>;
	selectWorkspace: (workspaceId: string) => void;
	addWorkspace: (workspace: WorkspaceResponse) => void;
	clearError: () => void;
}

export const useC3WorkspaceStore = create<C3WorkspaceState>((set, get) => ({
	workspaces: [],
	activeWorkspaceId: null,
	activeWorkspace: null,
	isLoading: false,
	error: null,

	fetchWorkspaces: async () => {
		set({ isLoading: true, error: null });
		try {
			const data = await listWorkspaces();
			const currentActiveId = get().activeWorkspaceId;

			let nextActive: WorkspaceResponse | null = null;

			if (currentActiveId && data.some((w) => w.id === currentActiveId)) {
				nextActive = data.find((w) => w.id === currentActiveId) || null;
			} else if (data.length > 0) {
				nextActive = data[0];
			}

			set({
				workspaces: data,
				activeWorkspace: nextActive,
				activeWorkspaceId: nextActive ? nextActive.id : null,
				isLoading: false,
				error: null,
			});
		} catch (err: any) {
			console.error("Failed to fetch C3 workspaces:", err);
			set({
				isLoading: false,
				error: err?.message || "Failed to load workspaces from cloud backend.",
			});
		}
	},

	selectWorkspace: (workspaceId: string) => {
		const { workspaces } = get();
		const target = workspaces.find((w) => w.id === workspaceId);
		if (target) {
			set({
				activeWorkspaceId: target.id,
				activeWorkspace: target,
			});
		}
	},

	addWorkspace: (workspace: WorkspaceResponse) => {
		set((state) => {
			const existingIndex = state.workspaces.findIndex((w) => w.id === workspace.id);
			let updatedWorkspaces: WorkspaceResponse[];

			if (existingIndex >= 0) {
				updatedWorkspaces = [...state.workspaces];
				updatedWorkspaces[existingIndex] = workspace;
			} else {
				updatedWorkspaces = [workspace, ...state.workspaces];
			}

			return {
				workspaces: updatedWorkspaces,
				activeWorkspaceId: workspace.id,
				activeWorkspace: workspace,
			};
		});
	},

	clearError: () => set({ error: null }),
}));

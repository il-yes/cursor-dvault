// frontend/src/stores/c3DialogStore.ts
import { create } from 'zustand';

type C3DialogState = {
  channelId: string | null;
  isCreateC3DialogOpen: boolean;
  openC3CreateDialog: (open: boolean, channelId?: string) => void;
  closeCreateDialog: () => void;
  setCreateDialogOpen: (open: boolean) => void;
};

export const useC3DialogStore = create<C3DialogState>((set) => ({
  isCreateC3DialogOpen: false,
  channelId: null,
  openC3CreateDialog: (open: boolean, channelId?: string) => set({ isCreateC3DialogOpen: open, channelId: channelId || null }),
  closeCreateDialog: () => set({ isCreateC3DialogOpen: false, channelId: null }),
  setCreateDialogOpen: (open) => set({ isCreateC3DialogOpen: open }),
}));
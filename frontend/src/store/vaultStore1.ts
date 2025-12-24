import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { VaultContext, VaultEntry } from '@/types/vault';
import { CreateShareEntryPayload, Recipient, SharedEntry } from '@/types/sharing';
import { toast } from '@/hooks/use-toast';
import * as AppAPI from "../../wailsjs/go/main/App";
import mockVaultPayload from '@/data/vault-payload.json';
import { useAuthStore } from '@/store/useAuthStore';
import { wailsBridge } from '@/services/wailsBridge';

interface VaultStoreState {
  vault: VaultContext | null;
  shared: { status: 'idle' | 'loading' | 'loaded'; items: SharedEntry[] };
  sharedWithMe: { status: 'idle' | 'loading' | 'loaded'; items: SharedEntry[] };
  isLoading: boolean;
  lastSyncTime: string | null;

  loadVault: (preloaded?: PreloadedVaultResponse) => Promise<void>;
  setVault: (vault: VaultContext) => void;
  clearVault: () => void;
  setSharedEntries: (sharedEntries: SharedEntry[]) => void;
  addSharedEntry: (entry: CreateShareEntryPayload) => string | null;
  updateSharedEntry: (entryId: string, updates: Partial<SharedEntry>) => void;
  removeSharedEntry: (entryId: string) => void;
  updateSharedEntryRecipients: (entryId: string, recipients: Recipient[]) => void;
  setSharedWithMe: (sharedWithMe: SharedEntry[]) => void;

  hydrateVault: (context: VaultContext) => void;
  refreshVault: () => Promise<void>;

  toggleFavorite: (entryId: string) => Promise<void>;

  addEntry: (entry: VaultEntry) => Promise<void>;
  updateEntry: (entryId: string, updates: Partial<VaultEntry>) => Promise<void>;
  deleteEntry: (entryId: string) => Promise<void>;
  restoreEntry: (entryId: string) => Promise<void>;

  addFolder: (name: string) => Promise<void>;
  sync: (jwtToken: string) => Promise<void>;
  syncVault: (jwtToken: string, vaultPassword: string) => Promise<string>;

  encryptFile: (jwtToken: string, filePath: Uint8Array, vaultPassword: string) => Promise<string>;
  encryptVault: (jwtToken: string, vaultPassword: string) => Promise<string>;
  uploadToIPFS: (jwtToken: string, filePath: string | Uint8Array) => Promise<string>;
  createStellarCommit: (jwtToken: string, cid: string) => Promise<string>;
}

interface PreloadedVaultResponse {
  User: any;
  Vault: any;
  Tokens?: { access_token: string; refresh_token: string };
  SharedEntries?: any[];
  VaultRuntimeContext?: any;
  LastCID?: string;
  Dirty?: boolean;
  vault_runtime_context?: any;
  last_cid?: string;
  dirty?: boolean;
}

const USE_MOCK = import.meta.env.VITE_MOCK_VAULT === 'true';

export const useVaultStore = create<VaultStoreState>()(
  persist(
    (set, get) => ({
      vault: null,
      shared: { status: 'idle', items: [] },
      sharedWithMe: { status: 'idle', items: [] },
      isLoading: false,
      lastSyncTime: null,

      loadVault: async (preloaded?: PreloadedVaultResponse) => {
        set({ isLoading: true });
        try {
          let data: any;

          if (preloaded) data = preloaded;
          else if (USE_MOCK) data = { User: { id: 'mock-user', role: 'user' }, Vault: mockVaultPayload, SharedEntries: [] };
          else {
            const auth = useAuthStore.getState();
            const response = await AppAPI.GetSession(auth.user?.id || '');
            if (response.Error) throw new Error(response.Error);
            const vaultData = response.Data;
            data = { User: vaultData.User, Vault: vaultData.Vault, SharedEntries: vaultData.SharedEntries || [] };
          }

          const vaultObject: VaultContext = {
            user_id: data.User.id,
            role: data.User.role,
            Vault: data.Vault,
            LastCID: data.last_cid || data.LastCID || 'main',
            Dirty: data.dirty || data.Dirty || false,
            LastSynced: data.LastSynced || new Date().toISOString(),
            LastUpdated: data.LastUpdated || new Date().toISOString(),
            vault_runtime_context: data.vault_runtime_context || data.VaultRuntimeContext || data.Vault.vault_runtime_context || {},
          };

          set({
            vault: vaultObject,
            shared: { status: 'loaded', items: data.SharedEntries || [] },
            sharedWithMe: { status: 'loaded', items: [] },
            lastSyncTime: new Date().toISOString(),
            isLoading: false,
          });

          toast({
            title: preloaded ? 'Vault loaded' : USE_MOCK ? 'Vault (Mock)' : 'Vault loaded',
            description: `Last synced: ${new Date().toLocaleString()}`,
          });
        } catch (err) {
          console.error('❌ Failed to load vault:', err);
          set({ isLoading: false });
        }
      },

      setVault: (vault) => set({ vault, lastSyncTime: new Date().toISOString() }),
      clearVault: () => set({ vault: null, shared: { status: 'idle', items: [] }, lastSyncTime: null }),
      setSharedEntries: (sharedEntries) => set({ shared: { status: 'loaded', items: sharedEntries } }),
      setSharedWithMe: (sharedWithMe) => set({ sharedWithMe: { status: 'loaded', items: sharedWithMe } }),

      addSharedEntry: (payload) => {
        const { shared } = get();
        const now = new Date().toISOString();
        const tempShareId = `local-${Date.now()}`;
        const entry: SharedEntry = { ...payload, id: tempShareId, created_at: now, updated_at: now, shared_at: now, audit_log: [], recipients: [] };
        set({ shared: { ...shared, items: [...shared.items, entry] } });
        return tempShareId;
      },

      updateSharedEntry: (entryId, updates) => {
        const { shared } = get();
        set({ shared: { ...shared, items: shared.items.map((item) => (item.id === entryId ? { ...item, ...updates } : item)) } });
      },

      removeSharedEntry: (entryId) => {
        const { shared } = get();
        set({ shared: { ...shared, items: shared.items.filter((item) => item.id !== entryId) } });
      },

      updateSharedEntryRecipients: (entryId, recipients) => {
        const { shared } = get();
        set({ shared: { ...shared, items: shared.items.map((item) => (item.id === entryId ? { ...item, recipients } : item)) } });
      },

      hydrateVault: (context) => get().setVault(context),
      refreshVault: async () => { const context = await wailsBridge.requestContext(); if (context) get().setVault(context); },

      addEntry: async (entry) => {
        set((state) => {
          if (!state.vault) return state;
          const type = entry.type as keyof typeof state.vault.Vault.entries;
          return { vault: { ...state.vault, Vault: { ...state.vault.Vault, entries: { ...state.vault.Vault.entries, [type]: [...(state.vault.Vault.entries[type] || []), entry] } } }, lastSyncTime: new Date().toISOString() };
        });
      },

      updateEntry: async (entryId, updates) => {
        set((state) => {
          if (!state.vault) return state;
          let updated = false;
          const newEntries: any = {};
          for (const type of Object.keys(state.vault.Vault.entries)) {
            newEntries[type] = state.vault.Vault.entries[type].map((e: VaultEntry) =>
              e.id === entryId ? ((updated = true), { ...e, ...updates, updated_at: new Date().toISOString() }) : e
            );
          }
          if (!updated) return state;
          return { vault: { ...state.vault, Vault: { ...state.vault.Vault, entries: newEntries } }, lastSyncTime: new Date().toISOString() };
        });
      },

      deleteEntry: async (entryId) => get().updateEntry(entryId, { trashed: true }),
      restoreEntry: async (entryId) => get().updateEntry(entryId, { trashed: false }),
      toggleFavorite: async (entryId) => {
        set((state) => {
          if (!state.vault) return state;
          const newEntries: any = {};
          for (const type of Object.keys(state.vault.Vault.entries)) {
            newEntries[type] = state.vault.Vault.entries[type].map((e: VaultEntry) =>
              e.id === entryId ? { ...e, is_favorite: !e.is_favorite } : e
            );
          }
          return { vault: { ...state.vault, Vault: { ...state.vault.Vault, entries: newEntries } } };
        });
      },

      addFolder: async (name) => {
        const { jwtToken } = useAuthStore.getState();
        if (!jwtToken) throw new Error("Authentication token not found");
        await AppAPI.CreateFolder(name, jwtToken);
        await get().refreshVault();
      },

      sync: async (jwtToken) => { await AppAPI.SynchronizeVault(jwtToken, ""); await get().refreshVault(); },
      encryptFile: async (jwtToken, fileData, vaultPassword) => AppAPI.EncryptFile(jwtToken, fileData as any, vaultPassword),
      encryptVault: async (jwtToken, vaultPassword) => AppAPI.EncryptVault(jwtToken, vaultPassword),
      uploadToIPFS: async (jwtToken, filePath) => AppAPI.UploadToIPFS(jwtToken, filePath as any),
      createStellarCommit: async (jwtToken, cid) => AppAPI.CreateStellarCommit(jwtToken, cid),
      syncVault: async (jwtToken, vaultPassword) => AppAPI.SynchronizeVault(jwtToken, vaultPassword),

      onRehydrateStorage: () => (state) => console.log("💾 REHYDRATED STATE:", state?.vault),
    }),
    {
      name: "vault-storage",
      merge: (persistedState, currentState) => {
        const persisted = persistedState as Partial<VaultStoreState>;
        return {
          ...currentState,
          ...persisted,
          vault: currentState.vault ?? persisted.vault ?? null,
          shared: currentState.shared ?? persisted.shared ?? { status: 'idle', items: [] },
          sharedWithMe: currentState.sharedWithMe ?? persisted.sharedWithMe ?? { status: 'idle', items: [] },
          lastSyncTime: currentState.lastSyncTime ?? persisted.lastSyncTime ?? null,
        };
      },
    }
  )
);

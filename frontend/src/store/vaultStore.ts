import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { VaultContext } from '@/types/vault';
import { SharedEntry } from '@/types/sharing';
import { toast } from '@/hooks/use-toast';

// Import or paste your mock payload JSON here
import mockVaultPayload from '@/data/vault-payload.json';

interface VaultStoreState {
  vault: VaultContext | null;
  shared: {
    status: 'idle' | 'loading' | 'loaded';
    items: SharedEntry[];
  };
  isLoading: boolean;
  lastSyncTime: string | null;

  // Actions
  loadVault: () => Promise<void>;
  setVault: (vault: VaultContext) => void;
  clearVault: () => void;
  addSharedEntry: (entry: SharedEntry) => void;
  updateSharedEntry: (entryId: string, updates: Partial<SharedEntry>) => void;
  removeSharedEntry: (entryId: string) => void;
}

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';
const USE_MOCK = import.meta.env.VITE_MOCK_VAULT === 'true';

export const useVaultStore = create<VaultStoreState>()(
  persist(
    (set, get) => ({
      vault: null,
      shared: {
        status: 'idle',
        items: [],
      },
      isLoading: false,
      lastSyncTime: null,

      loadVault: async () => {
        set({ isLoading: true });

        try {
          let data: any;

          if (USE_MOCK) {
            try {
              // 🔹 Use local mock payload
              data = mockVaultPayload;
              console.log('✅ Using mock vault payload');
              await new Promise((res) => setTimeout(res, 500)); // simulate delay
            } catch (err) {
              console.error('❌ Failed to load vault mock:', err);
              return null;
            }

          } else {
            // 🔹 Normal backend fetch
            const response = await fetch(`${API_BASE_URL}/api/vault`, {
              method: 'GET',
              headers: { 'Content-Type': 'application/json' },
              credentials: 'include',
            });

            if (!response.ok) throw new Error(`Failed to load vault: ${response.status}`);
            data = await response.json();
          }

          console.log(data.SharedEntries)

          // ✅ Update store state
          set({
            vault: {
              user_id: data.user_id,
              role: data.role,
              Vault: data.Vault,
              LastCID: data.LastCID,
              Dirty: data.Dirty,
              LastSynced: data.LastSynced,
              LastUpdated: data.LastUpdated,
              vault_runtime_context: data.vault_runtime_context,
            },
            shared: {
              status: 'loaded',
              items: data.SharedEntries || [],
            },
            lastSyncTime: new Date().toISOString(),
            isLoading: false,
          });

          toast({
            title: USE_MOCK ? 'Vault (Mock)' : 'Vault loaded',
            description: `Last synced: ${data.LastSynced || 'Just now'}`,
          });
        } catch (error) {
          console.error('❌ Failed to load vault:', error);
          set({ isLoading: false });

          toast({
            title: 'Failed to load vault',
            description: 'Could not connect to backend. Using cached data.',
            variant: 'destructive',
          });
        }
      },

      setVault: (vault) => {
        set({ vault, lastSyncTime: new Date().toISOString() });
      },

      clearVault: () => {
        set({
          vault: null,
          shared: { status: 'idle', items: [] },
          lastSyncTime: null,
        });
      },

      addSharedEntry: (entry) => {
        const { shared } = get();
        set({
          shared: {
            ...shared,
            items: [...shared.items, entry],
          },
        });
      },

      updateSharedEntry: (entryId, updates) => {
        const { shared } = get();
        set({
          shared: {
            ...shared,
            items: shared.items.map((item) =>
              item.id === entryId ? { ...item, ...updates } : item
            ),
          },
        });
      },

      removeSharedEntry: (entryId) => {
        const { shared } = get();
        set({
          shared: {
            ...shared,
            items: shared.items.filter((item) => item.id !== entryId),
          },
        });
      },
    }),
    {
      name: 'vault-storage',
      partialize: (state) => ({
        vault: state.vault,
        shared: state.shared,
        lastSyncTime: state.lastSyncTime,
      }),
    }
  )
);


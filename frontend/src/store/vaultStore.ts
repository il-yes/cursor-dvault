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
  loadVault: (preloaded?: PreloadedVaultResponse) => Promise<void>;
  setVault: (vault: VaultContext) => void;
  clearVault: () => void;
  addSharedEntry: (entry: SharedEntry) => void;
  updateSharedEntry: (entryId: string, updates: Partial<SharedEntry>) => void;
  removeSharedEntry: (entryId: string) => void;
}
interface PreloadedVaultResponse {
  User: any;
  Vault: any;
  Tokens?: {
    access_token: string;
    refresh_token: string;
  };
  SharedEntries?: any[];
  VaultRuntimeContext?: any;
  LastCID?: string;
  Dirty?: boolean;
  vault_runtime_context?: any;
  last_cid?: string;
  dirty?: boolean;
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

      loadVault1: async () => {
        set({ isLoading: true });

        try {
          let data: any;

          if (USE_MOCK) {
            try {
              // 🔹 Use local mock payload
              data = mockVaultPayload;
              console.log('✅ Using mock vault payload');
              await new Promise((res) => setTimeout(res, 500)); // simulate delay
              console.log("Mock vault ")
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
            console.log("Normal backend fetch")
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

      loadVault2: async (preloadedData?: any) => {
        set({ isLoading: true });

        try {
          let data: any;

          if (preloadedData) {
            // ✅ Use already returned session from login
            data = preloadedData;
            console.log("✔ Using preloaded vault from SignIn");
            await new Promise((res) => setTimeout(res, 200)); // optional delay
          } else if (USE_MOCK) {
            data = mockVaultPayload;
            console.log("✔ Using mock vault payload");
            await new Promise((res) => setTimeout(res, 500));
          } else {
            const response = await fetch(`${API_BASE_URL}/api/vault`, {
              method: 'GET',
              headers: { 'Content-Type': 'application/json' },
              credentials: 'include',
            });

            if (!response.ok) throw new Error(`Failed to load vault: ${response.status}`);
            data = await response.json();
            console.log("✔ Normal backend fetch");
          }

          // Update the store
          set({
            vault: data.Vault,
            shared: {
              status: 'loaded',
              items: data.Vault.SharedEntries || [],
            },
            lastSyncTime: new Date().toISOString(),
            isLoading: false,
          });

        } catch (err) {
          console.error("❌ Failed to load vault:", err);
          set({ isLoading: false });
        }
      },

      // 🔹 New: load vault from SignIn response or fetch from backend
      loadVault: async (preloaded?: PreloadedVaultResponse) => {
        set({ isLoading: true });
        try {
          let data: any;

          if (preloaded) {
            // 🔥 Hard validation to avoid silent failures
            if (!preloaded.Vault || !preloaded.User) {
              console.error("❌ loadVault: Invalid preload object", preloaded);
              throw new Error("Preloaded vault invalid — missing Vault or User");
            }

            data = preloaded;
            console.log("✅ Using preloaded vault from SignIn");
          }
          else if (USE_MOCK) {
            // ✅ Use mock data
            data = {
              User: { id: (mockVaultPayload as any).user_id || 'mock-user', role: 'user' },
              Vault: (mockVaultPayload as any).Vault || mockVaultPayload,
              SharedEntries: (mockVaultPayload as any).SharedEntries || [],
            };
            console.log('✅ Using mock vault payload');
            await new Promise((res) => setTimeout(res, 500)); // simulate delay
          } else {
            // ✅ Fetch from backend
            const response = await fetch(`${API_BASE_URL}/api/vault`, {
              method: 'GET',
              headers: { 'Content-Type': 'application/json' },
              credentials: 'include',
            });
            console.log("response", response);

            if (!response.ok) {
              throw new Error(`Failed to load vault: ${response.status}`);
            }

            const vaultData = await response.json();
            data = {
              User: vaultData.User || { id: vaultData.user_id, role: vaultData.role },
              Vault: vaultData.Vault || vaultData,
              SharedEntries: vaultData.SharedEntries || [],
            };
            console.log('✅ Loaded vault from backend');
          }

          const vaultObject = {
            user_id: data.User.id,
            role: data.User.role,
            Vault: data.Vault,
            LastCID: data.last_cid || data.LastCID || 'main',
            Dirty: data.dirty || data.Dirty || false,
            LastSynced: data.LastSynced || new Date().toISOString(),
            LastUpdated: data.LastUpdated || new Date().toISOString(),
            vault_runtime_context: data.vault_runtime_context || data.VaultRuntimeContext || data.Vault.vault_runtime_context || {},
          };

          console.log('📦 vaultStore.loadVault: Setting vault object:', {
            user_id: vaultObject.user_id,
            hasVault: !!vaultObject.Vault,
            hasEntries: !!vaultObject.Vault?.entries,
            entriesKeys: vaultObject.Vault?.entries ? Object.keys(vaultObject.Vault.entries) : [],
            hasRuntimeContext: !!vaultObject.vault_runtime_context,
          });

          set({
            vault: vaultObject,
            shared: {
              status: 'loaded',
              items: data.SharedEntries || [],
            },
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



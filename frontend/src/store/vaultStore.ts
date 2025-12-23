import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { VaultContext } from '@/types/vault';
import { CreateShareEntryPayload, Recipient, SharedEntry } from '@/types/sharing';
import { toast } from '@/hooks/use-toast';
import * as AppAPI from "../../wailsjs/go/main/App";

// Import or paste your mock payload JSON here
import mockVaultPayload from '@/data/vault-payload.json';
import { listSharedEntries, listSharedWithMe } from '@/services/api';
import { useAuthStore } from '@/store/useAuthStore';

interface VaultStoreState {
  vault: VaultContext | null;
  shared: {
    status: 'idle' | 'loading' | 'loaded';
    items: SharedEntry[];
  };
  sharedWithMe: {
    status: 'idle' | 'loading' | 'loaded';
    items: SharedEntry[];
  };
  isLoading: boolean;
  lastSyncTime: string | null;

  // Actions
  loadVault: (preloaded?: PreloadedVaultResponse) => Promise<void>;
  setVault: (vault: VaultContext) => void;
  clearVault: () => void;
  setSharedEntries: (sharedEntries: SharedEntry[]) => void;
  addSharedEntry: (entry: CreateShareEntryPayload) => string;
  updateSharedEntry: (entryId: string, updates: Partial<SharedEntry>) => void;
  removeSharedEntry: (entryId: string) => void;
  updateSharedEntryRecipients: (entryId: string, recipients: Recipient[]) => void;

  setSharedWithMe: (sharedWithMe: SharedEntry[]) => void;

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

const CLOUD_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:4001';
const USE_MOCK = import.meta.env.VITE_MOCK_VAULT === 'true';

export const useVaultStore = create<VaultStoreState>()(
  persist(
    (set, get) => ({
      vault: null,
      shared: {
        status: 'idle',
        items: [],
      },
      sharedWithMe: {
        status: 'idle',
        items: [],
      },
      isLoading: false,
      lastSyncTime: null,

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
            // const response = await fetch(`${CLOUD_BASE_URL}/vault`, {
            //   method: 'GET',
            //   headers: { 'Content-Type': 'application/json' },
            //   credentials: 'include',
            // });

            const auth = useAuthStore.getState();
            const response = await AppAPI.GetSession(auth.user?.id || '') 
            console.log("VaultStore - loadVault response", response);

            if (response.Error) {
              throw new Error(`Failed to load vault: ${response.Error}`);
            }

            const vaultData = response.Data;
            data = {
              User: vaultData?.User || { id: vaultData?.user_id, role: vaultData?.role },
              Vault: vaultData?.Vault || vaultData,
              SharedEntries: vaultData?.SharedEntries || [],
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

          // share by me
          const sharedEntries = await listSharedEntries();
          console.log('✅ Listed shared entries:', sharedEntries);
          // Use the fetched sharedEntries if available, otherwise fall back to preloaded data.SharedEntries
          const finalSharedEntries = (sharedEntries && sharedEntries.length > 0)
            ? sharedEntries
            : (data.SharedEntries || []);

          // share with me
          const sharedWithMe = await listSharedWithMe();
          console.log('✅ Listed shared with me:', sharedWithMe);
          // Use the fetched sharedEntries if available, otherwise fall back to preloaded data.SharedEntries
          const finalSharedWithMe = (sharedWithMe && sharedWithMe.length > 0)
            ? sharedWithMe
            : (data.Sha || []);

          set({
            vault: vaultObject,
            shared: {
              status: 'loaded',
              items: finalSharedEntries,
            },
            sharedWithMe: {
              status: 'loaded',
              items: finalSharedWithMe,
            },
            lastSyncTime: new Date().toISOString(),
            isLoading: false,
          });

          toast({
            title: preloaded ? 'Vault loaded' : USE_MOCK ? 'Vault (Mock)' : 'Vault loaded',
            description: `Last synced: ${new Date().toLocaleString()}`,
          });
        } catch (err) {
          err && console.error('❌ Failed to load vault:', err);
          set({ isLoading: false });

          // toast({
          //   title: 'Failed to load vault',
          //   description: 'Could not connect to backend. Using cached data.',
          //   variant: 'destructive',
          // });
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

      addSharedEntry2: (payload: CreateShareEntryPayload) => {
        try {
          const { shared } = get();

          const now = new Date().toISOString();
          const tempShareId = `local-${Date.now()}`;

          const entry: SharedEntry = {
            id: tempShareId,     // temporary ID until backend returns real one
            created_at: now,
            updated_at: now,
            shared_at: now,

            audit_log: [],

            entry_name: payload.entry_name,
            entry_type: payload.entry_type,
            status: payload.status,
            access_mode: payload.access_mode,
            encryption: payload.encryption,
            entry_snapshot: payload.entry_snapshot,
            expires_at: payload.expires_at,

            // Map simplified recipients to full Recipient objects
            recipients: payload.recipients.map((r, index) => ({
              id: `rec-${Date.now()}-${index}`,
              share_id: tempShareId,
              name: r.name,
              email: r.email,
              role: r.role,
              created_at: now,
              updated_at: now,
            })),
          };

          set({
            shared: {
              ...shared,
              items: [...shared.items, entry],
            },
          });

        } catch (err) {
          console.error('❌ Failed to add shared entry:', err);
        }
      },

      addSharedEntry: (payload: CreateShareEntryPayload) => {
        try {
          const { shared } = get();

          const now = new Date().toISOString();
          const tempShareId = `local-${Date.now()}`;

          const entry: SharedEntry = {
            id: tempShareId,
            created_at: now,
            updated_at: now,
            shared_at: now,

            audit_log: [],

            entry_name: payload.entry_name,
            entry_type: payload.entry_type,

            // always pending when optimistic
            status: "pending",
            access_mode: payload.access_mode,
            encryption: payload.encryption,
            entry_snapshot: payload.entry_snapshot,
            expires_at: payload.expires_at,

            // IMPORTANT: no recipients during optimistic state
            recipients: [],
          };

          set({
            shared: {
              ...shared,
              items: [...shared.items, entry],
            },
          });

          // return temporary ID so caller can replace it later
          return tempShareId;

        } catch (err) {
          console.error("❌ Failed to add shared entry:", err);
        }
      },

      setSharedEntries: (sharedEntries: SharedEntry[]) => {
        set({
          shared: {
            status: 'loaded',
            items: sharedEntries,
          },
        });
      },

      setSharedWithMe: (sharedWithMe: SharedEntry[]) => {
        set({
          sharedWithMe: {
            status: 'loaded',
            items: sharedWithMe,
          },
        });
      },

      updateSharedEntry: (entryId: string, updates) => {
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
      updateSharedEntryRecipients: (entryId: string, recipients) => {
        const { shared } = get();
        set({
          shared: {
            ...shared,
            items: shared.items.map((item) =>
              item.id === entryId
                ? { ...item, recipients: [...recipients] }
                : item
            ),
          },
        });
      }
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



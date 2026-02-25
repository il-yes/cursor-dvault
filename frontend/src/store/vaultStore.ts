import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { VaultContext, VaultEntry } from '@/types/vault';
import { CreateLinkShareEntryPayload, CreateShareEntryPayload, LinkShareEntry, Recipient, SharedEntry } from '@/types/sharing';
import { toast } from '@/hooks/use-toast';
import * as AppAPI from "../../wailsjs/go/main/App";

// Import or paste your mock payload JSON here
import mockVaultPayload from '@/data/vault-payload.json';
import { listSharedEntries, listSharedWithMe, listLinkSharesByMe, listLinkSharesWithMe, updateEntry as updateEntryApi } from '@/services/api';
import { useAuthStore } from '@/store/useAuthStore';
import { wailsBridge } from '@/services/wailsBridge';
import { get } from 'http';
import { withAuth } from '@/hooks/withAuth';

// TODO: Backend expects file path string, not Uint8Array.


interface VaultStoreState {
  vault: VaultContext | null;
  linkSharedByMe: {
    status: 'idle' | 'loading' | 'loaded';
    items: LinkShareEntry[];
  };
  linkSharedWithMe: {
    status: 'idle' | 'loading' | 'loaded';
    items: LinkShareEntry[];
  };
  shared: { // cryptographic share by me
    status: 'idle' | 'loading' | 'loaded';
    items: SharedEntry[];
  };
  sharedWithMe: { // cryptographic share with me
    status: 'idle' | 'loading' | 'loaded';
    items: SharedEntry[];
  };
  isLoading: boolean;
  lastSyncTime: string | null;

  // Actions
  loadVault: (preloaded?: PreloadedVaultResponse) => Promise<void>;
  setVault: (vault: VaultContext) => void;
  clearVault: () => void;
  // Cryptographic share
  setSharedEntries: (sharedEntries: SharedEntry[]) => void;
  addSharedEntry: (entry: CreateShareEntryPayload) => string | null;
  updateSharedEntry: (entryId: string, updates: Partial<SharedEntry>) => void;
  removeSharedEntry: (entryId: string) => void;
  updateSharedEntryRecipients: (entryId: string, recipients: Recipient[]) => void;
  setSharedWithMe: (sharedWithMe: SharedEntry[]) => void;
  // Link share
  setLinkSharedEntries: (linkSharedEntries: LinkShareEntry[]) => void;
  addLinkSharedEntry: (entry: LinkShareEntry) => void;
  updateLinkSharedEntry: (entryId: string, updates: Partial<LinkShareEntry>) => void;
  removeLinkSharedEntry: (entryId: string) => void;

  // From useVault

  hydrateVault: (context: VaultContext) => void;
  refreshVault: () => Promise<void>;

  toggleFavorite: (entryId: string) => Promise<void>;

  addEntry: (entry: VaultEntry) => Promise<void>;
  updateEntry: (entryId: string, updates: Partial<VaultEntry>) => Promise<void>;
  deleteEntry: (entryId: string) => Promise<void>;
  restoreEntry: (entryId: string) => Promise<void>;

  addFolder: (name: string) => Promise<void>;
  removeFolder: (folderId: string) => Promise<void>;


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
      linkSharedByMe: {
        status: 'idle',
        items: [],
      },
      linkSharedWithMe: {
        status: 'idle',
        items: [],
      },
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

      // 🔹 Vault
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
            console.log("✅ Using preloaded vault from SignIn", data);
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
            user_id: data.User.id || data?.User?.ID,
            role: data.User.role,
            Vault: {
              ...data.Vault,
              entries: data.Vault.entries
            },
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

          // Cryptographic share by me
          const sharedEntries = await listSharedEntries();
          // Use the fetched sharedEntries if available, otherwise fall back to preloaded data.SharedEntries
          const finalSharedEntries = (sharedEntries && sharedEntries.length > 0)
            ? sharedEntries
            : (data.SharedEntries || []);

          // Cryptographic share with me
          const sharedWithMe = await listSharedWithMe();
          const finalSharedWithMe = (sharedWithMe && sharedWithMe.length > 0)
            ? sharedWithMe
            : (data.SharedWithMe || []);

          // Link share by me
          const linkSharedByMe = await listLinkSharesByMe();
          const finalLinkSharedByMe = (linkSharedByMe && linkSharedByMe.length > 0)
            ? linkSharedByMe
            : (data.linkSharedByMe || []);

          // Link share with me
          const linkSharedWithMe = await listLinkSharesWithMe();
          const finalLinkSharedWithMe = (linkSharedWithMe && linkSharedWithMe.length > 0)
            ? linkSharedWithMe
            : (data.linkSharedWithMe || []);

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
            linkSharedByMe: {
              status: 'loaded',
              items: finalLinkSharedByMe,
            },
            linkSharedWithMe: {
              status: 'loaded',
              items: finalLinkSharedWithMe,
            },
            lastSyncTime: new Date().toISOString(),
            isLoading: false,
          });

          console.log("🧪 STORE IMMEDIATE CHECK:", get().vault);


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

      // V0
      clearVault: () => {
        set({
          vault: null,
          shared: { status: 'idle', items: [] },
          lastSyncTime: null,
        });
      },

      // 🔹 Cryptographic Share
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
      },

      // 🔹 Link Share
      setLinkSharedEntries: (linkSharedEntries: LinkShareEntry[]) => {
        set({
          linkSharedByMe: {
            status: 'loaded',
            items: linkSharedEntries,
          },
        });
      },

      addLinkSharedEntry: (entry: LinkShareEntry) => {
        const { linkSharedByMe } = get();

        set({
          linkSharedByMe: {
            ...linkSharedByMe,
            items: [...linkSharedByMe.items, entry],
          },
        });
      },

      updateLinkSharedEntry: (entryId: string, updates: Partial<LinkShareEntry>) => {
        const { linkSharedByMe } = get();
        set({
          linkSharedByMe: {
            ...linkSharedByMe,
            items: linkSharedByMe.items.map((item) =>
              item.id === entryId ? { ...item, ...updates } : item
            ),
          },
        });
      },

      removeLinkSharedEntry: (entryId: string) => {
        const { linkSharedByMe } = get();
        set({
          linkSharedByMe: {
            ...linkSharedByMe,
            items: linkSharedByMe.items.filter((item) => item.id !== entryId),
          },
        });
      },

      /* ------------------------------------------------------------------ */
      /* Backend → Store hydration                                           */
      /* ------------------------------------------------------------------ */


      loadInitialContext: async () => {
        try {
          const context = await wailsBridge.requestContext();
          if (context) {
            get().hydrateVault(context);
          }
        } catch (err) {
          console.error("❌ Failed to load initial vault context:", err);
        }
      },

      hydrateVault: (context: VaultContext) => {
        useVaultStore.getState().setVault(context);

        toast({
          title: "Vault ready",
          description: `Last synced: ${context.LastSynced || "Just now"}`,
        });
      },

      refreshVault: async () => {
        try {
          const context = await wailsBridge.requestContext();
          if (context) {
            useVaultStore.getState().setVault(context);
          }
        } catch (err) {
          console.error("❌ Failed to refresh vault:", err);
          toast({
            title: "Sync failed",
            description: "Could not refresh vault data.",
            variant: "destructive",
          });
        }
      },


      /* ------------------------------------------------------------------ */
      /* Mutations — ALL go through Zustand                                  */
      /* ------------------------------------------------------------------ */

      addEntry: async (entry: VaultEntry): Promise<void> => {
        useVaultStore.setState((state) => {
          if (!state.vault) return state;
          state.vault.Dirty = true;

          const type = entry.type as keyof typeof state.vault.Vault.entries;

          return {
            vault: {
              ...state.vault,
              Vault: {
                ...state.vault.Vault,
                entries: {
                  ...state.vault.Vault.entries,
                  [type]: [...(state.vault.Vault.entries[type] || []), entry],
                },
              },
            },
            lastSyncTime: new Date().toISOString(),
          };
        });
      },

      updateEntryV0: async (entryId: string, updates: Partial<VaultEntry>) => {
        useVaultStore.setState((state) => {
          if (!state.vault) return state;

          const newEntries: any = {};
          let updated = false;

          for (const type of Object.keys(state.vault.Vault.entries)) {
            const entries = state.vault.Vault.entries[type];
            if (!Array.isArray(entries)) {
              newEntries[type] = entries;
              continue;
            }

            newEntries[type] = entries.map((e) =>
              e.id === entryId
                ? ((updated = true),
                  { ...e, ...updates, updated_at: new Date().toISOString() })
                : e
            );
          }

          if (!updated) return state;

          return {
            vault: {
              ...state.vault,
              Vault: {
                ...state.vault.Vault,
                entries: newEntries,
              },
            },
          };
        });
      },
      updateEntry: async (entryId: string, updates: Partial<VaultEntry>) => {
        useVaultStore.setState((state) => {
          if (!state.vault) return state;

          const newEntries: any = {};
          let updated = false;

          for (const type of Object.keys(state.vault.Vault.entries)) {
            const entries = state.vault.Vault.entries[type];
            if (!entries) continue; // <-- safeguard for null/undefined
            newEntries[type] = entries.map((e) =>
              e.id === entryId
                ? ((updated = true), { ...e, ...updates, updated_at: new Date().toISOString() })
                : e
            );
          }

          if (!updated) return state;

          state.vault.Dirty = true;

          return {
            vault: {
              ...state.vault,
              Vault: {
                ...state.vault.Vault,
                entries: newEntries,
              },
            },
            lastSyncTime: new Date().toISOString(),
          };
        });
      },


      deleteEntry: async (entryId: string) => {
        await get().updateEntry(entryId, { trashed: true });
      },

      restoreEntry: async (entryId: string) => {
        await get().updateEntry(entryId, { trashed: false });
      },

      toggleFavorite: async (entryId: string) => {
        useVaultStore.setState((state) => {
          if (!state.vault) return state;

          const newEntries: any = {};

          for (const type of Object.keys(state.vault.Vault.entries)) {
            const entries = state.vault.Vault.entries[type];
            newEntries[type] = entries.map((e) =>
              e.id === entryId ? { ...e, is_favorite: !e.is_favorite } : e
            );
          }
          console.log("🚀 ~ toggleFavorite ~ newEntries:", newEntries)

          return {
            vault: {
              ...state.vault,
              Vault: {
                ...state.vault.Vault,
                entries: newEntries,
              },
            },
          };
        });
      },

      /* ------------------------------------------------------------------ */
      /* Additional Methods Implementation                                   */
      /* ------------------------------------------------------------------ */

      addFolder: async (name: string) => {
        useVaultStore.setState((state) => {
          if (!state.vault) return state;
          state.vault.Dirty = true;


          return {
            vault: {
              ...state.vault,
              Vault: {
                ...state.vault.Vault,
                folders: [...state.vault.Vault.folders, { id: name, name }],
              },
            },
            lastSyncTime: new Date().toISOString(),
          };
        });
        // try {
        //     await withAuth(async (token) => {
        //       await AppAPI.CreateFolder(name, token);
        //     });
        //   await get().refreshVault();
        // } catch (err) {
        //   console.error("❌ Failed to add folder:", err);
        //   toast({
        //     title: "Error",
        //     description: `Failed to create folder: ${(err as Error).message}`,
        //     variant: "destructive",
        //   });
        //   throw err;
        // }
      },

      removeFolder: async (folderId: string) => {
        useVaultStore.setState((state) => {
          if (!state.vault) return state;
          state.vault.Dirty = true;


          return {
            vault: {
              ...state.vault,
              Vault: {
                ...state.vault.Vault,
                folders: state.vault.Vault.folders.filter((folder) => folder.id !== folderId),
              },
            },
            lastSyncTime: new Date().toISOString(),
          };
        });
      },

      sync: async (jwtToken: string) => {
        try {
          await AppAPI.SynchronizeVault(jwtToken, ""); // Password might be needed depending on implementation
          await get().refreshVault();
        } catch (err) {
          console.error("❌ Failed to sync:", err);
          throw err;
        }
      },

      encryptFile: async (jwtToken: string, fileData: Uint8Array, vaultPassword: string) => {
        // Note: AppAPI.EncryptFile expects (jwt, filePath, password). 
        // If we have Uint8Array, we might need to handle it differently 
        // or the backend expects a path. ProfileBeta.tsx uses readFileAsBuffer then calls this.
        // Based on App.d.ts: export function EncryptFile(arg1:string,arg2:string,arg3:string):Promise<string>;
        // arg2 is string (filePath). If ProfileBeta passes Uint8Array, there might be a mismatch 
        // or a different API needed. However, ProfileBeta.tsx:L110 shows it being used.
        // Actually, ProfileBeta.tsx:L110: const encryptedData = await encryptFile(jwtToken, filePath, vaultPassword);
        // where filePath is Uint8Array from readFileAsBuffer. 
        // This implies VaultContextValue's encryptFile takes Uint8Array.
        // But AppAPI.EncryptFile takes string (path).
        // I will implementation it to call AppAPI.EncryptFile and cast for now to satisfy types,
        // but this might need refinement if the backend actually wants a path.
        return await AppAPI.EncryptFile(jwtToken, fileData as any, vaultPassword);
      },

      encryptVault: async (jwtToken: string, vaultPassword: string) => {
        return await AppAPI.EncryptVault(jwtToken, vaultPassword);
      },

      uploadToIPFS: async (jwtToken: string, filePath: string | Uint8Array) => {
        return await AppAPI.UploadToIPFS(jwtToken, filePath as any);
      },

      createStellarCommit: async (jwtToken: string, cid: string) => {
        return await AppAPI.CreateStellarCommit(jwtToken, cid);
      },

      syncVault: async (jwtToken: string, vaultPassword: string) => {
        return await AppAPI.SynchronizeVault(jwtToken, vaultPassword);
      },

      /* ------------------------------------------------------------------ */
      onRehydrateStorage: () => (state) => {
        console.log("💾 REHYDRATED STATE:", state?.vault);
      },
    }),
    {
      name: "vault-storage",

      merge: (persistedState, currentState) => {
        const persisted = persistedState as Partial<VaultStoreState>;

        return {
          ...currentState,
          ...persisted,
          vault: persisted.vault ?? currentState.vault ?? null,         // <-- fix here
          shared: persisted.shared ?? currentState.shared ?? { status: 'idle', items: [] },
          sharedWithMe: persisted.sharedWithMe ?? currentState.sharedWithMe ?? { status: 'idle', items: [] },
          lastSyncTime: persisted.lastSyncTime ?? currentState.lastSyncTime ?? null,
        };
      }





    }

  )
);



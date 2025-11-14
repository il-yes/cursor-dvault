import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { VaultContext, VaultEntry, Folder } from '@/types/vault';
import { wailsBridge } from '@/services/wailsBridge';
import { toast } from '@/hooks/use-toast';
import { mockFolders, getMockVaultEntries } from '@/data/mockVaultData';

interface VaultContextValue {
  vaultContext: VaultContext | null;
  isLoading: boolean;
  hydrateVault: (context: VaultContext) => void;
  refreshVault: () => Promise<void>;
  clearVault: () => void;
  addEntry: (entry: VaultEntry) => void;
  updateEntry: (entryId: string, updates: Partial<VaultEntry>) => void;
  deleteEntry: (entryId: string) => void;
  toggleFavorite: (entryId: string) => void;
  addFolder: (name: string) => void;
}

const VaultContextProvider = createContext<VaultContextValue | undefined>(undefined);

export function VaultProvider({ children }: { children: ReactNode }) {
  const [vaultContext, setVaultContext] = useState<VaultContext | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    // Setup Wails bridge listener
    const unsubscribe = wailsBridge.onContextReceived((context) => {
      hydrateVault(context);
    });

    // Try to load context on mount or initialize with mock data
    loadInitialContext();

    return () => {
      unsubscribe();
    };
  }, []);

  const loadInitialContext = async () => {
    setIsLoading(true);
    try {
      const context = await wailsBridge.requestContext();
      if (context) {
        hydrateVault(context);
      } else {
        // Initialize with mock data if no context
        initializeMockVault();
      }
    } catch (error) {
      console.error('Failed to load initial vault context:', error);
      // Fallback to mock data
      initializeMockVault();
    } finally {
      setIsLoading(false);
    }
  };

  const initializeMockVault = () => {
    const mockContext: VaultContext = {
      user_id: "mock-user-123",
      role: "owner",
      Vault: {
        version: "1.0.0",
        name: "Sovereign Vault",
        folders: mockFolders,
        entries: getMockVaultEntries(),
      },
      Dirty: false,
      LastSynced: new Date().toISOString(),
      LastUpdated: new Date().toISOString(),
      vault_runtime_context: {
        CurrentUser: {
          id: "mock-user-123",
          role: "owner",
          name: "John",
          last_name: "Doe",
          email: "john.doe@example.com",
          stellar_account: {
            public_key: "GXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
          },
        },
        AppSettings: {
          id: "mock-app-settings-123",
          repo_id: "mock-repo-123",
          branch: "main",
          tracecore_enabled: false,
          commit_rules: [],
          branching_model: "single",
          encryption_policy: "AES-256-GCM",
          actors: ["user"],
          federated_providers: null,
          default_phase: "vault_entry",
          default_vault_path: "",
          vault_settings: {
            max_entries: 1000,
            encryption_scheme: "AES-256-GCM",
          },
          blockchain: {
            stellar: {
              network: "testnet",
              horizon_url: "https://horizon-testnet.stellar.org",
              fee: 100,
            },
            ipfs: {
              gateway_url: "https://ipfs.io",
              api_endpoint: "http://localhost:5001",
            },
          },
          auto_sync_enabled: true,
        },
        WorkingBranch: "main",
        LoadedEntries: [],
      },
    };
    setVaultContext(mockContext);
  };

  const hydrateVault = (context: VaultContext) => {
    setVaultContext(context);
    
    // Show sync status
    if (context.Dirty && !context.vault_runtime_context.AppSettings.auto_sync_enabled) {
      toast({
        title: "Vault loaded",
        description: `Last synced: ${context.LastSynced || 'Never'}. Local changes pending.`,
      });
    } else {
      toast({
        title: "Vault ready",
        description: `Synced locally. Last synced: ${context.LastSynced || 'Just now'}`,
      });
    }
  };

  const refreshVault = async () => {
    setIsLoading(true);
    try {
      const context = await wailsBridge.requestContext();
      if (context) {
        setVaultContext(context);
      }
    } catch (error) {
      console.error('Failed to refresh vault:', error);
      toast({
        title: "Sync failed",
        description: "Could not refresh vault data.",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  const clearVault = () => {
    setVaultContext(null);
  };

  const addEntry = (entry: VaultEntry) => {
    if (!vaultContext) return;
    
    const updatedContext = { ...vaultContext };
    updatedContext.Vault.entries[entry.type].push(entry);
    updatedContext.Dirty = true;
    updatedContext.LastUpdated = new Date().toISOString();
    
    setVaultContext(updatedContext);
    toast({
      title: "Entry created",
      description: `${entry.entry_name} has been added to your vault.`,
    });
  };

  const updateEntry = (entryId: string, updates: Partial<VaultEntry>) => {
    if (!vaultContext) return;
    
    const updatedContext = { ...vaultContext };
    let found = false;
    
    // Find and update the entry across all types
    Object.keys(updatedContext.Vault.entries).forEach((type) => {
      const entries = updatedContext.Vault.entries[type as keyof typeof updatedContext.Vault.entries];
      const index = entries.findIndex(e => e.id === entryId);
      if (index !== -1) {
        entries[index] = { ...entries[index], ...updates, updated_at: new Date().toISOString() } as VaultEntry;
        found = true;
      }
    });
    
    if (found) {
      updatedContext.Dirty = true;
      updatedContext.LastUpdated = new Date().toISOString();
      setVaultContext(updatedContext);
      toast({
        title: "Entry updated",
        description: "Changes saved successfully.",
      });
    }
  };

  const deleteEntry = (entryId: string) => {
    if (!vaultContext) return;
    
    updateEntry(entryId, { trashed: true } as Partial<VaultEntry>);
    toast({
      title: "Entry moved to trash",
      description: "You can restore it from the Trash folder.",
    });
  };

  const toggleFavorite = (entryId: string) => {
    if (!vaultContext) return;
    
    let currentFavorite = false;
    Object.keys(vaultContext.Vault.entries).forEach((type) => {
      const entries = vaultContext.Vault.entries[type as keyof typeof vaultContext.Vault.entries];
      const entry = entries.find(e => e.id === entryId);
      if (entry) {
        currentFavorite = entry.is_favorite || false;
      }
    });
    
    updateEntry(entryId, { is_favorite: !currentFavorite });
  };

  const addFolder = (name: string) => {
    if (!vaultContext) return;
    
    const newFolder: Folder = {
      id: `folder-${Date.now()}`,
      name,
      icon: "📁",
    };
    
    const updatedContext = { ...vaultContext };
    updatedContext.Vault.folders.push(newFolder);
    updatedContext.Dirty = true;
    
    setVaultContext(updatedContext);
    toast({
      title: "Folder created",
      description: `${name} has been added to your vault.`,
    });
  };

  return (
    <VaultContextProvider.Provider
      value={{
        vaultContext,
        isLoading,
        hydrateVault,
        refreshVault,
        clearVault,
        addEntry,
        updateEntry,
        deleteEntry,
        toggleFavorite,
        addFolder,
      }}
    >
      {children}
    </VaultContextProvider.Provider>
  );
}

export function useVault(): VaultContextValue {
  const context = useContext(VaultContextProvider);
  if (context === undefined) {
    throw new Error('useVault must be used within a VaultProvider');
  }
  return context;
}

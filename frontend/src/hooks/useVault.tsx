import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { VaultContext, VaultEntry, Folder } from '@/types/vault';
import { wailsBridge } from '@/services/wailsBridge';
import { toast } from '@/hooks/use-toast';
import { mockFolders, getMockVaultEntries } from '@/data/mockVaultData';
import * as AppAPI from "../../wailsjs/go/main/App";
import { useVaultStore } from "@/store/vaultStore";
import { useAuthStore } from "@/store/useAuthStore";


interface VaultContextValue {
	vaultContext: VaultContext | null;
	isLoading: boolean;
	hydrateVault: (context: VaultContext) => void;
	refreshVault: () => Promise<void>;
	clearVault: () => void;
	addEntry: (entry: VaultEntry) => Promise<void>;
	updateEntry: (entryId: string, updates: Partial<VaultEntry>) => Promise<void>;
	deleteEntry: (entryId: string) => Promise<void>;
	restoreEntry: (entryId: string) => Promise<void>;
	toggleFavorite: (entryId: string) => Promise<void>;
	addFolder: (name: string) => Promise<void>;
}

const VaultContextProvider = createContext<VaultContextValue | undefined>(undefined);

export function VaultProvider({ children }: { children: ReactNode }) {
	const [vaultContext, setVaultContext] = useState<VaultContext | null>(null);
	const [isLoading, setIsLoading] = useState(true);
	const vaultStoreData = useVaultStore((state) => state.vault);

	// Sync vaultStore data to vaultContext
	useEffect(() => {
		console.log('🔄 VaultProvider: useEffect triggered', {
			hasVaultStoreData: !!vaultStoreData,
			vaultStoreData: vaultStoreData ? {
				user_id: vaultStoreData.user_id,
				hasVault: !!vaultStoreData.Vault,
				hasEntries: !!vaultStoreData.Vault?.entries,
				lastUpdated: vaultStoreData.LastUpdated
			} : null
		});

		if (vaultStoreData) {
			console.log('✅ VaultProvider: Syncing vaultStoreData to vaultContext');
			setVaultContext(vaultStoreData);
			setIsLoading(false);
		} else {
			console.log('⚠️ VaultProvider: No vaultStoreData to sync');
		}
	}, [vaultStoreData]);

	useEffect(() => {
		// Listen to context pushes from backend
		const unsubscribe = wailsBridge.onContextReceived((context) => {
			hydrateVault(context);
		});

		// Load context at startup only if vaultStore is empty
		if (!vaultStoreData) {
			loadInitialContext();
		} else {
			setIsLoading(false);
		}

		return () => {
			unsubscribe();
		};
	}, []);

	const loadInitialContext = async () => {
		try {
			const context = await wailsBridge.requestContext();

			if (context) {
				hydrateVault(context);
			} else {
				// No mock! User is probably not logged in
				setVaultContext(null);
			}
		} catch (error) {
			console.error("Failed to load initial vault context:", error);
			setVaultContext(null);
		} finally {
			setIsLoading(false);
		}
	};

	const hydrateVault = (context: VaultContext) => {
		// Update vaultStore - this will trigger the useEffect to update vaultContext
		useVaultStore.getState().setVault(context);

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
				// Update vaultStore - this will trigger the useEffect to update vaultContext
				useVaultStore.getState().setVault(context);
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

	// Add a new entry
	const addEntry = async (entry: VaultEntry): Promise<void> => {
		if (!vaultContext) return;

		const type = entry.type as keyof typeof vaultContext.Vault.entries;

		// ✅ Entry already has backend-generated ID from DashboardLayout
		// Create updated context
		const updatedContext = {
			...vaultContext,
			Vault: {
				...vaultContext.Vault,
				entries: {
					...vaultContext.Vault.entries,
					[type]: [...(vaultContext.Vault.entries[type] || []), entry],
				},
			},
			Dirty: true,
			LastUpdated: new Date().toISOString(),
		};

		// Update vaultStore - this will trigger the useEffect to update vaultContext
		useVaultStore.getState().setVault(updatedContext);
	};

	// Update an entry
	const updateEntry = async (entryId: string, updates: Partial<VaultEntry>): Promise<void> => {
		if (!vaultContext) {
			console.error("❌ Cannot update entry: vaultContext is null");
			toast({ title: "Error", description: "Vault not loaded. Please refresh the page.", variant: "destructive" });
			return;
		}

		const { jwtToken } = useAuthStore.getState();
		let updated = false;
		let updatedEntry: VaultEntry | null = null;
		let entryType: string = '';

		// Create a deep copy of entries with new array references
		const newEntries: any = {};

		Object.keys(vaultContext.Vault.entries).forEach((type) => {
			const entries = vaultContext.Vault.entries[type as keyof typeof vaultContext.Vault.entries];

			// ✅ Safety check: ensure entries array exists
			if (!entries || !Array.isArray(entries)) {
				console.warn(`⚠️ No entries array for type: ${type}`);
				newEntries[type] = entries;
				return;
			}

			const index = entries.findIndex((e) => e.id === entryId);
			if (index !== -1) {
				// Create a new array with the updated entry
				const updatedEntryData = { ...entries[index], ...updates, updated_at: new Date().toISOString() } as VaultEntry;
				newEntries[type] = [
					...entries.slice(0, index),
					updatedEntryData,
					...entries.slice(index + 1)
				];
				updatedEntry = updatedEntryData;
				entryType = updatedEntryData.type.toLowerCase();
				updated = true;
			} else {
				// Keep the same array if no update
				newEntries[type] = entries;
			}
		});

		if (!updated || !updatedEntry) {
			console.error(`❌ Entry not found: ${entryId}`);
			toast({ title: "Error", description: "Entry not found in vault.", variant: "destructive" });
			return;
		}

		// 1️⃣ Update context immediately (optimistic update)
		const updatedContext = {
			...vaultContext,
			Vault: {
				...vaultContext.Vault,
				entries: newEntries // New entries object with new array references
			},
			Dirty: true,
			LastUpdated: new Date().toISOString()
		};

		// Update vaultStore - this will trigger the useEffect to update vaultContext
		useVaultStore.getState().setVault(updatedContext);

		if (!jwtToken) {
			toast({ title: "Warning", description: "Changes saved locally only (not authenticated).", variant: "default" });
			return;
		}

		try {
			// 2️⃣ Persist to backend with correct signature: EditEntry(entryType, fullEntry, jwtToken)
			// Ensure the entry type in the payload is also lowercase
			const entryPayload = { ...updatedEntry, type: entryType };
			await AppAPI.EditEntry(entryType, entryPayload, jwtToken);
			toast({ title: "Entry updated", description: "Changes saved successfully." });
		} catch (err) {
			console.error(err);
			toast({ title: "Error", description: "Failed to update entry on backend.", variant: "destructive" });
		}
	};

	// Delete (move to trash)
	const deleteEntry = async (entryId: string): Promise<void> => {
		if (!vaultContext) {
			console.error("❌ Cannot delete entry: vaultContext is null");
			toast({ title: "Error", description: "Vault not loaded. Please refresh the page.", variant: "destructive" });
			return;
		}

		// Find the entry to get its type and data
		let entryToTrash: VaultEntry | null = null;
		let entryType: string = '';

		Object.keys(vaultContext.Vault.entries).forEach((type) => {
			const entries = vaultContext.Vault.entries[type as keyof typeof vaultContext.Vault.entries];

			// ✅ Safety check: ensure entries array exists
			if (!entries || !Array.isArray(entries)) {
				return;
			}

			const entry = entries.find(e => e.id === entryId);
			if (entry) {
				entryToTrash = entry;
				// ✅ Ensure type is lowercase for backend compatibility
				entryType = entry.type.toLowerCase();
			}
		});

		if (!entryToTrash) {
			console.error(`❌ Entry not found: ${entryId}`);
			toast({ title: "Error", description: "Entry not found.", variant: "destructive" });
			return;
		}

		updateEntry(entryId, { trashed: true });
		const { jwtToken } = useAuthStore.getState();
		if (!jwtToken) return;

		try {
			// Ensure the entry type in the payload is also lowercase
			const entryPayload = { ...entryToTrash, type: entryType };
			await AppAPI.TrashEntry(entryType, entryPayload, jwtToken);
			toast({ title: "Entry moved to trash", description: "You can restore it from the Trash folder." });
		} catch (err) {
			console.error(err);
			toast({ title: "Error", description: "Failed to delete entry on backend.", variant: "destructive" });
		}
	};

	// Restore a trashed entry
	const restoreEntry = async (entryId: string): Promise<void> => {
		if (!vaultContext) {
			console.error("❌ Cannot restore entry: vaultContext is null");
			toast({ title: "Error", description: "Vault not loaded. Please refresh the page.", variant: "destructive" });
			return;
		}

		// Find the entry to get its type and data
		let entryToRestore: VaultEntry | null = null;
		let entryType: string = '';

		Object.keys(vaultContext.Vault.entries).forEach((type) => {
			const entries = vaultContext.Vault.entries[type as keyof typeof vaultContext.Vault.entries];

			// ✅ Safety check: ensure entries array exists
			if (!entries || !Array.isArray(entries)) {
				return;
			}

			const entry = entries.find(e => e.id === entryId);
			if (entry) {
				entryToRestore = entry;
				// ✅ Ensure type is lowercase for backend compatibility
				entryType = entry.type.toLowerCase();
			}
		});

		if (!entryToRestore) {
			console.error(`❌ Entry not found: ${entryId}`);
			toast({ title: "Error", description: "Entry not found.", variant: "destructive" });
			return;
		}

		updateEntry(entryId, { trashed: false });
		const { jwtToken } = useAuthStore.getState();
		if (!jwtToken) return;

		try {
			// Ensure the entry type in the payload is also lowercase
			const entryPayload = { ...entryToRestore, type: entryType };
			await AppAPI.RestoreEntry(entryType, entryPayload, jwtToken);
			toast({ title: "Entry restored", description: "Entry restored successfully." });
		} catch (err) {
			console.error(err);
			toast({ title: "Error", description: "Failed to restore entry on backend.", variant: "destructive" });
		}
	};

	// Toggle favorite
	const toggleFavorite = async (entryId: string): Promise<void> => {
		if (!vaultContext) {
			console.error("❌ Cannot toggle favorite: vaultContext is null");
			toast({ title: "Error", description: "Vault not loaded. Please refresh the page.", variant: "destructive" });
			return;
		}

		// Find the entry to get its type and current favorite status
		let entryToUpdate: VaultEntry | null = null;
		let entryType: string = '';

		Object.keys(vaultContext.Vault.entries).forEach((type) => {
			const entries = vaultContext.Vault.entries[type as keyof typeof vaultContext.Vault.entries];

			// ✅ Safety check: ensure entries array exists
			if (!entries || !Array.isArray(entries)) {
				return;
			}

			const entry = entries.find(e => e.id === entryId);
			if (entry) {
				entryToUpdate = entry;
				// ✅ Ensure type is lowercase for backend compatibility
				entryType = entry.type.toLowerCase();
			}
		});

		if (!entryToUpdate) {
			console.error(`❌ Entry not found: ${entryId}`);
			toast({ title: "Error", description: "Entry not found.", variant: "destructive" });
			return;
		}

		const newFavoriteStatus = !entryToUpdate.is_favorite;
		updateEntry(entryId, { is_favorite: newFavoriteStatus });

		const { jwtToken } = useAuthStore.getState();
		if (!jwtToken) return;

		try {
			const updatedEntry = { ...entryToUpdate, is_favorite: newFavoriteStatus };
			await AppAPI.EditEntry(entryType, updatedEntry, jwtToken);
		} catch (err) {
			console.error(err);
			toast({ title: "Error", description: "Failed to toggle favorite on backend.", variant: "destructive" });
		}
	};

	// Add folder
	const addFolder = async (name: string): Promise<void> => {
		if (!vaultContext) return;

		const newFolder: Folder = { id: `folder-${Date.now()}`, name, icon: "📁" };
		setVaultContext({ ...vaultContext, Vault: { ...vaultContext.Vault, folders: [...vaultContext.Vault.folders, newFolder] }, Dirty: true, LastUpdated: new Date().toISOString() });

		const { jwtToken } = useAuthStore.getState();
		if (!jwtToken) return;

		try {
			await AppAPI.CreateFolder(name, jwtToken);
			toast({ title: "Folder created", description: `${name} added successfully.` });
		} catch (err) {
			console.error(err);
			toast({ title: "Error", description: "Failed to create folder on backend.", variant: "destructive" });
		}
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
				restoreEntry,
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

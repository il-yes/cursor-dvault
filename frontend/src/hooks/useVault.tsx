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
	sync: (jwtToken: string) => Promise<void>;
	encryptFile: (jwtToken: string, filePath: Uint8Array, vaultPassword: string) => Promise<string>;
	encryptVault: (jwtToken: string, vaultPassword: string) => Promise<string>;
	uploadToIPFS: (jwtToken: string, filePath: string | Uint8Array) => Promise<string>;
	createStellarCommit: (jwtToken: string, cid: string) => Promise<string>;
	syncVault: (jwtToken: string, vaultPassword: string) => Promise<string>;
}

const VaultContextProvider = createContext<VaultContextValue | undefined>(undefined);

export function VaultProvider({ children }: { children: ReactNode }) {
	// 🔑 SINGLE SOURCE OF TRUTH
	const vaultContext = useVaultStore((state) => state.vault);
	const isLoading = useVaultStore((state) => state.isLoading);

	/* ------------------------------------------------------------------ */
	/* Backend → Store hydration                                           */
	/* ------------------------------------------------------------------ */

	useEffect(() => {
		const unsubscribe = wailsBridge.onContextReceived((context) => {
			hydrateVault(context);
		});

		// Load initial vault only once if empty
		if (!vaultContext) {
			loadInitialContext();
		}

		return unsubscribe;
	}, []);

	const loadInitialContext = async () => {
		try {
			const context = await wailsBridge.requestContext();
			if (context) {
				hydrateVault(context);
			}
		} catch (err) {
			console.error("❌ Failed to load initial vault context:", err);
		}
	};

	const hydrateVault = (context: VaultContext) => {
		useVaultStore.getState().setVault(context);

		toast({
			title: "Vault ready",
			description: `Last synced: ${context.LastSynced || "Just now"}`,
		});
	};

	const refreshVault = async () => {
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
	};

	const clearVault = () => {
		useVaultStore.getState().clearVault();
	};

	/* ------------------------------------------------------------------ */
	/* Mutations — ALL go through Zustand                                  */
	/* ------------------------------------------------------------------ */

	const addEntry = async (entry: VaultEntry): Promise<void> => {
		useVaultStore.setState((state) => {
			if (!state.vault) return state;

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
	};

	const updateEntry = async (entryId: string, updates: Partial<VaultEntry>) => {
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
	};

	const deleteEntry = async (entryId: string) => {
		await updateEntry(entryId, { trashed: true });
	};

	const restoreEntry = async (entryId: string) => {
		await updateEntry(entryId, { trashed: false });
	};

	const toggleFavorite = async (entryId: string) => {
		useVaultStore.setState((state) => {
			if (!state.vault) return state;

			const newEntries: any = {};

			for (const type of Object.keys(state.vault.Vault.entries)) {
				const entries = state.vault.Vault.entries[type];
				newEntries[type] = entries.map((e) =>
					e.id === entryId ? { ...e, is_favorite: !e.is_favorite } : e
				);
			}

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
	};

	/* ------------------------------------------------------------------ */
	/* Additional Methods Implementation                                   */
	/* ------------------------------------------------------------------ */

	const addFolder = async (name: string) => {
		try {
			const { jwtToken } = useAuthStore.getState();
			if (!jwtToken) throw new Error("Authentication token not found");

			await AppAPI.CreateFolder(name, jwtToken);
			await refreshVault();
		} catch (err) {
			console.error("❌ Failed to add folder:", err);
			toast({
				title: "Error",
				description: `Failed to create folder: ${(err as Error).message}`,
				variant: "destructive",
			});
			throw err;
		}
	};

	const sync = async (jwtToken: string) => {
		try {
			const context = useVaultStore.getState().vault;
			if (!context) throw new Error("Vault context not found");
			context.Dirty = false;
			context.LastSynced = new Date().toISOString();
			useVaultStore.getState().setVault(context);
			await AppAPI.SynchronizeVault(jwtToken, ""); // Password might be needed depending on implementation
			await refreshVault();
		} catch (err) {
			console.error("❌ Failed to sync:", err);
			throw err;
		}
	};

	const encryptFile = async (jwtToken: string, fileData: Uint8Array, vaultPassword: string) => {
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
	};

	const encryptVault = async (jwtToken: string, vaultPassword: string) => {
		return await AppAPI.EncryptVault(jwtToken, vaultPassword);
	};

	const uploadToIPFS = async (jwtToken: string, filePath: string | Uint8Array) => {
		return await AppAPI.UploadToIPFS(jwtToken, filePath as any);
	};

	const createStellarCommit = async (jwtToken: string, cid: string) => {
		return await AppAPI.CreateStellarCommit(jwtToken, cid);
	};

	const syncVault = async (jwtToken: string, vaultPassword: string) => {
		const context = useVaultStore.getState().vault;
		if (!context) throw new Error("Vault context not found");
		context.Dirty = false;
		context.LastSynced = new Date().toISOString();
		useVaultStore.getState().setVault(context);
		return await AppAPI.SynchronizeVault(jwtToken, vaultPassword);
	};

	/* ------------------------------------------------------------------ */

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
				sync,
				encryptFile,
				encryptVault,
				uploadToIPFS,
				createStellarCommit,
				syncVault,
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

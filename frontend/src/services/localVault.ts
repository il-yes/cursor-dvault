/**
 * Local vault persistence using IndexedDB
 * Fallback for offline scenarios
 */

import { VaultContext } from '@/types/vault';

const DB_NAME = 'SovereignVault';
const DB_VERSION = 1;
const STORE_NAME = 'vault_drafts';

interface VaultDraft extends VaultContext {
  draft_id: string;
  saved_at: string;
  synced: boolean;
}

class LocalVaultService {
  private db: IDBDatabase | null = null;

  async init(): Promise<void> {
    return new Promise((resolve, reject) => {
      const request = indexedDB.open(DB_NAME, DB_VERSION);

      request.onerror = () => reject(request.error);
      request.onsuccess = () => {
        this.db = request.result;
        resolve();
      };

      request.onupgradeneeded = (event) => {
        const db = (event.target as IDBOpenDBRequest).result;
        if (!db.objectStoreNames.contains(STORE_NAME)) {
          db.createObjectStore(STORE_NAME, { keyPath: 'draft_id' });
        }
      };
    });
  }

  async saveDraft(vaultContext: Partial<VaultContext>): Promise<string> {
    await this.ensureDB();
    
    const draft: VaultDraft = {
      ...vaultContext as VaultContext,
      draft_id: `draft_${Date.now()}`,
      saved_at: new Date().toISOString(),
      synced: false,
      Dirty: true,
    };

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([STORE_NAME], 'readwrite');
      const store = transaction.objectStore(STORE_NAME);
      const request = store.add(draft);

      request.onsuccess = () => resolve(draft.draft_id);
      request.onerror = () => reject(request.error);
    });
  }

  async loadDraft(draftId: string): Promise<VaultDraft | null> {
    await this.ensureDB();

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([STORE_NAME], 'readonly');
      const store = transaction.objectStore(STORE_NAME);
      const request = store.get(draftId);

      request.onsuccess = () => resolve(request.result || null);
      request.onerror = () => reject(request.error);
    });
  }

  async getAllDrafts(): Promise<VaultDraft[]> {
    await this.ensureDB();

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([STORE_NAME], 'readonly');
      const store = transaction.objectStore(STORE_NAME);
      const request = store.getAll();

      request.onsuccess = () => resolve(request.result || []);
      request.onerror = () => reject(request.error);
    });
  }

  async markSynced(draftId: string): Promise<void> {
    await this.ensureDB();

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([STORE_NAME], 'readwrite');
      const store = transaction.objectStore(STORE_NAME);
      const getRequest = store.get(draftId);

      getRequest.onsuccess = () => {
        const draft = getRequest.result;
        if (draft) {
          draft.synced = true;
          const updateRequest = store.put(draft);
          updateRequest.onsuccess = () => resolve();
          updateRequest.onerror = () => reject(updateRequest.error);
        } else {
          reject(new Error('Draft not found'));
        }
      };
      getRequest.onerror = () => reject(getRequest.error);
    });
  }

  async deleteDraft(draftId: string): Promise<void> {
    await this.ensureDB();

    return new Promise((resolve, reject) => {
      const transaction = this.db!.transaction([STORE_NAME], 'readwrite');
      const store = transaction.objectStore(STORE_NAME);
      const request = store.delete(draftId);

      request.onsuccess = () => resolve();
      request.onerror = () => reject(request.error);
    });
  }

  async exportDraft(draftId: string): Promise<Blob> {
    const draft = await this.loadDraft(draftId);
    if (!draft) {
      throw new Error('Draft not found');
    }

    const json = JSON.stringify(draft, null, 2);
    return new Blob([json], { type: 'application/json' });
  }

  private async ensureDB(): Promise<void> {
    if (!this.db) {
      await this.init();
    }
  }
}

export const localVaultService = new LocalVaultService();

/**
 * API Service for D-Vault
 * 
 * This service provides a modular interface to connect with the Golang backend.
 * All IPFS, Stellar, and Tracecore operations are handled by the backend.
 * 
 * Backend Endpoints (to be implemented):
 * 
 * POST   /api/vault/entries          - Create new vault entry
 * GET    /api/vault/entries          - List all vault entries
 * GET    /api/vault/entries/:id      - Get specific entry
 * PUT    /api/vault/entries/:id      - Update entry
 * DELETE /api/vault/entries/:id      - Delete entry
 * POST   /api/vault/entries/:id/share - Share entry (generate access token)
 * 
 * POST   /api/ipfs/upload            - Upload to IPFS
 * GET    /api/ipfs/:cid              - Retrieve from IPFS
 * 
 * POST   /api/stellar/anchor          - Anchor hash to Stellar
 * GET    /api/stellar/verify/:tx     - Verify Stellar transaction
 * 
 * POST   /api/tracecore/commit       - Create Tracecore commit
 * GET    /api/tracecore/verify/:id   - Verify commit integrity
 */

export interface VaultEntry {
  id: string;
  title: string;
  content: string;
  category: string;
  ipfsHash: string;
  stellarTxHash: string;
  tracecoreCommitId: string;
  createdAt: string;
  updatedAt: string;
  encrypted: boolean;
}

export interface CreateEntryPayload {
  title: string;
  content: string;
  category: string;
}

export interface ShareEntryResponse {
  shareUrl: string;
  expiresAt: string;
}

// Backend API base URL (configure based on environment)
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

/**
 * Mock data for preview/development mode
 * Replace with actual API calls when backend is connected
 */
const MOCK_ENTRIES: VaultEntry[] = [
  {
    id: "1",
    title: "Personal Identity Documents",
    content: "Passport, Driver's License, Birth Certificate",
    category: "Identity",
    ipfsHash: "QmX7fKRxC3wvJ8PyF4bEqLHxQzPnZ9KuLpYvWmCdRqX1aZ",
    stellarTxHash: "0xstellar123...abc",
    tracecoreCommitId: "tc_commit_001",
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    encrypted: true,
  },
  {
    id: "2",
    title: "Financial Records Q1 2025",
    content: "Tax documents, bank statements, investment records",
    category: "Finance",
    ipfsHash: "QmY8gLRyD4xwK9QzG5cFsMJyRqYwXnDpSqY2bVnEsY2bC",
    stellarTxHash: "0xstellar456...def",
    tracecoreCommitId: "tc_commit_002",
    createdAt: new Date(Date.now() - 86400000).toISOString(),
    updatedAt: new Date(Date.now() - 86400000).toISOString(),
    encrypted: true,
  },
];

/**
 * List all vault entries
 */
export async function listEntries(): Promise<VaultEntry[]> {
  try {
    // TODO: Replace with actual API call
    // const response = await fetch(`${API_BASE_URL}/api/vault/entries`);
    // const data = await response.json();
    // return data.entries;
    
    // Mock implementation for preview
    return new Promise((resolve) => {
      setTimeout(() => resolve(MOCK_ENTRIES), 500);
    });
  } catch (error) {
    console.error('Failed to list entries:', error);
    throw error;
  }
}

/**
 * Get a specific vault entry by ID
 */
export async function getEntry(id: string): Promise<VaultEntry> {
  try {
    // TODO: Replace with actual API call
    // const response = await fetch(`${API_BASE_URL}/api/vault/entries/${id}`);
    // const data = await response.json();
    // return data.entry;
    
    // Mock implementation
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        const entry = MOCK_ENTRIES.find(e => e.id === id);
        if (entry) {
          resolve(entry);
        } else {
          reject(new Error('Entry not found'));
        }
      }, 300);
    });
  } catch (error) {
    console.error('Failed to get entry:', error);
    throw error;
  }
}

/**
 * Create a new vault entry
 * Backend handles encryption, IPFS upload, Stellar anchoring, and Tracecore commit
 */
export async function createEntry(payload: CreateEntryPayload): Promise<VaultEntry> {
  try {
    // TODO: Replace with actual API call
    // const response = await fetch(`${API_BASE_URL}/api/vault/entries`, {
    //   method: 'POST',
    //   headers: { 'Content-Type': 'application/json' },
    //   body: JSON.stringify(payload),
    // });
    // const data = await response.json();
    // return data.entry;
    
    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        const newEntry: VaultEntry = {
          id: Math.random().toString(36).substr(2, 9),
          ...payload,
          ipfsHash: `Qm${Math.random().toString(36).substr(2, 44)}`,
          stellarTxHash: `0xstellar${Math.random().toString(36).substr(2, 10)}`,
          tracecoreCommitId: `tc_commit_${Math.random().toString(36).substr(2, 8)}`,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          encrypted: true,
        };
        MOCK_ENTRIES.push(newEntry);
        resolve(newEntry);
      }, 800);
    });
  } catch (error) {
    console.error('Failed to create entry:', error);
    throw error;
  }
}

/**
 * Update an existing vault entry
 */
export async function updateEntry(id: string, payload: Partial<CreateEntryPayload>): Promise<VaultEntry> {
  try {
    // TODO: Replace with actual API call
    // const response = await fetch(`${API_BASE_URL}/api/vault/entries/${id}`, {
    //   method: 'PUT',
    //   headers: { 'Content-Type': 'application/json' },
    //   body: JSON.stringify(payload),
    // });
    // const data = await response.json();
    // return data.entry;
    
    // Mock implementation
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        const index = MOCK_ENTRIES.findIndex(e => e.id === id);
        if (index !== -1) {
          MOCK_ENTRIES[index] = {
            ...MOCK_ENTRIES[index],
            ...payload,
            updatedAt: new Date().toISOString(),
          };
          resolve(MOCK_ENTRIES[index]);
        } else {
          reject(new Error('Entry not found'));
        }
      }, 600);
    });
  } catch (error) {
    console.error('Failed to update entry:', error);
    throw error;
  }
}

/**
 * Delete a vault entry
 */
export async function deleteEntry(id: string): Promise<void> {
  try {
    // TODO: Replace with actual API call
    // await fetch(`${API_BASE_URL}/api/vault/entries/${id}`, {
    //   method: 'DELETE',
    // });
    
    // Mock implementation
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        const index = MOCK_ENTRIES.findIndex(e => e.id === id);
        if (index !== -1) {
          MOCK_ENTRIES.splice(index, 1);
          resolve();
        } else {
          reject(new Error('Entry not found'));
        }
      }, 400);
    });
  } catch (error) {
    console.error('Failed to delete entry:', error);
    throw error;
  }
}

/**
 * Share a vault entry (generate temporary access token)
 */
export async function shareEntry(id: string, expirationHours: number = 24): Promise<ShareEntryResponse> {
  try {
    // TODO: Replace with actual API call
    // const response = await fetch(`${API_BASE_URL}/api/vault/entries/${id}/share`, {
    //   method: 'POST',
    //   headers: { 'Content-Type': 'application/json' },
    //   body: JSON.stringify({ expirationHours }),
    // });
    // const data = await response.json();
    // return data;
    
    // Mock implementation
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({
          shareUrl: `https://dvault.app/shared/${id}?token=${Math.random().toString(36).substr(2)}`,
          expiresAt: new Date(Date.now() + expirationHours * 3600000).toISOString(),
        });
      }, 500);
    });
  } catch (error) {
    console.error('Failed to share entry:', error);
    throw error;
  }
}

/**
 * Vault creation types and functions
 */
export interface CreateVaultPayload {
  name: string;
  plan: "freemium" | "pro" | "organization";
  stellarPublicKey?: string;
  stellarPrivateKey?: string;
  payment?: {
    name: string;
    email: string;
    cardNumber: string;
  };
}

/**
 * Create a new vault with selected plan
 */
export async function createVault(payload: CreateVaultPayload): Promise<{ success: boolean; vaultContext?: any }> {
  try {
    // TODO: Replace with actual API call
    const response = await fetch(`${API_BASE_URL}/api/vault/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }
    
    const data = await response.json();
    return { success: true, vaultContext: data.vault_context };
    
    // Mock implementation - simulate API call
    // return new Promise((resolve) => {
    //   setTimeout(() => {
    //     console.log('Creating vault:', payload.name, 'with plan:', payload.plan);
    //     resolve({ success: true });
    //   }, 1500);
    // });
  } catch (error) {
    console.error('Failed to create vault:', error);
    throw error;
  }
}

/**
 * Stellar account setup
 */
export async function setupStellarAccount(): Promise<{ publicKey: string; privateKey: string }> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/stellar/setup`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
    });
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Failed to setup Stellar account:', error);
    throw error;
  }
}

/**
 * Decrypt a sensitive field
 */
export async function decryptField(payload: { entry_id: string; field_name: string; challenge?: string }): Promise<{ plaintext: string; expires_in: number }> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/entry/decrypt`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Failed to decrypt field:', error);
    throw error;
  }
}

/**
 * Log audit event
 */
export async function logAuditEvent(event: {
  event_type: 'view' | 'decrypt' | 'create' | 'update' | 'delete';
  entry_id: string;
  field_name?: string;
  timestamp: string;
  user_id: string;
}): Promise<void> {
  try {
    await fetch(`${API_BASE_URL}/api/audit/log`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(event),
    });
  } catch (error) {
    console.error('Failed to log audit event:', error);
  }
}

/**
 * Upgrade vault plan
 */
export async function upgradeVaultPlan(payload: {
  plan: "pro" | "organization";
  payment: {
    name: string;
    email: string;
    cardNumber: string;
  };
}): Promise<{ success: boolean }> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/payment/upgrade`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Failed to upgrade plan:', error);
    throw error;
  }
}

/**
 * Get full vault context (entries + shared entries + runtime context)
 */
export async function getVaultContext(): Promise<any> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/vault`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
    });
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Failed to get vault context:', error);
    throw error;
  }
}

/**
 * Create a new shared entry
 */
export async function createSharedEntry(payload: {
  entry_id: string;
  recipients: string[];
  permission: 'read' | 'edit' | 'temporary';
  expiration_date?: string;
  custom_message?: string;
}): Promise<any> {
  try {
    console.log("Frontend origin =", window.location.origin);

    const response = await fetch(`http://localhost:4001/api/shares`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(payload),
    });
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Failed to create shared entry:', error);
    throw error;
  }
}
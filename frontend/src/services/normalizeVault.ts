// src/utils/normalizeVault.ts
export function normalizePreloadedVault(data: any) {
  // Backend returns User, Vault, vault_runtime_context, dirty, last_cid, Tokens...
  const User = data.User || data.user;
  const Vault = data.Vault || data.vault || {};
  const vault_runtime_context = data.vault_runtime_context || data.Vault?.vault_runtime_context || {};
  const last_cid = data.last_cid || data.LastCID || data.LastCid || data.LastCid || 'main';
  const dirty = typeof data.dirty !== 'undefined' ? data.dirty : data.Dirty || false;
  const SharedEntries = data.SharedEntries || data.sharedEntries || [];

  // Ensure entries keys exist and are arrays
  const entries = Vault.entries || {};
  const normalizedEntries: Record<string, any[]> = {};
  const expectedTypes = ['login', 'card', 'identity', 'note', 'sshkey'];

  // keep any other types present too
  Object.keys(entries).forEach((k) => {
    const v = entries[k];
    // if null -> convert to []
    normalizedEntries[k] = Array.isArray(v) ? v : v ? (Array.isArray(v) ? v : [v]) : [];
  });

  // ensure expected types exist
  expectedTypes.forEach((t) => {
    if (!normalizedEntries[t]) normalizedEntries[t] = [];
  });

  const normalizedVault = {
    User,
    Vault: {
      ...Vault,
      entries: normalizedEntries,
      folders: Vault.folders || [],
      version: Vault.version || '1.0.0',
      name: Vault.name || '',
      created_at: Vault.created_at,
      updated_at: Vault.updated_at,
    },
    vault_runtime_context,
    last_cid,
    dirty,
    SharedEntries,
    Tokens: data.Tokens || data.tokens || {},
    CloudTokens: data.cloud_token || data.cloud_tokens || {},
  };

  return normalizedVault;
}

/**
 * Vault types matching backend VaultContext structure
 */

export interface StellarAccount {
  public_key: string;
  private_key?: string; // Only in dev mode
}

export interface CurrentUser {
  id: string;
  role: string;
  name: string;
  last_name: string;
  email: string;
  stellar_account: StellarAccount;
}

export interface BlockchainConfig {
  stellar: {
    network: string;
    horizon_url: string;
  };
  ipfs: {
    gateway_url: string;
    pinning_service?: string;
  };
}
export interface AppSettings {
  id: string;
  repo_id: string;
  branch: string;
  tracecore_enabled: boolean;
  commit_rules: CommitRule[];
  branching_model: string;
  encryption_policy: string;
  actors: string[];
  federated_providers: FederatedProvider[] | null;
  default_phase: string;
  default_vault_path: string;
  vault_settings: VaultSettings;
  blockchain: BlockchainConfig2;
  user_id?: number;
  auto_sync_enabled: boolean;
}

export interface CommitRule {
  id: number;
  rule: string;
  actors: string[] | null;
}

export interface VaultSettings {
  max_entries: number;
  encryption_scheme: string;
}

export interface BlockchainConfig2 {
  stellar: StellarConfig;
  ipfs: IPFSConfig;
}

export interface StellarConfig {
  network: string;
  horizon_url: string;
  fee: number;
}

export interface IPFSConfig {
  api_endpoint: string;
  gateway_url: string;
}

export interface FederatedProvider {
  name: string;
  client_id: string;
  client_secret?: string;
  scopes?: string[];
}


// Base Entry matching Go BaseEntry struct
export interface BaseEntry {
  id: string;
  entry_name: string;
  folder_id?: string;
  type: 'login' | 'card' | 'note' | 'sshkey' | 'identity';
  additionnal_note?: string;
  custom_fields?: Record<string, any>;
  trashed: boolean;
  is_draft: boolean;
  created_at: string;
  updated_at: string;
  is_favorite?: boolean;
}

// Login Entry
export interface LoginEntry extends BaseEntry {
  type: 'login';
  user_name: string;
  password: string;
  web_site?: string;
}

// Card Entry
export interface CardEntry extends BaseEntry {
  type: 'card';
  owner: string;
  number: string;
  expiration: string;
  cvc: string;
}

// Identity Entry
export interface IdentityEntry extends BaseEntry {
  type: 'identity';
  genre?: string;
  firstname?: string;
  second_firstname?: string;
  lastname?: string;
  username?: string;
  company?: string;
  social_security_number?: string;
  ID_number?: string;
  driver_license?: string;
  mail?: string;
  telephone?: string;
  address_one?: string;
  address_two?: string;
  address_three?: string;
  city?: string;
  state?: string;
  postal_code?: string;
  country?: string;
}

// Note Entry
export interface NoteEntry extends BaseEntry {
  type: 'note';
}

// SSH Key Entry
export interface SSHKeyEntry extends BaseEntry {
  type: 'sshkey';
  private_key: string;
  public_key: string;
  e_fingerprint: string;
}

// Union type for all entries
export type VaultEntry = LoginEntry | CardEntry | IdentityEntry | NoteEntry | SSHKeyEntry;

export interface Folder {
  id: string;
  name: string;
  icon?: string;
  parent_id?: string;
}

export interface Vault {
  version: string;
  name: string;
  folders: Folder[];
  entries: {
    login: VaultEntry[];
    card: VaultEntry[];
    note: VaultEntry[];
    sshkey: VaultEntry[];
    identity: VaultEntry[];
  };
  created_at?: string;
  updated_at?: string;

}

export interface VaultRuntimeContext {
  CurrentUser: CurrentUser;
  AppSettings: AppSettings;
  WorkingBranch: string;
  LoadedEntries: string[];
}

export interface VaultContext {
  user_id: string;
  role: string;
  Vault: Vault;
  LastCID?: string;
  Dirty: boolean;
  LastSynced?: string;
  LastUpdated: string;
  vault_runtime_context: VaultRuntimeContext;
}

export interface DecryptRequest {
  entry_id: string;
  field_name: string;
  challenge?: string;
}

export interface DecryptResponse {
  field_name: string;
  plaintext: string;
  expires_in: number; // seconds
}

export interface AuditEvent {
  event_type: 'view' | 'decrypt' | 'create' | 'update' | 'delete';
  entry_id: string;
  field_name?: string;
  timestamp: string;
  user_id: string;
}

export interface VaultPayload {
  version: string;
  name: string;
  folders: Folder[];
  entries: {
    login: Entry[];      // 👈 Separate out types
    card: Entry[];
    note: Entry[];
    identity: Entry[]
    sshkey: Entry[]
  };
  created_at: string;
  updated_at: string;
}
export interface LoginResponse {
  User: User;
  Vault: VaultPayload;
}

export interface Folder {
  id: string;
  name: string;
  type?: string;
}

export interface LoginRequest {
  email?: string;
  password?: string;
  publicKey?: string;
  signedMessage?: string;
  signature?: string;
}
export interface User {
  id: string;
  email: string;
  password: string;
  username: string;
  role: string;
}
export type Entry =
  | LoginEntry
  | CardEntry
  | IdentityEntry
  | NoteEntry
  | SSHKeyEntry;


export type WailsResponse<T> = {
  result: T | null;
  error: string | null;
  callbackid: string;
};
export interface PreloadedVaultResponse {
  User: {
    id: string | number;
    role: string;
    email?: string;
    username?: string;
  };
  Vault: any; // You can type strictly if you want
  Tokens?: {
    access_token: string;
    refresh_token: string;
  };
  SharedEntries?: any[];
  VaultRuntimeContext?: any;
  LastCID?: string;
  Dirty?: boolean;
}

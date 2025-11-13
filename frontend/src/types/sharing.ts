export interface Recipient {
  id: string;
  name: string;
  email: string;
  role: "viewer" | "editor" | "owner";
  joined_at: string;
}

export interface AuditEvent {
  id: string;
  action: string;
  actor: string;
  timestamp: string;
  details?: string;
}

export interface SharedEntry {
  id: string;
  entry_name: string;
  entry_type: "login" | "card" | "note" | "identity" | "sshkey";
  description?: string;
  folder?: string;
  status: "active" | "expired" | "revoked" | "pending";
  access_mode: "read" | "edit";
  recipients: Recipient[];
  shared_at: string;
  expires_at?: string;
  encryption: "AES-256-GCM";
  blockchain_hash?: string;
  ipfs_anchor?: string;
  created_at: string;
  updated_at: string;
  audit_log: AuditEvent[];
}

export type ShareFilter = "all" | "sent" | "received" | "pending" | "revoked";
export type DetailView = "recipients" | "audit" | "encryption" | "metadata";


export interface AuditEvent {
  id: string;
  action: string;
  actor?: string;
  timestamp: string;
  details?: string;
}

export interface SharedEntryV0 {
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
export interface SharedEntry {
  id: string; // backend-generated
  entry_name: string;
  entry_type: "login" | "card" | "note" | "identity" | "sshkey";

  status: "active" | "expired" | "revoked" | "pending";
  access_mode: "read" | "edit";
  encryption: "AES-256-GCM";

  entry_snapshot: EntrySnapshot;

  recipients: Recipient[]; // only name, email, role when sending

  shared_at: string;
  expires_at?: string | null;

  created_at: string;
  updated_at: string;

  audit_log: AuditEvent[]; // always empty for now, backend returns nothing
}

export type ShareFilter = "all" | "sent" | "received" | "pending" | "revoked" | "withme";
export type DetailView = "recipients" | "audit" | "encryption" | "metadata";
export interface EntrySnapshot {
  entry_name: string;
  type: "login" | "card" | "note" | "identity" | "sshkey";

  user_name?: string;
  password?: string;
  website?: string;

  cardholder_name?: string;
  card_number?: string;
  expiration?: string;
  cvv?: string;

  private_key?: string;
  public_key?: string;

  note?: string;

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

  extra_fields?: Record<string, unknown>;
}
export interface Recipient {
  id: string;        // backend-generated
  share_id: string;  // backend-generated
  name: string;
  email: string;
  public_key: string;
  role: "viewer" | "read" | "editor" | "owner";
  joined_at?: string;  // backend-generated
  created_at?: string;
  updated_at?: string;
}

export interface AuditEvent {
  id: string;
  action: string;
  timestamp: string;
  actor?: string;
  metadata?: Record<string, any>;
}
export interface CreateShareEntryPayload {
  entry_name: string;
  entry_type: "login" | "card" | "note" | "identity" | "sshkey";
  status: "active" | "pending"  | "expired" | "revoked";
  access_mode: "read" | "edit";
  encryption: "AES-256-GCM";
  entry_snapshot: EntrySnapshot;
  recipients: {
    name: string;
    email: string;
    public_key: string;
    role: "viewer" | "editor";
  }[];
  expires_at?: string | null;
}


import { Folder, LoginEntry, CardEntry, IdentityEntry, NoteEntry, SSHKeyEntry, VaultEntry } from "@/types/vault";

export const mockFolders: Folder[] = [
  { id: "folder-1", name: "Stone", icon: "🪨" },
  { id: "folder-2", name: "Minerals", icon: "💎" },
  { id: "folder-3", name: "Projects", icon: "📁" },
  { id: "folder-4", name: "Personal", icon: "👤" },
];

// Mock login entries
const mockLoginEntries: LoginEntry[] = [
  {
    id: "login-1",
    type: "login",
    entry_name: "GitHub Account",
    folder_id: "folder-3",
    user_name: "developer@example.com",
    password: "••••••••••••",
    web_site: "https://github.com",
    trashed: false,
    is_draft: false,
    created_at: "2024-01-15T10:30:00Z",
    updated_at: "2024-01-20T14:45:00Z",
    is_favorite: true,
  },
  {
    id: "login-2",
    type: "login",
    entry_name: "AWS Console",
    folder_id: "folder-3",
    user_name: "admin@company.com",
    password: "••••••••••••",
    web_site: "https://aws.amazon.com",
    additionnal_note: "Production environment access",
    trashed: false,
    is_draft: false,
    created_at: "2024-01-10T09:15:00Z",
    updated_at: "2024-01-22T16:20:00Z",
    is_favorite: false,
  },
  {
    id: "login-3",
    type: "login",
    entry_name: "Personal Email",
    folder_id: "folder-4",
    user_name: "myemail@gmail.com",
    password: "••••••••••••",
    web_site: "https://mail.google.com",
    trashed: false,
    is_draft: false,
    created_at: "2024-01-05T08:00:00Z",
    updated_at: "2024-01-18T11:30:00Z",
    is_favorite: true,
  },
];

// Mock card entries
const mockCardEntries: CardEntry[] = [
  {
    id: "card-1",
    type: "card",
    entry_name: "Visa Gold Card",
    folder_id: "folder-4",
    owner: "John Doe",
    number: "•••• •••• •••• 1234",
    expiration: "12/26",
    cvc: "•••",
    trashed: false,
    is_draft: false,
    created_at: "2024-01-12T12:00:00Z",
    updated_at: "2024-01-12T12:00:00Z",
    is_favorite: false,
  },
  {
    id: "card-2",
    type: "card",
    entry_name: "Business Mastercard",
    folder_id: "folder-3",
    owner: "Jane Smith",
    number: "•••• •••• •••• 5678",
    expiration: "08/25",
    cvc: "•••",
    additionnal_note: "Company expenses only",
    trashed: false,
    is_draft: false,
    created_at: "2024-01-08T14:30:00Z",
    updated_at: "2024-01-19T10:15:00Z",
    is_favorite: true,
  },
];

// Mock note entries
const mockNoteEntries: NoteEntry[] = [
  {
    id: "note-1",
    type: "note",
    entry_name: "Project Ideas",
    folder_id: "folder-3",
    additionnal_note: "1. Implement AI-powered search\n2. Add blockchain verification\n3. Enhance mobile responsiveness\n4. Create API documentation",
    trashed: false,
    is_draft: false,
    created_at: "2024-01-14T16:45:00Z",
    updated_at: "2024-01-21T09:30:00Z",
    is_favorite: false,
  },
  {
    id: "note-2",
    type: "note",
    entry_name: "Recovery Codes",
    folder_id: "folder-4",
    additionnal_note: "Backup codes for 2FA:\n1. ABC-DEF-GHI\n2. JKL-MNO-PQR\n3. STU-VWX-YZ1",
    trashed: false,
    is_draft: false,
    created_at: "2024-01-09T11:20:00Z",
    updated_at: "2024-01-09T11:20:00Z",
    is_favorite: true,
  },
];

// Mock SSH key entries
const mockSSHKeyEntries: SSHKeyEntry[] = [
  {
    id: "sshkey-1",
    type: "sshkey",
    entry_name: "Production Server",
    folder_id: "folder-3",
    private_key: "••••••••••••••••••••",
    public_key: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC...",
    e_fingerprint: "SHA256:nThbg6kXUpJWGl7E1IGOCspRomTxdCARLviKw6E5SY8",
    additionnal_note: "Main production server access",
    trashed: false,
    is_draft: false,
    created_at: "2024-01-11T13:00:00Z",
    updated_at: "2024-01-11T13:00:00Z",
    is_favorite: true,
  },
  {
    id: "sshkey-2",
    type: "sshkey",
    entry_name: "Development Environment",
    folder_id: "folder-3",
    private_key: "••••••••••••••••••••",
    public_key: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm...",
    e_fingerprint: "SHA256:4Vw6fK2zLx5bHZRr8t9p3qWe7yUiOkNhGfDsSaQpLm1",
    trashed: false,
    is_draft: false,
    created_at: "2024-01-13T15:30:00Z",
    updated_at: "2024-01-20T12:00:00Z",
    is_favorite: false,
  },
];

// Mock identity entries
const mockIdentityEntries: IdentityEntry[] = [
  {
    id: "identity-1",
    type: "identity",
    entry_name: "Personal Identity",
    folder_id: "folder-4",
    genre: "Male",
    firstname: "John",
    lastname: "Doe",
    username: "johndoe",
    mail: "john.doe@example.com",
    telephone: "+1-555-0123",
    address_one: "123 Main Street",
    city: "San Francisco",
    state: "California",
    postal_code: "94102",
    country: "United States",
    trashed: false,
    is_draft: false,
    created_at: "2024-01-07T10:00:00Z",
    updated_at: "2024-01-16T14:30:00Z",
    is_favorite: false,
  },
  {
    id: "identity-2",
    type: "identity",
    entry_name: "Business Identity",
    folder_id: "folder-3",
    firstname: "Jane",
    lastname: "Smith",
    company: "TechCorp Inc.",
    username: "janesmith",
    mail: "jane.smith@techcorp.com",
    telephone: "+1-555-9876",
    address_one: "456 Business Ave",
    address_two: "Suite 200",
    city: "New York",
    state: "New York",
    postal_code: "10001",
    country: "United States",
    trashed: false,
    is_draft: false,
    created_at: "2024-01-06T09:00:00Z",
    updated_at: "2024-01-17T11:45:00Z",
    is_favorite: true,
  },
];

export const mockLinkShares = [
  {
    id: "1",
    entry_name: "aws keys",
    status: "active", // or "expired", "revoked"
    expiry: "2026-01-28",
    uses_left: 2,
    link: "https://ankhora.app/share/abcdef123456",
    audit_log: [/* ... */],
  },
  {
    id: "2",
    entry_name: "google keys",
    status: "active", // or "expired", "revoked"
    expiry: "2026-01-28",
    uses_left: 2,
    link: "https://ankhora.app/share/abcdef123456",
    audit_log: [/* ... */],
  },
  {
    id: "3",
    entry_name: "github keys",
    status: "expired", // or "expired", "revoked"
    expiry: "2025-12-28",
    uses_left: 0,
    link: "https://ankhora.app/share/abcdef123456",
    audit_log: [/* ... */],
  },
  // ...
];


export const getMockVaultEntries = () => ({
  login: mockLoginEntries as VaultEntry[],
  card: mockCardEntries as VaultEntry[],
  note: mockNoteEntries as VaultEntry[],
  sshkey: mockSSHKeyEntries as VaultEntry[],
  identity: mockIdentityEntries as VaultEntry[],
});

export const getAllMockEntries = (): VaultEntry[] => [
  ...mockLoginEntries,
  ...mockCardEntries,
  ...mockNoteEntries,
  ...mockSSHKeyEntries,
  ...mockIdentityEntries,
];

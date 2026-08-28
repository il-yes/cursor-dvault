import {
  C3TrustGroupSummary,
  C3ChannelSummary,
  C3CollaborationConfiguration,
  C3SecurityState,
  C3UseCaseTemplate,
  C3CollaborationMetrics,
} from "../domain/configuration";

/**
 * Isolated Mock Configuration Provider Data
 *
 * This module is isolated from presentation components and service layers.
 * In Phase 2, this provider interface will be substituted with Wails RPC adapters.
 */

export const mockTrustGroups: C3TrustGroupSummary[] = [
  {
    id: "tg_91a0c4f1",
    name: "Engineering Operations & Architecture",
    description: "Cryptographic trust boundary for technical leads, system architects, and release reviewers.",
    memberCount: 8,
    deviceCount: 12,
    kekVersion: 4,
    status: "active",
    createdDate: "2026-01-15T09:30:00Z",
    associatedChannelIds: ["ch_881a029f", "ch_110b99c4"],
    cryptographicState: {
      kekVersion: 4,
      envelopeCount: 16,
      activeEnvelopes: 12,
      revokedEnvelopes: 4,
      lastRotation: "2026-08-10T14:22:00Z",
    },
    members: [
      {
        id: "usr_lead_01",
        identityName: "Alice Vance (Principal Architect)",
        email: "alice.vance@ankhora.io",
        role: "admin",
        status: "active",
        deviceCount: 2,
        joinedAt: "2026-01-15T09:30:00Z",
      },
      {
        id: "usr_eng_02",
        identityName: "Bob Miller (Core Engineer)",
        email: "bob.miller@ankhora.io",
        role: "member",
        status: "active",
        deviceCount: 1,
        joinedAt: "2026-02-01T10:15:00Z",
      },
      {
        id: "usr_eng_03",
        identityName: "Carol Danvers (Security Auditor)",
        email: "carol.danvers@ankhora.io",
        role: "auditor",
        status: "active",
        deviceCount: 2,
        joinedAt: "2026-03-12T11:00:00Z",
      },
    ],
    devices: [
      {
        id: "dev_macbook_pro_m3",
        deviceName: "Alice's MacBook Pro M3 Max",
        identityId: "usr_lead_01",
        status: "active",
        lastSeenAt: "2026-08-28T13:45:00Z",
        keyState: "valid",
      },
      {
        id: "dev_linux_workstation",
        deviceName: "Bob's Arch Linux Workstation",
        identityId: "usr_eng_02",
        status: "active",
        lastSeenAt: "2026-08-28T12:10:00Z",
        keyState: "valid",
      },
      {
        id: "dev_sec_yubikey",
        deviceName: "Carol's Hardware Security Device",
        identityId: "usr_eng_03",
        status: "active",
        lastSeenAt: "2026-08-27T18:00:00Z",
        keyState: "rotation_required",
      },
    ],
  },
  {
    id: "tg_33f8b002",
    name: "Legal Counsel & Regulatory Compliance",
    description: "Sovereign access boundary for legal advisors, compliance managers, and external auditors.",
    memberCount: 5,
    deviceCount: 7,
    kekVersion: 2,
    status: "active",
    createdDate: "2026-02-10T11:00:00Z",
    associatedChannelIds: ["ch_44921f00"],
    cryptographicState: {
      kekVersion: 2,
      envelopeCount: 8,
      activeEnvelopes: 7,
      revokedEnvelopes: 1,
      lastRotation: "2026-06-01T09:00:00Z",
    },
    members: [
      {
        id: "usr_legal_01",
        identityName: "David Wright (General Counsel)",
        email: "david.wright@ankhora.io",
        role: "admin",
        status: "active",
        deviceCount: 2,
        joinedAt: "2026-02-10T11:00:00Z",
      },
      {
        id: "usr_legal_02",
        identityName: "Elena Rostova (Compliance Officer)",
        email: "elena.rostova@ankhora.io",
        role: "member",
        status: "active",
        deviceCount: 1,
        joinedAt: "2026-02-15T14:30:00Z",
      },
    ],
    devices: [
      {
        id: "dev_david_ipad",
        deviceName: "David's iPad Pro (Secure Air)",
        identityId: "usr_legal_01",
        status: "active",
        lastSeenAt: "2026-08-28T10:00:00Z",
        keyState: "valid",
      },
      {
        id: "dev_elena_thinkpad",
        deviceName: "Elena's ThinkPad X1 Carbon",
        identityId: "usr_legal_02",
        status: "active",
        lastSeenAt: "2026-08-28T09:15:00Z",
        keyState: "valid",
      },
    ],
  },
  {
    id: "tg_7718e910",
    name: "Executive Review & Governance",
    description: "High-sovereignty sign-off boundary for board members and executive stakeholders.",
    memberCount: 4,
    deviceCount: 5,
    kekVersion: 1,
    status: "pending",
    createdDate: "2026-07-20T16:00:00Z",
    associatedChannelIds: [],
    cryptographicState: {
      kekVersion: 1,
      envelopeCount: 5,
      activeEnvelopes: 5,
      revokedEnvelopes: 0,
      lastRotation: "2026-07-20T16:00:00Z",
    },
    members: [
      {
        id: "usr_exec_01",
        identityName: "Frank Sterling (Chief Executive)",
        email: "frank.sterling@ankhora.io",
        role: "admin",
        status: "active",
        deviceCount: 2,
        joinedAt: "2026-07-20T16:00:00Z",
      },
    ],
    devices: [
      {
        id: "dev_frank_macbook",
        deviceName: "Frank's MacBook Air M2",
        identityId: "usr_exec_01",
        status: "active",
        lastSeenAt: "2026-08-26T15:30:00Z",
        keyState: "valid",
      },
    ],
  },
];

export const mockChannels: C3ChannelSummary[] = [
  {
    id: "ch_881a029f",
    name: "Architecture & Infrastructure Review",
    status: "active",
    trustGroupId: "tg_91a0c4f1",
    participantCount: 8,
    threadCount: 14,
    createdDate: "2026-01-20T10:00:00Z",
    description: "Governance surface for technical specs, RFC approvals, and key rotation notifications.",
  },
  {
    id: "ch_110b99c4",
    name: "Sprint Deliverables & Quality Gates",
    status: "active",
    trustGroupId: "tg_91a0c4f1",
    participantCount: 12,
    threadCount: 22,
    createdDate: "2026-02-05T09:00:00Z",
    description: "Daily collaboration channel for pull request approvals and release verification.",
  },
  {
    id: "ch_44921f00",
    name: "Contract Audits & Regulatory Compliance",
    status: "active",
    trustGroupId: "tg_33f8b002",
    participantCount: 6,
    threadCount: 9,
    createdDate: "2026-02-18T14:00:00Z",
    description: "Legal review channel for counter-signing contract drafts and NDAs.",
  },
  {
    id: "ch_5500c21a",
    name: "Public Announcements & General Discussion",
    status: "active",
    trustGroupId: undefined, // Explicitly no Trust Group associated!
    participantCount: 24,
    threadCount: 31,
    createdDate: "2026-01-10T08:00:00Z",
    description: "Unencrypted public communication surface for general workspace discussions.",
  },
  {
    id: "ch_9921e488",
    name: "Draft Specs (Unassigned Boundary)",
    status: "draft",
    trustGroupId: undefined, // Explicitly no Trust Group associated!
    participantCount: 3,
    threadCount: 2,
    createdDate: "2026-08-01T11:00:00Z",
    description: "Draft channel pending assignment to a sovereign Trust Group.",
  },
];

export const mockCollaborationMetrics: C3CollaborationMetrics = {
  pendingApprovals: 4,
  recentRejections: 1,
  activeTransfers: 3,
  totalShares: 18,
};

export const mockCollaborationConfig: C3CollaborationConfiguration = {
  approvalsEnabled: true,
  rejectionEnabled: true,
  transfersEnabled: true,
  sharesEnabled: true,
  allowedTransferTrustGroupIds: ["tg_91a0c4f1", "tg_33f8b002", "tg_7718e910"],
  defaultTransferMessage: "Requesting formal transfer of responsibility to designated Trust Group.",
  requireApprovalForTransfer: true,
  requireReasonForRejection: true,
};

export const mockSecurityState: C3SecurityState = {
  deviceIdentityStatus: "configured",
  sovereignKeyringStatus: "healthy",
  trustGroupKeysState: "current",
  currentKekVersion: 4,
  activeEnvelopes: 24,
  revokedEnvelopes: 5,
  lastRotation: "2026-08-10T14:22:00Z",
  localDeviceId: "dev_macbook_pro_m3",
  localVaultAddress: "vault_local_sovereign_01",
};

export const mockUseCaseTemplates: C3UseCaseTemplate[] = [
  {
    id: "tpl_generic",
    name: "Generic C3 Collaboration Framework",
    category: "Generic",
    description: "Base collaboration model providing Approval, Reject, Transfer, and C3 Share primitives across any resource type.",
    icon: "⚡",
    recommendedTrustGroups: ["tg_91a0c4f1"],
    rules: [
      {
        resourceType: "vault_entry",
        allowedActions: ["approval", "reject", "transfer", "c3_share"],
        approvalRequired: false,
        transferRule: "Transfer allowed between active member Trust Groups",
      },
    ],
  },
  {
    id: "tpl_construction",
    name: "Construction & Infrastructure Quality Assurance",
    category: "Construction",
    description: "Configures C3 primitives for structural safety sign-offs, site inspection approvals, and blueprint transfer.",
    icon: "🏗️",
    recommendedTrustGroups: ["tg_91a0c4f1"],
    rules: [
      {
        resourceType: "blueprint_document",
        allowedActions: ["approval", "reject", "transfer"],
        approvalRequired: true,
        transferRule: "Transfer from Structural Design TG to Site Engineering TG requires formal Approval",
      },
      {
        resourceType: "inspection_report",
        allowedActions: ["approval", "reject"],
        approvalRequired: true,
        transferRule: "Non-transferable until signed off by Lead Auditor",
      },
    ],
  },
  {
    id: "tpl_pharma",
    name: "Pharmaceutical Clinical Batch Verification",
    category: "Pharmaceutical",
    description: "Enforces strict dual-control approval, mandatory rejection reasons, and audit logging for trial data.",
    icon: "🧪",
    recommendedTrustGroups: ["tg_33f8b002"],
    rules: [
      {
        resourceType: "batch_record",
        allowedActions: ["approval", "reject", "transfer"],
        approvalRequired: true,
        transferRule: "Transfer to Quality Assurance TG requires unanimous Approval by Lab Directors",
      },
    ],
  },
  {
    id: "tpl_logistics",
    name: "Supply Chain & Chain of Custody",
    category: "Logistics",
    description: "Tracks physical/digital responsibility transfer between logistics hubs and receiving agents.",
    icon: "📦",
    recommendedTrustGroups: ["tg_91a0c4f1", "tg_33f8b002"],
    rules: [
      {
        resourceType: "manifest",
        allowedActions: ["transfer", "approval"],
        approvalRequired: false,
        transferRule: "Sequential transfer lifecycle: REQUESTED -> APPROVED -> COMPLETED",
      },
    ],
  },
];

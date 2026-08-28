/**
 * C3 Configuration Domain Models & Frontend Types
 *
 * Defines the authoritative information architecture for C3 Configuration,
 * Trust Groups, Channels, Collaboration Primitives, Sovereignty, and Templates.
 */

export interface C3ResourceReference {
  resourceType: string;
  resourceId: string;
  sourceEventId?: string;
}

export interface C3TrustGroupMember {
  id: string;
  identityName: string;
  email?: string;
  role: "admin" | "member" | "auditor";
  status: "active" | "pending" | "revoked";
  deviceCount: number;
  joinedAt: string;
}

export interface C3TrustGroupDevice {
  id: string;
  deviceName: string;
  identityId: string;
  status: "active" | "pending" | "revoked";
  lastSeenAt: string;
  keyState: "valid" | "rotation_required" | "revoked";
}

export interface C3CryptographicState {
  kekVersion: number;
  envelopeCount: number;
  activeEnvelopes: number;
  revokedEnvelopes: number;
  lastRotation: string;
}

export interface C3TrustGroupSummary {
  id: string;
  name: string;
  description?: string;
  memberCount: number;
  deviceCount: number;
  kekVersion: number;
  status: "active" | "pending" | "revoked";
  createdDate: string;
  members: C3TrustGroupMember[];
  devices: C3TrustGroupDevice[];
  cryptographicState: C3CryptographicState;
  associatedChannelIds: string[];
}

export interface C3ChannelSummary {
  id: string;
  name: string;
  status: "active" | "archived" | "draft";
  trustGroupId?: string;
  participantCount: number;
  threadCount: number;
  createdDate: string;
  description?: string;
}

export interface C3CollaborationMetrics {
  pendingApprovals: number;
  recentRejections: number;
  activeTransfers: number;
  totalShares: number;
}

export interface C3CollaborationConfiguration {
  approvalsEnabled: boolean;
  rejectionEnabled: boolean;
  transfersEnabled: boolean;
  sharesEnabled: boolean;
  allowedTransferTrustGroupIds: string[];
  defaultTransferMessage?: string;
  requireApprovalForTransfer: boolean;
  requireReasonForRejection: boolean;
}

export interface C3SecurityState {
  deviceIdentityStatus: "configured" | "unconfigured" | "error";
  sovereignKeyringStatus: "healthy" | "degraded" | "locked";
  trustGroupKeysState: "current" | "rotation_pending" | "compromised";
  currentKekVersion: number;
  activeEnvelopes: number;
  revokedEnvelopes: number;
  lastRotation: string;
  localDeviceId: string;
  localVaultAddress: string;
}

export interface C3UseCaseTemplate {
  id: string;
  name: string;
  category: "Generic" | "Construction" | "Pharmaceutical" | "Logistics" | "Finance";
  description: string;
  icon: string;
  recommendedTrustGroups: string[];
  rules: {
    resourceType: string;
    allowedActions: ("approval" | "reject" | "transfer" | "c3_share")[];
    approvalRequired: boolean;
    transferRule: string;
  }[];
}

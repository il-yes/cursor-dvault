export interface ResourceReference {
  resourceType: string;
  resourceId: string;
  sourceEventId?: string;
}

export const C3ResourceTypes = {
  VaultEntry: "vault_entry",
  ShareEntry: "share_entry",
  Message: "message",
  Document: "document",
  Event: "event",
} as const;

export type C3ResourceType = (typeof C3ResourceTypes)[keyof typeof C3ResourceTypes] | string;

/**
 * Desktop Resource Service
 * 
 * Central dispatcher for protected timeline resources (Classical Shares & C3 Collaborative Shares).
 * Keeps React UI completely isolated from cryptographic logic and raw Go error sentinels.
 */

import { useAuthStore } from "@/store/useAuthStore";
import { useC3ConfigurationStore } from "@/components/C3/configuration/store/useC3ConfigurationStore";
import * as AppAPI from "../../wailsjs/go/main/App";

export type ResourceKind = "share_entry" | "storage_asset";

export interface ResourceRefInput {
  refType: ResourceKind | string;
  shareEntryId?: string;
  trustGroupId?: string;
  cid?: string;
  deviceId?: string;
}

export interface OpenedResourceResult {
  resourceId: string;
  title?: string;
  content: string;
  createdBy?: string;
  createdAt?: string;
  metadata?: Record<string, string>;
  kind: ResourceKind;
}

/**
 * User-facing presentation error messages (abstracting Go error sentinels)
 */
export function translateResourceError(error: any): string {
  const errMsg = String(error?.message || error || "");

  if (errMsg.includes("share entry has been revoked") || errMsg.includes("ErrShareEntryRevoked")) {
    return "Resource no longer available";
  }
  if (errMsg.includes("not an authorized member") || errMsg.includes("ErrUnauthorizedMember")) {
    return "You no longer have access to this resource";
  }
  if (errMsg.includes("no active device key envelope") || errMsg.includes("ErrKeyEnvelopeNotFound")) {
    return "This device is not authorized to access this resource";
  }
  if (errMsg.includes("share entry not found") || errMsg.includes("ErrShareEntryNotFound")) {
    return "Resource descriptor not found";
  }
  return "Unable to open protected resource";
}

class DesktopResourceService {
  /**
   * Dispatch resource resolution to appropriate backend application service
   */
  async openResource(input: ResourceRefInput): Promise<OpenedResourceResult> {
    const jwtToken = useAuthStore.getState().jwtToken;
    const localDeviceIdFromStore = useC3ConfigurationStore.getState().securityState?.localDeviceId;
    const deviceId = input.deviceId ||
      localDeviceIdFromStore ||
      (typeof localStorage !== "undefined" ? localStorage.getItem("device_id") : null) ||
      "dev_local_01";

    if (input.refType === "share_entry" && input.shareEntryId) {
      console.log(`[C3-FORENSIC][02] desktopResourceService.openResource refType=${input.refType} shareEntryID=${input.shareEntryId}`);
      try {
        const response = await AppAPI.ResolveCollaborativeShare(jwtToken, input.shareEntryId);
        console.log(`[C3-FORENSIC][12] ResolveCollaborativeShare response received shareEntryID=${response.share_entry_id} trustGroupID=${response.trust_group_id} createdBy=${response.created_by}`);
        
        let contentText = "";
        if (typeof response.plaintext === "string") {
          contentText = response.plaintext;
        } else if (response.plaintext instanceof Uint8Array) {
          contentText = new TextDecoder().decode(response.plaintext);
        } else if (Array.isArray(response.plaintext)) {
          contentText = new TextDecoder().decode(new Uint8Array(response.plaintext));
        }

        return {
          resourceId: response.share_entry_id,
          title: response.metadata?.title || response.metadata?.name || "Collaborative Protected Resource",
          content: contentText,
          createdBy: response.created_by,
          createdAt: response.created_at,
          metadata: response.metadata || {},
          kind: "share_entry",
        };
      } catch (err: any) {
        console.error(`[C3-FORENSIC][12] ResolveCollaborativeShare IPC error=${err?.message || err}`);
        throw new Error(translateResourceError(err));
      }
    } else if (input.refType === "storage_asset" && input.shareEntryId) {
      try {
        const user = useAuthStore.getState().user;
        const result = await AppAPI.AccessDecryptVaultEntry(jwtToken, {
          share_id: input.shareEntryId,
          recipient_email: user?.email || user?.Email || "",
          challenge: "",
          signature: "",
          ip_address: "",
        });

        return {
          resourceId: input.shareEntryId,
          title: "Classical Storage Resource",
          content: result.data?.payload || "",
          kind: "storage_asset",
        };
      } catch (err) {
        throw new Error(translateResourceError(err));
      }
    }

    throw new Error("Invalid resource reference type");
  }
}

export const desktopResourceService = new DesktopResourceService();

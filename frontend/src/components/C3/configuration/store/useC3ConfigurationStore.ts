import { create } from "zustand";
import {
  C3TrustGroupSummary,
  C3ChannelSummary,
  C3CollaborationConfiguration,
  C3SecurityState,
  C3UseCaseTemplate,
  C3CollaborationMetrics,
} from "../domain/configuration";
import {
  mockTrustGroups,
  mockChannels,
  mockCollaborationConfig,
  mockSecurityState,
  mockUseCaseTemplates,
  mockCollaborationMetrics,
} from "../mock/configuration";
import {
  ChannelResponse,
  listTrustGroups,
  createTrustGroup as createTrustGroupApi,
  updateTrustGroup as updateTrustGroupApi,
  deleteTrustGroup as deleteTrustGroupApi,
} from "@/services/api";
import { useC3ChannelStore } from "../../infrastructure/store/useC3ChannelStore";
import { useC3WorkspaceStore } from "../../infrastructure/store/useC3WorkspaceStore";

export type C3ConfigurationTab =
  | "overview"
  | "trust-groups"
  | "trust"
  | "channels"
  | "slots"
  | "collaboration"
  | "assignments"
  | "security"
  | "properties"
  | "templates";


export interface SlotDraft {
  id: string;
  name: string;
  role: string;
  vault_id: string;
  gated: boolean;
  order: number;
}

export interface AssignmentDraft {
  slot_id: string;
  owner_id: string;
  public_key: string;
  vault_address: string;
  role?: "primary" | "backup" | "observer";
  since?: string;
  by?: string;
}

export interface PropertyDraft {
  key: string;
  value: string;
  type?: "String" | "Number" | "Enum" | "Boolean";
  required?: boolean;
  locked?: boolean;
  overridable?: boolean;
  defaultValue?: string;
  desc?: string;
}

export interface PolicyDraft {
  requireAllSlots: boolean;
  allowParallelCommits: boolean;
  maxThreadDuration: string;
  allowAmendments: boolean;
  requireReceiptAck: boolean;
  receiptTimeout: string;
  allowRejectionWithReason: boolean;
  retainThreadData: string;
  autoArchive: boolean;
  deleteOnArchive: boolean;
  disputeTimeout: string;
  notifyAllOnDispute: boolean;
  anchorAllCommits: boolean;
  anchorAllReceipts: boolean;
  anchorTrustGroupChanges: boolean;
}

export const defaultPolicyDraft: PolicyDraft = {
  requireAllSlots: true,
  allowParallelCommits: false,
  maxThreadDuration: "30 days",
  allowAmendments: true,
  requireReceiptAck: true,
  receiptTimeout: "72 hours",
  allowRejectionWithReason: true,
  retainThreadData: "365 days",
  autoArchive: true,
  deleteOnArchive: false,
  disputeTimeout: "14 days",
  notifyAllOnDispute: true,
  anchorAllCommits: true,
  anchorAllReceipts: true,
  anchorTrustGroupChanges: true,
};

interface C3ConfigurationState {
  activeTab: C3ConfigurationTab;
  selectedTrustGroupId: string | null;
  selectedChannelId: string | null;
  isCreateTrustGroupOpen: boolean;

  trustGroups: C3TrustGroupSummary[];
  channels: C3ChannelSummary[];
  collaborationConfig: C3CollaborationConfiguration;
  metrics: C3CollaborationMetrics;
  securityState: C3SecurityState;
  templates: C3UseCaseTemplate[];

  isLoading: boolean;
  error: string | null;

  // Unified Channel Draft State
  originalChannel: ChannelResponse | null;
  draftSlots: SlotDraft[];
  draftAssignments: AssignmentDraft[];
  draftProperties: PropertyDraft[];
  draftPolicy: PolicyDraft;
  isDirty: boolean;
  isSaving: boolean;
  saveStatus: "idle" | "saving" | "success" | "error";
  saveFeedback: string | null;

  // Actions
  setActiveTab: (tab: C3ConfigurationTab) => void;
  selectTrustGroup: (trustGroupId: string | null) => void;
  selectChannel: (channelId: string | null) => void;
  setCreateTrustGroupOpen: (open: boolean) => void;

  loadConfiguration: () => Promise<void>;
  createTrustGroup: (payload: {
    name: string;
    description?: string;
    initialMembers?: string[];
    vaultName?: string;
  }) => Promise<C3TrustGroupSummary>;
  updateTrustGroup: (payload: {
    id: string;
    name: string;
    description?: string;
  }) => Promise<C3TrustGroupSummary>;
  deleteTrustGroup: (trustGroupId: string) => Promise<void>;
  updateCollaborationConfig: (updates: Partial<C3CollaborationConfiguration>) => void;

  // Draft Mutations
  setDraftSlots: (slots: SlotDraft[]) => void;
  setDraftAssignments: (assignments: AssignmentDraft[]) => void;
  setDraftProperties: (properties: PropertyDraft[]) => void;
  setDraftPolicy: (policyUpdates: Partial<PolicyDraft>) => void;
  discardChanges: () => void;
  saveChannelConfig: (correlationId?: string) => Promise<boolean>;

}

const initialProperties: PropertyDraft[] = [
  { key: "vendor", value: "", type: "String", required: true, locked: false, overridable: true, desc: "Name of the counterparty or supplier" },
  { key: "contract_value", value: "", type: "Number", required: true, locked: false, overridable: true, desc: "Total contract amount in base currency" },
  { key: "currency", value: "USD", type: "String", defaultValue: "USD", required: false, locked: true, overridable: false, desc: "Currency code (ISO 4217). Locked by template." },
  { key: "contract_type", value: "external", type: "Enum", defaultValue: "external", required: false, locked: false, overridable: true, desc: "Classification of contract scope" },
  { key: "requires_board_approval", value: "false", type: "Boolean", defaultValue: "false", required: false, locked: true, overridable: false, desc: "Whether board sign-off is mandatory. Locked by template." },
];

export const defaultCollaborationConfig: C3CollaborationConfiguration = {
  approvalsEnabled: true,
  rejectionEnabled: true,
  transfersEnabled: true,
  sharesEnabled: true,
  allowedTransferTrustGroupIds: [],
  defaultTransferMessage: "",
  requireApprovalForTransfer: true,
  requireReasonForRejection: true,
};

export const defaultCollaborationMetrics: C3CollaborationMetrics = {
  pendingApprovals: 0,
  recentRejections: 0,
  activeTransfers: 0,
  totalShares: 0,
};

export const defaultSecurityState: C3SecurityState = {
  deviceIdentityStatus: "configured",
  sovereignKeyringStatus: "healthy",
  trustGroupKeysState: "current",
  currentKekVersion: 1,
  activeEnvelopes: 0,
  revokedEnvelopes: 0,
  lastRotation: "Never",
  localDeviceId: "dev_local_01",
  localVaultAddress: "vault_local",
};

export const useC3ConfigurationStore = create<C3ConfigurationState>((set, get) => ({
  activeTab: "overview",
  selectedTrustGroupId: null,
  selectedChannelId: null,
  isCreateTrustGroupOpen: false,

  trustGroups: [],
  channels: [],
  collaborationConfig: defaultCollaborationConfig,
  metrics: defaultCollaborationMetrics,
  securityState: defaultSecurityState,
  templates: [],

  isLoading: false,
  error: null,

  // Draft channel state
  originalChannel: null,
  draftSlots: [],
  draftAssignments: [],
  draftProperties: [],
  draftPolicy: defaultPolicyDraft,
  isDirty: false,
  isSaving: false,
  saveStatus: "idle",
  saveFeedback: null,

  setActiveTab: (tab) => set({ activeTab: tab }),
  selectTrustGroup: (trustGroupId) => set({ selectedTrustGroupId: trustGroupId }),
  selectChannel: (channelId) => set({ selectedChannelId: channelId }),
  setCreateTrustGroupOpen: (open) => set({ isCreateTrustGroupOpen: open }),

  loadConfiguration: async () => {
    // Do not overwrite un-saved user edits if store is dirty
    if (get().isDirty) {
      return;
    }

    set({ isLoading: true, error: null });
    try {
      let channelStore = useC3ChannelStore.getState();
      let channelId = channelStore.activeChannelId || channelStore.activeChannel?.id || channelStore.channels[0]?.id;

      if (!channelId && channelStore.activeWorkspaceId) {
        try {
          await channelStore.fetchChannels(channelStore.activeWorkspaceId);
          channelStore = useC3ChannelStore.getState();
          channelId = channelStore.activeChannelId || channelStore.activeChannel?.id || channelStore.channels[0]?.id;
        } catch (fetchErr: any) {
          console.error(`[C3_LOAD] Failed to fetch workspace channels:`, fetchErr);
        }
      }

      console.log(`[BOUNDARIES][READ] loadConfiguration enter channelId=${channelId || "NONE_SELECTED"}`);
      let activeChannel = channelStore.activeChannel;

      if (channelId) {
        try {
          console.log(`[BOUNDARIES][READ] loadConfiguration calling refreshChannel(${channelId})`);
          activeChannel = await channelStore.refreshChannel(channelId);
          console.log(`[BOUNDARIES][READ] loadConfiguration refreshChannel returned activeChannel.slots=`, JSON.stringify(activeChannel?.slots));
        } catch (refreshErr: any) {
          console.error(`[C3_LOAD] refreshChannel failed for channel ${channelId}:`, refreshErr);
          throw new Error(`Failed to load channel details for ${channelId}: ${refreshErr?.message || refreshErr}`);
        }
      }

      if (activeChannel) {
        if (!activeChannel.slots) {
          console.warn("[SLOTS] Existing channel has missing/undefined slots data; preserving explicit empty array:", activeChannel.id);
        }

        const slotsFromBackend: SlotDraft[] = (activeChannel.slots || []).map((s) => ({
          id: s.id,
          name: s.name,
          role: s.role,
          vault_id: s.vault_id,
          gated: s.gated,
          order: s.order,
        }));

        const assignmentsFromBackend: AssignmentDraft[] = (activeChannel.assignments || []).map((a) => ({
          slot_id: a.slot_id,
          owner_id: a.owner_id,
          public_key: a.public_key,
          vault_address: a.vault_address,
        }));

        const propertiesFromBackend: PropertyDraft[] = (activeChannel.properties || []).map((p) => {
          const meta = initialProperties.find((ip) => ip.key === p.key);
          return {
            key: p.key,
            value: p.value,
            type: meta?.type || "String",
            required: meta?.required ?? false,
            locked: meta?.locked ?? false,
            overridable: meta?.overridable ?? true,
            defaultValue: meta?.defaultValue || p.value,
            desc: meta?.desc || "Channel property",
          };
        });

        const policyFromBackend: PolicyDraft = activeChannel.policy
          ? { ...defaultPolicyDraft, ...(activeChannel.policy as unknown as Partial<PolicyDraft>) }
          : defaultPolicyDraft;

        let trustGroupsFromBackend: C3TrustGroupSummary[] = [];
        try {
          const activeWsId = useC3WorkspaceStore.getState().activeWorkspaceId || useC3WorkspaceStore.getState().activeWorkspace?.id || activeChannel?.workspace_id || "";
          console.log(`[C3][TRACE][ADD_MEMBER][trace=tgcrud-001][15] layer=FRONTEND_STORE file=frontend/src/components/C3/configuration/store/useC3ConfigurationStore.ts function=loadConfiguration operation=listTrustGroups activeWsId=${activeWsId}`);
          const rawGroups = await listTrustGroups(activeWsId);
          const groupsArray = Array.isArray(rawGroups) ? rawGroups : (rawGroups as any)?.data || (rawGroups as any)?.Data || [];
          if (Array.isArray(groupsArray)) {
            trustGroupsFromBackend = groupsArray.map((tg: any) => {
              const memberCidsList = tg.member_cids || tg.MemberCIDs;
              const deserializedMembers = (() => {
                if (Array.isArray(memberCidsList) && memberCidsList.length > 0) {
                  return memberCidsList.map((cidItem: any) => {
                    const cidStr = typeof cidItem === "string" ? cidItem : (cidItem.id || cidItem.ID || cidItem.member_id || cidItem.MemberID || "vault_member");
                    return {
                      id: cidStr,
                      identityName: cidStr,
                      role: "member" as const,
                      status: "active" as const,
                      deviceCount: 1,
                      joinedAt: tg.created_at || tg.CreatedAt || new Date().toISOString(),
                    };
                  });
                }
                const rawList = tg.members || tg.Members || tg.key_envelopes || tg.KeyEnvelopes || [];
                return rawList.map((env: any) => ({
                  id: typeof env === "string" ? env : (env.id || env.ID || env.member_id || env.MemberID || "usr_member"),
                  identityName: typeof env === "string" ? env : (env.identityName || env.member_id || env.MemberID || env.id || env.ID || "vault_member"),
                  role: typeof env === "string" ? "member" : (env.role || "member"),
                  status: typeof env === "string" ? "active" : (env.status || (env.revoked_at ? "revoked" : "active")),
                  deviceCount: typeof env === "string" ? 1 : (env.deviceCount || 1),
                  joinedAt: typeof env === "string" ? new Date().toISOString() : (env.joinedAt || env.created_at || env.CreatedAt || new Date().toISOString()),
                }));
              })();

              console.log(`[C3][TRACE][ADD_MEMBER][trace=tgcrud-001][24] layer=DESKTOP_DESERIALIZATION file=frontend/src/components/C3/configuration/store/useC3ConfigurationStore.ts function=loadConfiguration input.member_cids=${JSON.stringify(memberCidsList)} output.members.length=${deserializedMembers.length}`);

              return {
                id: tg.id || tg.ID,
                name: tg.name || tg.Name || "Unnamed Trust Group",
                description: tg.description || tg.Description || "Sovereign Trust Group",
                memberCount: deserializedMembers.length,
                deviceCount: (tg.devices || []).length || 1,
                status: "active",
                kekVersion: tg.kek_version || tg.KEKVersion || 1,
                createdDate: tg.created_at || tg.CreatedAt || new Date().toISOString(),
                associatedChannelIds: activeChannel?.id ? [activeChannel.id] : [],
                members: deserializedMembers,
                devices: [],
                cryptographicState: {
                  kekVersion: tg.kek_version || tg.KEKVersion || 1,
                  envelopeCount: (tg.key_envelopes || tg.KeyEnvelopes || []).length || 1,
                  activeEnvelopes: (tg.key_envelopes || tg.KeyEnvelopes || []).length || 1,
                  revokedEnvelopes: 0,
                  lastRotation: tg.created_at || tg.CreatedAt || new Date().toISOString(),
                },
              };
            });
          }
        } catch (tgErr) {
          console.warn("[C3_STORE] Failed to fetch trust groups from backend, preserving empty list:", tgErr);
        }

        console.log(`[C3][TRACE][ADD_MEMBER][trace=tgcrud-001][25] layer=ZUSTAND_STORE file=frontend/src/components/C3/configuration/store/useC3ConfigurationStore.ts function=loadConfiguration output.trustGroupsCount=${trustGroupsFromBackend.length}`);

        console.log(`[C3_SLOTS_SOURCE]\nsource=LOAD_CONFIGURATION_BACKEND\nchannelId=${activeChannel.id}\nslotsCount=${slotsFromBackend.length}\nslots=${JSON.stringify(slotsFromBackend)}`);

        set({
          originalChannel: activeChannel,
          draftSlots: slotsFromBackend,
          draftAssignments: assignmentsFromBackend,
          draftProperties: propertiesFromBackend,
          draftPolicy: policyFromBackend,
          trustGroups: trustGroupsFromBackend,
          isDirty: false,
          isLoading: false,
          error: null,
        });
      } else {
        console.log(`[C3_SLOTS_SOURCE]\nsource=LOAD_CONFIGURATION_NO_ACTIVE_CHANNEL\nchannelId=none\nslotsCount=0\nslots=[]`);
        set({
          originalChannel: null,
          draftSlots: [],
          draftAssignments: [],
          draftProperties: [],
          draftPolicy: defaultPolicyDraft,
          trustGroups: [],
          channels: [],
          collaborationConfig: defaultCollaborationConfig,
          metrics: defaultCollaborationMetrics,
          securityState: defaultSecurityState,
          isDirty: false,
          isLoading: false,
          error: channelId ? "Channel not found." : "No active channel selected.",
        });
      }
    } catch (err: any) {
      console.error(`[C3_SLOTS_SOURCE]\nsource=LOAD_CONFIGURATION_ERROR\nchannelId=none\nslotsCount=0\nslots=[]`);
      set({
        originalChannel: null,
        draftSlots: [],
        draftAssignments: [],
        draftProperties: [],
        collaborationConfig: defaultCollaborationConfig,
        metrics: defaultCollaborationMetrics,
        securityState: defaultSecurityState,
        isLoading: false,
        error: err?.message || "Failed to load C3 Configuration from backend.",
      });
    }
  },

  setDraftSlots: (slots) => {
    console.log(`[C3_SLOTS_SOURCE]\nsource=SET_DRAFT_SLOTS_MUTATION\nchannelId=${get().originalChannel?.id || "none"}\nslotsCount=${slots.length}\nslots=${JSON.stringify(slots)}`);
    set({ draftSlots: slots, isDirty: true, saveStatus: "idle" });
  },
  setDraftAssignments: (assignments) => set({ draftAssignments: assignments, isDirty: true, saveStatus: "idle" }),
  setDraftProperties: (properties) => set({ draftProperties: properties, isDirty: true, saveStatus: "idle" }),
  setDraftPolicy: (policyUpdates) =>
    set((state) => ({
      draftPolicy: { ...state.draftPolicy, ...policyUpdates },
      isDirty: true,
      saveStatus: "idle",
    })),

  discardChanges: () => {
    const orig = get().originalChannel;
    if (orig) {
      if (!orig.slots) {
        console.warn("[SLOTS] Discarding changes for existing channel with missing/undefined slots:", orig.id);
      }
      const slots: SlotDraft[] = (orig.slots || []).map((s) => ({
        id: s.id,
        name: s.name,
        role: s.role,
        vault_id: s.vault_id,
        gated: s.gated,
        order: s.order,
      }));

      const assignments: AssignmentDraft[] = (orig.assignments || []).map((a) => ({
        slot_id: a.slot_id,
        owner_id: a.owner_id,
        public_key: a.public_key,
        vault_address: a.vault_address,
      }));

      const properties: PropertyDraft[] = (orig.properties || []).map((p) => ({
        key: p.key,
        value: p.value,
      }));

      const policyFromOrig: PolicyDraft = orig.policy
        ? { ...defaultPolicyDraft, ...(orig.policy as unknown as Partial<PolicyDraft>) }
        : defaultPolicyDraft;

      console.log(`[C3_SLOTS_SOURCE]\nsource=DISCARD_CHANGES_ORIGINAL\nchannelId=${orig.id}\nslotsCount=${slots.length}\nslots=${JSON.stringify(slots)}`);

      set({
        draftSlots: slots,
        draftAssignments: assignments,
        draftProperties: properties,
        draftPolicy: policyFromOrig,
        isDirty: false,
        saveStatus: "idle",
        saveFeedback: null,
      });
    } else {
      console.log(`[C3_SLOTS_SOURCE]\nsource=DISCARD_CHANGES_NULL\nchannelId=none\nslotsCount=0\nslots=[]`);
      set({
        draftSlots: [],
        draftAssignments: [],
        draftProperties: [],
        draftPolicy: defaultPolicyDraft,
        isDirty: false,
        saveStatus: "idle",
        saveFeedback: null,
      });
    }
  },

  saveChannelConfig: async (correlationId?: string) => {
    const id = correlationId || `save_${Date.now()}_${Math.random().toString(36).substring(2,6)}`;
    console.log(`[SLOTS][SAVE][id=${id}] STEP=03 EVENT=STORE_SAVE_ENTER channelId=active`);

    const { draftSlots, draftAssignments, draftProperties, draftPolicy, originalChannel } = get();
    const channelStore = useC3ChannelStore.getState();
    const channelId = originalChannel?.id || channelStore.activeChannelId || channelStore.activeChannel?.id || channelStore.channels[0]?.id;

    set({ isSaving: true, saveStatus: "saving", saveFeedback: null });

    try {
      if (channelId) {
        const payload = {
          channel_id: channelId,
          title: originalChannel?.title || "contract-execution",
          slots: draftSlots.map((s, idx) => ({
            id: s.id || String(idx + 1),
            name: s.name,
            role: s.role || "participant",
            vault_id: s.vault_id,
            gated: s.gated,
            order: s.order || idx + 1,
          })),
          assignments: draftAssignments.map((a) => ({
            slot_id: a.slot_id,
            owner_id: a.owner_id,
            public_key: a.public_key,
            vault_address: a.vault_address,
          })),
          properties: draftProperties.map((p) => ({
            key: p.key,
            value: p.value || "",
          })),
          policy: draftPolicy,
        };

        console.log(`[SLOTS][SAVE][id=${id}] STEP=04 EVENT=API_PAYLOAD payload=`, JSON.stringify(payload));
        console.log(`[SLOTS][SAVE][id=${id}] STEP=05 EVENT=API_CALL_START channelId=${channelId}`);

        const updated = await channelStore.updateChannelDetails(payload);

        console.log(`[SLOTS][SAVE][id=${id}] STEP=06 EVENT=API_CALL_RETURN success=true response=`, JSON.stringify(updated));

        const slotsFromBackend: SlotDraft[] = (updated.slots || []).map((s) => ({
          id: s.id,
          name: s.name,
          role: s.role,
          vault_id: s.vault_id,
          gated: s.gated,
          order: s.order,
        }));

        const assignmentsFromBackend: AssignmentDraft[] = (updated.assignments || []).map((a) => ({
          slot_id: a.slot_id,
          owner_id: a.owner_id,
          public_key: a.public_key,
          vault_address: a.vault_address,
        }));

        const propertiesFromBackend: PropertyDraft[] = (updated.properties || []).map((p) => {
          const meta = draftProperties.find((dp) => dp.key === p.key) || initialProperties.find((ip) => ip.key === p.key);
          return {
            key: p.key,
            value: p.value,
            type: meta?.type || "String",
            required: meta?.required ?? false,
            locked: meta?.locked ?? false,
            overridable: meta?.overridable ?? true,
            defaultValue: meta?.defaultValue || p.value,
            desc: meta?.desc || "Channel property",
          };
        });

        const policyFromUpdated: PolicyDraft = updated.policy
          ? { ...defaultPolicyDraft, ...(updated.policy as unknown as Partial<PolicyDraft>) }
          : get().draftPolicy;

        console.log(`[SLOTS][SAVE][id=${id}] STEP=20 EVENT=STORE_STATE_AFTER_API draftSlots=`, JSON.stringify(slotsFromBackend));

        console.log(`[C3_SLOTS_SOURCE]\nsource=SAVE_CHANNEL_CONFIG_POST_API\nchannelId=${updated.id}\nslotsCount=${slotsFromBackend.length}\nslots=${JSON.stringify(slotsFromBackend)}`);

        set({
          originalChannel: updated,
          draftSlots: slotsFromBackend,
          draftAssignments: assignmentsFromBackend,
          draftProperties: propertiesFromBackend,
          draftPolicy: policyFromUpdated,
          isDirty: false,
          isSaving: false,
          saveStatus: "success",
          saveFeedback: "✓ Slot configuration saved successfully.",
        });
      } else {
        throw new Error("Cannot save slots: no active backend channel is available.");
      }
      return true;
    } catch (err: any) {
      console.log(`[SLOTS][SAVE][id=${id}] STEP=06 EVENT=API_CALL_ERROR error=`, err?.message || err);
      set({
        isSaving: false,
        saveStatus: "error",
        saveFeedback: `Failed to save slot configuration: ${err?.message || "Unknown error"}`,
      });
      return false;
    }
  },

  createTrustGroup: async ({ name, description, vaultName }) => {
    set({ isLoading: true, error: null });
    try {
      const activeWsId = useC3WorkspaceStore.getState().activeWorkspaceId || useC3WorkspaceStore.getState().activeWorkspace?.id;
      const orig = get().originalChannel;
      const workspaceId = activeWsId || orig?.workspace_id || "";
      if (!workspaceId) {
        throw new Error("Cannot create Trust Group: no active workspace selected.");
      }
      const res = await createTrustGroupApi(name.trim(), workspaceId, vaultName);
      const tgData = res?.data || res?.Data || res;

      const newId = tgData?.id || tgData?.ID || `tg_${Math.random().toString(36).substring(2, 10)}`;
      const newGroup: C3TrustGroupSummary = {
        id: newId,
        name: tgData?.name || tgData?.Name || name.trim(),
        description: description?.trim() || "Sovereign Trust Group",
        memberCount: (tgData?.member_cids || tgData?.MemberCIDs || tgData?.members || []).length || 1,
        deviceCount: 1,
        kekVersion: tgData?.kek_version || tgData?.KEKVersion || 1,
        status: "active",
        createdDate: tgData?.created_at || tgData?.CreatedAt || new Date().toISOString(),
        associatedChannelIds: orig?.id ? [orig.id] : [],
        members: (tgData?.key_envelopes || tgData?.KeyEnvelopes || []).map((env: any) => ({
          id: env.id || env.ID || env.member_id,
          identityName: env.member_id || env.MemberID || "vault_member",
          role: "member",
          status: env.revoked_at ? "revoked" : "active",
          deviceCount: 1,
          joinedAt: env.created_at || new Date().toISOString(),
        })),
        devices: [],
        cryptographicState: {
          kekVersion: tgData?.kek_version || tgData?.KEKVersion || 1,
          envelopeCount: 1,
          activeEnvelopes: 1,
          revokedEnvelopes: 0,
          lastRotation: tgData?.created_at || tgData?.CreatedAt || new Date().toISOString(),
        },
      };

      set((state) => ({
        trustGroups: [newGroup, ...state.trustGroups],
        isLoading: false,
        isCreateTrustGroupOpen: false,
        selectedTrustGroupId: newId,
      }));

      return newGroup;
    } catch (err: any) {
      console.error("[C3_STORE] Failed to create Trust Group:", err);
      set({
        isLoading: false,
        error: err?.message || "Failed to create Trust Group.",
      });
      throw err;
    }
  },

  updateTrustGroup: async ({ id, name, description }) => {
    set({ isLoading: true, error: null });
    try {
      const res = await updateTrustGroupApi(id, name.trim());
      const tgData = res?.data || res?.Data || res;

      const updatedName = tgData?.name || tgData?.Name || name.trim();
      const updatedKek = tgData?.kek_version || tgData?.KEKVersion;

      let updatedGroup: C3TrustGroupSummary | null = null;

      set((state) => {
        const nextGroups = state.trustGroups.map((group) => {
          if (group.id !== id) return group;
          const updated: C3TrustGroupSummary = {
            ...group,
            name: updatedName,
            description: description?.trim() ?? group.description,
            kekVersion: updatedKek ?? group.kekVersion,
          };
          updatedGroup = updated;
          return updated;
        });

        return {
          trustGroups: nextGroups,
          isLoading: false,
        };
      });

      if (!updatedGroup) {
        throw new Error(`Trust group ${id} not found in store state.`);
      }

      return updatedGroup;
    } catch (err: any) {
      console.error("[C3_STORE] Failed to update Trust Group:", err);
      set({
        isLoading: false,
        error: err?.message || "Failed to update Trust Group.",
      });
      throw err;
    }
  },

  deleteTrustGroup: async (trustGroupId) => {
    set({ isLoading: true, error: null });
    try {
      await deleteTrustGroupApi(trustGroupId);

      set((state) => ({
        trustGroups: state.trustGroups.filter((tg) => tg.id !== trustGroupId),
        selectedTrustGroupId: state.selectedTrustGroupId === trustGroupId ? null : state.selectedTrustGroupId,
        isLoading: false,
      }));
    } catch (err: any) {
      console.error("[C3_STORE] Failed to delete Trust Group:", err);
      set({
        isLoading: false,
        error: err?.message || "Failed to delete Trust Group.",
      });
      throw err;
    }
  },

  updateCollaborationConfig: (updates) => {
    set((state) => ({
      collaborationConfig: {
        ...state.collaborationConfig,
        ...updates,
      },
    }));
  },
}));

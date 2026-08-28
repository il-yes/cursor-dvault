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

export type C3ConfigurationTab =
  | "overview"
  | "trust-groups"
  | "channels"
  | "collaboration"
  | "security"
  | "templates";

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
  }) => Promise<C3TrustGroupSummary>;
  updateCollaborationConfig: (updates: Partial<C3CollaborationConfiguration>) => void;
}

export const useC3ConfigurationStore = create<C3ConfigurationState>((set, get) => ({
  activeTab: "overview",
  selectedTrustGroupId: null,
  selectedChannelId: null,
  isCreateTrustGroupOpen: false,

  trustGroups: mockTrustGroups,
  channels: mockChannels,
  collaborationConfig: mockCollaborationConfig,
  metrics: mockCollaborationMetrics,
  securityState: mockSecurityState,
  templates: mockUseCaseTemplates,

  isLoading: false,
  error: null,

  setActiveTab: (tab) => set({ activeTab: tab }),
  selectTrustGroup: (trustGroupId) => set({ selectedTrustGroupId: trustGroupId }),
  selectChannel: (channelId) => set({ selectedChannelId: channelId }),
  setCreateTrustGroupOpen: (open) => set({ isCreateTrustGroupOpen: open }),

  loadConfiguration: async () => {
    set({ isLoading: true, error: null });
    try {
      // Phase 1 resolves through isolated mock data provider.
      // Phase 2 will execute Wails RPC adapters (e.g. ListTrustGroups, ListChannels).
      set({
        trustGroups: mockTrustGroups,
        channels: mockChannels,
        collaborationConfig: mockCollaborationConfig,
        metrics: mockCollaborationMetrics,
        securityState: mockSecurityState,
        templates: mockUseCaseTemplates,
        isLoading: false,
      });
    } catch (err: any) {
      set({
        isLoading: false,
        error: err?.message || "Failed to load C3 Configuration.",
      });
    }
  },

  createTrustGroup: async ({ name, description }) => {
    set({ isLoading: true, error: null });
    try {
      const newId = `tg_${Math.random().toString(36).substring(2, 10)}`;
      const newGroup: C3TrustGroupSummary = {
        id: newId,
        name: name.trim(),
        description: description?.trim() || "Sovereign Trust Group",
        memberCount: 1,
        deviceCount: 1,
        kekVersion: 1,
        status: "active",
        createdDate: new Date().toISOString(),
        associatedChannelIds: [],
        cryptographicState: {
          kekVersion: 1,
          envelopeCount: 1,
          activeEnvelopes: 1,
          revokedEnvelopes: 0,
          lastRotation: new Date().toISOString(),
        },
        members: [
          {
            id: "usr_creator_me",
            identityName: "Current Vault Owner (You)",
            role: "admin",
            status: "active",
            deviceCount: 1,
            joinedAt: new Date().toISOString(),
          },
        ],
        devices: [
          {
            id: "dev_current_desktop",
            deviceName: "Current Desktop Device",
            identityId: "usr_creator_me",
            status: "active",
            lastSeenAt: new Date().toISOString(),
            keyState: "valid",
          },
        ],
      };

      set((state) => ({
        trustGroups: [newGroup, ...state.trustGroups],
        collaborationConfig: {
          ...state.collaborationConfig,
          allowedTransferTrustGroupIds: [
            ...state.collaborationConfig.allowedTransferTrustGroupIds,
            newId,
          ],
        },
        isLoading: false,
        isCreateTrustGroupOpen: false,
        selectedTrustGroupId: newId,
      }));

      return newGroup;
    } catch (err: any) {
      set({
        isLoading: false,
        error: err?.message || "Failed to create Trust Group.",
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

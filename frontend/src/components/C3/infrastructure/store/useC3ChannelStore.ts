import { create } from "zustand";
import {
	listChannels,
	activateChannel as apiActivateChannel,
	revokeChannel as apiRevokeChannel,
	getChannel as apiGetChannel,
	updateChannel as apiUpdateChannel,
	deleteChannel as apiDeleteChannel,
	listParticipants as apiListParticipants,
	addParticipant as apiAddParticipant,
	inviteToChannel as apiInviteToChannel,
	acceptChannelInvitation as apiAcceptChannelInvitation,
	ChannelResponse,
	ChannelParticipantResponse,
	ChannelInvitationResponse,
	AddParticipantPayload,
	InviteToChannelPayload,
	UpdateChannelPayload,
} from "@/services/api";

export interface C3ChannelState {
	channels: ChannelResponse[];
	activeChannel: ChannelResponse | null;
	activeChannelId: string | null;
	activeWorkspaceId: string | null;
	isLoading: boolean;
	error: string | null;

	participants: ChannelParticipantResponse[];
	participantsChannelId: string | null;
	participantsLoading: boolean;
	participantsError: string | null;

	invitations: ChannelInvitationResponse[];
	invitationsChannelId: string | null;
	invitationsLoading: boolean;
	invitationsError: string | null;

	fetchChannels: (workspaceId: string) => Promise<void>;
	selectChannel: (channelId: string) => void;
	addChannel: (channel: ChannelResponse) => void;
	updateChannel: (channel: ChannelResponse) => void;
	activateChannel: (channelId: string) => Promise<ChannelResponse>;
	revokeChannel: (channelId: string) => Promise<void>;
	refreshChannel: (channelId: string) => Promise<ChannelResponse>;
	updateChannelDetails: (payload: UpdateChannelPayload) => Promise<ChannelResponse>;
	deleteChannel: (channelId: string) => Promise<void>;
	fetchParticipants: (channelId: string) => Promise<void>;
	addParticipant: (channelId: string, payload: AddParticipantPayload) => Promise<ChannelParticipantResponse>;
	inviteToChannel: (channelId: string, payload: InviteToChannelPayload) => Promise<ChannelInvitationResponse>;
	acceptInvitation: (invitationId: string) => Promise<ChannelInvitationResponse>;
	setInvitationsChannel: (channelId: string) => void;
	clearChannels: () => void;
}

export const useC3ChannelStore = create<C3ChannelState>((set, get) => ({
	channels: [],
	activeChannel: null,
	activeChannelId: null,
	activeWorkspaceId: null,
	isLoading: false,
	error: null,

	participants: [],
	participantsChannelId: null,
	participantsLoading: false,
	participantsError: null,

	invitations: [],
	invitationsChannelId: null,
	invitationsLoading: false,
	invitationsError: null,

	fetchChannels: async (workspaceId: string) => {
		if (!workspaceId) {
			set({
				channels: [],
				activeChannel: null,
				activeChannelId: null,
				activeWorkspaceId: null,
				isLoading: false,
				error: null,
				participants: [],
				participantsChannelId: null,
				participantsError: null,
				invitations: [],
				invitationsChannelId: null,
				invitationsError: null,
			});
			return;
		}

		// Workspace -> Channel Invariant: Clear channel-scoped state on a real
		// workspace change. On a refresh of the same workspace (e.g. the
		// LedgerLayout re-mounts after login while a Channel page is open), the
		// channel-scoped participant/invitation state is preserved so that
		// in-flight panel fetches are not dropped by the stale-response guard.
		const previousWorkspaceId = get().activeWorkspaceId;
		const currentActiveId = get().activeChannelId;
		set({
			activeWorkspaceId: workspaceId,
			channels: [],
			activeChannel: previousWorkspaceId === workspaceId ? get().activeChannel : null,
			activeChannelId: previousWorkspaceId === workspaceId ? currentActiveId : currentActiveId,
			isLoading: true,
			error: null,
			...(previousWorkspaceId !== workspaceId
				? {
						participants: [],
						participantsChannelId: null,
						participantsError: null,
						invitations: [],
						invitationsChannelId: null,
						invitationsError: null,
					}
				: {}),
		});

		try {
			const fetched = await listChannels(workspaceId);
			if (get().activeWorkspaceId !== workspaceId) return;

			const targetActiveId = get().activeChannelId || currentActiveId;
			const active =
				fetched.find((c) => c.id === targetActiveId) ??
				(fetched.length > 0 ? fetched[0] : null);

			set({
				channels: fetched,
				activeChannel: active,
				activeChannelId: active ? active.id : targetActiveId,
				isLoading: false,
				error: null,
			});
		} catch (err: unknown) {
			if (get().activeWorkspaceId !== workspaceId) return;
			console.error(`Failed to fetch channels for workspace ${workspaceId}:`, err);
			set({
				channels: [],
				isLoading: false,
				error: err instanceof Error ? err.message : "Failed to load channels.",
			});
		}
	},

	selectChannel: (channelId: string) => {
		if (get().activeChannelId === channelId && get().activeChannel?.id === channelId) {
			return;
		}
		const found = get().channels.find((c) => c.id === channelId) || null;
		set({
			activeChannel: found,
			activeChannelId: found ? found.id : channelId,
		});
	},

	addChannel: (newChannel: ChannelResponse) => {
		const currentWorkspaceId = get().activeWorkspaceId;

		if (!currentWorkspaceId || newChannel.workspace_id === currentWorkspaceId) {
			set((state) => ({
				activeWorkspaceId: newChannel.workspace_id || currentWorkspaceId,
				channels: [newChannel, ...state.channels.filter((c) => c.id !== newChannel.id)],
				activeChannel: newChannel,
				activeChannelId: newChannel.id,
			}));
		}
	},

	updateChannel: (updated: ChannelResponse) => {
		set((state) => {
			const channels = state.channels.map((c) => (c.id === updated.id ? { ...c, ...updated } : c));

			return {
				channels,
				activeChannel:
					state.activeChannel?.id === updated.id
						? { ...state.activeChannel, ...updated }
						: state.activeChannel,
			};
		});
	},

	activateChannel: async (channelId: string) => {
		const updated = await apiActivateChannel(channelId);
		set((state) => {
			const channels = state.channels.map((c) => (c.id === updated.id ? { ...c, ...updated } : c));

			return {
				channels,
				activeChannel:
					state.activeChannel?.id === updated.id
						? { ...state.activeChannel, ...updated }
						: state.activeChannel,
			};
		});
		return updated;
	},

	revokeChannel: async (channelId: string) => {
		const workspaceId = get().activeWorkspaceId;

		// Cloud is authoritative for the revocation. The revoke response carries
		// no Channel data, so no local status mutation happens here. After a
		// successful revoke the workspace channel list is refreshed and the
		// store is updated from the Cloud-provided (revoked) Channel.
		await apiRevokeChannel(channelId);

		if (!workspaceId) return;

		const refreshed = await listChannels(workspaceId);
		if (get().activeWorkspaceId !== workspaceId) return;

		set((state) => {
			const activeChannel =
				state.activeChannel?.id === channelId
					? refreshed.find((r) => r.id === channelId) ?? state.activeChannel
					: state.activeChannel;

			return {
				channels: refreshed,
				activeChannel,
			};
		});
	},

	// refreshChannel refetches a single Channel from the authoritative Cloud
	// backend (GET /channels/{id}) and replaces the local copy with the
	// Cloud-provided aggregate. The backend is the single source of truth for
	// channel existence; a missing channel surfaces as an error.
	refreshChannel: async (channelId: string) => {
		const fetched = await apiGetChannel(channelId);
		set((state) => {
			const channels = state.channels.map((c) => (c.id === channelId ? { ...c, ...fetched } : c));
			return {
				channels,
				activeChannel:
					state.activeChannel?.id === channelId
						? { ...state.activeChannel, ...fetched }
						: state.activeChannel,
			};
		});
		return fetched;
	},

	// updateChannelDetails persists a Channel update through the authoritative
	// Cloud backend (PUT /channels/{id}) and replaces the local copy with the
	// Cloud-returned Channel, which is authoritative for what was applied.
	updateChannelDetails: async (payload: UpdateChannelPayload) => {
		const updated = await apiUpdateChannel(payload);
		set((state) => {
			const channels = state.channels.map((c) => (c.id === updated.id ? { ...c, ...updated } : c));
			return {
				channels,
				activeChannel:
					state.activeChannel?.id === updated.id
						? { ...state.activeChannel, ...updated }
						: state.activeChannel,
			};
		});
		return updated;
	},

	// deleteChannel deletes a Channel through the authoritative Cloud backend
	// (DELETE /channels/{id}). The Cloud delete response carries no Channel
	// data, so the store removes the channel locally after a successful delete;
	// if the active channel was deleted the next channel becomes active.
	deleteChannel: async (channelId: string) => {
		await apiDeleteChannel(channelId);
		set((state) => {
			const channels = state.channels.filter((c) => c.id !== channelId);
			const removedActive = state.activeChannel?.id === channelId;
			return {
				channels,
				activeChannel: removedActive ? channels[0] ?? null : state.activeChannel,
				activeChannelId: removedActive ? channels[0]?.id ?? null : state.activeChannelId,
			};
		});
	},

	// fetchParticipants loads the vaults Cloud has persisted as participants for
	// a channel. The result is channel-scoped and dropped if the active channel
	// changed while the request was in flight (stale-response protection).
	fetchParticipants: async (channelId: string) => {
		if (!channelId) {
			return;
		}

		set({
			participants: [],
			participantsChannelId: channelId,
			participantsLoading: true,
			participantsError: null,
		});

		try {
			const fetched = await apiListParticipants(channelId);
			if (get().participantsChannelId !== channelId) return;

			set({
				participants: fetched,
				participantsLoading: false,
				participantsError: null,
			});
		} catch (err: unknown) {
			if (get().participantsChannelId !== channelId) return;
			console.error(`Failed to fetch participants for channel ${channelId}:`, err);
			set({
				participants: [],
				participantsLoading: false,
				participantsError: err instanceof Error ? err.message : "Failed to load participants.",
			});
		}
	},

	addParticipant: async (channelId: string, payload: AddParticipantPayload) => {
		const added = await apiAddParticipant(channelId, payload);
		if (get().participantsChannelId === channelId) {
			set((state) => {
				const existing = state.participants.some((p) => p.vault_id === added.vault_id);
				return {
					participants: existing
						? state.participants.map((p) =>
								p.vault_id === added.vault_id ? { ...p, ...added } : p,
							)
						: [...state.participants, added],
				};
			});
		}
		return added;
	},

	// inviteToChannel creates a channel invitation through the authoritative
	// Cloud backend. The inviter is the Desktop user's own vault. Cloud persists
	// the pending invitation and dedupes pending invitations for the same
	// channel + invitee; the store records the returned invitation for the
	// active channel.
	inviteToChannel: async (channelId: string, payload: InviteToChannelPayload) => {
		const invited = await apiInviteToChannel(channelId, payload);
		if (get().invitationsChannelId === channelId) {
			set((state) => {
				const existing = state.invitations.some((i) => i.id === invited.id);
				return {
					invitations: existing
						? state.invitations.map((i) => (i.id === invited.id ? { ...i, ...invited } : i))
						: [...state.invitations, invited],
				};
			});
		}
		return invited;
	},

	// acceptInvitation accepts a pending channel invitation through the
	// authoritative Cloud backend. Cloud validates that the accepting vault (the
	// Desktop user's own vault) is the invitation's invitee and persists the
	// resulting participant. The accept response carries the accepted
	// Invitation, not the participant — the store records the accepted
	// invitation and, when the channel is known, refreshes participants via
	// ListParticipants so the Cloud-created participant is observed without
	// being fabricated locally.
	acceptInvitation: async (invitationId: string) => {
		const accepted = await apiAcceptChannelInvitation(invitationId);
		if (get().invitationsChannelId === accepted.channel_id) {
			set((state) => ({
				invitations: state.invitations.map((i) => (i.id === accepted.id ? { ...i, ...accepted } : i)),
			}));
		}

		if (accepted.channel_id) {
			await get().fetchParticipants(accepted.channel_id);
		}

		return accepted;
	},

	// setInvitationsChannel scopes the session-created invitations to a
	// channel. Cloud has no invitation-list endpoint; the invitation list is the
	// set of invitations created (and accepted) during this session, dropped
	// when the active channel changes.
	setInvitationsChannel: (channelId: string) => {
		if (get().invitationsChannelId === channelId) return;
		set({
			invitations: [],
			invitationsChannelId: channelId,
			invitationsError: null,
		});
	},

	clearChannels: () => {
		set({
			channels: [],
			activeChannel: null,
			activeChannelId: null,
			activeWorkspaceId: null,
			isLoading: false,
			error: null,
			participants: [],
			participantsChannelId: null,
			participantsLoading: false,
			participantsError: null,
			invitations: [],
			invitationsChannelId: null,
			invitationsLoading: false,
			invitationsError: null,
		});
	},
}));

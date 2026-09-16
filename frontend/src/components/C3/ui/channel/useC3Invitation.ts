import { useState, useCallback } from "react";
import { useC3ChannelStore } from "../../infrastructure/store/useC3ChannelStore";
import { ChannelInvitationResponse } from "@/services/api";

export async function sendC3Invitation(
	inviteeVaultId: string,
	targetChannelId?: string | null,
): Promise<ChannelInvitationResponse> {
	const id = inviteeVaultId.trim();
	if (!id) {
		throw new Error("Invitee vault ID is required");
	}
	const effectiveChannelId = targetChannelId || useC3ChannelStore.getState().activeChannelId;
	if (!effectiveChannelId) {
		throw new Error("Channel ID is required for sending an invitation");
	}

	return await useC3ChannelStore.getState().inviteToChannel(effectiveChannelId, { invitee_vault_id: id });
}

export function useC3Invitation(channelId?: string | null) {
	const [inviting, setInviting] = useState(false);
	const [inviteError, setInviteError] = useState<string | null>(null);
	const [inviteSuccess, setInviteSuccess] = useState<string | null>(null);

	const sendInvitation = useCallback(
		async (inviteeVaultId: string, targetChannelId?: string): Promise<ChannelInvitationResponse> => {
			setInviting(true);
			setInviteError(null);
			setInviteSuccess(null);
			try {
				const invited = await sendC3Invitation(inviteeVaultId, targetChannelId || channelId);
				setInviteSuccess(`Invitation sent — ${invited.status} for ${invited.invitee_vault_id || inviteeVaultId}.`);
				return invited;
			} catch (err: unknown) {
				const msg = err instanceof Error ? err.message : "Failed to invite vault.";
				setInviteError(msg);
				throw err;
			} finally {
				setInviting(false);
			}
		},
		[channelId],
	);

	return {
		sendInvitation,
		inviting,
		inviteError,
		inviteSuccess,
		setInviteError,
		setInviteSuccess,
		acceptInvitation: (invitationId: string) => useC3ChannelStore.getState().acceptInvitation(invitationId),
		setInvitationsChannel: (targetChannelId: string) => useC3ChannelStore.getState().setInvitationsChannel(targetChannelId),
	};
}

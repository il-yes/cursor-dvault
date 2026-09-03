import { describe, it, expect, beforeEach, vi } from "vitest";
import { useC3ChannelStore } from "../../../infrastructure/store/useC3ChannelStore";
import { sendC3Invitation } from "../useC3Invitation";

describe("useC3Invitation & Unified Member Invitation Path Invariants", () => {
	beforeEach(() => {
		useC3ChannelStore.setState({
			channels: [],
			activeChannel: null,
			activeChannelId: "ch_authoritative_test_100",
			invitations: [],
			invitationsChannelId: null,
			invitationsError: null,
		});
	});

	it("invokes the exact same inviteToChannel store method with identity_vault_id passed unchanged", async () => {
		let capturedChannelId = "";
		let capturedPayload: any = null;

		// Mock store inviteToChannel method
		const mockInviteToChannel = vi.fn().mockImplementation(async (channelId: string, payload: any) => {
			capturedChannelId = channelId;
			capturedPayload = payload;
			return {
				id: "inv_998877",
				channel_id: channelId,
				inviter_vault_id: "ankhora://my_vault",
				invitee_vault_id: payload.invitee_vault_id,
				status: "pending",
				created_at: new Date().toISOString(),
			};
		});

		useC3ChannelStore.setState({ inviteToChannel: mockInviteToChannel });

		const targetVaultId = "ankhora://vault_alice_authoritative_555";
		const targetChannelId = "ch_authoritative_test_100";

		// Direct call to exported invitation helper used by InvitationsPanel and TrustGroupDetailView
		const result = await sendC3Invitation(targetVaultId, targetChannelId);

		expect(mockInviteToChannel).toHaveBeenCalledTimes(1);
		expect(capturedChannelId).toBe("ch_authoritative_test_100");
		expect(capturedPayload).toEqual({ invitee_vault_id: targetVaultId });
		expect(result.invitee_vault_id).toBe(targetVaultId);

		// Invariant checks
		expect(capturedPayload.invitee_vault_id).not.toContain("tg_"); // Must not be channel/trustgroup id
		expect(capturedChannelId).not.toBe(targetVaultId); // Must preserve distinction
	});

	it("preserves distinction between Channel participation (addParticipant) and TrustGroup invitation (sendC3Invitation)", async () => {
		const mockAddParticipant = vi.fn().mockResolvedValue({
			vault_id: "ankhora://vault_bob",
			role: "participant",
			joined_at: Math.floor(Date.now() / 1000),
		});

		const mockInviteToChannel = vi.fn().mockResolvedValue({
			id: "inv_12345",
			invitee_vault_id: "ankhora://vault_bob",
			status: "pending",
		});

		useC3ChannelStore.setState({
			addParticipant: mockAddParticipant,
			inviteToChannel: mockInviteToChannel,
		});

		// 1. Add Participant (Channel participation)
		await useC3ChannelStore.getState().addParticipant("ch_100", { vault_id: "ankhora://vault_bob" });
		expect(mockAddParticipant).toHaveBeenCalledWith("ch_100", { vault_id: "ankhora://vault_bob" });

		// 2. Add Member (TrustGroup invitation)
		await sendC3Invitation("ankhora://vault_bob", "ch_100");
		expect(mockInviteToChannel).toHaveBeenCalledWith("ch_100", { invitee_vault_id: "ankhora://vault_bob" });
	});
});

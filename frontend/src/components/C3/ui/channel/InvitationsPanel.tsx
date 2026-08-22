import { useEffect, useState } from "react";
import { ChannelInvitationResponse } from "@/services/api";
import { useC3ChannelStore } from "../../infrastructure/store/useC3ChannelStore";

const rowStyle: React.CSSProperties = {
	display: "flex",
	alignItems: "center",
	justifyContent: "space-between",
	padding: "8px 14px",
	borderBottom: "1px solid rgba(255, 255, 255, 0.06)",
	fontSize: "13px",
};

const labelStyle: React.CSSProperties = {
	fontSize: "11px",
	fontWeight: 600,
	letterSpacing: "0.08em",
	textTransform: "uppercase",
	opacity: 0.6,
};

const inputStyle: React.CSSProperties = {
	background: "rgba(255,255,255,0.05)",
	border: "1px solid rgba(255,255,255,0.12)",
	borderRadius: "6px",
	padding: "8px 10px",
	fontSize: "13px",
	color: "inherit",
};

const buttonStyle: React.CSSProperties = {
	background: "rgba(255,255,255,0.08)",
	border: "1px solid rgba(255,255,255,0.18)",
	borderRadius: "6px",
	padding: "8px 14px",
	fontSize: "13px",
	cursor: "pointer",
	color: "inherit",
};

function formatCreatedAt(createdAt: string): string {
	if (!createdAt) return "";
	const d = new Date(createdAt);
	if (Number.isNaN(d.getTime())) return "";
	return d.toLocaleDateString();
}

// InvitationsPanel exposes the Desktop -> Cloud channel invitation lifecycle on
// the Channel Detail page. Cloud is authoritative: an invitation is created
// through the Cloud backend and accepted through the Cloud backend. Acceptance
// records the accepted Invitation (never a fabricated participant); the
// Cloud-created participant is observed through the participants panel, which
// the store refreshes after acceptance. Cloud has no invitation-list endpoint,
// so the panel lists only the invitations created/accepted during this session.
export const InvitationsPanel = ({ channelId }: { channelId: string }) => {
	const {
		invitations,
		invitationsError,
		inviteToChannel,
		acceptInvitation,
		setInvitationsChannel,
	} = useC3ChannelStore();

	const [inviteeVaultId, setInviteeVaultId] = useState("");
	const [inviting, setInviting] = useState(false);
	const [inviteError, setInviteError] = useState<string | null>(null);
	const [inviteSuccess, setInviteSuccess] = useState<string | null>(null);

	const [acceptId, setAcceptId] = useState("");
	const [accepting, setAccepting] = useState(false);
	const [acceptError, setAcceptError] = useState<string | null>(null);

	useEffect(() => {
		if (channelId) {
			setInvitationsChannel(channelId);
		}
	}, [channelId, setInvitationsChannel]);

	const handleInvite = async () => {
		const id = inviteeVaultId.trim();
		if (!id || inviting) return;

		setInviting(true);
		setInviteError(null);
		setInviteSuccess(null);
		try {
			const invited = await inviteToChannel(channelId, { invitee_vault_id: id });
			setInviteeVaultId("");
			setInviteSuccess(`Invitation sent — ${invited.status} for ${invited.invitee_vault_id || id}.`);
		} catch (err: unknown) {
			setInviteError(err instanceof Error ? err.message : "Failed to invite vault.");
		} finally {
			setInviting(false);
		}
	};

	const handleAccept = async (invitation: ChannelInvitationResponse) => {
		if (accepting) return;
		setAccepting(true);
		setAcceptError(null);
		try {
			await acceptInvitation(invitation.id);
			const workspaceStore = (await import("../../infrastructure/store/useC3WorkspaceStore")).useC3WorkspaceStore;
			await workspaceStore.getState().fetchWorkspaces();
			const freshWorkspaces = workspaceStore.getState().workspaces;
			if (freshWorkspaces.length > 0) {
				const target = freshWorkspaces.find((w) => w.id === channelId) || freshWorkspaces[0];
				workspaceStore.getState().selectWorkspace(target.id);
				await useC3ChannelStore.getState().fetchChannels(target.id);
			}
		} catch (err: unknown) {
			setAcceptError(err instanceof Error ? err.message : "Failed to accept invitation.");
		} finally {
			setAccepting(false);
		}
	};

	const handleAcceptById = async () => {
		const id = acceptId.trim();
		if (!id || accepting) return;

		setAccepting(true);
		setAcceptError(null);
		try {
			await acceptInvitation(id);
			setAcceptId("");
			const workspaceStore = (await import("../../infrastructure/store/useC3WorkspaceStore")).useC3WorkspaceStore;
			await workspaceStore.getState().fetchWorkspaces();
			const freshWorkspaces = workspaceStore.getState().workspaces;
			if (freshWorkspaces.length > 0) {
				const target = freshWorkspaces.find((w) => w.id === channelId) || freshWorkspaces[0];
				workspaceStore.getState().selectWorkspace(target.id);
				await useC3ChannelStore.getState().fetchChannels(target.id);
			}
		} catch (err: unknown) {
			setAcceptError(err instanceof Error ? err.message : "Failed to accept invitation.");
		} finally {
			setAccepting(false);
		}
	};

	return (
		<div className="channel-panel">
			<div style={{ marginBottom: "8px" }}>
				<span style={labelStyle}>Invitations</span>
				<span style={{ opacity: 0.5, marginLeft: "8px", fontSize: "11px" }}>
					{invitations.length} this session
				</span>
			</div>

			{invitationsError && (
				<div style={{ color: "#EF4444", fontSize: "12px", padding: "6px 0" }}>
					⚠️ {invitationsError}
				</div>
			)}

			{inviteSuccess && (
				<div style={{ color: "#4ADE80", fontSize: "12px", padding: "6px 0" }}>
					{inviteSuccess}
				</div>
			)}

			{invitations.length === 0 && !inviting && !inviteError && !inviteSuccess && (
				<div style={{ opacity: 0.6, fontSize: "12px", padding: "6px 0" }}>
					No invitations created in this session.
				</div>
			)}

			{invitations.map((i) => (
				<div key={i.id} style={rowStyle}>
					<div>
						<div style={{ fontWeight: 600 }}>{i.invitee_vault_id || "—"}</div>
						<div style={{ opacity: 0.5, fontSize: "11px" }}>
							{i.status}
							{i.inviter_vault_id ? ` · by ${i.inviter_vault_id}` : ""}
							{formatCreatedAt(i.created_at) ? ` · ${formatCreatedAt(i.created_at)}` : ""}
							{` · ${i.id.slice(0, 8)}…`}
						</div>
					</div>
					{i.status === "pending" ? (
						<button
							onClick={() => handleAccept(i)}
							disabled={accepting}
							style={{ ...buttonStyle, opacity: accepting ? 0.5 : 1, cursor: accepting ? "default" : "pointer" }}
						>
							{accepting ? "Accepting…" : "Accept"}
						</button>
					) : (
						<span style={{ opacity: 0.6, fontSize: "11px", fontWeight: 600 }}>{i.status}</span>
					)}
				</div>
			))}

			<div style={{ display: "flex", gap: "8px", marginTop: "12px" }}>
				<input
					value={inviteeVaultId}
					onChange={(e) => setInviteeVaultId(e.target.value)}
					placeholder="Invite vault address (ankhora://vault_id)"
					style={{ flex: 1, ...inputStyle }}
					onKeyDown={(e) => {
						if (e.key === "Enter") handleInvite();
					}}
				/>
				<button
					onClick={handleInvite}
					disabled={inviting || !inviteeVaultId.trim()}
					style={{
						...buttonStyle,
						cursor: inviting || !inviteeVaultId.trim() ? "default" : "pointer",
						opacity: inviting || !inviteeVaultId.trim() ? 0.5 : 1,
					}}
				>
					{inviting ? "Inviting…" : "Invite vault"}
				</button>
			</div>

			{inviteError && (
				<div style={{ color: "#EF4444", fontSize: "12px", marginTop: "8px" }}>
					⚠️ {inviteError}
				</div>
			)}

			<div style={{ display: "flex", gap: "8px", marginTop: "12px", borderTop: "1px solid rgba(255,255,255,0.08)", paddingTop: "12px" }}>
				<input
					value={acceptId}
					onChange={(e) => setAcceptId(e.target.value)}
					placeholder="Accept an invitation by ID (invitee is this vault)"
					style={{ flex: 1, ...inputStyle }}
					onKeyDown={(e) => {
						if (e.key === "Enter") handleAcceptById();
					}}
				/>
				<button
					onClick={handleAcceptById}
					disabled={accepting || !acceptId.trim()}
					style={{
						...buttonStyle,
						cursor: accepting || !acceptId.trim() ? "default" : "pointer",
						opacity: accepting || !acceptId.trim() ? 0.5 : 1,
					}}
				>
					{accepting ? "Accepting…" : "Accept by ID"}
				</button>
			</div>

			{acceptError && (
				<div style={{ color: "#EF4444", fontSize: "12px", marginTop: "8px" }}>
					⚠️ {acceptError}
				</div>
			)}

			<div style={{ opacity: 0.5, fontSize: "11px", marginTop: "10px" }}>
				Acceptance runs through the Cloud backend: the accepting vault must be the
				invitation's invitee, and the resulting participant is created and owned by Cloud.
			</div>
		</div>
	);
};

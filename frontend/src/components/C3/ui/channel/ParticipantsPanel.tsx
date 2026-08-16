import { useEffect, useState } from "react";
import { ChannelParticipantResponse } from "@/services/api";
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

function formatJoinedAt(joinedAt: number): string {
	if (!joinedAt) return "";
	const d = new Date(joinedAt * 1000);
	if (Number.isNaN(d.getTime())) return "";
	return d.toLocaleDateString();
}

export const ParticipantsPanel = ({ channelId }: { channelId: string }) => {
	const { participants, participantsChannelId, participantsLoading, participantsError, fetchParticipants, addParticipant } =
		useC3ChannelStore();

	const [vaultId, setVaultId] = useState("");
	const [adding, setAdding] = useState(false);
	const [addError, setAddError] = useState<string | null>(null);

	// Participants are channel-scoped. Fetch whenever the requested channel
	// differs from the channel the store is currently scoped to (including the
	// null case after a store reset), so the Cloud-persisted list is always
	// re-fetched on mount and after any channel-scope invalidation.
	useEffect(() => {
		if (channelId && participantsChannelId !== channelId) {
			fetchParticipants(channelId);
		}
	}, [channelId, participantsChannelId, fetchParticipants]);

	const handleAdd = async () => {
		const id = vaultId.trim();
		if (!id || adding) return;

		setAdding(true);
		setAddError(null);
		try {
			const added = await addParticipant(channelId, { vault_id: id });
			setVaultId("");
			if (added) {
				await fetchParticipants(channelId);
			}
		} catch (err: unknown) {
			setAddError(err instanceof Error ? err.message : "Failed to add participant.");
		} finally {
			setAdding(false);
		}
	};

	return (
		<div className="channel-panel">
			<div style={{ marginBottom: "8px" }}>
				<span style={labelStyle}>Participants</span>
				<span style={{ opacity: 0.5, marginLeft: "8px", fontSize: "11px" }}>
					{participants.length} from Cloud
				</span>
			</div>

			{participantsLoading && (
				<div style={{ opacity: 0.6, fontSize: "12px", padding: "6px 0" }}>
					Loading participants…
				</div>
			)}

			{!participantsLoading && participantsError && (
				<div style={{ color: "#EF4444", fontSize: "12px", padding: "6px 0" }}>
					⚠️ {participantsError}
				</div>
			)}

			{!participantsLoading && !participantsError && participants.length === 0 && (
				<div style={{ opacity: 0.6, fontSize: "12px", padding: "6px 0" }}>
					No external vaults joined yet.
				</div>
			)}

			{participants.map((p) => (
				<div key={p.vault_id} style={rowStyle}>
					<div>
						<div style={{ fontWeight: 600 }}>{p.vault_id}</div>
						<div style={{ opacity: 0.5, fontSize: "11px" }}>
							{p.role || "—"}
							{p.direction ? ` · ${p.direction}` : ""}
							{formatJoinedAt(p.joined_at) ? ` · joined ${formatJoinedAt(p.joined_at)}` : ""}
						</div>
					</div>
					{(p.permissions?.length ?? 0) > 0 && (
						<div style={{ opacity: 0.6, fontSize: "11px" }}>
							{p.permissions.join(", ")}
						</div>
					)}
				</div>
			))}

			<div style={{ display: "flex", gap: "8px", marginTop: "12px" }}>
				<input
					value={vaultId}
					onChange={(e) => setVaultId(e.target.value)}
					placeholder="Paste vault address (ankhora://vault_id)"
					style={{
						flex: 1,
						background: "rgba(255,255,255,0.05)",
						border: "1px solid rgba(255,255,255,0.12)",
						borderRadius: "6px",
						padding: "8px 10px",
						fontSize: "13px",
						color: "inherit",
					}}
					onKeyDown={(e) => {
						if (e.key === "Enter") handleAdd();
					}}
				/>
				<button
					onClick={handleAdd}
					disabled={adding || !vaultId.trim()}
					style={{
						background: "rgba(255,255,255,0.08)",
						border: "1px solid rgba(255,255,255,0.18)",
						borderRadius: "6px",
						padding: "8px 14px",
						fontSize: "13px",
						cursor: adding || !vaultId.trim() ? "default" : "pointer",
						opacity: adding || !vaultId.trim() ? 0.5 : 1,
						color: "inherit",
					}}
				>
					{adding ? "Adding…" : "Add external vault"}
				</button>
			</div>

			{addError && (
				<div style={{ color: "#EF4444", fontSize: "12px", marginTop: "8px" }}>
					⚠️ {addError}
				</div>
			)}
		</div>
	);
};

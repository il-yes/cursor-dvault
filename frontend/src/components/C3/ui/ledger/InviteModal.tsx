import React, { useState } from "react";
import { inviteToChannel } from "@/services/api";
import { useC3ChannelStore } from "../../infrastructure/store/useC3ChannelStore";

interface InviteModalProps {
	isOpen: boolean;
	onClose: () => void;
	defaultChannelId?: string;
}

export const InviteModal: React.FC<InviteModalProps> = ({
	isOpen,
	onClose,
	defaultChannelId,
}) => {
	const { channels, setInvitationsChannel } = useC3ChannelStore();

	const [selectedChannelId, setSelectedChannelId] = useState(
		defaultChannelId || channels[0]?.id || ""
	);
	const [inviteeVaultId, setInviteeVaultId] = useState("");
	const [submitting, setSubmitting] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [successMsg, setSuccessMsg] = useState<string | null>(null);

	if (!isOpen) return null;

	const effectiveChannelId = selectedChannelId || defaultChannelId || channels[0]?.id || "";

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		const id = inviteeVaultId.trim();
		if (!id) {
			setError("Please enter a valid Vault ID or Email");
			return;
		}
		if (!effectiveChannelId) {
			setError("Please select a target channel");
			return;
		}

		setSubmitting(true);
		setError(null);
		setSuccessMsg(null);

		try {
			const res = await inviteToChannel(effectiveChannelId, {
				invitee_vault_id: id,
			});
			setSuccessMsg(`Invitation sent successfully to ${id}`);
			setInviteeVaultId("");
			setInvitationsChannel(effectiveChannelId);
			setTimeout(() => {
				setSuccessMsg(null);
				onClose();
			}, 1500);
		} catch (err: any) {
			setError(err?.message || "Failed to send invitation");
		} finally {
			setSubmitting(false);
		}
	};

	return (
		<div
			className="scrim"
			style={{
				position: "fixed",
				inset: 0,
				background: "rgba(0,0,0,0.5)",
				zIndex: 9999,
				display: "flex",
				alignItems: "center",
				justifyContent: "center",
			}}
		>
			<div
				className="modal"
				style={{
					width: "480px",
					background: "#ffffff",
					borderRadius: "12px",
					boxShadow: "0 20px 60px rgba(0,0,0,0.3)",
					padding: "24px",
					display: "flex",
					flexDirection: "column",
					gap: "16px",
				}}
			>
				<div
					style={{
						display: "flex",
						alignItems: "center",
						justifyContent: "space-between",
						borderBottom: "1px solid #eee",
						paddingBottom: "12px",
					}}
				>
					<h3
						style={{
							margin: 0,
							fontSize: "16px",
							fontWeight: 700,
							color: "#1a1a1a",
						}}
					>
						Invite Member to Channel
					</h3>
					<button
						onClick={onClose}
						style={{
							background: "transparent",
							border: "none",
							fontSize: "18px",
							cursor: "pointer",
							color: "#888",
						}}
					>
						✕
					</button>
				</div>

				<form onSubmit={handleSubmit} style={{ display: "flex", flexDirection: "column", gap: "14px" }}>
					{channels.length > 0 && (
						<div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
							<label style={{ fontSize: "12px", fontWeight: 600, color: "#555" }}>
								Target Channel
							</label>
							<select
								value={effectiveChannelId}
								onChange={(e) => setSelectedChannelId(e.target.value)}
								style={{
									padding: "8px 12px",
									borderRadius: "6px",
									border: "1px solid #ccc",
									fontSize: "13px",
									background: "#fff",
								}}
							>
								{channels.map((ch) => (
									<option key={ch.id} value={ch.id}>
										{ch.name || ch.title || ch.id}
									</option>
								))}
							</select>
						</div>
					)}

					<div style={{ display: "flex", flexDirection: "column", gap: "6px" }}>
						<label style={{ fontSize: "12px", fontWeight: 600, color: "#555" }}>
							Invitee Email / Vault ID
						</label>
						<input
							type="text"
							placeholder="bob@ankhora.test or vault_123"
							value={inviteeVaultId}
							onChange={(e) => setInviteeVaultId(e.target.value)}
							style={{
								padding: "9px 12px",
								borderRadius: "6px",
								border: "1px solid #ccc",
								fontSize: "13px",
								outline: "none",
							}}
							disabled={submitting}
							autoFocus
						/>
					</div>

					{error && (
						<div
							style={{
								padding: "8px 12px",
								borderRadius: "6px",
								background: "#fef2f2",
								border: "1px solid #fecaca",
								color: "#dc2626",
								fontSize: "12px",
							}}
						>
							{error}
						</div>
					)}

					{successMsg && (
						<div
							style={{
								padding: "8px 12px",
								borderRadius: "6px",
								background: "#f0fdf4",
								border: "1px solid #bbf7d0",
								color: "#16a34a",
								fontSize: "12px",
							}}
						>
							{successMsg}
						</div>
					)}

					<div
						style={{
							display: "flex",
							justifyContent: "flex-end",
							gap: "10px",
							marginTop: "8px",
							borderTop: "1px solid #eee",
							paddingTop: "12px",
						}}
					>
						<button
							type="button"
							onClick={onClose}
							disabled={submitting}
							style={{
								padding: "8px 14px",
								borderRadius: "6px",
								border: "1px solid #ccc",
								background: "#f5f5f5",
								color: "#333",
								fontSize: "13px",
								cursor: "pointer",
							}}
						>
							Cancel
						</button>
						<button
							type="submit"
							disabled={submitting}
							style={{
								padding: "8px 16px",
								borderRadius: "6px",
								border: "none",
								background: "#10b981",
								color: "#fff",
								fontWeight: 600,
								fontSize: "13px",
								cursor: submitting ? "not-allowed" : "pointer",
							}}
						>
							{submitting ? "Sending..." : "Send Invitation"}
						</button>
					</div>
				</form>
			</div>
		</div>
	);
};

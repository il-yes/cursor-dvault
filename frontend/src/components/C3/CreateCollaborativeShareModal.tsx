import React, { useState } from "react";
import { useC3WorkspaceStore } from "@/components/C3/infrastructure/store/useC3WorkspaceStore";
import { useC3ChannelStore } from "@/components/C3/infrastructure/store/useC3ChannelStore";
import { useC3ThreadStore } from "@/components/C3/infrastructure/store/useC3ThreadStore";
import { useC3CollaborationStore } from "@/components/C3/infrastructure/store/useC3CollaborationStore";
import { ShareEntryRefResponse } from "@/services/api";

interface CreateCollaborativeShareModalProps {
	isOpen: boolean;
	onClose: () => void;
	onShareCreated?: (shareRef: ShareEntryRefResponse) => void;
}

export const CreateCollaborativeShareModal: React.FC<CreateCollaborativeShareModalProps> = ({
	isOpen,
	onClose,
	onShareCreated,
}) => {
	const { activeWorkspace } = useC3WorkspaceStore();
	const { activeChannel } = useC3ChannelStore();
	const { activeThread, activeThreadId } = useC3ThreadStore();
	const { createShare, isLoading, error } = useC3CollaborationStore();

	const [trustGroupId, setTrustGroupId] = useState("tg_legal_counsel");
	const [targetVaultId, setTargetVaultId] = useState("vault_partner_02");
	const [assetCid, setAssetCid] = useState("QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco");
	const [notes, setNotes] = useState("");
	const [localError, setLocalError] = useState<string | null>(null);

	if (!isOpen) return null;

	const handleReset = () => {
		setTrustGroupId("tg_legal_counsel");
		setTargetVaultId("vault_partner_02");
		setAssetCid("QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco");
		setNotes("");
		setLocalError(null);
	};

	const handleClose = () => {
		if (isLoading) return;
		handleReset();
		onClose();
	};

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		setLocalError(null);

		if (!activeThreadId) {
			setLocalError("No active thread selected. Please select a thread before creating a share.");
			return;
		}

		if (!trustGroupId.trim()) {
			setLocalError("TrustGroup ID is required.");
			return;
		}

		if (!targetVaultId.trim()) {
			setLocalError("Target Vault ID is required.");
			return;
		}

		if (!assetCid.trim()) {
			setLocalError("Asset CID is required.");
			return;
		}

		try {
			const shareRef = await createShare({
				thread_id: activeThreadId,
				trust_group_id: trustGroupId.trim(),
				asset_cid: assetCid.trim(),
				target_vault_id: targetVaultId.trim(),
				notes: notes.trim(),
			});

			handleReset();
			if (onShareCreated) {
				onShareCreated(shareRef);
			}
			onClose();
		} catch (err: any) {
			setLocalError(err?.message || "Failed to create collaborative share.");
		}
	};

	const displayError = localError || error;

	return (
		<div
			className="modal-backdrop"
			style={{
				position: "fixed",
				top: 0,
				left: 0,
				right: 0,
				bottom: 0,
				backgroundColor: "rgba(0, 0, 0, 0.75)",
				backdropFilter: "blur(4px)",
				display: "flex",
				alignItems: "center",
				justifyContent: "center",
				zIndex: 1000,
			}}
		>
			<div
				className="modal-content"
				style={{
					backgroundColor: "#161B22",
					border: "1px solid #30363D",
					borderRadius: "8px",
					width: "100%",
					maxWidth: "520px",
					boxShadow: "0 20px 25px -5px rgba(0, 0, 0, 0.5), 0 10px 10px -5px rgba(0, 0, 0, 0.04)",
					overflow: "hidden",
					display: "flex",
					flexDirection: "column",
				}}
			>
				{/* Modal Header */}
				<div
					className="modal-header"
					style={{
						padding: "16px 20px",
						borderBottom: "1px solid rgba(255, 255, 255, 0.08)",
						display: "flex",
						alignItems: "center",
						justifyContent: "space-between",
					}}
				>
					<h3 style={{ margin: 0, fontSize: "16px", fontWeight: 600, color: "#F0F6FC" }}>
						🤝 Create Collaborative Share
					</h3>
					<button
						onClick={handleClose}
						disabled={isLoading}
						style={{
							background: "none",
							border: "none",
							color: "#8B949E",
							fontSize: "18px",
							cursor: "pointer",
							padding: "4px",
						}}
					>
						✕
					</button>
				</div>

				{/* Modal Body */}
				<form onSubmit={handleSubmit} style={{ padding: "20px", display: "flex", flexDirection: "column", gap: "16px" }}>
					{displayError && (
						<div
							style={{
								padding: "10px 12px",
								backgroundColor: "rgba(239, 68, 68, 0.1)",
								border: "1px solid rgba(239, 68, 68, 0.3)",
								borderRadius: "6px",
								color: "#F87171",
								fontSize: "13px",
							}}
						>
							⚠️ {displayError}
						</div>
					)}

					{/* Read-only Context Banner */}
					<div
						style={{
							padding: "10px 12px",
							backgroundColor: "rgba(37, 99, 235, 0.08)",
							border: "1px solid rgba(37, 99, 235, 0.2)",
							borderRadius: "6px",
							fontSize: "12px",
							color: "#C9D1D9",
							display: "flex",
							flexDirection: "column",
							gap: "4px",
						}}
					>
						<div>
							<span style={{ color: "#8B949E" }}>Workspace:</span> <strong>{activeWorkspace?.name || "N/A"}</strong>
						</div>
						<div>
							<span style={{ color: "#8B949E" }}>Channel:</span> <strong>{activeChannel?.title || "N/A"}</strong>
						</div>
						<div>
							<span style={{ color: "#8B949E" }}>Thread:</span> <strong>{activeThread?.title || "N/A"}</strong>
						</div>
					</div>

					{/* TrustGroup ID Field */}
					<div>
						<label style={{ display: "block", fontSize: "12px", color: "#8B949E", marginBottom: "6px" }}>
							TrustGroup ID *
						</label>
						<input
							type="text"
							placeholder="tg_legal_counsel"
							value={trustGroupId}
							onChange={(e) => setTrustGroupId(e.target.value)}
							disabled={isLoading}
							style={{
								width: "100%",
								padding: "8px 12px",
								backgroundColor: "#0D1117",
								border: "1px solid #30363D",
								borderRadius: "6px",
								color: "#F0F6FC",
								fontSize: "13px",
								fontFamily: "monospace",
								boxSizing: "border-box",
							}}
						/>
					</div>

					{/* Target Vault ID Field */}
					<div>
						<label style={{ display: "block", fontSize: "12px", color: "#8B949E", marginBottom: "6px" }}>
							Target Member Vault ID *
						</label>
						<input
							type="text"
							placeholder="vault_partner_02"
							value={targetVaultId}
							onChange={(e) => setTargetVaultId(e.target.value)}
							disabled={isLoading}
							style={{
								width: "100%",
								padding: "8px 12px",
								backgroundColor: "#0D1117",
								border: "1px solid #30363D",
								borderRadius: "6px",
								color: "#F0F6FC",
								fontSize: "13px",
								fontFamily: "monospace",
								boxSizing: "border-box",
							}}
						/>
					</div>

					{/* Asset CID Field */}
					<div>
						<label style={{ display: "block", fontSize: "12px", color: "#8B949E", marginBottom: "6px" }}>
							Asset CID (Reference Hash) *
						</label>
						<input
							type="text"
							placeholder="QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco"
							value={assetCid}
							onChange={(e) => setAssetCid(e.target.value)}
							disabled={isLoading}
							style={{
								width: "100%",
								padding: "8px 12px",
								backgroundColor: "#0D1117",
								border: "1px solid #30363D",
								borderRadius: "6px",
								color: "#F0F6FC",
								fontSize: "13px",
								fontFamily: "monospace",
								boxSizing: "border-box",
							}}
						/>
					</div>

					{/* Notes Field */}
					<div>
						<label style={{ display: "block", fontSize: "12px", color: "#8B949E", marginBottom: "6px" }}>
							Collaboration Notes (Optional)
						</label>
						<textarea
							placeholder="Initiating joint contract draft review with legal partner..."
							value={notes}
							onChange={(e) => setNotes(e.target.value)}
							disabled={isLoading}
							rows={3}
							style={{
								width: "100%",
								padding: "8px 12px",
								backgroundColor: "#0D1117",
								border: "1px solid #30363D",
								borderRadius: "6px",
								color: "#F0F6FC",
								fontSize: "13px",
								boxSizing: "border-box",
								resize: "vertical",
							}}
						/>
					</div>

					{/* Modal Footer */}
					<div
						style={{
							marginTop: "8px",
							display: "flex",
							alignItems: "center",
							justifyContent: "flex-end",
							gap: "10px",
						}}
					>
						<button
							type="button"
							onClick={handleClose}
							disabled={isLoading}
							style={{
								padding: "8px 16px",
								backgroundColor: "transparent",
								border: "1px solid #30363D",
								borderRadius: "6px",
								color: "#C9D1D9",
								fontSize: "13px",
								fontWeight: 500,
								cursor: "pointer",
							}}
						>
							Cancel
						</button>
						<button
							type="submit"
							disabled={isLoading || !activeThreadId}
							style={{
								padding: "8px 16px",
								backgroundColor: "#238636",
								border: "none",
								borderRadius: "6px",
								color: "#FFFFFF",
								fontSize: "13px",
								fontWeight: 600,
								cursor: isLoading || !activeThreadId ? "not-allowed" : "pointer",
								opacity: isLoading || !activeThreadId ? 0.6 : 1,
							}}
						>
							{isLoading ? "Creating Share..." : "Create Collaborative Share"}
						</button>
					</div>
				</form>
			</div>
		</div>
	);
};

import React, { useState } from "react";
import { appendThreadEvent, ThreadEventResponse } from "@/services/api";

interface AppendThreadEventModalProps {
	isOpen: boolean;
	activeWorkspaceName?: string;
	activeChannelTitle?: string;
	activeThreadId: string | null;
	activeThreadTitle?: string;
	onClose: () => void;
	onEventAppended?: (event: ThreadEventResponse) => void;
}

export const AppendThreadEventModal: React.FC<AppendThreadEventModalProps> = ({
	isOpen,
	activeWorkspaceName = "Active Workspace",
	activeChannelTitle = "Active Channel",
	activeThreadId,
	activeThreadTitle = "Active Thread",
	onClose,
	onEventAppended,
}) => {
	const [eventType, setEventType] = useState("entry.shared");
	const [customType, setCustomType] = useState("");
	const [notes, setNotes] = useState("");
	const [attachPayloadRef, setAttachPayloadRef] = useState(false);
	const [cid, setCid] = useState("");
	const [contentHash, setContentHash] = useState("");
	const [sizeBytes, setSizeBytes] = useState<number | "">(1048576);
	const [assetName, setAssetName] = useState("");

	const [attachShareEntryRef, setAttachShareEntryRef] = useState(false);
	const [shareEntryId, setShareEntryId] = useState("");
	const [trustGroupId, setTrustGroupId] = useState("");

	const [isLoading, setIsLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);

	if (!isOpen) return null;

	const handleReset = () => {
		setEventType("entry.shared");
		setCustomType("");
		setNotes("");
		setAttachPayloadRef(false);
		setCid("");
		setContentHash("");
		setSizeBytes(1048576);
		setAssetName("");
		setAttachShareEntryRef(false);
		setShareEntryId("");
		setTrustGroupId("");
		setError(null);
		setIsLoading(false);
	};

	const handleClose = () => {
		if (isLoading) return;
		handleReset();
		onClose();
	};

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		setError(null);

		const finalType = (eventType === "custom" ? customType : eventType).trim();
		if (!finalType) {
			setError("Event type is required.");
			return;
		}

		if (!activeThreadId) {
			setError("No active thread selected. Please select a thread before appending an event.");
			return;
		}

		setIsLoading(true);

		try {
			const payloadData: Record<string, any> = {
				notes: notes.trim(),
				timestamp: new Date().toISOString(),
			};

			if (attachPayloadRef && cid.trim()) {
				payloadData.payload_ref = {
					cid: cid.trim(),
					content_hash: contentHash.trim() || "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
					size: Number(sizeBytes) || 0,
					asset_type: finalType,
					name: assetName.trim() || undefined,
				};
			}

			if (attachShareEntryRef) {
				payloadData.share_entry_ref = {
					share_entry_id: shareEntryId.trim() || `se_${Date.now()}`,
					trust_group_id: trustGroupId.trim() || "tg_legal_counsel",
					asset_cid: cid.trim() || undefined,
					created_by: "active_user",
					status: "active",
					created_at: new Date().toISOString(),
				};
			}

			const newEvent = await appendThreadEvent({
				thread_id: activeThreadId,
				type: finalType,
				payload: payloadData,
			});

			handleReset();
			if (onEventAppended) {
				onEventAppended(newEvent);
			}
			onClose();
		} catch (err: any) {
			console.error("Failed to append thread event:", err);
			setError(err?.message || "An unexpected error occurred while appending the thread event.");
		} finally {
			setIsLoading(false);
		}
	};

	return (
		<div
			className="modal-overlay"
			style={{
				position: "fixed",
				top: 0,
				left: 0,
				right: 0,
				bottom: 0,
				backgroundColor: "rgba(0, 0, 0, 0.65)",
				backdropFilter: "blur(4px)",
				display: "flex",
				alignItems: "center",
				justifyContent: "center",
				zIndex: 1000,
			}}
		>
			<div
				className="modal"
				style={{
					width: "100%",
					maxWidth: "500px",
					backgroundColor: "#161B22",
					borderRadius: "8px",
					border: "1px solid rgba(255, 255, 255, 0.1)",
					boxShadow: "0 20px 25px -5px rgba(0, 0, 0, 0.5), 0 10px 10px -5px rgba(0, 0, 0, 0.4)",
					color: "#F0F6FC",
					fontFamily: "Inter, -apple-system, BlinkMacSystemFont, sans-serif",
					overflow: "hidden",
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
					<div>
						<div
							className="mh-title"
							style={{
								fontSize: "16px",
								fontWeight: 600,
								color: "#F0F6FC",
							}}
						>
							Append Thread Event
						</div>
						<div
							className="mh-sub"
							style={{
								fontSize: "12px",
								color: "#8B949E",
								marginTop: "2px",
							}}
						>
							Append an event to the authoritative thread timeline
						</div>
					</div>
					<button
						type="button"
						className="mh-close"
						onClick={handleClose}
						disabled={isLoading}
						style={{
							background: "none",
							border: "none",
							color: "#8B949E",
							fontSize: "16px",
							cursor: isLoading ? "not-allowed" : "pointer",
							padding: "4px",
						}}
					>
						✕
					</button>
				</div>

				{/* Modal Body */}
				<form onSubmit={handleSubmit}>
					<div
						className="modal-body"
						style={{
							padding: "20px",
							display: "flex",
							flexDirection: "column",
							gap: "16px",
						}}
					>
						{/* Error Alert */}
						{error && (
							<div
								style={{
									backgroundColor: "rgba(239, 68, 68, 0.15)",
									border: "1px solid rgba(239, 68, 68, 0.4)",
									borderRadius: "6px",
									padding: "10px 12px",
									color: "#F87171",
									fontSize: "13px",
									display: "flex",
									alignItems: "flex-start",
									gap: "8px",
								}}
							>
								<span>⚠️</span>
								<div style={{ flex: 1 }}>{error}</div>
							</div>
						)}

						{/* Read-only Context */}
						<div
							style={{
								padding: "12px",
								backgroundColor: "rgba(255, 255, 255, 0.04)",
								border: "1px solid rgba(255, 255, 255, 0.08)",
								borderRadius: "6px",
								fontSize: "13px",
								display: "flex",
								flexDirection: "column",
								gap: "6px",
							}}
						>
							<div style={{ display: "flex", justifyContent: "space-between", color: "#8B949E" }}>
								<span>Workspace:</span>
								<strong style={{ color: "#C9D1D9", fontWeight: 600 }}>{activeWorkspaceName}</strong>
							</div>
							<div style={{ display: "flex", justifyContent: "space-between", color: "#8B949E" }}>
								<span>Channel:</span>
								<strong style={{ color: "#C9D1D9", fontWeight: 600 }}>{activeChannelTitle}</strong>
							</div>
							<div style={{ display: "flex", justifyContent: "space-between", color: "#8B949E" }}>
								<span>Target Thread:</span>
								<strong style={{ color: "#2563EB", fontWeight: 600 }}>{activeThreadTitle}</strong>
							</div>
						</div>

						{/* Event Type Selection */}
						<div>
							<label
								htmlFor="event-type-select"
								style={{
									display: "block",
									fontSize: "12px",
									fontWeight: 600,
									color: "#C9D1D9",
									marginBottom: "6px",
									textTransform: "uppercase",
									letterSpacing: "0.5px",
								}}
							>
								Event Type <span style={{ color: "#EF4444" }}>*</span>
							</label>
							<select
								id="event-type-select"
								value={eventType}
								onChange={(e) => setEventType(e.target.value)}
								disabled={isLoading}
								style={{
									width: "100%",
									padding: "10px 12px",
									backgroundColor: "#0D1117",
									border: "1px solid #30363D",
									borderRadius: "6px",
									color: "#F0F6FC",
									fontSize: "14px",
									outline: "none",
									boxSizing: "border-box",
								}}
							>
								<option value="entry.shared">entry.shared</option>
								<option value="invoice.created">invoice.created</option>
								<option value="finance.approved">finance.approved</option>
								<option value="payment.released">payment.released</option>
								<option value="receipt.issued">receipt.issued</option>
								<option value="thread.event.appended">thread.event.appended</option>
								<option value="custom">Custom Event Type…</option>
							</select>
						</div>

						{/* Custom Event Type Input */}
						{eventType === "custom" && (
							<div>
								<label
									htmlFor="custom-event-type-input"
									style={{
										display: "block",
										fontSize: "12px",
										fontWeight: 600,
										color: "#C9D1D9",
										marginBottom: "6px",
										textTransform: "uppercase",
										letterSpacing: "0.5px",
									}}
								>
									Custom Type Identifier
								</label>
								<input
									id="custom-event-type-input"
									type="text"
									placeholder="e.g. audit.completed"
									value={customType}
									onChange={(e) => setCustomType(e.target.value)}
									disabled={isLoading}
									autoFocus
									style={{
										width: "100%",
										padding: "10px 12px",
										backgroundColor: "#0D1117",
										border: "1px solid #30363D",
										borderRadius: "6px",
										color: "#F0F6FC",
										fontSize: "14px",
										outline: "none",
										boxSizing: "border-box",
									}}
								/>
							</div>
						)}

						{/* Event Notes / Summary Input */}
						<div>
							<label
								htmlFor="event-notes-input"
								style={{
									display: "block",
									fontSize: "12px",
									fontWeight: 600,
									color: "#C9D1D9",
									marginBottom: "6px",
									textTransform: "uppercase",
									letterSpacing: "0.5px",
								}}
							>
								Notes / Non-Sensitive Summary
							</label>
							<textarea
								id="event-notes-input"
								placeholder="e.g. Countersigned contract draft committed to vault."
								value={notes}
								onChange={(e) => setNotes(e.target.value)}
								disabled={isLoading}
								rows={2}
								style={{
									width: "100%",
									padding: "10px 12px",
									backgroundColor: "#0D1117",
									border: "1px solid #30363D",
									borderRadius: "6px",
									color: "#F0F6FC",
									fontSize: "14px",
									outline: "none",
									boxSizing: "border-box",
									resize: "vertical",
								}}
							/>
						</div>

						{/* PayloadRef Toggle & Inputs */}
						<div
							style={{
								padding: "12px",
								backgroundColor: "rgba(255, 255, 255, 0.02)",
								border: "1px solid rgba(255, 255, 255, 0.08)",
								borderRadius: "6px",
								display: "flex",
								flexDirection: "column",
								gap: "10px",
							}}
						>
							<label style={{ display: "flex", alignItems: "center", gap: "8px", cursor: "pointer", fontSize: "13px", color: "#C9D1D9", fontWeight: 600 }}>
								<input
									type="checkbox"
									checked={attachPayloadRef}
									onChange={(e) => setAttachPayloadRef(e.target.checked)}
									disabled={isLoading}
								/>
								<span>Attach Safe Asset Reference (PayloadRef)</span>
							</label>

							{attachPayloadRef && (
								<div style={{ display: "flex", flexDirection: "column", gap: "10px", marginTop: "4px" }}>
									<div>
										<label style={{ display: "block", fontSize: "11px", color: "#8B949E", marginBottom: "4px" }}>
											Asset CID (IPFS / Storage Reference) *
										</label>
										<input
											type="text"
											placeholder="e.g. QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco"
											value={cid}
											onChange={(e) => setCid(e.target.value)}
											disabled={isLoading}
											style={{
												width: "100%",
												padding: "8px 10px",
												backgroundColor: "#0D1117",
												border: "1px solid #30363D",
												borderRadius: "4px",
												color: "#F0F6FC",
												fontSize: "13px",
												fontFamily: "monospace",
												boxSizing: "border-box",
											}}
										/>
									</div>

									<div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "10px" }}>
										<div>
											<label style={{ display: "block", fontSize: "11px", color: "#8B949E", marginBottom: "4px" }}>
												File Name (Optional)
											</label>
											<input
												type="text"
												placeholder="contract_draft.pdf"
												value={assetName}
												onChange={(e) => setAssetName(e.target.value)}
												disabled={isLoading}
												style={{
													width: "100%",
													padding: "8px 10px",
													backgroundColor: "#0D1117",
													border: "1px solid #30363D",
													borderRadius: "4px",
													color: "#F0F6FC",
													fontSize: "13px",
													boxSizing: "border-box",
												}}
											/>
										</div>
										<div>
											<label style={{ display: "block", fontSize: "11px", color: "#8B949E", marginBottom: "4px" }}>
												Size (Bytes)
											</label>
											<input
												type="number"
												placeholder="1048576"
												value={sizeBytes}
												onChange={(e) => setSizeBytes(e.target.value === "" ? "" : Number(e.target.value))}
												disabled={isLoading}
												style={{
													width: "100%",
													padding: "8px 10px",
													backgroundColor: "#0D1117",
													border: "1px solid #30363D",
													borderRadius: "4px",
													color: "#F0F6FC",
													fontSize: "13px",
													boxSizing: "border-box",
												}}
											/>
										</div>
									</div>
								</div>
							)}
						</div>

						{/* ShareEntryRef Toggle & Inputs */}
						<div
							style={{
								padding: "12px",
								backgroundColor: "rgba(255, 255, 255, 0.02)",
								border: "1px solid rgba(255, 255, 255, 0.08)",
								borderRadius: "6px",
								display: "flex",
								flexDirection: "column",
								gap: "10px",
							}}
						>
							<label style={{ display: "flex", alignItems: "center", gap: "8px", cursor: "pointer", fontSize: "13px", color: "#C9D1D9", fontWeight: 600 }}>
								<input
									type="checkbox"
									checked={attachShareEntryRef}
									onChange={(e) => setAttachShareEntryRef(e.target.checked)}
									disabled={isLoading}
								/>
								<span>Attach ShareEntry Collaboration Reference</span>
							</label>

							{attachShareEntryRef && (
								<div style={{ display: "flex", flexDirection: "column", gap: "10px", marginTop: "4px" }}>
									<div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "10px" }}>
										<div>
											<label style={{ display: "block", fontSize: "11px", color: "#8B949E", marginBottom: "4px" }}>
												ShareEntry ID (Optional)
											</label>
											<input
												type="text"
												placeholder="se_contract_01"
												value={shareEntryId}
												onChange={(e) => setShareEntryId(e.target.value)}
												disabled={isLoading}
												style={{
													width: "100%",
													padding: "8px 10px",
													backgroundColor: "#0D1117",
													border: "1px solid #30363D",
													borderRadius: "4px",
													color: "#F0F6FC",
													fontSize: "13px",
													fontFamily: "monospace",
													boxSizing: "border-box",
												}}
											/>
										</div>
										<div>
											<label style={{ display: "block", fontSize: "11px", color: "#8B949E", marginBottom: "4px" }}>
												TrustGroup ID (Optional)
											</label>
											<input
												type="text"
												placeholder="tg_legal_counsel"
												value={trustGroupId}
												onChange={(e) => setTrustGroupId(e.target.value)}
												disabled={isLoading}
												style={{
													width: "100%",
													padding: "8px 10px",
													backgroundColor: "#0D1117",
													border: "1px solid #30363D",
													borderRadius: "4px",
													color: "#F0F6FC",
													fontSize: "13px",
													fontFamily: "monospace",
													boxSizing: "border-box",
												}}
											/>
										</div>
									</div>
								</div>
							)}
						</div>
					</div>

					{/* Modal Footer */}
					<div
						className="modal-footer"
						style={{
							padding: "12px 20px",
							borderTop: "1px solid rgba(255, 255, 255, 0.08)",
							display: "flex",
							alignItems: "center",
							justifyContent: "flex-end",
							gap: "10px",
						}}
					>
						<button
							type="button"
							className="btn btn-ghost"
							onClick={handleClose}
							disabled={isLoading}
							style={{
								padding: "8px 16px",
								borderRadius: "6px",
								background: "transparent",
								border: "1px solid rgba(255, 255, 255, 0.15)",
								color: "#C9D1D9",
								fontSize: "13px",
								cursor: isLoading ? "not-allowed" : "pointer",
							}}
						>
							Cancel
						</button>
						<button
							type="submit"
							className="btn btn-primary"
							disabled={isLoading}
							style={{
								padding: "8px 18px",
								borderRadius: "6px",
								background: isLoading ? "#238636a0" : "#238636",
								border: "none",
								color: "#FFFFFF",
								fontSize: "13px",
								fontWeight: 600,
								cursor: isLoading ? "not-allowed" : "pointer",
								display: "flex",
								alignItems: "center",
								gap: "6px",
							}}
						>
							{isLoading ? (
								<>
									<span
										style={{
											display: "inline-block",
											width: "12px",
											height: "12px",
											border: "2px solid #ffffff",
											borderTopColor: "transparent",
											borderRadius: "50%",
											animation: "spin 0.8s linear infinite",
										}}
									/>
									<span>Appending…</span>
								</>
							) : (
								"Append Event"
							)}
						</button>
					</div>
				</form>
			</div>
		</div>
	);
};

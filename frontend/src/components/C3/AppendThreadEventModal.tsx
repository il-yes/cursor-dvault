import React, { useState } from "react";
import { appendThreadEvent, ThreadEventResponse } from "@/services/api";

export interface AppendThreadEventSlidingViewProps {
	isOpen: boolean;
	activeWorkspaceName?: string;
	activeChannelTitle?: string;
	activeThreadId: string | null;
	activeThreadTitle?: string;
	onClose: () => void;
	onEventAppended?: (event: ThreadEventResponse) => void;
}

export const AppendThreadEventSlidingView: React.FC<AppendThreadEventSlidingViewProps> = ({
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

	// PayloadRef state
	const [attachPayloadRef, setAttachPayloadRef] = useState(false);
	const [cid, setCid] = useState("");
	const [contentHash, setContentHash] = useState("");
	const [sizeBytes, setSizeBytes] = useState<number | "">(1048576);
	const [assetName, setAssetName] = useState("");

	// ShareEntryRef state
	const [attachShareEntryRef, setAttachShareEntryRef] = useState(false);
	const [shareEntryId, setShareEntryId] = useState("");
	const [trustGroupId, setTrustGroupId] = useState("");

	// TrustGroupRef state
	const [attachTrustGroupRef, setAttachTrustGroupRef] = useState(false);
	const [tgName, setTgName] = useState("");
	const [memberVaults, setMemberVaults] = useState("vault_legal_01, vault_finance_02");

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
		setAttachTrustGroupRef(false);
		setTgName("");
		setMemberVaults("vault_legal_01, vault_finance_02");
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

			if (attachTrustGroupRef) {
				const membersList = memberVaults
					.split(",")
					.map((s) => s.trim())
					.filter(Boolean)
					.map((vId) => ({ vault_id: vId, role: "Member", joined_at: new Date().toISOString() }));

				payloadData.trust_group_ref = {
					id: trustGroupId.trim() || `tg_${Date.now()}`,
					name: tgName.trim() || "Legal & Compliance Trust Group",
					status: "active",
					member_count: membersList.length,
					members: membersList,
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

	const sectionBoxStyle: React.CSSProperties = {
		padding: "12px 14px",
		backgroundColor: "#fafafa",
		border: "1px solid #ebebeb",
		borderRadius: "8px",
		display: "flex",
		flexDirection: "column",
		gap: "10px",
	};

	return (
		<>
			{/* Backdrop */}
			<div className="c3-sliding-view-backdrop" onClick={handleClose} />

			{/* Fixed Right-Side 460px Panel */}
			<div className="c3-sliding-view-container">
				<div className="slide-panel">
					{/* Header */}
					<div className="sp-header">
						<div className="sp-header-row">
							<div>
								<div className="sp-title">Append Thread Event</div>
								<div className="sp-subtitle">
									Append an event to the authoritative thread timeline
								</div>
							</div>
							<div
								className="sp-close"
								onClick={handleClose}
								role="button"
								tabIndex={0}
							>
								✕
							</div>
						</div>
					</div>

					{/* Form Body */}
					<form onSubmit={handleSubmit} style={{ display: "contents" }}>
						<div className="sp-body">
							{/* Error Alert */}
							{error && (
								<div
									style={{
										backgroundColor: "rgba(239, 68, 68, 0.08)",
										border: "1px solid rgba(239, 68, 68, 0.3)",
										borderRadius: "6px",
										padding: "10px 12px",
										color: "#DC2626",
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

							{/* Hierarchy Context Display */}
							<div>
								<div className="fl">Target Thread Context</div>
								<div className="channel-flow-box">
									<div className="cfb-row" style={{ fontSize: "12px", color: "#333" }}>
										<span style={{ color: "#888" }}>Workspace:</span>
										<strong>{activeWorkspaceName}</strong>
										<span style={{ color: "#ccc" }}>→</span>
										<span style={{ color: "#888" }}>Channel:</span>
										<strong>{activeChannelTitle}</strong>
									</div>
									<div style={{ marginTop: "6px", fontSize: "13px", fontWeight: 600, color: "#C8922A" }}>
										Thread: {activeThreadTitle}
									</div>
								</div>
							</div>

							{/* Event Type Select */}
							<div>
								<div className="fl">
									Event Type <span style={{ color: "#EF4444" }}>*</span>
								</div>
								<select
									className="prop-input"
									value={eventType}
									onChange={(e) => setEventType(e.target.value)}
									disabled={isLoading}
									style={{ height: "38px" }}
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
									<div className="fl">Custom Type Identifier</div>
									<input
										className="prop-input"
										type="text"
										placeholder="e.g. audit.completed"
										value={customType}
										onChange={(e) => setCustomType(e.target.value)}
										disabled={isLoading}
										autoFocus
									/>
								</div>
							)}

							{/* Notes Input */}
							<div>
								<div className="fl">
									Notes / Non-Sensitive Summary
								</div>
								<textarea
									className="prop-input"
									placeholder="e.g. Countersigned contract draft committed to vault."
									value={notes}
									onChange={(e) => setNotes(e.target.value)}
									disabled={isLoading}
									rows={3}
									style={{ resize: "vertical" }}
								/>
							</div>

							{/* PayloadRef Toggle & Section */}
							<div style={sectionBoxStyle}>
								<label style={{ display: "flex", alignItems: "center", gap: "8px", cursor: "pointer", fontSize: "12px", color: "#333", fontWeight: 600 }}>
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
											<div className="prop-label">Asset CID (IPFS / Storage Ref) *</div>
											<input
												className="prop-input"
												type="text"
												placeholder="e.g. QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco"
												value={cid}
												onChange={(e) => setCid(e.target.value)}
												disabled={isLoading}
												style={{ fontFamily: "SF Mono, monospace" }}
											/>
										</div>

										<div style={{ display: "flex", gap: "10px" }}>
											<div style={{ flex: 1 }}>
												<div className="prop-label">File Name (Optional)</div>
												<input
													className="prop-input"
													type="text"
													placeholder="contract_draft.pdf"
													value={assetName}
													onChange={(e) => setAssetName(e.target.value)}
													disabled={isLoading}
												/>
											</div>
											<div style={{ flex: 1 }}>
												<div className="prop-label">Size (Bytes)</div>
												<input
													className="prop-input"
													type="number"
													placeholder="1048576"
													value={sizeBytes}
													onChange={(e) => setSizeBytes(e.target.value === "" ? "" : Number(e.target.value))}
													disabled={isLoading}
												/>
											</div>
										</div>
									</div>
								)}
							</div>

							{/* ShareEntryRef Toggle & Section */}
							<div style={sectionBoxStyle}>
								<label style={{ display: "flex", alignItems: "center", gap: "8px", cursor: "pointer", fontSize: "12px", color: "#333", fontWeight: 600 }}>
									<input
										type="checkbox"
										checked={attachShareEntryRef}
										onChange={(e) => setAttachShareEntryRef(e.target.checked)}
										disabled={isLoading}
									/>
									<span>Attach ShareEntry Collaboration Reference</span>
								</label>

								{attachShareEntryRef && (
									<div style={{ display: "flex", gap: "10px", marginTop: "4px" }}>
										<div style={{ flex: 1 }}>
											<div className="prop-label">ShareEntry ID (Optional)</div>
											<input
												className="prop-input"
												type="text"
												placeholder="se_contract_01"
												value={shareEntryId}
												onChange={(e) => setShareEntryId(e.target.value)}
												disabled={isLoading}
												style={{ fontFamily: "SF Mono, monospace" }}
											/>
										</div>
										<div style={{ flex: 1 }}>
											<div className="prop-label">TrustGroup ID (Optional)</div>
											<input
												className="prop-input"
												type="text"
												placeholder="tg_legal_counsel"
												value={trustGroupId}
												onChange={(e) => setTrustGroupId(e.target.value)}
												disabled={isLoading}
												style={{ fontFamily: "SF Mono, monospace" }}
											/>
										</div>
									</div>
								)}
							</div>

							{/* TrustGroupRef Toggle & Section */}
							<div style={sectionBoxStyle}>
								<label style={{ display: "flex", alignItems: "center", gap: "8px", cursor: "pointer", fontSize: "12px", color: "#333", fontWeight: 600 }}>
									<input
										type="checkbox"
										checked={attachTrustGroupRef}
										onChange={(e) => setAttachTrustGroupRef(e.target.checked)}
										disabled={isLoading}
									/>
									<span>Attach TrustGroup Governance Context</span>
								</label>

								{attachTrustGroupRef && (
									<div style={{ display: "flex", flexDirection: "column", gap: "10px", marginTop: "4px" }}>
										<div>
											<div className="prop-label">TrustGroup Name</div>
											<input
												className="prop-input"
												type="text"
												placeholder="Legal & Compliance Trust Group"
												value={tgName}
												onChange={(e) => setTgName(e.target.value)}
												disabled={isLoading}
											/>
										</div>

										<div>
											<div className="prop-label">Member Vault IDs (Comma Separated)</div>
											<input
												className="prop-input"
												type="text"
												placeholder="vault_legal_01, vault_finance_02"
												value={memberVaults}
												onChange={(e) => setMemberVaults(e.target.value)}
												disabled={isLoading}
												style={{ fontFamily: "SF Mono, monospace" }}
											/>
										</div>
									</div>
								)}
							</div>
						</div>

						{/* Footer */}
						<div className="sp-footer">
							<button
								type="submit"
								className="start-btn"
								disabled={isLoading}
								style={{
									opacity: isLoading ? 0.6 : 1,
									cursor: isLoading ? "not-allowed" : "pointer",
								}}
							>
								{isLoading ? (
									<>
										<span
											style={{
												display: "inline-block",
												width: "14px",
												height: "14px",
												border: "2px solid #fff",
												borderTopColor: "transparent",
												borderRadius: "50%",
												animation: "spin 0.8s linear infinite",
											}}
										/>
										<span>Appending Event…</span>
									</>
								) : (
									"▶ Append Event"
								)}
							</button>
							<div className="footer-note">
								Event will be appended to the thread timeline in cursor order.
							</div>
						</div>
					</form>
				</div>
			</div>
		</>
	);
};

export const AppendThreadEventModal = AppendThreadEventSlidingView;

import React, { useState, useMemo, useEffect } from "react";
import { createPortal } from "react-dom";
import { appendThreadEvent, ThreadEventResponse } from "@/services/api";
import { useVaultStore } from "@/store/vaultStore";
import { VaultEntry } from "@/types/vault";

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
	const [notes, setNotes] = useState("");

	// Source of Available Entries: Session Vault (useVaultStore)
	const vaultContext = useVaultStore((state) => state.vault);

	const allVaultEntries = useMemo(() => {
		if (!vaultContext?.Vault?.entries) return [];
		const entries: VaultEntry[] = [
			...(vaultContext.Vault.entries.login || []),
			...(vaultContext.Vault.entries.card || []),
			...(vaultContext.Vault.entries.note || []),
			...(vaultContext.Vault.entries.sshkey || []),
			...(vaultContext.Vault.entries.identity || []),
		];
		return entries;
	}, [vaultContext]);

	const [selectedEntryId, setSelectedEntryId] = useState<string>("");

	// Auto-select first available entry if none selected
	useEffect(() => {
		if (allVaultEntries.length > 0 && !selectedEntryId) {
			setSelectedEntryId(allVaultEntries[0].id);
		}
	}, [allVaultEntries, selectedEntryId]);

	// Status States
	const [isLoading, setIsLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [successMsg, setSuccessMsg] = useState<string | null>(null);

	if (!isOpen) return null;

	const handleClose = () => {
		if (isLoading) return;
		setError(null);
		setSuccessMsg(null);
		onClose();
	};

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		setError(null);
		setSuccessMsg(null);

		const finalType = "entry.shared";

		if (!activeThreadId) {
			setError("No active thread selected. Please select a thread before appending an event.");
			return;
		}

		if (!selectedEntryId) {
			setError("Please select a Vault Entry from your session vault.");
			return;
		}

		const selectedEntry = allVaultEntries.find((entry) => entry.id === selectedEntryId);

		const payloadData: Record<string, any> = {
			ref_type: "vault_entry",
			entry_id: selectedEntryId,
			entry_name: selectedEntry?.entry_name || "Vault Entry",
			entry_type: selectedEntry?.type || "note",
			notes: notes.trim(),
		};

		setIsLoading(true);

		try {
			const newEvent = await appendThreadEvent({
				thread_id: activeThreadId,
				type: finalType,
				payload: payloadData,
			});

			setSuccessMsg("✓ Event created successfully");
			if (onEventAppended) {
				onEventAppended(newEvent);
			}
			// NOTE: Do NOT close panel automatically. Keep open and usable!
		} catch (err: any) {
			console.error("Failed to append thread event:", err);
			setError(err?.message || "Event could not be created");
		} finally {
			setIsLoading(false);
		}
	};

	// Use React Portal to mount under document.body outside parent vaul drawer transform!
	return createPortal(
		<>
			{/* Backdrop */}
			<div className="c3-sliding-view-backdrop" onClick={handleClose} />

			{/* Fixed Right-Side Viewport-Anchored Panel */}
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
							{/* Success Message Alert */}
							{successMsg && (
								<div
									style={{
										backgroundColor: "rgba(16, 185, 129, 0.08)",
										border: "1px solid rgba(16, 185, 129, 0.3)",
										borderRadius: "6px",
										padding: "10px 12px",
										color: "#059669",
										fontSize: "13px",
										fontWeight: 600,
										display: "flex",
										alignItems: "center",
										gap: "8px",
									}}
								>
									<span>{successMsg}</span>
								</div>
							)}

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

							{/* Session Vault Entry Select */}
							<div>
								<div className="fl">
									Vault Entry <span style={{ color: "#EF4444" }}>*</span>
								</div>
								{allVaultEntries.length > 0 ? (
									<select
										className="prop-input"
										value={selectedEntryId}
										onChange={(e) => setSelectedEntryId(e.target.value)}
										disabled={isLoading}
										style={{ height: "38px" }}
									>
										{allVaultEntries.map((entry) => (
											<option key={entry.id} value={entry.id}>
												[{entry.type.toUpperCase()}] {entry.entry_name || "Untitled Entry"}
											</option>
										))}
									</select>
								) : (
									<div style={{ fontSize: "12px", color: "#B45309", backgroundColor: "rgba(245, 158, 11, 0.08)", border: "1px solid rgba(245, 158, 11, 0.3)", borderRadius: "6px", padding: "8px 10px" }}>
										No Vault Entries found in current session vault. Please unlock or sync your vault.
									</div>
								)}
							</div>

							{/* Notes Input */}
							<div>
								<div className="fl">
									Notes / Non-Sensitive Summary <span style={{ color: "#999", fontWeight: 400 }}>(optional)</span>
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
						</div>

						{/* Footer */}
						<div className="sp-footer">
							<button
								type="submit"
								className="start-btn"
								disabled={isLoading || allVaultEntries.length === 0}
								style={{
									opacity: isLoading || allVaultEntries.length === 0 ? 0.6 : 1,
									cursor: isLoading || allVaultEntries.length === 0 ? "not-allowed" : "pointer",
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
										<span>Appending...</span>
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
		</>,
		document.body
	);
};

export const AppendThreadEventModal = AppendThreadEventSlidingView;

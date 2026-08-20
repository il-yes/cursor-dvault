import React, { useState } from "react";
import { createThread, ThreadResponse } from "@/services/api";

export interface ThreadSlidingViewProps {
	isOpen: boolean;
	activeWorkspaceName?: string;
	activeChannelId: string | null;
	activeChannelTitle?: string;
	onClose: () => void;
	onThreadCreated?: (thread: ThreadResponse) => void;
}

export const ThreadSlidingView: React.FC<ThreadSlidingViewProps> = ({
	isOpen,
	activeWorkspaceName = "Active Workspace",
	activeChannelId,
	activeChannelTitle = "Contract Execution",
	onClose,
	onThreadCreated,
}) => {
	const [title, setTitle] = useState("");
	const [subtitle, setSubtitle] = useState("");
	const [assetType, setAssetType] = useState("note");
	const [isLoading, setIsLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [validationError, setValidationError] = useState<string | null>(null);

	if (!isOpen) return null;

	const handleReset = () => {
		setTitle("");
		setSubtitle("");
		setAssetType("note");
		setError(null);
		setValidationError(null);
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
		setValidationError(null);

		const trimmedTitle = title.trim();
		if (!trimmedTitle) {
			setValidationError("Thread title is required.");
			return;
		}

		if (!activeChannelId) {
			setError("No active channel selected. Please select a channel before creating a thread.");
			return;
		}

		setIsLoading(true);

		try {
			console.log("[BOUNDARY_LOG] CREATE THREAD CALLING", { channel_id: activeChannelId, title: trimmedTitle });
			const createdThread = await createThread({
				channel_id: activeChannelId,
				title: trimmedTitle,
				subtitle: subtitle.trim(),
				asset_type: assetType,
			});
			console.log("[BOUNDARY_LOG] CREATE THREAD RETURNED", createdThread);

			handleReset();
			if (onThreadCreated) {
				console.log("[BOUNDARY_LOG] CALLING addThread (via onThreadCreated)", createdThread);
				onThreadCreated(createdThread);
			}
			onClose();
		} catch (err: any) {
			console.error("[BOUNDARY_LOG] CREATE THREAD FAILED with error:", err);
			const msg = typeof err === "string" ? err : err?.message || "An unexpected error occurred while creating the thread.";
			setError(msg);
		} finally {
			setIsLoading(false);
		}
	};

	const channelDisplayName = activeChannelId ? (activeChannelTitle || "Selected Channel") : "No Active Channel Selected";
	const channelSlug = activeChannelTitle
		? activeChannelTitle.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/(^-|-$)/g, "")
		: "none";

	return (
		<>
			{/* Static Backdrop Layer */}
			<div className="c3-sliding-view-backdrop" onClick={handleClose} />

			{/* Fixed Right-Side 460px Panel */}
			<div className="c3-sliding-view-container">
				<div className="slide-panel">
					{/* Header */}
					<div className="sp-header">
						<div className="sp-header-row">
							<div>
								<div className="sp-title">
									New Thread <span style={{ background: "#FF0055", color: "#FFF", padding: "2px 8px", borderRadius: "4px", fontSize: "11px", fontWeight: 700 }}>[RUNTIME-MARKER-ACTIVE]</span>
								</div>
								<div className="sp-subtitle">
									Instantiate a channel into a new thread
								</div>
								<div style={{ background: "#1E293B", color: "#38BDF8", padding: "4px 8px", borderRadius: "4px", fontSize: "11px", marginTop: "6px", fontFamily: "monospace" }}>
									activeChannelId: {activeChannelId ? `"${activeChannelId}"` : "NULL (NOT SELECTED)"}
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
							{/* Missing Channel Alert Banner */}
							{!activeChannelId && (
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
									<div style={{ flex: 1 }}>
										<strong>Channel Required:</strong> Please select an active channel from the workspace before creating a thread.
									</div>
								</div>
							)}

							{/* General Error Alert */}
							{error && activeChannelId && (
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

							{/* 1. Channel Display */}
							<div>
								<div className="fl">
									Channel{" "}
									<span className="fl-hint">defines slots, gates, vault flow</span>
								</div>
								<div className="channel-selected">
									<div className="cs-icon-wrap">📄</div>
									<div className="cs-content">
										<div className="cs-name">{channelDisplayName}</div>
										<div className="cs-desc">
											{channelSlug || "contract-execution"} · workspace: {activeWorkspaceName}
										</div>
									</div>
								</div>
							</div>

							{/* 2. Channel Flow Preview */}
							<div>
								<div className="fl">Flow</div>
								<div className="channel-flow-box">
									<div className="cfb-row">
										<div className="cfb-vault">
											<div className="cfb-dot" style={{ background: "#7C3AED" }} />
											vault_legal
										</div>
										<div className="cfb-arrow">→</div>
										<div className="cfb-vault">
											<div className="cfb-dot" style={{ background: "#2563EB" }} />
											vault_finance
										</div>
										<div className="cfb-arrow">→</div>
										<div className="cfb-vault">
											<div className="cfb-dot" style={{ background: "#444" }} />
											vault_direction
										</div>
									</div>
									<div className="cfb-meta">
										<div className="cfb-metaitem">
											<strong>3</strong> slots
										</div>
										<div className="cfb-metaitem">
											<strong>2</strong> gated
										</div>
										<div className="cfb-metaitem">
											first slot: <strong>contract_draft</strong>
										</div>
									</div>
								</div>
							</div>

							{/* 3. Thread Name Input */}
							<div>
								<div className="fl">
									Thread Name <span style={{ color: "#EF4444" }}>*</span>
								</div>
								<div className="thread-name-wrap">
									<div className="thread-name-prefix">{channelSlug || "contract"} —</div>
									<input
										className="thread-name-input"
										type="text"
										placeholder="e.g. Accenture Partnership"
										value={title}
										onChange={(e) => {
											setTitle(e.target.value);
											if (validationError) setValidationError(null);
										}}
										disabled={isLoading}
										autoFocus
									/>
								</div>
								{validationError && (
									<div style={{ color: "#EF4444", fontSize: "12px", marginTop: "4px" }}>
										{validationError}
									</div>
								)}
							</div>

							{/* 4. Subtitle / Summary */}
							<div>
								<div className="fl">
									Subtitle / Summary{" "}
									<span className="fl-hint">optional descriptor</span>
								</div>
								<input
									className="prop-input"
									type="text"
									placeholder="e.g. Master Services Agreement & SOW #4"
									value={subtitle}
									onChange={(e) => setSubtitle(e.target.value)}
									disabled={isLoading}
								/>
							</div>

							{/* 5. Asset Type */}
							<div>
								<div className="fl">Asset Type</div>
								<select
									className="prop-input"
									value={assetType}
									onChange={(e) => setAssetType(e.target.value)}
									disabled={isLoading}
									style={{ height: "38px" }}
								>
									<option value="note">Note / Document</option>
									<option value="login">Login / Identity Credentials</option>
									<option value="card">Card / Payment Instrument</option>
									<option value="ssh_key">SSH Key / Access Key</option>
								</select>
							</div>

							{/* Stellar Anchor Info */}
							<div className="stellar-info">
								<span className="si-icon">✦</span>
								<span className="si-text">
									A <strong>genesis transaction</strong> will be anchored on Stellar
									the moment this thread starts. Every subsequent commit is
									automatically anchored.
								</span>
								<div className="si-status">
									<div className="si-dot" />
									<span className="si-label">Active</span>
								</div>
							</div>
						</div>

						{/* Footer */}
						<div className="sp-footer">
							<button
								type="submit"
								className="start-btn"
								disabled={isLoading || !title.trim() || !activeChannelId}
								style={{
									opacity: isLoading || !title.trim() || !activeChannelId ? 0.6 : 1,
									cursor: isLoading || !title.trim() || !activeChannelId ? "not-allowed" : "pointer",
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
										<span>Starting Thread…</span>
									</>
								) : (
									"▶ Start Thread"
								)}
							</button>
							<div className="footer-note">
								Thread appears in the ledger immediately.{" "}
								<strong>vault_legal</strong> can commit{" "}
								<strong>contract_draft</strong> right away. C3 extension can be
								added at any time.
							</div>
						</div>
					</form>
				</div>
			</div>
		</>
	);
};

export const CreateThreadSlidingView = ThreadSlidingView;
export const CreateThreadModal = ThreadSlidingView;

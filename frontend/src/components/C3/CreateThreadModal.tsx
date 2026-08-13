import React, { useState } from "react";
import { createThread, ThreadResponse } from "@/services/api";

interface CreateThreadModalProps {
	isOpen: boolean;
	activeWorkspaceName?: string;
	activeChannelId: string | null;
	activeChannelTitle?: string;
	onClose: () => void;
	onThreadCreated?: (thread: ThreadResponse) => void;
}

export const CreateThreadModal: React.FC<CreateThreadModalProps> = ({
	isOpen,
	activeWorkspaceName = "Active Workspace",
	activeChannelId,
	activeChannelTitle = "Active Channel",
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
			const createdThread = await createThread({
				channel_id: activeChannelId,
				title: trimmedTitle,
				subtitle: subtitle.trim(),
				asset_type: assetType,
			});

			handleReset();
			if (onThreadCreated) {
				onThreadCreated(createdThread);
			}
			onClose();
		} catch (err: any) {
			console.error("Failed to create thread:", err);
			setError(err?.message || "An unexpected error occurred while creating the thread.");
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
							Create C3 Thread
						</div>
						<div
							className="mh-sub"
							style={{
								fontSize: "12px",
								color: "#8B949E",
								marginTop: "2px",
							}}
						>
							Start a new collaborative thread under your active channel
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
								padding: "10px 12px",
								backgroundColor: "rgba(255, 255, 255, 0.04)",
								border: "1px solid rgba(255, 255, 255, 0.08)",
								borderRadius: "6px",
								fontSize: "13px",
								display: "flex",
								flexDirection: "column",
								gap: "4px",
							}}
						>
							<div style={{ display: "flex", justifyContent: "space-between", color: "#8B949E" }}>
								<span>Workspace:</span>
								<strong style={{ color: "#C9D1D9", fontWeight: 600 }}>{activeWorkspaceName}</strong>
							</div>
							<div style={{ display: "flex", justifyContent: "space-between", color: "#8B949E" }}>
								<span>Target Channel:</span>
								<strong style={{ color: "#2563EB", fontWeight: 600 }}>{activeChannelTitle}</strong>
							</div>
						</div>

						{/* Thread Title Input */}
						<div>
							<label
								htmlFor="thread-title-input"
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
								Thread Title <span style={{ color: "#EF4444" }}>*</span>
							</label>
							<input
								id="thread-title-input"
								type="text"
								placeholder="e.g. Accenture Partnership Contract Execution"
								value={title}
								onChange={(e) => {
									setTitle(e.target.value);
									if (validationError) setValidationError(null);
								}}
								disabled={isLoading}
								autoFocus
								style={{
									width: "100%",
									padding: "10px 12px",
									backgroundColor: "#0D1117",
									border: validationError ? "1px solid #EF4444" : "1px solid #30363D",
									borderRadius: "6px",
									color: "#F0F6FC",
									fontSize: "14px",
									outline: "none",
									boxSizing: "border-box",
								}}
							/>
							{validationError && (
								<div style={{ color: "#EF4444", fontSize: "12px", marginTop: "4px" }}>
									{validationError}
								</div>
							)}
						</div>

						{/* Thread Subtitle Input */}
						<div>
							<label
								htmlFor="thread-subtitle-input"
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
								Subtitle / Summary
							</label>
							<input
								id="thread-subtitle-input"
								type="text"
								placeholder="e.g. Master Services Agreement & Statement of Work #4"
								value={subtitle}
								onChange={(e) => setSubtitle(e.target.value)}
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
							/>
						</div>

						{/* Asset Type Input */}
						<div>
							<label
								htmlFor="thread-asset-type-select"
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
								Asset Type
							</label>
							<select
								id="thread-asset-type-select"
								value={assetType}
								onChange={(e) => setAssetType(e.target.value)}
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
								<option value="note">Note / Document</option>
								<option value="login">Login / Identity Credentials</option>
								<option value="card">Card / Payment Instrument</option>
								<option value="ssh_key">SSH Key / Access Key</option>
							</select>
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
							disabled={isLoading || !title.trim()}
							style={{
								padding: "8px 18px",
								borderRadius: "6px",
								background: isLoading || !title.trim() ? "#238636a0" : "#238636",
								border: "none",
								color: "#FFFFFF",
								fontSize: "13px",
								fontWeight: 600,
								cursor: isLoading || !title.trim() ? "not-allowed" : "pointer",
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
									<span>Creating…</span>
								</>
							) : (
								"Create Thread"
							)}
						</button>
					</div>
				</form>
			</div>
		</div>
	);
};

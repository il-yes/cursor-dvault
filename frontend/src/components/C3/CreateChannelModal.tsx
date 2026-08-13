import React, { useState } from "react";
import { createChannel, ChannelResponse } from "@/services/api";

interface CreateChannelModalProps {
	isOpen: boolean;
	activeWorkspaceId: string | null;
	activeWorkspaceName?: string;
	onClose: () => void;
	onChannelCreated?: (channel: ChannelResponse) => void;
}

export const CreateChannelModal: React.FC<CreateChannelModalProps> = ({
	isOpen,
	activeWorkspaceId,
	activeWorkspaceName = "Active Workspace",
	onClose,
	onChannelCreated,
}) => {
	const [title, setTitle] = useState("");
	const [isLoading, setIsLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [validationError, setValidationError] = useState<string | null>(null);

	if (!isOpen) return null;

	const handleReset = () => {
		setTitle("");
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
			setValidationError("Channel title is required.");
			return;
		}

		if (!activeWorkspaceId) {
			setError("No active workspace selected. Please select or create a workspace first.");
			return;
		}

		setIsLoading(true);

		try {
			const createdChannel = await createChannel({
				title: trimmedTitle,
				workspace_id: activeWorkspaceId,
			});

			handleReset();
			if (onChannelCreated) {
				onChannelCreated(createdChannel);
			}
			onClose();
		} catch (err: any) {
			console.error("Failed to create channel:", err);
			setError(err?.message || "An unexpected error occurred while creating the channel.");
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
					maxWidth: "480px",
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
							Create C3 Channel
						</div>
						<div
							className="mh-sub"
							style={{
								fontSize: "12px",
								color: "#8B949E",
								marginTop: "2px",
							}}
						>
							Add a collaboration channel under your active workspace
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

						{/* Read-only Workspace Context */}
						<div
							style={{
								padding: "10px 12px",
								backgroundColor: "rgba(255, 255, 255, 0.04)",
								border: "1px solid rgba(255, 255, 255, 0.08)",
								borderRadius: "6px",
								fontSize: "13px",
								color: "#8B949E",
								display: "flex",
								alignItems: "center",
								justifyContent: "space-between",
							}}
						>
							<span>Target Workspace:</span>
							<strong style={{ color: "#C9D1D9", fontWeight: 600 }}>
								{activeWorkspaceName}
							</strong>
						</div>

						{/* Channel Title Input */}
						<div>
							<label
								htmlFor="channel-title-input"
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
								Channel Title <span style={{ color: "#EF4444" }}>*</span>
							</label>
							<input
								id="channel-title-input"
								type="text"
								placeholder="e.g. Battery System Engineering, Financial Auditing"
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
								"Create Channel"
							)}
						</button>
					</div>
				</form>
			</div>
		</div>
	);
};

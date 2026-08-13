import React, { useState } from "react";
import { createWorkspace, WorkspaceResponse } from "@/services/api";

interface CreateWorkspaceModalProps {
	isOpen: boolean;
	onClose: () => void;
	onWorkspaceCreated?: (workspace: WorkspaceResponse) => void;
}

export const CreateWorkspaceModal: React.FC<CreateWorkspaceModalProps> = ({
	isOpen,
	onClose,
	onWorkspaceCreated,
}) => {
	const [name, setName] = useState("");
	const [description, setDescription] = useState("");
	const [isLoading, setIsLoading] = useState(false);
	const [error, setError] = useState<string | null>(null);
	const [validationError, setValidationError] = useState<string | null>(null);

	if (!isOpen) return null;

	const handleReset = () => {
		setName("");
		setDescription("");
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

		const trimmedName = name.trim();
		if (!trimmedName) {
			setValidationError("Workspace name is required.");
			return;
		}

		setIsLoading(true);

		try {
			const createdWorkspace = await createWorkspace({
				name: trimmedName,
				description: description.trim(),
			});

			handleReset();
			if (onWorkspaceCreated) {
				onWorkspaceCreated(createdWorkspace);
			}
			onClose();
		} catch (err: any) {
			console.error("Failed to create workspace:", err);
			setError(err?.message || "An unexpected error occurred while creating the workspace.");
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
							Create C3 Workspace
						</div>
						<div
							className="mh-sub"
							style={{
								fontSize: "12px",
								color: "#8B949E",
								marginTop: "2px",
							}}
						>
							Add a top-level workspace for channels and collaborative threads
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

						{/* Workspace Name Input */}
						<div>
							<label
								htmlFor="workspace-name-input"
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
								Workspace Name <span style={{ color: "#EF4444" }}>*</span>
							</label>
							<input
								id="workspace-name-input"
								type="text"
								placeholder="e.g. EVTOL Development Program"
								value={name}
								onChange={(e) => {
									setName(e.target.value);
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

						{/* Workspace Description Input */}
						<div>
							<label
								htmlFor="workspace-desc-input"
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
								Description
							</label>
							<textarea
								id="workspace-desc-input"
								placeholder="e.g. Substrate workspace for joint engineering collaboration."
								value={description}
								onChange={(e) => setDescription(e.target.value)}
								disabled={isLoading}
								rows={3}
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
							disabled={isLoading || !name.trim()}
							style={{
								padding: "8px 18px",
								borderRadius: "6px",
								background: isLoading || !name.trim() ? "#238636a0" : "#238636",
								border: "none",
								color: "#FFFFFF",
								fontSize: "13px",
								fontWeight: 600,
								cursor: isLoading || !name.trim() ? "not-allowed" : "pointer",
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
								"Create Workspace"
							)}
						</button>
					</div>
				</form>
			</div>
		</div>
	);
};

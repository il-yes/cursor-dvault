import React from "react";

interface RevokeChannelConfirmModalProps {
	isOpen: boolean;
	channelTitle: string;
	isRevoking: boolean;
	onCancel: () => void;
	onConfirm: () => void;
}

export const RevokeChannelConfirmModal: React.FC<RevokeChannelConfirmModalProps> = ({
	isOpen,
	channelTitle,
	isRevoking,
	onCancel,
	onConfirm,
}) => {
	if (!isOpen) return null;

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
					maxWidth: "440px",
					backgroundColor: "#161B22",
					borderRadius: "8px",
					border: "1px solid rgba(255, 255, 255, 0.1)",
					boxShadow: "0 20px 25px -5px rgba(0, 0, 0, 0.5), 0 10px 10px -5px rgba(0, 0, 0, 0.4)",
					color: "#F0F6FC",
					fontFamily: "Inter, -apple-system, BlinkMacSystemFont, sans-serif",
					overflow: "hidden",
				}}
			>
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
							style={{ fontSize: "16px", fontWeight: 600, color: "#F0F6FC" }}
						>
							Revoke Channel
						</div>
						<div
							className="mh-sub"
							style={{ fontSize: "12px", color: "#8B949E", marginTop: "2px" }}
						>
							Lifecycle action enforced by the Cloud backend
						</div>
					</div>
					<button
						type="button"
						className="mh-close"
						onClick={onCancel}
						disabled={isRevoking}
						style={{
							background: "none",
							border: "none",
							color: "#8B949E",
							fontSize: "16px",
							cursor: isRevoking ? "not-allowed" : "pointer",
							padding: "4px",
						}}
					>
						✕
					</button>
				</div>

				<div
					className="modal-body"
					style={{
						padding: "20px",
						display: "flex",
						flexDirection: "column",
						gap: "16px",
					}}
				>
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
						<span>Channel:</span>
						<strong style={{ color: "#C9D1D9", fontWeight: 600 }}>{channelTitle}</strong>
					</div>

					<div
						style={{
							padding: "10px 12px",
							backgroundColor: "rgba(239, 68, 68, 0.12)",
							border: "1px solid rgba(239, 68, 68, 0.3)",
							borderRadius: "6px",
							fontSize: "13px",
							color: "#F87171",
							display: "flex",
							alignItems: "flex-start",
							gap: "8px",
						}}
					>
						<span>⚠️</span>
						<div style={{ flex: 1 }}>
							This channel will become <strong>revoked</strong>. The Cloud backend is
							authoritative for the revocation; the channel status will be refreshed
							from the Cloud after the action completes.
						</div>
					</div>
				</div>

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
						onClick={onCancel}
						disabled={isRevoking}
						style={{
							padding: "8px 16px",
							borderRadius: "6px",
							background: "transparent",
							border: "1px solid rgba(255, 255, 255, 0.15)",
							color: "#C9D1D9",
							fontSize: "13px",
							cursor: isRevoking ? "not-allowed" : "pointer",
						}}
					>
						Cancel
					</button>
					<button
						type="button"
						className="btn btn-danger"
						onClick={onConfirm}
						disabled={isRevoking}
						style={{
							padding: "8px 18px",
							borderRadius: "6px",
							background: isRevoking ? "#B91C1Ca0" : "#DC2626",
							border: "none",
							color: "#FFFFFF",
							fontSize: "13px",
							fontWeight: 600,
							cursor: isRevoking ? "not-allowed" : "pointer",
							display: "flex",
							alignItems: "center",
							gap: "6px",
						}}
					>
						{isRevoking ? (
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
								<span>Revoking…</span>
							</>
						) : (
							"Revoke Channel"
						)}
					</button>
				</div>
			</div>
		</div>
	);
};

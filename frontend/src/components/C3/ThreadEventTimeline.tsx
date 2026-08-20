import React from "react";
import { ThreadEventResponse } from "@/services/api";
import { C3ResourceCard } from "./C3ResourceCard";

interface ThreadEventTimelineProps {
	events: ThreadEventResponse[];
	isLoading: boolean;
	error: string | null;
	activeThreadTitle?: string;
	onOpenAppendModal: () => void;
}

export const ThreadEventTimeline: React.FC<ThreadEventTimelineProps> = ({
	events,
	isLoading,
	error,
	activeThreadTitle = "Selected Thread",
	onOpenAppendModal,
}) => {
	if (isLoading) {
		return (
			<div
				style={{
					padding: "24px 16px",
					fontSize: "13px",
					color: "#8B949E",
					textAlign: "center",
				}}
			>
				<span>Loading timeline events…</span>
			</div>
		);
	}

	if (error) {
		return (
			<div
				style={{
					padding: "16px",
					backgroundColor: "rgba(239, 68, 68, 0.12)",
					border: "1px solid rgba(239, 68, 68, 0.3)",
					borderRadius: "6px",
					fontSize: "13px",
					color: "#EF4444",
					margin: "16px",
				}}
			>
				⚠️ {error}
			</div>
		);
	}

	return (
		<div
			className="c3-thread-timeline-panel"
			style={{
				display: "flex",
				flexDirection: "column",
				gap: "16px",
				padding: "20px",
				backgroundColor: "#0D1117",
				borderRadius: "8px",
				border: "1px solid rgba(255, 255, 255, 0.08)",
				fontFamily: "Inter, -apple-system, BlinkMacSystemFont, sans-serif",
			}}
		>
			{/* Timeline Header */}
			<div
				style={{
					display: "flex",
					alignItems: "center",
					justifyContent: "space-between",
					paddingBottom: "12px",
					borderBottom: "1px solid rgba(255, 255, 255, 0.08)",
				}}
			>
				<div>
					<h4 style={{ margin: 0, fontSize: "14px", fontWeight: 600, color: "#F0F6FC" }}>
						Collaborative Timeline
					</h4>
					<span style={{ fontSize: "12px", color: "#8B949E" }}>
						Authoritative event history for {activeThreadTitle} ({events.length})
					</span>
				</div>
				<button
					type="button"
					onClick={onOpenAppendModal}
					className="btn btn-primary"
					style={{
						padding: "6px 14px",
						fontSize: "12px",
						borderRadius: "6px",
						fontWeight: 600,
						background: "#238636",
						border: "none",
						color: "#FFFFFF",
						cursor: "pointer",
					}}
				>
					+ Add Event
				</button>
			</div>

			{/* Empty Event List */}
			{events.length === 0 ? (
				<div
					style={{
						padding: "32px 16px",
						textAlign: "center",
						color: "#8B949E",
						fontSize: "13px",
						backgroundColor: "rgba(255, 255, 255, 0.02)",
						borderRadius: "6px",
						border: "1px dashed rgba(255, 255, 255, 0.08)",
					}}
				>
					<div style={{ fontSize: "24px", marginBottom: "8px" }}>📜</div>
					<p style={{ margin: "0 0 12px 0", color: "#C9D1D9" }}>No events appended to this thread yet.</p>
					<button
						type="button"
						onClick={onOpenAppendModal}
						style={{
							background: "transparent",
							border: "1px solid rgba(255, 255, 255, 0.15)",
							color: "#58A6FF",
							borderRadius: "4px",
							padding: "4px 12px",
							fontSize: "12px",
							cursor: "pointer",
						}}
					>
						+ Append First Event
					</button>
				</div>
			) : (
				/* Event Stream in Cursor Order */
				<div
					style={{
						display: "flex",
						flexDirection: "column",
						gap: "12px",
					}}
				>
					{events.map((evt, idx) => (
						<div
							key={evt.id || idx}
							style={{
								display: "flex",
								alignItems: "flex-start",
								gap: "12px",
								padding: "12px 14px",
								backgroundColor: "rgba(255, 255, 255, 0.03)",
								border: "1px solid rgba(255, 255, 255, 0.06)",
								borderRadius: "6px",
								transition: "background-color 0.15s ease",
							}}
						>
							{/* Cursor Badge */}
							<div
								style={{
									display: "flex",
									alignItems: "center",
									justifyContent: "center",
									width: "24px",
									height: "24px",
									borderRadius: "50%",
									backgroundColor: "rgba(37, 99, 235, 0.2)",
									color: "#58A6FF",
									fontSize: "11px",
									fontWeight: 700,
									flexShrink: 0,
								}}
							>
								{evt.cursor || idx + 1}
							</div>

							{/* Event Content */}
							<div style={{ flex: 1, overflow: "hidden" }}>
								<div style={{ display: "flex", alignItems: "center", justifyContent: "space-between" }}>
									<strong style={{ fontSize: "13px", color: "#F0F6FC" }}>{evt.type}</strong>
									<span style={{ fontSize: "11px", color: "#8B949E" }}>
										{evt.created_at ? new Date(evt.created_at).toLocaleString() : "Just now"}
									</span>
								</div>

								<div style={{ fontSize: "11px", color: "#8B949E", marginTop: "2px", fontFamily: "monospace" }}>
									Event ID: {evt.id ? evt.id.slice(0, 16) + "…" : "generated"}
								</div>

								{/* Payload details if safe metadata exists */}
								{evt.payload && Object.keys(evt.payload).length > 0 && (
									<div
										style={{
											marginTop: "8px",
											padding: "8px 10px",
											backgroundColor: "#161B22",
											border: "1px solid rgba(255, 255, 255, 0.06)",
											borderRadius: "4px",
											fontSize: "12px",
											color: "#C9D1D9",
											display: "flex",
											flexDirection: "column",
											gap: "6px",
										}}
									>
										{evt.payload.notes && (
											<div>
												<span style={{ color: "#8B949E" }}>Note:</span> {evt.payload.notes}
											</div>
										)}

										{/* Safe PayloadRef Reference Metadata Card */}
										{(evt.payload.payload_ref || evt.payload_ref || evt.payload.cid) && (() => {
											const ref = evt.payload.payload_ref || evt.payload_ref || {
												cid: evt.payload.cid,
												content_hash: evt.payload.content_hash,
												size: evt.payload.size,
												name: evt.payload.name,
											};

											const formattedSize = ref.size
												? (ref.size / 1024 / 1024).toFixed(2) + " MB"
												: "Unknown size";

											return (
												<div
													style={{
														marginTop: "4px",
														padding: "8px",
														backgroundColor: "rgba(37, 99, 235, 0.08)",
														border: "1px solid rgba(37, 99, 235, 0.2)",
														borderRadius: "4px",
														display: "flex",
														flexDirection: "column",
														gap: "4px",
														fontSize: "11px",
													}}
												>
													<div style={{ display: "flex", alignItems: "center", justifyContent: "space-between" }}>
														<strong style={{ color: "#58A6FF" }}>
															📦 PayloadRef: {ref.name || "Asset Reference"}
														</strong>
														<span style={{ color: "#8B949E" }}>{formattedSize}</span>
													</div>
													{ref.cid && (
														<div style={{ fontFamily: "monospace", color: "#C9D1D9", wordBreak: "break-all" }}>
															<span style={{ color: "#8B949E" }}>CID:</span> {ref.cid}
														</div>
													)}
													{ref.content_hash && (
														<div style={{ fontFamily: "monospace", color: "#8B949E", wordBreak: "break-all", fontSize: "10px" }}>
															SHA-256: {ref.content_hash}
														</div>
													)}
												</div>
											);
										})()}

										{/* Safe ShareEntry Collaboration Reference Card with Interactive Open Action */}
										{(evt.payload.share_entry_ref || evt.share_entry_ref || evt.payload.share_entry_id) && (() => {
											const shareRef = evt.payload.share_entry_ref || evt.share_entry_ref || {
												share_entry_id: evt.payload.share_entry_id,
												trust_group_id: evt.payload.trust_group_id,
												asset_cid: evt.payload.asset_cid,
												created_by: evt.payload.created_by,
												status: evt.payload.status || "active",
											};

											return (
												<C3ResourceCard
													refType={evt.payload.ref_type || "share_entry"}
													shareEntryId={shareRef.share_entry_id}
													trustGroupId={shareRef.trust_group_id}
													cid={shareRef.asset_cid}
													author={shareRef.created_by}
													createdAt={evt.created_at}
												/>
											);
										})()}

										{/* Safe TrustGroup Governance Context Card */}
										{(evt.payload.trust_group_ref || evt.trust_group_ref || evt.payload.trust_group_id) && (() => {
											const tgRef = evt.payload.trust_group_ref || evt.trust_group_ref || {
												id: evt.payload.trust_group_id,
												name: evt.payload.trust_group_name || "TrustGroup Governance",
												status: evt.payload.status || "active",
												member_count: evt.payload.member_count || (evt.payload.members ? evt.payload.members.length : 1),
												members: evt.payload.members,
											};

											return (
												<div
													style={{
														marginTop: "4px",
														padding: "8px",
														backgroundColor: "rgba(124, 58, 237, 0.08)",
														border: "1px solid rgba(124, 58, 237, 0.2)",
														borderRadius: "4px",
														display: "flex",
														flexDirection: "column",
														gap: "4px",
														fontSize: "11px",
													}}
												>
													<div style={{ display: "flex", alignItems: "center", justifyContent: "space-between" }}>
														<strong style={{ color: "#A78BFA" }}>
															🛡️ TrustGroup: {tgRef.name || tgRef.id}
														</strong>
														<span style={{ color: "#8B949E", fontSize: "10px" }}>
															{tgRef.member_count || 1} Member(s)
														</span>
													</div>
													<div style={{ fontFamily: "monospace", color: "#C9D1D9" }}>
														<span style={{ color: "#8B949E" }}>ID:</span> {tgRef.id}
													</div>
													{tgRef.members && tgRef.members.length > 0 && (
														<div style={{ color: "#8B949E", fontSize: "10px", marginTop: "2px" }}>
															Vaults: {tgRef.members.map((m: any) => m.vault_id || m).join(", ")}
														</div>
													)}
												</div>
											);
										})()}
									</div>
								)}
							</div>
						</div>
					))}
				</div>
			)}
		</div>
	);
};
